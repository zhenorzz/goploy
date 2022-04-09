// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const namespaceTable = "`namespace`"

type Namespace struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	UserID        int64   `json:"-"`
	RoleID        int64   `json:"role_id"`
	PermissionIDs []int64 `json:"permissionIds"`
	InsertTime    string  `json:"insertTime,omitempty"`
	UpdateTime    string  `json:"updateTime,omitempty"`
}

type Namespaces []Namespace

func (ns Namespace) AddRow() (int64, error) {
	result, err := sq.
		Insert(namespaceTable).
		Columns("name").
		Values(ns.Name).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (ns Namespace) EditRow() error {
	_, err := sq.
		Update(namespaceTable).
		Set("name", ns.Name).
		Where(sq.Eq{"id": ns.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (ns Namespace) GetAllByUserID() (Namespaces, error) {
	rows, err := sq.
		Select("namespace.id, namespace.name, role_id").
		From(namespaceTable).
		Join(namespaceUserTable + " ON namespace_user.namespace_id = namespace.id").
		Where(sq.Eq{"user_id": ns.UserID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	namespaces := Namespaces{}
	for rows.Next() {
		var namespace Namespace
		if err := rows.Scan(&namespace.ID, &namespace.Name, &namespace.RoleID); err != nil {
			return nil, err
		}
		namespaces = append(namespaces, namespace)
	}
	return namespaces, nil
}

func (ns Namespace) GetList() (Namespaces, error) {
	rows, err := sq.
		Select("id, name, insert_time, update_time").
		From(namespaceTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}

	namespaces := Namespaces{}
	for rows.Next() {
		var namespace Namespace

		if err := rows.Scan(&namespace.ID, &namespace.Name, &namespace.InsertTime, &namespace.UpdateTime); err != nil {
			return nil, err
		}
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

func (ns Namespace) GetData() (Namespace, error) {
	var namespace Namespace
	err := sq.
		Select("name, insert_time, update_time").
		From(namespaceTable).
		Where(sq.Eq{"id": ns.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&namespace.Name, &namespace.InsertTime, &namespace.UpdateTime)
	return namespace, err
}
