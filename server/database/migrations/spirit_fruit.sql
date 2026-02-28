/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 90500 (9.5.0)
 Source Host           : localhost:3306
 Source Schema         : spirit_fruit

 Target Server Type    : MySQL
 Target Server Version : 90500 (9.5.0)
 File Encoding         : 65001

 Date: 01/03/2026 00:19:06
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admins
-- ----------------------------
DROP TABLE IF EXISTS `admins`;
CREATE TABLE `admins` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
  `mobile` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '手机号',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '密码',
  `email` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '邮箱',
  `authority_id` bigint unsigned DEFAULT '888' COMMENT '用户角色ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_admins_username` (`username`) USING BTREE,
  UNIQUE KEY `idx_admins_mobile` (`mobile`) USING BTREE,
  UNIQUE KEY `idx_admins_email` (`email`) USING BTREE,
  KEY `idx_admins_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of admins
-- ----------------------------
BEGIN;
INSERT INTO `admins` (`id`, `username`, `mobile`, `password`, `email`, `authority_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'admin', '18888888888', '$2a$14$xWzUdjlTLT8IWUI8cnGnTO3oJ4ELbFx0jIngLw49G2dJtqdgiaa8C', 'admin@example.com', 666, '2026-02-01 17:10:45', '2026-02-01 17:10:45', NULL);
INSERT INTO `admins` (`id`, `username`, `mobile`, `password`, `email`, `authority_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 'manager', '18666666666', '$2a$14$OBASAVoCLqAqZ0JdpSKCoO70ri5Z8ydJMZK0EllBABC4Li9at3X52', 'manager@example.com', 999, '2026-02-01 17:10:46', '2026-02-01 17:10:46', NULL);
COMMIT;

-- ----------------------------
-- Table structure for async_tasks
-- ----------------------------
DROP TABLE IF EXISTS `async_tasks`;
CREATE TABLE `async_tasks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `project_id` int DEFAULT '0',
  `rel_id` int DEFAULT '0',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `status` tinyint DEFAULT NULL,
  `payload` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `process` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `error_msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `started_at` datetime DEFAULT NULL,
  `finished_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of async_tasks
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for characters
-- ----------------------------
DROP TABLE IF EXISTS `characters`;
CREATE TABLE `characters` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `project_id` bigint unsigned NOT NULL COMMENT '所属项目ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '角色名',
  `role_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT 'main' COMMENT '角色类型: main/supporting/minor',
  `gender` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '性别(需从appearance解析或留空)',
  `age_group` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '年龄段',
  `personality` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '性格描述',
  `appearance_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '外貌长文本描述(原appearance)',
  `visual_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT 'AI绘画专用Prompt(从appearance提取)',
  `avatar_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '头像/定妆照',
  `voice_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'TTS音色ID',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_characters_project_id` (`project_id`) USING BTREE,
  KEY `idx_characters_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='角色';

-- ----------------------------
-- Records of characters
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for projects
-- ----------------------------
DROP TABLE IF EXISTS `projects`;
CREATE TABLE `projects` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `admin_id` bigint unsigned NOT NULL DEFAULT '1' COMMENT '归属用户ID(默认1)',
  `serial_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '业务流水号',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '项目名称/短剧标题',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '项目简介',
  `style` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '画面风格',
  `genre` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '题材',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态 0-Draft 1-Generating 2-Completed',
  `image` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '封面图',
  `total_duration` int DEFAULT '0' COMMENT '总时长(秒)',
  `settings` json DEFAULT NULL COMMENT '生成配置快照',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_projects_user_id` (`admin_id`) USING BTREE,
  KEY `idx_projects_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='短剧项目';

-- ----------------------------
-- Records of projects
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for props
-- ----------------------------
DROP TABLE IF EXISTS `props`;
CREATE TABLE `props` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint unsigned DEFAULT NULL COMMENT '所属项目ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '道具名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '道具类型',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '道具描述',
  `image_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'AI绘画提示词',
  `image_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '道具图片URL',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_props_project_id` (`project_id`) USING BTREE,
  KEY `idx_props_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='道具表';

-- ----------------------------
-- Records of props
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for scenes
-- ----------------------------
DROP TABLE IF EXISTS `scenes`;
CREATE TABLE `scenes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `project_id` int DEFAULT '0' COMMENT '所属项目ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '场景名称',
  `location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '地点',
  `time` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '时间',
  `atmosphere` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '氛围描述',
  `visual_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT 'AI绘画Prompt',
  `image_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '图片路径',
  `status` tinyint DEFAULT NULL COMMENT '状态 1-待生成 2-生成中 3-已完成',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of scenes
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for scripts
-- ----------------------------
DROP TABLE IF EXISTS `scripts`;
CREATE TABLE `scripts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `project_id` bigint unsigned NOT NULL COMMENT '所属项目ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '分集标题',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '剧本正文',
  `outline` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '大纲/简介',
  `episode_no` int DEFAULT '1' COMMENT '第几集',
  `is_locked` tinyint(1) DEFAULT '0' COMMENT '是否定稿 0-否 1-是',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_scripts_project_id` (`project_id`) USING BTREE,
  KEY `idx_scripts_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='剧本';

-- ----------------------------
-- Records of scripts
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for shot_characters
-- ----------------------------
DROP TABLE IF EXISTS `shot_characters`;
CREATE TABLE `shot_characters` (
  `shot_id` int NOT NULL,
  `character_id` int NOT NULL,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`shot_id`,`character_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of shot_characters
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for shot_frame_images
-- ----------------------------
DROP TABLE IF EXISTS `shot_frame_images`;
CREATE TABLE `shot_frame_images` (
  `id` int NOT NULL AUTO_INCREMENT,
  `project_id` int DEFAULT NULL,
  `shot_id` int DEFAULT NULL,
  `frame_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `image_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `image_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of shot_frame_images
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for shot_frame_prompts
-- ----------------------------
DROP TABLE IF EXISTS `shot_frame_prompts`;
CREATE TABLE `shot_frame_prompts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `shot_id` int DEFAULT NULL,
  `frame_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `layout` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of shot_frame_prompts
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for shot_props
-- ----------------------------
DROP TABLE IF EXISTS `shot_props`;
CREATE TABLE `shot_props` (
  `shot_id` int NOT NULL,
  `props_id` int NOT NULL,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`shot_id`,`props_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of shot_props
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for shots
-- ----------------------------
DROP TABLE IF EXISTS `shots`;
CREATE TABLE `shots` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `project_id` bigint unsigned NOT NULL COMMENT '所属项目ID',
  `script_id` bigint unsigned NOT NULL COMMENT '所属剧本/分集ID',
  `scene_id` bigint DEFAULT NULL COMMENT '场景id',
  `sequence_no` int NOT NULL DEFAULT '0' COMMENT '镜头序号',
  `shot_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '景别: 全景/特写/中景',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '标题',
  `action` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '人物动作描述',
  `time` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '具体时间描述',
  `location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '具体地点描述',
  `camera_movement` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '运镜: 推/拉/摇/移',
  `angle` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '视角: 俯视/平视',
  `dialogue` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '台词/旁白',
  `visual_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '画面描述',
  `atmosphere` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '氛围/环境描述',
  `image_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '绘画Prompt',
  `video_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '视频生成Prompt',
  `audio_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '音效/BGM提示词',
  `image_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '分镜图',
  `video_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '最终视频片段',
  `audio_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '配音/音效',
  `duration_ms` int DEFAULT '3000' COMMENT '时长(毫秒, 原duration*1000)',
  `status` tinyint DEFAULT '0' COMMENT '状态 0-Pending 1-Done 2-Fail',
  `created_at` datetime DEFAULT NULL COMMENT '添加时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_shots_project_id` (`project_id`) USING BTREE,
  KEY `idx_shots_script_id` (`script_id`) USING BTREE,
  KEY `idx_shots_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='镜头表';

-- ----------------------------
-- Records of shots
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for sys_base_menus
-- ----------------------------
DROP TABLE IF EXISTS `sys_base_menus`;
CREATE TABLE `sys_base_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父菜单ID',
  `path` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '路由path',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '路由name',
  `hidden` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否在列表隐藏(0=显示,1=隐藏)',
  `component` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '对应前端文件路径',
  `sort` bigint NOT NULL DEFAULT '0' COMMENT '排序标记',
  `title` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '菜单名',
  `icon` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '菜单图标',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_sys_base_menus_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of sys_base_menus
-- ----------------------------
BEGIN;
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, NULL, '/dashboard', 'Dashboard', 1, 'Layout', 1, '仪表盘', 'dashboard', '2026-02-01 17:10:46', '2026-02-01 18:19:53', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 1, '/dashboard/base', 'DashboardStatistics', 0, '/dashboard/base/index.vue', 2, '统计报表', NULL, '2026-02-01 17:10:46', '2026-02-01 17:10:46', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, NULL, '/admin/characters', 'AdminCharactersModule', 1, 'Layout', 10, '角色管理', 'app', '2026-02-01 17:10:46', '2026-02-02 16:40:00', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 3, '/admin/characters/list', 'AdminCharactersList', 0, '/characters/index.vue', 1, '角色列表', NULL, '2026-02-01 17:10:46', '2026-02-01 17:10:46', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, NULL, '/admin/projects', 'AdminProjectsModule', 0, 'Layout', 11, '短剧项目', 'folder-add', '2026-02-01 17:10:46', '2026-02-02 16:31:45', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 5, '/admin/projects/list', 'AdminProjectsList', 1, '/projects/index.vue', 0, '短剧项目', NULL, '2026-02-01 17:10:46', '2026-02-06 09:58:10', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, NULL, '/admin/scripts', 'AdminScriptsModule', 1, 'Layout', 12, '剧本管理', 'app', '2026-02-01 17:10:46', '2026-02-02 16:40:12', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 7, '/admin/scripts/list', 'AdminScriptsList', 0, '/scripts/index.vue', 1, '剧本列表', NULL, '2026-02-01 17:10:46', '2026-02-01 17:10:46', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, NULL, '/admin/shots', 'AdminShotsModule', 1, 'Layout', 13, '镜头管理', 'app', '2026-02-01 17:10:46', '2026-02-02 16:40:24', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 9, '/admin/shots/list', 'AdminShotsList', 0, '/shots/index.vue', 1, '镜头列表', NULL, '2026-02-01 17:10:46', '2026-02-02 16:32:00', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, NULL, '/admin/admins', 'AdminAdminsModule', 1, 'Layout', 14, '系统管理员', 'user-setting', '2026-02-01 17:10:46', '2026-02-02 16:32:04', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 11, '/admin/admins/list', 'AdminAdminsList', 0, '/admins/index.vue', 1, '系统管理员', NULL, '2026-02-01 17:10:46', '2026-02-02 16:32:10', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, NULL, '/admin/system', 'AdminSystemModule', 1, 'Layout', 900, '系统管理', 'setting', '2026-02-01 17:10:46', '2026-02-02 16:39:50', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 13, '/admin/system/menus', 'AdminSysBaseMenuList', 0, '/sys_base_menus/index.vue', 1, '菜单管理', NULL, '2026-02-01 17:10:46', '2026-02-01 17:10:46', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, NULL, '/user', 'UserCenter', 1, 'Layout', 999, '个人中心', 'user-circle', '2026-02-01 17:10:46', '2026-02-01 17:10:46', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 15, '/user/index', 'UserProfile', 1, '/user/index', 1, '个人信息', NULL, '2026-02-01 17:10:46', '2026-02-01 17:10:46', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, 5, 'detail/:id', 'ProjectDetail', 1, '/projects/detail.vue', 20, '项目工作台', NULL, '2026-02-02 17:23:05', '2026-02-06 09:58:18', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (18, 5, 'chapter/:id/:episodeNumber', 'ProjectChapterCreate', 1, '/projects/createChapter.vue', 0, '章节创作', NULL, '2026-02-02 17:47:44', '2026-02-02 17:47:44', NULL);
INSERT INTO `sys_base_menus` (`id`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `created_at`, `updated_at`, `deleted_at`) VALUES (19, 5, 'editor/:dramaId/:episodeNumber', 'ScriptEditor', 1, '/projects/scriptEditor.vue', 0, '视频创作', NULL, '2026-02-08 16:54:37', '2026-02-08 16:54:54', NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
