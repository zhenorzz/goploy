package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"github.com/zhenorzz/goploy/ws"
	"golang.org/x/crypto/ssh"
)

// Deploy struct
type Deploy Controller

// GetList deploy list
func (deploy Deploy) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Project model.Projects `json:"projectList"`
	}

	projects, err := model.ProjectUser{UserID: gp.TokenInfo.ID}.GetDepolyListByUserID()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	response := core.Response{Data: RepData{Project: projects}}
	response.JSON(w)
}

// GetDetail deploy detail
func (deploy Deploy) GetDetail(w http.ResponseWriter, gp *core.Goploy) {

	type RepData struct {
		GitTrace        model.GitTrace     `json:"gitTrace"`
		RemoteTraceList model.RemoteTraces `json:"remoteTraceList"`
	}

	id, err := strconv.Atoi(gp.URLQuery.Get("id"))
	if err != nil {
		response := core.Response{Code: 1, Message: "id参数错误"}
		response.JSON(w)
		return
	}

	gitTrace, err := model.GitTrace{ProjectID: uint32(id)}.GetLatestRow()
	if err == sql.ErrNoRows {
		response := core.Response{Code: 1, Message: "项目尚无构建记录"}
		response.JSON(w)
		return
	} else if err != nil {
		response := core.Response{Code: 1, Message: "GitTrace.GetLatestRow失败"}
		response.JSON(w)
		return
	}

	remoteTracesList, err := model.RemoteTrace{GitTraceID: gitTrace.ID}.GetListByGitTraceID()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	response := core.Response{Data: RepData{GitTrace: gitTrace, RemoteTraceList: remoteTracesList}}
	response.JSON(w)
}

// Sync the publish information in websocket
func (deploy Deploy) Sync(w http.ResponseWriter, gp *core.Goploy) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if strings.Contains(r.Header.Get("origin"), strings.Split(r.Host, ":")[0]) {
				return true
			}
			return false
		},
	}
	c, err := upgrader.Upgrade(w, gp.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	projectUsers, err := model.ProjectUser{UserID: gp.TokenInfo.ID}.GetListByUserID()
	if err != nil || len(projectUsers) == 0 {
		c.WriteJSON(&ws.SyncBroadcast{
			DataType: 0,
			Message:  "没有绑定服务器",
		})
		c.Close()
		return
	}
	projectMap := make(map[uint32]struct{})
	for _, projectUser := range projectUsers {
		projectMap[projectUser.ProjectID] = struct{}{}
	}
	ws.GetSyncHub().Register <- &ws.SyncClient{
		Conn:       c,
		UserID:     gp.TokenInfo.ID,
		UserName:   gp.TokenInfo.Name,
		ProjectMap: projectMap,
	}
}

// Publish the project
func (deploy Deploy) Publish(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)

	if err := json.Unmarshal(body, &reqData); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	project, err := model.Project{
		ID: reqData.ID,
	}.GetData()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	projectServers, err := model.ProjectServer{ProjectID: reqData.ID}.GetBindServerListByProjectID()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	go execSync(gp.TokenInfo, project, projectServers)
	project.PublisherID = gp.TokenInfo.ID
	project.PublisherName = gp.TokenInfo.Name
	project.UpdateTime = time.Now().Unix()
	_ = project.Publish()

	response := core.Response{Message: "部署中，请稍后"}
	response.JSON(w)
}
func execSync(tokenInfo core.TokenInfo, project model.Project, projectServers model.ProjectServers) {
	gitTraceModel := model.GitTrace{
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   tokenInfo.ID,
		PublisherName: tokenInfo.Name,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}
	if err := gitCreate(tokenInfo, project); err != nil {
		gitTraceModel.Detail = err.Error()
		gitTraceModel.State = 0
		gitTraceModel.AddRow()
		return
	}
	stdout, err := gitSync(tokenInfo, project)
	if err != nil {
		gitTraceModel.Detail = err.Error()
		gitTraceModel.State = 0
		gitTraceModel.AddRow()
		return
	}

	gitTraceModel.Detail = stdout
	gitTraceModel.State = 1
	gitTraceID, _ := gitTraceModel.AddRow()

	for _, projectServer := range projectServers {
		go remoteSync(tokenInfo, gitTraceID, project, projectServer)
	}
}

