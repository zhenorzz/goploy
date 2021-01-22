# 安装问题

# Error during Websocket handshake(检测到Websocket与服务器断开连接)

原因是反向代理设置不正确
```nginx
# nginx需要设置升级http协议
proxy_set_header Upgrade         $http_upgrade;
proxy_set_header Connection      "upgrade";
```

# Illegal request

请检查反向代理是否正确，查看network是否有带cookie