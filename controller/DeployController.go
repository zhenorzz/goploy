// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package controller

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/permission"
	"github.com/zhenorzz/goploy/repository"
	"github.com/zhenorzz/goploy/response"
	"github.com/zhenorzz/goploy/service"
	"github.com/zhenorzz/goploy/service/cmd"
	"github.com/zhenorzz/goploy/task"
	"github.com/zhenorzz/goploy/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Deploy Controller

func (d Deploy) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/deploy/getList", http.MethodGet, d.GetList).Permissions(permission.ShowDeployPage),
		core.NewRoute("/deploy/getPublishTrace", http.MethodGet, d.GetPublishTrace).Permissions(permission.DeployDetail),
		core.NewRoute("/deploy/getPublishTraceDetail", http.MethodGet, d.GetPublishTraceDetail).Permissions(permission.DeployDetail),
		core.NewRoute("/deploy/getPreview", http.MethodGet, d.GetPreview).Permissions(permission.DeployDetail),
		core.NewRoute("/deploy/review", http.MethodPut, d.Review).Permissions(permission.DeployReview),
		core.NewRoute("/deploy/resetState", http.MethodPut, d.ResetState).Permissions(permission.DeployResetState),
		core.NewRoute("/deploy/publish", http.MethodPost, d.Publish).Permissions(permission.DeployProject).Middleware(middleware.HasProjectPermission),
		core.NewRoute("/deploy/rebuild", http.MethodPost, d.Rebuild).Permissions(permission.DeployRollback).Middleware(middleware.HasProjectPermission),
		core.NewRoute("/deploy/greyPublish", http.MethodPost, d.GreyPublish).Permissions(permission.GreyDeploy).Middleware(middleware.HasProjectPermission),
		core.NewWhiteRoute("/deploy/webhook", http.MethodPost, d.Webhook).Middleware(middleware.FilterEvent),
		core.NewWhiteRoute("/deploy/callback", http.MethodGet, d.Callback),
		core.NewRoute("/deploy/fileCompare", http.MethodPost, d.FileCompare).Permissions(permission.FileCompare),
		core.NewRoute("/deploy/fileDiff", http.MethodPost, d.FileDiff).Permissions(permission.FileCompare),
		core.NewRoute("/deploy/manageProcess", http.MethodPost, d.ManageProcess).Permissions(permission.ProcessManager),
	}
}

func (Deploy) GetList(gp *core.Goploy) core.Response {
	var projects model.Projects
	var err error
	if _, ok := gp.Namespace.PermissionIDs[permission.GetAllDeployList]; ok {
		projects, err = model.Project{NamespaceID: gp.Namespace.ID}.GetDeployList()
	} else {
		projects, err = model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID}.GetDeployList()
	}
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			Project model.Projects `json:"list"`
		}{Project: projects},
	}
}

