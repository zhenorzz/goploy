// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"database/sql"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/service"
	"github.com/zhenorzz/goploy/ws"
	"sync/atomic"
	"time"
)

var monitorTick = time.Tick(time.Minute)

func startMonitorTask() {
	atomic.AddInt32(&counter, 1)
	go func() {
		for {
			select {
			case <-monitorTick:
				monitorTask()
			case <-stop:
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

type MonitorCache struct {
	errorTimes  int
	time        int64
	silentCycle int
}

var monitorCaches = map[int64]MonitorCache{}

func monitorTask() {
	monitors, err := model.Monitor{State: model.Enable}.GetAllByState()
	if err != nil && err != sql.ErrNoRows {
		pkg.Log(pkg.ERROR, "get monitor list error, detail:"+err.Error())
	}
	monitorIDs := map[int64]struct{}{}
	for _, monitor := range monitors {
		monitorIDs[monitor.ID] = struct{}{}
		monitorCache, ok := monitorCaches[monitor.ID]
		if !ok {
			monitorCaches[monitor.ID] = MonitorCache{}
			monitorCache = monitorCaches[monitor.ID]
		}

		if monitorCache.silentCycle > 0 {
			monitorCache.silentCycle++
			if monitorCache.silentCycle >= monitor.SilentCycle {
				monitorCache.silentCycle = 0
			}
			monitorCaches[monitor.ID] = monitorCache
			continue
		}

		now := time.Now().Unix()
		println(monitor.Name, "detect", time.Now().String())
		if int(now-monitorCache.time) >= monitor.Second {
			println(monitor.Name, "in", time.Now().String())
			monitorCache.time = now
			ms, err := service.NewMonitorFromTarget(monitor.Type, monitor.Target)
			if err != nil {
				_ = monitor.TurnOff(err.Error())
				pkg.Log(pkg.ERROR, "monitor "+monitor.Name+" encounter error, "+err.Error())
				ws.GetHub().Data <- &ws.Data{
					Type:    ws.TypeMonitor,
					Message: ws.MonitorMessage{MonitorID: monitor.ID, State: ws.MonitorTurnOff, ErrorContent: err.Error()},
				}
			} else if err := ms.Check(); err != nil {
				monitorErrorContent := err.Error()
				monitorCache.errorTimes++
				pkg.Log(pkg.ERROR, "monitor "+monitor.Name+" encounter error, "+monitorErrorContent)
				if monitor.Times <= uint16(monitorCache.errorTimes) {
					if body, err := monitor.Notify(monitorErrorContent); err != nil {
						pkg.Log(pkg.ERROR, "monitor "+monitor.Name+" notify error, "+err.Error())
					} else {
						monitorCache.errorTimes = 0
						monitorCache.silentCycle = 1
						pkg.Log(pkg.TRACE, "monitor "+monitor.Name+" notify return "+body)
						_ = monitor.TurnOff(monitorErrorContent)
						ws.GetHub().Data <- &ws.Data{
							Type:    ws.TypeMonitor,
							Message: ws.MonitorMessage{MonitorID: monitor.ID, State: ws.MonitorTurnOff, ErrorContent: monitorErrorContent},
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
