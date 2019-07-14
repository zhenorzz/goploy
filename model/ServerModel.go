package model

import (
	"errors"
)

// Server mysql table server
type Server struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Owner      string `json:"owner"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// Servers many project
type Servers []Server

// AddRow add one row to table server and add id to s.ID
func (s Server) AddRow() (uint32, error) {
	db := NewDB()
	result, err := db.Exec(
		"INSERT INTO server (name, ip, owner, create_time, update_time) VALUES (?, ?, ?, ?, ?)",
		s.Name,
		s.IP,
		s.Owner,
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
	db := NewDB()
	_, err := db.Exec(
		`UPDATE server SET 
		  name = ?,
		  ip = ?,
		  owner = ?
		WHERE
		 id = ?`,
		s.Name,
		s.IP,
		s.Owner,
		s.ID,
	)
	return err
}

// GetList server row
func (s Server) GetList() (Servers, error) {
	db := NewDB()
	rows, err := db.Query("SELECT id, name, ip, owner, create_time, update_time FROM server")
	if err != nil {
		return nil, err
	}
	var servers Servers
	for rows.Next() {
		var server Server

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Owner, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

// GetAll server row
func (s Server) GetAll() (Servers, error) {
	db := NewDB()
	rows, err := db.Query("SELECT id, name, ip, owner, create_time, update_time FROM server")
	if err != nil {
		return nil, err
	}
	var servers Servers
	for rows.Next() {
		var server Server

		if err := rows.Scan(&server.ID, &server.Name, &server.IP, &server.Owner, &server.CreateTime, &server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

// QueryRow add server information to s *Server
func (s *Server) QueryRow() error {
	db := NewDB()
	err := db.QueryRow("SELECT name, ip, owner, create_time, update_time FROM server WHERE id = ?", s.ID).Scan(&s.Name, &s.IP, &s.Owner, &s.CreateTime, &s.UpdateTime)
	if err != nil {
		return errors.New("数据查询失败")
	}
	return nil
}
