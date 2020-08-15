package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/route"
	"github.com/zhenorzz/goploy/task"
	"github.com/zhenorzz/goploy/utils"
	"github.com/zhenorzz/goploy/ws"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	println(`
   ______            __           
  / ____/___  ____  / /___  __  __
 / / __/ __ \/ __ \/ / __ \/ / / /
/ /_/ / /_/ / /_/ / / /_/ / /_/ / 
\____/\____/ .___/_/\____/\__, /  
          /_/            /____/   
`)
	install()
	godotenv.Load(core.GlobalPath + ".env")
	core.CreateValidator()
	model.Init()
	ws.Init()
	route.Init()
	task.Init()
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func install() {
	println("Check if it is installed for the first time")
	_, err := os.Stat(".env")
	if err == nil || os.IsExist(err) {
		println("The configuration file already exists, no need to reinstall (if you need to reinstall, please back up the database goploy first, delete the .env file)")
		return
	}
	inputReader := bufio.NewReader(os.Stdin)
	println("Installation guidelines (Enter to confirm input)")
	println("Please enter the mysql user:")
	mysqlUser, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlUser = utils.ClearNewline(mysqlUser)
	println("Please enter the mysql password:")
	mysqlPassword, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPassword = utils.ClearNewline(mysqlPassword)
	if len(mysqlPassword) != 0 {
		mysqlPassword = ":" + mysqlPassword
	}
	println("Please enter the mysql host(default 127.0.0.1, without port):")
	mysqlHost, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlHost = utils.ClearNewline(mysqlHost)
	if len(mysqlHost) == 0 {
		mysqlHost = "127.0.0.1"
	}
	println("Please enter the mysql port(default 3306):")
	mysqlPort, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPort = utils.ClearNewline(mysqlPort)
	if len(mysqlPort) == 0 {
		mysqlPort = "3306"
	}
	println("Please enter the absolute path of the log directory(default /tmp/):")
	logPath, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	logPath = utils.ClearNewline(logPath)
	if len(logPath) == 0 {
		logPath = "/tmp/"
	}
	println("Please enter the absolute path of the ssh-key directory(default /root/.ssh/id_rsa):")
	sshFile, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	sshFile = utils.ClearNewline(sshFile)
	if len(sshFile) == 0 {
		sshFile = "/root/.ssh/id_rsa"
	}
	println("Please enter the listening port(default 80):")
	port, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	port = utils.ClearNewline(port)
	if len(port) == 0 {
		port = "80"
	}
	println("Start to install the database...")

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s%s@tcp(%s:%s)/?charset=utf8mb4,utf8\n",
		mysqlUser,
		mysqlPassword,
		mysqlHost,
		mysqlPort))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := model.ImportSQL(db); err != nil {
		panic(err)
	}
	println("Database installation is complete")
	envContent := ""
	envContent += "DB_TYPE=mysql\n"
	envContent += fmt.Sprintf(
		"DB_CONN=%s%s@tcp(%s:%s)/goploy?charset=utf8mb4,utf8\n",
		mysqlUser,
		mysqlPassword,
		mysqlHost,
		mysqlPort)
	envContent += fmt.Sprintf("SIGN_KEY=%d\n", time.Now().Unix())
	envContent += fmt.Sprintf("LOG_PATH=%s\n", logPath)
	envContent += fmt.Sprintf("SSHKEY_PATH=%s\n", sshFile)
	envContent += "ENV=production\n"
	envContent += fmt.Sprintf("PORT=%s\n", port)
	println("Start writing configuration file...")
	file, err := os.Create(".env")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(envContent)
	println("Write configuration file completed")
}
