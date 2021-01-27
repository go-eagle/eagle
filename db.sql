-- 创建数据库
CREATE DATABASE IF NOT EXISTS `snake`;
USE `snake`;

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table user_fans
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_fans`;

CREATE TABLE `user_fans` (
     `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
     `follower_uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝的uid',
     `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态 1:已关注 0:取消关注',
     `created_at` datetime DEFAULT NULL,
     `updated_at` datetime DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `idx_uid_fid` (`user_id`,`follower_uid`),
     KEY `idx_status_uid` (`status`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户粉丝表';

LOCK TABLES `user_fans` WRITE;
/*!40000 ALTER TABLE `user_fans` DISABLE KEYS */;

INSERT INTO `user_fans` (`id`, `user_id`, `follower_uid`, `status`, `created_at`, `updated_at`)
VALUES
(1,2,1,1,'2020-05-23 00:12:30',NULL),
(2,4,1,1,'2020-05-23 00:23:10',NULL),
(3,12,1,1,'2020-05-23 00:25:48','2020-05-23 00:27:03'),
(5,13,1,1,'2020-05-29 12:50:54',NULL);

/*!40000 ALTER TABLE `user_fans` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_follow
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_follow`;

CREATE TABLE `user_follow` (
   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
   `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发起关注的人',
   `followed_uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '被关注用户的uid',
   `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '关注状态 1:已关注 0:取消关注',
   `created_at` datetime DEFAULT NULL,
   `updated_at` datetime DEFAULT NULL,
   PRIMARY KEY (`id`),
   UNIQUE KEY `uniq_uid_fuid` (`user_id`,`followed_uid`),
   KEY `idx_status_uid` (`status`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户关注表';

LOCK TABLES `user_follow` WRITE;
/*!40000 ALTER TABLE `user_follow` DISABLE KEYS */;

INSERT INTO `user_follow` (`id`, `user_id`, `followed_uid`, `status`, `created_at`, `updated_at`)
VALUES
(1,1,2,1,'2020-05-23 00:12:30',NULL),
(2,1,4,1,'2020-05-23 00:23:10',NULL),
(3,1,12,1,'2020-05-23 00:25:48','2020-05-23 00:27:03'),
(5,1,13,1,'2020-05-29 12:50:54',NULL);

/*!40000 ALTER TABLE `user_follow` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_stat
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_stat`;

CREATE TABLE `user_stat` (
 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
 `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
 `follow_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '关注数',
 `follower_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝数',
 `status` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '状态  1:正常',
 `created_at` timestamp NULL DEFAULT NULL,
 `updated_at` timestamp NULL DEFAULT NULL,
 PRIMARY KEY (`id`),
 UNIQUE KEY `uniq_uid` (`user_id`),
 KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户统计表';

LOCK TABLES `user_stat` WRITE;
/*!40000 ALTER TABLE `user_stat` DISABLE KEYS */;

INSERT INTO `user_stat` (`id`, `user_id`, `follow_count`, `follower_count`, `status`, `created_at`, `updated_at`)
VALUES
(1,1,3,0,1,'2020-05-23 00:12:30','2020-05-29 12:50:54'),
(2,2,0,0,1,'2020-05-23 00:12:30','2020-05-23 00:20:09'),
(8,4,0,1,1,'2020-05-23 00:23:10',NULL),
(10,12,0,1,1,'2020-05-23 00:25:48','2020-05-23 00:27:03'),
(16,13,0,1,1,'2020-05-29 12:50:54',NULL);

/*!40000 ALTER TABLE `user_stat` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_base`;

CREATE TABLE `user_base` (
     `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `username` varchar(255) NOT NULL DEFAULT '',
     `password` varchar(60) NOT NULL DEFAULT '',
     `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
     `phone` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '手机号',
     `email` varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱',
     `sex` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '性别 0:未知 1:男 2:女',
     `deleted_at` timestamp NULL DEFAULT NULL,
     `created_at` timestamp NULL DEFAULT NULL,
     `updated_at` timestamp NULL DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `uniq_username` (`username`),
     UNIQUE KEY `uniq_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

LOCK TABLES `user_base` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `user_base` (`id`, `username`, `password`, `avatar`, `phone`, `email`, `sex`, `deleted_at`, `created_at`, `updated_at`)
VALUES
(1,'test-name','$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym','/uploads/avatar.jpg',13010102020,'123@cc.com',1,NULL,'2020-02-09 10:23:33','2020-05-09 10:23:33'),
(2,'admin','$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym','',1,'1234@cc.com',0,NULL,'2020-05-20 22:42:18','2020-05-20 22:42:18'),
(4,'admin2','$2a$10$Dps9oN3Oe3ZDMACih3DCGeTvR.jW/I8WD1NqapCJ6Vq3PzjnusI9i','',0,'12345@cc.com',0,NULL,'2020-05-20 22:43:21','2020-05-20 22:43:21'),
(12,'user001','123456','',13810002000,'',0,NULL,'0000-00-00 00:00:00','0000-00-00 00:00:00'),
(13,'user002','123456','',13810002001,'',0,NULL,'0000-00-00 00:00:00','0000-00-00 00:00:00');

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
