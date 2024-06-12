package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/repo"
	"io"
	"net/http"
	"text/template"
)

type DeployData struct {
	DeployState    uint8
	Project        model.Project
	ProjectServers model.ProjectServers
	CommitInfo     repo.CommitInfo
	DeployDetail   string
}

type MonitorData struct {
	Monitor  model.Monitor
	ErrorMsg string
}

type ServerMonitorData struct {
	Server        model.Server
	ServerMonitor model.ServerMonitor
	CycleValue    string
}

const (
	UseByDeploy        = "deploy"
	UseByMonitor       = "monitor"
	UseByServerMonitor = "server_monitor"
)

func Send(name string, useBy string, data any, notifyType uint8, notifyTarget string) error {
	if notifyType == model.NotifyCustom {
		type message struct {
			Data any `json:"data"`
		}
		msg := message{
			Data: data,
		}
		b, _ := json.Marshal(msg)
		_, err := http.Post(notifyTarget, "application/json", bytes.NewBuffer(b))
		if err != nil {
			log.Error(fmt.Sprintf("%s notify exec err: %s", name, err))
			return err
		}
		return nil
	}

	notificationData, err := model.NotificationTemplate{UseBy: useBy, Type: notifyType}.GetTemplate()
	if err != nil {
		log.Error(fmt.Sprintf("%s could not find notification template: %s", name, err))
		return err
	}
	var buf bytes.Buffer

	tmpl, err := template.New(name + "title").Parse(notificationData.Title)
	if err != nil {
		log.Error(fmt.Sprintf("%s parse notification title error: %s", name, err))
		return err
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s execute notification title error: %s", name, err))
		return err
	}

	title := buf.String()

	tmpl, err = template.New(name + "template").Parse(notificationData.Template)
	if err != nil {
		log.Error(fmt.Sprintf("%s parse notification template error: %s", name, err))
		return err
	}

	buf.Reset()

	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s execute notification template error: %s", name, err))
		return err
	}

	text := buf.String()

	println(title)

	println(text)

	var resp *http.Response
	if notifyType == model.NotifyWeiXin {
		type markdown struct {
			Content string `json:"content"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}

		text = title + "\n" + text
		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Content: text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(notifyTarget, "application/json", bytes.NewBuffer(b))
	} else if notifyType == model.NotifyDingTalk {
		type markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Title: title,
				Text:  text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(notifyTarget, "application/json", bytes.NewBuffer(b))
	} else if notifyType == model.NotifyFeiShu {
		type content struct {
			Text string `json:"text"`
		}
		type message struct {
			MsgType string  `json:"msg_type"`
			Content content `json:"content"`
		}
		text = title + "\n" + text
		msg := message{
			MsgType: "text",
			Content: content{
				Text: text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(notifyTarget, "application/json", bytes.NewBuffer(b))
	}

	if err != nil {
		log.Error(fmt.Sprintf("%s notify exec err: %s", name, err))
		return err
	} else if resp != nil {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(fmt.Sprintf("%s notify read body err: %s", name, err))
		} else {
			log.Trace(fmt.Sprintf("%s notify success: %s", name, string(responseData)))
		}
		_ = resp.Body.Close()
	}
	return nil
}
