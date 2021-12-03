package controller

import (
	"bytes"
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"io"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Server struct
type Server Controller

// GetList -
func (Server) GetList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

// GetTotal -
func (Server) GetTotal(gp *core.Goploy) *core.Response {
	total, err := model.Server{NamespaceID: gp.Namespace.ID}.GetTotal()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

// GetOption -
func (Server) GetOption(gp *core.Goploy) *core.Response {
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

// GetPublicKey -
func (Server) GetPublicKey(gp *core.Goploy) *core.Response {
	publicKeyPath := gp.URLQuery.Get("path")

	contentByte, err := ioutil.ReadFile(publicKeyPath + ".pub")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Key string `json:"key"`
		}{Key: string(contentByte)},
	}
}

// Check server
func (Server) Check(gp *core.Goploy) *core.Response {
	type ReqData struct {
		IP       string `json:"ip" validate:"required,ip|hostname"`
		Port     int    `json:"port" validate:"min=0,max=65535"`
		Owner    string `json:"owner" validate:"required,max=255"`
		Path     string `json:"path" validate:"required,max=255"`
		Password string `json:"password" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if Conn, err := utils.DialSSH(reqData.Owner, reqData.Password, reqData.Path, reqData.IP, reqData.Port); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	} else {
		_ = Conn.Close()
	}
	return &core.Response{Message: "Connected"}
}

// Add server
func (s Server) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		NamespaceID int64  `json:"namespaceId" validate:"gte=0"`
		IP          string `json:"ip" validate:"ip|hostname"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required,max=255"`
		Path        string `json:"path" validate:"required,max=255"`
		Password    string `json:"password"`
		Description string `json:"description" validate:"max=255"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.Server{
		NamespaceID: reqData.NamespaceID,
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		Path:        reqData.Path,
		Password:    reqData.Password,
		Description: reqData.Description,
		OSInfo:      s.getOSInfo(reqData.Owner, reqData.IP, strconv.Itoa(reqData.Port), reqData.Path),
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

// Edit server
func (s Server) Edit(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		NamespaceID int64  `json:"namespaceId" validate:"gte=0"`
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"required,ip|hostname"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required,max=255"`
		Path        string `json:"path" validate:"required,max=255"`
		Password    string `json:"password" validate:"max=255"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Server{
		ID:          reqData.ID,
		NamespaceID: reqData.NamespaceID,
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		Path:        reqData.Path,
		Password:    reqData.Password,
		Description: reqData.Description,
		OSInfo:      s.getOSInfo(reqData.Owner, reqData.IP, strconv.Itoa(reqData.Port), reqData.Path),
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Toggle server
func (Server) Toggle(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID    int64 `json:"id" validate:"gt=0"`
		State int8  `json:"state" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Server{ID: reqData.ID, State: reqData.State}).ToggleRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// DownloadFile sftp download file
func (Server) DownloadFile(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: "invalid server id"}
	}

	server, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	client, err := utils.DialSSH(server.Owner, server.Password, server.Path, server.IP, server.Port)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer client.Close()
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	filename := gp.URLQuery.Get("file")
	srcFile, _ := sftpClient.Open(filename) //远程
	FileStat, _ := srcFile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)
	gp.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+path.Base(filename))
	gp.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	gp.ResponseWriter.Header().Set("Content-Length", FileSize)
	_, _ = io.Copy(gp.ResponseWriter, srcFile)
	_ = srcFile.Close()
	_ = sftpClient.Close()
	return nil
}

// UploadFile sftp upload file
func (Server) UploadFile(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: "invalid server id"}
	}
	server, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	file, fileHandler, err := gp.Request.FormFile("file")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	client, err := utils.DialSSH(server.Owner, server.Password, server.Path, server.IP, server.Port)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	filePath := gp.URLQuery.Get("filePath")
	remoteFile, err := sftpClient.Create(filePath + "/" + fileHandler.Filename)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	_, err = remoteFile.Write(fileBytes)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return nil
}

