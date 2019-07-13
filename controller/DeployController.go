package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Deploy struct
type Deploy struct{}

// Get deploy list
func (deploy *Deploy) Get(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Deploys model.Deploys `json:"deployList"`
	}

	model := model.Deploys{}
	err := model.Query()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Deploys: model}}
	response.Json(w)
}

// Publish the project
func (deploy *Deploy) Publish(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	projectID := reqData.ID
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

	srcPath := "./repository/" + projectModel.Name

	if err := gitSync(projectID, srcPath); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	for _, projectServer := range projectServersModel {
		destPath := projectServer.ServerOwner + "@" + projectServer.ServerIP + ":" + projectModel.Path
		go rsync(projectID, srcPath, destPath)
	}

	// deployModel.Status = 1
	// _ = deployModel.ChangeStatus()
	response := core.Response{Message: "部署中，请稍后"}
	response.Json(w)
}

func gitSync(projectID uint32, srcPath string) error {
	clean := exec.Command("git", "clean", "-f")
	clean.Dir = srcPath
	var cleanOutbuf, cleanErrbuf bytes.Buffer
	clean.Stdout = &cleanOutbuf
	clean.Stderr = &cleanErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(projectID), 10)+" git clean -f")
	if err := clean.Run(); err != nil {
		// deployModel.Status = 3
		// _ = deployModel.ChangeStatus()
		core.Log(core.ERROR, cleanErrbuf.String())
		return err
	}
	pull := exec.Command("git", "pull")
	pull.Dir = srcPath
	var pullOutbuf, pullErrbuf bytes.Buffer
	pull.Stdout = &pullOutbuf
	pull.Stderr = &pullErrbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(projectID), 10)+" git pull")
	if err := pull.Run(); err != nil {
		// deployModel.Status = 3
		// _ = deployModel.ChangeStatus()
		core.Log(core.ERROR, pullErrbuf.String())
		return err
	}
	return nil
}

func rsync(projectID uint32, srcPath, destPath string) {
	cmd := exec.Command("rsync", "-rtv", "--delete", srcPath, destPath)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	core.Log(core.TRACE, "projectID:"+strconv.FormatUint(uint64(projectID), 10)+" rsync -rtv --delete "+srcPath+destPath)
	if err := cmd.Run(); err != nil {
		// deployModel.Status = 3
		// _ = deployModel.ChangeStatus()
		core.Log(core.ERROR, errbuf.String())
		return
	}
}
