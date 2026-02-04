-- MySQL dump 10.13  Distrib 9.5.0, for macos26.0 (arm64)
--
-- Host: 127.0.0.1    Database: spirit_fruit
-- ------------------------------------------------------
-- Server version	9.5.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '034c39fc-b0d7-11f0-907a-bf226fc6928a:1-124733';

--
-- Table structure for table `admins`
--

DROP TABLE IF EXISTS `admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admins` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(120) COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
  `mobile` char(11) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '手机号',
  `password` char(64) COLLATE utf8mb4_bin NOT NULL COMMENT '密码',
  `email` varchar(80) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '邮箱',
  `authority_id` bigint unsigned DEFAULT '888' COMMENT '用户角色ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  UNIQUE KEY `idx_admins_mobile` (`mobile`),
  UNIQUE KEY `idx_admins_email` (`email`),
  KEY `idx_admins_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admins`
--

LOCK TABLES `admins` WRITE;
/*!40000 ALTER TABLE `admins` DISABLE KEYS */;
INSERT INTO `admins` VALUES (1,'admin','18888888888','$2a$14$xWzUdjlTLT8IWUI8cnGnTO3oJ4ELbFx0jIngLw49G2dJtqdgiaa8C','admin@example.com',666,'2026-02-01 17:10:45','2026-02-01 17:10:45',NULL),(2,'manager','18666666666','$2a$14$OBASAVoCLqAqZ0JdpSKCoO70ri5Z8ydJMZK0EllBABC4Li9at3X52','manager@example.com',999,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL);
/*!40000 ALTER TABLE `admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `characters`
--

DROP TABLE IF EXISTS `characters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `characters` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `project_id` bigint unsigned NOT NULL COMMENT '所属项目ID',
  `name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '角色名',
  `role_type` varchar(50) COLLATE utf8mb4_bin DEFAULT 'main' COMMENT '角色类型: main/supporting/minor',
  `gender` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '性别(需从appearance解析或留空)',
  `age_group` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '年龄段',
  `personality` text COLLATE utf8mb4_bin COMMENT '性格描述',
  `appearance_desc` text COLLATE utf8mb4_bin COMMENT '外貌长文本描述(原appearance)',
  `visual_prompt` text COLLATE utf8mb4_bin COMMENT 'AI绘画专用Prompt(从appearance提取)',
  `avatar_url` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '头像/定妆照',
  `voice_id` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'TTS音色ID',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_characters_project_id` (`project_id`),
  KEY `idx_characters_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='角色';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `characters`
--

LOCK TABLES `characters` WRITE;
/*!40000 ALTER TABLE `characters` DISABLE KEYS */;
/*!40000 ALTER TABLE `characters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `projects`
--

DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `projects` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `admin_id` bigint unsigned NOT NULL DEFAULT '1' COMMENT '归属用户ID(默认1)',
  `serial_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '业务流水号',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '项目名称/短剧标题',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '项目简介',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态 0-Draft 1-Generating 2-Completed',
  `image` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '封面图',
  `total_duration` int DEFAULT '0' COMMENT '总时长(秒)',
  `settings` json DEFAULT NULL COMMENT '生成配置快照',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_projects_user_id` (`admin_id`),
  KEY `idx_projects_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='短剧项目';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projects`
--

LOCK TABLES `projects` WRITE;
/*!40000 ALTER TABLE `projects` DISABLE KEYS */;
/*!40000 ALTER TABLE `projects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scripts`
--

DROP TABLE IF EXISTS `scripts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `scripts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `project_id` bigint unsigned NOT NULL COMMENT '所属项目ID',
  `title` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '分集标题',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '剧本正文',
  `outline` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '大纲/简介',
  `episode_no` int DEFAULT '1' COMMENT '第几集',
  `is_locked` tinyint(1) DEFAULT '0' COMMENT '是否定稿 0-否 1-是',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_scripts_project_id` (`project_id`),
  KEY `idx_scripts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='剧本';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scripts`
--

LOCK TABLES `scripts` WRITE;
/*!40000 ALTER TABLE `scripts` DISABLE KEYS */;
/*!40000 ALTER TABLE `scripts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `shots`
--

DROP TABLE IF EXISTS `shots`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `shots` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `project_id` bigint unsigned NOT NULL COMMENT '所属项目ID',
  `script_id` bigint unsigned NOT NULL COMMENT '所属剧本/分集ID',
  `sequence_no` int NOT NULL DEFAULT '0' COMMENT '镜头序号',
  `shot_type` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '景别: 全景/特写/中景',
  `camera_movement` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '运镜: 推/拉/摇/移',
  `angle` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '视角: 俯视/平视',
  `dialogue` text COLLATE utf8mb4_bin COMMENT '台词/旁白',
  `visual_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '画面描述',
  `atmosphere` text COLLATE utf8mb4_bin COMMENT '氛围/环境描述',
  `image_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '绘画Prompt',
  `video_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '视频生成Prompt',
  `audio_prompt` text COLLATE utf8mb4_bin COMMENT '音效/BGM提示词',
  `image_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '分镜图',
  `video_url` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '最终视频片段',
  `audio_url` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '配音/音效',
  `duration_ms` int DEFAULT '3000' COMMENT '时长(毫秒, 原duration*1000)',
  `status` tinyint DEFAULT '0' COMMENT '状态 0-Pending 1-Done 2-Fail',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_shots_project_id` (`project_id`),
  KEY `idx_shots_script_id` (`script_id`),
  KEY `idx_shots_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='镜头表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `shots`
--

LOCK TABLES `shots` WRITE;
/*!40000 ALTER TABLE `shots` DISABLE KEYS */;
/*!40000 ALTER TABLE `shots` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_base_menus`
--

DROP TABLE IF EXISTS `sys_base_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_base_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父菜单ID',
  `path` varchar(120) COLLATE utf8mb4_bin NOT NULL COMMENT '路由path',
  `name` varchar(120) COLLATE utf8mb4_bin NOT NULL COMMENT '路由name',
  `hidden` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否在列表隐藏(0=显示,1=隐藏)',
  `component` varchar(120) COLLATE utf8mb4_bin NOT NULL COMMENT '对应前端文件路径',
  `sort` bigint NOT NULL DEFAULT '0' COMMENT '排序标记',
  `title` varchar(120) COLLATE utf8mb4_bin NOT NULL COMMENT '菜单名',
  `icon` varchar(120) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '菜单图标',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_sys_base_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_base_menus`
--

LOCK TABLES `sys_base_menus` WRITE;
/*!40000 ALTER TABLE `sys_base_menus` DISABLE KEYS */;
INSERT INTO `sys_base_menus` VALUES (1,NULL,'/dashboard','Dashboard',0,'Layout',1,'仪表盘','dashboard','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(2,1,'/dashboard/base','DashboardStatistics',0,'/dashboard/base/index.vue',2,'统计报表',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(3,NULL,'/admin/characters','AdminCharactersModule',0,'Layout',10,'角色管理','app','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(4,3,'/admin/characters/list','AdminCharactersList',0,'/characters/index.vue',1,'角色列表',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(5,NULL,'/admin/projects','AdminProjectsModule',0,'Layout',11,'短剧项目管理','folder-add','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(6,5,'/admin/projects/list','AdminProjectsList',0,'/projects/index.vue',1,'短剧项目列表',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(7,NULL,'/admin/scripts','AdminScriptsModule',0,'Layout',12,'剧本管理','app','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(8,7,'/admin/scripts/list','AdminScriptsList',0,'/scripts/index.vue',1,'剧本列表',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(9,NULL,'/admin/shots','AdminShotsModule',0,'Layout',13,'镜头表管理','app','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(10,9,'/admin/shots/list','AdminShotsList',0,'/shots/index.vue',1,'镜头表列表',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(11,NULL,'/admin/admins','AdminAdminsModule',0,'Layout',14,'系统管理员管理','user-setting','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(12,11,'/admin/admins/list','AdminAdminsList',0,'/admins/index.vue',1,'系统管理员列表',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(13,NULL,'/admin/system','AdminSystemModule',0,'Layout',900,'系统管理','setting','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(14,13,'/admin/system/menus','AdminSysBaseMenuList',0,'/sys_base_menus/index.vue',1,'菜单管理',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(15,NULL,'/user','UserCenter',1,'Layout',999,'个人中心','user-circle','2026-02-01 17:10:46','2026-02-01 17:10:46',NULL),(16,15,'/user/index','UserProfile',1,'/user/index',1,'个人信息',NULL,'2026-02-01 17:10:46','2026-02-01 17:10:46',NULL);
/*!40000 ALTER TABLE `sys_base_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'spirit_fruit'
--
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-02-01 17:10:46
