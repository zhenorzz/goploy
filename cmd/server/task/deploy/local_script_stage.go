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
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (gsync *Gsync) afterPullScriptStage() error {
	if gsync.Project.Script.AfterPull.Content == "" {
		return nil
	}

	gsync.PublishTrace.Type = model.AfterPull
	gsync.PublishTrace.InsertTime = time.Now().Format("20060102150405")
	return gsync.runLocalScript()
}

func (gsync *Gsync) deployFinishScriptStage() error {
	if gsync.Project.Script.DeployFinish.Content == "" {
		return nil
	}

	gsync.PublishTrace.Type = model.DeployFinish
	gsync.PublishTrace.InsertTime = time.Now().Format("20060102150405")
	return gsync.runLocalScript()
}

func (gsync *Gsync) runLocalScript() error {
	var mode = ""
	var content = ""
	var scriptName = ""
	project := gsync.Project
	switch gsync.PublishTrace.Type {
	case model.AfterPull:
		ws.Send(ws.Data{
			Type:    ws.TypeProject,
			Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Run pull script"},
		})
		mode = project.Script.AfterPull.Mode
		content = project.Script.AfterPull.Content
		scriptName = "goploy-after-pull"

	case model.DeployFinish:
		ws.Send(ws.Data{
			Type:    ws.TypeProject,
			Message: deployMessage{ProjectID: gsync.Project.ID, ProjectName: gsync.Project.Name, State: Deploying, Message: "Run finish script"},
		})
		mode = project.Script.DeployFinish.Mode
		content = project.Script.DeployFinish.Content
		scriptName = "goploy-deploy-finish"

	default:
		return errors.New("not support stage")
	}

	log.Tracef(projectLogFormat, gsync.Project.ID, content)

	commitInfo := gsync.CommitInfo
	srcPath := config.GetProjectPath(project.ID)
	scriptMode := "bash"
	if mode != "" {
		scriptMode = mode
	}
	scriptText := commitInfo.ReplaceVars(content)
	scriptText = project.ReplaceVars(scriptText)
	scriptText = gsync.ProjectServers.ReplaceVars(scriptText)
	scriptText = project.ReplaceCustomVars(scriptText)

	// run yaml script by docker
	if mode == "yaml" {
		var dockerScript docker.Script
		_ = yaml.Unmarshal([]byte(scriptText), &dockerScript)

		projectPath, err := filepath.Abs(config.GetProjectPath(project.ID))
		if err != nil {
			return fmt.Errorf("get repository abs path err: %s", err)
		}

		if len(dockerScript.Steps) == 0 {
			return nil
		}

		dockerConfig := docker.Config{
			ProjectID:   project.ID,
			ProjectPath: projectPath,
		}

		if err = dockerConfig.Setup(); err != nil {
			return err
		}

		for stepIndex, step := range dockerScript.Steps {
			gsync.PublishTrace.InsertTime = time.Now().Format("20060102150405")
			scriptText = strings.Join(step.Commands, "\n")
			tmpScriptName := scriptName + fmt.Sprintf("-y%d", stepIndex)
			scriptFullName := path.Join(srcPath, tmpScriptName)
			step.ScriptName = tmpScriptName

			_ = os.WriteFile(scriptFullName, []byte(scriptText), 0755)

			dockerOutput, dockerErr := dockerConfig.Run(step)

			ext, _ := json.Marshal(struct {
				Script string `json:"script"`
				Step   string `json:"step"`
			}{scriptText, step.Name})
			gsync.PublishTrace.Ext = string(ext)

			_ = os.Remove(scriptFullName)

			if dockerErr != nil {
				gsync.PublishTrace.Detail = dockerErr.Error()
				gsync.PublishTrace.State = model.Fail
				if _, err := gsync.PublishTrace.AddRow(); err != nil {
					log.Errorf(projectLogFormat, gsync.Project.ID, err)
				}
				return fmt.Errorf("run docker script err: %s", dockerErr)
			} else {
				gsync.PublishTrace.Detail = dockerOutput
				gsync.PublishTrace.State = model.Success
				if _, err := gsync.PublishTrace.AddRow(); err != nil {
					log.Errorf(projectLogFormat, gsync.Project.ID, err)
				}
			}
		}
	} else {
		scriptName += fmt.Sprintf(".%s", pkg.GetScriptExt(mode))
		scriptFullName := path.Join(srcPath, scriptName)
		_ = os.WriteFile(scriptFullName, []byte(scriptText), 0755)

		var commandOptions []string
		if scriptMode == "cmd" {
			commandOptions = append(commandOptions, "/C")
			scriptFullName, _ = filepath.Abs(scriptFullName)
		}
		commandOptions = append(commandOptions, scriptFullName)

		ext, _ := json.Marshal(struct {
			Script string `json:"script"`
		}{Script: scriptText})
		gsync.PublishTrace.Ext = string(ext)

		handler := exec.Command(scriptMode, commandOptions...)
		handler.Dir = srcPath

		if output, err := handler.CombinedOutput(); err != nil {
			gsync.PublishTrace.Detail = fmt.Sprintf("err: %s\noutput: %s", err, string(output))
			gsync.PublishTrace.State = model.Fail
			if _, err := gsync.PublishTrace.AddRow(); err != nil {
				log.Errorf(projectLogFormat, gsync.Project.ID, err)
			}
			return fmt.Errorf("err: %s, output: %s", err, string(output))
		} else {
			_ = os.Remove(scriptFullName)
			gsync.PublishTrace.Detail = string(output)
			gsync.PublishTrace.State = model.Success
			if _, err := gsync.PublishTrace.AddRow(); err != nil {
				log.Errorf(projectLogFormat, gsync.Project.ID, err)
			}
		}
	}

	return nil
}
