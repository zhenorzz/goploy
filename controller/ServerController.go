// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package controller

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/permission"
	"github.com/zhenorzz/goploy/response"
	"github.com/zhenorzz/goploy/utils"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
)

// Server struct
type Server Controller

func (s Server) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/server/getList", http.MethodGet, s.GetList).Permissions(permission.ShowServerPage),
		core.NewRoute("/server/getOption", http.MethodGet, s.GetOption),
		core.NewRoute("/server/getPublicKey", http.MethodGet, s.GetPublicKey).Permissions(permission.AddServer, permission.EditServer).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/check", http.MethodPost, s.Check).Permissions(permission.AddServer, permission.EditServer).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/import", http.MethodPost, s.Import).Permissions(permission.ImportCSV).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/add", http.MethodPost, s.Add).Permissions(permission.AddServer).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/edit", http.MethodPut, s.Edit).Permissions(permission.EditServer).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/toggle", http.MethodPut, s.Toggle).Permissions(permission.EditServer).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/installAgent", http.MethodPost, s.InstallAgent).Permissions(permission.InstallAgent).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/previewFile", http.MethodGet, s.PreviewFile).Permissions(permission.SFTPPreviewFile).LogFunc(middleware.AddPreviewLog),
		core.NewRoute("/server/downloadFile", http.MethodGet, s.DownloadFile).Permissions(permission.SFTPDownloadFile).LogFunc(middleware.AddDownloadLog),
		core.NewRoute("/server/uploadFile", http.MethodPost, s.UploadFile).Permissions(permission.SFTPTransferFile).LogFunc(middleware.AddUploadLog),
		core.NewRoute("/server/deleteFile", http.MethodDelete, s.DeleteFile).Permissions(permission.SFTPDeleteFile).LogFunc(middleware.AddDeleteLog),
		core.NewRoute("/server/transferFile", http.MethodPost, s.TransferFile).Permissions(permission.SFTPUploadFile),
		core.NewRoute("/server/report", http.MethodGet, s.Report).Permissions(permission.ShowServerMonitorPage),
		core.NewRoute("/server/getAllMonitor", http.MethodGet, s.GetAllMonitor).Permissions(permission.ShowServerMonitorPage),
		core.NewRoute("/server/addMonitor", http.MethodPost, s.AddMonitor).Permissions(permission.AddServerWarningRule).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/editMonitor", http.MethodPut, s.EditMonitor).Permissions(permission.EditServerWarningRule).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/deleteMonitor", http.MethodDelete, s.DeleteMonitor).Permissions(permission.DeleteServerWarningRule).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/getProcessList", http.MethodGet, s.GetProcessList).Permissions(permission.ShowServerProcessPage),
		core.NewRoute("/server/addProcess", http.MethodPost, s.AddProcess).Permissions(permission.AddServerProcess).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/editProcess", http.MethodPut, s.EditProcess).Permissions(permission.EditServerProcess).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/deleteProcess", http.MethodDelete, s.DeleteProcess).Permissions(permission.DeleteServerProcess).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/execProcess", http.MethodPost, s.ExecProcess).Permissions(permission.ShowServerProcessPage).LogFunc(middleware.AddOPLog),
		core.NewRoute("/server/getRemoteCrontabList", http.MethodGet, s.GetRemoteCrontabList).Permissions(permission.ShowServerProcessPage),
	}
}

func (Server) GetList(gp *core.Goploy) core.Response {
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

func (Server) GetOption(gp *core.Goploy) core.Response {
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

func (Server) GetPublicKey(gp *core.Goploy) core.Response {
	publicKeyPath := gp.URLQuery.Get("path")

	contentByte, err := ioutil.ReadFile(publicKeyPath + ".pub")
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Key string `json:"key"`
		}{Key: string(contentByte)},
	}
}

