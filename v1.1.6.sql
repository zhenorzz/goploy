ALTER TABLE `goploy`.`server`
ADD COLUMN `path` varchar(255) NOT NULL DEFAULT '' AFTER `owner`,
ADD COLUMN `password` varchar(255) NOT NULL DEFAULT '' AFTER `path`;