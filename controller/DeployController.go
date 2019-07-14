package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"golang.org/x/crypto/ssh"
)

// Deploy struct
type Deploy struct{}

// Get deploy list
func (deploy *Deploy) GetList(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Project model.Projects `json:"projectList"`
	}

	projects, err := model.Project{Status: 2}.GetDepolyList()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	response := core.Response{Data: RepData{Project: projects}}
	response.Json(w)
}

// GetDetail deploy detail
func (deploy *Deploy) GetDetail(w http.ResponseWriter, r *http.Request) {

	type RepData struct {
		GitTrace        model.GitTrace     `json:"gitTrace"`
		RemoteTraceList model.RemoteTraces `json:"remoteTraceList"`
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response := core.Response{Code: 1, Message: "id参数错误"}
		response.Json(w)
		return
	}

	gitTraceModel := model.GitTrace{}

	if err := gitTraceModel.QueryLatestRow(uint32(id)); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	remoteTracesList, err := model.RemoteTrace{}.GetListByGitTraceID(gitTraceModel.ID)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	response := core.Response{Data: RepData{GitTrace: gitTraceModel, RemoteTraceList: remoteTracesList}}
	response.Json(w)
}

// Publish the project
func (deploy *Deploy) Publish(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(body, &reqData); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	project, err := model.Project{
		ID: reqData.ID,
	}.GetData()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if project.Status != 2 {
		response := core.Response{Code: 1, Message: "项目尚未初始化"}
		response.Json(w)
		return
	}

	projectServers, err := model.ProjectServer{}.GetBindServerListByProjectID(reqData.ID)

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	gitTraceModel := model.GitTrace{
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   core.GolbalUserID,
		PublisherName: core.GolbalUserName,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	stdout, err := gitSync(project)
	if err != nil {
		gitTraceModel.Detail = err.Error()
		gitTraceModel.State = 0
		_, _ = gitTraceModel.AddRow()
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	gitTraceModel.Detail = stdout
	gitTraceModel.State = 1
	gitTraceID, _ := gitTraceModel.AddRow()

	for _, projectServer := range projectServers {
		go remoteExec(gitTraceID, project, projectServer)
	}

	project.PublisherID = core.GolbalUserID
	project.PublisherName = core.GolbalUserName
	project.UpdateTime = time.Now().Unix()
	_ = project.Publish()

	response := core.Response{Message: "部署中，请稍后"}
	response.Json(w)
}

func gitSync(project model.Project) (string, error) {
	srcPath := "./repository/" + project.Name
	clean := exec.Command("git", "clean", "-f")
	clean.Dir = srcPath
	var cleanOutbuf, cleanErrbuf bytes.Buffer
	clean.Stdout = &cleanOutbuf
	clean.Stderr = &cleanErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" git clean -f")
	if err := clean.Run(); err != nil {
		core.Log(core.ERROR, cleanErrbuf.String())
		return "", err
	}
	pull := exec.Command("git", "pull")
	pull.Dir = srcPath
	var pullOutbuf, pullErrbuf bytes.Buffer
	pull.Stdout = &pullOutbuf
	pull.Stderr = &pullErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" git pull")
	if err := pull.Run(); err != nil {
		core.Log(core.ERROR, pullErrbuf.String())
		return "", err
	}

	core.Log(core.TRACE, pullOutbuf.String())
	return pullOutbuf.String(), nil
}

func remoteExec(gitTraceID uint32, project model.Project, projectServer model.ProjectServer) {
	srcPath := "./repository/" + project.Name
	remoteMachine := projectServer.ServerOwner + "@" + projectServer.ServerIP
	destPath := remoteMachine + ":" + project.Path
	cmd := exec.Command("rsync", "-rtv", "--exclude", ".git", "--delete", srcPath, destPath)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" rsync -rtv --delete "+srcPath+destPath)
	remoteTraceModel := model.RemoteTrace{
		GitTraceID:    gitTraceID,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		ServerID:      projectServer.ServerID,
		ServerName:    projectServer.ServerName,
		PublisherID:   core.GolbalUserID,
		PublisherName: core.GolbalUserName,
		Type:          1,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, errbuf.String())
		remoteTraceModel.Detail = errbuf.String()
		remoteTraceModel.State = 0
		remoteTraceModel.AddRow()
	} else {
		remoteTraceModel.Detail = outbuf.String()
		remoteTraceModel.State = 1
		remoteTraceModel.AddRow()
	}
	// 没有脚本就不运行
	if project.Script == "" {
		return
	}
	remoteTraceModel.Type = 2
	// 执行ssh脚本
	session, err := connect(projectServer.ServerOwner, "", projectServer.ServerIP, 22)
	if err != nil {
		core.Log(core.ERROR, err.Error())
		remoteTraceModel.Detail = err.Error()
		remoteTraceModel.State = 0
		remoteTraceModel.AddRow()
		return
	}
	defer session.Close()
	var sshOutbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	// 需要更改脚本的权限 以免执行不了
	if err := session.Run("chmod u+x " + project.Script + "&&" + project.Script); err != nil {
		core.Log(core.ERROR, err.Error())
		remoteTraceModel.Detail = err.Error()
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
