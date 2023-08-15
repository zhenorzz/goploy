package memory

import (
	"sync"
	"time"
)

const (
	Key           = "login_error_times_"
	LockKey       = "login_lock_"
	MaxErrorTimes = 5
	ExpireTime    = 5 * time.Minute
	LockTime      = 15 * time.Minute
)

type UserCache struct {
	c     map[string]user
	mutex sync.RWMutex
}

type user struct {
	times    int
	expireIn time.Time
}

var userCache *UserCache
var userOnce sync.Once

func (uc *UserCache) IncErrorTimes(account string) int {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	cacheKey := getCacheKey(account)

	times := 0
	v, ok := uc.c[cacheKey]
	if ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) {
		times = v.times
	}

	times += 1

	uc.c[cacheKey] = user{
		times:    times,
		expireIn: time.Now().Add(ExpireTime),
	}

	return times
}

func (uc *UserCache) LockAccount(account string) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	lockKey := getLockKey(account)

	uc.c[lockKey] = user{
		times:    1,
		expireIn: time.Now().Add(LockTime),
	}

	cacheKey := getCacheKey(account)

	_, ok := uc.c[cacheKey]
	if ok {
		delete(uc.c, cacheKey)
	}
}

func (uc *UserCache) IsLock(account string) bool {
	uc.mutex.RLock()
	defer uc.mutex.RUnlock()

	lockKey := getLockKey(account)
	v, ok := uc.c[lockKey]

	return ok && !v.expireIn.IsZero() && v.expireIn.After(time.Now()) && v.times > 0
}

func (uc *UserCache) CleanExpired() {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	for key, val := range uc.c {
		if val.expireIn.Before(time.Now()) {
			delete(uc.c, key)
		}
	}
}

func getCacheKey(account string) string {
	return Key + account
}

func getLockKey(account string) string {
	return LockKey + account
}

func GetUserCache() *UserCache {
	userOnce.Do(func() {
		userCache = &UserCache{
			c: make(map[string]user),
		}
	})

	return userCache
}
