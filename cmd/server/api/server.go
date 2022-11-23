// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/model"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
)

// Server struct
type Server API

func (s Server) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/server/getList", http.MethodGet, s.GetList).Permissions(config.ShowServerPage),
		server.NewRoute("/server/getOption", http.MethodGet, s.GetOption),
		server.NewRoute("/server/getPublicKey", http.MethodGet, s.GetPublicKey).Permissions(config.AddServer, config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/check", http.MethodPost, s.Check).Permissions(config.AddServer, config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/import", http.MethodPost, s.Import).Permissions(config.ImportCSV).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/add", http.MethodPost, s.Add).Permissions(config.AddServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/edit", http.MethodPut, s.Edit).Permissions(config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/toggle", http.MethodPut, s.Toggle).Permissions(config.EditServer).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/installAgent", http.MethodPost, s.InstallAgent).Permissions(config.InstallAgent).LogFunc(middleware.AddOPLog),
		server.NewRoute("/server/previewFile", http.MethodGet, s.PreviewFile).Permissions(config.SFTPPreviewFile).LogFunc(middleware.AddPreviewLog),
		server.NewRoute("/server/downloadFile", http.MethodGet, s.DownloadFile).Permissions(config.SFTPDownloadFile).LogFunc(middleware.AddDownloadLog),
		server.NewRoute("/server/uploadFile", http.MethodPost, s.UploadFile).Permissions(config.SFTPTransferFile).LogFunc(middleware.AddUploadLog),
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

func (Server) GetPublicKey(gp *server.Goploy) server.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
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
				err = Validate.Var(srv.Name, "required")
				if err != nil {
					errMsg += "name,"
				}

				srv.OS = record[headerIdx["os"]]
				err = Validate.Var(srv.OS, "oneof=linux windows")
				if err != nil {
					errMsg += "os,"
				}

				srv.IP = record[headerIdx["host"]]
				err = Validate.Var(srv.IP, "ip|hostname")
				if err != nil {
					errMsg += "host,"
				}

				srv.Port, err = strconv.Atoi(record[headerIdx["port"]])
				if err != nil {
					errMsg += "port,"
				}

				srv.Owner = record[headerIdx["owner"]]
				err = Validate.Var(srv.Owner, "required,max=255")
				if err != nil {
					errMsg += "owner,"
				}

				srv.Path = record[headerIdx["path"]]
				err = Validate.Var(record[headerIdx["path"]], "max=255")
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
					pkg.Log(pkg.ERROR, fmt.Sprintf("Error on No.%d line %s, field validation on %s failed", i, record, errMsg))
				} else {
					srv.OSInfo = srv.ToSSHConfig().GetOSInfo()
					if _, err := srv.AddRow(); err != nil {
						errOccur = true
						pkg.Log(pkg.ERROR, fmt.Sprintf("Error on No.%d line %s, %s", i, record, err.Error()))
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Server) Toggle(gp *server.Goploy) server.Response {
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

func (Server) InstallAgent(gp *server.Goploy) server.Response {
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
			srv, err := (model.Server{ID: id}).GetData()
			if err != nil {
				pkg.Log(pkg.ERROR, fmt.Sprintf("Error on %d server, %s", id, err.Error()))
				return
			}
			client, err := srv.ToSSHConfig().Dial()
			if err != nil {
				pkg.Log(pkg.ERROR, fmt.Sprintf("Error on %d server, %s", id, err.Error()))
				return
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				pkg.Log(pkg.ERROR, fmt.Sprintf("Error on %d server, %s", id, err.Error()))
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
				pkg.Log(pkg.ERROR, fmt.Sprintf("Error on %d server, %s, detail: %s", id, err.Error(), sshErrbuf.String()))
				return
			}
			pkg.Log(pkg.INFO, sshErrbuf.String())
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
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
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

func (Server) DeleteFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		File     string `json:"file" validate:"required"`
		ServerID int64  `json:"serverId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Server) Report(gp *server.Goploy) server.Response {
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

func (s Server) DeleteMonitor(gp *server.Goploy) server.Response {
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

func (Server) EditProcess(gp *server.Goploy) server.Response {
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

func (Server) DeleteProcess(gp *server.Goploy) server.Response {
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

	if err := decodeJson(gp.Body, &reqData); err != nil {
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
