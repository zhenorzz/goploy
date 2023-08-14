// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package server

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/internal/validator"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Server struct
type Server api.API

func (s Server) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/server/getList", http.MethodGet, s.GetList).Permissions(config.ShowServerPage),
		server.NewRoute("/server/getOption", http.MethodGet, s.GetOption),
		server.NewRoute("/server/getBindProjectList", http.MethodGet, s.GetBindProjectList).Permissions(config.ShowServerPage),
		server.NewRoute("/server/getPublicKey", http.MethodGet, s.GetPublicKey).Permissions(config.AddServer, config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/check", http.MethodPost, s.Check).Permissions(config.AddServer, config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/import", http.MethodPost, s.Import).Permissions(config.ImportCSV).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/add", http.MethodPost, s.Add).Permissions(config.AddServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/edit", http.MethodPut, s.Edit).Permissions(config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/toggle", http.MethodPut, s.Toggle).Permissions(config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/unbindProject", http.MethodDelete, s.UnbindProject).Permissions(config.UnbindServerProject).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/installAgent", http.MethodPost, s.InstallAgent).Permissions(config.InstallAgent).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/previewFile", http.MethodGet, s.PreviewFile).Permissions(config.SFTPPreviewFile).LogFunc(middleware.AddPreviewLog),
		server.NewRoute("/server/downloadFile", http.MethodGet, s.DownloadFile).Permissions(config.SFTPDownloadFile).LogFunc(middleware.AddDownloadLog),
		server.NewRoute("/server/uploadFile", http.MethodPost, s.UploadFile).Permissions(config.SFTPTransferFile).LogFunc(middleware.AddUploadLog),
		server.NewRoute("/server/editFile", http.MethodPut, s.EditFile).Permissions(config.SFTPEditFile).LogFunc(middleware.AddEditLog),
		server.NewRoute("/server/renameFile", http.MethodPut, s.RenameFile).Permissions(config.SFTPEditFile).LogFunc(middleware.AddRenameLog),
		server.NewRoute("/server/copyFile", http.MethodPut, s.CopyFile).Permissions(config.SFTPUploadFile).LogFunc(middleware.AddCopyLog),
		server.NewRoute("/server/deleteFile", http.MethodDelete, s.DeleteFile).Permissions(config.SFTPDeleteFile).LogFunc(middleware.AddDeleteLog),
		server.NewRoute("/server/transferFile", http.MethodPost, s.TransferFile).Permissions(config.SFTPUploadFile),
		server.NewRoute("/server/report", http.MethodGet, s.Report).Permissions(config.ShowServerMonitorPage),
		server.NewRoute("/server/getAllMonitor", http.MethodGet, s.GetAllMonitor).Permissions(config.ShowServerMonitorPage),
		server.NewRoute("/server/addMonitor", http.MethodPost, s.AddMonitor).Permissions(config.AddServerWarningRule).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/editMonitor", http.MethodPut, s.EditMonitor).Permissions(config.EditServerWarningRule).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/deleteMonitor", http.MethodDelete, s.DeleteMonitor).Permissions(config.DeleteServerWarningRule).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/getProcessList", http.MethodGet, s.GetProcessList).Permissions(config.ShowServerProcessPage),
		server.NewRoute("/server/addProcess", http.MethodPost, s.AddProcess).Permissions(config.AddServerProcess).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/editProcess", http.MethodPut, s.EditProcess).Permissions(config.EditServerProcess).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/deleteProcess", http.MethodDelete, s.DeleteProcess).Permissions(config.DeleteServerProcess).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/execProcess", http.MethodPost, s.ExecProcess).Permissions(config.ShowServerProcessPage).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/execScript", http.MethodPost, s.ExecScript).Permissions(config.ShowServerScriptPage).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/getNginxConfigList", http.MethodGet, s.GetNginxConfigList).Permissions(config.ShowServerNginxPage),
		server.NewRoute("/server/getNginxConfigContent", http.MethodGet, s.GetNginxConfContent).Permissions(config.ShowServerNginxPage),
		server.NewRoute("/server/getNginxPath", http.MethodGet, s.GetNginxPath).Permissions(config.ShowServerNginxPage),
		server.NewRoute("/server/manageNginx", http.MethodPost, s.ManageNginx).Permissions(config.ManageServerNginx).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/copyNginxConfig", http.MethodPost, s.CopyNginxConfig).Permissions(config.AddNginxConfig).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/editNginxConfig", http.MethodPut, s.EditNginxConfig).Permissions(config.EditNginxConfig).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/renameNginxConfig", http.MethodPut, s.RenameNginxConfig).Permissions(config.EditNginxConfig).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/deleteNginxConfig", http.MethodDelete, s.DeleteNginxConfig).Permissions(config.DeleteNginxConfig).LogFunc(middleware.AddOPLog),
	}
}

