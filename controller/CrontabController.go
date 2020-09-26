package controller

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"strconv"
	"strings"
)

// Crontab struct
type Crontab Controller

// GetList crontab list
func (Crontab) GetList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	crontabs, err := model.Crontab{NamespaceID: gp.Namespace.ID, Command: gp.URLQuery.Get("command")}.GetList(pagination)

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Crontabs model.Crontabs `json:"list"`
		}{Crontabs: crontabs},
	}
}

// GetTotal crontab total
func (Crontab) GetTotal(gp *core.Goploy) *core.Response {
	total, err := model.Crontab{NamespaceID: gp.Namespace.ID, Command: gp.URLQuery.Get("command")}.GetTotal()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

// GetList crontab list
func (Crontab) GetRemoteServerList(gp *core.Goploy) *core.Response {
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	server, err := model.Server{
		ID: serverID,
	}.GetData()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	session, err := utils.ConnectSSH(server.Owner, "", server.IP, server.Port)

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	sshOutbuf.Reset()
	if err = session.Run("crontab -l"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	crontabList := strings.Split(sshOutbuf.String(), "\n")

	var crontabs []string
	for _, crontab := range crontabList {
		// windows \r\n
		crontab = strings.TrimRight(crontab, "\r")
		if len(crontab) == 0 {
			continue
		}
		// skip error format
		if len(strings.Split(crontab, " ")) < 5 {
			continue
		}
		// skip comment
		if strings.Index("#", crontab) == 0 {
			continue
		}
		crontabs = append(crontabs, crontab)
	}

	return &core.Response{
		Data: struct {
			Crontabs []string `json:"list"`
		}{Crontabs: crontabs},
	}
}

// GetBindServerList project detail
func (Crontab) GetBindServerList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	crontabServers, err := model.CrontabServer{CrontabID: id}.GetBindServerListByProjectID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			CrontabServers model.CrontabServers `json:"list"`
		}{CrontabServers: crontabServers},
	}
}

// Add one crontab
func (Crontab) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Command   string  `json:"command" validate:"required"`
		ServerIDs []int64 `json:"serverIds"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	crontabID, err := model.Crontab{
		NamespaceID: gp.Namespace.ID,
		Command:     reqData.Command,
		Creator:     gp.UserInfo.Name,
		CreatorID:   gp.UserInfo.ID,
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if len(reqData.ServerIDs) != 0 {
		crontabServersModel := model.CrontabServers{}
		for _, serverID := range reqData.ServerIDs {
			crontabServerModel := model.CrontabServer{
				CrontabID: crontabID,
				ServerID:  serverID,
			}

			go addCrontab(serverID, reqData.Command)

			crontabServersModel = append(crontabServersModel, crontabServerModel)
		}

		if err := crontabServersModel.AddMany(); err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}
	}

	return &core.Response{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: crontabID},
	}
}

// Edit one crontab
func (Crontab) Edit(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID      int64  `json:"id" validate:"gt=0"`
		Command string `json:"command" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	crontabInfo, err := model.Crontab{ID: reqData.ID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	err = model.Crontab{
		ID:       reqData.ID,
		Command:  reqData.Command,
		Editor:   gp.UserInfo.Name,
		EditorID: gp.UserInfo.ID,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	// 命令没修改过 不需要修改服务器的定时任务
	if crontabInfo.Command == reqData.Command {
		return &core.Response{}
	}

	crontabServers, err := model.CrontabServer{CrontabID: reqData.ID}.GetAllByCrontabID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	for _, crontabServer := range crontabServers {
		go deleteCrontab(crontabServer.ServerID, crontabInfo.Command)
		go addCrontab(crontabServer.ServerID, reqData.Command)
	}

	return &core.Response{}
}

// import many crontab
func (Crontab) Import(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Commands []string `json:"commands" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	commands := make(map[string]string)
	var commandMD5s []string
	for _, command := range reqData.Commands {
		if len(command) == 0 {
			continue
		}
		h := md5.New()
		h.Write([]byte(command))
		commandMD5 := hex.EncodeToString(h.Sum(nil))
		commands[commandMD5] = command
		commandMD5s = append(commandMD5s, commandMD5)
	}

	crontabList, err := model.Crontab{}.GetAllInCommandMD5(commandMD5s)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	for _, crontab := range crontabList {
		if _, ok := commands[crontab.CommandMD5]; ok {
			delete(commands, crontab.CommandMD5)
		}
	}

	var addCommands []string
	for _, command := range commands {
		addCommands = append(addCommands, command)
	}
	if len(addCommands) != 0 {
		err := model.Crontab{Creator: gp.UserInfo.Name, CreatorID: gp.UserInfo.ID}.AddRowsInCommand(addCommands)
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}
	}

	return &core.Response{}
}

// DeleteRow one Crontab
func (Crontab) Remove(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID    int64 `json:"id" validate:"gt=0"`
		Radio int8  `json:"radio" validate:"min=0,max=1"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if reqData.Radio == 1 {
		crontabInfo, err := model.Crontab{ID: reqData.ID}.GetData()
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}

		crontabServers, err := model.CrontabServer{CrontabID: reqData.ID}.GetAllByCrontabID()
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}

		for _, crontabServer := range crontabServers {
			go deleteCrontab(crontabServer.ServerID, crontabInfo.Command)
		}
	}

	if err := (model.Crontab{ID: reqData.ID}).DeleteRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.CrontabServer{CrontabID: reqData.ID}).DeleteByCrontabID(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{}
}

