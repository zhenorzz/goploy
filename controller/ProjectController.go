// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package controller

import (
	"bytes"
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/repo"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/permission"
	"github.com/zhenorzz/goploy/response"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Project Controller

func (p Project) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/project/getList", http.MethodGet, p.GetList).Permissions(permission.ShowProjectPage),
		core.NewRoute("/project/pingRepos", http.MethodGet, p.PingRepos),
		core.NewRoute("/project/getRemoteBranchList", http.MethodGet, p.GetRemoteBranchList),
		core.NewRoute("/project/getBindServerList", http.MethodGet, p.GetBindServerList),
		core.NewRoute("/project/getBindUserList", http.MethodGet, p.GetBindUserList),
		core.NewRoute("/project/getProjectFileList", http.MethodGet, p.GetProjectFileList).Permissions(permission.FileSync),
		core.NewRoute("/project/getProjectFileContent", http.MethodGet, p.GetProjectFileContent).Permissions(permission.FileSync),
		core.NewRoute("/project/getReposFileList", http.MethodGet, p.GetReposFileList).Permissions(permission.FileCompare),
		core.NewRoute("/project/add", http.MethodPost, p.Add).Permissions(permission.AddProject).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/edit", http.MethodPut, p.Edit).Permissions(permission.EditProject).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/setAutoDeploy", http.MethodPut, p.SetAutoDeploy).Permissions(permission.SwitchProjectWebhook).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/remove", http.MethodDelete, p.Remove).Permissions(permission.DeleteProject).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/uploadFile", http.MethodPost, p.UploadFile).Permissions(permission.FileSync).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/removeFile", http.MethodDelete, p.RemoveFile).Permissions(permission.FileSync).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/addFile", http.MethodPost, p.AddFile).Permissions(permission.FileSync).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/editFile", http.MethodPut, p.EditFile).Permissions(permission.FileSync).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/addTask", http.MethodPost, p.AddTask).Permissions(permission.DeployTask).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/removeTask", http.MethodDelete, p.RemoveTask).Permissions(permission.DeployTask).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/getTaskList", http.MethodGet, p.GetTaskList).Permissions(permission.DeployTask),
		core.NewRoute("/project/getReviewList", http.MethodGet, p.GetReviewList).Permissions(permission.DeployReview),
		core.NewRoute("/project/getProcessList", http.MethodGet, p.GetProcessList).Permissions(permission.ProcessManager),
		core.NewRoute("/project/addProcess", http.MethodPost, p.AddProcess).Permissions(permission.ProcessManager).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/editProcess", http.MethodPut, p.EditProcess).Permissions(permission.ProcessManager).LogFunc(middleware.AddOPLog),
		core.NewRoute("/project/deleteProcess", http.MethodDelete, p.DeleteProcess).Permissions(permission.ProcessManager).LogFunc(middleware.AddOPLog),
	}
}

func (Project) GetList(gp *core.Goploy) core.Response {
	var projectList model.Projects
	var err error
	if _, ok := gp.Namespace.PermissionIDs[permission.GetAllProjectList]; ok {
		projectList, err = model.Project{NamespaceID: gp.Namespace.ID}.GetList()
	} else {
		projectList, err = model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID}.GetList()
	}
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Projects model.Projects `json:"list"`
		}{Projects: projectList},
	}
}

