package model

import sq "github.com/Masterminds/squirrel"

const projectUserTable = "`project_user`"

// ProjectUser project user relationship
type ProjectUser struct {
	Project
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`
	UserID      int64  `json:"userId"`
	UserName    string `json:"userName"`
	InsertTime  string  `json:"insertTime"`
	UpdateTime  string  `json:"updateTime"`
}

// ProjectUsers project user relationship
type ProjectUsers []ProjectUser

// GetBindUserListByProjectID user row
func (pu ProjectUser) GetBindUserListByProjectID() (ProjectUsers, error) {
	rows, err := sq.
		Select("project_user.id, project_id, user_id, user.name, project_user.insert_time, project_user.update_time").
		From(projectUserTable).
		LeftJoin(userTable + " ON project_user.user_id = user.id").
		Where(sq.Eq{"project_id": pu.ProjectID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var projectUsers ProjectUsers
	for rows.Next() {
		var projectUser ProjectUser

		if err := rows.Scan(&projectUser.ID, &projectUser.ProjectID, &projectUser.UserID, &projectUser.UserName, &projectUser.InsertTime, &projectUser.UpdateTime); err != nil {
			return projectUsers, err
		}
		projectUsers = append(projectUsers, projectUser)
	}
	return projectUsers, nil
}

// GetBindUserListByProjectID user row
func (pu ProjectUser) GetBindProjectListByUserID() (ProjectUsers, error) {
	rows, err := sq.
		Select("project_user.id, project_id, user_id, project.name, project.group_id, project_user.insert_time, project_user.update_time").
		From(projectUserTable).
		LeftJoin(projectTable + " ON project_user.project_id = project.id").
		Where(sq.Eq{"user_id": pu.UserID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var projectUsers ProjectUsers
	for rows.Next() {
		var projectUser ProjectUser

		if err := rows.Scan(&projectUser.ID, &projectUser.ProjectID, &projectUser.UserID, &projectUser.ProjectName, &projectUser.GroupID, &projectUser.InsertTime, &projectUser.UpdateTime); err != nil {
			return projectUsers, err
		}
		projectUsers = append(projectUsers, projectUser)
	}
	return projectUsers, nil
}

// GetListByUserID user row
func (pu ProjectUser) GetListByUserID() (ProjectUsers, error) {
	rows, err := sq.
		Select("id, project_id, user_id, insert_time, update_time").
		From(projectUserTable).
		Where(sq.Eq{"user_id": pu.UserID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var projectUsers ProjectUsers
	for rows.Next() {
		var projectUser ProjectUser

		if err := rows.Scan(&projectUser.ID, &projectUser.ProjectID, &projectUser.UserID, &projectUser.InsertTime, &projectUser.UpdateTime); err != nil {
			return projectUsers, err
		}
		projectUsers = append(projectUsers, projectUser)
	}
	return projectUsers, nil
}

// GetListLeftJoinProjectByUserID user row by status
func (pu ProjectUser) GetListLeftJoinProjectByUserID() (Projects, error) {
	builder := sq.
		Select("project_id, project.name, publisher_id, publisher_name, project.group_id, project.environment, project.branch, project.last_publish_token, project.deploy_state, project.update_time").
		From(projectUserTable).
		LeftJoin(projectTable + " ON project_user.project_id = project.id").
		Where(sq.Eq{
			"project_user.user_id": pu.UserID,
			"project.state":        Enable,
		})
	if pu.GroupID != 0 {
		builder = builder.Where(sq.Eq{"project.group_id": pu.GroupID})
	}

	if len(pu.Name) > 0 {
		builder = builder.Where(sq.Like{"project.name": "%" + pu.Name + "%"})
	}

	rows, err := builder.RunWith(DB).Query()
	if err != nil {
		return nil, err
	}
	var projects Projects
	for rows.Next() {
		var project Project

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.PublisherID,
			&project.PublisherName,
			&project.GroupID,
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

// AddMany add many row to table project_server
func (pu ProjectUsers) AddMany() error {
	if len(pu) == 0 {
		return nil
	}
	builder := sq.
		Insert(projectUserTable).
		Columns("project_id", "user_id")

	for _, row := range pu {
		builder = builder.Values(row.ProjectID, row.UserID, row.InsertTime, row.UpdateTime)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

// DeleteRow edit one row to table ProjectUser
func (pu ProjectUser) DeleteRow() error {
	_, err := sq.
		Delete(projectUserTable).
		Where(sq.Eq{"id": pu.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeleteRow edit one row to table ProjectUser
func (pu ProjectUser) DeleteByUserID() error {
	_, err := sq.
		Delete(projectUserTable).
		Where(sq.Eq{"user_id": pu.UserID}).
		RunWith(DB).
		Exec()
	return err
}
