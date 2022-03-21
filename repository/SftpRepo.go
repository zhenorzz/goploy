package repository

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"golang.org/x/crypto/ssh"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type SftpRepo struct{}

func (sftpRepo SftpRepo) Ping(url string) error {
	sshClient, err := sftpRepo.dial(url)
	if err != nil {
		return err
	}
	_ = sshClient.Close()
	return nil
}

func (sftpRepo SftpRepo) Create(projectID int64) error {
	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project does not exist, projectID:%d", projectID))
		return err
	}
	return sftpRepo.Follow(project, "")
}

func (sftpRepo SftpRepo) Follow(project model.Project, _ string) error {
	projectID := project.ID
	srcPath := core.GetProjectPath(projectID)
	_ = os.RemoveAll(srcPath)
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to mkdir, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}

	sshClient, err := sftpRepo.dial(project.URL)
	if err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to connect ftp, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	defer sshClient.Close()
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	var downloadFromSFTP func(localDir, remoteDir string) error
	downloadFromSFTP = func(localDir, remoteDir string) error {
		remoteEntries, err := sftpClient.ReadDir(remoteDir)
		if err != nil {
			return err
		}
		for _, entry := range remoteEntries {
			nextLocalDir := path.Join(localDir, entry.Name())
			nextRemoteDir := path.Join(remoteDir, entry.Name())
			if entry.Mode()&os.ModeSymlink != 0 {
				entry, err = sftpClient.Stat(nextRemoteDir)
				if err != nil {
					return err
				}
			}

			if entry.IsDir() {
				if err = os.Mkdir(nextLocalDir, 0755); err != nil {
					return err
				}
				if err = downloadFromSFTP(nextLocalDir, nextRemoteDir); err != nil {
					return err
				}
			} else {
				remoteFile, err := sftpClient.Open(nextRemoteDir)
				if err != nil {
					return err
				}
				localFile, err := os.Create(nextLocalDir)
				if err != nil {
					remoteFile.Close()
					return err
				}
				_, err = io.Copy(localFile, remoteFile)
				if err != nil {
					println(nextLocalDir, nextRemoteDir)
					return err
				}
				localFile.Close()
				remoteFile.Close()
			}
		}
		return nil
	}
	_url, _, _ := sftpRepo.parseURL(project.URL)
	u, err := url.Parse(_url)
	if err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to parse url, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	if err := downloadFromSFTP(srcPath, u.Path); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to download file, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	core.Log(core.TRACE, fmt.Sprintf("The project success to download, projectID:%d", projectID))
	return nil
}

func (SftpRepo) RemoteBranchList(url string) ([]string, error) {
	return []string{"virtual"}, nil
}

func (SftpRepo) BranchList(projectID int64) ([]string, error) {
	return []string{"virtual"}, nil
}

func (SftpRepo) CommitLog(projectID int64, rows int) ([]CommitInfo, error) {
	commitInfo := CommitInfo{
		Branch:    "virtual",
		Commit:    "",
		Author:    "",
		Timestamp: time.Now().Unix(),
		Message:   "",
		Tag:       "",
		Diff:      "",
	}
	return []CommitInfo{commitInfo}, nil
}

func (sftpRepo SftpRepo) BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error) {
	return []CommitInfo{{Commit: "virtual"}}, nil
}

func (sftpRepo SftpRepo) TagLog(projectID int64, rows int) ([]CommitInfo, error) {
	return []CommitInfo{}, nil
}

func (sftpRepo SftpRepo) dial(rawURL string) (*ssh.Client, error) {
	var (
		_url    = rawURL
		host    string
		port    = 22
		user    = "root"
		keyFile = "/root/.ssh/id_rsa"
	)
	_url, user, keyFile = sftpRepo.parseURL(rawURL)
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	h := strings.Split(u.Host, ":")
	if len(h) == 1 {
		host = u.Host
	} else {
		host = h[0]
		port, _ = strconv.Atoi(h[1])
	}

	client, err := utils.SSHConfig{
		User: user,
		Path: keyFile,
		Host: host,
		Port: port,
	}.Dial()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (SftpRepo) parseURL(rawURL string) (_url, user, keyFile string) {
	_url = rawURL
	user = "root"
	keyFile = "/root/.ssh/id_rsa"
	urlSplit := strings.Split(rawURL, " ")
	for _, item := range urlSplit {
		if item == "" {
		} else if strings.HasPrefix(item, "--user=") {
			userSplit := strings.Split(item, "=")
			user = userSplit[1]
		} else if strings.HasPrefix(item, "--keyFile=") {
			keyFileSplit := strings.Split(item, "=")
			keyFile = keyFileSplit[1]
		} else {
			_url = item
		}
	}
	return
}
