package response

import (
	"errors"
	"fmt"
)

type OldResponse struct {
	Code      int    `json:"errcode"`
	RequestId string `json:"request_id"`
	Message   string `json:"errmsg"`
	Result    bool   `json:"result"`
}

type MobileUserId struct {
	OldResponse
	Result struct {
		Userid string `json:"userid"`
	} `json:"result"`
}

func (r OldResponse) CheckError() (err error) {
	if r.Code != 0 {
		err = errors.New(fmt.Sprintf("api return error, code: %v, message: %s, request_id: %s", r.Code, r.Message, r.RequestId))
	}
	return err
}
