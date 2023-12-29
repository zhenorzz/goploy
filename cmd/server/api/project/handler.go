// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package project

import (
	"bytes"
	"fmt"
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pipeline/docker"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/repo"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Project api.API

func (p Project) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/project/getList", http.MethodGet, p.GetList).Permissions(config.ShowProjectPage),
		server.NewRoute("/project/getLabelList", http.MethodGet, p.GetLabelList),
		server.NewRoute("/project/pingRepos", http.MethodGet, p.PingRepos),
		server.NewRoute("/project/getRemoteBranchList", http.MethodGet, p.GetRemoteBranchList),
		server.NewRoute("/project/getBindServerList", http.MethodGet, p.GetBindServerList),
		server.NewRoute("/project/getBindUserList", http.MethodGet, p.GetBindUserList),
		server.NewRoute("/project/getProjectFileList", http.MethodGet, p.GetProjectFileList).Permissions(config.FileSync),
		server.NewRoute("/project/getProjectFileContent", http.MethodGet, p.GetProjectFileContent).Permissions(config.FileSync),
		server.NewRoute("/project/getReposFileList", http.MethodGet, p.GetReposFileList).Permissions(config.FileCompare),
		server.NewRoute("/project/add", http.MethodPost, p.Add).Permissions(config.AddProject).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/edit", http.MethodPut, p.Edit).Permissions(config.EditProject).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/setAutoDeploy", http.MethodPut, p.SetAutoDeploy).Permissions(config.SwitchProjectWebhook).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/remove", http.MethodDelete, p.Remove).Permissions(config.DeleteProject).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/uploadFile", http.MethodPost, p.UploadFile).Permissions(config.FileSync).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/removeFile", http.MethodDelete, p.RemoveFile).Permissions(config.FileSync).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/addFile", http.MethodPost, p.AddFile).Permissions(config.FileSync).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/editFile", http.MethodPut, p.EditFile).Permissions(config.FileSync).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/addTask", http.MethodPost, p.AddTask).Permissions(config.DeployTask).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/removeTask", http.MethodDelete, p.RemoveTask).Permissions(config.DeployTask).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/getTaskList", http.MethodGet, p.GetTaskList).Permissions(config.DeployTask),
		server.NewRoute("/project/getReviewList", http.MethodGet, p.GetReviewList).Permissions(config.DeployReview),
		server.NewRoute("/project/getProcessList", http.MethodGet, p.GetProcessList).Permissions(config.ProcessManager),
		server.NewRoute("/project/addProcess", http.MethodPost, p.AddProcess).Permissions(config.ProcessManager).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/editProcess", http.MethodPut, p.EditProcess).Permissions(config.ProcessManager).LogFunc(middleware.AddOPLog),
		server.NewRoute("/project/deleteProcess", http.MethodDelete, p.DeleteProcess).Permissions(config.ProcessManager).LogFunc(middleware.AddOPLog),
	}
}

// GetLabelList lists all project's labels
// @Summary List all project's labels
// @Tags Project
// @Produce json
// @Success 200 {object} response.JSON{data=project.GetLabelList.RespData}
// @Router /project/getLabelList [get]
func (Project) GetLabelList(gp *server.Goploy) server.Response {
	var list []string
	var err error
	if _, ok := gp.Namespace.PermissionIDs[config.GetAllProjectList]; ok {
		list, err = model.Project{NamespaceID: gp.Namespace.ID}.GetLabelList()
	} else {
		list, err = model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID}.GetLabelList()
	}
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	type RespData struct {
		List []string `json:"list"`
	}
	return response.JSON{
		Data: RespData{List: list},
	}
}

// GetList lists all projects
// @Summary List all projects
// @Tags Project
// @Produce json
// @Success 200 {object} response.JSON{data=project.GetList.RespData}
// @Router /project/getList [get]
func (Project) GetList(gp *server.Goploy) server.Response {
	var projectList model.Projects
	var err error
	if _, ok := gp.Namespace.PermissionIDs[config.GetAllProjectList]; ok {
		projectList, err = model.Project{NamespaceID: gp.Namespace.ID}.GetList()
	} else {
		projectList, err = model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID}.GetList()
	}
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		Projects model.Projects `json:"list"`
	}
	return response.JSON{
		Data: RespData{Projects: projectList},
	}
}

