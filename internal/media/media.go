package media

import (
	"github.com/zhenorzz/goploy/internal/media/dingtalk"
	"github.com/zhenorzz/goploy/internal/media/feishu"
	"strings"
)

type Media interface {
	Login(authCode string, redirectUri string) (string, error)
}

func GetMedia(state string) Media {
	if strings.Contains(state, "dingtalk") {
		return &dingtalk.Dingtalk{}
	} else {
		return &feishu.Feishu{}
	}
}
