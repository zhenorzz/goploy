// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/cmd/server/ws"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/monitor"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var monitorTick = time.Tick(time.Second)
var monitorMutex sync.Mutex

func startMonitorTask() {
	atomic.AddInt32(&counter, 1)
	go func() {
		for {
			select {
			case <-monitorTick:
				if monitorMutex.TryLock() {
					monitorTask()
					monitorMutex.Unlock()
				}
			case <-stop:
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

type monitorMessage struct {
	MonitorID    int64  `json:"monitorId"`
	State        uint8  `json:"state"`
	ErrorContent string `json:"errorContent"`
}

func (monitorMessage) CanSendTo(*ws.Client) error {
	return nil
}

type MonitorCache struct {
	errorTimes     int
	time           int64
	silentCycle    int
	itemEditedTime string
}

var monitorCaches = map[int64]MonitorCache{}

func monitorTask() {
	monitors, err := model.Monitor{State: model.Enable}.GetAllByState()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("get monitor list error, detail:" + err.Error())
	}
	monitorIDs := map[int64]struct{}{}
	for _, m := range monitors {
		monitorIDs[m.ID] = struct{}{}
		monitorCache, ok := monitorCaches[m.ID]
		if !ok || monitorCache.itemEditedTime != m.UpdateTime {
			monitorCache = MonitorCache{itemEditedTime: m.UpdateTime}
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

		if int(now-monitorCache.time) < m.Second {
			continue
		}

		monitorCache.time = now
		ms := monitor.NewMonitorFromTarget(
			m.Type,
			m.Target,
			monitor.WithSuccessScript(m.SuccessServerID, m.SuccessScript),
			monitor.WithFailScript(m.FailServerID, m.FailScript),
		)

		if err := ms.Check(); err != nil {
			monitorErrorContent := err.Error()
			monitorCache.errorTimes++
			log.Error("m " + m.Name + " encounter error, " + monitorErrorContent)
			if m.Times <= uint16(monitorCache.errorTimes) {
				if body, err := m.Notify(monitorErrorContent); err != nil {
					log.Error("m " + m.Name + " notify error, " + err.Error())
				} else {
					monitorCache.errorTimes = 0
					monitorCache.silentCycle = 1
					log.Trace("m " + m.Name + " notify return " + body)
					_ = m.TurnOff(monitorErrorContent)
					ws.Send(ws.Data{
						Type:    ws.TypeMonitor,
						Message: monitorMessage{MonitorID: m.ID, State: model.Disable, ErrorContent: monitorErrorContent},
					})
				}
			}
			var serverID int64
			var se monitor.ScriptError
			if errors.As(err, &se) {
				serverID = se.ServerID
			}

			if err = ms.RunFailScript(serverID); err != nil {
				log.Error("Failed to run fail script " + err.Error())
			}
		} else {
			for _, item := range ms.Items {
				serverID, err := strconv.ParseInt(item, 10, 64)
				if err != nil {
					err = ms.RunSuccessScript(-1)
				} else {
					err = ms.RunSuccessScript(serverID)
				}
				if err != nil {
					log.Error("Failed to run successful script " + err.Error())
				}

			}
			monitorCache.errorTimes = 0

		}
		monitorCaches[m.ID] = monitorCache
	}

	for cacheID := range monitorCaches {
		if _, ok := monitorIDs[cacheID]; !ok {
			delete(monitorCaches, cacheID)
		}
	}
}
