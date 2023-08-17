package memory

import (
	"sync"
	"time"
)

type CaptchaCache struct {
	data map[string]captcha
	sync.RWMutex
}

type captcha struct {
	dots     interface{}
	expireIn time.Time
}

var captchaCache = &CaptchaCache{
	data: make(map[string]captcha),
}

func (c *CaptchaCache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if !v.expireIn.IsZero() && v.expireIn.After(time.Now()) {
		return v.dots, true
	}

	return nil, false
}

func (c *CaptchaCache) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	var expireIn time.Time

	if ttl > 0 {
		expireIn = time.Now().Add(ttl)
	}

	c.data[key] = captcha{
		dots:     value,
		expireIn: expireIn,
	}

	time.AfterFunc(ttl, func() {
		delete(c.data, key)
	})
}

func (c *CaptchaCache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.data, key)
}

func (c *CaptchaCache) IsChecked(key string) bool {
	if key == "" {
		return false
	}

	if check, ok := c.Get(key); ok {
		_check, _ := check.(bool)
		return _check
	}

	return false
}

func GetCaptchaCache() *CaptchaCache {
	return captchaCache
}
