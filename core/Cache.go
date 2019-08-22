package core

import (
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/model"

	cache "github.com/patrickmn/go-cache"
)

// Cache uint
var Cache = cache.New(24*time.Hour, 48*time.Hour)

// GetUserCache return model.User, model.Permissions, error
func GetUserCache(UserID uint32) (model.User, model.Permissions, error) {
	var userData model.User
	var permissions model.Permissions
	var err error

	userData, err = GetUserData(UserID)
	if err != nil {
		return userData, permissions, err
	}

	if x, found := Cache.Get("permissions:" + strconv.Itoa(int(UserID))); found {
		permissions = *x.(*model.Permissions)
	} else {
		// 超级管理员获取全部权限
		if userData.RoleID == 1 {
			permissions, err = model.Permission{}.GetAll()
			if err != nil {
				return userData, permissions, err
			}
		} else {
			role, err := model.Role{ID: userData.RoleID}.GetData()
			if err != nil {
				return userData, permissions, err
			}
			permissions, err = model.Permission{}.GetAllByPermissionList(role.PermissionList)
			if err != nil {
				return userData, permissions, err
			}
			Cache.Set("permissions:"+strconv.Itoa(int(UserID)), &permissions, cache.DefaultExpiration)
		}
	}
	return userData, permissions, nil
}

// GetUserData return model.User and error
func GetUserData(UserID uint32) (model.User, error) {
	var userData model.User
	var err error
	if x, found := Cache.Get("userInfo:" + strconv.Itoa(int(UserID))); found {
		userData = *x.(*model.User)
	} else {
		userData, err = model.User{ID: UserID}.GetData()
		if err != nil {
			return userData, err
		}

		Cache.Set("userInfo:"+strconv.Itoa(int(UserID)), &userData, cache.DefaultExpiration)
	}
	return userData, nil
}
