package api

const (
	Api  = "https://api.dingtalk.com/"
	Oapi = "https://oapi.dingtalk.com/"
)

type CommonResponse struct {
	Code         string `json:"code"`
	RequestId    string `json:"requestid"`
	Message      string `json:"message"`
	ErrCode      int    `json:"errcode"`
	OldRequestId string `json:"request_id"`
	ErrMsg       string `json:"errmsg"`
}