// PingRepos ping the repository
// @Summary Ping the repository
// @Tags Project
// @Produce json
// @Param request query project.PingRepos.ReqData true "query params"
// @Success 200 {object} response.JSON
// @Router /project/pingRepos [get]
func (Project) PingRepos(gp *server.Goploy) server.Response {
	type ReqData struct {
		URL      string `schema:"url" validate:"required"`
		RepoType string `schema:"repoType" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

// GetRemoteBranchList lists all remote branches
// @Summary List all remote branches
// @Tags Project
// @Produce json
// @Param request query project.GetRemoteBranchList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=project.GetRemoteBranchList.RespData}
// @Router /project/getRemoteBranchList [get]
func (Project) GetRemoteBranchList(gp *server.Goploy) server.Response {
	type ReqData struct {
		URL      string `schema:"url" validate:"required"`
		RepoType string `schema:"repoType" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

	type RespData struct {
		Branch []string `json:"branch"`
	}
	return response.JSON{Data: RespData{Branch: list}}
}

// GetBindServerList lists all binding servers
// @Summary List all binding servers
// @Tags Project
// @Produce json
// @Param id path int true "project id"
// @Success 200 {object} response.JSON{data=project.GetBindServerList.RespData}
// @Router /project/getBindServerList [get]
func (Project) GetBindServerList(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectServers, err := model.ProjectServer{ProjectID: id}.GetBindServerListByProjectID()
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

// GetBindUserList lists all binding users
// @Summary List all binding users
// @Tags Project
// @Produce json
// @Param id path int true "project id"
// @Success 200 {object} response.JSON{data=project.GetBindUserList.RespData}
// @Router /project/getBindUserList [get]
func (Project) GetBindUserList(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectUsers, err := model.ProjectUser{ProjectID: id}.GetBindUserListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ProjectUsers model.ProjectUsers `json:"list"`
	}
	return response.JSON{
		Data: RespData{ProjectUsers: projectUsers},
	}
}

// GetProjectFileList lists all project files
// @Summary List all project files
// @Tags Project
// @Produce json
// @Param id path int true "project id"
// @Success 200 {object} response.JSON{data=project.GetProjectFileList.RespData}
// @Router /project/getProjectFileList [get]
func (Project) GetProjectFileList(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectFiles, err := model.ProjectFile{ProjectID: id}.GetListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ProjectFiles model.ProjectFiles `json:"list"`
	}
	return response.JSON{
		Data: RespData{ProjectFiles: projectFiles},
	}
}

