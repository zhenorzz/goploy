package controller

import (
	"bytes"
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Project struct
type Project Controller

// GetList -
func (Project) GetList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectName := gp.URLQuery.Get("projectName")
	projectList, err := model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID, Name: projectName}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Projects model.Projects `json:"list"`
		}{Projects: projectList},
	}
}

// GetTotal -
func (Project) GetTotal(gp *core.Goploy) *core.Response {
	var total int64
	var err error
	projectName := gp.URLQuery.Get("projectName")
	total, err = model.Project{NamespaceID: gp.Namespace.ID, UserID: gp.UserInfo.ID, Name: projectName}.GetTotal()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

// GetRemoteBranchList -
func (Project) GetRemoteBranchList(gp *core.Goploy) *core.Response {
	url := gp.URLQuery.Get("url")

	if strings.Contains(url, "git@") {
		host := strings.Split(url, "git@")[1]
		host = strings.Split(host, ":")[0]
		knownHostsPath := os.ExpandEnv("$HOME/.ssh/known_hosts")
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
				return &core.Response{Code: core.Error, Message: cmdErrbuf.String()}
			}
			f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return &core.Response{Code: core.Error, Message: err.Error()}
			}
			defer f.Close()
			if _, err := f.Write(cmdOutbuf.Bytes()); err != nil {
				return &core.Response{Code: core.Error, Message: err.Error()}
			}
		}
	}

	cmd := exec.Command("git", "ls-remote", "-h", url)
	var cmdOutbuf, cmdErrbuf bytes.Buffer
	cmd.Stdout = &cmdOutbuf
	cmd.Stderr = &cmdErrbuf
	if err := cmd.Run(); err != nil {
		return &core.Response{Code: core.Error, Message: cmdErrbuf.String()}
	}
	var branch []string
	for _, branchWithSha := range strings.Split(cmdOutbuf.String(), "\n") {
		if len(branchWithSha) != 0 {
			branchWithShaSlice := strings.Fields(branchWithSha)
			branchWithHead := branchWithShaSlice[len(branchWithShaSlice)-1]
			branchWithHeadSlice := strings.Split(branchWithHead, "/")
			branch = append(branch, branchWithHeadSlice[len(branchWithHeadSlice)-1])
		}
	}

	return &core.Response{Data: struct {
		Branch []string `json:"branch"`
	}{Branch: branch}}
}

// GetBindServerList -
func (Project) GetBindServerList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectServers, err := model.ProjectServer{ProjectID: id}.GetBindServerListByProjectID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			ProjectServers model.ProjectServers `json:"list"`
		}{ProjectServers: projectServers},
	}
}

// GetBindUserList -
func (Project) GetBindUserList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectUsers, err := model.ProjectUser{ProjectID: id, NamespaceID: gp.Namespace.ID}.GetBindUserListByProjectID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			ProjectUsers model.ProjectUsers `json:"list"`
		}{ProjectUsers: projectUsers},
	}
}

// GetProjectFileList -
func (Project) GetProjectFileList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectFiles, err := model.ProjectFile{ProjectID: id}.GetListByProjectID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			ProjectFiles model.ProjectFiles `json:"list"`
		}{ProjectFiles: projectFiles},
	}
}