func (Deploy) GetPreview(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID  int64  `schema:"projectId" validate:"gt=0"`
		UserID     int64  `schema:"userId"`
		State      int    `schema:"state"`
		CommitDate string `schema:"commitDate"`
		DeployDate string `schema:"deployDate"`
		Branch     string `schema:"branch"`
		Commit     string `schema:"commit"`
		Filename   string `schema:"filename"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	commitDate := strings.Split(reqData.CommitDate, ",")
	for i, date := range commitDate {
		tm2, _ := time.Parse("2006-01-02 15:04:05", date)
		commitDate[i] = strconv.FormatInt(tm2.Unix(), 10)
	}
	gitTraceList, pagination, err := model.PublishTrace{
		ProjectID:    reqData.ProjectID,
		PublisherID:  reqData.UserID,
		PublishState: reqData.State,
	}.GetPreview(
		reqData.Branch,
		reqData.Commit,
		reqData.Filename,
		commitDate,
		strings.Split(reqData.DeployDate, ","),
		pagination,
	)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			GitTraceList model.PublishTraces `json:"list"`
			Pagination   model.Pagination    `json:"pagination"`
		}{GitTraceList: gitTraceList, Pagination: pagination},
	}
}

func (Deploy) GetPublishTrace(gp *core.Goploy) core.Response {
	lastPublishToken := gp.URLQuery.Get("lastPublishToken")
	publishTraceList, err := model.PublishTrace{Token: lastPublishToken}.GetListByToken()
	if err == sql.ErrNoRows {
		return response.JSON{Code: response.Error, Message: "No deploy record"}
	} else if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			PublishTraceList model.PublishTraces `json:"list"`
		}{PublishTraceList: publishTraceList},
	}
}

func (Deploy) GetPublishTraceDetail(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	detail, err := model.PublishTrace{ID: id}.GetDetail()
	if err == sql.ErrNoRows {
		return response.JSON{Code: response.Error, Message: "No deploy record"}
	} else if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Detail string `json:"detail"`
		}{Detail: detail},
	}
}

func (Deploy) ResetState(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64 `json:"projectId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Project{ID: reqData.ProjectID}).ResetState(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Deploy) FileCompare(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		FilePath  string `json:"filePath" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: reqData.ProjectID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srcPath := path.Join(core.GetProjectPath(reqData.ProjectID), reqData.FilePath)
	file, err := os.Open(srcPath)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer file.Close()
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	srcMD5 := hex.EncodeToString(hash.Sum(nil))
	projectServers, err := model.ProjectServer{ProjectID: reqData.ProjectID}.GetBindServerListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if len(projectServers) == 0 {
		return response.JSON{Code: response.Error, Message: "project have no server"}
	}

	type FileCompareData struct {
		ServerName string `json:"serverName"`
		ServerIP   string `json:"serverIP"`
		ServerID   int64  `json:"serverId"`
		Status     string `json:"status"`
		IsModified bool   `json:"isModified"`
	}
	var fileCompareList []FileCompareData
	ch := make(chan FileCompareData, len(projectServers))

	distPath := path.Join(project.Path, reqData.FilePath)
	for _, server := range projectServers {
		go func(server model.ProjectServer) {
			fileCompare := FileCompareData{server.ServerName, server.ServerIP, server.ServerID, "no change", false}
			client, err := server.ToSSHConfig().Dial()
			if err != nil {
				fileCompare.Status = "client error"
				ch <- fileCompare
				return
			}
			defer client.Close()

			//此时获取了sshClient，下面使用sshClient构建sftpClient
			sftpClient, err := sftp.NewClient(client)
			if err != nil {
				fileCompare.Status = "sftp error"
				ch <- fileCompare
				return
			}
			defer sftpClient.Close()
			file, err := sftpClient.Open(distPath)
			if err != nil {
				fileCompare.Status = "remote file not exists"
				ch <- fileCompare
				return
			}
			defer file.Close()
			hash := md5.New()
			_, _ = io.Copy(hash, file)
			distMD5 := hex.EncodeToString(hash.Sum(nil))
			if srcMD5 != distMD5 {
				fileCompare.Status = "modified"
				fileCompare.IsModified = true
				ch <- fileCompare
				return
			}
			ch <- fileCompare
		}(server)
	}

	for i := 0; i < len(projectServers); i++ {
		fileCompareList = append(fileCompareList, <-ch)
	}
	close(ch)
	return response.JSON{Data: fileCompareList}
}

