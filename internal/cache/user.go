package cache

import "time"

type User interface {
	IncErrorTimes(account string) int
	LockAccount(account string)
	IsLock(account string) bool
	IsShowCaptcha(account string) bool
	DeleteShowCaptcha(account string)
}

const (
	UserCacheKey             = "login_error_times_"
	UserCacheLockKey         = "login_lock_"
	UserCacheShowCaptchaKey  = "login_show_captcha_"
	UserCacheMaxErrorTimes   = 5
	UserCacheExpireTime      = 5 * time.Minute
	UserCacheLockTime        = 15 * time.Minute
	UserCacheShowCaptchaTime = 15 * time.Minute
)

type UserData struct {
	Times int
}
