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

INSERT IGNORE INTO `permission` (`id`, `pid`, `name`, `sort`, `description`) VALUES (77, 23, 'ShowServerScriptPage', 0, '');