// GetProjectFileContent shows the project file content
// @Summary Show the project file content
// @Tags Project
// @Produce json
// @Param id path int true "project id"
// @Success 200 {object} response.JSON{data=project.GetProjectFileContent.RespData}
// @Router /project/getProjectFileContent [get]
func (Project) GetProjectFileContent(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectFileData, err := model.ProjectFile{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	fileBytes, err := os.ReadFile(path.Join(config.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename))
	if err != nil {
		fmt.Println("read fail", err)
	}

	type RespData struct {
		Content string `json:"content"`
	}
	return response.JSON{
		Data: RespData{Content: string(fileBytes)},
	}
}

// GetReposFileList lists repository files
// @Summary List repository files
// @Tags Project
// @Produce json
// @Param request query project.GetReposFileList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=[]project.GetReposFileList.fileInfo}
// @Router /project/getReposFileList [get]
func (Project) GetReposFileList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID   int64  `schema:"id" validate:"required,gt=0"`
		Path string `schema:"path" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	files, err := os.ReadDir(path.Join(config.GetProjectPath(reqData.ID), reqData.Path))
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
	for _, file := range files {
		fileDetail, err := file.Info()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}

		fileList = append(fileList, fileInfo{
			Name:    file.Name(),
			Size:    fileDetail.Size(),
			Mode:    file.Type().String(),
			ModTime: fileDetail.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:   file.IsDir(),
		})
	}

	return response.JSON{Data: fileList}
}

// Add adds a project
// @Summary Add a project
// @Tags Project
// @Produce json
// @Param request body project.Add.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/add [post]
func (Project) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name                string              `json:"name" validate:"required"`
		RepoType            string              `json:"repoType" validate:"required"`
		URL                 string              `json:"url" validate:"required"`
		Label               string              `json:"label"`
		Path                string              `json:"path" validate:"required"`
		Environment         uint8               `json:"environment" validate:"required"`
		Branch              string              `json:"branch" validate:"required"`
		SymlinkPath         string              `json:"symlinkPath"`
		SymlinkBackupNumber uint8               `json:"symlinkBackupNumber"`
		Review              uint8               `json:"review"`
		ReviewURL           string              `json:"reviewURL"`
		Script              model.ProjectScript `json:"script"`
		TransferType        string              `json:"transferType"`
		TransferOption      string              `json:"transferOption"`
		DeployServerMode    string              `json:"deployServerMode"`
		ServerIDs           []int64             `json:"serverIds"`
		UserIDs             []int64             `json:"userIds"`
		NotifyType          uint8               `json:"notifyType"`
		NotifyTarget        string              `json:"notifyTarget"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if _, err := pkg.ParseCommandLine(reqData.TransferOption); err != nil {
		return response.JSON{Code: response.Error, Message: "Invalid transfer option format"}
	}

	if err := verifyProjectScript(reqData.Script); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectID, err := model.Project{
		NamespaceID:         gp.Namespace.ID,
		Name:                reqData.Name,
		RepoType:            reqData.RepoType,
		URL:                 reqData.URL,
		Label:               reqData.Label,
		Path:                reqData.Path,
		Environment:         reqData.Environment,
		Branch:              reqData.Branch,
		SymlinkPath:         reqData.SymlinkPath,
		SymlinkBackupNumber: reqData.SymlinkBackupNumber,
		Review:              reqData.Review,
		ReviewURL:           reqData.ReviewURL,
		Script:              reqData.Script,
		TransferType:        reqData.TransferType,
		TransferOption:      reqData.TransferOption,
		DeployServerMode:    reqData.DeployServerMode,
		NotifyType:          reqData.NotifyType,
		NotifyTarget:        reqData.NotifyTarget,
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

// Edit edits the project
// @Summary Edit the project
// @Tags Project
// @Produce json
// @Param request body project.Edit.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/edit [put]
func (Project) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID                  int64               `json:"id" validate:"required,gt=0"`
		Name                string              `json:"name"`
		RepoType            string              `json:"repoType"`
		URL                 string              `json:"url"`
		Label               string              `json:"label"`
		Path                string              `json:"path"`
		SymlinkPath         string              `json:"symlinkPath"`
		SymlinkBackupNumber uint8               `json:"symlinkBackupNumber"`
		Review              uint8               `json:"review"`
		ReviewURL           string              `json:"reviewURL"`
		Environment         uint8               `json:"environment"`
		Branch              string              `json:"branch"`
		ServerIDs           []int64             `json:"serverIds"`
		UserIDs             []int64             `json:"userIds"`
		Script              model.ProjectScript `json:"script"`
		TransferType        string              `json:"transferType"`
		TransferOption      string              `json:"transferOption"`
		DeployServerMode    string              `json:"deployServerMode"`
		NotifyType          uint8               `json:"notifyType"`
		NotifyTarget        string              `json:"notifyTarget"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if _, err := pkg.ParseCommandLine(reqData.TransferOption); err != nil {
		return response.JSON{Code: response.Error, Message: "Invalid option format"}
	}

	if err := verifyProjectScript(reqData.Script); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectData, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err = model.Project{
		ID:                  reqData.ID,
		Name:                reqData.Name,
		RepoType:            reqData.RepoType,
		URL:                 reqData.URL,
		Label:               reqData.Label,
		Path:                reqData.Path,
		Environment:         reqData.Environment,
		Branch:              reqData.Branch,
		SymlinkPath:         reqData.SymlinkPath,
		SymlinkBackupNumber: reqData.SymlinkBackupNumber,
		Review:              reqData.Review,
		ReviewURL:           reqData.ReviewURL,
		Script:              reqData.Script,
		TransferType:        reqData.TransferType,
		TransferOption:      reqData.TransferOption,
		DeployServerMode:    reqData.DeployServerMode,
		NotifyType:          reqData.NotifyType,
		NotifyTarget:        reqData.NotifyTarget,
	}.EditRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if reqData.URL != projectData.URL {
		srcPath := config.GetProjectPath(projectData.ID)
		_, err := os.Stat(srcPath)
		if err == nil || os.IsNotExist(err) == false {
			cmd := exec.Command("git", "remote", "set-url", "origin", reqData.URL)
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

// SetAutoDeploy updates project set auto_deploy
// @Summary Update project set auto_deploy
// @Tags Project
// @Produce json
// @Param request body project.SetAutoDeploy.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/setAutoDeploy [put]
func (Project) SetAutoDeploy(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID         int64 `json:"id" validate:"required,gt=0"`
		AutoDeploy uint8 `json:"autoDeploy" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

// Remove removes the project
// @Summary Remove the project
// @Tags Project
// @Produce json
// @Param request body project.Remove.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/remove [delete]
func (Project) Remove(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

// UploadFile uploads a file to the project
// @Summary Upload a file to the project
// @Tags Project
// @Produce json
// @Param request query project.UploadFile.ReqData true "query params"
// @Param file formData file true "file"
// @Success 200 {object} response.JSON{data=project.UploadFile.RespData}
// @Router /project/uploadFile [post]
func (Project) UploadFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ProjectFileID int64  `schema:"projectFileId" validate:"required,gte=0"`
		ProjectID     int64  `schema:"projectId"  validate:"required,gt=0"`
		Filename      string `schema:"filename"  validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

	// read all the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := os.WriteFile(filePath, fileBytes, 0755); err != nil {
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

	type RespData struct {
		ID int64 `json:"id"`
	}
	return response.JSON{
		Data: RespData{ID: reqData.ProjectFileID},
	}
}

// AddFile adds a file to the project
// @Summary Add a file to the project
// @Tags Project
// @Produce json
// @Param request body project.AddFile.ReqData true "body params"
// @Success 200 {object} response.JSON{data=project.AddFile.RespData}
// @Router /project/addFile [post]
func (Project) AddFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"required,gt=0"`
		Content   string `json:"content" validate:"required"`
		Filename  string `json:"filename" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer file.Close()

	_, err = file.WriteString(reqData.Content)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.ProjectFile{
		ProjectID: reqData.ProjectID,
		Filename:  reqData.Filename,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ID int64 `json:"id"`
	}
	return response.JSON{
		Data: RespData{ID: id},
	}
}

// EditFile edits the file to the project
// @Summary Edit the file to the project
// @Tags Project
// @Produce json
// @Param request body project.EditFile.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/editFile [put]
func (Project) EditFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"required,gt=0"`
		Content string `json:"content" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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
	_, err = file.WriteString(reqData.Content)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

// RemoveFile removes the file from the project
// @Summary Remove the file from the project
// @Tags Project
// @Produce json
// @Param request body project.RemoveFile.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/removeFile [delete]
func (Project) RemoveFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ProjectFileID int64 `json:"projectFileId" validate:"required,gt=0"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

// GetReviewList lists all reviews in the project
// @Summary List all reviews in the project
// @Tags Project
// @Produce json
// @Param request query project.GetReviewList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=project.GetReviewList.RespData}
// @Router /project/getReviewList [get]
func (Project) GetReviewList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `schema:"id" validate:"required,gt=0"`
		model.Pagination
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	ProjectReviews, pagination, err := model.ProjectReview{ProjectID: reqData.ID}.GetListByProjectID(reqData.Pagination)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ProjectReviews model.ProjectReviews `json:"list"`
		Pagination     model.Pagination     `json:"pagination"`
	}
	return response.JSON{
		Data: RespData{ProjectReviews: ProjectReviews, Pagination: pagination},
	}
}

// GetTaskList lists all tasks in the project
// @Summary List all tasks in the project
// @Tags Project
// @Produce json
// @Param request query project.GetTaskList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=project.GetTaskList.RespData}
// @Router /project/getTaskList [get]
func (Project) GetTaskList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `schema:"id" validate:"required,gt=0"`
		model.Pagination
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectTaskList, pagination, err := model.ProjectTask{ProjectID: reqData.ID}.GetListByProjectID(reqData.Pagination)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ProjectTasks model.ProjectTasks `json:"list"`
		Pagination   model.Pagination   `json:"pagination"`
	}
	return response.JSON{
		Data: RespData{ProjectTasks: projectTaskList, Pagination: pagination},
	}
}

