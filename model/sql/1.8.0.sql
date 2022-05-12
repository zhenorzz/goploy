ALTER TABLE `project`
    CHANGE COLUMN `rsync_option` `transfer_option` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `branch`,
    ADD COLUMN `transfer_type` varchar(255) NOT NULL DEFAULT '' AFTER `branch`;
UPDATE `project` SET `transfer_type` = 'rsync';