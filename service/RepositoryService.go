package service

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"os"
	"strconv"
)

// Repository -
type Repository struct {
	ProjectID int64
}

func (repository Repository) Create() error {
	project, err := model.Project{ID: repository.ProjectID}.GetData()
	if err != nil {
		core.Log(core.TRACE, "The project does not exist, projectID:"+strconv.FormatInt(repository.ProjectID, 10))
		return err
	}
	srcPath := core.GetProjectPath(project.Name)
	if _, err := os.Stat(srcPath); err == nil {
		return nil
	}
	if err := os.RemoveAll(srcPath); err != nil {
		core.Log(core.TRACE, "The project fail to remove, projectID:"+strconv.FormatInt(project.ID, 10)+" ,error: "+err.Error())
		return err
	}
	git := utils.GIT{}
	if err := git.Clone([]string{project.URL, srcPath}); err != nil {
		core.Log(core.ERROR, "The project fail to initialize, projectID:"+strconv.FormatInt(project.ID, 10)+" ,error: "+err.Error()+", detail: "+git.Err.String())
		return err
	}

	if project.Branch != "master" {
		if err := git.Checkout([]string{"-b", project.Branch, "origin/" + project.Branch}); err != nil {
			core.Log(core.ERROR, "The project fail to switch branch, projectID:"+strconv.FormatInt(project.ID, 10)+" ,error: "+err.Error()+", detail: "+git.Err.String())
			os.RemoveAll(srcPath)
			return err
		}
	}
	core.Log(core.TRACE, "The project success to initialize, projectID:"+strconv.FormatInt(project.ID, 10))
	return nil
}
