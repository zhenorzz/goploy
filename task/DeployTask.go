// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package task

import (
	"container/list"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/service"
	"github.com/zhenorzz/goploy/ws"
	"sync"
	"sync/atomic"
	"time"
)

var deployList = list.New()
var deployTick = time.Tick(time.Millisecond)

func startDeployTask() {
	atomic.AddInt32(&counter, 1)
	var deployingNumber int32
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case <-deployTick:
				if atomic.LoadInt32(&deployingNumber) < config.Toml.APP.DeployLimit {
					atomic.AddInt32(&deployingNumber, 1)
					if deployElem := deployList.Front(); deployElem != nil {
						wg.Add(1)
						go func(gsync service.Gsync) {
							gsync.Exec()
							atomic.AddInt32(&deployingNumber, -1)
							wg.Done()
						}(deployList.Remove(deployElem).(service.Gsync))
					} else {
						atomic.AddInt32(&deployingNumber, -1)
					}
				}
			case <-stop:
				wg.Wait()
				atomic.AddInt32(&counter, -1)
				return
			}
		}
	}()
}

func AddDeployTask(gsync service.Gsync) {
	ws.GetHub().Data <- &ws.Data{
		Type: ws.TypeProject,
		Message: ws.ProjectMessage{
			ProjectID:   gsync.Project.ID,
			ProjectName: gsync.Project.Name,
			State:       ws.TaskWaiting,
			Message:     "Task waiting",
			Ext: struct {
				LastPublishToken string `json:"lastPublishToken"`
			}{gsync.Project.LastPublishToken},
		},
	}
	model.PublishTrace{
		Token:         gsync.Project.LastPublishToken,
		ProjectID:     gsync.Project.ID,
		ProjectName:   gsync.Project.Name,
		PublisherID:   gsync.UserInfo.ID,
		PublisherName: gsync.UserInfo.Name,
		Type:          model.QUEUE,
		State:         model.Success,
	}.AddRow()
	deployList.PushBack(gsync)
}