// AddTask adds a task to the project
// @Summary Add a task to the project
// @Tags Project
// @Produce json
// @Param request body project.AddTask.ReqData true "body params"
// @Success 200 {object} response.JSON{data=project.AddTask.RespData}
// @Router /project/addTask [post]
func (Project) AddTask(gp *server.Goploy) server.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"required,gt=0"`
		Branch    string `json:"branch" validate:"required"`
		Commit    string `json:"commit" validate:"required"`
		Date      string `json:"date" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

	type RespData struct {
		ID int64 `json:"id"`
	}
	return response.JSON{
		Data: RespData{ID: id},
	}
}

// RemoveTask removes the task from the project
// @Summary Remove the task from the project
// @Tags Project
// @Produce json
// @Param request body project.RemoveTask.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/removeTask [delete]
func (Project) RemoveTask(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectTask{ID: reqData.ID}).RemoveRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

// GetProcessList lists all processes in the project
// @Summary List all processes in the project
// @Tags Project
// @Produce json
// @Param request query project.GetProcessList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=project.GetProcessList.RespData}
// @Router /project/getProcessList [get]
func (Project) GetProcessList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ProjectID int64 `schema:"projectId" validate:"required,gt=0"`
		model.Pagination
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := model.ProjectProcess{ProjectID: reqData.ProjectID}.GetListByProjectID(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		List model.ProjectProcesses `json:"list"`
	}
	return response.JSON{
		Data: RespData{List: list},
	}
}

