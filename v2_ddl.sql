ALTER TABLE `goploy`.`group` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`install_trace` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`package` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`project` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`project_server` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`project_user` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`publish_trace` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`server` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`template` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();

ALTER TABLE `goploy`.`user` 
DROP COLUMN `create_time`,
DROP COLUMN `update_time`,
DROP COLUMN `last_login_time`,
ADD COLUMN `last_login_time` datetime DEFAULT NULL,
ADD COLUMN `insert_time` datetime NOT NULL DEFAULT now(),
ADD COLUMN `update_time` datetime NOT NULL DEFAULT now();