func gitCreate(tokenInfo core.TokenInfo, project model.Project) error {
	srcPath := core.GolbalPath + "repository/" + project.Name
	if _, err := os.Stat(srcPath); err != nil {
		if err := os.RemoveAll(srcPath); err != nil {
			return err
		}
		repo := project.URL
		cmd := exec.Command("git", "clone", repo, srcPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化 git clone")
		ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.GitType,
			State:     ws.Success,
			Message:   "项目初始化 git clone",
		}
		if err := cmd.Run(); err != nil {
			core.Log(core.ERROR, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化失败:"+err.Error())
			ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.GitType,
				State:     ws.Fail,
				Message:   "项目初始化失败",
			}
			return err
		}
		core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化成功")
		ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.GitType,
			State:     ws.Success,
			Message:   "项目初始化成功",
		}
	}
	return nil
}

func gitSync(tokenInfo core.TokenInfo, project model.Project) (string, error) {
	srcPath := core.GolbalPath + "repository/" + project.Name

	clean := exec.Command("git", "clean", "-f")
	clean.Dir = srcPath
	var cleanOutbuf, cleanErrbuf bytes.Buffer
	clean.Stdout = &cleanOutbuf
	clean.Stderr = &cleanErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" git clean -f")
	ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
		ProjectID: project.ID,
		UserID:    tokenInfo.ID,
		DataType:  ws.GitType,
		State:     ws.Success,
		Message:   "git clean -f",
	}
	if err := clean.Run(); err != nil {
		core.Log(core.ERROR, cleanErrbuf.String())
		ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.GitType,
			State:     ws.Fail,
			Message:   cleanErrbuf.String(),
		}
		return "", err
	}
	pull := exec.Command("git", "pull")
	pull.Dir = srcPath
	var pullOutbuf, pullErrbuf bytes.Buffer
	pull.Stdout = &pullOutbuf
	pull.Stderr = &pullErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" git pull")
	ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
		ProjectID: project.ID,
		UserID:    tokenInfo.ID,
		DataType:  ws.GitType,
		State:     ws.Success,
		Message:   "git pull",
	}
	if err := pull.Run(); err != nil {
		core.Log(core.ERROR, pullErrbuf.String())
		ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.GitType,
			State:     ws.Fail,
			Message:   pullErrbuf.String(),
		}
		return "", err
	}

	core.Log(core.TRACE, pullOutbuf.String())
	ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{
		ProjectID: project.ID,
		UserID:    tokenInfo.ID,
		DataType:  ws.GitType,
		State:     ws.Success,
		Message:   pullOutbuf.String(),
	}
	return pullOutbuf.String(), nil
}

