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

// Query server row
func (p *Permissions) Query(permissionList string) error {
	db := NewDB()
	query := fmt.Sprintf("SELECT id, title, uri, pid FROM permission WHERE id IN (%s)", permissionList)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	for rows.Next() {
		var permission Permission

		if err := rows.Scan(&permission.ID, &permission.Title, &permission.URI, &permission.PID); err != nil {
			return err
		}
		*p = append(*p, permission)
	}
	return nil
}
