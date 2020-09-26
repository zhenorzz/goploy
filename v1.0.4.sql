ALTER TABLE `goploy`.`project`
ADD COLUMN `review` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '0.disable 1.enable' AFTER `symlink_path`,
ADD COLUMN `review_url` varchar(1000) NOT NULL DEFAULT '' COMMENT 'review notification link' AFTER `review`;

CREATE TABLE IF NOT EXISTS `goploy`.`project_review` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0',
  `commit_id` char(40) NOT NULL DEFAULT '',
  `review_url` varchar(1000) NOT NULL DEFAULT '' COMMENT 'review notification link',
  `state` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '0.Pending 1.Approve 2.Deny',
  `creator_id` int(10) unsigned NOT NULL DEFAULT '0',
  `creator` varchar(255) NOT NULL DEFAULT '',
  `editor_id` int(10) unsigned NOT NULL DEFAULT '0',
  `editor` varchar(255) NOT NULL DEFAULT '',
  `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
