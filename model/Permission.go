package model

import "fmt"

// Permission mysql table server
type Permission struct {
	ID       uint32       `json:"id"`
	Title    string       `json:"title"`
	URI      string       `json:"uri"`
	PID      uint32       `json:"pid"`
	Children []Permission `json:"children"`
}

// Permissions many project
type Permissions []Permission

// GetAll Permission row
func (p Permission) GetAll() (Permissions, error) {
	rows, err := DB.Query("SELECT id, title, uri, pid FROM permission")
	if err != nil {
		return nil, err
	}
	var permissions Permissions
	for rows.Next() {
		var permission Permission

		if err := rows.Scan(&permission.ID, &permission.Title, &permission.URI, &permission.PID); err != nil {
			return permissions, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

// GetAllByPermissionList server row
func (p Permission) GetAllByPermissionList(permissionList string) (Permissions, error) {
	query := fmt.Sprintf("SELECT id, title, uri, pid FROM permission WHERE id IN (%s)", permissionList)
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	var permissions Permissions
	for rows.Next() {
		var permission Permission

		if err := rows.Scan(&permission.ID, &permission.Title, &permission.URI, &permission.PID); err != nil {
			return permissions, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}