func (Project) PingRepos(gp *core.Goploy) core.Response {
	type ReqData struct {
		URL      string `schema:"url" validate:"required"`
		RepoType string `schema:"repoType" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if strings.Contains(reqData.URL, "git@") {
		host := strings.Split(reqData.URL, "git@")[1]
		host = strings.Split(host, ":")[0]
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		knownHostsPath := homeDir + "/.ssh/known_hosts"
		var cmdOutbuf, cmdErrbuf bytes.Buffer
		cmd := exec.Command("ssh-keygen", "-F", host, "-f", knownHostsPath)
		cmd.Stdout = &cmdOutbuf
		cmd.Stderr = &cmdErrbuf
		if err := cmd.Run(); err != nil {
			cmdOutbuf.Reset()
			cmdErrbuf.Reset()
			cmd := exec.Command("ssh-keyscan", host)
			cmd.Stdout = &cmdOutbuf
			cmd.Stderr = &cmdErrbuf
			if err := cmd.Run(); err != nil {
				return response.JSON{Code: response.Error, Message: cmdErrbuf.String()}
			}
			f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
			defer f.Close()
			if _, err := f.Write(cmdOutbuf.Bytes()); err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		}
	}

	r, err := repo.GetRepo(reqData.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := r.Ping(reqData.URL); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) GetRemoteBranchList(gp *core.Goploy) core.Response {
	type ReqData struct {
		URL      string `schema:"url" validate:"required"`
		RepoType string `schema:"repoType" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if strings.Contains(reqData.URL, "git@") {
		host := strings.Split(reqData.URL, "git@")[1]
		host = strings.Split(host, ":")[0]
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		knownHostsPath := homeDir + "/.ssh/known_hosts"
		var cmdOutbuf, cmdErrbuf bytes.Buffer
		cmd := exec.Command("ssh-keygen", "-F", host, "-f", knownHostsPath)
		cmd.Stdout = &cmdOutbuf
		cmd.Stderr = &cmdErrbuf
		if err := cmd.Run(); err != nil {
			cmdOutbuf.Reset()
			cmdErrbuf.Reset()
			cmd := exec.Command("ssh-keyscan", host)
			cmd.Stdout = &cmdOutbuf
			cmd.Stderr = &cmdErrbuf
			if err := cmd.Run(); err != nil {
				return response.JSON{Code: response.Error, Message: cmdErrbuf.String()}
			}
			f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
			defer f.Close()
			if _, err := f.Write(cmdOutbuf.Bytes()); err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		}
	}

	r, err := repo.GetRepo(reqData.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := r.RemoteBranchList(reqData.URL)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{Data: struct {
		Branch []string `json:"branch"`
	}{Branch: list}}
}

func (Project) GetBindServerList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectServers, err := model.ProjectServer{ProjectID: id}.GetBindServerListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectServers model.ProjectServers `json:"list"`
		}{ProjectServers: projectServers},
	}
}

func (Project) GetBindUserList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectUsers, err := model.ProjectUser{ProjectID: id, NamespaceID: gp.Namespace.ID}.GetBindUserListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectUsers model.ProjectUsers `json:"list"`
		}{ProjectUsers: projectUsers},
	}
}

func (Project) GetProjectFileList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectFiles, err := model.ProjectFile{ProjectID: id}.GetListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectFiles model.ProjectFiles `json:"list"`
		}{ProjectFiles: projectFiles},
	}
}

