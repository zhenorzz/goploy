// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/cmd/server/ws"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/pkg/cmd"
	"github.com/zhenorzz/goploy/internal/repo"
	"github.com/zhenorzz/goploy/internal/transmitter"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var deployList = list.New()
var deployTick = time.Tick(time.Millisecond)

type Gsync struct {
	UserInfo       model.User
	Project        model.Project
	ProjectServers model.ProjectServers
	PublishTrace   model.PublishTrace
	CommitInfo     repo.CommitInfo
	CommitID       string
	Branch         string
}

type syncMessage struct {
	serverName string
	projectID  int64
	detail     string
	state      int
}

type deployMessage struct {
	ProjectID   int64       `json:"projectId"`
	ProjectName string      `json:"projectName"`
	State       uint8       `json:"state"`
	Message     string      `json:"message"`
	Ext         interface{} `json:"ext"`
}

const (
	Queue = iota
	Deploying
	Success
	Fail
)

func (deployMessage) CanSendTo(*ws.Client) error {
	return nil
}

var projectLogFormat = "projectID: %d %s"

func startDeployTask() {
	atomic.AddInt32(&counter, 1)
	var deployingNumber int32
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case <-deployTick:
				if atomic.LoadInt32(&deployingNumber) < config.Toml.APP.DeployLimit {
					atomic.AddInt32(&deployingNumber, 1)
					if deployElem := deployList.Front(); deployElem != nil {
						wg.Add(1)
						go func(gsync *Gsync) {
							gsync.Exec()
							atomic.AddInt32(&deployingNumber, -1)
							wg.Done()
						}(deployList.Remove(deployElem).(*Gsync))
					} else {
						atomic.AddInt32(&deployingNumber, -1)
					}
				}
			case <-stop:
				wg.Wait()
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

func AddDeployTask(gsync Gsync) {
	ws.Send(ws.Data{
		Type: ws.TypeProject,
		Message: deployMessage{
			ProjectID:   gsync.Project.ID,
			ProjectName: gsync.Project.Name,
			State:       Queue,
			Message:     "Task waiting",
			Ext: struct {
				LastPublishToken string `json:"lastPublishToken"`
			}{gsync.Project.LastPublishToken},
		},
	})
	gsync.PublishTrace = model.PublishTrace{
		Token:         gsync.Project.LastPublishToken,
		ProjectID:     gsync.Project.ID,
		ProjectName:   gsync.Project.Name,
		PublisherID:   gsync.UserInfo.ID,
		PublisherName: gsync.UserInfo.Name,
		Type:          model.QUEUE,
		State:         model.Success,
	}
	if _, err := gsync.PublishTrace.AddRow(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, "insert trace error, "+err.Error())
	}
	deployList.PushBack(&gsync)
}

func (gsync *Gsync) Exec() {
	log.Tracef(projectLogFormat, gsync.Project.ID, "deploy start")
	var err error
	defer func() {
		if err == nil {
			return
		}
		ws.Send(ws.Data{
			Type:    ws.TypeProject,
			Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Fail, Message: err.Error()},
		})
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
		if err := gsync.Project.DeployFail(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		gsync.notify(model.ProjectFail, err.Error())

	}()

	err = gsync.repoStage()
	if err != nil {
		return
	}

	err = gsync.copyLocalFileStage()
	if err != nil {
		return
	}

	err = gsync.afterPullScriptStage()
	if err != nil {
		return
	}

	err = gsync.serverStage()
	if err != nil {
		return
	}

	err = gsync.deployFinishScriptStage()
	if err != nil {
		return
	}

	if err := gsync.Project.DeploySuccess(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
	}

	gsync.PublishTrace.Type = model.PublishFinish
	gsync.PublishTrace.State = model.Success
	if _, err := gsync.PublishTrace.AddRow(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
	}

	log.Tracef(projectLogFormat, gsync.Project.ID, "deploy success")
	ws.Send(ws.Data{
		Type:    ws.TypeProject,
		Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Success, Message: "Success", Ext: gsync.CommitInfo},
	})
	gsync.notify(model.ProjectSuccess, "")

	gsync.removeExpiredBackup()
}

