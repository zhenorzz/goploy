package request

type UserAccessTokenReq struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Code         string `json:"code"`
	GrandType    string `json:"grantType"`
}

type GetUserIdByMobileReq struct {
	Mobile string `json:"mobile"`
}

type AccessTokenReq struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}
