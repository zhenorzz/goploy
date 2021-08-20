ALTER TABLE `goploy`.`project`
    ADD COLUMN `repo_type` varchar(255) NOT NULL DEFAULT '' COMMENT 'git | svn' AFTER `name`,
    ADD COLUMN `symlink_backup_number` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT 'symlink backup number' AFTER `symlink_path`;

UPDATE `goploy`.`project` SET `repo_type` = 'git', `symlink_backup_number` = 10;