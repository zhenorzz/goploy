// Copyright 2022 The Goploy Authors. All rights reserved.

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/zhenorzz/goploy/cmd/server/api/agent"
	"github.com/zhenorzz/goploy/cmd/server/api/cron"
	"github.com/zhenorzz/goploy/cmd/server/api/deploy"
	logApi "github.com/zhenorzz/goploy/cmd/server/api/log"
	"github.com/zhenorzz/goploy/cmd/server/api/monitor"
	"github.com/zhenorzz/goploy/cmd/server/api/namespace"
	"github.com/zhenorzz/goploy/cmd/server/api/project"
	"github.com/zhenorzz/goploy/cmd/server/api/repository"
	"github.com/zhenorzz/goploy/cmd/server/api/role"
	server2 "github.com/zhenorzz/goploy/cmd/server/api/server"
	"github.com/zhenorzz/goploy/cmd/server/api/template"
	"github.com/zhenorzz/goploy/cmd/server/api/user"
	"github.com/zhenorzz/goploy/cmd/server/task"
	"github.com/zhenorzz/goploy/cmd/server/ws"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/database"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	help bool
	v    bool
	s    string
)

const appVersion = "1.16.2"

func init() {
	flag.StringVar(&config.AssetDir, "asset-dir", "", "default: ./")
	flag.StringVar(&s, "s", "", "stop")
	flag.BoolVar(&help, "help", false, "list available subcommands and some concept guides")
	flag.BoolVar(&v, "version", false, "show goploy version")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	go checkUpdate()
}

// @title Goploy
// @version 1.16.2
// @description A web deployment system tool!
// @contact.name zhenorzz
// @contact.url https://github.com/zhenorzz/goploy
// @contact.email zhenorzz@gmail.com
// @license.name GPLv3
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @host 127.0.0.1:3001
// @externalDocs.description Get started
// @externalDocs.url https://docs.goploy.icu/
// @securityDefinitions.apikey ApiKeyHeader
// @in header
// @name X-API-KEY
// @securityDefinitions.apikey ApiKeyQueryParam
// @in query
// @name X-API-KEY
// @securityDefinitions.apikey NamespaceHeader
// @in query
// @name G-N-ID
// @securityDefinitions.apikey NamespaceQueryParam
// @in query
// @name G-N-ID
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
	config.Init()
	model.Init()
	if err := model.Update(appVersion); err != nil {
		println(err.Error())
	}
	task.Init()
	srv := server.Server{
		Server: http.Server{
			Addr: ":" + config.Toml.Web.Port,
		},
		Router: server.NewRouter(),
	}

	srv.Router.Register(user.User{})
	srv.Router.Register(namespace.Namespace{})
	srv.Router.Register(role.Role{})
	srv.Router.Register(project.Project{})
	srv.Router.Register(repository.Repository{})
	srv.Router.Register(monitor.Monitor{})
	srv.Router.Register(deploy.Deploy{})
	srv.Router.Register(server2.Server{})
	srv.Router.Register(logApi.Log{})
	srv.Router.Register(cron.Cron{})
	srv.Router.Register(agent.Agent{})
	srv.Router.Register(template.Template{})
	srv.Router.Register(ws.GetHub())
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
	srv.Spin()
}

func install() {
	_, err := os.Stat(config.GetConfigFile())
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

	if !runningInDocker() {
		cfg = readStdin(cfg)
	} else {
		cfg.DB.Host = os.Getenv("DB_HOST")
		cfg.DB.User = os.Getenv("DB_USER")
		cfg.DB.Password = os.Getenv("DB_USER_PASSWORD")
		cfg.DB.Database = os.Getenv("DB_NAME")
		cfg.DB.Port = os.Getenv("DB_PORT")
		if cfg.DB.Port == "" {
			cfg.DB.Port = "3306"
		}
	}

	println("Start to install the database...")

	runner, err := model.Open(cfg.DB)
	if err != nil {
		panic(err)
	}
	defer runner.Close()
	if err := runner.CreateDB(cfg.DB.Database); err != nil {
		panic(err)
	}
	if err := runner.UseDB(cfg.DB.Database); err != nil {
		panic(err)
	}
	if err := runner.ImportSQL(database.GoploySQL); err != nil {
		panic(err)
	}
	println("Database installation is complete")
	println("Start writing configuration file...")
	err = config.Write(cfg)
	if err != nil {
		panic("Write config file error, " + err.Error())
	}
	println("Write configuration file completed")
}

func handleClientSignal() {
	switch s {
	case "stop":
		pidFile := config.GetPidFile()
		pidStr, err := os.ReadFile(pidFile)
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

func readStdin(cfg config.Config) config.Config {
	println("Installation guide ↓")
	inputReader := bufio.NewReader(os.Stdin)
	println("Installation guidelines (Enter to confirm input)")

	println("Please enter the mysql user:")
	mysqlUser, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	cfg.DB.User = pkg.ClearNewline(mysqlUser)

	println("Please enter the mysql password:")
	mysqlPassword, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPassword = pkg.ClearNewline(mysqlPassword)
	if len(mysqlPassword) != 0 {
		cfg.DB.Password = mysqlPassword
	}

	println("Please enter the mysql host(default 127.0.0.1, without port):")
	mysqlHost, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlHost = pkg.ClearNewline(mysqlHost)
	if len(mysqlHost) != 0 {
		cfg.DB.Host = mysqlHost
	}

	println("Please enter the mysql port(default 3306):")
	mysqlPort, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPort = pkg.ClearNewline(mysqlPort)
	if len(mysqlPort) != 0 {
		cfg.DB.Port = mysqlPort
	}

	println("Please enter the database name(default goploy):")
	mysqlDB, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlDB = pkg.ClearNewline(mysqlDB)
	if len(mysqlDB) != 0 {
		cfg.DB.Database = mysqlDB
	}

	println("Please enter the absolute path of the log directory(default stdout):")
	logPath, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	logPath = pkg.ClearNewline(logPath)
	if len(logPath) != 0 {
		cfg.Log.Path = logPath
	}

	println("Please enter the listening port(default 80):")
	port, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	port = pkg.ClearNewline(port)
	if len(port) != 0 {
		cfg.Web.Port = port
	}
	return cfg
}

func checkUpdate() {
	resp, err := http.Get("https://api.github.com/repos/zhenorzz/goploy/releases/latest")
	if err != nil {
		println("Check failed")
		println(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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

func runningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}

	if _, err := os.Stat("/proc/self/cgroup"); err == nil {
		data, err := os.ReadFile("/proc/self/cgroup")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.Contains(line, "docker") {
					return true
				}
			}
		}
	}

	return false
}
