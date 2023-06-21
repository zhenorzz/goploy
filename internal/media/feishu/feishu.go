package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/media/feishu/request"
	"github.com/zhenorzz/goploy/internal/media/feishu/response"
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

func (c *Client) Request(method, api string, query url.Values, body interface{}, data Response) (err error) {
	var (
		req          *http.Request
		resp         *http.Response
		token        string
		responseData []byte
	)

	uri, _ := url.Parse(api)

	token = query.Get("access_token")
	if token != "" {
		query.Del("access_token")
	}

	uri.RawQuery = query.Encode()

	if body != nil {
		b, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, uri.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, uri.String(), nil)
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

	if err = json.Unmarshal(responseData, &data); err != nil {
		return err
	}

	if err = data.CheckError(); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetUserAccessToken(authCode string, redirectUri string) (resp response.UserAccessTokenResp, err error) {
	req := request.UserAccessTokenReq{
		ClientId:     c.Key,
		ClientSecret: c.Secret,
		Code:         authCode,
		GrantType:    "authorization_code",
		RedirectUri:  redirectUri,
	}
	return resp, c.Request(http.MethodPost, UserAccessToken, nil, req, &resp)
}

func (c *Client) GetUserInfo(userAccessToken string) (resp response.UserInfoResp, err error) {
	query := url.Values{}
	query.Set("access_token", userAccessToken)

	return resp, c.Request(http.MethodGet, UserInfo, query, nil, &resp)
}
