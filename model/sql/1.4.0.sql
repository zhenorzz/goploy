CREATE TABLE IF NOT EXISTS `goploy`.`cron` (
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

CREATE TABLE IF NOT EXISTS `goploy`.`cron_log` (
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