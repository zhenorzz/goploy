CREATE TABLE IF NOT EXISTS `server_process` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `server_id` int(10) unsigned NOT NULL DEFAULT '0',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `start` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `stop` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `restart` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_server` (`server_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (70, 23, 'ShowServerProcessPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (71, 23, 'AddServerProcess', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (72, 23, 'EditServerProcess', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (73, 23, 'DeleteServerProcess', 0, '');
