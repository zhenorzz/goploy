package model

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const projectTable = "`project`"

// Project mysql table project
type Project struct {
	ID                int64  `json:"id"`
	GroupID           int64  `json:"groupId"`
	Name              string `json:"name"`
	URL               string `json:"url"`
	Path              string `json:"path"`
	Environment       string `json:"environment"`
	Branch            string `json:"branch"`
	AfterPullScript   string `json:"afterPullScript"`
	AfterDeployScript string `json:"afterDeployScript"`
	RsyncOption       string `json:"rsyncOption"`
	PublisherID       int64  `json:"publisherId"`
	PublisherName     string `json:"publisherName"`
	DeployState       uint8  `json:"deployState"`
	LastPublishToken  string `json:"lastPublishToken"`
	State             uint8  `json:"state"`
	CreateTime        int64  `json:"createTime"`
	UpdateTime        int64  `json:"updateTime"`
}

const (
	ProjectNotDeploy = 0
	ProjectDeploying = 1
	ProjectSuccess   = 2
	ProjectFail      = 3
)

// Projects many project
type Projects []Project

// AddRow add one row to table project and add id to p.ID
func (p Project) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectTable).
		Columns("group_id", "name", "url", "path", "environment", "branch", "after_pull_script", "after_deploy_script", "rsync_option", "create_time", "update_time").
		Values(p.GroupID, p.Name, p.URL, p.Path, p.Environment, p.Branch, p.AfterPullScript, p.AfterDeployScript, p.RsyncOption, p.CreateTime, p.UpdateTime).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table project
func (p Project) EditRow() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"group_id":            p.GroupID,
			"name":                p.Name,
			"url":                 p.URL,
			"path":                p.Path,
			"environment":         p.Environment,
			"branch":              p.Branch,
			"after_pull_script":   p.AfterPullScript,
			"after_deploy_script": p.AfterDeployScript,
			"rsync_option":        p.RsyncOption,
			"update_time":         p.UpdateTime,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()
	return err
}

// RemoveRow project
func (p Project) RemoveRow() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"state":       Disable,
			"update_time": p.UpdateTime,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()
	return err
}

// Publish for project
func (p Project) Publish() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"publisher_id":       p.PublisherID,
			"publisher_name":     p.PublisherName,
			"deploy_state":       p.DeployState,
			"last_publish_token": p.LastPublishToken,
			"update_time":        p.UpdateTime,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeploySuccess return err
func (p Project) DeploySuccess() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"deploy_state": ProjectSuccess,
			"update_time":  p.UpdateTime,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()

	return err
}

// DeploySuccess return err
func (p Project) DeployFail() error {
	_, err := sq.
		Update(projectTable).
		SetMap(sq.Eq{
			"deploy_state": ProjectFail,
			"update_time":  p.UpdateTime,
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()

	return err
}

// GetList project row
func (p Project) GetList(pagination Pagination) (Projects, Pagination, error) {
	rows, err := sq.
		Select("id, group_id, name, url, path, environment, branch, after_pull_script, after_deploy_script, rsync_option, create_time, update_time").
		From(projectTable).
		Where(sq.Eq{"state": Enable}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()

	if err != nil {
		return nil, pagination, err
	}
	var projects Projects
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.GroupID, &project.Name, &project.URL, &project.Path, &project.Environment, &project.Branch, &project.AfterPullScript, &project.AfterDeployScript, &project.RsyncOption, &project.CreateTime, &project.UpdateTime); err != nil {
			return nil, pagination, err
		}
		projects = append(projects, project)
	}
	err = sq.
		Select("COUNT(*) AS count").
		From(projectTable).
		Where(sq.Eq{"state": Enable}).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}
	return projects, pagination, nil
}

// GetListByManagerGroupStr project row
func (p Project) GetListByManagerGroupStr(pagination Pagination, managerGroupStr string) (Projects, Pagination, error) {
	if managerGroupStr == "" {
		return nil, pagination, nil
	}
	builder := sq.
		Select("id, group_id, name, url, path, environment, branch, after_pull_script, after_deploy_script, rsync_option, create_time, update_time").
		From(projectTable).
		Where(sq.Eq{"state": Enable}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC")
	pageBuilder := sq.
		Select("COUNT(*) AS count").
		Where(sq.Eq{"state": Enable}).
		From(projectTable)
	if managerGroupStr != "all" {
		builder = builder.Where(sq.Eq{"group_id": strings.Split(managerGroupStr, ",")})
		pageBuilder = pageBuilder.Where(sq.Eq{"group_id": strings.Split(managerGroupStr, ",")})
	}
	rows, err := builder.RunWith(DB).Query()
	if err != nil {
		return nil, pagination, err
	}
	var projects Projects
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.GroupID, &project.Name, &project.URL, &project.Path, &project.Environment, &project.Branch, &project.AfterPullScript, &project.AfterDeployScript, &project.RsyncOption, &project.CreateTime, &project.UpdateTime); err != nil {
			return nil, pagination, err
		}
		projects = append(projects, project)
	}
	err = pageBuilder.
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}
	return projects, pagination, nil
}

// GetData add project information to p *Project
func (p Project) GetData() (Project, error) {
	var project Project
	err := sq.
		Select("id, group_id, name, url, path, environment, branch, after_pull_script, after_deploy_script, rsync_option, create_time, update_time").
		From(projectTable).
		Where(sq.Eq{"id": p.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(
			&project.ID,
			&project.GroupID,
			&project.Name,
			&project.URL,
			&project.Path,
			&project.Environment,
			&project.Branch,
			&project.AfterPullScript,
			&project.AfterDeployScript,
			&project.RsyncOption,
			&project.CreateTime,
			&project.UpdateTime)
	if err != nil {
		return project, err
	}
	return project, nil
}
