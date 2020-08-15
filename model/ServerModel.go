package model

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

const serverTable = "`server`"

// Server -
type Server struct {
	ID               int64  `json:"id"`
	LastInstallToken string `json:"lastInstallToken"`
	Name             string `json:"name"`
	IP               string `json:"ip"`
	Port             int    `json:"port"`
	Owner            string `json:"owner"`
	NamespaceID      int64  `json:"namespaceId"`
	Description      string `json:"description"`
	InsertTime       string `json:"insertTime"`
	UpdateTime       string `json:"updateTime"`
}

// Servers -
type Servers []Server

// GetList -
func (s Server) GetList(pagination Pagination) (Servers, error) {
	rows, err := sq.
		Select("id, name, ip, port, owner, description, insert_time, update_time").
		From(serverTable).
		Where(sq.Eq{
			"namespace_id": s.NamespaceID,
			"state":        Enable,
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

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.Description, &server.InsertTime, &server.UpdateTime); err != nil {
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
			"namespace_id": s.NamespaceID,
			"state":        Enable,
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
			"namespace_id": s.NamespaceID,
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
		Select("id, name, ip, port, owner, namespace_id").
		From(serverTable).
		Where(sq.Eq{"id": s.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&server.ID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.NamespaceID)
	if err != nil {
		return server, errors.New("数据查询失败")
	}
	return server, nil
}

// AddRow return LastInsertId
func (s Server) AddRow() (int64, error) {
	result, err := sq.
		Insert(serverTable).
		Columns("name", "ip", "port", "owner", "namespace_id", "description").
		Values(s.Name, s.IP, s.Port, s.Owner, s.NamespaceID, s.Description).
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
			"name":        s.Name,
			"ip":          s.IP,
			"port":        s.Port,
			"owner":       s.Owner,
			"description": s.Description,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(DB).
		Exec()
	return err
}

// RemoveRow -
func (s Server) RemoveRow() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}
	_, err = sq.
		Update(serverTable).
		SetMap(sq.Eq{
			"state": Disable,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sq.
		Delete(projectServerTable).
		Where(sq.Eq{"server_id": s.ID}).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.New("事务提交失败")
	}
	return nil
}

// Install server
func (s Server) Install() error {
	_, err := sq.
		Update(serverTable).
		SetMap(sq.Eq{
			"last_install_token": s.LastInstallToken,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(DB).
		Exec()
	return err
}
