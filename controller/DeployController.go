package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"goploy/core"
	"goploy/model"
	"goploy/utils"
	"goploy/ws"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

// Deploy struct
type Deploy Controller

// GetList deploy list
func (deploy Deploy) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		Project model.Projects `json:"projectList"`
	}
	groupIDStr := gp.URLQuery.Get("groupId")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: "groupId参数错误"}
		response.JSON(w)
		return
	}

	projectName := gp.URLQuery.Get("projectName")

	var projects model.Projects
	projects, err = model.Project{
		GroupID: groupID,
		Name:    projectName,
	}.GetUserProjectList(gp.UserInfo.ID, gp.UserInfo.Role, gp.UserInfo.ManageGroupStr)

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	response := core.Response{Data: RespData{Project: projects}}
	response.JSON(w)
}

// GetPreview deploy detail
func (deploy Deploy) GetPreview(w http.ResponseWriter, gp *core.Goploy) {

	type RespData struct {
		GitTraceList model.PublishTraces `json:"gitTraceList"`
	}

	projectID, err := strconv.ParseInt(gp.URLQuery.Get("projectId"), 10, 64)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: "id参数错误"}
		response.JSON(w)
		return
	}

	gitTraceList, err := model.PublishTrace{ProjectID: projectID}.GetPreviewByProjectID()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{GitTraceList: gitTraceList}}
	response.JSON(w)
}

// GetDetail deploy detail
func (deploy Deploy) GetDetail(w http.ResponseWriter, gp *core.Goploy) {

	type RespData struct {
		PublishTraceList model.PublishTraces `json:"publishTraceList"`
	}

	lastPublishToken := gp.URLQuery.Get("lastPublishToken")

	publishTraceList, err := model.PublishTrace{Token: lastPublishToken}.GetListByToken()
	if err == sql.ErrNoRows {
		response := core.Response{Code: core.Deny, Message: "项目尚无构建记录"}
		response.JSON(w)
		return
	} else if err != nil {
		response := core.Response{Code: core.Deny, Message: "GitTrace.GetListByToken失败"}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{PublishTraceList: publishTraceList}}
	response.JSON(w)
}

// GetCommitList get latest 10 commit list
func (deploy Deploy) GetCommitList(w http.ResponseWriter, gp *core.Goploy) {

	type RespData struct {
		CommitList []string `json:"commitList"`
	}

	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: "id参数错误"}
		response.JSON(w)
		return
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	srcPath := core.RepositoryPath + project.Name

	log := exec.Command("git", "log", "--pretty=format:%H`%an`%ar`%s", "-n", "10")
	log.Dir = srcPath
	var logOutbuf, logErrbuf bytes.Buffer
	log.Stdout = &logOutbuf
	log.Stderr = &logErrbuf
	if err := log.Run(); err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{CommitList: strings.Split(logOutbuf.String(), "\n")}}
	response.JSON(w)
}

// Publish the project
func (deploy Deploy) Publish(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ProjectID int64 `json:"projectId"`
	}
	var reqData ReqData
	if err := json.Unmarshal(gp.Body, &reqData); err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	project, err := model.Project{
		ID: reqData.ProjectID,
	}.GetData()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	if project.DeployState == model.ProjectDeploying {
		response := core.Response{Code: core.Deny, Message: "项目正在构建，请稍后再试"}
		response.JSON(w)
		return
	}

	projectServers, err := model.ProjectServer{ProjectID: reqData.ProjectID}.GetBindServerListByProjectID()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = gp.UserInfo.Name
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	project.UpdateTime = time.Now().Unix()
	err = project.Publish()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "部署中，请稍后"}
	response.JSON(w)

	go execSync(gp.UserInfo, project, projectServers)
	return
}

// Rollback the project
func (deploy Deploy) Rollback(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ProjectID int64  `json:"projectId"`
		Commit    string `json:"commit"`
	}
	var reqData ReqData
	if err := json.Unmarshal(gp.Body, &reqData); err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	project, err := model.Project{
		ID: reqData.ProjectID,
	}.GetData()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	if project.DeployState == model.ProjectDeploying {
		response := core.Response{Code: core.Deny, Message: "项目正在构建，请稍后再试"}
		response.JSON(w)
		return
	}

	projectServers, err := model.ProjectServer{ProjectID: reqData.ProjectID}.GetBindServerListByProjectID()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = gp.UserInfo.Name
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	project.UpdateTime = time.Now().Unix()
	err = project.Publish()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "重新构建中，请稍后"}
	response.JSON(w)
	go execRollback(gp.UserInfo, reqData.Commit, project, projectServers)
	return
}

