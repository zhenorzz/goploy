ALTER TABLE `goploy`.`project_task`
ADD COLUMN `branch` varchar(255) NOT NULL DEFAULT '' AFTER `project_id`;
ALTER TABLE `goploy`.`project_review`
ADD COLUMN `branch` varchar(255) NOT NULL DEFAULT '' AFTER `project_id`;