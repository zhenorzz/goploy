package controller

import (
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"io"
	"io/ioutil"
	"path"
	"strconv"
)

// Server struct
type Server Controller

// GetList -
func (Server) GetList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

// GetTotal -
func (Server) GetTotal(gp *core.Goploy) *core.Response {
	total, err := model.Server{NamespaceID: gp.Namespace.ID}.GetTotal()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

// GetOption -
func (Server) GetOption(gp *core.Goploy) *core.Response {
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

// GetPublicKey -
func (Server) GetPublicKey(gp *core.Goploy) *core.Response {
	path := gp.URLQuery.Get("path")

	contentByte, err := ioutil.ReadFile(path + ".pub")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Key string `json:"key"`
		}{Key: string(contentByte)},
	}
}

// Check server
func (Server) Check(gp *core.Goploy) *core.Response {
	type ReqData struct {
		IP       string `json:"ip" validate:"required,ip|hostname"`
		Port     int    `json:"port" validate:"min=0,max=65535"`
		Owner    string `json:"owner" validate:"required,max=255"`
		Path     string `json:"path" validate:"required,max=255"`
		Password string `json:"password" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if _, err := utils.DialSSH(reqData.Owner, reqData.Password, reqData.Path, reqData.IP, reqData.Port); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Message: "Connected"}
}

// Add server
func (Server) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		NamespaceID int64  `json:"namespaceId" validate:"gte=0"`
		IP          string `json:"ip" validate:"ip|hostname"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required,max=255"`
		Path        string `json:"path" validate:"required,max=255"`
		Password    string `json:"password"`
		Description string `json:"description" validate:"max=255"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.Server{
		NamespaceID: reqData.NamespaceID,
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		Path:        reqData.Path,
		Password:    reqData.Password,
		Description: reqData.Description,
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

// Edit server
func (Server) Edit(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		NamespaceID int64  `json:"namespaceId" validate:"gte=0"`
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"required,ip|hostname"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required,max=255"`
		Path        string `json:"path" validate:"required,max=255"`
		Password    string `json:"password" validate:"max=255"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Server{
		ID:          reqData.ID,
		NamespaceID: reqData.NamespaceID,
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		Path:        reqData.Path,
		Password:    reqData.Password,
		Description: reqData.Description,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Toggle server
func (Server) Toggle(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID    int64 `json:"id" validate:"gt=0"`
		State int8  `json:"state" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Server{ID: reqData.ID, State: reqData.State}).ToggleRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// RemoteFile sftp download file
func (Server) RemoteFile(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: "invalid server id"}
	}

	filename := gp.URLQuery.Get("file")
	server, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	client, err := utils.DialSSH(server.Owner, server.Password, server.Path, server.IP, server.Port)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	//defer client.Close()
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	srcFile, _ := sftpClient.Open(filename) //远程
	FileStat, _ := srcFile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)
	gp.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+path.Base(filename))
	gp.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	gp.ResponseWriter.Header().Set("Content-Length", FileSize)
	_, _ = io.Copy(gp.ResponseWriter, srcFile)
	_ = srcFile.Close()
	_ = sftpClient.Close()
	return nil
}
