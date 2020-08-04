package model

import (
	sq "github.com/Masterminds/squirrel"
)

const crontabServerTable = "`crontab_server`"

// CrontabServer crontab server relationship
type CrontabServer struct {
	ID                int64  `json:"id"`
	CrontabID         int64  `json:"crontabId"`
	ServerID          int64  `json:"serverId"`
	ServerName        string `json:"serverName"`
	ServerIP          string `json:"serverIP"`
	ServerPort        int64  `json:"serverPort"`
	ServerOwner       string `json:"serverOwner"`
	ServerDescription string `json:"serverDescription"`
	InsertTime        string `json:"insertTime"`
	UpdateTime        string `json:"updateTime"`
}

// CrontabServers crontab server relationship
type CrontabServers []CrontabServer

// GetAllByCrontabID -
func (cs CrontabServer) GetAllByCrontabID() (CrontabServers, error) {
	rows, err := sq.
		Select("id, crontab_id, server_id").
		From(crontabServerTable).
		Where(sq.Eq{"crontab_id": cs.CrontabID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	crontabServers := CrontabServers{}
	for rows.Next() {
		var crontabServer CrontabServer

		if err := rows.Scan(
			&crontabServer.ID,
			&crontabServer.CrontabID,
			&crontabServer.ServerID); err != nil {
			return nil, err
		}
		crontabServers = append(crontabServers, crontabServer)
	}
	return crontabServers, nil
}

// GetBindServerListByProjectID server row
func (cs CrontabServer) GetBindServerListByProjectID() (CrontabServers, error) {
	rows, err := sq.
		Select("crontab_server.id, crontab_id, server_id, server.name, server.ip, server.port, server.owner, server.description, crontab_server.insert_time, crontab_server.update_time").
		From(crontabServerTable).
		LeftJoin(serverTable + " ON crontab_server.server_id = server.id").
		Where(sq.Eq{"crontab_id": cs.CrontabID}).
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	crontabServers := CrontabServers{}
	for rows.Next() {
		var crontabServer CrontabServer

		if err := rows.Scan(
			&crontabServer.ID,
			&crontabServer.CrontabID,
			&crontabServer.ServerID,
			&crontabServer.ServerName,
			&crontabServer.ServerIP,
			&crontabServer.ServerPort,
			&crontabServer.ServerOwner,
			&crontabServer.ServerDescription,
			&crontabServer.InsertTime,
			&crontabServer.UpdateTime); err != nil {
			return nil, err
		}
		crontabServers = append(crontabServers, crontabServer)
	}
	return crontabServers, nil
}

// AddMany add many row to table project_server
func (cs CrontabServers) AddMany() error {
	if len(cs) == 0 {
		return nil
	}

	builder := sq.
		Insert(crontabServerTable).
		Columns("crontab_id", "server_id")

	for _, row := range cs {
		builder = builder.Values(row.CrontabID, row.ServerID)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

// DeleteRow edit one row to table server
func (cs CrontabServer) DeleteRow() error {
	_, err := sq.
		Delete(crontabServerTable).
		Where(sq.Eq{"id": cs.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeleteRowByCrontabID -
func (cs CrontabServer) DeleteRowByCrontabID() error {
	_, err := sq.
		Delete(crontabServerTable).
		Where(sq.Eq{"crontab_id": cs.CrontabID}).
		RunWith(DB).
		Exec()
	return err
}
