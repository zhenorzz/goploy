ALTER TABLE `server_process`
    ADD COLUMN `namespace_id` int UNSIGNED NOT NULL,
    ADD COLUMN `items` json NULL COMMENT '{name: string, command: string}[]';

UPDATE `server_process` SET `items` = '[]';

INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (75, 23, 'SFTPTransferFile', 0, '');
