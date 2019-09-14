package model

import sq "github.com/Masterminds/squirrel"

const projectServerTable = "`project_server`"

// ProjectServer project server relationship
type ProjectServer struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"projectId"`
	ServerID    int64  `json:"serverId"`
	ServerName  string `json:"serverName"`
	ServerIP    string `json:"serverIP"`
	ServerPort  int64  `json:"serverPort"`
	ServerOwner string `json:"serverOwner"`
	CreateTime  int64  `json:"createTime"`
	UpdateTime  int64  `json:"updateTime"`
}

// ProjectServers project server relationship
type ProjectServers []ProjectServer

// GetBindServerListByProjectID server row
func (ps ProjectServer) GetBindServerListByProjectID() (ProjectServers, error) {
	rows, err := sq.
		Select("project_server.id, project_id, server_id, server.name, server.ip, server.port, server.owner, project_server.create_time, project_server.update_time").
		From(projectServerTable).
		LeftJoin(serverTable + " project_server.server_id = server.id").
		Where(sq.Eq{"project_id": ps.ProjectID}).
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	var projectServers ProjectServers
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
			&projectServer.CreateTime,
			&projectServer.UpdateTime); err != nil {
			return projectServers, err
		}
		projectServers = append(projectServers, projectServer)
	}
	return projectServers, nil
}

// AddMany add many row to table project_server
func (ps ProjectServers) AddMany() error {
	if len(ps) == 0 {
		return nil
	}

	builder := sq.
		Insert(projectServerTable).
		Columns("project_id", "server_id", "create_time", "update_time")

	for _, row := range ps {
		builder.Values(row.ProjectID, row.ServerID, row.CreateTime, row.UpdateTime)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

// DeleteRow edit one row to table server
func (ps ProjectServer) DeleteRow() error {
	_, err := sq.
		Delete(projectServerTable).
		Where(sq.Eq{"id": ps.ID}).
		RunWith(DB).
		Exec()
	return err
}
