package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Deploy struct
type Deploy struct{}

// Get deploy list
func (deploy *Deploy) Get(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Project model.Projects `json:"projectList"`
	}

	model := model.Projects{}
	err := model.QueryByStatus(2)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	response := core.Response{Data: RepData{Project: model}}
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

	projectModel := model.Project{
		ID: reqData.ID,
	}

	if err := projectModel.QueryRow(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if projectModel.Status != 2 {
		response := core.Response{Code: 1, Message: "项目尚未初始化"}
		response.Json(w)
		return
	}

	projectServersModel := model.ProjectServers{}

	if err := projectServersModel.Query(reqData.ID); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	gitTraceModel := model.GitTrace{
		ProjectID:     projectModel.ID,
		ProjectName:   projectModel.Name,
		PublisherID:   core.GolbalUserID,
		PublisherName: core.GolbalUserName,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	stdout, err := gitSync(projectModel)
	if err != nil {
		gitTraceModel.Detail = err.Error()
		gitTraceModel.State = 0
		_ = gitTraceModel.AddRow()
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	} else {
		gitTraceModel.Detail = stdout
		gitTraceModel.State = 1
		_ = gitTraceModel.AddRow()
	}

	for _, projectServer := range projectServersModel {
		go rsync(gitTraceModel.ID, projectModel, projectServer)
	}

	projectModel.PublisherID = core.GolbalUserID
	projectModel.PublisherName = core.GolbalUserName
	projectModel.UpdateTime = time.Now().Unix()
	_ = projectModel.Publish()

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

func rsync(gitTraceID uint32, project model.Project, projectServer model.ProjectServer) {
	srcPath := "./repository/" + project.Name
	destPath := projectServer.ServerOwner + "@" + projectServer.ServerIP + ":" + project.Path
	cmd := exec.Command("rsync", "-rtv", "--delete", srcPath, destPath)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(project.ID), 10)+" rsync -rtv --delete "+srcPath+destPath)
	rsyncTraceModel := model.RsyncTrace{
		GitTraceID:    gitTraceID,
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		ServerID:      projectServer.ServerID,
		ServerName:    projectServer.ServerName,
		PublisherID:   core.GolbalUserID,
		PublisherName: core.GolbalUserName,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, errbuf.String())
		rsyncTraceModel.Detail = errbuf.String()
		rsyncTraceModel.State = 0
		_ = rsyncTraceModel.AddRow()
	} else {
		rsyncTraceModel.Detail = outbuf.String()
		rsyncTraceModel.State = 1
		_ = rsyncTraceModel.AddRow()
	}
}
