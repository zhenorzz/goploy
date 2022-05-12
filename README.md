<p align=center>
    <img src="./banner.png" alt="logo" title="logo" />
</p>

<p align="center">
  <a href="#">
      <img src="https://img.shields.io/badge/readme%20style-standard-brightgreen.svg">
  </a>
  <a href="#">
      <img src="https://img.shields.io/badge/give%20me-a%20star-green.svg">
    </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg">
  </a>
</p>

English | [简体中文](./README.zh-CN.md)

Name: go + deploy

A web deployment system tool!

Support all kinds of code release and rollback, which can be done through the web with one click!

Complete installation instructions, no difficulty in getting started!

> Important note: The master branch may be in an unstable or unavailable state during the development process. Please use releases instead of master to obtain stable binary files.

[DEMO](http://demo.goploy.icu) admin:admin!@# (It may not be able to open, depending on the mood)

[Docker](https://hub.docker.com/r/zhenorzz/goploy)

[Dockerfile](./docker/Dockerfile)

[Document](https://docs.goploy.icu)

[Goploy-Agent](https://github.com/zhenorzz/goploy-agent) Monitor server performance

[Jetbrains](https://www.jetbrains.com/?from=zhenorzz/goploy) supports this project with GoLand licenses. We appreciate their support for free and open source software!

## Content

- [Feature](#Feature)
- [Install](#Install)
- [Use](#Use)
- [Preview](#Preview)
- [Diagram](#Diagram)
- [Backend](#Backend)
- [Frontend](#Frontend)
- [Repository](#Repository)
- [Contribute](#Contribute)
- [License](#License)

## Feature

Use Goploy to automate your development workflow, so you can focus on work that matters most. 

Goploy is commonly used for:

- Building projects
- Support git svn ftp sftp
- Deployment across os
- RBAC
- Monitor http tcp ping process script server
- Second cron 
- Xterm
- Sftp
- LDAP

## Install
1. Install mysql
2. Download the latest release

## Use
1. Run ./goploy or goploy.exe or goploy.mac
2. Follow the installation guide
3. web http://ip:port  (Account:Password admin:admin!@#)

## Preview
![Preview](./preview.png)

## Diagram
![Diagram](./goploy.png)

## Backend
1. Install go >= 1.16
2. go mod required
3. edit .env ENV=dev   
4. go run main.go
5. use gin (hot reload)

## Frontend
1. cd web
2. vi .env.development
3. npm run dev

## Repository

- [gin](https://github.com/codegangsta/gin) - GO hot reload。
- [element-ui](https://github.com/ElemeFE/element) - UI。

## Contribute

[Issue](https://github.com/zhenorzz/goploy/issues/new) 

Create a pull request.

## License

[GPLv3](LICENSE) © zhenorzz
