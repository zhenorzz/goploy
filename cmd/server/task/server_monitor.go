// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/internal/model"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var serverMonitorTick = time.Tick(time.Minute)
var serverMonitorMutex sync.Mutex

func startServerMonitorTask() {
	atomic.AddInt32(&counter, 1)
	go func() {
		for {
			select {
			case <-serverMonitorTick:
				if serverMonitorMutex.TryLock() {
					serverMonitorTask()
					serverMonitorMutex.Unlock()
				}
			case <-stop:
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

type ServerMonitorCache struct {
	lastCycle   int
	silentCycle int
}

var serverMonitorCaches = map[int64]ServerMonitorCache{}

var loop = 0

func serverMonitorTask() {
	loop++
	var serverCaches = map[int64]model.Server{}
	serverMonitorTasks, err := model.ServerMonitor{}.GetAllModBy(loop, time.Now().Format("15:04"))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("get server monitor list error, detail:" + err.Error())
	}
	for _, serverMonitor := range serverMonitorTasks {
		monitorCache, ok := serverMonitorCaches[serverMonitor.ID]
		if !ok {
			monitorCache = ServerMonitorCache{}
		}

		if monitorCache.silentCycle > 0 {
			monitorCache.silentCycle++
			if monitorCache.silentCycle >= serverMonitor.SilentCycle {
				monitorCache.silentCycle = 0
			}
			serverMonitorCaches[serverMonitor.ID] = monitorCache
			continue
		}

		cycleValue, err := model.ServerAgentLog{
			ServerID: serverMonitor.ServerID,
			Item:     serverMonitor.Item,
		}.GetCycleValue(serverMonitor.GroupCycle, serverMonitor.Formula)
		if err != nil {
			log.Error("get cycle value failed, detail:" + err.Error())
			continue
		}

		compareRes := false
		switch serverMonitor.Operator {
		case ">=":
			compareRes = strings.Compare(cycleValue, serverMonitor.Value) >= 0
		case ">":
			compareRes = strings.Compare(cycleValue, serverMonitor.Value) > 0
		case "<=":
			compareRes = strings.Compare(cycleValue, serverMonitor.Value) <= 0
		case "<":
			compareRes = strings.Compare(cycleValue, serverMonitor.Value) < 0
		case "!=":
			compareRes = strings.Compare(cycleValue, serverMonitor.Value) != 0
		}
		if compareRes {
			monitorCache.lastCycle++
		} else {
			monitorCache.lastCycle = 0
		}

		if monitorCache.lastCycle >= serverMonitor.LastCycle {
			monitorCache.silentCycle = 1
			monitorCache.lastCycle = 0

			if _, ok := serverCaches[serverMonitor.ServerID]; !ok {
				server, err := model.Server{ID: serverMonitor.ServerID}.GetData()
				if err != nil {
					log.Error(fmt.Sprintf("monitor task %d has no server, detail: %s", serverMonitor.ID, err.Error()))
					continue
				}
				serverCaches[serverMonitor.ServerID] = server
			}
			body, err := serverMonitor.Notify(serverCaches[serverMonitor.ServerID], cycleValue)
			if err != nil {
				log.Error(fmt.Sprintf("monitor task %d notify error, %s", serverMonitor.ID, err.Error()))
			} else {
				log.Trace(fmt.Sprintf("monitor task %d notify return %s", serverMonitor.ID, body))
			}
		}
		serverMonitorCaches[serverMonitor.ID] = monitorCache
	}
}
