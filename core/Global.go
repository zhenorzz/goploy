package core

import "github.com/zhenorzz/goploy/utils"

// GolbalPath current path end with /
var GolbalPath = utils.GetCurrentPath()

// RepositoryPath repository path end with /
var RepositoryPath = GolbalPath + "repository/"

// TemplatePath template path end with /
var TemplatePath = RepositoryPath + "template-package/"

//role
const (
	RoleAdmin        = "admin"
	RoleManager      = "manager"
	RoleGroupManager = "group-manager"
	RoleMember       = "member"
)
