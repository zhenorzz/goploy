CREATE TABLE IF NOT EXISTS `project`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT 'project name',
  `repo_type` varchar(255) NOT NULL DEFAULT '' COMMENT 'repository type (git | svn)',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT 'repository url',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT 'project deploy path',
  `symlink_path` varchar(255) NOT NULL DEFAULT '' COMMENT '(ln -sfn symlink_path/uuid project_path)',
  `symlink_backup_number` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'symlink backup number',
  `environment` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '1.production 2.pre-release 3.test 4.development',
  `branch` varchar(255) NOT NULL DEFAULT 'master' COMMENT 'repository branch',
  `label` varchar(6382) NOT NULL DEFAULT '',
  `review` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '0.disable 1.enable',
  `review_url` varchar(1000) NOT NULL DEFAULT '' COMMENT 'review notification link',
  `script` json NOT NULL COMMENT 'script',
  `transfer_type` varchar(255) NOT NULL DEFAULT '',
  `transfer_option` varchar(255) NOT NULL DEFAULT '',
  `deploy_server_mode` varchar(255) NOT NULL DEFAULT '' COMMENT 'serial | parallel',
  `auto_deploy` tinyint(4) UNSIGNED NOT NULL DEFAULT 1 COMMENT '0.disable 1.webhook',
  `state` tinyint(4) UNSIGNED NOT NULL DEFAULT 1 COMMENT '0.disable 1.enable',
  `deploy_state` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '0.not deploy 1.deploying 2.success 3.fail',
  `publisher_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `publisher_name` varchar(255) NOT NULL DEFAULT '',
  `last_publish_token` char(36) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
  `notify_type` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '1.weixin 2.ding talk 3.feishu 255.custom',
  `notify_target` varchar(255) NOT NULL DEFAULT '' COMMENT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `project_file` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL,
  `filename` varchar(255) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_project_id` (`project_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `project_server`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `server_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_project_server`(`project_id`, `server_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `project_user`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_project_user`(`project_id`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `project_task` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0',
  `commit` char(40) NOT NULL DEFAULT '',
  `branch` varchar(255) NOT NULL DEFAULT '',
  `date` datetime DEFAULT NULL,
  `state` tinyint(4) unsigned NOT NULL DEFAULT '1',
  `is_run` tinyint(4) unsigned NOT NULL DEFAULT '0',
  `creator_id` int(10) unsigned NOT NULL DEFAULT '0',
  `creator` varchar(255) NOT NULL DEFAULT '',
  `editor_id` int(10) unsigned NOT NULL DEFAULT '0',
  `editor` varchar(255) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_project_update` (`project_id`,`update_time`) USING BTREE COMMENT 'project_id,update_time'
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `project_review` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0',
  `commit_id` char(40) NOT NULL DEFAULT '',
  `branch` varchar(255) NOT NULL DEFAULT '',
  `review_url` varchar(1000) NOT NULL DEFAULT '' COMMENT 'review notification link',
  `state` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '0.Pending 1.Approve 2.Deny',
  `creator_id` int(10) unsigned NOT NULL DEFAULT '0',
  `creator` varchar(255) NOT NULL DEFAULT '',
  `editor_id` int(10) unsigned NOT NULL DEFAULT '0',
  `editor` varchar(255) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `project_process` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `project_id` int(10) unsigned NOT NULL DEFAULT '0',
    `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `start` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `stop` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `status` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `restart` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `publish_trace` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `token` char(36) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
  `project_id` int(10) unsigned NOT NULL DEFAULT '0',
  `project_name` varchar(255) NOT NULL DEFAULT '',
  `detail` longtext NOT NULL,
  `state` tinyint(4) unsigned NOT NULL DEFAULT '1',
  `publisher_id` int(10) unsigned NOT NULL DEFAULT '0',
  `publisher_name` varchar(255) NOT NULL DEFAULT '',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1.before pull 2.pulled 3.after pull 4.before deploy 5.deploy 6.after deploy',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `ext` longtext NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_project_id` (`project_id`) USING BTREE COMMENT 'project_id'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `monitor` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `namespace_id` int(10) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `type` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '1.http 2.dial 3.ping 4.process 5.script',
  `target` json NOT NULL,
  `second` int(10) unsigned NOT NULL DEFAULT '1' COMMENT 'How many seconds to run',
  `times` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT 'How many times of failures',
  `silent_cycle` int(10) unsigned NOT NULL DEFAULT '0',
  `description` varchar(255) NOT NULL DEFAULT '',
  `notify_type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '1.wecom 2.ding talk 3.feishu 255.custom',
  `notify_target` varchar(255) NOT NULL DEFAULT '',
  `state` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '0.disable 1.enable',
  `error_content` varchar(1000) NOT NULL DEFAULT '',
  `success_server_id` int(10) unsigned NOT NULL DEFAULT '0',
  `success_script` text NOT NULL,
  `fail_server_id` int(10) unsigned NOT NULL DEFAULT '0',
  `fail_script` text NOT NULL,
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `server`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL DEFAULT '',
  `ip` varchar(255) NOT NULL DEFAULT '',
  `port` smallint(10) UNSIGNED NOT NULL DEFAULT 22,
  `owner` varchar(255) NOT NULL DEFAULT '',
  `path` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `jump_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `jump_port` smallint(5) unsigned NOT NULL DEFAULT '0',
  `jump_owner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `jump_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `jump_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `os` varchar(255) NOT NULL DEFAULT '' COMMENT 'linux|windows',
  `os_info` varchar(255) NOT NULL DEFAULT '' COMMENT 'os|cpu cores|mem',
  `description` varchar(255) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `state` tinyint(4) UNSIGNED NOT NULL DEFAULT 1 COMMENT '0.disable 1.enable',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_namespace_ip` (`namespace_id`,`ip`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `server_agent_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `server_id` int(10) unsigned NOT NULL DEFAULT '0',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `item` varchar(255) NOT NULL DEFAULT '',
  `value` varchar(255) NOT NULL DEFAULT '',
  `report_time` datetime NOT NULL,
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_server_type_item_time` (`server_id`,`type`,`item`,`report_time`) USING BTREE,
  KEY `idx_server_item_time` (`server_id`,`item`,`report_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `server_monitor` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `server_id` int(10) unsigned NOT NULL DEFAULT '0',
  `item` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `formula` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'avg|max|min',
  `operator` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `value` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `group_cycle` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'uint minute',
  `last_cycle` int(10) unsigned NOT NULL DEFAULT '0',
  `silent_cycle` int(10) unsigned NOT NULL DEFAULT '0',
  `start_time` char(5) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '00:00',
  `end_time` char(5) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '23:59',
  `notify_type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '1=企业微信 2=钉钉 3=飞书 255=custom',
  `notify_target` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_server_item` (`server_id`,`item`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `cron` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `server_id` int(10) unsigned NOT NULL DEFAULT '0',
    `expression` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `command` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `single_mode` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1:wait the current run completed',
    `log_level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0:none 1:stdout 2: 1+stderr ',
    `description` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `creator` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `editor` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `state` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '0.disable 1.enable',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `server_process` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `namespace_id` int(10) unsigned NOT NULL,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `items` json DEFAULT NULL COMMENT '{name: string, command: string}[]',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_namespace` (`namespace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `cron_log` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `cron_id` int(10) unsigned NOT NULL DEFAULT '0',
    `server_id` int(10) unsigned NOT NULL DEFAULT '0',
    `exec_code` int(10) NOT NULL DEFAULT '0' COMMENT 'shell exec code',
    `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `report_time` datetime NOT NULL,
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_cron` (`cron_id`),
    KEY `idx_server_cron` (`server_id`,`cron_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `user`  (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(30) NOT NULL DEFAULT '',
  `password` varchar(60) NOT NULL DEFAULT '',
  `password_update_time` datetime DEFAULT NULL,
  `name` varchar(30) NOT NULL DEFAULT '',
  `contact` varchar(255) NOT NULL DEFAULT '',
  `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0.disable 1.enable',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `last_login_time` datetime DEFAULT NULL,
  `super_manager` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'the mark of super admin',
  `api_key` char(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `namespace` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `namespace_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `namespace_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `role_id` int(10) unsigned NOT NULL DEFAULT '0',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_namespace_user` (`namespace_id`,`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `role` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `permission` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `pid` int(10) unsigned NOT NULL DEFAULT '0',
    `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `sort` int(10) NOT NULL DEFAULT '0',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_pid` (`pid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `role_permission` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `role_id` int(10) unsigned NOT NULL,
    `permission_id` int(10) unsigned NOT NULL,
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `template` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `type` tinyint(3) unsigned NOT NULL DEFAULT '0',
    `name` varchar(255) NOT NULL DEFAULT '',
    `content` text NOT NULL,
    `description` varchar(2047) NOT NULL DEFAULT '',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `system_config` (
    `id`    int(10) unsigned NOT NULL AUTO_INCREMENT,
    `key`   varchar(255) NOT NULL DEFAULT '',
    `value` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `udx_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `login_log` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `account` varchar(30) NOT NULL DEFAULT '',
    `remote_addr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `user_agent` varchar(255) NOT NULL DEFAULT '',
    `referer` varchar(255) NOT NULL DEFAULT '',
    `reason` varchar(2555) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `login_time` datetime NOT NULL,
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `operation_log` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` int(10) unsigned NOT NULL COMMENT '',
    `namespace_id` int(10) unsigned NOT NULL DEFAULT '0',
    `router` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'request router',
    `api` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'api',
    `request_time` datetime COMMENT 'request time',
    `request_data` json COMMENT 'request data',
    `response_time` datetime COMMENT 'response time',
    `response_data` json COMMENT 'response data',
    PRIMARY KEY (`id`),
    KEY `idx_user_namespace` (`user_id`, `namespace_id`),
    KEY `idx_router` (`router`),
    KEY `idx_api` (`api`),
    KEY `idx_request_time` (`request_time`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `sftp_log` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `namespace_id` int(10) unsigned NOT NULL DEFAULT '0',
    `user_id` int(10) unsigned NOT NULL DEFAULT '0',
    `server_id` int(10) unsigned NOT NULL DEFAULT '0',
    `remote_addr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `user_agent` varchar(255) NOT NULL DEFAULT '',
    `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'READ|PREVIEW|DOWNLOAD|UPLOAD',
    `path` varchar(255) NOT NULL DEFAULT '',
    `reason` varchar(2555) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `terminal_log` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `namespace_id` int(10) unsigned NOT NULL DEFAULT '0',
    `user_id` int(10) unsigned NOT NULL DEFAULT '0',
    `server_id` int(10) unsigned NOT NULL DEFAULT '0',
    `remote_addr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `user_agent` varchar(255) NOT NULL DEFAULT '',
    `start_time` datetime NOT NULL,
    `end_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT IGNORE INTO `user`(`id`, `account`, `password`, `name`, `contact`, `state`, `super_manager`) VALUES (1, 'admin', '$2a$10$89ZJ2xeJj35GOw11Qiucr.phaEZP4.kBX6aKTs7oWFp1xcGBBgijm', '超管', '', 1, 1);
INSERT IGNORE INTO `namespace`(`id`, `name`) VALUES (1, 'goploy');
INSERT IGNORE INTO `namespace_user`(`id`, `namespace_id`, `user_id`, `role_id`) VALUES (1, 1, 1, 0);
INSERT IGNORE INTO `system_config` (`id`, `key`, `value`) VALUES (1, 'version', '1.16.2');
INSERT IGNORE INTO `role`(`id`, `name`, `description`) VALUES (1, 'manager', '');
INSERT IGNORE INTO `role`(`id`, `name`, `description`) VALUES (2, 'member', '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (1, 0, 'Log', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (2, 1, 'ShowLoginLogPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (3, 1, 'ShowPublishLogPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (4, 1, 'ShowSFTPLogPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (5, 1, 'ShowTerminalLogPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (6, 1, 'ShowTerminalRecord', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (7, 0, 'Member', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (8, 7, 'ShowMemberPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (9, 7, 'AddMember', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (10, 7, 'EditMember', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (11, 7, 'DeleteMember', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (12, 0, 'Namespace', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (13, 12, 'ShowNamespacePage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (14, 12, 'AddNamespace', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (15, 12, 'EditNamespace', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (16, 12, 'AddNamespaceUser', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (17, 12, 'DeleteNamespaceUser', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (18, 12, 'ShowRolePage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (19, 12, 'AddRole', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (20, 12, 'EditRole', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (21, 12, 'DeleteRole', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (22, 12, 'EditPermission', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (23, 0, 'Server', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (24, 23, 'ShowServerPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (25, 23, 'AddServer', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (26, 23, 'EditServer', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (27, 23, 'SwitchServerState', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (28, 23, 'InstallAgent', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (29, 23, 'ImportCSV', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (30, 23, 'ShowServerMonitorPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (31, 23, 'AddServerWarningRule', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (32, 23, 'EditServerWarningRule', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (33, 23, 'DeleteServerWarningRule', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (34, 23, 'ShowTerminalPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (35, 23, 'ShowSftpFilePage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (36, 23, 'SFTPUploadFile', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (37, 23, 'SFTPPreviewFile', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (38, 23, 'SFTPDownloadFile', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (39, 23, 'ShowCronPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (40, 23, 'AddCron', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (41, 23, 'EditCron', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (42, 23, 'DeleteCron', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (43, 0, 'Project', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (44, 43, 'ShowProjectPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (45, 43, 'GetAllProjectList', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (46, 43, 'GetBindProjectList', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (47, 43, 'AddProject', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (48, 43, 'EditProject', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (49, 43, 'DeleteProject', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (50, 43, 'SwitchProjectWebhook', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (51, 0, 'Monitor', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (52, 51, 'ShowMonitorPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (53, 51, 'AddMonitor', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (54, 51, 'EditMonitor', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (55, 51, 'DeleteMonitor', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (56, 0, 'Deploy', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (57, 56, 'ShowDeployPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (58, 56, 'GetAllDeployList', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (59, 56, 'GetBindDeployList', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (60, 56, 'DeployDetail', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (61, 56, 'DeployProject', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (62, 56, 'DeployResetState', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (63, 56, 'GreyDeploy', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (64, 56, 'DeployRollback', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (65, 56, 'DeployReview', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (66, 56, 'DeployTask', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (67, 56, 'FileCompare', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (68, 56, 'FileSync', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (69, 56, 'ProcessManager', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (70, 23, 'ShowServerProcessPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (71, 23, 'AddServerProcess', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (72, 23, 'EditServerProcess', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (73, 23, 'DeleteServerProcess', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (74, 1, 'ShowOperationLogPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (75, 23, 'SFTPTransferFile', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (76, 23, 'SFTPDeleteFile', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (77, 23, 'ShowServerScriptPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (78, 23, 'SFTPEditFile', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (79, 23, 'ShowServerNginxPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (80, 23, 'ManageServerNginx', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (81, 23, 'AddNginxConfig', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (82, 23, 'EditNginxConfig', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (83, 23, 'DeleteNginxConfig', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (84, 23, 'UnbindServerProject', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (85, 43, 'ManageRepository',  0, '');
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 14);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 15);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 16);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 17);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 18);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 19);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 20);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 21);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 22);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 24);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 25);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 26);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 27);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 28);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 29);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 30);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 31);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 32);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 33);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 34);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 35);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 36);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 37);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 38);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 39);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 40);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 41);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 42);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 44);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 45);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 46);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 47);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 48);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 49);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 50);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 52);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 53);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 54);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 55);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 57);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 58);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 59);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 60);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 61);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 62);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 63);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 64);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 65);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 66);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 67);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 68);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 69);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 57);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 59);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 60);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 61);
INSERT IGNORE INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 67);