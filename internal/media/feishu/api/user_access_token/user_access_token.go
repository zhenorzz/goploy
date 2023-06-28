package user_access_token

import "github.com/zhenorzz/goploy/internal/media/feishu/api"

const Url = api.OauthApi + "token"

type Request struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectUri  string `json:"redirect_uri"`
}

type Response struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
}
