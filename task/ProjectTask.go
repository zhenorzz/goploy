package task

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/service"
	"time"
)

func projectTask() {
	date := time.Now().Format("2006-01-02 15:04:05")
	projectTasks, err := model.ProjectTask{}.GetNotRunListLTDate(date)
	if err != nil && err != sql.ErrNoRows {
		core.Log(core.ERROR, "get project task list error, detail:"+err.Error())
	}
	for _, projectTask := range projectTasks {
		project, err := model.Project{ID: projectTask.ProjectID}.GetData()

		if err != nil {
			core.Log(core.ERROR, "publish task has no project, detail:"+err.Error())
			continue
		}

		if project.DeployState == model.ProjectDeploying {
			core.Log(core.ERROR, "The project in publish is being build by other")
			continue
		}

		if err := projectTask.SetRun(); err != nil {
			core.Log(core.ERROR, "publish task set run fail, detail:"+err.Error())
			continue
		}

		projectServers, err := model.ProjectServer{ProjectID: projectTask.ProjectID}.GetBindServerListByProjectID()

		if err != nil {
			core.Log(core.ERROR, "publish task has no server, detail:"+err.Error())
			continue
		}

		userInfo, err := model.User{ID: 1}.GetData()
		if err != nil {
			core.Log(core.ERROR, "publish task has no user, detail:"+err.Error())
			continue
		}

		project.PublisherID = userInfo.ID
		project.PublisherName = userInfo.Name
		project.DeployState = model.ProjectDeploying
		project.LastPublishToken = uuid.New().String()
		err = project.Publish()
		if err != nil {
			core.Log(core.ERROR, "publish task change state error, detail:"+err.Error())
			continue
		}
		go service.Gsync{
			UserInfo:       userInfo,
			Project:        project,
			ProjectServers: projectServers,
			CommitID:       projectTask.CommitID,
			Branch:         projectTask.Branch,
		}.Exec()
	}
}
