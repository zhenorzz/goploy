// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/route"
	"github.com/zhenorzz/goploy/task"
	"github.com/zhenorzz/goploy/utils"
	"github.com/zhenorzz/goploy/ws"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	help bool
	v    bool
	s    string
)

const appVersion = "1.8.1"

func init() {
	flag.StringVar(&core.AssetDir, "asset-dir", "", "default: ./")
	flag.StringVar(&s, "s", "", "stop")
	flag.BoolVar(&help, "help", false, "list available subcommands and some concept guides")
	flag.BoolVar(&v, "version", false, "show goploy version")
	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if v {
		println(appVersion)
		return
	}
	handleClientSignal()
	println(`
   ______            __           
  / ____/___  ____  / /___  __  __
 / / __/ __ \/ __ \/ / __ \/ / / /
/ /_/ / /_/ / /_/ / / /_/ / /_/ / 
\____/\____/ .___/_/\____/\__, /  
          /_/            /____/   ` + appVersion + "\n")
	install()
	config.Create(core.GetConfigFile())
	model.Init()
	if err := model.Update(appVersion); err != nil {
		println(err.Error())
	}
	pid := strconv.Itoa(os.Getpid())
	_ = ioutil.WriteFile(path.Join(core.GetAssetDir(), "goploy.pid"), []byte(pid), 0755)
	println("Start at " + time.Now().String())
	println("goploy -h for more help")
	println("Current pid:   " + pid)
	println("Config Loaded: " + core.GetConfigFile())
	println("Env:           " + config.Toml.Env)
	println("Log:           " + config.Toml.Log.Path)
	println("Listen:        " + config.Toml.Web.Port)
	println("Running...")
	core.CreateValidator()
	ws.Init()
	route.Init()
	task.Init()
	go checkUpdate()
	// server
	srv := http.Server{
		Addr: ":" + config.Toml.Web.Port,
	}
	go func() {
		c := make(chan os.Signal, 1)

		// SIGINT Ctrl+C
		// SIGTERM A generic signal used to cause program termination
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		println("Received the signal: " + (<-c).String())

		println("Server is trying to shutdown, wait for a minute")
		ctx, cancel := context.WithTimeout(context.Background(), config.Toml.APP.ShutdownTimeout*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			println("Server shutdown timeout, err: %v\n", err)
		}
		println("Server shutdown gracefully")

		println("Task is trying to shutdown, wait for a minute")
		if err := task.Shutdown(ctx); err != nil {
			println("Task shutdown timeout, err: %v\n", err)
		}
		println("Task shutdown gracefully")
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err.Error())
	}
	_ = os.Remove(path.Join(core.GetAssetDir(), "goploy.pid"))
	println("shutdown success")
	return
}

func install() {
	_, err := os.Stat(core.GetConfigFile())
	if err == nil || os.IsExist(err) {
		println("The configuration file already exists, no need to reinstall (if you need to reinstall, please back up the database `goploy` first, delete the .env file, then restart.)")
		return
	}
	cfg := config.Config{
		Env:    "production",
		APP:    config.APPConfig{DeployLimit: int32(runtime.NumCPU()), ShutdownTimeout: 10},
		Cookie: config.CookieConfig{Name: "goploy_token", Expire: 86400},
		JWT:    config.JWTConfig{Key: time.Now().String()},
		DB:     config.DBConfig{Type: "mysql", Host: "127.0.0.1", Port: "3306", Database: "goploy"},
		Log:    config.LogConfig{Path: "stdout"},
		Web:    config.WebConfig{Port: "80"},
	}
	println("Installation guide ↓")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("rsync", "--version")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		println(err.Error() + ", detail: " + stderr.String())
		panic("Please check if rsync is installed correctly, see https://rsync.samba.org/download.html")
	}
	git := utils.GIT{}
	if err := git.Run("--version"); err != nil {
		println(err.Error() + ", detail: " + git.Err.String())
		panic("Please check if git is installed correctly, see https://git-scm.com/downloads")
	}
	inputReader := bufio.NewReader(os.Stdin)
	println("Installation guidelines (Enter to confirm input)")
	println("Please enter the mysql user:")
	mysqlUser, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	cfg.DB.User = utils.ClearNewline(mysqlUser)
	println("Please enter the mysql password:")
	mysqlPassword, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPassword = utils.ClearNewline(mysqlPassword)
	if len(mysqlPassword) != 0 {
		cfg.DB.Password = mysqlPassword
	}
	println("Please enter the mysql host(default 127.0.0.1, without port):")
	mysqlHost, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlHost = utils.ClearNewline(mysqlHost)
	if len(mysqlHost) != 0 {
		cfg.DB.Host = mysqlHost
	}
	println("Please enter the mysql port(default 3306):")
	mysqlPort, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPort = utils.ClearNewline(mysqlPort)
	if len(mysqlPort) != 0 {
		cfg.DB.Port = mysqlPort
	}
	println("Please enter the absolute path of the log directory(default stdout):")
	logPath, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	logPath = utils.ClearNewline(logPath)
	if len(logPath) != 0 {
		cfg.Log.Path = logPath
	}
	println("Please enter the listening port(default 80):")
	port, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	port = utils.ClearNewline(port)
	if len(port) != 0 {
		cfg.Web.Port = port
	}
	println("Start to install the database...")

	db, err := sql.Open(cfg.DB.Type, fmt.Sprintf(
		"%s:%s@(%s:%s)/?charset=utf8mb4,utf8\n",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := model.CreateDB(db, cfg.DB.Database); err != nil {
		panic(err)
	}
	if err := model.UserDB(db, cfg.DB.Database); err != nil {
		panic(err)
	}
	if err := model.ImportSQL(db, "sql/goploy.sql"); err != nil {
		panic(err)
	}
	println("Database installation is complete")
	println("Start writing configuration file...")
	err = config.Write(core.GetConfigFile(), cfg)
	if err != nil {
		panic("Write config file error, " + err.Error())
	}
	println("Write configuration file completed")
}

func handleClientSignal() {
	switch s {
	case "stop":
		pidFile := path.Join(core.GetAssetDir(), "goploy.pid")
		pidStr, err := ioutil.ReadFile(pidFile)
		if err != nil {
			log.Fatal("handle stop, ", err.Error(), ", may be the server not start")
		}
		pid, _ := strconv.Atoi(string(pidStr))
		process, err := os.FindProcess(pid)
		if err != nil {
			log.Fatal("handle stop, ", err.Error(), ", may be the server not start")
		}
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Fatal("handle stop, ", err.Error())
		}
		println("App is trying to shutdown, wait for a minute")
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			if _, err := os.Stat(pidFile); errors.Is(err, os.ErrNotExist) {
				println("Success")
				break
			} else if err != nil {
				log.Fatal("handle stop, ", err.Error())
			}
		}

		os.Exit(1)
	}
}

func checkUpdate() {
	resp, err := http.Get("https://api.github.com/repos/zhenorzz/goploy/releases/latest")
	if err != nil {
		println("Check failed")
		println(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Check failed")
		println(err.Error())
		return
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		println("Check failed")
		println(err.Error())
		return
	}
	tagName, ok := result["tag_name"].(string)
	if !ok {
		println("Check failed")
		return
	}
	tagVer, err := version.NewVersion(tagName)
	if err != nil {
		println("Check version error")
		println(err.Error())
		return
	}
	currentVer, _ := version.NewVersion(appVersion)
	if tagVer.GreaterThan(currentVer) {
		println("New release available")
		println(result["html_url"].(string))
	}
}
