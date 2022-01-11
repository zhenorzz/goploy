package controller

import (
	"bytes"
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/repository"
	"github.com/zhenorzz/goploy/response"
	"github.com/zhenorzz/goploy/utils"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Project Controller

func (p Project) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/project/getList", http.MethodGet, p.GetList),
		core.NewRoute("/project/getTotal", http.MethodGet, p.GetTotal),
		core.NewRoute("/project/pingRepos", http.MethodGet, p.PingRepos),
		core.NewRoute("/project/getRemoteBranchList", http.MethodGet, p.GetRemoteBranchList),
		core.NewRoute("/project/getBindServerList", http.MethodGet, p.GetBindServerList),
		core.NewRoute("/project/getBindUserList", http.MethodGet, p.GetBindUserList),
		core.NewRoute("/project/getProjectFileList", http.MethodGet, p.GetProjectFileList),
		core.NewRoute("/project/getProjectFileContent", http.MethodGet, p.GetProjectFileContent),
		core.NewRoute("/project/getReposFileList", http.MethodGet, p.GetReposFileList),
		core.NewRoute("/project/add", http.MethodPost, p.Add).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/edit", http.MethodPut, p.Edit).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/setAutoDeploy", http.MethodPut, p.SetAutoDeploy).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/remove", http.MethodDelete, p.Remove).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/uploadFile", http.MethodPost, p.UploadFile).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/removeFile", http.MethodDelete, p.RemoveFile).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/addServer", http.MethodPost, p.AddServer).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/addUser", http.MethodPost, p.AddUser).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/removeServer", http.MethodDelete, p.RemoveServer).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/removeUser", http.MethodDelete, p.RemoveUser).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/addFile", http.MethodPost, p.AddFile).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/editFile", http.MethodPut, p.EditFile).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/addTask", http.MethodPost, p.AddTask).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/removeTask", http.MethodDelete, p.RemoveTask).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/getTaskList", http.MethodGet, p.GetTaskList).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/project/getReviewList", http.MethodGet, p.GetReviewList),
		core.NewRoute("/project/getProcessList", http.MethodGet, p.GetProcessList).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/project/addProcess", http.MethodPost, p.AddProcess).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/project/editProcess", http.MethodPut, p.EditProcess).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/project/deleteProcess", http.MethodDelete, p.DeleteProcess).Roles(core.RoleAdmin, core.RoleManager),
	}
}

func (Project) GetList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectName := gp.URLQuery.Get("projectName")
	projectList, err := model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID, Name: projectName}.GetList(pagination)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Projects model.Projects `json:"list"`
		}{Projects: projectList},
	}
}

