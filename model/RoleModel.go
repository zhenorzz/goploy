package model

import (
	"errors"
)

// Role mysql table server
type Role struct {
	ID             uint32 `json:"id"`
	Name           string `json:"name"`
	PermissionList string `json:"permissionList"`
	Remark         string `json:"remark"`
	CreateTime     int64  `json:"createTime"`
	UpdateTime     int64  `json:"updateTime"`
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
	rows, err := DB.Query("SELECT id, name, remark, permission_list FROM role")
	if err != nil {
		return nil, err
	}
	var roles Roles
	for rows.Next() {
		var role Role

		if err := rows.Scan(&role.ID, &role.Name, &role.Remark, &role.PermissionList); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// AddRow add one row to table server and add id to s.ID
func (r Role) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO role (name, remark, permission_list, create_time, update_time) VALUES (?, ?, ?, ?, ?)",
		r.Name,
		r.Remark,
		r.PermissionList,
		r.CreateTime,
		r.UpdateTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table Role
func (r Role) EditRow() error {
	_, err := DB.Exec(
		`UPDATE role SET 
			name = ?,
			remark = ?,
			permission_list = ?
		WHERE
		 id = ?`,
		r.Name,
		r.Remark,
		r.PermissionList,
		r.ID,
	)
	return err
}
