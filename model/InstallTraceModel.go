package model

import sq "github.com/Masterminds/squirrel"

const installTraceTable = "`install_trace`"

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
	InsertTime   string `json:"insertTime"`
	UpdateTime   string `json:"updateTime"`
}

// InstallTraces InstallTrace slice
type InstallTraces []InstallTrace

// install trace type
const (
	Rsync  = 1
	SSH    = 2
	Script = 3
)

// AddRow add one row to table deploy and add id to deploy.ID
func (it InstallTrace) AddRow() (int64, error) {
	result, err := sq.
		Insert(installTraceTable).
		Columns("token", "server_id", "server_name", "detail", "state", "operator_id", "operator_name", "type", "ext").
		Values(it.Token, it.ServerID, it.ServerName, it.Detail, it.State, it.OperatorID, it.OperatorName, it.Type, it.Ext).
		RunWith(DB).
		Exec()

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

// GetListByToken InstallTrace row
func (it InstallTrace) GetListByToken() (InstallTraces, error) {
	rows, err := sq.
		Select("id, token, server_id, server_name, detail, state, operator_id, operator_name, type, ext, insert_time, update_time").
		From(installTraceTable).
		Where(sq.Eq{"token": it.Token}).
		RunWith(DB).
		Query()

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
			&installTrace.InsertTime,
			&installTrace.UpdateTime); err != nil {
			return nil, err
		}
		installTraces = append(installTraces, installTrace)
	}
	return installTraces, nil
}

// GetPreviewByProjectID InstallTrace row
func (it InstallTrace) GetPreviewByProjectID() (InstallTraces, error) {
	rows, err := sq.
		Select("id,token,server_id,server_name,detail,state,operator_id,operator_name,type,ext,create_time,update_time").
		Column("!EXISTS (SELECT id FROM " + installTraceTable + " AS it where it.state = 0 AND it.token = install_trace.token) as install_state").
		From(installTraceTable).
		Where(sq.Eq{"project_id": it.ServerID, "type": Rsync}).
		Limit(15).
		RunWith(DB).
		Query()

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
			&installTrace.InsertTime,
			&installTrace.UpdateTime,
			&installTrace.InstallState); err != nil {
			return nil, err
		}
		installTraces = append(installTraces, installTrace)
	}
	return installTraces, nil
}

// GetListGroupByToken InstallTrace token list
func (it InstallTrace) GetListGroupByToken() (InstallTraces, error) {
	rows, err := sq.
		Select("token, MIN(state) as install_state").
		From(installTraceTable).
		Where(sq.Eq{"server_id": it.ServerID}).
		GroupBy("token").
		RunWith(DB).
		Query()
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
