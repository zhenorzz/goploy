package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel is log level
type LogLevel string

// log level
const (
	TRACE   LogLevel = "TRACE: "
	WARNING LogLevel = "WARNING: "
	INFO    LogLevel = "INFO: "
	ERROR   LogLevel = "ERROR: "
)

// Log information to file with logged day
func Log(lv LogLevel, content string) {
	logPath, err := filepath.Abs(os.Getenv("LOG_PATH"))
	if err != nil {
		fmt.Println(err.Error())
	}
	if _, err := os.Stat(logPath); err != nil && os.IsNotExist(err) {
		err := os.Mkdir(logPath, os.ModePerm)
		if nil != err {
			fmt.Println(err.Error())
		}
	}
	file := logPath + "/" + time.Now().Format("20060102") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		fmt.Println(err.Error())
	}
	logger := log.New(logFile, string(lv), log.LstdFlags|log.Llongfile)
	logger.Output(2, content)
}
