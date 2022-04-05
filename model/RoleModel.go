package model

import (
	sq "github.com/Masterminds/squirrel"
)

const roleTable = "`role`"

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	InsertTime  string `json:"insertTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

type Roles []Role

func (r Role) AddRow() (int64, error) {
	result, err := sq.
		Insert(roleTable).
		Columns("name", "description").
		Values(r.Name, r.Description).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (r Role) EditRow() error {
	_, err := sq.
		Update(roleTable).
		Set("name", r.Name).
		Set("description", r.Description).
		Where(sq.Eq{"id": r.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (r Role) DeleteRow() error {
	_, err := sq.
		Delete(roleTable).
		Where(sq.Eq{"id": r.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (r Role) GetList(pagination Pagination) (Roles, error) {
	rows, err := sq.
		Select("id, name, description, insert_time, update_time").
		From(roleTable).
		OrderBy("id DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}

	roles := Roles{}
	for rows.Next() {
		var role Role

		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.InsertTime, &role.UpdateTime); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r Role) GetTotal() (int64, error) {
	var total int64
	err := sq.
		Select("COUNT(*) AS count").
		From(roleTable).
		RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r Role) GetData() (Role, error) {
	var role Role
	err := sq.
		Select("name, description").
		From(roleTable).
		Where(sq.Eq{"id": r.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&role.Name, &role.Description)
	return role, err
}
