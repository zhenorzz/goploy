package controller

import (
	"database/sql"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Package struct
type Package Controller

// GetList list
func (Package) GetList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	packageList, err := model.Package{}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Packages model.Packages `json:"list"`
		}{Packages: packageList},
	}
}

// GetTotal total
func (Package) GetTotal(gp *core.Goploy) *core.Response {
	total, err := model.Package{}.GetTotal()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

// GetOption list
func (Package) GetOption(gp *core.Goploy) *core.Response {
	packageList, err := model.Package{}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Packages model.Packages `json:"list"`
		}{Packages: packageList},
	}
}

// Upload file
func (Package) Upload(gp *core.Goploy) *core.Response {
	file, handler, err := gp.Request.FormFile("file")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer file.Close()
	_, err = model.Package{Name: handler.Filename}.GetDataByName()
	if err != sql.ErrNoRows {
		return &core.Response{Code: core.Deny, Message: "The same file already exists"}
	}
	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	filePath := path.Join(core.GetPackagePath(), handler.Filename)
	if err := ioutil.WriteFile(filePath, fileBytes, 0755); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	packageIDStr := gp.URLQuery.Get("packageId")

	if packageIDStr == "" {
		_, err = model.Package{
			Name: handler.Filename,
			Size: handler.Size,
		}.AddRow()
	} else {
		packageID, err := strconv.ParseInt(packageIDStr, 10, 64)
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}

		err = model.Package{
			ID:   packageID,
			Name: handler.Filename,
			Size: handler.Size,
		}.EditRow()
	}

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{
		Data: struct {
			Filename string `json:"filename"`
		}{Filename: handler.Filename},
	}
}
