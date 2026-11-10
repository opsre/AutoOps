/*
 Navicat Premium Dump SQL

 Source Server         : 172.16.7.137
 Source Server Type    : MySQL
 Source Server Version : 80045 (8.0.45)
 Source Host           : 172.16.7.137:3306
 Source Schema         : autoops

 Target Server Type    : MySQL
 Target Server Version : 80045 (8.0.45)
 File Encoding         : 65001

 Date: 10/11/2026 23:14:27
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for aiops_agent_feedback
-- ----------------------------
DROP TABLE IF EXISTS `aiops_agent_feedback`;
CREATE TABLE `aiops_agent_feedback`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `run_id` bigint UNSIGNED NOT NULL,
  `session_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `rating` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `intent` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `username` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `raw_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_aiops_agent_feedback_run_id`(`run_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_feedback_session_id`(`session_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_feedback_rating`(`rating` ASC) USING BTREE,
  INDEX `idx_aiops_agent_feedback_intent`(`intent` ASC) USING BTREE,
  INDEX `idx_aiops_agent_feedback_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_feedback_username`(`username` ASC) USING BTREE,
  INDEX `idx_aiops_agent_feedback_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_agent_feedback
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_agent_message
-- ----------------------------
DROP TABLE IF EXISTS `aiops_agent_message`;
CREATE TABLE `aiops_agent_message`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `run_id` bigint UNSIGNED NULL DEFAULT NULL,
  `role` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `intent` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `entities_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `tool_plan_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `timeline_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `diagnosis_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `raw_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_aiops_agent_message_session_id`(`session_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_message_run_id`(`run_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_message_role`(`role` ASC) USING BTREE,
  INDEX `idx_aiops_agent_message_intent`(`intent` ASC) USING BTREE,
  INDEX `idx_aiops_agent_message_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_agent_message
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_agent_run
-- ----------------------------
DROP TABLE IF EXISTS `aiops_agent_run`;
CREATE TABLE `aiops_agent_run`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `run_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `provider_id` bigint UNSIGNED NULL DEFAULT NULL,
  `provider_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `provider_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `model_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `summary` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `steps_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `result_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `used_llm` tinyint(1) NULL DEFAULT NULL,
  `error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `started_at` datetime NOT NULL,
  `finished_at` datetime NULL DEFAULT NULL,
  `duration_ms` bigint NULL DEFAULT NULL,
  `generated_at` datetime NULL DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `alert_id` bigint UNSIGNED NULL DEFAULT NULL,
  `fingerprint` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_aiops_agent_run_run_type`(`run_type` ASC) USING BTREE,
  INDEX `idx_aiops_agent_run_status`(`status` ASC) USING BTREE,
  INDEX `idx_aiops_agent_run_provider_id`(`provider_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_run_started_at`(`started_at` ASC) USING BTREE,
  INDEX `idx_aiops_agent_run_generated_at`(`generated_at` ASC) USING BTREE,
  INDEX `idx_aiops_agent_run_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_agent_run
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_agent_session
-- ----------------------------
DROP TABLE IF EXISTS `aiops_agent_session`;
CREATE TABLE `aiops_agent_session`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `channel` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `username` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `last_run_id` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_aiops_agent_session_session_id`(`session_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_session_channel`(`channel` ASC) USING BTREE,
  INDEX `idx_aiops_agent_session_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_session_username`(`username` ASC) USING BTREE,
  INDEX `idx_aiops_agent_session_status`(`status` ASC) USING BTREE,
  INDEX `idx_aiops_agent_session_last_run_id`(`last_run_id` ASC) USING BTREE,
  INDEX `idx_aiops_agent_session_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_agent_session
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_evidence
-- ----------------------------
DROP TABLE IF EXISTS `aiops_evidence`;
CREATE TABLE `aiops_evidence`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `run_id` bigint UNSIGNED NOT NULL,
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `reference` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `raw_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_aiops_evidence_run_id`(`run_id` ASC) USING BTREE,
  INDEX `idx_aiops_evidence_type`(`type` ASC) USING BTREE,
  INDEX `idx_aiops_evidence_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_evidence
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_knowledge_chunk
-- ----------------------------
DROP TABLE IF EXISTS `aiops_knowledge_chunk`;
CREATE TABLE `aiops_knowledge_chunk`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `document_id` bigint UNSIGNED NOT NULL,
  `chunk_index` bigint NOT NULL,
  `heading` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `content_hash` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `keywords_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `embedding_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending',
  `embedding_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `section_path` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_aiops_knowledge_chunk_document_id`(`document_id` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_chunk_heading`(`heading` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_chunk_content_hash`(`content_hash` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_chunk_embedding_status`(`embedding_status` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_chunk_embedding_id`(`embedding_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 171 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_knowledge_chunk
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_knowledge_document
-- ----------------------------
DROP TABLE IF EXISTS `aiops_knowledge_document`;
CREATE TABLE `aiops_knowledge_document`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `doc_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'incident',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'draft',
  `group_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `category` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `service_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `env` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `tags_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `keywords_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `alert_names_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `log_keywords_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `owners_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `content_markdown` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `summary` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `symptoms_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `root_causes_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `resolution_steps_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `verification_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `rollback_plan` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `source_filename` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_by` bigint UNSIGNED NULL DEFAULT NULL,
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `group_id` bigint UNSIGNED NOT NULL DEFAULT 0,
  `summary_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending',
  `summary_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `search_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'ready',
  `search_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `search_version` bigint NOT NULL DEFAULT 1,
  `chunk_count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_aiops_knowledge_document_title`(`title` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_doc_type`(`doc_type` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_status`(`status` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_group_name`(`group_name` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_category`(`category` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_severity`(`severity` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_service_name`(`service_name` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_env`(`env` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_created_by`(`created_by` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_document_updated_by`(`updated_by` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_knowledge_document
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_knowledge_group
-- ----------------------------
DROP TABLE IF EXISTS `aiops_knowledge_group`;
CREATE TABLE `aiops_knowledge_group`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` bigint UNSIGNED NOT NULL DEFAULT 0,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `sort_order` bigint NOT NULL DEFAULT 0,
  `created_by` bigint UNSIGNED NULL DEFAULT NULL,
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_aiops_knowledge_group_parent_name`(`parent_id` ASC, `name` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_group_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_group_sort_order`(`sort_order` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_group_created_by`(`created_by` ASC) USING BTREE,
  INDEX `idx_aiops_knowledge_group_updated_by`(`updated_by` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_knowledge_group
-- ----------------------------
INSERT INTO `aiops_knowledge_group` VALUES (12, 0, '运维知识库', '运维制度、维护记录、故障经验和技术文档的统一知识库', 10, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (13, 12, '01-运维管理制度', '运维管理制度、规范和流程', 10, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (14, 12, '02-部署维护记录', '系统部署、变更和日常维护记录', 20, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (15, 12, '03-运维故障记录', '故障现象、排查过程、根因和复盘记录', 30, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (16, 12, '04-中间件维护记录', '中间件配置、部署、升级和故障处理记录', 40, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (17, 12, '05-数据库管理', '数据库部署、配置、备份、恢复和优化记录', 50, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (18, 12, '06-运维开发', '运维工具、脚本、平台和自动化开发记录', 60, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (19, 12, '07-运维监控', '监控指标、告警规则和可观测性建设记录', 70, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (20, 12, '08-运维体系建设', '运维流程、标准化和体系建设记录', 80, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (21, 12, '09-运维安全管理', '账号权限、审计、漏洞和安全管理记录', 90, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');
INSERT INTO `aiops_knowledge_group` VALUES (22, 12, '10-总结记录', '阶段总结、经验沉淀和改进记录', 100, 0, 0, '2026-08-03 01:20:03.394', '2026-08-03 01:20:03.394');

-- ----------------------------
-- Table structure for aiops_llm_provider
-- ----------------------------
DROP TABLE IF EXISTS `aiops_llm_provider`;
CREATE TABLE `aiops_llm_provider`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `provider_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `protocol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `base_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `chat_path` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `api_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `api_secret` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `auth_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'bearer',
  `organization` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `project` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `default_model` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `models` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `headers` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `timeout_seconds` bigint NOT NULL DEFAULT 60,
  `max_retries` bigint NOT NULL DEFAULT 2,
  `max_context_messages` bigint NOT NULL DEFAULT 20,
  `temperature` decimal(5, 3) NOT NULL DEFAULT 0.700,
  `max_tokens` bigint NOT NULL DEFAULT 0,
  `status` bigint NOT NULL DEFAULT 1,
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_aiops_llm_provider_name`(`name` ASC) USING BTREE,
  INDEX `idx_aiops_llm_provider_provider_type`(`provider_type` ASC) USING BTREE,
  INDEX `idx_aiops_llm_provider_status`(`status` ASC) USING BTREE,
  INDEX `idx_aiops_llm_provider_is_default`(`is_default` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_llm_provider
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_mcp_server
-- ----------------------------
DROP TABLE IF EXISTS `aiops_mcp_server`;
CREATE TABLE `aiops_mcp_server`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `enabled` tinyint(1) NOT NULL,
  `config_version` bigint UNSIGNED NOT NULL DEFAULT 1,
  `created_by` bigint UNSIGNED NULL DEFAULT NULL,
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_aiops_mcp_server_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_mcp_server
-- ----------------------------

-- ----------------------------
-- Table structure for aiops_mcp_tool_config
-- ----------------------------
DROP TABLE IF EXISTS `aiops_mcp_tool_config`;
CREATE TABLE `aiops_mcp_tool_config`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `server_code` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `tool_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT 0,
  `risk_level` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `required_confirmation` tinyint(1) NOT NULL DEFAULT 0,
  `timeout_seconds` bigint NOT NULL DEFAULT 30,
  `rate_limit_per_minute` bigint NOT NULL DEFAULT 60,
  `max_result_rows` bigint NOT NULL DEFAULT 1000,
  `max_result_bytes` bigint NOT NULL DEFAULT 1048576,
  `allowed_environments` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `config_version` bigint UNSIGNED NOT NULL DEFAULT 1,
  `created_by` bigint UNSIGNED NULL DEFAULT NULL,
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_aiops_mcp_tool`(`server_code` ASC, `tool_name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of aiops_mcp_tool_config
-- ----------------------------

-- ----------------------------
-- Table structure for alert_active_events
-- ----------------------------
DROP TABLE IF EXISTS `alert_active_events`;
CREATE TABLE `alert_active_events`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `fingerprint` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `rule_id` bigint UNSIGNED NULL DEFAULT NULL,
  `rule_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `datasource_id` bigint UNSIGNED NULL DEFAULT NULL,
  `datasource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `labels` json NULL,
  `annotations` json NULL,
  `value` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `trigger_query` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `trigger_condition` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `event_detail` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `asset_id` bigint UNSIGNED NULL DEFAULT NULL,
  `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `first_triggered_at` datetime(3) NULL DEFAULT NULL,
  `last_evaluated_at` datetime(3) NULL DEFAULT NULL,
  `last_notified_at` datetime(3) NULL DEFAULT NULL,
  `recovered_at` datetime(3) NULL DEFAULT NULL,
  `notify_count` int NULL DEFAULT 0,
  `suppressed` tinyint NULL DEFAULT 0,
  `silenced` tinyint NULL DEFAULT 0,
  `claim_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `claimed_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_alert_active_events_fingerprint`(`fingerprint` ASC) USING BTREE,
  INDEX `idx_alert_active_events_rule_id`(`rule_id` ASC) USING BTREE,
  INDEX `idx_alert_active_events_rule_name`(`rule_name` ASC) USING BTREE,
  INDEX `idx_alert_active_events_datasource_id`(`datasource_id` ASC) USING BTREE,
  INDEX `idx_alert_active_events_datasource_type`(`datasource_type` ASC) USING BTREE,
  INDEX `idx_alert_active_events_severity`(`severity` ASC) USING BTREE,
  INDEX `idx_alert_active_events_status`(`status` ASC) USING BTREE,
  INDEX `idx_alert_active_events_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_alert_active_events_business_id`(`business_id` ASC) USING BTREE,
  INDEX `idx_alert_active_events_suppressed`(`suppressed` ASC) USING BTREE,
  INDEX `idx_alert_active_events_silenced`(`silenced` ASC) USING BTREE,
  INDEX `idx_alert_active_events_claim_user`(`claim_user` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_active_events
-- ----------------------------

-- ----------------------------
-- Table structure for alert_builtin_rule_template_groups
-- ----------------------------
DROP TABLE IF EXISTS `alert_builtin_rule_template_groups`;
CREATE TABLE `alert_builtin_rule_template_groups`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` bigint UNSIGNED NULL DEFAULT 0,
  `level` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `resource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `sort` int NULL DEFAULT 0,
  `builtin` tinyint NOT NULL DEFAULT 0,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_builtin_rule_template_groups_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_template_groups_level`(`level` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_template_groups_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_template_groups_value`(`value` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_template_groups_resource_type`(`resource_type` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_template_groups_builtin`(`builtin` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_template_groups_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 32 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_builtin_rule_template_groups
-- ----------------------------
INSERT INTO `alert_builtin_rule_template_groups` VALUES (1, 0, 'component', 'Linux', 'Linux', 'host', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (2, 1, 'collector', 'node_exporter', 'node_exporter', 'host', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (3, 2, 'category', 'cpu', 'cpu', 'host', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (4, 2, 'category', 'memory', 'memory', 'host', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (5, 2, 'category', 'disk', 'disk', 'host', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (6, 2, 'category', 'availability', 'availability', 'host', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (7, 0, 'component', 'Kubernetes', 'Kubernetes', 'k8s', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (8, 7, 'collector', 'kube-state-metrics', 'kube-state-metrics', 'k8s', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (9, 8, 'category', 'node', 'node', 'k8s', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (10, 8, 'category', 'deployment', 'deployment', 'k8s', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (11, 8, 'category', 'pod', 'pod', 'k8s', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (12, 0, 'component', 'MySQL', 'MySQL', 'mysql', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (13, 12, 'collector', 'mysqld_exporter', 'mysqld_exporter', 'mysql', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (14, 13, 'category', 'connection', 'connection', 'mysql', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (15, 13, 'category', 'file', 'file', 'mysql', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (16, 13, 'category', 'replication', 'replication', 'mysql', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (17, 0, 'component', 'Prometheus', 'Prometheus', 'prometheus', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (18, 17, 'collector', 'prometheus', 'prometheus', 'prometheus', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (19, 18, 'category', 'target', 'target', 'prometheus', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (20, 0, 'component', 'VictoriaMetrics', 'VictoriaMetrics', 'victoriametrics', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (21, 20, 'collector', 'victoriametrics', 'victoriametrics', 'victoriametrics', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (22, 21, 'category', 'availability', 'availability', 'victoriametrics', '', 0, 1, 1, 'system', 'system', '2026-07-10 23:50:48.978', '2026-07-10 23:50:48.978');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (23, 1, 'collector', 'Agent 监控模板', 'agent', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (24, 23, 'category', 'CPU', 'cpu', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (25, 23, 'category', '内存', 'memory', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (26, 23, 'category', '负载', 'load', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (27, 23, 'category', '磁盘', 'disk', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (28, 23, 'category', '磁盘 IO', 'disk_io', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (29, 23, 'category', '网络', 'network', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (30, 23, 'category', '连接数', 'connection', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');
INSERT INTO `alert_builtin_rule_template_groups` VALUES (31, 23, 'category', '服务', 'service', 'host', '', 0, 1, 1, 'system', 'system', '2026-08-02 15:22:31.372', '2026-08-02 15:22:31.372');

-- ----------------------------
-- Table structure for alert_builtin_rule_templates
-- ----------------------------
DROP TABLE IF EXISTS `alert_builtin_rule_templates`;
CREATE TABLE `alert_builtin_rule_templates`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `resource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'host',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `collector` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `category` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `datasource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `source` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'manual',
  `source_uuid` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `query_template` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `default_condition_operator` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `default_condition_value` double NULL DEFAULT NULL,
  `default_eval_interval_seconds` bigint NOT NULL DEFAULT 60,
  `default_for_duration_seconds` bigint NOT NULL DEFAULT 0,
  `default_severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'warning',
  `default_enabled` tinyint NOT NULL DEFAULT 0,
  `labels` json NULL,
  `annotations` json NULL,
  `variables` json NULL,
  `tags` json NULL,
  `rule_config` json NULL,
  `sort` int NULL DEFAULT 0,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_resource_type`(`resource_type` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_component`(`component` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_collector`(`collector` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_category`(`category` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_datasource_type`(`datasource_type` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_source`(`source` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_source_uuid`(`source_uuid` ASC) USING BTREE,
  INDEX `idx_alert_builtin_rule_templates_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_builtin_rule_templates
-- ----------------------------
INSERT INTO `alert_builtin_rule_templates` VALUES (1, 'host', 'Linux', 'node_exporter', 'cpu', '主机 CPU 使用率过高', '参考 n9e Linux 内置规则，检测 CPU 使用率持续过高', 'prometheus', 'platform', 'platform.host.cpu.high', '100 - (avg by(instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100) > {{threshold}}', '>', 85, 60, 300, 'warning', 0, '{\"category\": \"cpu\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"主机 CPU 使用率过高\", \"description\": \"实例 {{ $labels.instance }} CPU 使用率持续超过 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 85}}', NULL, '{\"inhibit\": false, \"queries\": [{\"prom_ql\": \"100 - (avg by(instance) (irate(node_cpu_seconds_total{mode=\\\"idle\\\"}[5m])) * 100) > 75\", \"severity\": \"info\"}, {\"prom_ql\": \"100 - (avg by(instance) (irate(node_cpu_seconds_total{mode=\\\"idle\\\"}[5m])) * 100) > 85\", \"severity\": \"warning\"}, {\"prom_ql\": \"100 - (avg by(instance) (irate(node_cpu_seconds_total{mode=\\\"idle\\\"}[5m])) * 100) > 95\", \"severity\": \"critical\"}]}', 100, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (2, 'host', 'Linux', 'node_exporter', 'memory', '主机内存可用率过低', '检测 Linux 主机可用内存比例过低', 'prometheus', 'platform', 'platform.host.memory.available.low', 'node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < {{threshold}}', '<', 10, 60, 300, 'warning', 0, '{\"category\": \"memory\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"主机内存可用率过低\", \"description\": \"实例 {{ $labels.instance }} 可用内存低于 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 10}}', NULL, NULL, 95, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (3, 'host', 'Linux', 'node_exporter', 'disk', '主机磁盘使用率过高', '参考 n9e Linux 内置规则，检测文件系统使用率过高', 'prometheus', 'platform', 'platform.host.disk.usage.high', '(100 - ((node_filesystem_avail_bytes{fstype!~\"tmpfs|overlay\"} * 100) / node_filesystem_size_bytes{fstype!~\"tmpfs|overlay\"})) > {{threshold}}', '>', 92, 60, 300, 'critical', 0, '{\"category\": \"disk\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"主机磁盘使用率过高\", \"description\": \"实例 {{ $labels.instance }} 挂载点 {{ $labels.mountpoint }} 使用率超过 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 92}}', NULL, NULL, 90, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (4, 'host', 'Linux', 'node_exporter', 'availability', '主机 Exporter 离线', '检测 Prometheus 抓取目标不可用', 'prometheus', 'platform', 'platform.host.exporter.down', 'up{job=~\"node.*|node_exporter\"} == 0', '==', 0, 60, 180, 'critical', 0, '{\"category\": \"availability\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"主机 Exporter 离线\", \"description\": \"实例 {{ $labels.instance }} Exporter 无法抓取。\"}', NULL, NULL, NULL, 85, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (5, 'k8s', 'Kubernetes', 'kube-state-metrics', 'node', 'K8S Node NotReady', '参考 n9e Kubernetes 内置规则，检测节点 Ready 状态异常', 'prometheus', 'platform', 'platform.k8s.node.notready', 'kube_node_status_condition{job=\"kube-state-metrics\",condition=\"Ready\",status=\"true\"} == 0', '==', 0, 60, 300, 'critical', 0, '{\"category\": \"node\", \"component\": \"Kubernetes\", \"resourceType\": \"k8s\"}', '{\"summary\": \"K8S Node NotReady\", \"description\": \"节点 {{ $labels.node }} 处于 NotReady 状态。\"}', NULL, NULL, NULL, 100, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (6, 'k8s', 'Kubernetes', 'kube-state-metrics', 'deployment', 'Deployment 副本异常', '参考 n9e Kubernetes k8s-workloads 规则，检测 Deployment 可用副本不足', 'prometheus', 'platform', 'platform.k8s.deployment.replicas.mismatch', 'kube_deployment_spec_replicas{job=\"kube-state-metrics\"} != kube_deployment_status_replicas_available{job=\"kube-state-metrics\"}', '!=', 0, 60, 600, 'warning', 0, '{\"category\": \"deployment\", \"component\": \"Kubernetes\", \"resourceType\": \"k8s\"}', '{\"summary\": \"Deployment 副本异常\", \"description\": \"命名空间 {{ $labels.namespace }} 的 Deployment {{ $labels.deployment }} 可用副本与期望副本不一致。\"}', NULL, NULL, NULL, 95, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (7, 'k8s', 'Kubernetes', 'kube-state-metrics', 'pod', 'Pod 失败', '检测 Pod phase 为 Failed', 'prometheus', 'platform', 'platform.k8s.pod.failed', 'kube_pod_status_phase{job=\"kube-state-metrics\",phase=\"Failed\"} > 0', '>', 0, 60, 300, 'warning', 0, '{\"category\": \"pod\", \"component\": \"Kubernetes\", \"resourceType\": \"k8s\"}', '{\"summary\": \"Pod 失败\", \"description\": \"命名空间 {{ $labels.namespace }} 的 Pod {{ $labels.pod }} 处于 Failed 状态。\"}', NULL, NULL, NULL, 90, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (8, 'mysql', 'MySQL', 'mysqld_exporter', 'connection', 'MySQL 运行连接占比过高', '参考 n9e MySQL exporter 规则，检测 Threads_running 占最大连接数比例', 'prometheus', 'platform', 'platform.mysql.threads.running.high', 'avg by (instance) (mysql_global_status_threads_running) / avg by (instance) (mysql_global_variables_max_connections) * 100 > {{threshold}}', '>', 60, 60, 120, 'warning', 0, '{\"category\": \"connection\", \"component\": \"MySQL\", \"resourceType\": \"mysql\"}', '{\"summary\": \"MySQL 运行连接占比过高\", \"description\": \"实例 {{ $labels.instance }} 运行中连接占比超过 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 60}}', NULL, NULL, 100, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (9, 'mysql', 'MySQL', 'mysqld_exporter', 'file', 'MySQL 打开文件数占比过高', '参考 n9e MySQL exporter 规则，检测打开文件数接近限制', 'prometheus', 'platform', 'platform.mysql.open.files.high', 'avg by (instance) (mysql_global_status_innodb_num_open_files) / avg by (instance)(mysql_global_variables_open_files_limit) * 100 > {{threshold}}', '>', 80, 60, 120, 'warning', 0, '{\"category\": \"file\", \"component\": \"MySQL\", \"resourceType\": \"mysql\"}', '{\"summary\": \"MySQL 打开文件数占比过高\", \"description\": \"实例 {{ $labels.instance }} 打开文件数占比超过 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 80}}', NULL, NULL, 95, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (10, 'mysql', 'MySQL', 'mysqld_exporter', 'replication', 'MySQL 主从复制中断', '检测 MySQL Slave IO 或 SQL 线程异常', 'prometheus', 'platform', 'platform.mysql.replication.down', 'mysql_slave_status_slave_io_running == 0 or mysql_slave_status_slave_sql_running == 0', '==', 0, 60, 120, 'critical', 0, '{\"category\": \"replication\", \"component\": \"MySQL\", \"resourceType\": \"mysql\"}', '{\"summary\": \"MySQL 主从复制中断\", \"description\": \"实例 {{ $labels.instance }} 主从复制线程异常。\"}', NULL, NULL, NULL, 90, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (11, 'prometheus', 'Prometheus', 'prometheus', 'target', 'Prometheus Target Down', '检测 Prometheus 抓取目标不可用', 'prometheus', 'platform', 'platform.prometheus.target.down', 'up == 0', '==', 0, 60, 180, 'critical', 0, '{\"category\": \"target\", \"component\": \"Prometheus\", \"resourceType\": \"prometheus\"}', '{\"summary\": \"Prometheus Target Down\", \"description\": \"抓取目标 {{ $labels.job }}/{{ $labels.instance }} 不可用。\"}', NULL, NULL, NULL, 100, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (12, 'victoriametrics', 'VictoriaMetrics', 'victoriametrics', 'availability', 'VictoriaMetrics 实例不可用', '检测 VictoriaMetrics 自身采集目标不可用', 'victoriametrics', 'platform', 'platform.victoriametrics.down', 'up{job=~\"victoriametrics|vm.*\"} == 0', '==', 0, 60, 180, 'critical', 0, '{\"category\": \"availability\", \"component\": \"VictoriaMetrics\", \"resourceType\": \"victoriametrics\"}', '{\"summary\": \"VictoriaMetrics 实例不可用\", \"description\": \"VictoriaMetrics 目标 {{ $labels.instance }} 不可用。\"}', NULL, NULL, NULL, 100, 1, '2026-07-10 23:50:48.938', '2026-07-10 23:50:48.938');
INSERT INTO `alert_builtin_rule_templates` VALUES (13, 'host', 'Linux', 'agent', 'cpu', 'Agent CPU 使用率过高', 'CPU 使用率（system_cpu_usage_percent）', 'prometheus', 'platform', 'platform.agent.system_cpu_usage_percent', 'avg by(instance) (system_cpu_usage_percent) > {{threshold}}', '>', 85, 60, 300, 'warning', 0, '{\"category\": \"cpu\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent CPU 使用率过高\", \"description\": \"实例 {{ $labels.instance }} CPU 使用率，触发阈值 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 85}}', '{\"metric\": \"system_cpu_usage_percent\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (14, 'host', 'Linux', 'agent', 'cpu', 'Agent 用户态 CPU 使用率过高', '用户空间 CPU 占用比率（system_cpu_user_percent）', 'prometheus', 'platform', 'platform.agent.system_cpu_user_percent', 'system_cpu_user_percent > {{threshold}}', '>', 80, 60, 300, 'warning', 0, '{\"category\": \"cpu\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 用户态 CPU 使用率过高\", \"description\": \"实例 {{ $labels.instance }} 用户空间 CPU 占用比率，触发阈值 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 80}}', '{\"metric\": \"system_cpu_user_percent\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (15, 'host', 'Linux', 'agent', 'cpu', 'Agent 内核态 CPU 使用率过高', '内核空间 CPU 占用比率（system_cpu_system_percent）', 'prometheus', 'platform', 'platform.agent.system_cpu_system_percent', 'system_cpu_system_percent > {{threshold}}', '>', 40, 60, 300, 'warning', 0, '{\"category\": \"cpu\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 内核态 CPU 使用率过高\", \"description\": \"实例 {{ $labels.instance }} 内核空间 CPU 占用比率，触发阈值 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 40}}', '{\"metric\": \"system_cpu_system_percent\", \"collector\": \"agent\"}', NULL, 90, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (16, 'host', 'Linux', 'agent', 'memory', 'Agent 内存使用率过高', '内存使用率（system_memory_usage_percent）', 'prometheus', 'platform', 'platform.agent.system_memory_usage_percent', 'system_memory_usage_percent > {{threshold}}', '>', 85, 60, 300, 'warning', 0, '{\"category\": \"memory\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 内存使用率过高\", \"description\": \"实例 {{ $labels.instance }} 内存使用率，触发阈值 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 85}}', '{\"metric\": \"system_memory_usage_percent\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (17, 'host', 'Linux', 'agent', 'memory', 'Agent 已用内存过高', '已用内存实际大小（system_memory_used_gb）', 'prometheus', 'platform', 'platform.agent.system_memory_used_gb', 'system_memory_used_gb > {{threshold}}', '>', 8, 60, 300, 'warning', 0, '{\"category\": \"memory\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 已用内存过高\", \"description\": \"实例 {{ $labels.instance }} 已用内存实际大小，触发阈值 {{threshold}}GB。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"GB\", \"default\": 8}}', '{\"metric\": \"system_memory_used_gb\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (18, 'host', 'Linux', 'agent', 'memory', 'Agent 可用内存过低', '可用内存实际大小（system_memory_free_gb）', 'prometheus', 'platform', 'platform.agent.system_memory_free_gb', 'system_memory_free_gb < {{threshold}}', '<', 1, 60, 300, 'warning', 0, '{\"category\": \"memory\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 可用内存过低\", \"description\": \"实例 {{ $labels.instance }} 可用内存实际大小，触发阈值 {{threshold}}GB。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"GB\", \"default\": 1}}', '{\"metric\": \"system_memory_free_gb\", \"collector\": \"agent\"}', NULL, 90, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (19, 'host', 'Linux', 'agent', 'memory', 'Agent Swap 使用量过高', '已使用交换分区大小（system_swap_used_gb）', 'prometheus', 'platform', 'platform.agent.system_swap_used_gb', 'system_swap_used_gb > {{threshold}}', '>', 1, 60, 300, 'warning', 0, '{\"category\": \"memory\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent Swap 使用量过高\", \"description\": \"实例 {{ $labels.instance }} 已使用交换分区大小，触发阈值 {{threshold}}GB。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"GB\", \"default\": 1}}', '{\"metric\": \"system_swap_used_gb\", \"collector\": \"agent\"}', NULL, 85, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (20, 'host', 'Linux', 'agent', 'load', 'Agent 1 分钟平均负载过高', '系统 1 分钟平均负载（system_load_average_1min）', 'prometheus', 'platform', 'platform.agent.system_load_average_1min', 'system_load_average_1min > {{threshold}}', '>', 4, 60, 300, 'warning', 0, '{\"category\": \"load\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 1 分钟平均负载过高\", \"description\": \"实例 {{ $labels.instance }} 系统 1 分钟平均负载，触发阈值 {{threshold}}。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"\", \"default\": 4}}', '{\"metric\": \"system_load_average_1min\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (21, 'host', 'Linux', 'agent', 'load', 'Agent 5 分钟平均负载过高', '系统 5 分钟平均负载（system_load_average_5min）', 'prometheus', 'platform', 'platform.agent.system_load_average_5min', 'system_load_average_5min > {{threshold}}', '>', 4, 60, 300, 'warning', 0, '{\"category\": \"load\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 5 分钟平均负载过高\", \"description\": \"实例 {{ $labels.instance }} 系统 5 分钟平均负载，触发阈值 {{threshold}}。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"\", \"default\": 4}}', '{\"metric\": \"system_load_average_5min\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (22, 'host', 'Linux', 'agent', 'load', 'Agent 15 分钟平均负载过高', '系统 15 分钟平均负载（system_load_average_15min）', 'prometheus', 'platform', 'platform.agent.system_load_average_15min', 'system_load_average_15min > {{threshold}}', '>', 4, 60, 300, 'warning', 0, '{\"category\": \"load\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 15 分钟平均负载过高\", \"description\": \"实例 {{ $labels.instance }} 系统 15 分钟平均负载，触发阈值 {{threshold}}。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"\", \"default\": 4}}', '{\"metric\": \"system_load_average_15min\", \"collector\": \"agent\"}', NULL, 90, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (23, 'host', 'Linux', 'agent', 'load', 'Agent 1 分钟负载比过高', '系统 1 分钟负载比（system_load_ratio_1min）', 'prometheus', 'platform', 'platform.agent.system_load_ratio_1min', 'system_load_ratio_1min > {{threshold}}', '>', 1, 60, 300, 'warning', 0, '{\"category\": \"load\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 1 分钟负载比过高\", \"description\": \"实例 {{ $labels.instance }} 系统 1 分钟负载比，触发阈值 {{threshold}}。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"\", \"default\": 1}}', '{\"metric\": \"system_load_ratio_1min\", \"collector\": \"agent\"}', NULL, 85, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (24, 'host', 'Linux', 'agent', 'load', 'Agent 5 分钟负载比过高', '系统 5 分钟负载比（system_load_ratio_5min）', 'prometheus', 'platform', 'platform.agent.system_load_ratio_5min', 'system_load_ratio_5min > {{threshold}}', '>', 1, 60, 300, 'warning', 0, '{\"category\": \"load\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 5 分钟负载比过高\", \"description\": \"实例 {{ $labels.instance }} 系统 5 分钟负载比，触发阈值 {{threshold}}。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"\", \"default\": 1}}', '{\"metric\": \"system_load_ratio_5min\", \"collector\": \"agent\"}', NULL, 80, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (25, 'host', 'Linux', 'agent', 'load', 'Agent 15 分钟负载比过高', '系统 15 分钟负载比（system_load_ratio_15min）', 'prometheus', 'platform', 'platform.agent.system_load_ratio_15min', 'system_load_ratio_15min > {{threshold}}', '>', 1, 60, 300, 'warning', 0, '{\"category\": \"load\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 15 分钟负载比过高\", \"description\": \"实例 {{ $labels.instance }} 系统 15 分钟负载比，触发阈值 {{threshold}}。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"\", \"default\": 1}}', '{\"metric\": \"system_load_ratio_15min\", \"collector\": \"agent\"}', NULL, 75, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (26, 'host', 'Linux', 'agent', 'disk', 'Agent 根分区磁盘使用率过高', '根分区磁盘使用率（system_disk_usage_percent）', 'prometheus', 'platform', 'platform.agent.system_disk_usage_percent', 'system_disk_usage_percent{mountpoint=\"/\"} > {{threshold}}', '>', 85, 60, 300, 'warning', 0, '{\"category\": \"disk\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 根分区磁盘使用率过高\", \"description\": \"实例 {{ $labels.instance }} 根分区磁盘使用率，触发阈值 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 85}}', '{\"metric\": \"system_disk_usage_percent\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (27, 'host', 'Linux', 'agent', 'disk', 'Agent 根分区已用空间过高', '根分区已用磁盘空间（system_disk_used_bytes）', 'prometheus', 'platform', 'platform.agent.system_disk_used_bytes', 'system_disk_used_bytes{mountpoint=\"/\"} > {{threshold}}', '>', 100, 60, 300, 'warning', 0, '{\"category\": \"disk\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 根分区已用空间过高\", \"description\": \"实例 {{ $labels.instance }} 根分区已用磁盘空间，触发阈值 {{threshold}}GB。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"GB\", \"default\": 100}}', '{\"metric\": \"system_disk_used_bytes\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (28, 'host', 'Linux', 'agent', 'disk', 'Agent 根分区 inode 使用率过高', '根分区磁盘 inode 使用率（system_disk_inode_usage_percent）', 'prometheus', 'platform', 'platform.agent.system_disk_inode_usage_percent', 'system_disk_inode_usage_percent{mountpoint=\"/\"} > {{threshold}}', '>', 85, 60, 300, 'warning', 0, '{\"category\": \"disk\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 根分区 inode 使用率过高\", \"description\": \"实例 {{ $labels.instance }} 根分区磁盘 inode 使用率，触发阈值 {{threshold}}%。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"%\", \"default\": 85}}', '{\"metric\": \"system_disk_inode_usage_percent\", \"collector\": \"agent\"}', NULL, 90, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (29, 'host', 'Linux', 'agent', 'disk_io', 'Agent 磁盘读取速率过高', '磁盘 IO 读取速率（system_disk_read_kb_per_second）', 'prometheus', 'platform', 'platform.agent.system_disk_read_kb_per_second', 'sum by(instance) (system_disk_read_kb_per_second) > {{threshold}}', '>', 102400, 60, 300, 'warning', 0, '{\"category\": \"disk_io\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 磁盘读取速率过高\", \"description\": \"实例 {{ $labels.instance }} 磁盘 IO 读取速率，触发阈值 {{threshold}}KB/s。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"KB/s\", \"default\": 102400}}', '{\"metric\": \"system_disk_read_kb_per_second\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (30, 'host', 'Linux', 'agent', 'disk_io', 'Agent 磁盘写入速率过高', '磁盘 IO 写入速率（system_disk_write_kb_per_second）', 'prometheus', 'platform', 'platform.agent.system_disk_write_kb_per_second', 'sum by(instance) (system_disk_write_kb_per_second) > {{threshold}}', '>', 102400, 60, 300, 'warning', 0, '{\"category\": \"disk_io\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 磁盘写入速率过高\", \"description\": \"实例 {{ $labels.instance }} 磁盘 IO 写入速率，触发阈值 {{threshold}}KB/s。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"KB/s\", \"default\": 102400}}', '{\"metric\": \"system_disk_write_kb_per_second\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (31, 'host', 'Linux', 'agent', 'disk_io', 'Agent 磁盘读 IOPS 过高', '磁盘 IO 每秒读操作数（system_disk_read_iops）', 'prometheus', 'platform', 'platform.agent.system_disk_read_iops', 'sum by(instance) (system_disk_read_iops) > {{threshold}}', '>', 1000, 60, 300, 'warning', 0, '{\"category\": \"disk_io\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 磁盘读 IOPS 过高\", \"description\": \"实例 {{ $labels.instance }} 磁盘 IO 每秒读操作数，触发阈值 {{threshold}}IOPS。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"IOPS\", \"default\": 1000}}', '{\"metric\": \"system_disk_read_iops\", \"collector\": \"agent\"}', NULL, 90, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (32, 'host', 'Linux', 'agent', 'disk_io', 'Agent 磁盘写 IOPS 过高', '磁盘 IO 每秒写操作数（system_disk_write_iops）', 'prometheus', 'platform', 'platform.agent.system_disk_write_iops', 'sum by(instance) (system_disk_write_iops) > {{threshold}}', '>', 1000, 60, 300, 'warning', 0, '{\"category\": \"disk_io\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 磁盘写 IOPS 过高\", \"description\": \"实例 {{ $labels.instance }} 磁盘 IO 每秒写操作数，触发阈值 {{threshold}}IOPS。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"IOPS\", \"default\": 1000}}', '{\"metric\": \"system_disk_write_iops\", \"collector\": \"agent\"}', NULL, 85, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (33, 'host', 'Linux', 'agent', 'network', 'Agent 网络接收速率过高', '网络接收速率（system_network_receive_kb_per_second）', 'prometheus', 'platform', 'platform.agent.system_network_receive_kb_per_second', 'system_network_receive_kb_per_second > {{threshold}}', '>', 102400, 60, 300, 'warning', 0, '{\"category\": \"network\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 网络接收速率过高\", \"description\": \"实例 {{ $labels.instance }} 网络接收速率，触发阈值 {{threshold}}KB/s。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"KB/s\", \"default\": 102400}}', '{\"metric\": \"system_network_receive_kb_per_second\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (34, 'host', 'Linux', 'agent', 'network', 'Agent 网络发送速率过高', '网络发送速率（system_network_send_kb_per_second）', 'prometheus', 'platform', 'platform.agent.system_network_send_kb_per_second', 'system_network_send_kb_per_second > {{threshold}}', '>', 102400, 60, 300, 'warning', 0, '{\"category\": \"network\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 网络发送速率过高\", \"description\": \"实例 {{ $labels.instance }} 网络发送速率，触发阈值 {{threshold}}KB/s。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"KB/s\", \"default\": 102400}}', '{\"metric\": \"system_network_send_kb_per_second\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (35, 'host', 'Linux', 'agent', 'connection', 'Agent TCP 连接数过高', 'TCP 连接总数（system_tcp_connections）', 'prometheus', 'platform', 'platform.agent.system_tcp_connections', 'system_tcp_connections > {{threshold}}', '>', 2000, 60, 300, 'warning', 0, '{\"category\": \"connection\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent TCP 连接数过高\", \"description\": \"实例 {{ $labels.instance }} TCP 连接总数，触发阈值 {{threshold}}个。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"个\", \"default\": 2000}}', '{\"metric\": \"system_tcp_connections\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (36, 'host', 'Linux', 'agent', 'connection', 'Agent UDP 连接数过高', 'UDP 连接总数（system_udp_connections）', 'prometheus', 'platform', 'platform.agent.system_udp_connections', 'system_udp_connections > {{threshold}}', '>', 1000, 60, 300, 'warning', 0, '{\"category\": \"connection\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent UDP 连接数过高\", \"description\": \"实例 {{ $labels.instance }} UDP 连接总数，触发阈值 {{threshold}}个。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"个\", \"default\": 1000}}', '{\"metric\": \"system_udp_connections\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (37, 'host', 'Linux', 'agent', 'service', 'Agent 服务端口未监听', '服务端口监听状态（tcp_port_listening）', 'prometheus', 'platform', 'platform.agent.tcp_port_listening', 'tcp_port_listening == {{threshold}}', '==', 0, 60, 180, 'critical', 0, '{\"category\": \"service\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 服务端口未监听\", \"description\": \"实例 {{ $labels.instance }} 服务端口监听状态，触发阈值 {{threshold}}。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"\", \"default\": 0}}', '{\"metric\": \"tcp_port_listening\", \"collector\": \"agent\"}', NULL, 100, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (38, 'host', 'Linux', 'agent', 'service', 'Agent 系统进程数过高', '系统进程总数（system_total_processes）', 'prometheus', 'platform', 'platform.agent.system_total_processes', 'system_total_processes > {{threshold}}', '>', 500, 60, 300, 'warning', 0, '{\"category\": \"service\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 系统进程数过高\", \"description\": \"实例 {{ $labels.instance }} 系统进程总数，触发阈值 {{threshold}}个。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"个\", \"default\": 500}}', '{\"metric\": \"system_total_processes\", \"collector\": \"agent\"}', NULL, 95, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (39, 'host', 'Linux', 'agent', 'service', 'Agent 系统最近重启', '系统持续运行天数（system_uptime_days）', 'prometheus', 'platform', 'platform.agent.system_uptime_days', 'system_uptime_days < {{threshold}}', '<', 1, 60, 60, 'info', 0, '{\"category\": \"service\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 系统最近重启\", \"description\": \"实例 {{ $labels.instance }} 系统持续运行天数，触发阈值 {{threshold}}天。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"天\", \"default\": 1}}', '{\"metric\": \"system_uptime_days\", \"collector\": \"agent\"}', NULL, 90, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');
INSERT INTO `alert_builtin_rule_templates` VALUES (40, 'host', 'Linux', 'agent', 'service', 'Agent 打开文件句柄数过高', '系统当前打开的文件句柄数（system_open_file_descriptors）', 'prometheus', 'platform', 'platform.agent.system_open_file_descriptors', 'system_open_file_descriptors > {{threshold}}', '>', 10000, 60, 300, 'warning', 0, '{\"category\": \"service\", \"component\": \"Linux\", \"resourceType\": \"host\"}', '{\"summary\": \"Agent 打开文件句柄数过高\", \"description\": \"实例 {{ $labels.instance }} 系统当前打开的文件句柄数，触发阈值 {{threshold}}个。\"}', '{\"threshold\": {\"type\": \"number\", \"unit\": \"个\", \"default\": 10000}}', '{\"metric\": \"system_open_file_descriptors\", \"collector\": \"agent\"}', NULL, 85, 1, '2026-08-02 15:22:31.292', '2026-08-02 15:22:31.292');

-- ----------------------------
-- Table structure for alert_datasources
-- ----------------------------
DROP TABLE IF EXISTS `alert_datasources`;
CREATE TABLE `alert_datasources`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `datasource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `auth_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'none',
  `username` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `password` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `token` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `config` json NULL,
  `labels` json NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `builtin` tinyint NOT NULL DEFAULT 0,
  `is_default` tinyint NOT NULL DEFAULT 0,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_datasources_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_datasources_datasource_type`(`datasource_type` ASC) USING BTREE,
  INDEX `idx_alert_datasources_builtin`(`builtin` ASC) USING BTREE,
  INDEX `idx_alert_datasources_default`(`is_default` ASC) USING BTREE,
  INDEX `idx_alert_datasources_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_datasources
-- ----------------------------
INSERT INTO `alert_datasources` VALUES (1, '平台内置告警源', 'builtin', 'http://112.74.57.93:8428', 'none', '', '', '', '{\"backend\": \"victoriametrics\", \"queryType\": \"promql\"}', '{\"source\": \"platform\", \"backend\": \"victoriametrics\"}', '平台内置 VictoriaMetrics 告警查询数据源', 1, 1, 1, 'system', 'system', '2026-07-10 23:50:48.926', '2026-11-10 21:12:16.165');
INSERT INTO `alert_datasources` VALUES (2, '平台内置日志源', 'victorialogs', 'http://112.74.57.93:9428', 'none', '', '', '', '{\"backend\": \"victorialogs\", \"queryType\": \"logsql\"}', '{\"source\": \"platform\", \"backend\": \"victorialogs\"}', '平台内置 VictoriaLogs 日志告警查询数据源', 1, 0, 1, 'system', 'system', '2026-07-10 23:50:48.926', '2026-11-10 21:12:16.165');

-- ----------------------------
-- Table structure for alert_domains
-- ----------------------------
DROP TABLE IF EXISTS `alert_domains`;
CREATE TABLE `alert_domains`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `domain` varchar(253) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `protocol` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'https',
  `port` int NOT NULL DEFAULT 0,
  `path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '/',
  `tags` json NULL,
  `expected_status` int NOT NULL DEFAULT 0,
  `check_interval_seconds` bigint NOT NULL DEFAULT 300,
  `timeout_seconds` bigint NOT NULL DEFAULT 10,
  `ssl_expire_warn_days` int NOT NULL DEFAULT 30,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `remarks` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `alert_rule_id` bigint UNSIGNED NOT NULL DEFAULT 0,
  `notify_rule_ids` json NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'unknown',
  `status_code` int NOT NULL DEFAULT 0,
  `response_time_ms` bigint NOT NULL DEFAULT 0,
  `certificate_expires_at` datetime(3) NULL DEFAULT NULL,
  `certificate_remaining_days` int NOT NULL DEFAULT 0,
  `certificate_issuer` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `failure_reason` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `last_checked_at` datetime(3) NULL DEFAULT NULL,
  `next_check_at` datetime(3) NULL DEFAULT NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_alert_domains_domain`(`domain` ASC) USING BTREE,
  INDEX `idx_alert_domains_enabled`(`enabled` ASC) USING BTREE,
  INDEX `idx_alert_domains_alert_rule_id`(`alert_rule_id` ASC) USING BTREE,
  INDEX `idx_alert_domains_status`(`status` ASC) USING BTREE,
  INDEX `idx_alert_domains_last_checked_at`(`last_checked_at` ASC) USING BTREE,
  INDEX `idx_alert_domains_next_check_at`(`next_check_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_domains
-- ----------------------------

-- ----------------------------
-- Table structure for alert_event_timeline
-- ----------------------------
DROP TABLE IF EXISTS `alert_event_timeline`;
CREATE TABLE `alert_event_timeline`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id` bigint UNSIGNED NULL DEFAULT NULL,
  `action` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `message` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `operator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `metadata` json NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_event_timeline_event_id`(`event_id` ASC) USING BTREE,
  INDEX `idx_alert_event_timeline_action`(`action` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_event_timeline
-- ----------------------------

-- ----------------------------
-- Table structure for alert_history_events
-- ----------------------------
DROP TABLE IF EXISTS `alert_history_events`;
CREATE TABLE `alert_history_events`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `fingerprint` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `active_event_id` bigint UNSIGNED NULL DEFAULT NULL,
  `rule_id` bigint UNSIGNED NULL DEFAULT NULL,
  `rule_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `datasource_id` bigint UNSIGNED NULL DEFAULT NULL,
  `datasource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `event_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `labels` json NULL,
  `annotations` json NULL,
  `value` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `trigger_query` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `trigger_condition` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `event_detail` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `asset_id` bigint UNSIGNED NULL DEFAULT NULL,
  `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `first_triggered_at` datetime(3) NULL DEFAULT NULL,
  `recovered_at` datetime(3) NULL DEFAULT NULL,
  `claim_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `occurred_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_history_events_fingerprint`(`fingerprint` ASC) USING BTREE,
  INDEX `idx_alert_history_events_active_event_id`(`active_event_id` ASC) USING BTREE,
  INDEX `idx_alert_history_events_rule_id`(`rule_id` ASC) USING BTREE,
  INDEX `idx_alert_history_events_rule_name`(`rule_name` ASC) USING BTREE,
  INDEX `idx_alert_history_events_datasource_id`(`datasource_id` ASC) USING BTREE,
  INDEX `idx_alert_history_events_datasource_type`(`datasource_type` ASC) USING BTREE,
  INDEX `idx_alert_history_events_event_type`(`event_type` ASC) USING BTREE,
  INDEX `idx_alert_history_events_severity`(`severity` ASC) USING BTREE,
  INDEX `idx_alert_history_events_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_alert_history_events_business_id`(`business_id` ASC) USING BTREE,
  INDEX `idx_alert_history_events_claim_user`(`claim_user` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_history_events
-- ----------------------------

-- ----------------------------
-- Table structure for alert_notify_channels
-- ----------------------------
DROP TABLE IF EXISTS `alert_notify_channels`;
CREATE TABLE `alert_notify_channels`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `config` json NULL,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_notify_channels_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_notify_channels_type`(`type` ASC) USING BTREE,
  INDEX `idx_alert_notify_channels_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_notify_channels
-- ----------------------------

-- ----------------------------
-- Table structure for alert_notify_groups
-- ----------------------------
DROP TABLE IF EXISTS `alert_notify_groups`;
CREATE TABLE `alert_notify_groups`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `channel_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'leliao',
  `user_ids` json NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_alert_notify_groups_business_id`(`business_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_groups_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_notify_groups_channel_type`(`channel_type` ASC) USING BTREE,
  INDEX `idx_alert_notify_groups_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_notify_groups
-- ----------------------------

-- ----------------------------
-- Table structure for alert_notify_records
-- ----------------------------
DROP TABLE IF EXISTS `alert_notify_records`;
CREATE TABLE `alert_notify_records`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id` bigint UNSIGNED NULL DEFAULT NULL,
  `event_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `notify_group_key` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `notify_rule_id` bigint UNSIGNED NULL DEFAULT NULL,
  `rule_id` bigint UNSIGNED NULL DEFAULT NULL,
  `rule_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `channel_id` bigint UNSIGNED NULL DEFAULT NULL,
  `channel_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `channel_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `datasource_id` bigint UNSIGNED NULL DEFAULT NULL,
  `datasource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `receiver` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `error_message` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `event_count` int NULL DEFAULT 0,
  `resource_count` int NULL DEFAULT 0,
  `rule_names` json NULL,
  `request` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `response` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `operator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `notify_group_id` bigint UNSIGNED NULL DEFAULT NULL,
  `notify_group_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_notify_records_event_id`(`event_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_event_type`(`event_type` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_notify_group_key`(`notify_group_key` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_notify_rule_id`(`notify_rule_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_rule_id`(`rule_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_rule_name`(`rule_name` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_channel_id`(`channel_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_channel_type`(`channel_type` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_datasource_id`(`datasource_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_datasource_type`(`datasource_type` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_severity`(`severity` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_receiver`(`receiver` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_status`(`status` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_business_id`(`business_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_records_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_notify_records
-- ----------------------------

-- ----------------------------
-- Table structure for alert_notify_route_rules
-- ----------------------------
DROP TABLE IF EXISTS `alert_notify_route_rules`;
CREATE TABLE `alert_notify_route_rules`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `conditions` json NULL,
  `notify_rule_id` bigint UNSIGNED NULL DEFAULT NULL,
  `notify_group_id` bigint UNSIGNED NOT NULL,
  `override_group` tinyint NOT NULL DEFAULT 0,
  `priority` int NOT NULL DEFAULT 0,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_alert_notify_route_rules_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_notify_route_rules_notify_rule_id`(`notify_rule_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_route_rules_notify_group_id`(`notify_group_id` ASC) USING BTREE,
  INDEX `idx_alert_notify_route_rules_priority`(`priority` ASC) USING BTREE,
  INDEX `idx_alert_notify_route_rules_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_notify_route_rules
-- ----------------------------

-- ----------------------------
-- Table structure for alert_notify_rules
-- ----------------------------
DROP TABLE IF EXISTS `alert_notify_rules`;
CREATE TABLE `alert_notify_rules`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `matchers` json NULL,
  `channels` json NULL,
  `receiver_groups` json NULL,
  `repeat_notify` tinyint NOT NULL DEFAULT 1,
  `repeat_mode` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'bySeverity',
  `repeat_interval_seconds` bigint NOT NULL DEFAULT 3600,
  `repeat_intervals` json NULL,
  `max_repeat_times` int NOT NULL DEFAULT 0,
  `recovery_notify` tinyint NOT NULL DEFAULT 1,
  `aggregation_enabled` tinyint NOT NULL DEFAULT 1,
  `group_wait_seconds` bigint NOT NULL DEFAULT 30,
  `group_interval_seconds` bigint NOT NULL DEFAULT 300,
  `aggregation_window_seconds` bigint NOT NULL DEFAULT 60,
  `priority` int NOT NULL DEFAULT 0,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_notify_rules_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_notify_rules_priority`(`priority` ASC) USING BTREE,
  INDEX `idx_alert_notify_rules_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_notify_rules
-- ----------------------------

-- ----------------------------
-- Table structure for alert_notify_templates
-- ----------------------------
DROP TABLE IF EXISTS `alert_notify_templates`;
CREATE TABLE `alert_notify_templates`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ident` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `channel_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `template_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'markdown',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `content` json NULL,
  `builtin` tinyint NOT NULL DEFAULT 0,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_notify_templates_ident`(`ident` ASC) USING BTREE,
  INDEX `idx_alert_notify_templates_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_notify_templates_channel_type`(`channel_type` ASC) USING BTREE,
  INDEX `idx_alert_notify_templates_template_type`(`template_type` ASC) USING BTREE,
  INDEX `idx_alert_notify_templates_builtin`(`builtin` ASC) USING BTREE,
  INDEX `idx_alert_notify_templates_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_notify_templates
-- ----------------------------
INSERT INTO `alert_notify_templates` VALUES (1, 'builtin_feishu_markdown', '飞书通知模板', 'feishu', 'markdown', '发送飞书机器人的普通文本/Markdown 告警通知模板', '{\"title\": \"{{ if .IsRecovered }}【已恢复】{{ else }}【报警中】{{ end }}{{ .RuleName }}\", \"footer\": \"AIOps 运维平台\", \"content\": \"{{ if .IsRecovered }}【已恢复】{{ else }}【报警中】{{ end }}：{{ .RuleName }}\\n【级别】: {{ .Severity }} {{ .SeverityIcon }}\\n{{ if ne .PolicyName .RuleName }}【策略】: {{ .PolicyName }}\\n{{ end }}【规则】: {{ .RuleName }}\\n{{ if gt .EventCount 1 }}【资源】: {{ .AffectedResources }}{{ else }}【目标】: {{ .ResourceName }}{{ end }}\\n{{ if .Condition }}【条件】: {{ .Condition }}{{ end }}\\n{{ if .ValueSummary }}【当前值】: {{ .ValueSummary }}{{ else if .Threshold }}【阈值】: {{ .Threshold }}{{ end }}\\n【时间】: {{ if .IsRecovered }}{{ .RecoveredAt }}{{ else }}{{ .FirstTriggeredAt }}{{ end }}\\n【持续】: {{ .Duration }}\\n{{ if gt .EventCount 1 }}【聚合】: {{ .EventCount }} 个事件 / {{ .ResourceCount }} 个资源{{ end }}\\n{{ if .Description }}【描述】: {{ .Description }}{{ end }}\\n{{ if .LabelLines }}【标签】：\\n{{ .LabelLines }}{{ end }}\\n【详情】: {{ .DetailURL }}\\n【静默】: {{ .SilenceURL }}\"}', 1, 1, 'system', 'system', '2026-07-10 23:50:49.108', '2026-11-10 21:12:16.367');
INSERT INTO `alert_notify_templates` VALUES (2, 'builtin_feishu_card', '飞书卡片通知模板', 'feishu_card', 'card', '发送飞书应用机器人的交互式卡片告警通知模板', '{\"raw\": \"{\\n  \\\"schema\\\": \\\"2.0\\\",\\n  \\\"header\\\": {\\n    \\\"template\\\": \\\"{{ if .IsRecovered }}green{{ else }}red{{ end }}\\\",\\n    \\\"title\\\": {\\n      \\\"content\\\": \\\"{{ if .IsRecovered }}【已恢复】{{ else }}【报警中】{{ end }}- {{ .RuleName }}\\\",\\n      \\\"tag\\\": \\\"plain_text\\\"\\n    }\\n  },\\n  \\\"body\\\": {\\n    \\\"elements\\\": [\\n      {\\n        \\\"tag\\\": \\\"div\\\",\\n        \\\"text\\\": {\\n          \\\"tag\\\": \\\"lark_md\\\",\\n          \\\"content\\\": \\\"**级别**: {{ .Severity }} {{ .SeverityIcon }}{{ if ne .PolicyName .RuleName }}\\\\n**策略**: {{ .PolicyName }}{{ end }}\\\\n**规则**: {{ .RuleName }}{{ if gt .EventCount 1 }}\\\\n**资源**: {{ .AffectedResources }}{{ else }}\\\\n**目标**: {{ .ResourceName }}{{ end }}{{ if .ConditionCard }}\\\\n**条件**: {{ .ConditionCard }}{{ end }}{{ if .ValueSummary }}\\\\n**当前值**: {{ .ValueSummary }}{{ else if .Threshold }}\\\\n**阈值**: {{ .Threshold }}{{ end }}\\\\n**时间**: {{ if .IsRecovered }}{{ .RecoveredAt }}{{ else }}{{ .FirstTriggeredAt }}{{ end }}\\\\n**持续**: {{ .Duration }}{{ if gt .EventCount 1 }}\\\\n**聚合**: {{ .EventCount }} 个事件 / {{ .ResourceCount }} 个资源{{ end }}{{ if .Description }}\\\\n**描述**: {{ .Description }}{{ end }}{{ if .LabelLinesCard }}\\\\n**标签**:\\\\n{{ .LabelLinesCard }}{{ end }}\\\"\\n        }\\n      }{{ if not .IsRecovered }},\\n      {\\n        \\\"tag\\\": \\\"hr\\\"\\n      },\\n      {\\n        \\\"tag\\\": \\\"column_set\\\",\\n        \\\"flex_mode\\\": \\\"none\\\",\\n        \\\"horizontal_spacing\\\": \\\"16px\\\",\\n        \\\"columns\\\": [\\n          {\\n            \\\"tag\\\": \\\"column\\\",\\n            \\\"width\\\": \\\"weighted\\\",\\n            \\\"elements\\\": [\\n              {\\n                \\\"tag\\\": \\\"button\\\",\\n                \\\"element_id\\\": \\\"claim_button_{{ .EventID }}\\\",\\n                \\\"text\\\": {\\n                  \\\"tag\\\": \\\"plain_text\\\",\\n                  \\\"content\\\": \\\"认领\\\"\\n                },\\n                \\\"type\\\": \\\"primary_filled\\\",\\n                \\\"width\\\": \\\"default\\\",\\n                \\\"behaviors\\\": [\\n                  {\\n                    \\\"type\\\": \\\"callback\\\",\\n                    \\\"value\\\": {\\n                      \\\"action\\\": \\\"claim\\\",\\n                      \\\"event_id\\\": \\\"{{ .EventID }}\\\",\\n                      \\\"event_ids\\\": \\\"{{ .EventIDs }}\\\",\\n                      \\\"fingerprint\\\": \\\"{{ .Fingerprint }}\\\"\\n                    }\\n                  }\\n                ]\\n              }\\n            ]\\n          },\\n          {\\n            \\\"tag\\\": \\\"column\\\",\\n            \\\"width\\\": \\\"weighted\\\",\\n            \\\"elements\\\": [\\n              {\\n                \\\"tag\\\": \\\"button\\\",\\n                \\\"element_id\\\": \\\"silence_button_{{ .EventID }}\\\",\\n                \\\"text\\\": {\\n                  \\\"tag\\\": \\\"plain_text\\\",\\n                  \\\"content\\\": \\\"静默一小时\\\"\\n                },\\n                \\\"type\\\": \\\"danger_filled\\\",\\n                \\\"width\\\": \\\"default\\\",\\n                \\\"behaviors\\\": [\\n                  {\\n                    \\\"type\\\": \\\"callback\\\",\\n                    \\\"value\\\": {\\n                      \\\"action\\\": \\\"silence\\\",\\n                      \\\"duration\\\": \\\"1h\\\",\\n                      \\\"event_id\\\": \\\"{{ .EventID }}\\\",\\n                      \\\"event_ids\\\": \\\"{{ .EventIDs }}\\\",\\n                      \\\"fingerprint\\\": \\\"{{ .Fingerprint }}\\\"\\n                    }\\n                  }\\n                ]\\n              }\\n            ]\\n          }\\n        ]\\n      }{{ end }}\\n    ]\\n  }\\n}\"}', 1, 1, 'system', 'system', '2026-07-10 23:50:49.108', '2026-11-10 21:12:16.367');
INSERT INTO `alert_notify_templates` VALUES (3, 'builtin_dingtalk_markdown', '钉钉通知模板', 'dingtalk', 'markdown', '发送钉钉机器人的 Markdown 告警通知模板', '{\"title\": \"{{ if .IsRecovered }}【已恢复】{{ else }}【报警中】{{ end }}{{ .RuleName }}\", \"footer\": \"AIOps 运维平台\", \"content\": \"### {{ if .IsRecovered }}【已恢复】{{ else }}【报警中】{{ end }}{{ .RuleName }}\\n- **级别**: {{ .Severity }} {{ .SeverityIcon }}\\n{{ if ne .PolicyName .RuleName }}- **策略**: {{ .PolicyName }}\\n{{ end }}- **规则**: {{ .RuleName }}\\n{{ if gt .EventCount 1 }}- **资源**: {{ .AffectedResources }}\\n{{ else }}- **目标**: {{ .ResourceName }}\\n{{ end }}{{ if .Condition }}- **条件**: {{ .Condition }}\\n{{ end }}{{ if .ValueSummary }}- **当前值**: {{ .ValueSummary }}\\n{{ else if .Threshold }}- **阈值**: {{ .Threshold }}\\n{{ end }}- **时间**: {{ if .IsRecovered }}{{ .RecoveredAt }}{{ else }}{{ .FirstTriggeredAt }}{{ end }}\\n- **持续**: {{ .Duration }}\\n{{ if gt .EventCount 1 }}- **聚合**: {{ .EventCount }} 个事件 / {{ .ResourceCount }} 个资源\\n{{ end }}{{ if .Description }}- **描述**: {{ .Description }}\\n{{ end }}{{ if .LabelLines }}- **标签**:\\n{{ .LabelLines }}\\n{{ end }}[事件详情]({{ .DetailURL }}) | [静默1小时]({{ .SilenceURL }})\"}', 1, 1, 'system', 'system', '2026-07-10 23:50:49.108', '2026-11-10 21:12:16.367');
INSERT INTO `alert_notify_templates` VALUES (4, 'builtin_wecom_markdown', '企业微信通知模板', 'wecom', 'markdown', '发送企业微信机器人的 Markdown 告警通知模板', '{\"title\": \"{{ if .IsRecovered }}【已恢复】{{ else }}【报警中】{{ end }}{{ .RuleName }}\", \"footer\": \"AIOps 运维平台\", \"content\": \"{{ if .IsRecovered }}<font color=\\\"info\\\">【已恢复】{{ .RuleName }}</font>{{ else }}<font color=\\\"warning\\\">【报警中】{{ .RuleName }}</font>{{ end }}\\n>级别: {{ .Severity }} {{ .SeverityIcon }}\\n{{ if ne .PolicyName .RuleName }}>策略: {{ .PolicyName }}\\n{{ end }}>规则: {{ .RuleName }}\\n{{ if gt .EventCount 1 }}>资源: {{ .AffectedResources }}\\n{{ else }}>目标: {{ .ResourceName }}\\n{{ end }}{{ if .Condition }}>条件: {{ .Condition }}\\n{{ end }}{{ if .ValueSummary }}>当前值: {{ .ValueSummary }}\\n{{ else if .Threshold }}>阈值: {{ .Threshold }}\\n{{ end }}>时间: {{ if .IsRecovered }}{{ .RecoveredAt }}{{ else }}{{ .FirstTriggeredAt }}{{ end }}\\n>持续: {{ .Duration }}\\n{{ if gt .EventCount 1 }}>聚合: {{ .EventCount }} 个事件 / {{ .ResourceCount }} 个资源\\n{{ end }}{{ if .Description }}>描述: {{ .Description }}\\n{{ end }}{{ if .LabelLines }}>标签:\\n{{ .LabelLines }}\\n{{ end }}[事件详情]({{ .DetailURL }}) | [静默1小时]({{ .SilenceURL }})\"}', 1, 1, 'system', 'system', '2026-07-10 23:50:49.108', '2026-11-10 21:12:16.367');
INSERT INTO `alert_notify_templates` VALUES (6, 'builtin_leliao_textcard', '乐聊通知模板', 'leliao', 'markdown', '发送公司内部乐聊通知提醒的文本卡片模板', '{\"title\": \"{{ if .IsRecovered }}【已恢复】{{ else }}【告警】{{ end }}AIOps告警中心-{{ .RuleName }}\", \"footer\": \"AIOps 运维平台\", \"content\": \"{{ if .IsRecovered }}【已恢复】{{ else }}【报警中】{{ end }}：{{ .RuleName }}\\n【级别】: {{ .Severity }} {{ .SeverityIcon }}\\n{{ if ne .PolicyName .RuleName }}【策略】: {{ .PolicyName }}\\n{{ end }}【规则】: {{ .RuleName }}\\n{{ if gt .EventCount 1 }}【资源】: {{ .AffectedResources }}{{ else }}【目标】: {{ .ResourceName }}{{ end }}\\n{{ if .Condition }}【条件】: {{ .Condition }}{{ end }}\\n{{ if .ValueSummary }}【当前值】: {{ .ValueSummary }}{{ else if .Threshold }}【阈值】: {{ .Threshold }}{{ end }}\\n【时间】: {{ if .IsRecovered }}{{ .RecoveredAt }}{{ else }}{{ .FirstTriggeredAt }}{{ end }}\\n【持续】: {{ .Duration }}\\n{{ if gt .EventCount 1 }}【聚合】: {{ .EventCount }} 个事件 / {{ .ResourceCount }} 个资源{{ end }}\\n{{ if .Description }}【描述】: {{ .Description }}{{ end }}\\n{{ if .LabelLines }}【标签】：\\n{{ .LabelLines }}{{ end }}\"}', 1, 1, 'system', 'system', '2026-08-03 01:07:46.228', '2026-11-10 21:12:16.367');

-- ----------------------------
-- Table structure for alert_rule_groups
-- ----------------------------
DROP TABLE IF EXISTS `alert_rule_groups`;
CREATE TABLE `alert_rule_groups`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` bigint UNSIGNED NULL DEFAULT 0,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `sort` int NULL DEFAULT 0,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_rule_groups_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_alert_rule_groups_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_rule_groups_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_rule_groups
-- ----------------------------
INSERT INTO `alert_rule_groups` VALUES (4, 0, '日志告警', 'VictoriaLogs LogsQL 日志告警规则', 0, 1, 'system', 'system', '2026-11-10 20:06:35.154', '2026-11-10 20:06:35.154');

-- ----------------------------
-- Table structure for alert_rules
-- ----------------------------
DROP TABLE IF EXISTS `alert_rules`;
CREATE TABLE `alert_rules`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `group_id` bigint UNSIGNED NULL DEFAULT NULL,
  `datasource_id` bigint UNSIGNED NULL DEFAULT NULL,
  `datasource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `datasource_queries` json NULL,
  `query` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `condition_operator` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `condition_value` double NULL DEFAULT NULL,
  `rule_config` json NULL,
  `append_tags` json NULL,
  `cron_pattern` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '@every 60s',
  `eval_interval_seconds` bigint NOT NULL DEFAULT 60,
  `for_duration_seconds` bigint NOT NULL DEFAULT 0,
  `recover_duration_seconds` bigint NOT NULL DEFAULT 0,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'warning',
  `enabled` tinyint NOT NULL DEFAULT 1,
  `enable_in_bg` tinyint NOT NULL DEFAULT 0,
  `time_zone` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'Asia/Shanghai',
  `labels` json NULL,
  `annotations` json NULL,
  `asset_scope` json NULL,
  `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `no_data_policy` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'ignore',
  `error_policy` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'ignore',
  `effective_time` json NULL,
  `notify_rule_ids` json NULL,
  `notify_recovered` tinyint NOT NULL DEFAULT 1,
  `notify_repeat_step_minutes` bigint NOT NULL DEFAULT 60,
  `notify_max_number` int NOT NULL DEFAULT 0,
  `pipeline_configs` json NULL,
  `runbook_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `callbacks` json NULL,
  `self_heal_config` json NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `internal` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_rules_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_rules_group_id`(`group_id` ASC) USING BTREE,
  INDEX `idx_alert_rules_datasource_id`(`datasource_id` ASC) USING BTREE,
  INDEX `idx_alert_rules_datasource_type`(`datasource_type` ASC) USING BTREE,
  INDEX `idx_alert_rules_severity`(`severity` ASC) USING BTREE,
  INDEX `idx_alert_rules_enabled`(`enabled` ASC) USING BTREE,
  INDEX `idx_alert_rules_enable_in_bg`(`enable_in_bg` ASC) USING BTREE,
  INDEX `idx_alert_rules_business_id`(`business_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_rules
-- ----------------------------

-- ----------------------------
-- Table structure for alert_silence_rules
-- ----------------------------
DROP TABLE IF EXISTS `alert_silence_rules`;
CREATE TABLE `alert_silence_rules`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `matchers` json NULL,
  `start_time` datetime(3) NULL DEFAULT NULL,
  `end_time` datetime(3) NULL DEFAULT NULL,
  `repeat_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'once',
  `reason` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_silence_rules_name`(`name` ASC) USING BTREE,
  INDEX `idx_alert_silence_rules_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_silence_rules
-- ----------------------------

-- ----------------------------
-- Table structure for alert_ssl_certificate_deploy_logs
-- ----------------------------
DROP TABLE IF EXISTS `alert_ssl_certificate_deploy_logs`;
CREATE TABLE `alert_ssl_certificate_deploy_logs`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `certificate_id` bigint UNSIGNED NOT NULL,
  `domain` varchar(253) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `host_id` bigint UNSIGNED NOT NULL,
  `host_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `host_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `auth_id` bigint UNSIGNED NOT NULL,
  `deploy_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `certificate_path` varchar(768) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `private_key_path` varchar(768) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `backup_path` varchar(768) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `status` tinyint NOT NULL,
  `log_content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_ssl_certificate_deploy_logs_certificate_id`(`certificate_id` ASC) USING BTREE,
  INDEX `idx_alert_ssl_certificate_deploy_logs_domain`(`domain` ASC) USING BTREE,
  INDEX `idx_alert_ssl_certificate_deploy_logs_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_alert_ssl_certificate_deploy_logs_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_ssl_certificate_deploy_logs
-- ----------------------------

-- ----------------------------
-- Table structure for alert_ssl_certificates
-- ----------------------------
DROP TABLE IF EXISTS `alert_ssl_certificates`;
CREATE TABLE `alert_ssl_certificates`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `domain` varchar(253) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `key_manage_id` bigint UNSIGNED NOT NULL,
  `ca_provider` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `algorithm` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(254) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint NOT NULL,
  `certificate` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `private_key` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `issuer_certificate` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `certificate_url` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `issued_at` datetime(3) NULL DEFAULT NULL,
  `expires_at` datetime(3) NULL DEFAULT NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_ssl_certificates_domain`(`domain` ASC) USING BTREE,
  INDEX `idx_alert_ssl_certificates_key_manage_id`(`key_manage_id` ASC) USING BTREE,
  INDEX `idx_alert_ssl_certificates_status`(`status` ASC) USING BTREE,
  INDEX `idx_alert_ssl_certificates_expires_at`(`expires_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of alert_ssl_certificates
-- ----------------------------

-- ----------------------------
-- Table structure for bastion_asset_authorization
-- ----------------------------
DROP TABLE IF EXISTS `bastion_asset_authorization`;
CREATE TABLE `bastion_asset_authorization`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `subject_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'授权对象类型:user/group\'',
  `subject_id` bigint UNSIGNED NOT NULL COMMENT '\'授权对象ID(用户ID或用户组ID)\'',
  `asset_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'资产类型:host/k8s_cluster等\'',
  `asset_id` bigint UNSIGNED NOT NULL COMMENT '\'资产ID\'',
  `scope_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '\'资产子范围类型，空值表示整个资产\'',
  `scope_value` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '\'资产子范围值，如K8S命名空间\'',
  `effective_time` datetime(3) NULL DEFAULT NULL COMMENT '\'生效时间\'',
  `expire_time` datetime(3) NULL DEFAULT NULL COMMENT '\'失效时间\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_asset_subject_resource`(`subject_type` ASC, `subject_id` ASC, `asset_type` ASC, `asset_id` ASC, `scope_type` ASC, `scope_value` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bastion_asset_authorization
-- ----------------------------

-- ----------------------------
-- Table structure for bastion_host_authorization
-- ----------------------------
DROP TABLE IF EXISTS `bastion_host_authorization`;
CREATE TABLE `bastion_host_authorization`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `subject_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'授权对象类型:user/group\'',
  `subject_id` bigint UNSIGNED NOT NULL COMMENT '\'授权对象ID(用户ID或用户组ID)\'',
  `host_id` bigint UNSIGNED NOT NULL COMMENT '\'主机ID\'',
  `auth_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'授权登录账号ID列表，逗号分隔\'',
  `effective_time` datetime(3) NULL DEFAULT NULL COMMENT '\'生效时间\'',
  `expire_time` datetime(3) NULL DEFAULT NULL COMMENT '\'失效时间\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  `asset_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'host' COMMENT '\'资产类型:host/k8s/log\'',
  `asset_id` bigint UNSIGNED NOT NULL COMMENT '\'资产ID:主机ID/K8S集群ID/日志项目ID\'',
  `scope` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'授权范围:host/cluster/namespace/project/logStore\'',
  `scope_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'范围资源ID:日志库ID等\'',
  `scope_value` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'范围资源值:命名空间等\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_subject_host`(`subject_type` ASC, `subject_id` ASC, `host_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bastion_host_authorization
-- ----------------------------

-- ----------------------------
-- Table structure for bastion_user_group
-- ----------------------------
DROP TABLE IF EXISTS `bastion_user_group`;
CREATE TABLE `bastion_user_group`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'用户组名称\'',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'用户组描述\'',
  `member_ids` json NULL COMMENT '\'成员ID列表(JSON数组)\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'状态:1->启用,2->禁用\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_bastion_user_group_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bastion_user_group
-- ----------------------------
INSERT INTO `bastion_user_group` VALUES (1, 'test', '', '[102]', 1, '2026-07-14 00:10:18.967', '2026-07-14 00:10:18.967');

-- ----------------------------
-- Table structure for bill_cloud_account
-- ----------------------------
DROP TABLE IF EXISTS `bill_cloud_account`;
CREATE TABLE `bill_cloud_account`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `operator` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `account_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `account_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `key_manage_id` bigint UNSIGNED NOT NULL DEFAULT 0,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `bill_sync_enabled` tinyint NOT NULL DEFAULT 1,
  `balance_sync_enabled` tinyint NOT NULL DEFAULT 1,
  `balance_sync_interval_minutes` bigint NOT NULL DEFAULT 60,
  `sync_cycle` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'daily',
  `lookback_months` bigint NOT NULL DEFAULT 12,
  `sync_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending',
  `last_sync_time` datetime(3) NULL DEFAULT NULL,
  `last_sync_error` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `current_month_cost` decimal(20, 6) NOT NULL DEFAULT 0.000000,
  `available_balance` decimal(20, 6) NOT NULL DEFAULT 0.000000,
  `currency` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'CNY',
  `balance_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'unknown',
  `last_balance_time` datetime(3) NULL DEFAULT NULL,
  `last_balance_error` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `balance_alert_rule_id` bigint UNSIGNED NOT NULL DEFAULT 0,
  `balance_alert_enabled` tinyint NOT NULL DEFAULT 0,
  `balance_alert_threshold` decimal(20, 6) NOT NULL DEFAULT 0.000000,
  `balance_alert_threshold_days` bigint NOT NULL DEFAULT 0,
  `balance_alert_on_query_error` tinyint NOT NULL DEFAULT 1,
  `balance_alert_severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'warning',
  `balance_alert_repeat_minutes` bigint NOT NULL DEFAULT 1440,
  `balance_alert_notify_recovered` tinyint NOT NULL DEFAULT 1,
  `balance_alert_notify_rule_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_bill_cloud_account_operator`(`operator` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_account_name`(`account_name` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_account_id`(`account_id` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_key_manage_id`(`key_manage_id` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_enabled`(`enabled` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_bill_sync_enabled`(`bill_sync_enabled` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_balance_sync_enabled`(`balance_sync_enabled` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_sync_status`(`sync_status` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_last_sync_time`(`last_sync_time` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_balance_status`(`balance_status` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_last_balance_time`(`last_balance_time` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_balance_alert_rule_id`(`balance_alert_rule_id` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_balance_alert_enabled`(`balance_alert_enabled` ASC) USING BTREE,
  INDEX `idx_bill_cloud_account_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bill_cloud_account
-- ----------------------------

-- ----------------------------
-- Table structure for bill_datasource
-- ----------------------------
DROP TABLE IF EXISTS `bill_datasource`;
CREATE TABLE `bill_datasource`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `operator` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `account_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `access_key` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `secret_key` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `sync_cycle` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'daily',
  `lookback_months` bigint NULL DEFAULT 3,
  `sync_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'success',
  `last_sync_time` datetime(3) NULL DEFAULT NULL,
  `current_month_cost` decimal(20, 2) NULL DEFAULT 0.00,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_bill_datasource_operator`(`operator` ASC) USING BTREE,
  INDEX `idx_bill_datasource_account_id`(`account_id` ASC) USING BTREE,
  INDEX `idx_bill_datasource_sync_status`(`sync_status` ASC) USING BTREE,
  INDEX `idx_bill_datasource_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bill_datasource
-- ----------------------------

-- ----------------------------
-- Table structure for bill_detail
-- ----------------------------
DROP TABLE IF EXISTS `bill_detail`;
CREATE TABLE `bill_detail`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `resource_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '资源ID',
  `resource_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '资源名称',
  `product_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '产品名称',
  `sub_product_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '子产品名称',
  `component_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '组件名称',
  `item_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '计费项编码',
  `usage` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用量',
  `usage_unit` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用量单位',
  `cost` decimal(20, 6) NULL DEFAULT NULL COMMENT '原价',
  `real_cost` decimal(20, 6) NULL DEFAULT NULL COMMENT '实际花费',
  `discount` decimal(10, 6) NULL DEFAULT NULL COMMENT '折扣',
  `single_price` decimal(20, 6) NULL DEFAULT NULL COMMENT '单价',
  `price_unit` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '价格单位',
  `time_span` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '时长',
  `pay_time` datetime NULL DEFAULT NULL,
  `fee_begin_time` datetime NULL DEFAULT NULL COMMENT '开始使用时间',
  `fee_end_time` datetime NULL DEFAULT NULL COMMENT '结束使用时间',
  `pay_mode` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '付费模式',
  `action_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '交易类型',
  `region` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '地域',
  `zone` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '可用区',
  `project_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '项目名称',
  `order_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '订单ID',
  `bill_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '交易ID',
  `record_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '账单记录ID',
  `operator` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '运营商',
  `month` varchar(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '账单月份',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `cloud_account_id` bigint UNSIGNED NOT NULL,
  `product_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `billing_item_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `billing_item_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `usage_quantity` decimal(24, 8) NOT NULL DEFAULT 0.00000000,
  `original_amount` decimal(20, 6) NOT NULL DEFAULT 0.000000,
  `discount_amount` decimal(20, 6) NOT NULL DEFAULT 0.000000,
  `payable_amount` decimal(20, 6) NOT NULL DEFAULT 0.000000,
  `currency` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'CNY',
  `charge_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `transaction_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `transaction_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `deducted_at` datetime(3) NULL DEFAULT NULL,
  `billing_start_at` datetime(3) NULL DEFAULT NULL,
  `billing_end_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_bill_detail_account_month_record`(`cloud_account_id` ASC, `month` ASC, `record_id` ASC) USING BTREE,
  INDEX `idx_bill_detail_resource_id`(`resource_id` ASC) USING BTREE,
  INDEX `idx_bill_detail_product_name`(`product_name` ASC) USING BTREE,
  INDEX `idx_bill_detail_item_code`(`item_code` ASC) USING BTREE,
  INDEX `idx_bill_detail_pay_time`(`pay_time` ASC) USING BTREE,
  INDEX `idx_bill_detail_region`(`region` ASC) USING BTREE,
  INDEX `idx_bill_detail_order_id`(`order_id` ASC) USING BTREE,
  INDEX `idx_bill_detail_month`(`month` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23091 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bill_detail
-- ----------------------------

-- ----------------------------
-- Table structure for bill_monthly_summary
-- ----------------------------
DROP TABLE IF EXISTS `bill_monthly_summary`;
CREATE TABLE `bill_monthly_summary`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `cloud_account_id` bigint UNSIGNED NOT NULL,
  `operator` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `month` char(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `total_cost` decimal(20, 4) NOT NULL DEFAULT 0.0000,
  `currency` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'CNY',
  `bill_status` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'estimating',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_bill_monthly_account`(`cloud_account_id` ASC, `month` ASC) USING BTREE,
  INDEX `idx_bill_monthly_summary_cloud_account_id`(`cloud_account_id` ASC) USING BTREE,
  INDEX `idx_bill_monthly_summary_operator`(`operator` ASC) USING BTREE,
  INDEX `idx_bill_monthly_summary_month`(`month` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bill_monthly_summary
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_agent
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_agent`;
CREATE TABLE `cmdb_agent`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `host_id` bigint UNSIGNED NOT NULL COMMENT '\'主机ID\'',
  `host_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'主机名称\'',
  `version` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '1.0.0' COMMENT '\'Agent版本\'',
  `status` bigint NULL DEFAULT NULL COMMENT '\'状态:1->部署中,2->部署失败,3->运行中,4->已停止\'',
  `install_path` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'安装路径\'',
  `port` bigint NULL DEFAULT 12580 COMMENT '\'监听端口\'',
  `pid` bigint NULL DEFAULT NULL COMMENT '\'进程ID\'',
  `last_heartbeat` datetime(3) NULL DEFAULT NULL COMMENT '\'最后心跳时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'错误信息\'',
  `install_progress` bigint NULL DEFAULT 0 COMMENT '\'安装进度(0-100)\'',
  `runtime_token` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_monitor_agent_host_id`(`host_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 348 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_agent
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_business_model
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_business_model`;
CREATE TABLE `cmdb_business_model`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '业务模型名称',
  `display_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '展示名称',
  `env` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '环境',
  `owner` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '负责人',
  `biz_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '业务编码',
  `properties` json NULL COMMENT '业务属性',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'active' COMMENT '状态',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_business_model_name`(`name` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_display_name`(`display_name` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_env`(`env` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_owner`(`owner` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_biz_code`(`biz_code` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_business_model
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_business_model_binding
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_business_model_binding`;
CREATE TABLE `cmdb_business_model_binding`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `business_id` bigint UNSIGNED NOT NULL COMMENT '业务模型ID',
  `parent_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '父绑定节点ID',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '节点名称',
  `display_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '展示名称',
  `resource_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '资源类型',
  `source_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '来源模块类型',
  `source_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '来源模块资源ID',
  `source_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '来源模块资源唯一键',
  `env` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '环境',
  `owner` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '负责人',
  `biz_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '业务编码',
  `properties` json NULL COMMENT '节点属性',
  `instance_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '实例唯一ID',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'active' COMMENT '状态',
  `tags` json NULL COMMENT '标签',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_business_id`(`business_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_name`(`name` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_display_name`(`display_name` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_resource_type`(`resource_type` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_source_type`(`source_type` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_source_id`(`source_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_source_key`(`source_key` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_env`(`env` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_owner`(`owner` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_biz_code`(`biz_code` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_instance_id`(`instance_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_binding_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_business_model_binding
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_business_model_health_snapshot
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_business_model_health_snapshot`;
CREATE TABLE `cmdb_business_model_health_snapshot`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `business_id` bigint UNSIGNED NOT NULL COMMENT '业务模型ID',
  `resource_id` bigint UNSIGNED NOT NULL COMMENT '绑定节点ID',
  `dimension` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '健康维度',
  `score` bigint NOT NULL DEFAULT 100 COMMENT '健康评分',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'healthy' COMMENT '健康状态',
  `source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '数据来源',
  `message` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '摘要信息',
  `evidence` json NULL COMMENT '健康证据',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_cmdb_business_health_resource_dimension`(`business_id` ASC, `resource_id` ASC, `dimension` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_health_snapshot_business_id`(`business_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_health_snapshot_resource_id`(`resource_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_health_snapshot_dimension`(`dimension` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_health_snapshot_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_business_model_health_snapshot
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_business_model_relation
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_business_model_relation`;
CREATE TABLE `cmdb_business_model_relation`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `business_id` bigint UNSIGNED NOT NULL COMMENT '业务模型ID',
  `src_id` bigint UNSIGNED NOT NULL COMMENT '源绑定节点ID',
  `dst_id` bigint UNSIGNED NOT NULL COMMENT '目标绑定节点ID',
  `relation_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '关系类型',
  `properties` json NULL COMMENT '关系属性',
  `status` tinyint NULL DEFAULT 1 COMMENT '状态 1-有效 0-无效',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_cmdb_business_relation`(`business_id` ASC, `src_id` ASC, `dst_id` ASC, `relation_type` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_relation_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_relation_business_id`(`business_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_relation_src_id`(`src_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_relation_dst_id`(`dst_id` ASC) USING BTREE,
  INDEX `idx_cmdb_business_model_relation_relation_type`(`relation_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_business_model_relation
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_db_management
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_db_management`;
CREATE TABLE `cmdb_db_management`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `instance_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `type` int NOT NULL,
  `version` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `environment` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `status` int NOT NULL DEFAULT 1,
  `auth_status` int NOT NULL DEFAULT 0,
  `account_id` bigint UNSIGNED NOT NULL,
  `group_id` bigint UNSIGNED NOT NULL,
  `host_id` bigint UNSIGNED NULL DEFAULT NULL,
  `host_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `port` int NULL DEFAULT NULL,
  `endpoint` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `region` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `idc_room_id` bigint UNSIGNED NULL DEFAULT NULL,
  `network_area` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `public_access` tinyint(1) NULL DEFAULT 0,
  `business_system` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `project` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `owner` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `dba_owner` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `ops_owner` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `department` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `cpu` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `memory` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `storage_size_gb` int NULL DEFAULT NULL,
  `used_storage_gb` int NULL DEFAULT NULL,
  `storage_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `deploy_mode` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `node_count` int NULL DEFAULT NULL,
  `cluster_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `role` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `primary_id` bigint UNSIGNED NULL DEFAULT NULL,
  `auth_method` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `ssl_enabled` tinyint(1) NULL DEFAULT 0,
  `audit_enabled` tinyint(1) NULL DEFAULT 0,
  `security_level` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `risk_tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `backup_policy` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `last_backup_time` datetime(3) NULL DEFAULT NULL,
  `monitor_status` int NULL DEFAULT 0,
  `alert_status` int NULL DEFAULT 0,
  `slow_log_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `prometheus_account_id` bigint UNSIGNED NULL DEFAULT NULL,
  `prometheus_label_selector` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `last_inspection_time` datetime(3) NULL DEFAULT NULL,
  `expire_time` datetime(3) NULL DEFAULT NULL,
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cmdb_db_management_name`(`name` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_instance_name`(`instance_name` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_type`(`type` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_environment`(`environment` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_status`(`status` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_auth_status`(`auth_status` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_account_id`(`account_id` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_group_id`(`group_id` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_host_ip`(`host_ip` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_idc_room_id`(`idc_room_id` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_business_system`(`business_system` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_owner`(`owner` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_cluster_name`(`cluster_name` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_primary_id`(`primary_id` ASC) USING BTREE,
  INDEX `idx_cmdb_db_management_prometheus_account_id`(`prometheus_account_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_db_management
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_group
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_group`;
CREATE TABLE `cmdb_group`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `parent_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '\'父级分组ID\'',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'分组名称\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `remark` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'备注\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 89 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_group
-- ----------------------------
INSERT INTO `cmdb_group` VALUES (85, 0, '默认业务组', '2026-06-21 15:14:17.299', NULL, '2026-06-21 15:14:17.299');
INSERT INTO `cmdb_group` VALUES (86, 0, '深圳科技有限公司', '2026-06-21 16:09:19.842', NULL, '2026-06-21 16:09:28.161');
INSERT INTO `cmdb_group` VALUES (87, 86, 'AI-运维业务', '2026-06-21 16:09:58.205', NULL, '2026-06-21 16:09:58.205');

-- ----------------------------
-- Table structure for cmdb_group_authorization
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_group_authorization`;
CREATE TABLE `cmdb_group_authorization`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'用户组名称\'',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'用户组描述\'',
  `member_ids` json NULL COMMENT '\'成员ID列表(JSON数组)\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'状态:1->启用,2->禁用\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_group_authorization_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_group_authorization
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_host
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_host`;
CREATE TABLE `cmdb_host`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `host_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'名称\'',
  `group_id` bigint UNSIGNED NOT NULL COMMENT '\'分组ID\'',
  `private_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'私网IP\'',
  `public_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'公网IP\'',
  `ssh_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'SSH用户名\'',
  `ssh_key_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'SSH凭据ID\'',
  `ssh_port` bigint NULL DEFAULT 22 COMMENT '\'SSH端口\'',
  `remark` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'备注\'',
  `vendor` bigint NULL DEFAULT NULL COMMENT '\'1->自建,2->阿里云,3->腾讯云,4->百度云,5->华为云,6->AWS\'',
  `region` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'区域\'',
  `instance_id` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'实例ID\'',
  `os` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'操作系统\'',
  `status` bigint NULL DEFAULT NULL COMMENT '\'状态:1->认证成功,2->未认证,3->认证失败\'',
  `cpu` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'CPU信息\'',
  `memory` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'内存信息\'',
  `disk` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'磁盘信息\'',
  `billing_type` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'计费方式\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `expire_time` datetime(3) NULL DEFAULT NULL COMMENT '\'到期时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  `ssh_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'SSH连接IP\'',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'ecs主机名称\'',
  `ssh_gateway_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'中转网关凭据ID\'',
  `tag` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'标签(格式:key=value,key=value)\'',
  `host_type` bigint NOT NULL DEFAULT 1 COMMENT '\'主机类型:1->Linux,2->Windows\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_cmdb_group_hosts`(`group_id` ASC) USING BTREE,
  CONSTRAINT `fk_cmdb_group_hosts` FOREIGN KEY (`group_id`) REFERENCES `cmdb_group` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 127 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_host
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_host_auth
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_host_auth`;
CREATE TABLE `cmdb_host_auth`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `host_id` bigint UNSIGNED NOT NULL,
  `auth_id` bigint UNSIGNED NOT NULL,
  `is_default` bigint NOT NULL DEFAULT 0,
  `create_time` datetime(3) NOT NULL,
  `update_time` datetime(3) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_cmdb_host_auth`(`host_id` ASC, `auth_id` ASC) USING BTREE,
  INDEX `idx_cmdb_host_auth_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_cmdb_host_auth_auth_id`(`auth_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 112 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_host_auth
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_host_authorization
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_host_authorization`;
CREATE TABLE `cmdb_host_authorization`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `subject_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'授权对象类型:user/group\'',
  `subject_id` bigint UNSIGNED NOT NULL COMMENT '\'授权对象ID(用户ID或用户组ID)\'',
  `host_id` bigint UNSIGNED NOT NULL COMMENT '\'主机ID\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_subject_host`(`subject_type` ASC, `subject_id` ASC, `host_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_host_authorization
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_host_inspection_run
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_host_inspection_run`;
CREATE TABLE `cmdb_host_inspection_run`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `scope_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `total` bigint NULL DEFAULT NULL,
  `completed` bigint NULL DEFAULT NULL,
  `started_by` bigint UNSIGNED NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `summary_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `results_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `error_message` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cmdb_host_inspection_run_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_host_inspection_run
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_idc_asset
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_idc_asset`;
CREATE TABLE `cmdb_idc_asset`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `room_id` bigint UNSIGNED NOT NULL,
  `cabinet_id` bigint UNSIGNED NULL DEFAULT 0,
  `cabinet_position` bigint NULL DEFAULT 0,
  `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `brand` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `model` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `cpu` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `disk` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `memory` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `purpose` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `remote_management_card` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `remark` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `purchase_date` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `warranty_start_date` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `warranty_end_date` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `usage_years` decimal(10, 2) NULL DEFAULT 0.00,
  `asset_evaluation` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` datetime(3) NOT NULL,
  `update_time` datetime(3) NOT NULL,
  `cabinet_u_size` bigint NULL DEFAULT 1,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_cmdb_idc_asset_ip`(`ip` ASC) USING BTREE,
  INDEX `idx_cmdb_idc_asset_room_id`(`room_id` ASC) USING BTREE,
  INDEX `idx_cmdb_idc_asset_cabinet_id`(`cabinet_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 154 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_idc_asset
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_idc_cabinet
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_idc_cabinet`;
CREATE TABLE `cmdb_idc_cabinet`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `room_id` bigint UNSIGNED NOT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` datetime(3) NOT NULL,
  `update_time` datetime(3) NOT NULL,
  `total_u` bigint NOT NULL DEFAULT 42,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_cmdb_idc_cabinet_room_name`(`room_id` ASC, `name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_idc_cabinet
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_idc_room
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_idc_room`;
CREATE TABLE `cmdb_idc_room`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` datetime(3) NOT NULL,
  `update_time` datetime(3) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_cmdb_idc_room_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_idc_room
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_idc_vm
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_idc_vm`;
CREATE TABLE `cmdb_idc_vm`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `physical_asset_id` bigint UNSIGNED NOT NULL,
  `vm_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `vm_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `remark` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` datetime(3) NOT NULL,
  `update_time` datetime(3) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cmdb_idc_vm_physical_asset_id`(`physical_asset_id` ASC) USING BTREE,
  INDEX `idx_cmdb_idc_vm_vm_ip`(`vm_ip` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 303 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_idc_vm
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_middleware_asset
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_middleware_asset`;
CREATE TABLE `cmdb_middleware_asset`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `version` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `environment` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `host_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `port` bigint NOT NULL,
  `host_id` bigint UNSIGNED NULL DEFAULT NULL,
  `host_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `business_system` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `owner` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `cluster_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `status` bigint NOT NULL DEFAULT 1,
  `last_check_status` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unknown',
  `last_check_at` datetime(3) NULL DEFAULT NULL,
  `last_check_latency_ms` bigint NOT NULL DEFAULT 0,
  `last_check_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `config` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_middleware_type_name`(`name` ASC, `type` ASC) USING BTREE,
  INDEX `idx_cmdb_middleware_asset_type`(`type` ASC) USING BTREE,
  INDEX `idx_cmdb_middleware_asset_environment`(`environment` ASC) USING BTREE,
  INDEX `idx_cmdb_middleware_asset_host_ip`(`host_ip` ASC) USING BTREE,
  INDEX `idx_cmdb_middleware_asset_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_cmdb_middleware_asset_business_system`(`business_system` ASC) USING BTREE,
  INDEX `idx_cmdb_middleware_asset_cluster_name`(`cluster_name` ASC) USING BTREE,
  INDEX `idx_cmdb_middleware_asset_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_middleware_asset
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_model_definition
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_model_definition`;
CREATE TABLE `cmdb_model_definition`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `type_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '资源类型标识',
  `display_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '显示名称',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '图标',
  `color` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '拓扑颜色',
  `aliases` json NULL COMMENT '别名列表',
  `fields` json NULL COMMENT '字段定义',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '描述',
  `status` tinyint NULL DEFAULT 1 COMMENT '状态',
  `sort_order` bigint NULL DEFAULT 0 COMMENT '排序',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '分组',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_model_definition_type_name`(`type_name` ASC) USING BTREE,
  INDEX `idx_cmdb_model_definition_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_model_definition
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_model_health_snapshot
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_model_health_snapshot`;
CREATE TABLE `cmdb_model_health_snapshot`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `resource_id` bigint UNSIGNED NOT NULL COMMENT '资源ID',
  `dimension` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '健康维度',
  `score` bigint NOT NULL DEFAULT 100 COMMENT '健康评分',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'healthy' COMMENT '健康状态',
  `source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '数据来源',
  `message` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '摘要信息',
  `evidence` json NULL COMMENT '健康证据',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_cmdb_model_health_resource_dimension`(`resource_id` ASC, `dimension` ASC) USING BTREE,
  INDEX `idx_cmdb_model_health_snapshot_resource_id`(`resource_id` ASC) USING BTREE,
  INDEX `idx_cmdb_model_health_snapshot_dimension`(`dimension` ASC) USING BTREE,
  INDEX `idx_cmdb_model_health_snapshot_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14535 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_model_health_snapshot
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_model_relation
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_model_relation`;
CREATE TABLE `cmdb_model_relation`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `src_id` bigint UNSIGNED NOT NULL COMMENT '源资源ID',
  `dst_id` bigint UNSIGNED NOT NULL COMMENT '目标资源ID',
  `relation_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '关系类型',
  `properties` json NULL COMMENT '关系属性',
  `status` tinyint NULL DEFAULT 1 COMMENT '状态 1-有效 0-无效',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cmdb_model_relation_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_cmdb_model_relation_src_id`(`src_id` ASC) USING BTREE,
  INDEX `idx_cmdb_model_relation_dst_id`(`dst_id` ASC) USING BTREE,
  INDEX `idx_cmdb_model_relation_relation_type`(`relation_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3973 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_model_relation
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_model_resource
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_model_resource`;
CREATE TABLE `cmdb_model_resource`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '资源名称',
  `resource_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '资源类型',
  `properties` json NULL COMMENT '所有属性',
  `instance_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '实例唯一ID（UUID/SN）允许为空',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'active' COMMENT '状态',
  `tags` json NULL COMMENT '标签',
  `display_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '展示名称',
  `source_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '来源模块类型',
  `source_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '来源模块资源ID',
  `source_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '来源模块资源唯一键',
  `env` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '环境',
  `owner` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '负责人',
  `biz_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '业务编码',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cmdb_model_resource_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_name`(`name` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_resource_type`(`resource_type` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_status`(`status` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_instance_id`(`instance_id` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_display_name`(`display_name` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_source_type`(`source_type` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_source_id`(`source_id` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_source_key`(`source_key` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_env`(`env` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_owner`(`owner` ASC) USING BTREE,
  INDEX `idx_cmdb_model_resource_biz_code`(`biz_code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1082 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_model_resource
-- ----------------------------
INSERT INTO `cmdb_model_resource` VALUES (1081, '2026-06-21 15:17:34.862', '2026-06-21 15:17:34.862', NULL, 'test', 'business_system', NULL, '', 'active', NULL, 'test', '', 0, '', 'prod', '', '');

-- ----------------------------
-- Table structure for cmdb_mysql_monitor_task
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_mysql_monitor_task`;
CREATE TABLE `cmdb_mysql_monitor_task`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `asset_id` bigint UNSIGNED NOT NULL,
  `agent_host_id` bigint UNSIGNED NOT NULL,
  `instance_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `username` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `interval_seconds` bigint NOT NULL DEFAULT 30,
  `timeout_seconds` bigint NOT NULL DEFAULT 5,
  `tls` tinyint(1) NOT NULL DEFAULT 0,
  `enabled` tinyint(1) NOT NULL DEFAULT 0,
  `config_version` bigint UNSIGNED NOT NULL DEFAULT 1,
  `collector_status` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unknown',
  `last_collect_at` datetime(3) NULL DEFAULT NULL,
  `last_success_at` datetime(3) NULL DEFAULT NULL,
  `last_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_mysql_monitor_task_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_cmdb_mysql_monitor_task_agent_host_id`(`agent_host_id` ASC) USING BTREE,
  INDEX `idx_cmdb_mysql_monitor_task_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_mysql_monitor_task
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_postgresql_monitor_task
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_postgresql_monitor_task`;
CREATE TABLE `cmdb_postgresql_monitor_task`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `asset_id` bigint UNSIGNED NOT NULL,
  `agent_host_id` bigint UNSIGNED NOT NULL,
  `instance_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `username` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `interval_seconds` bigint NOT NULL DEFAULT 30,
  `timeout_seconds` bigint NOT NULL DEFAULT 5,
  `tls` tinyint(1) NOT NULL DEFAULT 0,
  `enabled` tinyint(1) NOT NULL DEFAULT 0,
  `config_version` bigint UNSIGNED NOT NULL DEFAULT 1,
  `collector_status` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unknown',
  `last_collect_at` datetime(3) NULL DEFAULT NULL,
  `last_success_at` datetime(3) NULL DEFAULT NULL,
  `last_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_postgresql_monitor_task_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_cmdb_postgresql_monitor_task_agent_host_id`(`agent_host_id` ASC) USING BTREE,
  INDEX `idx_cmdb_postgresql_monitor_task_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_postgresql_monitor_task
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_redis_monitor_task
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_redis_monitor_task`;
CREATE TABLE `cmdb_redis_monitor_task`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `asset_id` bigint UNSIGNED NOT NULL,
  `agent_host_id` bigint UNSIGNED NOT NULL,
  `instance_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `username` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `interval_seconds` bigint NOT NULL DEFAULT 30,
  `timeout_seconds` bigint NOT NULL DEFAULT 5,
  `tls` tinyint(1) NOT NULL DEFAULT 0,
  `enabled` tinyint(1) NOT NULL DEFAULT 0,
  `config_version` bigint UNSIGNED NOT NULL DEFAULT 1,
  `collector_status` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unknown',
  `last_collect_at` datetime(3) NULL DEFAULT NULL,
  `last_success_at` datetime(3) NULL DEFAULT NULL,
  `last_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_redis_monitor_task_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_cmdb_redis_monitor_task_agent_host_id`(`agent_host_id` ASC) USING BTREE,
  INDEX `idx_cmdb_redis_monitor_task_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_redis_monitor_task
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_snmp_devices
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_snmp_devices`;
CREATE TABLE `cmdb_snmp_devices`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'primary key\'',
  `device_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'device name\'',
  `device_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'device type\'',
  `device_vendor` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'device vendor\'',
  `device_model` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'device model\'',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'device ip address\'',
  `snmp_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'external snmp monitor api url\'',
  `snmp_version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'v2c' COMMENT '\'snmp version\'',
  `snmp_community` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'snmp community\'',
  `snmp_port` bigint NULL DEFAULT 161 COMMENT '\'snmp port\'',
  `remote_port` bigint NULL DEFAULT NULL COMMENT '\'remote login port\'',
  `remote_username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'remote login username\'',
  `remote_password` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'remote login password\'',
  `remote_domain` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'remote login domain\'',
  `exporter_installed` bigint NULL DEFAULT 0 COMMENT '\'windows exporter install status\'',
  `exporter_port` bigint NULL DEFAULT 9182 COMMENT '\'windows exporter port\'',
  `exporter_version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'windows exporter version\'',
  `monitor_enabled` bigint NULL DEFAULT 1 COMMENT '\'monitor enabled\'',
  `location` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'device location\'',
  `owner` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'device owner\'',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'remark\'',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'offline' COMMENT '\'device status\'',
  `last_seen` datetime(3) NULL DEFAULT NULL COMMENT '\'last seen time\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'create time\'',
  `update_time` datetime(3) NOT NULL COMMENT '\'update time\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_snmp_devices_ip_address`(`ip_address` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_snmp_devices
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_sql
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_sql`;
CREATE TABLE `cmdb_sql`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` int NOT NULL,
  `account_id` bigint UNSIGNED NOT NULL,
  `group_id` bigint UNSIGNED NOT NULL,
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_sql
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_sql_log
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_sql_log`;
CREATE TABLE `cmdb_sql_log`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `instance_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `database` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `operation_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `sql_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `exec_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `scanned_rows` bigint NULL DEFAULT 0,
  `affected_rows` bigint NULL DEFAULT 0,
  `execution_time` bigint NULL DEFAULT 0,
  `returned_rows` bigint NULL DEFAULT 0,
  `result` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `query_time` datetime(3) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cmdb_sql_log_query_time`(`query_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cmdb_sql_log
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_windows_exporter_status
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_windows_exporter_status`;
CREATE TABLE `cmdb_windows_exporter_status`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `host_id` bigint UNSIGNED NOT NULL,
  `status` bigint NOT NULL DEFAULT 0,
  `port` bigint NOT NULL DEFAULT 9182,
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `last_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_cmdb_windows_exporter_status_host_id`(`host_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cmdb_windows_exporter_status
-- ----------------------------

-- ----------------------------
-- Table structure for config_account
-- ----------------------------
DROP TABLE IF EXISTS `config_account`;
CREATE TABLE `config_account`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `alias` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `host` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `port` bigint NOT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` bigint NOT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 45 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of config_account
-- ----------------------------
INSERT INTO `config_account` VALUES (44, 'mysql', '127.0.0.1:3306', 0, 'root', '0LRBKy+cktFQV7noo2MYGts1qRqnRBEI0brEIQpQYu8=', 3, '', '2026-07-14 01:06:26.813', '2026-07-14 01:06:26.813');

-- ----------------------------
-- Table structure for config_ecsauth
-- ----------------------------
DROP TABLE IF EXISTS `config_ecsauth`;
CREATE TABLE `config_ecsauth`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'凭证名称\'',
  `type` bigint NOT NULL COMMENT '\'认证类型:1->密码,2->私钥,3->公钥(免认证),4->跳板机\'',
  `username` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'用户名\'',
  `password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'密码(type=1时使用)\'',
  `public_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'私钥内容(type=2时使用，字段名历史原因)\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `remark` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'备注\'',
  `bastion_host` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'跳板机IP\'',
  `bastion_port` bigint NULL DEFAULT 22 COMMENT '\'跳板机端口\'',
  `bastion_username` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'跳板机用户名\'',
  `bastion_auth_type` bigint NULL DEFAULT NULL COMMENT '\'跳板机认证类型:1->密码,2->私钥,3->免密\'',
  `bastion_password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'跳板机密码\'',
  `bastion_private_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'跳板机私钥\'',
  `target_auth_type` bigint NULL DEFAULT NULL COMMENT '\'目标机器认证类型:1->密码,2->私钥,3->免密\'',
  `target_password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'目标机器密码\'',
  `target_private_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'目标机器私钥\'',
  `port` bigint NULL DEFAULT 22 COMMENT '\'端口号\'',
  `update_time` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 151 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of config_ecsauth
-- ----------------------------
INSERT INTO `config_ecsauth` VALUES (149, '免密认证', 3, 'root', '', '', '2026-07-13 21:49:25.922', '', '', 22, '', 0, '', '', 0, '', '', 22, '2026-07-13 21:49:25.922');

-- ----------------------------
-- Table structure for config_keymanage
-- ----------------------------
DROP TABLE IF EXISTS `config_keymanage`;
CREATE TABLE `config_keymanage`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `key_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `key_type` bigint NOT NULL,
  `key_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `key_secret` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of config_keymanage
-- ----------------------------

-- ----------------------------
-- Table structure for config_sync_schedule
-- ----------------------------
DROP TABLE IF EXISTS `config_sync_schedule`;
CREATE TABLE `config_sync_schedule`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `cron_expr` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `key_types` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` bigint NOT NULL DEFAULT 1,
  `last_run_time` datetime(3) NULL DEFAULT NULL,
  `next_run_time` datetime(3) NULL DEFAULT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `sync_log` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of config_sync_schedule
-- ----------------------------

-- ----------------------------
-- Table structure for db_asset_review_policy
-- ----------------------------
DROP TABLE IF EXISTS `db_asset_review_policy`;
CREATE TABLE `db_asset_review_policy`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `asset_id` bigint UNSIGNED NOT NULL,
  `environment` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `dml_rule_group_id` bigint UNSIGNED NULL DEFAULT NULL,
  `ddl_rule_group_id` bigint UNSIGNED NULL DEFAULT NULL,
  `default_approver_id` bigint UNSIGNED NULL DEFAULT NULL,
  `default_approver_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `default_executor_id` bigint UNSIGNED NULL DEFAULT NULL,
  `default_executor_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `require_backup` tinyint(1) NOT NULL DEFAULT 0,
  `require_rollback` tinyint(1) NOT NULL DEFAULT 0,
  `enabled` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_db_asset_review_policy_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_db_asset_review_policy_environment`(`environment` ASC) USING BTREE,
  INDEX `idx_db_asset_review_policy_enabled`(`enabled` ASC) USING BTREE,
  INDEX `idx_db_asset_review_policy_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_asset_review_policy
-- ----------------------------

-- ----------------------------
-- Table structure for db_operation_audit
-- ----------------------------
DROP TABLE IF EXISTS `db_operation_audit`;
CREATE TABLE `db_operation_audit`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `asset_id` bigint UNSIGNED NULL DEFAULT NULL,
  `engine` int NULL DEFAULT NULL,
  `database_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `schema_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `object_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `object_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `action` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `request_body` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `where_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `before_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `after_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `rows_affected` bigint NULL DEFAULT NULL,
  `success` tinyint(1) NULL DEFAULT NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `duration_ms` bigint NULL DEFAULT NULL,
  `operator_id` bigint UNSIGNED NULL DEFAULT NULL,
  `operator_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `client_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_db_operation_audit_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_engine`(`engine` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_database_name`(`database_name` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_schema_name`(`schema_name` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_object_name`(`object_name` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_object_type`(`object_type` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_action`(`action` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_success`(`success` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_operator_id`(`operator_id` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_operator_name`(`operator_name` ASC) USING BTREE,
  INDEX `idx_db_operation_audit_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1465 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_operation_audit
-- ----------------------------
INSERT INTO `db_operation_audit` VALUES (1464, 1, 1, '', '', '', 'asset', 'metadata', '{\"AssetID\":1,\"Level\":\"root\",\"Database\":\"\",\"Schema\":\"\",\"ObjectName\":\"\",\"ObjectType\":\"asset\",\"Keyword\":\"\"}', '', '', '', 0, 0, 'dial tcp 127.0.0.1:3306: connect: connection refused', 0, 89, 'admin', '14.127.222.121', '2026-11-10 20:51:35.940');

-- ----------------------------
-- Table structure for db_slow_log_tasks
-- ----------------------------
DROP TABLE IF EXISTS `db_slow_log_tasks`;
CREATE TABLE `db_slow_log_tasks`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `asset_id` bigint UNSIGNED NOT NULL,
  `engine` int NULL DEFAULT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `progress` int NULL DEFAULT NULL,
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `auth_id` bigint UNSIGNED NULL DEFAULT NULL,
  `request_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `result_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `slow_log_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `host_id` bigint UNSIGNED NULL DEFAULT NULL,
  `host_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `host_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `total_bytes` bigint NULL DEFAULT NULL,
  `read_bytes` bigint NULL DEFAULT NULL,
  `total_entries` bigint NULL DEFAULT NULL,
  `matched_entries` bigint NULL DEFAULT NULL,
  `unique_sql` bigint NULL DEFAULT NULL,
  `created_by` bigint UNSIGNED NULL DEFAULT NULL,
  `created_by_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `completed_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_db_slow_log_tasks_task_id`(`task_id` ASC) USING BTREE,
  INDEX `idx_db_slow_log_tasks_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_db_slow_log_tasks_engine`(`engine` ASC) USING BTREE,
  INDEX `idx_db_slow_log_tasks_status`(`status` ASC) USING BTREE,
  INDEX `idx_db_slow_log_tasks_auth_id`(`auth_id` ASC) USING BTREE,
  INDEX `idx_db_slow_log_tasks_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_db_slow_log_tasks_created_by`(`created_by` ASC) USING BTREE,
  INDEX `idx_db_slow_log_tasks_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_slow_log_tasks
-- ----------------------------

-- ----------------------------
-- Table structure for db_sql_approval_record
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_approval_record`;
CREATE TABLE `db_sql_approval_record`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint UNSIGNED NOT NULL,
  `operator_id` bigint UNSIGNED NULL DEFAULT NULL,
  `operator_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `action` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_db_sql_approval_record_ticket_id`(`ticket_id` ASC) USING BTREE,
  INDEX `idx_db_sql_approval_record_operator_id`(`operator_id` ASC) USING BTREE,
  INDEX `idx_db_sql_approval_record_operator_name`(`operator_name` ASC) USING BTREE,
  INDEX `idx_db_sql_approval_record_action`(`action` ASC) USING BTREE,
  INDEX `idx_db_sql_approval_record_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_sql_approval_record
-- ----------------------------

-- ----------------------------
-- Table structure for db_sql_execution_record
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_execution_record`;
CREATE TABLE `db_sql_execution_record`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint UNSIGNED NOT NULL,
  `statement_id` bigint UNSIGNED NULL DEFAULT NULL,
  `sql_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `rows_affected` bigint NULL DEFAULT NULL,
  `duration_ms` bigint NULL DEFAULT NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `executor_id` bigint UNSIGNED NULL DEFAULT NULL,
  `executor_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `result_columns` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `result_rows` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `exec_log` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_db_sql_execution_record_ticket_id`(`ticket_id` ASC) USING BTREE,
  INDEX `idx_db_sql_execution_record_statement_id`(`statement_id` ASC) USING BTREE,
  INDEX `idx_db_sql_execution_record_status`(`status` ASC) USING BTREE,
  INDEX `idx_db_sql_execution_record_executor_id`(`executor_id` ASC) USING BTREE,
  INDEX `idx_db_sql_execution_record_executor_name`(`executor_name` ASC) USING BTREE,
  INDEX `idx_db_sql_execution_record_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_sql_execution_record
-- ----------------------------

-- ----------------------------
-- Table structure for db_sql_review_result
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_review_result`;
CREATE TABLE `db_sql_review_result`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint UNSIGNED NOT NULL,
  `statement_id` bigint UNSIGNED NULL DEFAULT NULL,
  `rule_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `rule_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `level` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `suggestion` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `start_line` bigint NULL DEFAULT NULL,
  `end_line` bigint NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_db_sql_review_result_ticket_id`(`ticket_id` ASC) USING BTREE,
  INDEX `idx_db_sql_review_result_statement_id`(`statement_id` ASC) USING BTREE,
  INDEX `idx_db_sql_review_result_rule_code`(`rule_code` ASC) USING BTREE,
  INDEX `idx_db_sql_review_result_level`(`level` ASC) USING BTREE,
  INDEX `idx_db_sql_review_result_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 39 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_sql_review_result
-- ----------------------------

-- ----------------------------
-- Table structure for db_sql_rollback_record
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_rollback_record`;
CREATE TABLE `db_sql_rollback_record`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint UNSIGNED NOT NULL,
  `statement_id` bigint UNSIGNED NULL DEFAULT NULL,
  `statement_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `original_database` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `original_table` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `backup_database` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `backup_table` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `backup_rows` bigint NULL DEFAULT NULL,
  `rollback_sql` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `rollback_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'manual',
  `available` tinyint(1) NOT NULL DEFAULT 0,
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `binlog_start_file` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `binlog_start_pos` bigint UNSIGNED NULL DEFAULT NULL,
  `binlog_end_file` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `binlog_end_pos` bigint UNSIGNED NULL DEFAULT NULL,
  `flashback_tool` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `generated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_db_sql_rollback_record_ticket_id`(`ticket_id` ASC) USING BTREE,
  INDEX `idx_db_sql_rollback_record_statement_id`(`statement_id` ASC) USING BTREE,
  INDEX `idx_db_sql_rollback_record_statement_type`(`statement_type` ASC) USING BTREE,
  INDEX `idx_db_sql_rollback_record_original_database`(`original_database` ASC) USING BTREE,
  INDEX `idx_db_sql_rollback_record_original_table`(`original_table` ASC) USING BTREE,
  INDEX `idx_db_sql_rollback_record_backup_database`(`backup_database` ASC) USING BTREE,
  INDEX `idx_db_sql_rollback_record_backup_table`(`backup_table` ASC) USING BTREE,
  INDEX `idx_db_sql_rollback_record_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 21 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_sql_rollback_record
-- ----------------------------

-- ----------------------------
-- Table structure for db_sql_rule_group
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_rule_group`;
CREATE TABLE `db_sql_rule_group`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `engine` int NULL DEFAULT NULL,
  `environment` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT 1,
  `rules_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_db_sql_rule_group_name`(`name` ASC) USING BTREE,
  INDEX `idx_db_sql_rule_group_engine`(`engine` ASC) USING BTREE,
  INDEX `idx_db_sql_rule_group_environment`(`environment` ASC) USING BTREE,
  INDEX `idx_db_sql_rule_group_enabled`(`enabled` ASC) USING BTREE,
  INDEX `idx_db_sql_rule_group_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of db_sql_rule_group
-- ----------------------------

-- ----------------------------
-- Table structure for db_sql_statement
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_statement`;
CREATE TABLE `db_sql_statement`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint UNSIGNED NOT NULL,
  `statement_index` bigint NOT NULL,
  `sql_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `statement_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `object_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `object_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `database_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `schema_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `risk_level` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `estimated_rows` bigint NULL DEFAULT NULL,
  `review_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_db_sql_statement_ticket_id`(`ticket_id` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_statement_type`(`statement_type` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_object_name`(`object_name` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_object_type`(`object_type` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_database_name`(`database_name` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_schema_name`(`schema_name` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_risk_level`(`risk_level` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_review_status`(`review_status` ASC) USING BTREE,
  INDEX `idx_db_sql_statement_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of db_sql_statement
-- ----------------------------
INSERT INTO `db_sql_statement` VALUES (1, 1, 1, 'SELECT * FROM `autoops`.`cmdb_host` LIMIT 20', 'select', 'autoops.cmdb_host', '', 'autoops', '', 'low', 0, 'warn', '查询SQL建议包含WHERE条件', '2026-07-26 16:35:32.523');

-- ----------------------------
-- Table structure for db_sql_ticket
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_ticket`;
CREATE TABLE `db_sql_ticket`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `title` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `ticket_type` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `asset_id` bigint UNSIGNED NOT NULL,
  `use_gh_ost` tinyint(1) NOT NULL DEFAULT 0,
  `gh_ost_script` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `engine` int NULL DEFAULT NULL,
  `database_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `schema_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `environment` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `sql_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `sql_fingerprint` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `risk_level` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `backup_required` tinyint(1) NOT NULL DEFAULT 0,
  `rollback_required` tinyint(1) NOT NULL DEFAULT 0,
  `review_passed` tinyint(1) NOT NULL DEFAULT 0,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `submitter_id` bigint UNSIGNED NULL DEFAULT NULL,
  `submitter_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `submitted_at` datetime(3) NULL DEFAULT NULL,
  `approver_id` bigint UNSIGNED NULL DEFAULT NULL,
  `approver_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `reviewed_at` datetime(3) NULL DEFAULT NULL,
  `review_comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `executor_id` bigint UNSIGNED NULL DEFAULT NULL,
  `executor_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `executed_at` datetime(3) NULL DEFAULT NULL,
  `execute_comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `scheduled_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `workflow_instance_id` bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_db_sql_ticket_ticket_no`(`ticket_no` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_title`(`title` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_ticket_type`(`ticket_type` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_engine`(`engine` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_database_name`(`database_name` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_schema_name`(`schema_name` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_environment`(`environment` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_sql_fingerprint`(`sql_fingerprint` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_status`(`status` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_risk_level`(`risk_level` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_review_passed`(`review_passed` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_submitter_id`(`submitter_id` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_submitter_name`(`submitter_name` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_approver_id`(`approver_id` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_approver_name`(`approver_name` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_executor_id`(`executor_id` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_executor_name`(`executor_name` ASC) USING BTREE,
  INDEX `idx_db_sql_ticket_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of db_sql_ticket
-- ----------------------------
INSERT INTO `db_sql_ticket` VALUES (1, 'DB-QUERY-20260726163532', '查询主机资产', 'query', 1, 0, '', 1, 'autoops', '', 'prod', 'SELECT * FROM `autoops`.`cmdb_host` LIMIT 20', '6cd6aa9df520d4958572f3433242d19106c13626', 'succeeded', 'medium', 0, 0, 1, '', 89, 'admin', '2026-07-26 16:35:32.523', 89, 'admin', '2026-07-26 16:35:53.451', '', 89, 'admin', '2026-07-26 16:36:06.222', '', NULL, '2026-07-26 16:35:32.523', '2026-07-26 16:36:06.222', NULL);

-- ----------------------------
-- Table structure for host_ssh_upload_tasks
-- ----------------------------
DROP TABLE IF EXISTS `host_ssh_upload_tasks`;
CREATE TABLE `host_ssh_upload_tasks`  (
  `task_id` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `host_id` bigint UNSIGNED NULL DEFAULT NULL,
  `file_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `file_size` bigint NULL DEFAULT NULL,
  `dest_path` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `status` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `progress` double NULL DEFAULT NULL,
  `uploaded_bytes` bigint NULL DEFAULT NULL,
  `total_bytes` bigint NULL DEFAULT NULL,
  `message` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `error` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_by` bigint UNSIGNED NULL DEFAULT NULL,
  `created_by_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `completed_at` datetime(3) NULL DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of host_ssh_upload_tasks
-- ----------------------------

-- ----------------------------
-- Table structure for k8s_cluster
-- ----------------------------
DROP TABLE IF EXISTS `k8s_cluster`;
CREATE TABLE `k8s_cluster`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键ID\'',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'集群名称\'',
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'集群版本\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'集群状态:1-创建中,2-运行中,3-离线\'',
  `credential` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'集群凭证(kubeconfig)\'',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'集群描述\'',
  `cluster_type` bigint NOT NULL DEFAULT 1 COMMENT '\'集群类型:1-自建,2-导入\'',
  `node_count` bigint NULL DEFAULT 0 COMMENT '\'节点数量\'',
  `ready_nodes` bigint NULL DEFAULT 0 COMMENT '\'就绪节点数\'',
  `master_nodes` bigint NULL DEFAULT 0 COMMENT '\'Master节点数\'',
  `worker_nodes` bigint NULL DEFAULT 0 COMMENT '\'Worker节点数\'',
  `last_sync_at` datetime(3) NULL DEFAULT NULL COMMENT '\'最后同步时间\'',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `monitor_datasource_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '\'绑定监控数据源ID\'',
  `monitor_datasource_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'监控数据源类型\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_k8s_cluster_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of k8s_cluster
-- ----------------------------

-- ----------------------------
-- Table structure for k8s_cluster_inspections
-- ----------------------------
DROP TABLE IF EXISTS `k8s_cluster_inspections`;
CREATE TABLE `k8s_cluster_inspections`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `cluster_id` bigint UNSIGNED NOT NULL,
  `cluster_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `progress` bigint NULL DEFAULT 0,
  `current_step` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `completed_clusters` bigint NULL DEFAULT 0,
  `total_clusters` bigint NULL DEFAULT 1,
  `score` bigint NULL DEFAULT NULL,
  `check_count` bigint NULL DEFAULT NULL,
  `pass_count` bigint NULL DEFAULT NULL,
  `warning_count` bigint NULL DEFAULT NULL,
  `fail_count` bigint NULL DEFAULT NULL,
  `duration` bigint NULL DEFAULT NULL,
  `report_data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `start_time` datetime(3) NULL DEFAULT NULL,
  `end_time` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_k8s_cluster_inspections_cluster_id`(`cluster_id` ASC) USING BTREE,
  INDEX `idx_k8s_cluster_inspections_status`(`status` ASC) USING BTREE,
  INDEX `idx_k8s_cluster_inspections_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of k8s_cluster_inspections
-- ----------------------------

-- ----------------------------
-- Table structure for log_agent_operation
-- ----------------------------
DROP TABLE IF EXISTS `log_agent_operation`;
CREATE TABLE `log_agent_operation`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `agent_id` bigint UNSIGNED NOT NULL,
  `host_id` bigint UNSIGNED NOT NULL,
  `op_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `op_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending',
  `progress` bigint NOT NULL DEFAULT 0,
  `progress_text` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `target_version` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `config_version` bigint NOT NULL DEFAULT 0,
  `request_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `result_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `operator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_log_agent_operation_agent_id`(`agent_id` ASC) USING BTREE,
  INDEX `idx_log_agent_operation_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_log_agent_operation_op_type`(`op_type` ASC) USING BTREE,
  INDEX `idx_log_agent_operation_op_status`(`op_status` ASC) USING BTREE,
  INDEX `idx_log_agent_operation_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_agent_operation
-- ----------------------------

-- ----------------------------
-- Table structure for log_agent_status
-- ----------------------------
DROP TABLE IF EXISTS `log_agent_status`;
CREATE TABLE `log_agent_status`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `agent_id` bigint UNSIGNED NOT NULL,
  `host_id` bigint UNSIGNED NOT NULL,
  `runtime_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'host',
  `cluster_id` bigint UNSIGNED NULL DEFAULT NULL,
  `node_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `vector_version` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `vector_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'not_installed',
  `vector_p_id` bigint NULL DEFAULT NULL,
  `install_path` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `config_path` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `desired_config_version` bigint NOT NULL DEFAULT 0,
  `applied_config_version` bigint NOT NULL DEFAULT 0,
  `desired_config_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `applied_config_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `config_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending',
  `last_heartbeat_at` datetime(3) NULL DEFAULT NULL,
  `last_reload_at` datetime(3) NULL DEFAULT NULL,
  `last_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `last_error_at` datetime(3) NULL DEFAULT NULL,
  `metrics_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `vector_pid` bigint NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_agent_status_agent_id`(`agent_id` ASC) USING BTREE,
  INDEX `idx_log_agent_status_agent_id`(`agent_id` ASC) USING BTREE,
  INDEX `idx_log_agent_status_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_log_agent_status_cluster_id`(`cluster_id` ASC) USING BTREE,
  INDEX `idx_log_agent_status_node_name`(`node_name` ASC) USING BTREE,
  INDEX `idx_log_agent_status_vector_status`(`vector_status` ASC) USING BTREE,
  INDEX `idx_log_agent_status_config_status`(`config_status` ASC) USING BTREE,
  INDEX `idx_log_agent_status_last_heartbeat_at`(`last_heartbeat_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_agent_status
-- ----------------------------

-- ----------------------------
-- Table structure for log_alert_event
-- ----------------------------
DROP TABLE IF EXISTS `log_alert_event`;
CREATE TABLE `log_alert_event`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `rule_id` bigint UNSIGNED NOT NULL,
  `project_id` bigint UNSIGNED NOT NULL,
  `store_id` bigint UNSIGNED NOT NULL,
  `source_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `trigger_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `summary` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `match_count` bigint NOT NULL DEFAULT 0,
  `first_match_time` datetime(3) NULL DEFAULT NULL,
  `last_match_time` datetime(3) NULL DEFAULT NULL,
  `fired_at` datetime(3) NULL DEFAULT NULL,
  `recovered_at` datetime(3) NULL DEFAULT NULL,
  `sample_query_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `sample_event_ids_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `diagnosis_report_id` bigint UNSIGNED NULL DEFAULT NULL,
  `notify_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `notify_result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_log_alert_event_rule_id`(`rule_id` ASC) USING BTREE,
  INDEX `idx_log_alert_event_project_id`(`project_id` ASC) USING BTREE,
  INDEX `idx_log_alert_event_store_id`(`store_id` ASC) USING BTREE,
  INDEX `idx_log_alert_event_status`(`status` ASC) USING BTREE,
  INDEX `idx_log_alert_event_fired_at`(`fired_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_alert_event
-- ----------------------------

-- ----------------------------
-- Table structure for log_alert_rule
-- ----------------------------
DROP TABLE IF EXISTS `log_alert_rule`;
CREATE TABLE `log_alert_rule`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` bigint UNSIGNED NOT NULL,
  `store_id` bigint UNSIGNED NOT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `rule_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `scan_interval_minutes` bigint NOT NULL DEFAULT 5,
  `window_minutes` bigint NOT NULL DEFAULT 5,
  `threshold` bigint NOT NULL DEFAULT 1,
  `silence_minutes` bigint NOT NULL DEFAULT 30,
  `keywords_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `field_conditions_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `query_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `notify_channels_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `labels_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `annotations_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `status` bigint NOT NULL DEFAULT 1,
  `last_eval_at` datetime(3) NULL DEFAULT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_alert_rule_code`(`code` ASC) USING BTREE,
  INDEX `idx_log_alert_rule_project_id`(`project_id` ASC) USING BTREE,
  INDEX `idx_log_alert_rule_store_id`(`store_id` ASC) USING BTREE,
  INDEX `idx_log_alert_rule_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_alert_rule
-- ----------------------------

-- ----------------------------
-- Table structure for log_collect_task
-- ----------------------------
DROP TABLE IF EXISTS `log_collect_task`;
CREATE TABLE `log_collect_task`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` bigint UNSIGNED NOT NULL,
  `store_id` bigint UNSIGNED NOT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `collect_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `execute_mode` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `execute_targets_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `path_pattern` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `path_patterns_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `k8s_log_source` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `parse_template_id` bigint UNSIGNED NULL DEFAULT NULL,
  `parse_mode_snapshot` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `sample_log` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `regex_pattern` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `multiline_pattern_snapshot` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `labels_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `filters_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `config_version` bigint NOT NULL DEFAULT 1,
  `status` bigint NOT NULL DEFAULT 1,
  `last_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name_key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `config_key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `multiline_enabled_snapshot` tinyint(1) NULL DEFAULT NULL,
  `fingerprint_strategy` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `read_from` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_collect_task_code`(`code` ASC) USING BTREE,
  UNIQUE INDEX `uk_log_collect_task_name_key`(`name_key` ASC) USING BTREE,
  UNIQUE INDEX `uk_log_collect_task_config_key`(`config_key` ASC) USING BTREE,
  INDEX `idx_log_collect_task_project_id`(`project_id` ASC) USING BTREE,
  INDEX `idx_log_collect_task_store_id`(`store_id` ASC) USING BTREE,
  INDEX `idx_log_collect_task_parse_template_id`(`parse_template_id` ASC) USING BTREE,
  INDEX `idx_log_collect_task_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_collect_task
-- ----------------------------

-- ----------------------------
-- Table structure for log_k8s_cluster_collector_operation
-- ----------------------------
DROP TABLE IF EXISTS `log_k8s_cluster_collector_operation`;
CREATE TABLE `log_k8s_cluster_collector_operation`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `cluster_id` bigint UNSIGNED NOT NULL,
  `cluster_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `op_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `op_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending',
  `progress` bigint NOT NULL DEFAULT 0,
  `progress_text` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `request_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `result_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `operator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_log_k8s_cluster_collector_op_cluster_id`(`cluster_id` ASC) USING BTREE,
  INDEX `idx_log_k8s_cluster_collector_op_type`(`op_type` ASC) USING BTREE,
  INDEX `idx_log_k8s_cluster_collector_op_status`(`op_status` ASC) USING BTREE,
  INDEX `idx_log_k8s_cluster_collector_op_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_k8s_cluster_collector_operation
-- ----------------------------

-- ----------------------------
-- Table structure for log_k8s_collector_operation
-- ----------------------------
DROP TABLE IF EXISTS `log_k8s_collector_operation`;
CREATE TABLE `log_k8s_collector_operation`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_id` bigint UNSIGNED NOT NULL,
  `target_index` bigint NOT NULL DEFAULT 0,
  `cluster_id` bigint UNSIGNED NOT NULL,
  `namespace` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `workload_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `workload_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `k8_s_log_source` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `op_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `op_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending',
  `progress` bigint NOT NULL DEFAULT 0,
  `progress_text` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `request_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `result_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `operator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_log_k8s_collector_op_task_id`(`task_id` ASC) USING BTREE,
  INDEX `idx_log_k8s_collector_op_target`(`target_index` ASC) USING BTREE,
  INDEX `idx_log_k8s_collector_op_cluster_id`(`cluster_id` ASC) USING BTREE,
  INDEX `idx_log_k8s_collector_op_namespace`(`namespace` ASC) USING BTREE,
  INDEX `idx_log_k8s_collector_op_workload`(`workload_name` ASC) USING BTREE,
  INDEX `idx_log_k8s_collector_op_type`(`op_type` ASC) USING BTREE,
  INDEX `idx_log_k8s_collector_op_status`(`op_status` ASC) USING BTREE,
  INDEX `idx_log_k8s_collector_op_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_k8s_collector_operation
-- ----------------------------

-- ----------------------------
-- Table structure for log_mysql_slow_daily_report
-- ----------------------------
DROP TABLE IF EXISTS `log_mysql_slow_daily_report`;
CREATE TABLE `log_mysql_slow_daily_report`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `binding_key` varchar(160) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `asset_id` bigint UNSIGNED NOT NULL,
  `host_id` bigint UNSIGNED NOT NULL,
  `project_id` bigint UNSIGNED NOT NULL,
  `store_id` bigint UNSIGNED NOT NULL,
  `collect_task_id` bigint UNSIGNED NOT NULL,
  `report_date` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `task_token` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `progress` bigint NOT NULL DEFAULT 0,
  `analysis_total` bigint NOT NULL DEFAULT 0,
  `analysis_completed` bigint NOT NULL DEFAULT 0,
  `analysis_failed` bigint NOT NULL DEFAULT 0,
  `current_digest` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `cancel_requested` tinyint(1) NOT NULL DEFAULT 0,
  `result_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `generated_by` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `heartbeat_at` datetime(3) NULL DEFAULT NULL,
  `generated_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_mysql_slow_daily_report_binding`(`binding_key` ASC, `report_date` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_asset_id`(`asset_id` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_project_id`(`project_id` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_store_id`(`store_id` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_collect_task_id`(`collect_task_id` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_report_date`(`report_date` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_status`(`status` ASC) USING BTREE,
  INDEX `idx_log_mysql_slow_daily_report_task_token`(`task_token` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_mysql_slow_daily_report
-- ----------------------------

-- ----------------------------
-- Table structure for log_mysql_slow_explain_cache
-- ----------------------------
DROP TABLE IF EXISTS `log_mysql_slow_explain_cache`;
CREATE TABLE `log_mysql_slow_explain_cache`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `scope_key` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `digest` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `indexes_fingerprint` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `plan_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `hit_count` bigint NOT NULL DEFAULT 0,
  `last_hit_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_mysql_slow_explain_cache_entry`(`scope_key` ASC, `digest` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_mysql_slow_explain_cache
-- ----------------------------

-- ----------------------------
-- Table structure for log_mysql_slow_llm_advice_cache
-- ----------------------------
DROP TABLE IF EXISTS `log_mysql_slow_llm_advice_cache`;
CREATE TABLE `log_mysql_slow_llm_advice_cache`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `scope_key` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `digest` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `prompt_version` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `evidence_fingerprint` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `risk_level` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `optimization_source` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `analysis_model` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `analysis_provider` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `advice_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `hit_count` bigint NOT NULL DEFAULT 0,
  `last_hit_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_mysql_slow_llm_advice_cache_entry`(`scope_key` ASC, `digest` ASC, `prompt_version` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_mysql_slow_llm_advice_cache
-- ----------------------------

-- ----------------------------
-- Table structure for log_nginx_domain_binding
-- ----------------------------
DROP TABLE IF EXISTS `log_nginx_domain_binding`;
CREATE TABLE `log_nginx_domain_binding`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` bigint UNSIGNED NOT NULL,
  `store_id` bigint UNSIGNED NOT NULL,
  `domain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `display_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `status` bigint NOT NULL DEFAULT 1,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_nginx_binding_domain`(`project_id` ASC, `domain` ASC) USING BTREE,
  INDEX `idx_log_nginx_binding_project`(`project_id` ASC) USING BTREE,
  INDEX `idx_log_nginx_binding_store`(`store_id` ASC) USING BTREE,
  INDEX `idx_log_nginx_binding_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_nginx_domain_binding
-- ----------------------------

-- ----------------------------
-- Table structure for log_parse_template
-- ----------------------------
DROP TABLE IF EXISTS `log_parse_template`;
CREATE TABLE `log_parse_template`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` bigint UNSIGNED NOT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `parse_mode` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `input_format` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `multiline_enabled` tinyint(1) NOT NULL DEFAULT 0,
  `multiline_pattern` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `time_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `time_format` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `level_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `message_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `trace_id_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `span_id_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `fields_mapping_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `labels_mapping_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `vrl_script` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `status` bigint NOT NULL DEFAULT 1,
  `version` bigint NOT NULL DEFAULT 1,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_parse_template_code`(`code` ASC) USING BTREE,
  INDEX `idx_log_parse_template_project_id`(`project_id` ASC) USING BTREE,
  INDEX `idx_log_parse_template_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_parse_template
-- ----------------------------

-- ----------------------------
-- Table structure for log_project
-- ----------------------------
DROP TABLE IF EXISTS `log_project`;
CREATE TABLE `log_project`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `owner_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `owner_id` bigint UNSIGNED NULL DEFAULT NULL,
  `owner_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `env` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `status` bigint NOT NULL DEFAULT 1,
  `extra_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_project_code`(`code` ASC) USING BTREE,
  INDEX `idx_log_project_owner_id`(`owner_id` ASC) USING BTREE,
  INDEX `idx_log_project_env`(`env` ASC) USING BTREE,
  INDEX `idx_log_project_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_project
-- ----------------------------

-- ----------------------------
-- Table structure for log_store
-- ----------------------------
DROP TABLE IF EXISTS `log_store`;
CREATE TABLE `log_store`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` bigint UNSIGNED NOT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `store_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `retention_days` bigint NOT NULL DEFAULT 7,
  `status` bigint NOT NULL DEFAULT 1,
  `sort` bigint NOT NULL DEFAULT 0,
  `vlogs_tenant` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `extra_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_log_store_code`(`code` ASC) USING BTREE,
  INDEX `idx_log_store_project_id`(`project_id` ASC) USING BTREE,
  INDEX `idx_log_store_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log_store
-- ----------------------------
INSERT INTO `log_store` VALUES (1, 1, 'K8S', 'k8s_00b22a42', '', 'app', 7, 1, 1, '', '', 'admin', 'admin', '2026-07-12 21:49:43.965', '2026-07-12 21:49:43.965');

-- ----------------------------
-- Table structure for monitor_alerts_groups
-- ----------------------------
DROP TABLE IF EXISTS `monitor_alerts_groups`;
CREATE TABLE `monitor_alerts_groups`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `filters` json NULL,
  `sort_order` bigint NULL DEFAULT 0,
  `status` bigint NULL DEFAULT 1,
  `created_by` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of monitor_alerts_groups
-- ----------------------------

-- ----------------------------
-- Table structure for monitor_alerts_history
-- ----------------------------
DROP TABLE IF EXISTS `monitor_alerts_history`;
CREATE TABLE `monitor_alerts_history`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `fingerprint` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `alert_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `labels` json NULL,
  `annotations` json NULL,
  `source` enum('discover','receive') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `event_type` enum('firing','resolved','reminder') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `source_detail` json NULL,
  `source_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `is_latest_in_source` tinyint(1) NOT NULL DEFAULT 1,
  `query` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `triggered_at` datetime NULL DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_alert_history_source_fingerprint_latest_id`(`source_name` ASC, `fingerprint` ASC, `is_latest_in_source` ASC, `id` ASC) USING BTREE,
  INDEX `idx_fingerprint`(`fingerprint` ASC) USING BTREE,
  INDEX `idx_alert_name`(`alert_name` ASC) USING BTREE,
  INDEX `idx_severity`(`severity` ASC) USING BTREE,
  INDEX `idx_source`(`source` ASC) USING BTREE,
  INDEX `idx_event_type`(`event_type` ASC) USING BTREE,
  INDEX `idx_alert_history_source_event_triggered`(`source_name` ASC, `event_type` ASC, `triggered_at` ASC) USING BTREE,
  INDEX `idx_alert_history_latest_event_triggered`(`is_latest_in_source` ASC, `event_type` ASC, `triggered_at` ASC) USING BTREE,
  INDEX `idx_alert_history_source_created`(`source_name` ASC, `created_at` ASC) USING BTREE,
  INDEX `idx_alert_history_source_latest_event_triggered`(`source_name` ASC, `is_latest_in_source` ASC, `triggered_at` ASC) USING BTREE,
  INDEX `idx_triggered_at`(`triggered_at` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE,
  INDEX `idx_alert_history_source_fingerprint_id`(`source_name` ASC, `fingerprint` ASC, `id` ASC) USING BTREE,
  INDEX `idx_alert_history_fingerprint_id`(`fingerprint` ASC, `id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of monitor_alerts_history
-- ----------------------------

-- ----------------------------
-- Table structure for sys_activity_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_activity_log`;
CREATE TABLE `sys_activity_log`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `activity_type` bigint NOT NULL COMMENT '\'动态类型：1=密钥同步，2=域名检查，3=服务器巡检，4=定时任务，5=其他\'',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'动态标题\'',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'动态详细内容\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'状态：1=成功，2=失败，3=部分成功\'',
  `related_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'关联ID（如同步任务ID、域名ID等）\'',
  `summary` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'摘要信息\'',
  `duration` bigint NULL DEFAULT NULL COMMENT '\'执行耗时（秒）\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_activity_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_admin
-- ----------------------------
DROP TABLE IF EXISTS `sys_admin`;
CREATE TABLE `sys_admin`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `post_id` int NULL DEFAULT NULL COMMENT '岗位id',
  `dept_id` int NULL DEFAULT NULL COMMENT '部门id',
  `username` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '账号',
  `password` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '密码',
  `nickname` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '昵称',
  `icon` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '头像',
  `email` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '手机',
  `note` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '备注信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `status` int NOT NULL DEFAULT 1 COMMENT '帐号启用状态：1->启用,2->禁用',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE,
  INDEX `create_time`(`create_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 116 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '后台管理员表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_admin
-- ----------------------------
INSERT INTO `sys_admin` VALUES (89, 1, 5, 'admin', 'e10adc3949ba59abbe56e057f20f883e', '管理员', 'http://autoops.xin:8000/api/v1/upload/avatar/20261110/719809375.png', '123456789@qq.com', '13754354536', '后端研发', '2023-05-23 22:15:50', 1);
INSERT INTO `sys_admin` VALUES (102, 12, 3, 'test', 'e10adc3949ba59abbe56e057f20f883e', '游客', 'http://43.251.227.249:8000/api/v1/upload/avatar/20261110/659215805.png', 'zfwh1024@163.com', '13826541511', '游客', '2025-09-24 12:49:06', 1);

-- ----------------------------
-- Table structure for sys_admin_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_admin_role`;
CREATE TABLE `sys_admin_role`  (
  `admin_id` int NOT NULL COMMENT '管理员id',
  `role_id` int NOT NULL COMMENT '角色id'
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '管理员和角色关系表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_admin_role
-- ----------------------------
INSERT INTO `sys_admin_role` VALUES (0, 1);
INSERT INTO `sys_admin_role` VALUES (0, 1);
INSERT INTO `sys_admin_role` VALUES (0, 1);
INSERT INTO `sys_admin_role` VALUES (113, 1);
INSERT INTO `sys_admin_role` VALUES (114, 1);
INSERT INTO `sys_admin_role` VALUES (115, 1);
INSERT INTO `sys_admin_role` VALUES (89, 1);
INSERT INTO `sys_admin_role` VALUES (102, 14);

-- ----------------------------
-- Table structure for sys_blocking_policy
-- ----------------------------
DROP TABLE IF EXISTS `sys_blocking_policy`;
CREATE TABLE `sys_blocking_policy`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `enabled` tinyint(1) NULL DEFAULT 1,
  `priority` bigint NULL DEFAULT 0,
  `block_mode` bigint NULL DEFAULT 1,
  `enable_alias_resolution` tinyint(1) NULL DEFAULT 1,
  `enable_script_scanning` tinyint(1) NULL DEFAULT 1,
  `custom_rules` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `whitelist_cmds` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_blocking_policy
-- ----------------------------

-- ----------------------------
-- Table structure for sys_command_audit
-- ----------------------------
DROP TABLE IF EXISTS `sys_command_audit`;
CREATE TABLE `sys_command_audit`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键ID\'',
  `recording_id` bigint UNSIGNED NOT NULL COMMENT '\'录制记录ID\'',
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'会话ID\'',
  `command` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'执行的命令\'',
  `timestamp` double NOT NULL COMMENT '\'相对时间戳(秒)\'',
  `sequence` bigint NOT NULL COMMENT '\'命令序号\'',
  `is_sensitive` tinyint(1) NULL DEFAULT 0 COMMENT '\'是否为敏感命令\'',
  `risk_level` bigint NULL DEFAULT 0 COMMENT '\'风险等级:0-正常 1-可疑 2-高危\'',
  `risk_reason` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'风险原因\'',
  `execute_time` datetime(3) NOT NULL COMMENT '\'执行时间(绝对时间)\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_command_audit_recording_id`(`recording_id` ASC) USING BTREE,
  INDEX `idx_sys_command_audit_session_id`(`session_id` ASC) USING BTREE,
  INDEX `idx_sys_command_audit_is_sensitive`(`is_sensitive` ASC) USING BTREE,
  INDEX `idx_sys_command_audit_risk_level`(`risk_level` ASC) USING BTREE,
  INDEX `idx_sys_command_audit_execute_time`(`execute_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1419 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_command_audit
-- ----------------------------
INSERT INTO `sys_command_audit` VALUES (1300, 493, '1783864622_045c6cdf31eebdc5', 'ls', 2.822665946, 1, 0, 0, '', '2026-07-12 21:57:05.703', '2026-07-12 21:57:22.578');
INSERT INTO `sys_command_audit` VALUES (1301, 493, '1783864622_045c6cdf31eebdc5', 'pwd', 7.535494526, 2, 0, 0, '', '2026-07-12 21:57:10.415', '2026-07-12 21:57:22.578');
INSERT INTO `sys_command_audit` VALUES (1302, 493, '1783864622_045c6cdf31eebdc5', 'ls', 8.405747648, 3, 0, 0, '', '2026-07-12 21:57:11.286', '2026-07-12 21:57:22.578');
INSERT INTO `sys_command_audit` VALUES (1303, 495, '1783865320_5bca49db52f33df9', 'ls', 1.583180301, 1, 0, 0, '', '2026-07-12 22:08:41.613', '2026-07-12 22:09:36.173');
INSERT INTO `sys_command_audit` VALUES (1304, 495, '1783865320_5bca49db52f33df9', 'k3s kubectl create clusterrolebinding default-admin-cluster-admin \\', 12.045703169, 2, 0, 0, '', '2026-07-12 22:08:52.076', '2026-07-12 22:09:36.173');
INSERT INTO `sys_command_audit` VALUES (1305, 495, '1783865320_5bca49db52f33df9', 'k3s kubectl auth can-i delete daemonsets.apps -n aiops-vector \\', 18.598756288, 3, 0, 0, '', '2026-07-12 22:08:58.629', '2026-07-12 22:09:36.173');
INSERT INTO `sys_command_audit` VALUES (1306, 496, '1783865629_e5157f5f717f2045', 'k3s kubectl get ds aiops-vector -n aiops-vector \\', 4.673116395, 1, 0, 0, '', '2026-07-12 22:13:53.775', '2026-07-12 22:14:39.828');
INSERT INTO `sys_command_audit` VALUES (1307, 497, '1783865691_36ab1c2e578dacd0', 'crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/vector:latest-debian[Hdockercrpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/vector:latest-debian  pull', 7.731174886, 1, 1, 1, '包含敏感操作关键词', '2026-07-12 22:14:59.636', '2026-07-12 22:16:45.080');
INSERT INTO `sys_command_audit` VALUES (1308, 500, '1783866471_cfef2a4dc8d8a420', 'ls', 1.336389636, 1, 0, 0, '', '2026-07-12 22:27:53.259', '2026-07-12 22:29:05.880');
INSERT INTO `sys_command_audit` VALUES (1309, 500, '1783866471_cfef2a4dc8d8a420', 'kubectl get pods', 7.743061318, 2, 0, 0, '', '2026-07-12 22:27:59.666', '2026-07-12 22:29:05.880');
INSERT INTO `sys_command_audit` VALUES (1310, 500, '1783866471_cfef2a4dc8d8a420', 'svc', 19.066740975, 3, 0, 0, '', '2026-07-12 22:28:10.990', '2026-07-12 22:29:05.880');
INSERT INTO `sys_command_audit` VALUES (1311, 500, '1783866471_cfef2a4dc8d8a420', 'hostname -I', 34.580573766, 4, 0, 0, '', '2026-07-12 22:28:26.504', '2026-07-12 22:29:05.880');
INSERT INTO `sys_command_audit` VALUES (1312, 500, '1783866471_cfef2a4dc8d8a420', 'curl 172.31.6.35:30320', 47.104326098, 5, 1, 1, '包含敏感操作关键词', '2026-07-12 22:28:39.027', '2026-07-12 22:29:05.881');
INSERT INTO `sys_command_audit` VALUES (1313, 501, '1783866560_1d3848fabf8a88f3', 'curl  http://127.0.0.1', 68.720307256, 1, 1, 1, '包含敏感操作关键词', '2026-07-12 22:30:28.948', '2026-07-12 22:30:40.870');
INSERT INTO `sys_command_audit` VALUES (1314, 501, '1783866560_1d3848fabf8a88f3', 'ifconfig.io', 77.18500021, 2, 0, 0, '', '2026-07-12 22:30:37.413', '2026-07-12 22:30:40.870');
INSERT INTO `sys_command_audit` VALUES (1315, 503, '1783866651_bd10b30192f040f2', '[A[A[A[B[B', 5.914663595, 1, 0, 0, '', '2026-07-12 22:30:56.935', '2026-07-12 22:31:38.913');
INSERT INTO `sys_command_audit` VALUES (1316, 503, '1783866651_bd10b30192f040f2', '[A[A[A[139.9.205.38', 13.694669962, 2, 0, 0, '', '2026-07-12 22:31:04.715', '2026-07-12 22:31:38.921');
INSERT INTO `sys_command_audit` VALUES (1317, 507, '1783868481_f1da48f3bdf52c87', 'for i in $(seq 1 10); do', 2.570566858, 1, 0, 0, '', '2026-07-12 23:01:24.491', '2026-07-12 23:04:49.719');
INSERT INTO `sys_command_audit` VALUES (1318, 507, '1783868481_f1da48f3bdf52c87', '[A', 15.784858136, 2, 0, 0, '', '2026-07-12 23:01:37.705', '2026-07-12 23:04:49.719');
INSERT INTO `sys_command_audit` VALUES (1319, 507, '1783868481_f1da48f3bdf52c87', '[A', 31.033889368, 3, 0, 0, '', '2026-07-12 23:01:52.954', '2026-07-12 23:04:49.719');
INSERT INTO `sys_command_audit` VALUES (1320, 507, '1783868481_f1da48f3bdf52c87', '[A', 53.360820742, 4, 0, 0, '', '2026-07-12 23:02:15.281', '2026-07-12 23:04:49.719');
INSERT INTO `sys_command_audit` VALUES (1321, 507, '1783868481_f1da48f3bdf52c87', '[A', 127.943860936, 5, 0, 0, '', '2026-07-12 23:03:29.864', '2026-07-12 23:04:49.733');
INSERT INTO `sys_command_audit` VALUES (1322, 507, '1783868481_f1da48f3bdf52c87', '[A', 134.943868589, 6, 0, 0, '', '2026-07-12 23:03:36.864', '2026-07-12 23:04:49.733');
INSERT INTO `sys_command_audit` VALUES (1323, 506, '1783867925_7b86c665c5720cd4', 'k3s kubectl logs -n aiops-vector ds/aiops-vector --tail=100 | egrep -i \'kafka|rdkafka|healthcheck|error|timeout\'', 2.272579999, 1, 0, 0, '', '2026-07-12 22:52:07.522', '2026-07-12 23:04:49.725');
INSERT INTO `sys_command_audit` VALUES (1324, 506, '1783867925_7b86c665c5720cd4', 'k3s kubectl get cm aiops-vector -n aiops-vector \\', 30.048028841, 2, 0, 0, '', '2026-07-12 22:52:35.297', '2026-07-12 23:04:49.725');
INSERT INTO `sys_command_audit` VALUES (1325, 506, '1783867925_7b86c665c5720cd4', 'kubectl get pods', 327.449668746, 3, 0, 0, '', '2026-07-12 22:57:32.699', '2026-07-12 23:04:49.726');
INSERT INTO `sys_command_audit` VALUES (1326, 506, '1783867925_7b86c665c5720cd4', 'kubectl get log -f my-nginx-6844ccb58c-8qqh2', 345.537906618, 4, 0, 0, '', '2026-07-12 22:57:50.787', '2026-07-12 23:04:49.726');
INSERT INTO `sys_command_audit` VALUES (1327, 506, '1783867925_7b86c665c5720cd4', '[A[B[A[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[Ds', 356.45011394, 5, 0, 0, '', '2026-07-12 22:58:01.699', '2026-07-12 23:04:49.726');
INSERT INTO `sys_command_audit` VALUES (1328, 506, '1783867925_7b86c665c5720cd4', 'kubectl logs -f my-nginx-6844ccb58c-8qqh2', 401.178137404, 6, 0, 0, '', '2026-07-12 22:58:46.427', '2026-07-12 23:04:49.727');
INSERT INTO `sys_command_audit` VALUES (1329, 506, '1783867925_7b86c665c5720cd4', '[A', 464.161971363, 7, 0, 0, '', '2026-07-12 22:59:49.411', '2026-07-12 23:04:49.732');
INSERT INTO `sys_command_audit` VALUES (1330, 506, '1783867925_7b86c665c5720cd4', 'kubectl get pod', 508.30772728, 8, 0, 0, '', '2026-07-12 23:00:33.557', '2026-07-12 23:04:49.732');
INSERT INTO `sys_command_audit` VALUES (1331, 506, '1783867925_7b86c665c5720cd4', 'svc', 514.603598448, 9, 0, 0, '', '2026-07-12 23:00:39.853', '2026-07-12 23:04:49.732');
INSERT INTO `sys_command_audit` VALUES (1332, 506, '1783867925_7b86c665c5720cd4', 'kubectl logs deploy/my-nginx --tail=20', 569.146306006, 10, 0, 0, '', '2026-07-12 23:01:34.395', '2026-07-12 23:04:49.734');
INSERT INTO `sys_command_audit` VALUES (1333, 506, '1783867925_7b86c665c5720cd4', '[A', 584.427132063, 11, 0, 0, '', '2026-07-12 23:01:49.676', '2026-07-12 23:04:49.734');
INSERT INTO `sys_command_audit` VALUES (1334, 506, '1783867925_7b86c665c5720cd4', '[A', 607.401450838, 12, 0, 0, '', '2026-07-12 23:02:12.650', '2026-07-12 23:04:49.735');
INSERT INTO `sys_command_audit` VALUES (1335, 506, '1783867925_7b86c665c5720cd4', 'kubectl exec -n aiops-vector ds/aiops-vector -- sh -c \'', 679.860287628, 13, 0, 0, '', '2026-07-12 23:03:25.109', '2026-07-12 23:04:49.735');
INSERT INTO `sys_command_audit` VALUES (1336, 506, '1783867925_7b86c665c5720cd4', 'kubectl exec -n aiops-vector ds/aiops-vector -- sh -c \'', 703.847784645, 14, 0, 0, '', '2026-07-12 23:03:49.097', '2026-07-12 23:04:49.735');
INSERT INTO `sys_command_audit` VALUES (1337, 506, '1783867925_7b86c665c5720cd4', 'kubectl logs -n aiops-vector ds/aiops-vector --tail=300 | egrep -i \'k8s_stdout_file|file|kafka|rdkafka|error|warn|parse|abort\'', 716.680788622, 15, 0, 0, '', '2026-07-12 23:04:01.930', '2026-07-12 23:04:49.736');
INSERT INTO `sys_command_audit` VALUES (1338, 505, '1783867199_3f4e5bc015295334', 'ls', 24.390633254, 1, 0, 0, '', '2026-07-12 22:40:23.546', '2026-07-12 23:04:49.735');
INSERT INTO `sys_command_audit` VALUES (1339, 505, '1783867199_3f4e5bc015295334', 'docker ps', 28.686574804, 2, 0, 0, '', '2026-07-12 22:40:27.842', '2026-07-12 23:04:49.736');
INSERT INTO `sys_command_audit` VALUES (1340, 505, '1783867199_3f4e5bc015295334', 'docker inspect kafka-broker | grep KAFKA_ADVERTISED_LISTENERS', 244.423620135, 3, 0, 0, '', '2026-07-12 22:44:03.579', '2026-07-12 23:04:49.737');
INSERT INTO `sys_command_audit` VALUES (1341, 505, '1783867199_3f4e5bc015295334', 'KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://112.74.57.93:9092', 342.191468482, 4, 0, 0, '', '2026-07-12 22:45:41.346', '2026-07-12 23:04:49.738');
INSERT INTO `sys_command_audit` VALUES (1342, 505, '1783867199_3f4e5bc015295334', 'docker rm -f kafka-broker', 348.832634783, 5, 1, 1, '包含敏感操作关键词', '2026-07-12 22:45:47.988', '2026-07-12 23:04:49.738');
INSERT INTO `sys_command_audit` VALUES (1343, 505, '1783867199_3f4e5bc015295334', 'docker exec -it kafka-broker /opt/kafka/bin/kafka-topics.sh \\', 363.877518398, 6, 0, 0, '', '2026-07-12 22:46:03.033', '2026-07-12 23:04:49.738');
INSERT INTO `sys_command_audit` VALUES (1344, 505, '1783867199_3f4e5bc015295334', 'docker inspect kafka-broker | grep KAFKA_ADVERTISED_LISTENERS', 545.414167761, 7, 0, 0, '', '2026-07-12 22:49:04.569', '2026-07-12 23:04:49.739');
INSERT INTO `sys_command_audit` VALUES (1345, 505, '1783867199_3f4e5bc015295334', 'k3s kubectl get cm aiops-vector -n aiops-vector \\', 562.287158527, 8, 0, 0, '', '2026-07-12 22:49:21.442', '2026-07-12 23:04:49.739');
INSERT INTO `sys_command_audit` VALUES (1346, 505, '1783867199_3f4e5bc015295334', '.', 564.430730567, 9, 0, 0, '', '2026-07-12 22:49:23.586', '2026-07-12 23:04:49.739');
INSERT INTO `sys_command_audit` VALUES (1347, 505, '1783867199_3f4e5bc015295334', '[', 566.888508661, 10, 0, 0, '', '2026-07-12 22:49:26.044', '2026-07-12 23:04:49.739');
INSERT INTO `sys_command_audit` VALUES (1348, 505, '1783867199_3f4e5bc015295334', '[A[H[', 578.045657356, 11, 0, 0, '', '2026-07-12 22:49:37.201', '2026-07-12 23:04:49.739');
INSERT INTO `sys_command_audit` VALUES (1349, 505, '1783867199_3f4e5bc015295334', 'k3s kubectl get cm aiops-vector -n aiops-vector \\', 585.391489289, 12, 0, 0, '', '2026-07-12 22:49:44.546', '2026-07-12 23:04:49.739');
INSERT INTO `sys_command_audit` VALUES (1350, 505, '1783867199_3f4e5bc015295334', 'k3s kubectl delete pod -n aiops-vector -l app=aiops-vector', 597.846727584, 13, 0, 0, '', '2026-07-12 22:49:57.002', '2026-07-12 23:04:49.739');
INSERT INTO `sys_command_audit` VALUES (1351, 505, '1783867199_3f4e5bc015295334', 'k3s kubectl get cm aiops-vector -n aiops-vector \\', 704.191034987, 14, 0, 0, '', '2026-07-12 22:51:43.346', '2026-07-12 23:04:49.740');
INSERT INTO `sys_command_audit` VALUES (1352, 505, '1783867199_3f4e5bc015295334', 'k3s kubectl logs -n aiops-vector ds/aiops-vector --tail=100 | egrep -i \'kafka|rdkafka|healthcheck|error|timeout\'', 717.478452709, 15, 0, 0, '', '2026-07-12 22:51:56.633', '2026-07-12 23:04:49.740');
INSERT INTO `sys_command_audit` VALUES (1353, 505, '1783867199_3f4e5bc015295334', 'docker exec -it kafka-broker /opt/kafka/bin/kafka-consumer-groups.sh \\', 744.053529782, 16, 0, 0, '', '2026-07-12 22:52:23.209', '2026-07-12 23:04:49.740');
INSERT INTO `sys_command_audit` VALUES (1354, 505, '1783867199_3f4e5bc015295334', 'docker ps | grep aiops-vector-aggregator', 805.24778042, 17, 0, 0, '', '2026-07-12 22:53:24.403', '2026-07-12 23:04:49.740');
INSERT INTO `sys_command_audit` VALUES (1355, 505, '1783867199_3f4e5bc015295334', '[A', 877.870771434, 18, 0, 0, '', '2026-07-12 22:54:37.026', '2026-07-12 23:04:49.740');
INSERT INTO `sys_command_audit` VALUES (1356, 505, '1783867199_3f4e5bc015295334', 'docker rm -f aiops-vector-aggregator', 880.412181747, 19, 1, 1, '包含敏感操作关键词', '2026-07-12 22:54:39.567', '2026-07-12 23:04:49.740');
INSERT INTO `sys_command_audit` VALUES (1357, 505, '1783867199_3f4e5bc015295334', '[A[A[A[A[A[A[B', 886.958437423, 20, 0, 0, '', '2026-07-12 22:54:46.113', '2026-07-12 23:04:49.741');
INSERT INTO `sys_command_audit` VALUES (1358, 505, '1783867199_3f4e5bc015295334', '[A[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D f', 905.877799016, 21, 0, 0, '', '2026-07-12 22:55:05.033', '2026-07-12 23:04:49.741');
INSERT INTO `sys_command_audit` VALUES (1359, 505, '1783867199_3f4e5bc015295334', 'docker exec -it kafka-broker /opt/kafka/bin/kafka-consumer-groups.sh \\', 924.791272991, 22, 0, 0, '', '2026-07-12 22:55:23.946', '2026-07-12 23:04:49.741');
INSERT INTO `sys_command_audit` VALUES (1360, 505, '1783867199_3f4e5bc015295334', 'docker exec -it kafka-broker /opt/kafka/bin/kafka-consumer-groups.sh \\', 924.915889046, 23, 0, 0, '', '2026-07-12 22:55:24.071', '2026-07-12 23:04:49.741');
INSERT INTO `sys_command_audit` VALUES (1361, 505, '1783867199_3f4e5bc015295334', '[A[A', 935.214417297, 24, 0, 0, '', '2026-07-12 22:55:34.369', '2026-07-12 23:04:49.741');
INSERT INTO `sys_command_audit` VALUES (1362, 505, '1783867199_3f4e5bc015295334', 'docker exec -it kafka-broker /opt/kafka/bin/kafka-consumer-groups.sh \\', 974.732542029, 25, 0, 0, '', '2026-07-12 22:56:13.888', '2026-07-12 23:04:49.741');
INSERT INTO `sys_command_audit` VALUES (1363, 505, '1783867199_3f4e5bc015295334', 'docker exec -it kafka-broker /opt/kafka/bin/kafka-console-consumer.sh \\', 1035.790958575, 26, 0, 0, '', '2026-07-12 22:57:14.946', '2026-07-12 23:04:49.741');
INSERT INTO `sys_command_audit` VALUES (1364, 505, '1783867199_3f4e5bc015295334', '[A', 1168.494847528, 27, 0, 0, '', '2026-07-12 22:59:27.650', '2026-07-12 23:04:49.742');
INSERT INTO `sys_command_audit` VALUES (1365, 505, '1783867199_3f4e5bc015295334', '[A', 1326.973999508, 28, 0, 0, '', '2026-07-12 23:02:06.129', '2026-07-12 23:04:49.742');
INSERT INTO `sys_command_audit` VALUES (1366, 508, '1783868704_ea1cde6cd8cc01d7', '[A', 3.077660416, 1, 0, 0, '', '2026-07-12 23:05:08.058', '2026-07-12 23:06:35.147');
INSERT INTO `sys_command_audit` VALUES (1367, 508, '1783868704_ea1cde6cd8cc01d7', '[A[A', 6.717615015, 2, 0, 0, '', '2026-07-12 23:05:11.698', '2026-07-12 23:06:35.148');
INSERT INTO `sys_command_audit` VALUES (1368, 508, '1783868704_ea1cde6cd8cc01d7', '[A[A', 22.75014956, 3, 0, 0, '', '2026-07-12 23:05:27.730', '2026-07-12 23:06:35.149');
INSERT INTO `sys_command_audit` VALUES (1369, 511, '1783869848_f215a702e0fece56', 'kubectl exec -n aiops-vector ds/aiops-vector -- vector validate /etc/vector/vector.yaml', 2.894735502, 1, 0, 0, '', '2026-07-12 23:24:11.365', '2026-07-12 23:27:35.008');
INSERT INTO `sys_command_audit` VALUES (1370, 510, '1783869836_9d692a20fe9af254', 'kubectl exec -n aiops-vector ds/aiops-vector -- vector validate /etc/vector/vector.yaml', 3.145290772, 1, 0, 0, '', '2026-07-12 23:23:59.818', '2026-07-12 23:27:35.008');
INSERT INTO `sys_command_audit` VALUES (1371, 512, '1783869857_67f4a076e3850f31', '[A[A', 2.708907725, 1, 0, 0, '', '2026-07-12 23:24:20.270', '2026-07-12 23:27:35.059');
INSERT INTO `sys_command_audit` VALUES (1372, 512, '1783869857_67f4a076e3850f31', '[A[A', 16.331322704, 2, 0, 0, '', '2026-07-12 23:24:33.892', '2026-07-12 23:27:35.059');
INSERT INTO `sys_command_audit` VALUES (1373, 512, '1783869857_67f4a076e3850f31', 'kubectl exec -n aiops-vector ds/aiops-vector -- sh -c \'', 28.688533014, 3, 0, 0, '', '2026-07-12 23:24:46.250', '2026-07-12 23:27:35.059');
INSERT INTO `sys_command_audit` VALUES (1374, 520, '1784127518_d3b0e61902506324', 'docker ps', 9.297641667, 1, 0, 0, '', '2026-07-15 22:58:47.774', '2026-07-15 22:59:27.568');
INSERT INTO `sys_command_audit` VALUES (1375, 520, '1784127518_d3b0e61902506324', '[A[A', 14.362020525, 2, 0, 0, '', '2026-07-15 22:58:52.838', '2026-07-15 22:59:27.569');
INSERT INTO `sys_command_audit` VALUES (1376, 520, '1784127518_d3b0e61902506324', '[15;5Rls', 18.449199436, 3, 0, 0, '', '2026-07-15 22:58:56.925', '2026-07-15 22:59:27.569');
INSERT INTO `sys_command_audit` VALUES (1377, 520, '1784127518_d3b0e61902506324', '[20;5Rps -ef', 21.754481388, 4, 0, 0, '', '2026-07-15 22:59:00.230', '2026-07-15 22:59:27.569');
INSERT INTO `sys_command_audit` VALUES (1378, 519, '1784127431_3386d9eab93f3680', 'ls', 1.740415161, 1, 0, 0, '', '2026-07-15 22:57:13.112', '2026-07-15 22:59:27.568');
INSERT INTO `sys_command_audit` VALUES (1379, 519, '1784127431_3386d9eab93f3680', 'ls', 3.308146154, 2, 0, 0, '', '2026-07-15 22:57:14.680', '2026-07-15 22:59:27.569');
INSERT INTO `sys_command_audit` VALUES (1380, 519, '1784127431_3386d9eab93f3680', 'kubectl get nodesl', 13.844595074, 3, 0, 0, '', '2026-07-15 22:57:25.217', '2026-07-15 22:59:27.569');
INSERT INTO `sys_command_audit` VALUES (1381, 519, '1784127431_3386d9eab93f3680', '[', 17.622116682, 4, 0, 0, '', '2026-07-15 22:57:28.994', '2026-07-15 22:59:27.575');
INSERT INTO `sys_command_audit` VALUES (1382, 519, '1784127431_3386d9eab93f3680', 'pods   --all', 32.396786367, 5, 0, 0, '', '2026-07-15 22:57:43.769', '2026-07-15 22:59:27.575');
INSERT INTO `sys_command_audit` VALUES (1383, 519, '1784127431_3386d9eab93f3680', 'kubectl get pods -A', 134.012342698, 6, 0, 0, '', '2026-07-15 22:59:25.384', '2026-07-15 22:59:27.576');
INSERT INTO `sys_command_audit` VALUES (1384, 521, '1784127575_68d833e663c85bce', '[A', 271.362621623, 1, 0, 0, '', '2026-07-15 23:04:06.729', '2026-07-15 23:20:18.039');
INSERT INTO `sys_command_audit` VALUES (1385, 521, '1784127575_68d833e663c85bce', 'nodes', 297.937702117, 2, 0, 0, '', '2026-07-15 23:04:33.304', '2026-07-15 23:20:18.039');
INSERT INTO `sys_command_audit` VALUES (1386, 521, '1784127575_68d833e663c85bce', 'apt  -y install  helm repo update', 1141.567040851, 3, 0, 0, '', '2026-07-15 23:18:36.933', '2026-07-15 23:20:18.039');
INSERT INTO `sys_command_audit` VALUES (1387, 521, '1784127575_68d833e663c85bce', 'lldx', 1238.298600137, 4, 0, 0, '', '2026-07-15 23:20:13.665', '2026-07-15 23:20:18.040');
INSERT INTO `sys_command_audit` VALUES (1388, 524, '1784129732_6dd0f99f22500852', 'kubectl get nodes', 8.280205389, 1, 0, 0, '', '2026-07-15 23:35:41.035', '2026-07-16 00:20:57.149');
INSERT INTO `sys_command_audit` VALUES (1389, 524, '1784129732_6dd0f99f22500852', '[A[pod  -A', 2207.554619386, 2, 0, 0, '', '2026-07-16 00:12:20.310', '2026-07-16 00:20:57.149');
INSERT INTO `sys_command_audit` VALUES (1390, 524, '1784129732_6dd0f99f22500852', '[A', 2348.860612502, 3, 0, 0, '', '2026-07-16 00:14:41.616', '2026-07-16 00:20:57.149');
INSERT INTO `sys_command_audit` VALUES (1391, 524, '1784129732_6dd0f99f22500852', 'kubectl -n monitoring get pods,pvc -o wide', 2380.058863694, 4, 0, 0, '', '2026-07-16 00:15:12.814', '2026-07-16 00:20:57.149');
INSERT INTO `sys_command_audit` VALUES (1392, 524, '1784129732_6dd0f99f22500852', 'docdoc  pull registry.cn-hangzhou.aliyuncs.com/rancher/mirrored-metrics-server:v0.8.1', 2495.759017618, 5, 1, 1, '包含敏感操作关键词', '2026-07-16 00:17:08.514', '2026-07-16 00:20:57.149');
INSERT INTO `sys_command_audit` VALUES (1393, 524, '1784129732_6dd0f99f22500852', 'docdockerdocker                         dockerd                        dockerd-rootless-setuptool.sh  dockerd-rootless.sh            docker-proxyroot@hw-ops:~# docker pull crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/node-exporter:v1.8.2', 2681.910263449, 6, 1, 1, '包含敏感操作关键词', '2026-07-16 00:20:14.665', '2026-07-16 00:20:57.149');
INSERT INTO `sys_command_audit` VALUES (1394, 525, '1784132830_0a93fd90ece4cd5a', 'free -h', 262.832103538, 1, 0, 0, '', '2026-07-16 00:31:33.752', '2026-07-16 01:07:45.659');
INSERT INTO `sys_command_audit` VALUES (1395, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring set resources deployment/kube-state-metrics \\', 334.388696662, 2, 0, 0, '', '2026-07-16 00:32:45.309', '2026-07-16 01:07:45.659');
INSERT INTO `sys_command_audit` VALUES (1396, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring set resources deployment/kube-state-metrics \\', 580.326129109, 3, 0, 0, '', '2026-07-16 00:36:51.246', '2026-07-16 01:07:45.659');
INSERT INTO `sys_command_audit` VALUES (1397, 525, '1784132830_0a93fd90ece4cd5a', 'free -h', 607.42466936, 4, 0, 0, '', '2026-07-16 00:37:18.345', '2026-07-16 01:07:45.659');
INSERT INTO `sys_command_audit` VALUES (1398, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl top pods -A', 666.651158845, 5, 0, 0, '', '2026-07-16 00:38:17.571', '2026-07-16 01:07:45.659');
INSERT INTO `sys_command_audit` VALUES (1399, 525, '1784132830_0a93fd90ece4cd5a', 'curl -sS -G http://127.0.0.1:8428/api/v1/query \\', 687.818396973, 6, 1, 1, '包含敏感操作关键词', '2026-07-16 00:38:38.738', '2026-07-16 01:07:45.659');
INSERT INTO `sys_command_audit` VALUES (1400, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring logs deployment/vmagent --tail=100 | grep -iE \'error|failed\'', 701.398446033, 7, 0, 0, '', '2026-07-16 00:38:52.318', '2026-07-16 01:07:45.660');
INSERT INTO `sys_command_audit` VALUES (1401, 525, '1784132830_0a93fd90ece4cd5a', '[Acurl -sS -G http://112.74.57.93:8428/api/v1/query \\', 794.237102463, 8, 1, 1, '包含敏感操作关键词', '2026-07-16 00:40:25.157', '2026-07-16 01:07:45.660');
INSERT INTO `sys_command_audit` VALUES (1402, 525, '1784132830_0a93fd90ece4cd5a', '[A', 873.869625986, 9, 0, 0, '', '2026-07-16 00:41:44.790', '2026-07-16 01:07:45.660');
INSERT INTO `sys_command_audit` VALUES (1403, 525, '1784132830_0a93fd90ece4cd5a', 'ls', 887.556872615, 10, 0, 0, '', '2026-07-16 00:41:58.477', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1404, 525, '1784132830_0a93fd90ece4cd5a', 'vim test.yaml', 890.614068842, 11, 0, 0, '', '2026-07-16 00:42:01.534', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1405, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl apply -f test.yaml', 907.857387547, 12, 0, 0, '', '2026-07-16 00:42:18.777', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1406, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring rollout restart deployment/vmagent', 919.687987688, 13, 0, 0, '', '2026-07-16 00:42:30.608', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1407, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring logs deployment/vmagent --since=2m \\', 950.016499238, 14, 1, 1, '包含敏感操作关键词', '2026-07-16 00:43:00.936', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1408, 525, '1784132830_0a93fd90ece4cd5a', 'curl -sS -G http://127.0.0.1:8428/api/v1/query \\', 962.818122239, 15, 1, 1, '包含敏感操作关键词', '2026-07-16 00:43:13.738', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1409, 525, '1784132830_0a93fd90ece4cd5a', '[A[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[192.168.3.7', 982.443515889, 16, 0, 0, '', '2026-07-16 00:43:33.364', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1410, 525, '1784132830_0a93fd90ece4cd5a', '[A[D[D[D[D[D[H[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[C[112.74.57.93', 1010.468126196, 17, 0, 0, '', '2026-07-16 00:44:01.388', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1411, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring logs deployment/vmagent --since=2m \\', 1038.061576816, 18, 1, 1, '包含敏感操作关键词', '2026-07-16 00:44:28.982', '2026-07-16 01:07:45.661');
INSERT INTO `sys_command_audit` VALUES (1412, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring get configmap vmagent-config \\', 1101.344306426, 19, 0, 0, '', '2026-07-16 00:45:32.264', '2026-07-16 01:07:45.662');
INSERT INTO `sys_command_audit` VALUES (1413, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring rollout restart deployment/vmagent', 1346.10068562, 20, 0, 0, '', '2026-07-16 00:49:37.021', '2026-07-16 01:07:45.662');
INSERT INTO `sys_command_audit` VALUES (1414, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring patch deployment vmagent \\', 1511.760886971, 21, 0, 0, '', '2026-07-16 00:52:22.681', '2026-07-16 01:07:45.662');
INSERT INTO `sys_command_audit` VALUES (1415, 525, '1784132830_0a93fd90ece4cd5a', 'kubectl -n monitoring rollout restart deployment/vmagent', 1520.981736457, 22, 0, 0, '', '2026-07-16 00:52:31.902', '2026-07-16 01:07:45.662');
INSERT INTO `sys_command_audit` VALUES (1416, 526, '1784437387_af2414ea8f1e2a5b', 'df -h', 11.316920567, 1, 0, 0, '', '2026-07-19 13:03:18.597', '2026-07-19 13:03:26.236');
INSERT INTO `sys_command_audit` VALUES (1417, 527, '1784438766_121fffc6b16f97e6', 'baidu-ops[D[D[D[D[D[D[D[D[D[D[D[D[Dhost   [D[D[Dnamec[2@tl set-hostname', 24.212106127, 1, 0, 0, '', '2026-07-19 13:26:31.165', '2026-07-19 13:26:39.810');
INSERT INTO `sys_command_audit` VALUES (1418, 527, '1784438766_121fffc6b16f97e6', 'hostname', 29.450596585, 2, 0, 0, '', '2026-07-19 13:26:36.404', '2026-07-19 13:26:39.810');

-- ----------------------------
-- Table structure for sys_command_blocking
-- ----------------------------
DROP TABLE IF EXISTS `sys_command_blocking`;
CREATE TABLE `sys_command_blocking`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `recording_id` bigint UNSIGNED NULL DEFAULT NULL,
  `original_cmd` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `resolved_cmd` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `blocked` tinyint(1) NULL DEFAULT NULL,
  `block_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `risk_level` bigint NULL DEFAULT NULL,
  `is_alias` tinyint(1) NULL DEFAULT 0,
  `is_script` tinyint(1) NULL DEFAULT 0,
  `detected_issues` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `policy_id` bigint UNSIGNED NULL DEFAULT NULL,
  `policy_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `host_id` bigint UNSIGNED NULL DEFAULT NULL,
  `host_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_command_blocking_session_id`(`session_id` ASC) USING BTREE,
  INDEX `idx_sys_command_blocking_recording_id`(`recording_id` ASC) USING BTREE,
  INDEX `idx_sys_command_blocking_blocked`(`blocked` ASC) USING BTREE,
  INDEX `idx_sys_command_blocking_risk_level`(`risk_level` ASC) USING BTREE,
  INDEX `idx_sys_command_blocking_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_sys_command_blocking_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_sys_command_blocking_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_command_blocking
-- ----------------------------

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `config_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'配置键\'',
  `config_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'配置类型(ldap,email,sms等)\'',
  `config_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'配置数据(JSON格式)\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'状态:1->启用,2->禁用\'',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'备注\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NOT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_config_config_key`(`config_key` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `sys_config` VALUES (1, 'ldap', 'ldap', '{\"enable\":true,\"host\":\"ldap-out.dding.net\",\"port\":389,\"baseDn\":\"dc=dding,dc=cn\",\"bindUser\":\"cn=reader,dc=dding,dc=cn\",\"bindPass\":\"fsyunding2018\",\"authFilter\":\"(\\u0026(cn=%s))\",\"coverAttributes\":true,\"tls\":false,\"startTLS\":false,\"defaultRoles\":null,\"defaultRoleId\":13,\"attributes\":{\"nickname\":\"cn\",\"phone\":\"mobile\",\"email\":\"mail\"},\"remark\":\"\"}', 1, '', '2025-12-09 13:19:56.671', '2025-12-09 14:35:37.377');
INSERT INTO `sys_config` VALUES (2, 'license', 'license', '{\"activated\":true,\"licenseKey\":\"eyJjdXN0b21lck5hbWUiOiJvcHMiLCJtYWNoaW5lQ29kZSI6IjBjZjlkZTBkZmMwMmIxZTQ1OTlkNGU3MGExMzQ3MjEwYWFhMzY5Y2I2MDhiMTNhODFlZWUyNmZjMTcyZWFiNTciLCJsaWNlbnNlVHlwZSI6ImFubnVhbCIsImlzc3VlZEF0IjoiMjAyNi0wMy0yMyAxNzoxMjowMCIsImV4cGlyZUF0IjoiMjA5OS0wMy0yMyAxNzoyMDowMCJ9.nmfhYO2APkza54nAVHUKoBMR7APotObX+5SgsYN+cBniSNbY5mdV/wW4k9FxPY3qm5qxHc4OviZzhS2rwWu24fRjhUGL+g7yIl+5WyZ2HcHHcEn+JN+fJflytmgdf0lGuxGEgI46ADQUuzc/+m7xixxfxH6YibVnKBwyQsruE+IVjyUpRbbp/xYkEYvCa09ej6UmDdJiaHWSre4lLPi9MMWrgvApXluD+7W8GfUuyxfx3ceAw1+ahabVRauRJqUhIMkhSZ6LolSKi0ADi1Gh6mlMLtL35YtI2bYXgsocLOpVSkd4CCGQ4ZyyGCg+rNaP2ASvXTOnxddUOIbJfrt1Ug==\",\"customerName\":\"ops\",\"machineCode\":\"0cf9de0dfc02b1e4599d4e70a1347210aaa369cb608b13a81eee26fc172eab57\",\"licenseType\":\"annual\",\"issuedAt\":\"2026-03-23 17:12:00\",\"activatedAt\":\"2026-03-23 17:22:56\",\"expireAt\":\"2099-03-23 17:20:00\"}', 1, '系统激活信息', '2026-03-23 16:51:47.717', '2026-03-23 17:22:56.703');

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` int NOT NULL COMMENT '父id',
  `dept_type` int NOT NULL COMMENT '部门类型（1->公司, 2->中心，3->部门）',
  `dept_name` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '部门名称',
  `dept_status` int NOT NULL DEFAULT 1 COMMENT '部门状态（1->正常 2->停用）',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `dept_name`(`dept_name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_bin COMMENT = '部门表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_dept
-- ----------------------------
INSERT INTO `sys_dept` VALUES (1, 0, 1, '神州科技有限公司', 1, '2023-06-14 17:53:23');
INSERT INTO `sys_dept` VALUES (2, 1, 2, '深圳研发中心', 1, '2023-06-14 17:53:55');
INSERT INTO `sys_dept` VALUES (3, 2, 3, '运维部', 1, '2023-06-14 17:54:15');
INSERT INTO `sys_dept` VALUES (5, 2, 3, '研发部', 1, '2023-06-14 17:55:25');

-- ----------------------------
-- Table structure for sys_login_info
-- ----------------------------
DROP TABLE IF EXISTS `sys_login_info`;
CREATE TABLE `sys_login_info`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '用户账号',
  `ip_address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '登录IP地址',
  `login_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '登录地点',
  `browser` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '浏览器类型',
  `os` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '操作系统',
  `login_status` int NULL DEFAULT NULL COMMENT '登录状态（1-成功 2-失败）',
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '提示消息',
  `login_time` datetime NULL DEFAULT NULL COMMENT '访问时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '登录日志记录' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_login_info
-- ----------------------------
INSERT INTO `sys_login_info` VALUES (1, 'admin', '58.152.20.185', '香港 电讯盈科有限公司', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:12:09');
INSERT INTO `sys_login_info` VALUES (2, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:25:50');
INSERT INTO `sys_login_info` VALUES (3, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:25:50');
INSERT INTO `sys_login_info` VALUES (4, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:25:50');
INSERT INTO `sys_login_info` VALUES (5, 'test', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:40:14');
INSERT INTO `sys_login_info` VALUES (6, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:40:39');
INSERT INTO `sys_login_info` VALUES (7, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 2, '密码不正确', '2026-11-10 20:41:32');
INSERT INTO `sys_login_info` VALUES (8, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:49:21');
INSERT INTO `sys_login_info` VALUES (9, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 20:50:21');
INSERT INTO `sys_login_info` VALUES (10, 'admin', '14.127.222.121', '广东省深圳市 电信', 'Chrome/156.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 21:05:39');
INSERT INTO `sys_login_info` VALUES (11, 'admin', '172.16.11.112', '局域网', 'Chrome/135.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 23:03:08');
INSERT INTO `sys_login_info` VALUES (12, 'test', '172.16.11.112', '局域网', 'Chrome/135.0.0.0', 'Windows 10', 1, '登录成功', '2026-11-10 23:03:20');

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` int NULL DEFAULT NULL COMMENT '父级菜单id',
  `menu_name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '菜单名称',
  `icon` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '图标',
  `value` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '接口权限值',
  `menu_type` int NULL DEFAULT NULL COMMENT '菜单类型：1->目录；2->菜单；3->按钮（接口绑定权限）',
  `url` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '菜单url',
  `menu_status` int NULL DEFAULT 2 COMMENT '启用状态；1->禁用；2->启用',
  `sort` int NULL DEFAULT NULL COMMENT '排序',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 416 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (4, 0, '系统管理', 'StarFilled', '', 1, '', 2, 12, '2022-09-04 13:57:39');
INSERT INTO `sys_menu` VALUES (6, 4, '用户信息', 'Avatar', 'base:admin:list', 2, 'system/admin', 2, 1, '2022-09-04 13:59:39');
INSERT INTO `sys_menu` VALUES (7, 4, '角色信息', 'InfoFilled', 'base:role:list', 2, 'system/role', 2, 2, '2022-09-04 14:00:12');
INSERT INTO `sys_menu` VALUES (8, 4, '菜单信息', 'Histogram', 'base:menu:list', 2, 'system/menu', 2, 3, '2022-09-04 14:00:17');
INSERT INTO `sys_menu` VALUES (10, 4, '资产授权', 'Promotion', 'base:post:list', 2, 'system/post', 2, 5, '2022-09-04 14:02:06');
INSERT INTO `sys_menu` VALUES (16, 6, '新增用户', '', 'base:admin:add', 3, '', 2, 1, '2022-09-04 18:32:55');
INSERT INTO `sys_menu` VALUES (17, 6, '修改用户', '', 'base:admin:edit', 3, '', 2, 2, '2022-09-04 18:33:29');
INSERT INTO `sys_menu` VALUES (18, 6, '删除用户', '', 'base:admin:delete', 3, '', 2, 3, '2022-09-04 18:33:51');
INSERT INTO `sys_menu` VALUES (21, 7, '新增角色', '', 'base:role:add', 3, '', 2, 1, '2022-09-04 18:44:47');
INSERT INTO `sys_menu` VALUES (22, 7, '修改角色', '', 'base:role:edit', 3, '', 2, 2, '2022-09-04 18:45:10');
INSERT INTO `sys_menu` VALUES (23, 7, '删除角色', '', 'base:role:delete', 3, '', 2, 3, '2022-09-04 18:45:46');
INSERT INTO `sys_menu` VALUES (24, 7, '分配权限', '', 'base:role:assign', 3, '', 2, 4, '2022-09-04 18:46:20');
INSERT INTO `sys_menu` VALUES (26, 8, '新增菜单', '', 'base:menu:add', 3, '', 2, 1, '2022-09-04 18:49:51');
INSERT INTO `sys_menu` VALUES (27, 8, '修改菜单', '', 'base:menu:edit', 3, '', 2, 2, '2022-09-04 18:50:24');
INSERT INTO `sys_menu` VALUES (28, 8, '删除菜单', '', 'base:menu:delete', 3, '', 2, 3, '2022-09-04 18:50:53');
INSERT INTO `sys_menu` VALUES (32, 10, '新增岗位', '', 'base:post:add', 3, '', 2, 1, '2022-09-04 18:53:28');
INSERT INTO `sys_menu` VALUES (33, 10, '修改岗位', '', 'base:post:edit', 3, '', 2, 2, '2022-09-04 18:53:48');
INSERT INTO `sys_menu` VALUES (34, 10, '删除岗位', '', 'base:post:delete', 3, '', 2, 3, '2022-09-04 18:54:00');
INSERT INTO `sys_menu` VALUES (44, 0, '审计中心', 'BellFilled', '', 1, '', 2, 11, '2022-09-05 11:06:57');
INSERT INTO `sys_menu` VALUES (45, 44, '操作日志', 'User', 'audit:operator:list', 2, 'audit/operator', 2, 2, '2022-09-05 11:10:54');
INSERT INTO `sys_menu` VALUES (46, 44, '登录日志', 'DocumentRemove', 'audit:loginLog:list', 2, 'audit/loginlog', 2, 1, '2022-09-05 11:11:31');
INSERT INTO `sys_menu` VALUES (47, 45, '清空操作日志', '', 'monitor:operator:clean', 3, '', 2, 1, '2022-09-05 11:12:36');
INSERT INTO `sys_menu` VALUES (49, 46, '清空登录日志', '', 'monitor:loginLog:clean', 3, '', 2, 1, '2022-09-05 11:16:01');
INSERT INTO `sys_menu` VALUES (60, 6, '重置密码', NULL, 'base:admin:reset', 3, NULL, 2, 6, '2022-12-01 16:33:34');
INSERT INTO `sys_menu` VALUES (62, 46, '删除登录日志', '', 'monitor:loginLog:delete', 3, '', 2, 2, '2022-12-02 15:41:56');
INSERT INTO `sys_menu` VALUES (73, 45, '删除操作日志', '', 'monitor:operator:delete', 3, '', 2, 3, '2023-06-02 10:09:38');
INSERT INTO `sys_menu` VALUES (78, 80, '主机管理', 'Platform', 'cmdb:ecs:list', 2, 'cmdb/ecs', 2, 2, '2025-06-29 00:30:35');
INSERT INTO `sys_menu` VALUES (80, 0, '资产管理', 'TrendCharts', '', 1, '', 2, 1, '2025-07-03 11:47:07');
INSERT INTO `sys_menu` VALUES (84, 0, '配置中心', 'Tools', '', 1, '', 2, 5, '2025-07-04 17:00:01');
INSERT INTO `sys_menu` VALUES (85, 84, '主机凭据', 'Setting', 'config:ecs:key', 2, 'config/ecskey', 2, 1, '2025-07-04 17:03:10');
INSERT INTO `sys_menu` VALUES (86, 84, '通用凭据', 'User', 'config:accountauth:key', 2, 'config/accountauth', 2, 2, '2025-07-04 17:08:20');
INSERT INTO `sys_menu` VALUES (88, 80, '网络设备', 'Shop', 'cmdb:switch', 2, 'cmdb/switch', 2, 3, '2025-07-16 15:17:14');
INSERT INTO `sys_menu` VALUES (89, 88, '创建资产', '', 'cmdb:snmp:add', 3, '', 2, 1, '2025-07-18 15:24:31');
INSERT INTO `sys_menu` VALUES (90, 88, '修改资产', '', 'cmdb:snmp:update', 3, '', 2, 2, '2025-07-18 15:25:49');
INSERT INTO `sys_menu` VALUES (91, 88, '删除资产', '', 'cmdb:snmp:delete', 3, '', 2, 3, '2025-07-18 15:26:21');
INSERT INTO `sys_menu` VALUES (104, 84, '密钥管理', 'Phone', 'config:keymanage:key', 2, 'config/keymanage', 2, 3, '2025-09-08 13:24:40');
INSERT INTO `sys_menu` VALUES (113, 45, '批量删除', '', 'monitor:operator:delete', 3, '', 2, 2, '2025-09-17 20:55:13');
INSERT INTO `sys_menu` VALUES (115, 104, '删除密钥', '', 'config:keymanage:delete', 3, '', 2, 2, '2025-09-18 10:53:44');
INSERT INTO `sys_menu` VALUES (117, 104, '同步主机', '', 'config:keymanage:rsync', 3, '', 2, 3, '2025-09-18 10:57:25');
INSERT INTO `sys_menu` VALUES (118, 104, '创建密钥', '', 'config:keymanage:create', 3, '', 2, 4, '2025-09-18 11:01:12');
INSERT INTO `sys_menu` VALUES (119, 86, '修改账号', '', 'config:common:edit', 3, '', 2, 2, '2025-09-18 11:47:33');
INSERT INTO `sys_menu` VALUES (120, 86, '删除账号', '', 'config:common:delete', 3, '', 2, 3, '2025-09-18 11:48:17');
INSERT INTO `sys_menu` VALUES (121, 86, '解密账号', '', 'config:common:decrypt', 3, '', 2, 4, '2025-09-18 11:48:57');
INSERT INTO `sys_menu` VALUES (122, 86, '创建账号', '', 'config:common:add', 3, '', 2, 1, '2025-09-18 11:49:30');
INSERT INTO `sys_menu` VALUES (123, 85, '修改凭据', '', 'config:ecs:edit', 3, '', 2, 1, '2025-09-18 11:54:16');
INSERT INTO `sys_menu` VALUES (124, 85, '删除凭据', '', 'config:ecs:delete', 3, '', 2, 2, '2025-09-18 11:54:51');
INSERT INTO `sys_menu` VALUES (125, 85, '创建凭据', '', 'config:ecs:create', 3, '', 2, 3, '2025-09-18 11:55:21');
INSERT INTO `sys_menu` VALUES (149, 78, '创建主机', '', 'cmdb:ecs:add', 3, '', 2, 1, '2025-09-18 14:47:42');
INSERT INTO `sys_menu` VALUES (150, 78, '主机终端', '', 'cmdb:ecs:terminal', 3, '', 2, 2, '2025-09-18 14:48:36');
INSERT INTO `sys_menu` VALUES (151, 78, '修改主机信息', '', 'cmdb:ecs:edit', 3, '', 2, 3, '2025-09-18 14:49:43');
INSERT INTO `sys_menu` VALUES (152, 78, '上传文件到主机', '', 'cmdb:ecs:upload', 3, '', 2, 4, '2025-09-18 14:50:38');
INSERT INTO `sys_menu` VALUES (153, 78, '执行主机命令', '', 'cmdb:ecs:shell', 3, '', 2, 5, '2025-09-18 14:51:10');
INSERT INTO `sys_menu` VALUES (154, 78, '监控主机', '', 'cmdb:ecs:monitor', 3, '', 2, 6, '2025-09-18 14:51:52');
INSERT INTO `sys_menu` VALUES (155, 78, '删除主机', '', 'cmdb:ecs:delete', 3, '', 2, 7, '2025-09-18 14:52:20');
INSERT INTO `sys_menu` VALUES (165, 78, '连接主机终端', '', 'cmdb:ecs:connecthost', 3, '', 2, 8, '2025-09-18 21:11:43');
INSERT INTO `sys_menu` VALUES (166, 78, '同步主机信息', '', 'cmdb:ecs:rsync', 3, '', 2, 9, '2025-09-19 21:35:06');
INSERT INTO `sys_menu` VALUES (241, 78, '下载文件', '', 'cmdb:ecs:download', 3, '', 2, 10, '2025-12-13 16:13:17');
INSERT INTO `sys_menu` VALUES (243, 78, '批量操作', '', 'cmdb:ecs:batch', 3, '', 2, 10, '2025-12-26 19:04:56');
INSERT INTO `sys_menu` VALUES (244, 104, '查看密钥', '', 'config:keymanage:decrypt', 3, '', 2, 5, '2025-12-27 00:20:50');
INSERT INTO `sys_menu` VALUES (245, 88, '监控大盘', '', 'cmdb:snmp:collect', 3, '', 2, 4, '2026-01-01 15:48:05');
INSERT INTO `sys_menu` VALUES (246, 88, '终端按钮', '', 'cmdb:snmp:ssh', 3, '', 2, 5, '2026-01-01 22:18:54');
INSERT INTO `sys_menu` VALUES (247, 44, '会话录制', 'Coin', 'audit:video', 2, 'audit/video', 2, 3, '2026-01-08 19:57:24');
INSERT INTO `sys_menu` VALUES (250, 86, '测试连接', '', 'config:common:test', 3, '', 2, 5, '2026-01-18 23:27:29');
INSERT INTO `sys_menu` VALUES (279, 78, '查看主机详情', '', 'cmdb:ecs:details', 3, '', 2, 11, '2026-03-23 15:04:04');
INSERT INTO `sys_menu` VALUES (293, 80, '凭据管理', 'Key', 'cmdb:credentials', 2, 'cmdb/credentials', 2, 4, '2026-04-03 18:01:20');
INSERT INTO `sys_menu` VALUES (294, 80, '机房管理', 'User', 'cmdb:idc', 2, 'cmdb/idc', 2, 1, '2026-04-03 18:12:31');
INSERT INTO `sys_menu` VALUES (295, 80, '探针管理', 'CameraFilled', 'cmdb:agent', 2, 'cmdb/agent', 2, 5, '2026-04-07 14:08:09');
INSERT INTO `sys_menu` VALUES (296, 80, '仪表盘', 'CameraFilled', 'cmdb:dashboard', 2, 'cmdb/dashboard', 2, 0, '2026-04-07 14:14:22');
INSERT INTO `sys_menu` VALUES (297, 4, '仪表盘', 'Promotion', 'system:dashboard', 2, 'system/dashboard', 2, 0, '2026-04-07 14:17:30');
INSERT INTO `sys_menu` VALUES (298, 295, '部署探针', '', 'ops:agent:create', 3, '', 2, 1, '2026-04-08 14:57:39');
INSERT INTO `sys_menu` VALUES (299, 295, '卸载探针', '', 'ops:agent:delete', 3, '', 2, 2, '2026-04-08 15:00:42');
INSERT INTO `sys_menu` VALUES (300, 295, '查看探针', '', 'ops:agent:get', 3, '', 2, 3, '2026-04-08 15:01:08');
INSERT INTO `sys_menu` VALUES (301, 295, '批量卸载探针', '', 'ops:agent:deleteall', 3, '', 2, 4, '2026-04-08 15:02:26');
INSERT INTO `sys_menu` VALUES (303, 0, '容器管理', 'Shop', '', 1, '', 2, 2, '2026-04-15 10:33:16');
INSERT INTO `sys_menu` VALUES (306, 303, '集群管理', 'Platform', 'k8s/cluster', 2, 'k8s/cluster', 2, 1, '2026-04-20 14:24:53');
INSERT INTO `sys_menu` VALUES (313, 0, 'AIOps运维', 'Share', '', 1, '', 2, 8, '2026-04-27 11:55:55');
INSERT INTO `sys_menu` VALUES (314, 313, '智能助手', 'MagicStick', 'ai:aiasset', 2, 'ai/aiasset', 2, 1, '2026-04-27 11:57:41');
INSERT INTO `sys_menu` VALUES (315, 313, '运维智库', 'Printer', 'ai:thinktank', 2, 'ai/thinktank', 2, 2, '2026-04-30 17:05:15');
INSERT INTO `sys_menu` VALUES (316, 313, '运维模型', 'Help', 'ai:model', 2, 'ai/model', 2, 10, '2026-04-30 17:05:44');
INSERT INTO `sys_menu` VALUES (317, 313, '日志分析', 'Help', 'ai:logs_analysis', 2, 'ai/logs_analysis', 2, 4, '2026-05-06 15:48:56');
INSERT INTO `sys_menu` VALUES (322, 303, '节点管理', 'Monitor', 'k8s/node', 2, 'k8s/node', 2, 2, '2026-05-13 09:22:17');
INSERT INTO `sys_menu` VALUES (323, 303, '命名空间', 'Collection', 'k8s/namespace', 2, 'k8s/namespace', 2, 3, '2026-05-13 09:22:42');
INSERT INTO `sys_menu` VALUES (324, 303, '工作负载', 'Operation', 'k8s/workload', 2, 'k8s/workload', 2, 4, '2026-05-13 09:23:08');
INSERT INTO `sys_menu` VALUES (325, 303, '网络管理', 'Connection', 'k8s/network', 2, 'k8s/network', 2, 5, '2026-05-13 09:23:52');
INSERT INTO `sys_menu` VALUES (326, 303, '存储管理', 'FolderOpened', 'k8s/storage', 2, 'k8s/storage', 2, 6, '2026-05-13 09:24:11');
INSERT INTO `sys_menu` VALUES (327, 303, '配置管理', 'Setting', 'k8s/config', 2, 'k8s/config', 2, 7, '2026-05-13 09:24:39');
INSERT INTO `sys_menu` VALUES (328, 303, '应用诊断', 'FirstAidKit', 'k8s/arthas', 2, 'k8s/arthas', 2, 8, '2026-05-13 09:25:09');
INSERT INTO `sys_menu` VALUES (329, 306, '注册集群', '', 'cloud:k8s:register', 3, '', 2, 1, '2026-05-13 09:26:22');
INSERT INTO `sys_menu` VALUES (330, 306, '修改集群', '', 'cloud:k8s:edit', 3, '', 2, 2, '2026-05-13 09:26:39');
INSERT INTO `sys_menu` VALUES (331, 306, '同步集群', '', 'cloud:k8s:rsync', 3, '', 2, 3, '2026-05-13 09:26:56');
INSERT INTO `sys_menu` VALUES (332, 306, '删除集群', '', 'cloud:k8s:delete', 3, '', 2, 4, '2026-05-13 09:27:11');
INSERT INTO `sys_menu` VALUES (333, 306, '查看凭据', '', 'cloud:k8s:get', 3, '', 2, 5, '2026-05-13 09:27:32');
INSERT INTO `sys_menu` VALUES (334, 322, '查看监控仪表盘', '', 'k8s:node:monitor', 3, '', 2, 1, '2026-05-13 09:28:51');
INSERT INTO `sys_menu` VALUES (335, 322, '查看节点资源详情', '', 'k8s:node:details', 3, '', 2, 2, '2026-05-13 09:29:10');
INSERT INTO `sys_menu` VALUES (336, 322, '节点更多操作', '', 'k8s:node:actions', 3, '', 2, 3, '2026-05-13 09:29:28');
INSERT INTO `sys_menu` VALUES (337, 323, '创建命名空间', '', 'k8s:namespace:add', 3, '', 2, 1, '2026-05-13 09:30:06');
INSERT INTO `sys_menu` VALUES (338, 323, '查看命名空间详情', '', 'k8s:namespace:details', 3, '', 2, 2, '2026-05-13 09:30:24');
INSERT INTO `sys_menu` VALUES (339, 323, '查看命名空间配置', '', 'k8s:namespace:setup', 3, '', 2, 3, '2026-05-13 09:30:46');
INSERT INTO `sys_menu` VALUES (340, 323, '查看限制命名空间', '', 'k8s:namespace:restriction', 3, '', 2, 4, '2026-05-13 09:31:45');
INSERT INTO `sys_menu` VALUES (341, 323, '命名空间更多操作', '', 'k8s:namespace:actions', 3, '', 2, 5, '2026-05-13 09:32:14');
INSERT INTO `sys_menu` VALUES (342, 324, '新增工作负载', '', 'k8s:workload:add', 3, '', 2, 1, '2026-05-13 09:33:08');
INSERT INTO `sys_menu` VALUES (343, 324, '工作负载更多操作', '', 'k8s:workload:actions', 3, '', 2, 2, '2026-05-13 09:33:38');
INSERT INTO `sys_menu` VALUES (344, 324, '回滚工作负载', '', 'k8s:workload:rollback_version', 3, '', 2, 3, '2026-05-13 09:34:09');
INSERT INTO `sys_menu` VALUES (345, 324, '查看pod日志', '', 'k8s:workload:podlog', 3, '', 2, 4, '2026-05-13 09:34:44');
INSERT INTO `sys_menu` VALUES (346, 324, '删除pod', '', 'k8s:workload:poddelete', 3, '', 2, 5, '2026-05-13 09:35:04');
INSERT INTO `sys_menu` VALUES (347, 324, '登陆pod终端', '', 'k8s:workload:terminal', 3, '', 2, 6, '2026-05-13 09:35:22');
INSERT INTO `sys_menu` VALUES (348, 324, '编辑pod-yaml', '', 'k8s:workload:edityaml', 3, '', 2, 7, '2026-05-13 09:35:49');
INSERT INTO `sys_menu` VALUES (349, 325, '新增service', '', 'k8s:network:addservice', 3, '', 2, 1, '2026-05-13 09:36:45');
INSERT INTO `sys_menu` VALUES (350, 325, '新增ingress', '', 'k8s:network:addingress', 3, '', 2, 2, '2026-05-13 09:37:07');
INSERT INTO `sys_menu` VALUES (351, 325, '更多操作', '', 'k8s:network:actions', 3, '', 2, 3, '2026-05-13 09:37:22');
INSERT INTO `sys_menu` VALUES (352, 326, 'k8s存储操作', '', 'k8s:storage:actions', 3, '', 2, 1, '2026-05-13 09:38:07');
INSERT INTO `sys_menu` VALUES (353, 326, '查看k8s存储配置', '', 'k8s:storage:get', 3, '', 2, 2, '2026-05-13 09:38:31');
INSERT INTO `sys_menu` VALUES (354, 327, 'k8s配置更多操作', '', 'k8s:config:actions', 3, '', 2, 1, '2026-05-13 09:39:04');
INSERT INTO `sys_menu` VALUES (355, 327, '查看配置详情', '', 'k8s:config:get', 3, '', 2, 2, '2026-05-13 09:39:22');
INSERT INTO `sys_menu` VALUES (357, 303, '集群巡检', 'Finished', 'k8s:inspection', 2, 'k8s/inspection', 2, 9, '2026-05-18 17:08:57');
INSERT INTO `sys_menu` VALUES (358, 80, '数据库管理', 'Eleme', 'cmdb:db-management', 2, 'cmdb/db-management', 2, 7, '2026-05-19 11:43:47');
INSERT INTO `sys_menu` VALUES (359, 0, 'DB数据库', '', '', 1, '', 2, 3, '2026-05-19 15:02:29');
INSERT INTO `sys_menu` VALUES (360, 359, '数据资产', 'Coin', 'db:management', 2, 'db/management', 2, 1, '2026-05-19 15:09:36');
INSERT INTO `sys_menu` VALUES (361, 359, 'SQL查询', 'Coin', 'db:crud', 2, 'db/crud', 2, 2, '2026-05-19 15:14:05');
INSERT INTO `sys_menu` VALUES (362, 359, 'SQL审核', 'Coin', 'db:review', 2, 'db/review', 2, 3, '2026-05-19 15:15:59');
INSERT INTO `sys_menu` VALUES (363, 359, '性能分析', 'Coin', 'db:backup', 2, 'db/backup', 2, 5, '2026-05-19 15:17:16');
INSERT INTO `sys_menu` VALUES (364, 359, '监控概览', 'Coin', 'db:monitoring', 2, 'db/monitoring', 2, 4, '2026-05-19 15:18:41');
INSERT INTO `sys_menu` VALUES (365, 359, '慢日志分析', 'Coin', 'db:slowlog', 2, 'db/slowlog', 2, 6, '2026-05-19 15:20:13');
INSERT INTO `sys_menu` VALUES (366, 359, '数据库AI助手', 'Coin', 'db:operation_log', 2, 'db/operation_log', 2, 7, '2026-05-19 15:27:25');
INSERT INTO `sys_menu` VALUES (367, 303, '事件管理', 'Bell', 'k8s:events:list', 2, 'k8s/events', 2, 10, '2026-05-20 11:31:29');
INSERT INTO `sys_menu` VALUES (369, 4, '菜单导航', 'Promotion', 'base:nav', 2, 'system/nav', 2, 6, '2026-05-29 18:43:53');
INSERT INTO `sys_menu` VALUES (370, 80, '主机巡检', 'Wallet', 'cmdb:inspection', 2, 'cmdb/inspection', 2, 6, '2026-06-08 17:27:05');
INSERT INTO `sys_menu` VALUES (374, 80, '中间件管理', 'Eleme', 'cmdb:middleware', 2, 'cmdb/middleware', 2, 8, '2026-07-11 00:10:18');
INSERT INTO `sys_menu` VALUES (383, 0, '日志中心', 'Shop', '', 1, '', 2, 4, '2026-07-11 00:21:23');
INSERT INTO `sys_menu` VALUES (384, 383, '日志管理', 'Menu', 'logs:project', 2, 'logs/project', 2, 1, '2026-07-11 00:22:00');
INSERT INTO `sys_menu` VALUES (385, 383, '日志查询', 'Menu', 'logs:query', 2, 'logs/query', 2, 2, '2026-07-11 00:22:37');
INSERT INTO `sys_menu` VALUES (386, 383, '日志采集', 'Menu', 'logs:collect', 2, 'logs/collect', 2, 3, '2026-07-11 00:23:33');
INSERT INTO `sys_menu` VALUES (387, 383, '日志告警', 'Menu', 'logs:alert', 2, 'logs/alert', 2, 4, '2026-07-11 00:24:18');
INSERT INTO `sys_menu` VALUES (388, 0, '费用中心', 'Share', '', 1, '', 2, 10, '2026-07-11 00:25:31');
INSERT INTO `sys_menu` VALUES (389, 388, '账单总览', 'House', 'cost:bill_dashboard', 2, 'cost/bill_dashboard', 2, 1, '2026-07-11 00:26:25');
INSERT INTO `sys_menu` VALUES (390, 388, '账单明细', 'House', 'cost:cloud_bill', 2, 'cost/cloud_bill', 2, 2, '2026-07-11 00:27:13');
INSERT INTO `sys_menu` VALUES (391, 388, '账单分析', 'House', 'cost:cloud_view', 2, 'cost/cloud_view', 2, 3, '2026-07-11 00:28:08');
INSERT INTO `sys_menu` VALUES (393, 0, '告警中心', '', '', 1, '', 2, 6, '2026-07-11 00:30:09');
INSERT INTO `sys_menu` VALUES (394, 393, '告警中心', 'List', 'alert:view', 2, 'alert/view', 2, 1, '2026-07-11 00:31:32');
INSERT INTO `sys_menu` VALUES (395, 393, '告警事件', 'List', 'alert:event', 2, 'alert/event', 2, 2, '2026-07-11 00:33:00');
INSERT INTO `sys_menu` VALUES (396, 393, '告警管理', 'List', 'alert:rule', 2, 'alert/rule', 2, 3, '2026-07-11 00:34:05');
INSERT INTO `sys_menu` VALUES (397, 393, '告警模板', 'List', 'alert:template', 2, 'alert/template', 2, 4, '2026-07-11 00:34:38');
INSERT INTO `sys_menu` VALUES (398, 393, '告警通知', 'List', 'alert:notice', 2, 'alert/notice', 2, 5, '2026-07-11 00:35:26');
INSERT INTO `sys_menu` VALUES (399, 393, '告警数据源', 'List', 'alert:source', 2, 'alert/source', 2, 7, '2026-07-11 00:36:27');
INSERT INTO `sys_menu` VALUES (400, 0, '工单中心', '', '', 1, '', 2, 9, '2026-07-14 01:23:02');
INSERT INTO `sys_menu` VALUES (401, 400, '我的工单', 'StarFilled', 'ticket:management', 2, 'ticket/management', 2, 1, '2026-07-14 01:23:32');
INSERT INTO `sys_menu` VALUES (402, 0, '服务作业', 'StarFilled', '', 1, '', 2, 7, '2026-07-14 01:23:55');
INSERT INTO `sys_menu` VALUES (403, 402, '作业调度', 'DocumentCopy', 'task:schedule', 2, 'task/schedule', 2, 1, '2026-07-14 01:24:24');
INSERT INTO `sys_menu` VALUES (404, 393, '域名告警', 'Link', 'alert:domain:list', 2, 'alert/domain', 2, 6, '2026-07-30 23:13:51');
INSERT INTO `sys_menu` VALUES (405, 400, '事务工单', 'StarFilled', 'ticket:daily', 2, 'ticket/daily', 2, 2, '2026-07-30 23:26:02');
INSERT INTO `sys_menu` VALUES (406, 400, '发布工单', 'StarFilled', 'ticket:release', 2, 'ticket/release', 2, 3, '2026-07-30 23:26:27');
INSERT INTO `sys_menu` VALUES (407, 400, '审批流配置', 'StarFilled', 'ticket:workflow', 2, 'ticket/workflow', 2, 4, '2026-07-30 23:27:30');
INSERT INTO `sys_menu` VALUES (408, 383, '日志链路', 'Menu', 'logs:link', 2, 'logs/link', 2, 5, '2026-07-30 23:28:34');
INSERT INTO `sys_menu` VALUES (409, 402, 'Ansible作业', 'Document', 'task:ansible', 2, 'task/ansible', 2, 2, '2026-07-31 22:10:29');
INSERT INTO `sys_menu` VALUES (410, 402, '作业模板', 'DocumentDelete', 'task:templates', 2, 'task/templates', 2, 3, '2026-07-31 22:11:07');
INSERT INTO `sys_menu` VALUES (411, 402, '作业中心', 'Files', 'task:monitor', 2, 'task/monitor', 2, 4, '2026-07-31 22:12:28');
INSERT INTO `sys_menu` VALUES (412, 4, '授权管理', 'Key', 'system:license:view', 2, 'system/license', 2, 7, '2026-08-02 16:20:24');
INSERT INTO `sys_menu` VALUES (413, 383, '数据慢日志', 'Menu', 'logs:mysql-slow', 2, 'logs/query/mysql-slow', 2, 6, '2026-11-10 20:46:03');
INSERT INTO `sys_menu` VALUES (414, 383, '业务访问日志', 'Menu', 'logs:nginx', 2, 'logs/nginx', 2, 7, '2026-11-10 20:46:20');
INSERT INTO `sys_menu` VALUES (415, 393, '监控中心', 'List', 'alert:monitoring', 2, 'alert/monitoring', 2, 0, '2026-11-10 20:47:00');

-- ----------------------------
-- Table structure for sys_nav_category
-- ----------------------------
DROP TABLE IF EXISTS `sys_nav_category`;
CREATE TABLE `sys_nav_category`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'分类名称\'',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'分类编码\'',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'分类图标\'',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'分类描述\'',
  `sort` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'排序\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'状态：1->启用,2->禁用\'',
  `create_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'创建人ID\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_nav_category_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_nav_category
-- ----------------------------
INSERT INTO `sys_nav_category` VALUES (1, '公共服务', 'public_service', 'Grid', '公司内部公共系统、项目系统和运维平台入口', 10, 1, 0, '2026-11-10 20:06:35.387', '2026-11-10 20:06:35.387');
INSERT INTO `sys_nav_category` VALUES (2, '个人工具', 'personal_tools', 'Star', '个人常用工具网站和个人链接', 20, 1, 0, '2026-11-10 20:06:35.387', '2026-11-10 20:06:35.387');

-- ----------------------------
-- Table structure for sys_nav_item
-- ----------------------------
DROP TABLE IF EXISTS `sys_nav_item`;
CREATE TABLE `sys_nav_item`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键\'',
  `category_id` bigint UNSIGNED NOT NULL COMMENT '\'分类ID\'',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'导航名称\'',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'导航地址\'',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'图标地址或图标名\'',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'说明\'',
  `environment` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'环境：dev/test/pre/prod/none\'',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'标签，逗号分隔\'',
  `scope` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'可见范围：public/private/acl\'',
  `owner_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'所属用户ID\'',
  `owner_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'所属用户名\'',
  `visible_user_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'可见用户ID，逗号分隔\'',
  `visible_role_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'可见角色ID，逗号分隔\'',
  `visible_dept_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'可见部门ID，逗号分隔\'',
  `sso_enabled` tinyint(1) NULL DEFAULT NULL COMMENT '\'是否支持单点登录\'',
  `sso_mode` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'SSO模式：oidc/cas/saml/ticket/proxy/none\'',
  `open_mode` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'blank' COMMENT '\'打开方式：blank/self\'',
  `sort` bigint UNSIGNED NULL DEFAULT NULL COMMENT '\'排序\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'状态：1->启用,2->禁用\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_nav_item_category_id`(`category_id` ASC) USING BTREE,
  INDEX `idx_sys_nav_item_scope`(`scope` ASC) USING BTREE,
  INDEX `idx_sys_nav_item_owner_id`(`owner_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_nav_item
-- ----------------------------

-- ----------------------------
-- Table structure for sys_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_log`;
CREATE TABLE `sys_operation_log`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `admin_id` bigint UNSIGNED NOT NULL COMMENT '\'管理员id\'',
  `username` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '\'管理员账号\'',
  `method` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '\'请求方式\'',
  `ip` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '\'IP\'',
  `url` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '\'URL\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `description` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '\'操作描述\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `username`(`username` ASC) USING BTREE,
  INDEX `create_time`(`create_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 36 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '操作日志记录表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_operation_log
-- ----------------------------
INSERT INTO `sys_operation_log` VALUES (1, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/sysOperationLog/clean', '2026-08-03 01:09:18.309', '清空操作日志接口');
INSERT INTO `sys_operation_log` VALUES (2, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/recordings/batch', '2026-08-03 01:09:24.603', 'swagger未定义: DELETE /api/v1/recordings/batch');
INSERT INTO `sys_operation_log` VALUES (3, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/recordings/batch', '2026-08-03 01:09:29.413', 'swagger未定义: DELETE /api/v1/recordings/batch');
INSERT INTO `sys_operation_log` VALUES (4, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/recordings/batch', '2026-08-03 01:09:33.114', 'swagger未定义: DELETE /api/v1/recordings/batch');
INSERT INTO `sys_operation_log` VALUES (5, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/recordings/batch', '2026-08-03 01:09:37.166', 'swagger未定义: DELETE /api/v1/recordings/batch');
INSERT INTO `sys_operation_log` VALUES (6, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/aiops/llm/providers', '2026-08-03 01:09:57.092', '删除 LLM 供应商');
INSERT INTO `sys_operation_log` VALUES (7, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/config/ecsauthdelete', '2026-08-03 01:10:31.889', '删除凭据');
INSERT INTO `sys_operation_log` VALUES (8, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/config/ecsauthdelete', '2026-08-03 01:10:33.644', '删除凭据');
INSERT INTO `sys_operation_log` VALUES (9, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/config/ecsauthdelete', '2026-08-03 01:10:38.966', '删除凭据');
INSERT INTO `sys_operation_log` VALUES (10, 89, 'admin', 'put', '192.168.3.7', '/api/v1/role/assignPermissions', '2026-08-03 01:11:46.102', '分配权限接口');
INSERT INTO `sys_operation_log` VALUES (11, 89, 'admin', 'post', '14.127.222.121', '/api/v1/system/license/activate', '2026-11-10 20:14:51.309', 'swagger未定义: POST /api/v1/system/license/activate');
INSERT INTO `sys_operation_log` VALUES (12, 89, 'admin', 'post', '14.127.222.121', '/api/v1/upload', '2026-11-10 20:17:37.720', 'swagger未定义: POST /api/v1/upload');
INSERT INTO `sys_operation_log` VALUES (13, 89, 'admin', 'put', '14.127.222.121', '/api/v1/admin/update', '2026-11-10 20:17:47.469', 'swagger未定义: PUT /api/v1/admin/update');
INSERT INTO `sys_operation_log` VALUES (14, 89, 'admin', 'post', '14.127.222.121', '/api/v1/upload', '2026-11-10 20:18:40.760', 'swagger未定义: POST /api/v1/upload');
INSERT INTO `sys_operation_log` VALUES (15, 89, 'admin', 'post', '14.127.222.121', '/api/v1/upload', '2026-11-10 20:18:59.058', 'swagger未定义: POST /api/v1/upload');
INSERT INTO `sys_operation_log` VALUES (16, 89, 'admin', 'post', '14.127.222.121', '/api/v1/upload', '2026-11-10 20:26:20.144', 'swagger未定义: POST /api/v1/upload');
INSERT INTO `sys_operation_log` VALUES (17, 89, 'admin', 'put', '14.127.222.121', '/api/v1/admin/update', '2026-11-10 20:26:22.377', 'swagger未定义: PUT /api/v1/admin/update');
INSERT INTO `sys_operation_log` VALUES (18, 89, 'admin', 'post', '14.127.222.121', '/api/v1/upload', '2026-11-10 20:27:04.660', 'swagger未定义: POST /api/v1/upload');
INSERT INTO `sys_operation_log` VALUES (19, 89, 'admin', 'put', '14.127.222.121', '/api/v1/admin/update', '2026-11-10 20:27:07.302', 'swagger未定义: PUT /api/v1/admin/update');
INSERT INTO `sys_operation_log` VALUES (20, 89, 'admin', 'put', '14.127.222.121', '/api/v1/admin/updatePassword', '2026-11-10 20:29:32.159', 'swagger未定义: PUT /api/v1/admin/updatePassword');
INSERT INTO `sys_operation_log` VALUES (21, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:42:59.153', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (22, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:43:06.786', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (23, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:43:16.776', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (24, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:43:32.646', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (25, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:43:48.051', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (26, 89, 'admin', 'post', '14.127.222.121', '/api/v1/menu/add', '2026-11-10 20:46:02.736', 'swagger未定义: POST /api/v1/menu/add');
INSERT INTO `sys_operation_log` VALUES (27, 89, 'admin', 'post', '14.127.222.121', '/api/v1/menu/add', '2026-11-10 20:46:20.031', 'swagger未定义: POST /api/v1/menu/add');
INSERT INTO `sys_operation_log` VALUES (28, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:46:27.466', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (29, 89, 'admin', 'post', '14.127.222.121', '/api/v1/menu/add', '2026-11-10 20:47:00.254', 'swagger未定义: POST /api/v1/menu/add');
INSERT INTO `sys_operation_log` VALUES (30, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:47:17.035', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (31, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:47:43.539', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (32, 89, 'admin', 'put', '14.127.222.121', '/api/v1/menu/update', '2026-11-10 20:47:49.888', 'swagger未定义: PUT /api/v1/menu/update');
INSERT INTO `sys_operation_log` VALUES (33, 89, 'admin', 'put', '14.127.222.121', '/api/v1/role/assignPermissions', '2026-11-10 20:50:07.681', 'swagger未定义: PUT /api/v1/role/assignPermissions');
INSERT INTO `sys_operation_log` VALUES (34, 89, 'admin', 'delete', '14.127.222.121', '/api/v1/cmdb/database-asset/delete', '2026-11-10 20:51:48.102', 'swagger未定义: DELETE /api/v1/cmdb/database-asset/delete');
INSERT INTO `sys_operation_log` VALUES (35, 89, 'admin', 'put', '14.127.222.121', '/api/v1/admin/updatePassword', '2026-11-10 21:05:27.652', 'swagger未定义: PUT /api/v1/admin/updatePassword');

-- ----------------------------
-- Table structure for sys_post
-- ----------------------------
DROP TABLE IF EXISTS `sys_post`;
CREATE TABLE `sys_post`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `post_code` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '岗位编码',
  `post_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '岗位名称',
  `post_status` int NOT NULL DEFAULT 1 COMMENT '状态（1->正常 2->停用）',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `remark` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_bin COMMENT = '岗位信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_post
-- ----------------------------
INSERT INTO `sys_post` VALUES (1, 'AAA', '研发总监', 1, '2023-06-14 20:08:22', '主管各个部门');
INSERT INTO `sys_post` VALUES (10, 'ops', '运维工程师', 1, '2025-06-28 22:46:33', '运维工程师');
INSERT INTO `sys_post` VALUES (11, 'dev', '研发工程师', 1, '2025-06-28 22:50:29', '研发工程师');
INSERT INTO `sys_post` VALUES (12, 'test', '测试工程师', 1, '2025-06-28 22:52:57', '测试工程师');
INSERT INTO `sys_post` VALUES (14, 'cccc', '运维总监', 1, '2026-01-14 11:12:04', '');
INSERT INTO `sys_post` VALUES (15, 'bbb', '测试经理', 1, '2026-01-14 11:12:44', '');
INSERT INTO `sys_post` VALUES (16, 'technical_manager', '技术经理', 1, '2026-07-23 00:40:50', '技术管理岗位');

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `role_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '角色名称',
  `role_key` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '角色权限字符串',
  `status` int NOT NULL DEFAULT 1 COMMENT '启用状态：1->启用；2->禁用',
  `description` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '描述',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `name`(`role_name` ASC) USING BTREE,
  INDEX `create_time`(`create_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 18 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '后台角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES (1, '超级管理员', 'admin', 1, '最大权限', '2023-06-12 20:04:53');
INSERT INTO `sys_role` VALUES (14, 'test', 'test', 1, '123', '2026-04-22 14:44:53');
INSERT INTO `sys_role` VALUES (17, '测试组', 'test_group', 1, '测试组角色', '2026-07-23 00:40:50');

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
  `role_id` bigint UNSIGNED NOT NULL COMMENT '\'角色id\'',
  `menu_id` bigint UNSIGNED NOT NULL COMMENT '\'菜单id\''
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '角色和菜单关系表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO `sys_role_menu` VALUES (14, 78);
INSERT INTO `sys_role_menu` VALUES (14, 149);
INSERT INTO `sys_role_menu` VALUES (14, 150);
INSERT INTO `sys_role_menu` VALUES (14, 151);
INSERT INTO `sys_role_menu` VALUES (14, 152);
INSERT INTO `sys_role_menu` VALUES (14, 153);
INSERT INTO `sys_role_menu` VALUES (14, 154);
INSERT INTO `sys_role_menu` VALUES (14, 155);
INSERT INTO `sys_role_menu` VALUES (14, 165);
INSERT INTO `sys_role_menu` VALUES (14, 166);
INSERT INTO `sys_role_menu` VALUES (14, 241);
INSERT INTO `sys_role_menu` VALUES (14, 243);
INSERT INTO `sys_role_menu` VALUES (14, 279);
INSERT INTO `sys_role_menu` VALUES (14, 303);
INSERT INTO `sys_role_menu` VALUES (14, 306);
INSERT INTO `sys_role_menu` VALUES (14, 329);
INSERT INTO `sys_role_menu` VALUES (14, 330);
INSERT INTO `sys_role_menu` VALUES (14, 331);
INSERT INTO `sys_role_menu` VALUES (14, 332);
INSERT INTO `sys_role_menu` VALUES (14, 333);
INSERT INTO `sys_role_menu` VALUES (14, 322);
INSERT INTO `sys_role_menu` VALUES (14, 334);
INSERT INTO `sys_role_menu` VALUES (14, 335);
INSERT INTO `sys_role_menu` VALUES (14, 336);
INSERT INTO `sys_role_menu` VALUES (14, 323);
INSERT INTO `sys_role_menu` VALUES (14, 337);
INSERT INTO `sys_role_menu` VALUES (14, 338);
INSERT INTO `sys_role_menu` VALUES (14, 339);
INSERT INTO `sys_role_menu` VALUES (14, 340);
INSERT INTO `sys_role_menu` VALUES (14, 341);
INSERT INTO `sys_role_menu` VALUES (14, 324);
INSERT INTO `sys_role_menu` VALUES (14, 342);
INSERT INTO `sys_role_menu` VALUES (14, 343);
INSERT INTO `sys_role_menu` VALUES (14, 344);
INSERT INTO `sys_role_menu` VALUES (14, 345);
INSERT INTO `sys_role_menu` VALUES (14, 346);
INSERT INTO `sys_role_menu` VALUES (14, 347);
INSERT INTO `sys_role_menu` VALUES (14, 348);
INSERT INTO `sys_role_menu` VALUES (14, 325);
INSERT INTO `sys_role_menu` VALUES (14, 349);
INSERT INTO `sys_role_menu` VALUES (14, 350);
INSERT INTO `sys_role_menu` VALUES (14, 351);
INSERT INTO `sys_role_menu` VALUES (14, 326);
INSERT INTO `sys_role_menu` VALUES (14, 352);
INSERT INTO `sys_role_menu` VALUES (14, 353);
INSERT INTO `sys_role_menu` VALUES (14, 327);
INSERT INTO `sys_role_menu` VALUES (14, 354);
INSERT INTO `sys_role_menu` VALUES (14, 355);
INSERT INTO `sys_role_menu` VALUES (14, 328);
INSERT INTO `sys_role_menu` VALUES (14, 357);
INSERT INTO `sys_role_menu` VALUES (14, 367);
INSERT INTO `sys_role_menu` VALUES (14, 383);
INSERT INTO `sys_role_menu` VALUES (14, 384);
INSERT INTO `sys_role_menu` VALUES (14, 385);
INSERT INTO `sys_role_menu` VALUES (14, 386);
INSERT INTO `sys_role_menu` VALUES (14, 387);
INSERT INTO `sys_role_menu` VALUES (14, 80);
INSERT INTO `sys_role_menu` VALUES (1, 80);
INSERT INTO `sys_role_menu` VALUES (1, 294);
INSERT INTO `sys_role_menu` VALUES (1, 78);
INSERT INTO `sys_role_menu` VALUES (1, 149);
INSERT INTO `sys_role_menu` VALUES (1, 150);
INSERT INTO `sys_role_menu` VALUES (1, 151);
INSERT INTO `sys_role_menu` VALUES (1, 152);
INSERT INTO `sys_role_menu` VALUES (1, 153);
INSERT INTO `sys_role_menu` VALUES (1, 154);
INSERT INTO `sys_role_menu` VALUES (1, 155);
INSERT INTO `sys_role_menu` VALUES (1, 165);
INSERT INTO `sys_role_menu` VALUES (1, 166);
INSERT INTO `sys_role_menu` VALUES (1, 241);
INSERT INTO `sys_role_menu` VALUES (1, 243);
INSERT INTO `sys_role_menu` VALUES (1, 279);
INSERT INTO `sys_role_menu` VALUES (1, 88);
INSERT INTO `sys_role_menu` VALUES (1, 89);
INSERT INTO `sys_role_menu` VALUES (1, 90);
INSERT INTO `sys_role_menu` VALUES (1, 91);
INSERT INTO `sys_role_menu` VALUES (1, 245);
INSERT INTO `sys_role_menu` VALUES (1, 246);
INSERT INTO `sys_role_menu` VALUES (1, 293);
INSERT INTO `sys_role_menu` VALUES (1, 295);
INSERT INTO `sys_role_menu` VALUES (1, 298);
INSERT INTO `sys_role_menu` VALUES (1, 299);
INSERT INTO `sys_role_menu` VALUES (1, 300);
INSERT INTO `sys_role_menu` VALUES (1, 301);
INSERT INTO `sys_role_menu` VALUES (1, 370);
INSERT INTO `sys_role_menu` VALUES (1, 358);
INSERT INTO `sys_role_menu` VALUES (1, 374);
INSERT INTO `sys_role_menu` VALUES (1, 296);
INSERT INTO `sys_role_menu` VALUES (1, 303);
INSERT INTO `sys_role_menu` VALUES (1, 306);
INSERT INTO `sys_role_menu` VALUES (1, 329);
INSERT INTO `sys_role_menu` VALUES (1, 330);
INSERT INTO `sys_role_menu` VALUES (1, 331);
INSERT INTO `sys_role_menu` VALUES (1, 332);
INSERT INTO `sys_role_menu` VALUES (1, 333);
INSERT INTO `sys_role_menu` VALUES (1, 322);
INSERT INTO `sys_role_menu` VALUES (1, 334);
INSERT INTO `sys_role_menu` VALUES (1, 335);
INSERT INTO `sys_role_menu` VALUES (1, 336);
INSERT INTO `sys_role_menu` VALUES (1, 323);
INSERT INTO `sys_role_menu` VALUES (1, 337);
INSERT INTO `sys_role_menu` VALUES (1, 338);
INSERT INTO `sys_role_menu` VALUES (1, 339);
INSERT INTO `sys_role_menu` VALUES (1, 340);
INSERT INTO `sys_role_menu` VALUES (1, 341);
INSERT INTO `sys_role_menu` VALUES (1, 324);
INSERT INTO `sys_role_menu` VALUES (1, 342);
INSERT INTO `sys_role_menu` VALUES (1, 343);
INSERT INTO `sys_role_menu` VALUES (1, 344);
INSERT INTO `sys_role_menu` VALUES (1, 345);
INSERT INTO `sys_role_menu` VALUES (1, 346);
INSERT INTO `sys_role_menu` VALUES (1, 347);
INSERT INTO `sys_role_menu` VALUES (1, 348);
INSERT INTO `sys_role_menu` VALUES (1, 325);
INSERT INTO `sys_role_menu` VALUES (1, 349);
INSERT INTO `sys_role_menu` VALUES (1, 350);
INSERT INTO `sys_role_menu` VALUES (1, 351);
INSERT INTO `sys_role_menu` VALUES (1, 326);
INSERT INTO `sys_role_menu` VALUES (1, 352);
INSERT INTO `sys_role_menu` VALUES (1, 353);
INSERT INTO `sys_role_menu` VALUES (1, 327);
INSERT INTO `sys_role_menu` VALUES (1, 354);
INSERT INTO `sys_role_menu` VALUES (1, 355);
INSERT INTO `sys_role_menu` VALUES (1, 328);
INSERT INTO `sys_role_menu` VALUES (1, 357);
INSERT INTO `sys_role_menu` VALUES (1, 367);
INSERT INTO `sys_role_menu` VALUES (1, 359);
INSERT INTO `sys_role_menu` VALUES (1, 360);
INSERT INTO `sys_role_menu` VALUES (1, 361);
INSERT INTO `sys_role_menu` VALUES (1, 362);
INSERT INTO `sys_role_menu` VALUES (1, 364);
INSERT INTO `sys_role_menu` VALUES (1, 363);
INSERT INTO `sys_role_menu` VALUES (1, 365);
INSERT INTO `sys_role_menu` VALUES (1, 366);
INSERT INTO `sys_role_menu` VALUES (1, 383);
INSERT INTO `sys_role_menu` VALUES (1, 384);
INSERT INTO `sys_role_menu` VALUES (1, 385);
INSERT INTO `sys_role_menu` VALUES (1, 386);
INSERT INTO `sys_role_menu` VALUES (1, 387);
INSERT INTO `sys_role_menu` VALUES (1, 408);
INSERT INTO `sys_role_menu` VALUES (1, 413);
INSERT INTO `sys_role_menu` VALUES (1, 414);
INSERT INTO `sys_role_menu` VALUES (1, 84);
INSERT INTO `sys_role_menu` VALUES (1, 85);
INSERT INTO `sys_role_menu` VALUES (1, 123);
INSERT INTO `sys_role_menu` VALUES (1, 124);
INSERT INTO `sys_role_menu` VALUES (1, 125);
INSERT INTO `sys_role_menu` VALUES (1, 86);
INSERT INTO `sys_role_menu` VALUES (1, 122);
INSERT INTO `sys_role_menu` VALUES (1, 119);
INSERT INTO `sys_role_menu` VALUES (1, 120);
INSERT INTO `sys_role_menu` VALUES (1, 121);
INSERT INTO `sys_role_menu` VALUES (1, 250);
INSERT INTO `sys_role_menu` VALUES (1, 104);
INSERT INTO `sys_role_menu` VALUES (1, 115);
INSERT INTO `sys_role_menu` VALUES (1, 117);
INSERT INTO `sys_role_menu` VALUES (1, 118);
INSERT INTO `sys_role_menu` VALUES (1, 244);
INSERT INTO `sys_role_menu` VALUES (1, 393);
INSERT INTO `sys_role_menu` VALUES (1, 394);
INSERT INTO `sys_role_menu` VALUES (1, 395);
INSERT INTO `sys_role_menu` VALUES (1, 396);
INSERT INTO `sys_role_menu` VALUES (1, 397);
INSERT INTO `sys_role_menu` VALUES (1, 398);
INSERT INTO `sys_role_menu` VALUES (1, 404);
INSERT INTO `sys_role_menu` VALUES (1, 399);
INSERT INTO `sys_role_menu` VALUES (1, 415);
INSERT INTO `sys_role_menu` VALUES (1, 402);
INSERT INTO `sys_role_menu` VALUES (1, 403);
INSERT INTO `sys_role_menu` VALUES (1, 409);
INSERT INTO `sys_role_menu` VALUES (1, 410);
INSERT INTO `sys_role_menu` VALUES (1, 411);
INSERT INTO `sys_role_menu` VALUES (1, 313);
INSERT INTO `sys_role_menu` VALUES (1, 314);
INSERT INTO `sys_role_menu` VALUES (1, 315);
INSERT INTO `sys_role_menu` VALUES (1, 317);
INSERT INTO `sys_role_menu` VALUES (1, 316);
INSERT INTO `sys_role_menu` VALUES (1, 400);
INSERT INTO `sys_role_menu` VALUES (1, 401);
INSERT INTO `sys_role_menu` VALUES (1, 405);
INSERT INTO `sys_role_menu` VALUES (1, 406);
INSERT INTO `sys_role_menu` VALUES (1, 407);
INSERT INTO `sys_role_menu` VALUES (1, 388);
INSERT INTO `sys_role_menu` VALUES (1, 389);
INSERT INTO `sys_role_menu` VALUES (1, 390);
INSERT INTO `sys_role_menu` VALUES (1, 391);
INSERT INTO `sys_role_menu` VALUES (1, 44);
INSERT INTO `sys_role_menu` VALUES (1, 46);
INSERT INTO `sys_role_menu` VALUES (1, 49);
INSERT INTO `sys_role_menu` VALUES (1, 62);
INSERT INTO `sys_role_menu` VALUES (1, 45);
INSERT INTO `sys_role_menu` VALUES (1, 47);
INSERT INTO `sys_role_menu` VALUES (1, 113);
INSERT INTO `sys_role_menu` VALUES (1, 73);
INSERT INTO `sys_role_menu` VALUES (1, 247);
INSERT INTO `sys_role_menu` VALUES (1, 6);
INSERT INTO `sys_role_menu` VALUES (1, 16);
INSERT INTO `sys_role_menu` VALUES (1, 17);
INSERT INTO `sys_role_menu` VALUES (1, 18);
INSERT INTO `sys_role_menu` VALUES (1, 60);
INSERT INTO `sys_role_menu` VALUES (1, 7);
INSERT INTO `sys_role_menu` VALUES (1, 21);
INSERT INTO `sys_role_menu` VALUES (1, 22);
INSERT INTO `sys_role_menu` VALUES (1, 23);
INSERT INTO `sys_role_menu` VALUES (1, 24);
INSERT INTO `sys_role_menu` VALUES (1, 8);
INSERT INTO `sys_role_menu` VALUES (1, 26);
INSERT INTO `sys_role_menu` VALUES (1, 27);
INSERT INTO `sys_role_menu` VALUES (1, 28);
INSERT INTO `sys_role_menu` VALUES (1, 10);
INSERT INTO `sys_role_menu` VALUES (1, 32);
INSERT INTO `sys_role_menu` VALUES (1, 33);
INSERT INTO `sys_role_menu` VALUES (1, 34);
INSERT INTO `sys_role_menu` VALUES (1, 369);
INSERT INTO `sys_role_menu` VALUES (1, 412);
INSERT INTO `sys_role_menu` VALUES (1, 4);

-- ----------------------------
-- Table structure for sys_session_recording
-- ----------------------------
DROP TABLE IF EXISTS `sys_session_recording`;
CREATE TABLE `sys_session_recording`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键ID\'',
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'会话唯一标识\'',
  `admin_id` bigint UNSIGNED NOT NULL COMMENT '\'操作用户ID\'',
  `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'用户名\'',
  `host_id` bigint UNSIGNED NOT NULL COMMENT '\'目标主机ID\'',
  `host_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'主机名称\'',
  `host_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'主机IP\'',
  `ssh_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'SSH登录用户\'',
  `start_time` datetime(3) NOT NULL COMMENT '\'会话开始时间\'',
  `end_time` datetime(3) NULL DEFAULT NULL COMMENT '\'会话结束时间\'',
  `duration` bigint NULL DEFAULT NULL COMMENT '\'会话时长(秒)\'',
  `terminal_width` bigint NULL DEFAULT 80 COMMENT '\'终端宽度\'',
  `terminal_height` bigint NULL DEFAULT 24 COMMENT '\'终端高度\'',
  `file_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'录制文件路径\'',
  `file_size` bigint NULL DEFAULT NULL COMMENT '\'文件大小(字节)\'',
  `storage_type` bigint NULL DEFAULT 1 COMMENT '\'存储类型:1-本地 2-OSS\'',
  `oss_key` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'OSS对象key\'',
  `input_count` bigint NULL DEFAULT 0 COMMENT '\'输入事件数量\'',
  `output_count` bigint NULL DEFAULT 0 COMMENT '\'输出事件数量\'',
  `resize_count` bigint NULL DEFAULT 0 COMMENT '\'窗口调整次数\'',
  `command_count` bigint NULL DEFAULT 0 COMMENT '\'命令执行次数(估算)\'',
  `client_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'客户端IP\'',
  `user_agent` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'用户代理\'',
  `risk_level` bigint NULL DEFAULT 0 COMMENT '\'风险等级:0-正常 1-可疑 2-高危\'',
  `has_sensitive_cmd` tinyint(1) NULL DEFAULT 0 COMMENT '\'是否包含敏感命令\'',
  `status` bigint NULL DEFAULT 1 COMMENT '\'状态:1-录制中 2-已完成 3-异常终止\'',
  `error_msg` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'错误信息\'',
  `create_time` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '\'更新时间\'',
  `delete_time` datetime(3) NULL DEFAULT NULL COMMENT '\'删除时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_session_recording_session_id`(`session_id` ASC) USING BTREE,
  INDEX `idx_sys_session_recording_host_id`(`host_id` ASC) USING BTREE,
  INDEX `idx_sys_session_recording_host_ip`(`host_ip` ASC) USING BTREE,
  INDEX `idx_sys_session_recording_start_time`(`start_time` ASC) USING BTREE,
  INDEX `idx_sys_session_recording_risk_level`(`risk_level` ASC) USING BTREE,
  INDEX `idx_sys_session_recording_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 528 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_session_recording
-- ----------------------------

-- ----------------------------
-- Table structure for system_migration_meta
-- ----------------------------
DROP TABLE IF EXISTS `system_migration_meta`;
CREATE TABLE `system_migration_meta`  (
  `meta_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `meta_value` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`meta_key`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of system_migration_meta
-- ----------------------------
INSERT INTO `system_migration_meta` VALUES ('alert_history_latest_in_source_backfilled_v1', 'done', '2026-05-20 09:20:17');
INSERT INTO `system_migration_meta` VALUES ('auto_migrate_models_v1', 'c5a0f66b5ed5d9571942869152e74beb176e5acaccb1399309537983013a0c4d', '2026-11-10 20:06:35');
INSERT INTO `system_migration_meta` VALUES ('autoops_docker_bootstrap_v1', 'complete', '2026-11-10 20:06:10');
INSERT INTO `system_migration_meta` VALUES ('default_knowledge_groups_v1', 'completed', '2026-11-10 21:12:16');

-- ----------------------------
-- Table structure for task_ansible
-- ----------------------------
DROP TABLE IF EXISTS `task_ansible`;
CREATE TABLE `task_ansible`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键ID\'',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'任务名称\'',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'任务描述\'',
  `type` bigint NOT NULL DEFAULT 1 COMMENT '\'任务类型:1-手动,2-Git,3-K8s\'',
  `git_repo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'Git仓库地址\'',
  `host_groups` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'主机分组JSON\'',
  `all_host_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'所有主机ID JSON数组\'',
  `global_vars` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'全局变量JSON\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'任务状态:1-等待中,2-运行中,3-成功,4-异常\'',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'错误信息\'',
  `task_count` bigint NOT NULL DEFAULT 0 COMMENT '\'任务数量(Type=1时为上传文件数,Type=2时为解析的playbook数,Type=3时固定为1)\'',
  `total_duration` bigint NOT NULL DEFAULT 0 COMMENT '\'任务执行总耗时(秒,所有子任务耗时总和)\'',
  `created_at` datetime(3) NOT NULL COMMENT '\'创建时间\'',
  `updated_at` datetime(3) NOT NULL COMMENT '\'更新时间\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_task_ansible_name`(`name` ASC) USING BTREE,
  INDEX `idx_task_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of task_ansible
-- ----------------------------

-- ----------------------------
-- Table structure for task_ansiblework
-- ----------------------------
DROP TABLE IF EXISTS `task_ansiblework`;
CREATE TABLE `task_ansiblework`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'主键ID\'',
  `task_id` bigint UNSIGNED NOT NULL COMMENT '\'父任务ID\'',
  `entry_file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'入口文件名\'',
  `entry_file_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'入口文件路径\'',
  `log_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'日志路径\'',
  `status` bigint NOT NULL DEFAULT 1 COMMENT '\'子任务状态:1-等待中,2-运行中,3-成功,4-异常\'',
  `start_time` datetime(3) NULL DEFAULT NULL COMMENT '\'开始时间\'',
  `end_time` datetime(3) NULL DEFAULT NULL COMMENT '\'结束时间\'',
  `duration` bigint NULL DEFAULT NULL COMMENT '\'执行耗时(秒)\'',
  `exit_code` bigint NULL DEFAULT NULL COMMENT '\'退出代码\'',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '\'错误信息\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_task_id`(`task_id` ASC) USING BTREE,
  INDEX `idx_task_work_composite`(`task_id` ASC, `status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of task_ansiblework
-- ----------------------------

-- ----------------------------
-- Table structure for task_job
-- ----------------------------
DROP TABLE IF EXISTS `task_job`;
CREATE TABLE `task_job`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '任务标题',
  `type` bigint NULL DEFAULT NULL COMMENT '任务类型 1=普通任务,2=定时任务,3=ansible任务',
  `shell` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '任务内容(任务模板ID,多个用逗号分隔)',
  `host_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '主机ID(多个用逗号分隔)',
  `cron_expr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '定时表达式(* * * * *)',
  `tasklog` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '任务执行日志',
  `status` bigint NULL DEFAULT NULL COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常,5=已暂停',
  `duration` bigint NULL DEFAULT NULL COMMENT '执行耗时(秒)',
  `task_count` bigint NULL DEFAULT NULL COMMENT '任务数量',
  `execute_count` bigint NULL DEFAULT 0 COMMENT '执行次数',
  `next_run_time` datetime(3) NULL DEFAULT NULL COMMENT '下次执行时间',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '任务备注',
  `start_time` datetime(3) NULL DEFAULT NULL COMMENT '任务开始时间',
  `end_time` datetime(3) NULL DEFAULT NULL COMMENT '任务结束时间',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '任务创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_task_job_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of task_job
-- ----------------------------

-- ----------------------------
-- Table structure for task_template
-- ----------------------------
DROP TABLE IF EXISTS `task_template`;
CREATE TABLE `task_template`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` bigint NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `updated_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_task_template_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of task_template
-- ----------------------------

-- ----------------------------
-- Table structure for task_work
-- ----------------------------
DROP TABLE IF EXISTS `task_work`;
CREATE TABLE `task_work`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '关联的任务ID',
  `template_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '任务模板ID',
  `host_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '执行主机ID',
  `type` bigint NULL DEFAULT NULL COMMENT '任务类型 1=普通任务,2=定时任务',
  `status` bigint NULL DEFAULT NULL COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常',
  `log` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '任务日志',
  `log_path` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '日志文件路径',
  `start_time` datetime(3) NULL DEFAULT NULL COMMENT '任务开始时间',
  `end_time` datetime(3) NULL DEFAULT NULL COMMENT '任务结束时间',
  `duration` bigint NULL DEFAULT NULL COMMENT '执行耗时(秒)',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `scheduled_time` datetime(3) NULL DEFAULT NULL COMMENT '定时任务执行时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_task_work_task_id`(`task_id` ASC) USING BTREE,
  INDEX `idx_task_work_template_id`(`template_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of task_work
-- ----------------------------

-- ----------------------------
-- Table structure for ticket_actions
-- ----------------------------
DROP TABLE IF EXISTS `ticket_actions`;
CREATE TABLE `ticket_actions`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint UNSIGNED NOT NULL,
  `operator_id` bigint UNSIGNED NOT NULL,
  `operator_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `action` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `from_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `to_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `comment` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ticket_actions_ticket_id`(`ticket_id` ASC) USING BTREE,
  INDEX `idx_ticket_actions_operator_id`(`operator_id` ASC) USING BTREE,
  INDEX `idx_ticket_actions_action`(`action` ASC) USING BTREE,
  INDEX `idx_ticket_actions_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ticket_actions
-- ----------------------------

-- ----------------------------
-- Table structure for ticket_approval_tasks
-- ----------------------------
DROP TABLE IF EXISTS `ticket_approval_tasks`;
CREATE TABLE `ticket_approval_tasks`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `instance_id` bigint UNSIGNED NOT NULL,
  `node_id` bigint UNSIGNED NOT NULL,
  `node_order` bigint NOT NULL,
  `node_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `assignee_id` bigint UNSIGNED NOT NULL,
  `assignee_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `comment` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `operator_id` bigint UNSIGNED NULL DEFAULT NULL,
  `operator_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `acted_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ticket_task_instance_status`(`instance_id` ASC, `status` ASC) USING BTREE,
  INDEX `idx_ticket_approval_tasks_node_id`(`node_id` ASC) USING BTREE,
  INDEX `idx_ticket_task_assignee_status`(`assignee_id` ASC, `status` ASC) USING BTREE,
  INDEX `idx_ticket_approval_tasks_operator_id`(`operator_id` ASC) USING BTREE,
  INDEX `idx_ticket_approval_tasks_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ticket_approval_tasks
-- ----------------------------

-- ----------------------------
-- Table structure for ticket_orders
-- ----------------------------
DROP TABLE IF EXISTS `ticket_orders`;
CREATE TABLE `ticket_orders`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `title` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `business_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `category` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `environment` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `priority` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `impact_scope` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `form_data` json NULL,
  `attachments` json NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `workflow_instance_id` bigint UNSIGNED NULL DEFAULT NULL,
  `applicant_id` bigint UNSIGNED NOT NULL,
  `applicant_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `handler_id` bigint UNSIGNED NULL DEFAULT NULL,
  `handler_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `strategy` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `rollback_plan` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `executor_id` bigint UNSIGNED NULL DEFAULT NULL,
  `executor_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `expected_at` datetime(3) NULL DEFAULT NULL,
  `resolved_at` datetime(3) NULL DEFAULT NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `closed_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ticket_orders_ticket_no`(`ticket_no` ASC) USING BTREE,
  INDEX `idx_ticket_orders_title`(`title` ASC) USING BTREE,
  INDEX `idx_ticket_type_status`(`business_type` ASC) USING BTREE,
  INDEX `idx_ticket_orders_category`(`category` ASC) USING BTREE,
  INDEX `idx_ticket_orders_environment`(`environment` ASC) USING BTREE,
  INDEX `idx_ticket_orders_priority`(`priority` ASC) USING BTREE,
  INDEX `idx_ticket_orders_status`(`status` ASC) USING BTREE,
  INDEX `idx_ticket_orders_workflow_instance_id`(`workflow_instance_id` ASC) USING BTREE,
  INDEX `idx_ticket_orders_applicant_id`(`applicant_id` ASC) USING BTREE,
  INDEX `idx_ticket_orders_applicant_name`(`applicant_name` ASC) USING BTREE,
  INDEX `idx_ticket_orders_handler_id`(`handler_id` ASC) USING BTREE,
  INDEX `idx_ticket_orders_handler_name`(`handler_name` ASC) USING BTREE,
  INDEX `idx_ticket_orders_executor_id`(`executor_id` ASC) USING BTREE,
  INDEX `idx_ticket_orders_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ticket_orders
-- ----------------------------

-- ----------------------------
-- Table structure for ticket_release_applications
-- ----------------------------
DROP TABLE IF EXISTS `ticket_release_applications`;
CREATE TABLE `ticket_release_applications`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `system_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `environment` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `jenkins_account_id` bigint UNSIGNED NOT NULL,
  `jenkins_job_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `default_parameters` json NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ticket_release_applications_code`(`code` ASC) USING BTREE,
  INDEX `idx_ticket_release_applications_name`(`name` ASC) USING BTREE,
  INDEX `idx_ticket_release_applications_system_name`(`system_name` ASC) USING BTREE,
  INDEX `idx_ticket_release_applications_environment`(`environment` ASC) USING BTREE,
  INDEX `idx_ticket_release_applications_jenkins_account_id`(`jenkins_account_id` ASC) USING BTREE,
  INDEX `idx_ticket_release_applications_enabled`(`enabled` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ticket_release_applications
-- ----------------------------

-- ----------------------------
-- Table structure for ticket_release_items
-- ----------------------------
DROP TABLE IF EXISTS `ticket_release_items`;
CREATE TABLE `ticket_release_items`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `release_id` bigint UNSIGNED NOT NULL,
  `application_id` bigint UNSIGNED NOT NULL,
  `application_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `version` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `target` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `jenkins_account_id` bigint UNSIGNED NOT NULL,
  `jenkins_job_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `parameters` json NULL,
  `execute_order` bigint NOT NULL DEFAULT 0,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `queue_id` bigint NULL DEFAULT NULL,
  `build_number` bigint NULL DEFAULT NULL,
  `build_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ticket_release_items_release_id`(`release_id` ASC) USING BTREE,
  INDEX `idx_ticket_release_items_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_ticket_release_items_jenkins_account_id`(`jenkins_account_id` ASC) USING BTREE,
  INDEX `idx_ticket_release_items_status`(`status` ASC) USING BTREE,
  INDEX `idx_ticket_release_items_queue_id`(`queue_id` ASC) USING BTREE,
  INDEX `idx_ticket_release_items_build_number`(`build_number` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ticket_release_items
-- ----------------------------

-- ----------------------------
-- Table structure for ticket_workflow_definitions
-- ----------------------------
DROP TABLE IF EXISTS `ticket_workflow_definitions`;
CREATE TABLE `ticket_workflow_definitions`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `business_types` json NULL,
  `environments` json NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `enabled` tinyint NOT NULL DEFAULT 1,
  `version` bigint NOT NULL DEFAULT 1,
  `prevent_self_approval` tinyint NOT NULL DEFAULT 1,
  `maintainer_id` bigint UNSIGNED NULL DEFAULT NULL,
  `maintainer_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `nodes` json NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ticket_workflow_definitions_code`(`code` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_definitions_name`(`name` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_definitions_enabled`(`enabled` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_definitions_maintainer_id`(`maintainer_id` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_definitions_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ticket_workflow_definitions
-- ----------------------------

-- ----------------------------
-- Table structure for ticket_workflow_instances
-- ----------------------------
DROP TABLE IF EXISTS `ticket_workflow_instances`;
CREATE TABLE `ticket_workflow_instances`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `definition_id` bigint UNSIGNED NOT NULL,
  `definition_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `definition_version` bigint NOT NULL,
  `business_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `business_id` bigint UNSIGNED NOT NULL,
  `business_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `applicant_id` bigint UNSIGNED NOT NULL,
  `applicant_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `current_node_order` bigint NOT NULL DEFAULT 0,
  `current_node_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `node_snapshot` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `started_at` datetime(3) NULL DEFAULT NULL,
  `finished_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ticket_flow_business`(`business_type` ASC, `business_id` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_instances_definition_id`(`definition_id` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_instances_business_no`(`business_no` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_instances_applicant_id`(`applicant_id` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_instances_status`(`status` ASC) USING BTREE,
  INDEX `idx_ticket_workflow_instances_started_at`(`started_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ticket_workflow_instances
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
