package model

// ProjectServer project server relationship
type ProjectServer struct {
	ID          uint32 `json:"id"`
	ProjectID   uint32 `json:"projectId"`
	ServerID    uint32 `json:"serverId"`
	ServerName  string `json:"serverName"`
	ServerIP    string `json:"serverIP"`
	ServerOwner string `json:"serverOwner"`
	CreateTime  int64  `json:"createTime"`
	UpdateTime  int64  `json:"updateTime"`
}

// ProjectServers project server relationship
type ProjectServers []ProjectServer

// GetBindServerListByProjectID server row
func (ps ProjectServer) GetBindServerListByProjectID(projectID uint32) (ProjectServers, error) {
	db := NewDB()
	rows, err := db.Query(
		`SELECT 
			project_server.id,
			project_id,
			server_id,
			server.name,
			server.ip,
			server.owner,
			project_server.create_time,
			project_server.update_time
		FROM project_server
		LEFT JOIN server 
		ON project_server.server_id = server.id
		WHERE project_id = ?`, projectID)
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
	db := NewDB()
	sqlStr := "INSERT INTO project_server (project_id, server_id, create_time, update_time) VALUES "
	vals := []interface{}{}

	for _, row := range ps {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, row.ProjectID, row.ServerID, row.CreateTime, row.UpdateTime)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	//format all vals at once
	_, err = stmt.Exec(vals...)
	return err
}

// DeleteRow edit one row to table server
func (ps ProjectServer) DeleteRow() error {
	db := NewDB()
	_, err := db.Exec(`DELETE FROM project_server WHERE id = ?`, ps.ID)
	return err
}
