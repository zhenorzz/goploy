package model

import (
	"database/sql"
	"errors"
)

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

	if err != nil {
		println(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	return uint32(id), err
}

// GetListByGitTraceID RemoteTrace row
func (rt RemoteTrace) GetListByGitTraceID() (RemoteTraces, error) {
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
		WHERE git_trace_id = ?`, rt.GitTraceID)
	if err != nil {
		return nil, err
	}
	var remoteTraces RemoteTraces
	for rows.Next() {
		var remoteTrace RemoteTrace

		if err := rows.Scan(
			&remoteTrace.ID,
			&remoteTrace.GitTraceID,
			&remoteTrace.ProjectID,
			&remoteTrace.ProjectName,
			&remoteTrace.ServerID,
			&remoteTrace.ServerName,
			&remoteTrace.Detail,
			&remoteTrace.State,
			&remoteTrace.PublisherID,
			&remoteTrace.PublisherName,
			&remoteTrace.Type,
			&remoteTrace.CreateTime,
			&remoteTrace.UpdateTime); err == sql.ErrNoRows {
			return remoteTraces, errors.New("项目尚无远程同步记录")
		} else if err != nil {
			return nil, errors.New("RemoteTrace.GetListByGitTraceID数据查询失败")
		}
		remoteTraces = append(remoteTraces, remoteTrace)
	}
	return remoteTraces, nil
}

// CountFailStateByGitTraceID fail num
func (rt RemoteTrace) CountFailStateByGitTraceID() (uint, error) {
	var num uint

	err := DB.QueryRow("SELECT count(*) FROM remote_trace WHERE git_trace_id = ? AND state = 0", rt.GitTraceID).Scan(&num)
	if err != nil {
		return 0, errors.New("数据查询失败")
	}
	return num, nil
}
