# 安装Goploy

#（QQ群903750786有视频教程）

# 安装准备

安装前必须保证系统已经安装了下列软件

- go >= 1.16
- Git >= 2.10
- MySQL >= 5.7
- MySQL客户端
- Rsync(linux mac自带，windows需要安装[cwRsync](https://www.itefix.net/cwrsync))

# 安装 Goploy

下载releases
- https://github.com/zhenorzz/goploy/releases
- https://gitee.com/goploy/goploy/releases

[Docker](https://hub.docker.com/r/zhenorzz/goploy)

# 启动

```shell
# 新手推荐用root启动，避免不必要的问题
# 运行不了可能需要 chmod a+x, 再不行就发issue或加群
# Windows
goploy.exe
# Linux
./goploy
# Mac
./goploy.mac
```

# 配置

```shell
请输入mysql的用户:
***
请输入mysql的密码:
******
请输入mysql的主机(默认127.0.0.1，不带端口):

请输入mysql的端口(默认3306):

请输入日志目录的绝对路径(默认stdout):

请输入监听端口(默认80，打开网页时的端口):

#输入完成稍等片刻即可安装完成
```
    
# 访问

http://host:port (账号:密码 admin:admin!@#)

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

# 守护进程

推荐使用systemd

```shell
[Unit]
Description=The Goploy
After=network.target

[Service]
Environment="HOME=/root"
WorkingDirectory=/var/www/goploy
ExecStart=/var/www/goploy/goploy

[Install]
WantedBy=multi-user.target
```