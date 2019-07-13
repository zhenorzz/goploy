package model

// RsyncTrace mysql table for rsync trace
type RsyncTrace struct {
	ID            uint32 `json:"id"`
	GitTraceID    uint32 `json:"gitTraceId"`
	ProjectID     uint32 `json:"projectId"`
	ProjectName   string `json:"projectName"`
	ServerID      uint32 `json:"serverId"`
	ServerName    string `json:"serverName"`
	Detail        string `json:"detail"`
	State         uint8  `json:"state"`
	PublisherID   uint32 `json:"publisherId"`
	PublisherName string `json:"publisherName"`
	CreateTime    int64  `json:"createTime"`
	UpdateTime    int64  `json:"updateTime"`
}

// RsyncTraces RsyncTrace slice
type RsyncTraces []RsyncTrace

// AddRow add one row to table deploy and add id to deploy.ID
func (rt *RsyncTrace) AddRow() error {
	db := NewDB()
	result, err := db.Exec(
		"INSERT INTO rsync_trace (git_trace_id, project_id, project_name, server_id, server_name, detail, state, publisher_id, publisher_name, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rt.GitTraceID,
		rt.ProjectID,
		rt.ProjectName,
		rt.ServerID,
		rt.ServerName,
		rt.Detail,
		rt.State,
		rt.PublisherID,
		rt.PublisherName,
		rt.CreateTime,
		rt.UpdateTime,
	)
	id, err := result.LastInsertId()
	rt.ID = uint32(id)
	return err
}

// Query RsyncTrace row
func (rts *RsyncTraces) Query(gitTraceID uint32) error {
	db := NewDB()
	rows, err := db.Query(
		`SELECT 
			id,
			git_trace_id,
			project_id,
			project_name,
			server_id,
			server_name,
			detail,
			state,
			publisher_id,
			publisher_name,
			create_time,
			update_time
		FROM rsync_trace
		WHERE git_trace_id = ?`, gitTraceID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var rsyncTrace RsyncTrace

		if err := rows.Scan(
			&rsyncTrace.ID,
			&rsyncTrace.GitTraceID,
			&rsyncTrace.ProjectID,
			&rsyncTrace.ProjectName,
			&rsyncTrace.ServerID,
			&rsyncTrace.ServerName,
			&rsyncTrace.Detail,
			&rsyncTrace.State,
			&rsyncTrace.PublisherID,
			&rsyncTrace.PublisherName,
			&rsyncTrace.CreateTime,
			&rsyncTrace.UpdateTime); err != nil {
			return err
		}
		*rts = append(*rts, rsyncTrace)
	}
	return nil
}
