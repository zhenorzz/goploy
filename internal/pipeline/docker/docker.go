package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"path"
	"path/filepath"
)

type Script struct {
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name         string   `yaml:"name"`
	Commands     []string `yaml:"commands"`
	Image        string   `yaml:"image"`
	ImageOptions struct {
		Registry   string `yaml:"registry"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Dockerfile string `yaml:"dockerfile"`
	} `yaml:"imageOptions"`
	ScriptName    string
	ContainerName string
}

type Config struct {
	ProjectID   int64
	ProjectPath string
	Server      model.Server
	Client      client.APIClient
}

func GetDockerProjectPath(projectID int64) string {
	return fmt.Sprintf("/data/www/repository/project_%d", projectID)
}

func GetDockerProjectScriptPath(projectID int64, scriptName string) string {
	return path.Join(GetDockerProjectPath(projectID), scriptName)
}

func (c *Config) Setup() (err error) {
	if c.Server.IP != "" {
		c.Client, err = client.NewClientWithOpts(client.WithHost(fmt.Sprintf("tcp://%s:2375", c.Server.IP)), client.WithAPIVersionNegotiation())
	} else {
		c.Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	}
	if err != nil {
		return fmt.Errorf("connect docker err: %s", err)
	}

	return nil
}

func (c *Config) Run(step Step) (outStr string, errStr string) {
	defer func() {
		c.Destroy(step)
	}()

	step.ContainerName = uuid.New().String()
	ctx := context.Background()
	pullOptions := types.ImagePullOptions{}

	if step.ImageOptions.Registry != "" {
		step.Image = fmt.Sprintf("%s%s", step.ImageOptions.Registry, step.Image)
	}

	// build image by dockerfile
	if step.ImageOptions.Dockerfile != "" {
		localProjectPath, err := filepath.Abs(config.GetProjectPath(c.ProjectID))
		if err != nil {
			errStr = fmt.Sprintf("get local repository abs path err: %s", err)
			return
		}

		tar, err := archive.TarWithOptions(filepath.Join(localProjectPath, step.ImageOptions.Dockerfile), &archive.TarOptions{})

		_, err = c.Client.ImageBuild(ctx, tar, types.ImageBuildOptions{
			Tags:        []string{step.Image},
			Dockerfile:  "Dockerfile",
			Remove:      true,
			ForceRemove: true,
		})
		if err != nil {
			errStr = fmt.Sprintf("build image err: %s", err)
			return
		}
	} else {
		// pull image from private registry
		if step.ImageOptions.Username != "" && step.ImageOptions.Password != "" {
			authConfig := registry.AuthConfig{
				Username: step.ImageOptions.Username,
				Password: step.ImageOptions.Password,
			}
			authConfigBytes, _ := json.Marshal(authConfig)
			authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
			pullOptions.RegistryAuth = authConfigEncoded
		}

		_, err := c.Client.ImagePull(ctx, step.Image, pullOptions)
		if err != nil {
			errStr = fmt.Sprintf("pull docker image err: %s", err)
			return
		}
	}

	dockerProjectPath := GetDockerProjectPath(c.ProjectID)

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: c.ProjectPath,
				Target: dockerProjectPath,
			},
		},
	}

	_, err := c.Client.ContainerCreate(ctx, &container.Config{
		Image:        step.Image,
		Cmd:          []string{GetDockerProjectScriptPath(c.ProjectID, step.ScriptName)},
		Entrypoint:   []string{"/bin/sh"},
		WorkingDir:   dockerProjectPath,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
	}, hostConfig, nil, nil, step.ContainerName)

	// if err is image does not exist, re-pull image and re-create container
	if client.IsErrNotFound(err) && step.ImageOptions.Dockerfile == "" {
		_, err = c.Client.ImagePull(ctx, step.Image, pullOptions)
		if err != nil {
			errStr = fmt.Sprintf("pull docker image twice err: %s", err)
			return
		}

		_, err = c.Client.ContainerCreate(ctx, &container.Config{
			Image:        step.Image,
			Cmd:          []string{GetDockerProjectScriptPath(c.ProjectID, step.ScriptName)},
			Entrypoint:   []string{"/bin/sh"},
			WorkingDir:   dockerProjectPath,
			AttachStdin:  false,
			AttachStdout: true,
			AttachStderr: true,
		}, hostConfig, nil, nil, step.ContainerName)
	}

	if err != nil {
		errStr = fmt.Sprintf("create docker container err: %s", err)
		return
	}

	if err := c.Client.ContainerStart(ctx, step.ContainerName, types.ContainerStartOptions{}); err != nil {
		errStr = fmt.Sprintf("start docker container err: %s", err)
		return
	}

	statusCh, errCh := c.Client.ContainerWait(ctx, step.ContainerName, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			errStr = fmt.Sprintf("wait docker container err: %s", err)
			break
		}
	case <-statusCh:
	}

	if errStr != "" {
		return
	}

	out, err := c.Client.ContainerLogs(ctx, step.ContainerName, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		errStr = fmt.Sprintf("logs docker container err: %s", err)
		return
	}

	var dockerOutbuf, dockerErrbuf bytes.Buffer
	stdcopy.StdCopy(&dockerOutbuf, &dockerErrbuf, out)
	defer out.Close()

	if dockerErrbuf.Len() > 0 {
		errStr = fmt.Sprintf("run docker script err: %s", dockerErrbuf.String())
		return
	}

	return dockerOutbuf.String(), errStr
}

func (c *Config) Destroy(step Step) {
	ctx := context.Background()
	_ = c.Client.ContainerKill(ctx, step.ContainerName, "9")
	_ = c.Client.ContainerRemove(ctx, step.ContainerName, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	})
}
