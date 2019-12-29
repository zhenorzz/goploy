package controller

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"goploy/core"
	"goploy/model"
)

// Package struct
type Package Controller

// GetList  list
func (Package) GetList(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Package    model.Packages   `json:"packageList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	packageList, pagination, err := model.Package{}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Package: packageList, Pagination: pagination}}
}

// GetOption  list
func (Package) GetOption(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Package model.Packages `json:"packageList"`
	}

	packageList, err := model.Package{}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Package: packageList}}
}

// Upload file
func (Package) Upload(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Filename string `json:"filename"`
	}
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
	filePath := core.PackagePath + handler.Filename
	if err := ioutil.WriteFile(filePath, fileBytes, 0755); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	packageIDStr := gp.URLQuery.Get("packageId")

	if packageIDStr == "" {
		_, err = model.Package{
			Name:       handler.Filename,
			Size:       handler.Size,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}.AddRow()
	} else {
		packageID, err := strconv.ParseInt(packageIDStr, 10, 64)
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}

		err = model.Package{
			ID:         packageID,
			Name:       handler.Filename,
			Size:       handler.Size,
			UpdateTime: time.Now().Unix(),
		}.EditRow()
	}

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{Data: RespData{Filename: handler.Filename}}
}
