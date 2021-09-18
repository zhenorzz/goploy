package core

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

//role
const (
	RoleAdmin        = "admin"
	RoleManager      = "manager"
	RoleGroupManager = "group-manager"
	RoleMember       = "member"
)

//Roles all role
var Roles = [...]string{RoleAdmin, RoleManager, RoleGroupManager, RoleMember}

// LoginCookieName jwt cookie name
const LoginCookieName = "goploy_token"

// NamespaceHeaderName namespace cookie name
const NamespaceHeaderName = "G-N-ID"

var (
	AssetDir string
	Gwg      sync.WaitGroup
)

// GetAssetDir if env = 'production' return absolute else return relative
func GetAssetDir() string {
	if AssetDir != "" {
		return AssetDir
	}
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	app, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	i := strings.LastIndex(app, "/")
	if i < 0 {
		i = strings.LastIndex(app, "\\")
	}
	if i < 0 {
		panic(err)
	}
	return app[0 : i+1]
}

func GetEnvFile() string {
	return path.Join(GetAssetDir(), ".env")
}

func GetRepositoryPath() string {
	return path.Join(GetAssetDir(), "repository")
}

func GetProjectFilePath(projectID int64) string {
	return path.Join(GetRepositoryPath(), "project-file", "project_"+strconv.FormatInt(projectID, 10))
}

func GetProjectPath(projectID int64) string {
	return path.Join(GetRepositoryPath(), "project_"+strconv.FormatInt(projectID, 10))
}
