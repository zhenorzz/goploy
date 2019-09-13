package core

import "goploy/utils"

// GolbalPath current path end with /
var GolbalPath = utils.GetCurrentPath()

// RepositoryPath repository path end with /
var RepositoryPath = GolbalPath + "repository/"

// PackagePath template path end with /
var PackagePath = RepositoryPath + "template-package/"

//role
const (
	RoleAdmin        = "admin"
	RoleManager      = "manager"
	RoleGroupManager = "group-manager"
	RoleMember       = "member"
)

// LoginCookieName jwt cookie name
const LoginCookieName = "goploy_token"
