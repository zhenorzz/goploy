// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const serverProcessTable = "`server_process`"

type ServerProcess struct {
	ID          int64  `json:"id"`
	NamespaceID int64  `json:"namespaceId"`
	Name        string `json:"name"`
	Items       string `json:"items"`
	InsertTime  string `json:"insertTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

type ServerProcesses []ServerProcess

type ServerProcessItem struct {
	Name    string
	Command string
}

type ServerProcessItems []ServerProcessItem

func (sp ServerProcess) GetData() (ServerProcess, error) {
	var serverProcess ServerProcess
	err := sq.
		Select("id, name, items").
		From(serverProcessTable).
		Where(sq.Eq{"id": sp.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&serverProcess.ID,
			&serverProcess.Name,
			&serverProcess.Items)
	if err != nil {
		return serverProcess, err
	}
	return serverProcess, nil
}

func (sp ServerProcess) GetList() (ServerProcesses, error) {
	rows, err := sq.
		Select("id, name, items, insert_time, update_time").
		From(serverProcessTable).
		Where(sq.Eq{"namespace_id": sp.NamespaceID}).
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
			&serverProcess.Name,
			&serverProcess.Items,
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
		Columns("namespace_id", "name", "items").
		Values(sp.NamespaceID, sp.Name, sp.Items).
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
			"name":  sp.Name,
			"items": sp.Items,
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
