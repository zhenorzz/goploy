package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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
	type RespData struct {
		Project model.Projects `json:"projectList"`
	}
	groupID, err := strconv.ParseInt(gp.URLQuery.Get("groupId"), 10, 64)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: "groupId参数错误"}
		response.JSON(w)
		return
	}

	projectName := gp.URLQuery.Get("projectName")

	projects, err := model.ProjectUser{
		UserID: gp.TokenInfo.ID,
		Project: model.Project{
			GroupID: groupID,
			Name:    projectName,
		}}.GetDepolyListByUserID()
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

	projectServers, err := model.ProjectServer{ProjectID: reqData.ProjectID}.GetBindServerListByProjectID()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	project.PublisherID = gp.TokenInfo.ID
	project.PublisherName = gp.TokenInfo.Name
	project.LastPublishToken = uuid.New().String()
	project.UpdateTime = time.Now().Unix()
	_ = project.Publish()

	go execSync(gp.TokenInfo, project, projectServers)

	response := core.Response{Message: "部署中，请稍后"}
	response.JSON(w)
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

	projectServers, err := model.ProjectServer{ProjectID: reqData.ProjectID}.GetBindServerListByProjectID()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	project.PublisherID = gp.TokenInfo.ID
	project.PublisherName = gp.TokenInfo.Name
	project.LastPublishToken = uuid.New().String()
	project.UpdateTime = time.Now().Unix()
	_ = project.Publish()

	go execRollback(gp.TokenInfo, reqData.Commit, project, projectServers)

	response := core.Response{Message: "重新构建中，请稍后"}
	response.JSON(w)
}

func execSync(tokenInfo core.TokenInfo, project model.Project, projectServers model.ProjectServers) {
	publishTraceModel := model.PublishTrace{
		Token:         project.LastPublishToken,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   tokenInfo.ID,
		PublisherName: tokenInfo.Name,
		Type:          model.Pull,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	if err := gitCreate(tokenInfo, project); err != nil {
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		return
	}

	stdout, err := gitSync(tokenInfo, project)
	if err != nil {
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		return
	}

	commit, err := gitCommitID(tokenInfo, project)
	if err != nil {
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
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
		if err := runAfterPullScript(tokenInfo, project); err != nil {
			return
		}
	}

	for _, projectServer := range projectServers {
		go remoteSync(tokenInfo, project, projectServer)
	}
}

func execRollback(tokenInfo core.TokenInfo, commit string, project model.Project, projectServers model.ProjectServers) {
	publishTraceModel := model.PublishTrace{
		Token:         project.LastPublishToken,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   tokenInfo.ID,
		PublisherName: tokenInfo.Name,
		Type:          model.Pull,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}
	stdout, err := gitRollback(tokenInfo, commit, project)
	if err != nil {
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = 0
		publishTraceModel.AddRow()
		return
	}
	ext, _ := json.Marshal(struct {
		Commit string `json:"commit"`
	}{commit})
	publishTraceModel.Ext = string(ext)
	publishTraceModel.Detail = stdout
	publishTraceModel.State = 1

	if _, err := publishTraceModel.AddRow(); err != nil {
		core.Log(core.ERROR, err.Error())
	}

	if project.AfterPullScript != "" {
		if err := runAfterPullScript(tokenInfo, project); err != nil {
			return
		}
	}

	for _, projectServer := range projectServers {
		go remoteSync(tokenInfo, project, projectServer)
	}
}

func gitCreate(tokenInfo core.TokenInfo, project model.Project) error {
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
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.LocalType,
				State:     ws.Success,
				Message:   "项目初始化 git clone",
			},
		}
		if err := cmd.Run(); err != nil {
			core.Log(core.ERROR, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化失败:"+err.Error())
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.ProjectMessage{
					ProjectID: project.ID,
					UserID:    tokenInfo.ID,
					DataType:  ws.LocalType,
					State:     ws.Fail,
					Message:   "项目初始化失败",
				},
			}
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
				ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
					ToUserID: tokenInfo.ID,
					Message: ws.ProjectMessage{
						ProjectID: project.ID,
						UserID:    tokenInfo.ID,
						DataType:  ws.LocalType,
						State:     ws.Fail,
						Message:   checkoutErrbuf.String(),
					},
				}
				os.RemoveAll(srcPath)
				return errors.New(checkoutErrbuf.String())
			}
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.ProjectMessage{
					ProjectID: project.ID,
					UserID:    tokenInfo.ID,
					DataType:  ws.LocalType,
					State:     ws.Success,
					Message:   checkoutOutbuf.String(),
				},
			}
		}

		core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" 项目初始化成功")
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.LocalType,
				State:     ws.Success,
				Message:   "项目初始化成功",
			},
		}
	}
	return nil
}

