package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/cmd/server/ws"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pipeline/docker"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/pkg/cmd"
	"github.com/zhenorzz/goploy/internal/transmitter"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func (gsync *Gsync) serverStage() error {
	ws.Send(ws.Data{
		Type:    ws.TypeProject,
		Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Sync"},
	})
	ch := make(chan syncMessage, len(gsync.ProjectServers))
	gsync.PublishTrace.Type = model.Deploy
	var serverSync = func(projectServer model.ProjectServer, index int) {
		project := gsync.Project
		publishTraceModel := gsync.PublishTrace
		publishTraceModel.InsertTime = time.Now().Format("20060102150405")
		// write after deploy script for rsync
		var dockerScript docker.Script
		scriptName := fmt.Sprintf("goploy-after-deploy-p%d-s%d.%s", project.ID, projectServer.ServerID, pkg.GetScriptExt(project.Script.AfterDeploy.Mode))
		scriptContent := ""
		if project.Script.AfterDeploy.Content != "" {
			scriptContent = project.ReplaceVars(project.Script.AfterDeploy.Content)
			scriptContent = project.ReplaceCustomVars(scriptContent)
			scriptContent = projectServer.ReplaceVars(scriptContent)
			scriptContent = strings.Replace(scriptContent, "${SERVER_TOTAL_NUMBER}", strconv.Itoa(len(gsync.ProjectServers)), -1)
			scriptContent = strings.Replace(scriptContent, "${SERVER_SERIAL_NUMBER}", strconv.Itoa(index), -1)

			if project.Script.AfterDeploy.Mode == "yaml" {
				_ = yaml.Unmarshal([]byte(scriptContent), &dockerScript)

				for stepIndex, step := range dockerScript.Steps {
					scriptName = fmt.Sprintf("goploy-after-deploy-p%d-s%d-y%d", project.ID, projectServer.ServerID, stepIndex)
					// delete the script
					step.Commands = append(step.Commands, fmt.Sprintf("rm -f %s", docker.GetDockerProjectScriptPath(project.ID, scriptName)))
					scriptContent = strings.Join(step.Commands, "\n")
					project.Script.AfterDeploy.ScriptNames = append(project.Script.AfterDeploy.ScriptNames, scriptName)
					_ = os.WriteFile(path.Join(config.GetProjectPath(project.ID), scriptName), []byte(scriptContent), 0755)
				}
			} else {
				project.Script.AfterDeploy.ScriptNames = append(project.Script.AfterDeploy.ScriptNames, scriptName)
				_ = os.WriteFile(path.Join(config.GetProjectPath(project.ID), scriptName), []byte(scriptContent), 0755)
			}
		}

		transmitterEntity := transmitter.New(project, projectServer)
		logCmd := transmitterEntity.String()
		log.Trace("projectID: " + strconv.FormatInt(project.ID, 10) + " " + logCmd)
		ext, _ := json.Marshal(struct {
			ServerID   int64  `json:"serverId"`
			ServerName string `json:"serverName"`
			Command    string `json:"command"`
		}{projectServer.ServerID, projectServer.Server.Name, logCmd})
		publishTraceModel.Ext = string(ext)

		if transmitterOutput, err := transmitterEntity.Exec(); err != nil {
			log.Error(fmt.Sprintf("projectID: %d transmit exec err: %s, output: %s", project.ID, err, transmitterOutput))
			publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", err, transmitterOutput)
			publishTraceModel.State = model.Fail
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Errorf(projectLogFormat, project.ID, err)
			}

			ch <- syncMessage{
				serverName: projectServer.Server.Name,
				projectID:  project.ID,
				detail:     err.Error(),
				state:      model.ProjectFail,
			}
			return
		} else {
			publishTraceModel.Detail = transmitterOutput
			publishTraceModel.State = model.Success
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Errorf(projectLogFormat, project.ID, err)
			}
		}

		if project.Script.AfterDeploy.Mode == "yaml" && len(dockerScript.Steps) > 0 {
			dockerConfig := docker.Config{
				ProjectID:   project.ID,
				ProjectPath: project.Path,
				Server:      projectServer.Server,
			}

			if err := dockerConfig.Setup(); err != nil {
				log.Error(fmt.Sprintf("projectID: %d err: %s", project.ID, err))
				publishTraceModel.Detail = err.Error()
				publishTraceModel.State = model.Fail
				if _, err := publishTraceModel.AddRow(); err != nil {
					log.Errorf(projectLogFormat, project.ID, err)
				}
				ch <- syncMessage{
					serverName: projectServer.Server.Name,
					projectID:  project.ID,
					detail:     err.Error(),
					state:      model.ProjectFail,
				}
				return
			}

			for stepIndex, step := range dockerScript.Steps {
				publishTraceModel.Type = model.AfterDeploy
				publishTraceModel.InsertTime = time.Now().Format("20060102150405")

				step.ScriptName = fmt.Sprintf("goploy-after-deploy-p%d-s%d-y%d", project.ID, projectServer.ServerID, stepIndex)
				dockerOutput, dockerErr := dockerConfig.Run(step)

				scriptContent = strings.Join(step.Commands, "\n")
				ext, _ = json.Marshal(struct {
					ServerID   int64  `json:"serverId"`
					ServerName string `json:"serverName"`
					Script     string `json:"script"`
					Step       string `json:"step"`
				}{projectServer.ServerID, projectServer.Server.Name, scriptContent, step.Name})
				publishTraceModel.Ext = string(ext)

				scriptFullName := path.Join(config.GetProjectPath(project.ID), step.ScriptName)
				_ = os.Remove(scriptFullName)

				if dockerErr != nil {
					log.Error(fmt.Sprintf("projectID: %d run docker script err: %s", project.ID, dockerErr))
					publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", dockerErr, dockerOutput)
					publishTraceModel.State = model.Fail
					if _, err := publishTraceModel.AddRow(); err != nil {
						log.Errorf(projectLogFormat, project.ID, err)
					}
					ch <- syncMessage{
						serverName: projectServer.Server.Name,
						projectID:  project.ID,
						detail:     fmt.Sprintf("err: %s\noutput: %s", dockerErr, dockerOutput),
						state:      model.ProjectFail,
					}
					return
				} else {
					publishTraceModel.Detail = dockerOutput
					publishTraceModel.State = model.Success
					if _, err := publishTraceModel.AddRow(); err != nil {
						log.Error("projectID: " + strconv.FormatInt(project.ID, 10) + " " + err.Error())
					}
				}
			}
		} else {
			var afterDeployCommands []string
			cmdEntity := cmd.New(projectServer.Server.OS)
			if len(project.SymlinkPath) != 0 {
				destDir := cmd.Join(project.SymlinkPath, project.LastPublishToken)
				afterDeployCommands = append(afterDeployCommands, cmdEntity.Symlink(destDir, project.Path))
			}

			if project.Script.AfterDeploy.Content != "" {
				afterDeployScriptPath := cmd.Join(project.Path, scriptName)
				afterDeployCommands = append(afterDeployCommands, cmdEntity.Script(project.Script.AfterDeploy.Mode, afterDeployScriptPath))
				afterDeployCommands = append(afterDeployCommands, cmdEntity.Remove(afterDeployScriptPath))
			}

			// no symlink and deploy script
			if len(afterDeployCommands) == 0 {
				ch <- syncMessage{
					serverName: projectServer.Server.Name,
					projectID:  project.ID,
					state:      model.ProjectSuccess,
				}
				return
			}
			completeAfterDeployCmd := strings.Join(afterDeployCommands, "&&")
			publishTraceModel.Type = model.AfterDeploy
			publishTraceModel.InsertTime = time.Now().Format("20060102150405")
			ext, _ = json.Marshal(struct {
				ServerID   int64  `json:"serverId"`
				ServerName string `json:"serverName"`
				Script     string `json:"script"`
			}{projectServer.ServerID, projectServer.Server.Name, scriptContent})
			publishTraceModel.Ext = string(ext)

			client, err := projectServer.ToSSHConfig().Dial()
			if err != nil {
				log.Error(err.Error())
				publishTraceModel.Detail = err.Error()
				publishTraceModel.State = model.Fail
				if _, err := publishTraceModel.AddRow(); err != nil {
					log.Errorf(projectLogFormat, project.ID, err)
				}
				ch <- syncMessage{
					serverName: projectServer.Server.Name,
					projectID:  project.ID,
					detail:     err.Error(),
					state:      model.ProjectFail,
				}
				return
			}
			defer client.Close()

			session, sessionErr := client.NewSession()
			if sessionErr != nil {
				log.Error(sessionErr.Error())
				publishTraceModel.Detail = sessionErr.Error()
				publishTraceModel.State = model.Fail
				if _, err := publishTraceModel.AddRow(); err != nil {
					log.Errorf(projectLogFormat, project.ID, err)
				}
				ch <- syncMessage{
					serverName: projectServer.Server.Name,
					projectID:  project.ID,
					detail:     sessionErr.Error(),
					state:      model.ProjectFail,
				}
				return
			}
			defer session.Close()

			log.Trace(fmt.Sprintf("projectID: %d ssh exec: %s", project.ID, completeAfterDeployCmd))

			output, err := session.CombinedOutput(completeAfterDeployCmd)
			if err != nil {
				log.Error(fmt.Sprintf("projectID: %d ssh exec err: %s, output: %s", project.ID, err, output))
				publishTraceModel.Detail = fmt.Sprintf("err: %s\noutput: %s", err, output)
				publishTraceModel.State = model.Fail
				if _, err := publishTraceModel.AddRow(); err != nil {
					log.Errorf(projectLogFormat, project.ID, err)
				}
				ch <- syncMessage{
					serverName: projectServer.Server.Name,
					projectID:  project.ID,
					detail:     fmt.Sprintf("%s\noutput: %s", err.Error(), output),
					state:      model.ProjectFail,
				}
				return
			}

			publishTraceModel.Detail = string(output)
			publishTraceModel.State = model.Success
			if _, err := publishTraceModel.AddRow(); err != nil {
				log.Error("projectID: " + strconv.FormatInt(project.ID, 10) + " " + err.Error())
			}
		}

		ch <- syncMessage{
			serverName: projectServer.Server.Name,
			projectID:  project.ID,
			state:      model.ProjectSuccess,
		}
	}

	for index, projectServer := range gsync.ProjectServers {
		if gsync.Project.DeployServerMode == "serial" {
			serverSync(projectServer, index+1)
		} else {
			go serverSync(projectServer, 0)
		}
	}

	message := ""
	for i := 0; i < len(gsync.ProjectServers); i++ {
		msg := <-ch
		if msg.state == model.ProjectFail {
			message += msg.serverName + " error message: " + msg.detail
		}
	}
	close(ch)

	if message != "" {
		return errors.New(message)
	}
	return nil
}
