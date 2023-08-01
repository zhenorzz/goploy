// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/internal/model"
	"sync/atomic"
	"time"
)

var projectTick = time.Tick(time.Minute)

func startProjectTask() {
	atomic.AddInt32(&counter, 1)
	go func() {
		for {
			select {
			case <-projectTick:
				projectTask()
			case <-stop:
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

func projectTask() {
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

		AddDeployTask(Gsync{
			UserInfo:       userInfo,
			Project:        project,
			ProjectServers: projectServers,
			CommitID:       projectTask.CommitID,
			Branch:         projectTask.Branch,
		})
	}
}
