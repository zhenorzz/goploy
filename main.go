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
	install()
	godotenv.Load(core.GlobalPath + ".env")
	println("应用启动")
	println("http://localhost:" + os.Getenv("PORT"))
	core.CreateValidator()
	model.Init()
	ws.Init()
	route.Init()
	task.Init()
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func install() {
	println("检测是否第一次安装")
	_, err := os.Stat(".env")
	if err == nil || os.IsExist(err) {
		println("配置文件已存在，无需重新安装(如果需要重新安装，请先备份好数据库goploy，删除.env文件)")
		return
	}
	inputReader := bufio.NewReader(os.Stdin)
	println("安装指引(回车确认输入)")
	println("请输入mysql的用户:")
	mysqlUser, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlUser = utils.ClearNewline(mysqlUser)
	println("请输入mysql的密码:")
	mysqlPassword, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPassword = utils.ClearNewline(mysqlPassword)
	if len(mysqlPassword) != 0 {
		mysqlPassword = ":" + mysqlPassword
	}
	println("请输入mysql的主机(默认127.0.0.1，不带端口):")
	mysqlHost, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlHost = utils.ClearNewline(mysqlHost)
	if len(mysqlHost) == 0 {
		mysqlHost = "127.0.0.1"
	}
	println("请输入mysql的端口(默认3306):")
	mysqlPort, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	mysqlPort = utils.ClearNewline(mysqlPort)
	if len(mysqlPort) == 0 {
		mysqlPort = "3306"
	}
	println("请输入日志目录的绝对路径(默认/tmp/):")
	logPath, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	logPath = utils.ClearNewline(logPath)
	if len(logPath) == 0 {
		logPath = "/tmp/"
	}
	println("请输入sshkey的绝对路径(默认/root/.ssh/id_rsa):")
	sshFile, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	sshFile = utils.ClearNewline(sshFile)
	if len(sshFile) == 0 {
		sshFile = "/root/.ssh/id_rsa"
	}
	println("请输入监听端口(默认80，打开网页时的端口):")
	port, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	port = utils.ClearNewline(port)
	if len(port) == 0 {
		port = "80"
	}
	println("开始安装数据库...")

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
	println("安装数据库完成")
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
	println("开始写入配置文件")
	file, err := os.Create(".env")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(envContent)
	println("写入配置文件完成")
}
