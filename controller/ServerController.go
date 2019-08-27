package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/ws"
	"golang.org/x/crypto/ssh"
)

// Server struct
type Server Controller

// GetList server list
func (server Server) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Server model.Servers `json:"serverList"`
	}
	userData, err := core.GetUserData(gp.TokenInfo.ID)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	serverList, err := model.Server{}.GetListByManagerGroupStr(userData.ManageGroupStr)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Server: serverList}}
	response.JSON(w)
}

// getInstallPreview server install token list
func (server Server) GetInstallPreview(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		InstallTraceList model.InstallTraces `json:"installTraceList"`
	}
	serverID, err := strconv.Atoi(gp.URLQuery.Get("serverId"))
	if err != nil {
		response := core.Response{Code: 1, Message: "serverId参数错误"}
		response.JSON(w)
		return
	}
	installTraceList, err := model.InstallTrace{ServerID: uint32(serverID)}.GetListGroupByToken()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{InstallTraceList: installTraceList}}
	response.JSON(w)
}

// GetInstallList server install list by token
func (server Server) GetInstallList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		InstallTraceList model.InstallTraces `json:"installTraceList"`
	}
	token := gp.URLQuery.Get("token")
	installTraceList, err := model.InstallTrace{Token: token}.GetListByToken()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{InstallTraceList: installTraceList}}
	response.JSON(w)
}

// GetOption server list
func (server Server) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Server model.Servers `json:"serverList"`
	}

	serverList, err := model.Server{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Server: serverList}}
	response.JSON(w)
}

