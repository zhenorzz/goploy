package utils

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type GIT struct {
	Dir    string
	Output bytes.Buffer
	Err    bytes.Buffer
}

func (git *GIT) Run(operator string, options ...string) error {
	git.Output.Reset()
	git.Err.Reset()
	cmd := exec.Command("git", append([]string{operator}, options...)...)
	if len(git.Dir) != 0 {
		cmd.Dir = git.Dir
	}
	cmd.Stdout = &git.Output
	cmd.Stderr = &git.Err
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Clone(options ...string) error {
	if err := git.Run("clone", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Clean(options ...string) error {
	if err := git.Run("clean", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Checkout(options ...string) error {
	if err := git.Run("checkout", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Add(options ...string) error {
	if err := git.Run("add", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Pull(options ...string) error {
	if err := git.Run("pull", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Fetch(options ...string) error {
	if err := git.Run("fetch", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Log(options ...string) error {
	if err := git.Run("log", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Branch(options ...string) error {
	if err := git.Run("branch", options...); err != nil {
		return err
	}
	return nil
}

func (git *GIT) Reset(options ...string) error {
	if err := git.Run("reset", options...); err != nil {
		return err
	}
	return nil
}

type Commit struct {
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	Author    string `json:"author"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	Tag       string `json:"tag"`
	Diff      string `json:"diff"`
}

func ParseGITLog(rawCommitLog string) []Commit {
	unformatCommitList := strings.Split(rawCommitLog, "`start`")
	unformatCommitList = unformatCommitList[1:]
	var commitList []Commit
	for _, commitRow := range unformatCommitList {
		commitRowSplit := strings.Split(commitRow, "`")
		timestamp, _ := strconv.Atoi(commitRowSplit[2])
		commitList = append(commitList, Commit{
			Commit:    commitRowSplit[0],
			Author:    commitRowSplit[1],
			Timestamp: timestamp,
			Message:   commitRowSplit[3],
			Tag:       commitRowSplit[4],
			Diff:      strings.Trim(commitRowSplit[5], "\n"),
		})
	}
	return commitList
}
