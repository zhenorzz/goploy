ALTER TABLE `goploy`.`project`
ADD COLUMN `review` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '0.disable 1.enable' AFTER `symlink_path`,
ADD COLUMN `review_url` varchar(1000) NOT NULL DEFAULT '' COMMENT 'review notification link' AFTER `review`;