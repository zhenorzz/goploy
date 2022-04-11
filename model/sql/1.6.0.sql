CREATE TABLE IF NOT EXISTS `role` (
        `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
        `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
        `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
        `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`) USING BTREE,
        UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `permission` (
      `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
      `pid` int(10) unsigned NOT NULL DEFAULT '0',
      `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
      `sort` int(10) NOT NULL DEFAULT '0',
      `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
      `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
      `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY (`id`) USING BTREE,
      UNIQUE KEY `uk_name` (`name`),
      KEY `idx_pid` (`pid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `role_permission` (
       `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
       `role_id` int(10) unsigned NOT NULL,
       `permission_id` int(10) unsigned NOT NULL,
       `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
       `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       PRIMARY KEY (`id`) USING BTREE,
       UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE `namespace_user` ADD COLUMN `role_id` int(10) UNSIGNED NOT NULL DEFAULT 0 AFTER `user_id`;

UPDATE `namespace_user` SET `role_id` = 1 WHERE `role` = 'manager';
UPDATE `namespace_user` SET `role_id` = 2 WHERE `role` = 'group-manager';
UPDATE `namespace_user` SET `role_id` = 3 WHERE `role` = 'member';

ALTER TABLE `namespace_user` DROP COLUMN `role`;

INSERT INTO `role`(`id`, `name`, `description`) VALUES (1, 'manager', '');
INSERT INTO `role`(`id`, `name`, `description`) VALUES (2, 'group-manager', '');
INSERT INTO `role`(`id`, `name`, `description`) VALUES (3, 'member', '');

INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (1, 0, 'Log', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (2, 1, 'ShowLoginLogPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (3, 1, 'ShowPublishLogPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (4, 1, 'ShowSFTPLogPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (5, 1, 'ShowTerminalLogPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (6, 1, 'ShowTerminalRecord', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (7, 0, 'Member', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (8, 7, 'ShowMemberPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (9, 7, 'AddMember', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (10, 7, 'EditMember', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (11, 7, 'DeleteMember', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (12, 0, 'Namespace', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (13, 12, 'ShowNamespacePage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (14, 12, 'AddNamespace', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (15, 12, 'EditNamespace', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (16, 12, 'AddNamespaceUser', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (17, 12, 'DeleteNamespaceUser', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (18, 12, 'ShowRolePage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (19, 12, 'AddRole', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (20, 12, 'EditRole', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (21, 12, 'DeleteRole', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (22, 12, 'EditPermission', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (23, 0, 'Server', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (24, 23, 'ShowServerPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (25, 23, 'AddServer', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (26, 23, 'EditServer', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (27, 23, 'SwitchServerState', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (28, 23, 'InstallAgent', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (29, 23, 'ImportCSV', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (30, 23, 'ShowServerMonitorPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (31, 23, 'AddServerWarningRule', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (32, 23, 'EditServerWarningRule', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (33, 23, 'DeleteServerWarningRule', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (34, 23, 'ShowTerminalPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (35, 23, 'ShowSftpFilePage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (36, 23, 'SFTPUploadFile', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (37, 23, 'SFTPPreviewFile', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (38, 23, 'SFTPDownloadFile', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (39, 23, 'ShowCronPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (40, 23, 'AddCron', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (41, 23, 'EditCron', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (42, 23, 'DeleteCron', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (43, 0, 'Project', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (44, 43, 'ShowProjectPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (45, 43, 'GetAllProjectList', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (46, 43, 'GetBindProjectList', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (47, 43, 'AddProject', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (48, 43, 'EditProject', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (49, 43, 'DeleteProject', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (50, 43, 'SwitchProjectWebhook', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (51, 0, 'Monitor', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (52, 51, 'ShowMonitorPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (53, 51, 'AddMonitor', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (54, 51, 'EditMonitor', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (55, 51, 'DeleteMonitor', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (56, 0, 'Deploy', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (57, 56, 'ShowDeployPage', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (58, 56, 'GetAllDeployList', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (59, 56, 'GetBindDeployList', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (60, 56, 'DeployDetail', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (61, 56, 'DeployProject', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (62, 56, 'DeployResetState', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (63, 56, 'GreyDeploy', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (64, 56, 'DeployRollback', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (65, 56, 'DeployReview', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (66, 56, 'DeployTask', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (67, 56, 'FileCompare', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (68, 56, 'FileSync', 0, '');
INSERT INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (69, 56, 'ProcessManager', 0, '');

INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 13);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 14);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 15);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 16);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 17);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 18);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 19);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 20);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 21);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 22);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 24);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 25);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 26);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 27);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 28);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 29);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 30);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 31);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 32);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 33);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 34);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 35);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 36);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 37);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 38);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 39);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 40);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 41);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 42);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 44);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 45);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 46);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 47);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 48);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 49);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 50);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 52);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 53);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 54);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 55);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 57);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 58);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 59);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 60);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 61);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 62);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 63);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 64);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 65);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 66);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 67);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 68);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (1, 69);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 50);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 49);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 48);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 47);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 46);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 44);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 55);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 54);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 53);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 52);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 69);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 68);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 67);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 66);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 65);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 64);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 63);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 62);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 61);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 60);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 59);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (2, 57);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (3, 57);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (3, 59);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (3, 60);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (3, 61);
INSERT INTO `role_permission`(`role_id`, `permission_id`) VALUES (3, 67);
