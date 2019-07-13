package model

import "errors"

// Project mysql table project
type Project struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Path       string `json:"path"`
	Status     uint8  `json:"status"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// ServerDetail mysql table server join project_server
type ServerDetail struct {
	ProjectServerID uint32 `json:"projectServerid"`
	Name            string `json:"name"`
	ServerID        string `json:"serverId"`
}

// ProjectDetail mysql table project join server
type ProjectDetail struct {
	Project
	Servers []ServerDetail `json:"serverDetail"`
}

// Projects many project
type Projects []Project

// AddRow add one row to table project and add id to p.ID
func (p *Project) AddRow() error {
	db := NewDB()
	result, err := db.Exec(
		"INSERT INTO project (name, url, path, create_time, update_time) VALUES (?, ?, ?, ?, ?)",
		p.Name,
		p.URL,
		p.Path,
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
	_, err := db.Exec(
		"UPDATE project SET status = ? where id = ?",
		p.Status,
		p.ID,
	)
	return err
}

// Query user row
func (p *Projects) Query() error {
	db := NewDB()
	rows, err := db.Query("SELECT id, name, url, path, status, create_time, update_time FROM project")
	if err != nil {
		return err
	}
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.Name, &project.URL, &project.Path, &project.Status, &project.CreateTime, &project.UpdateTime); err != nil {
			return err
		}
		*p = append(*p, project)
	}
	return nil
}

// QueryRow add project information to p *Project
func (p *Project) QueryRow() error {
	db := NewDB()
	err := db.QueryRow("SELECT name, url, path, status, create_time, update_time FROM project WHERE id = ?", p.ID).Scan(&p.Name, &p.URL, &p.Path, &p.Status, &p.CreateTime, &p.UpdateTime)
	if err != nil {
		return errors.New("数据查询失败")
	}
	return nil
}

// Detail return ProjectDetail
func (p *ProjectDetail) Detail() error {
	if err := p.QueryRow(); err != nil {
		return err
	}
	db := NewDB()
	rows, err := db.Query(
		`SELECT 
			project_server.id as project_server_id,
			server.id as server_id,
			server.name
		FROM server 
		LEFT JOIN project_server 
		ON project_server.server_id = server.id
		WHERE project_id = ?`, p.ID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var serverDetail ServerDetail

		if err := rows.Scan(&serverDetail.ProjectServerID, &serverDetail.ServerID, &serverDetail.Name); err != nil {
			return err
		}
		p.Servers = append(p.Servers, serverDetail)
	}
	return nil
}