func (Server) Check(gp *core.Goploy) core.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	sshConfig := utils.SSHConfig{
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

func (Server) Import(gp *core.Goploy) core.Response {
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
				server := model.Server{
					NamespaceID: gp.Namespace.ID,
				}
				server.Name = record[headerIdx["name"]]
				err = core.Validate.Var(server.Name, "required")
				if err != nil {
					errMsg += "name,"
				}

				server.OS = record[headerIdx["os"]]
				err = core.Validate.Var(server.OS, "oneof=linux windows")
				if err != nil {
					errMsg += "os,"
				}

				server.IP = record[headerIdx["host"]]
				err = core.Validate.Var(server.IP, "ip|hostname")
				if err != nil {
					errMsg += "host,"
				}

				server.Port, err = strconv.Atoi(record[headerIdx["port"]])
				if err != nil {
					errMsg += "port,"
				}

				server.Owner = record[headerIdx["owner"]]
				err = core.Validate.Var(server.Owner, "required,max=255")
				if err != nil {
					errMsg += "owner,"
				}

				server.Path = record[headerIdx["path"]]
				err = core.Validate.Var(record[headerIdx["path"]], "max=255")
				if err != nil {
					errMsg += "path,"
				}

				if headerIdx["password"] != -1 {
					server.Password = record[headerIdx["password"]]
				}
				if headerIdx["description"] != -1 {
					server.Description = record[headerIdx["description"]]
				}
				if headerIdx["jumpHost"] != -1 {
					server.JumpIP = record[headerIdx["jumpHost"]]
				}
				if headerIdx["jumpPort"] != -1 {
					server.JumpPort, _ = strconv.Atoi(record[headerIdx["jumpPort"]])
				}
				if headerIdx["jumpOwner"] != -1 {
					server.JumpOwner = record[headerIdx["jumpOwner"]]
				}
				if headerIdx["jumpPath"] != -1 {
					server.JumpPath = record[headerIdx["jumpPath"]]
				}
				if headerIdx["jumpPassword"] != -1 {
					server.JumpPassword = record[headerIdx["jumpPassword"]]
				}
				errMsg = strings.TrimRight(errMsg, ",")
				if errMsg != "" {
					errOccur = true
					core.Log(core.ERROR, fmt.Sprintf("Error on No.%d line %s, field validation on %s failed", i, record, errMsg))
				} else {
					server.OSInfo = server.ToSSHConfig().GetOSInfo()
					if _, err := server.AddRow(); err != nil {
						errOccur = true
						core.Log(core.ERROR, fmt.Sprintf("Error on No.%d line %s, %s", i, record, err.Error()))
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

func (s Server) Add(gp *core.Goploy) core.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	server := model.Server{
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
	server.OSInfo = server.ToSSHConfig().GetOSInfo()

	id, err := server.AddRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (s Server) Edit(gp *core.Goploy) core.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	server := model.Server{
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
	server.OSInfo = server.ToSSHConfig().GetOSInfo()

	if err := server.EditRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) Toggle(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID    int64 `json:"id" validate:"gt=0"`
		State int8  `json:"state" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Server{ID: reqData.ID, State: reqData.State}).ToggleRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) InstallAgent(gp *core.Goploy) core.Response {
	type ReqData struct {
		IDs         []int64 `json:"ids" validate:"min=1"`
		InstallPath string  `json:"installPath" validate:"required"`
		Tool        string  `json:"tool" validate:"required"`
		ReportURL   string  `json:"reportURL" validate:"required"`
		WebPort     string  `json:"webPort" validate:"omitempty"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	downloadURL := "https://github.com/goploy-devops/goploy-agent/releases/latest/download/goploy-agent"
	downloadCommand := fmt.Sprintf("wget -N %s", downloadURL)
	if reqData.Tool == "curl" {
		downloadCommand = fmt.Sprintf("curl %s -o goploy-agent", downloadURL)
	}

	for _, id := range reqData.IDs {
		go func(id int64) {
			server, err := (model.Server{ID: id}).GetData()
			if err != nil {
				core.Log(core.ERROR, fmt.Sprintf("Error on %d server, %s", id, err.Error()))
				return
			}
			client, err := server.ToSSHConfig().Dial()
			if err != nil {
				core.Log(core.ERROR, fmt.Sprintf("Error on %d server, %s", id, err.Error()))
				return
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				core.Log(core.ERROR, fmt.Sprintf("Error on %d server, %s", id, err.Error()))
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
				"echo env = 'production' > ./goploy-agent.toml",
				"echo [goploy] >> ./goploy-agent.toml",
				fmt.Sprintf("echo reportURL = '%s' >> ./goploy-agent.toml", reqData.ReportURL),
				fmt.Sprintf("echo key = '%s' >> ./goploy-agent.toml", config.Toml.JWT.Key),
				"echo uidType = 'id' >> ./goploy-agent.toml",
				fmt.Sprintf("echo uid = '%d' >> ./goploy-agent.toml", id),
				"echo [log] >> ./goploy-agent.toml",
				"echo path = 'stdout' >> ./goploy-agent.toml",
				"echo [web] >> ./goploy-agent.toml",
				fmt.Sprintf("echo port = '%s' >> ./goploy-agent.toml", reqData.WebPort),
				"chmod a+x ./goploy-agent",
				"nohup ./goploy-agent &",
			}
			if err := session.Run(strings.Join(commands, "&&")); err != nil {
				core.Log(core.ERROR, fmt.Sprintf("Error on %d server, %s, detail: %s", id, err.Error(), sshErrbuf.String()))
				return
			}
			core.Log(core.INFO, sshErrbuf.String())
		}(id)
	}

	return response.JSON{}
}

func (Server) PreviewFile(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: "invalid server id"}
	}
	server, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := server.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.SftpFile{Filename: gp.URLQuery.Get("file"), Client: client, Disposition: "inline"}
}

func (Server) DownloadFile(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: "invalid server id"}
	}
	server, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := server.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.SftpFile{Filename: gp.URLQuery.Get("file"), Client: client, Disposition: "attachment"}
}

func (Server) UploadFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID       int64  `schema:"id" validate:"gt=0"`
		FilePath string `schema:"filePath"  validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	file, fileHandler, err := gp.Request.FormFile("file")
	if err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}
	defer file.Close()

	server, err := (model.Server{ID: reqData.ID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := server.ToSSHConfig().Dial()
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

func (Server) DeleteFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		File     string `json:"file" validate:"required"`
		ServerID int64  `json:"serverId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	server, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := server.ToSSHConfig().Dial()
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

func (Server) TransferFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		SourceServerID int64   `json:"sourceServerId" validate:"required"`
		SourceFile     string  `json:"sourceFile" validate:"required"`
		SourceIsDir    bool    `json:"sourceIsDir"`
		DestServerIDs  []int64 `json:"destServerIds" validate:"min=1"`
		DestDir        string  `json:"destDir" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Server) Report(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID      int64  `schema:"serverId" validate:"gt=0"`
		Type          int    `schema:"type" validate:"gt=0"`
		DatetimeRange string `schema:"datetimeRange"  validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
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

	for _, log := range serverAgentLogs {
		if _, ok := flagMap[log.Item]; !ok {
			flagMap[log.Item] = Flag{}
		}
		flagMap[log.Item] = Flag{Count: flagMap[log.Item].Count + 1}
	}

	serverAgentMap := map[string]model.ServerAgentLogs{}
	for _, log := range serverAgentLogs {
		flagMap[log.Item] = Flag{
			Count: flagMap[log.Item].Count,
			Curr:  flagMap[log.Item].Curr + 1,
		}
		step := flagMap[log.Item].Count / 60
		if flagMap[log.Item].Count <= 60 ||
			flagMap[log.Item].Curr%step == 0 ||
			flagMap[log.Item].Count-1 == flagMap[log.Item].Curr {
			serverAgentMap[log.Item] = append(serverAgentMap[log.Item], log)
		}
	}

	return response.JSON{
		Data: struct {
			ServerAgentMap map[string]model.ServerAgentLogs `json:"map"`
		}{ServerAgentMap: serverAgentMap},
	}
}

