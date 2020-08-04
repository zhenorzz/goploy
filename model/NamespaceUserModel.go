package model

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

const namespaceUserTable = "`namespace_user`"

// NamespaceUser namespace user relationship
type NamespaceUser struct {
	ID            int64  `json:"id,omitempty"`
	NamespaceID   int64  `json:"namespaceId,omitempty"`
	NamespaceName string `json:"namespaceName,omitempty"`
	UserID        int64  `json:"userId,omitempty"`
	UserName      string `json:"userName,omitempty"`
	Role          string `json:"role,omitempty"`
	InsertTime    string `json:"insertTime,omitempty"`
	UpdateTime    string `json:"updateTime,omitempty"`
}

// NamespaceUsers namespace user relationship
type NamespaceUsers []NamespaceUser

// GetBindUserListByNamespaceID user row
func (nu NamespaceUser) GetBindUserListByNamespaceID() (NamespaceUsers, error) {
	rows, err := sq.
		Select("namespace_user.id, namespace_id, user_id, user.name, namespace_user.role, namespace_user.insert_time, namespace_user.update_time").
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

		if err := rows.Scan(&namespaceUser.ID, &namespaceUser.NamespaceID, &namespaceUser.UserID, &namespaceUser.UserName, &namespaceUser.Role, &namespaceUser.InsertTime, &namespaceUser.UpdateTime); err != nil {
			return nil, err
		}
		namespaceUsers = append(namespaceUsers, namespaceUser)
	}
	return namespaceUsers, nil
}

// GetAllUserByNamespaceID user row
func (nu NamespaceUser) GetAllUserByNamespaceID() (NamespaceUsers, error) {
	rows, err := sq.
		Select("user_id, user.name, namespace_user.role").
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

		if err := rows.Scan(&namespaceUser.UserID, &namespaceUser.UserName, &namespaceUser.Role); err != nil {
			return nil, err
		}
		namespaceUsers = append(namespaceUsers, namespaceUser)
	}
	return namespaceUsers, nil
}

// GetAllGteManagerByNamespaceID user row
func (nu NamespaceUser) GetAllGteManagerByNamespaceID() (NamespaceUsers, error) {
	rows, err := sq.
		Select("user_id, role").
		From(namespaceUserTable).
		Where(sq.Eq{
			"namespace_id": nu.NamespaceID,
			"role":         []string{"admin", "manager"},
		}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	namespaceUsers := NamespaceUsers{}
	for rows.Next() {
		var namespaceUser NamespaceUser

		if err := rows.Scan(&namespaceUser.UserID, &namespaceUser.Role); err != nil {
			return nil, err
		}
		namespaceUsers = append(namespaceUsers, namespaceUser)
	}
	return namespaceUsers, nil
}

// AddMany add many row to table namespace_user
func (nu NamespaceUsers) AddMany() error {
	if len(nu) == 0 {
		return nil
	}
	builder := sq.
		Replace(namespaceUserTable).
		Columns("namespace_id", "user_id", "role")

	for _, row := range nu {
		builder = builder.Values(row.NamespaceID, row.UserID, row.Role)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

// AddAdminByNamespaceID add many admin by namespace id to table namespace_user
func (nu NamespaceUser) AddAdminByNamespaceID() error {

	builder := sq.
		Insert(namespaceUserTable).
		Columns("namespace_id", "user_id", "role").
		Select(sq.
			Select(fmt.Sprintf("%d as namespace_id, id as user_id, 'admin' as role", nu.NamespaceID)).
			From(userTable).
			Where(sq.Eq{"super_manager": SuperManager}))
	_, err := builder.RunWith(DB).Exec()
	return err
}

// AddAdminByUserID add many admin by user id to table namespace_user
func (nu NamespaceUser) AddAdminByUserID() error {
	builder := sq.
		Replace(namespaceUserTable).
		Columns("namespace_id", "user_id", "role").
		Select(sq.
			Select(fmt.Sprintf("id as namespace_id, %d as user_id, 'admin' as role", nu.UserID)).
			From(namespaceTable))
	_, err := builder.RunWith(DB).Exec()
	return err
}

// DeleteRow edit one row to table NamespaceUser
func (nu NamespaceUser) DeleteRow() error {
	_, err := sq.
		Delete(namespaceUserTable).
		Where(sq.Eq{"id": nu.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeleteByUserID delete row by user id
func (nu NamespaceUser) DeleteByUserID() error {
	_, err := sq.
		Delete(namespaceUserTable).
		Where(sq.Eq{"user_id": nu.UserID}).
		RunWith(DB).
		Exec()
	return err
}
