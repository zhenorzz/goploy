package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"github.com/zhenorzz/goploy/ws"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Sync -
type Sync struct {
	UserInfo       model.User
	Project        model.Project
	ProjectServers model.ProjectServers
	CommitID       string
}

type syncMessage struct {
	serverName string
	projectID  int64
	detail     string
	state      int
}

// Exec Sync
func (sync Sync) Exec() {
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(sync.Project.ID, 10)+" deploy start")
	publishTraceModel := model.PublishTrace{
		Token:         sync.Project.LastPublishToken,
		ProjectID:     sync.Project.ID,
		ProjectName:   sync.Project.Name,
		PublisherID:   sync.UserInfo.ID,
		PublisherName: sync.UserInfo.Name,
		Type:          model.Pull,
	}
	var gitCommitInfo utils.Commit
	var err error
	if len(sync.CommitID) == 0 {
		gitCommitInfo, err = gitSync(sync.Project)
	} else {
		gitCommitInfo, err = gitRollback(sync.CommitID, sync.Project)
	}
	if err != nil {
		sync.Project.DeployFail()
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: sync.Project.ID, ProjectName: sync.Project.Name, State: ws.ProjectFail, Message: err.Error()},
		}
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
		go notify(sync.Project, model.ProjectFail, err.Error())
		return
	}
	ext, _ := json.Marshal(gitCommitInfo)
	publishTraceModel.Ext = string(ext)
	publishTraceModel.State = model.Success
	if _, err := publishTraceModel.AddRow(); err != nil {
		core.Log(core.ERROR, err.Error())
	}
	if sync.Project.AfterPullScript != "" {
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: sync.Project.ID, ProjectName: sync.Project.Name, State: ws.AfterPullScript, Message: "Run pull script"},
		}
		outputString, err := runAfterPullScript(sync.Project)
		publishTraceModel.Type = model.AfterPull
		ext, _ := json.Marshal(struct {
			Script string `json:"script"`
		}{sync.Project.AfterPullScript})
		publishTraceModel.Ext = string(ext)
		if err != nil {
			sync.Project.DeployFail()
			publishTraceModel.Detail = err.Error()
			publishTraceModel.State = model.Fail
			ws.GetHub().Data <- &ws.Data{
				Type:    ws.TypeProject,
				Message: ws.ProjectMessage{ProjectID: sync.Project.ID, ProjectName: sync.Project.Name, State: ws.ProjectFail, Message: err.Error()},
			}
			if _, err := publishTraceModel.AddRow(); err != nil {
				core.Log(core.ERROR, err.Error())
			}
			go notify(sync.Project, model.ProjectFail, err.Error())
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
		Message: ws.ProjectMessage{ProjectID: sync.Project.ID, ProjectName: sync.Project.Name, State: ws.Rsync, Message: "Rsync"},
	}
	ch := make(chan syncMessage, len(sync.ProjectServers))
	for _, projectServer := range sync.ProjectServers {
		go remoteSync(ch, sync.UserInfo, sync.Project, projectServer)
	}

	message := ""
	for i := 0; i < len(sync.ProjectServers); i++ {
		syncMessage := <-ch
		if syncMessage.state == model.ProjectFail {
			message += syncMessage.serverName + " error message: " + syncMessage.detail
		}
	}
	if message == "" {
		sync.Project.DeploySuccess()
		core.Log(core.TRACE, "projectID:"+strconv.FormatInt(sync.Project.ID, 10)+" deploy success")
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: sync.Project.ID, ProjectName: sync.Project.Name, State: ws.ProjectSuccess, Message: "Success"},
		}
		go notify(sync.Project, model.ProjectSuccess, message)

	} else {
		sync.Project.DeployFail()
		core.Log(core.TRACE, "projectID:"+strconv.FormatInt(sync.Project.ID, 10)+" deploy fail")
		ws.GetHub().Data <- &ws.Data{
			Type:    ws.TypeProject,
			Message: ws.ProjectMessage{ProjectID: sync.Project.ID, ProjectName: sync.Project.Name, State: ws.ProjectFail, Message: message},
		}
		go notify(sync.Project, model.ProjectFail, message)

	}

	clean(sync.Project, sync.ProjectServers)
	return
}

func gitSync(project model.Project) (utils.Commit, error) {
	if err := gitCreate(project); err != nil {
		return utils.Commit{}, err
	}

	if err := gitPull(project); err != nil {
		return utils.Commit{}, err
	}

	commit, err := gitCommitLog(project)
	if err != nil {
		return utils.Commit{}, err
	}
	ws.GetHub().Data <- &ws.Data{
		Type: ws.TypeProject,
		Message: ws.ProjectMessage{
			ProjectID:   project.ID,
			ProjectName: project.Name,
			State:       ws.GitPull,
			Message:     "Get pull info",
			Ext:         commit,
		},
	}
	return commit, err
}

