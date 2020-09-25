ALTER TABLE `goploy`.`monitor`
ADD COLUMN `notify_times` smallint(5) UNSIGNED NOT NULL DEFAULT 1 AFTER `notify_target`,
ADD COLUMN `error_content` varchar(1000) NOT NULL DEFAULT '' AFTER `notify_times`;