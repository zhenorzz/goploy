package request

type UserAccessTokenReq struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectUri  string `json:"redirect_uri"`
}
