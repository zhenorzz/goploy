// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	Env      string         `toml:"env"`
	APP      APPConfig      `toml:"app"`
	Cookie   CookieConfig   `toml:"cookie"`
	JWT      JWTConfig      `toml:"jwt"`
	DB       DBConfig       `toml:"db"`
	Log      LogConfig      `toml:"log"`
	Web      WebConfig      `toml:"web"`
	LDAP     LDAPConfig     `toml:"ldap"`
	Dingtalk DingtalkConfig `toml:"dingtalk"`
	Feishu   FeishuConfig   `toml:"feishu"`
}

type APPConfig struct {
	DeployLimit     int32         `toml:"deployLimit"`
	ShutdownTimeout time.Duration `toml:"shutdownTimeout"`
	RepositoryPath  string        `toml:"repositoryPath"`
}

type CookieConfig struct {
	Name   string `toml:"name"`
	Expire int    `toml:"expire"` // second
}

type JWTConfig struct {
	Key string `toml:"key"`
}

type DBConfig struct {
	Type     string `toml:"type"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
}

type LogConfig struct {
	Path string `toml:"path"`
}

type WebConfig struct {
	Port string `toml:"port"`
}

type LDAPConfig struct {
	Enabled    bool   `toml:"enabled"`
	URL        string `toml:"url"`
	BindDN     string `toml:"bindDN"`
	Password   string `toml:"password"`
	BaseDN     string `toml:"baseDN"`
	UID        string `toml:"uid"`
	UserFilter string `toml:"userFilter"`
}

type DingtalkConfig struct {
	AppKey    string `toml:"appKey"`
	AppSecret string `toml:"appSecret"`
}

type FeishuConfig struct {
	AppKey    string `toml:"appKey"`
	AppSecret string `toml:"appSecret"`
}

var Toml Config

func InitToml() {
	config, err := os.ReadFile(GetConfigFile())
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(config, &Toml)
	if err != nil {
		panic(err)
	}
	setAPPDefault()
	setDBDefault()
	setLogger()
}

func setAPPDefault() {
	if Toml.APP.ShutdownTimeout == 0 {
		Toml.APP.ShutdownTimeout = 10
	}
}

func setDBDefault() {
	if Toml.DB.Type == "" {
		Toml.DB.Type = "mysql"
	}
	if Toml.DB.Host == "" {
		Toml.DB.Host = "127.0.0.1"
	}
	if Toml.DB.Port == "" {
		Toml.DB.Port = "3306"
	}
	if Toml.DB.Database == "" {
		Toml.DB.Database = "goploy"
	}
}

func setLogger() {
	var logFile io.Writer
	logPathEnv := Toml.Log.Path
	if strings.ToLower(logPathEnv) == "stdout" {
		logFile = os.Stdout
	} else {
		logPath, err := filepath.Abs(logPathEnv)
		if err != nil {
			fmt.Println(err.Error())
		}
		if _, err := os.Stat(logPath); err != nil && os.IsNotExist(err) {
			if err := os.Mkdir(logPath, os.ModePerm); nil != err {
				panic(err.Error())
			}
		}
		logFile, err = os.OpenFile(logPath+"/goploy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if nil != err {
			panic(err.Error())
		}
	}
	log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s()", path.Base(f.Function)), fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	})

	log.SetOutput(logFile)

	log.SetLevel(log.TraceLevel)

}

func Write(cfg Config) error {
	yamlData, err := toml.Marshal(&cfg)

	if err != nil {
		return err
	}

	err = os.WriteFile(GetConfigFile(), yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}
