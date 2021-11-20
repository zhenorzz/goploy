package model

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

const serverTable = "`server`"

// Server -
type Server struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Owner       string `json:"owner"`
	Path        string `json:"path"`
	Password    string `json:"password"`
	NamespaceID int64  `json:"namespaceId"`
	Description string `json:"description"`
	OSInfo      string `json:"osInfo"`
	State       int8   `json:"state"`
	InsertTime  string `json:"insertTime"`
	UpdateTime  string `json:"updateTime"`
}

// Servers -
type Servers []Server

// GetList -
func (s Server) GetList(pagination Pagination) (Servers, error) {
	rows, err := sq.
		Select("id, namespace_id, name, ip, port, owner, path, password, description, os_info, state, insert_time, update_time").
		From(serverTable).
		Where(sq.Eq{
			"namespace_id": []int64{0, s.NamespaceID},
		}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	servers := Servers{}
	for rows.Next() {
		var server Server

		if err := rows.Scan(
			&server.ID,
			&server.NamespaceID,
			&server.Name,
			&server.IP,
			&server.Port,
			&server.Owner,
			&server.Path,
			&server.Password,
			&server.Description,
			&server.OSInfo,
			&server.State,
			&server.InsertTime,
			&server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}

	return servers, nil
}

// GetTotal -
func (s Server) GetTotal() (int64, error) {
	var total int64
	err := sq.
		Select("COUNT(*) AS count").
		From(serverTable).
		Where(sq.Eq{
			"namespace_id": []int64{0, s.NamespaceID},
		}).
		RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetAll -
func (s Server) GetAll() (Servers, error) {
	rows, err := sq.
		Select("id, name, ip, owner, description, insert_time, update_time").
		From(serverTable).
		Where(sq.Eq{
			"namespace_id": []int64{0, s.NamespaceID},
			"state":        Enable,
		}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	servers := Servers{}
	for rows.Next() {
		var server Server
		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Owner, &server.Description, &server.InsertTime, &server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

// GetData -
func (s Server) GetData() (Server, error) {
	var server Server
	err := sq.
		Select("id, namespace_id, name, ip, port, owner, path, password").
		From(serverTable).
		Where(sq.Eq{"id": s.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&server.ID, &server.NamespaceID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.Path, &server.Password)
	if err != nil {
		return server, err
	}
	return server, nil
}

// AddRow return LastInsertId
func (s Server) AddRow() (int64, error) {
	result, err := sq.
		Insert(serverTable).
		Columns("namespace_id", "name", "ip", "port", "owner", "password", "path", "description", "os_info").
		Values(s.NamespaceID, s.Name, s.IP, s.Port, s.Owner, s.Password, s.Path, s.Description, s.OSInfo).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow -
func (s Server) EditRow() error {
	_, err := sq.
		Update(serverTable).
		SetMap(sq.Eq{
			"namespace_id": s.NamespaceID,
			"name":         s.Name,
			"ip":           s.IP,
			"port":         s.Port,
			"owner":        s.Owner,
			"password":     s.Password,
			"path":         s.Path,
			"description":  s.Description,
			"os_info":      s.OSInfo,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(DB).
		Exec()
	return err
}

// ToggleRow -
func (s Server) ToggleRow() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}
	_, err = sq.
		Update(serverTable).
		SetMap(sq.Eq{
			"state": s.State,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	if s.State == Disable {
		_, err = sq.
			Delete(projectServerTable).
			Where(sq.Eq{"server_id": s.ID}).
			RunWith(tx).
			Exec()
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return errors.New("事务提交失败")
	}
	return nil
}
