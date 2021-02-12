package task

import (
	"context"
	"sync/atomic"
	"time"
)

var counter int32
var stop = make(chan struct{})

func Init() {
	atomic.AddInt32(&counter, 1)
	go ticker(stop)
}

func ticker(stop <-chan struct{}) {
	defer atomic.AddInt32(&counter, -1)
	// create ticker
	minute := time.Tick(time.Minute)
	second := time.Tick(time.Second)
	for {
		select {
		case <-second:
			monitorTask()
		case <-minute:
			projectTask()
		case <-stop:
			return
		}
	}
}

func Shutdown(ctx context.Context) error {
	close(stop)
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if atomic.LoadInt32(&counter) == 0 {
				return nil
			}
		}
	}
}
