package task

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/service"
	"github.com/zhenorzz/goploy/ws"
	"net/http"
	"time"
)

type MonitorCache struct {
	errorTimes  int
	notifyTimes int
	time        int64
}

var monitorCaches = map[int64]MonitorCache{}

func monitorTask() {
	monitors, err := model.Monitor{State: model.Enable}.GetAllByState()
	if err != nil && err != sql.ErrNoRows {
		core.Log(core.ERROR, "get monitor list error, detail:"+err.Error())
	}
	monitorIDs := map[int64]struct{}{}
	for _, monitor := range monitors {
		monitorIDs[monitor.ID] = struct{}{}
		monitorCache, ok := monitorCaches[monitor.ID]
		if !ok {
			monitorCaches[monitor.ID] = MonitorCache{}
		}

		now := time.Now().Unix()

		if int(now-monitorCache.time) > monitor.Second {
			monitorCache.time = now
			if err := (service.Gnet{URL: monitor.URL}.Ping()); err != nil {
				monitorCache.errorTimes++
				core.Log(core.ERROR, "monitor "+monitor.Name+" encounter error, "+err.Error())
				if monitor.Times == uint16(monitorCache.errorTimes) {
					monitorCache.errorTimes = 0
					monitorCache.notifyTimes++
					notice(monitor, err)
					if monitor.NotifyTimes == uint16(monitorCache.notifyTimes) {
						monitorCache.notifyTimes = 0
						_ = monitor.TurnOff(err.Error())
						ws.GetHub().Data <- &ws.Data{
							Type:    ws.TypeMonitor,
							Message: ws.MonitorMessage{MonitorID: monitor.ID, State: ws.MonitorTurnOff, ErrorContent: err.Error()},
						}
					}
				}
			} else {
				monitorCache.errorTimes = 0
			}
			monitorCaches[monitor.ID] = monitorCache
		}
	}

	for cacheID := range monitorCaches {
		if _, ok := monitorIDs[cacheID]; !ok {
			delete(monitorCaches, cacheID)
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
		_, _ = http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if monitor.NotifyType == model.NotifyDingTalk {
		type markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		text := "#### Monitor: " + monitor.Name + " can not access \n >" + err.Error()

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Title: monitor.Name,
				Text:  text,
			},
		}
		b, _ := json.Marshal(msg)
		_, _ = http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
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
		_, _ = http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if monitor.NotifyType == model.NotifyCustom {
		type message struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				MonitorName string `json:"monitorName"`
				URL         string `json:"url"`
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
		msg.Data.URL = monitor.URL
		msg.Data.Second = monitor.Second
		msg.Data.Times = monitor.Times
		b, _ := json.Marshal(msg)
		_, _ = http.Post(monitor.NotifyTarget, "application/json", bytes.NewBuffer(b))
	}
}
