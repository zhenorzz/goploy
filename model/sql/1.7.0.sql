ALTER TABLE `monitor`
DROP COLUMN `url`,
DROP COLUMN `notify_times`,
ADD COLUMN `type` smallint UNSIGNED NOT NULL DEFAULT 0 AFTER `name`,
ADD COLUMN `target` json NULL AFTER `type`,
ADD COLUMN `silent_cycle` int UNSIGNED NOT NULL DEFAULT 0 AFTER `times`;