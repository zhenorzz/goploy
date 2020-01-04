package model

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
)

const serverTable = "`server`"

// Server mysql table server
type Server struct {
	ID               int64  `json:"id"`
	LastInstallToken string `json:"lastInstallToken"`
	Name             string `json:"name"`
	IP               string `json:"ip"`
	Port             int    `json:"port"`
	Owner            string `json:"owner"`
	GroupID          int64  `json:"groupId"`
	Description      string `json:"description"`
	CreateTime       int64  `json:"createTime"`
	UpdateTime       int64  `json:"updateTime"`
}

// Servers many server
type Servers []Server

// GetList server row
func (s Server) GetList(pagination Pagination) (Servers, Pagination, error) {
	rows, err := sq.
		Select("id, name, ip, port, owner, group_id, description, create_time, update_time").
		From(serverTable).
		Where(sq.Eq{"state": Enable}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, pagination, err
	}
	var servers Servers
	for rows.Next() {
		var server Server

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.GroupID, &server.Description, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, pagination, err
		}
		servers = append(servers, server)
	}
	err = sq.
		Select("COUNT(*) AS count").
		From(serverTable).
		Where(sq.Eq{"state": Enable}).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}
	return servers, pagination, nil
}

// GetListByManagerGroupStr server row
func (s Server) GetListInGroupIDs(groupIDs []string, pagination Pagination) (Servers, Pagination, error) {
	builder := sq.
		Select("id, name, ip, port, owner, group_id, description, create_time, update_time").
		From(serverTable).
		Where(sq.Eq{"state": Enable}).
		Where(sq.Eq{"group_id": groupIDs}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC")
	pageBuilder := sq.
		Select("COUNT(*) AS count").
		Where(sq.Eq{"state": Enable}).
		Where(sq.Eq{"group_id": groupIDs}).
		From(serverTable)
	rows, err := builder.RunWith(DB).Query()
	if err != nil {
		return nil, pagination, err
	}
	var servers Servers
	for rows.Next() {
		var server Server

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.GroupID, &server.Description, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, pagination, err
		}
		servers = append(servers, server)
	}
	err = pageBuilder.
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}

	return servers, pagination, nil
}

// GetAll server row
func (s Server) GetAll() (Servers, error) {
	rows, err := sq.
		Select("id, name, ip, owner, group_id, description, create_time, update_time").
		From(serverTable).
		Where(sq.Eq{"state": Enable}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var servers Servers
	for rows.Next() {
		var server Server
		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Owner, &server.GroupID, &server.Description, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

// GetData add server information to s *Server
func (s Server) GetData() (Server, error) {
	var server Server
	err := sq.
		Select("id, name, ip, port, owner, group_id, create_time, update_time").
		From(serverTable).
		Where(sq.Eq{"id": s.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&server.ID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.GroupID, &server.CreateTime, &server.UpdateTime)
	if err != nil {
		return server, errors.New("数据查询失败")
	}
	return server, nil
}

// AddRow add one row to table server and add id to s.ID
func (s Server) AddRow() (int64, error) {
	result, err := sq.
		Insert(serverTable).
		Columns("name", "ip", "port", "owner", "group_id", "description", "create_time", "update_time").
		Values(s.Name, s.IP, s.Port, s.Owner, s.GroupID, s.Description, s.CreateTime, s.UpdateTime).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table server
func (s Server) EditRow() error {
	_, err := sq.
		Update(serverTable).
		SetMap(sq.Eq{
			"name":        s.Name,
			"ip":          s.IP,
			"port":        s.Port,
			"owner":       s.Owner,
			"group_id":    s.GroupID,
			"description": s.Description,
			"update_time": s.UpdateTime,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(DB).
		Exec()
	return err
}

// Remove Server
func (s Server) Remove() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}
	_, err = sq.
		Update(serverTable).
		SetMap(sq.Eq{
			"state":       Disable,
			"update_time": s.UpdateTime,
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
			"update_time":        s.UpdateTime,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(DB).
		Exec()
	return err
}
