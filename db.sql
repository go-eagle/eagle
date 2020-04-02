-- 创建数据库
CREATE DATABASE IF NOT EXISTS `snake`;
USE `snake`;

--
-- Table structure for table `tb_users`
--

DROP TABLE IF EXISTS `tb_users`;

CREATE TABLE `tb_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(30) NOT NULL,
  `password` varchar(32) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_username` (`username`),
  KEY `idx_tb_users_deletedAt` (`deleted_at`)
) ENGINE=Innodb DEFAULT CHARSET=utf8;
