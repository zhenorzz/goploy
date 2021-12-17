package task

import (
	"context"
	"sync/atomic"
	"time"
)

var counter int32

func Init() {
	startMonitorTask()
	startProjectTask()
	startServerMonitorTask()
	startDeployTask()
}

func Shutdown(ctx context.Context) error {
	shutdownMonitorTask()
	shutdownProjectTask()
	shutdownServerMonitorTask()
	shutdownDeployTask()
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
