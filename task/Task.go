package task

import "time"

func Init() {
	go ticker()
}

func ticker() {
	// create ticker
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
