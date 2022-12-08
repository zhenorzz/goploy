// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/zhenorzz/goploy/config"
	"path"
	"strconv"
	"strings"
)

const projectTable = "`project`"

type Project struct {
	ID                    int64  `json:"id"`
	NamespaceID           int64  `json:"namespaceId"`
	UserID                int64  `json:"userId,omitempty"`
	RepoType              string `json:"repoType"`
	Name                  string `json:"name"`
	URL                   string `json:"url"`
	Path                  string `json:"path"`
	Environment           uint8  `json:"environment"`
	Branch                string `json:"branch"`
	SymlinkPath           string `json:"symlinkPath"`
	SymlinkBackupNumber   uint8  `json:"symlinkBackupNumber"`
	Review                uint8  `json:"review"`
	ReviewURL             string `json:"reviewURL"`
	AfterPullScriptMode   string `json:"afterPullScriptMode"`
	AfterPullScript       string `json:"afterPullScript"`
	AfterDeployScriptMode string `json:"afterDeployScriptMode"`
	AfterDeployScript     string `json:"afterDeployScript"`
	TransferType          string `json:"transferType"`
	TransferOption        string `json:"transferOption"`
	DeployServerMode      string `json:"deployServerMode"`
	AutoDeploy            uint8  `json:"autoDeploy"`
	PublisherID           int64  `json:"publisherId"`
	PublisherName         string `json:"publisherName"`
	PublishExt            string `json:"publishExt"`
	DeployState           uint8  `json:"deployState"`
	LastPublishToken      string `json:"lastPublishToken"`
	NotifyType            uint8  `json:"notifyType"`
	NotifyTarget          string `json:"notifyTarget"`
	State                 uint8  `json:"state"`
	InsertTime            string `json:"insertTime"`
	UpdateTime            string `json:"updateTime"`
}

// Project deploy state
const (
	ProjectNotDeploy = 0
	ProjectDeploying = 1
	ProjectSuccess   = 2
	ProjectFail      = 3
)

// Project deploy type
const (
	ProjectManualDeploy  = 0
	ProjectWebhookDeploy = 1
)

const (
	RepoGit  = "git"
	RepoSVN  = "svn"
	RepoFTP  = "ftp"
	RepoSFTP = "sftp"
)

// Project notify type
const (
	NotifyWeiXin   = 1
	NotifyDingTalk = 2
	NotifyFeiShu   = 3
	NotifyCustom   = 255
)

type Projects []Project

func (p Project) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectTable).
		Columns(
			"namespace_id",
			"name",
			"repo_type",
			"url",
			"path",
			"environment",
			"branch",
			"symlink_path",
			"symlink_backup_number",
			"review",
			"review_url",
			"after_pull_script_mode",
			"after_pull_script",
			"after_deploy_script_mode",
			"after_deploy_script",
			"transfer_type",
			"transfer_option",
			"deploy_server_mode",
			"notify_type",
			"notify_target",
		).
		Values(
			p.NamespaceID,
			p.Name,
			p.RepoType,
			p.URL,
			p.Path,
			p.Environment,
			p.Branch,
			p.SymlinkPath,
			p.SymlinkBackupNumber,
			p.Review,
			p.ReviewURL,
			p.AfterPullScriptMode,
			p.AfterPullScript,
			p.AfterDeployScriptMode,
			p.AfterDeployScript,
			p.TransferType,
			p.TransferOption,
			p.DeployServerMode,
			p.NotifyType,
			p.NotifyTarget,
		).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (p Project) EditRow() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"name":                     p.Name,
			"repo_type":                p.RepoType,
			"url":                      p.URL,
			"path":                     p.Path,
			"environment":              p.Environment,
			"branch":                   p.Branch,
			"symlink_path":             p.SymlinkPath,
			"symlink_backup_number":    p.SymlinkBackupNumber,
			"review":                   p.Review,
			"review_url":               p.ReviewURL,
			"after_pull_script_mode":   p.AfterPullScriptMode,
			"after_pull_script":        p.AfterPullScript,
			"after_deploy_script_mode": p.AfterDeployScriptMode,
			"after_deploy_script":      p.AfterDeployScript,
			"transfer_type":            p.TransferType,
			"transfer_option":          p.TransferOption,
			"deploy_server_mode":       p.DeployServerMode,
			"notify_type":              p.NotifyType,
			"notify_target":            p.NotifyTarget,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (p Project) SetAutoDeploy() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"auto_deploy": p.AutoDeploy,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()

	return err
}

func (p Project) RemoveRow() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"state": Disable,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (p Project) Publish() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"publisher_id":       p.PublisherID,
			"publisher_name":     p.PublisherName,
			"deploy_state":       p.DeployState,
			"last_publish_token": p.LastPublishToken,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()
	return err
}

// ResetState set deploy_state to NotDeploy
func (p Project) ResetState() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"deploy_state": ProjectNotDeploy,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()

	return err
}

// DeploySuccess set deploy_state to success
func (p Project) DeploySuccess() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"deploy_state": ProjectSuccess,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()

	return err
}

// DeployFail set deploy_state to fail
func (p Project) DeployFail() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"deploy_state": ProjectFail,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()

	return err
}

