<p align=center>
    <img src="./banner.png" alt="logo" title="logo" />
</p>

<p align="center">
  <a href="#">
      <img src="https://img.shields.io/badge/readme%20style-standard-brightgreen.svg" alt="readme">
  </a>
  <a href="#">
      <img src="https://img.shields.io/badge/give%20me-a%20star-green.svg" alt="give me a star">
    </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-GPLv3-yellow.svg" alt="License-GPLv3">
  </a>
</p>

[English](./README.md)| 简体中文

名字来源：go + deploy

一个web发布代码工具!

支持各种代码发布、回滚，可以通过web来一键完成！

完整的安装指引，上手无难度！

> 重要提示: master 分支在开发过程中可能处于不稳定或者不可用状态，请通过releases而非master去获取稳定的二进制文件。

[DEMO](http://demo.goploy.icu) admin:admin!@# (不一定能打开，看心情)

[Docker](https://hub.docker.com/r/zhenorzz/goploy)

[Dockerfile](./docker/Dockerfile)

[帮助文档](https://docs.goploy.icu)

[OpenAPI](https://api-doc.goploy.icu)

[Goploy-Agent](https://github.com/zhenorzz/goploy-agent) 监控服务器性能

[Goploy-Vscode](https://github.com/goploy-devops/goploy-vscode) 在Vscode插件市场搜索Goploy

[Goploy-Jetbrains](https://github.com/goploy-devops/goploy-jetbrains) 在Jetbrains插件市场搜索Goploy

## 内容列表

- [功能](#功能)
- [安装](#安装)
- [使用](#使用)
- [预览](#预览)
- [原理图](#原理图)
- [后端开发说明](#后端开发说明)
- [前端开发说明](#前端开发说明)
- [如何贡献](#如何贡献)
- [使用许可](#使用许可)

## 功能

Use Goploy to automate your development workflow, so you can focus on work that matters most.

Goploy 用于:

- 构建项目
- 发布项目到unix、windows
- 支持git、svn、ftp、sftp
- 基于角色的访问控制RBAC
- 监控http、tcp、process
- 服务器性能
- 秒级定时任务
- Xterm
- Sftp
- LDAP

## 视频

[环境配置](https://www.zhihu.com/zvideo/1512006837884469248)
[应用配置](https://www.zhihu.com/zvideo/1513073009618952192)

## 安装

下载 [release](https://github.com/zhenorzz/goploy/releases)

或者从tag构建
```
git clone https://github.com/zhenorzz/goploy.git
cd ./goploy
git checkout -b tagname
cd ./web
npm install && npm run build
cd ../
go build -o goploy cmd/server/main.go 
```

## 使用
1. 运行./goploy or goploy.exe or goploy.mac (运行不了可能需要 chmod a+x, 再不行就发issue或加群)
2. 根据安装指引填入相关的配置
3. web http://ip:port  (账号:密码 admin:admin!@#)

## 预览
![预览](./preview.gif)

## 原理图
![原理图](./goploy.png)

## 后端开发说明
1. 安装go，必须是1.16以上
2. 项目使用 go mod 管理
3. 修改 goploy.toml `cp goploy.example.toml goploy.toml`
4. 需要编译一次前端 [前端开发说明](#前端开发说明)
5. 运行 `cd cmd/server && go run main.go --asset-dir=../../`
6. 或者使用 gin（可以热更新代码，改变就生效）

## 前端开发说明
1. `cd web`
2. `npm install` 或者 `yarn install`
3. 修改 .env.development
4. `npm run dev` 或者 `yarn dev`
5. 编译 `npm run build` 或者 `yarn build`

## 如何贡献

非常欢迎你的加入!

[提一个Issue](https://github.com/zhenorzz/goploy/issues/new) 

提交一个 Pull Request.

## 使用许可

[GPLv3](LICENSE) © zhenorzz
