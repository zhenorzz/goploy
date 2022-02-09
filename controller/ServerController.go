package controller

import (
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"github.com/zhenorzz/goploy/utils"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// Server struct
type Server Controller

func (s Server) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/server/getList", http.MethodGet, s.GetList),
		core.NewRoute("/server/getTotal", http.MethodGet, s.GetTotal),
		core.NewRoute("/server/getOption", http.MethodGet, s.GetOption),
		core.NewRoute("/server/getPublicKey", http.MethodGet, s.GetPublicKey),
		core.NewRoute("/server/getTerminalRecord", http.MethodGet, s.GetTerminalRecord),
		core.NewRoute("/server/check", http.MethodPost, s.Check).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/add", http.MethodPost, s.Add).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/edit", http.MethodPut, s.Edit).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/toggle", http.MethodPut, s.Toggle).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/previewFile", http.MethodGet, s.PreviewFile).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/downloadFile", http.MethodGet, s.DownloadFile).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/uploadFile", http.MethodPost, s.UploadFile).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/report", http.MethodGet, s.Report).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/getAllMonitor", http.MethodGet, s.GetAllMonitor),
		core.NewRoute("/server/addMonitor", http.MethodPost, s.AddMonitor).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/editMonitor", http.MethodPut, s.EditMonitor).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/server/deleteMonitor", http.MethodDelete, s.DeleteMonitor).Roles(core.RoleAdmin, core.RoleManager),
	}
}

func (Server) GetList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetList(pagination)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

func (Server) GetTotal(gp *core.Goploy) core.Response {
	total, err := model.Server{NamespaceID: gp.Namespace.ID}.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Server) GetOption(gp *core.Goploy) core.Response {
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

func (Server) GetTerminalRecord(*core.Goploy) core.Response {
	return response.File{Filename: path.Join(core.GetLogPath(), "demo.cast")}
}

func (Server) GetPublicKey(gp *core.Goploy) core.Response {
	publicKeyPath := gp.URLQuery.Get("path")

	contentByte, err := ioutil.ReadFile(publicKeyPath + ".pub")
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Key string `json:"key"`
		}{Key: string(contentByte)},
	}
}

