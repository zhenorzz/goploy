ALTER TABLE `goploy`.`monitor`
DROP COLUMN `port`,
CHANGE COLUMN `domain` `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL AFTER `name`;