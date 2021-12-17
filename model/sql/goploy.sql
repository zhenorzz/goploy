CREATE DATABASE IF NOT EXISTS `goploy`;

use `goploy`;

CREATE TABLE IF NOT EXISTS `log`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `type` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT 'log type',
  `ip` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `desc` varchar(30) NOT NULL DEFAULT '' COMMENT 'description',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '',
  `create_time` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_create_time`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `project`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT 'project name',
  `repo_type` varchar(255) NOT NULL DEFAULT '' COMMENT 'repository type (git | svn)',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT 'repository url',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT 'project deploy path',
  `symlink_path` varchar(255) NOT NULL DEFAULT '' COMMENT '(ln -sfn symlink_path/uuid project_path)',
  `symlink_backup_number` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '软链备份数量',
  `environment` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '1.production 2.pre-release 3.test 4.development',
  `branch` varchar(255) NOT NULL DEFAULT 'master' COMMENT 'repository branch',
  `review` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '0.disable 1.enable',
  `review_url` varchar(1000) NOT NULL DEFAULT '' COMMENT 'review notification link',
  `after_pull_script_mode` varchar(20) NOT NULL DEFAULT '' COMMENT 'sh|php|py|...',
  `after_pull_script` text NOT NULL COMMENT '',
  `after_deploy_script_mode` varchar(20) NOT NULL DEFAULT '' COMMENT 'sh|php|py|...',
  `after_deploy_script` text NOT NULL COMMENT '',
  `rsync_option` varchar(2000) NOT NULL DEFAULT '' COMMENT 'rsync options',
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

CREATE TABLE IF NOT EXISTS `publish_trace` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `token` char(36) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
  `project_id` int(10) unsigned NOT NULL DEFAULT '0',
  `project_group_id` int(10) unsigned NOT NULL DEFAULT '0',
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
  `url` varchar(255) NOT NULL,
  `second` int(10) unsigned NOT NULL DEFAULT '1' COMMENT 'How many seconds to run',
  `times` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT 'How many times of failures',
  `description` varchar(255) NOT NULL DEFAULT '',
  `notify_type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '1.weixin 2.ding talk 3.feishu 255.custom',
  `notify_target` varchar(255) NOT NULL DEFAULT '',
  `notify_times` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT 'Notify times',
  `state` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '0.disable 1.enable',
  `error_content` varchar(1000) NOT NULL DEFAULT '',
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
  `os_info` varchar(255) NOT NULL DEFAULT '' COMMENT 'os|cpu cores|mem',
  `description` varchar(255) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `state` tinyint(4) UNSIGNED NOT NULL DEFAULT 1 COMMENT '0.disable 1.enable',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_namespace_ip` (`namespace_id`,`ip`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

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
  `name` varchar(30) NOT NULL DEFAULT '',
  `contact` varchar(255) NOT NULL DEFAULT '',
  `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0.disable 1.enable',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `last_login_time` datetime DEFAULT NULL,
  `super_manager` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'the mark of super admin',
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
  `role` varchar(20) NOT NULL,
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_namespace_user` (`namespace_id`,`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `system_config` (
    `id`    int(10) unsigned NOT NULL AUTO_INCREMENT,
    `key`   varchar(255) NOT NULL DEFAULT '',
    `value` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `udx_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT IGNORE INTO `user`(`id`, `account`, `password`, `name`, `contact`, `state`, `super_manager`) VALUES (1, 'admin', '$2a$10$89ZJ2xeJj35GOw11Qiucr.phaEZP4.kBX6aKTs7oWFp1xcGBBgijm', '超管', '', 1, 1);
INSERT IGNORE INTO `namespace`(`id`, `name`) VALUES (1, 'goploy');
INSERT IGNORE INTO `namespace_user`(`id`, `namespace_id`, `user_id`, `role`) VALUES (1, 1, 1, 'admin');
INSERT IGNORE INTO `system_config` (`id`, `key`, `value`) VALUES (1, 'version', '1.4.3');