func remoteSync(tokenInfo core.TokenInfo, gitTraceID uint32, project model.Project, projectServer model.ProjectServer) {
	srcPath := core.GolbalPath + "repository/" + project.Name + "/"
	remoteMachine := projectServer.ServerOwner + "@" + projectServer.ServerIP
	destPath := remoteMachine + ":" + project.Path
	remoteTraceModel := model.RemoteTrace{
		GitTraceID:    gitTraceID,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		ServerID:      projectServer.ServerID,
		ServerName:    projectServer.ServerName,
		PublisherID:   tokenInfo.ID,
		PublisherName: tokenInfo.Name,
		Type:          1,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}
	rsyncOption, err := utils.ParseCommandLine(project.RsyncOption)
	if err != nil {
		core.Log(core.ERROR, err.Error())
		remoteTraceModel.Detail = err.Error()
		remoteTraceModel.State = 0
		remoteTraceModel.AddRow()
	}
	rsyncOption = append(rsyncOption, "-e", "ssh -o StrictHostKeyChecking=no", srcPath, destPath)
	cmd := exec.Command("rsync", rsyncOption...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" rsync "+strings.Join(rsyncOption, " "))
	ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
		DataType: ws.RsyncType,
		State:    ws.Success,
		Message:  "rsync " + strings.Join(rsyncOption, " "),
	}
	var rsyncError error
	// 失败重试三次
	for attempt := 0; attempt < 3; attempt++ {
		rsyncError = cmd.Run()
		if rsyncError != nil {
			core.Log(core.ERROR, errbuf.String())
			ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.RsyncType,
				State:    ws.Fail,
				Message:  errbuf.String(),
			}
		} else {
			ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.RsyncType,
				State:    ws.Success,
				Message:  outbuf.String(),
			}
			remoteTraceModel.Detail = outbuf.String()
			remoteTraceModel.State = 1
			remoteTraceModel.AddRow()
			break
		}
	}

	if rsyncError != nil {
		ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
			DataType: ws.RsyncType,
			State:    ws.Fail,
			Message:  "rsync重试失败",
		}
		remoteTraceModel.Detail = errbuf.String()
		remoteTraceModel.State = 0
		remoteTraceModel.AddRow()
		return
	}
	// 没有脚本就不运行
	if project.Script == "" {
		return
	}
	remoteTraceModel.Type = 2
	// 执行ssh脚本
	ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
		DataType: ws.ScriptType,
		State:    ws.Success,
		Message:  "开始连接ssh",
	}
	var session *ssh.Session
	var connectError error
	for attempt := 0; attempt < 3; attempt++ {
		session, connectError = connect(projectServer.ServerOwner, "", projectServer.ServerIP, 22)
		if connectError != nil {
			core.Log(core.ERROR, connectError.Error())
			ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.ScriptType,
				State:    ws.Fail,
				Message:  connectError.Error(),
			}
		} else {
			ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.ScriptType,
				State:    ws.Success,
				Message:  "开始连接成功",
			}
			break
		}

	}

	if connectError != nil {
		core.Log(core.ERROR, connectError.Error())
		ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
			DataType: ws.ScriptType,
			State:    ws.Fail,
			Message:  "ssh重连失败",
		}
		remoteTraceModel.Detail = connectError.Error()
		remoteTraceModel.State = 0
		remoteTraceModel.AddRow()
		return
	}

	defer session.Close()
	var sshOutbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
		DataType: ws.ScriptType,
		State:    ws.Success,
		Message:  "运行:" + project.Script,
	}
	var scriptError error
	for attempt := 0; attempt < 3; attempt++ {
		sshOutbuf.Reset()
		if scriptError = session.Run(project.Script); scriptError != nil {
			core.Log(core.ERROR, scriptError.Error())
			ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.ScriptType,
				State:    ws.Fail,
				Message:  scriptError.Error(),
			}
		} else {
			ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.ScriptType,
				State:    ws.Success,
				Message:  sshOutbuf.String(),
			}
			break
		}
	}

	if scriptError != nil {
		core.Log(core.ERROR, scriptError.Error())
		ws.GetSyncHub().Broadcast <- &ws.SyncBroadcast{ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
			DataType: ws.ScriptType,
			State:    ws.Fail,
			Message:  "脚本运行失败",
		}
		remoteTraceModel.Detail = scriptError.Error()
		remoteTraceModel.State = 0
		remoteTraceModel.AddRow()
		return
	}

	remoteTraceModel.Detail = sshOutbuf.String()
	remoteTraceModel.State = 1
	remoteTraceModel.AddRow()
	return
}

func connect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)

	pemBytes, err := ioutil.ReadFile(os.Getenv("SSHKEY_PATH"))
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if password == "" {
		signer, err = ssh.ParsePrivateKey(pemBytes)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
	}
	if err != nil {
		return nil, err
	}
	auth = append(auth, ssh.PublicKeys(signer))

	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}
