package deploy

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/internal/model"
	"strconv"
	"sync"
)

// keep the latest 10 project
func (gsync *Gsync) removeExpiredBackup() {
	if gsync.Project.SymlinkPath == "" {
		return
	}
	var wg sync.WaitGroup
	for _, projectServer := range gsync.ProjectServers {
		wg.Add(1)
		go func(projectServer model.ProjectServer) {
			defer wg.Done()
			client, err := projectServer.ToSSHConfig().Dial()
			if err != nil {
				log.Error(err.Error())
				return
			}
			defer func() {
				if err := client.Close(); err != nil {
					log.Error(err.Error())
				}
			}()
			session, err := client.NewSession()
			if err != nil {
				log.Error(err.Error())
				return
			}
			defer func() {
				if err := session.Close(); err != nil {
					log.Error(err.Error())
				}
			}()
			var sshOutbuf, sshErrbuf bytes.Buffer
			session.Stdout = &sshOutbuf
			session.Stderr = &sshErrbuf
			if err = session.Run("cd " + gsync.Project.SymlinkPath + ";ls -t | awk 'NR>" + strconv.Itoa(int(gsync.Project.SymlinkBackupNumber)) + "' | xargs rm -rf"); err != nil {
				log.Error(err.Error())
			}
		}(projectServer)
	}
	wg.Wait()
}
