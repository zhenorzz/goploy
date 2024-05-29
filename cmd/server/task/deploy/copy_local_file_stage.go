package deploy

import (
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
)

func (gsync *Gsync) copyLocalFileStage() error {
	if totalFileNumber, err := (model.ProjectFile{ProjectID: gsync.Project.ID}).GetTotalByProjectID(); err != nil {
		return err
	} else if totalFileNumber > 0 {
		if err := pkg.CopyDir(config.GetProjectFilePath(gsync.Project.ID), config.GetProjectPath(gsync.Project.ID)); err != nil {
			return err
		}
	}
	return nil
}