func (gsync *Gsync) repoStage() error {
	ws.Send(ws.Data{
		Type:    ws.TypeProject,
		Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Repo follow"},
	})

	gsync.PublishTrace.Type = model.Pull
	var err error
	r, _ := repo.GetRepo(gsync.Project.RepoType)
	if len(gsync.CommitID) == 0 {
		err = r.Follow(gsync.Project, "origin/"+gsync.Project.Branch)
	} else {
		err = r.Follow(gsync.Project, gsync.CommitID)
	}

	if err != nil {
		gsync.PublishTrace.Detail = err.Error()
		gsync.PublishTrace.State = model.Fail
		if _, err := gsync.PublishTrace.AddRow(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		return err
	}

	commitList, err := r.CommitLog(gsync.Project.ID, 1)
	if err != nil {
		gsync.PublishTrace.Detail = err.Error()
		gsync.PublishTrace.State = model.Fail
		if _, err := gsync.PublishTrace.AddRow(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		return err
	}

	gsync.CommitInfo = commitList[0]
	if gsync.Branch != "" {
		gsync.CommitInfo.Branch = gsync.Branch
	} else {
		gsync.CommitInfo.Branch = "origin/" + gsync.Project.Branch
	}

	ext, _ := json.Marshal(gsync.CommitInfo)
	gsync.PublishTrace.Ext = string(ext)
	gsync.PublishTrace.State = model.Success
	if _, err := gsync.PublishTrace.AddRow(); err != nil {
		log.Errorf(projectLogFormat, gsync.Project.ID, err)
	}
	return nil
}

func (gsync *Gsync) copyLocalFileStage() error {
	if totalFileNumber, err := (model.ProjectFile{ProjectID: gsync.Project.ID}).GetTotalByProjectID(); err != nil {
		return err
	} else if totalFileNumber > 0 {
		if err := pkg.CopyDir(config.GetProjectFilePath(gsync.Project.ID), config.GetProjectPath(gsync.Project.ID)); err != nil {
			return err
		}
	}
	return nil
}

func (gsync *Gsync) afterPullScriptStage() error {
	if gsync.Project.Script.AfterPull.Content == "" {
		return nil
	}

	gsync.PublishTrace.Type = model.AfterPull
	ext, _ := json.Marshal(struct {
		Script string `json:"script"`
	}{gsync.Project.Script.AfterPull.Content})
	gsync.PublishTrace.Ext = string(ext)
	return gsync.runLocalScript()
}

func (gsync *Gsync) serverStage() error {
	ws.Send(ws.Data{
		Type:    ws.TypeProject,
		Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Sync"},
	})
	ch := make(chan syncMessage, len(gsync.ProjectServers))
	gsync.PublishTrace.Type = model.Deploy
	var serverSync = func(projectServer model.ProjectServer, index int) {
		project := gsync.Project
		publishTraceModel := gsync.PublishTrace
		// write after deploy script for rsync
		scriptName := fmt.Sprintf("goploy-after-deploy-p%d-s%d.%s", project.ID, projectServer.ServerID, pkg.GetScriptExt(project.Script.AfterDeploy.Mode))
		if project.Script.AfterDeploy.Content != "" {
			scriptContent := project.ReplaceVars(project.Script.AfterDeploy.Content)
			scriptContent = projectServer.ReplaceVars(scriptContent)
			scriptContent = strings.Replace(scriptContent, "${SERVER_TOTAL_NUMBER}", strconv.Itoa(len(gsync.ProjectServers)), -1)
			scriptContent = strings.Replace(scriptContent, "${SERVER_SERIAL_NUMBER}", strconv.Itoa(index), -1)
			_ = os.WriteFile(path.Join(config.GetProjectPath(project.ID), scriptName), []byte(scriptContent), 0755)
		}

		transmitterEntity := transmitter.New(project, projectServer)
		logCmd := transmitterEntity.String()
		log.Trace("projectID: " + strconv.FormatInt(project.ID, 10) + " " + logCmd)
		ext, _ := json.Marshal(struct {
			ServerID   int64  `json:"serverId"`
			ServerName string `json:"serverName"`
			Command    string `json:"command"`
		}{projectServer.ServerID, projectServer.Server.Name, logCmd})
		publishTraceModel.Ext = string(ext)

		if transmitterOutput, err := transmitterEntity.Exec(); err != nil {
			log.Error(fmt.Sprintf("projectID: %d transmit exec err: %s, output: %s", project.ID, err, transmitterOutput))
			publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", err, transmitterOutput)
			publishTraceModel.State = model.Fail
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Errorf(projectLogFormat, project.ID, err)
			}

			ch <- syncMessage{
				serverName: projectServer.Server.Name,
				projectID:  project.ID,
				detail:     err.Error(),
				state:      model.ProjectFail,
			}
			return
		} else {
			publishTraceModel.Detail = transmitterOutput
			publishTraceModel.State = model.Success
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Errorf(projectLogFormat, project.ID, err)
			}
		}

		var afterDeployCommands []string
		cmdEntity := cmd.New(projectServer.Server.OS)
		if len(project.SymlinkPath) != 0 {
			destDir := path.Join(project.SymlinkPath, project.LastPublishToken)
			afterDeployCommands = append(afterDeployCommands, cmdEntity.Symlink(destDir, project.Path))
		}

		if project.Script.AfterDeploy.Content != "" {
			afterDeployScriptPath := path.Join(project.Path, scriptName)
			afterDeployCommands = append(afterDeployCommands, cmdEntity.Script(project.Script.AfterDeploy.Mode, afterDeployScriptPath))
			afterDeployCommands = append(afterDeployCommands, cmdEntity.Remove(afterDeployScriptPath))
		}

		// no symlink and deploy script
		if len(afterDeployCommands) == 0 {
			ch <- syncMessage{
				serverName: projectServer.Server.Name,
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
		}{projectServer.ServerID, projectServer.Server.Name, completeAfterDeployCmd})
		publishTraceModel.Ext = string(ext)

		client, err := projectServer.ToSSHConfig().Dial()
		if err != nil {
			log.Error(err.Error())
			publishTraceModel.Detail = err.Error()
			publishTraceModel.State = model.Fail
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Errorf(projectLogFormat, project.ID, err)
			}
			ch <- syncMessage{
				serverName: projectServer.Server.Name,
				projectID:  project.ID,
				detail:     err.Error(),
				state:      model.ProjectFail,
			}
			return
		}
		defer client.Close()

		session, sessionErr := client.NewSession()
		if sessionErr != nil {
			log.Error(sessionErr.Error())
			publishTraceModel.Detail = sessionErr.Error()
			publishTraceModel.State = model.Fail
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Errorf(projectLogFormat, project.ID, err)
			}
			ch <- syncMessage{
				serverName: projectServer.Server.Name,
				projectID:  project.ID,
				detail:     sessionErr.Error(),
				state:      model.ProjectFail,
			}
			return
		}
		defer session.Close()
		log.Trace(fmt.Sprintf("projectID: %d ssh exec: %s", project.ID, completeAfterDeployCmd))
		output, err := session.CombinedOutput(completeAfterDeployCmd)
		if err != nil {
			log.Error(fmt.Sprintf("projectID: %d ssh exec err: %s, output: %s", project.ID, err, output))
			publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", err, output)
			publishTraceModel.State = model.Fail
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Errorf(projectLogFormat, project.ID, err)
			}
			ch <- syncMessage{
				serverName: projectServer.Server.Name,
				projectID:  project.ID,
				detail:     fmt.Sprintf("%s\noutput: %s", err.Error(), output),
				state:      model.ProjectFail,
			}
			return
		}

		publishTraceModel.Detail = string(output)
		publishTraceModel.State = model.Success
		if _, err := publishTraceModel.AddRow(); err != nil {
			log.Error("projectID: " + strconv.FormatInt(project.ID, 10) + " " + err.Error())
		}
		ch <- syncMessage{
			serverName: projectServer.Server.Name,
			projectID:  project.ID,
			state:      model.ProjectSuccess,
		}
	}

	for index, projectServer := range gsync.ProjectServers {
		if gsync.Project.DeployServerMode == "serial" {
			serverSync(projectServer, index+1)
		} else {
			go serverSync(projectServer, 0)
		}
	}

	message := ""
	for i := 0; i < len(gsync.ProjectServers); i++ {
		msg := <-ch
		if msg.state == model.ProjectFail {
			message += msg.serverName + " error message: " + msg.detail
		}
	}
	close(ch)

	if message != "" {
		return errors.New(message)
	}
	return nil
}