func (Server) GetList(gp *server.Goploy) server.Response {
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

func (Server) GetOption(gp *server.Goploy) server.Response {
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

// GetBindProjectList lists all binding projects
// @Summary List all binding projects
// @Tags Server
// @Produce json
// @Param request query server.GetBindUserList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=project.GetBindProjectList.RespData}
// @Router /server/getBindProjectList [get]
func (Server) GetBindProjectList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectServers, err := model.ProjectServer{ServerID: reqData.ID}.GetBindProjectListByServerID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ProjectServers model.ProjectServers `json:"list"`
	}
	return response.JSON{
		Data: RespData{ProjectServers: projectServers},
	}
}

func (Server) GetPublicKey(gp *server.Goploy) server.Response {
	publicKeyPath := gp.URLQuery.Get("path")

	contentByte, err := os.ReadFile(publicKeyPath + ".pub")
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Key string `json:"key"`
		}{Key: string(contentByte)},
	}
}

func (Server) Check(gp *server.Goploy) server.Response {
	type ReqData struct {
		IP           string `json:"ip" validate:"required,ip|hostname"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Owner        string `json:"owner" validate:"required,max=255"`
		Path         string `json:"path" validate:"max=255"`
		Password     string `json:"password" validate:"max=255"`
		JumpIP       string `json:"jumpIP" validate:"omitempty,ip|hostname"`
		JumpPort     int    `json:"jumpPort" validate:"min=0,max=65535"`
		JumpOwner    string `json:"jumpOwner" validate:"max=255"`
		JumpPath     string `json:"jumpPath" validate:"max=255"`
		JumpPassword string `json:"jumpPassword"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	sshConfig := pkg.SSHConfig{
		User:         reqData.Owner,
		Password:     reqData.Password,
		Path:         reqData.Path,
		Host:         reqData.IP,
		Port:         reqData.Port,
		JumpUser:     reqData.JumpOwner,
		JumpPassword: reqData.JumpPassword,
		JumpPath:     reqData.JumpPath,
		JumpHost:     reqData.JumpIP,
		JumpPort:     reqData.JumpPort,
	}

	if Conn, err := sshConfig.Dial(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	} else {
		_ = Conn.Close()
	}
	return response.JSON{Message: "Connected"}
}

func (Server) Import(gp *server.Goploy) server.Response {
	file, _, err := gp.Request.FormFile("file")
	if err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}
	defer file.Close()
	r := csv.NewReader(file)
	i := 0
	headerIdx := map[string]int{
		"name":         -1,
		"os":           -1,
		"host":         -1,
		"port":         -1,
		"owner":        -1,
		"path":         -1,
		"password":     -1,
		"description":  -1,
		"jumpHost":     -1,
		"jumpPort":     -1,
		"jumpOwner":    -1,
		"jumpPath":     -1,
		"jumpPassword": -1,
	}
	errOccur := false
	var wg sync.WaitGroup
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		i++
		if i == 1 {
			for index, header := range record {
				if _, ok := headerIdx[header]; !ok {
					return response.JSON{Code: response.Error, Message: fmt.Sprintf("%s does not match the csv field", header)}
				} else {
					headerIdx[header] = index
				}
			}
			requiredFields := []string{"name", "host", "port", "owner"}
			missingFields := ""
			for _, field := range requiredFields {
				if headerIdx[field] == -1 {
					missingFields += field + ","
				}
			}
			missingFields = strings.TrimRight(missingFields, ",")
			if missingFields != "" {
				return response.JSON{Code: response.Error, Message: fmt.Sprintf("missing field %s", missingFields)}
			}

		} else {
			wg.Add(1)
			go func() {
				errMsg := ""
				srv := model.Server{
					NamespaceID: gp.Namespace.ID,
				}
				srv.Name = record[headerIdx["name"]]
				err = validator.Validate.Var(srv.Name, "required")
				if err != nil {
					errMsg += "name,"
				}

				srv.OS = record[headerIdx["os"]]
				err = validator.Validate.Var(srv.OS, "oneof=linux windows")
				if err != nil {
					errMsg += "os,"
				}

				srv.IP = record[headerIdx["host"]]
				err = validator.Validate.Var(srv.IP, "ip|hostname")
				if err != nil {
					errMsg += "host,"
				}

				srv.Port, err = strconv.Atoi(record[headerIdx["port"]])
				if err != nil {
					errMsg += "port,"
				}

				srv.Owner = record[headerIdx["owner"]]
				err = validator.Validate.Var(srv.Owner, "required,max=255")
				if err != nil {
					errMsg += "owner,"
				}

				srv.Path = record[headerIdx["path"]]
				err = validator.Validate.Var(record[headerIdx["path"]], "max=255")
				if err != nil {
					errMsg += "path,"
				}

				if headerIdx["password"] != -1 {
					srv.Password = record[headerIdx["password"]]
				}
				if headerIdx["description"] != -1 {
					srv.Description = record[headerIdx["description"]]
				}
				if headerIdx["jumpHost"] != -1 {
					srv.JumpIP = record[headerIdx["jumpHost"]]
				}
				if headerIdx["jumpPort"] != -1 {
					srv.JumpPort, _ = strconv.Atoi(record[headerIdx["jumpPort"]])
				}
				if headerIdx["jumpOwner"] != -1 {
					srv.JumpOwner = record[headerIdx["jumpOwner"]]
				}
				if headerIdx["jumpPath"] != -1 {
					srv.JumpPath = record[headerIdx["jumpPath"]]
				}
				if headerIdx["jumpPassword"] != -1 {
					srv.JumpPassword = record[headerIdx["jumpPassword"]]
				}
				errMsg = strings.TrimRight(errMsg, ",")
				if errMsg != "" {
					errOccur = true
					log.Error(fmt.Sprintf("Error on No.%d line %s, field validation on %s failed", i, record, errMsg))
				} else {
					srv.OSInfo = srv.ToSSHConfig().GetOSInfo()
					if _, err := srv.AddRow(); err != nil {
						errOccur = true
						log.Error(fmt.Sprintf("Error on No.%d line %s, %s", i, record, err.Error()))
					}
				}

				wg.Done()
			}()
		}
	}
	wg.Wait()

	if errOccur {
		return response.JSON{Code: response.Error, Message: "Encountered some unknown errors, please check the log details"}
	}

	return response.JSON{}
}