// GetProjectFileContent -
func (Project) GetProjectFileContent(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectFileData, err := model.ProjectFile{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	fileBytes, err := ioutil.ReadFile(path.Join(core.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename))
	if err != nil {
		fmt.Println("read fail", err)
	}
	return &core.Response{
		Data: struct {
			Content string `json:"content"`
		}{Content: string(fileBytes)},
	}
}

// Add project
func (Project) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name                  string  `json:"name" validate:"required"`
		URL                   string  `json:"url" validate:"required"`
		Path                  string  `json:"path" validate:"required"`
		Environment           uint8   `json:"Environment" validate:"required"`
		Branch                string  `json:"branch" validate:"required"`
		SymlinkPath           string  `json:"symlinkPath"`
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
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if _, err := utils.ParseCommandLine(reqData.RsyncOption); err != nil {
		return &core.Response{Code: core.Error, Message: "Invalid rsync option format"}
	}

	projectID, err := model.Project{
		NamespaceID:           gp.Namespace.ID,
		Name:                  reqData.Name,
		URL:                   reqData.URL,
		Path:                  reqData.Path,
		Environment:           reqData.Environment,
		Branch:                reqData.Branch,
		SymlinkPath:           reqData.SymlinkPath,
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
		return &core.Response{Code: core.Error, Message: err.Error()}
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
		return &core.Response{Code: core.Error, Message: err.Error()}
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
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	for _, namespaceUser := range namespaceUsers {
		projectUserModel := model.ProjectUser{
			ProjectID: projectID,
			UserID:    namespaceUser.UserID,
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}

	if err := projectUsersModel.AddMany(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Edit project
func (Project) Edit(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID                    int64  `json:"id" validate:"gt=0"`
		Name                  string `json:"name"`
		URL                   string `json:"url"`
		Path                  string `json:"path"`
		SymlinkPath           string `json:"symlinkPath"`
		Review                uint8  `json:"review"`
		ReviewURL             string `json:"reviewURL"`
		Environment           uint8  `json:"Environment"`
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
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if _, err := utils.ParseCommandLine(reqData.RsyncOption); err != nil {
		return &core.Response{Code: core.Error, Message: "Invalid rsync option format"}
	}

	projectData, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	err = model.Project{
		ID:                    reqData.ID,
		Name:                  reqData.Name,
		URL:                   reqData.URL,
		Path:                  reqData.Path,
		Environment:           reqData.Environment,
		Branch:                reqData.Branch,
		SymlinkPath:           reqData.SymlinkPath,
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
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if reqData.URL != projectData.URL {
		srcPath := core.GetProjectPath(projectData.ID)
		_, err := os.Stat(srcPath)
		if err == nil || os.IsNotExist(err) == false {
			repo := reqData.URL
			cmd := exec.Command("git", "remote", "set-url", "origin", repo)
			cmd.Dir = srcPath
			if err := cmd.Run(); err != nil {
				return &core.Response{Code: core.Error, Message: "Project change url fail, you can do it manually, reason: " + err.Error()}
			}
		}
	}

	return &core.Response{}
}

// SetAutoDeploy -
func (Project) SetAutoDeploy(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID         int64 `json:"id" validate:"gt=0"`
		AutoDeploy uint8 `json:"autoDeploy" validate:"gte=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Project{
		ID:         reqData.ID,
		AutoDeploy: reqData.AutoDeploy,
	}.SetAutoDeploy()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// RemoveRow Project
func (Project) Remove(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	projectData, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	srcPath := core.GetProjectPath(projectData.ID)
	if err := os.RemoveAll(srcPath); err != nil {
		return &core.Response{Code: core.Error, Message: "Delete folder fail, Detail: " + err.Error()}
	}

	if err := (model.Project{ID: reqData.ID}).RemoveRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{}
}

// UploadFile -
func (Project) UploadFile(gp *core.Goploy) *core.Response {
	file, _, err := gp.Request.FormFile("file")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer file.Close()
	id, err := strconv.ParseInt(gp.URLQuery.Get("projectFileId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectID, err := strconv.ParseInt(gp.URLQuery.Get("projectId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	filename := gp.URLQuery.Get("filename")
	filePath := path.Join(core.GetProjectFilePath(projectID), filename)

	if _, err := os.Stat(path.Dir(filePath)); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path.Dir(filePath), 0755)
			if err != nil {
				return &core.Response{Code: core.Error, Message: err.Error()}
			}
		} else {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := ioutil.WriteFile(filePath, fileBytes, 0755); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if id == 0 {
		id, err = model.ProjectFile{
			Filename:  filename,
			ProjectID: projectID,
		}.AddRow()
	} else {
		err = model.ProjectFile{
			ID:        id,
			Filename:  filename,
			ProjectID: projectID,
		}.EditRow()
	}

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

// AddFile to project
func (Project) AddFile(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Content   string `json:"content" validate:"required"`
		Filename  string `json:"filename" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
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
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

// EditFile to project
func (Project) EditFile(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"gt=0"`
		Content string `json:"content" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	projectFileData, err := model.ProjectFile{ID: reqData.ID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	file, err := os.Create(path.Join(core.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(reqData.Content)

	return &core.Response{}
}

// RemoveFile from Project
func (Project) RemoveFile(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectFileID int64 `json:"projectFileId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	projectFileData, err := model.ProjectFile{ID: reqData.ProjectFileID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := os.Remove(path.Join(core.GetProjectFilePath(projectFileData.ProjectID), projectFileData.Filename)); err != nil {
		return &core.Response{Code: core.Error, Message: "Delete file fail, Detail: " + err.Error()}
	}

	if err := (model.ProjectFile{ID: reqData.ProjectFileID}).DeleteRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{}
}

// AddServer to project
func (Project) AddServer(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64   `json:"projectId" validate:"gt=0"`
		ServerIDs []int64 `json:"serverIds" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
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
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{}
}

// RemoveServer from Project
func (Project) RemoveServer(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectServerID int64 `json:"projectServerId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.ProjectServer{ID: reqData.ProjectServerID}).DeleteRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// AddUser to project
func (Project) AddUser(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64   `json:"projectId" validate:"gt=0"`
		UserIDs   []int64 `json:"userIds" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
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
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// RemoveUser from Project
func (Project) RemoveUser(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectUserID int64 `json:"projectUserId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}

	if err := (model.ProjectUser{ID: reqData.ProjectUserID}).DeleteRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// GetTaskList -
func (Project) GetTaskList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectTaskList, pagination, err := model.ProjectTask{ProjectID: id}.GetListByProjectID(pagination)

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			ProjectTasks model.ProjectTasks `json:"list"`
			Pagination   model.Pagination   `json:"pagination"`
		}{ProjectTasks: projectTaskList, Pagination: pagination},
	}
}

// GetReviewList -
func (Project) GetReviewList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	ProjectReviews, pagination, err := model.ProjectReview{ProjectID: id}.GetListByProjectID(pagination)

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			ProjectReviews model.ProjectReviews `json:"list"`
			Pagination     model.Pagination     `json:"pagination"`
		}{ProjectReviews: ProjectReviews, Pagination: pagination},
	}
}

// AddTask to project
func (Project) AddTask(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Branch    string `json:"branch" validate:"required"`
		CommitID  string `json:"commitId" validate:"len=40"`
		Date      string `json:"date" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.ProjectTask{
		ProjectID: reqData.ProjectID,
		CommitID:  reqData.CommitID,
		Branch:    reqData.Branch,
		Date:      reqData.Date,
		Creator:   gp.UserInfo.Name,
		CreatorID: gp.UserInfo.ID,
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

// RemoveTask from project
func (Project) RemoveTask(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.ProjectTask{ID: reqData.ID}).RemoveRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{}
}
