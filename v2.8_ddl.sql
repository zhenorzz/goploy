CREATE TABLE `goploy`.`project_task` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0',
  `commit_id` char(40) NOT NULL DEFAULT '',
  `date` datetime DEFAULT NULL,
  `state` tinyint(4) unsigned NOT NULL DEFAULT '1',
  `is_run` tinyint(4) unsigned NOT NULL DEFAULT '0',
  `creator_id` int(10) unsigned NOT NULL DEFAULT '0',
  `creator` varchar(255) NOT NULL DEFAULT '',
  `editor_id` int(10) unsigned NOT NULL DEFAULT '0',
  `editor` varchar(255) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;
