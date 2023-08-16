package cache

import "time"

type Captcha interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	IsChecked(key string) bool
}
