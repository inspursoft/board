drop database if exists board;
create database board charset = utf8;

use board;

DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `project`;
DROP TABLE IF EXISTS `project_member`;
DROP TABLE IF EXISTS `role`;


CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(45) NOT NULL,
  `password` varchar(45) NOT NULL,
  `email` varchar(45) NOT NULL,
  `realname` varchar(45) DEFAULT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `deleted` int(11) DEFAULT NULL,
  `system_admin` int(11) DEFAULT NULL,
  `project_admin` int(11) DEFAULT NULL,
  `reset_uuid` varchar(255) DEFAULT NULL,
  `salt` varchar(255) DEFAULT NULL,
  `creation_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


INSERT INTO `board`.`user` (`username`, `password`, `email`, `realname`, `comment`, `creation_time`, `update_time`, `deleted`, `system_admin`, `project_admin`)
  VALUES ('admin', 'Board12345', 'admin@inspur.com', 'admin', 'admin user', now(), now(), 0, 1, 1);

CREATE TABLE `board`.`project` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(50) NULL,
  `comment` VARCHAR(255) NULL,
  `creation_time` DATETIME NULL,
  `update_time` DATETIME NULL,
  `deleted` SMALLINT(1) NULL,
  `owner_id` INT NULL,
  `owner_name` VARCHAR(45) NULL,
  `public` SMALLINT(1) NULL,
  `toggleable` SMALLINT(1) NULL,
  `current_user_role_id` INT NULL,
  `service_count` INT NULL,
  PRIMARY KEY (`id`));

INSERT INTO `board`.`project`
 (`id`, `name`, `comment`, `creation_time`, `update_time`, `deleted`, `owner_id`, 
  `owner_name`, `public`, `toggleable`, `current_user_role_id`, `service_count`)
 VALUES
 (1, 'library', 'library comment', now(), now(), 0, 1,'admin', 1, 1, 1, 0);


CREATE TABLE `project_member` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `project_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`,`user_id`,`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `board`.`project_member`
 (`id`, `user_id`, `project_id`, `role_id`)
 VALUES
 (2, 1, 1, 2);

CREATE TABLE `board`.`role` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL,
  `comment` VARCHAR(45) NULL,
  PRIMARY KEY (`id`));

INSERT INTO `board`.`role` (id, name, comment) 
  VALUES (1, 'projectAdmin', 'Project Admin'),
         (2, 'developer', 'Developer'),
         (3, 'visitor', 'Visitor');

-- --------------------------------------------------
--  Table Structure for `model/get_resource.Pods`
-- --------------------------------------------------
DROP TABLE IF EXISTS `pod`;
DROP TABLE IF EXISTS `node`;
DROP TABLE IF EXISTS `service`;
DROP TABLE IF EXISTS `dashboard_service_second`;
DROP TABLE IF EXISTS `dashboard_service_minute`;
DROP TABLE IF EXISTS `dashboard_service_hour`;
DROP TABLE IF EXISTS `dashboard_service_day`;
DROP TABLE IF EXISTS `log`;
-- --------------------------------------------------
--  Table Structure for log
-- --------------------------------------------------
CREATE TABLE `log_collector` (
  `id`                    BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `pod_name`              VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_hostIP`            VARCHAR(30)           NOT NULL DEFAULT '',
  `containers_numbers`    VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_creat_time`        TIMESTAMP             NOT NULL DEFAULT NOW(),
  `all_pod_numbers`       VARCHAR(30)           NOT NULL DEFAULT '',
  `all_container_numbers` VARCHAR(30)           NOT NULL DEFAULT '',
  `service_name`          VARCHAR(30)           NOT NULL DEFAULT '',
  `service_numbers`       VARCHAR(30)           NOT NULL DEFAULT '',
  `service_creat_time`    TIMESTAMP             NOT NULL DEFAULT NOW(),
  `record_time`           TIMESTAMP             NOT NULL DEFAULT NOW()
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

-- --------------------------------------------------
--  Table Structure for `model/get_resource.Pods`
-- --------------------------------------------------
CREATE TABLE `pod` (
  `id`                    BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `pod_name`              VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_hostIP`            VARCHAR(30)           NOT NULL DEFAULT '',
  `containers_numbers`    VARCHAR(30)           NOT NULL DEFAULT '',
  `creat_time`            TIMESTAMP             NOT NULL DEFAULT NOW(),
  `record_time`           TIMESTAMP             NOT NULL DEFAULT NOW(),
  `all_pod_numbers`       VARCHAR(30)           NOT NULL DEFAULT '',
  `all_container_numbers` VARCHAR(30)           NOT NULL DEFAULT ''
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

-- --------------------------------------------------
--  Table Structure for `model/get_resource.Nodes`
-- --------------------------------------------------
CREATE TABLE `node` (
  `id`               BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `node_name`        VARCHAR(30)           NOT NULL DEFAULT '',
  `numbers_cpu_core` VARCHAR(30)           NOT NULL DEFAULT '',
  `numbers_gpu_core` VARCHAR(30)           NOT NULL DEFAULT '',
  `memory_size`      VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_limit`        VARCHAR(30)           NOT NULL DEFAULT '',
  `create_time`      TIMESTAMP             NOT NULL DEFAULT NOW(),
  `record_time`      TIMESTAMP             NOT NULL DEFAULT NOW()
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

-- --------------------------------------------------
--  Table Structure for `model/get_resource.Services`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `service` (
  `uid`             BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `service_name`    VARCHAR(30)           NOT NULL DEFAULT '',
  `service_numbers` VARCHAR(30)           NOT NULL DEFAULT '',
  `creat_time`      TIMESTAMP             NOT NULL DEFAULT NOW(),
  `record_time`     TIMESTAMP             NOT NULL DEFAULT NOW()
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;
-- --------------------------------------------------
--  Table Structure for `model/dashboard`
-- --------------------------------------------------

CREATE TABLE `dashboard_service_second` (
  `id`               BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `service_numbers`  VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_number`       VARCHAR(30)           NOT NULL DEFAULT '',
  `container_number` VARCHAR(30)           NOT NULL DEFAULT '',
  `record_time`      TIMESTAMP             NOT NULL DEFAULT NOW()
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

-- --------------------------------------------------
--  Table Structure for `model/dashboard`
-- --------------------------------------------------

CREATE TABLE `dashboard_service_minute` (
  `id`               BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `service_numbers`  VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_number`       VARCHAR(30)           NOT NULL DEFAULT '',
  `container_number` VARCHAR(30)           NOT NULL DEFAULT '',
  `record_time`      TIMESTAMP             NOT NULL DEFAULT NOW()
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

-- --------------------------------------------------
--  Table Structure for `model/dashboard`
-- --------------------------------------------------

CREATE TABLE `dashboard_service_hour` (
  `id`               BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `service_numbers`  VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_number`       VARCHAR(30)           NOT NULL DEFAULT '',
  `container_number` VARCHAR(30)           NOT NULL DEFAULT '',
  `record_time`      TIMESTAMP             NOT NULL DEFAULT NOW()
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

-- --------------------------------------------------
--  Table Structure for `model/dashboard`
-- --------------------------------------------------

CREATE TABLE `dashboard_service_day` (
  `id`               BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `service_numbers`  VARCHAR(30)           NOT NULL DEFAULT '',
  `pod_number`       VARCHAR(30)           NOT NULL DEFAULT '',
  `container_number` VARCHAR(30)           NOT NULL DEFAULT '',
  `record_time`      TIMESTAMP             NOT NULL DEFAULT NOW()
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;


