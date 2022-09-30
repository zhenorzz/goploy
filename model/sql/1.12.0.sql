ALTER TABLE `server_process`
    ADD COLUMN `namespace_id` int UNSIGNED NOT NULL,
    ADD COLUMN `items` json NULL COMMENT '{name: string, command: string}[]';

UPDATE `server_process` SET `items` = '[]';