// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"io"
	log1 "log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TRACE   = "TRACE: "
	WARNING = "WARNING: "
	INFO    = "INFO: "
	ERROR   = "ERROR: "
)

func log(lv string, content string) {
	var logFile io.Writer
	logPathEnv := config.Toml.Log.Path
	if strings.ToLower(logPathEnv) == "stdout" {
		logFile = os.Stdout
	} else {
		logPath, err := filepath.Abs(logPathEnv)
		if err != nil {
			fmt.Println(err.Error())
		}
		if _, err := os.Stat(logPath); err != nil && os.IsNotExist(err) {
			err := os.Mkdir(logPath, os.ModePerm)
			if nil != err {
				fmt.Println(err.Error())
			}
		}
		file := logPath + "/"
		if config.Toml.Log.Split {
			file += time.Now().Format("20060102") + ".log"
		} else {
			file += "goploy.log"
		}
		logFile, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if nil != err {
			fmt.Println(err.Error())
		}
	}

	logger := log1.New(logFile, lv, log1.LstdFlags|log1.Lshortfile)
	logger.Output(2, content)
}

func logf(lv string, format string, a ...interface{}) {
	log(lv, fmt.Sprintf(format, a...))
}

func Error(s string) {
	log(ERROR, s)
}

func Errorf(s string, a ...interface{}) {
	logf(ERROR, s, a...)
}

func Warning(s string) {
	log(WARNING, s)
}

func Warningf(s string, a ...interface{}) {
	logf(WARNING, s, a...)
}

func Trace(s string) {
	log(TRACE, s)
}

func Tracef(s string, a ...interface{}) {
	logf(TRACE, s, a...)
}

func Info(s string) {
	log(INFO, s)
}

func Infof(s string, a ...interface{}) {
	logf(INFO, s, a...)
}