func (Server) Check(gp *core.Goploy) core.Response {
	type ReqData struct {
		IP           string `json:"ip" validate:"required,ip|hostname"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Owner        string `json:"owner" validate:"required,max=255"`
		Path         string `json:"path" validate:"required,max=255"`
		Password     string `json:"password" validate:"max=255"`
		JumpIP       string `json:"jumpIP" validate:"omitempty,ip|hostname"`
		JumpPort     int    `json:"jumpPort" validate:"min=0,max=65535"`
		JumpOwner    string `json:"jumpOwner" validate:"max=255"`
		JumpPath     string `json:"jumpPath" validate:"max=255"`
		JumpPassword string `json:"jumpPassword"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	sshConfig := utils.SSHConfig{
		User:         reqData.Owner,
		Password:     reqData.Password,
		Path:         reqData.Path,
		Host:         reqData.IP,
		Port:         reqData.Port,
		JumpUser:     reqData.JumpOwner,
		JumpPassword: reqData.JumpPassword,
		JumpPath:     reqData.JumpPath,
		JumpHost:     reqData.JumpIP,
		JumpPort:     reqData.JumpPort,
	}

	if Conn, err := sshConfig.Dial(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	} else {
		_ = Conn.Close()
	}
	return response.JSON{Message: "Connected"}
}

func (s Server) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name         string `json:"name" validate:"required"`
		NamespaceID  int64  `json:"namespaceId" validate:"gte=0"`
		IP           string `json:"ip" validate:"ip|hostname"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Owner        string `json:"owner" validate:"required,max=255"`
		Path         string `json:"path" validate:"required,max=255"`
		Password     string `json:"password"`
		Description  string `json:"description" validate:"max=255"`
		JumpIP       string `json:"jumpIP" validate:"omitempty,ip|hostname"`
		JumpPort     int    `json:"jumpPort" validate:"min=0,max=65535"`
		JumpOwner    string `json:"jumpOwner" validate:"max=255"`
		JumpPath     string `json:"jumpPath" validate:"max=255"`
		JumpPassword string `json:"jumpPassword"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	server := model.Server{
		NamespaceID:  reqData.NamespaceID,
		Name:         reqData.Name,
		IP:           reqData.IP,
		Port:         reqData.Port,
		Owner:        reqData.Owner,
		Path:         reqData.Path,
		Password:     reqData.Password,
		JumpIP:       reqData.JumpIP,
		JumpPort:     reqData.JumpPort,
		JumpOwner:    reqData.JumpOwner,
		JumpPath:     reqData.JumpPath,
		JumpPassword: reqData.JumpPassword,
		Description:  reqData.Description,
	}
	server.OSInfo = server.Convert2SSHConfig().GetOSInfo()

	id, err := server.AddRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (s Server) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"gt=0"`
		NamespaceID  int64  `json:"namespaceId" validate:"gte=0"`
		Name         string `json:"name" validate:"required"`
		IP           string `json:"ip" validate:"required,ip|hostname"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Owner        string `json:"owner" validate:"required,max=255"`
		Path         string `json:"path" validate:"required,max=255"`
		Password     string `json:"password" validate:"max=255"`
		Description  string `json:"description" validate:"max=255"`
		JumpIP       string `json:"jumpIP" validate:"omitempty,ip|hostname"`
		JumpPort     int    `json:"jumpPort" validate:"min=0,max=65535"`
		JumpOwner    string `json:"jumpOwner" validate:"max=255"`
		JumpPath     string `json:"jumpPath" validate:"max=255"`
		JumpPassword string `json:"jumpPassword"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	server := model.Server{
		ID:           reqData.ID,
		NamespaceID:  reqData.NamespaceID,
		Name:         reqData.Name,
		IP:           reqData.IP,
		Port:         reqData.Port,
		Owner:        reqData.Owner,
		Path:         reqData.Path,
		Password:     reqData.Password,
		JumpIP:       reqData.JumpIP,
		JumpPort:     reqData.JumpPort,
		JumpOwner:    reqData.JumpOwner,
		JumpPath:     reqData.JumpPath,
		JumpPassword: reqData.JumpPassword,
		Description:  reqData.Description,
	}
	server.OSInfo = server.Convert2SSHConfig().GetOSInfo()

	if err := server.EditRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) Toggle(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID    int64 `json:"id" validate:"gt=0"`
		State int8  `json:"state" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Server{ID: reqData.ID, State: reqData.State}).ToggleRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Server) PreviewFile(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: "invalid server id"}
	}
	server, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := server.Convert2SSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.SftpFile{Filename: gp.URLQuery.Get("file"), Client: client, Disposition: "inline"}
}

func (Server) DownloadFile(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: "invalid server id"}
	}
	server, err := (model.Server{ID: id}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	client, err := server.Convert2SSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.SftpFile{Filename: gp.URLQuery.Get("file"), Client: client, Disposition: "attachment"}
}

func (Server) UploadFile(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID       int64  `schema:"id" validate:"gt=0"`
		FilePath string `schema:"filePath"  validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	server, err := (model.Server{ID: reqData.ID}).GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	file, fileHandler, err := gp.Request.FormFile("file")
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer file.Close()

	client, err := server.Convert2SSHConfig().Dial()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Create(reqData.FilePath + "/" + fileHandler.Filename)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, file)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Server) Report(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID      int64  `schema:"serverId" validate:"gt=0"`
		Type          int    `schema:"type" validate:"gt=0"`
		DatetimeRange string `schema:"datetimeRange"  validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	datetimeRange := strings.Split(reqData.DatetimeRange, ",")
	if len(datetimeRange) != 2 {
		return response.JSON{Code: response.Error, Message: "invalid datetime range"}
	}
	serverAgentLogs, err := (model.ServerAgentLog{ServerID: reqData.ServerID, Type: reqData.Type}).GetListBetweenTime(datetimeRange[0], datetimeRange[1])
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
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

	return response.JSON{
		Data: struct {
			ServerAgentMap map[string]model.ServerAgentLogs `json:"map"`
		}{ServerAgentMap: serverAgentMap},
	}
}

func (Server) GetAllMonitor(gp *core.Goploy) core.Response {
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	serverMonitorList, err := model.ServerMonitor{ServerID: serverID}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.ServerMonitors `json:"list"`
		}{List: serverMonitorList},
	}
}

func (s Server) AddMonitor(gp *core.Goploy) core.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
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
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (s Server) EditMonitor(gp *core.Goploy) core.Response {
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
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
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
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{}
}

func (s Server) DeleteMonitor(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.ServerMonitor{
		ID: reqData.ID,
	}.DeleteRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{}
}
