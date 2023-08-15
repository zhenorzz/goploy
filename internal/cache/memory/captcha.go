package memory

import (
	"sync"
	"time"
)

type CaptchaCache struct {
	c     map[string]captcha
	mutex sync.RWMutex
}

type captcha struct {
	dots     interface{}
	expireIn time.Time
}

var captchaCache *CaptchaCache
var captchaOnce sync.Once

func (c *CaptchaCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	v, ok := c.c[key]
	if !ok {
		return nil, false
	}

	if !v.expireIn.IsZero() && v.expireIn.After(time.Now()) {
		return v.dots, true
	}

	return nil, false
}

func (c *CaptchaCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expireIn time.Time

	if ttl > 0 {
		expireIn = time.Now().Add(ttl)
	}

	c.c[key] = captcha{
		dots:     value,
		expireIn: expireIn,
	}
}

func (c *CaptchaCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.c, key)
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

func (c *CaptchaCache) CleanExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, val := range c.c {
		if val.expireIn.Before(time.Now()) {
			delete(c.c, key)
		}
	}
}

func GetCaptchaCache() *CaptchaCache {
	captchaOnce.Do(func() {
		captchaCache = &CaptchaCache{
			c: make(map[string]captcha),
		}
	})

	return captchaCache
}
