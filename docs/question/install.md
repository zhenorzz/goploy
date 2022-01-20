# 安装问题

# panic: Error 1071: Specified key was too long;max key length is 767 bytes

需要升级数据库版本 MySQL >= 5.7

# Error during Websocket handshake(检测到Websocket与服务器断开连接)

原因是反向代理设置不正确
```nginx
# nginx需要设置升级http协议
proxy_set_header Upgrade         $http_upgrade;
proxy_set_header Connection      "upgrade";
```

# Illegal request

请检查反向代理是否正确，查看network是否有带cookie

# centos6安装rsync源
```shell
cd /etc/yum.repos.d/
mkdir bak
mv * bak/
curl -o /etc/yum.repos.d/CentOS-Base.repo http://file.kangle.odata.cc/repo/Centos-6.repo
curl -o /etc/yum.repos.d/epel.repo http://file.kangle.odata.cc/repo/epel-6.repo
yum clean all
yum makecache
```

# git 自建证书

设置环境变量GIT_SSL_NO_VERIFY=1

