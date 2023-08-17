package cache

import "time"

type User interface {
	IncErrorTimes(account string, expireTime time.Duration, showCaptchaTime time.Duration) int
	LockAccount(account string, lockTime time.Duration)
	IsLock(account string) bool
	IsShowCaptcha(account string) bool
	DeleteShowCaptcha(account string)
}

const (
	UserCacheMaxErrorTimes   = 5
	UserCacheExpireTime      = 5 * time.Minute
	UserCacheLockTime        = 15 * time.Minute
	UserCacheShowCaptchaTime = 15 * time.Minute
)
