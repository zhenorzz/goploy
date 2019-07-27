package model

import (
	"errors"
	"time"
)

// Project mysql table project
type Project struct {
	ID            uint32 `json:"id"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	Path          string `json:"path"`
	Environment   string `json:"environment"`
	Branch        string `json:"branch"`
	Script        string `json:"script"`
	RsyncOption   string `json:"rsyncOption"`
	PublisherID   uint32 `json:"publisherId"`
	PublisherName string `json:"publisherName"`
	PublishState  uint8  `json:"publishState"`
	State         uint8  `json:"state"`
	CreateTime    int64  `json:"createTime"`
	UpdateTime    int64  `json:"updateTime"`
}

// Projects many project
type Projects []Project

// AddRow add one row to table project and add id to p.ID
func (p Project) AddRow() (uint32, error) {
	result, err := DB.Exec(
		`INSERT INTO project (
			name, 
			url, 
			path, 
			environment, 
			branch, 
			script, 
			rsync_option, 
			create_time, 
			update_time) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name,
		p.URL,
		p.Path,
		p.Environment,
		p.Branch,
		p.Script,
		p.RsyncOption,
		p.CreateTime,
		p.UpdateTime,
	)
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table project
func (p Project) EditRow() error {
	_, err := DB.Exec(
		`UPDATE project SET 
		  name = ?,
		  url = ?,
		  path = ?,
		  environment = ?,
		  branch = ?,
		  script = ?,
		  rsync_option = ?
		WHERE
		 id = ?`,
		p.Name,
		p.URL,
		p.Path,
		p.Environment,
		p.Branch,
		p.Script,
		p.RsyncOption,
		p.ID,
	)
	return err
}

// EditPublishState update publish state
func (p Project) EditPublishState() error {
	_, err := DB.Exec(
		`UPDATE project SET 
		  publish_state = ?
		WHERE
		 id = ?`,
		p.PublishState,
		p.ID,
	)
	return err
}

// RemoveRow project
func (p Project) RemoveRow() error {
	_, err := DB.Exec(
		`UPDATE project SET 
		  state = 0,
		  update_time = ?
		WHERE
		 id = ?`,
		p.UpdateTime,
		p.ID,
	)
	return err
}

// Publish for project
func (p Project) Publish() error {
	_, err := DB.Exec(
		"UPDATE project SET publisher_id = ?, publisher_name = ?, update_time = ? WHERE id = ?",
		p.PublisherID,
		p.PublisherName,
		p.UpdateTime,
		p.ID,
	)
	return err
}

// GetList project row
func (p Project) GetList() (Projects, error) {
	rows, err := DB.Query("SELECT id, name, url, path, environment, branch, script, rsync_option, create_time, update_time FROM project WHERE state = 1 ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	var projects Projects
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.Name, &project.URL, &project.Path, &project.Environment, &project.Branch, &project.Script, &project.RsyncOption, &project.CreateTime, &project.UpdateTime); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// GetData add project information to p *Project
func (p Project) GetData() (Project, error) {
	var project Project
	err := DB.QueryRow("SELECT id, name, url, path, environment, branch, script, rsync_option, create_time, update_time FROM project WHERE id = ?", p.ID).Scan(&project.ID, &project.Name, &project.URL, &project.Path, &project.Environment, &project.Branch, &project.Script, &project.RsyncOption, &project.CreateTime, &project.UpdateTime)
	if err != nil {
		return project, errors.New("数据查询失败")
	}
	return project, nil
}

// FindNeedToUpdateProjectList find the project need to update its publish state
func (p Project) FindNeedToUpdateProjectList() (Projects, error) {
	rows, err := DB.Query("SELECT id, name, url, path, script, rsync_option, create_time, update_time FROM project WHERE update_time >= ?", time.Now().Unix()-30*60)
	if err != nil {
		return nil, err
	}
	var projects Projects
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.Name, &project.URL, &project.Path, &project.Script, &project.RsyncOption, &project.CreateTime, &project.UpdateTime); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
