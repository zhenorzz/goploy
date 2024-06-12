// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const notificationTemplateTable = "`notification_template`"

type NotificationTemplate struct {
	ID         int64  `json:"id"`
	Type       uint8  `json:"type"`
	UseBy      string `json:"useBy"`
	Title      string `json:"title"`
	Template   string `json:"template"`
	InsertTime string `json:"insertTime,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

type NotificationTemplates []NotificationTemplate

func (nt NotificationTemplate) EditRow() error {
	_, err := sq.
		Update(notificationTemplateTable).
		Set("template", nt.Template).
		Set("title", nt.Title).
		Where(sq.Eq{"id": nt.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (nt NotificationTemplate) GetList() (NotificationTemplates, error) {
	rows, err := sq.
		Select("id, use_by, type, title, template, insert_time, update_time").
		From(notificationTemplateTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}

	notificationTemplates := NotificationTemplates{}
	for rows.Next() {
		var notificationTemplate NotificationTemplate

		if err := rows.Scan(
			&notificationTemplate.ID,
			&notificationTemplate.UseBy,
			&notificationTemplate.Type,
			&notificationTemplate.Title,
			&notificationTemplate.Template,
			&notificationTemplate.InsertTime,
			&notificationTemplate.UpdateTime); err != nil {
			return nil, err
		}
		notificationTemplates = append(notificationTemplates, notificationTemplate)
	}

	return notificationTemplates, nil
}

func (nt NotificationTemplate) GetAll() (NotificationTemplates, error) {
	rows, err := sq.
		Select("id, use_by, type, template").
		From(notificationTemplateTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	notificationTemplates := NotificationTemplates{}
	for rows.Next() {
		var notificationTemplate NotificationTemplate

		if err := rows.Scan(&notificationTemplate.ID, &notificationTemplate.UseBy, &notificationTemplate.Type, &notificationTemplate.Template); err != nil {
			return notificationTemplates, err
		}
		notificationTemplates = append(notificationTemplates, notificationTemplate)
	}
	return notificationTemplates, nil
}

func (nt NotificationTemplate) GetData() (NotificationTemplate, error) {
	var notificationTemplate NotificationTemplate
	err := sq.
		Select("use_by, type, template").
		From(notificationTemplateTable).
		Where(sq.Eq{"id": nt.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&notificationTemplate.UseBy, &notificationTemplate.Type, &notificationTemplate.Template)
	return notificationTemplate, err
}

func (nt NotificationTemplate) GetTemplate() (NotificationTemplate, error) {
	var notificationTemplate NotificationTemplate
	err := sq.
		Select("title, template").
		From(notificationTemplateTable).
		Where(sq.Eq{"use_by": nt.UseBy}).
		Where(sq.Eq{"type": nt.Type}).
		RunWith(DB).
		QueryRow().
		Scan(&notificationTemplate.Title, &notificationTemplate.Template)
	return notificationTemplate, err
}
