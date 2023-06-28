package get_user_by_mobile

import (
	"github.com/zhenorzz/goploy/internal/media/dingtalk/api"
)

const Url = api.Oapi + "topapi/v2/user/getbymobile"

type Request struct {
	Mobile string `json:"mobile"`
}

type Response struct {
	Result struct {
		Userid string `json:"userid"`
	} `json:"result"`
}
