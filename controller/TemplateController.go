package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Template struct
type Template Controller

// GetList template list
func (template Template) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Template model.Templates `json:"templateList"`
	}

	templateList, err := model.Template{}.GetList()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Template: templateList}}
	response.JSON(w)
}

// GetOption template list
func (template Template) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Template model.Templates `json:"templateList"`
	}

	templateList, err := model.Template{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Template: templateList}}
	response.JSON(w)
}

// Upload file
func (template Template) Upload(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		FileName string `json:"fileName"`
	}
	file, handler, err := gp.Request.FormFile("file")
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	defer file.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	filePath := core.TemplatePath + handler.Filename
	templateID := gp.URLQuery.Get("templateId")
	if templateID != "" {
		filePath = core.TemplatePath + templateID + "/" + handler.Filename
	}
	err = ioutil.WriteFile(filePath, fileBytes, 0755)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{FileName: handler.Filename}}
	response.JSON(w)
}

// Add one template
func (template Template) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Name    string `json:"name"`
		Remark  string `json:"remark"`
		Package string `json:"package"`
		Script  string `json:"script"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	id, err := model.Template{
		Name:       reqData.Name,
		Remark:     reqData.Remark,
		Package:    reqData.Package,
		Script:     reqData.Script,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	if err := os.Mkdir(core.TemplatePath+strconv.Itoa(int(id)), os.ModePerm); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	for _, fileName := range strings.Split(reqData.Package, ",") {
		err := os.Rename(core.TemplatePath+fileName, core.TemplatePath+strconv.Itoa(int(id))+"/"+fileName)
		if err != nil {
			response := core.Response{Code: 1, Message: err.Error()}
			response.JSON(w)
			return
		}
	}

	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one template
func (template Template) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID      uint32 `json:"id"`
		Name    string `json:"name"`
		Remark  string `json:"remark"`
		Package string `json:"package"`
		Script  string `json:"script"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Template{
		ID:         reqData.ID,
		Name:       reqData.Name,
		Remark:     reqData.Remark,
		Package:    reqData.Package,
		Script:     reqData.Script,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}

// RemovePackage one Template file
func (template Template) RemovePackage(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		TemplateID uint32 `json:"templateId"`
		Filename   string `json:"filename"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	filePath := core.TemplatePath + reqData.Filename

	if reqData.TemplateID != 0 {
		filePath = core.TemplatePath + strconv.Itoa(int(reqData.TemplateID)) + "/" + reqData.Filename
	}

	if err := os.Remove(filePath); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}

// Remove one Template
func (template Template) Remove(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Template{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.Remove()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	os.Remove(core.TemplatePath + strconv.Itoa(int(reqData.ID)))
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}
