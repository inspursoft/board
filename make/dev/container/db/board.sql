drop database if exists board;
create database board charset = utf8;

use board;

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


INSERT INTO `board`.`user` (`username`, `password`, `email`, `realname`, `comment`, `creation_time`)
  VALUES ('admin', 'Board12345', 'admin@inspur.com', 'admin', 'admin user', now());

