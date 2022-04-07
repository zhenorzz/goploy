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

INSERT INTO `role`(`id`, `name`, `description`, `insert_time`, `update_time`) VALUES (1, 'manager', '', '2021-05-09 10:03:17', '2022-04-03 13:05:27');
INSERT INTO `role`(`id`, `name`, `description`, `insert_time`, `update_time`) VALUES (2, 'group-manager', '', '2021-05-16 17:10:57', '2022-04-03 13:05:29');
INSERT INTO `role`(`id`, `name`, `description`, `insert_time`, `update_time`) VALUES (3, 'member', '', '2022-04-03 23:23:59', '2022-04-05 18:41:49');

INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (1, 0, 'Log', 0, '', '2021-05-09 10:03:17', '2022-04-03 14:22:21');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (2, 1, 'ShowLoginLogPage', 0, '', '2021-05-16 17:10:57', '2022-04-03 14:22:41');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (3, 1, 'ShowPublishLogPage', 0, '', '2022-04-03 13:05:35', '2022-04-03 14:23:05');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (4, 1, 'ShowSFTPLogPage', 0, '', '2022-04-03 14:23:14', '2022-04-03 14:23:14');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (5, 1, 'ShowTerminalLogPage', 0, '', '2022-04-03 14:23:28', '2022-04-03 14:23:28');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (6, 1, 'ShowTerminalRecord', 0, '', '2022-04-03 14:23:38', '2022-04-03 14:23:38');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (7, 0, 'Member', 0, '', '2022-04-03 14:24:08', '2022-04-03 14:24:08');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (8, 7, 'ShowMemberPage', 0, '', '2022-04-03 14:24:32', '2022-04-03 14:24:32');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (9, 7, 'AddMember', 0, '', '2022-04-03 14:24:51', '2022-04-03 14:27:35');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (10, 7, 'EditMember', 0, '', '2022-04-03 14:24:58', '2022-04-03 14:27:35');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (11, 7, 'DeleteMember', 0, '', '2022-04-03 14:25:08', '2022-04-03 14:27:35');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (12, 0, 'Namespace', 0, '', '2022-04-03 14:25:45', '2022-04-03 14:38:33');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (13, 12, 'ShowNamespacePage', 0, '', '2022-04-03 14:25:45', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (14, 12, 'AddNamespace', 0, '', '2022-04-03 14:39:05', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (15, 12, 'EditNamespace', 0, '', '2022-04-03 14:39:18', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (16, 12, 'GetNamespaceUserList', 0, '', '2022-04-03 14:39:52', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (17, 12, 'AddNamespaceUser', 0, '', '2022-04-03 14:39:25', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (18, 12, 'DeleteNamespaceUser', 0, '', '2022-04-03 14:39:47', '2022-04-05 19:11:05');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (19, 12, 'ShowRolePage', 0, '', '2022-04-03 14:39:57', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (20, 12, 'AddRole', 0, '', '2022-04-03 14:40:00', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (21, 12, 'EditRole', 0, '', '2022-04-03 14:40:07', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (22, 12, 'DeleteRole', 0, '', '2022-04-03 14:40:23', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (23, 12, 'EditPermission', 0, '', '2022-04-03 18:03:41', '2022-04-05 16:46:02');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (24, 0, 'Server', 0, '', '2022-04-04 15:23:28', '2022-04-04 15:23:28');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (25, 24, 'ShowServerPage', 0, '', '2022-04-04 15:23:41', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (26, 24, 'AddServer', 0, '', '2022-04-04 15:23:49', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (27, 24, 'EditServer', 0, '', '2022-04-04 15:23:54', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (28, 24, 'SwitchServerState', 0, '', '2022-04-04 15:24:45', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (29, 24, 'InstallAgent', 0, '', '2022-04-04 15:24:54', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (30, 24, 'ImportCSV', 0, '', '2022-04-04 15:25:02', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (31, 24, 'ShowServerMonitorPage', 0, '', '2022-04-04 15:25:40', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (32, 24, 'AddServerWarningRule', 0, '', '2022-04-04 15:25:40', '2022-04-06 18:32:11');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (33, 24, 'EditServerWarningRule', 0, '', '2022-04-04 15:25:40', '2022-04-06 18:32:12');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (34, 24, 'DeleteServerWarningRule', 0, '', '2022-04-04 15:25:40', '2022-04-06 18:32:14');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (35, 24, 'ShowTerminalPage', 0, '', '2022-04-04 15:25:59', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (36, 24, 'ShowSftpFilePage', 0, '', '2022-04-04 15:26:07', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (37, 24, 'SFTPUploadFile', 0, '', '2022-04-04 15:26:26', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (38, 24, 'SFTPPreviewFile', 0, '', '2022-04-04 15:26:26', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (39, 24, 'SFTPDownloadFile', 0, '', '2022-04-04 15:26:35', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (40, 24, 'ShowCronPage', 0, '', '2022-04-04 15:27:49', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (41, 24, 'AddCron', 0, '', '2022-04-04 15:27:52', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (42, 24, 'EditCron', 0, '', '2022-04-04 15:28:11', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (43, 24, 'DeleteCron', 0, '', '2022-04-04 15:28:16', '2022-04-05 19:13:29');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (44, 0, 'Project', 0, '', '2022-04-04 15:29:54', '2022-04-04 15:29:54');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (45, 44, 'ShowProjectPage', 0, '', '2022-04-04 15:29:54', '2022-04-06 18:28:44');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (46, 44, 'GetAllProjectList', 0, '', '2022-04-04 15:30:10', '2022-04-06 18:28:44');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (47, 44, 'GetBindProjectList', 0, '', '2022-04-04 15:30:18', '2022-04-06 18:28:44');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (48, 44, 'AddProject', 0, '', '2022-04-04 15:30:23', '2022-04-06 18:28:44');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (49, 44, 'EditProject', 0, '', '2022-04-04 15:30:28', '2022-04-06 18:28:44');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (50, 44, 'DeleteProject', 0, '', '2022-04-04 15:31:40', '2022-04-06 18:28:44');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (51, 44, 'SwitchProjectWebhook', 0, '', '2022-04-04 15:32:29', '2022-04-06 18:28:44');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (52, 0, 'Monitor', 0, '', '2022-04-04 15:34:42', '2022-04-04 15:34:42');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (53, 52, 'ShowMonitorPage', 0, '', '2022-04-04 15:35:12', '2022-04-06 18:28:47');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (54, 52, 'AddMonitor', 0, '', '2022-04-04 15:35:20', '2022-04-06 18:28:47');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (55, 52, 'EditMonitor', 0, '', '2022-04-04 15:35:25', '2022-04-06 18:28:48');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (56, 52, 'DeleteMonitor', 0, '', '2022-04-04 15:35:30', '2022-04-06 18:28:48');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (57, 0, 'Deploy', 0, '', '2022-04-04 15:37:04', '2022-04-04 15:37:04');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (58, 57, 'ShowDeployPage', 0, '', '2022-04-04 15:37:04', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (59, 57, 'GetAllDeployList', 0, '', '2022-04-04 15:37:19', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (60, 57, 'GetBindDeployList', 0, '', '2022-04-04 15:37:27', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (61, 57, 'DeployDetail', 0, '', '2022-04-04 15:40:14', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (62, 57, 'DeployProject', 0, '', '2022-04-04 15:37:52', '2022-04-07 14:27:25');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (63, 57, 'DeployResetState', 0, '', '2022-04-07 14:13:46', '2022-04-07 14:27:46');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (64, 57, 'GreyDeploy', 0, '', '2022-04-07 14:15:18', '2022-04-07 14:27:51');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (65, 57, 'DeployRollback', 0, '', '2022-04-04 15:38:24', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (66, 57, 'DeployReview', 0, '', '2022-04-04 15:38:24', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (67, 57, 'DeployTask', 0, '', '2022-04-04 15:38:37', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (68, 57, 'FileCompare', 0, '', '2022-04-04 15:38:46', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (69, 57, 'FileSync', 0, '', '2022-04-04 15:38:49', '2022-04-06 18:28:59');
INSERT INTO `permission` (`id`, `pid`, `name`, `sort`, `description`, `insert_time`, `update_time`) VALUES (70, 57, 'ProcessManager', 0, '', '2022-04-04 15:39:02', '2022-04-06 18:28:59');
