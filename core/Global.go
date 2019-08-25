package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GolbalPath current path end with /
var GolbalPath = getCurrentPath()

// RepositoryPath repository path end with /
var RepositoryPath = GolbalPath + "repository/"

// TemplatePath template path end with /
var TemplatePath = RepositoryPath + "template-package/"

// getCurrentPath if env = 'production' return absolute else return relative
func getCurrentPath() string {
	if os.Getenv("ENV") != "production" {
		return "./"
	}
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	path, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		panic(err)
	}
	return string(path[0 : i+1])
}
