package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/media/feishu/api"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	OauthApi        = "https://passport.feishu.cn/suite/passport/oauth/"
	UserAccessToken = OauthApi + "token"
	UserInfo        = OauthApi + "userinfo"
)

type Feishu struct {
}

type Client struct {
	Key    string
	Secret string
	Client *http.Client
	Method string
	Api    string
	Query  url.Values
	Body   interface{}
	Resp   Response
}

type Response interface {
	CheckError() error
}

func (f Feishu) Login(authCode string, redirectUri string) (string, error) {
	feishuClient := Client{
		Key:    config.Toml.Feishu.AppKey,
		Secret: config.Toml.Feishu.AppSecret,
		Client: &http.Client{},
	}

	userAccessTokenInfo, err := feishuClient.GetUserAccessToken(authCode, redirectUri)
	if err != nil {
		return "", err
	}

	userInfo, err := feishuClient.GetUserInfo(userAccessTokenInfo.AccessToken)
	if err != nil {
		return "", err
	}

	if userInfo.UserId == "" {
		return "", errors.New("please scan the code again after joining the feishu company")
	}

	return strings.Trim(userInfo.Mobile, "+86"), nil
}

func (c *Client) Request() (err error) {
	var (
		req          *http.Request
		resp         *http.Response
		token        string
		responseData []byte
	)

	uri, _ := url.Parse(c.Api)

	token = c.Query.Get("access_token")
	if token != "" {
		c.Query.Del("access_token")
	}

	uri.RawQuery = c.Query.Encode()

	if c.Body != nil {
		b, _ := json.Marshal(c.Body)
		req, _ = http.NewRequest(c.Method, uri.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(c.Method, uri.String(), nil)
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer ", token))
	}

	if resp, err = c.Client.Do(req); err != nil {
		return err
	}

	defer resp.Body.Close()

	if responseData, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(responseData, c.Resp); err != nil {
		return err
	}

	if err = c.Resp.CheckError(); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetUserAccessToken(authCode string, redirectUri string) (resp api.UserAccessTokenResp, err error) {
	param := api.UserAccessToken{
		Request: api.UserAccessTokenReq{
			ClientId:     c.Key,
			ClientSecret: c.Secret,
			Code:         authCode,
			GrantType:    "authorization_code",
			RedirectUri:  redirectUri,
		},
		Response: resp,
	}

	c.Method = http.MethodPost
	c.Api = UserAccessToken
	c.Query = nil
	c.Body = param.Request
	c.Resp = &param.Response

	return param.Response, c.Request()
}

func (c *Client) GetUserInfo(userAccessToken string) (resp api.UserInfoResp, err error) {
	param := api.UserInfo{
		Request:  nil,
		Response: resp,
	}

	c.Method = http.MethodGet
	c.Api = UserInfo
	c.Query = url.Values{}
	c.Query.Set("access_token", userAccessToken)
	c.Body = param.Request
	c.Resp = &param.Response

	return param.Response, c.Request()
}
