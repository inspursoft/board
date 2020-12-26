drop database if exists board;
create database board charset = utf8;

use board;

DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `project`;
DROP TABLE IF EXISTS `project_member`;
DROP TABLE IF EXISTS `role`;


CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(45) NOT NULL,
  `email` varchar(255) NOT NULL,
  `realname` varchar(255) DEFAULT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `deleted` SMALLINT(1) DEFAULT NULL,
  `system_admin` SMALLINT(1) DEFAULT NULL,
  `reset_uuid` varchar(255) DEFAULT NULL,
  `salt` varchar(255) DEFAULT NULL,
  `repo_token` VARCHAR(127) NULL DEFAULT NULL,  
  `creation_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `failed_times` INT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


INSERT INTO `board`.`user` (`username`, `password`, `email`, `realname`, `comment`, `creation_time`, `update_time`, `deleted`, `system_admin`)
  VALUES ('admin', 'Board12345', 'admin@inspur.com', 'admin', 'admin user', now(), now(), 0, 1);

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
  `istio_support` SMALLINT(1) NULL,
  PRIMARY KEY (`id`));

INSERT INTO `board`.`project`
 (`id`, `name`, `comment`, `creation_time`, `update_time`, `deleted`, `owner_id`, 
  `owner_name`, `public`, `toggleable`, `current_user_role_id`, `service_count`, `istio_support`)
 VALUES
 (1, 'library', 'library comment', now(), now(), 0, 1,'admin', 1, 1, 1, 0, 0);


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
 (1, 1, 1, 1);

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
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect.Node`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `node` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `node_name` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_cpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_gpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `memory_size` varchar(255) NOT NULL DEFAULT '' ,
        `pod_limit` varchar(255) NOT NULL DEFAULT '' ,
        `create_time` varchar(255) NOT NULL DEFAULT '' ,
        `ip` varchar(255) NOT NULL DEFAULT '' ,
        `cpu_usage` double precision NOT NULL DEFAULT 0 ,
        `mem_usage` double precision NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0 ,
        `storage_total` bigint NOT NULL DEFAULT 0 ,
        `storage_use` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect.Pod`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `pod` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `pod_name` varchar(255) NOT NULL DEFAULT '' ,
        `pod_hostIP` varchar(255) NOT NULL DEFAULT '' ,
        `create_time` varchar(255) NOT NULL DEFAULT '' ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;

    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect.Service`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `service` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `service_name` varchar(255) NOT NULL DEFAULT '' ,
        `create_time` varchar(255) NOT NULL DEFAULT '' ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect.ServiceKvMap`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `service_kv_map` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `name` varchar(255) NOT NULL DEFAULT '' ,
        `value` varchar(255) NOT NULL DEFAULT '' ,
        `belong` varchar(255) NOT NULL DEFAULT '' ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect.PodKvMap`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `pod_kv_map` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `name` varchar(255) NOT NULL DEFAULT '' ,
        `value` varchar(255) NOT NULL DEFAULT '' ,
        `belong` varchar(255) NOT NULL DEFAULT '' ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.ServiceDashboardSecond`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `service_dashboard_second` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `service_name` varchar(255) NOT NULL DEFAULT '' ,
        `pod_number` bigint NOT NULL DEFAULT 0 ,
        `container_number` bigint NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.ServiceDashboardMinute`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `service_dashboard_minute` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `service_name` varchar(255) NOT NULL DEFAULT '' ,
        `pod_number` bigint NOT NULL DEFAULT 0 ,
        `container_number` bigint NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.ServiceDashboardHour`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `service_dashboard_hour` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `service_name` varchar(255) NOT NULL DEFAULT '' ,
        `pod_number` bigint NOT NULL DEFAULT 0 ,
        `container_number` bigint NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.ServiceDashboardDay`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `service_dashboard_day` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `service_name` varchar(255) NOT NULL DEFAULT '' ,
        `pod_number` bigint NOT NULL DEFAULT 0 ,
        `container_number` bigint NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.TimeListLog`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `time_list_log` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `record_time` bigint NOT NULL DEFAULT 0
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.NodeDashboardMinute`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `node_dashboard_minute` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `node_name` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_cpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_gpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `memory_size` varchar(255) NOT NULL DEFAULT '' ,
        `pod_limit` varchar(255) NOT NULL DEFAULT '' ,
        `create_time` varchar(255) NOT NULL DEFAULT '' ,
        `ip` varchar(255) NOT NULL DEFAULT '' ,
        `cpu_usage` double precision NOT NULL DEFAULT 0 ,
        `mem_usage` double precision NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0 ,
        `storage_total` bigint NOT NULL DEFAULT 0 ,
        `storage_use` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.NodeDashboardHour`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `node_dashboard_hour` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `node_name` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_cpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_gpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `memory_size` varchar(255) NOT NULL DEFAULT '' ,
        `pod_limit` varchar(255) NOT NULL DEFAULT '' ,
        `create_time` varchar(255) NOT NULL DEFAULT '' ,
        `ip` varchar(255) NOT NULL DEFAULT '' ,
        `cpu_usage` double precision NOT NULL DEFAULT 0 ,
        `mem_usage` double precision NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0 ,
        `storage_total` bigint NOT NULL DEFAULT 0 ,
        `storage_use` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;

    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/collector/model/collect/dashboard.NodeDashboardDay`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `node_dashboard_day` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `node_name` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_cpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `numbers_gpu_core` varchar(255) NOT NULL DEFAULT '' ,
        `memory_size` varchar(255) NOT NULL DEFAULT '' ,
        `pod_limit` varchar(255) NOT NULL DEFAULT '' ,
        `create_time` varchar(255) NOT NULL DEFAULT '' ,
        `ip` varchar(255) NOT NULL DEFAULT '' ,
        `cpu_usage` double precision NOT NULL DEFAULT 0 ,
        `mem_usage` double precision NOT NULL DEFAULT 0 ,
        `time_list_id` bigint NOT NULL DEFAULT 0 ,
        `storage_total` bigint NOT NULL DEFAULT 0 ,
        `storage_use` bigint NOT NULL DEFAULT 0,
        KEY `idx_time_list_id` (`time_list_id`)
    ) ENGINE=InnoDB;

    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/common/model/image`
    -- --------------------------------------------------
    CREATE TABLE `board`.`image` (
        `id` INT AUTO_INCREMENT NOT NULL,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `comment` VARCHAR(255) NULL,
        `deleted` SMALLINT(1) NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

    CREATE TABLE `board`.`image_tag` (
        `id` INT AUTO_INCREMENT NOT NULL,
        `image_name` VARCHAR(255) NOT NULL DEFAULT '',
        `tag` VARCHAR(255) NOT NULL DEFAULT '',
        `deleted` SMALLINT(1) NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
    -- --------------------------------------------------
    --  Table Structure for `git/inspursoft/board/src/common/model/yaml/serviceconfig`
    -- --------------------------------------------------
    CREATE TABLE `board`.`service_config` (
        `id` INT AUTO_INCREMENT NOT NULL,
        `project_id` INT NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

    CREATE TABLE `board`.`service_config_image` (
        `service_id` INT NOT NULL,
        `image_tag_id` INT NOT NULL
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

    CREATE TABLE `board`.`service_status` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `project_id` INT NOT NULL,
        `project_name` VARCHAR(255) NOT NULL DEFAULT '',
        `comment` VARCHAR(255) NOT NULL DEFAULT '',
        `owner_id` INT NOT NULL,
        `owner_name` VARCHAR(255) DEFAULT NULL,
        `status` SMALLINT(1) NOT NULL,
		`type` SMALLINT(1) NOT NULL DEFAULT 0,
        `public` SMALLINT(1) NULL,
        `deleted` SMALLINT(1) NOT NULL DEFAULT 0,
        `creation_time` datetime DEFAULT NULL,
        `update_time` datetime DEFAULT NULL,
        `source` SMALLINT(1) NOT NULL,
        `source_id` INT NOT NULL DEFAULT 0,
        `service_yaml` TEXT,
        `deployment_yaml` TEXT,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;	

    CREATE TABLE `board`.`node_group` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `comment` VARCHAR(255) NOT NULL DEFAULT '',
        `owner_id` INT NOT NULL,
        `creation_time` datetime DEFAULT NULL,
        `update_time` datetime DEFAULT NULL,
        `deleted` SMALLINT(1) NOT NULL DEFAULT 0,
        `project_name` VARCHAR(255) NOT NULL DEFAULT '',
        `project_id` INT NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;	

    CREATE TABLE `board`.`config` (
        `name` varchar(50) NOT NULL DEFAULT '',
        `value` varchar(255) DEFAULT NULL,
        `comment` varchar(50) DEFAULT NULL,
        PRIMARY KEY (`name`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	
    CREATE TABLE `board`.`operation` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `creation_time` timestamp DEFAULT 0,
        `update_time` timestamp DEFAULT 0,
        `deleted` SMALLINT(1) NOT NULL DEFAULT 0,
        `project_name` VARCHAR(255) DEFAULT '',
        `project_id` INT DEFAULT 0,
        `user_name` VARCHAR(255) DEFAULT '',
        `user_id` INT DEFAULT 0,
        `object_type` VARCHAR(255) DEFAULT '',
        `object_name` VARCHAR(255) DEFAULT '',
        `action` VARCHAR(255) DEFAULT '',
        `status` VARCHAR(255) DEFAULT '',
        `path` VARCHAR(255) DEFAULT '',		
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;	

    CREATE TABLE `board`.`service_auto_scale` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `service_id` INT NOT NULL DEFAULT 0,
        `min_pod` INT NOT NULL DEFAULT 0,
        `max_pod` INT NOT NULL DEFAULT 0,
        `cpu_percent` INT NOT NULL DEFAULT 0,
        `status` INT NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;	
	
	CREATE TABLE `board`.`persistent_volume` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `type` INT NOT NULL DEFAULT 0,
        `state` INT NOT NULL DEFAULT 0,
		`capacity` VARCHAR(255) NOT NULL DEFAULT '',
		`accessmode` VARCHAR(255) NOT NULL DEFAULT '',
	    `class` VARCHAR(255) NOT NULL DEFAULT '',	
        `readonly` SMALLINT(1) NULL,		
	    `reclaim` VARCHAR(255) NOT NULL DEFAULT '',		
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;	
	
	CREATE TABLE `board`.`persistent_volume_option_nfs` (
        `id` INT NOT NULL,
		`path` VARCHAR(255) NOT NULL DEFAULT '',
		`server` VARCHAR(255) NOT NULL DEFAULT '',	
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	CREATE TABLE `board`.`persistent_volume_option_cephrbd` (
        `id` INT NOT NULL,
		`user` VARCHAR(255) NOT NULL DEFAULT '',
		`keyring` VARCHAR(255) NOT NULL DEFAULT '',	
		`pool` VARCHAR(255) NOT NULL DEFAULT '',
		`image` VARCHAR(255) NOT NULL DEFAULT '',	
		`fstype` VARCHAR(255) NOT NULL DEFAULT '',
		`secretname` VARCHAR(255) NOT NULL DEFAULT '',		
		`secretnamespace` VARCHAR(255) NOT NULL DEFAULT '',
		`monitors` VARCHAR(255) NOT NULL DEFAULT '',			
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;	
	
	CREATE TABLE `board`.`persistent_volume_claim_m` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `projectid` INT NOT NULL DEFAULT 0,
		`capacity` VARCHAR(255) NOT NULL DEFAULT '',
		`accessmode` VARCHAR(255) NOT NULL DEFAULT '',
	    `class` VARCHAR(255) NOT NULL DEFAULT '',		
	    `pvname` VARCHAR(255) NOT NULL DEFAULT '',		
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;		

    CREATE TABLE `board`.`helm_repository` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `url` VARCHAR(255) NOT NULL DEFAULT '',
        `type` INT NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;	

INSERT INTO `board`.`helm_repository`
 (`id`, `name`, `url`, `type`)
 VALUES
 (1, 'chartmuseum', 'http://chartmuseum:8080/', 1);

    CREATE TABLE `board`.`helm_release` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `project_id` INT NOT NULL,
        `project_name` VARCHAR(255) NOT NULL,
        `repository_id` INT NOT NULL,
        `repository` VARCHAR(255) NOT NULL,
        `workloads` TEXT,
        `owner_id` INT NOT NULL,
        `owner_name` VARCHAR(255) NOT NULL,
        `creation_time` datetime DEFAULT NULL,
        `update_time` datetime DEFAULT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

    CREATE TABLE `board`.`job_status` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL DEFAULT '',
        `project_id` INT NOT NULL,
        `project_name` VARCHAR(255) NOT NULL DEFAULT '',
        `comment` VARCHAR(255) NOT NULL DEFAULT '',
        `owner_id` INT NOT NULL,
        `owner_name` VARCHAR(255) DEFAULT NULL,
        `status` SMALLINT(1) NOT NULL,
        `deleted` SMALLINT(1) NOT NULL DEFAULT 0,
        `creation_time` datetime DEFAULT NULL,
        `update_time` datetime DEFAULT NULL,
        `source` SMALLINT(1) NOT NULL,
        `yaml` TEXT,
        PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;	
