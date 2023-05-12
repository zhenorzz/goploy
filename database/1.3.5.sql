ALTER TABLE `goploy`.`server`
    ADD COLUMN `os_info` varchar(255) NOT NULL DEFAULT '' COMMENT 'os|cpu cores|mem' AFTER `description`;