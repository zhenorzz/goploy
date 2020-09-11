ALTER TABLE `goploy`.`monitor`
ADD COLUMN `notify_times` smallint(5) UNSIGNED NOT NULL DEFAULT 1 AFTER `notify_target`,
ADD COLUMN `error_content` text NOT NULL AFTER `notify_times`;