func gitRollback(commitSha string, project model.Project) (utils.Commit, error) {
	if err := gitReset(commitSha, project); err != nil {
		return utils.Commit{}, err
	}

	commit, err := gitCommitLog(project)
	if err != nil {
		return utils.Commit{}, err
	}
	ws.GetHub().Data <- &ws.Data{
		Type: ws.TypeProject,
		Message: ws.ProjectMessage{
			ProjectID:   project.ID,
			ProjectName: project.Name,
			State:       ws.GitReset,
			Message:     "Get pull info",
			Ext:         commit,
		},
	}
	return commit, err
}

func gitCreate(project model.Project) error {
	if err := (Repository{ProjectID: project.ID}.Create()); err != nil {
		return err
	}
	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: project.ID, ProjectName: project.Name, State: ws.GitCreate, Message: "git create"},
	}
	return nil
}

func gitPull(project model.Project) error {
	git := utils.GIT{Dir: core.GetProjectPath(project.Name)}
	// git clean removes all not tracked files
	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: project.ID, ProjectName: project.Name, State: ws.GitClean, Message: "git clean"},
	}
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git clean -f")
	if err := git.Clean([]string{"-f"}); err != nil {
		core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
		return errors.New(git.Err.String())
	}

	// git checkout clears all not staged changes.
	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: project.ID, ProjectName: project.Name, State: ws.GitCheckout, Message: "git checkout"},
	}
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git checkout -- .")
	if err := git.Checkout([]string{"--", "."}); err != nil {
		core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
		return errors.New(git.Err.String())
	}

	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: project.ID, ProjectName: project.Name, State: ws.GitPull, Message: "git pull"},
	}
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git pull")
	if err := git.Pull([]string{}); err != nil {
		core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
		return errors.New(git.Err.String())
	}
	return nil
}

func gitReset(commit string, project model.Project) error {
	srcPath := core.GetProjectPath(project.Name)
	ws.GetHub().Data <- &ws.Data{
		Type:    ws.TypeProject,
		Message: ws.ProjectMessage{ProjectID: project.ID, ProjectName: project.Name, State: ws.GitReset, Message: "git reset"},
	}
	resetCmd := exec.Command("git", "reset", "--hard", commit)
	resetCmd.Dir = srcPath
	var resetOutbuf, resetErrbuf bytes.Buffer
	resetCmd.Stdout = &resetOutbuf
	resetCmd.Stderr = &resetErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git reset --hard "+commit)
	if err := resetCmd.Run(); err != nil {
		core.Log(core.ERROR, resetErrbuf.String())
		return errors.New(resetErrbuf.String())
	}

	core.Log(core.TRACE, resetOutbuf.String())
	return nil
}

func gitCommitLog(project model.Project) (utils.Commit, error) {
	git := utils.GIT{Dir: core.GetProjectPath(project.Name)}

	if err := git.Log([]string{"--stat", "--pretty=format:`start`%H`%an`%at`%s`", "-n", "1"}); err != nil {
		core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
		return utils.Commit{}, errors.New(git.Err.String())
	}
	commitList := utils.ParseGITLog(git.Output.String())
	return commitList[0], nil
}

func runAfterPullScript(project model.Project) (string, error) {
	srcPath := core.GetProjectPath(project.Name)
	scriptName := "goploy-after-pull." + utils.GetScriptExt(project.AfterPullScriptMode)
	scriptFullName := path.Join(srcPath, scriptName)
	scriptMode := "bash"
	if len(project.AfterPullScriptMode) != 0 {
		scriptMode = project.AfterPullScriptMode
	}
	ioutil.WriteFile(scriptFullName, []byte(project.AfterPullScript), 0755)
	handler := exec.Command(scriptMode, path.Join(".", scriptName))
	handler.Dir = srcPath
	var outbuf, errbuf bytes.Buffer
	handler.Stdout = &outbuf
	handler.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+project.AfterPullScript)
	if err := handler.Run(); err != nil {
		core.Log(core.ERROR, errbuf.String())
		return "", errors.New(errbuf.String())
	}

	os.Remove(scriptName)
	return outbuf.String(), nil
}

