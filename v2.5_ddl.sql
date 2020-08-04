
ALTER TABLE `goploy`.`project`
ADD COLUMN `after_deploy_script_mode` varchar(20) NOT NULL DEFAULT '' COMMENT '脚本类型(默认bash)' AFTER `symlink_path`,
ADD COLUMN `after_pull_script_mode` varchar(20) NOT NULL DEFAULT '' COMMENT '脚本类型(默认bash)' AFTER `after_deploy_script`;