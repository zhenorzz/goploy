package repository

import (
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"io"
	"os"
	"path"
	"time"
)

type FtpRepo struct{}

func (FtpRepo) Ping(url string) error {
	c, err := utils.FTPConfig{}.DialFromURL(url)
	if err != nil {
		return err
	}
	_ = c.Quit()
	return nil
}

// Create -
func (ftpRepo FtpRepo) Create(projectID int64) error {
	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project does not exist, projectID:%d", projectID))
		return err
	}
	return ftpRepo.Follow(project, "")
}

func (ftpRepo FtpRepo) Follow(project model.Project, _ string) error {
	projectID := project.ID
	srcPath := core.GetProjectPath(projectID)
	_ = os.RemoveAll(srcPath)
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to mkdir, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}

	c, err := utils.FTPConfig{}.DialFromURL(project.URL)
	if err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to connect ftp, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	var downloadFromFTP func(localDir, remoteDir string) error
	downloadFromFTP = func(localDir, remoteDir string) error {
		remoteEntries, err := c.List(remoteDir)
		if err != nil {
			return err
		}
		for _, entry := range remoteEntries {
			if entry.Type == 1 {
				nextLocalDir := path.Join(localDir, entry.Name)
				if err := os.Mkdir(nextLocalDir, 0755); err != nil {
					return err
				}
				if err := downloadFromFTP(nextLocalDir, path.Join(remoteDir, entry.Name)); err != nil {
					return err
				}
			} else {
				remoteFile, err := c.Retr(path.Join(remoteDir, entry.Name))
				if err != nil {
					return err
				}
				localFile, err := os.Create(path.Join(localDir, entry.Name))
				if err != nil {
					remoteFile.Close()
					return err
				}
				_, err = io.Copy(localFile, remoteFile)
				if err != nil {
					return err
				}
				localFile.Close()
				remoteFile.Close()
			}
		}
		return nil
	}
	if err := downloadFromFTP(srcPath, ""); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to download file, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	_ = c.Quit()
	core.Log(core.TRACE, fmt.Sprintf("The project success to download, projectID:%d", projectID))
	return nil
}

func (FtpRepo) RemoteBranchList(url string) ([]string, error) {
	return []string{"master"}, nil
}

func (FtpRepo) BranchList(projectID int64) ([]string, error) {
	return []string{"master"}, nil
}

func (FtpRepo) CommitLog(projectID int64, rows int) ([]CommitInfo, error) {
	commitInfo := CommitInfo{
		Branch:    "master",
		Commit:    "",
		Author:    "",
		Timestamp: time.Now().Unix(),
		Message:   "",
		Tag:       "",
		Diff:      "",
	}
	return []CommitInfo{commitInfo}, nil
}

func (ftpRepo FtpRepo) BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error) {
	return []CommitInfo{}, nil
}

func (ftpRepo FtpRepo) TagLog(projectID int64, rows int) ([]CommitInfo, error) {
	return []CommitInfo{}, nil
}
