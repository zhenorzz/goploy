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

ALTER TABLE `server`
ADD COLUMN `jump_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `password`,
ADD COLUMN `jump_port` smallint UNSIGNED NOT NULL DEFAULT 0 AFTER `jump_ip`,
ADD COLUMN `jump_owner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `jump_port`,
ADD COLUMN `jump_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `jump_owner`,
ADD COLUMN `jump_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `jump_path`;
