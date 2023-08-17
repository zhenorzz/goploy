package cache

import "time"

type DingtalkAccessToken interface {
	Get(key string) (string, bool)
	Set(key string, value string, ttl time.Duration)
}
