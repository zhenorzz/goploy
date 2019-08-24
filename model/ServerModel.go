package model

import (
	"errors"
)

// Server mysql table server
type Server struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Port       uint32 `json:"port"`
	Owner      string `json:"owner"`
	GroupID    uint32 `json:"groupId"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// Servers many server
type Servers []Server

// AddRow add one row to table server and add id to s.ID
func (s Server) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO server (name, ip, port, owner, group_id, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?)",
		s.Name,
		s.IP,
		s.Port,
		s.Owner,
		s.GroupID,
		s.CreateTime,
		s.UpdateTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table server
func (s Server) EditRow() error {
	_, err := DB.Exec(
		`UPDATE server SET 
		  name = ?,
		  ip = ?,
		  port = ?,
		  owner = ?, 
		  group_id = ?,
		  update_time = ?
		WHERE
		 id = ?`,
		s.Name,
		s.IP,
		s.Port,
		s.Owner,
		s.GroupID,
		s.UpdateTime,
		s.ID,
	)
	return err
}

// Remove Server
func (s Server) Remove() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}
	_, err = tx.Exec(
		`UPDATE server SET 
		  state = 0,
		  update_time = ?
		WHERE
		 id = ?`,
		s.UpdateTime,
		s.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		`DELETE FROM 
			project_server 
		WHERE
		 server_id = ?`,
		s.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.New("事务提交失败")
	}
	return nil
}

// GetList server row
func (s Server) GetList() (Servers, error) {
	rows, err := DB.Query("SELECT id, name, ip, port, owner, group_id, create_time, update_time FROM server WHERE state = 1")
	if err != nil {
		return nil, err
	}
	var servers Servers
	for rows.Next() {
		var server Server

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.GroupID, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

// GetListByManagerGroupStr server row
func (s Server) GetListByManagerGroupStr(managerGroupStr string) (Servers, error) {
	sql := "SELECT id, name, ip, port, owner, group_id, create_time, update_time FROM server WHERE state = 1"
	if managerGroupStr == "" {
		return nil, nil
	} else if managerGroupStr != "all" {
		sql += " AND group_id IN (" + managerGroupStr + ")"
	}
	sql += " ORDER BY id DESC"
	rows, err := DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var servers Servers
	for rows.Next() {
		var server Server

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Port, &server.Owner, &server.GroupID, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

// GetAll server row
func (s Server) GetAll() (Servers, error) {
	rows, err := DB.Query("SELECT id, name, ip, owner, group_id, create_time, update_time FROM server WHERE state = 1 ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	var servers Servers
	for rows.Next() {
		var server Server

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Owner, &server.GroupID, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

// GetData add server information to s *Server
func (s Server) GetData() (Server, error) {
	var server Server
	err := DB.QueryRow("SELECT name, ip, owner, group_id, create_time, update_time FROM server WHERE id = ?", s.ID).Scan(&server.Name, &server.IP, &server.Owner, &server.GroupID, &server.CreateTime, &server.UpdateTime)
	if err != nil {
		return server, errors.New("数据查询失败")
	}
	return server, nil
}
