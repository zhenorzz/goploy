# 构建流程

## 前提

场景：部署php代码到服务器。

A服务器-安装goploy的服务器（一般都是内网服务器）

B服务器-生产环境服务器（一般都是外网服务器）


## 准备工作

* A服务器使用[ssh-keygen](https://www.jianshu.com/p/dd053c18e5ee)生成public key添加到B服务器上（目的是使A服务器能连上B服务器）
* A服务器安装git(记得测试是否能clone和pull代码)
* 检查A、B服务器是否有[rsync](http://www.ruanyifeng.com/blog/2020/08/rsync.html)

## 配置

1. 服务器管理配置服务器
    - ssh key所有者填写ssh-keygen生成时的所有者
2. 项目管理配置项目
    - 部署路径为B服务器的路径，部署时A服务器rsync到B服务器该目录;
    - 开启软链部署，部署时A服务器rsync到B服务器的软链目录，然后ln -sf到部署目录;
    - 拉取后运行脚本是指，git pull后运行的脚本，通常用于打包(npm run build、mvn package之类);
    - 部署后运行脚本是指，rsync到B服务器后运行的脚本，通常用于重启应用(app restart);
    - 开启构建通知后，项目构建会推送到相应的app，如何配置参考相应app的webhook或机器人;
    
## 构建

1. 检测是否有项目目录，无则git clone;
2. git pull;
3. git pull完成后运行脚本;
4. rsync到B服务器
5. rsync完成后运行脚本;
6. 触发构建通知    