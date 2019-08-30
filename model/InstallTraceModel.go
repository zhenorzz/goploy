package model

// InstallTrace mysql table for install trace
type InstallTrace struct {
	ID           int64  `json:"id"`
	Token        string `json:"token"`
	ServerID     int64  `json:"serverId"`
	ServerName   string `json:"serverName"`
	Detail       string `json:"detail"`
	State        uint8  `json:"state"`
	OperatorID   int64  `json:"operatorId"`
	OperatorName string `json:"operatorName"`
	Type         int64  `json:"type"`
	Ext          string `json:"ext"`
	InstallState int64  `json:"installState"`
	CreateTime   int64  `json:"createTime"`
	UpdateTime   int64  `json:"updateTime"`
}

// InstallTraces InstallTrace slice
type InstallTraces []InstallTrace

// install trace type
const (
	Rsync = 1

	Script = 2
)

// AddRow add one row to table deploy and add id to deploy.ID
func (it InstallTrace) AddRow() (int64, error) {
	result, err := DB.Exec(
		"INSERT INTO install_trace (token, server_id, server_name, detail, state, operator_id, operator_name, type, ext, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		it.Token,
		it.ServerID,
		it.ServerName,
		it.Detail,
		it.State,
		it.OperatorID,
		it.OperatorName,
		it.Type,
		it.Ext,
		it.CreateTime,
		it.UpdateTime,
	)

	if err != nil {
		println(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

// GetListByToken InstallTrace row
func (it InstallTrace) GetListByToken() (InstallTraces, error) {
	rows, err := DB.Query(
		`SELECT 
			id,
			token,
			server_id,
			server_name,
			detail,
			state,
			operator_id,
			operator_name,
			type,
			ext,
			create_time,
			update_time
		FROM install_trace
		WHERE token = ?`, it.Token)
	if err != nil {
		return nil, err
	}
	var installTraces InstallTraces
	for rows.Next() {
		var installTrace InstallTrace

		if err := rows.Scan(
			&installTrace.ID,
			&installTrace.Token,
			&installTrace.ServerID,
			&installTrace.ServerName,
			&installTrace.Detail,
			&installTrace.State,
			&installTrace.OperatorID,
			&installTrace.OperatorName,
			&installTrace.Type,
			&installTrace.Ext,
			&installTrace.CreateTime,
			&installTrace.UpdateTime); err != nil {
			return nil, err
		}
		installTraces = append(installTraces, installTrace)
	}
	return installTraces, nil
}

// GetPreviewByProjectID InstallTrace row
func (it InstallTrace) GetPreviewByProjectID() (InstallTraces, error) {
	rows, err := DB.Query(
		`SELECT 
			id,
			token,
			server_id,
			server_name,
			detail,
			state,
			operator_id,
			operator_name,
			type,
			ext,
			!EXISTS (SELECT id FROM install_trace AS it where it.state = 0 AND it.token = install_trace.token) as publish_state,
			create_time,
			update_time
		FROM 
			install_trace
		WHERE 
			project_id = ?
		AND
			type = ?
		ORDER BY id DESC
		LIMIT 15`, it.ServerID, Rsync)
	if err != nil {
		return nil, err
	}
	var installTraces InstallTraces
	for rows.Next() {
		var installTrace InstallTrace

		if err := rows.Scan(
			&installTrace.ID,
			&installTrace.Token,
			&installTrace.ServerID,
			&installTrace.ServerName,
			&installTrace.Detail,
			&installTrace.State,
			&installTrace.OperatorID,
			&installTrace.OperatorName,
			&installTrace.Type,
			&installTrace.Ext,
			&installTrace.InstallState,
			&installTrace.CreateTime,
			&installTrace.UpdateTime); err != nil {
			return nil, err
		}
		installTraces = append(installTraces, installTrace)
	}
	return installTraces, nil
}

// GetListGroupByToken InstallTrace token list
func (it InstallTrace) GetListGroupByToken() (InstallTraces, error) {
	rows, err := DB.Query(
		`SELECT 
			token,
			MIN(state) as install_state
		FROM 
			install_trace
		WHERE 
			server_id = ?
		GROUP BY token`, it.ServerID)
	if err != nil {
		return nil, err
	}
	var installTraces InstallTraces
	for rows.Next() {
		var installTrace InstallTrace

		if err := rows.Scan(
			&installTrace.Token,
			&installTrace.InstallState); err != nil {
			return nil, err
		}
		installTraces = append(installTraces, installTrace)
	}
	return installTraces, nil
}
