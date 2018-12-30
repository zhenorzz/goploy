package model

import "errors"

// Project mysql table project
type Project struct {
	ID         uint32 `json:"id"`
	Project    string `json:"project"`
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	Status     uint8  `json:"status"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// Projects many project
type Projects []Project

// AddRow add one row to table project and add id to p.ID
func (p *Project) AddRow() error {
	db := NewDB()
	result, err := db.Exec(
		"INSERT INTO project (project, owner, repository, create_time, update_time) VALUES (?, ?, ?, ?, ?)",
		p.Project,
		p.Owner,
		p.Repository,
		p.CreateTime,
		p.UpdateTime,
	)
	id, err := result.LastInsertId()
	p.ID = uint32(id)
	return err
}

// ChangeStatus for project
func (p *Project) ChangeStatus() error {
	db := NewDB()
	result, err := db.Exec(
		"UPDATE project SET status = ? where id = ?",
		p.Status,
		p.ID,
	)
	id, err := result.LastInsertId()
	p.ID = uint32(id)
	return err
}

// Query user row
func (p *Projects) Query() error {
	db := NewDB()
	rows, err := db.Query("SELECT id, project, owner, repository, status, create_time, update_time FROM project")
	if err != nil {
		return err
	}
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.Project, &project.Owner, &project.Repository, &project.Status, &project.CreateTime, &project.UpdateTime); err != nil {
			return err
		}
		*p = append(*p, project)
	}
	return nil
}

// QueryRow add project information to p *Project
func (p *Project) QueryRow() error {
	db := NewDB()
	err := db.QueryRow("SELECT project, owner, repository, status, create_time, update_time FROM project WHERE id = ?", p.ID).Scan(&p.Project, &p.Owner, &p.Repository, &p.Status, &p.CreateTime, &p.UpdateTime)
	if err != nil {
		return errors.New("数据查询失败")
	}
	return nil
}
