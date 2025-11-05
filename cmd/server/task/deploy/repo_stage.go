package deploy

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/cmd/server/ws"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/repo"
)

func (gsync *Gsync) repoStage() error {
	ws.Send(ws.Data{
		Type:    ws.TypeProject,
		Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Repo follow"},
	})
	gsync.PublishTrace.Type = model.Pull
	gsync.PublishTrace.InsertTime = time.Now().Format("20060102150405")
	var err error
	r, _ := repo.GetRepo(gsync.Project.RepoType)
	if len(gsync.CommitID) == 0 {
		err = r.Follow(gsync.Project.ID, "origin/"+gsync.Project.Branch, gsync.Project.URL, gsync.Project.Branch)
	} else {
		err = r.Follow(gsync.Project.ID, gsync.CommitID, gsync.Project.URL, gsync.Project.Branch)
	}

	if err != nil {
		gsync.PublishTrace.Detail = err.Error()
		gsync.PublishTrace.State = model.Fail
		if _, err := gsync.PublishTrace.AddRow(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		return err
	}

	commitList, err := r.CommitLog(gsync.Project.ID, 1)
	if err != nil {
		gsync.PublishTrace.Detail = err.Error()
		gsync.PublishTrace.State = model.Fail
		if _, err := gsync.PublishTrace.AddRow(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		return err
	}

	gsync.CommitInfo = commitList[0]
	if gsync.Branch != "" {
		gsync.CommitInfo.Branch = gsync.Branch
	} else {
		gsync.CommitInfo.Branch = "origin/" + gsync.Project.Branch
	}

	ext, _ := json.Marshal(gsync.CommitInfo)
	gsync.PublishTrace.Ext = string(ext)
	gsync.PublishTrace.State = model.Success
	if _, err := gsync.PublishTrace.AddRow(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
	}
	return nil
}
