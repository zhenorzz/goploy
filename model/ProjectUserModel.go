package model

// ProjectUser project user relationship
type ProjectUser struct {
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
func (pu ProjectUser) GetBindUserListByProjectID(projectID uint32) (ProjectUsers, error) {
	db := NewDB()
	rows, err := db.Query(
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
		WHERE project_id = ?`, projectID)
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

// AddMany add many row to table project_server
func (pu ProjectUsers) AddMany() error {
	db := NewDB()
	sqlStr := "INSERT INTO project_user (project_id, user_id, create_time, update_time) VALUES "
	vals := []interface{}{}

	for _, row := range pu {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, row.ProjectID, row.UserID, row.CreateTime, row.UpdateTime)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	//format all vals at once
	_, err = stmt.Exec(vals...)
	return err
}

// DeleteRow edit one row to table ProjectUser
func (pu ProjectUser) DeleteRow() error {
	db := NewDB()
	_, err := db.Exec(`DELETE FROM project_user WHERE id = ?`, pu.ID)
	return err
}
