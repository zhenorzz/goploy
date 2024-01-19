package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type LogConfig struct {
	Path string `toml:"path"`
}

func (l *LogConfig) OnChange() error {
	setLogger()
	return nil
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
