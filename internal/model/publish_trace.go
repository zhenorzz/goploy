// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

const publishTraceTable = "`publish_trace`"

// PublishTrace -
type PublishTrace struct {
	ID            int64  `json:"id"`
	Token         string `json:"token"`
	NamespaceID   int64  `json:"namespaceId"`
	ProjectID     int64  `json:"projectId"`
	ProjectName   string `json:"projectName"`
	Detail        string `json:"detail"`
	State         int    `json:"state"`
	PublisherID   int64  `json:"publisherId"`
	PublisherName string `json:"publisherName"`
	Type          int    `json:"type"`
	Ext           string `json:"ext"`
	InsertTime    string `json:"insertTime"`
	UpdateTime    string `json:"updateTime"`
}

// PublishTraces -
type PublishTraces []PublishTrace

// publish trace state
const (
	Queue = iota
	BeforePull
	Pull
	AfterPull
	BeforeDeploy
	Deploy
	AfterDeploy
	DeployFinish
	PublishFinish
)

func (pt PublishTrace) GetList(page, limit uint64) (PublishTraces, error) {
	builder := sq.
		Select(
			fmt.Sprintf("%s.token", publishTraceTable),
			fmt.Sprintf("min(%s.project_id)", publishTraceTable),
			fmt.Sprintf("min(%s.project_name)", publishTraceTable),
			fmt.Sprintf("min(%s.publisher_id)", publishTraceTable),
			fmt.Sprintf("min(%s.publisher_name)", publishTraceTable),
			fmt.Sprintf("min(%s.state)", publishTraceTable),
			fmt.Sprintf("IFNULL(GROUP_CONCAT(IF(%s.state = 0, %[1]s.detail, NULL)), '') as detail", publishTraceTable),
			fmt.Sprintf("min(%s.insert_time) as insert_time", publishTraceTable),
		).
		From(publishTraceTable).
		GroupBy("token")

	if pt.NamespaceID > 0 {
		builder = builder.
			Join(fmt.Sprintf("%s ON %[1]s.id = %s.project_id", projectTable, publishTraceTable)).
			Where(sq.Eq{projectTable + ".namespace_id": pt.NamespaceID})
	}
	if pt.PublisherName != "" {
		builder = builder.Where(sq.Eq{publishTraceTable + ".publisher_name": pt.PublisherName})
	}
	if pt.ProjectName != "" {
		builder = builder.Where(sq.Eq{publishTraceTable + ".project_name": pt.ProjectName})
	}

	rows, err := builder.RunWith(DB).
		OrderBy("insert_time DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Query()
	if err != nil {
		return nil, err
	}
	publishTraces := PublishTraces{}
	for rows.Next() {
		var publishTrace PublishTrace

		if err := rows.Scan(
			&publishTrace.Token,
			&publishTrace.ProjectID,
			&publishTrace.ProjectName,
			&publishTrace.PublisherID,
			&publishTrace.PublisherName,
			&publishTrace.State,
			&publishTrace.Detail,
			&publishTrace.InsertTime); err != nil {
			return nil, err
		}
		publishTraces = append(publishTraces, publishTrace)
	}

	return publishTraces, nil
}

func (pt PublishTrace) GetTotal() (int64, error) {
	var total int64
	builder := sq.Select("COUNT(distinct token) AS count").
		From(publishTraceTable)
	if pt.NamespaceID > 0 {
		builder = builder.
			Join(fmt.Sprintf("%s ON %[1]s.id = %s.project_id", projectTable, publishTraceTable)).
			Where(sq.Eq{projectTable + ".namespace_id": pt.NamespaceID})
	}
	if pt.PublisherName != "" {
		builder = builder.Where(sq.Eq{publishTraceTable + ".publisher_name": pt.PublisherName})
	}
	if pt.ProjectName != "" {
		builder = builder.Where(sq.Eq{publishTraceTable + ".project_name": pt.ProjectName})
	}
	err := builder.RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (pt PublishTrace) GetListByToken() (PublishTraces, error) {
	rows, err := sq.
		Select(
			"id",
			"token",
			"project_id",
			"project_name",
			"if(state = 0,detail, '') as detail",
			"state",
			"publisher_id",
			"publisher_name",
			"type",
			"ext",
			"insert_time",
			"update_time").
		From(publishTraceTable).
		Where(sq.Eq{"token": pt.Token}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	publishTraces := PublishTraces{}
	for rows.Next() {
		var publishTrace PublishTrace

		if err := rows.Scan(
			&publishTrace.ID,
			&publishTrace.Token,
			&publishTrace.ProjectID,
			&publishTrace.ProjectName,
			&publishTrace.Detail,
			&publishTrace.State,
			&publishTrace.PublisherID,
			&publishTrace.PublisherName,
			&publishTrace.Type,
			&publishTrace.Ext,
			&publishTrace.InsertTime,
			&publishTrace.UpdateTime); err != nil {
			return nil, err
		}
		publishTraces = append(publishTraces, publishTrace)
	}
	return publishTraces, nil
}

func (pt PublishTrace) GetPreview(
	branch string,
	commit string,
	filename string,
	commitDate []string,
	deployDate []string,
	pagination Pagination,
) (PublishTraces, Pagination, error) {
	builder := sq.
		Select(
			"token",
			"MIN(publisher_name) publisher_name",
			"MIN(state) state",
			"GROUP_CONCAT(IF(type = 2 and ext != '', JSON_EXTRACT(ext, '$.commit') , '') SEPARATOR '') as ext",
			"MIN(insert_time) insert_time",
			"MAX(update_time) update_time",
		).
		From(publishTraceTable)
	if pt.ProjectID != 0 {
		builder = builder.Where(sq.Eq{"project_id": pt.ProjectID})
	}
	if pt.PublisherID != 0 {
		builder = builder.Where(sq.Eq{"publisher_id": pt.PublisherID})
	}
	if pt.Token != "" {
		builder = builder.Where(sq.Eq{"token": pt.Token})
	}
	if branch != "" {
		builder = builder.Where(sq.Like{"ext": "%" + branch + "%"})
	}
	if commit != "" {
		builder = builder.Where(sq.Like{"ext": "%" + commit + "%"})
	}
	if filename != "" {
		builder = builder.Where(sq.Like{"ext": "%" + filename + "%"})
	}
	if len(commitDate) > 1 {
		builder = builder.Where(`substring(ext, POSITION('"timestamp":' IN ext) + 12, 10) between ? and ?`, commitDate[0], commitDate[1])
	}
	if len(deployDate) > 1 {
		builder = builder.Where("insert_time between ? and ?", deployDate[0], deployDate[1])
	}
	if pt.State != -1 {
		builder = builder.Having(sq.Eq{"state": pt.State})
	}
	rows, err := builder.RunWith(DB).
		GroupBy("token").
		OrderBy("update_time DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		Query()
	if err != nil {
		return nil, pagination, err
	}
	publishTraces := PublishTraces{}
	for rows.Next() {
		var publishTrace PublishTrace

		if err := rows.Scan(
			&publishTrace.Token,
			&publishTrace.PublisherName,
			&publishTrace.State,
			&publishTrace.Ext,
			&publishTrace.InsertTime,
			&publishTrace.UpdateTime); err != nil {
			return nil, pagination, err
		}
		publishTraces = append(publishTraces, publishTrace)
	}

	builder = sq.
		Select("COUNT(*) AS count").
		From(publishTraceTable).
		Where(sq.Eq{"type": Pull})
	if pt.ProjectID != 0 {
		builder = builder.Where(sq.Eq{"project_id": pt.ProjectID})
	}
	if pt.PublisherID != 0 {
		builder = builder.Where(sq.Eq{"publisher_id": pt.PublisherID})
	}
	if branch != "" {
		builder = builder.Where(sq.Like{"ext": "%" + branch + "%"})
	}
	if commit != "" {
		builder = builder.Where(sq.Like{"ext": "%" + commit + "%"})
	}
	if filename != "" {
		builder = builder.Where(sq.Like{"ext": "%" + filename + "%"})
	}
	if len(commitDate) > 1 {
		builder = builder.Where(`substring(ext, POSITION('"timestamp":' IN ext) + 12, 10) between ? and ?`, commitDate[0], commitDate[1])
	}
	if len(deployDate) > 1 {
		builder = builder.Where("insert_time between ? and ?", deployDate[0], deployDate[1])
	}
	if pt.State == 0 {
		builder = builder.Where("EXISTS (SELECT id FROM " + publishTraceTable + " AS pt where pt.state = 0 AND pt.token = publish_trace.token)")
	} else if pt.State == 1 {
		builder = builder.Where("! EXISTS (SELECT id FROM " + publishTraceTable + " AS pt where pt.state = 0 AND pt.token = publish_trace.token)")
	}

	err = builder.RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return publishTraces, pagination, err
	}
	return publishTraces, pagination, nil
}

func (pt PublishTrace) GetDetail() (string, error) {
	var detail string
	err := sq.
		Select("detail").
		From(publishTraceTable).
		Where(sq.Eq{"id": pt.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&detail)
	if err != nil {
		return detail, err
	}
	return detail, nil
}

// AddRow return LastInsertId
func (pt PublishTrace) AddRow() (int64, error) {
	result, err := sq.
		Insert(publishTraceTable).
		Columns("token", "project_id", "project_name", "detail", "state", "publisher_id", "publisher_name", "type", "ext", "insert_time").
		Values(pt.Token, pt.ProjectID, pt.ProjectName, pt.Detail, pt.State, pt.PublisherID, pt.PublisherName, pt.Type, pt.Ext, pt.InsertTime).
		RunWith(DB).
		Exec()

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

// EditUpdateTimeByToken -
func (pt PublishTrace) EditUpdateTimeByToken() error {
	_, err := sq.
		Update(publishTraceTable).
		SetMap(sq.Eq{
			"update_time": pt.UpdateTime,
		}).
		Where(sq.Eq{"token": pt.Token}).
		RunWith(DB).
		Exec()
	return err
}
