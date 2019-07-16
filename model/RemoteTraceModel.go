package model

// RemoteTrace mysql table for rsync trace
type RemoteTrace struct {
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
	Type          uint32 `json:"type"`
	CreateTime    int64  `json:"createTime"`
	UpdateTime    int64  `json:"updateTime"`
}

// RemoteTraces RemoteTrace slice
type RemoteTraces []RemoteTrace

// AddRow add one row to table deploy and add id to deploy.ID
func (rt RemoteTrace) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO remote_trace (git_trace_id, project_id, project_name, server_id, server_name, detail, state, publisher_id, publisher_name, type, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rt.GitTraceID,
		rt.ProjectID,
		rt.ProjectName,
		rt.ServerID,
		rt.ServerName,
		rt.Detail,
		rt.State,
		rt.PublisherID,
		rt.PublisherName,
		rt.Type,
		rt.CreateTime,
		rt.UpdateTime,
	)
	id, err := result.LastInsertId()
	return uint32(id), err
}

// GetListByGitTraceID RemoteTrace row
func (rt RemoteTrace) GetListByGitTraceID(gitTraceID uint32) (RemoteTraces, error) {
	rows, err := DB.Query(
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
			type,
			create_time,
			update_time
		FROM remote_trace
		WHERE git_trace_id = ?`, gitTraceID)
	if err != nil {
		return nil, err
	}
	var RemoteTraces RemoteTraces
	for rows.Next() {
		var RemoteTrace RemoteTrace

		if err := rows.Scan(
			&RemoteTrace.ID,
			&RemoteTrace.GitTraceID,
			&RemoteTrace.ProjectID,
			&RemoteTrace.ProjectName,
			&RemoteTrace.ServerID,
			&RemoteTrace.ServerName,
			&RemoteTrace.Detail,
			&RemoteTrace.State,
			&RemoteTrace.PublisherID,
			&RemoteTrace.PublisherName,
			&RemoteTrace.Type,
			&RemoteTrace.CreateTime,
			&RemoteTrace.UpdateTime); err != nil {
			return RemoteTraces, err
		}
		RemoteTraces = append(RemoteTraces, RemoteTrace)
	}
	return RemoteTraces, nil
}
