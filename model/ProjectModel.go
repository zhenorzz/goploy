package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

const projectTable = "`project`"

// Project mysql table project
type Project struct {
	ID                    int64  `json:"id"`
	NamespaceID           int64  `json:"namespaceId"`
	UserID                int64  `json:"userId,omitempty"`
	Name                  string `json:"name"`
	URL                   string `json:"url"`
	Path                  string `json:"path"`
	SymlinkPath           string `json:"symlinkPath"`
	Environment           string `json:"environment"`
	Branch                string `json:"branch"`
	AfterPullScriptMode   string `json:"afterPullScriptMode"`
	AfterPullScript       string `json:"afterPullScript"`
	AfterDeployScriptMode string `json:"afterDeployScriptMode"`
	AfterDeployScript     string `json:"afterDeployScript"`
	RsyncOption           string `json:"rsyncOption"`
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

const (
	ProjectNotDeploy = 0
	ProjectDeploying = 1
	ProjectSuccess   = 2
	ProjectFail      = 3
)

const (
	ProjectManualDeploy  = 0
	ProjectWebhookDeploy = 1
)

const (
	NotifyWeiXin  = 1
	NotifyDingTalk  = 2
	NotifyFeiShu  = 3
	NotifyCustom  = 255
)


// Projects many project
type Projects []Project

// AddRow add one row to table project and add id to p.ID
func (p Project) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectTable).
		Columns("namespace_id", "name", "url", "path", "symlink_path", "environment", "branch", "after_pull_script_mode", "after_pull_script", "after_deploy_script_mode", "after_deploy_script", "rsync_option", "auto_deploy", "notify_type", "notify_target").
		Values(p.NamespaceID, p.Name, p.URL, p.Path, p.SymlinkPath, p.Environment, p.Branch, p.AfterPullScriptMode, p.AfterPullScript, p.AfterDeployScriptMode, p.AfterDeployScript, p.RsyncOption, p.AutoDeploy, p.NotifyType, p.NotifyTarget).
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
			"name":                     p.Name,
			"url":                      p.URL,
			"path":                     p.Path,
			"symlink_path":             p.SymlinkPath,
			"environment":              p.Environment,
			"branch":                   p.Branch,
			"after_pull_script_mode":   p.AfterPullScriptMode,
			"after_pull_script":        p.AfterPullScript,
			"after_deploy_script_mode": p.AfterDeployScriptMode,
			"after_deploy_script":      p.AfterDeployScript,
			"rsync_option":             p.RsyncOption,
			"auto_deploy":              p.AutoDeploy,
			"notify_type":              p.NotifyType,
			"notify_target":            p.NotifyTarget,
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
			"state": Disable,
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
		}).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()

	return err
}