// AddProcess adds a process to the project
// @Summary Add process to the project
// @Tags Project
// @Produce json
// @Param request body project.AddProcess.ReqData true "body params"
// @Success 200 {object} response.JSON{data=project.AddProcess.RespData}
// @Router /project/addProcess [post]
func (Project) AddProcess(gp *server.Goploy) server.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"required,gt=0"`
		Name      string `json:"name" validate:"required"`
		Status    string `json:"status"`
		Start     string `json:"start"`
		Stop      string `json:"stop"`
		Restart   string `json:"restart"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

	type RespData struct {
		ID int64 `json:"id"`
	}
	return response.JSON{
		Data: RespData{ID: id},
	}
}

// EditProcess edits the process in the project
// @Summary Edit the process in the project
// @Tags Project
// @Produce json
// @Param request body project.EditProcess.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/editProcess [post]
func (Project) EditProcess(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"required,gt=0"`
		Name    string `json:"name" validate:"required"`
		Status  string `json:"status"`
		Start   string `json:"start"`
		Stop    string `json:"stop"`
		Restart string `json:"restart"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

// DeleteProcess deletes the process from the project
// @Summary Delete the process from the project
// @Tags Project
// @Produce json
// @Param request body project.DeleteProcess.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /project/deleteProcess [delete]
func (Project) DeleteProcess(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectProcess{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func verifyProjectScript(projectScript model.ProjectScript) error {
	var dockerScript docker.Script
	errStr := "unmarshal %s yaml script fail, please check it"

	if projectScript.AfterPull.Mode == "yaml" && projectScript.AfterPull.Content != "" {
		if err := yaml.Unmarshal([]byte(projectScript.AfterPull.Content), &dockerScript); err != nil {
			return fmt.Errorf(errStr, "After Pull")
		}
	}

	if projectScript.AfterDeploy.Mode == "yaml" && projectScript.AfterDeploy.Content != "" {
		if err := yaml.Unmarshal([]byte(projectScript.AfterDeploy.Content), &dockerScript); err != nil {
			return fmt.Errorf(errStr, "After Deploy")
		}
	}

	if projectScript.DeployFinish.Mode == "yaml" && projectScript.DeployFinish.Content != "" {
		if err := yaml.Unmarshal([]byte(projectScript.DeployFinish.Content), &dockerScript); err != nil {
			return fmt.Errorf(errStr, "Deploy Finish")
		}
	}

	return nil
}
