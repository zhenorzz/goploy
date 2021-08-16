package repository

import (
	"fmt"
	"github.com/zhenorzz/goploy/model"
)

type Repo interface {
	Create(projectID int64) error
	Follow(project model.Project, target string) error
	BranchList(projectID int64) ([]string, error)
	CommitLog(projectID int64, rows int) ([]CommitInfo, error)
	BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error)
	TagLog(projectID int64, rows int) ([]CommitInfo, error)
}

type CommitInfo struct {
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	Author    string `json:"author"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	Tag       string `json:"tag"`
	Diff      string `json:"diff"`
}

func GetRepo(repoType string) (Repo, error) {
	if repoType == "git" {
		return GitRepo{}, nil
	}
	return nil, fmt.Errorf("wrong repo type passed")
}
