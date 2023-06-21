package response

import (
	"errors"
	"fmt"
)

type CommonResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Code             int    `json:"code"`
	Message          string `json:"message"`
}

type UserAccessTokenResp struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
	CommonResponse
}

type UserInfoResp struct {
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
	CommonResponse
}

func (r *CommonResponse) CheckError() (err error) {
	if r.Error != "" || r.Message != "" {
		err = errors.New(fmt.Sprintf("api return error, code: %v, message: %s, error: %s", r.Code, r.Message, r.Error))
	}
	return err
}
