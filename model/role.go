// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

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
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = sq.
		Delete(roleTable).
		Where(sq.Eq{"id": r.ID}).
		RunWith(tx).
		Exec()

	if err == nil {
		_, err = sq.
			Delete(rolePermissionTable).
			Where(sq.Eq{"role_id": r.ID}).
			RunWith(tx).
			Exec()
	}

	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func (r Role) GetList() (Roles, error) {
	rows, err := sq.
		Select("id, name, description, insert_time, update_time").
		From(roleTable).
		OrderBy("id DESC").
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

func (r Role) GetAll() (Roles, error) {
	rows, err := sq.
		Select("id, name, description").
		From(roleTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	roles := Roles{}
	for rows.Next() {
		var role Role

		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return roles, err
		}
		roles = append(roles, role)
	}
	return roles, nil
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
