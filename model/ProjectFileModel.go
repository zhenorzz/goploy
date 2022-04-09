// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import sq "github.com/Masterminds/squirrel"

const projectFileTable = "`project_file`"

// ProjectFile -
type ProjectFile struct {
	ID         int64  `json:"id"`
	ProjectID  int64  `json:"projectId"`
	Filename   string `json:"filename"`
	InsertTime string `json:"insertTime"`
	UpdateTime string `json:"updateTime"`
}

// ProjectFiles -
type ProjectFiles []ProjectFile

// GetListByProjectID -
func (pf ProjectFile) GetListByProjectID() (ProjectFiles, error) {
	rows, err := sq.
		Select("id, project_id, filename, insert_time, update_time").
		From(projectFileTable).
		Where(sq.Eq{"project_id": pf.ProjectID}).
		RunWith(DB).
		Query()

	if err != nil {
		return nil, err
	}
	projectFiles := ProjectFiles{}
	for rows.Next() {
		var projectFile ProjectFile

		if err := rows.Scan(
			&projectFile.ID,
			&projectFile.ProjectID,
			&projectFile.Filename,
			&projectFile.InsertTime,
			&projectFile.UpdateTime); err != nil {
			return nil, err
		}
		projectFiles = append(projectFiles, projectFile)
	}
	return projectFiles, nil
}

// GetTotalByProjectID -
func (pf ProjectFile) GetTotalByProjectID() (int64, error) {
	var total int64
	err := sq.
		Select("COUNT(*) AS count").
		From(projectFileTable).
		Where(sq.Eq{"project_id": pf.ProjectID}).
		RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetData -
func (pf ProjectFile) GetData() (ProjectFile, error) {
	var projectFile ProjectFile
	err := sq.
		Select("id, project_id, filename, insert_time, update_time").
		From(projectFileTable).
		Where(sq.Eq{"id": pf.ID}).
		RunWith(DB).
		QueryRow().
		Scan(
			&projectFile.ID,
			&projectFile.ProjectID,
			&projectFile.Filename,
			&projectFile.InsertTime,
			&projectFile.UpdateTime)
	if err != nil {
		return projectFile, err
	}
	return projectFile, nil
}

// AddRow return LastInsertId
func (pf ProjectFile) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectFileTable).
		Columns("project_id", "filename").
		Values(pf.ProjectID, pf.Filename).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow -
func (pf ProjectFile) EditRow() error {
	_, err := sq.
		Update(projectFileTable).
		SetMap(sq.Eq{
			"filename": pf.Filename,
		}).
		Where(sq.Eq{"id": pf.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeleteRow -
func (pf ProjectFile) DeleteRow() error {
	_, err := sq.
		Delete(projectFileTable).
		Where(sq.Eq{"id": pf.ID}).
		RunWith(DB).
		Exec()
	return err
}
