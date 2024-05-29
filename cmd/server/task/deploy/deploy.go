// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package deploy

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/cmd/server/ws"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/repo"
	"sync"
	"sync/atomic"
	"time"
)

var deployList = list.New()
var deploySecondTick = time.Tick(time.Millisecond)
var cronMinuteTick = time.Tick(time.Minute)

type Gsync struct {
	UserInfo       model.User
	Project        model.Project
	ProjectServers model.ProjectServers
	PublishTrace   model.PublishTrace
	CommitInfo     repo.CommitInfo
	CommitID       string
	Branch         string
}

type syncMessage struct {
	serverName string
	projectID  int64
	detail     string
	state      int
}

type deployMessage struct {
	ProjectID   int64       `json:"projectId"`
	ProjectName string      `json:"projectName"`
	State       uint8       `json:"state"`
	Message     string      `json:"message"`
	Ext         interface{} `json:"ext"`
}

const (
	Queue = iota
	Deploying
	Success
	Fail
)

func (deployMessage) CanSendTo(*ws.Client) error {
	return nil
}

var projectLogFormat = "projectID: %d %s"

func Run(counter *int32, stop <-chan struct{}) {
	atomic.AddInt32(counter, 1)
	var deployingNumber int32
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case <-deploySecondTick:
				if atomic.LoadInt32(&deployingNumber) < config.Toml.APP.DeployLimit {
					atomic.AddInt32(&deployingNumber, 1)
					if deployElem := deployList.Front(); deployElem != nil {
						wg.Add(1)
						go func(gsync *Gsync) {
							gsync.exec()
							atomic.AddInt32(&deployingNumber, -1)
							wg.Done()
						}(deployList.Remove(deployElem).(*Gsync))
					} else {
						atomic.AddInt32(&deployingNumber, -1)
					}
				}
			case <-cronMinuteTick:
				cronTask()
			case <-stop:
				wg.Wait()
				atomic.AddInt32(counter, -1)
				return
			}
		}
	}()
}

func cronTask() {
	date := time.Now().Format("2006-01-02 15:04:05")
	projectTasks, err := model.ProjectTask{}.GetNotRunListLTDate(date)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("get project task list error, detail:" + err.Error())
	}
	for _, projectTask := range projectTasks {
		project, err := model.Project{ID: projectTask.ProjectID}.GetData()

		if err != nil {
			log.Error("publish task has no project, detail:" + err.Error())
			continue
		}

		if err := projectTask.SetRun(); err != nil {
			log.Error("publish task set run fail, detail:" + err.Error())
			continue
		}

		projectServers, err := model.ProjectServer{ProjectID: projectTask.ProjectID}.GetBindServerListByProjectID()

		if err != nil {
			log.Error("publish task has no server, detail:" + err.Error())
			continue
		}

		userInfo, err := model.User{ID: 1}.GetData()
		if err != nil {
			log.Error("publish task has no user, detail:" + err.Error())
			continue
		}

		project.PublisherID = userInfo.ID
		project.PublisherName = userInfo.Name
		project.DeployState = model.ProjectDeploying
		project.LastPublishToken = uuid.New().String()
		err = project.Publish()
		if err != nil {
			log.Error("publish task change state error, detail:" + err.Error())
			continue
		}

		AddTask(Gsync{
			UserInfo:       userInfo,
			Project:        project,
			ProjectServers: projectServers,
			CommitID:       projectTask.CommitID,
			Branch:         projectTask.Branch,
		})
	}
}

func AddTask(gsync Gsync) {
	ws.Send(ws.Data{
		Type: ws.TypeProject,
		Message: deployMessage{
			ProjectID:   gsync.Project.ID,
			ProjectName: gsync.Project.Name,
			State:       Queue,
			Message:     "Task waiting",
			Ext: struct {
				LastPublishToken string `json:"lastPublishToken"`
			}{gsync.Project.LastPublishToken},
		},
	})

	queueExt := model.PublishTraceQueueExt{
		Script: gsync.Project.Script,
	}
	ext, _ := json.Marshal(queueExt)
	gsync.PublishTrace = model.PublishTrace{
		Token:         gsync.Project.LastPublishToken,
		ProjectID:     gsync.Project.ID,
		ProjectName:   gsync.Project.Name,
		PublisherID:   gsync.UserInfo.ID,
		PublisherName: gsync.UserInfo.Name,
		Ext:           string(ext),
		InsertTime:    time.Now().Format("20060102150405"),
		Type:          model.Queue,
		State:         model.Success,
	}
	if _, err := gsync.PublishTrace.AddRow(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, "insert trace error, "+err.Error())
	}
	deployList.PushBack(&gsync)
}

func (gsync *Gsync) exec() {
	log.Tracef(projectLogFormat, gsync.Project.ID, "deploy start")
	var err error
	defer func() {
		if err == nil {
			return
		}
		ws.Send(ws.Data{
			Type:    ws.TypeProject,
			Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Fail, Message: err.Error()},
		})
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
		if err := gsync.Project.DeployFail(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		gsync.notify(model.ProjectFail, err.Error())

	}()

	err = gsync.repoStage()
	if err != nil {
		return
	}

	err = gsync.copyLocalFileStage()
	if err != nil {
		return
	}

	err = gsync.afterPullScriptStage()
	if err != nil {
		return
	}

	err = gsync.serverStage()
	if err != nil {
		return
	}

	err = gsync.deployFinishScriptStage()
	if err != nil {
		return
	}

	if err := gsync.Project.DeploySuccess(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
	}

	gsync.PublishTrace.Type = model.PublishFinish
	gsync.PublishTrace.State = model.Success
	gsync.PublishTrace.InsertTime = time.Now().Format("20060102150405")
	if _, err := gsync.PublishTrace.AddRow(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
	}

	log.Tracef(projectLogFormat, gsync.Project.ID, "deploy success")
	ws.Send(ws.Data{
		Type:    ws.TypeProject,
		Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Success, Message: "Success", Ext: gsync.CommitInfo},
	})
	gsync.notify(model.ProjectSuccess, "")

	gsync.removeExpiredBackup()
}
