package model

import sq "github.com/Masterminds/squirrel"

const templateTable = "`template`"

// Template mysql table template
type Template struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Remark       string `json:"remark"`
	PackageIDStr string `json:"packageIdStr"`
	Script       string `json:"script"`
	CreateTime   int64  `json:"createTime"`
	UpdateTime   int64  `json:"updateTime"`
}

// Templates many template
type Templates []Template

// AddRow add one row to table template and add id to tpl.ID
func (tpl Template) AddRow() (int64, error) {
	result, err := sq.
		Insert(templateTable).
		Columns("name", "remark", "script", "package_id_str", "create_time", "update_time").
		Values(tpl.Name, tpl.Remark, tpl.Script, tpl.PackageIDStr, tpl.CreateTime, tpl.UpdateTime).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table template
func (tpl Template) EditRow() error {
	_, err := sq.
		Update(templateTable).
		SetMap(sq.Eq{
			"name":           tpl.Name,
			"remark":         tpl.Remark,
			"script":         tpl.Script,
			"package_id_str": tpl.PackageIDStr,
		}).
		Where(sq.Eq{"id": tpl.ID}).
		RunWith(DB).
		Exec()
	return err
}

// Remove Template
func (tpl Template) Remove() error {
	_, err := sq.
		Delete(templateTable).
		Where(sq.Eq{"server_id": tpl.ID}).
		RunWith(DB).
		Exec()
	return err
}

// GetList template row
func (tpl Template) GetList(pagination Pagination) (Templates, Pagination, error) {
	rows, err := sq.
		Select("id, name, remark, script, package_id_str, create_time, update_time").
		From(templateTable).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, pagination, err
	}
	var templates Templates
	for rows.Next() {
		var template Template

		if err := rows.Scan(&template.ID, &template.Name, &template.Remark, &template.Script, &template.PackageIDStr, &template.CreateTime, &template.UpdateTime); err != nil {
			return nil, pagination, err
		}
		templates = append(templates, template)
	}
	err = sq.
		Select("COUNT(*) AS count").
		From(templateTable).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}
	return templates, pagination, nil
}

// GetAll template row
func (tpl Template) GetAll() (Templates, error) {
	rows, err := sq.
		Select("id, name, remark, script, package_id_str, create_time, update_time").
		From(templateTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var templates Templates
	for rows.Next() {
		var template Template

		if err := rows.Scan(&template.ID, &template.Name, &template.Remark, &template.Script, &template.PackageIDStr, &template.CreateTime, &template.UpdateTime); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

// GetData add template information to tpl *Template
func (tpl Template) GetData() (Template, error) {
	var template Template
	err := sq.
		Select("id, name, remark, script, package_id_str, create_time, update_time").
		From(templateTable).
		Where(sq.Eq{"id": tpl.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&template.ID, &template.Name, &template.Remark, &template.Script, &template.PackageIDStr, &template.CreateTime, &template.UpdateTime)
	if err != nil {
		return template, err
	}
	return template, nil
}