func gitSync(tokenInfo core.TokenInfo, project model.Project) (string, error) {
	srcPath := core.RepositoryPath + project.Name

	clean := exec.Command("git", "clean", "-f")
	clean.Dir = srcPath
	var cleanOutbuf, cleanErrbuf bytes.Buffer
	clean.Stdout = &cleanOutbuf
	clean.Stderr = &cleanErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git clean -f")
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   "git clean -f",
		},
	}
	if err := clean.Run(); err != nil {
		core.Log(core.ERROR, cleanErrbuf.String())
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.LocalType,
				State:     ws.Fail,
				Message:   cleanErrbuf.String(),
			},
		}
		return "", errors.New(cleanErrbuf.String())
	}
	pull := exec.Command("git", "pull")
	pull.Dir = srcPath
	var pullOutbuf, pullErrbuf bytes.Buffer
	pull.Stdout = &pullOutbuf
	pull.Stderr = &pullErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git pull")
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   "git pull",
		},
	}
	if err := pull.Run(); err != nil {
		core.Log(core.ERROR, pullErrbuf.String())
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.LocalType,
				State:     ws.Fail,
				Message:   pullErrbuf.String(),
			},
		}
		return "", errors.New(pullErrbuf.String())
	}
	core.Log(core.TRACE, pullOutbuf.String())
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   pullOutbuf.String(),
		},
	}
	return pullOutbuf.String(), nil
}

func gitRollback(tokenInfo core.TokenInfo, commit string, project model.Project) (string, error) {
	srcPath := core.RepositoryPath + project.Name

	resetCmd := exec.Command("git", "reset", "--hard", commit)
	resetCmd.Dir = srcPath
	var resetOutbuf, resetErrbuf bytes.Buffer
	resetCmd.Stdout = &resetOutbuf
	resetCmd.Stderr = &resetErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git reset --hard "+commit)
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   "git reset --hard " + commit,
		},
	}
	if err := resetCmd.Run(); err != nil {
		core.Log(core.ERROR, resetErrbuf.String())
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.LocalType,
				State:     ws.Fail,
				Message:   resetErrbuf.String(),
			},
		}
		return "", errors.New(resetErrbuf.String())
	}

	core.Log(core.TRACE, resetOutbuf.String())
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   resetOutbuf.String(),
		},
	}
	return resetOutbuf.String(), nil
}

func gitCommitID(tokenInfo core.TokenInfo, project model.Project) (string, error) {
	srcPath := core.RepositoryPath + project.Name

	git := exec.Command("git", "rev-parse", "HEAD")
	git.Dir = srcPath
	var gitOutbuf, gitErrbuf bytes.Buffer
	git.Stdout = &gitOutbuf
	git.Stderr = &gitErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git rev-parse HEAD")
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   "git rev-parse HEAD",
		},
	}
	if err := git.Run(); err != nil {
		core.Log(core.ERROR, gitErrbuf.String())
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.LocalType,
				State:     ws.Success,
				Message:   gitErrbuf.String(),
			},
		}
		return "", errors.New(gitErrbuf.String())
	}
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   "commitSHA: " + gitOutbuf.String(),
		},
	}
	return gitOutbuf.String(), nil
}

func runAfterPullScript(tokenInfo core.TokenInfo, project model.Project) error {
	publishTraceModel := model.PublishTrace{
		Token:         project.LastPublishToken,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   tokenInfo.ID,
		PublisherName: tokenInfo.Name,
		Type:          model.AfterPull,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}
	ext, _ := json.Marshal(struct {
		Script string `json:"script"`
	}{project.AfterPullScript})
	publishTraceModel.Ext = string(ext)
	srcPath := core.RepositoryPath + project.Name
	scriptName := srcPath + "/after-pull.sh"
	ioutil.WriteFile(scriptName, []byte(project.AfterPullScript), 0755)
	handler := exec.Command("bash", "./after-pull.sh")
	handler.Dir = srcPath
	var outbuf, errbuf bytes.Buffer
	handler.Stdout = &outbuf
	handler.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+project.AfterPullScript)
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   project.AfterPullScript,
		},
	}
	if err := handler.Run(); err != nil {
		core.Log(core.ERROR, err.Error())
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID,
				UserID:    tokenInfo.ID,
				DataType:  ws.LocalType,
				State:     ws.Success,
				Message:   errbuf.String(),
			},
		}

		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		if _, err := publishTraceModel.AddRow(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
		return err
	}
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID,
			UserID:    tokenInfo.ID,
			DataType:  ws.LocalType,
			State:     ws.Success,
			Message:   outbuf.String(),
		},
	}
	publishTraceModel.Detail = outbuf.String()
	publishTraceModel.State = model.Success
	os.Remove(scriptName)
	if _, err := publishTraceModel.AddRow(); err != nil {
		core.Log(core.ERROR, err.Error())
	}

	return nil
}

