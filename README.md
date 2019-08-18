# goploy
## 使用说明
1. 导入根目录goploy.sql
2. cp .env.example .env
3. 修改.env里面的配置
4. 运行./goploy or goploy.exe
5. web http://ip:port  (账号:密码 admin:admin!@#)

## 后端开发说明
1. 安装go，推荐1.11以上
2. 项目使用dep管理依赖，安装好dep之后，dep ensure -v
3. 运行go run main.go

## 前端开发说明
1. cd web
2. 修改.env.development
3. npm run dev

## 原理图
![原理图](https://github.com/zhenorzz/goploy/blob/master/goploy.png)

## 预览
![预览](https://github.com/zhenorzz/goploy/blob/master/snapshot.gif)
有不懂的看可以发issue

欢迎大家合作开发