func (Deploy) FileDiff(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		ServerID  int64  `json:"serverId" validate:"gt=0"`
		FilePath  string `json:"filePath" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: reqData.ProjectID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srcText, err := ioutil.ReadFile(path.Join(core.GetProjectPath(reqData.ProjectID), reqData.FilePath))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	server, err := model.Server{ID: reqData.ServerID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	client, err := server.ToSSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	distFile, err := sftpClient.Open(path.Join(project.Path, reqData.FilePath))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer distFile.Close()
	distText, err := ioutil.ReadAll(distFile)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{Data: struct {
		SrcText  string `json:"srcText"`
		DistText string `json:"distText"`
	}{SrcText: string(srcText), DistText: string(distText)}}
}

func (Deploy) ManageProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID         int64  `json:"serverId" validate:"gt=0"`
		ProjectProcessID int64  `json:"projectProcessId" validate:"gt=0"`
		Command          string `json:"command" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectProcess, err := model.ProjectProcess{ID: reqData.ProjectProcessID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	project, err := (model.Project{ID: projectProcess.ProjectID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	server, err := (model.Server{ID: reqData.ServerID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	script := ""
	switch reqData.Command {
	case "status":
		script = projectProcess.Status
	case "start":
		script = projectProcess.Start
	case "stop":
		script = projectProcess.Stop
	case "restart":
		script = projectProcess.Restart
	default:
		return response.JSON{Code: response.Error, Message: "Command error"}
	}
	if script == "" {
		return response.JSON{Code: response.Error, Message: "Command empty"}
	}

	script = service.ReplaceProjectVars(script, project)

	client, err := server.ToSSHConfig().Dial()
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
	err = session.Run(script)
	core.Log(core.TRACE, fmt.Sprintf("%s exec cmd %s, result %t, stdout: %s, stderr: %s", gp.UserInfo.Name, script, err == nil, sshOutbuf.String(), sshErrbuf.String()))
	return response.JSON{
		Data: struct {
			ExecRes bool   `json:"execRes"`
			Stdout  string `json:"stdout"`
			Stderr  string `json:"stderr"`
		}{ExecRes: err == nil, Stdout: sshOutbuf.String(), Stderr: sshErrbuf.String()},
	}
}

func (Deploy) Publish(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Commit    string `json:"commit"`
		Branch    string `json:"branch"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	project, err := model.Project{ID: reqData.ProjectID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if project.Review == model.Enable {
		err = projectReview(gp, project, reqData.Commit, reqData.Branch)
	} else {
		err = projectDeploy(gp, project, reqData.Commit, reqData.Branch)
	}
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Deploy) Rebuild(gp *core.Goploy) core.Response {
	type ReqData struct {
		Token string `json:"token"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	var err error
	publishTraceList, err := model.PublishTrace{Token: reqData.Token}.GetListByToken()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectID := publishTraceList[0].ProjectID
	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectServers, err := model.ProjectServer{ProjectID: projectID}.GetBindServerListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	needToPublish := project.SymlinkPath == ""
	var commitInfo repository.CommitInfo
	publishTraceServerCount := 0
	for _, publishTrace := range publishTraceList {
		// publish failed
		if publishTrace.State == 0 {
			needToPublish = true
			break
		}

		if publishTrace.Type == model.Pull {
			err := json.Unmarshal([]byte(publishTrace.Ext), &commitInfo)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		} else if publishTrace.Type == model.Deploy {
			for _, projectServer := range projectServers {
				if strings.Contains(publishTrace.Ext, projectServer.ServerIP) {
					publishTraceServerCount++
					break
				}
			}
		}
	}

	// project server has changed
	if publishTraceServerCount != len(projectServers) {
		needToPublish = true
	}
	if needToPublish == false {
		if len(project.AfterDeployScript) != 0 {
			scriptName := path.Join(core.GetProjectPath(project.ID), "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode))
			ioutil.WriteFile(scriptName, []byte(service.ReplaceProjectVars(project.AfterDeployScript, project)), 0755)
		}
		ch := make(chan bool, len(projectServers))
		for _, projectServer := range projectServers {
			go func(projectServer model.ProjectServer) {
				client, err := projectServer.ToSSHConfig().Dial()
				if err != nil {
					core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" dial err: "+err.Error())
					ch <- false
					return
				}
				defer client.Close()
				session, err := client.NewSession()
				if err != nil {
					core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" new session err: "+err.Error())
					ch <- false
					return
				}

				destDir := path.Join(project.SymlinkPath, project.LastPublishToken)
				// check if the path is existed or not
				if output, err := session.CombinedOutput("cd " + destDir); err != nil {
					core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" check symlink path err: "+err.Error()+", detail: "+string(output))
					ch <- false
					return
				}
				session, err = client.NewSession()
				if err != nil {
					core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" new session err: "+err.Error())
					ch <- false
					return
				}
				var afterDeployCommands []string
				cmdEntity := cmd.New("linux")
				afterDeployCommands = append(afterDeployCommands, cmdEntity.Symlink(destDir, project.Path))
				afterDeployCommands = append(afterDeployCommands, cmdEntity.ChangeDirTime(destDir))
				if len(project.AfterDeployScript) != 0 {
					scriptMode := "bash"
					if len(project.AfterDeployScriptMode) != 0 {
						scriptMode = project.AfterDeployScriptMode
					}
					afterDeployScriptPath := path.Join(project.Path, "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode))
					afterDeployCommands = append(afterDeployCommands, scriptMode+" "+cmdEntity.Path(afterDeployScriptPath))
					afterDeployCommands = append(afterDeployCommands, cmdEntity.Remove(afterDeployScriptPath))
				}
				// redirect to project path
				if output, err := session.CombinedOutput(strings.Join(afterDeployCommands, "&&")); err != nil {
					core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" symlink err: "+err.Error()+", detail: "+string(output))
					ch <- false
					return
				}
				ch <- true
			}(projectServer)
		}

		for i := 0; i < len(projectServers); i++ {
			if <-ch == false {
				needToPublish = true
				break
			}
		}
		close(ch)
		if needToPublish == false {
			model.PublishTrace{
				Token:      reqData.Token,
				UpdateTime: time.Now().Format("20060102150405"),
			}.EditUpdateTimeByToken()
			project.PublisherID = gp.UserInfo.ID
			project.PublisherName = gp.UserInfo.Name
			project.LastPublishToken = reqData.Token
			project.Publish()
			return response.JSON{Data: "symlink"}
		}
	}

	if needToPublish == true {
		repoEntity, _ := repository.GetRepo(project.RepoType)
		if !repoEntity.CanRollback() {
			return response.JSON{Code: response.Error, Message: fmt.Sprintf("plesae enable symlink to rollback the %s repo", project.RepoType)}
		}

		project.PublisherID = gp.UserInfo.ID
		project.PublisherName = gp.UserInfo.Name
		project.DeployState = model.ProjectDeploying
		project.LastPublishToken = uuid.New().String()
		err = project.Publish()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		task.AddDeployTask(service.Gsync{
			UserInfo:       gp.UserInfo,
			Project:        project,
			ProjectServers: projectServers,
			CommitID:       commitInfo.Commit,
			Branch:         commitInfo.Branch,
		})
	}
	return response.JSON{Data: "publish"}
}

func (Deploy) GreyPublish(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64   `json:"projectId" validate:"gt=0"`
		Commit    string  `json:"commit"`
		Branch    string  `json:"branch"`
		ServerIDs []int64 `json:"serverIds"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	project, err := model.Project{ID: reqData.ProjectID}.GetData()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	bindProjectServers, err := model.ProjectServer{ProjectID: project.ID}.GetBindServerListByProjectID()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectServers := model.ProjectServers{}

	for _, projectServer := range bindProjectServers {
		for _, serverID := range reqData.ServerIDs {
			if projectServer.ServerID == serverID {
				projectServers = append(projectServers, projectServer)
			}
		}
	}

	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = gp.UserInfo.Name
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	err = project.Publish()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	task.AddDeployTask(service.Gsync{
		UserInfo:       gp.UserInfo,
		Project:        project,
		ProjectServers: projectServers,
		CommitID:       reqData.Commit,
		Branch:         reqData.Branch,
	})

	return response.JSON{}
}

func (Deploy) Review(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectReviewID int64 `json:"projectReviewId" validate:"gt=0"`
		State           uint8 `json:"state" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectReviewModel := model.ProjectReview{
		ID:       reqData.ProjectReviewID,
		State:    reqData.State,
		Editor:   gp.UserInfo.Name,
		EditorID: gp.UserInfo.ID,
	}
	projectReview, err := projectReviewModel.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if projectReview.State != model.PENDING {
		return response.JSON{Code: response.Error, Message: "Project review state is invalid"}
	}

	if reqData.State == model.APPROVE {
		project, err := model.Project{ID: projectReview.ProjectID}.GetData()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		if err := projectDeploy(gp, project, projectReview.CommitID, projectReview.Branch); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}
	projectReviewModel.EditRow()

	return response.JSON{}
}

func (Deploy) Webhook(gp *core.Goploy) core.Response {
	projectID, err := strconv.ParseInt(gp.URLQuery.Get("project_id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	type ReqData struct {
		Ref string `json:"ref" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if project.State != model.Disable {
		return response.JSON{Code: response.Deny, Message: "Project is disabled"}
	}

	if project.AutoDeploy != model.ProjectWebhookDeploy {
		return response.JSON{Code: response.Deny, Message: "Webhook auto deploy turn off, go to project setting turn on"}
	}

	branch := ""
	if project.RepoType == model.RepoSVN {
		branch = reqData.Ref
	} else {
		branch = strings.Split(reqData.Ref, "/")[2]
	}

	if project.Branch != branch {
		return response.JSON{Code: response.Deny, Message: "Receive branch:" + branch + " push event, not equal to current branch"}
	}

	gp.UserInfo, err = model.User{ID: 1}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectServers, err := model.ProjectServer{ProjectID: project.ID}.GetBindServerListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = "webhook"
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	err = project.Publish()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	task.AddDeployTask(service.Gsync{
		UserInfo:       gp.UserInfo,
		Project:        project,
		ProjectServers: projectServers,
	})
	return response.JSON{Message: "receive push signal"}
}

func (Deploy) Callback(gp *core.Goploy) core.Response {
	projectReviewID, err := strconv.ParseInt(gp.URLQuery.Get("project_review_id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectReviewModel := model.ProjectReview{
		ID:       projectReviewID,
		State:    model.APPROVE,
		Editor:   "admin",
		EditorID: 1,
	}
	projectReview, err := projectReviewModel.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if projectReview.State != model.PENDING {
		return response.JSON{Code: response.Error, Message: "Project review state is invalid"}
	}

	project, err := model.Project{ID: projectReview.ProjectID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if err := projectDeploy(gp, project, projectReview.CommitID, projectReview.Branch); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectReviewModel.EditRow()

	return response.JSON{}
}

func projectDeploy(gp *core.Goploy, project model.Project, commitID string, branch string) error {
	projectServers, err := model.ProjectServer{ProjectID: project.ID}.GetBindServerListByProjectID()
	if err != nil {
		return err
	}
	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = gp.UserInfo.Name
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	err = project.Publish()
	if err != nil {
		return err
	}
	task.AddDeployTask(service.Gsync{
		UserInfo:       gp.UserInfo,
		Project:        project,
		ProjectServers: projectServers,
		CommitID:       commitID,
		Branch:         branch,
	})
	return nil
}

func projectReview(gp *core.Goploy, project model.Project, commitID string, branch string) error {
	if len(commitID) == 0 {
		return errors.New("commit id is required")
	}
	projectReviewModel := model.ProjectReview{
		ProjectID: project.ID,
		Branch:    branch,
		CommitID:  commitID,
		Creator:   gp.UserInfo.Name,
		CreatorID: gp.UserInfo.ID,
	}
	reviewURL := project.ReviewURL
	if len(reviewURL) > 0 {
		reviewURL = strings.Replace(reviewURL, "__PROJECT_ID__", strconv.FormatInt(project.ID, 10), 1)
		reviewURL = strings.Replace(reviewURL, "__PROJECT_NAME__", project.Name, 1)
		reviewURL = strings.Replace(reviewURL, "__BRANCH__", branch, 1)
		reviewURL = strings.Replace(reviewURL, "__ENVIRONMENT__", strconv.Itoa(int(project.Environment)), 1)
		reviewURL = strings.Replace(reviewURL, "__COMMIT_ID__", commitID, 1)
		reviewURL = strings.Replace(reviewURL, "__PUBLISH_TIME__", strconv.FormatInt(time.Now().Unix(), 10), 1)
		reviewURL = strings.Replace(reviewURL, "__PUBLISHER_ID__", gp.UserInfo.Name, 1)
		reviewURL = strings.Replace(reviewURL, "__PUBLISHER_NAME__", strconv.FormatInt(gp.UserInfo.ID, 10), 1)

		projectReviewModel.ReviewURL = reviewURL
	}
	id, err := projectReviewModel.AddRow()
	if err != nil {
		return err
	}
	if len(reviewURL) > 0 {
		callback := "http://"
		if gp.Request.TLS != nil {
			callback = "https://"
		}
		callback += gp.Request.Host + "/deploy/callback?project_review_id=" + strconv.FormatInt(id, 10)
		callback = url.QueryEscape(callback)
		reviewURL = strings.Replace(reviewURL, "__CALLBACK__", callback, 1)

		resp, err := http.Get(reviewURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	return nil
}
