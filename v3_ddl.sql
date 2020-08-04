ALTER TABLE `goploy`.`user`
DROP COLUMN `role`,
DROP COLUMN `manage_group_str`,
ADD COLUMN `super_manager` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 AFTER `mobile`;

ALTER TABLE `goploy`.`server`
CHANGE COLUMN `group_id` `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0 AFTER `id`,
MODIFY COLUMN `ip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `name`,
ADD UNIQUE INDEX `uk_namespace_ip`(`namespace_id`, `ip`);

ALTER TABLE `goploy`.`project`
CHANGE COLUMN `group_id` `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0 AFTER `id`;

ALTER TABLE `goploy`.`crontab`
ADD COLUMN `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0 AFTER `id`,
DROP INDEX `uk_command_md5`,
ADD UNIQUE INDEX `uk_command_md5`(`namespace_id`, `command_md5`) USING BTREE;

UPDATE `goploy`.`user` SET `super_manager` = 1 WHERE `id` = 1

CREATE TABLE `namespace` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `namespace_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `namespace_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `role` varchar(20) NOT NULL,
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_namespace_user` (`namespace_id`,`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `monitor` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `namespace_id` int(10) unsigned NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `domain` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `port` smallint(5) unsigned NOT NULL DEFAULT '80',
  `second` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '间隔',
  `times` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT '连续失败次数',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `notify_type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '1=企业微信 2=钉钉 3=飞书 255=自定义',
  `notify_target` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `state` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '0=暂停  1=开启',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


INSERT INTO `goploy`.`namespace` SELECT * FROM `goploy`.`group`;
INSERT INTO `goploy`.`namespace_user` (`user_id`, `namespace_id`, `role`) SELECT 1 as `user_id`, `id` as `namespace_id`, 'admin' as `role` FROM `goploy`.`namespace`;
INSERT INTO `goploy`.`project_user` (`project_id`, `user_id`) SELECT `project`.`id` , `user_id` FROM `goploy`.`namespace_user` JOIN `goploy`.`project` USING(`namespace_id`) WHERE `namespace_user`.`role` IN ('admin', 'manager');
