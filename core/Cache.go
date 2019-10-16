package core

import (
	"strconv"
	"time"

	"goploy/model"

	cache "github.com/patrickmn/go-cache"
)

// Cache uint
var Cache = cache.New(24*time.Hour, 48*time.Hour)

// GetUserInfo return model.User and error
func GetUserInfo(userID int64) (model.User, error) {
	var userData model.User
	var err error
	if x, found := Cache.Get("userInfo:" + strconv.Itoa(int(userID))); found {
		userData = *x.(*model.User)
	} else {
		userData, err = model.User{ID: userID}.GetData()
		if err != nil {
			return userData, err
		}

		Cache.Set("userInfo:"+strconv.Itoa(int(userID)), &userData, cache.DefaultExpiration)
	}
	return userData, nil
}

// GetUserInfo return model.User and error
func GetUserProject(userID int64) (model.ProjectUsers, error) {
	var projectUsers model.ProjectUsers
	var err error
	if x, found := Cache.Get("projectUser:" + strconv.Itoa(int(userID))); found {
		projectUsers = *x.(*model.ProjectUsers)
	} else {
		projectUsers, err = model.ProjectUser{ID: userID}.GetListByUserID()
		if err != nil {
			return projectUsers, err
		}

		Cache.Set("projectUser:"+strconv.Itoa(int(userID)), &projectUsers, cache.DefaultExpiration)
	}
	return projectUsers, nil
}

func UserHasProject(userID int64, projectID int64) bool {
	projectUsers, err := GetUserProject(userID)
	if err != nil {
		return false
	}

	for _, projectUser := range projectUsers {
		if projectUser.UserID == userID {
			return true
		}
	}

	return false
}
