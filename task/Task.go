package task

import "time"

func Init() {
	go ticker()
}

func ticker() {
	// 创建一个计时器
	minute := time.Tick(time.Minute)
	second := time.Tick(time.Second)
	for {
		select {
		case <-second:
			monitorTask()
		case <-minute:
			projectTask()
		}
	}
}
