package utils

import (
	"bytes"
	"os/exec"
)

type SVN struct {
	Dir    string
	Output bytes.Buffer
	Err    bytes.Buffer
}

func (svn *SVN) Run(operator string, options ...string) error {
	svn.Output.Reset()
	svn.Err.Reset()
	cmd := exec.Command("svn", append([]string{operator}, options...)...)
	if len(svn.Dir) != 0 {
		cmd.Dir = svn.Dir
	}
	cmd.Stdout = &svn.Output
	cmd.Stderr = &svn.Err
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (svn *SVN) Clone(options ...string) error {
	if err := svn.Run("co", options...); err != nil {
		return err
	}
	return nil
}

func (svn *SVN) Pull(options ...string) error {
	if err := svn.Run("up", options...); err != nil {
		return err
	}
	return nil
}

func (svn *SVN) Log(options ...string) error {
	if err := svn.Run("log", options...); err != nil {
		return err
	}
	return nil
}

func (svn *SVN) LS(options ...string) error {
	if err := svn.Run("ls", options...); err != nil {
		return err
	}
	return nil
}
