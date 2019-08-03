package model

// ProjectUser project user relationship
type ProjectUser struct {
	Project
	ID         uint32 `json:"id"`
	ProjectID  uint32 `json:"projectId"`
	UserID     uint32 `json:"userId"`
	UserName   string `json:"userName"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// ProjectUsers project user relationship
type ProjectUsers []ProjectUser

// GetBindUserListByProjectID user row
func (pu ProjectUser) GetBindUserListByProjectID() (ProjectUsers, error) {
	rows, err := DB.Query(
		`SELECT 
		    project_user.id,
			project_id,
			user_id,
			user.name,
			project_user.create_time,
			project_user.update_time
		FROM project_user
		LEFT JOIN user 
		ON project_user.user_id = user.id
		WHERE project_id = ?`, pu.ProjectID)
	if err != nil {
		return nil, err
	}
	var projectUsers ProjectUsers
	for rows.Next() {
		var projectUser ProjectUser

		if err := rows.Scan(&projectUser.ID, &projectUser.ProjectID, &projectUser.UserID, &projectUser.UserName, &projectUser.CreateTime, &projectUser.UpdateTime); err != nil {
			return projectUsers, err
		}
		projectUsers = append(projectUsers, projectUser)
	}
	return projectUsers, nil
}

// GetListByUserID user row
func (pu ProjectUser) GetListByUserID() (ProjectUsers, error) {
	rows, err := DB.Query(
		`SELECT 
		    id,
			project_id,
			user_id,
			create_time,
			update_time
		FROM project_user
		WHERE user_id = ?`, pu.UserID)
	if err != nil {
		return nil, err
	}
	var projectUsers ProjectUsers
	for rows.Next() {
		var projectUser ProjectUser

		if err := rows.Scan(&projectUser.ID, &projectUser.ProjectID, &projectUser.UserID, &projectUser.CreateTime, &projectUser.UpdateTime); err != nil {
			return projectUsers, err
		}
		projectUsers = append(projectUsers, projectUser)
	}
	return projectUsers, nil
}

// GetDepolyListByUserID user row by status
func (pu ProjectUser) GetDepolyListByUserID() (Projects, error) {
	rows, err := DB.Query(`
		SELECT 
			project_id, 
			project.name, 
			publisher_id, 
			publisher_name, 
			project.project_group_id, 
			project.environment, 
			project.branch, 
			project.publish_state, 
			project.update_time 
		FROM 
			project_user 
		LEFT JOIN 
			project 
		ON 
			project_user.project_id = project.id
		WHERE 
			project_user.user_id = ?
		AND 
			project.project_group_id = ?
		AND 
			project.state = 1`,
		pu.UserID, pu.ProjectGroupID)
	if err != nil {
		return nil, err
	}
	var projects Projects
	for rows.Next() {
		var project Project

		if err := rows.Scan(&project.ID, &project.Name, &project.PublisherID, &project.PublisherName, &project.ProjectGroupID, &project.Environment, &project.Branch, &project.PublishState, &project.UpdateTime); err != nil {
			return projects, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// GetDataByProjectUser  by  projectid and userid
func (pu ProjectUser) GetDataByProjectUser() (ProjectUser, error) {
	var projectUser ProjectUser
	err := DB.QueryRow(`
		SELECT 
			id, project_id, user_id 
		FROM 
			project_user
		WHERE 
			project_id = ? 
		AND 
			user_id = ?`, pu.ProjectID, pu.UserID).Scan(&projectUser.ID, &projectUser.ProjectID, &projectUser.UserID)
	if err != nil {
		return projectUser, err
	}
	return projectUser, nil
}

// AddMany add many row to table project_server
func (pu ProjectUsers) AddMany() error {
	if len(pu) == 0 {
		return nil
	}
	sqlStr := "INSERT INTO project_user (project_id, user_id, create_time, update_time) VALUES "
	vals := []interface{}{}

	for _, row := range pu {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, row.ProjectID, row.UserID, row.CreateTime, row.UpdateTime)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		return err
	}
	//format all vals at once
	_, err = stmt.Exec(vals...)
	return err
}

// DeleteRow edit one row to table ProjectUser
func (pu ProjectUser) DeleteRow() error {
	_, err := DB.Exec(`DELETE FROM project_user WHERE id = ?`, pu.ID)
	return err
}
