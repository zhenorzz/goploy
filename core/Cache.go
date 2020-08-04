package core

import (
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/model"

	"github.com/patrickmn/go-cache"
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
func GetNamespace(userID int64) (model.Namespaces, error) {
	var namespaceList model.Namespaces
	var err error
	if x, found := Cache.Get("namespace:" + strconv.Itoa(int(userID))); found {
		namespaceList = *x.(*model.Namespaces)
	} else {
		namespaceList, err = model.Namespace{UserID: userID}.GetAllByUserID()
		if err != nil {
			return namespaceList, err
		}

		Cache.Set("namespace:"+strconv.Itoa(int(userID)), &namespaceList, cache.DefaultExpiration)
	}
	return namespaceList, nil
}
