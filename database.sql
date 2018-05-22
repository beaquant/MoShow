-- MySQL dump 10.13  Distrib 5.7.17, for Win64 (x86_64)
--
-- Host: 47.96.177.91    Database: MoShow
-- ------------------------------------------------------
-- Server version	5.7.21-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `admin`
--

DROP TABLE IF EXISTS `admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `wechat_id` varchar(45) NOT NULL,
  `alias` varchar(45) NOT NULL DEFAULT '大佬',
  `last_login_info` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `wechat_id_UNIQUE` (`wechat_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `balance_chg`
--

DROP TABLE IF EXISTS `balance_chg`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `balance_chg` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `from_user_id` bigint(20) unsigned DEFAULT NULL,
  `chg_type` int(10) NOT NULL DEFAULT '0',
  `chg_info` json DEFAULT NULL,
  `amount` int(11) NOT NULL DEFAULT '0',
  `time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4783 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `config`
--

DROP TABLE IF EXISTS `config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `config` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `conf_key` varchar(45) NOT NULL,
  `val` json NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dial`
--

DROP TABLE IF EXISTS `dial`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dial` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `from_user_id` bigint(20) unsigned NOT NULL,
  `to_user_id` bigint(20) unsigned NOT NULL,
  `duration` int(10) unsigned NOT NULL DEFAULT '0',
  `create_at` bigint(20) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `tag` json DEFAULT NULL,
  `clearing` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1397 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `feedback`
--

DROP TABLE IF EXISTS `feedback`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `feedback` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `content` json DEFAULT NULL,
  `type` int(11) NOT NULL DEFAULT '0',
  `time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `guest`
--

DROP TABLE IF EXISTS `guest`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `guest` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `guest_id` bigint(20) unsigned NOT NULL,
  `time` bigint(20) NOT NULL DEFAULT '0',
  `count` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=780 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order`
--

DROP TABLE IF EXISTS `order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `coin_count` int(11) NOT NULL DEFAULT '0',
  `success` tinyint(4) NOT NULL DEFAULT '0',
  `pay_type` int(11) NOT NULL DEFAULT '0',
  `create_at` bigint(20) NOT NULL DEFAULT '0',
  `pay_time` bigint(20) DEFAULT '0',
  `pay_info` json DEFAULT NULL,
  `product_info` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=343 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `profile_chg`
--

DROP TABLE IF EXISTS `profile_chg`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `profile_chg` (
  `id` int(11) NOT NULL,
  `cover_pic` varchar(256) NOT NULL,
  `cover_pic_check` int(11) NOT NULL,
  `video` varchar(256) NOT NULL,
  `video_check` int(11) NOT NULL,
  `update_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `subscribe`
--

