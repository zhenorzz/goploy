// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const namespaceUserTable = "`namespace_user`"

// NamespaceUser -
type NamespaceUser struct {
	ID            int64              `json:"id,omitempty"`
	NamespaceID   int64              `json:"namespaceId,omitempty"`
	NamespaceName string             `json:"namespaceName,omitempty"`
	UserID        int64              `json:"userId,omitempty"`
	UserName      string             `json:"userName,omitempty"`
	RoleID        int64              `json:"roleId,omitempty"`
	RoleName      string             `json:"roleName,omitempty"`
	PermissionIDs map[int64]struct{} `json:"permissionIds,omitempty"`
	InsertTime    string             `json:"insertTime,omitempty"`
	UpdateTime    string             `json:"updateTime,omitempty"`
}

// NamespaceUsers -
type NamespaceUsers []NamespaceUser

func (nu NamespaceUser) GetUserNamespaceList() (NamespaceUsers, error) {
	rows, err := sq.
		Select("namespace_id, namespace.name").
		From(namespaceUserTable).
		Join(namespaceTable + " ON namespace_user.namespace_id = namespace.id").
		Where(sq.Eq{"user_id": nu.UserID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	namespaceUsers := NamespaceUsers{}
	for rows.Next() {
		var namespaceUser NamespaceUser

		if err := rows.Scan(&namespaceUser.NamespaceID, &namespaceUser.NamespaceName); err != nil {
			return nil, err
		}
		namespaceUsers = append(namespaceUsers, namespaceUser)
	}
	return namespaceUsers, nil
}

func (nu NamespaceUser) GetBindUserListByNamespaceID() (NamespaceUsers, error) {
	rows, err := sq.
		Select(
			"namespace_user.id",
			"namespace_user.user_id",
			"user.name",
			"namespace_user.role_id",
			"IFNULL(role.name, '')",
			"namespace_user.insert_time",
			"namespace_user.update_time",
		).
		From(namespaceUserTable).
		LeftJoin(userTable + " ON namespace_user.user_id = user.id").
		LeftJoin(roleTable + " ON namespace_user.role_id = role.id").
		Where(sq.Eq{"namespace_id": nu.NamespaceID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	namespaceUsers := NamespaceUsers{}
	for rows.Next() {
		var namespaceUser NamespaceUser

		if err := rows.Scan(
			&namespaceUser.ID,
			&namespaceUser.UserID,
			&namespaceUser.UserName,
			&namespaceUser.RoleID,
			&namespaceUser.RoleName,
			&namespaceUser.InsertTime,
			&namespaceUser.UpdateTime,
		); err != nil {
			return nil, err
		}
		namespaceUsers = append(namespaceUsers, namespaceUser)
	}
	return namespaceUsers, nil
}

func (nu NamespaceUser) GetAllUserByNamespaceID() (NamespaceUsers, error) {
	rows, err := sq.
		Select("user_id, user.name, namespace_user.role_id").
		From(namespaceUserTable).
		LeftJoin(userTable + " ON namespace_user.user_id = user.id").
		Where(sq.Eq{"namespace_id": nu.NamespaceID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	namespaceUsers := NamespaceUsers{}
	for rows.Next() {
		var namespaceUser NamespaceUser

		if err := rows.Scan(&namespaceUser.UserID, &namespaceUser.UserName, &namespaceUser.RoleID); err != nil {
			return nil, err
		}
		namespaceUsers = append(namespaceUsers, namespaceUser)
	}
	return namespaceUsers, nil
}

func (nu NamespaceUser) GetDataByUserNamespace() (NamespaceUser, error) {
	var namespaceUser NamespaceUser
	var permissionIDs string
	err := sq.Select("namespace_user.role_id, IFNULL(GROUP_CONCAT(permission_id), '')").
		From(namespaceUserTable).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.role_id = %s.role_id", rolePermissionTable, namespaceUserTable)).
		Where(sq.Eq{"user_id": nu.UserID, "namespace_id": nu.NamespaceID}).
		GroupBy("namespace_user.role_id").
		RunWith(DB).
		QueryRow().
		Scan(&namespaceUser.RoleID, &permissionIDs)

	if err != nil {
		return namespaceUser, err
	}
	namespaceUser.PermissionIDs = map[int64]struct{}{}
	if permissionIDs == "" {
		return namespaceUser, err
	}
	for _, permissionID := range strings.Split(permissionIDs, ",") {
		v, err := strconv.ParseInt(permissionID, 10, 64)
		if err != nil {
			return namespaceUser, err
		}
		namespaceUser.PermissionIDs[v] = struct{}{}
	}
	return namespaceUser, err
}

func (nu NamespaceUser) GetDataByRoleID() (NamespaceUser, error) {
	var namespaceUser NamespaceUser
	err := sq.Select("id, namespace_id, role_id").
		From(namespaceUserTable).
		Where(sq.Eq{"role_id": nu.RoleID}).
		RunWith(DB).
		QueryRow().
		Scan(&namespaceUser.ID, &namespaceUser.NamespaceID, &namespaceUser.RoleID)
	if err != nil {
		return namespaceUser, err
	}

	return namespaceUser, err
}

func (nu NamespaceUsers) AddMany() error {
	if len(nu) == 0 {
		return nil
	}
	builder := sq.
		Replace(namespaceUserTable).
		Columns("namespace_id", "user_id", "role_id")

	for _, row := range nu {
		builder = builder.Values(row.NamespaceID, row.UserID, row.RoleID)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

func (nu NamespaceUser) AddAdminByNamespaceID() error {

	builder := sq.
		Insert(namespaceUserTable).
		Columns("namespace_id", "user_id", "role_id").
		Select(sq.
			Select(fmt.Sprintf("%d as namespace_id, id as user_id, 0 as role_id", nu.NamespaceID)).
			From(userTable).
			Where(sq.Eq{"super_manager": SuperManager}))
	_, err := builder.RunWith(DB).Exec()
	return err
}

func (nu NamespaceUser) AddAdminByUserID() error {
	builder := sq.
		Replace(namespaceUserTable).
		Columns("namespace_id", "user_id", "role_id").
		Select(sq.
			Select(fmt.Sprintf("id as namespace_id, %d as user_id, 0 as role_id", nu.UserID)).
			From(namespaceTable))
	_, err := builder.RunWith(DB).Exec()
	return err
}

func (nu NamespaceUser) DeleteRow() error {
	_, err := sq.
		Delete(namespaceUserTable).
		Where(sq.Eq{"id": nu.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (nu NamespaceUser) DeleteByUserID() error {
	_, err := sq.
		Delete(namespaceUserTable).
		Where(sq.Eq{"user_id": nu.UserID}).
		RunWith(DB).
		Exec()
	return err
}