// Add adds the servers
// @Summary Add the servers
// @Tags Server
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body server.Add.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /server/add [post]
func (s Server) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name         string `json:"name" validate:"required"`
		NamespaceID  int64  `json:"namespaceId" validate:"gte=0"`
		OS           string `json:"os" validate:"oneof=linux windows"`
		IP           string `json:"ip" validate:"ip|hostname"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Owner        string `json:"owner" validate:"required,max=255"`
		Path         string `json:"path" validate:"max=255"`
		Password     string `json:"password"`
		Description  string `json:"description" validate:"max=255"`
		JumpIP       string `json:"jumpIP" validate:"omitempty,ip|hostname"`
		JumpPort     int    `json:"jumpPort" validate:"min=0,max=65535"`
		JumpOwner    string `json:"jumpOwner" validate:"max=255"`
		JumpPath     string `json:"jumpPath" validate:"max=255"`
		JumpPassword string `json:"jumpPassword"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv := model.Server{
		NamespaceID:  reqData.NamespaceID,
		Name:         reqData.Name,
		OS:           reqData.OS,
		IP:           reqData.IP,
		Port:         reqData.Port,
		Owner:        reqData.Owner,
		Path:         reqData.Path,
		Password:     reqData.Password,
		JumpIP:       reqData.JumpIP,
		JumpPort:     reqData.JumpPort,
		JumpOwner:    reqData.JumpOwner,
		JumpPath:     reqData.JumpPath,
		JumpPassword: reqData.JumpPassword,
		Description:  reqData.Description,
	}
	srv.OSInfo = srv.ToSSHConfig().GetOSInfo()

	id, err := srv.AddRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

// Edit edits the servers
// @Summary Edit the servers
// @Tags Server
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body server.Edit.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /server/edit [put]
func (s Server) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"gt=0"`
		NamespaceID  int64  `json:"namespaceId" validate:"gte=0"`
		Name         string `json:"name" validate:"required"`
		OS           string `json:"os" validate:"oneof=linux windows"`
		IP           string `json:"ip" validate:"required,ip|hostname"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Owner        string `json:"owner" validate:"required,max=255"`
		Path         string `json:"path" validate:"max=255"`
		Password     string `json:"password" validate:"max=255"`
		Description  string `json:"description" validate:"max=255"`
		JumpIP       string `json:"jumpIP" validate:"omitempty,ip|hostname"`
		JumpPort     int    `json:"jumpPort" validate:"min=0,max=65535"`
		JumpOwner    string `json:"jumpOwner" validate:"max=255"`
		JumpPath     string `json:"jumpPath" validate:"max=255"`
		JumpPassword string `json:"jumpPassword"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	srv := model.Server{
		ID:           reqData.ID,
		NamespaceID:  reqData.NamespaceID,
		Name:         reqData.Name,
		OS:           reqData.OS,
		IP:           reqData.IP,
		Port:         reqData.Port,
		Owner:        reqData.Owner,
		Path:         reqData.Path,
		Password:     reqData.Password,
		JumpIP:       reqData.JumpIP,
		JumpPort:     reqData.JumpPort,
		JumpOwner:    reqData.JumpOwner,
		JumpPath:     reqData.JumpPath,
		JumpPassword: reqData.JumpPassword,
		Description:  reqData.Description,
	}
	srv.OSInfo = srv.ToSSHConfig().GetOSInfo()

	if err := srv.EditRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

// Toggle toggles the server state
// @Summary Toggle the server state
// @Tags Server
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body server.Toggle.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /server/toggle [put]
func (Server) Toggle(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID    int64 `json:"id" validate:"gt=0"`
		State int8  `json:"state" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Server{ID: reqData.ID, State: reqData.State}).ToggleRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

