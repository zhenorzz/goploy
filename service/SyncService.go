// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/repository"
	"github.com/zhenorzz/goploy/service/cmd"
	"github.com/zhenorzz/goploy/service/transmitter"
	"github.com/zhenorzz/goploy/utils"
	"github.com/zhenorzz/goploy/ws"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Gsync struct {
	UserInfo       model.User
	Project        model.Project
	ProjectServers model.ProjectServers
	CommitInfo     repository.CommitInfo
	CommitID       string
	Branch         string
}

type syncMessage struct {
	serverName string
	projectID  int64
	detail     string
	state      int
}

var projectLogFormat = "projectID: %d %s"

func (gsync Gsync) Exec() {
	core.Logf(core.TRACE, "projectID: %d deploy start", gsync.Project.ID)
	publishTraceModel := model.PublishTrace{
		Token:         gsync.Project.LastPublishToken,
		ProjectID:     gsync.Project.ID,
		ProjectName:   gsync.Project.Name,
		PublisherID:   gsync.UserInfo.ID,
		PublisherName: gsync.UserInfo.Name,
		Type:          model.Pull,
	}
	var err error

	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.RepoFollow, Message: "Repo follow"},
	}
	repo, _ := repository.GetRepo(gsync.Project.RepoType)
	if len(gsync.CommitID) == 0 {
		err = repo.Follow(gsync.Project, "origin/"+gsync.Project.Branch)
	} else {
		err = repo.Follow(gsync.Project, gsync.CommitID)
	}
	if err != nil {
		core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, err)
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.DeployFail, Message: err.Error()},
		}
		_ = gsync.Project.DeployFail()
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, err)
		}
		gsync.notify(model.ProjectFail, err.Error())
		return
	}

	commitList, _ := repo.CommitLog(gsync.Project.ID, 1)
	gsync.CommitInfo = commitList[0]
	if gsync.Branch != "" {
		gsync.CommitInfo.Branch = gsync.Branch
	} else {
		gsync.CommitInfo.Branch = "origin/" + gsync.Project.Branch
	}

	ext, _ := json.Marshal(gsync.CommitInfo)
	publishTraceModel.Ext = string(ext)
	publishTraceModel.State = model.Success
	if _, err := publishTraceModel.AddRow(); err != nil {
		core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, err)
		return
	}

	if totalFileNumber, err := (model.ProjectFile{ProjectID: gsync.Project.ID}).GetTotalByProjectID(); err != nil {
		core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, err)
		return
	} else if totalFileNumber > 0 {
		if err := utils.CopyDir(core.GetProjectFilePath(gsync.Project.ID), core.GetProjectPath(gsync.Project.ID)); err != nil {
			core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, err)
			return
		}
	}

	if gsync.Project.AfterPullScript != "" {
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.AfterPullScript, Message: "Run pull script"},
		}

		publishTraceModel.Type = model.AfterPull
		ext, _ = json.Marshal(struct {
			Script string `json:"script"`
		}{gsync.Project.AfterPullScript})
		publishTraceModel.Ext = string(ext)
		core.Logf(core.TRACE, projectLogFormat, gsync.Project.ID, gsync.Project.AfterPullScript)
		if outputString, err := gsync.runAfterPullScript(); err != nil {
			core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, fmt.Sprintf("err: %s, output: %s", err, outputString))
			ws.GetHub().Data <- &ws.Data{
				Type:    ws.TypeProject,
				Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.DeployFail, Message: err.Error()},
			}
			_ = gsync.Project.DeployFail()
			publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", err, outputString)
			publishTraceModel.State = model.Fail
			if _, err := publishTraceModel.AddRow(); err != nil {
				core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, err)
			}
			gsync.notify(model.ProjectFail, err.Error())
			return
		} else {
			publishTraceModel.Detail = outputString
			publishTraceModel.State = model.Success
			if _, err := publishTraceModel.AddRow(); err != nil {
				core.Logf(core.ERROR, projectLogFormat, gsync.Project.ID, err)
			}
		}
	}

	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.Rsync, Message: "Sync"},
	}
	ch := make(chan syncMessage, len(gsync.ProjectServers))

	gsync.remoteSync(ch)
	message := ""
	for i := 0; i < len(gsync.ProjectServers); i++ {
		syncMessage := <-ch
		if syncMessage.state == model.ProjectFail {
			message += syncMessage.serverName + " error message: " + syncMessage.detail
		}
	}
	close(ch)
	if message == "" {
		_ = gsync.Project.DeploySuccess()
		core.Logf(core.TRACE, projectLogFormat, gsync.Project.ID, "deploy success")
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.DeploySuccess, Message: "Success", Ext: gsync.CommitInfo},
		}
		gsync.notify(model.ProjectSuccess, message)
	} else {
		_ = gsync.Project.DeployFail()
		core.Logf(core.TRACE, projectLogFormat, gsync.Project.ID, "deploy fail")
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.DeployFail, Message: message},
		}
		gsync.notify(model.ProjectFail, message)
	}

	if gsync.Project.SymlinkPath != "" {
		gsync.removeExpiredBackup()
	}
	return
}