DROP TABLE IF EXISTS `subscribe`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `subscribe` (
  `id` int(11) NOT NULL,
  `following` json DEFAULT NULL COMMENT '正在关注',
  `follower` json DEFAULT NULL COMMENT '关注者',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary view structure for view `time_line`
--

DROP TABLE IF EXISTS `time_line`;
/*!50001 DROP VIEW IF EXISTS `time_line`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `time_line` AS SELECT 
 1 AS `id`,
 1 AS `alias`,
 1 AS `gender`,
 1 AS `cover`,
 1 AS `description`,
 1 AS `birthday`,
 1 AS `location`,
 1 AS `price`,
 1 AS `user_type`,
 1 AS `user_status`,
 1 AS `online_status`,
 1 AS `anchor_auth_status`,
 1 AS `dial_accept`,
 1 AS `dial_deny`,
 1 AS `update_at`,
 1 AS `dial_duration`,
 1 AS `create_at`,
 1 AS `recent_duration`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `user_extra`
--

DROP TABLE IF EXISTS `user_extra`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_extra` (
  `id` int(11) NOT NULL,
  `gift_his` json DEFAULT NULL COMMENT '收到的所有礼物',
  `income_his` int(11) NOT NULL DEFAULT '0' COMMENT '历史总收益',
  `invite_count` int(11) NOT NULL DEFAULT '0' COMMENT '邀请的总人数',
  `invite_income_his` int(11) NOT NULL DEFAULT '0',
  `balance_his` int(11) NOT NULL DEFAULT '0' COMMENT '历史总充值',
  `video_view_pay` json DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_profile`
--

DROP TABLE IF EXISTS `user_profile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_profile` (
  `id` bigint(20) unsigned NOT NULL,
  `alias` varchar(64) NOT NULL DEFAULT 'MoShow',
  `gender` int(11) unsigned NOT NULL DEFAULT '0',
  `cover` json DEFAULT NULL,
  `description` varchar(256) DEFAULT NULL,
  `birthday` bigint(20) NOT NULL DEFAULT '725817600',
  `location` varchar(45) NOT NULL DEFAULT '浙江杭州',
  `balance` int(11) unsigned NOT NULL DEFAULT '0',
  `income` int(11) unsigned NOT NULL DEFAULT '0',
  `price` int(11) unsigned NOT NULL DEFAULT '0',
  `user_type` int(11) NOT NULL DEFAULT '0',
  `im_token` varchar(64) DEFAULT NULL,
  `user_status` int(11) NOT NULL DEFAULT '0',
  `online_status` int(11) NOT NULL DEFAULT '0',
  `anchor_auth_status` int(11) NOT NULL DEFAULT '0',
  `dial_accept` int(11) NOT NULL DEFAULT '0',
  `dial_deny` int(11) NOT NULL DEFAULT '0',
  `update_at` bigint(20) NOT NULL DEFAULT '0',
  `dial_duration` int(11) NOT NULL DEFAULT '0',
  `alipay_acct` json DEFAULT NULL,
  `tag` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(128) DEFAULT NULL,
  `wechat_id` varchar(128) DEFAULT NULL,
  `acct_type` int(11) NOT NULL DEFAULT '0' COMMENT '0手机\n1微信',
  `acct_status` int(11) NOT NULL DEFAULT '0' COMMENT '0正常\n1注销\n2屏蔽',
  `create_at` bigint(20) NOT NULL DEFAULT '0',
  `invited_by` bigint(20) unsigned DEFAULT NULL,
  `invited_award` bigint(20) DEFAULT '0',
  `last_login_info` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=169221 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `withdraw`
--

DROP TABLE IF EXISTS `withdraw`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `withdraw` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `amount` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `create_at` bigint(20) NOT NULL,
  `tag` json DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Final view structure for view `time_line`
--

/*!50001 DROP VIEW IF EXISTS `time_line`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `time_line` AS (select `user_profile`.`id` AS `id`,`user_profile`.`alias` AS `alias`,`user_profile`.`gender` AS `gender`,`user_profile`.`cover` AS `cover`,`user_profile`.`description` AS `description`,`user_profile`.`birthday` AS `birthday`,`user_profile`.`location` AS `location`,`user_profile`.`price` AS `price`,`user_profile`.`user_type` AS `user_type`,`user_profile`.`user_status` AS `user_status`,`user_profile`.`online_status` AS `online_status`,`user_profile`.`anchor_auth_status` AS `anchor_auth_status`,`user_profile`.`dial_accept` AS `dial_accept`,`user_profile`.`dial_deny` AS `dial_deny`,`user_profile`.`update_at` AS `update_at`,`user_profile`.`dial_duration` AS `dial_duration`,`users`.`create_at` AS `create_at`,sum(`dial`.`duration`) AS `recent_duration` from ((`user_profile` left join `users` on((`user_profile`.`id` = `users`.`id`))) left join `dial` on((((`user_profile`.`id` = `dial`.`from_user_id`) or (`user_profile`.`id` = `dial`.`to_user_id`)) and (`dial`.`create_at` > unix_timestamp((now() - interval 3 day)))))) where ((`user_profile`.`id` <> 1) and (`users`.`acct_status` <> 1) and (`user_profile`.`user_status` <> 2)) group by `user_profile`.`id`) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-05-22 10:33:17
