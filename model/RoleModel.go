package model

import (
	"errors"
)

// Role mysql table server
type Role struct {
	ID             uint32 `json:"id"`
	Name           string `json:"name"`
	PermissionList string `json:"permissionList"`
}

// Roles Role slice
type Roles []Role

// GetData role information to r *Role
func (r Role) GetData() (Role, error) {
	var role Role
	err := DB.QueryRow("SELECT id, name, permission_list FROM role WHERE id = ?", r.ID).Scan(&role.ID, &role.Name, &role.PermissionList)
	if err != nil {
		return role, errors.New("数据查询失败")
	}
	return role, nil
}

// GetAll Role row
func (r Role) GetAll() (Roles, error) {
	rows, err := DB.Query("SELECT id, name, permission_list FROM role")
	if err != nil {
		return nil, err
	}
	var roles Roles
	for rows.Next() {
		var role Role

		if err := rows.Scan(&role.ID, &role.Name, &role.PermissionList); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
