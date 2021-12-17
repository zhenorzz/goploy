package task

import (
	"database/sql"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/service"
	"github.com/zhenorzz/goploy/ws"
	"sync/atomic"
	"time"
)

var monitorTick = time.Tick(time.Second)
var monitorTaskDone = make(chan struct{})

func startMonitorTask() {
	atomic.AddInt32(&counter, 1)
	go func() {
		for {
			select {
			case <-monitorTick:
				monitorTask()
			case <-monitorTaskDone:
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

func shutdownMonitorTask() {
	close(monitorTaskDone)
}

type MonitorCache struct {
	errorTimes  int
	notifyTimes int
	time        int64
}

var monitorCaches = map[int64]MonitorCache{}

func monitorTask() {
	monitors, err := model.Monitor{State: model.Enable}.GetAllByState()
	if err != nil && err != sql.ErrNoRows {
		core.Log(core.ERROR, "get monitor list error, detail:"+err.Error())
	}
	monitorIDs := map[int64]struct{}{}
	for _, monitor := range monitors {
		monitorIDs[monitor.ID] = struct{}{}
		monitorCache, ok := monitorCaches[monitor.ID]
		if !ok {
			monitorCaches[monitor.ID] = MonitorCache{}
			monitorCache = monitorCaches[monitor.ID]
		}

		now := time.Now().Unix()

		if int(now-monitorCache.time) > monitor.Second {
			monitorCache.time = now
			if err := (service.Gnet{URL: monitor.URL}.Ping()); err != nil {
				monitorCache.errorTimes++
				core.Log(core.ERROR, "monitor "+monitor.Name+" encounter error, "+err.Error())
				if monitor.Times == uint16(monitorCache.errorTimes) {
					monitorCache.errorTimes = 0
					monitorCache.notifyTimes++
					body, err := monitor.Notify(err.Error())
					if err != nil {
						core.Log(core.ERROR, "monitor "+monitor.Name+" notify error, "+err.Error())
					} else {
						core.Log(core.TRACE, "monitor "+monitor.Name+" notify return "+body)
					}
					if monitor.NotifyTimes == uint16(monitorCache.notifyTimes) {
						monitorCache.notifyTimes = 0
						_ = monitor.TurnOff(err.Error())
						ws.GetHub().Data <- &ws.Data{
							Type:    ws.TypeMonitor,
							Message: ws.MonitorMessage{MonitorID: monitor.ID, State: ws.MonitorTurnOff, ErrorContent: err.Error()},
						}
					}
				}
			} else {
				monitorCache.errorTimes = 0
			}
			monitorCaches[monitor.ID] = monitorCache
		}
	}

	for cacheID := range monitorCaches {
		if _, ok := monitorIDs[cacheID]; !ok {
			delete(monitorCaches, cacheID)
		}
	}
}
