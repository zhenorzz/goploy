package core

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