// Add one server
func (server Server) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Name    string `json:"name"`
		IP      string `json:"ip"`
		Port    uint32 `json:"port"`
		Owner   string `json:"owner"`
		GroupID uint32 `json:"groupId"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	_, err = model.Server{
		Name:       reqData.Name,
		IP:         reqData.IP,
		Port:       reqData.Port,
		Owner:      reqData.Owner,
		GroupID:    reqData.GroupID,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one server
func (server Server) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID      uint32 `json:"id"`
		Name    string `json:"name"`
		IP      string `json:"ip"`
		Port    uint32 `json:"port"`
		Owner   string `json:"owner"`
		GroupID uint32 `json:"groupId"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Server{
		ID:         reqData.ID,
		Name:       reqData.Name,
		IP:         reqData.IP,
		Port:       reqData.Port,
		Owner:      reqData.Owner,
		GroupID:    reqData.GroupID,
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

// Remove one Server
func (server Server) Remove(w http.ResponseWriter, gp *core.Goploy) {
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
	err = model.Server{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.Remove()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}

// Install Server Enviroment
func (server Server) Install(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ServerID   uint32 `json:"serverId"`
		TemplateID uint32 `json:"templateId"`
	}
	var reqData ReqData
	if err := json.Unmarshal(gp.Body, &reqData); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	serverInfo, err := model.Server{
		ID: reqData.ServerID,
	}.GetData()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	templateInfo, err := model.Template{
		ID: reqData.TemplateID,
	}.GetData()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	serverInfo.LastInstallToken = uuid.New().String()
	serverInfo.UpdateTime = time.Now().Unix()
	serverInfo.Install()

	go remoteInstall(gp.TokenInfo, serverInfo, templateInfo)

	response := core.Response{Message: "正在安装"}
	response.JSON(w)
}

func remoteInstall(tokenInfo core.TokenInfo, server model.Server, template model.Template) {
	installTraceModel := model.InstallTrace{
		Token:        server.LastInstallToken,
		ServerID:     server.ID,
		ServerName:   server.Name,
		OperatorID:   tokenInfo.ID,
		OperatorName: tokenInfo.Name,
		Type:         model.Rsync,
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
	}
	srcPath := core.TemplatePath + strconv.Itoa(int(template.ID)) + "/"
	remoteMachine := server.Owner + "@" + server.IP
	destPath := remoteMachine + ":/tmp"
	rsyncOption := []string{
		"-rtv",
		"-e",
		"ssh -p " + strconv.Itoa(int(server.Port)) + " -o StrictHostKeyChecking=no",
		srcPath,
		destPath,
	}
	cmd := exec.Command("rsync", rsyncOption...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(server.ID), 10)+" rsync "+strings.Join(rsyncOption, " "))
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.InstallMessage{
			ServerID: server.ID,
			UserID:   tokenInfo.ID,
			DataType: ws.RsyncType,
			State:    ws.Success,
			Message:  "rsync " + strings.Join(rsyncOption, " "),
		},
	}
	var rsyncError error
	// 失败重试三次
	for attempt := 0; attempt < 3; attempt++ {
		rsyncError = cmd.Run()
		if rsyncError != nil {
			core.Log(core.ERROR, errbuf.String())
		} else {
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.InstallMessage{
					ServerID: server.ID,
					UserID:   tokenInfo.ID,
					DataType: ws.RsyncType,
					State:    ws.Success,
					Message:  outbuf.String(),
				},
			}
			ext, _ := json.Marshal(struct {
				Command string `json:"command"`
				Package string `json:"package"`
			}{"rsync " + strings.Join(rsyncOption, " "), template.Package})
			installTraceModel.Ext = string(ext)
			installTraceModel.Detail = outbuf.String()
			installTraceModel.State = model.Success
			installTraceModel.AddRow()
			break
		}
	}

	if rsyncError != nil {
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.InstallMessage{
				ServerID: server.ID,
				UserID:   tokenInfo.ID,
				DataType: ws.RsyncType,
				State:    ws.Fail,
				Message:  "rsync重试失败",
			},
		}
		ext, _ := json.Marshal(struct {
			Command string `json:"command"`
			Package string `json:"package"`
		}{"rsync " + strings.Join(rsyncOption, " "), template.Package})
		installTraceModel.Ext = string(ext)
		installTraceModel.Detail = errbuf.String()
		installTraceModel.State = model.Fail
		installTraceModel.AddRow()
		return
	}

	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.InstallMessage{
			ServerID: server.ID,
			UserID:   tokenInfo.ID,
			DataType: ws.ScriptType,
			State:    ws.Success,
			Message:  "开始连接ssh",
		},
	}
	var session *ssh.Session
	var connectError error
	var scriptError error
	for attempt := 0; attempt < 3; attempt++ {
		session, connectError = connect(server.Owner, "", server.IP, int(server.Port))
		if connectError != nil {
			core.Log(core.ERROR, connectError.Error())
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.InstallMessage{
					ServerID: server.ID,
					UserID:   tokenInfo.ID,
					DataType: ws.ScriptType,
					State:    ws.Fail,
					Message:  connectError.Error(),
				},
			}
		} else {
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.InstallMessage{
					ServerID: server.ID,
					UserID:   tokenInfo.ID,
					DataType: ws.ScriptType,
					State:    ws.Success,
					Message:  "开始连接成功",
				},
			}
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.InstallMessage{
					ServerID: server.ID,
					UserID:   tokenInfo.ID,
					DataType: ws.ScriptType,
					State:    ws.Success,
					Message:  "运行:" + template.Script,
				},
			}
			defer session.Close()
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			sshOutbuf.Reset()
			templateInstallScript := "echo '" + template.Script + "' > /tmp/template-install.sh;bash /tmp/template-install.sh"
			if scriptError = session.Run(templateInstallScript); scriptError != nil {
				core.Log(core.ERROR, scriptError.Error())
				ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
					ToUserID: tokenInfo.ID,
					Message: ws.InstallMessage{
						ServerID: server.ID,
						UserID:   tokenInfo.ID,
						DataType: ws.ScriptType,
						State:    ws.Fail,
						Message:  scriptError.Error(),
					},
				}
			} else {
				ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
					ToUserID: tokenInfo.ID,
					Message: ws.InstallMessage{
						ServerID: server.ID,
						UserID:   tokenInfo.ID,
						DataType: ws.ScriptType,
						State:    ws.Success,
						Message:  sshOutbuf.String(),
					},
				}
				ext, _ := json.Marshal(struct {
					Package string `json:"package"`
					Script  string `json:"script"`
				}{template.Package, template.Script})
				installTraceModel.Ext = string(ext)
				installTraceModel.Type = model.Script
				installTraceModel.State = model.Success
				installTraceModel.Detail = sshOutbuf.String()
				installTraceModel.AddRow()
				break
			}
		}

	}

	if connectError != nil {
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.InstallMessage{
				ServerID: server.ID,
				UserID:   tokenInfo.ID,
				DataType: ws.ScriptType,
				State:    ws.Fail,
				Message:  "ssh重连失败",
			},
		}
		return
	} else if scriptError != nil {
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.InstallMessage{
				ServerID: server.ID,
				UserID:   tokenInfo.ID,
				DataType: ws.ScriptType,
				State:    ws.Fail,
				Message:  "脚本重试失败",
			},
		}
		ext, _ := json.Marshal(struct {
			Package string `json:"package"`
			Script  string `json:"script"`
		}{template.Package, template.Script})
		installTraceModel.Ext = string(ext)
		installTraceModel.Type = model.Script
		installTraceModel.State = model.Fail
		installTraceModel.Detail = scriptError.Error()
		installTraceModel.AddRow()
		return
	}
	return
}
