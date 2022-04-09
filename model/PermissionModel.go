// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
	"strconv"
	"strings"
)

const permissionTable = "`permission`"

type Permission struct {
	ID          int64  `json:"id"`
	PID         int64  `json:"pid"`
	Name        string `json:"name"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
	InsertTime  string `json:"insertTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

type Permissions []Permission

func (p Permission) GetList() (Permissions, error) {
	rows, err := sq.
		Select("id, pid, name, sort, description").
		From(permissionTable).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}

	permissions := Permissions{}
	for rows.Next() {
		var permission Permission

		if err := rows.Scan(&permission.ID, &permission.PID, &permission.Name, &permission.Sort, &permission.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (p Permission) GetIDs() (map[int64]struct{}, error) {
	var permissionIDs string
	err := sq.Select("GROUP_CONCAT(id)").
		From(permissionTable).
		RunWith(DB).
		QueryRow().
		Scan(&permissionIDs)

	if err != nil {
		return nil, err
	}

	ids := map[int64]struct{}{}
	for _, permissionID := range strings.Split(permissionIDs, ",") {
		v, err := strconv.ParseInt(permissionID, 10, 64)
		if err != nil {
			return nil, err
		}
		ids[v] = struct{}{}
	}
	return ids, nil
}
