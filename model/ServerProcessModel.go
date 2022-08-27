// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const serverProcessTable = "`server_process`"

type ServerProcess struct {
	ID         int64  `json:"id"`
	ServerID   int64  `json:"serverId"`
	Name       string `json:"name"`
	Start      string `json:"start"`
	Stop       string `json:"stop"`
	Status     string `json:"status"`
	Restart    string `json:"restart"`
	InsertTime string `json:"insertTime,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

type ServerProcesses []ServerProcess

func (sp ServerProcess) GetData() (ServerProcess, error) {
	var serverProcess ServerProcess
	err := sq.
		Select("id, server_id, name, start, stop, status, restart").
		From(serverProcessTable).
		Where(sq.Eq{"id": sp.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&serverProcess.ID,
			&serverProcess.ServerID,
			&serverProcess.Name,
			&serverProcess.Start,
			&serverProcess.Stop,
			&serverProcess.Status,
			&serverProcess.Restart)
	if err != nil {
		return serverProcess, err
	}
	return serverProcess, nil
}

func (sp ServerProcess) GetListByServerID() (ServerProcesses, error) {
	rows, err := sq.
		Select("id, server_id, name, start, stop, status, restart, insert_time, update_time").
		From(serverProcessTable).
		Where(sq.Eq{"server_id": sp.ServerID}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	serverProcesses := ServerProcesses{}
	for rows.Next() {
		var serverProcess ServerProcess

		if err := rows.Scan(
			&serverProcess.ID,
			&serverProcess.ServerID,
			&serverProcess.Name,
			&serverProcess.Start,
			&serverProcess.Stop,
			&serverProcess.Status,
			&serverProcess.Restart,
			&serverProcess.InsertTime,
			&serverProcess.UpdateTime,
		); err != nil {
			return serverProcesses, err
		}
		serverProcesses = append(serverProcesses, serverProcess)
	}
	return serverProcesses, nil
}

func (sp ServerProcess) AddRow() (int64, error) {
	result, err := sq.
		Insert(serverProcessTable).
		Columns("server_id", "name", "start", "stop", "status", "restart").
		Values(sp.ServerID, sp.Name, sp.Start, sp.Stop, sp.Status, sp.Restart).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (sp ServerProcess) EditRow() error {
	_, err := sq.
		Update(serverProcessTable).
		SetMap(sq.Eq{
			"name":    sp.Name,
			"start":   sp.Start,
			"stop":    sp.Stop,
			"status":  sp.Status,
			"restart": sp.Restart,
		}).
		Where(sq.Eq{"id": sp.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (sp ServerProcess) DeleteRow() error {
	_, err := sq.
		Delete(serverProcessTable).
		Where(sq.Eq{"id": sp.ID}).
		RunWith(DB).
		Exec()
	return err
}
