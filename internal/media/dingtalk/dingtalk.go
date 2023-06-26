package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/api"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/cache"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Api               = "https://api.dingtalk.com/"
	Oapi              = "https://oapi.dingtalk.com/"
	UserAccessToken   = Api + "v1.0/oauth2/userAccessToken"
	ContactUserMe     = Api + "v1.0/contact/users/me"
	GetUserIdByMobile = Oapi + "topapi/v2/user/getbymobile"
	GetAccessToken    = Api + "v1.0/oauth2/accessToken"
)

type Dingtalk struct {
}

type Client struct {
	Key    string
	Secret string
	Client *http.Client
	Cache  *cache.AccessTokenCache
	Method string
	Api    string
	Query  url.Values
	Body   interface{}
	Resp   Response
}

type Response interface {
	CheckError() error
}

func (d Dingtalk) Login(authCode string, redirectUri string) (string, error) {
	dingtalkClient := Client{
		Key:    config.Toml.Dingtalk.AppKey,
		Secret: config.Toml.Dingtalk.AppSecret,
		Client: &http.Client{},
		Cache:  cache.GetCache(),
	}

	userAccessTokenInfo, err := dingtalkClient.GetUserAccessToken(authCode)
	if err != nil {
		return "", err
	}

	contactUserInfo, err := dingtalkClient.GetContactUser(userAccessTokenInfo.AccessToken)
	if err != nil {
		return "", err
	}

	mobileUserId, err := dingtalkClient.GetUserIdByMobile(contactUserInfo.Mobile)
	if err != nil {
		return "", err
	}

	if mobileUserId.Result.Userid == "" {
		return "", errors.New("please scan the code again after joining the dingtalk company")
	}

	return contactUserInfo.Mobile, nil
}

func (c *Client) Request() (err error) {
	var (
		req          *http.Request
		resp         *http.Response
		token        string
		responseData []byte
	)

	uri, _ := url.Parse(c.Api)
	if strings.Contains(c.Api, Api) {
		token = c.Query.Get("access_token")
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

	if strings.Contains(c.Api, Api) {
		req.Header.Set("x-acs-dingtalk-access-token", token)
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

func (c *Client) GetUserAccessToken(authCode string) (resp api.UserAccessTokenResp, err error) {
	param := api.UserAccessToken{
		Request: api.UserAccessTokenReq{
			ClientId:     c.Key,
			ClientSecret: c.Secret,
			Code:         authCode,
			GrandType:    "authorization_code",
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

func (c *Client) GetContactUser(userAccessToken string) (resp api.ContactUserResp, err error) {
	param := api.ContactUser{
		Request:  nil,
		Response: resp,
	}

	c.Method = http.MethodGet
	c.Api = ContactUserMe
	c.Query = url.Values{}
	c.Query.Set("access_token", userAccessToken)
	c.Body = param.Request
	c.Resp = &param.Response

	return param.Response, c.Request()
}

func (c *Client) GetUserIdByMobile(mobile string) (resp api.MobileUserIdResp, err error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return resp, err
	}

	param := api.MobileUserId{
		Request: api.GetUserIdByMobileReq{
			Mobile: mobile,
		},
		Response: resp,
	}

	c.Method = http.MethodPost
	c.Api = GetUserIdByMobile
	c.Query = url.Values{}
	c.Query.Set("access_token", accessToken)
	c.Body = param.Request
	c.Resp = &param.Response

	return param.Response, c.Request()
}

func (c *Client) GetAccessToken() (accessToken string, err error) {
	accessToken, ok := c.Cache.Get(c.Key)
	if !ok {
		param := api.AccessToken{
			Request: api.AccessTokenReq{
				AppKey:    c.Key,
				AppSecret: c.Secret,
			},
			Response: api.AccessTokenResp{},
		}

		c.Method = http.MethodPost
		c.Api = GetAccessToken
		c.Query = nil
		c.Body = param.Request
		c.Resp = &param.Response

		if err = c.Request(); err != nil {
			return "", err
		}

		c.Cache.Set(c.Key, param.Response.AccessToken, time.Duration(param.Response.ExpireIn))

		accessToken = param.Response.AccessToken
	}

	return accessToken, nil
}
