// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repo

import (
	"crypto/tls"
	"fmt"
	"github.com/jlaffaye/ftp"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type FtpRepo struct{}

func (FtpRepo) CanRollback() bool {
	return false
}

func (ftpRepo FtpRepo) Ping(url string) error {
	c, err := ftpRepo.dial(url)
	if err != nil {
		return err
	}
	_ = c.Quit()
	return nil
}

// Create -
func (ftpRepo FtpRepo) Create(projectID int64) error {
	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		log.Error(fmt.Sprintf("The project does not exist, projectID:%d", projectID))
		return err
	}
	return ftpRepo.Follow(project, "")
}

func (ftpRepo FtpRepo) Follow(project model.Project, _ string) error {
	projectID := project.ID
	srcPath := config.GetProjectPath(projectID)
	_ = os.RemoveAll(srcPath)
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		log.Error(fmt.Sprintf("The project fail to mkdir, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}

	c, err := ftpRepo.dial(project.URL)
	if err != nil {
		log.Error(fmt.Sprintf("The project fail to connect ftp, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	var downloadFromFTP func(localDir, remoteDir string) error
	downloadFromFTP = func(localDir, remoteDir string) error {
		remoteEntries, err := c.List(remoteDir)
		if err != nil {
			return err
		}
		for _, entry := range remoteEntries {
			if entry.Type == 1 {
				nextLocalDir := path.Join(localDir, entry.Name)
				if err := os.Mkdir(nextLocalDir, 0755); err != nil {
					return err
				}
				if err := downloadFromFTP(nextLocalDir, path.Join(remoteDir, entry.Name)); err != nil {
					return err
				}
			} else {
				remoteFile, err := c.Retr(path.Join(remoteDir, entry.Name))
				if err != nil {
					return err
				}
				localFile, err := os.Create(path.Join(localDir, entry.Name))
				if err != nil {
					remoteFile.Close()
					return err
				}
				_, err = io.Copy(localFile, remoteFile)
				if err != nil {
					return err
				}
				localFile.Close()
				remoteFile.Close()
			}
		}
		return nil
	}
	if err := downloadFromFTP(srcPath, ""); err != nil {
		log.Error(fmt.Sprintf("The project fail to download file, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	_ = c.Quit()
	log.Trace(fmt.Sprintf("The project success to download, projectID:%d", projectID))
	return nil
}

func (FtpRepo) RemoteBranchList(url string) ([]string, error) {
	return []string{"virtual"}, nil
}

func (FtpRepo) BranchList(projectID int64) ([]string, error) {
	return []string{"virtual"}, nil
}

func (FtpRepo) CommitLog(projectID int64, rows int) ([]CommitInfo, error) {
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

func (ftpRepo FtpRepo) BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error) {
	return []CommitInfo{{Commit: "virtual"}}, nil
}

func (ftpRepo FtpRepo) TagLog(projectID int64, rows int) ([]CommitInfo, error) {
	return []CommitInfo{}, nil
}

func (FtpRepo) dial(_url string) (*ftp.ServerConn, error) {
	var (
		enableTLS bool
		addr      string
		username  = "anonymous"
		password  = "anonymous@domain.com"
		dir       string
	)
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "ftps" {
		enableTLS = true
	}
	h := strings.Split(u.Host, ":")
	if len(h) == 1 {
		u.Host += ":21"
	}
	addr = u.Host
	dir = u.Path
	if u.User.Username() != "" {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	opts := []ftp.DialOption{
		ftp.DialWithTimeout(10 * time.Second),
	}

	if enableTLS == true {
		opts = append(opts, ftp.DialWithTLS(&tls.Config{
			InsecureSkipVerify: true,
		}))
	}

	c, err := ftp.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	if err = c.Login(username, password); err != nil {
		return nil, err
	}
	if dir != "" {
		if err = c.ChangeDir(dir); err != nil {
			return nil, err
		}
	}
	return c, nil
}
