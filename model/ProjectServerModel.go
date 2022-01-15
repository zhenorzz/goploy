package model

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/zhenorzz/goploy/utils"
)

const projectServerTable = "`project_server`"

// ProjectServer -
type ProjectServer struct {
	ID                 int64  `json:"id"`
	ProjectID          int64  `json:"projectId"`
	ServerID           int64  `json:"serverId"`
	ServerName         string `json:"serverName"`
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

// ProjectServers -
type ProjectServers []ProjectServer

// GetBindServerListByProjectID -
func (ps ProjectServer) GetBindServerListByProjectID() (ProjectServers, error) {
	rows, err := sq.
		Select(`
			project_server.id, 
			project_id, 
			server_id, 
			server.name, 
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

// AddMany -
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

// DeleteRow -
func (ps ProjectServer) DeleteRow() error {
	_, err := sq.
		Delete(projectServerTable).
		Where(sq.Eq{"id": ps.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (ps ProjectServer) Convert2SSHConfig() utils.SSHConfig {
	return utils.SSHConfig{
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