// UnbindProject unbinds the projects
// @Summary Unbind the projects
// @Tags Server
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body server.UnbindProject.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /server/unbindProject [delete]
func (Server) UnbindProject(gp *server.Goploy) server.Response {
	type ReqData struct {
		IDs []int64 `json:"ids" validate:"required,min=1"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectServer{}).DeleteInID(reqData.IDs); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) InstallAgent(gp *server.Goploy) server.Response {
	type ReqData struct {
		IDs         []int64 `json:"ids" validate:"required,min=1"`
		InstallPath string  `json:"installPath" validate:"required"`
		Tool        string  `json:"tool" validate:"required"`
		ReportURL   string  `json:"reportURL" validate:"required"`
		WebPort     string  `json:"webPort" validate:"omitempty"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	downloadURL := "https://github.com/goploy-devops/goploy-agent/releases/latest/download/goploy-agent"
	downloadCommand := fmt.Sprintf("wget %s --connect-timeout=7 --dns-timeout=7 -nv -O goploy-agent >> /dev/null 2>&1", downloadURL)
	if reqData.Tool == "curl" {
		downloadCommand = fmt.Sprintf("curl %s -o goploy-agent", downloadURL)
	}

	for _, id := range reqData.IDs {
		go func(id int64) {
			srv, err := (model.Server{ID: id}).GetData()
			if err != nil {
				log.Error(fmt.Sprintf("Error on %d server, %s", id, err.Error()))
				return
			}
			client, err := srv.ToSSHConfig().Dial()
			if err != nil {
				log.Error(fmt.Sprintf("Error on %d server, %s", id, err.Error()))
				return
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				log.Error(fmt.Sprintf("Error on %d server, %s", id, err.Error()))
				return
			}
			defer session.Close()
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			commands := []string{
				fmt.Sprintf("mkdir -p %s", reqData.InstallPath),
				fmt.Sprintf("cd %s", reqData.InstallPath),
				fmt.Sprintf("[./goploy-agent -s stop || true"),
				downloadCommand,
				"touch ./goploy-agent.toml",
				"echo env = \\'production\\' > ./goploy-agent.toml",
				"echo [goploy] >> ./goploy-agent.toml",
				fmt.Sprintf("echo reportURL = \\'%s\\' >> ./goploy-agent.toml", reqData.ReportURL),
				fmt.Sprintf("echo key = \\'%s\\' >> ./goploy-agent.toml", config.Toml.JWT.Key),
				"echo uidType = \\'id\\' >> ./goploy-agent.toml",
				fmt.Sprintf("echo uid = '%d' >> ./goploy-agent.toml", id),
				"echo [log] >> ./goploy-agent.toml",
				"echo path = \\'stdout\\' >> ./goploy-agent.toml",
				"echo [web] >> ./goploy-agent.toml",
				fmt.Sprintf("echo port = '%s' >> ./goploy-agent.toml", reqData.WebPort),
				"chmod a+x ./goploy-agent",
				"nohup ./goploy-agent &",
			}
			if err := session.Run(strings.Join(commands, "&&")); err != nil {
				log.Error(fmt.Sprintf("Error on %d server, %s, detail: %s", id, err.Error(), sshErrbuf.String()))
				return
			}
			log.Info(sshErrbuf.String())
		}(id)
	}

	return response.JSON{}
}

func (Server) PreviewFile(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: "invalid server id"}
	}
	srv, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.SftpFile{Filename: gp.URLQuery.Get("file"), Client: client, Disposition: "inline"}
}

func (Server) DownloadFile(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: "invalid server id"}
	}
	srv, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.SftpFile{Filename: gp.URLQuery.Get("file"), Client: client, Disposition: "attachment"}
}

