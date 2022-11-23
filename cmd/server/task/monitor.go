// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"database/sql"
	"github.com/zhenorzz/goploy/cmd/server/ws"
	"github.com/zhenorzz/goploy/internal/monitor"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/model"
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
		pkg.Log(pkg.ERROR, "get m list error, detail:"+err.Error())
	}
	monitorIDs := map[int64]struct{}{}
	for _, m := range monitors {
		monitorIDs[m.ID] = struct{}{}
		monitorCache, ok := monitorCaches[m.ID]
		if !ok {
			monitorCaches[m.ID] = MonitorCache{}
			monitorCache = monitorCaches[m.ID]
		}

		if monitorCache.silentCycle > 0 {
			monitorCache.silentCycle++
			if monitorCache.silentCycle >= m.SilentCycle {
				monitorCache.silentCycle = 0
			}
			monitorCaches[m.ID] = monitorCache
			continue
		}

		now := time.Now().Unix()
		println(m.Name, "detect", time.Now().String())
		if int(now-monitorCache.time) >= m.Second {
			println(m.Name, "in", time.Now().String())
			monitorCache.time = now
			ms, err := monitor.NewMonitorFromTarget(m.Type, m.Target)
			if err != nil {
				_ = m.TurnOff(err.Error())
				pkg.Log(pkg.ERROR, "m "+m.Name+" encounter error, "+err.Error())
				ws.GetHub().Data <- &ws.Data{
					Type:    ws.TypeMonitor,
					Message: ws.MonitorMessage{MonitorID: m.ID, State: ws.MonitorTurnOff, ErrorContent: err.Error()},
				}
			} else if err := ms.Check(); err != nil {
				monitorErrorContent := err.Error()
				monitorCache.errorTimes++
				pkg.Log(pkg.ERROR, "m "+m.Name+" encounter error, "+monitorErrorContent)
				if m.Times <= uint16(monitorCache.errorTimes) {
					if body, err := m.Notify(monitorErrorContent); err != nil {
						pkg.Log(pkg.ERROR, "m "+m.Name+" notify error, "+err.Error())
					} else {
						monitorCache.errorTimes = 0
						monitorCache.silentCycle = 1
						pkg.Log(pkg.TRACE, "m "+m.Name+" notify return "+body)
						_ = m.TurnOff(monitorErrorContent)
						ws.GetHub().Data <- &ws.Data{
							Type:    ws.TypeMonitor,
							Message: ws.MonitorMessage{MonitorID: m.ID, State: ws.MonitorTurnOff, ErrorContent: monitorErrorContent},
						}
					}
				}
			} else {
				monitorCache.errorTimes = 0
			}
			monitorCaches[m.ID] = monitorCache
		}
	}

	for cacheID := range monitorCaches {
		if _, ok := monitorIDs[cacheID]; !ok {
			delete(monitorCaches, cacheID)
		}
	}
}
