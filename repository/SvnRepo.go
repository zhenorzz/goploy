// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repository

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

type SvnRepo struct{}

// Ping -
func (SvnRepo) Ping(url string) error {
	svn := utils.SVN{}
	if err := svn.LS(strings.Split(url, " ")...); err != nil {
		return errors.New(svn.Err.String())
	}
	return nil
}

// Create -
func (SvnRepo) Create(projectID int64) error {
	srcPath := core.GetProjectPath(projectID)
	if _, err := os.Stat(srcPath); err == nil {
		return nil
	}
	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		core.Log(core.TRACE, "The project does not exist, projectID:"+strconv.FormatInt(projectID, 10))
		return err
	}
	if err := os.RemoveAll(srcPath); err != nil {
		core.Log(core.TRACE, "The project fail to remove, projectID:"+strconv.FormatInt(project.ID, 10)+" ,error: "+err.Error())
		return err
	}
	svn := utils.SVN{}
	options := strings.Split(project.URL, " ")
	options = append(options, srcPath)
	if err := svn.Clone(options...); err != nil {
		core.Log(core.ERROR, "The project fail to initialize, projectID:"+strconv.FormatInt(project.ID, 10)+" ,error: "+err.Error()+", detail: "+svn.Err.String())
		return err
	}
	core.Log(core.TRACE, "The project success to initialize, projectID:"+strconv.FormatInt(project.ID, 10))
	return nil
}

func (svnRepo SvnRepo) Follow(project model.Project, target string) error {
	if err := svnRepo.Create(project.ID); err != nil {
		return err
	}
	svn := utils.SVN{Dir: core.GetProjectPath(project.ID)}

	// the length of commit id is 40
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" svn up")
	if strings.Index(target, "r") == 0 {
		if err := svn.Pull("-r", target); err != nil {
			core.Log(core.ERROR, err.Error()+", detail: "+svn.Err.String())
			return errors.New(svn.Err.String())
		}
	} else {
		if err := svn.Pull(); err != nil {
			core.Log(core.ERROR, err.Error()+", detail: "+svn.Err.String())
			return errors.New(svn.Err.String())
		}
	}

	return nil
}

func (SvnRepo) RemoteBranchList(url string) ([]string, error) {
	return []string{"master"}, nil
}

func (SvnRepo) BranchList(projectID int64) ([]string, error) {
	return []string{"master"}, nil
}

func (SvnRepo) CommitLog(projectID int64, rows int) ([]CommitInfo, error) {
	svn := utils.SVN{Dir: core.GetProjectPath(projectID)}

	if err := svn.Log("-v", "--xml", "-l", strconv.Itoa(rows)); err != nil {
		return []CommitInfo{}, errors.New(svn.Err.String())
	}

	list := parseSVNLog(svn.Output.String())
	return list, nil
}

func (SvnRepo) BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error) {
	svn := utils.SVN{Dir: core.GetProjectPath(projectID)}

	if err := svn.Log("-v", "--xml", "-l", strconv.Itoa(rows)); err != nil {
		return []CommitInfo{}, errors.New(svn.Err.String())
	}

	list := parseSVNLog(svn.Output.String())
	return list, nil
}

func (SvnRepo) TagLog(projectID int64, rows int) ([]CommitInfo, error) {
	return []CommitInfo{}, nil
}

func parseSVNLog(rawCommitLog string) []CommitInfo {
	type path struct {
		Action   string `xml:"action,attr"`
		PropMods bool   `xml:"prop-mods,attr"`
		TextMods bool   `xml:"text-mods,attr"`
		Kind     string `xml:"kind,attr"`
		Path     string `xml:",chardata"`
	}

	type logEntry struct {
		Revision string `xml:"revision,attr"`
		Author   string `xml:"author"`
		Date     string `xml:"date"`
		Msg      string `xml:"msg"`
		Paths    []path `xml:"paths>path"`
	}

	type log struct {
		XMLName  xml.Name   `xml:"log"`
		LogEntry []logEntry `xml:"logentry"`
	}

	logs := new(log)
	err := xml.Unmarshal([]byte(rawCommitLog), logs)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	var commitInfos []CommitInfo

	for _, entry := range logs.LogEntry {
		formatTime, err := time.Parse("2006-01-02T15:04:05.999999999Z", entry.Date)
		var timestamp int64
		if err == nil {
			timestamp = formatTime.Unix()
		}
		commitInfo := CommitInfo{
			Branch:    "master",
			Commit:    "r" + entry.Revision,
			Author:    entry.Author,
			Timestamp: timestamp,
			Message:   entry.Msg,
		}

		diff := "Changed paths:\n"

		for _, iPath := range entry.Paths {
			diff += iPath.Action + "\t" + iPath.Path
		}

		commitInfo.Diff = diff

		commitInfos = append(commitInfos, commitInfo)
	}

	return commitInfos
}
