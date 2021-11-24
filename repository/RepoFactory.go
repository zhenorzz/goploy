package repository

import (
	"fmt"
	"github.com/zhenorzz/goploy/model"
)

type Repo interface {
	Ping(url string) error
	// Create one repository
	Create(projectID int64) error
	// Follow the repository code and update to latest
	Follow(project model.Project, target string) error
	// RemoteBranchList list remote branches in the url
	RemoteBranchList(url string) ([]string, error)
	// BranchList list the local repository's branches
	BranchList(projectID int64) ([]string, error)
	// CommitLog list the local commit log
	CommitLog(projectID int64, rows int) ([]CommitInfo, error)
	// BranchLog list the local commit log from specific branch
	BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error)
	// TagLog list the local commit log from all tag
	TagLog(projectID int64, rows int) ([]CommitInfo, error)
}

type CommitInfo struct {
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	Author    string `json:"author"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Tag       string `json:"tag"`
	Diff      string `json:"diff"`
}

func GetRepo(repoType string) (Repo, error) {
	if repoType == "git" {
		return GitRepo{}, nil
	} else if repoType == "svn" {
		return SvnRepo{}, nil
	}
	return nil, fmt.Errorf("wrong repo type passed")
}