// AddServer one crontab
func (Crontab) AddServer(gp *core.Goploy) *core.Response {
	type ReqData struct {
		CrontabID int64   `json:"crontabId" validate:"gt=0"`
		ServerIDs []int64 `json:"serverIds" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	crontabID := reqData.CrontabID

	crontabInfo, err := model.Crontab{ID: crontabID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	crontabServersModel := model.CrontabServers{}
	for _, serverID := range reqData.ServerIDs {
		crontabServerModel := model.CrontabServer{
			CrontabID: crontabID,
			ServerID:  serverID,
		}
		go addCrontab(serverID, crontabInfo.Command)
		crontabServersModel = append(crontabServersModel, crontabServerModel)
	}

	if err := crontabServersModel.AddMany(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{}
}

// RemoveCrontabServer one crontab
func (Crontab) RemoveCrontabServer(gp *core.Goploy) *core.Response {
	type ReqData struct {
		CrontabServerID int64 `json:"crontabServerId" validate:"gt=0"`
		CrontabID       int64 `json:"crontabId" validate:"gt=0"`
		ServerID        int64 `json:"serverId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	crontabInfo, err := model.Crontab{ID: reqData.CrontabID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.CrontabServer{ID: reqData.CrontabServerID}).DeleteRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	go deleteCrontab(reqData.ServerID, crontabInfo.Command)

	return &core.Response{}
}

func addCrontab(serverID int64, command string) {
	server, err := model.Server{
		ID: serverID,
	}.GetData()

	if err != nil {
		core.Log(core.TRACE, "serverID:"+strconv.FormatUint(uint64(serverID), 10)+" get server fail, detail:"+err.Error())
		return
	}

	session, err := utils.ConnectSSH(server.Owner, "", server.IP, server.Port)

	if err != nil {
		core.Log(core.TRACE, "serverID:"+strconv.FormatUint(uint64(serverID), 10)+" connect server fail, detail:"+err.Error())
		return
	}

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	sshOutbuf.Reset()
	if err = session.Run(`(crontab -l ; echo "` + command + `") 2>&1 | grep -v "no crontab" | sort | uniq | crontab -`); err != nil {
		core.Log(core.TRACE, "serverID:"+strconv.FormatUint(uint64(serverID), 10)+" add "+command+" fail, detail:"+sshErrbuf.String())
		return
	}
}

func deleteCrontab(serverID int64, command string) {
	server, err := model.Server{ID: serverID}.GetData()

	if err != nil {
		core.Log(core.TRACE, "serverID:"+strconv.FormatUint(uint64(serverID), 10)+" get server fail, detail:"+err.Error())
		return
	}

	session, err := utils.ConnectSSH(server.Owner, "", server.IP, server.Port)

	if err != nil {
		core.Log(core.TRACE, "serverID:"+strconv.FormatUint(uint64(serverID), 10)+" connect server fail, detail:"+err.Error())
		return
	}

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	sshOutbuf.Reset()
	if err = session.Run(`(crontab -l ; echo "` + command + `") 2>&1 | grep -v "no crontab" | grep -v -F "` + command + `" | sort | uniq | crontab -`); err != nil {
		core.Log(core.TRACE, "serverID:"+strconv.FormatUint(uint64(serverID), 10)+" delete "+command+" fail, detail:"+sshErrbuf.String())
		return
	}
}
