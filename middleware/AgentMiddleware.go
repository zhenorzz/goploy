package middleware

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/core"
	"strconv"
	"time"
)

func CheckSign(gp *core.Goploy) error {
	sign := gp.URLQuery.Get("sign")
	if sign == "" {
		return errors.New("sign missing")
	}

	timestampStr := gp.URLQuery.Get("timestamp")
	if timestampStr == "" {
		return errors.New("timestamp missing")
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return errors.New("parse timestamp error")
	}

	currenTimestamp := time.Now().Unix()
	if currenTimestamp > timestamp+30 {
		return errors.New("request expired")
	}

	unsignedStr := string(gp.Body) + timestampStr + config.Toml.JWT.Key
	h := sha256.New()
	h.Write([]byte(unsignedStr))
	if sign != base64.URLEncoding.EncodeToString(h.Sum(nil)) {
		return errors.New("sign error")
	}
	return nil
}
