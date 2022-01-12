package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/repository"
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

// Gsync -
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

func (gsync Gsync) Exec() {
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(gsync.Project.ID, 10)+" deploy start")
	publishTraceModel := model.PublishTrace{
		Token:         gsync.Project.LastPublishToken,
		ProjectID:     gsync.Project.ID,
		ProjectName:   gsync.Project.Name,
		PublisherID:   gsync.UserInfo.ID,
		PublisherName: gsync.UserInfo.Name,
		Type:          model.Pull,
	}
	var err error

	repo, err := repository.GetRepo(gsync.Project.RepoType)
	if err != nil {
		_ = gsync.Project.DeployFail()
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.DeployFail, Message: err.Error()},
		}
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
		return
	}

	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.RepoFollow, Message: "Repo follow"},
	}

	if len(gsync.CommitID) == 0 {
		err = repo.Follow(gsync.Project, "origin/"+gsync.Project.Branch)
	} else {
		err = repo.Follow(gsync.Project, gsync.CommitID)
	}
	if err != nil {
		_ = gsync.Project.DeployFail()
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.DeployFail, Message: err.Error()},
		}
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
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
		core.Log(core.ERROR, err.Error())
		return
	}

	if totalFileNumber, err := (model.ProjectFile{ProjectID: gsync.Project.ID}).GetTotalByProjectID(); err != nil {
		core.Log(core.ERROR, err.Error())
		return
	} else if totalFileNumber > 0 {
		if err := utils.CopyDir(core.GetProjectFilePath(gsync.Project.ID), core.GetProjectPath(gsync.Project.ID)); err != nil {
			core.Log(core.ERROR, err.Error())
			return
		}
	}

	if gsync.Project.AfterPullScript != "" {
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.AfterPullScript, Message: "Run pull script"},
		}
		outputString, err := gsync.runAfterPullScript()
		publishTraceModel.Type = model.AfterPull
		ext, _ := json.Marshal(struct {
			Script string `json:"script"`
		}{gsync.Project.AfterPullScript})
		publishTraceModel.Ext = string(ext)
		if err != nil {
			_ = gsync.Project.DeployFail()
			publishTraceModel.Detail = err.Error()
			publishTraceModel.State = model.Fail
			ws.GetHub().Data <- &ws.Data{
				Type:    ws.TypeProject,
				Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.DeployFail, Message: err.Error()},
			}
			if _, err := publishTraceModel.AddRow(); err != nil {
				core.Log(core.ERROR, err.Error())
			}
			gsync.notify(model.ProjectFail, err.Error())
			return
		}
		publishTraceModel.Detail = outputString
		publishTraceModel.State = model.Success
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
	}

	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: ws.Rsync, Message: "Rsync"},
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
		core.Log(core.TRACE, "projectID:"+strconv.FormatInt(gsync.Project.ID, 10)+" deploy success")
		ws.GetHub().Data <- &ws.Data{
			Type: ws.TypeProject,
			Message: ws.ProjectMessage{
				ProjectID:   gsync.Project.ID,
				ProjectName: gsync.Project.Name,
				State:       ws.DeploySuccess,
				Message:     "Success",
				Ext:         gsync.CommitInfo,
			},
		}
		gsync.notify(model.ProjectSuccess, message)

	} else {
		_ = gsync.Project.DeployFail()
		core.Log(core.TRACE, "projectID:"+strconv.FormatInt(gsync.Project.ID, 10)+" deploy fail")
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
	var outbuf, errbuf bytes.Buffer
	handler.Stdout = &outbuf
	handler.Stderr = &errbuf
	core.Log(core.TRACE, fmt.Sprintf("projectID:%d %s", project.ID, project.AfterPullScript))
	if err := handler.Run(); err != nil {
		core.Log(core.ERROR, errbuf.String())
		return "", errors.New(errbuf.String())
	}
	_ = os.Remove(scriptName)
	return outbuf.String(), nil
}