type SyncMessage struct {
	serverName string
	ProjectID  int64
	Detail     string
	State      int
}

func execSync(userInfo model.User, project model.Project, projectServers model.ProjectServers) {
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" deploy start")
	ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
		Type: ws.TypeProject,
		Message: ws.ProjectMessage{
			ProjectID:   project.ID,
			ProjectName: project.Name,
			UserID:      userInfo.ID,
			Username:    userInfo.Name,
			State:       model.ProjectDeploying,
		},
	}
	publishTraceModel := model.PublishTrace{
		Token:         project.LastPublishToken,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   userInfo.ID,
		PublisherName: userInfo.Name,
		Type:          model.Pull,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	gitPullMessage, gitCommitID, err := gitSync(project)
	if err != nil {
		project.DeployFail()
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
			Type: ws.TypeProject,
			Message: ws.ProjectMessage{
				ProjectID:   project.ID,
				ProjectName: project.Name,
				UserID:      userInfo.ID,
				Username:    userInfo.Name,
				State:       model.ProjectFail,
				Message:     err.Error(),
			},
		}
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
		go notify(project, model.ProjectFail, err.Error())
		return
	} else {
		ext, _ := json.Marshal(struct {
			Commit string `json:"commit"`
		}{gitCommitID})
		publishTraceModel.Ext = string(ext)
		publishTraceModel.Detail = gitPullMessage
		publishTraceModel.State = model.Success
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
	}
	if project.AfterPullScript != "" {
		outputString, err := runAfterPullScript(project)
		publishTraceModel.Type = model.AfterPull
		ext, _ := json.Marshal(struct {
			Script string `json:"script"`
		}{project.AfterPullScript})
		publishTraceModel.Ext = string(ext)
		if err != nil {
			project.DeployFail()
			publishTraceModel.Detail = err.Error()
			publishTraceModel.State = model.Fail
			ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
				Type: ws.TypeProject,
				Message: ws.ProjectMessage{
					ProjectID:   project.ID,
					ProjectName: project.Name,
					UserID:      userInfo.ID,
					State:       model.ProjectFail,
					Message:     err.Error(),
				},
			}
			if _, err := publishTraceModel.AddRow(); err != nil {
				core.Log(core.ERROR, err.Error())
			}
			go notify(project, model.ProjectFail, err.Error())
			return
		}
		publishTraceModel.Detail = outputString
		publishTraceModel.State = model.Success
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
	}

	ch := make(chan SyncMessage, len(projectServers))
	for _, projectServer := range projectServers {
		go remoteSync(ch, userInfo, project, projectServer)
	}

	message := ""
	for i := 0; i < len(projectServers); i++ {
		syncMessage := <-ch
		if syncMessage.State == model.ProjectFail {
			message += syncMessage.serverName + " error message: " + syncMessage.Detail
		}
	}
	state := model.ProjectFail
	if message == "" {
		project.DeploySuccess()
		state = model.ProjectSuccess
		core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" deploy success")

	} else {
		project.DeployFail()
		core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" deploy fail")

	}

	ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
		Type: ws.TypeProject,
		Message: ws.ProjectMessage{
			ProjectID:   project.ID,
			ProjectName: project.Name,
			UserID:      userInfo.ID,
			Username:    userInfo.Name,
			State:       uint8(state),
			Message:     message,
		},
	}
	go notify(project, state, message)
	return
}

