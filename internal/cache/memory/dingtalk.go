package memory

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
var dingtalkOnce sync.Once

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

func (ac *AccessTokenCache) CleanExpired() {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	for key, val := range ac.c {
		if val.expireIn.Before(time.Now()) {
			delete(ac.c, key)
		}
	}
}

func GetDingTalkAccessTokenCache() *AccessTokenCache {
	dingtalkOnce.Do(func() {
		accessTokenCache = &AccessTokenCache{
			c: make(map[string]accessToken),
		}
	})

	return accessTokenCache
}