func (Server) UploadFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID       int64  `schema:"id" validate:"gt=0"`
		FilePath string `schema:"filePath"  validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	file, fileHandler, err := gp.Request.FormFile("file")
	if err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}
	defer file.Close()

	srv, err := (model.Server{ID: reqData.ID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Create(reqData.FilePath + "/" + fileHandler.Filename)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, file)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Server) EditFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64  `json:"serverId" validate:"gt=0"`
		File     string `json:"file" validate:"required"`
		Content  string `json:"content" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Create(reqData.File)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	_, err = remoteFile.Write([]byte(reqData.Content))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Server) CopyFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64  `json:"serverId" validate:"gt=0"`
		IsDir    bool   `json:"isDir"`
		Dir      string `json:"dir" validate:"required"`
		SrcName  string `json:"srcName" validate:"required"`
		DstName  string `json:"dstName" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer session.Close()

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	optionR := ""
	if reqData.IsDir {
		optionR = "-r"
	}
	if err = session.Run(fmt.Sprintf("cp %s %s %s", optionR, path.Join(reqData.Dir, reqData.SrcName), path.Join(reqData.Dir, reqData.DstName))); err != nil {
		return response.JSON{Code: response.Error, Message: "err: " + err.Error() + ", detail: " + sshErrbuf.String()}
	}
	return response.JSON{}
}

func (Server) RenameFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID    int64  `json:"serverId" validate:"gt=0"`
		Dir         string `json:"dir" validate:"required"`
		NewName     string `json:"newName" validate:"required"`
		CurrentName string `json:"currentName" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()
	err = sftpClient.Rename(path.Join(reqData.Dir, reqData.CurrentName), path.Join(reqData.Dir, reqData.NewName))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) DeleteFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		File     string `json:"file" validate:"required"`
		ServerID int64  `json:"serverId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	fi, err := sftpClient.Stat(reqData.File)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if fi.IsDir() == true {
		err = sftpClient.RemoveDirectory(reqData.File)
	} else {
		err = sftpClient.Remove(reqData.File)
	}

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Server) TransferFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		SourceServerID int64   `json:"sourceServerId" validate:"required"`
		SourceFile     string  `json:"sourceFile" validate:"required"`
		SourceIsDir    bool    `json:"sourceIsDir"`
		DestServerIDs  []int64 `json:"destServerIds" validate:"min=1"`
		DestDir        string  `json:"destDir" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	sourceServer, err := (model.Server{ID: reqData.SourceServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := sourceServer.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	for _, destServerID := range reqData.DestServerIDs {
		destServer, err := (model.Server{ID: destServerID}).GetData()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		err = func() error {
			destSSHClient, err := destServer.ToSSHConfig().Dial()
			if err != nil {
				return err
			}
			defer destSSHClient.Close()

			destSFTPClient, err := sftp.NewClient(destSSHClient)
			if err != nil {
				return err
			}
			defer destSFTPClient.Close()

			if reqData.SourceIsDir == false {
				if err := destSFTPClient.MkdirAll(reqData.DestDir); err != nil {
					return err
				}

				srcFile, err := sftpClient.Open(reqData.SourceFile)
				if err != nil {
					return err
				}
				defer srcFile.Close()

				destFile, err := destSFTPClient.Create(reqData.DestDir + "/" + path.Base(reqData.SourceFile))
				if err != nil {
					return err
				}
				defer destFile.Close()

				if _, err = io.Copy(destFile, srcFile); err != nil {
					return err
				}
			} else {
				w := sftpClient.Walk(reqData.SourceFile)
				// skip root dir
				if w.Step(); w.Err() != nil {
					return w.Err()
				}
				for w.Step() {
					if w.Err() != nil {
						continue
					}
					fileInfo := w.Stat()
					filePath := w.Path()
					destTarget := path.Join(reqData.DestDir, filePath[len(reqData.SourceFile):])
					if fileInfo.IsDir() {
						if err := destSFTPClient.MkdirAll(destTarget); err != nil {
							return err
						}
					} else {
						err := func() error {
							srcFile, err := sftpClient.Open(filePath)
							if err != nil {
								return err
							}
							defer srcFile.Close()

							destFile, err := destSFTPClient.Create(destTarget)
							if err != nil {
								return err
							}
							defer destFile.Close()

							if _, err := io.Copy(destFile, srcFile); err != nil {
								return err
							}
							return nil
						}()
						if err != nil {
							return err
						}
					}
				}
			}
			return nil
		}()

		if err != nil {
			return response.JSON{
				Code:    response.Error,
				Message: err.Error(),
				Data: struct {
					ServerID   int64  `json:"serverId"`
					ServerName string `json:"serverName"`
				}{destServerID, destServer.Name},
			}
		}

	}

	return response.JSON{}
}

func (Server) Report(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID      int64  `schema:"serverId" validate:"gt=0"`
		Type          int    `schema:"type" validate:"gt=0"`
		DatetimeRange string `schema:"datetimeRange"  validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	datetimeRange := strings.Split(reqData.DatetimeRange, ",")
	if len(datetimeRange) != 2 {
		return response.JSON{Code: response.Error, Message: "invalid datetime range"}
	}
	serverAgentLogs, err := (model.ServerAgentLog{ServerID: reqData.ServerID, Type: reqData.Type}).GetListBetweenTime(datetimeRange[0], datetimeRange[1])
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type Flag struct {
		Count int
		Curr  int
	}

	flagMap := map[string]Flag{}

	for _, serverAgentLog := range serverAgentLogs {
		if _, ok := flagMap[serverAgentLog.Item]; !ok {
			flagMap[serverAgentLog.Item] = Flag{}
		}
		flagMap[serverAgentLog.Item] = Flag{Count: flagMap[serverAgentLog.Item].Count + 1}
	}

	serverAgentMap := map[string]model.ServerAgentLogs{}
	for _, serverAgentLog := range serverAgentLogs {
		flagMap[serverAgentLog.Item] = Flag{
			Count: flagMap[serverAgentLog.Item].Count,
			Curr:  flagMap[serverAgentLog.Item].Curr + 1,
		}
		step := flagMap[serverAgentLog.Item].Count / 60
		if flagMap[serverAgentLog.Item].Count <= 60 ||
			flagMap[serverAgentLog.Item].Curr%step == 0 ||
			flagMap[serverAgentLog.Item].Count-1 == flagMap[serverAgentLog.Item].Curr {
			serverAgentMap[serverAgentLog.Item] = append(serverAgentMap[serverAgentLog.Item], serverAgentLog)
		}
	}

	return response.JSON{
		Data: struct {
			ServerAgentMap map[string]model.ServerAgentLogs `json:"map"`
		}{ServerAgentMap: serverAgentMap},
	}
}

func (Server) GetAllMonitor(gp *server.Goploy) server.Response {
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	serverMonitorList, err := model.ServerMonitor{ServerID: serverID}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.ServerMonitors `json:"list"`
		}{List: serverMonitorList},
	}
}

func (s Server) AddMonitor(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID     int64  `json:"serverId" validate:"required"`
		Item         string `json:"item" validate:"required"`
		Formula      string `json:"formula" validate:"required"`
		Operator     string `json:"operator" validate:"required"`
		Value        string `json:"value" validate:"required"`
		GroupCycle   int    `json:"groupCycle" validate:"required"`
		LastCycle    int    `json:"lastCycle" validate:"required"`
		SilentCycle  int    `json:"silentCycle" validate:"required"`
		StartTime    string `json:"startTime" validate:"required,len=5"`
		EndTime      string `json:"endTime" validate:"required,len=5"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.ServerMonitor{
		ServerID:     reqData.ServerID,
		Item:         reqData.Item,
		Formula:      reqData.Formula,
		Operator:     reqData.Operator,
		Value:        reqData.Value,
		GroupCycle:   reqData.GroupCycle,
		LastCycle:    reqData.LastCycle,
		SilentCycle:  reqData.SilentCycle,
		StartTime:    reqData.StartTime,
		EndTime:      reqData.EndTime,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (s Server) EditMonitor(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"required"`
		Item         string `json:"item" validate:"required"`
		Formula      string `json:"formula" validate:"required"`
		Operator     string `json:"operator" validate:"required"`
		Value        string `json:"value" validate:"required"`
		GroupCycle   int    `json:"groupCycle" validate:"required"`
		LastCycle    int    `json:"lastCycle" validate:"required"`
		SilentCycle  int    `json:"silentCycle" validate:"required"`
		StartTime    string `json:"startTime" validate:"required,len=5"`
		EndTime      string `json:"endTime" validate:"required,len=5"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.ServerMonitor{
		ID:           reqData.ID,
		Item:         reqData.Item,
		Formula:      reqData.Formula,
		Operator:     reqData.Operator,
		Value:        reqData.Value,
		GroupCycle:   reqData.GroupCycle,
		LastCycle:    reqData.LastCycle,
		SilentCycle:  reqData.SilentCycle,
		StartTime:    reqData.StartTime,
		EndTime:      reqData.EndTime,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{}
}

func (s Server) DeleteMonitor(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.ServerMonitor{
		ID: reqData.ID,
	}.DeleteRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{}
}

func (Server) GetProcessList(gp *server.Goploy) server.Response {
	list, err := model.ServerProcess{NamespaceID: gp.Namespace.ID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.ServerProcesses `json:"list"`
		}{List: list},
	}
}

func (Server) AddProcess(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name  string `json:"name" validate:"required"`
		Items string `json:"items"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.ServerProcess{
		NamespaceID: gp.Namespace.ID,
		Name:        reqData.Name,
		Items:       reqData.Items,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Server) EditProcess(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID    int64  `json:"id" validate:"gt=0"`
		Name  string `json:"name" validate:"required"`
		Items string `json:"items"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.ServerProcess{
		ID:    reqData.ID,
		Name:  reqData.Name,
		Items: reqData.Items,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) DeleteProcess(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ServerProcess{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Server) ExecProcess(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID       int64  `json:"id" validate:"gt=0"`
		ServerID int64  `json:"serverId" validate:"gt=0"`
		Name     string `json:"name" validate:"required"`
	}

	type RespData struct {
		ServerID int64  `json:"serverId"`
		ExecRes  bool   `json:"execRes"`
		Stdout   string `json:"stdout"`
		Stderr   string `json:"stderr"`
	}

	var reqData ReqData

	var respData RespData
	respData.ExecRes = false
	respData.ServerID = reqData.ServerID

	if err := gp.Decode(&reqData); err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}

	serverProcess, err := model.ServerProcess{ID: reqData.ID}.GetData()
	if err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}
	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}

	var processItems model.ServerProcessItems
	if err := json.Unmarshal([]byte(serverProcess.Items), &processItems); err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}

	script := ""
	for _, processItem := range processItems {
		if processItem.Name == reqData.Name {
			script = processItem.Command
			break
		}
	}

	if script == "" {
		return response.JSON{Code: response.Error, Message: "Command empty"}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}
	defer session.Close()

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	err = session.Run(script)
	respData.ExecRes = err == nil
	respData.Stdout = sshOutbuf.String()
	respData.Stderr = sshErrbuf.String()
	return response.JSON{Data: respData}
}