func remoteSync(chInput chan<- syncMessage, userInfo model.User, project model.Project, projectServer model.ProjectServer) {
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

	if len(project.AfterDeployScript) != 0 {
		scriptName := path.Join(core.GetProjectPath(project.Name), "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode))
		ioutil.WriteFile(scriptName, []byte(project.AfterDeployScript), 0755)
	}

	rsyncOption, _ := utils.ParseCommandLine(project.RsyncOption)
	rsyncOption = append(rsyncOption, "-e", "ssh -p "+strconv.Itoa(int(projectServer.ServerPort))+" -o StrictHostKeyChecking=no")
	if len(project.SymlinkPath) != 0 {
		destDir = path.Join(project.SymlinkPath, project.Name, project.LastPublishToken)
		rsyncOption = append(rsyncOption, "--rsync-path=mkdir -p "+destDir+" && rsync")
	}
	srcPath := core.GetProjectPath(project.Name) + "/"
	destPath := remoteMachine + ":" + destDir
	rsyncOption = append(rsyncOption, srcPath, destPath)
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
		chInput <- syncMessage{
			serverName: projectServer.ServerName,
			projectID:  project.ID,
			detail:     errbuf.String(),
			state:      model.ProjectFail,
		}
		return
	}

	var afterDeployCommands []string
	if len(project.SymlinkPath) != 0 {
		afterDeployCommands = append(afterDeployCommands, "ln -sfn "+destDir+" "+project.Path)
		// change the destination folder time, make sure it can not be clean
		afterDeployCommands = append(afterDeployCommands, "touch -m "+destDir)
	}

	if len(project.AfterDeployScript) != 0 {
		scriptMode := "bash"
		if len(project.AfterDeployScript) != 0 {
			scriptMode = project.AfterDeployScriptMode
		}
		afterDeployCommands = append(afterDeployCommands, scriptMode+" "+path.Join(project.Path, "goploy-after-deploy."+utils.GetScriptExt(project.AfterDeployScriptMode)))
	}

	// no symlink and deploy script
	if len(afterDeployCommands) == 0 {
		chInput <- syncMessage{
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
			if scriptError = session.Run(strings.Join(afterDeployCommands, ";")); scriptError != nil {
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
		chInput <- syncMessage{
			serverName: projectServer.ServerName,
			projectID:  project.ID,
			detail:     connectError.Error(),
			state:      model.ProjectFail,
		}
		return
	} else if scriptError != nil {
		publishTraceModel.Detail = scriptError.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		chInput <- syncMessage{
			serverName: projectServer.ServerName,
			projectID:  project.ID,
			detail:     scriptError.Error(),
			state:      model.ProjectFail,
		}
		return
	}
	chInput <- syncMessage{
		serverName: projectServer.ServerName,
		projectID:  project.ID,
		state:      model.ProjectSuccess,
	}
	return
}

func notify(project model.Project, deployState int, detail string) {
	if project.NotifyType == 0 {
		return
	} else if project.NotifyType == model.NotifyWeiXin {
		type markdown struct {
			Content string `json:"content"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		content := "Deploy: <font color=\"warning\">" + project.Name + "</font>\n "

		if deployState == model.ProjectFail {
			content += "> State: <font color=\"red\">fail</font> \n "
			content += "> Detail: <font color=\"comment\">" + detail + "</font>"
		} else {
			content += "> State: <font color=\"green\">success</font>"
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
			core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" "+err.Error())
		}
	} else if project.NotifyType == model.NotifyDingTalk {
		type message struct {
			Msgtype string `json:"msgtype"`
			Title   string `json:"title"`
			Text    string `json:"text"`
		}
		text := ""
		if deployState == model.ProjectFail {
			text += "> State: <font color=\"red\">fail</font> \n "
			text += "> Detail: <font color=\"comment\">" + detail + "</font>"
		} else {
			text += "> State: <font color=\"green\">success</font>"
		}

		msg := message{
			Msgtype: "markdown",
			Title:   "Deploy:" + project.Name,
			Text:    text,
		}
		b, _ := json.Marshal(msg)
		_, err := http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
		if err != nil {
			core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" "+err.Error())
		}
	} else if project.NotifyType == model.NotifyFeiShu {
		type message struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}
		text := ""
		if deployState == model.ProjectFail {
			text += "State: fail\n "
			text += "Detail: " + detail
		} else {
			text += "State: success"
		}

		msg := message{
			Title: "Deploy:" + project.Name,
			Text:  text,
		}
		b, _ := json.Marshal(msg)
		_, err := http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
		if err != nil {
			core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" "+err.Error())
		}
	} else if project.NotifyType == model.NotifyCustom {
		type message struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				ProjectID   int64  `json:"projectId"`
				ProjectName string `json:"projectName"`
				Branch      string `json:"branch"`
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
		msg.Data.Branch = project.Branch
		b, _ := json.Marshal(msg)
		_, err := http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
		if err != nil {
			core.Log(core.ERROR, "projectID:"+strconv.FormatInt(project.ID, 10)+" "+err.Error())
		}
	}
}

//clean the expired backup
func clean(project model.Project, projectServers model.ProjectServers) {
	if len(project.SymlinkPath) == 0 {
		return
	}
	for _, projectServer := range projectServers {
		go removeExpiredBackup(project, projectServer)
	}
}

//keep the latest 10 project
func removeExpiredBackup(project model.Project, projectServer model.ProjectServer) {
	var session *ssh.Session
	var connectError error
	var scriptError error
	session, connectError = utils.ConnectSSH(projectServer.ServerOwner, "", projectServer.ServerIP, int(projectServer.ServerPort))
	if connectError != nil {
		core.Log(core.ERROR, connectError.Error())
		return
	}
	defer session.Close()
	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	destDir := path.Join(project.SymlinkPath, project.Name)
	if scriptError = session.Run("cd " + destDir + ";ls -t | awk 'NR>10' | xargs rm -rf"); scriptError != nil {
		core.Log(core.ERROR, scriptError.Error())
	}
}
