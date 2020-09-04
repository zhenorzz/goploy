package core

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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

// NamespaceCookieName namespace cookie name
const NamespaceCookieName = "goploy_namespace"

// GetCurrentPath if env = 'production' return absolute else return relative
func GetAppPath() string {
	if os.Getenv("ENV") != "production" {
		return "./"
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

func GetRepositoryPath() string {
	return path.Join(GetAppPath(), "repository")
}

func GetPackagePath() string {
	return path.Join(GetRepositoryPath(), "template-package")
}

func GetProjectPath(projectID int64) string {
	return path.Join(GetRepositoryPath(), "project_"+strconv.FormatInt(projectID, 10))
}
