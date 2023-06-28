package user_access_token

import "github.com/zhenorzz/goploy/internal/media/dingtalk/api"

const Url = api.Api + "v1.0/oauth2/userAccessToken"

type Request struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Code         string `json:"code"`
	GrandType    string `json:"grantType"`
}

type Response struct {
	ExpireIn     int    `json:"expireIn"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
