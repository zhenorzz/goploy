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

INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`) VALUES (74, 1, 'ShowOperationLogPage', 0, '');

ALTER TABLE `project` ADD COLUMN `deploy_server_mode` varchar(255) NOT NULL DEFAULT '' COMMENT 'serial | parallel' AFTER `transfer_option`;

UPDATE `project` SET `deploy_server_mode` = 'parallel';