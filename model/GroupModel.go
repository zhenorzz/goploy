package model

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

const groupTable = "`group`"

// Group mysql table group
type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// Groups many Group
type Groups []Group

// AddRow add one row to table Group
func (g Group) AddRow() (int64, error) {
	result, err := sq.
		Insert(groupTable).
		Columns("name", "create_time", "update_time").
		Values(g.Name, g.CreateTime, g.UpdateTime).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table group
func (g Group) EditRow() error {
	_, err := sq.
		Update(groupTable).
		Set("name", g.Name).
		Where(sq.Eq{"id": g.ID}).
		RunWith(DB).
		Exec()
	return err
}

// Remove Server
func (g Group) Remove() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}

	_, err = sq.
		Delete(groupTable).
		Where(sq.Eq{"id": g.ID}).
		RunWith(tx).
		Exec()

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sq.
		Update(projectTable).
		Set("group_id", 0).
		Set("update_time", g.UpdateTime).
		Where(sq.Eq{"group_id": g.ID}).
		RunWith(tx).
		Exec()

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
func (g Group) GetList(pagination Pagination) (Groups, Pagination, error) {
	rows, err := sq.
		Select("id, name, create_time, update_time").
		From(groupTable).
		OrderBy("id DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, pagination, err
	}

	var groups Groups
	for rows.Next() {
		var group Group

		if err := rows.Scan(&group.ID, &group.Name, &group.CreateTime, &group.UpdateTime); err != nil {
			return nil, pagination, err
		}
		groups = append(groups, group)
	}

	err = sq.
		Select("COUNT(*) AS count").
		From(groupTable).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}

	return groups, pagination, nil
}

// GetList Group row
func (g Group) GetListInGroupIDs(groupIDs []string, pagination Pagination) (Groups, Pagination, error) {
	rows, err := sq.
		Select("id, name, create_time, update_time").
		From(groupTable).
		Where(sq.Eq{"id": groupIDs}).
		OrderBy("id DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, pagination, err
	}

	var groups Groups
	for rows.Next() {
		var group Group

		if err := rows.Scan(&group.ID, &group.Name, &group.CreateTime, &group.UpdateTime); err != nil {
			return nil, pagination, err
		}
		groups = append(groups, group)
	}

	err = sq.
		Select("COUNT(*) AS count").
		From(groupTable).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}

	return groups, pagination, nil
}

// GetAll Group row
func (g Group) GetAll() (Groups, error) {
	rows, err := sq.
		Select("id, name, create_time, update_time").
		From(groupTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
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
	var group Group
	err := sq.
		Select("name, create_time, update_time").
		From(groupTable).
		Where(sq.Eq{"id": g.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&group.Name, &group.CreateTime, &group.UpdateTime)
	if err != nil {
		return group, errors.New("数据查询失败")
	}
	return group, nil
}
