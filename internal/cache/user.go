package cache

import "time"

type User interface {
	IncrErrorTimes(account string, expireTime time.Duration) int
	LockAccount(account string, lockTime time.Duration)
	IsLock(account string) bool
	IsShowCaptcha(account string) bool
	DeleteErrorTimes(account string)
}

const (
	UserCacheMaxErrorTimes = 5
	UserCacheExpireTime    = 5 * time.Minute
	UserCacheLockTime      = 15 * time.Minute
)