func (gsync Gsync) runAfterPullScript() (string, error) {
	project := gsync.Project
	commitInfo := gsync.CommitInfo
	srcPath := core.GetProjectPath(project.ID)
	scriptName := "goploy-after-pull." + utils.GetScriptExt(project.AfterPullScriptMode)
	scriptFullName := path.Join(srcPath, scriptName)
	scriptMode := "bash"
	if len(project.AfterPullScriptMode) != 0 {
		scriptMode = project.AfterPullScriptMode
	}
	scriptText := ReplaceProjectVars(ReplaceCommitVars(project.AfterPullScript, commitInfo), project)
	_ = ioutil.WriteFile(scriptFullName, []byte(scriptText), 0755)
	var commandOptions []string
	if project.AfterPullScriptMode == "cmd" {
		commandOptions = append(commandOptions, "/C")
		scriptFullName, _ = filepath.Abs(scriptFullName)
	}
	commandOptions = append(commandOptions, scriptFullName)

	handler := exec.Command(scriptMode, commandOptions...)
	handler.Dir = srcPath

	if output, err := handler.CombinedOutput(); err != nil {
		return "", err
	} else {
		_ = os.Remove(scriptName)
		return string(output), nil
	}
}

func (gsync Gsync) remoteSync(msgChIn chan<- syncMessage) {
	for _, projectServer := range gsync.ProjectServers {
		go func(projectServer model.ProjectServer) {
			project := gsync.Project
			publishTraceModel := model.PublishTrace{
				Token:         project.LastPublishToken,
				ProjectID:     project.ID,
				ProjectName:   project.Name,
				PublisherID:   gsync.UserInfo.ID,
				PublisherName: gsync.UserInfo.Name,
				Type:          model.Deploy,
			}
			// write after deploy script for rsync
			if len(project.AfterDeployScript) != 0 {
				scriptName := path.Join(core.GetProjectPath(project.ID), "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode))
				_ = ioutil.WriteFile(scriptName, []byte(ReplaceProjectVars(project.AfterDeployScript, project)), 0755)
			}

			transmitterEntity := transmitter.New(project, projectServer)
			logCmd := transmitterEntity.String()
			core.Log(core.TRACE, "projectID: "+strconv.FormatInt(project.ID, 10)+" "+logCmd)
			ext, _ := json.Marshal(struct {
				ServerID   int64  `json:"serverId"`
				ServerName string `json:"serverName"`
				Command    string `json:"command"`
			}{projectServer.ServerID, projectServer.ServerName, logCmd})
			publishTraceModel.Ext = string(ext)

			if transmitterOutput, err := transmitterEntity.Exec(); err != nil {
				core.Log(core.ERROR, fmt.Sprintf("projectID: %d rsync exec err: %s, output: %s", project.ID, err, transmitterOutput))
				publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", err, transmitterOutput)
				publishTraceModel.State = model.Fail
				publishTraceModel.AddRow()
				msgChIn <- syncMessage{
					serverName: projectServer.ServerName,
					projectID:  project.ID,
					detail:     err.Error(),
					state:      model.ProjectFail,
				}
				return
			} else {
				publishTraceModel.Detail = transmitterOutput
				publishTraceModel.State = model.Success
				publishTraceModel.AddRow()
			}

			var afterDeployCommands []string
			cmdEntity := cmd.New(projectServer.ServerOS)
			if len(project.SymlinkPath) != 0 {
				destDir := path.Join(project.SymlinkPath, project.LastPublishToken)
				afterDeployCommands = append(afterDeployCommands, cmdEntity.Symlink(destDir, project.Path))
			}

			if len(project.AfterDeployScript) != 0 {
				afterDeployScriptPath := path.Join(project.Path, "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode))
				afterDeployCommands = append(afterDeployCommands, cmdEntity.Script(project.AfterDeployScriptMode, afterDeployScriptPath))
				afterDeployCommands = append(afterDeployCommands, cmdEntity.Remove(afterDeployScriptPath))
			}

			// no symlink and deploy script
			if len(afterDeployCommands) == 0 {
				msgChIn <- syncMessage{
					serverName: projectServer.ServerName,
					projectID:  project.ID,
					state:      model.ProjectSuccess,
				}
				return
			}
			completeAfterDeployCmd := strings.Join(afterDeployCommands, "&&")
			publishTraceModel.Type = model.AfterDeploy
			ext, _ = json.Marshal(struct {
				ServerID   int64  `json:"serverId"`
				ServerName string `json:"serverName"`
				Script     string `json:"script"`
			}{projectServer.ServerID, projectServer.ServerName, completeAfterDeployCmd})
			publishTraceModel.Ext = string(ext)

			client, err := projectServer.ToSSHConfig().Dial()
			if err != nil {
				core.Log(core.ERROR, err.Error())
				publishTraceModel.Detail = err.Error()
				publishTraceModel.State = model.Fail
				publishTraceModel.AddRow()
				msgChIn <- syncMessage{
					serverName: projectServer.ServerName,
					projectID:  project.ID,
					detail:     err.Error(),
					state:      model.ProjectFail,
				}
				return
			}
			defer client.Close()

			session, sessionErr := client.NewSession()
			if sessionErr != nil {
				core.Log(core.ERROR, sessionErr.Error())
				publishTraceModel.Detail = sessionErr.Error()
				publishTraceModel.State = model.Fail
				publishTraceModel.AddRow()
				msgChIn <- syncMessage{
					serverName: projectServer.ServerName,
					projectID:  project.ID,
					detail:     sessionErr.Error(),
					state:      model.ProjectFail,
				}
				return
			}
			defer session.Close()
			core.Log(core.TRACE, fmt.Sprintf("projectID: %d ssh exec: %s", project.ID, completeAfterDeployCmd))
			if output, err := session.CombinedOutput(completeAfterDeployCmd); err != nil {
				core.Log(core.ERROR, fmt.Sprintf("projectID: %d ssh exec err: %s, output: %s", project.ID, err, output))
				publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", err, output)
				publishTraceModel.State = model.Fail
				if _, err := publishTraceModel.AddRow(); err != nil {
					core.Log(core.ERROR, "projectID: "+strconv.FormatInt(project.ID, 10)+" "+err.Error())
				}
				msgChIn <- syncMessage{
					serverName: projectServer.ServerName,
					projectID:  project.ID,
					detail:     fmt.Sprintf("%s\noutput: %s", err.Error(), output),
					state:      model.ProjectFail,
				}
				return
			} else {
				publishTraceModel.Detail = string(output)
				publishTraceModel.State = model.Success
				if _, err := publishTraceModel.AddRow(); err != nil {
					core.Log(core.ERROR, "projectID: "+strconv.FormatInt(project.ID, 10)+" "+err.Error())
				}
			}

			msgChIn <- syncMessage{
				serverName: projectServer.ServerName,
				projectID:  project.ID,
				state:      model.ProjectSuccess,
			}
			return
		}(projectServer)
	}
}

// commit id
// commit message
// server ip & name
// deploy user name
// deploy time
func (gsync Gsync) notify(deployState int, detail string) {
	if gsync.Project.NotifyType == 0 {
		return
	}
	serverList := ""
	for _, projectServer := range gsync.ProjectServers {
		if projectServer.ServerName != projectServer.ServerIP {
			serverList += projectServer.ServerName + "(" + projectServer.ServerIP + ")"
		} else {
			serverList += projectServer.ServerIP
		}
		serverList += ", "
	}
	serverList = strings.TrimRight(serverList, ", ")
	project := gsync.Project
	commitInfo := gsync.CommitInfo
	var err error
	var resp *http.Response
	if project.NotifyType == model.NotifyWeiXin {
		type markdown struct {
			Content string `json:"content"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		content := "Deploy: <font color=\"warning\">" + project.Name + "</font>\n"
		content += "Publisher: <font color=\"comment\">" + project.PublisherName + "</font>\n"
		content += "Author: <font color=\"comment\">" + commitInfo.Author + "</font>\n"
		content += "Branch: <font color=\"comment\">" + commitInfo.Branch + "</font>\n"
		content += "CommitSHA: <font color=\"comment\">" + commitInfo.Commit + "</font>\n"
		content += "CommitMessage: <font color=\"comment\">" + commitInfo.Message + "</font>\n"
		content += "ServerList: <font color=\"comment\">" + serverList + "</font>\n"
		if deployState == model.ProjectFail {
			content += "State: <font color=\"red\">fail</font> \n"
			content += "> Detail: <font color=\"comment\">" + detail + "</font>"
		} else {
			content += "State: <font color=\"green\">success</font>"
		}

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Content: content,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if project.NotifyType == model.NotifyDingTalk {
		type markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		text := "#### Deploy：" + project.Name + "  \n  "
		text += "#### Publisher：" + project.PublisherName + "  \n  "
		text += "#### Author：" + commitInfo.Author + "  \n  "
		text += "#### Branch：" + commitInfo.Branch + "  \n  "
		text += "#### CommitSHA：" + commitInfo.Commit + "  \n  "
		text += "#### CommitMessage：" + commitInfo.Message + "  \n  "
		text += "#### ServerList：" + serverList + "  \n  "
		if deployState == model.ProjectFail {
			text += "#### State： <font color=\"red\">fail</font>  \n  "
			text += "> Detail: <font color=\"comment\">" + detail + "</font>"
		} else {
			text += "#### State： <font color=\"green\">success</font>"
		}

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Title: project.Name,
				Text:  text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if project.NotifyType == model.NotifyFeiShu {
		type content struct {
			Text string `json:"text"`
		}
		type message struct {
			MsgType string  `json:"msg_type"`
			Content content `json:"content"`
		}
		text := ""
		text += "Publisher: " + project.PublisherName + "\n"
		text += "Author: " + commitInfo.Author + "\n"
		text += "Branch: " + commitInfo.Branch + "\n"
		text += "CommitSHA: " + commitInfo.Commit + "\n"
		text += "CommitMessage: " + commitInfo.Message + "\n"
		text += "ServerList: " + serverList + "\n"
		if deployState == model.ProjectFail {
			text += "State: fail\n "
			text += "Detail: " + detail
		} else {
			text += "State: success"
		}

		msg := message{
			MsgType: "text",
			Content: content{
				Text: text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if project.NotifyType == model.NotifyCustom {
		type message struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				ProjectID     int64  `json:"projectId"`
				ProjectName   string `json:"projectName"`
				Publisher     string `json:"publisher"`
				Author        string `json:"author"`
				Branch        string `json:"branch"`
				CommitSHA     string `json:"commitSHA"`
				CommitMessage string `json:"commitMessage"`
				ServerList    string `json:"serverList"`
			} `json:"data"`
		}
		code := 0
		if deployState == model.ProjectFail {
			code = 1
		}
		msg := message{
			Code:    code,
			Message: detail,
		}
		msg.Data.ProjectID = project.ID
		msg.Data.ProjectName = project.Name
		msg.Data.Publisher = project.PublisherName
		msg.Data.Author = commitInfo.Author
		msg.Data.Branch = commitInfo.Branch
		msg.Data.CommitSHA = commitInfo.Commit
		msg.Data.CommitMessage = commitInfo.Message
		msg.Data.ServerList = serverList
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	}

	if err != nil {
		core.Log(core.ERROR, fmt.Sprintf("projectID: %d notify exec err: %s", project.ID, err))
	} else {
		defer resp.Body.Close()
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			core.Log(core.ERROR, fmt.Sprintf("projectID: %d notify read body err: %s", project.ID, err))
		} else {
			core.Log(core.TRACE, fmt.Sprintf("projectID: %d notify success: %s", project.ID, string(responseData)))
		}
	}
}

//keep the latest 10 project
func (gsync Gsync) removeExpiredBackup() {
	var wg sync.WaitGroup
	for _, projectServer := range gsync.ProjectServers {
		wg.Add(1)
		go func(projectServer model.ProjectServer) {
			defer wg.Done()
			client, err := projectServer.ToSSHConfig().Dial()
			if err != nil {
				core.Log(core.ERROR, err.Error())
				return
			}
			defer client.Close()
			session, err := client.NewSession()
			if err != nil {
				core.Log(core.ERROR, err.Error())
				return
			}
			defer session.Close()
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			if err = session.Run("cd " + gsync.Project.SymlinkPath + ";ls -t | awk 'NR>" + strconv.Itoa(int(gsync.Project.SymlinkBackupNumber)) + "' | xargs rm -rf"); err != nil {
				core.Log(core.ERROR, err.Error())
			}
		}(projectServer)
	}
	wg.Wait()
}

func ReplaceCommitVars(script string, commitInfo repository.CommitInfo) string {
	scriptVars := map[string]string{
		"${COMMIT_TAG}":       commitInfo.Tag,
		"${COMMIT_BRANCH}":    commitInfo.Branch,
		"${COMMIT_ID}":        commitInfo.Commit,
		"${COMMIT_SHORT_ID}":  commitInfo.Commit,
		"${COMMIT_AUTHOR}":    commitInfo.Author,
		"${COMMIT_TIMESTAMP}": strconv.FormatInt(commitInfo.Timestamp, 10),
		"${COMMIT_MESSAGE}":   commitInfo.Message,
	}

	if len(commitInfo.Commit) > 6 {
		scriptVars["${COMMIT_SHORT_ID}"] = commitInfo.Commit[0:6]
	}

	for key, value := range scriptVars {
		script = strings.Replace(script, key, value, -1)
	}
	return script
}

func ReplaceProjectVars(script string, project model.Project) string {
	scriptVars := map[string]string{
		"${PROJECT_PATH}":         project.Path,
		"${PROJECT_SYMLINK_PATH}": path.Join(project.SymlinkPath, project.LastPublishToken),
		"${PROJECT_NAME}":         project.Name,
		"${REPOSITORY_PATH}":      core.GetProjectPath(project.ID),
	}
	for key, value := range scriptVars {
		script = strings.Replace(script, key, value, -1)
	}
	return script
}
