package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/media/feishu/api"
	"github.com/zhenorzz/goploy/internal/media/feishu/api/user_access_token"
	"github.com/zhenorzz/goploy/internal/media/feishu/api/user_info"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Feishu struct {
	Key    string
	Secret string
	Client *http.Client
	Method string
	Api    string
	Query  url.Values
	Body   interface{}
	Resp   interface{}
	Token  string
}

func (f *Feishu) Login(authCode string, redirectUri string) (string, error) {
	f.Key = config.Toml.Feishu.AppKey
	f.Secret = config.Toml.Feishu.AppSecret
	f.Client = &http.Client{}

	userAccessTokenInfo, err := f.GetUserAccessToken(authCode, redirectUri)
	if err != nil {
		return "", err
	}

	userInfo, err := f.GetUserInfo(userAccessTokenInfo.AccessToken)
	if err != nil {
		return "", err
	}

	if userInfo.UserId == "" {
		return "", errors.New("please scan the code again after joining the feishu company")
	}

	return strings.Trim(userInfo.Mobile, "+86"), nil
}

func (f *Feishu) Request() (err error) {
	var (
		req            *http.Request
		resp           *http.Response
		responseData   []byte
		commonResponse api.CommonResponse
	)

	uri, _ := url.Parse(f.Api)

	uri.RawQuery = f.Query.Encode()

	if f.Body != nil {
		b, _ := json.Marshal(f.Body)
		req, _ = http.NewRequest(f.Method, uri.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(f.Method, uri.String(), nil)
	}

	if f.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer ", f.Token))
	}

	if resp, err = f.Client.Do(req); err != nil {
		return err
	}

	defer resp.Body.Close()

	if responseData, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(responseData, &commonResponse); err != nil {
		return err
	}

	if commonResponse.Error != "" || commonResponse.Message != "" {
		return errors.New(fmt.Sprintf("api return error, code: %v, message: %s, error: %s", commonResponse.Code, commonResponse.Message, commonResponse.Error))
	}

	if err = json.Unmarshal(responseData, f.Resp); err != nil {
		return err
	}

	return nil
}

func (f *Feishu) GetUserAccessToken(authCode string, redirectUri string) (resp user_access_token.Response, err error) {
	f.Method = http.MethodPost
	f.Api = user_access_token.Url
	f.Query = nil
	f.Token = ""
	f.Body = user_access_token.Request{
		ClientId:     f.Key,
		ClientSecret: f.Secret,
		Code:         authCode,
		GrantType:    "authorization_code",
		RedirectUri:  redirectUri,
	}
	f.Resp = &resp

	return resp, f.Request()
}

func (f *Feishu) GetUserInfo(userAccessToken string) (resp user_info.Response, err error) {
	f.Method = http.MethodGet
	f.Api = user_info.Url
	f.Query = nil
	f.Token = userAccessToken
	f.Body = nil
	f.Resp = &resp

	return resp, f.Request()
}
