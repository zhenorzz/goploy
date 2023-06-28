package contact

import "github.com/zhenorzz/goploy/internal/media/dingtalk/api"

const Url = api.Api + "v1.0/contact/users/me"

type Response struct {
	Nick      string `json:"nick"`
	UnionId   string `json:"unionId"`
	OpenId    string `json:"openId"`
	Mobile    string `json:"mobile"`
	StateCode string `json:"stateCode"`
}
