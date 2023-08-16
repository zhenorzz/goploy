package memory

import (
	"github.com/zhenorzz/goploy/internal/cache"
	"sync"
	"time"
)

type UserCache struct {
	data map[string]user
	sync.RWMutex
}

type user struct {
	cache.UserData
	expireIn time.Time
}

var userCache = &UserCache{
	data: make(map[string]user),
}

func (uc *UserCache) IncErrorTimes(account string) int {
	uc.Lock()
	defer uc.Unlock()

	cacheKey := getCacheKey(account)

	times := 0
	v, ok := uc.data[cacheKey]
	if ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) {
		times = v.Times
	}

	times += 1

	uc.data[cacheKey] = user{
		UserData: cache.UserData{
			Times: times,
		},
		expireIn: time.Now().Add(cache.UserCacheExpireTime),
	}

	// show captcha
	showCaptchaKey := getShowCaptchaKey(account)
	uc.data[showCaptchaKey] = user{
		UserData: cache.UserData{
			Times: 1,
		},
		expireIn: time.Now().Add(cache.UserCacheShowCaptchaTime),
	}

	time.AfterFunc(cache.UserCacheShowCaptchaTime, func() {
		delete(uc.data, showCaptchaKey)
	})

	time.AfterFunc(cache.UserCacheExpireTime, func() {
		delete(uc.data, cacheKey)
	})

	return times
}

func (uc *UserCache) LockAccount(account string) {
	uc.Lock()
	defer uc.Unlock()

	lockKey := getLockKey(account)

	uc.data[lockKey] = user{
		UserData: cache.UserData{
			Times: 1,
		},
		expireIn: time.Now().Add(cache.UserCacheLockTime),
	}

	time.AfterFunc(cache.UserCacheLockTime, func() {
		delete(uc.data, lockKey)
	})

	cacheKey := getCacheKey(account)

	_, ok := uc.data[cacheKey]
	if ok {
		delete(uc.data, cacheKey)
	}
}

func (uc *UserCache) IsLock(account string) bool {
	uc.RLock()
	defer uc.RUnlock()

	lockKey := getLockKey(account)
	v, ok := uc.data[lockKey]

	return ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) && v.Times > 0
}

func (uc *UserCache) IsShowCaptcha(account string) bool {
	uc.RLock()
	defer uc.RUnlock()

	showCaptchaKey := getShowCaptchaKey(account)
	v, ok := uc.data[showCaptchaKey]

	return ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) && v.Times > 0
}

func (uc *UserCache) DeleteShowCaptcha(account string) {
	uc.Lock()
	defer uc.Unlock()

	delete(uc.data, getShowCaptchaKey(account))
}

func getCacheKey(account string) string {
	return cache.UserCacheKey + account
}

func getLockKey(account string) string {
	return cache.UserCacheLockKey + account
}

func getShowCaptchaKey(account string) string {
	return cache.UserCacheShowCaptchaKey + account
}

func GetUserCache() *UserCache {
	return userCache
}
