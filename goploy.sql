# Host: localhost  (Version: 5.5.53)
# Date: 2019-07-15 16:11:00
# Generator: MySQL-Front 5.3  (Build 4.234)

/*!40101 SET NAMES utf8 */;

#
# Structure for table "git_trace"
#

DROP TABLE IF EXISTS `git_trace`;
CREATE TABLE `git_trace` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL,
  `project_name` varchar(255) NOT NULL,
  `detail` text NOT NULL,
  `state` tinyint(4) unsigned NOT NULL,
  `publisher_id` int(10) unsigned NOT NULL,
  `publisher_name` varchar(255) NOT NULL,
  `create_time` int(10) unsigned NOT NULL,
  `update_time` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `project_id` (`project_id`) USING BTREE COMMENT 'project_id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#
# Structure for table "log"
#

DROP TABLE IF EXISTS `log`;
CREATE TABLE `log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '日志类型',
  `ip` int(10) unsigned NOT NULL DEFAULT '0',
  `desc` varchar(30) NOT NULL DEFAULT '' COMMENT '备注',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#
# Structure for table "permission"
#

DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '权限id',
  `title` varchar(64) NOT NULL DEFAULT '' COMMENT '权限标题',
  `uri` varchar(1000) NOT NULL DEFAULT '' COMMENT '权限uri',
  `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '该记录是否有效1：有效、0：无效',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父级权限ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `permission` VALUES (1, '主页', '/dashboard', 1, 0);
INSERT INTO `permission` VALUES (2, '主页列表', '/dashboard/list', 1, 1);
INSERT INTO `permission` VALUES (3, '项目', '/project', 1, 0);
INSERT INTO `permission` VALUES (4, '项目部署', '/project/deploy', 1, 3);
INSERT INTO `permission` VALUES (5, '项目设置', '/project/setting', 1, 3);
INSERT INTO `permission` VALUES (6, '项目详情', '/project/detail', 1, 3);
INSERT INTO `permission` VALUES (7, '服务器管理', '/project/server', 1, 3);
INSERT INTO `permission` VALUES (8, '成员', '/member', 1, 0);
INSERT INTO `permission` VALUES (9, '成员列表', '/member/list', 1, 8);
INSERT INTO `permission` VALUES (10, '部署', '/deploy', 1, 0);
INSERT INTO `permission` VALUES (11, '构建发布', '/deploy/publish', 1, 10);
#
# Structure for table "project"
#

DROP TABLE IF EXISTS `project`;
CREATE TABLE `project` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '项目名称',
  `url` varchar(255) NOT NULL COMMENT '项目仓库地址',
  `path` varchar(255) NOT NULL COMMENT '项目部署路径',
  `script` varchar(255) NOT NULL COMMENT '脚本路径',
  `rsync_option` varchar(255) NOT NULL COMMENT 'rsync 参数',
  `publisher_id` int(10) unsigned NOT NULL,
  `publisher_name` varchar(255) NOT NULL,
  `create_time` int(11) unsigned NOT NULL,
  `update_time` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#
# Structure for table "project_server"
#

DROP TABLE IF EXISTS `project_server`;
CREATE TABLE `project_server` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL,
  `server_id` int(10) unsigned NOT NULL,
  `create_time` int(11) unsigned NOT NULL,
  `update_time` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `project_id` (`project_id`,`server_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#
# Structure for table "project_user"
#

DROP TABLE IF EXISTS `project_user`;
CREATE TABLE `project_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `create_time` int(11) unsigned NOT NULL,
  `update_time` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `project_id` (`project_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#
# Structure for table "remote_trace"
#

DROP TABLE IF EXISTS `remote_trace`;
CREATE TABLE `remote_trace` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `git_trace_id` int(10) unsigned NOT NULL,
  `project_id` int(10) unsigned NOT NULL,
  `project_name` varchar(255) NOT NULL,
  `server_id` int(11) unsigned NOT NULL,
  `server_name` varchar(255) NOT NULL,
  `detail` text NOT NULL,
  `state` tinyint(4) unsigned NOT NULL,
  `publisher_id` int(10) unsigned NOT NULL,
  `publisher_name` varchar(255) NOT NULL,
  `type` tinyint(3) unsigned NOT NULL COMMENT '1同步文件2运行脚本',
  `create_time` int(10) unsigned NOT NULL,
  `update_time` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `project_id` (`project_id`) USING BTREE COMMENT 'project_id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#
# Structure for table "role"
#

DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '角色名',
  `permission_list` varchar(255) NOT NULL COMMENT '权限列表',
  `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '该记录是否有效1：有效、0：无效',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `desc` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `role` VALUES (1, '超级管理员', '1,2,3,4,5,6,7,8,9,10,11', 1, 1527927486, 1527927486, '超级管理员是最高权限者');
INSERT INTO `role` VALUES (2, '管理员', '10,11', 1, 1527928486, 1552642448, '包含除呼叫功能以外所有客服和管理权限');
INSERT INTO `role` VALUES (3, '普通客服', '10,11', 1, 1527928486, 1552979330, '负责一线在线咨询的接待, 留言的处理(若购买工单中心, 则还有工单功能, 未购买则无)');
INSERT INTO `role` VALUES (4, '高级管理员', '10,11', 1, 1546413701, 1551148083, '大佬就是不一样');
#
# Structure for table "server"
#

DROP TABLE IF EXISTS `server`;
CREATE TABLE `server` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `ip` varchar(255) NOT NULL,
  `owner` varchar(255) NOT NULL,
  `create_time` int(10) unsigned NOT NULL,
  `update_time` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#
# Structure for table "user"
#

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(30) NOT NULL,
  `password` varchar(60) NOT NULL,
  `name` varchar(30) NOT NULL,
  `mobile` varchar(15) NOT NULL,
  `role_id` int(10) unsigned NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0:=被禁用  1=正常',
  `create_time` int(11) unsigned NOT NULL,
  `update_time` int(11) unsigned NOT NULL,
  `last_login_time` int(11) unsigned NOT NULL COMMENT '最后登陆时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `user` VALUES (1, 'admin', '$2a$10$89ZJ2xeJj35GOw11Qiucr.phaEZP4.kBX6aKTs7oWFp1xcGBBgijm', '超管', '18034562122', 1, 1, 1540869597, 1561461294, 1562320680);