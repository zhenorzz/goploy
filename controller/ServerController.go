package controller

import (
	"bytes"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"github.com/zhenorzz/goploy/ws"

	"github.com/google/uuid"
)

// Server struct
type Server Controller

// GetList -
func (server Server) GetList(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Servers model.Servers `json:"list"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Servers: serverList}}
}

// GetTotal -
func (server Server) GetTotal(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Total int64 `json:"total"`
	}
	total, err := model.Server{NamespaceID: gp.Namespace.ID}.GetTotal()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Total: total}}
}

// GetInstallPreview server install preview list
func (server Server) GetInstallPreview(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		InstallTraceList model.InstallTraces `json:"installTraceList"`
	}
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	installTraceList, err := model.InstallTrace{ServerID: serverID}.GetListGroupByToken()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{InstallTraceList: installTraceList}}
}

// GetInstallList server install list by token
func (server Server) GetInstallList(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		InstallTraceList model.InstallTraces `json:"installTraceList"`
	}
	token := gp.URLQuery.Get("token")
	if err := core.Validate.Var(token, "uuid4"); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &core.Response{Code: core.Error, Message: "Token" + err.Translate(core.Trans)}
		}
	}
	installTraceList, err := model.InstallTrace{Token: token}.GetListByToken()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{InstallTraceList: installTraceList}}
}

// GetOption -
func (server Server) GetOption(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Servers model.Servers `json:"list"`
	}

	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Servers: serverList}}
}

// Check server
func (server Server) Check(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		IP    string `json:"ip" validate:"ip4_addr"`
		Port  int    `json:"port" validate:"min=0,max=65535"`
		Owner string `json:"owner" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if _, err := utils.ConnectSSH(reqData.Owner, "", reqData.IP, reqData.Port); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Message: "Connected"}
}

