package task

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

var counter int32
var stop = make(chan struct{})

func Init() {
	atomic.AddInt32(&counter, 1)
	go ticker()
}

func ticker() {
	defer atomic.AddInt32(&counter, -1)
	// create ticker
	minute := time.Tick(time.Minute)
	second := time.Tick(time.Second)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			select {
			case <-second:
				monitorTask()
			case <-stop:
				wg.Done()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-minute:
				projectTask()
				serverMonitorTask()
			case <-stop:
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
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
