// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/zhenorzz/goploy/internal/pkg"
)

const projectServerTable = "`project_server`"

type ProjectServer struct {
	ID         int64   `json:"id"`
	ProjectID  int64   `json:"projectId"`
	ServerID   int64   `json:"serverId"`
	Project    Project `json:"project,omitempty"`
	Server     Server  `json:"server,omitempty"`
	InsertTime string  `json:"insertTime"`
	UpdateTime string  `json:"updateTime"`
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
			&projectServer.Server.Name,
			&projectServer.Server.OS,
			&projectServer.Server.IP,
			&projectServer.Server.Port,
			&projectServer.Server.Owner,
			&projectServer.Server.Password,
			&projectServer.Server.Path,
			&projectServer.Server.JumpIP,
			&projectServer.Server.JumpPort,
			&projectServer.Server.JumpOwner,
			&projectServer.Server.JumpPassword,
			&projectServer.Server.JumpPath,
			&projectServer.Server.Description,
			&projectServer.InsertTime,
			&projectServer.UpdateTime); err != nil {
			return nil, err
		}
		projectServers = append(projectServers, projectServer)
	}
	return projectServers, nil
}

func (ps ProjectServer) GetBindProjectListByServerID() (ProjectServers, error) {
	rows, err := sq.
		Select(`
			project_server.id, 
			project_id, 
			server_id, 
			project.name, 
			project.branch, 
			project.environment, 
			project_server.insert_time, 
			project_server.update_time`).
		From(projectServerTable).
		LeftJoin(projectTable + " ON project_server.project_id = project.id").
		Where(sq.Eq{"server_id": ps.ServerID}).
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
			&projectServer.Project.Name,
			&projectServer.Project.Branch,
			&projectServer.Project.Environment,
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

func (ps ProjectServer) DeleteInID(id []int64) error {
	_, err := sq.
		Delete(projectServerTable).
		Where(sq.Eq{"id": id}).
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
		User:         ps.Server.Owner,
		Password:     ps.Server.Password,
		Path:         ps.Server.Path,
		Host:         ps.Server.IP,
		Port:         ps.Server.Port,
		JumpUser:     ps.Server.JumpOwner,
		JumpPassword: ps.Server.JumpPassword,
		JumpPath:     ps.Server.JumpPath,
		JumpHost:     ps.Server.JumpIP,
		JumpPort:     ps.Server.JumpPort,
	}
}

func (ps ProjectServer) ToSSHOption() string {
	proxyCommand := ""
	if ps.Server.JumpIP != "" {
		if ps.Server.JumpPath != "" {
			if ps.Server.JumpPassword != "" {
				proxyCommand = fmt.Sprintf(`-o ProxyCommand="sshpass -p %s -P assphrase ssh -o StrictHostKeyChecking=no -W %%h:%%p -i %s -p %d %s@%s" `, quote(ps.Server.JumpPassword), ps.Server.JumpPath, ps.Server.JumpPort, ps.Server.JumpOwner, ps.Server.JumpIP)
			} else {
				proxyCommand = fmt.Sprintf(`-o ProxyCommand="ssh -o StrictHostKeyChecking=no -W %%h:%%p -i %s -p %d %s@%s" `, ps.Server.JumpPath, ps.Server.JumpPort, ps.Server.JumpOwner, ps.Server.JumpIP)
			}
		} else {
			proxyCommand = fmt.Sprintf(`-o ProxyCommand="sshpass -p %s ssh -o StrictHostKeyChecking=no -W %%h:%%p -p %d %s@%s" `, quote(ps.Server.JumpPassword), ps.Server.JumpPort, ps.Server.JumpOwner, ps.Server.JumpIP)
		}
	}
	if ps.Server.Path != "" {
		if ps.Server.Password != "" {
			return fmt.Sprintf(`sshpass -p %s -P assphrase ssh -o StrictHostKeyChecking=no %s -p %d -i %s`, quote(ps.Server.Password), proxyCommand, ps.Server.Port, ps.Server.Path)
		} else {
			return fmt.Sprintf(`ssh -o StrictHostKeyChecking=no %s -p %d -i %s`, proxyCommand, ps.Server.Port, ps.Server.Path)
		}
	} else {
		return fmt.Sprintf(`sshpass -p %s ssh -o StrictHostKeyChecking=no %s -p %d`, quote(ps.Server.Password), proxyCommand, ps.Server.Port)
	}
}

func (ps ProjectServer) ReplaceVars(script string) string {
	return ps.Server.ReplaceVars(script)
}

func (ps ProjectServers) ReplaceVars(script string) string {
	var servers []Server
	for _, p := range ps {
		servers = append(servers, p.Server)
	}
	serverJson, _ := json.Marshal(servers)
	script = strings.Replace(script, "${PROJECT_SERVERS}", string(serverJson), -1)
	return script
}

// Note that doubling a single-quote inside a single-quoted string gives you a single-quote
// see `man rsync`
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
