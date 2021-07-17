CREATE TABLE IF NOT EXISTS `goploy`.`system_config` (
    `id`    int(10) unsigned NOT NULL AUTO_INCREMENT,
    `key`   varchar(255) NOT NULL DEFAULT '',
    `value` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `udx_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

REPLACE INTO `goploy`.`system_config` (`key`, `value`) VALUES ('version', '');