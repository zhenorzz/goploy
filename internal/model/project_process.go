// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const projectProcessTable = "`project_process`"

type ProjectProcess struct {
	ID         int64  `json:"id"`
	ProjectID  int64  `json:"projectId"`
	Name       string `json:"name"`
	Start      string `json:"start"`
	Stop       string `json:"stop"`
	Status     string `json:"status"`
	Restart    string `json:"restart"`
	InsertTime string `json:"insertTime,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

type ProjectProcesses []ProjectProcess

// GetData -
func (pp ProjectProcess) GetData() (ProjectProcess, error) {
	var projectProcess ProjectProcess
	err := sq.
		Select("id, project_id, name, start, stop, status, restart").
		From(projectProcessTable).
		Where(sq.Eq{"id": pp.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&projectProcess.ID,
			&projectProcess.ProjectID,
			&projectProcess.Name,
			&projectProcess.Start,
			&projectProcess.Stop,
			&projectProcess.Status,
			&projectProcess.Restart)
	if err != nil {
		return projectProcess, err
	}
	return projectProcess, nil
}

func (pp ProjectProcess) GetListByProjectID(page, limit uint64) (ProjectProcesses, error) {
	rows, err := sq.
		Select("id, project_id, name, start, stop, status, restart, insert_time, update_time").
		From(projectProcessTable).
		Where(sq.Eq{"project_id": pp.ProjectID}).
		Limit(limit).
		Offset((page - 1) * limit).
		OrderBy("id DESC").
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	projectProcesses := ProjectProcesses{}
	for rows.Next() {
		var projectProcess ProjectProcess

		if err := rows.Scan(
			&projectProcess.ID,
			&projectProcess.ProjectID,
			&projectProcess.Name,
			&projectProcess.Start,
			&projectProcess.Stop,
			&projectProcess.Status,
			&projectProcess.Restart,
			&projectProcess.InsertTime,
			&projectProcess.UpdateTime,
		); err != nil {
			return projectProcesses, err
		}
		projectProcesses = append(projectProcesses, projectProcess)
	}
	return projectProcesses, nil
}

func (pp ProjectProcess) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectProcessTable).
		Columns("project_id", "name", "start", "stop", "status", "restart").
		Values(pp.ProjectID, pp.Name, pp.Start, pp.Stop, pp.Status, pp.Restart).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (pp ProjectProcess) EditRow() error {
	_, err := sq.
		Update(projectProcessTable).
		SetMap(sq.Eq{
			"name":    pp.Name,
			"start":   pp.Start,
			"stop":    pp.Stop,
			"status":  pp.Status,
			"restart": pp.Restart,
		}).
		Where(sq.Eq{"id": pp.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (pp ProjectProcess) DeleteRow() error {
	_, err := sq.
		Delete(projectProcessTable).
		Where(sq.Eq{"id": pp.ID}).
		RunWith(DB).
		Exec()
	return err
}