func execRollback(userInfo model.User, commit string, project model.Project, projectServers model.ProjectServers) {
	ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
		Type: ws.TypeProject,
		Message: ws.ProjectMessage{
			ProjectID:   project.ID,
			ProjectName: project.Name,
			UserID:      userInfo.ID,
			Username:    userInfo.Name,
			State:       model.ProjectDeploying,
		},
	}
	publishTraceModel := model.PublishTrace{
		Token:         project.LastPublishToken,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   userInfo.ID,
		PublisherName: userInfo.Name,
		Type:          model.Pull,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}
	stdout, err := gitRollback(commit, project)
	if err != nil {
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
			Type: ws.TypeProject,
			Message: ws.ProjectMessage{
				ProjectID:   project.ID,
				ProjectName: project.Name,
				UserID:      userInfo.ID,
				Username:    userInfo.Name,
				State:       model.ProjectFail,
				Message:     err.Error(),
			},
		}
		go notify(project, model.ProjectFail, err.Error())
		return
	}
	ext, _ := json.Marshal(struct {
		Commit string `json:"commit"`
	}{commit})
	publishTraceModel.Ext = string(ext)
	publishTraceModel.Detail = stdout
	publishTraceModel.State = model.Success

	if _, err := publishTraceModel.AddRow(); err != nil {
		core.Log(core.ERROR, err.Error())
	}

	if project.AfterPullScript != "" {
		outputString, err := runAfterPullScript(project)
		publishTraceModel.Type = model.AfterPull
		ext, _ := json.Marshal(struct {
			Script string `json:"script"`
		}{project.AfterPullScript})
		publishTraceModel.Ext = string(ext)
		if err != nil {
			project.DeployFail()
			publishTraceModel.Detail = err.Error()
			publishTraceModel.State = model.Fail
			ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
				Type: ws.TypeProject,
				Message: ws.ProjectMessage{
					ProjectID:   project.ID,
					ProjectName: project.Name,
					UserID:      userInfo.ID,
					Username:    userInfo.Name,
					State:       model.ProjectFail,
					Message:     err.Error(),
				},
			}
			if _, err := publishTraceModel.AddRow(); err != nil {
				core.Log(core.ERROR, err.Error())
			}
			go notify(project, model.ProjectFail, err.Error())
			return
		}
		publishTraceModel.Detail = outputString
		publishTraceModel.State = model.Success
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
	}

	ch := make(chan SyncMessage, len(projectServers))
	for _, projectServer := range projectServers {
		go remoteSync(ch, userInfo, project, projectServer)
	}

	message := ""
	for i := 0; i < len(projectServers); i++ {
		syncMessage := <-ch
		if syncMessage.State == model.ProjectFail {
			message += syncMessage.serverName + " error message: " + syncMessage.Detail
		}
	}
	state := model.ProjectFail
	if message == "" {
		project.DeploySuccess()
		state = model.ProjectSuccess
	} else {
		project.DeployFail()
	}
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" rollback success")
	ws.GetBroadcastHub().BroadcastData <- &ws.BroadcastData{
		Type: ws.TypeProject,
		Message: ws.ProjectMessage{
			ProjectID:   project.ID,
			ProjectName: project.Name,
			UserID:      userInfo.ID,
			Username:    userInfo.Name,
			State:       uint8(state),
			Message:     message,
		},
	}
	go notify(project, state, message)
	return
}
func gitSync(project model.Project) (string, string, error) {
	if err := gitCreate(project); err != nil {
		return "", "", err
	}
	stdout, err := gitPull(project)
	if err != nil {
		return "", "", err
	}

	commit, err := gitCommitID(project)
	if err != nil {
		return "", "", err
	}
	return stdout, commit, err
}

