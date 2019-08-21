package model

import (
	"errors"
)

// Group mysql table group
type Group struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// Groups many Group
type Groups []Group

// AddRow add one row to table Group
func (g Group) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO `group` (name, create_time, update_time) VALUES (?, ?, ?)",
		g.Name,
		g.CreateTime,
		g.UpdateTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table group
func (g Group) EditRow() error {
	_, err := DB.Exec(
		"UPDATE `group` SET name = ? WHERE id = ?",
		g.Name,
		g.ID,
	)
	return err
}

// Remove Server
func (g Group) Remove() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}
	_, err = tx.Exec(
		"DELETE FROM `group` WHERE id = ?",
		g.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		`UPDATE project SET 
		  group_id = 0,
		  update_time = ?
		WHERE
		 group_id = ?`,
		g.UpdateTime,
		g.ID,
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

// GetList Group row
func (g Group) GetList() (Groups, error) {
	rows, err := DB.Query("SELECT id, name, create_time, update_time FROM `group`")
	if err != nil {
		return nil, err
	}
	var projectGroups Groups
	for rows.Next() {
		var projectGroup Group

		if err := rows.Scan(&projectGroup.ID, &projectGroup.Name, &projectGroup.CreateTime, &projectGroup.UpdateTime); err != nil {
			return nil, err
		}
		projectGroups = append(projectGroups, projectGroup)
	}
	return projectGroups, nil
}

// GetAll Group row
func (g Group) GetAll() (Groups, error) {
	rows, err := DB.Query("SELECT id, name, create_time, update_time FROM `group` ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	var projectGroups Groups
	for rows.Next() {
		var projectGroup Group

		if err := rows.Scan(&projectGroup.ID, &projectGroup.Name, &projectGroup.CreateTime, &projectGroup.UpdateTime); err != nil {
			return nil, err
		}
		projectGroups = append(projectGroups, projectGroup)
	}
	return projectGroups, nil
}

// GetData to Group
func (g Group) GetData() (Group, error) {
	var projectGroup Group
	err := DB.QueryRow("SELECT name, create_time, update_time FROM `group` WHERE id = ?", g.ID).Scan(&projectGroup.Name, &projectGroup.CreateTime, &projectGroup.UpdateTime)
	if err != nil {
		return projectGroup, errors.New("数据查询失败")
	}
	return projectGroup, nil
}
