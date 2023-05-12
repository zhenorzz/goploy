// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const rolePermissionTable = "`role_permission`"

type RolePermission struct {
	ID           int64  `json:"id"`
	RoleID       int64  `json:"roleId"`
	PermissionID int64  `json:"permissionId"`
	InsertTime   string `json:"insertTime,omitempty"`
	UpdateTime   string `json:"updateTime,omitempty"`
}

type RolePermissions []RolePermission

func (rp RolePermission) GetData() (RolePermission, error) {
	var rolePermission RolePermission
	err := sq.
		Select("id, role_id, permission_id").
		From(rolePermissionTable).
		Where(sq.Eq{"role_id": rp.RoleID}).
		RunWith(DB).
		QueryRow().
		Scan(&rolePermission.RoleID, &rolePermission.PermissionID)
	return rolePermission, err
}

func (rp RolePermission) GetList() (RolePermissions, error) {
	rows, err := sq.
		Select("id, role_id, permission_id").
		From(rolePermissionTable).
		Where(sq.Eq{"role_id": rp.RoleID}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}

	rolePermissions := RolePermissions{}
	for rows.Next() {
		var rolePermission RolePermission

		if err := rows.Scan(&rolePermission.ID, &rolePermission.RoleID, &rolePermission.PermissionID); err != nil {
			return nil, err
		}
		rolePermissions = append(rolePermissions, rolePermission)
	}

	return rolePermissions, nil
}

func (rps RolePermissions) AddMany() error {
	if len(rps) == 0 {
		return nil
	}
	builder := sq.
		Insert(rolePermissionTable).
		Columns("role_id", "permission_id")

	for _, row := range rps {
		builder = builder.Values(row.RoleID, row.PermissionID)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

func (rp RolePermission) DeleteByRoleID() error {
	_, err := sq.
		Delete(rolePermissionTable).
		Where(sq.Eq{"role_id": rp.RoleID}).
		RunWith(DB).
		Exec()
	return err
}
