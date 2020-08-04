package task

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"net"
	"net/http"
	"strconv"
	"time"
)

func monitorTask() {
	monitors, err := model.Monitor{State: model.Enable}.GetAllByState()
	if err != nil && err != sql.ErrNoRows {
		core.Log(core.ERROR, "get monitor list error, detail:"+err.Error())
	}
	for _, monitor := range monitors {
		monitorCache := map[string]int64{"errorTimes": 0, "time": 0}
		if x, found := core.Cache.Get("monitor:" + strconv.Itoa(int(monitor.ID))); found {
			monitorCache = x.(map[string]int64)
		}
		now := time.Now().Unix()

		if int(now-monitorCache["time"]) > monitor.Second {
			monitorCache["time"] = now
			conn, err := net.DialTimeout("tcp", monitor.Domain+":"+strconv.Itoa(monitor.Port), 5*time.Second)
			if err != nil {
				monitorCache["errorTimes"]++
				core.Log(core.ERROR, "monitor "+monitor.Name+" encounter error, "+err.Error())
				if monitor.Times == uint16(monitorCache["errorTimes"]) {
					monitorCache["errorTimes"] = 0
					notice(monitor, err)
				}
			} else {
				conn.Close()
			}
			core.Cache.Set("monitor:"+strconv.Itoa(int(monitor.ID)), monitorCache, cache.DefaultExpiration)
		}
	}
}

func notice(monitor model.Monitor, err error) {
	if monitor.NotifyType == model.NotifyWeiXin {
		type markdown struct {
			Content string `json:"content"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		content := "应用监控<font color=\"warning\">" + monitor.Name + "</font>，请相关同事注意。\n "
		content += "> <font color=\"warning\">无法访问</font> \n "
		content += "> <font color=\"comment\">" + err.Error() + "</font> \n "

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Content: content,
			},
		}
		b, _ := json.Marshal(msg)
		http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if monitor.NotifyType == model.NotifyDingTalk {
		type message struct {
			Msgtype string `json:"msgtype"`
			Title   string `json:"title"`
			Text    string `json:"text"`
		}
		text := "> <font color=\"red\">无法访问</font> \n "
		text += "> <font color=\"comment\">" + err.Error() + "</font> \n "

		msg := message{
			Msgtype: "markdown",
			Title:   "应用监控:" + monitor.Name,
			Text:    text,
		}
		b, _ := json.Marshal(msg)
		http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if monitor.NotifyType == model.NotifyFeiShu {
		type message struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}

		text := "无法访问\n "
		text += "详情: " + err.Error()

		msg := message{
			Title: "应用监控:" + monitor.Name,
			Text:  text,
		}
		b, _ := json.Marshal(msg)
		http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if monitor.NotifyType == model.NotifyCustom {
		type message struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				MonitorName string `json:"monitorName"`
				Domain      string `json:"domain"`
				Port        int    `json:"port"`
				Second      int    `json:"second"`
				Times       uint16 `json:"times"`
			} `json:"data"`
		}
		code := 0
		msg := message{
			Code:    code,
			Message: "应用监控:" + monitor.Name + "无法访问",
		}
		msg.Data.MonitorName = monitor.Name
		msg.Data.Domain = monitor.Domain
		msg.Data.Port = monitor.Port
		msg.Data.Second = monitor.Second
		msg.Data.Times = monitor.Times
		b, _ := json.Marshal(msg)
		http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	}
}
