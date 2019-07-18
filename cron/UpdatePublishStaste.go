package cron

import (
	"database/sql"

	"github.com/zhenorzz/goploy/model"
)

func updatePublishState() {
	NeedToUpdateProjectList, _ := model.Project{}.FindNeedToUpdateProjectList()
	for _, project := range NeedToUpdateProjectList {
		println(project.Name)
		gitTrace, err := model.GitTrace{ProjectID: project.ID}.GetLatestRow()
		if err == sql.ErrNoRows {
			continue
		} else if err != nil {
			continue
		} else if gitTrace.State == 0 {
			model.Project{ID: project.ID, PublishState: 0}.EditPublishState()
			continue
		}
		num, err := model.RemoteTrace{GitTraceID: gitTrace.ID}.CountFailStateByGitTraceID()
		if err != nil {
			continue
		} else if num != 0 {
			model.Project{ID: project.ID, PublishState: 0}.EditPublishState()
			continue
		} else {
			model.Project{ID: project.ID, PublishState: 1}.EditPublishState()
			continue
		}
	}
}
