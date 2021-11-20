CREATE TABLE IF NOT EXISTS `goploy`.`server_monitor` (
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