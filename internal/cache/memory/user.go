package memory

import (
	"sync"
	"time"
)

const (
	UserCacheKey     = "login_error_times_"
	UserCacheLockKey = "login_lock_"
)

type UserCache struct {
	data map[string]user
	sync.RWMutex
}

type user struct {
	times    int
	data     any
	expireIn time.Time
}

var userCache = &UserCache{
	data: make(map[string]user),
}

func (uc *UserCache) IncrErrorTimes(account string, expireTime time.Duration) int {
	uc.Lock()
	defer uc.Unlock()

	cacheKey := getCacheKey(account)

	times := 0
	v, ok := uc.data[cacheKey]
	if ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) {
		times = v.times
	}

	times += 1

	uc.data[cacheKey] = user{
		times:    times,
		expireIn: time.Now().Add(expireTime),
	}
	time.AfterFunc(expireTime, func() {
		delete(uc.data, cacheKey)
	})

	return times
}

func (uc *UserCache) LockAccount(account string, lockTime time.Duration) {
	uc.Lock()
	defer uc.Unlock()

	lockKey := getLockKey(account)

	uc.data[lockKey] = user{
		times:    1,
		expireIn: time.Now().Add(lockTime),
	}

	time.AfterFunc(lockTime, func() {
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

	return ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) && v.times > 0
}

func (uc *UserCache) IsShowCaptcha(account string) bool {
	uc.RLock()
	defer uc.RUnlock()

	v, ok := uc.data[getCacheKey(account)]

	return ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) && v.times > 0
}

func (uc *UserCache) DeleteErrorTimes(account string) {
	uc.Lock()
	defer uc.Unlock()

	delete(uc.data, getCacheKey(account))
}

func getCacheKey(account string) string {
	return UserCacheKey + account
}

func getLockKey(account string) string {
	return UserCacheLockKey + account
}

func GetUserCache() *UserCache {
	return userCache
}