func (Server) ExecScript(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerIDs []int64 `json:"serverIds" validate:"gt=0"`
		Script    string  `json:"script" validate:"required"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type ServerResp struct {
		ServerID int64  `json:"serverId"`
		ExecRes  bool   `json:"execRes"`
		Stdout   string `json:"stdout"`
		Stderr   string `json:"stderr"`
	}

	ch := make(chan ServerResp, len(reqData.ServerIDs))

	for _, serverId := range reqData.ServerIDs {
		go func(serverId int64) {
			serverResp := ServerResp{
				ServerID: serverId,
			}

			srv, err := (model.Server{ID: serverId}).GetData()
			if err != nil {
				serverResp.ExecRes = false
				serverResp.Stderr = err.Error()
				ch <- serverResp
				return
			}
			client, err := srv.ToSSHConfig().Dial()
			if err != nil {
				serverResp.ExecRes = false
				serverResp.Stderr = err.Error()
				ch <- serverResp
				return
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				serverResp.ExecRes = false
				serverResp.Stderr = err.Error()
				ch <- serverResp
				return
			}
			defer session.Close()

			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			err = session.Run(srv.ReplaceVars(reqData.Script))
			serverResp.ExecRes = err == nil
			serverResp.Stdout = sshOutbuf.String()
			serverResp.Stderr = sshErrbuf.String()
			ch <- serverResp
			return
		}(serverId)
	}

	var respData []ServerResp
	for range reqData.ServerIDs {
		respData = append(respData, <-ch)
	}
	return response.JSON{Data: respData}
}

func (Server) GetNginxPath(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64 `schema:"serverId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer session.Close()

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf

	command := `ls -l /proc/$(ps -ef | grep nginx | grep master | grep -v ps | awk '{print $2}')/exe | awk '{print $NF}'`
	nginxPath := ""
	if err = session.Run(command); err == nil {
		nginxPath = strings.Trim(sshOutbuf.String(), "\n")
		if len(strings.Split(nginxPath, "\n")) > 1 {
			nginxPath = ""
		} else if !pkg.IsFilePath(nginxPath) {
			nginxPath = ""
		}
	}

	return response.JSON{
		Data: struct {
			Path string `json:"path"`
		}{Path: nginxPath},
	}
}

