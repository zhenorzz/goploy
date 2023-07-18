package dingtalk

import (
	"sync"
	"time"
)

type AccessTokenCache struct {
	c     map[string]accessToken
	mutex sync.RWMutex
}

type accessToken struct {
	accessToken string
	expireIn    time.Time
}

var accessTokenCache *AccessTokenCache

func (ac *AccessTokenCache) Get(key string) (string, bool) {
	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	v, ok := ac.c[key]
	if !ok {
		return "", false
	}

	if !v.expireIn.IsZero() && v.expireIn.After(time.Now()) {
		return v.accessToken, true
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

	ac.c[key] = accessToken{
		accessToken: value,
		expireIn:    expireIn,
	}
}

func GetCache() *AccessTokenCache {
	if accessTokenCache == nil {
		accessTokenCache = &AccessTokenCache{
			c: make(map[string]accessToken),
		}
	}

	return accessTokenCache
}
