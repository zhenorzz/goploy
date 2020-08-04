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

// ProjectTask project user relationship
type ProjectTask struct {
	ID         int64  `json:"id"`
	ProjectID  int64  `json:"projectId"`
	CommitID   string `json:"commitId"`
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

// ProjectTasks project task relationship
type ProjectTasks []ProjectTask

// GetListByProjectID project task row
func (pt ProjectTask) GetListByProjectID(pagination Pagination) (ProjectTasks, Pagination, error) {
	rows, err := sq.
		Select("id, project_id, commit_id, date, is_run, state, creator, creator_id, editor, editor_id, insert_time, update_time").
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
			return nil, pagination, err
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
		return nil, pagination, err
	}
	return projectTasks, pagination, nil
}

// GetNotRunListLTDate project task row
func (pt ProjectTask) GetNotRunListLTDate(date string) (ProjectTasks, error) {
	rows, err := sq.
		Select("id, project_id, commit_id, date").
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

// AddRow add one row to table server
func (pt ProjectTask) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectTaskTable).
		Columns("project_id", "commit_id", "date", "creator", "creator_id").
		Values(pt.ProjectID, pt.CommitID, pt.Date, pt.Creator, pt.CreatorID).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table server
func (pt ProjectTask) EditRow() error {
	_, err := sq.
		Update(projectTaskTable).
		SetMap(sq.Eq{
			"commit_id": pt.CommitID,
			"date":      pt.Date,
			"editor":    pt.Editor,
			"editor_id": pt.EditorID,
		}).
		Where(sq.Eq{"id": pt.ID}).
		RunWith(DB).
		Exec()
	return err
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

// RemoveRow project task
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
