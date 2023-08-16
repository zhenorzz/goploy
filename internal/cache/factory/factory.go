package factory

import (
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/cache"
	"github.com/zhenorzz/goploy/internal/cache/memory"
)

func GetUserCache() cache.User {
	cacheType := config.Toml.Cache.Type
	if cacheType == "memory" {
		return memory.GetUserCache()
	} else {
		return memory.GetUserCache()
	}
}

func GetCaptchaCache() cache.Captcha {
	cacheType := config.Toml.Cache.Type
	if cacheType == "memory" {
		return memory.GetCaptchaCache()
	} else {
		return memory.GetCaptchaCache()
	}
}

func GetDingTalkAccessTokenCache() cache.DingtalkAccessToken {
	cacheType := config.Toml.Cache.Type
	if cacheType == "memory" {
		return memory.GetDingTalkAccessTokenCache()
	} else {
		return memory.GetDingTalkAccessTokenCache()
	}
}
