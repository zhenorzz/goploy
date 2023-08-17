package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/cache"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/api"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/api/access_token"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/api/contact"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/api/get_user_by_mobile"
	"github.com/zhenorzz/goploy/internal/media/dingtalk/api/user_access_token"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Dingtalk struct {
	Key    string
	Secret string
	Client *http.Client
	Cache  cache.DingtalkAccessToken
	Method string
	Api    string
	Query  url.Values
	Body   interface{}
	Token  string
}

func (d *Dingtalk) Login(authCode string, _ string) (string, error) {
	d.Key = config.Toml.Dingtalk.AppKey
	d.Secret = config.Toml.Dingtalk.AppSecret
	d.Client = &http.Client{}
	d.Cache = cache.GetDingTalkAccessTokenCache()

	userAccessTokenInfo, err := d.GetUserAccessToken(authCode)
	if err != nil {
		return "", err
	}

	contactUserInfo, err := d.GetContactUser(userAccessTokenInfo.AccessToken)
	if err != nil {
		return "", err
	}

	mobileUserId, err := d.GetUserIdByMobile(contactUserInfo.Mobile)
	if err != nil {
		return "", err
	}

	if mobileUserId.Result.Userid == "" {
		return "", errors.New("please scan the code again after joining the dingtalk company")
	}

	return contactUserInfo.Mobile, nil
}

func request[T any](d *Dingtalk) (response T, err error) {
	var (
		req            *http.Request
		resp           *http.Response
		responseData   []byte
		commonResponse api.CommonResponse
	)

	uri, _ := url.Parse(d.Api)
	uri.RawQuery = d.Query.Encode()

	if d.Body != nil {
		b, _ := json.Marshal(d.Body)
		req, _ = http.NewRequest(d.Method, uri.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(d.Method, uri.String(), nil)
	}

	if d.Token != "" {
		req.Header.Set("x-acs-dingtalk-access-token", d.Token)
	}

	if resp, err = d.Client.Do(req); err != nil {
		return
	}

	defer resp.Body.Close()

	if responseData, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(responseData, &commonResponse); err != nil {
		return
	}

	if commonResponse.Code != "" {
		return response, errors.New(fmt.Sprintf("api return error, code: %s, message: %s, request_id: %s", commonResponse.Code, commonResponse.Message, commonResponse.RequestId))
	} else if commonResponse.ErrCode != 0 {
		return response, errors.New(fmt.Sprintf("api return error, code: %v, message: %s, request_id: %s", commonResponse.ErrCode, commonResponse.ErrMsg, commonResponse.OldRequestId))
	}

	if err = json.Unmarshal(responseData, &response); err != nil {
		return
	}

	return
}

func (d *Dingtalk) GetUserAccessToken(authCode string) (resp user_access_token.Response, err error) {
	d.Method = http.MethodPost
	d.Api = user_access_token.Url
	d.Query = nil
	d.Token = ""
	d.Body = user_access_token.Request{
		ClientId:     d.Key,
		ClientSecret: d.Secret,
		Code:         authCode,
		GrandType:    "authorization_code",
	}

	return request[user_access_token.Response](d)
}

func (d *Dingtalk) GetContactUser(userAccessToken string) (contact.Response, error) {
	d.Method = http.MethodGet
	d.Api = contact.Url
	d.Query = nil
	d.Token = userAccessToken
	d.Body = nil

	return request[contact.Response](d)
}

func (d *Dingtalk) GetUserIdByMobile(mobile string) (get_user_by_mobile.Response, error) {
	t, err := d.GetAccessToken()
	if err != nil {
		return get_user_by_mobile.Response{}, err
	}

	d.Method = http.MethodPost
	d.Api = get_user_by_mobile.Url
	d.Query = url.Values{}
	d.Query.Set("access_token", t)
	d.Token = ""
	d.Body = get_user_by_mobile.Request{
		Mobile: mobile,
	}

	return request[get_user_by_mobile.Response](d)
}

func (d *Dingtalk) GetAccessToken() (string, error) {
	t, ok := d.Cache.Get(d.Key)
	if !ok {
		var resp access_token.Response

		d.Method = http.MethodPost
		d.Api = access_token.Url
		d.Query = nil
		d.Body = access_token.Request{
			AppKey:    d.Key,
			AppSecret: d.Secret,
		}
		resp, err := request[access_token.Response](d)
		if err != nil {
			return "", err
		}

		d.Cache.Set(d.Key, resp.AccessToken, time.Duration(resp.ExpireIn)*time.Second)
		t = resp.AccessToken
	}

	return t, nil
}
