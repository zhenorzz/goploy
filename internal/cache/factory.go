package cache

import (
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/cache/memory"
)

const MemoryCache = "memory"

var cacheType = config.Toml.Cache.Type

func GetUserCache() User {
	switch cacheType {
	case MemoryCache:
		return memory.GetUserCache()
	default:
		return memory.GetUserCache()
	}
}

func GetCaptchaCache() Captcha {
	switch cacheType {
	case MemoryCache:
		return memory.GetCaptchaCache()
	default:
		return memory.GetCaptchaCache()
	}
}

func GetDingTalkAccessTokenCache() DingtalkAccessToken {
	switch cacheType {
	case MemoryCache:
		return memory.GetDingTalkAccessTokenCache()
	default:
		return memory.GetDingTalkAccessTokenCache()
	}
}
