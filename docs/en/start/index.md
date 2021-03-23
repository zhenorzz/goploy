# Install Goploy

# Installation preparation

Dependencies

- Git >= 2.10
- MySQL>=5.6
- MySQL Client
- Rsync(windows need to install [cwRsync](https://www.itefix.net/cwrsync))

# Install Goploy
Method 1:

```
go get -u github.com/zhenorzz/goploy
```

Method 2: 

Download releases
- https://github.com/zhenorzz/goploy/releases
- https://gitee.com/goploy/goploy/releases

# Start up

```shell
# Windows
goploy.exe
# Linux
./goploy
# Mac
./goploy.mac
```

# Configuration

There are two ways to configure goploy

- Manually copy the configuration file before starting
    1. cp .env.example .env
    2. Install database files goploy.sql

- Boot installation

```shell
Please enter the mysql user:
***
Please enter the mysql password:
******
Please enter the mysql host(default 127.0.0.1, without port):

Please enter the mysql port(default 3306):

Please enter the absolute path of the log directory(default /tmp/):

Please enter the absolute path of the ssh-key directory(default /root/.ssh/id_rsa):

Please enter the listening port(default 80):

#After the input is complete, wait a moment to complete the installation
```
    
# Access

http://host:port (Account:Password admin:admin!@#)

# Reverse proxy

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

# Daemon

systemd

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