package core

import "github.com/zhenorzz/goploy/utils"

// GlobalPath current path end with /
var GlobalPath = utils.GetCurrentPath()

// RepositoryPath repository path end with /
var RepositoryPath = GlobalPath + "repository/"

// PackagePath template path end with /
var PackagePath = RepositoryPath + "template-package/"

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