func remoteSync(tokenInfo core.TokenInfo, project model.Project, projectServer model.ProjectServer) {
	srcPath := core.RepositoryPath + project.Name + "/"
	remoteMachine := projectServer.ServerOwner + "@" + projectServer.ServerIP
	destPath := remoteMachine + ":" + project.Path
	publishTraceModel := model.PublishTrace{
		Token:         project.LastPublishToken,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		PublisherID:   tokenInfo.ID,
		PublisherName: tokenInfo.Name,
		Type:          model.Deploy,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	ext, _ := json.Marshal(struct {
		ServerID   int64  `json:"serverId"`
		ServerName string `json:"serverName"`
	}{projectServer.ServerID, projectServer.ServerName})
	publishTraceModel.Ext = string(ext)

	rsyncOption, err := utils.ParseCommandLine(project.RsyncOption)
	if err != nil {
		core.Log(core.ERROR, err.Error())
		publishTraceModel.Detail = err.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
	}
	rsyncOption = append(rsyncOption, "-e", "ssh -p "+strconv.Itoa(int(projectServer.ServerPort))+" -o StrictHostKeyChecking=no", srcPath, destPath)
	cmd := exec.Command("rsync", rsyncOption...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" rsync "+strings.Join(rsyncOption, " "))
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
			DataType: ws.RsyncType,
			State:    ws.Success,
			Message:  "rsync " + strings.Join(rsyncOption, " "),
		},
	}
	var rsyncError error
	// 失败重试三次
	for attempt := 0; attempt < 3; attempt++ {
		rsyncError = cmd.Run()
		if rsyncError != nil {
			core.Log(core.ERROR, errbuf.String())
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.ProjectMessage{
					ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
					DataType: ws.RsyncType,
					State:    ws.Fail,
					Message:  errbuf.String(),
				},
			}
		} else {
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.ProjectMessage{
					ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
					DataType: ws.RsyncType,
					State:    ws.Success,
					Message:  outbuf.String(),
				},
			}
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
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.RsyncType,
				State:    ws.Fail,
				Message:  "rsync重试失败",
			},
		}
		publishTraceModel.Detail = errbuf.String()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		return
	}
	// 没有脚本就不运行
	if project.AfterDeployScript == "" {
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
	ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
		ToUserID: tokenInfo.ID,
		Message: ws.ProjectMessage{
			ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
			DataType: ws.ScriptType,
			State:    ws.Success,
			Message:  "开始连接ssh",
		},
	}
	var session *ssh.Session
	var connectError error
	var scriptError error
	for attempt := 0; attempt < 3; attempt++ {
		session, connectError = connect(projectServer.ServerOwner, "", projectServer.ServerIP, int(projectServer.ServerPort))
		if connectError != nil {
			core.Log(core.ERROR, connectError.Error())
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.ProjectMessage{
					ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
					DataType: ws.ScriptType,
					State:    ws.Fail,
					Message:  connectError.Error(),
				},
			}
		} else {
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.ProjectMessage{
					ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
					DataType: ws.ScriptType,
					State:    ws.Success,
					Message:  "开始连接成功",
				},
			}
			ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
				ToUserID: tokenInfo.ID,
				Message: ws.ProjectMessage{
					ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
					DataType: ws.ScriptType,
					State:    ws.Success,
					Message:  "运行:" + project.AfterDeployScript,
				},
			}
			defer session.Close()
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			sshOutbuf.Reset()
			afterDeployScript := "echo '" + project.AfterDeployScript + "' > /tmp/after-deploy.sh;bash /tmp/after-deploy.sh"
			if scriptError = session.Run(afterDeployScript); scriptError != nil {
				core.Log(core.ERROR, scriptError.Error())
				ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
					ToUserID: tokenInfo.ID,
					Message: ws.ProjectMessage{
						ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
						DataType: ws.ScriptType,
						State:    ws.Fail,
						Message:  scriptError.Error(),
					},
				}
			} else {
				ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
					ToUserID: tokenInfo.ID,
					Message: ws.ProjectMessage{
						ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
						DataType: ws.ScriptType,
						State:    ws.Success,
						Message:  sshOutbuf.String(),
					},
				}
				publishTraceModel.Detail = sshOutbuf.String()
				publishTraceModel.State = model.Success
				publishTraceModel.AddRow()
				break
			}
		}

	}

	if connectError != nil {
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.ScriptType,
				State:    ws.Fail,
				Message:  "ssh重连失败",
			},
		}
		publishTraceModel.Detail = connectError.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		return
	} else if scriptError != nil {
		ws.GetUnicastHub().UnicastData <- &ws.UnicastData{
			ToUserID: tokenInfo.ID,
			Message: ws.ProjectMessage{
				ProjectID: project.ID, UserID: tokenInfo.ID, ServerID: projectServer.ServerID, ServerName: projectServer.ServerName,
				DataType: ws.ScriptType,
				State:    ws.Fail,
				Message:  "脚本运行失败",
			},
		}
		publishTraceModel.Detail = scriptError.Error()
		publishTraceModel.State = model.Fail
		publishTraceModel.AddRow()
		return
	}
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
