# 安装Goploy

# 安装准备

安装前必须保证系统已经安装了下列软件

- MySQL-5.6或更高版本
- MySQL客户端
- Rsync(linux mac自带，windows需要安装)

# 安装 Goploy
方法1:

```
go get -u github.com/zhenorzz/goploy
```

方法2: 

下载源码
- https://github.com/zhenorzz/goploy
- https://gitee.com/zhenorzz/goploy

# 启动

```shell
# 运行不了可能需要 chmod a+x, 再不行就发issue或加群
# Windows
goploy.exe
# Linux
./goploy
# Mac
./goploy.mac
```

# 配置

配置goploy有两种方法

- 启动前手动复制配置文件
    1. cp .env.example .env
    2. 安装数据库文件 goploy.sql

- 直接启动引导安装

```shell
请输入mysql的用户:
***
请输入mysql的密码:
******
请输入mysql的主机(默认127.0.0.1，不带端口):

请输入mysql的端口(默认3306):

请输入日志目录的绝对路径(默认/tmp/):

请输入sshkey的绝对路径(默认/root/.ssh/id_rsa):

请输入监听端口(默认80，打开网页时的端口):

#输入完成稍等片刻即可安装完成
```
    
# 访问

http://ip:port (账号:密码 admin:admin!@#)

# 反向代理

```nginx
server{
    listen       80;
    server_name  goploy.com;
    access_log   /data/nginx_logs/goploy.com.log main;

    location /{
        proxy_set_header X-Real-IP       $remote_addr;
        proxy_set_header Host            $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_http_version               1.1;
        proxy_set_header Upgrade         $http_upgrade;
        proxy_set_header Connection      "upgrade";
        proxy_pass                       http://{yourip}:{yourport};
    }
}
```