func (Server) GetNginxConfigList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64  `schema:"serverId" validate:"gt=0"`
		Dir      string `schema:"dir" validate:"filepath"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer session.Close()

	output, err := session.CombinedOutput(reqData.Dir + " -t")

	if err != nil {
		return response.JSON{Code: response.Error, Message: fmt.Sprintf("output: %s", output)}
	}

	configPath := ""
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "nginx: the configuration file ") && strings.Contains(line, "syntax is ok") {
			configPath = strings.TrimPrefix(line, "nginx: the configuration file ")
			configPath = strings.TrimSuffix(configPath, " syntax is ok")
		}
	}

	if configPath == "" {
		return response.JSON{Code: response.Error, Message: "can not find nginx config path or config error"}
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	configFile, err := sftpClient.Open(configPath)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer configFile.Close()

	configContent, err := io.ReadAll(configFile)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	configFileDir := path.Dir(configPath)

	includeConfigPaths := []string{configPath}
	// match the file path in the include directive
	re := regexp.MustCompile("(?i)include\\s+([\"']?)(.*?)([\"']?);")
	matches := re.FindAllSubmatch(configContent, -1)

	for _, match := range matches {
		tmpPath := match[2]
		includeConfigPaths = append(includeConfigPaths, path.Join(configFileDir, string(tmpPath)))
	}

	type fileInfo struct {
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Mode    string `json:"mode"`
		ModTime string `json:"modTime"`
		Dir     string `json:"dir"`
	}

	var fileList []fileInfo
	for _, includeConfigPath := range includeConfigPaths {
		filePaths, err := sftpClient.Glob(includeConfigPath)
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		for _, f := range filePaths {
			fileStat, err := sftpClient.Stat(f)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}

			if !fileStat.IsDir() && fileStat.Mode()&os.ModeSymlink != 0 {
				continue
			}

			fileList = append(fileList, fileInfo{
				Name:    fileStat.Name(),
				Size:    fileStat.Size(),
				Mode:    fileStat.Mode().String(),
				ModTime: fileStat.ModTime().Format("2006-01-02 15:04:05"),
				Dir:     path.Dir(includeConfigPath),
			})
		}
	}

	return response.JSON{
		Data: struct {
			List []fileInfo `json:"list"`
		}{List: fileList},
	}
}

func (Server) ManageNginx(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64  `json:"serverId" validate:"gt=0"`
		Path     string `json:"path" validate:"filepath"`
		Command  string `json:"command" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	script := ""
	switch reqData.Command {
	case "reload":
		script = reqData.Path + " -s reload"
	case "check":
		script = reqData.Path + " -t"
	default:
		return response.JSON{Code: response.Error, Message: "Command error"}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	if reqData.Command == "reload" {
		checkSession, err := client.NewSession()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		defer checkSession.Close()

		output, err := checkSession.CombinedOutput(reqData.Path + " -t")
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}

		outputString := string(output)
		if !strings.Contains(outputString, "syntax is ok") || !strings.Contains(outputString, "test is successful") {
			return response.JSON{Code: response.Error, Message: outputString}
		}
	}

	session, err := client.NewSession()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer session.Close()

	output, err := session.CombinedOutput(script)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	outputString := string(output)

	return response.JSON{
		Data: struct {
			ExecRes bool   `json:"execRes"`
			Output  string `json:"output"`
		}{ExecRes: err == nil, Output: outputString},
	}
}

