package api

const OauthApi = "https://passport.feishu.cn/suite/passport/oauth/"

type CommonResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Code             int    `json:"code"`
	Message          string `json:"message"`
}
