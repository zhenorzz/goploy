ALTER TABLE `user`
    ADD COLUMN `api_key` varchar(255) NOT NULL DEFAULT '' AFTER `super_manager`;