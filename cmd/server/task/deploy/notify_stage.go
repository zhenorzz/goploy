package deploy

import (
	"fmt"
	"github.com/zhenorzz/goploy/internal/notify"
)

// commit id
// commit message
// server ip & name
// deploy user name
// deploy time
func (gsync *Gsync) notify(deployState uint8, detail string) {
	if gsync.Project.NotifyType == 0 {
		return
	}
	project := gsync.Project

	_ = notify.Send(fmt.Sprintf("project%d-deploy", project.ID), notify.UseByDeploy, notify.DeployData{
		DeployState:    deployState,
		Project:        project,
		ProjectServers: gsync.ProjectServers,
		CommitInfo:     gsync.CommitInfo,
		DeployDetail:   detail,
	}, project.NotifyType, project.NotifyTarget)

}
