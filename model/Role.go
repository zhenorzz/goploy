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

// GetData role information to r *Role
func (r Role) GetData() (Role, error) {
	var role Role
	db := NewDB()
	err := db.QueryRow("SELECT id, name, permission_list FROM role WHERE id = ?", r.ID).Scan(&role.ID, &role.Name, &role.PermissionList)
	if err != nil {
		return role, errors.New("数据查询失败")
	}
	return role, nil
}
