package model

// Project mysql table project
type Project struct {
	ID         uint32 `json:"id"`
	Project    string `json:"project"`
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
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

// Query user row
func (p *Projects) Query() error {
	db := NewDB()
	rows, err := db.Query("SELECT id, project, owner, repository, create_time, update_time FROM project")
	if err != nil {
		return err
	}
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.Project, &project.Owner, &project.Repository, &project.CreateTime, &project.UpdateTime); err != nil {
			return err
		}
		*p = append(*p, project)
	}
	return nil
}
