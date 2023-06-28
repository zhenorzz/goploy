package access_token

import "github.com/zhenorzz/goploy/internal/media/dingtalk/api"

const Url = api.Api + "v1.0/oauth2/accessToken"

type Request struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}

type Response struct {
	ExpireIn    int    `json:"expireIn"`
	AccessToken string `json:"accessToken"`
}
