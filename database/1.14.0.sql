ALTER TABLE monitor ADD success_script text NOT NULL;
ALTER TABLE monitor ADD fail_script text NOT NULL;
ALTER TABLE monitor ADD success_server_id INT ( 10 ) UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE monitor ADD fail_server_id INT ( 10 ) UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE project ADD label VARCHAR ( 6382 ) NOT NULL DEFAULT '';

INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (78, 23, 'SFTPEditFile', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (79, 23, 'ShowServerNginxPage', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (80, 23, 'ManageServerNginx', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (81, 23, 'AddNginxConfig', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (82, 23, 'EditNginxConfig', 0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (83, 23, 'DeleteNginxConfig', 0, '');

