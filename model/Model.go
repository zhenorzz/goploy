package model

import (
	"database/sql"
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Pagination struct
type Pagination struct {
	Page  uint64 `json:"page"`
	Rows  uint64 `json:"rows"`
	Total uint64 `json:"total"`
}

// state type
const (
	Fail = iota
	Success
)

// state type
const (
	Disable = iota
	Enable
)

// DB init when the program start
var DB *sql.DB

// Init DB
func Init() {
	dbType := os.Getenv("DB_TYPE")
	dbConn := os.Getenv("DB_CONN")
	var err error
	DB, err = sql.Open(dbType, dbConn)
	if err != nil {
		log.Fatal(err)
	}
}

// PaginationFrom param return pagination struct
func PaginationFrom(param url.Values) (Pagination, error) {
	page, err := strconv.ParseUint(param.Get("page"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("invalid page")
	}
	rows, err := strconv.ParseUint(param.Get("rows"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("invalid rows")
	}
	pagination := Pagination{Page: page, Rows: rows}
	return pagination, nil
}

const ddl string = "CREATE DATABASE IF NOT EXISTS `goploy`; CREATE TABLE IF NOT EXISTS `goploy`.`log`( `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, `type` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT 'log type', `ip` int(10) UNSIGNED NOT NULL DEFAULT 0, `desc` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'description', `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '', `create_time` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '', PRIMARY KEY (`id`) USING BTREE, INDEX `idx_create_time`(`create_time`) USING BTREE) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`project` ( `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0, `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'project name', `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'repository url', `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'project deploy path', `symlink_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '(ln -sfn symlink_path/uuid project_path)', `environment` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '1.production 2.pre-release 3.test 4.development', `branch` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'master' COMMENT 'repository branch', `after_pull_script_mode` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'sh|php|py|...', `after_pull_script` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '', `after_deploy_script_mode` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'sh|php|py|...', `after_deploy_script` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '', `rsync_option` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'rsync options', `auto_deploy` tinyint(4) UNSIGNED NOT NULL DEFAULT 1 COMMENT '0.disable 1.webhook', `state` tinyint(4) UNSIGNED NOT NULL DEFAULT 1 COMMENT '0.disable 1.enable', `deploy_state` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '0.not deploy 1.deploying 2.success 3.fail', `publisher_id` int(10) UNSIGNED NOT NULL DEFAULT 0, `publisher_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `last_publish_token` char(36) CHARACTER SET utf8mb4 NOT NULL DEFAULT '', `notify_type` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '1.weixin 2.ding talk 3.feishu 255.custom', `notify_target` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE ) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`project_server` ( `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, `project_id` int(10) UNSIGNED NOT NULL DEFAULT 0, `server_id` int(10) UNSIGNED NOT NULL DEFAULT 0, `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE, UNIQUE INDEX `uk_project_server`(`project_id`, `server_id`) USING BTREE ) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`project_user` ( `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, `project_id` int(10) UNSIGNED NOT NULL DEFAULT 0, `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0, `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE, UNIQUE INDEX `uk_project_user`(`project_id`, `user_id`) USING BTREE ) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`project_task` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `project_id` int(10) unsigned NOT NULL DEFAULT '0', `commit_id` char(40) NOT NULL DEFAULT '', `date` datetime DEFAULT NULL, `state` tinyint(4) unsigned NOT NULL DEFAULT '1', `is_run` tinyint(4) unsigned NOT NULL DEFAULT '0', `creator_id` int(10) unsigned NOT NULL DEFAULT '0', `creator` varchar(255) NOT NULL DEFAULT '', `editor_id` int(10) unsigned NOT NULL DEFAULT '0', `editor` varchar(255) NOT NULL DEFAULT '', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE, KEY `index_project_update` (`project_id`,`update_time`) USING BTREE COMMENT 'project_id,update_time' ) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`publish_trace` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `token` char(36) CHARACTER SET utf8mb4 NOT NULL DEFAULT '', `project_id` int(10) unsigned NOT NULL DEFAULT '0', `project_group_id` int(10) unsigned NOT NULL DEFAULT '0', `project_name` varchar(255) NOT NULL DEFAULT '', `detail` longtext NOT NULL, `state` tinyint(4) unsigned NOT NULL DEFAULT '1', `publisher_id` int(10) unsigned NOT NULL DEFAULT '0', `publisher_name` varchar(255) NOT NULL DEFAULT '', `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1.before pull 2.pulled 3.after pull 4.before deploy 5.deploy 6.after deploy', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, `ext` longtext NOT NULL, PRIMARY KEY (`id`) USING BTREE, KEY `idx_project_id` (`project_id`) USING BTREE COMMENT 'project_id' ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4; CREATE TABLE IF NOT EXISTS `goploy`.`monitor` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `namespace_id` int(10) unsigned NOT NULL, `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL, `domain` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL, `port` smallint(5) unsigned NOT NULL DEFAULT '80', `second` int(10) unsigned NOT NULL DEFAULT '1' COMMENT 'How many seconds to run', `times` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT 'How many times of failures', `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `notify_type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '1.weixin 2.ding talk 3.feishu 255.custom', `notify_target` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `notify_times` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT 'Notify times', `state` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '0.disable 1.enable', `error_content` text COLLATE utf8mb4_general_ci NOT NULL, `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`server` ( `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, `namespace_id` int(10) UNSIGNED NOT NULL DEFAULT 0, `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `ip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `port` smallint(10) UNSIGNED NOT NULL DEFAULT 22, `owner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `last_publish_token` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, `state` tinyint(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '0.disable 1.enable', PRIMARY KEY (`id`) USING BTREE, UNIQUE KEY `uk_namespace_ip` (`namespace_id`,`ip`) USING BTREE ) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`crontab` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `namespace_id` int(10) unsigned NOT NULL DEFAULT 0, `command` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `command_md5` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'command md5 for replace', `creator_id` int(10) unsigned NOT NULL DEFAULT '0', `creator` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `editor_id` int(10) unsigned NOT NULL DEFAULT '0', `editor` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE, UNIQUE KEY `uk_command_md5` (`namespace_id`,`command_md5`) USING BTREE ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`crontab_server` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `crontab_id` int(10) unsigned NOT NULL, `server_id` int(10) unsigned NOT NULL, `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE, UNIQUE KEY `idx_crontab_server` (`crontab_id`,`server_id`) USING BTREE ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`template` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `package_id_str` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL, `script` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL, `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`package` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `size` int(10) unsigned NOT NULL DEFAULT '0', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`install_trace` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `token` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `server_id` int(10) unsigned NOT NULL DEFAULT '0', `server_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `detail` longtext NOT NULL, `state` tinyint(4) unsigned NOT NULL DEFAULT '1', `operator_id` int(10) unsigned NOT NULL DEFAULT '0', `operator_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1.rsync 2.ssh 3.script', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, `ext` text NOT NULL, PRIMARY KEY (`id`) USING BTREE, KEY `idx_project_id` (`server_id`) USING BTREE COMMENT 'project_id' ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`user` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `account` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `password` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `contact` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0.disable 1.enable', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, `last_login_time` datetime DEFAULT NULL, `super_manager` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'the mark of super admin', PRIMARY KEY (`id`) USING BTREE ) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`namespace` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '', `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE, UNIQUE KEY `uk_name` (`name`) ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; CREATE TABLE IF NOT EXISTS `goploy`.`namespace_user` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `namespace_id` int(10) unsigned NOT NULL, `user_id` int(10) unsigned NOT NULL, `role` varchar(20) NOT NULL, `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) USING BTREE, UNIQUE KEY `uk_namespace_user` (`namespace_id`,`user_id`) USING BTREE ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;"
const dml string = "INSERT INTO `goploy`.`user`(`id`, `account`, `password`, `name`, `contact`, `state`, `super_manager`) VALUES (1, 'admin', '$2a$10$89ZJ2xeJj35GOw11Qiucr.phaEZP4.kBX6aKTs7oWFp1xcGBBgijm', '超管', '', 1, 1); INSERT INTO `goploy`.`namespace`(`id`, `name`) VALUES (1, 'goploy'); INSERT INTO `goploy`.`namespace_user`(`id`, `namespace_id`, `user_id`, `role`) VALUES (1, 1, 1, 'admin');"

// ImportSQL -
func ImportSQL(db *sql.DB) error {
	for _, query := range strings.Split(ddl, ";") {
		if len(query) == 0 {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	for _, query := range strings.Split(dml, ";") {
		if len(query) == 0 {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
