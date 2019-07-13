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

// QueryRow role information to r *Role
func (r *Role) QueryRow() error {
	db := NewDB()
	err := db.QueryRow("SELECT id, name, permission_list FROM role WHERE id = ?", r.ID).Scan(&r.ID, &r.Name, &r.PermissionList)
	if err != nil {
		return errors.New("数据查询失败")
	}
	return nil
}
