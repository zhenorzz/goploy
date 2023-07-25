ALTER TABLE `user`
    ADD COLUMN `api_key` char(32) NOT NULL DEFAULT '' AFTER `super_manager`;