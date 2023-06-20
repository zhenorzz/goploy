ALTER TABLE `project`
MODIFY COLUMN `after_deploy_script_mode` varchar(20) NOT NULL DEFAULT '' COMMENT 'deprecated after v1.14.1' AFTER `review_url`,
MODIFY COLUMN `after_deploy_script` text NOT NULL COMMENT 'deprecated after v1.14.1' AFTER `after_deploy_script_mode`,
MODIFY COLUMN `after_pull_script_mode` varchar(20) NOT NULL DEFAULT '' COMMENT 'deprecated after v1.14.1' AFTER `after_deploy_script`,
MODIFY COLUMN `after_pull_script` text NOT NULL COMMENT 'deprecated after v1.14.1' AFTER `after_pull_script_mode`
ADD COLUMN `script` json NOT NULL COMMENT 'script' AFTER `review_url`;

UPDATE project
SET script = JSON_OBJECT(
    'after_pull',
    JSON_OBJECT( 'mode', after_pull_script_mode, 'content', after_pull_script ),
    'after_deploy',
    JSON_OBJECT( 'mode', after_deploy_script_mode, 'content', after_deploy_script ),
    'deploy_success',
    JSON_OBJECT( 'mode', '', 'content', '' )
)