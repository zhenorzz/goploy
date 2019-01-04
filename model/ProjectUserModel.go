package model

// ProjectUser project user relationship
type ProjectUser struct {
	ID         uint32 `json:"id"`
	ProjectID  uint32 `json:"projectId"`
	UserID     uint32 `json:"userId"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// ProjectUsers project user relationship
type ProjectUsers []ProjectUser

// AddMany add many row to table project_server
func (pu *ProjectUsers) AddMany() error {
	db := NewDB()
	sqlStr := "INSERT INTO project_user (project_id, user_id, create_time, update_time) VALUES "
	vals := []interface{}{}

	for _, row := range *pu {
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