func (gsync *Gsync) deployFinishScriptStage() error {
	if gsync.Project.Script.DeployFinish.Content == "" {
		return nil
	}

	gsync.PublishTrace.Type = model.DeployFinish
	ext, _ := json.Marshal(struct {
		Script string `json:"script"`
	}{gsync.Project.Script.DeployFinish.Content})
	gsync.PublishTrace.Ext = string(ext)

	return gsync.runLocalScript()
}

func (gsync *Gsync) runLocalScript() error {
	var mode = ""
	var content = ""
	var scriptName = ""
	project := gsync.Project
	switch gsync.PublishTrace.Type {
	case model.AfterPull:
		ws.Send(ws.Data{
			Type:    ws.TypeProject,
			Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Run pull script"},
		})
		mode = project.Script.AfterPull.Mode
		content = project.Script.AfterPull.Content
		scriptName = "goploy-after-pull." + pkg.GetScriptExt(project.Script.AfterPull.Mode)

	case model.DeployFinish:
		ws.Send(ws.Data{
			Type:    ws.TypeProject,
			Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Run finish script"},
		})
		mode = project.Script.DeployFinish.Mode
		content = project.Script.DeployFinish.Content
		scriptName = "goploy-deploy-finish." + pkg.GetScriptExt(project.Script.DeployFinish.Mode)

	default:
		return errors.New("not support stage")
	}

	log.Tracef(projectLogFormat, gsync.Project.ID, content)

	commitInfo := gsync.CommitInfo
	srcPath := config.GetProjectPath(project.ID)
	scriptFullName := path.Join(srcPath, scriptName)
	scriptMode := "bash"
	if mode != "" {
		scriptMode = mode
	}
	scriptText := project.ReplaceVars(commitInfo.ReplaceVars(content))
	_ = os.WriteFile(scriptFullName, []byte(scriptText), 0755)
	var commandOptions []string
	if scriptMode == "cmd" {
		commandOptions = append(commandOptions, "/C")
		scriptFullName, _ = filepath.Abs(scriptFullName)
	}
	commandOptions = append(commandOptions, scriptFullName)

	handler := exec.Command(scriptMode, commandOptions...)
	handler.Dir = srcPath

	if output, err := handler.CombinedOutput(); err != nil {
		gsync.PublishTrace.Detail = fmt.Sprintf("err: %s\noutput: %s", err, string(output))
		gsync.PublishTrace.State = model.Fail
		if _, err := gsync.PublishTrace.AddRow(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		return fmt.Errorf("err: %s, output: %s", err, string(output))
	} else {
		_ = os.Remove(scriptFullName)
		gsync.PublishTrace.Detail = string(output)
		gsync.PublishTrace.State = model.Success
		if _, err := gsync.PublishTrace.AddRow(); err != nil {
			log.Errorf(projectLogFormat, gsync.Project.ID, err)
		}
		return nil
	}
}

// commit id
// commit message
// server ip & name
// deploy user name
// deploy time
func (gsync *Gsync) notify(deployState int, detail string) {
	if gsync.Project.NotifyType == 0 {
		return
	}
	serverList := ""
	for _, projectServer := range gsync.ProjectServers {
		if projectServer.Server.Name != projectServer.Server.IP {
			serverList += projectServer.Server.Name + "(" + projectServer.Server.IP + ")"
		} else {
			serverList += projectServer.Server.IP
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
		if commitInfo.Tag != "" {
			content += "Tag: <font color=\"comment\">" + commitInfo.Tag + "</font>\n"
		}
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
		if commitInfo.Tag != "" {
			text += "#### Tag：" + commitInfo.Tag + "  \n  "
		}
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
		text += "Deploy：" + project.Name + "\n"
		text += "Publisher: " + project.PublisherName + "\n"
		text += "Author: " + commitInfo.Author + "\n"
		if commitInfo.Tag != "" {
			text += "Tag: " + commitInfo.Tag + "\n"
		}
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
				Tag           string `json:"tag"`
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
		msg.Data.Tag = commitInfo.Tag
		msg.Data.CommitSHA = commitInfo.Commit
		msg.Data.CommitMessage = commitInfo.Message
		msg.Data.ServerList = serverList
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	}

	if err != nil {
		log.Error(fmt.Sprintf("projectID: %d notify exec err: %s", project.ID, err))
	} else {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(fmt.Sprintf("projectID: %d notify read body err: %s", project.ID, err))
		} else {
			log.Trace(fmt.Sprintf("projectID: %d notify success: %s", project.ID, string(responseData)))
		}
		_ = resp.Body.Close()
	}
}

// keep the latest 10 project
func (gsync *Gsync) removeExpiredBackup() {
	if gsync.Project.SymlinkPath == "" {
		return
	}
	var wg sync.WaitGroup
	for _, projectServer := range gsync.ProjectServers {
		wg.Add(1)
		go func(projectServer model.ProjectServer) {
			defer wg.Done()
			client, err := projectServer.ToSSHConfig().Dial()
			if err != nil {
				log.Error(err.Error())
				return
			}
			defer func() {
				if err := client.Close(); err != nil {
					log.Error(err.Error())
				}
			}()
			session, err := client.NewSession()
			if err != nil {
				log.Error(err.Error())
				return
			}
			defer func() {
				if err := session.Close(); err != nil {
					log.Error(err.Error())
				}
			}()
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			if err = session.Run("cd " + gsync.Project.SymlinkPath + ";ls -t | awk 'NR>" + strconv.Itoa(int(gsync.Project.SymlinkBackupNumber)) + "' | xargs rm -rf"); err != nil {
				log.Error(err.Error())
			}
		}(projectServer)
	}
	wg.Wait()
}