// GetList project row
func (p Project) GetList(pagination Pagination) (Projects, error) {
	builder := sq.
		Select("project.id, name, url, path, symlink_path, environment, branch, after_pull_script_mode, after_pull_script, after_deploy_script_mode, after_deploy_script, rsync_option, auto_deploy, notify_type, notify_target, project.insert_time, project.update_time").
		From(projectTable).
		Join(projectUserTable + " ON project_user.project_id = project.id").
		Where(sq.Eq{
			"namespace_id": p.NamespaceID,
			"user_id":      p.UserID,
			"state":        Enable,
		})

	if len(p.Name) > 0 {
		builder = builder.Where(sq.Like{"name": "%" + p.Name + "%"})
	}

	rows, err := builder.Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
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
			&project.URL,
			&project.Path,
			&project.SymlinkPath,
			&project.Environment,
			&project.Branch,
			&project.AfterPullScriptMode,
			&project.AfterPullScript,
			&project.AfterDeployScriptMode,
			&project.AfterDeployScript,
			&project.RsyncOption,
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

// GetList project total
func (p Project) GetTotal() (int64, error) {
	var total int64
	builder := sq.
		Select("COUNT(*) AS count").
		From(projectTable).
		Join(projectUserTable + " ON project_user.project_id = project.id").
		Where(sq.Eq{
			"namespace_id": p.NamespaceID,
			"user_id":      p.UserID,
			"state":        Enable,
		})
	if len(p.Name) > 0 {
		builder = builder.Where(sq.Like{"name": "%" + p.Name + "%"})
	}

	err := builder.RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (p Project) GetUserProjectList(userID int64) (Projects, error) {
	builder := sq.
		Select(`
			project.id, 
			project.name,
			project.publisher_id,
			project.publisher_name,
			IFNULL(publish_trace.ext, '{}'),
			project.environment, 
			project.branch, 
			project.last_publish_token,
			project.deploy_state, 
			project.update_time`).
		From(projectUserTable).
		LeftJoin(projectTable + " ON project_user.project_id = project.id").
		LeftJoin(fmt.Sprintf("%[1]s ON %[1]s.token = %s.last_publish_token and type = %d", publishTraceTable, projectTable, Pull)).
		Where(sq.Eq{
			"project.namespace_id": p.NamespaceID,
			"project_user.user_id": userID,
			"project.state":        Enable,
		}).
		OrderBy("project.id DESC")

	if len(p.Name) > 0 {
		builder = builder.Where(sq.Like{"project.name": "%" + p.Name + "%"})
	}

	rows, err := builder.RunWith(DB).Query()
	if err != nil {
		return nil, err
	}
	projects := Projects{}
	for rows.Next() {
		var project Project

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.PublisherID,
			&project.PublisherName,
			&project.PublishExt,
			&project.Environment,
			&project.Branch,
			&project.LastPublishToken,
			&project.DeployState,
			&project.UpdateTime); err != nil {
			return projects, err
		}
		projects = append(projects, project)
	}

	return projects, nil

}

// GetAll Group row
func (p Project) GetAllByName() (Projects, error) {
	rows, err := sq.
		Select("id, name, url, path, environment, branch, rsync_option, deploy_state").
		From(projectTable).
		Where(sq.Eq{
			"namespace_id": p.NamespaceID,
			"name":         p.Name,
			"state":        Enable,
		}).
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
			&project.URL,
			&project.Path,
			&project.Environment,
			&project.Branch,
			&project.RsyncOption,
			&project.DeployState); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// GetData add project information to p *Project
func (p Project) GetData() (Project, error) {
	var project Project
	err := sq.
		Select("id, namespace_id, name, url, path, symlink_path, environment, branch, after_pull_script_mode, after_pull_script, after_deploy_script_mode, after_deploy_script, rsync_option, auto_deploy, deploy_state, notify_type, notify_target, insert_time, update_time").
		From(projectTable).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		QueryRow().
		Scan(
			&project.ID,
			&project.NamespaceID,
			&project.Name,
			&project.URL,
			&project.Path,
			&project.SymlinkPath,
			&project.Environment,
			&project.Branch,
			&project.AfterPullScriptMode,
			&project.AfterPullScript,
			&project.AfterDeployScriptMode,
			&project.AfterDeployScript,
			&project.RsyncOption,
			&project.AutoDeploy,
			&project.DeployState,
			&project.NotifyType,
			&project.NotifyTarget,
			&project.InsertTime,
			&project.UpdateTime)
	if err != nil {
		return project, err
	}
	return project, nil
}

// GetData add project information to p *Project
func (p Project) GetDataByName() (Project, error) {
	var project Project
	err := sq.
		Select("id, namespace_id, name, url, path, symlink_path, environment, branch, after_pull_script_mode, after_pull_script, after_deploy_script_mode, after_deploy_script, rsync_option, auto_deploy, deploy_state, notify_type, notify_target, insert_time, update_time").
		From(projectTable).
		Where(sq.Eq{"name": p.Name}).
		RunWith(DB).
		QueryRow().
		Scan(
			&project.ID,
			&project.NamespaceID,
			&project.Name,
			&project.URL,
			&project.Path,
			&project.SymlinkPath,
			&project.Environment,
			&project.Branch,
			&project.AfterPullScriptMode,
			&project.AfterPullScript,
			&project.AfterDeployScriptMode,
			&project.AfterDeployScript,
			&project.RsyncOption,
			&project.AutoDeploy,
			&project.DeployState,
			&project.NotifyType,
			&project.NotifyTarget,
			&project.InsertTime,
			&project.UpdateTime)
	if err != nil {
		return project, err
	}
	return project, nil
}

func (p Project) GetUserProjectData(userID int64) (Project, error) {
	var project Project
	err := sq.
		Select("project_id, project.name, publisher_id, publisher_name, project.environment, project.branch, project.last_publish_token, project.deploy_state").
		From(projectUserTable).
		LeftJoin(projectTable+" ON project_user.project_id = project.id").
		Where(sq.Eq{
			"project_user.project_id": p.ID,
			"project_user.user_id":    userID,
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
