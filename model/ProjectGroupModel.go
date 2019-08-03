package model

import (
	"errors"
)

// ProjectGroup mysql table project_group
type ProjectGroup struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// ProjectGroups many ProjectGroup
type ProjectGroups []ProjectGroup

// AddRow add one row to table ProjectGroup
func (pg ProjectGroup) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO project_group (name, create_time, update_time) VALUES (?, ?, ?)",
		pg.Name,
		pg.CreateTime,
		pg.UpdateTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table project_group
func (pg ProjectGroup) EditRow() error {
	_, err := DB.Exec(
		`UPDATE project_group SET 
		  name = ?
		WHERE
		 id = ?`,
		pg.Name,
		pg.ID,
	)
	return err
}

// Remove Server
func (pg ProjectGroup) Remove() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}
	_, err = tx.Exec(
		`DELETE FROM 
			project_group 
		WHERE
		    id = ?`,
		pg.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		`UPDATE project SET 
		  project_group_id = 0,
		  update_time = ?
		WHERE
		 project_group_id = ?`,
		pg.UpdateTime,
		pg.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.New("事务提交失败")
	}
	return nil
}

// GetList ProjectGroup row
func (pg ProjectGroup) GetList() (ProjectGroups, error) {
	rows, err := DB.Query("SELECT id, name, create_time, update_time FROM project_group")
	if err != nil {
		return nil, err
	}
	var projectGroups ProjectGroups
	for rows.Next() {
		var projectGroup ProjectGroup

		if err := rows.Scan(&projectGroup.ID, &projectGroup.Name, &projectGroup.CreateTime, &projectGroup.UpdateTime); err != nil {
			return nil, err
		}
		projectGroups = append(projectGroups, projectGroup)
	}
	return projectGroups, nil
}

// GetAll ProjectGroup row
func (pg ProjectGroup) GetAll() (ProjectGroups, error) {
	rows, err := DB.Query("SELECT id, name, create_time, update_time FROM project_group ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	var projectGroups ProjectGroups
	for rows.Next() {
		var projectGroup ProjectGroup

		if err := rows.Scan(&projectGroup.ID, &projectGroup.Name, &projectGroup.CreateTime, &projectGroup.UpdateTime); err != nil {
			return nil, err
		}
		projectGroups = append(projectGroups, projectGroup)
	}
	return projectGroups, nil
}

// GetData to ProjectGroup
func (pg ProjectGroup) GetData() (ProjectGroup, error) {
	var projectGroup ProjectGroup
	err := DB.QueryRow("SELECT name, create_time, update_time FROM project_group WHERE id = ?", pg.ID).Scan(&projectGroup.Name, &projectGroup.CreateTime, &projectGroup.UpdateTime)
	if err != nil {
		return projectGroup, errors.New("数据查询失败")
	}
	return projectGroup, nil
}
