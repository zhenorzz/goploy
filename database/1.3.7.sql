CREATE TABLE IF NOT EXISTS `goploy`.`server_agent_log` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `server_id` int(10) unsigned NOT NULL DEFAULT '0',
    `type` tinyint(3) unsigned NOT NULL DEFAULT '0',
    `item` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
    `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
    `report_time` datetime NOT NULL,
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_server_type_item_time` (`server_id`,`type`,`item`,`report_time`) USING BTREE,
    KEY `idx_server_item_time` (`server_id`,`item`,`report_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;