func (gsync Gsync) remoteSync(msgChIn chan<- syncMessage) {
	for _, projectServer := range gsync.ProjectServers {
		go func(projectServer model.ProjectServer) {
			project := gsync.Project
			userInfo := gsync.UserInfo
			remoteMachine := projectServer.ServerOwner + "@" + projectServer.ServerIP
			destDir := project.Path
			ext, _ := json.Marshal(struct {
				ServerID   int64  `json:"serverId"`
				ServerName string `json:"serverName"`
			}{projectServer.ServerID, projectServer.ServerName})
			publishTraceModel := model.PublishTrace{
				Token:         project.LastPublishToken,
				ProjectID:     project.ID,
				ProjectName:   project.Name,
				PublisherID:   userInfo.ID,
				PublisherName: userInfo.Name,
				Type:          model.Deploy,
				Ext:           string(ext),
			}
			// write after deploy script for rsync
			if len(project.AfterDeployScript) != 0 {
				scriptName := path.Join(core.GetProjectPath(project.ID), "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode))
				_ = ioutil.WriteFile(scriptName, []byte(ReplaceProjectVars(project.AfterDeployScript, project)), 0755)
			}
			rsyncOption, _ := utils.ParseCommandLine(project.RsyncOption)
			rsyncOption = append([]string{"--exclude", "goploy-after-pull.sh", "--include", "goploy-after-deploy.sh"}, rsyncOption...)
			rsyncOption = append(rsyncOption, "-e", fmt.Sprintf("ssh -p %d -o StrictHostKeyChecking=no -i %s", projectServer.ServerPort, projectServer.ServerPath))
			if len(project.SymlinkPath) != 0 {
				destDir = path.Join(project.SymlinkPath, project.LastPublishToken)
			}
			srcPath := core.GetProjectPath(project.ID) + "/"
			// rsync path can not contain colon
			// windows like C:\
			if strings.Contains(srcPath, ":\\") {
				srcPath = "/cygdrive/" + strings.Replace(srcPath, ":\\", "/", 1)
			}
			destPath := remoteMachine + ":" + destDir
			rsyncOption = append(rsyncOption, "--rsync-path=mkdir -p "+destDir+" && rsync", srcPath, destPath)
			cmd := exec.Command("rsync", rsyncOption...)
			var outbuf, errbuf bytes.Buffer
			cmd.Stdout = &outbuf
			cmd.Stderr = &errbuf
			core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" rsync "+strings.Join(rsyncOption, " "))

			if err := cmd.Run(); err != nil {
				core.Log(core.ERROR, "err: "+err.Error()+", detail: "+errbuf.String())
				publishTraceModel.Detail = "err: " + err.Error() + ", detail: " + errbuf.String()
				publishTraceModel.State = model.Fail
				publishTraceModel.AddRow()
				msgChIn <- syncMessage{
					serverName: projectServer.ServerName,
					projectID:  project.ID,
					detail:     "err: " + err.Error() + ", detail: " + errbuf.String(),
					state:      model.ProjectFail,
				}
				return
			}

			ext, _ = json.Marshal(struct {
				ServerID   int64  `json:"serverId"`
				ServerName string `json:"serverName"`
				Command    string `json:"command"`
			}{projectServer.ServerID, projectServer.ServerName, "rsync " + strings.Join(rsyncOption, " ")})
			publishTraceModel.Ext = string(ext)
			publishTraceModel.Detail = outbuf.String()
			publishTraceModel.State = model.Success
			publishTraceModel.AddRow()

			var afterDeployCommands []string
			if len(project.SymlinkPath) != 0 {
				// use relative path to fix docker symlink
				relativeDestDir := strings.Replace(destDir, path.Dir(project.Path), ".", 1)
				afterDeployCommands = append(afterDeployCommands, "ln -sfn "+relativeDestDir+" "+project.Path)
				// change the destination folder time, make sure it can not be clean
				afterDeployCommands = append(afterDeployCommands, "touch -m "+destDir)
			}

			if len(project.AfterDeployScript) != 0 {
				scriptMode := "bash"
				if len(project.AfterDeployScriptMode) != 0 {
					scriptMode = project.AfterDeployScriptMode
				}
				afterDeployScriptPath := path.Join(project.Path, "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode))
				afterDeployCommands = append(afterDeployCommands, scriptMode+" "+afterDeployScriptPath)
				afterDeployCommands = append(afterDeployCommands, "rm -f "+afterDeployScriptPath)
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

			publishTraceModel.Type = model.AfterDeploy
			ext, _ = json.Marshal(struct {
				ServerID   int64  `json:"serverId"`
				ServerName string `json:"serverName"`
				Script     string `json:"script"`
			}{projectServer.ServerID, projectServer.ServerName, strings.Join(afterDeployCommands, ";")})
			publishTraceModel.Ext = string(ext)

			client, dialError := utils.DialSSH(projectServer.ServerOwner, projectServer.ServerPassword, projectServer.ServerPath, projectServer.ServerIP, projectServer.ServerPort)
			if dialError != nil {
				core.Log(core.ERROR, dialError.Error())
				publishTraceModel.Detail = dialError.Error()
				publishTraceModel.State = model.Fail
				publishTraceModel.AddRow()
				msgChIn <- syncMessage{
					serverName: projectServer.ServerName,
					projectID:  project.ID,
					detail:     dialError.Error(),
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

			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			if err := session.Run(strings.Join(afterDeployCommands, "&&")); err != nil {
				core.Log(core.ERROR, "ssh exec err: "+err.Error())
				publishTraceModel.Detail = "err: " + err.Error() + ", detail: " + sshErrbuf.String()
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

			publishTraceModel.Detail = sshOutbuf.String()
			publishTraceModel.State = model.Success
			publishTraceModel.AddRow()

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
		core.Log(core.ERROR, fmt.Sprintf("notify projectID:%d %s", project.ID, err.Error()))
	} else {
		defer resp.Body.Close()
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			core.Log(core.ERROR, fmt.Sprintf("notify projectID:%d read body error", project.ID))
		} else {
			core.Log(core.TRACE, fmt.Sprintf("notify projectID:%d return %s", project.ID, string(responseData)))
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
			client, err := utils.DialSSH(projectServer.ServerOwner, projectServer.ServerPassword, projectServer.ServerPath, projectServer.ServerIP, projectServer.ServerPort)
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