func (Project) GetTotal(gp *core.Goploy) core.Response {
	var total int64
	var err error
	projectName := gp.URLQuery.Get("projectName")
	total, err = model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID, Name: projectName}.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Project) PingRepos(gp *core.Goploy) core.Response {
	type ReqData struct {
		URL      string `schema:"url" validate:"required"`
		RepoType string `schema:"repoType" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if strings.Contains(reqData.URL, "git@") {
		host := strings.Split(reqData.URL, "git@")[1]
		host = strings.Split(host, ":")[0]
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		knownHostsPath := homeDir + "/.ssh/known_hosts"
		var cmdOutbuf, cmdErrbuf bytes.Buffer
		cmd := exec.Command("ssh-keygen", "-F", host, "-f", knownHostsPath)
		cmd.Stdout = &cmdOutbuf
		cmd.Stderr = &cmdErrbuf
		if err := cmd.Run(); err != nil {
			cmdOutbuf.Reset()
			cmdErrbuf.Reset()
			cmd := exec.Command("ssh-keyscan", host)
			cmd.Stdout = &cmdOutbuf
			cmd.Stderr = &cmdErrbuf
			if err := cmd.Run(); err != nil {
				return response.JSON{Code: response.Error, Message: cmdErrbuf.String()}
			}
			f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
			defer f.Close()
			if _, err := f.Write(cmdOutbuf.Bytes()); err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		}
	}

	repo, err := repository.GetRepo(reqData.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := repo.Ping(reqData.URL); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) GetRemoteBranchList(gp *core.Goploy) core.Response {
	type ReqData struct {
		URL      string `schema:"url" validate:"required"`
		RepoType string `schema:"repoType" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if strings.Contains(reqData.URL, "git@") {
		host := strings.Split(reqData.URL, "git@")[1]
		host = strings.Split(host, ":")[0]
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		knownHostsPath := homeDir + "/.ssh/known_hosts"
		var cmdOutbuf, cmdErrbuf bytes.Buffer
		cmd := exec.Command("ssh-keygen", "-F", host, "-f", knownHostsPath)
		cmd.Stdout = &cmdOutbuf
		cmd.Stderr = &cmdErrbuf
		if err := cmd.Run(); err != nil {
			cmdOutbuf.Reset()
			cmdErrbuf.Reset()
			cmd := exec.Command("ssh-keyscan", host)
			cmd.Stdout = &cmdOutbuf
			cmd.Stderr = &cmdErrbuf
			if err := cmd.Run(); err != nil {
				return response.JSON{Code: response.Error, Message: cmdErrbuf.String()}
			}
			f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
			defer f.Close()
			if _, err := f.Write(cmdOutbuf.Bytes()); err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		}
	}

	repo, err := repository.GetRepo(reqData.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := repo.RemoteBranchList(reqData.URL)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{Data: struct {
		Branch []string `json:"branch"`
	}{Branch: list}}
}

func (Project) GetBindServerList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectServers, err := model.ProjectServer{ProjectID: id}.GetBindServerListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectServers model.ProjectServers `json:"list"`
		}{ProjectServers: projectServers},
	}
}

func (Project) GetBindUserList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectUsers, err := model.ProjectUser{ProjectID: id, NamespaceID: gp.Namespace.ID}.GetBindUserListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectUsers model.ProjectUsers `json:"list"`
		}{ProjectUsers: projectUsers},
	}
}

func (Project) GetProjectFileList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectFiles, err := model.ProjectFile{ProjectID: id}.GetListByProjectID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectFiles model.ProjectFiles `json:"list"`
		}{ProjectFiles: projectFiles},
	}
}

