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

func request[T any](f *Feishu) (response T, err error) {
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
		return
	}

	defer resp.Body.Close()

	if responseData, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(responseData, &commonResponse); err != nil {
		return
	}

	if commonResponse.Error != "" || commonResponse.Message != "" {
		return response, errors.New(fmt.Sprintf("api return error, code: %v, message: %s, error: %s", commonResponse.Code, commonResponse.Message, commonResponse.Error))
	}

	if err = json.Unmarshal(responseData, &response); err != nil {
		return
	}

	return
}

func (f *Feishu) GetUserAccessToken(authCode string, redirectUri string) (user_access_token.Response, error) {
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

	return request[user_access_token.Response](f)
}

func (f *Feishu) GetUserInfo(userAccessToken string) (user_info.Response, error) {
	f.Method = http.MethodGet
	f.Api = user_info.Url
	f.Query = nil
	f.Token = userAccessToken
	f.Body = nil

	return request[user_info.Response](f)
}
