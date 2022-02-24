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