package model

// Package mysql table package
type Package struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// Packages many package
type Packages []Package

// GetList package row
func (p Package) GetList() (Packages, error) {
	rows, err := DB.Query("SELECT id, name, size, create_time, update_time FROM package order by id desc")
	if err != nil {
		return nil, err
	}
	var packages Packages
	for rows.Next() {
		var pkg Package

		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.CreateTime, &pkg.UpdateTime); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// GetListInIDStr package row
func (p Package) GetListInIDStr(IDStr string) (Packages, error) {
	rows, err := DB.Query("SELECT id, name, size FROM package where id in (" + IDStr + ")")
	if err != nil {
		return nil, err
	}
	var packages Packages
	for rows.Next() {
		var pkg Package

		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Size); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// GetAll package row
func (p Package) GetAll() (Packages, error) {
	rows, err := DB.Query("SELECT id, name, size, create_time, update_time FROM package ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	var packages Packages
	for rows.Next() {
		var pkg Package

		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.CreateTime, &pkg.UpdateTime); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// GetData get package information
func (p Package) GetData() (Package, error) {
	var pkg Package
	err := DB.QueryRow("SELECT id, name, size, create_time, update_time FROM package WHERE id = ?", p.ID).Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.CreateTime, &pkg.UpdateTime)
	if err != nil {
		return pkg, err
	}
	return pkg, nil
}

// GetDataByName get package information
func (p Package) GetDataByName() (Package, error) {
	var pkg Package
	err := DB.QueryRow("SELECT id, name, size, create_time, update_time FROM package WHERE name = ?", p.Name).Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.CreateTime, &pkg.UpdateTime)
	if err != nil {
		return pkg, err
	}
	return pkg, nil
}

// AddRow add one row to table package and add id to p.ID
func (p Package) AddRow() (uint32, error) {
	result, err := DB.Exec(
		"INSERT INTO package (name, size, create_time, update_time) VALUES (?, ?, ?, ?)",
		p.Name,
		p.Size,
		p.CreateTime,
		p.UpdateTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table package
func (p Package) EditRow() error {
	_, err := DB.Exec(
		`UPDATE package SET 
		  name = ?,
		  size = ?
		WHERE
		 id = ?`,
		p.Name,
		p.Size,
		p.ID,
	)
	return err
}
