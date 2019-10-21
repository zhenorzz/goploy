package core

import (
	"strconv"
	"time"

	"goploy/model"

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
