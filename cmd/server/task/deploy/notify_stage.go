package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/internal/model"
	"io"
	"net/http"
	"strings"
)

// commit id
// commit message
// server ip & name
// deploy user name
// deploy time
func (gsync *Gsync) notify(deployState int, detail string) {
	if gsync.Project.NotifyType == 0 {
		return
	}
	serverList := ""
	for _, projectServer := range gsync.ProjectServers {
		if projectServer.Server.Name != projectServer.Server.IP {
			serverList += projectServer.Server.Name + "(" + projectServer.Server.IP + ")"
		} else {
			serverList += projectServer.Server.IP
		}
		serverList += ", "
	}
	serverList = strings.TrimRight(serverList, ", ")
	project := gsync.Project
	commitInfo := gsync.CommitInfo
	var err error
	var resp *http.Response
	if project.NotifyType == model.NotifyWeiXin {
		type markdown struct {
			Content string `json:"content"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		content := "Deploy: <font color=\"warning\">" + project.Name + "</font>\n"
		content += "Publisher: <font color=\"comment\">" + project.PublisherName + "</font>\n"
		content += "Author: <font color=\"comment\">" + commitInfo.Author + "</font>\n"
		if commitInfo.Tag != "" {
			content += "Tag: <font color=\"comment\">" + commitInfo.Tag + "</font>\n"
		}
		content += "Branch: <font color=\"comment\">" + commitInfo.Branch + "</font>\n"
		content += "CommitSHA: <font color=\"comment\">" + commitInfo.Commit + "</font>\n"
		content += "CommitMessage: <font color=\"comment\">" + commitInfo.Message + "</font>\n"
		content += "ServerList: <font color=\"comment\">" + serverList + "</font>\n"
		if deployState == model.ProjectFail {
			content += "State: <font color=\"red\">fail</font> \n"
			content += "> Detail: <font color=\"comment\">" + detail + "</font>"
		} else {
			content += "State: <font color=\"green\">success</font>"
		}

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Content: content,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if project.NotifyType == model.NotifyDingTalk {
		type markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		text := "#### Deploy：" + project.Name + "  \n  "
		text += "#### Publisher：" + project.PublisherName + "  \n  "
		text += "#### Author：" + commitInfo.Author + "  \n  "
		if commitInfo.Tag != "" {
			text += "#### Tag：" + commitInfo.Tag + "  \n  "
		}
		text += "#### Branch：" + commitInfo.Branch + "  \n  "
		text += "#### CommitSHA：" + commitInfo.Commit + "  \n  "
		text += "#### CommitMessage：" + commitInfo.Message + "  \n  "
		text += "#### ServerList：" + serverList + "  \n  "
		if deployState == model.ProjectFail {
			text += "#### State： <font color=\"red\">fail</font>  \n  "
			text += "> Detail: <font color=\"comment\">" + detail + "</font>"
		} else {
			text += "#### State： <font color=\"green\">success</font>"
		}

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Title: project.Name,
				Text:  text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if project.NotifyType == model.NotifyFeiShu {
		type content struct {
			Text string `json:"text"`
		}
		type message struct {
			MsgType string  `json:"msg_type"`
			Content content `json:"content"`
		}
		text := ""
		text += "Deploy：" + project.Name + "\n"
		text += "Publisher: " + project.PublisherName + "\n"
		text += "Author: " + commitInfo.Author + "\n"
		if commitInfo.Tag != "" {
			text += "Tag: " + commitInfo.Tag + "\n"
		}
		text += "Branch: " + commitInfo.Branch + "\n"
		text += "CommitSHA: " + commitInfo.Commit + "\n"
		text += "CommitMessage: " + commitInfo.Message + "\n"
		text += "ServerList: " + serverList + "\n"
		if deployState == model.ProjectFail {
			text += "State: fail\n "
			text += "Detail: " + detail
		} else {
			text += "State: success"
		}

		msg := message{
			MsgType: "text",
			Content: content{
				Text: text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if project.NotifyType == model.NotifyCustom {
		type message struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				ProjectID     int64  `json:"projectId"`
				ProjectName   string `json:"projectName"`
				Publisher     string `json:"publisher"`
				Author        string `json:"author"`
				Branch        string `json:"branch"`
				Tag           string `json:"tag"`
				CommitSHA     string `json:"commitSHA"`
				CommitMessage string `json:"commitMessage"`
				ServerList    string `json:"serverList"`
			} `json:"data"`
		}
		code := 0
		if deployState == model.ProjectFail {
			code = 1
		}
		msg := message{
			Code:    code,
			Message: detail,
		}
		msg.Data.ProjectID = project.ID
		msg.Data.ProjectName = project.Name
		msg.Data.Publisher = project.PublisherName
		msg.Data.Author = commitInfo.Author
		msg.Data.Branch = commitInfo.Branch
		msg.Data.Tag = commitInfo.Tag
		msg.Data.CommitSHA = commitInfo.Commit
		msg.Data.CommitMessage = commitInfo.Message
		msg.Data.ServerList = serverList
		b, _ := json.Marshal(msg)
		resp, err = http.Post(project.NotifyTarget, "application/json", bytes.NewBuffer(b))
	}

	if err != nil {
		log.Error(fmt.Sprintf("projectID: %d notify exec err: %s", project.ID, err))
	} else {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(fmt.Sprintf("projectID: %d notify read body err: %s", project.ID, err))
		} else {
			log.Trace(fmt.Sprintf("projectID: %d notify success: %s", project.ID, string(responseData)))
		}
		_ = resp.Body.Close()
	}
}
