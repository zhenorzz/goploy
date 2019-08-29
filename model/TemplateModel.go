package model

// Template mysql table template
type Template struct {
	ID           uint32 `json:"id"`
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
func (tpl Template) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO template (name, remark, script, package_id_str, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?)",
		tpl.Name,
		tpl.Remark,
		tpl.Script,
		tpl.PackageIDStr,
		tpl.CreateTime,
		tpl.UpdateTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table template
func (tpl Template) EditRow() error {
	_, err := DB.Exec(
		`UPDATE template SET 
		  name = ?,
		  remark = ?,
		  script = ?,
		  package_id_str = ?
		WHERE
		 id = ?`,
		tpl.Name,
		tpl.Remark,
		tpl.Script,
		tpl.PackageIDStr,
		tpl.ID,
	)
	return err
}

// Remove Template
func (tpl Template) Remove() error {
	_, err := DB.Exec(`DELETE FROM template WHERE id = ?`, tpl.ID)
	return err
}

// GetList template row
func (tpl Template) GetList() (Templates, error) {
	rows, err := DB.Query("SELECT id, name, remark, script, package_id_str, create_time, update_time FROM template")
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

// GetAll template row
func (tpl Template) GetAll() (Templates, error) {
	rows, err := DB.Query("SELECT id, name, remark, script, package_id_str, create_time, update_time FROM template ORDER BY id DESC")
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
	err := DB.QueryRow("SELECT id, name, remark, script, package_id_str, create_time, update_time FROM template WHERE id = ?", tpl.ID).Scan(&template.ID, &template.Name, &template.Remark, &template.Script, &template.PackageIDStr, &template.CreateTime, &template.UpdateTime)
	if err != nil {
		return template, err
	}
	return template, nil
}
