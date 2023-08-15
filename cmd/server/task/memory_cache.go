// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"github.com/zhenorzz/goploy/internal/cache/memory"
	"sync/atomic"
	"time"
)

var memoryCacheTick = time.Tick(time.Minute)

func startMemoryCacheTask() {
	atomic.AddInt32(&counter, 1)
	go func() {
		for {
			select {
			case <-memoryCacheTick:
				memoryCacheTask()
			case <-stop:
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

func memoryCacheTask() {
	memory.GetUserCache().CleanExpired()
	memory.GetCaptchaCache().CleanExpired()
	memory.GetDingTalkAccessTokenCache().CleanExpired()
}
