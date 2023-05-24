// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/zhenorzz/goploy/internal/pkg"
	"strconv"
	"strings"
)

const projectServerTable = "`project_server`"

type ProjectServer struct {
	ID                 int64  `json:"id"`
	ProjectID          int64  `json:"projectId"`
	ServerID           int64  `json:"serverId"`
	ServerName         string `json:"serverName"`
	ServerOS           string `json:"serverOS"`
	ServerIP           string `json:"serverIP"`
	ServerPort         int    `json:"serverPort"`
	ServerOwner        string `json:"serverOwner"`
	ServerPassword     string `json:"serverPassword"`
	ServerPath         string `json:"serverPath"`
	ServerJumpIP       string `json:"serverJumpIP"`
	ServerJumpPort     int    `json:"serverJumpPort"`
	ServerJumpOwner    string `json:"serverJumpOwner"`
	ServerJumpPassword string `json:"serverJumpPassword"`
	ServerJumpPath     string `json:"serverJumpPath"`
	ServerDescription  string `json:"serverDescription"`
	InsertTime         string `json:"insertTime"`
	UpdateTime         string `json:"updateTime"`
}

type ProjectServers []ProjectServer

func (ps ProjectServer) GetBindServerListByProjectID() (ProjectServers, error) {
	rows, err := sq.
		Select(`
			project_server.id, 
			project_id, 
			server_id, 
			server.name, 
			server.os, 
			server.ip, 
			server.port, 
			server.owner, 
			server.password, 
			server.path, 
			server.jump_ip, 
			server.jump_port, 
			server.jump_owner, 
			server.jump_password, 
			server.jump_path, 
			server.description,
			project_server.insert_time, 
			project_server.update_time`).
		From(projectServerTable).
		LeftJoin(serverTable + " ON project_server.server_id = server.id").
		Where(sq.Eq{"project_id": ps.ProjectID}).
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	projectServers := ProjectServers{}
	for rows.Next() {
		var projectServer ProjectServer

		if err := rows.Scan(
			&projectServer.ID,
			&projectServer.ProjectID,
			&projectServer.ServerID,
			&projectServer.ServerName,
			&projectServer.ServerOS,
			&projectServer.ServerIP,
			&projectServer.ServerPort,
			&projectServer.ServerOwner,
			&projectServer.ServerPassword,
			&projectServer.ServerPath,
			&projectServer.ServerJumpIP,
			&projectServer.ServerJumpPort,
			&projectServer.ServerJumpOwner,
			&projectServer.ServerJumpPassword,
			&projectServer.ServerJumpPath,
			&projectServer.ServerDescription,
			&projectServer.InsertTime,
			&projectServer.UpdateTime); err != nil {
			return nil, err
		}
		projectServers = append(projectServers, projectServer)
	}
	return projectServers, nil
}

func (ps ProjectServers) AddMany() error {
	if len(ps) == 0 {
		return nil
	}

	builder := sq.
		Replace(projectServerTable).
		Columns("project_id", "server_id")

	for _, row := range ps {
		builder = builder.Values(row.ProjectID, row.ServerID)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

func (ps ProjectServer) DeleteRow() error {
	_, err := sq.
		Delete(projectServerTable).
		Where(sq.Eq{"id": ps.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (ps ProjectServer) DeleteByProjectID() error {
	_, err := sq.
		Delete(projectServerTable).
		Where(sq.Eq{"project_id": ps.ProjectID}).
		RunWith(DB).
		Exec()
	return err
}

func (ps ProjectServer) ToSSHConfig() pkg.SSHConfig {
	return pkg.SSHConfig{
		User:         ps.ServerOwner,
		Password:     ps.ServerPassword,
		Path:         ps.ServerPath,
		Host:         ps.ServerIP,
		Port:         ps.ServerPort,
		JumpUser:     ps.ServerJumpOwner,
		JumpPassword: ps.ServerJumpPassword,
		JumpPath:     ps.ServerJumpPath,
		JumpHost:     ps.ServerJumpIP,
		JumpPort:     ps.ServerJumpPort,
	}
}

func (ps ProjectServer) ToSSHOption() string {
	proxyCommand := ""
	if ps.ServerJumpIP != "" {
		if ps.ServerJumpPath != "" {
			if ps.ServerJumpPassword != "" {
				proxyCommand = fmt.Sprintf("-o ProxyCommand='sshpass -p %s -P assphrase ssh -o StrictHostKeyChecking=no -W %%h:%%p -i %s -p %d %s@%s' ", ps.ServerJumpPassword, ps.ServerJumpPath, ps.ServerJumpPort, ps.ServerJumpOwner, ps.ServerJumpIP)
			} else {
				proxyCommand = fmt.Sprintf("-o ProxyCommand='ssh -o StrictHostKeyChecking=no -W %%h:%%p -i %s -p %d %s@%s' ", ps.ServerJumpPath, ps.ServerJumpPort, ps.ServerJumpOwner, ps.ServerJumpIP)
			}
		} else {
			proxyCommand = fmt.Sprintf("-o ProxyCommand='sshpass -p %s ssh -o StrictHostKeyChecking=no -W %%h:%%p -p %d %s@%s' ", ps.ServerJumpPassword, ps.ServerJumpPort, ps.ServerJumpOwner, ps.ServerJumpIP)
		}
	}
	if ps.ServerPath != "" {
		if ps.ServerPassword != "" {
			return fmt.Sprintf("sshpass -p %s -P assphrase ssh -o StrictHostKeyChecking=no %s -p %d -i %s", ps.ServerPassword, proxyCommand, ps.ServerPort, ps.ServerPath)
		} else {
			return fmt.Sprintf("ssh -o StrictHostKeyChecking=no %s -p %d -i %s", proxyCommand, ps.ServerPort, ps.ServerPath)
		}
	} else {
		return fmt.Sprintf("sshpass -p %s ssh -o StrictHostKeyChecking=no %s -p %d", ps.ServerPassword, proxyCommand, ps.ServerPort)
	}
}

func (ps ProjectServer) ReplaceVars(script string) string {
	scriptVars := map[string]string{
		"${SERVER_ID}":            strconv.FormatInt(ps.ServerID, 10),
		"${SERVER_NAME}":          ps.ServerName,
		"${SERVER_IP}":            ps.ServerIP,
		"${SERVER_PORT}":          strconv.Itoa(ps.ServerPort),
		"${SERVER_OWNER}":         ps.ServerOwner,
		"${SERVER_PASSWORD}":      ps.ServerPassword,
		"${SERVER_PATH}":          ps.ServerPath,
		"${SERVER_JUMP_IP}":       ps.ServerJumpIP,
		"${SERVER_JUMP_PORT}":     strconv.Itoa(ps.ServerJumpPort),
		"${SERVER_JUMP_OWNER}":    ps.ServerJumpOwner,
		"${SERVER_JUMP_PASSWORD}": ps.ServerJumpPassword,
		"${SERVER_JUMP_PATH}":     ps.ServerJumpPath,
	}
	for key, value := range scriptVars {
		script = strings.Replace(script, key, value, -1)
	}
	return script
}