func (Server) Report(gp *core.Goploy) *core.Response {
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: "invalid server id"}
	}
	logType, err := strconv.Atoi(gp.URLQuery.Get("type"))
	if err != nil {
		return &core.Response{Code: core.Error, Message: "invalid server id"}
	}
	datetimeRange := strings.Split(gp.URLQuery.Get("datetimeRange"), ",")
	if len(datetimeRange) != 2 {
		return &core.Response{Code: core.Error, Message: "invalid datetime range"}
	}
	serverAgentLogs, err := (model.ServerAgentLog{ServerID: serverID, Type: logType}).GetListBetweenTime(datetimeRange[0], datetimeRange[1])
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	type Flag struct {
		Count int
		Curr  int
	}

	flagMap := map[string]Flag{}

	for _, log := range serverAgentLogs {
		if _, ok := flagMap[log.Item]; !ok {
			flagMap[log.Item] = Flag{}
		}
		flagMap[log.Item] = Flag{Count: flagMap[log.Item].Count + 1}
	}

	serverAgentMap := map[string]model.ServerAgentLogs{}
	for _, log := range serverAgentLogs {
		flagMap[log.Item] = Flag{
			Count: flagMap[log.Item].Count,
			Curr:  flagMap[log.Item].Curr + 1,
		}
		step := flagMap[log.Item].Count / 60
		if flagMap[log.Item].Count <= 60 ||
			flagMap[log.Item].Curr%step == 0 ||
			flagMap[log.Item].Count-1 == flagMap[log.Item].Curr {
			serverAgentMap[log.Item] = append(serverAgentMap[log.Item], log)
		}
	}

	return &core.Response{
		Data: struct {
			ServerAgentMap map[string]model.ServerAgentLogs `json:"map"`
		}{ServerAgentMap: serverAgentMap},
	}
}

// GetAllMonitor -
func (Server) GetAllMonitor(gp *core.Goploy) *core.Response {
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	serverMonitorList, err := model.ServerMonitor{ServerID: serverID}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			List model.ServerMonitors `json:"list"`
		}{List: serverMonitorList},
	}
}

// AddMonitor server
func (s Server) AddMonitor(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ServerID     int64  `json:"serverId" validate:"required"`
		Item         string `json:"item" validate:"required"`
		Formula      string `json:"formula" validate:"required"`
		Operator     string `json:"operator" validate:"required"`
		Value        string `json:"value" validate:"required"`
		GroupCycle   int    `json:"groupCycle" validate:"required"`
		LastCycle    int    `json:"lastCycle" validate:"required"`
		SilentCycle  int    `json:"silentCycle" validate:"required"`
		StartTime    string `json:"startTime" validate:"required,len=5"`
		EndTime      string `json:"endTime" validate:"required,len=5"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.ServerMonitor{
		ServerID:     reqData.ServerID,
		Item:         reqData.Item,
		Formula:      reqData.Formula,
		Operator:     reqData.Operator,
		Value:        reqData.Value,
		GroupCycle:   reqData.GroupCycle,
		LastCycle:    reqData.LastCycle,
		SilentCycle:  reqData.SilentCycle,
		StartTime:    reqData.StartTime,
		EndTime:      reqData.EndTime,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
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

// EditMonitor server
func (s Server) EditMonitor(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"required"`
		Item         string `json:"item" validate:"required"`
		Formula      string `json:"formula" validate:"required"`
		Operator     string `json:"operator" validate:"required"`
		Value        string `json:"value" validate:"required"`
		GroupCycle   int    `json:"groupCycle" validate:"required"`
		LastCycle    int    `json:"lastCycle" validate:"required"`
		SilentCycle  int    `json:"silentCycle" validate:"required"`
		StartTime    string `json:"startTime" validate:"required,len=5"`
		EndTime      string `json:"endTime" validate:"required,len=5"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	err := model.ServerMonitor{
		ID:           reqData.ID,
		Item:         reqData.Item,
		Formula:      reqData.Formula,
		Operator:     reqData.Operator,
		Value:        reqData.Value,
		GroupCycle:   reqData.GroupCycle,
		LastCycle:    reqData.LastCycle,
		SilentCycle:  reqData.SilentCycle,
		StartTime:    reqData.StartTime,
		EndTime:      reqData.EndTime,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{}
}

// DeleteMonitor server
func (s Server) DeleteMonitor(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	err := model.ServerMonitor{
		ID: reqData.ID,
	}.DeleteRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{}
}

// version|cpu cores|mem
func (Server) getOSInfo(owner, ip, port, path string) string {
	osInfoScript := `cat /etc/os-release | grep "PRETTY_NAME" | awk -F\" '{print $2}' && cat /proc/cpuinfo  | grep "processor" | wc -l && cat /proc/meminfo | grep MemTotal | awk '{print $2}'`

	cmd := exec.Command("ssh",
		owner+"@"+ip,
		"-p", port,
		"-i", path,
		"-o", "StrictHostKeyChecking=no",
		osInfoScript)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return ""
	}
	// version|cpu cores|mem
	return strings.Replace(strings.Trim(out.String(), "\n"), "\n", "|", -1)
}