func gitCreate(project model.Project) error {
	srcPath := core.RepositoryPath + project.Name
	if _, err := os.Stat(srcPath); err != nil {
		if err := os.RemoveAll(srcPath); err != nil {
			return err
		}
		repo := project.URL
		cmd := exec.Command("git", "clone", repo, srcPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化 git clone")
		if err := cmd.Run(); err != nil {
			core.Log(core.ERROR, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化失败:"+err.Error())
			return errors.New("项目初始化失败")
		}

		if project.Branch != "master" {
			checkout := exec.Command("git", "checkout", "-b", project.Branch, "origin/"+project.Branch)
			checkout.Dir = srcPath
			var checkoutOutbuf, checkoutErrbuf bytes.Buffer
			checkout.Stdout = &checkoutOutbuf
			checkout.Stderr = &checkoutErrbuf
			if err := checkout.Run(); err != nil {
				core.Log(core.ERROR, checkoutErrbuf.String())
				os.RemoveAll(srcPath)
				return errors.New(checkoutErrbuf.String())
			}
		}
		core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化成功")
	}
	return nil
}

func gitPull(project model.Project) (string, error) {
	srcPath := core.RepositoryPath + project.Name

	clean := exec.Command("git", "clean", "-f")
	clean.Dir = srcPath
	var cleanOutbuf, cleanErrbuf bytes.Buffer
	clean.Stdout = &cleanOutbuf
	clean.Stderr = &cleanErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git clean -f")
	if err := clean.Run(); err != nil {
		core.Log(core.ERROR, cleanErrbuf.String())
		return "", errors.New(cleanErrbuf.String())
	}
	pull := exec.Command("git", "pull")
	pull.Dir = srcPath
	var pullOutbuf, pullErrbuf bytes.Buffer
	pull.Stdout = &pullOutbuf
	pull.Stderr = &pullErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git pull")
	if err := pull.Run(); err != nil {
		core.Log(core.ERROR, pullErrbuf.String())
		return "", errors.New(pullErrbuf.String())
	}
	core.Log(core.TRACE, pullOutbuf.String())
	return pullOutbuf.String(), nil
}

func gitRollback(commit string, project model.Project) (string, error) {
	srcPath := core.RepositoryPath + project.Name

	resetCmd := exec.Command("git", "reset", "--hard", commit)
	resetCmd.Dir = srcPath
	var resetOutbuf, resetErrbuf bytes.Buffer
	resetCmd.Stdout = &resetOutbuf
	resetCmd.Stderr = &resetErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git reset --hard "+commit)
	if err := resetCmd.Run(); err != nil {
		core.Log(core.ERROR, resetErrbuf.String())
		return "", errors.New(resetErrbuf.String())
	}

	core.Log(core.TRACE, resetOutbuf.String())
	return resetOutbuf.String(), nil
}

func gitCommitID(project model.Project) (string, error) {
	srcPath := core.RepositoryPath + project.Name

	git := exec.Command("git", "rev-parse", "HEAD")
	git.Dir = srcPath
	var gitOutbuf, gitErrbuf bytes.Buffer
	git.Stdout = &gitOutbuf
	git.Stderr = &gitErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git rev-parse HEAD")
	if err := git.Run(); err != nil {
		core.Log(core.ERROR, gitErrbuf.String())
		return "", errors.New(gitErrbuf.String())
	}
	return gitOutbuf.String(), nil
}

func runAfterPullScript(project model.Project) (string, error) {
	srcPath := core.RepositoryPath + project.Name
	scriptName := srcPath + "/after-pull.sh"
	ioutil.WriteFile(scriptName, []byte(project.AfterPullScript), 0755)
	handler := exec.Command("bash", "./after-pull.sh")
	handler.Dir = srcPath
	var outbuf, errbuf bytes.Buffer
	handler.Stdout = &outbuf
	handler.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+project.AfterPullScript)
	if err := handler.Run(); err != nil {
		core.Log(core.ERROR, err.Error())
		return "", err
	}

	os.Remove(scriptName)
	return outbuf.String(), nil
}

func remoteSync(chInput chan<- SyncMessage, userInfo model.User, project model.Project, projectServer model.ProjectServer) {
	srcPath := core.RepositoryPath + project.Name + "/"
	remoteMachine := projectServer.ServerOwner + "@" + projectServer.ServerIP
	destPath := remoteMachine + ":" + project.Path
	publishTraceModel := model.PublishTrace{
		Token:         project.LastPublishToken,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   userInfo.ID,
		PublisherName: userInfo.Name,
		Type:          model.Deploy,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	ext, _ := json.Marshal(struct {
		ServerID   int64  `json:"serverId"`
		ServerName string `json:"serverName"`
	}{projectServer.ServerID, projectServer.ServerName})
	publishTraceModel.Ext = string(ext)

	rsyncOption, _ := utils.ParseCommandLine(project.RsyncOption)
	rsyncOption = append(rsyncOption, "-e", "ssh -p "+strconv.Itoa(int(projectServer.ServerPort))+" -o StrictHostKeyChecking=no", srcPath, destPath)
	cmd := exec.Command("rsync", rsyncOption...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" rsync "+strings.Join(rsyncOption, " "))
	var rsyncError error
	// 失败重试三次
	for attempt := 0; attempt < 3; attempt++ {
		rsyncError = cmd.Run()
		if rsyncError != nil {
			core.Log(core.ERROR, errbuf.String())
		} else {
			ext, _ := json.Marshal(struct {
				ServerID   int64  `json:"serverId"`
				ServerName string `json:"serverName"`
				Command    string `json:"command"`
			}{projectServer.ServerID, projectServer.ServerName, "rsync " + strings.Join(rsyncOption, " ")})
			publishTraceModel.Ext = string(ext)
			publishTraceModel.Detail = outbuf.String()
			publishTraceModel.State = model.Success
			publishTraceModel.AddRow()
			break
		}
	}

	if rsyncError != nil {
		publishTraceModel.Detail = errbuf.String()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		chInput <- SyncMessage{
			serverName: projectServer.ServerName,
			ProjectID:  project.ID,
			Detail:     errbuf.String(),
			State:      model.ProjectFail,
		}
		return
	}
	// 没有脚本就不运行
	if project.AfterDeployScript == "" {
		chInput <- SyncMessage{
			serverName: projectServer.ServerName,
			ProjectID:  project.ID,
			State:      model.ProjectSuccess,
		}
		return
	}
	publishTraceModel.Type = model.AfterDeploy
	ext, _ = json.Marshal(struct {
		ServerID   int64  `json:"serverId"`
		ServerName string `json:"serverName"`
		Script     string `json:"script"`
	}{projectServer.ServerID, projectServer.ServerName, project.AfterDeployScript})
	publishTraceModel.Ext = string(ext)
	// 执行ssh脚本
	var session *ssh.Session
	var connectError error
	var scriptError error
	for attempt := 0; attempt < 3; attempt++ {
		session, connectError = utils.ConnectSSH(projectServer.ServerOwner, "", projectServer.ServerIP, int(projectServer.ServerPort))
		if connectError != nil {
			core.Log(core.ERROR, connectError.Error())
		} else {
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			sshOutbuf.Reset()
			afterDeployScript := "echo '" + project.AfterDeployScript + "' > /tmp/after-deploy.sh;bash /tmp/after-deploy.sh"
			if scriptError = session.Run(afterDeployScript); scriptError != nil {
				core.Log(core.ERROR, scriptError.Error())
			} else {
				publishTraceModel.Detail = sshOutbuf.String()
				publishTraceModel.State = model.Success
				publishTraceModel.AddRow()
				break
			}
		}

	}
	if session != nil {
		defer session.Close()
	}
	if connectError != nil {
		publishTraceModel.Detail = connectError.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		chInput <- SyncMessage{
			serverName: projectServer.ServerName,
			ProjectID:  project.ID,
			Detail:     connectError.Error(),
			State:      model.ProjectFail,
		}
		return
	} else if scriptError != nil {
		publishTraceModel.Detail = scriptError.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		chInput <- SyncMessage{
			serverName: projectServer.ServerName,
			ProjectID:  project.ID,
			Detail:     scriptError.Error(),
			State:      model.ProjectFail,
		}
		return
	}
	chInput <- SyncMessage{
		serverName: projectServer.ServerName,
		ProjectID:  project.ID,
		State:      model.ProjectSuccess,
	}
	return
}

func notify(project model.Project, deployState int, detail string) {
	if project.NotifyType == 0 {
		return
	} else if project.NotifyType == 1 {
		type markdown struct {
			Content string `json:"content"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		content := "构建项目<font color=\"warning\">" + project.Name + "</font>，请相关同事注意。\n "

		if deployState == model.ProjectFail {
			content += ">状态:<font color=\"red\"> 失败</font> \n "
			content += ">详情：<font color=\"comment\">" + detail + "</font>"
		} else {
			content += ">状态:<font color=\"green\"> 成功</font>"
		}

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Content: content,
			},
		}
		b, _ := json.Marshal(msg)
		_, err := http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
		if err != nil {
			core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10) + " " +err.Error())
		}
	}
}
