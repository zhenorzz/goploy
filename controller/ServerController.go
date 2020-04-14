package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"goploy/core"
	"goploy/model"
	"goploy/utils"
	"goploy/ws"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

// Server struct
type Server Controller

// GetList server list
func (server Server) GetList(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Server     model.Servers    `json:"serverList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	var serverList model.Servers
	if gp.UserInfo.Role == core.RoleAdmin || gp.UserInfo.Role == core.RoleManager {
		serverList, pagination, err = model.Server{}.GetList(pagination)
	} else {
		serverList, pagination, err = model.Server{}.GetListInGroupIDs(strings.Split(gp.UserInfo.ManageGroupStr, ","), pagination)
	}

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Server: serverList, Pagination: pagination}}
}

// GetInstallPreview server install token list
func (server Server) GetInstallPreview(w http.ResponseWriter, gp *core.Goploy) *core.Response {
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

// GetOption server list
func (server Server) GetOption(w http.ResponseWriter, _ *core.Goploy) *core.Response {
	type RespData struct {
		Server model.Servers `json:"serverList"`
	}

	serverList, err := model.Server{}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Server: serverList}}
}

// Check one server
func (server Server) Check(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		IP    string `json:"ip" validate:"ip4_addr"`
		Port  int    `json:"port" validate:"min=0,max=65535"`
		Owner string `json:"owner" validate:"required"`
	}
	type RespData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	_, err := utils.ConnectSSH(reqData.Owner, "", reqData.IP, reqData.Port)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Message: "Connected"}
}

// Add one server
func (server Server) Add(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"ip4_addr"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required"`
		GroupID     int64  `json:"groupId" validate:"min=0"`
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
		GroupID:     reqData.GroupID,
		Description: reqData.Description,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{Data: RespData{ID: id}}
}

// Edit one server
func (server Server) Edit(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"ip4_addr"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required"`
		GroupID     int64  `json:"groupId" validate:"min=0"`
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
		GroupID:     reqData.GroupID,
		Description: reqData.Description,
		UpdateTime:  time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := &core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Remove one Server
func (server Server) Remove(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Server{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.Remove()

	if err != nil {
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
	serverInfo.UpdateTime = time.Now().Unix()
	serverInfo.Install()

	go remoteInstall(gp.UserInfo, serverInfo, templateInfo)

	return &core.Response{Message: "Installing"}
}

func remoteInstall(userInfo model.User, server model.Server, template model.Template) {
	installTraceModel := model.InstallTrace{
		Token:        server.LastInstallToken,
		ServerID:     server.ID,
		ServerName:   server.Name,
		OperatorID:   userInfo.ID,
		OperatorName: userInfo.Name,
		Type:         model.Rsync,
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
	}
	if template.PackageIDStr != "" {
		packages, err := model.Package{}.GetAllInIDStr(template.PackageIDStr)
		if err != nil {
			core.Log(core.ERROR, server.LastInstallToken+":"+err.Error())
			return
		}
		srcPath := core.PackagePath
		remoteMachine := server.Owner + "@" + server.IP
		destPath := remoteMachine + ":/tmp"
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
		var rsyncError error
		// 失败重试三次
		for attempt := 0; attempt < 3; attempt++ {
			rsyncError = cmd.Run()
			if rsyncError != nil {
				core.Log(core.ERROR, errbuf.String())
			} else {
				ext, _ := json.Marshal(struct {
					Command string `json:"command"`
				}{"rsync " + strings.Join(rsyncOption, " ")})
				installTraceModel.Ext = string(ext)
				installTraceModel.Detail = outbuf.String()
				installTraceModel.State = model.Success
				installTraceModel.AddRow()

				ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
					ToUserID: userInfo.ID,
					Message:  installTraceModel,
				}
				break
			}
		}

		if rsyncError != nil {
			ext, _ := json.Marshal(struct {
				Command string `json:"command"`
			}{"rsync " + strings.Join(rsyncOption, " ")})
			installTraceModel.Ext = string(ext)
			installTraceModel.Detail = errbuf.String()
			installTraceModel.State = model.Fail
			installTraceModel.AddRow()
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: userInfo.ID,
				Message:  installTraceModel,
			}
			return
		}
	}

	var session *ssh.Session
	var connectError error
	var scriptError error
	for attempt := 0; attempt < 3; attempt++ {
		session, connectError = utils.ConnectSSH(server.Owner, "", server.IP, server.Port)
		if connectError != nil {
			core.Log(core.ERROR, connectError.Error())
		} else {
			ext, _ := json.Marshal(struct {
				SSH string `json:"ssh"`
			}{"ssh -p" + strconv.Itoa(server.Port) + " " + server.Owner + "@" + server.IP})
			installTraceModel.Ext = string(ext)
			installTraceModel.Type = model.SSH
			installTraceModel.State = model.Success
			installTraceModel.Detail = "connected"
			installTraceModel.AddRow()
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: userInfo.ID,
				Message:  installTraceModel,
			}
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			sshOutbuf.Reset()
			templateInstallScript := "echo '" + template.Script + "' > /tmp/template-install.sh;bash /tmp/template-install.sh"
			if scriptError = session.Run(templateInstallScript); scriptError != nil {
				core.Log(core.ERROR, scriptError.Error())
			} else {
				ext, _ := json.Marshal(struct {
					Script string `json:"script"`
				}{template.Script})
				installTraceModel.Ext = string(ext)
				installTraceModel.Type = model.Script
				installTraceModel.State = model.Success
				installTraceModel.Detail = sshOutbuf.String()
				installTraceModel.AddRow()
				ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
					ToUserID: userInfo.ID,
					Message:  installTraceModel,
				}
				break
			}
		}

	}
	if session != nil {
		defer session.Close()
	}
	if connectError != nil {
		ext, _ := json.Marshal(struct {
			SSH string `json:"ssh"`
		}{"ssh -p" + strconv.Itoa(server.Port) + " " + server.Owner + "@" + server.IP})
		installTraceModel.Ext = string(ext)
		installTraceModel.Type = model.SSH
		installTraceModel.State = model.Fail
		installTraceModel.Detail = connectError.Error()
		installTraceModel.AddRow()
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: userInfo.ID,
			Message:  installTraceModel,
		}
		return
	} else if scriptError != nil {
		ext, _ := json.Marshal(struct {
			Script string `json:"script"`
		}{template.Script})
		installTraceModel.Ext = string(ext)
		installTraceModel.Type = model.Script
		installTraceModel.State = model.Fail
		installTraceModel.Detail = scriptError.Error()
		installTraceModel.AddRow()
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: userInfo.ID,
			Message:  installTraceModel,
		}
		return
	}
	return
}
