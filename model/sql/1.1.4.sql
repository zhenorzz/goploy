ALTER TABLE `goploy`.`project`
    ADD COLUMN `repo_type` varchar(255) NOT NULL DEFAULT '' COMMENT 'git | svn' AFTER `name`;

UPDATE `goploy`.`project` SET `repo_type` = 'git';