// Add server
func (server Server) Add(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"ip4_addr"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	type RespData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.Server{
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		NamespaceID: gp.Namespace.ID,
		Description: reqData.Description,
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{Data: RespData{ID: id}}
}

// Edit server
func (server Server) Edit(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"ip4_addr"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Server{
		ID:          reqData.ID,
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		Description: reqData.Description,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// RemoveRow server
func (server Server) Remove(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Server{ID: reqData.ID}).RemoveRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Install Server Environment
func (server Server) Install(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ServerID   int64 `json:"serverId" validate:"gt=0"`
		TemplateID int64 `json:"templateId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	serverInfo, err := model.Server{
		ID: reqData.ServerID,
	}.GetData()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	templateInfo, err := model.Template{
		ID: reqData.TemplateID,
	}.GetData()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	serverInfo.LastInstallToken = uuid.New().String()
	serverInfo.Install()

	go remoteInstall(gp.UserInfo, serverInfo, templateInfo)

	return &core.Response{Message: "Installing"}
}

// remoteInstall -
func remoteInstall(userInfo model.User, server model.Server, template model.Template) {
	installTraceModel := model.InstallTrace{
		Token:        server.LastInstallToken,
		ServerID:     server.ID,
		ServerName:   server.Name,
		OperatorID:   userInfo.ID,
		OperatorName: userInfo.Name,
		Type:         model.Rsync,
	}
	if template.PackageIDStr != "" {
		packages, err := model.Package{}.GetAllInID(strings.Split(template.PackageIDStr, ","))
		if err != nil {
			core.Log(core.ERROR, server.LastInstallToken+":"+err.Error())
			return
		}
		srcPath := core.PackagePath
		remoteMachine := server.Owner + "@" + server.IP
		destPath := remoteMachine + ":/tmp/goploy"
		rsyncOption := []string{
			"-rtv",
			"-e",
			"ssh -p " + strconv.Itoa(server.Port) + " -o StrictHostKeyChecking=no",
			"--include",
			"'*/'",
		}

		for _, pkg := range packages {
			rsyncOption = append(rsyncOption, "--include", pkg.Name)
		}
		rsyncOption = append(rsyncOption, "--exclude", "'*'", srcPath, destPath)
		cmd := exec.Command("rsync", rsyncOption...)
		var outbuf, errbuf bytes.Buffer
		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf
		core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(server.ID), 10)+" rsync "+strings.Join(rsyncOption, " "))
		rsyncError := cmd.Run()
		ext, _ := json.Marshal(struct {
			Command string `json:"command"`
		}{"rsync " + strings.Join(rsyncOption, " ")})

		installTraceModel.Ext = string(ext)
		installTraceModel.Type = model.Rsync
		if rsyncError != nil {
			installTraceModel.Detail = errbuf.String()
			installTraceModel.State = model.Fail
			installTraceModel.AddRow()
			ws.GetHub().Data <- &ws.Data{
				Type:    ws.TypeServerTemplate,
				UserIDs: []int64{userInfo.ID},
				Message: ws.ServerTemplateMessage{
					ServerID:   installTraceModel.ServerID,
					ServerName: installTraceModel.ServerName,
					Detail:     installTraceModel.Detail,
					Ext:        string(ext),
					Type:       ws.ServerTemplateRsync,
				},
			}
			return
		}

		installTraceModel.Detail = outbuf.String()
		installTraceModel.State = model.Success
		installTraceModel.AddRow()

		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeServerTemplate,
			UserIDs: []int64{userInfo.ID},
			Message: ws.ServerTemplateMessage{
				ServerID:   installTraceModel.ServerID,
				ServerName: installTraceModel.ServerName,
				Detail:     installTraceModel.Detail,
				Ext:        string(ext),
				Type:       ws.ServerTemplateRsync,
			},
		}
	}

	var scriptError error
	session, connectError := utils.ConnectSSH(server.Owner, "", server.IP, server.Port)
	ext, _ := json.Marshal(struct {
		SSH string `json:"ssh"`
	}{"ssh -p" + strconv.Itoa(server.Port) + " " + server.Owner + "@" + server.IP})
	installTraceModel.Ext = string(ext)
	installTraceModel.Type = model.SSH
	if connectError != nil {
		installTraceModel.State = model.Fail
		installTraceModel.Detail = connectError.Error()
		installTraceModel.AddRow()
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeServerTemplate,
			UserIDs: []int64{userInfo.ID},
			Message: ws.ServerTemplateMessage{
				ServerID:   installTraceModel.ServerID,
				ServerName: installTraceModel.ServerName,
				Detail:     installTraceModel.Detail,
				Ext:        string(ext),
				Type:       ws.ServerTemplateSSH,
			},
		}
		return
	}

	installTraceModel.State = model.Success
	installTraceModel.Detail = "connected"
	installTraceModel.AddRow()
	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeServerTemplate,
		UserIDs: []int64{userInfo.ID},
		Message: ws.ServerTemplateMessage{
			ServerID:   installTraceModel.ServerID,
			ServerName: installTraceModel.ServerName,
			Detail:     installTraceModel.Detail,
			Ext:        string(ext),
			Type:       ws.ServerTemplateSSH,
		},
	}
	defer session.Close()
	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	templateInstallScript := "echo '" + template.Script + "' > /tmp/goploy/template-install.sh;bash /tmp/goploy/template-install.sh"
	ext, _ = json.Marshal(struct {
		Script string `json:"script"`
	}{template.Script})
	installTraceModel.Ext = string(ext)
	installTraceModel.Type = model.Script
	if scriptError = session.Run(templateInstallScript); scriptError != nil {
		installTraceModel.State = model.Fail
		installTraceModel.Detail = scriptError.Error()
		installTraceModel.AddRow()
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeServerTemplate,
			UserIDs: []int64{userInfo.ID},
			Message: ws.ServerTemplateMessage{
				ServerID:   installTraceModel.ServerID,
				ServerName: installTraceModel.ServerName,
				Detail:     installTraceModel.Detail,
				Ext:        string(ext),
				Type:       ws.ServerTemplateScript,
			},
		}
		return
	}
	installTraceModel.State = model.Success
	installTraceModel.Detail = sshOutbuf.String()
	installTraceModel.AddRow()
	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeServerTemplate,
		UserIDs: []int64{userInfo.ID},
		Message: ws.ServerTemplateMessage{
			ServerID:   installTraceModel.ServerID,
			ServerName: installTraceModel.ServerName,
			Detail:     installTraceModel.Detail,
			Ext:        string(ext),
			Type:       ws.ServerTemplateScript,
		},
	}
	return
}
