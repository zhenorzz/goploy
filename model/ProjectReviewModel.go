// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const projectReviewTable = "`project_review`"

// ProjectReview -
type ProjectReview struct {
	ID         int64  `json:"id"`
	ProjectID  int64  `json:"projectId"`
	CommitID   string `json:"commitId"`
	Branch     string `json:"branch"`
	ReviewURL  string `json:"reviewURL"`
	State      uint8  `json:"state"`
	Creator    string `json:"creator"`
	CreatorID  int64  `json:"creatorId"`
	Editor     string `json:"editor"`
	EditorID   int64  `json:"editorId"`
	InsertTime string `json:"insertTime"`
	UpdateTime string `json:"updateTime"`
}

// ProjectReviews -
type ProjectReviews []ProjectReview

// GetListByProjectID -
func (pr ProjectReview) GetListByProjectID(pagination Pagination) (ProjectReviews, Pagination, error) {
	rows, err := sq.
		Select("id, project_id, branch, commit_id, state, creator, creator_id, editor, editor_id, insert_time, update_time").
		From(projectReviewTable).
		Where(sq.Eq{"project_id": pr.ProjectID}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()

	if err != nil {
		return nil, pagination, err
	}
	projectReviews := ProjectReviews{}
	for rows.Next() {
		var projectReview ProjectReview

		if err := rows.Scan(
			&projectReview.ID,
			&projectReview.ProjectID,
			&projectReview.Branch,
			&projectReview.CommitID,
			&projectReview.State,
			&projectReview.Creator,
			&projectReview.CreatorID,
			&projectReview.Editor,
			&projectReview.EditorID,
			&projectReview.InsertTime,
			&projectReview.UpdateTime,
		); err != nil {
			return projectReviews, pagination, err
		}
		projectReviews = append(projectReviews, projectReview)
	}
	err = sq.
		Select("COUNT(*) AS count").
		From(projectReviewTable).
		Where(sq.Eq{"project_id": pr.ProjectID}).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return projectReviews, pagination, err
	}
	return projectReviews, pagination, nil
}

// GetData -
func (pr ProjectReview) GetData() (ProjectReview, error) {
	var projectReview ProjectReview
	err := sq.
		Select("id, project_id, branch, commit_id, state, insert_time, update_time").
		From(projectReviewTable).
		Where(sq.Eq{"id": pr.ID}).
		RunWith(DB).
		QueryRow().
		Scan(
			&projectReview.ID,
			&projectReview.ProjectID,
			&projectReview.Branch,
			&projectReview.CommitID,
			&projectReview.State,
			&projectReview.InsertTime,
			&projectReview.UpdateTime)
	if err != nil {
		return projectReview, err
	}
	return projectReview, nil
}

// AddRow -
func (pr ProjectReview) AddRow() (int64, error) {
	result, err := sq.
		Insert(projectReviewTable).
		Columns("project_id", "branch", "commit_id", "review_url", "creator", "creator_id").
		Values(pr.ProjectID, pr.Branch, pr.CommitID, pr.ReviewURL, pr.Creator, pr.CreatorID).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow -
func (pr ProjectReview) EditRow() error {
	_, err := sq.
		Update(projectReviewTable).
		SetMap(sq.Eq{
			"state":     pr.State,
			"editor":    pr.Editor,
			"editor_id": pr.EditorID,
		}).
		Where(sq.Eq{"id": pr.ID}).
		RunWith(DB).
		Exec()
	return err
}