func (p Project) GetList() (Projects, error) {
	builder := sq.
		Select(`
			project.id, 
			name, 
			repo_type, 
			url, 
			path, 
			environment,
			branch,
			symlink_path, 
			symlink_backup_number, 
			review, 
			review_url, 
			after_pull_script_mode,
			after_pull_script,
			after_deploy_script_mode,
			after_deploy_script,
			transfer_type,
			transfer_option,
			deploy_server_mode,
			auto_deploy,
			notify_type, 
			notify_target,
			project.insert_time,
			project.update_time
		`).
		From(projectTable).
		Where(sq.Eq{
			"namespace_id": p.NamespaceID,
			"state":        Enable,
		})

	if p.UserID > 0 {
		builder = builder.
			Join(projectUserTable + " ON project_user.project_id = project.id").
			Where(sq.Eq{"user_id": p.UserID})
	}

	rows, err := builder.
		OrderBy("id DESC").
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	projects := Projects{}
	for rows.Next() {
		var project Project

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.RepoType,
			&project.URL,
			&project.Path,
			&project.Environment,
			&project.Branch,
			&project.SymlinkPath,
			&project.SymlinkBackupNumber,
			&project.Review,
			&project.ReviewURL,
			&project.AfterPullScriptMode,
			&project.AfterPullScript,
			&project.AfterDeployScriptMode,
			&project.AfterDeployScript,
			&project.TransferType,
			&project.TransferOption,
			&project.DeployServerMode,
			&project.AutoDeploy,
			&project.NotifyType,
			&project.NotifyTarget,
			&project.InsertTime,
			&project.UpdateTime,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (p Project) GetDeployList() (Projects, error) {
	builder := sq.
		Select(`
			project.id, 
			project.name,
			project.repo_type,
			project.transfer_type,
			project.url,
			project.publisher_id,
			project.publisher_name,
			IFNULL(publish_trace.ext, '{}'),
			project.environment, 
			project.branch, 
			project.symlink_path, 
			project.review, 
			project.last_publish_token,
			project.auto_deploy,
			project.deploy_state, 
			project.update_time`).
		From(projectTable).
		LeftJoin(fmt.Sprintf("%[1]s ON %[1]s.token = %s.last_publish_token and type = %d", publishTraceTable, projectTable, Pull)).
		Where(sq.Eq{
			"project.namespace_id": p.NamespaceID,
			"project.state":        Enable,
		})

	if p.UserID > 0 {
		builder = builder.
			LeftJoin(projectUserTable + " ON project_user.project_id = project.id").
			Where(sq.Eq{"project_user.user_id": p.UserID})
	}

	rows, err := builder.
		OrderBy("project.id DESC").
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	projects := Projects{}
	for rows.Next() {
		var project Project

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.RepoType,
			&project.TransferType,
			&project.URL,
			&project.PublisherID,
			&project.PublisherName,
			&project.PublishExt,
			&project.Environment,
			&project.Branch,
			&project.SymlinkPath,
			&project.Review,
			&project.LastPublishToken,
			&project.AutoDeploy,
			&project.DeployState,
			&project.UpdateTime); err != nil {
			return projects, err
		}
		projects = append(projects, project)
	}

	return projects, nil

}

func (p Project) GetData() (Project, error) {
	var project Project
	err := sq.
		Select(`
			id, 
			namespace_id, 
			name, 
			repo_type, 
			url, 
			path, 
			environment, 
			branch, 
			symlink_path, 
			symlink_backup_number, 
			review, 
			review_url,
			after_pull_script_mode, 
			after_pull_script, 
			after_deploy_script_mode, 
			after_deploy_script, 
			transfer_type, 
			transfer_option, 
			deploy_server_mode, 
			auto_deploy, 
			deploy_state, 
			notify_type, 
			notify_target, 
			project.
			last_publish_token, 
			insert_time, 
			update_time`).
		From(projectTable).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		QueryRow().
		Scan(
			&project.ID,
			&project.NamespaceID,
			&project.Name,
			&project.RepoType,
			&project.URL,
			&project.Path,
			&project.Environment,
			&project.Branch,
			&project.SymlinkPath,
			&project.SymlinkBackupNumber,
			&project.Review,
			&project.ReviewURL,
			&project.AfterPullScriptMode,
			&project.AfterPullScript,
			&project.AfterDeployScriptMode,
			&project.AfterDeployScript,
			&project.TransferType,
			&project.TransferOption,
			&project.DeployServerMode,
			&project.AutoDeploy,
			&project.DeployState,
			&project.NotifyType,
			&project.NotifyTarget,
			&project.LastPublishToken,
			&project.InsertTime,
			&project.UpdateTime)
	if err != nil {
		return project, err
	}
	return project, nil
}

func (p Project) GetUserProjectData() (Project, error) {
	var project Project
	err := sq.
		Select("project_id, project.name, publisher_id, publisher_name, project.environment, project.branch, project.last_publish_token, project.deploy_state").
		From(projectUserTable).
		LeftJoin(projectTable+" ON project_user.project_id = project.id").
		Where(sq.Eq{
			"project_user.project_id": p.ID,
			"project_user.user_id":    p.UserID,
			"project.state":           Enable,
		}).
		RunWith(DB).
		QueryRow().
		Scan(&project.ID,
			&project.Name,
			&project.PublisherID,
			&project.PublisherName,
			&project.Environment,
			&project.Branch,
			&project.LastPublishToken,
			&project.DeployState)
	if err != nil {
		return project, err
	}
	return project, nil
}

func (p Project) ReplaceVars(script string) string {
	scriptVars := map[string]string{
		"${PROJECT_ID}":           strconv.FormatInt(p.ID, 10),
		"${PROJECT_PATH}":         p.Path,
		"${PROJECT_SYMLINK_PATH}": path.Join(p.SymlinkPath, p.LastPublishToken),
		"${PROJECT_NAME}":         p.Name,
		"${PROJECT_BRANCH}":       p.Branch,
		"${REPOSITORY_TYPE}":      p.RepoType,
		"${REPOSITORY_URL}":       p.URL,
		"${REPOSITORY_PATH}":      config.GetProjectPath(p.ID),
		"${PUBLISH_TOKEN}":        p.LastPublishToken,
	}
	for key, value := range scriptVars {
		script = strings.Replace(script, key, value, -1)
	}
	return script
}