func (Project) GetProjectFileContent(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectFileData, err := model.ProjectFile{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	fileBytes, err := ioutil.ReadFile(path.Join(config.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename))
	if err != nil {
		fmt.Println("read fail", err)
	}
	return response.JSON{
		Data: struct {
			Content string `json:"content"`
		}{Content: string(fileBytes)},
	}
}

func (Project) GetReposFileList(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID   int64  `schema:"id" validate:"gt=0"`
		Path string `schema:"path" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	files, err := ioutil.ReadDir(path.Join(config.GetProjectPath(reqData.ID), reqData.Path))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type fileInfo struct {
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Mode    string `json:"mode"`
		ModTime string `json:"modTime"`
		IsDir   bool   `json:"isDir"`
	}
	var fileList []fileInfo
	for _, f := range files {
		fileList = append(fileList, fileInfo{
			Name:    f.Name(),
			Size:    f.Size(),
			Mode:    f.Mode().String(),
			ModTime: f.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:   f.IsDir(),
		})
	}

	return response.JSON{Data: fileList}
}

func (Project) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name                  string  `json:"name" validate:"required"`
		RepoType              string  `json:"repoType" validate:"required"`
		URL                   string  `json:"url" validate:"required"`
		Path                  string  `json:"path" validate:"required"`
		Environment           uint8   `json:"environment" validate:"required"`
		Branch                string  `json:"branch" validate:"required"`
		SymlinkPath           string  `json:"symlinkPath"`
		SymlinkBackupNumber   uint8   `json:"symlinkBackupNumber"`
		Review                uint8   `json:"review"`
		ReviewURL             string  `json:"reviewURL"`
		AfterPullScriptMode   string  `json:"afterPullScriptMode"`
		AfterPullScript       string  `json:"afterPullScript"`
		AfterDeployScriptMode string  `json:"afterDeployScriptMode"`
		AfterDeployScript     string  `json:"afterDeployScript"`
		TransferType          string  `json:"transferType"`
		TransferOption        string  `json:"transferOption"`
		DeployServerMode      string  `json:"deployServerMode"`
		ServerIDs             []int64 `json:"serverIds"`
		UserIDs               []int64 `json:"userIds"`
		NotifyType            uint8   `json:"notifyType"`
		NotifyTarget          string  `json:"notifyTarget"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if _, err := pkg.ParseCommandLine(reqData.TransferOption); err != nil {
		return response.JSON{Code: response.Error, Message: "Invalid transfer option format"}
	}

	projectID, err := model.Project{
		NamespaceID:           gp.Namespace.ID,
		Name:                  reqData.Name,
		RepoType:              reqData.RepoType,
		URL:                   reqData.URL,
		Path:                  reqData.Path,
		Environment:           reqData.Environment,
		Branch:                reqData.Branch,
		SymlinkPath:           reqData.SymlinkPath,
		SymlinkBackupNumber:   reqData.SymlinkBackupNumber,
		Review:                reqData.Review,
		ReviewURL:             reqData.ReviewURL,
		AfterPullScriptMode:   reqData.AfterPullScriptMode,
		AfterPullScript:       reqData.AfterPullScript,
		AfterDeployScriptMode: reqData.AfterDeployScriptMode,
		AfterDeployScript:     reqData.AfterDeployScript,
		TransferType:          reqData.TransferType,
		TransferOption:        reqData.TransferOption,
		DeployServerMode:      reqData.DeployServerMode,
		NotifyType:            reqData.NotifyType,
		NotifyTarget:          reqData.NotifyTarget,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID: projectID,
			ServerID:  serverID,
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}
	if err := projectServersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectUsersModel := model.ProjectUsers{}
	for _, userID := range reqData.UserIDs {
		projectUserModel := model.ProjectUser{
			ProjectID: projectID,
			UserID:    userID,
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}
	if err := projectUsersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID                    int64   `json:"id" validate:"gt=0"`
		Name                  string  `json:"name"`
		RepoType              string  `json:"repoType"`
		URL                   string  `json:"url"`
		Path                  string  `json:"path"`
		SymlinkPath           string  `json:"symlinkPath"`
		SymlinkBackupNumber   uint8   `json:"symlinkBackupNumber"`
		Review                uint8   `json:"review"`
		ReviewURL             string  `json:"reviewURL"`
		Environment           uint8   `json:"environment"`
		Branch                string  `json:"branch"`
		ServerIDs             []int64 `json:"serverIds"`
		UserIDs               []int64 `json:"userIds"`
		AfterPullScriptMode   string  `json:"afterPullScriptMode"`
		AfterPullScript       string  `json:"afterPullScript"`
		AfterDeployScriptMode string  `json:"afterDeployScriptMode"`
		AfterDeployScript     string  `json:"afterDeployScript"`
		TransferType          string  `json:"transferType"`
		TransferOption        string  `json:"transferOption"`
		DeployServerMode      string  `json:"deployServerMode"`
		NotifyType            uint8   `json:"notifyType"`
		NotifyTarget          string  `json:"notifyTarget"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if _, err := pkg.ParseCommandLine(reqData.TransferOption); err != nil {
		return response.JSON{Code: response.Error, Message: "Invalid option format"}
	}

	projectData, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err = model.Project{
		ID:                    reqData.ID,
		Name:                  reqData.Name,
		RepoType:              reqData.RepoType,
		URL:                   reqData.URL,
		Path:                  reqData.Path,
		Environment:           reqData.Environment,
		Branch:                reqData.Branch,
		SymlinkPath:           reqData.SymlinkPath,
		SymlinkBackupNumber:   reqData.SymlinkBackupNumber,
		Review:                reqData.Review,
		ReviewURL:             reqData.ReviewURL,
		AfterPullScriptMode:   reqData.AfterPullScriptMode,
		AfterPullScript:       reqData.AfterPullScript,
		AfterDeployScriptMode: reqData.AfterDeployScriptMode,
		AfterDeployScript:     reqData.AfterDeployScript,
		TransferType:          reqData.TransferType,
		TransferOption:        reqData.TransferOption,
		DeployServerMode:      reqData.DeployServerMode,
		NotifyType:            reqData.NotifyType,
		NotifyTarget:          reqData.NotifyTarget,
	}.EditRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if reqData.URL != projectData.URL {
		srcPath := config.GetProjectPath(projectData.ID)
		_, err := os.Stat(srcPath)
		if err == nil || os.IsNotExist(err) == false {
			repo := reqData.URL
			cmd := exec.Command("git", "remote", "set-url", "origin", repo)
			cmd.Dir = srcPath
			if err := cmd.Run(); err != nil {
				return response.JSON{Code: response.Error, Message: "Project change url fail, you can do it manually, reason: " + err.Error()}
			}
		}
	}

	if err := (model.ProjectServer{ProjectID: projectData.ID}).DeleteByProjectID(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID: projectData.ID,
			ServerID:  serverID,
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}
	if err := projectServersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectUser{ProjectID: projectData.ID}).DeleteByProjectID(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectUsersModel := model.ProjectUsers{}
	for _, userID := range reqData.UserIDs {
		projectUserModel := model.ProjectUser{
			ProjectID: projectData.ID,
			UserID:    userID,
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}
	if err := projectUsersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) SetAutoDeploy(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID         int64 `json:"id" validate:"gt=0"`
		AutoDeploy uint8 `json:"autoDeploy" validate:"gte=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.Project{
		ID:         reqData.ID,
		AutoDeploy: reqData.AutoDeploy,
	}.SetAutoDeploy()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) Remove(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectData, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srcPath := config.GetProjectPath(projectData.ID)
	if err := os.RemoveAll(srcPath); err != nil {
		return response.JSON{Code: response.Error, Message: "Delete folder fail, Detail: " + err.Error()}
	}

	if err := (model.Project{ID: reqData.ID}).RemoveRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) UploadFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectFileID int64  `schema:"projectFileId" validate:"gte=0"`
		ProjectID     int64  `schema:"projectId"  validate:"gt=0"`
		Filename      string `schema:"filename"  validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	file, _, err := gp.Request.FormFile("file")
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer file.Close()

	filePath := path.Join(config.GetProjectFilePath(reqData.ProjectID), reqData.Filename)
	fileDir := path.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(fileDir, 0755)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		} else {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := ioutil.WriteFile(filePath, fileBytes, 0755); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if reqData.ProjectFileID == 0 {
		reqData.ProjectFileID, err = model.ProjectFile{
			Filename:  reqData.Filename,
			ProjectID: reqData.ProjectID,
		}.AddRow()
	} else {
		err = model.ProjectFile{
			ID:        reqData.ProjectFileID,
			Filename:  reqData.Filename,
			ProjectID: reqData.ProjectID,
		}.EditRow()
	}

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: reqData.ProjectFileID},
	}
}

func (Project) AddFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Content   string `json:"content" validate:"required"`
		Filename  string `json:"filename" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	filePath := path.Join(config.GetProjectFilePath(reqData.ProjectID), reqData.Filename)
	fileDir := path.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(fileDir, 0755)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		} else {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(reqData.Content)

	id, err := model.ProjectFile{
		ProjectID: reqData.ProjectID,
		Filename:  reqData.Filename,
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

func (Project) EditFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"gt=0"`
		Content string `json:"content" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectFileData, err := model.ProjectFile{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	_, err = os.Stat(config.GetProjectFilePath(projectFileData.ProjectID))
	if err != nil {
		err := os.MkdirAll(config.GetProjectFilePath(projectFileData.ProjectID), os.ModePerm)
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	file, err := os.Create(path.Join(config.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(reqData.Content)

	return response.JSON{}
}

func (Project) RemoveFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectFileID int64 `json:"projectFileId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectFileData, err := model.ProjectFile{ID: reqData.ProjectFileID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := os.Remove(path.Join(config.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename)); err != nil {
		if !os.IsNotExist(err) {
			return response.JSON{Code: response.Error, Message: "Delete file fail, Detail: " + err.Error()}
		}
	}

	if err := (model.ProjectFile{ID: reqData.ProjectFileID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) GetReviewList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	ProjectReviews, pagination, err := model.ProjectReview{ProjectID: id}.GetListByProjectID(pagination)

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectReviews model.ProjectReviews `json:"list"`
			Pagination     model.Pagination     `json:"pagination"`
		}{ProjectReviews: ProjectReviews, Pagination: pagination},
	}
}

func (Project) GetTaskList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectTaskList, pagination, err := model.ProjectTask{ProjectID: id}.GetListByProjectID(pagination)

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectTasks model.ProjectTasks `json:"list"`
			Pagination   model.Pagination   `json:"pagination"`
		}{ProjectTasks: projectTaskList, Pagination: pagination},
	}
}

func (Project) AddTask(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Branch    string `json:"branch" validate:"required"`
		Commit    string `json:"commit" validate:"required"`
		Date      string `json:"date" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.ProjectTask{
		ProjectID: reqData.ProjectID,
		CommitID:  reqData.Commit,
		Branch:    reqData.Branch,
		Date:      reqData.Date,
		Creator:   gp.UserInfo.Name,
		CreatorID: gp.UserInfo.ID,
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

func (Project) RemoveTask(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectTask{ID: reqData.ID}).RemoveRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) GetProcessList(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `schema:"projectId" validate:"gt=0"`
		Page      uint64 `schema:"page" validate:"gt=0"`
		Rows      uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := model.ProjectProcess{ProjectID: reqData.ProjectID}.GetListByProjectID(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.ProjectProcesses `json:"list"`
		}{List: list},
	}
}

func (Project) AddProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Name      string `json:"name" validate:"required"`
		Status    string `json:"status"`
		Start     string `json:"start"`
		Stop      string `json:"stop"`
		Restart   string `json:"restart"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.ProjectProcess{
		ProjectID: reqData.ProjectID,
		Name:      reqData.Name,
		Status:    reqData.Status,
		Start:     reqData.Start,
		Stop:      reqData.Stop,
		Restart:   reqData.Restart,
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

func (Project) EditProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"gt=0"`
		Name    string `json:"name" validate:"required"`
		Status  string `json:"status"`
		Start   string `json:"start"`
		Stop    string `json:"stop"`
		Restart string `json:"restart"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.ProjectProcess{
		ID:      reqData.ID,
		Name:    reqData.Name,
		Status:  reqData.Status,
		Start:   reqData.Start,
		Stop:    reqData.Stop,
		Restart: reqData.Restart,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) DeleteProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectProcess{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
