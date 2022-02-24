package middleware

import (
	"encoding/json"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"time"
)

func AddLoginLog(gp *core.Goploy, resp core.Response) {
	respJson := resp.(response.JSON)
	account := ""
	if respJson.Code != response.IllegalParam {
		type ReqData struct {
			Account string `json:"account"`
		}
		var reqData ReqData
		_ = json.Unmarshal(gp.Body, &reqData)
		account = reqData.Account
	}

	err := model.LoginLog{
		Account:    account,
		RemoteAddr: gp.Request.RemoteAddr,
		UserAgent:  gp.Request.UserAgent(),
		Referer:    gp.Request.Referer(),
		Reason:     respJson.Message,
		LoginTime:  time.Now().Format("20060102150405"),
	}.AddRow()
	if err != nil {
		core.Log(core.ERROR, err.Error())
	}
}
