// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/zhenorzz/goploy/internal/pkg"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const serverTable = "`server`"

type Server struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	Owner        string `json:"owner"`
	Path         string `json:"path"`
	Password     string `json:"password"`
	JumpIP       string `json:"jumpIP"`
	JumpPort     int    `json:"jumpPort"`
	JumpOwner    string `json:"jumpOwner"`
	JumpPath     string `json:"jumpPath"`
	JumpPassword string `json:"jumpPassword"`
	NamespaceID  int64  `json:"namespaceId"`
	Description  string `json:"description"`
	OS           string `json:"os"`
	OSInfo       string `json:"osInfo"`
	State        int8   `json:"state"`
	InsertTime   string `json:"insertTime"`
	UpdateTime   string `json:"updateTime"`
}

// Project deploy state
const (
	ServerOSWindows = "windows"
	ServerOSLinux   = "linux"
)

type Servers []Server

func (s Server) GetList() (Servers, error) {
	rows, err := sq.
		Select(
			"id",
			"namespace_id",
			"name",
			"ip",
			"port",
			"owner",
			"path",
			"password",
			"jump_ip",
			"jump_port",
			"jump_owner",
			"jump_path",
			"jump_password",
			"description",
			"os",
			"os_info",
			"state",
			"insert_time",
			"update_time").
		From(serverTable).
		Where(sq.Eq{
			"namespace_id": []int64{0, s.NamespaceID},
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
		if err := rows.Scan(
			&server.ID,
			&server.NamespaceID,
			&server.Name,
			&server.IP,
			&server.Port,
			&server.Owner,
			&server.Path,
			&server.Password,
			&server.JumpIP,
			&server.JumpPort,
			&server.JumpOwner,
			&server.JumpPath,
			&server.JumpPassword,
			&server.Description,
			&server.OS,
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

func (s Server) GetAll() (Servers, error) {
	rows, err := sq.
		Select(
			"id",
			"name",
			"ip",
			"owner",
			"description",
			"insert_time",
			"update_time").
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
		if err := rows.Scan(
			&server.ID,
			&server.Name,
			&server.IP,
			&server.Owner,
			&server.Description,
			&server.InsertTime,
			&server.UpdateTime); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

func (s Server) GetData() (Server, error) {
	var server Server
	builder := sq.
		Select(
			"id",
			"namespace_id",
			"name",
			"os",
			"ip",
			"port",
			"owner",
			"path",
			"password",
			"jump_ip",
			"jump_port",
			"jump_owner",
			"jump_path",
			"jump_password",
			"state",
		).
		From(serverTable)

	if s.ID > 0 {
		builder = builder.Where(sq.Eq{"id": s.ID})
	}

	if s.Name != "" {
		builder = builder.Where(sq.Eq{"name": s.Name})
	}

	if s.IP != "" {
		builder = builder.Where(sq.Eq{"ip": s.IP})
	}

	err := builder.
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&server.ID,
			&server.NamespaceID,
			&server.Name,
			&server.OS,
			&server.IP,
			&server.Port,
			&server.Owner,
			&server.Path,
			&server.Password,
			&server.JumpIP,
			&server.JumpPort,
			&server.JumpOwner,
			&server.JumpPath,
			&server.JumpPassword,
			&server.State,
		)
	if err != nil {
		return server, err
	}
	return server, nil
}

// AddRow return LastInsertId
func (s Server) AddRow() (int64, error) {
	result, err := sq.
		Insert(serverTable).
		Columns(
			"namespace_id",
			"name",
			"os",
			"ip",
			"port",
			"owner",
			"password",
			"path",
			"jump_ip",
			"jump_port",
			"jump_owner",
			"jump_password",
			"jump_path",
			"description",
			"os_info").
		Values(
			s.NamespaceID,
			s.Name,
			s.OS,
			s.IP,
			s.Port,
			s.Owner,
			s.Password,
			s.Path,
			s.JumpIP,
			s.JumpPort,
			s.JumpOwner,
			s.JumpPassword,
			s.JumpPath,
			s.Description,
			s.OSInfo).
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
			"namespace_id":  s.NamespaceID,
			"name":          s.Name,
			"os":            s.OS,
			"ip":            s.IP,
			"port":          s.Port,
			"owner":         s.Owner,
			"password":      s.Password,
			"path":          s.Path,
			"jump_ip":       s.JumpIP,
			"jump_port":     s.JumpPort,
			"jump_owner":    s.JumpOwner,
			"jump_password": s.JumpPassword,
			"jump_path":     s.JumpPath,
			"description":   s.Description,
			"os_info":       s.OSInfo,
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
		return err
	}
	_, err = sq.
		Update(serverTable).
		SetMap(sq.Eq{
			"state": s.State,
		}).
		Where(sq.Eq{"id": s.ID}).
		RunWith(tx).
		Exec()

	if err == nil && s.State == Disable {
		_, err = sq.
			Delete(projectServerTable).
			Where(sq.Eq{"server_id": s.ID}).
			RunWith(tx).
			Exec()
	}

	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func (s Server) ToSSHConfig() pkg.SSHConfig {
	return pkg.SSHConfig{
		User:         s.Owner,
		Password:     s.Password,
		Path:         s.Path,
		Host:         s.IP,
		Port:         s.Port,
		JumpUser:     s.JumpOwner,
		JumpPassword: s.JumpPassword,
		JumpPath:     s.JumpPath,
		JumpHost:     s.JumpIP,
		JumpPort:     s.JumpPort,
	}
}

func (s Server) ReplaceVars(script string) string {
	scriptVars := map[string]string{
		"${SERVER_ID}":            strconv.FormatInt(s.ID, 10),
		"${SERVER_NAME}":          s.Name,
		"${SERVER_IP}":            s.IP,
		"${SERVER_PORT}":          strconv.Itoa(s.Port),
		"${SERVER_OWNER}":         s.Owner,
		"${SERVER_PASSWORD}":      s.Password,
		"${SERVER_PATH}":          s.Path,
		"${SERVER_JUMP_IP}":       s.JumpIP,
		"${SERVER_JUMP_PORT}":     strconv.Itoa(s.JumpPort),
		"${SERVER_JUMP_OWNER}":    s.JumpOwner,
		"${SERVER_JUMP_PASSWORD}": s.JumpPassword,
		"${SERVER_JUMP_PATH}":     s.JumpPath,
	}
	for key, value := range scriptVars {
		script = strings.Replace(script, key, value, -1)
	}
	return script
}
