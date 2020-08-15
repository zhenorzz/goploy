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
		content := "Monitor: <font color=\"warning\">" + monitor.Name + "</font>\n "
		content += "> <font color=\"warning\">can not access</font> \n "
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
		text := "> <font color=\"red\">can not access</font> \n "
		text += "> <font color=\"comment\">" + err.Error() + "</font> \n "

		msg := message{
			Msgtype: "markdown",
			Title:   "Monitor: " + monitor.Name,
			Text:    text,
		}
		b, _ := json.Marshal(msg)
		http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if monitor.NotifyType == model.NotifyFeiShu {
		type message struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}

		text := "can not access\n "
		text += "detail:  " + err.Error()

		msg := message{
			Title: "Monitor:" + monitor.Name,
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
			Message: "Monitor:" + monitor.Name + "can not access",
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
