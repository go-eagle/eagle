-- 创建数据库
CREATE DATABASE IF NOT EXISTS `snake`;
USE `snake`;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
     `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `username` varchar(255) NOT NULL DEFAULT '',
     `password` varchar(60) NOT NULL DEFAULT '',
     `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
     `phone` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '手机号',
     `sex` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '性别 0:未知 1:男 2:女',
     `deleted_at` timestamp NULL DEFAULT NULL,
     `created_at` timestamp NULL DEFAULT NULL,
     `updated_at` timestamp NULL DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `uniq_username` (`username`),
     UNIQUE KEY `uniq_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

INSERT INTO `users` (`id`, `username`, `password`, `avatar`, `phone`, `sex`, `created_at`, `updated_at`,`deleted_at`)
VALUES
(null, 'test-name', '6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9', '/uploads/avatar.jpg', 13810002001, 1, '2020-02-09 10:23:33', '2020-05-09 10:23:33', NULL);

