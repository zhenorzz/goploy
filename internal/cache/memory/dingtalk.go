package memory

import (
	"github.com/zhenorzz/goploy/internal/cache"
	"sync"
	"time"
)

type AccessTokenCache struct {
	data  map[string]accessToken
	mutex sync.RWMutex
}

type accessToken struct {
	cache.AccessTokenData
	expireIn time.Time
}

var accessTokenCache = &AccessTokenCache{
	data: make(map[string]accessToken),
}

func (ac *AccessTokenCache) Get(key string) (string, bool) {
	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	v, ok := ac.data[key]
	if !ok {
		return "", false
	}

	if !v.expireIn.IsZero() && v.expireIn.After(time.Now()) {
		return v.AccessToken, true
	}

	return "", false
}

func (ac *AccessTokenCache) Set(key string, value string, ttl time.Duration) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	var expireIn time.Time

	if ttl > 0 {
		expireIn = time.Now().Add(ttl)
	}

	ac.data[key] = accessToken{
		AccessTokenData: cache.AccessTokenData{
			AccessToken: value,
		},
		expireIn: expireIn,
	}

	time.AfterFunc(ttl, func() {
		delete(ac.data, key)
	})
}

func GetDingTalkAccessTokenCache() *AccessTokenCache {
	return accessTokenCache
}