func (Project) GetProjectFileContent(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectFileData, err := model.ProjectFile{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	fileBytes, err := ioutil.ReadFile(path.Join(core.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename))
	if err != nil {
		fmt.Println("read fail", err)
	}
	return response.JSON{
		Data: struct {
			Content string `json:"content"`
		}{Content: string(fileBytes)},
	}
}

func (Project) GetReposFileList(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID   int64  `schema:"id" validate:"gt=0"`
		Path string `schema:"path" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	files, err := ioutil.ReadDir(path.Join(core.GetProjectPath(reqData.ID), reqData.Path))
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type fileInfo struct {
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Mode    string `json:"mode"`
		ModTime string `json:"modTime"`
		IsDir   bool   `json:"isDir"`
	}
	var fileList []fileInfo
	for _, f := range files {
		fileList = append(fileList, fileInfo{
			Name:    f.Name(),
			Size:    f.Size(),
			Mode:    f.Mode().String(),
			ModTime: f.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:   f.IsDir(),
		})
	}

	return response.JSON{Data: fileList}
}

func (Project) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name                  string  `json:"name" validate:"required"`
		RepoType              string  `json:"repoType" validate:"required"`
		URL                   string  `json:"url" validate:"required"`
		Path                  string  `json:"path" validate:"required"`
		Environment           uint8   `json:"environment" validate:"required"`
		Branch                string  `json:"branch" validate:"required"`
		SymlinkPath           string  `json:"symlinkPath"`
		SymlinkBackupNumber   uint8   `json:"symlinkBackupNumber"`
		Review                uint8   `json:"review"`
		ReviewURL             string  `json:"reviewURL"`
		AfterPullScriptMode   string  `json:"afterPullScriptMode"`
		AfterPullScript       string  `json:"afterPullScript"`
		AfterDeployScriptMode string  `json:"afterDeployScriptMode"`
		AfterDeployScript     string  `json:"afterDeployScript"`
		RsyncOption           string  `json:"rsyncOption"`
		ServerIDs             []int64 `json:"serverIds"`
		UserIDs               []int64 `json:"userIds"`
		NotifyType            uint8   `json:"notifyType"`
		NotifyTarget          string  `json:"notifyTarget"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if _, err := utils.ParseCommandLine(reqData.RsyncOption); err != nil {
		return response.JSON{Code: response.Error, Message: "Invalid rsync option format"}
	}

	projectID, err := model.Project{
		NamespaceID:           gp.Namespace.ID,
		Name:                  reqData.Name,
		RepoType:              reqData.RepoType,
		URL:                   reqData.URL,
		Path:                  reqData.Path,
		Environment:           reqData.Environment,
		Branch:                reqData.Branch,
		SymlinkPath:           reqData.SymlinkPath,
		SymlinkBackupNumber:   reqData.SymlinkBackupNumber,
		Review:                reqData.Review,
		ReviewURL:             reqData.ReviewURL,
		AfterPullScriptMode:   reqData.AfterPullScriptMode,
		AfterPullScript:       reqData.AfterPullScript,
		AfterDeployScriptMode: reqData.AfterDeployScriptMode,
		AfterDeployScript:     reqData.AfterDeployScript,
		RsyncOption:           reqData.RsyncOption,
		NotifyType:            reqData.NotifyType,
		NotifyTarget:          reqData.NotifyTarget,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID: projectID,
			ServerID:  serverID,
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}

	if err := projectServersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectUsersModel := model.ProjectUsers{}
	for _, userID := range reqData.UserIDs {
		projectUserModel := model.ProjectUser{
			ProjectID: projectID,
			UserID:    userID,
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}

	namespaceUsers, err := model.NamespaceUser{NamespaceID: gp.Namespace.ID}.GetAllGteManagerByNamespaceID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	for _, namespaceUser := range namespaceUsers {
		projectUserModel := model.ProjectUser{
			ProjectID: projectID,
			UserID:    namespaceUser.UserID,
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}

	if err := projectUsersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID                    int64  `json:"id" validate:"gt=0"`
		Name                  string `json:"name"`
		RepoType              string `json:"repoType"`
		URL                   string `json:"url"`
		Path                  string `json:"path"`
		SymlinkPath           string `json:"symlinkPath"`
		SymlinkBackupNumber   uint8  `json:"symlinkBackupNumber"`
		Review                uint8  `json:"review"`
		ReviewURL             string `json:"reviewURL"`
		Environment           uint8  `json:"environment"`
		Branch                string `json:"branch"`
		AfterPullScriptMode   string `json:"afterPullScriptMode"`
		AfterPullScript       string `json:"afterPullScript"`
		AfterDeployScriptMode string `json:"afterDeployScriptMode"`
		AfterDeployScript     string `json:"afterDeployScript"`
		RsyncOption           string `json:"rsyncOption"`
		NotifyType            uint8  `json:"notifyType"`
		NotifyTarget          string `json:"notifyTarget"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if _, err := utils.ParseCommandLine(reqData.RsyncOption); err != nil {
		return response.JSON{Code: response.Error, Message: "Invalid rsync option format"}
	}

	projectData, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err = model.Project{
		ID:                    reqData.ID,
		Name:                  reqData.Name,
		RepoType:              reqData.RepoType,
		URL:                   reqData.URL,
		Path:                  reqData.Path,
		Environment:           reqData.Environment,
		Branch:                reqData.Branch,
		SymlinkPath:           reqData.SymlinkPath,
		SymlinkBackupNumber:   reqData.SymlinkBackupNumber,
		Review:                reqData.Review,
		ReviewURL:             reqData.ReviewURL,
		AfterPullScriptMode:   reqData.AfterPullScriptMode,
		AfterPullScript:       reqData.AfterPullScript,
		AfterDeployScriptMode: reqData.AfterDeployScriptMode,
		AfterDeployScript:     reqData.AfterDeployScript,
		RsyncOption:           reqData.RsyncOption,
		NotifyType:            reqData.NotifyType,
		NotifyTarget:          reqData.NotifyTarget,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if reqData.URL != projectData.URL {
		srcPath := core.GetProjectPath(projectData.ID)
		_, err := os.Stat(srcPath)
		if err == nil || os.IsNotExist(err) == false {
			repo := reqData.URL
			cmd := exec.Command("git", "remote", "set-url", "origin", repo)
			cmd.Dir = srcPath
			if err := cmd.Run(); err != nil {
				return response.JSON{Code: response.Error, Message: "Project change url fail, you can do it manually, reason: " + err.Error()}
			}
		}
	}

	return response.JSON{}
}

func (Project) SetAutoDeploy(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID         int64 `json:"id" validate:"gt=0"`
		AutoDeploy uint8 `json:"autoDeploy" validate:"gte=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.Project{
		ID:         reqData.ID,
		AutoDeploy: reqData.AutoDeploy,
	}.SetAutoDeploy()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) Remove(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectData, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	srcPath := core.GetProjectPath(projectData.ID)
	if err := os.RemoveAll(srcPath); err != nil {
		return response.JSON{Code: response.Error, Message: "Delete folder fail, Detail: " + err.Error()}
	}

	if err := (model.Project{ID: reqData.ID}).RemoveRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) UploadFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectFileID int64  `schema:"projectFileId" validate:"gt=0"`
		ProjectId     int64  `schema:"projectId"  validate:"gt=0"`
		Filename      string `schema:"filename"  validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	file, _, err := gp.Request.FormFile("file")
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer file.Close()

	filePath := path.Join(core.GetProjectFilePath(reqData.ProjectFileID), reqData.Filename)

	if _, err := os.Stat(path.Dir(filePath)); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path.Dir(filePath), 0755)
			if err != nil {
				return response.JSON{Code: response.Error, Message: err.Error()}
			}
		} else {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := ioutil.WriteFile(filePath, fileBytes, 0755); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if reqData.ProjectFileID == 0 {
		reqData.ProjectFileID, err = model.ProjectFile{
			Filename:  reqData.Filename,
			ProjectID: reqData.ProjectId,
		}.AddRow()
	} else {
		err = model.ProjectFile{
			ID:        reqData.ProjectFileID,
			Filename:  reqData.Filename,
			ProjectID: reqData.ProjectId,
		}.EditRow()
	}

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: reqData.ProjectFileID},
	}
}

func (Project) AddFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Content   string `json:"content" validate:"required"`
		Filename  string `json:"filename" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	_, err := os.Stat(core.GetProjectFilePath(reqData.ProjectID))
	if err != nil {
		err := os.MkdirAll(core.GetProjectFilePath(reqData.ProjectID), os.ModePerm)
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}
	file, err := os.Create(path.Join(core.GetProjectFilePath(reqData.ProjectID), reqData.Filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(reqData.Content)

	id, err := model.ProjectFile{
		ProjectID: reqData.ProjectID,
		Filename:  reqData.Filename,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Project) EditFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"gt=0"`
		Content string `json:"content" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectFileData, err := model.ProjectFile{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	_, err = os.Stat(core.GetProjectFilePath(projectFileData.ProjectID))
	if err != nil {
		err := os.MkdirAll(core.GetProjectFilePath(projectFileData.ProjectID), os.ModePerm)
		if err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	file, err := os.Create(path.Join(core.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(reqData.Content)

	return response.JSON{}
}

func (Project) RemoveFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectFileID int64 `json:"projectFileId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	projectFileData, err := model.ProjectFile{ID: reqData.ProjectFileID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := os.Remove(path.Join(core.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename)); err != nil {
		return response.JSON{Code: response.Error, Message: "Delete file fail, Detail: " + err.Error()}
	}

	if err := (model.ProjectFile{ID: reqData.ProjectFileID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) AddServer(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64   `json:"projectId" validate:"gt=0"`
		ServerIDs []int64 `json:"serverIds" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectID := reqData.ProjectID

	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID: projectID,
			ServerID:  serverID,
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}

	if err := projectServersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{}
}

func (Project) RemoveServer(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectServerID int64 `json:"projectServerId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectServer{ID: reqData.ProjectServerID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) AddUser(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64   `json:"projectId" validate:"gt=0"`
		UserIDs   []int64 `json:"userIds" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectID := reqData.ProjectID

	projectUsersModel := model.ProjectUsers{}
	for _, userID := range reqData.UserIDs {
		projectUserModel := model.ProjectUser{
			ProjectID: projectID,
			UserID:    userID,
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}

	if err := projectUsersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) RemoveUser(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectUserID int64 `json:"projectUserId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectUser{ID: reqData.ProjectUserID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) GetReviewList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	ProjectReviews, pagination, err := model.ProjectReview{ProjectID: id}.GetListByProjectID(pagination)

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectReviews model.ProjectReviews `json:"list"`
			Pagination     model.Pagination     `json:"pagination"`
		}{ProjectReviews: ProjectReviews, Pagination: pagination},
	}
}

func (Project) GetTaskList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	projectTaskList, pagination, err := model.ProjectTask{ProjectID: id}.GetListByProjectID(pagination)

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ProjectTasks model.ProjectTasks `json:"list"`
			Pagination   model.Pagination   `json:"pagination"`
		}{ProjectTasks: projectTaskList, Pagination: pagination},
	}
}

func (Project) AddTask(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Branch    string `json:"branch" validate:"required"`
		Commit    string `json:"commit" validate:"required"`
		Date      string `json:"date" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.ProjectTask{
		ProjectID: reqData.ProjectID,
		CommitID:  reqData.Commit,
		Branch:    reqData.Branch,
		Date:      reqData.Date,
		Creator:   gp.UserInfo.Name,
		CreatorID: gp.UserInfo.ID,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Project) RemoveTask(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectTask{ID: reqData.ID}).RemoveRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Project) GetProcessList(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `schema:"projectId" validate:"gt=0"`
		Page      uint64 `schema:"page" validate:"gt=0"`
		Rows      uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := model.ProjectProcess{ProjectID: reqData.ProjectID}.GetListByProjectID(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.ProjectProcesses `json:"list"`
		}{List: list},
	}
}

func (Project) AddProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Name      string `json:"name" validate:"required"`
		Status    string `json:"status"`
		Start     string `json:"start"`
		Stop      string `json:"stop"`
		Restart   string `json:"restart"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.ProjectProcess{
		ProjectID: reqData.ProjectID,
		Name:      reqData.Name,
		Status:    reqData.Status,
		Start:     reqData.Start,
		Stop:      reqData.Stop,
		Restart:   reqData.Restart,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Project) EditProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"gt=0"`
		Name    string `json:"name" validate:"required"`
		Status  string `json:"status"`
		Start   string `json:"start"`
		Stop    string `json:"stop"`
		Restart string `json:"restart"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.ProjectProcess{
		ID:      reqData.ID,
		Name:    reqData.Name,
		Status:  reqData.Status,
		Start:   reqData.Start,
		Stop:    reqData.Stop,
		Restart: reqData.Restart,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Project) DeleteProcess(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.ProjectProcess{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