func (Server) GetNginxConfContent(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64  `schema:"serverId" validate:"gt=0"`
		Dir      string `schema:"dir" validate:"required"`
		Filename string `schema:"filename" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	configFile, err := sftpClient.Open(path.Join(reqData.Dir, reqData.Filename))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer configFile.Close()

	configContent, err := io.ReadAll(configFile)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			Content string `json:"content"`
		}{Content: string(configContent)},
	}
}

func (Server) EditNginxConfig(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64  `json:"serverId" validate:"gt=0"`
		Dir      string `json:"dir" validate:"required"`
		Filename string `json:"filename" validate:"required"`
		Content  string `json:"content" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	file, err := sftpClient.Create(path.Join(reqData.Dir, reqData.Filename))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer file.Close()

	_, err = file.Write([]byte(reqData.Content))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Server) CopyNginxConfig(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID int64  `json:"serverId" validate:"gt=0"`
		Dir      string `json:"dir" validate:"required"`
		SrcName  string `json:"srcName" validate:"required"`
		DstName  string `json:"dstName" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer session.Close()

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf

	if err = session.Run(fmt.Sprintf("cp %s %s", path.Join(reqData.Dir, reqData.SrcName), path.Join(reqData.Dir, reqData.DstName))); err != nil {
		return response.JSON{Code: response.Error, Message: "err: " + err.Error() + ", detail: " + sshErrbuf.String()}
	}
	return response.JSON{}
}

func (Server) RenameNginxConfig(gp *server.Goploy) server.Response {
	type ReqData struct {
		ServerID    int64  `json:"serverId" validate:"gt=0"`
		Dir         string `json:"dir" validate:"filepath"`
		NewName     string `json:"newName" validate:"required"`
		CurrentName string `json:"currentName" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	err = sftpClient.Rename(path.Join(reqData.Dir, reqData.CurrentName), path.Join(reqData.Dir, reqData.NewName))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) DeleteNginxConfig(gp *server.Goploy) server.Response {
	type ReqData struct {
		File     string `json:"file" validate:"required"`
		ServerID int64  `json:"serverId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srv, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	err = sftpClient.Remove(reqData.File)

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
