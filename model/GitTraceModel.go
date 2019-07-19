package model

import "errors"

// GitTrace mysql table for git trace
type GitTrace struct {
	ID            uint32 `json:"id"`
	Commit        string `json:"commit"`
	ProjectID     uint32 `json:"projectId"`
	ProjectName   string `json:"projectName"`
	Detail        string `json:"detail"`
	State         uint8  `json:"state"`
	RemoteState   uint8  `json:"remoteState"`
	PublisherID   uint32 `json:"publisherId"`
	PublisherName string `json:"publisherName"`
	CreateTime    int64  `json:"createTime"`
	UpdateTime    int64  `json:"updateTime"`
}

// GitTraces GitTrace slice
type GitTraces []GitTrace

// AddRow add one row to table deploy and add id to deploy.ID
func (gt GitTrace) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO git_trace (commit, project_id, project_name, detail, state, publisher_id, publisher_name, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		gt.Commit,
		gt.ProjectID,
		gt.ProjectName,
		gt.Detail,
		gt.State,
		gt.PublisherID,
		gt.PublisherName,
		gt.CreateTime,
		gt.UpdateTime,
	)
	id, err := result.LastInsertId()
	return uint32(id), err
}

// GetData get GitTrace information
func (gt GitTrace) GetData() (GitTrace, error) {
	var gitTrace GitTrace
	err := DB.QueryRow(`SELECT 
	        id,
	        commit,
			project_id, 
			project_name, 
			detail, 
			state, 
			publisher_id, 
			publisher_name, 
			create_time, 
			update_time
		FROM git_trace WHERE id = ?`, gt.ID).Scan(
		&gitTrace.ID,
		&gitTrace.Commit,
		&gitTrace.ProjectID,
		&gitTrace.ProjectName,
		&gitTrace.Detail,
		&gitTrace.State,
		&gitTrace.PublisherID,
		&gitTrace.PublisherName,
		&gitTrace.CreateTime,
		&gitTrace.UpdateTime)
	if err != nil {
		return gitTrace, err
	}
	return gitTrace, nil
}

// GetLatestRow add GitTrace information to gt *GitTrace
func (gt GitTrace) GetLatestRow() (GitTrace, error) {
	var gitTrace GitTrace
	err := DB.QueryRow(`SELECT 
	        id,
	        commit,
			project_id, 
			project_name, 
			detail, 
			state, 
			publisher_id, 
			publisher_name, 
			create_time, 
			update_time
		FROM git_trace WHERE project_id = ? ORDER BY id DESC Limit 1`, gt.ProjectID).Scan(
		&gitTrace.ID,
		&gitTrace.Commit,
		&gitTrace.ProjectID,
		&gitTrace.ProjectName,
		&gitTrace.Detail,
		&gitTrace.State,
		&gitTrace.PublisherID,
		&gitTrace.PublisherName,
		&gitTrace.CreateTime,
		&gitTrace.UpdateTime)
	if err != nil {
		return gitTrace, err
	}
	return gitTrace, nil
}

// GetListByProjectID GitTrace row
func (gt GitTrace) GetListByProjectID() (GitTraces, error) {
	rows, err := DB.Query(
		`SELECT 
			id,
			commit,
			project_id,
			project_name,
			detail,
			state,
			publisher_id,
			publisher_name,
			!EXISTS (SELECT id FROM remote_trace where remote_trace.state = 0 AND git_trace.id = remote_trace.git_trace_id) as remote_state,
			create_time,
			update_time
		FROM git_trace
		WHERE project_id = ?
		ORDER BY id DESC
		LIMIT 15`, gt.ProjectID)
	if err != nil {
		return nil, err
	}
	var gitTraces GitTraces
	for rows.Next() {
		var gitTrace GitTrace

		if err := rows.Scan(
			&gitTrace.ID,
			&gitTrace.Commit,
			&gitTrace.ProjectID,
			&gitTrace.ProjectName,
			&gitTrace.Detail,
			&gitTrace.State,
			&gitTrace.PublisherID,
			&gitTrace.PublisherName,
			&gitTrace.RemoteState,
			&gitTrace.CreateTime,
			&gitTrace.UpdateTime); err != nil {
			return nil, errors.New("GitTrace.GetListByProjectID数据查询失败")
		}
		gitTraces = append(gitTraces, gitTrace)
	}
	return gitTraces, nil
}
