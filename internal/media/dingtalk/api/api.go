package api

import (
	"errors"
	"fmt"
)

type UserAccessToken struct {
	Request  UserAccessTokenReq
	Response UserAccessTokenResp
}

type ContactUser struct {
	Request  interface{}
	Response ContactUserResp
}

type MobileUserId struct {
	Request  GetUserIdByMobileReq
	Response MobileUserIdResp
}

type AccessToken struct {
	Request  AccessTokenReq
	Response AccessTokenResp
}

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

type CommonResponse struct {
	Code         string `json:"code"`
	RequestId    string `json:"requestid"`
	Message      string `json:"message"`
	Result       bool   `json:"result"`
	ErrCode      int    `json:"errcode"`
	OldRequestId string `json:"request_id"`
	ErrMsg       string `json:"errmsg"`
}

func (r *CommonResponse) CheckError() (err error) {
	if r.Code != "" {
		err = errors.New(fmt.Sprintf("api return error, code: %s, message: %s, request_id: %s", r.Code, r.Message, r.RequestId))
	} else if r.ErrCode != 0 {
		err = errors.New(fmt.Sprintf("api return error, code: %v, message: %s, request_id: %s", r.ErrCode, r.ErrMsg, r.OldRequestId))
	}
	return err
}

type UserAccessTokenResp struct {
	CommonResponse
	ExpireIn     int    `json:"expireIn"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ContactUserResp struct {
	CommonResponse
	Nick      string `json:"nick"`
	UnionId   string `json:"unionId"`
	OpenId    string `json:"openId"`
	Mobile    string `json:"mobile"`
	StateCode string `json:"stateCode"`
}

type MobileUserIdResp struct {
	CommonResponse
	Result struct {
		Userid string `json:"userid"`
	} `json:"result"`
}

type AccessTokenResp struct {
	CommonResponse
	ExpireIn    int    `json:"expireIn"`
	AccessToken string `json:"accessToken"`
}
