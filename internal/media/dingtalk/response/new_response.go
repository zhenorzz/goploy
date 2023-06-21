package response

import (
	"errors"
	"fmt"
)

type NewResponse struct {
	Code      string `json:"code"`
	RequestId string `json:"requestid"`
	Message   string `json:"message"`
	Result    bool   `json:"result"`
}

type UserAccessTokenResp struct {
	ExpireIn     int    `json:"expireIn"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	NewResponse
}

type ContactUserResp struct {
	Nick      string `json:"nick"`
	UnionId   string `json:"unionId"`
	OpenId    string `json:"openId"`
	Mobile    string `json:"mobile"`
	StateCode string `json:"stateCode"`
	NewResponse
}

type AccessTokenResp struct {
	ExpireIn    int    `json:"expireIn"`
	AccessToken string `json:"accessToken"`
	NewResponse
}

func (r *NewResponse) CheckError() (err error) {
	if r.Code != "" {
		err = errors.New(fmt.Sprintf("api return error, code: %s, message: %s, request_id: %s", r.Code, r.Message, r.RequestId))
	}
	return err
}
