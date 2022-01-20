# 使用问题

# 部署后脚本无法重启应用服务

- 检查脚本是否正确kill应用服务
- 使用nohup重启应用

# windows双击exe后闪退

- 打开cmd或者powershell, 输入goploy.exe

# $HOME is not defined

- 检查是否存在环境变量$HOME
- 如果使用systemd或supervisord，需要把HOME目录写入配置文件

# 运行项目脚本提示 xx is not defined

- 如果命令是在启动goploy之后安装的，需要重启goploy
- 如果还是不能运行，可尝试使用全局路径或source /etc/profile

# 自定义文件无法上传

- 检查rsync有无排除该目录