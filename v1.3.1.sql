ALTER TABLE `goploy`.`project_task`
    CHANGE COLUMN `commit_id` `commit` char(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `branch`;
