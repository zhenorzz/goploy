// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const templateTable = "`template`"

type Template struct {
	ID          int64  `json:"id"`
	Type        uint8  `json:"type"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	Description string `json:"description"`
	InsertTime  string `json:"insertTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

type Templates []Template

func (t Template) AddRow() (int64, error) {
	result, err := sq.
		Insert(templateTable).
		Columns("type", "name", "content", "description").
		Values(t.Type, t.Name, t.Content, t.Description).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (t Template) DeleteRow() error {
	_, err := sq.
		Delete(templateTable).
		Where(sq.Eq{"id": t.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (t Template) GetAll() (Templates, error) {
	rows, err := sq.
		Select("id, name, content, description").
		From(templateTable).
		Where(sq.Eq{"type": t.Type}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	templates := Templates{}
	for rows.Next() {
		var template Template

		if err := rows.Scan(&template.ID, &template.Name, &template.Content, &template.Description); err != nil {
			return templates, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func (t Template) GetData() (Role, error) {
	var role Role
	err := sq.
		Select("name, description").
		From(templateTable).
		Where(sq.Eq{"id": t.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&role.Name, &role.Description)
	return role, err
}
