package user_info

import "github.com/zhenorzz/goploy/internal/media/feishu/api"

const Url = api.OauthApi + "userinfo"

type Response struct {
	Sub          string `json:"sub"`
	Picture      string `json:"picture"`
	Name         string `json:"name"`
	EnName       string `json:"en_name"`
	TenantKey    string `json:"tenant_key"`
	AvatarUrl    string `json:"avatar_url"`
	AvatarThumb  string `json:"avatar_thumb"`
	AvatarMiddle string `json:"avatar_middle"`
	AvatarBig    string `json:"avatar_big"`
	OpenId       string `json:"open_id"`
	UnionId      string `json:"union_id"`
	UserId       string `json:"user_id"`
	Mobile       string `json:"mobile"`
}
