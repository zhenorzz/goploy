package task

import (
	"database/sql"
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"strings"
	"time"
)

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
	if err != nil && err != sql.ErrNoRows {
		core.Log(core.ERROR, "get server monitor list error, detail:"+err.Error())
	}
	for _, serverMonitor := range serverMonitorTasks {
		monitorCache, ok := serverMonitorCaches[serverMonitor.ID]
		if !ok {
			serverMonitorCaches[serverMonitor.ID] = ServerMonitorCache{}
			monitorCache = serverMonitorCaches[serverMonitor.ID]
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
			core.Log(core.ERROR, "get cycle value failed, detail:"+err.Error())
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
					core.Log(core.ERROR, fmt.Sprintf("monitor task %d has no server, detail: %s", serverMonitor.ID, err.Error()))
					continue
				}
				serverCaches[serverMonitor.ServerID] = server
			}
			body, err := serverMonitor.Notify(serverCaches[serverMonitor.ServerID], cycleValue)
			if err != nil {
				core.Log(core.ERROR, fmt.Sprintf("monitor task %d notify error, %s", serverMonitor.ID, err.Error()))
			} else {
				core.Log(core.TRACE, fmt.Sprintf("monitor task %d notify return %s", serverMonitor.ID, body))
			}
		}
		serverMonitorCaches[serverMonitor.ID] = monitorCache
	}
}
