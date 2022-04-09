// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const projectTaskTable = "`project_task`"

// state type
const (
	NotRun = iota
	Run
)

// ProjectTask -
type ProjectTask struct {
	ID         int64  `json:"id"`
	ProjectID  int64  `json:"projectId"`
	Branch     string `json:"branch"`
	CommitID   string `json:"commit"`
	Date       string `json:"date"`
	State      uint8  `json:"state"`
	IsRun      uint8  `json:"isRun"`
	Creator    string `json:"creator"`
	CreatorID  int64  `json:"creatorId"`
	Editor     string `json:"editor"`
	EditorID   int64  `json:"editorId"`
	InsertTime string `json:"insertTime"`
	UpdateTime string `json:"updateTime"`
}

// ProjectTasks -
type ProjectTasks []ProjectTask

// GetListByProjectID -
func (pt ProjectTask) GetListByProjectID(pagination Pagination) (ProjectTasks, Pagination, error) {
	rows, err := sq.
		Select("id, project_id, branch, commit, date, is_run, state, creator, creator_id, editor, editor_id, insert_time, update_time").
		From(projectTaskTable).
		Where(sq.Eq{"project_id": pt.ProjectID}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()

	if err != nil {
		return nil, pagination, err
	}
	projectTasks := ProjectTasks{}
	for rows.Next() {
		var projectTask ProjectTask

		if err := rows.Scan(
			&projectTask.ID,
			&projectTask.ProjectID,
			&projectTask.Branch,
			&projectTask.CommitID,
			&projectTask.Date,
			&projectTask.IsRun,
			&projectTask.State,
			&projectTask.Creator,
			&projectTask.CreatorID,
			&projectTask.Editor,
			&projectTask.EditorID,
			&projectTask.InsertTime,
			&projectTask.UpdateTime,
		); err != nil {
			return projectTasks, pagination, err
		}
		projectTasks = append(projectTasks, projectTask)
	}
	err = sq.
		Select("COUNT(*) AS count").
		From(projectTaskTable).
		Where(sq.Eq{"project_id": pt.ProjectID}).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return projectTasks, pagination, err
	}
	return projectTasks, pagination, nil
}

// GetNotRunListLTDate -
func (pt ProjectTask) GetNotRunListLTDate(date string) (ProjectTasks, error) {
	rows, err := sq.
		Select("id, project_id, commit, date").
		From(projectTaskTable).
		Where(sq.LtOrEq{"date": date}).
		Where(sq.Eq{"state": Enable, "is_run": NotRun}).
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	projectTasks := ProjectTasks{}
	for rows.Next() {
		var projectTask ProjectTask
		if err := rows.Scan(
			&projectTask.ID,
			&projectTask.ProjectID,
			&projectTask.CommitID,
			&projectTask.Date,
		); err != nil {
			return nil, err
		}
		projectTasks = append(projectTasks, projectTask)
	}

	return projectTasks, nil
}

// AddRow -
func (pt ProjectTask) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectTaskTable).
		Columns("project_id", "branch", "commit", "date", "creator", "creator_id").
		Values(pt.ProjectID, pt.Branch, pt.CommitID, pt.Date, pt.Creator, pt.CreatorID).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// SetRun set project task to run
func (pt ProjectTask) SetRun() error {
	_, err := sq.
		Update(projectTaskTable).
		SetMap(sq.Eq{
			"is_run": Run,
		}).
		Where(sq.Eq{"id": pt.ID}).
		RunWith(DB).
		Exec()
	return err
}

// RemoveRow -
func (pt ProjectTask) RemoveRow() error {
	_, err := sq.
		Update(projectTaskTable).
		SetMap(sq.Eq{
			"state":     Disable,
			"editor":    pt.Editor,
			"editor_id": pt.EditorID,
		}).
		Where(sq.Eq{"id": pt.ID}).
		RunWith(DB).
		Exec()
	return err
}
