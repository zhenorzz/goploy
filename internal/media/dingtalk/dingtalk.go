package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/cache"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/request"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/response"
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

func (c *Client) Request(method, api string, query url.Values, body interface{}, data Response) (err error) {
	var (
		req          *http.Request
		resp         *http.Response
		token        string
		responseData []byte
	)

	uri, _ := url.Parse(api)
	if strings.Contains(api, Api) {
		token = query.Get("access_token")
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

	if strings.Contains(api, Api) {
		req.Header.Set("x-acs-dingtalk-access-token", token)
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

func (c *Client) GetUserAccessToken(authCode string) (resp response.UserAccessTokenResp, err error) {
	req := request.UserAccessTokenReq{
		ClientId:     c.Key,
		ClientSecret: c.Secret,
		Code:         authCode,
		GrandType:    "authorization_code",
	}
	return resp, c.Request(http.MethodPost, UserAccessToken, nil, req, &resp)
}

func (c *Client) GetContactUser(userAccessToken string) (resp response.ContactUserResp, err error) {
	query := url.Values{}
	query.Set("access_token", userAccessToken)

	return resp, c.Request(http.MethodGet, ContactUserMe, query, nil, &resp)
}

func (c *Client) GetUserIdByMobile(mobile string) (resp response.MobileUserId, err error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return resp, err
	}

	query := url.Values{}
	query.Set("access_token", accessToken)

	reqParam := request.GetUserIdByMobileReq{
		Mobile: mobile,
	}

	return resp, c.Request(http.MethodPost, GetUserIdByMobile, query, reqParam, &resp)
}

func (c *Client) GetAccessToken() (accessToken string, err error) {
	accessToken, ok := c.Cache.Get(c.Key)
	if !ok {
		reqParam := request.AccessTokenReq{
			AppKey:    c.Key,
			AppSecret: c.Secret,
		}

		var resp response.AccessTokenResp

		if err = c.Request(http.MethodPost, GetAccessToken, nil, reqParam, &resp); err != nil {
			return "", err
		}

		c.Cache.Set(c.Key, resp.AccessToken, time.Duration(resp.ExpireIn))

		accessToken = resp.AccessToken
	}

	return accessToken, nil
}
