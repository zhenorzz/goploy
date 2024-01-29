// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"os"
	"reflect"
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
	CORS     CORSConfig     `toml:"cors"`
	Captcha  CaptchaConfig  `toml:"captcha"`
	Cache    CacheConfig    `toml:"cache"`
}

type SetDefault interface {
	SetDefault()
}

var Toml Config
var Koanf = koanf.New(".")

func Init() {
	// If first time load config error, need to panic
	if err := setToml(); err != nil {
		panic(err)
	}

	v := reflect.ValueOf(&Toml)
	for i := 0; i < v.Elem().NumField(); i++ {
		fieldValue := v.Elem().Field(i)
		config := fieldValue.Addr().Interface()
		if c, ok := config.(SetDefault); ok {
			c.SetDefault()
		}
	}

	GetEventBus().Subscribe(APPEventTopic, &Toml.APP)
	GetEventBus().Subscribe(LogEventTopic, &Toml.Log)
	GetEventBus().Subscribe(DBEventTopic, &Toml.DB)

	err := file.Provider(GetConfigFile()).Watch(func(event interface{}, err error) {
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		oldToml := Toml

		if err = setToml(); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err = PublishEvents(Toml, getEventTopics(oldToml, Toml)); err != nil {
			// If new config publish events error, use the old config
			fmt.Printf("publish config events error: %v \n", err)
			errToml := Toml
			Toml = oldToml
			_ = PublishEvents(Toml, getEventTopics(errToml, Toml))
		}
	})

	if err != nil {
		panic(err)
	}
}

func getEventTopics(oldToml Config, newToml Config) (topics []string) {
	if oldToml.DB != newToml.DB {
		topics = append(topics, DBEventTopic)
	}

	if oldToml.Log != newToml.Log {
		topics = append(topics, LogEventTopic)
	}

	if oldToml.APP != newToml.APP {
		topics = append(topics, APPEventTopic)
	}

	return topics
}

func setToml() error {
	if err := Koanf.Load(file.Provider(GetConfigFile()), toml.Parser()); err != nil {
		return fmt.Errorf("load config file error: %s", err)
	}

	if err := Koanf.Unmarshal("", &Toml); err != nil {
		return fmt.Errorf("unmarshal config error: %s", err)
	}

	return nil
}

func Write(cfg Config) error {
	if err := Koanf.Load(structs.Provider(cfg, "toml"), nil); err != nil {
		return err
	}

	b, _ := Koanf.Marshal(toml.Parser())

	err := os.WriteFile(GetConfigFile(), b, 0644)
	if err != nil {
		return err
	}
	return nil
}