func (Server) GetAllMonitor(gp *core.Goploy) core.Response {
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

func (s Server) AddMonitor(gp *core.Goploy) core.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (s Server) EditMonitor(gp *core.Goploy) core.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (s Server) DeleteMonitor(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Server) GetProcessList(gp *core.Goploy) core.Response {
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

func (Server) AddProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name  string `json:"name" validate:"required"`
		Items string `json:"items"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Server) EditProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID    int64  `json:"id" validate:"gt=0"`
		Name  string `json:"name" validate:"required"`
		Items string `json:"items"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Server) DeleteProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ServerProcess{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Server) ExecProcess(gp *core.Goploy) core.Response {
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

	if err := decodeJson(gp.Body, &reqData); err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}

	serverProcess, err := model.ServerProcess{ID: reqData.ID}.GetData()
	if err != nil {
		respData.Stderr = err.Error()
		return response.JSON{Data: respData}
	}
	server, err := (model.Server{ID: reqData.ServerID}).GetData()
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

	client, err := server.ToSSHConfig().Dial()
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

func (Server) GetRemoteCrontabList(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID int64 `schema:"serverId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	server, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := server.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	b, err := session.CombinedOutput("crontab -l")
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	var crontabList []string
	for _, crontab := range strings.Split(string(b), "\n") {
		// windows \r\n
		crontab = strings.TrimRight(crontab, "\r")
		if len(crontab) == 0 {
			continue
		}
		// skip error format
		if len(strings.Split(crontab, " ")) < 5 {
			continue
		}
		// skip comment
		if strings.HasPrefix(crontab, "#") {
			continue
		}

		crontabList = append(crontabList, crontab)
	}

	fmt.Printf("%v", crontabList)
	return response.JSON{}
}
