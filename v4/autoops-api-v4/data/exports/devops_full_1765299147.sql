-- Database: devops
-- Export Time: 2025-12-10 00:52:28
-- Tables: 44
-- Export Type: full

-- Table structure for app_application
CREATE TABLE `app_application` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(255) NOT NULL,
 `code` varchar(255) NOT NULL,
 `business_group_id` bigint(20) NOT NULL,
 `business_dept_id` bigint(20) NOT NULL,
 `description` text,
 `repo_url` varchar(500),
 `dev_owners` json,
 `test_owners` json,
 `ops_owners` json,
 `programming_lang` varchar(100),
 `start_command` text,
 `stop_command` text,
 `health_api` varchar(500),
 `domains` json,
 `hosts` json,
 `databases` json,
 `other_res` json,
 `status` tinyint(3) DEFAULT 1,
 `created_at` datetime,
 `updated_at` datetime,
 `deleted_at` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `app_application` ADD unique INDEX `idx_app_application_code`(`code`) USING BTREE;
ALTER TABLE `app_application` ADD  INDEX `idx_app_application_deleted_at`(`deleted_at`) USING BTREE;

-- Data for table app_application
INSERT INTO `app_application` (`id`,`name`,`code`,`business_group_id`,`business_dept_id`,`description`,`repo_url`,`dev_owners`,`test_owners`,`ops_owners`,`programming_lang`,`start_command`,`stop_command`,`health_api`,`domains`,`hosts`,`databases`,`other_res`,`status`,`created_at`,`updated_at`,`deleted_at`) VALUES (14,'dev-lku-sass-pack','dev-lku-sass-pack',12,2,'dev-lku-sass-pack   服务镜像打包','git@code.dding.net/lockin/saas-toc-server.git','[98]','[98, 103, 104]','[98, 103, 104]','Java','','','',NULL,NULL,NULL,'{}',2,'2025-11-14 19:27:43','2025-11-24 14:59:49',NULL);
INSERT INTO `app_application` (`id`,`name`,`code`,`business_group_id`,`business_dept_id`,`description`,`repo_url`,`dev_owners`,`test_owners`,`ops_owners`,`programming_lang`,`start_command`,`stop_command`,`health_api`,`domains`,`hosts`,`databases`,`other_res`,`status`,`created_at`,`updated_at`,`deleted_at`) VALUES (15,'dev-lku-sass-deploy','dev-lku-sass-deploy',12,3,'dev-lku-sass-deploy k8s 发布','git@code.dding.net/lockin/saas-toc-server.git','[102, 98, 89]','[98, 102, 89]','[89]','Go','','','',NULL,NULL,NULL,'{}',2,'2025-11-14 19:30:29','2025-11-21 23:09:28',NULL);
INSERT INTO `app_application` (`id`,`name`,`code`,`business_group_id`,`business_dept_id`,`description`,`repo_url`,`dev_owners`,`test_owners`,`ops_owners`,`programming_lang`,`start_command`,`stop_command`,`health_api`,`domains`,`hosts`,`databases`,`other_res`,`status`,`created_at`,`updated_at`,`deleted_at`) VALUES (17,'test1-paas-pack','test1-paas-pack',12,2,'test1-paas-pack','git@code.dding.net/lockin/cloud-platform.git','[102]','[98]','[102]','Node.js','','','',NULL,NULL,NULL,'{}',2,'2025-11-21 23:10:42','2025-11-24 14:59:30',NULL);
INSERT INTO `app_application` (`id`,`name`,`code`,`business_group_id`,`business_dept_id`,`description`,`repo_url`,`dev_owners`,`test_owners`,`ops_owners`,`programming_lang`,`start_command`,`stop_command`,`health_api`,`domains`,`hosts`,`databases`,`other_res`,`status`,`created_at`,`updated_at`,`deleted_at`) VALUES (18,'test1-paas-deploy','test1-paas-deploy',12,2,'test1-paas-deploy','git@code.dding.net/lockin/cloud-platform.git','[102]','[98]','[103]','Node.js','','','',NULL,NULL,NULL,'{}',2,'2025-11-21 23:11:31','2025-11-24 14:59:57',NULL);
INSERT INTO `app_application` (`id`,`name`,`code`,`business_group_id`,`business_dept_id`,`description`,`repo_url`,`dev_owners`,`test_owners`,`ops_owners`,`programming_lang`,`start_command`,`stop_command`,`health_api`,`domains`,`hosts`,`databases`,`other_res`,`status`,`created_at`,`updated_at`,`deleted_at`) VALUES (19,'prod_saas3.0_data-export','prod-saas30-data-export',18,2,'prod_saas3.0_data-export','git@gitee.com:zhang_fan1024/zf-k8s.git','[104]','[103]','[89]','Java','','','',NULL,NULL,NULL,'{}',2,'2025-11-28 12:39:27','2025-12-05 20:43:02',NULL);

-- Table structure for app_jenkins_env
CREATE TABLE `app_jenkins_env` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `app_id` bigint(20) NOT NULL,
 `env_name` varchar(50) NOT NULL,
 `jenkins_server_id` bigint(20),
 `job_name` varchar(255),
 `job_url` varchar(500),
 `build_params` json,
 `deploy_config` json,
 `notification` json,
 `is_active` tinyint(3) DEFAULT 1,
 `created_at` datetime,
 `updated_at` datetime,
 `deleted_at` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `app_jenkins_env` ADD  INDEX `idx_app_jenkins_env_app_id`(`app_id`) USING BTREE;
ALTER TABLE `app_jenkins_env` ADD  INDEX `idx_app_jenkins_env_deleted_at`(`deleted_at`) USING BTREE;

-- Data for table app_jenkins_env
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (35,14,'prod',13,'dev-lku-sass-pack',NULL,NULL,NULL,NULL,1,'2025-11-14 19:27:43','2025-11-22 18:11:53',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (36,14,'test',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-14 19:27:43','2025-11-14 19:27:43',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (37,14,'dev',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-14 19:27:43','2025-11-14 19:27:43',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (38,14,'test2',14,'dev-lku-sass-pack',NULL,NULL,NULL,NULL,1,'2025-11-14 19:28:59','2025-11-14 19:28:59',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (39,15,'prod',13,'test1-saas-deploy',NULL,NULL,NULL,NULL,1,'2025-11-14 19:30:29','2025-11-21 23:06:29',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (40,15,'test',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-14 19:30:30','2025-11-14 19:30:30',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (41,15,'dev',14,'dev-lku-sass-deploy',NULL,NULL,NULL,NULL,1,'2025-11-14 19:30:30','2025-11-14 19:31:00',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (45,17,'prod',13,'test1-paas-pack',NULL,NULL,NULL,NULL,1,'2025-11-21 23:10:43','2025-11-21 23:12:06',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (46,17,'test',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-21 23:10:43','2025-11-21 23:10:43',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (47,17,'dev',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-21 23:10:43','2025-11-21 23:10:43',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (48,18,'prod',13,'test1-paas-deploy',NULL,NULL,NULL,NULL,1,'2025-11-21 23:11:32','2025-11-21 23:12:28',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (49,18,'test',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-21 23:11:32','2025-11-21 23:11:32',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (50,18,'dev',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-21 23:11:32','2025-11-21 23:11:32',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (51,19,'prod',13,'prod_saas3.0_data-export',NULL,NULL,NULL,NULL,1,'2025-11-28 12:39:27','2025-11-28 12:39:54',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (52,19,'test',13,'prod_saas3.0_data-export',NULL,NULL,NULL,NULL,1,'2025-11-28 12:39:28','2025-12-05 20:43:01',NULL);
INSERT INTO `app_jenkins_env` (`id`,`app_id`,`env_name`,`jenkins_server_id`,`job_name`,`job_url`,`build_params`,`deploy_config`,`notification`,`is_active`,`created_at`,`updated_at`,`deleted_at`) VALUES (53,19,'dev',NULL,'',NULL,NULL,NULL,NULL,1,'2025-11-28 12:39:28','2025-11-28 12:39:28',NULL);

-- Table structure for app_service_release
CREATE TABLE `app_service_release` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `title` varchar(255) NOT NULL COMMENT '发布标题',
 `business_group_id` bigint(20) NOT NULL COMMENT '业务组ID',
 `impact_feature` text COMMENT '影响功能描述',
 `applicant_id` bigint(20) NOT NULL COMMENT '申请人ID',
 `applicant_name` varchar(100) NOT NULL COMMENT '申请人姓名',
 `owner_approver_id` bigint(20) COMMENT '负责人审批人ID',
 `owner_approver_name` varchar(100) COMMENT '负责人审批人姓名',
 `security_approver_id` bigint(20) COMMENT '安全审批人ID',
 `security_approver_name` varchar(100) COMMENT '安全审批人姓名',
 `test_approver_id` bigint(20) COMMENT '测试审批人ID',
 `test_approver_name` varchar(100) COMMENT '测试审批人姓名',
 `owner_approval_status` bigint(19) DEFAULT 1 COMMENT '负责人审批状态',
 `security_approval_status` bigint(19) DEFAULT 1 COMMENT '安全审批状态',
 `test_approval_status` bigint(19) DEFAULT 1 COMMENT '测试审批状态',
 `owner_approval_time` datetime COMMENT '负责人审批时间',
 `security_approval_time` datetime COMMENT '安全审批时间',
 `test_approval_time` datetime COMMENT '测试审批时间',
 `owner_approval_remark` text COMMENT '负责人审批意见',
 `security_approval_remark` text COMMENT '安全审批意见',
 `test_approval_remark` text COMMENT '测试审批意见',
 `deploy_status` bigint(19) DEFAULT 1 COMMENT '运维发布状态',
 `regression_test_status` bigint(19) DEFAULT 1 COMMENT '回归测试状态',
 `status` bigint(19) DEFAULT 1 COMMENT '流程状态',
 `start_time` datetime COMMENT '发布开始时间',
 `end_time` datetime COMMENT '发布结束时间',
 `duration` bigint(19) DEFAULT 0 COMMENT '发布耗时(秒)',
 `service_count` bigint(19) DEFAULT 0 COMMENT '关联服务数量',
 `created_at` datetime,
 `updated_at` datetime,
 `deleted_at` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `app_service_release` ADD  INDEX `idx_app_service_release_deleted_at`(`deleted_at`) USING BTREE;

-- Data for table app_service_release
INSERT INTO `app_service_release` (`id`,`title`,`business_group_id`,`impact_feature`,`applicant_id`,`applicant_name`,`owner_approver_id`,`owner_approver_name`,`security_approver_id`,`security_approver_name`,`test_approver_id`,`test_approver_name`,`owner_approval_status`,`security_approval_status`,`test_approval_status`,`owner_approval_time`,`security_approval_time`,`test_approval_time`,`owner_approval_remark`,`security_approval_remark`,`test_approval_remark`,`deploy_status`,`regression_test_status`,`status`,`start_time`,`end_time`,`duration`,`service_count`,`created_at`,`updated_at`,`deleted_at`) VALUES (2,'测试---发布工单上线',12,'影响用户登录',1,'管理员',102,'用户_102',98,'用户_98',89,'用户_89',2,2,2,'2025-11-22 18:16:13','2025-11-22 18:16:27','2025-11-22 18:16:36','同意','同意','同意',3,3,2,'2025-11-22 18:18:03','2025-11-22 18:21:45',220,1,'2025-11-22 18:13:37','2025-11-22 23:01:55',NULL);
INSERT INTO `app_service_release` (`id`,`title`,`business_group_id`,`impact_feature`,`applicant_id`,`applicant_name`,`owner_approver_id`,`owner_approver_name`,`security_approver_id`,`security_approver_name`,`test_approver_id`,`test_approver_name`,`owner_approval_status`,`security_approval_status`,`test_approval_status`,`owner_approval_time`,`security_approval_time`,`test_approval_time`,`owner_approval_remark`,`security_approval_remark`,`test_approval_remark`,`deploy_status`,`regression_test_status`,`status`,`start_time`,`end_time`,`duration`,`service_count`,`created_at`,`updated_at`,`deleted_at`) VALUES (3,'测试多任务工单发布',12,'测试多任务工单发布',1,'管理员',103,'用户_103',98,'用户_98',104,'用户_104',2,2,2,'2025-11-22 19:33:08','2025-11-22 19:34:25','2025-11-22 19:34:41','','','',4,1,3,'2025-11-22 19:36:09','2025-11-22 19:41:18',591,4,'2025-11-22 19:31:57','2025-11-22 19:41:18',NULL);
INSERT INTO `app_service_release` (`id`,`title`,`business_group_id`,`impact_feature`,`applicant_id`,`applicant_name`,`owner_approver_id`,`owner_approver_name`,`security_approver_id`,`security_approver_name`,`test_approver_id`,`test_approver_name`,`owner_approval_status`,`security_approval_status`,`test_approval_status`,`owner_approval_time`,`security_approval_time`,`test_approval_time`,`owner_approval_remark`,`security_approval_remark`,`test_approval_remark`,`deploy_status`,`regression_test_status`,`status`,`start_time`,`end_time`,`duration`,`service_count`,`created_at`,`updated_at`,`deleted_at`) VALUES (4,'测试审批流程',12,'测试审批流程',1,'管理员',103,'用户_103',98,'用户_98',104,'用户_104',2,2,2,'2025-11-22 20:28:35','2025-11-22 20:28:48','2025-11-22 20:28:55','','','',4,1,3,'2025-11-22 20:29:16','2025-11-22 20:34:22',303,2,'2025-11-22 20:00:18','2025-11-22 20:34:22',NULL);
INSERT INTO `app_service_release` (`id`,`title`,`business_group_id`,`impact_feature`,`applicant_id`,`applicant_name`,`owner_approver_id`,`owner_approver_name`,`security_approver_id`,`security_approver_name`,`test_approver_id`,`test_approver_name`,`owner_approval_status`,`security_approval_status`,`test_approval_status`,`owner_approval_time`,`security_approval_time`,`test_approval_time`,`owner_approval_remark`,`security_approval_remark`,`test_approval_remark`,`deploy_status`,`regression_test_status`,`status`,`start_time`,`end_time`,`duration`,`service_count`,`created_at`,`updated_at`,`deleted_at`) VALUES (5,'测试---发布工单上线',12,'测试---发布工单上线',1,'管理员',104,'用户_104',103,'用户_103',98,'用户_98',2,2,2,'2025-11-22 23:13:56','2025-11-22 23:16:26','2025-11-22 23:16:34','','通过','',3,3,2,'2025-11-22 23:52:17','2025-11-22 23:55:51',212,1,'2025-11-22 23:13:39','2025-11-23 00:50:11',NULL);
INSERT INTO `app_service_release` (`id`,`title`,`business_group_id`,`impact_feature`,`applicant_id`,`applicant_name`,`owner_approver_id`,`owner_approver_name`,`security_approver_id`,`security_approver_name`,`test_approver_id`,`test_approver_name`,`owner_approval_status`,`security_approval_status`,`test_approval_status`,`owner_approval_time`,`security_approval_time`,`test_approval_time`,`owner_approval_remark`,`security_approval_remark`,`test_approval_remark`,`deploy_status`,`regression_test_status`,`status`,`start_time`,`end_time`,`duration`,`service_count`,`created_at`,`updated_at`,`deleted_at`) VALUES (6,'双任务上线',14,'123',1,'管理员',104,'用户_104',103,'用户_103',98,'用户_98',2,2,2,'2025-11-23 11:33:16','2025-11-23 11:33:22','2025-11-23 12:57:50','','','111',1,1,1,NULL,NULL,0,2,'2025-11-23 11:32:19','2025-11-23 12:57:50',NULL);
INSERT INTO `app_service_release` (`id`,`title`,`business_group_id`,`impact_feature`,`applicant_id`,`applicant_name`,`owner_approver_id`,`owner_approver_name`,`security_approver_id`,`security_approver_name`,`test_approver_id`,`test_approver_name`,`owner_approval_status`,`security_approval_status`,`test_approval_status`,`owner_approval_time`,`security_approval_time`,`test_approval_time`,`owner_approval_remark`,`security_approval_remark`,`test_approval_remark`,`deploy_status`,`regression_test_status`,`status`,`start_time`,`end_time`,`duration`,`service_count`,`created_at`,`updated_at`,`deleted_at`) VALUES (7,'测试生产环境服务上线',18,'影响用户登录',1,'管理员',104,'王五',98,'李四',98,'李四',2,2,2,'2025-11-24 00:44:47','2025-11-24 00:45:40','2025-11-24 00:45:45','同意','','',3,3,2,'2025-11-24 00:46:52','2025-11-24 00:52:38',344,1,'2025-11-24 00:42:24','2025-11-24 00:53:04',NULL);
INSERT INTO `app_service_release` (`id`,`title`,`business_group_id`,`impact_feature`,`applicant_id`,`applicant_name`,`owner_approver_id`,`owner_approver_name`,`security_approver_id`,`security_approver_name`,`test_approver_id`,`test_approver_name`,`owner_approval_status`,`security_approval_status`,`test_approval_status`,`owner_approval_time`,`security_approval_time`,`test_approval_time`,`owner_approval_remark`,`security_approval_remark`,`test_approval_remark`,`deploy_status`,`regression_test_status`,`status`,`start_time`,`end_time`,`duration`,`service_count`,`created_at`,`updated_at`,`deleted_at`) VALUES (8,'tesst123',12,'影响用户登录',1,'管理员',89,'admin',89,'admin',89,'admin',2,2,2,'2025-12-05 20:47:13','2025-12-05 20:47:18','2025-12-05 20:47:27','ok','','',3,2,1,'2025-12-05 20:47:39','2025-12-05 20:47:53',12,1,'2025-12-05 20:46:53','2025-12-05 20:47:53',NULL);

-- Table structure for app_service_release_item
CREATE TABLE `app_service_release_item` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `release_id` bigint(20) NOT NULL COMMENT '上线工单ID',
 `app_id` bigint(20) NOT NULL COMMENT '应用ID',
 `app_name` varchar(255) NOT NULL COMMENT '应用名称',
 `app_code` varchar(100) NOT NULL COMMENT '应用编码',
 `project_name` varchar(255) NOT NULL COMMENT '项目服务名称',
 `repo_url` varchar(500) NOT NULL COMMENT '项目地址',
 `branch` varchar(100) DEFAULT master COMMENT '发布分支',
 `commit_id` varchar(100) NOT NULL COMMENT 'Commit ID',
 `impact_feature` text COMMENT '影响功能',
 `function_module` text COMMENT '功能模块',
 `db_change` text COMMENT '数据库变更',
 `config_change` text COMMENT '配置变更',
 `remark` text COMMENT '备注信息',
 `jenkins_env_id` bigint(20) COMMENT 'Jenkins环境配置ID',
 `jenkins_job_url` varchar(500) COMMENT 'Jenkins任务URL',
 `parameters` text COMMENT 'Jenkins构建参数(JSON格式)',
 `build_number` bigint(19) DEFAULT 0 COMMENT '构建编号',
 `log_url` varchar(500) COMMENT '构建日志URL',
 `status` bigint(19) DEFAULT 1 COMMENT '发布状态',
 `start_time` datetime COMMENT '发布开始时间',
 `end_time` datetime COMMENT '发布结束时间',
 `duration` bigint(19) DEFAULT 0 COMMENT '发布耗时(秒)',
 `error_message` text COMMENT '错误信息',
 `execute_order` bigint(19) DEFAULT 0 COMMENT '执行顺序',
 `created_at` datetime,
 `updated_at` datetime,
 `deleted_at` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `app_service_release_item` ADD  INDEX `idx_app_service_release_item_deleted_at`(`deleted_at`) USING BTREE;
ALTER TABLE `app_service_release_item` ADD  INDEX `idx_app_service_release_item_release_id`(`release_id`) USING BTREE;

-- Data for table app_service_release_item
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (1,2,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','影响用户登录','','没有','没有','没有',35,'','{"commit_id":"083a2d6669effe63670b0d5ee899d29eeb187abf","compile":"true"}',408,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/408/console',3,'2025-11-22 18:18:03','2025-11-22 18:21:44',220,'',1,'2025-11-22 18:13:37','2025-11-22 18:21:44',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (2,3,15,'dev-lku-sass-deploy','dev-lku-sass-deploy','dev-lku-sass-deploy','git@code.dding.net/lockin/saas-toc-server.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',39,'','',19,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-saas-deploy/19/console',4,'2025-11-22 19:36:10','2025-11-22 19:36:37',27,'Jenkins构建失败',1,'2025-11-22 19:31:58','2025-11-22 19:36:37',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (3,3,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',35,'','',409,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/409/console',3,'2025-11-22 19:36:10','2025-11-22 19:39:59',229,'',2,'2025-11-22 19:31:58','2025-11-22 19:39:59',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (4,3,18,'test1-paas-deploy','test1-paas-deploy','test1-paas-deploy','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',48,'','',13,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-deploy/13/console',4,'2025-11-22 19:36:10','2025-11-22 19:36:38',28,'Jenkins构建失败',3,'2025-11-22 19:31:59','2025-11-22 19:36:38',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (5,3,17,'test1-paas-pack','test1-paas-pack','test1-paas-pack','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',45,'','',14,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-pack/14/console',3,'2025-11-22 19:36:10','2025-11-22 19:41:18',307,'',4,'2025-11-22 19:31:59','2025-11-22 19:41:18',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (6,4,17,'test1-paas-pack','test1-paas-pack','test1-paas-pack','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试审批流程','','测试审批流程','测试审批流程','测试审批流程',45,'','',15,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-pack/15/console',3,'2025-11-22 20:29:17','2025-11-22 20:34:08',291,'',1,'2025-11-22 20:00:19','2025-11-22 20:34:08',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (7,4,18,'test1-paas-deploy','test1-paas-deploy','test1-paas-deploy','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试审批流程','','测试审批流程','测试审批流程','测试审批流程',48,'','',14,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-deploy/14/console',4,'2025-11-22 20:34:09','2025-11-22 20:34:21',12,'Jenkins构建失败',2,'2025-11-22 20:00:19','2025-11-22 20:34:21',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (8,5,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','无','','无','无','无',35,'','',410,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/410/console',3,'2025-11-22 23:52:18','2025-11-22 23:55:51',212,'',1,'2025-11-22 23:13:40','2025-11-22 23:55:51',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (9,6,17,'test1-paas-pack','test1-paas-pack','test1-paas-pack','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','123','','123','123','123',45,'','{"commit_id":"c7417cbe54a11d89c819b4541316980dc1634687","compile":"true"}',0,'',1,NULL,NULL,0,'',1,'2025-11-23 11:32:20','2025-11-23 11:32:20',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (10,6,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','123','','123','123','123',35,'','{"commit_id":"083a2d6669effe63670b0d5ee899d29eeb187abf","compile":"true"}',0,'',1,NULL,NULL,0,'',2,'2025-11-23 11:32:20','2025-11-23 11:32:20',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (11,7,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','影响用户登录','','无','无','123',35,'','{"commit_id":"083a2d6669effe63670b0d5ee899d29eeb187abf","compile":"true"}',411,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/411/console',3,'2025-11-24 00:46:53','2025-11-24 00:52:37',344,'',1,'2025-11-24 00:42:25','2025-11-24 00:52:37',NULL);
INSERT INTO `app_service_release_item` (`id`,`release_id`,`app_id`,`app_name`,`app_code`,`project_name`,`repo_url`,`branch`,`commit_id`,`impact_feature`,`function_module`,`db_change`,`config_change`,`remark`,`jenkins_env_id`,`jenkins_job_url`,`parameters`,`build_number`,`log_url`,`status`,`start_time`,`end_time`,`duration`,`error_message`,`execute_order`,`created_at`,`updated_at`,`deleted_at`) VALUES (12,8,19,'prod_saas3.0_data-export','prod-saas30-data-export','prod_saas3.0_data-export','git@gitee.com:zhang_fan1024/zf-k8s.git','master','','影响用户登录','','无','无','无',51,'','{"commit_id":"123456789"}',10,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/10/console',3,'2025-12-05 20:47:40','2025-12-05 20:47:52',12,'',1,'2025-12-05 20:46:53','2025-12-05 20:47:53',NULL);

-- Table structure for app_sh_release
CREATE TABLE `app_sh_release` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `title` varchar(255) NOT NULL COMMENT '上线标题',
 `reason` text NOT NULL COMMENT '上线原因',
 `business_group_id` bigint(20) NOT NULL COMMENT '业务线ID',
 `app_id` bigint(20) NOT NULL COMMENT '服务ID',
 `app_name` varchar(255) NOT NULL COMMENT '服务名称',
 `app_code` varchar(100) NOT NULL COMMENT '服务编码',
 `applicant_id` bigint(20) NOT NULL COMMENT '申请人ID',
 `applicant_name` varchar(100) NOT NULL COMMENT '申请人姓名',
 `approver_id` bigint(20) COMMENT '审核人ID',
 `approver_name` varchar(100) COMMENT '审核人姓名',
 `executor_id` bigint(20) COMMENT '执行人ID',
 `executor_name` varchar(100) COMMENT '执行人姓名',
 `execute_dir` varchar(500) NOT NULL COMMENT '执行目录',
 `script_content` text NOT NULL COMMENT '脚本内容',
 `approval_status` bigint(19) DEFAULT 1 COMMENT '审核状态',
 `approval_time` datetime COMMENT '审核时间',
 `approval_remark` text COMMENT '审核意见',
 `execute_status` bigint(19) DEFAULT 1 COMMENT '执行状态',
 `status` bigint(19) DEFAULT 1 COMMENT '流程状态',
 `start_time` datetime COMMENT '脚本执行开始时间',
 `end_time` datetime COMMENT '脚本执行结束时间',
 `duration` bigint(19) DEFAULT 0 COMMENT '执行耗时(秒)',
 `jenkins_env_id` bigint(20) COMMENT 'Jenkins环境配置ID',
 `build_number` bigint(19) DEFAULT 0 COMMENT '构建编号',
 `log_url` varchar(500) COMMENT '构建日志URL',
 `error_message` text COMMENT '错误信息',
 `created_at` datetime,
 `updated_at` datetime,
 `deleted_at` datetime,
 `parameters` text COMMENT 'Jenkins构建参数(JSON格式)',
 `server_host_id` bigint(20) NOT NULL COMMENT '执行服务器主机ID(关联cmdb_host)',
 `pull_code_start_time` datetime COMMENT '拉取代码开始时间',
 `pull_code_end_time` datetime COMMENT '拉取代码结束时间',
 `script_output` longtext COMMENT '脚本执行输出', 
PRIMARY KEY (id)
);
ALTER TABLE `app_sh_release` ADD  INDEX `idx_app_sh_release_deleted_at`(`deleted_at`) USING BTREE;

-- Data for table app_sh_release
INSERT INTO `app_sh_release` (`id`,`title`,`reason`,`business_group_id`,`app_id`,`app_name`,`app_code`,`applicant_id`,`applicant_name`,`approver_id`,`approver_name`,`executor_id`,`executor_name`,`execute_dir`,`script_content`,`approval_status`,`approval_time`,`approval_remark`,`execute_status`,`status`,`start_time`,`end_time`,`duration`,`jenkins_env_id`,`build_number`,`log_url`,`error_message`,`created_at`,`updated_at`,`deleted_at`,`parameters`,`server_host_id`,`pull_code_start_time`,`pull_code_end_time`,`script_output`) VALUES (1,'prod_saas3.0_data-export-数据导出','prod_saas3.0_data-export-数据导出',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','ls',2,'2025-11-28 12:42:48','111',3,2,'2025-11-28 12:43:08','2025-11-28 12:43:24',14,51,3,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/3/console','','2025-11-28 12:41:29','2025-11-28 12:43:24',NULL,'{"commit_id":"123456789"}',0,NULL,NULL,NULL);
INSERT INTO `app_sh_release` (`id`,`title`,`reason`,`business_group_id`,`app_id`,`app_name`,`app_code`,`applicant_id`,`applicant_name`,`approver_id`,`approver_name`,`executor_id`,`executor_name`,`execute_dir`,`script_content`,`approval_status`,`approval_time`,`approval_remark`,`execute_status`,`status`,`start_time`,`end_time`,`duration`,`jenkins_env_id`,`build_number`,`log_url`,`error_message`,`created_at`,`updated_at`,`deleted_at`,`parameters`,`server_host_id`,`pull_code_start_time`,`pull_code_end_time`,`script_output`) VALUES (2,'测试001','测试001',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',98,'李四',89,'admin','/home/','ls  /root/',2,'2025-11-28 14:18:58','同意',3,1,NULL,NULL,0,51,4,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/4/console','','2025-11-28 13:37:27','2025-11-28 14:47:38',NULL,'{"commit_id":"123456789"}',0,'2025-11-28 14:47:22','2025-11-28 14:47:38',NULL);
INSERT INTO `app_sh_release` (`id`,`title`,`reason`,`business_group_id`,`app_id`,`app_name`,`app_code`,`applicant_id`,`applicant_name`,`approver_id`,`approver_name`,`executor_id`,`executor_name`,`execute_dir`,`script_content`,`approval_status`,`approval_time`,`approval_remark`,`execute_status`,`status`,`start_time`,`end_time`,`duration`,`jenkins_env_id`,`build_number`,`log_url`,`error_message`,`created_at`,`updated_at`,`deleted_at`,`parameters`,`server_host_id`,`pull_code_start_time`,`pull_code_end_time`,`script_output`) VALUES (3,'测试脚本执行','测试脚本执行',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','pwd
ls
hostname -I',2,'2025-11-28 15:19:15','111',6,2,'2025-11-28 15:24:27','2025-11-28 15:24:31',3,51,5,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/5/console','','2025-11-28 15:09:18','2025-11-28 15:24:31',NULL,'{"commit_id":"123456789"}',501,'2025-11-28 15:19:21','2025-11-28 15:19:34','bash: line 1: cd: /home/dingding/saas3-data-export/: No such file or directory
docker
elk
jdk11
jdk17
jdk18
luban-master
node
prometheus
snap
172.20.236.121 172.18.0.1 172.17.0.1 172.19.0.1 
');
INSERT INTO `app_sh_release` (`id`,`title`,`reason`,`business_group_id`,`app_id`,`app_name`,`app_code`,`applicant_id`,`applicant_name`,`approver_id`,`approver_name`,`executor_id`,`executor_name`,`execute_dir`,`script_content`,`approval_status`,`approval_time`,`approval_remark`,`execute_status`,`status`,`start_time`,`end_time`,`duration`,`jenkins_env_id`,`build_number`,`log_url`,`error_message`,`created_at`,`updated_at`,`deleted_at`,`parameters`,`server_host_id`,`pull_code_start_time`,`pull_code_end_time`,`script_output`) VALUES (4,'测试脚本002','测试脚本002',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','hostname
pwd
date',2,'2025-11-28 15:56:35','11',6,2,'2025-11-28 15:57:23','2025-11-28 15:57:27',3,51,6,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/6/console','','2025-11-28 15:56:13','2025-11-28 15:57:27',NULL,'{"commit_id":"123456789"}',501,'2025-11-28 15:56:44','2025-11-28 15:57:00','/root
Fri Nov 28 03:57:27 PM CST 2025
bash: line 1: cd: /home/dingding/saas3-data-export/: No such file or directory
');
INSERT INTO `app_sh_release` (`id`,`title`,`reason`,`business_group_id`,`app_id`,`app_name`,`app_code`,`applicant_id`,`applicant_name`,`approver_id`,`approver_name`,`executor_id`,`executor_name`,`execute_dir`,`script_content`,`approval_status`,`approval_time`,`approval_remark`,`execute_status`,`status`,`start_time`,`end_time`,`duration`,`jenkins_env_id`,`build_number`,`log_url`,`error_message`,`created_at`,`updated_at`,`deleted_at`,`parameters`,`server_host_id`,`pull_code_start_time`,`pull_code_end_time`,`script_output`) VALUES (5,'测试002','测试002',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/','pwd
ls
hostname
',2,'2025-11-28 16:12:04','ok',6,2,'2025-11-28 16:19:18','2025-11-28 16:19:22',3,51,7,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/7/console','','2025-11-28 16:10:22','2025-11-28 16:19:22',NULL,'{"commit_id":"123456789"}',501,'2025-11-28 16:12:11','2025-11-28 16:12:23','/home
devops
devops.tar.gz
go-ops
');
INSERT INTO `app_sh_release` (`id`,`title`,`reason`,`business_group_id`,`app_id`,`app_name`,`app_code`,`applicant_id`,`applicant_name`,`approver_id`,`approver_name`,`executor_id`,`executor_name`,`execute_dir`,`script_content`,`approval_status`,`approval_time`,`approval_remark`,`execute_status`,`status`,`start_time`,`end_time`,`duration`,`jenkins_env_id`,`build_number`,`log_url`,`error_message`,`created_at`,`updated_at`,`deleted_at`,`parameters`,`server_host_id`,`pull_code_start_time`,`pull_code_end_time`,`script_output`) VALUES (6,'test1111111111','test',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','pwd  
ls 
hostname',2,'2025-12-01 01:01:48','ok',6,2,'2025-12-01 01:02:36','2025-12-01 01:02:39',1,51,8,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/8/console','','2025-12-01 01:01:08','2025-12-01 01:02:39',NULL,'{"commit_id":"123456789"}',501,'2025-12-01 01:01:59','2025-12-01 01:02:16','/home/dingding/saas3-data-export
xlsx
go-ops
');

-- Table structure for cmdb_group
CREATE TABLE `cmdb_group` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '''父级分组ID''',
 `name` longtext NOT NULL COMMENT '''分组名称''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `remark` longtext COMMENT '''备注''',
 `update_time` datetime COMMENT '''更新时间''', 
PRIMARY KEY (id)
);

-- Data for table cmdb_group
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (1,0,'默认业务组','2025-07-10 11:02:07',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (2,0,'saas3业务线','2025-07-10 11:03:36',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (3,0,'saas4业务线','2025-07-10 11:03:56',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (4,11,'公寓业务','2025-07-10 11:05:45',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (9,11,'智能窗帘业务','2025-07-17 10:17:52',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (11,0,'saas5业务线','2025-07-17 10:21:44',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (12,2,'智能安防业务','2025-07-17 11:07:19',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (14,3,'智慧门锁业务','2025-07-17 11:19:14',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (15,0,'saas6业务组','2025-07-17 12:31:33',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (16,15,'智能家居业务','2025-07-22 10:32:14',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (17,0,'ops业务组','2025-07-23 10:39:44',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (18,2,'智控照明业务','2025-07-23 10:40:12',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (19,17,'运维组','2025-07-23 20:21:08',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (20,17,'安全组','2025-07-23 20:21:33',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (24,0,'lockin鹿客云','2025-09-07 17:49:17',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (26,24,'国内鹿客云','2025-09-07 18:04:56',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (27,24,'新加坡鹿客云','2025-09-07 18:05:13',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (28,24,'北美鹿客云','2025-09-07 18:05:31',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (29,17,'ai运维组','2025-10-01 15:17:35',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (30,0,'rental业务线','2025-12-05 20:17:08',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (31,30,'rental-test','2025-12-05 20:17:32',NULL,NULL);
INSERT INTO `cmdb_group` (`id`,`parent_id`,`name`,`create_time`,`remark`,`update_time`) VALUES (32,30,'rental-dev','2025-12-05 20:17:43',NULL,NULL);

-- Table structure for cmdb_host
CREATE TABLE `cmdb_host` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `host_name` longtext NOT NULL COMMENT '''名称''',
 `group_id` bigint(20) NOT NULL COMMENT '''分组ID''',
 `private_ip` longtext COMMENT '''私网IP''',
 `public_ip` longtext COMMENT '''公网IP''',
 `ssh_name` longtext COMMENT '''SSH用户名''',
 `ssh_key_id` bigint(20) COMMENT '''SSH凭据ID''',
 `ssh_port` bigint(19) DEFAULT 22 COMMENT '''SSH端口''',
 `remark` longtext COMMENT '''备注''',
 `vendor` bigint(19) COMMENT '''1->自建,2->阿里云,3->腾讯云''',
 `region` longtext COMMENT '''区域''',
 `instance_id` longtext COMMENT '''实例ID''',
 `os` longtext COMMENT '''操作系统''',
 `status` bigint(19) COMMENT '''状态:1->认证成功,2->未认证,3->认证失败''',
 `cpu` longtext COMMENT '''CPU信息''',
 `memory` longtext COMMENT '''内存信息''',
 `disk` longtext COMMENT '''磁盘信息''',
 `billing_type` longtext COMMENT '''计费方式''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `expire_time` datetime COMMENT '''到期时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `ssh_ip` longtext NOT NULL COMMENT '''SSH连接IP''',
 `name` longtext NOT NULL COMMENT '''ecs主机名称''',
 `ssh_gateway_id` bigint(20) COMMENT '''中转网关凭据ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `cmdb_host` ADD  INDEX `fk_cmdb_group_hosts`(`group_id`) USING BTREE;

-- Data for table cmdb_host
INSERT INTO `cmdb_host` (`id`,`host_name`,`group_id`,`private_ip`,`public_ip`,`ssh_name`,`ssh_key_id`,`ssh_port`,`remark`,`vendor`,`region`,`instance_id`,`os`,`status`,`cpu`,`memory`,`disk`,`billing_type`,`create_time`,`expire_time`,`update_time`,`ssh_ip`,`name`,`ssh_gateway_id`) VALUES (501,'8.130.14.34',12,'172.20.236.121','8.130.14.34','root',20,22,'123',2,'','','Ubuntu24.04.2',1,'2','2','40','','2025-10-18 16:14:25',NULL,'2025-12-05 20:18:59','8.130.14.34','go-ops',NULL);
INSERT INTO `cmdb_host` (`id`,`host_name`,`group_id`,`private_ip`,`public_ip`,`ssh_name`,`ssh_key_id`,`ssh_port`,`remark`,`vendor`,`region`,`instance_id`,`os`,`status`,`cpu`,`memory`,`disk`,`billing_type`,`create_time`,`expire_time`,`update_time`,`ssh_ip`,`name`,`ssh_gateway_id`) VALUES (506,'华为云ops',4,'172.31.6.35','139.9.205.38','root',22,22,'123',5,'','','Ubuntu24.04.2',1,'2','2','40','','2025-11-11 17:24:36',NULL,'2025-11-26 17:22:07','139.9.205.38','hw-ops',NULL);
INSERT INTO `cmdb_host` (`id`,`host_name`,`group_id`,`private_ip`,`public_ip`,`ssh_name`,`ssh_key_id`,`ssh_port`,`remark`,`vendor`,`region`,`instance_id`,`os`,`status`,`cpu`,`memory`,`disk`,`billing_type`,`create_time`,`expire_time`,`update_time`,`ssh_ip`,`name`,`ssh_gateway_id`) VALUES (510,'180.76.231.65',9,'192.168.16.2','180.76.231.65','root',20,22,'180.76.231.65',4,'','','Ubuntu22.04.1',1,'2','2','40','','2025-11-14 17:43:50',NULL,'2025-11-25 16:40:01','180.76.231.65','bd-git-cicd',NULL);
INSERT INTO `cmdb_host` (`id`,`host_name`,`group_id`,`private_ip`,`public_ip`,`ssh_name`,`ssh_key_id`,`ssh_port`,`remark`,`vendor`,`region`,`instance_id`,`os`,`status`,`cpu`,`memory`,`disk`,`billing_type`,`create_time`,`expire_time`,`update_time`,`ssh_ip`,`name`,`ssh_gateway_id`) VALUES (511,'虚拟机001',12,'172.16.226.16','120.231.244.158','root',13,22,'123',1,'','','CentOSLinux7(Core)',1,'2','4','17','','2025-11-23 23:44:35',NULL,'2025-11-23 23:45:32','172.16.226.16','jenkins',NULL);
INSERT INTO `cmdb_host` (`id`,`host_name`,`group_id`,`private_ip`,`public_ip`,`ssh_name`,`ssh_key_id`,`ssh_port`,`remark`,`vendor`,`region`,`instance_id`,`os`,`status`,`cpu`,`memory`,`disk`,`billing_type`,`create_time`,`expire_time`,`update_time`,`ssh_ip`,`name`,`ssh_gateway_id`) VALUES (512,'虚拟机002',4,'172.16.226.13','120.231.244.158','root',13,22,'123',1,'','','Ubuntu20.04.2',1,'2','3','19','','2025-11-23 23:46:03',NULL,'2025-11-23 23:46:07','172.16.226.13','k8s-node02',NULL);

-- Table structure for cmdb_sql
CREATE TABLE `cmdb_sql` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(100) NOT NULL,
 `type` int(10) NOT NULL,
 `account_id` bigint(20) NOT NULL,
 `group_id` bigint(20) NOT NULL,
 `tags` varchar(255),
 `description` varchar(500),
 `created_at` datetime,
 `updated_at` datetime, 
PRIMARY KEY (id)
);

-- Data for table cmdb_sql
INSERT INTO `cmdb_sql` (`id`,`name`,`type`,`account_id`,`group_id`,`tags`,`description`,`created_at`,`updated_at`) VALUES (1,'saas3-mysql',1,1,1,'prod','1111','2025-07-29 21:23:17','2025-07-29 21:23:17');
INSERT INTO `cmdb_sql` (`id`,`name`,`type`,`account_id`,`group_id`,`tags`,`description`,`created_at`,`updated_at`) VALUES (2,'saas3-redis-1',3,3,1,'prod','1111','2025-07-29 21:24:57','2025-09-06 15:12:22');
INSERT INTO `cmdb_sql` (`id`,`name`,`type`,`account_id`,`group_id`,`tags`,`description`,`created_at`,`updated_at`) VALUES (4,'saas3-pgsql',2,1,1,'prod','1111','2025-07-29 21:36:11','2025-09-06 15:12:44');

-- Table structure for cmdb_sql_log
CREATE TABLE `cmdb_sql_log` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `instance_id` varchar(64) NOT NULL,
 `database` varchar(128) NOT NULL,
 `operation_type` varchar(32) NOT NULL,
 `sql_content` text NOT NULL,
 `exec_user` varchar(64) NOT NULL,
 `ip` varchar(64) NOT NULL,
 `scanned_rows` bigint(19) DEFAULT 0,
 `affected_rows` bigint(19) DEFAULT 0,
 `execution_time` bigint(19) DEFAULT 0,
 `returned_rows` bigint(19) DEFAULT 0,
 `result` varchar(32) NOT NULL,
 `query_time` datetime NOT NULL, 
PRIMARY KEY (id)
);
ALTER TABLE `cmdb_sql_log` ADD  INDEX `idx_cmdb_sql_log_query_time`(`query_time`) USING BTREE;

-- Data for table cmdb_sql_log
INSERT INTO `cmdb_sql_log` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`ip`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`) VALUES (7,'8.130.14.34','gin-api','SELECT','select * from  cmdb_host;','admin','127.0.0.1',0,0,59,3,'SUCCESS','2025-08-25 10:17:24');

-- Table structure for cmdb_sql_records
CREATE TABLE `cmdb_sql_records` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `instance_id` varchar(64) NOT NULL,
 `database` varchar(128) NOT NULL,
 `operation_type` varchar(32) NOT NULL,
 `sql_content` text NOT NULL,
 `exec_user` varchar(64) NOT NULL,
 `scanned_rows` bigint(19) DEFAULT 0,
 `affected_rows` bigint(19) DEFAULT 0,
 `execution_time` bigint(19) DEFAULT 0,
 `returned_rows` bigint(19) DEFAULT 0,
 `result` varchar(32) NOT NULL,
 `query_time` datetime NOT NULL,
 `name` varchar(64) NOT NULL, 
PRIMARY KEY (id)
);
ALTER TABLE `cmdb_sql_records` ADD  INDEX `idx_cmdb_sql_records_query_time`(`query_time`) USING BTREE;

-- Data for table cmdb_sql_records
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (1,'8.130.14.34:3306','gin-api','SELECT','select * from cmdb_group;','anonymous',14,0,403,14,'SUCCESS','2025-07-29 11:14:13','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (2,'8.130.14.30','saas3-mysql','SELECT','select * from cmdb_group;','',0,0,50,10,'SUCCESS','2025-07-30 13:29:24','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (3,'8.130.14.34','gin-api','SELECT','select * from cmdb_group;','',0,0,54,14,'SUCCESS','2025-07-30 13:58:13','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (4,'8.130.14.34','saas3-mysql','INSERT','UPDATE `cmdb_group`SET `name` = ''sql测试组0000'' WHERE `parent_id` = 17AND `name` = ''sql测试组'';','',0,1,80,0,'SUCCESS','2025-07-30 14:00:44','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (5,'8.130.14.34','saas3-mysql','INSERT','UPDATE `cmdb_group` SET `name` = ''test123'' WHERE `id` = 22;','',0,1,80,0,'SUCCESS','2025-07-30 14:04:30','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (6,'8.130.14.34','saas3-mysql','INSERT','UPDATE `cmdb_group` SET `name` = ''test123111'' WHERE `id` = 22;','',0,1,80,0,'SUCCESS','2025-07-30 14:06:34','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (7,'8.130.14.34','gin-api','SELECT','UPDATE `cmdb_group` SET `name` = ''test123111'' WHERE `id` = 22;','',0,0,122,0,'SUCCESS','2025-07-30 14:07:44','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (8,'8.130.14.34','gin-api','SELECT','UPDATE `cmdb_group` SET `name` = ''test001'' WHERE `id` = 22;','',0,0,55,0,'SUCCESS','2025-07-30 14:11:42','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (9,'8.130.14.34','gin-api','SELECT','select * from cmdb_group;','',0,0,67,15,'SUCCESS','2025-07-30 15:24:57','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (10,'8.130.14.34','gin-api','SELECT','select * from cmdb_host;','',0,0,114,75,'SUCCESS','2025-07-30 15:25:19','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (11,'8.130.14.34','gin-api','SELECT','select * from cmdb_group;','',0,0,56,15,'SUCCESS','2025-07-30 15:29:16','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (12,'8.130.14.34','saas3-mysql','EXECUTE','create databases  db;','',0,1,100,0,'SUCCESS','2025-07-30 15:34:31','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (13,'8.130.14.34','gin-api','SELECT','select id,name from cmdb_group;','',0,0,51,15,'SUCCESS','2025-07-30 17:18:52','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (14,'8.130.14.34','gin-api','SELECT','select id,name from cmdb_group;','',0,0,53,15,'SUCCESS','2025-07-30 17:29:36','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (15,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','',0,0,56,15,'SUCCESS','2025-07-30 21:06:07','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (16,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','',0,0,64,15,'SUCCESS','2025-07-30 21:12:04','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (17,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','admin',0,0,51,15,'SUCCESS','2025-07-30 21:23:55','');
INSERT INTO `cmdb_sql_records` (`id`,`instance_id`,`database`,`operation_type`,`sql_content`,`exec_user`,`scanned_rows`,`affected_rows`,`execution_time`,`returned_rows`,`result`,`query_time`,`name`) VALUES (18,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','zhangsan',0,0,54,15,'SUCCESS','2025-07-30 21:26:51','');

-- Table structure for config_account
CREATE TABLE `config_account` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `alias` varchar(128) NOT NULL,
 `host` varchar(128) NOT NULL,
 `port` bigint(19) NOT NULL,
 `name` varchar(128) NOT NULL,
 `password` text NOT NULL,
 `type` bigint(19) NOT NULL,
 `remark` text,
 `created_at` datetime,
 `updated_at` datetime, 
PRIMARY KEY (id)
);

-- Data for table config_account
INSERT INTO `config_account` (`id`,`alias`,`host`,`port`,`name`,`password`,`type`,`remark`,`created_at`,`updated_at`) VALUES (13,'腾讯云测试jenkins','test-ops-jenkins-tc2.dding.net:8080',0,'zhangfan','FRXE96C2SKn2K8JE3abAY/CNEEDloyUYlKgpxFDlDNw=',4,'123','2025-11-07 20:16:22','2025-11-07 20:27:57');
INSERT INTO `config_account` (`id`,`alias`,`host`,`port`,`name`,`password`,`type`,`remark`,`created_at`,`updated_at`) VALUES (14,'测试环境jenkins','test-ops-jenkins-tc2.dding.net:8080',0,'zhangfan','s+2RZaFWCHNnA6EIAXJqtfek1DP7iE7HtbeHhxu6en4=',4,'test','2025-11-14 19:25:19','2025-11-14 19:25:19');

-- Table structure for config_ecsauth
CREATE TABLE `config_ecsauth` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `name` longtext NOT NULL COMMENT '''凭证名称''',
 `type` bigint(19) NOT NULL COMMENT '''认证类型:1->密码,2->私钥,3->公钥(免认证)''',
 `username` longtext COMMENT '''用户名''',
 `password` longtext COMMENT '''密码(type=1时使用)''',
 `public_key` text COMMENT '''私钥内容(type=2时使用，字段名历史原因)''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `remark` longtext COMMENT '''备注''',
 `port` bigint(19) DEFAULT 22 COMMENT '''端口号''', 
PRIMARY KEY (id)
);

-- Data for table config_ecsauth
INSERT INTO `config_ecsauth` (`id`,`name`,`type`,`username`,`password`,`public_key`,`create_time`,`remark`,`port`) VALUES (13,'ubuntu-1',1,'root','123456','-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn
NhAAAAAwEAAQAAAYEAqQG9d/UFU0eaa9rMd5l31mX1OLV8PKfuwiAy1vu4MDsIyWcQdMjx
n2Sy/Oj91RKUgnb/LNLQNpM7ARYlUJdyU5X7/PZVUF7bKV7RSdN5bcBTW36uOaX5K4mWNW
VWm/V0zFxllMb5tK5t3OlSyekEcFuKTSvsmZFfSaNBLkRY3qAC64RVRcxDiw0yHR68kipv
3QvJfj2uFTt9X+lKXYsHebTvxgIVgM+KtkEi5R7MZNFiqeTTWGOidwyooVD+wJPkxFvi9u
mIPkm7/ExR4A/WwQpIbyiA4r8k2hYFGRfbn2J2I7xq9t59xl57TtYRQ0VWbGSj+xxHFIOX
ovkk79fea1zgCUxbiTLESDm9rSIbQHWPC3uGhgoiwGn6Xb4VI+Xb6hieOdcbtftZq176Ko
YivJrw8eNGuMZf7/VbsVtQZhNovHTeXkVDBLzb45qU3oyhM/Bi0OkizR0NEDQGNZvCs3P3
kY6hLwSyC9FbOSPaG63PZpnUP/F2lNjbDx+w5qHBAAAFiKrffV+q331fAAAAB3NzaC1yc2
EAAAGBAKkBvXf1BVNHmmvazHeZd9Zl9Ti1fDyn7sIgMtb7uDA7CMlnEHTI8Z9ksvzo/dUS
lIJ2/yzS0DaTOwEWJVCXclOV+/z2VVBe2yle0UnTeW3AU1t+rjml+SuJljVlVpv1dMxcZZ
TG+bSubdzpUsnpBHBbik0r7JmRX0mjQS5EWN6gAuuEVUXMQ4sNMh0evJIqb90LyX49rhU7
fV/pSl2LB3m078YCFYDPirZBIuUezGTRYqnk01hjoncMqKFQ/sCT5MRb4vbpiD5Ju/xMUe
AP1sEKSG8ogOK/JNoWBRkX259idiO8avbefcZee07WEUNFVmxko/scRxSDl6L5JO/X3mtc
4AlMW4kyxEg5va0iG0B1jwt7hoYKIsBp+l2+FSPl2+oYnjnXG7X7Wate+iqGIrya8PHjRr
jGX+/1W7FbUGYTaLx03l5FQwS82+OalN6MoTPwYtDpIs0dDRA0BjWbwrNz95GOoS8EsgvR
Wzkj2hutz2aZ1D/xdpTY2w8fsOahwQAAAAMBAAEAAAGAPMbzbAAhOiG11rOIdDosbl2cIh
ih5O4/XnBV2BoN0spoDoWR1W2t9rQv6eOf5NAZuBEzEtt9JSMtRObB5ImSA50cUYwEgjKa
ffYrLXFvaZiQYYbdAL4/LPj9m5Yl31AWAKf5n9cuVVulBpyhuRqjDgTdZ4M0QsWUjta2yf
yqOOUyJ6AfSkp6v7avHum+xpGcSNHEVDG6OMh7/dHXfGxS8+GtmHeHZBs+tuwVsG+y8U6b
PJ2do1uxQT68JccK1hZYyOTCPyK/2jQxSXFJEG6FY1jndxkaSvsKFe39jzdxRic1D6R1bN
beZASTGrBzr734BaWtjQkzpO9kJ8kTLjbP8MVQYY4wCTfcnf1eAhFMlJJY4YGohrZ6FhAX
uOMrP+v6n0m8dXsrIMAgRWiYh1K8ZAK7tIIlM5SPyh7E8YbCGHuQNF2XARebLk0o13u8ff
5U58bVoEJBnXrOuiJdbvR/7QRcWcZUTV9kdRisF3ooo+TKWd7VOK2BAXHsam56PXwBAAAA
wQC43Px4FX19yPFXAu7zhw5IeH1XdfAoPzToYFe1VtvHYKAG3ivcVmNqo2EzUG1TJIIY7f
clepUAZqUj4c3R/GafwPY8QW3QbpMimpx/hVuyA3yDC0qjeaQ335yvdTtVu/Jv4yXAbjRi
2kgGDfoPQoHnNlNcfhcZI5c1e8yYqNy+cnPi90dnzvLACsIIrcw82IEN6cjf0XomplTdvM
R63qxvOAWOQZxcoUbZuTUPIqb8P1M/vzOS0mfvMqj/dqW8mv4AAADBANNix6vGRxWQoPcK
uY29BV2brwV3SKlFlTUtj8+YXIWdGVA4k94QJGeSr+mikiO8PySjeBZB/Ki8qdh+vPMKZ4
CZM9Uk38qkOdoAZTlJ2zCTIQoWGt1j4eSiAy38LBEiruAJPtXXOWTKSM0g/WaaX/8EhyTp
EuODV5UVJlsg6mXeH4OgcCo9B2Fmn5l8AKiexCucITGuRdzVhN3KNph/EIzd0Kn7zL2dMW
T0rtYL3Ja1y70UAqwDdVUyJjQFtrbbiQAAAMEAzK02BABVLbgcnEGRWeNaJQ7lA36hU3x1
6XgmESJWSPN4Yp5gNQBHsuFxB6ksZQkX76penNQG1cJPu7xURP4lrCQdyoCfdzcO+PpNKX
fPCBNLNGCefWKjcBdp961bkhNazH1B2p5g4lytZbS2NJduB3W+L7dcF18mnBCIlNdLxBm7
Z/PDjQfAhurtKZHvhWjg9qX32xqETldjePZUNHdt1TQ0Ih0Ei0J23trsKxBXkP0SzAIXiR
/z4Lox2qK0BG55AAAADnJvb3RAdWJ1bnR1LTAxAQIDBA==
-----END OPENSSH PRIVATE KEY-----','2025-07-28 16:11:42','这是自建虚拟机',22);
INSERT INTO `config_ecsauth` (`id`,`name`,`type`,`username`,`password`,`public_key`,`create_time`,`remark`,`port`) VALUES (20,'免密码认证',3,'root','','','2025-10-14 21:06:27','',22);
INSERT INTO `config_ecsauth` (`id`,`name`,`type`,`username`,`password`,`public_key`,`create_time`,`remark`,`port`) VALUES (21,'devops',1,'root','devops@2025','','2025-10-18 16:13:32','',22);
INSERT INTO `config_ecsauth` (`id`,`name`,`type`,`username`,`password`,`public_key`,`create_time`,`remark`,`port`) VALUES (22,'华为云ops',1,'root','zhangfan@2025','','2025-11-11 17:24:28','123',22);
INSERT INTO `config_ecsauth` (`id`,`name`,`type`,`username`,`password`,`public_key`,`create_time`,`remark`,`port`) VALUES (23,'百度云',1,'root','devops@2025','','2025-11-25 16:39:36','devops@2025',22);

-- Table structure for config_keymanage
CREATE TABLE `config_keymanage` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `key_type` bigint(19) NOT NULL,
 `key_id` text NOT NULL,
 `key_secret` text NOT NULL,
 `remark` text,
 `created_at` datetime,
 `updated_at` datetime, 
PRIMARY KEY (id)
);

-- Data for table config_keymanage
-- No data

-- Table structure for config_sync_schedule
CREATE TABLE `config_sync_schedule` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(100) NOT NULL,
 `cron_expr` varchar(100) NOT NULL,
 `key_types` text NOT NULL,
 `status` bigint(19) NOT NULL DEFAULT 1,
 `last_run_time` datetime,
 `next_run_time` datetime,
 `remark` text,
 `created_at` datetime,
 `updated_at` datetime,
 `sync_log` text, 
PRIMARY KEY (id)
);

-- Data for table config_sync_schedule
INSERT INTO `config_sync_schedule` (`id`,`name`,`cron_expr`,`key_types`,`status`,`last_run_time`,`next_run_time`,`remark`,`created_at`,`updated_at`,`sync_log`) VALUES (3,'阿里云定时同步','*/3 * * * *','[1]',0,'2025-10-01 13:33:00','2025-10-01 13:36:00','','2025-09-29 18:41:10','2025-10-01 13:34:41','[2025-10-01 13:33:00] 开始同步
- 阿里云: 同步成功
[2025-10-01 13:33:00] 同步完成，耗时: 827.31718ms
');

-- Table structure for db
CREATE TABLE `db` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `code` varchar(36) NOT NULL COMMENT '''数据库编码''',
 `name` varchar(100) NOT NULL COMMENT '''数据库名称''',
 `database` varchar(500) NOT NULL COMMENT '''数据库名(多个用空格分隔)''',
 `remark` varchar(500) COMMENT '''备注''',
 `instance_id` bigint(20) NOT NULL COMMENT '''实例ID''',
 `instance_code` varchar(36) COMMENT '''实例编码''',
 `status` bigint(19) DEFAULT 1 COMMENT '''状态:1->启用,2->禁用''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''',
 `modifier` varchar(64) COMMENT '''修改人''',
 `modifier_id` bigint(20) COMMENT '''修改人ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `db` ADD unique INDEX `idx_db_code`(`code`) USING BTREE;
ALTER TABLE `db` ADD  INDEX `idx_db_instance_id`(`instance_id`) USING BTREE;

-- Data for table db
-- No data

-- Table structure for db_es_instance
CREATE TABLE `db_es_instance` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `code` varchar(36) NOT NULL COMMENT '''实例编码''',
 `name` varchar(100) NOT NULL COMMENT '''实例名称''',
 `protocol` varchar(10) DEFAULT http COMMENT '''协议:http|https''',
 `host` varchar(255) NOT NULL COMMENT '''主机''',
 `port` bigint(19) DEFAULT 9200 COMMENT '''端口''',
 `username` varchar(100) COMMENT '''用户名(可选)''',
 `password` varchar(500) COMMENT '''密码(加密)''',
 `remark` varchar(500) COMMENT '''备注''',
 `ssh_tunnel_machine_id` bigint(20) DEFAULT 0 COMMENT '''SSH隧道机器ID''',
 `status` bigint(19) DEFAULT 1 COMMENT '''状态:1->启用,2->禁用''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''',
 `modifier` varchar(64) COMMENT '''修改人''',
 `modifier_id` bigint(20) COMMENT '''修改人ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `db_es_instance` ADD unique INDEX `idx_db_es_instance_code`(`code`) USING BTREE;

-- Data for table db_es_instance
INSERT INTO `db_es_instance` (`id`,`code`,`name`,`protocol`,`host`,`port`,`username`,`password`,`remark`,`ssh_tunnel_machine_id`,`status`,`create_time`,`update_time`,`creator`,`creator_id`,`modifier`,`modifier_id`) VALUES (1,'d6340479-aa87-4dc7-914b-e92cbc29af18','阿里云-测试环境es','http','192.168.2.123',9200,'','123','123',0,1,'2025-12-01 16:59:59','2025-12-02 17:09:47','',0,'',0);

-- Table structure for db_instance
CREATE TABLE `db_instance` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `code` varchar(36) NOT NULL COMMENT '''实例编码''',
 `name` varchar(100) NOT NULL COMMENT '''实例名称''',
 `type` varchar(20) NOT NULL COMMENT '''数据库类型:mysql,postgres,oracle等''',
 `host` varchar(100) NOT NULL COMMENT '''主机地址''',
 `port` bigint(19) NOT NULL COMMENT '''端口''',
 `network` varchar(20) DEFAULT tcp COMMENT '''网络类型:tcp,unix''',
 `params` varchar(500) COMMENT '''连接参数''',
 `username` varchar(100) NOT NULL COMMENT '''用户名''',
 `password` varchar(500) NOT NULL COMMENT '''密码(加密)''',
 `remark` varchar(500) COMMENT '''备注''',
 `ssh_tunnel_machine_id` bigint(20) DEFAULT 0 COMMENT '''SSH隧道机器ID''',
 `status` bigint(19) DEFAULT 1 COMMENT '''状态:1->启用,2->禁用''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''',
 `modifier` varchar(64) COMMENT '''修改人''',
 `modifier_id` bigint(20) COMMENT '''修改人ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `db_instance` ADD unique INDEX `idx_db_instance_code`(`code`) USING BTREE;

-- Data for table db_instance
INSERT INTO `db_instance` (`id`,`code`,`name`,`type`,`host`,`port`,`network`,`params`,`username`,`password`,`remark`,`ssh_tunnel_machine_id`,`status`,`create_time`,`update_time`,`creator`,`creator_id`,`modifier`,`modifier_id`) VALUES (1,'80774a9e-ab4b-4b93-89b7-fce449ab6a2c','阿里云-mysql','mysql','8.130.14.34',3306,'tcp','','root','devops@2025','123',0,1,'2025-11-29 17:10:50','2025-11-29 17:10:50','',0,'',0);
INSERT INTO `db_instance` (`id`,`code`,`name`,`type`,`host`,`port`,`network`,`params`,`username`,`password`,`remark`,`ssh_tunnel_machine_id`,`status`,`create_time`,`update_time`,`creator`,`creator_id`,`modifier`,`modifier_id`) VALUES (2,'e2a485e1-4a82-4875-848e-9e975c0f01a5','华为云pgsql','postgres','139.9.205.38',5432,'tcp','sslmode=disable','root','devops@123','123',0,1,'2025-11-30 00:14:57','2025-11-30 00:19:57','',0,'',0);
INSERT INTO `db_instance` (`id`,`code`,`name`,`type`,`host`,`port`,`network`,`params`,`username`,`password`,`remark`,`ssh_tunnel_machine_id`,`status`,`create_time`,`update_time`,`creator`,`creator_id`,`modifier`,`modifier_id`) VALUES (5,'096d5f73-3812-48e2-9fbd-f4214d4999ba','华为云mysql','mysql','139.9.205.38',3306,'tcp','','root','devops@2025','ops',0,1,'2025-12-06 22:16:56','2025-12-06 22:16:56','',0,'',0);

-- Table structure for db_mongo_instance
CREATE TABLE `db_mongo_instance` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `code` varchar(36) NOT NULL COMMENT '''实例编码''',
 `name` varchar(100) NOT NULL COMMENT '''实例名称''',
 `uri` varchar(500) NOT NULL COMMENT '''连接URI''',
 `ssh_tunnel_machine_id` bigint(20) DEFAULT 0 COMMENT '''SSH隧道机器ID''',
 `remark` varchar(500) COMMENT '''备注''',
 `status` bigint(19) DEFAULT 1 COMMENT '''状态:1->启用,2->禁用''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''',
 `modifier` varchar(64) COMMENT '''修改人''',
 `modifier_id` bigint(20) COMMENT '''修改人ID''',
 `type` varchar(20) DEFAULT mongodb COMMENT '''MongoDB类型:mongodb,mongodb-atlas等''', 
PRIMARY KEY (id)
);
ALTER TABLE `db_mongo_instance` ADD unique INDEX `idx_db_mongo_instance_code`(`code`) USING BTREE;

-- Data for table db_mongo_instance
INSERT INTO `db_mongo_instance` (`id`,`code`,`name`,`uri`,`ssh_tunnel_machine_id`,`remark`,`status`,`create_time`,`update_time`,`creator`,`creator_id`,`modifier`,`modifier_id`,`type`) VALUES (2,'mongo_1764642528321099000','mongo-test','mongodb://admin:password@192.168.2.123:27018',0,'123',1,'2025-12-02 10:28:48','2025-12-02 10:28:48','',0,'',0,'mongodb');

-- Table structure for db_redis_instance
CREATE TABLE `db_redis_instance` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `code` varchar(36) NOT NULL COMMENT '''实例编码''',
 `name` varchar(100) NOT NULL COMMENT '''实例名称''',
 `mode` varchar(20) NOT NULL COMMENT '''模式:standalone,cluster,sentinel''',
 `host` varchar(300) NOT NULL COMMENT '''主机: standalone为host:port, cluster为逗号分隔, sentinel为master=hosts''',
 `port` bigint(19) DEFAULT 0 COMMENT '''端口(standalone可用)''',
 `db` bigint(19) DEFAULT 0 COMMENT '''默认库号''',
 `username` varchar(100) COMMENT '''用户名(可选)''',
 `password` varchar(500) NOT NULL COMMENT '''密码(加密)''',
 `redis_node_password` varchar(500) COMMENT '''节点密码(仅sentinel)''',
 `remark` varchar(500) COMMENT '''备注''',
 `ssh_tunnel_machine_id` bigint(20) DEFAULT 0 COMMENT '''SSH隧道机器ID''',
 `status` bigint(19) DEFAULT 1 COMMENT '''状态:1->启用,2->禁用''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''',
 `modifier` varchar(64) COMMENT '''修改人''',
 `modifier_id` bigint(20) COMMENT '''修改人ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `db_redis_instance` ADD unique INDEX `idx_db_redis_instance_code`(`code`) USING BTREE;

-- Data for table db_redis_instance
INSERT INTO `db_redis_instance` (`id`,`code`,`name`,`mode`,`host`,`port`,`db`,`username`,`password`,`redis_node_password`,`remark`,`ssh_tunnel_machine_id`,`status`,`create_time`,`update_time`,`creator`,`creator_id`,`modifier`,`modifier_id`) VALUES (4,'2d8ac975-0a4b-4e10-8b92-d6ad781938f2','my-redis','standalone','8.130.14.34',6379,0,'','zhangfan@123','','123',0,1,'2025-11-30 13:39:44','2025-12-02 17:09:53','',0,'',0);

-- Table structure for db_sql
CREATE TABLE `db_sql` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `db_id` bigint(20) NOT NULL COMMENT '''数据库ID''',
 `db` varchar(100) NOT NULL COMMENT '''数据库名''',
 `name` varchar(100) NOT NULL COMMENT '''SQL名称''',
 `type` bigint(19) DEFAULT 1 COMMENT '''类型:1->查询,2->更新''',
 `sql` text NOT NULL COMMENT '''SQL语句''',
 `remark` varchar(500) COMMENT '''备注''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''',
 `modifier` varchar(64) COMMENT '''修改人''',
 `modifier_id` bigint(20) COMMENT '''修改人ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `db_sql` ADD  INDEX `idx_db_sql_db_id`(`db_id`) USING BTREE;

-- Data for table db_sql
-- No data

-- Table structure for db_sql_exec
CREATE TABLE `db_sql_exec` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `db_id` bigint(20) NOT NULL COMMENT '''数据库ID''',
 `db_name` varchar(100) NOT NULL COMMENT '''数据库名''',
 `table_name` varchar(100) COMMENT '''表名''',
 `type` tinyint(3) NOT NULL COMMENT '''类型:1->查询,2->插入,3->更新,4->删除,5->DDL''',
 `sql` text NOT NULL COMMENT '''SQL语句''',
 `old_value` longtext COMMENT '''旧值(用于回滚)''',
 `remark` varchar(500) COMMENT '''备注''',
 `status` tinyint(3) NOT NULL COMMENT '''状态:1->成功,2->失败''',
 `res` text COMMENT '''执行结果''',
 `exec_time` bigint(19) COMMENT '''执行时长(ms)''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `db_sql_exec` ADD  INDEX `idx_db_sql_exec_create_time`(`create_time`) USING BTREE;
ALTER TABLE `db_sql_exec` ADD  INDEX `idx_db_sql_exec_db_id`(`db_id`) USING BTREE;

-- Data for table db_sql_exec
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (1,1,'gin-api','APP_APPLICATION',0,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',2,'数据库不存在',110,'2025-11-29 18:05:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (2,1,'gin-api','APP_APPLICATION',0,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',2,'数据库不存在',112,'2025-11-29 18:05:05','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (3,1,'gin-api','CMDB_GROUP',0,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',2,'数据库不存在',119,'2025-11-29 18:05:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (4,1,'gin-api','CMDB_GROUP',0,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',2,'数据库不存在',116,'2025-11-29 18:05:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (5,1,'gin-api','CMDB_SQL',0,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',2,'数据库不存在',122,'2025-11-29 18:05:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (6,1,'gin-api','CMDB_SQL',0,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',2,'数据库不存在',136,'2025-11-29 18:05:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (7,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',119,'2025-11-29 18:11:06','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (8,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',493,'2025-11-29 18:11:06','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (9,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',127,'2025-11-29 18:11:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (10,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',110,'2025-11-29 18:11:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (11,1,'gin-api','APP_SERVICE_RELEASE_ITEM',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',548,'2025-11-29 18:28:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (12,1,'gin-api','APP_SERVICE_RELEASE_ITEM',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',1083,'2025-11-29 18:28:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (13,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 18:28:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (14,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',281,'2025-11-29 18:28:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (15,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',308,'2025-11-29 18:31:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (16,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',346,'2025-11-29 18:31:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (17,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',298,'2025-11-29 18:35:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (18,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-29 18:35:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (19,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',266,'2025-11-29 18:35:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (20,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',259,'2025-11-29 18:35:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (21,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 18:37:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (22,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',320,'2025-11-29 18:37:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (23,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-11-29 18:37:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (24,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',299,'2025-11-29 18:37:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (25,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',255,'2025-11-29 18:37:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (26,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',282,'2025-11-29 18:37:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (27,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-11-29 18:37:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (28,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',260,'2025-11-29 18:37:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (29,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-29 18:38:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (30,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 18:38:47','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (31,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',265,'2025-11-29 18:38:47','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (32,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',266,'2025-11-29 18:38:48','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (33,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-11-29 18:39:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (34,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-29 18:39:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (35,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',287,'2025-11-29 18:39:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (36,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',263,'2025-11-29 18:39:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (37,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',403,'2025-11-29 18:48:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (38,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',335,'2025-11-29 18:48:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (39,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',10295,'2025-11-29 18:48:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (40,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',2650,'2025-11-29 18:48:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (41,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',6494,'2025-11-29 18:48:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (42,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',293,'2025-11-29 18:48:16','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (43,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',3692,'2025-11-29 18:48:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (44,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',281,'2025-11-29 18:48:21','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (45,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',340,'2025-11-29 18:50:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (46,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',284,'2025-11-29 18:50:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (47,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-29 18:54:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (48,1,'devops','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',294,'2025-11-29 18:54:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (49,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',328,'2025-11-29 18:54:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (50,1,'devops','CMDB_GROUP',1,'SELECT COUNT(*) as total FROM cmdb_group','','',1,'执行成功',313,'2025-11-29 18:54:08','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (51,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 18:55:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (52,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',280,'2025-11-29 18:55:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (53,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-11-29 18:55:29','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (54,1,'devops','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',277,'2025-11-29 18:55:30','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (55,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 18:56:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (56,1,'devops','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',275,'2025-11-29 18:56:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (57,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',309,'2025-11-29 18:56:37','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (58,1,'devops','CMDB_GROUP',1,'SELECT COUNT(*) as total FROM cmdb_group','','',1,'执行成功',293,'2025-11-29 18:56:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (59,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',276,'2025-11-29 18:57:16','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (60,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',284,'2025-11-29 18:57:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (61,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',402,'2025-11-29 18:57:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (62,1,'devops','CMDB_SQL_LOG',1,'SELECT COUNT(*) as total FROM cmdb_sql_log','','',1,'执行成功',274,'2025-11-29 18:57:50','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (63,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',315,'2025-11-29 19:00:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (64,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',330,'2025-11-29 19:00:50','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (65,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',484,'2025-11-29 19:00:57','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (66,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',656,'2025-11-29 19:00:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (67,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',432,'2025-11-29 19:01:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (68,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',285,'2025-11-29 19:01:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (69,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',300,'2025-11-29 19:01:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (70,1,'devops','CONFIG_ACCOUNT',1,'SELECT COUNT(*) as total FROM config_account','','',1,'执行成功',282,'2025-11-29 19:01:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (71,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',442,'2025-11-29 19:02:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (72,1,'devops','CMDB_SQL_LOG',1,'SELECT COUNT(*) as total FROM cmdb_sql_log','','',1,'执行成功',388,'2025-11-29 19:02:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (73,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',321,'2025-11-29 19:04:24','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (74,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',291,'2025-11-29 19:04:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (75,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-11-29 19:05:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (76,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',309,'2025-11-29 19:05:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (77,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',361,'2025-11-29 19:06:09','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (78,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',291,'2025-11-29 19:06:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (79,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',1109,'2025-11-29 19:06:32','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (80,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',1586,'2025-11-29 19:06:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (81,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-29 19:06:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (82,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',552,'2025-11-29 19:06:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (83,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',307,'2025-11-29 19:06:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (84,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',254,'2025-11-29 19:06:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (85,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',553,'2025-11-29 19:06:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (86,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',302,'2025-11-29 19:06:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (87,1,'gin-api','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 19:06:36','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (88,1,'gin-api','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',276,'2025-11-29 19:06:37','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (89,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-29 19:08:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (90,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',331,'2025-11-29 19:08:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (91,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',324,'2025-11-29 19:08:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (92,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 19:08:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (93,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-29 19:09:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (94,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',259,'2025-11-29 19:09:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (95,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',360,'2025-11-29 19:08:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (96,1,'devops','CONFIG_ACCOUNT',1,'SELECT COUNT(*) as total FROM config_account','','',1,'执行成功',292,'2025-11-29 19:09:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (97,1,'devops','CMDB_SQL_RECORDS',1,'SELECT COUNT(*) as total FROM cmdb_sql_records','','',1,'执行成功',303,'2025-11-29 19:09:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (98,1,'devops','CONFIG_ACCOUNT',1,'SELECT COUNT(*) as total FROM config_account','','',1,'执行成功',326,'2025-11-29 19:09:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (99,1,'devops','CMDB_SQL_RECORDS',1,'SELECT COUNT(*) as total FROM cmdb_sql_records','','',1,'执行成功',295,'2025-11-29 19:09:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (100,1,'devops','CMDB_SQL_RECORDS',1,'SELECT COUNT(*) as total FROM cmdb_sql_records','','',1,'执行成功',334,'2025-11-29 19:09:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (101,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-11-29 19:09:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (102,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',278,'2025-11-29 19:09:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (103,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',1029,'2025-11-29 19:09:29','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (104,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',6538,'2025-11-29 19:09:36','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (105,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',264,'2025-11-29 19:23:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (106,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',296,'2025-11-29 19:23:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (107,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 19:23:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (108,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',279,'2025-11-29 19:24:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (109,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 100;','','',1,'执行成功',280,'2025-11-29 19:25:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (110,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-11-29 19:49:42','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (111,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',287,'2025-11-29 19:49:43','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (112,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',239,'2025-11-29 19:57:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (113,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 19:57:40','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (114,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 19:57:40','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (115,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',212,'2025-11-29 19:57:40','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (116,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',210,'2025-11-29 19:57:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (117,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-11-29 19:58:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (118,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',307,'2025-11-29 19:58:47','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (119,1,'devops','CMDB_HOST;',1,'SELECT * FROM cmdb_host;','','',1,'执行成功',212,'2025-11-29 20:00:09','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (120,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',263,'2025-11-29 20:03:26','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (121,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 20:03:44','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (122,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 100;','','',1,'执行成功',220,'2025-11-29 20:04:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (123,1,'devops','TASK_WORK',1,'SELECT * FROM task_work LIMIT 20 OFFSET 0','','',1,'执行成功',234,'2025-11-29 20:04:29','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (124,1,'devops','K8S_CLUSTER',1,'SELECT * FROM k8s_cluster LIMIT 20 OFFSET 0','','',1,'执行成功',247,'2025-11-29 20:05:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (125,1,'devops','CONFIG_ECSAUTH',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',515,'2025-11-29 20:07:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (126,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',219,'2025-11-29 20:08:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (127,1,'devops','CONFIG_ECSAUTH',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',386,'2025-11-29 20:08:56','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (128,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',212,'2025-11-29 20:09:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (129,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',241,'2025-11-29 20:09:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (130,1,'devops','SYS_MENU',1,'SELECT * FROM sys_menu LIMIT 20 OFFSET 0','','',1,'执行成功',211,'2025-11-29 20:09:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (131,1,'devops','CONFIG_KEYMANAGE',1,'SELECT * FROM config_keymanage LIMIT 20 OFFSET 0','','',1,'执行成功',218,'2025-11-29 20:10:16','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (132,1,'gin-api','APP_SERVICE_RELEASE',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',229,'2025-11-29 20:10:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (133,1,'gin-api','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',321,'2025-11-29 20:11:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (134,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 20:14:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (135,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 20:17:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (136,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',211,'2025-11-29 20:17:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (137,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',240,'2025-11-29 20:27:29','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (138,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-11-29 20:28:55','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (139,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',579,'2025-11-29 20:28:55','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (140,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-11-29 20:30:47','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (141,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',199,'2025-11-29 20:30:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (142,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-11-29 20:34:21','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (143,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',4246,'2025-11-29 20:45:39','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (144,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',582,'2025-11-29 20:49:42','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (145,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',212,'2025-11-29 20:59:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (146,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-11-29 21:02:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (147,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',242,'2025-11-29 21:05:56','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (148,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',269,'2025-11-29 21:06:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (149,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',327,'2025-11-29 21:08:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (150,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',646,'2025-11-29 21:10:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (151,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 100;','','',1,'执行成功',293,'2025-11-29 21:18:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (152,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',320,'2025-11-29 21:21:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (153,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-11-29 21:21:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (154,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',326,'2025-11-29 21:21:27','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (155,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:22:09','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (156,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:22:19','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (157,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',276,'2025-11-29 21:22:21','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (158,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',264,'2025-11-29 21:34:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (159,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',1520,'2025-11-29 21:35:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (160,1,'sys','sys_config',1,'SELECT * FROM sys_config LIMIT 20 OFFSET 0','','',1,'执行成功',293,'2025-11-29 21:36:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (161,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-11-29 21:37:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (162,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-11-29 21:37:26','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (163,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',298,'2025-11-29 21:37:48','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (164,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 21:38:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (165,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 21:40:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (166,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:42:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (167,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:45:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (168,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',285,'2025-11-29 21:45:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (169,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',270,'2025-11-29 21:45:42','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (170,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',491,'2025-11-29 21:46:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (171,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 21:50:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (172,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',270,'2025-11-29 21:50:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (173,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-11-29 21:51:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (174,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-11-29 21:54:42','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (175,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 21:57:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (176,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',408,'2025-11-29 21:57:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (177,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',279,'2025-11-29 21:57:43','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (178,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',341,'2025-11-29 21:57:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (179,1,'devops','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',309,'2025-11-29 21:58:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (180,1,'devops','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',579,'2025-11-29 21:58:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (181,1,'devops','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 21:58:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (182,1,'devops','monitor_agent',1,'SELECT * FROM monitor_agent LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 21:58:36','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (183,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',330,'2025-11-29 21:59:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (184,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',450,'2025-11-29 21:59:44','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (185,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 22:10:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (186,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-11-29 22:10:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (187,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',323,'2025-11-29 22:11:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (188,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-11-29 22:11:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (189,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',319,'2025-11-29 22:14:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (190,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',508,'2025-11-29 22:14:16','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (191,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-11-29 22:14:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (192,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-11-29 22:14:29','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (193,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',587,'2025-11-29 22:17:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (194,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-11-29 22:21:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (195,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-11-29 22:21:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (196,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',338,'2025-11-29 22:21:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (197,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',935,'2025-11-29 22:21:42','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (198,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-11-29 22:21:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (199,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',290,'2025-11-29 22:25:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (200,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-11-29 22:26:39','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (201,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 22:29:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (202,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-11-29 22:30:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (203,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 22:31:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (204,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',289,'2025-11-29 22:32:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (205,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',272,'2025-11-29 22:32:43','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (206,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',312,'2025-11-29 22:32:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (207,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',268,'2025-11-29 22:38:19','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (208,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',285,'2025-11-29 22:42:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (209,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 22:48:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (210,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',279,'2025-11-29 22:53:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (211,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',271,'2025-11-29 22:57:40','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (212,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-11-29 22:59:48','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (213,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-11-29 23:00:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (214,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',569,'2025-11-29 23:04:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (215,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',542,'2025-11-29 23:10:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (216,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',309,'2025-11-29 23:13:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (217,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',316,'2025-11-29 23:17:36','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (218,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',713,'2025-11-29 23:17:37','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (219,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',311,'2025-11-29 23:17:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (220,1,'gin-api','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',451,'2025-11-29 23:18:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (221,1,'gin-api','config_keymanage',1,'SELECT * FROM config_keymanage LIMIT 20 OFFSET 0','','',1,'执行成功',300,'2025-11-29 23:18:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (222,1,'gin-api','config_sync_schedule',1,'SELECT * FROM config_sync_schedule LIMIT 20 OFFSET 0','','',1,'执行成功',305,'2025-11-29 23:18:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (223,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-29 23:18:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (224,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',334,'2025-11-29 23:18:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (225,1,'gin-api','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-29 23:26:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (226,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',420,'2025-11-29 23:30:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (227,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',331,'2025-11-29 23:35:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (228,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',306,'2025-11-29 23:36:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (229,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',293,'2025-11-29 23:36:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (230,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-11-29 23:39:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (231,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',340,'2025-11-29 23:39:24','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (232,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',289,'2025-11-29 23:39:32','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (233,1,'devops','',5,'RENAME TABLE app_jenkins_env_copy_20251129233940 TO app_jenkins_env_copy_20251129233940_new111','','',1,'执行成功',298,'2025-11-29 23:39:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (234,1,'devops','app_jenkins_env_copy_20251129233940_new111',1,'SELECT * FROM app_jenkins_env_copy_20251129233940_new111 LIMIT 20 OFFSET 0','','',1,'执行成功',563,'2025-11-29 23:40:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (235,1,'devops','',5,'DROP TABLE app_jenkins_env_copy_20251129233940_new111','','',1,'执行成功',291,'2025-11-29 23:40:30','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (236,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',434,'2025-11-29 23:44:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (237,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',466,'2025-11-29 23:45:21','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (238,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-11-29 23:52:27','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (239,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',302,'2025-11-29 23:55:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (240,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',319,'2025-11-29 23:58:06','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (241,2,'postgres','',5,'CREATE DATABASE database_name;','','',1,'执行成功',599,'2025-11-30 00:56:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (242,2,'database_name','',5,'CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL
);','','',1,'执行成功',3259,'2025-11-30 01:10:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (243,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',4965,'2025-11-30 01:10:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (244,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',198,'2025-11-30 01:10:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (245,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',201,'2025-11-30 01:10:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (246,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',208,'2025-11-30 01:10:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (247,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',208,'2025-11-30 01:10:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (248,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',192,'2025-11-30 01:10:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (249,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-11-30 01:10:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (250,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',202,'2025-11-30 01:10:39','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (251,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (''张三'', ''zhangsan@123.com'')','','',1,'执行成功',376,'2025-11-30 01:11:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (252,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-30 01:11:21','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (253,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-30 11:38:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (254,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-30 11:38:44','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (255,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',272,'2025-11-30 11:38:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (256,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',342,'2025-11-30 11:39:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (257,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-30 11:39:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (258,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-11-30 11:40:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (259,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-11-30 14:35:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (260,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',219,'2025-11-30 15:43:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (261,1,'RECOVER_YOUR_DATA','',5,'CREATE DATABASE IF NOT EXISTS test 
  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',762,'2025-11-30 23:49:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (262,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',345,'2025-11-30 23:50:09','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (263,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-30 23:50:24','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (264,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',311,'2025-11-30 23:50:29','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (265,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',285,'2025-11-30 23:51:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (266,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',308,'2025-11-30 23:53:56','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (267,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',291,'2025-11-30 23:54:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (268,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-11-30 23:54:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (269,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-30 23:54:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (270,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',333,'2025-11-30 23:54:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (271,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',319,'2025-11-30 23:54:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (272,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',314,'2025-11-30 23:54:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (273,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',315,'2025-11-30 23:55:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (274,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',290,'2025-11-30 23:55:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (275,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-30 23:55:36','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (276,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-11-30 23:55:40','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (277,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 100;','','',1,'执行成功',288,'2025-11-30 23:56:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (278,1,'devops','cmdb_group;',1,'SELECT * FROM cmdb_group;','','',1,'执行成功',272,'2025-11-30 23:56:57','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (279,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',279,'2025-11-30 23:57:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (280,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-11-30 23:57:32','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (281,1,'devops','cmdb_group',2,'INSERT INTO cmdb_group (parent_id, name, create_time, remark, update_time) VALUES (''0'', ''test'', '''', ''123'', '''')','','',2,'Error 1292 (22007): Incorrect datetime value: '''' for column ''create_time'' at row 1',234,'2025-11-30 23:58:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (282,1,'devops','cmdb_host_copy_20251130235846',1,'SELECT * FROM cmdb_host_copy_20251130235846 LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-11-30 23:58:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (283,1,'devops','',5,'DROP TABLE cmdb_host_copy_20251130235846','','',1,'执行成功',295,'2025-11-30 23:58:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (284,1,'devops','',5,'RENAME TABLE cmdb_host_copy_20251130235914 TO cmdb_host_copy_123','','',1,'执行成功',274,'2025-11-30 23:59:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (285,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',207,'2025-11-30 23:59:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (286,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (''李四'', ''lisi@123.com'')','','',1,'执行成功',285,'2025-12-01 00:00:30','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (287,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',214,'2025-12-01 00:00:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (288,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-12-01 00:05:29','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (289,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-12-01 00:06:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (290,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',1184,'2025-12-01 00:06:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (291,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',611,'2025-12-01 00:15:50','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (292,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',291,'2025-12-01 00:17:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (293,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-12-01 00:19:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (294,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',417,'2025-12-01 00:19:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (295,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-12-01 00:20:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (296,1,'devops','cmdb_host',3,'UPDATE cmdb_host SET host_name = ''虚拟机001111'', group_id = 12, private_ip = ''172.16.226.16'', public_ip = ''120.231.244.158'', ssh_name = ''root'', ssh_key_id = 13, ssh_port = 22, remark = ''123'', vendor = 1, region = '''', instance_id = '''', os = ''CentOSLinux7(Core)'', status = 1, cpu = ''2'', memory = ''4'', disk = ''17'', billing_type = '''', create_time = ''2025-11-23 23:44:35'', expire_time = NULL, update_time = ''2025-11-23 23:45:32'', ssh_ip = ''172.16.226.16'', name = ''jenkins'', ssh_gateway_id = NULL WHERE id = 511;','','',1,'执行成功',254,'2025-12-01 00:20:08','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (297,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-12-01 00:20:08','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (298,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',301,'2025-12-01 00:20:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (299,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 100;','','',1,'执行成功',293,'2025-12-01 00:20:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (300,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',302,'2025-12-01 00:21:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (301,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 100;','','',1,'执行成功',569,'2025-12-01 00:21:57','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (302,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',367,'2025-12-01 00:23:39','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (303,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 100;','','',1,'执行成功',540,'2025-12-01 00:23:47','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (304,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',217,'2025-12-01 00:24:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (305,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',289,'2025-12-01 00:25:05','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (306,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (''test'', ''123@123.com'')','','',1,'执行成功',553,'2025-12-01 00:30:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (307,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-12-01 00:30:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (308,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-01 00:30:24','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (309,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',452,'2025-12-01 00:37:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (310,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',235,'2025-12-01 00:37:27','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (311,2,'database_name','users_copy_20251201003823',1,'SELECT * FROM users_copy_20251201003823 LIMIT 20 OFFSET 0','','',1,'执行成功',197,'2025-12-01 00:38:47','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (312,2,'database_name','users_copy_20251201003823',1,'SELECT * FROM users_copy_20251201003823 LIMIT 100;','','',1,'执行成功',227,'2025-12-01 00:38:56','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (313,1,'RECOVER_YOUR_DATA','',5,'CREATE DATABASE IF NOT EXISTS ops 
  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',729,'2025-12-01 00:43:24','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (314,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',296,'2025-12-01 00:43:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (315,1,'devops','cmdb_host_copy_123',1,'SELECT * FROM cmdb_host_copy_123 LIMIT 20 OFFSET 0','','',1,'执行成功',351,'2025-12-01 00:44:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (316,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',320,'2025-12-01 00:44:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (317,1,'devops','cmdb_host',3,'UPDATE cmdb_host SET host_name = ''虚拟机01'', group_id = 12, private_ip = ''172.16.226.16'', public_ip = ''120.231.244.158'', ssh_name = ''root'', ssh_key_id = 13, ssh_port = 22, remark = ''123'', vendor = 1, region = '''', instance_id = '''', os = ''CentOSLinux7(Core)'', status = 1, cpu = ''2'', memory = ''4'', disk = ''17'', billing_type = '''', create_time = ''2025-11-23 23:44:35'', expire_time = NULL, update_time = ''2025-11-23 23:45:32'', ssh_ip = ''172.16.226.16'', name = ''jenkins'', ssh_gateway_id = NULL WHERE id = 511;','','',1,'执行成功',722,'2025-12-01 00:44:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (318,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',333,'2025-12-01 00:44:26','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (319,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-12-01 00:44:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (320,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',316,'2025-12-01 00:44:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (321,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',475,'2025-12-01 00:46:30','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (322,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',326,'2025-12-01 00:46:48','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (323,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',290,'2025-12-01 00:46:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (324,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',296,'2025-12-01 00:46:57','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (325,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',306,'2025-12-01 00:47:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (326,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',296,'2025-12-01 00:47:09','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (327,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',327,'2025-12-01 00:47:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (328,1,'devops','cmdb_host;',1,'SELECT * FROM  cmdb_host;','','',1,'执行成功',297,'2025-12-01 00:48:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (329,1,'devops','',5,'DROP TABLE cmdb_host_copy_123','','',1,'执行成功',578,'2025-12-01 00:49:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (330,1,'devops','cmdb_host_copy_20251201004911',1,'SELECT * FROM cmdb_host_copy_20251201004911 LIMIT 20 OFFSET 0','','',1,'执行成功',321,'2025-12-01 00:49:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (331,1,'devops','',5,'RENAME TABLE cmdb_host_copy_20251201004911 TO cmdb_host_copy_123','','',1,'执行成功',253,'2025-12-01 00:49:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (332,1,'devops','',5,'DROP TABLE cmdb_host_copy_123','','',1,'执行成功',259,'2025-12-01 00:49:43','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (333,2,'database_name','',5,'CREATE TABLE students (
    id SERIAL PRIMARY KEY,                -- 学生ID，自增主键
    student_id VARCHAR(20) UNIQUE NOT NULL, -- 学号，唯一且非空
    name VARCHAR(50) NOT NULL,            -- 姓名
    gender CHAR(1) CHECK (gender IN (''M'', ''F'')), -- 性别：M 男，F 女
    birth_date DATE,                      -- 出生日期
    email VARCHAR(100) UNIQUE,            -- 邮箱，唯一
    phone VARCHAR(20),                    -- 电话
    enrollment_date DATE DEFAULT CURRENT_DATE, -- 入学日期，默认为当前日期
    major VARCHAR(100),                   -- 专业
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 记录创建时间
);','','',1,'执行成功',239,'2025-12-01 00:50:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (334,2,'database_name','students',1,'SELECT * FROM students LIMIT 20 OFFSET 0','','',1,'执行成功',222,'2025-12-01 00:50:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (335,2,'database_name','students',2,'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (''1'', ''张三'', ''m'', ''2004-05-15'', ''123@456.com'', ''12345678911'', ''2004-05-15'', ''软件工程'')','','',2,'pq: new row for relation "students" violates check constraint "students_gender_check"',186,'2025-12-01 00:52:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (336,2,'database_name','students',2,'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (''20230001'', ''张三'', ''m'', ''2004-05-15'', ''123@456.com'', ''13800138001'', ''2004-05-15'', ''软件工程'')','','',2,'pq: new row for relation "students" violates check constraint "students_gender_check"',182,'2025-12-01 00:52:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (337,2,'database_name','students',2,'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (''20230001'', ''张三'', ''''''M'''''', ''2004-05-15'', ''123@456.com'', ''13800138001'', ''2004-05-15'', ''软件工程'')','','',2,'pq: value too long for type character(1)',214,'2025-12-01 00:52:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (338,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',542,'2025-12-01 00:52:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (339,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (''test123'', ''123@123.com'')','','',1,'执行成功',183,'2025-12-01 00:53:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (340,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-01 00:53:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (341,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',207,'2025-12-01 00:53:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (342,2,'database_name','users',1,'SELECT * FROM users LIMIT 100;','','',1,'执行成功',196,'2025-12-01 00:53:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (343,2,'database_name','users_copy_20251201003823',1,'SELECT * FROM users_copy_20251201003823 LIMIT 20 OFFSET 0','','',1,'执行成功',231,'2025-12-01 00:53:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (344,2,'database_name','',5,'DROP TABLE users_copy_20251201003823','','',1,'执行成功',184,'2025-12-01 00:54:05','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (345,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-12-01 10:18:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (346,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 100;','','',1,'执行成功',257,'2025-12-01 10:19:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (347,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-01 12:05:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (348,1,'RECOVER_YOUR_DATA','RECOVER_YOUR_DATA',1,'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-12-01 14:07:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (349,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-01 14:07:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (350,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-01 17:10:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (351,1,'gin-api','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',245,'2025-12-01 17:11:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (352,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-01 17:11:39','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (353,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-12-01 17:11:57','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (354,1,'gin-api','quick_deployment_tasks',1,'SELECT * FROM quick_deployment_tasks LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-12-01 17:12:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (355,1,'gin-api','quick_deployments',1,'SELECT * FROM quick_deployments LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-01 17:12:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (356,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-12-01 19:50:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (357,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',294,'2025-12-02 11:13:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (358,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-02 11:34:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (359,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-02 11:35:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (360,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-12-02 12:01:05','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (361,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 100;','','',1,'执行成功',555,'2025-12-02 12:06:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (362,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-02 12:08:01','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (363,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',250,'2025-12-02 12:25:37','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (364,1,'devops','app_application',1,'SELECT *
FROM  app_application
LIMIT  100;','','',1,'执行成功',252,'2025-12-02 12:25:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (365,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-12-02 12:30:05','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (366,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-02 12:35:09','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (367,2,'database_name','students',1,'SELECT * FROM students LIMIT 20 OFFSET 0','','',1,'执行成功',180,'2025-12-02 12:35:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (368,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-12-02 12:38:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (369,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',251,'2025-12-02 12:38:39','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (370,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',250,'2025-12-02 14:59:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (371,1,'gin-api','db',1,'SELECT * FROM db LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-02 15:00:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (372,1,'gin-api','db_instance',1,'SELECT * FROM db_instance LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-12-02 15:01:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (373,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',255,'2025-12-02 15:09:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (374,1,'gin-api','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',247,'2025-12-02 15:41:50','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (375,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',251,'2025-12-02 15:42:08','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (376,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',252,'2025-12-02 15:42:47','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (377,1,'RECOVER_YOUR_DATA','RECOVER_YOUR_DATA',1,'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-12-02 16:00:08','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (378,1,'mayfly-go','t_db_backup_history',1,'SELECT * FROM t_db_backup_history LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-12-02 16:00:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (379,1,'mayfly-go','t_db_restore',1,'SELECT * FROM t_db_restore LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-12-02 16:00:16','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (380,3,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-02 16:59:07','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (381,3,'gin-api','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-02 16:59:08','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (382,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-03 10:26:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (383,1,'gin-api','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',323,'2025-12-03 10:27:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (384,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-03 10:27:21','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (385,1,'gin-api','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',548,'2025-12-03 10:31:37','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (386,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',263,'2025-12-03 10:36:57','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (387,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',340,'2025-12-03 10:37:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (388,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-12-03 10:40:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (389,1,'devops','cmdb_host_1',2,'INSERT INTO cmdb_host_1 (id, host_name, group_id, private_ip, public_ip, ssh_name, ssh_key_id, ssh_port, remark, vendor, region, instance_id, os, status, cpu, memory, disk, billing_type, create_time, expire_time, update_time, ssh_ip, name, ssh_gateway_id) VALUES (506, ''华为云ops'', 4, ''172.31.6.35'', ''139.9.205.38'', ''root'', 22, 22, ''123'', 5, '''', '''', ''Ubuntu24.04.2'', 1, ''2'', ''2'', ''40'', '''', ''2025-11-11 17:24:36'', NULL, ''2025-11-26 17:22:07'', ''139.9.205.38'', ''hw-ops'', NULL);','','',2,'Error 1146 (42S02): Table ''devops.cmdb_host_1'' doesn''t exist',202,'2025-12-03 10:43:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (390,1,'devops','cmdb_host_1',2,'INSERT INTO cmdb_host_1 (id, host_name, group_id, private_ip, public_ip, ssh_name, ssh_key_id, ssh_port, remark, vendor, region, instance_id, os, status, cpu, memory, disk, billing_type, create_time, expire_time, update_time, ssh_ip, name, ssh_gateway_id) VALUES (506, ''华为云ops'', 4, ''172.31.6.35'', ''139.9.205.38'', ''root'', 22, 22, ''123'', 5, '''', '''', ''Ubuntu24.04.2'', 1, ''2'', ''2'', ''40'', '''', ''2025-11-11 17:24:36'', NULL, ''2025-11-26 17:22:07'', ''139.9.205.38'', ''hw-ops'', NULL);','','',2,'Error 1146 (42S02): Table ''devops.cmdb_host_1'' doesn''t exist',216,'2025-12-03 10:43:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (391,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',785,'2025-12-03 11:09:19','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (392,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-03 11:24:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (393,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-12-03 11:35:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (394,1,'mayfly-go','t_db',1,'SELECT * FROM t_db LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-12-03 11:40:08','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (395,1,'mayfly-go','t_db_backup',1,'SELECT * FROM t_db_backup LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-03 11:40:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (396,1,'mayfly-go','t_db_data_sync_log',1,'SELECT * FROM t_db_data_sync_log LIMIT 20 OFFSET 0','','',1,'执行成功',416,'2025-12-03 11:41:35','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (397,1,'mayfly-go','t_db_instance',1,'SELECT * FROM t_db_instance LIMIT 20 OFFSET 0','','',1,'执行成功',259,'2025-12-03 11:41:38','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (398,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-03 11:41:41','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (399,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-03 11:41:59','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (400,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',269,'2025-12-03 11:43:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (401,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',317,'2025-12-03 11:45:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (402,1,'devops','',5,'CREATE TABLE `app_jenkins_env123` ( `id` bigint(20) NOT NULL AUTO_INCREMENT, `app_id` bigint(20) NOT NULL, `env_name` varchar(50) NOT NULL, `jenkins_server_id` bigint(20), `job_name` varchar(255), `job_url` varchar(500), `build_params` json, `deploy_config` json, `notification` json, `is_active` tinyint(3) DEFAULT 1, `created_at` datetime, `updated_at` datetime, `deleted_at` datetime, PRIMARY KEY (id) ); ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_app_id`(`app_id`) USING BTREE; ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_deleted_at`(`deleted_at`) USING BTREE','','',2,'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ''ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_app_id`(`app_id`) U'' at line 1',204,'2025-12-03 11:47:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (403,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',270,'2025-12-03 11:48:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (404,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-03 11:52:55','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (405,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-03 11:52:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (406,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',260,'2025-12-03 11:53:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (407,1,'gin-api','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-12-03 11:56:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (408,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-03 12:04:05','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (409,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-03 12:07:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (410,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-03 12:08:37','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (411,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',327,'2025-12-03 12:08:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (412,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-12-03 12:14:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (413,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',342,'2025-12-03 12:16:54','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (414,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-03 12:19:56','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (415,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-03 12:23:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (416,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-03 12:26:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (417,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-03 12:28:20','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (418,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',244,'2025-12-03 12:29:04','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (419,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',395,'2025-12-03 12:32:00','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (420,1,'devops','',5,'CREATE TABLE `cmdb_group123` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT ''''''主键'''''',
 `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT ''''''父级分组ID'''''',
 `name` longtext NOT NULL COMMENT ''''''分组名称'''''',
 `create_time` datetime NOT NULL COMMENT ''''''创建时间'''''',
 `remark` longtext COMMENT ''''''备注'''''',
 `update_time` datetime COMMENT ''''''更新时间'''''', 
PRIMARY KEY (id)
)','','',1,'执行成功',244,'2025-12-03 12:42:43','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (421,1,'devops','cmdb_group123',1,'SELECT * FROM cmdb_group123 LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-12-03 12:42:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (422,1,'devops','',5,'DROP TABLE cmdb_group123','','',1,'执行成功',228,'2025-12-03 12:43:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (423,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-03 12:46:17','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (424,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-12-03 15:54:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (425,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',252,'2025-12-03 15:55:18','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (426,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-03 15:55:30','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (427,1,'gin-api','sys_operation_log',1,'SELECT *
FROM  sys_operation_log
LIMIT  100;','','',1,'执行成功',266,'2025-12-03 15:55:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (428,1,'gin-api','sys_admin',1,'SELECT *
FROM  sys_admin
LIMIT  100;','','',1,'执行成功',250,'2025-12-03 15:57:53','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (429,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-05 20:29:40','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (430,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-12-05 20:32:28','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (431,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',392,'2025-12-05 20:33:31','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (432,1,'devops','app_application',1,'SELECT *
FROM  app_application
LIMIT  100;','','',1,'执行成功',308,'2025-12-05 20:33:51','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (433,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',424,'2025-12-05 20:34:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (434,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-12-05 20:34:30','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (435,1,'devops','cmdb_host',3,'UPDATE cmdb_host SET host_name = ''虚拟机2025'', group_id = 4, private_ip = ''172.16.226.13'', public_ip = ''120.231.244.158'', ssh_name = ''root'', ssh_key_id = 13, ssh_port = 22, remark = ''123'', vendor = 1, region = '''', instance_id = '''', os = ''Ubuntu20.04.2'', status = 1, cpu = ''2'', memory = ''3'', disk = ''19'', billing_type = '''', create_time = ''2025-11-23 23:46:03'', expire_time = NULL, update_time = ''2025-11-23 23:46:07'', ssh_ip = ''172.16.226.13'', name = ''k8s-node02'', ssh_gateway_id = NULL WHERE id = 512;','','',1,'执行成功',436,'2025-12-05 20:34:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (436,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-05 20:34:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (437,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-05 20:34:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (438,1,'devops','',5,'CREATE DATABASE IF NOT EXISTS test1
  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',237,'2025-12-05 20:35:34','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (439,1,'test1','',5,'CREATE TABLE test1.your_table_name (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  -- your columns here
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;','','',2,'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near '') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4'' at line 4',219,'2025-12-05 20:35:52','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (440,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',217,'2025-12-05 20:36:57','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (441,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (''王五'', ''xxxx@123.com'')','','',1,'执行成功',191,'2025-12-05 20:37:14','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (442,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-12-05 20:37:15','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (443,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',260,'2025-12-06 21:54:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (444,1,'devops','sys_menu',1,'SELECT * FROM sys_menu LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-12-06 22:34:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (445,1,'devops','sys_role_menu',1,'SELECT * FROM sys_role_menu LIMIT 20 OFFSET 0','','',1,'执行成功',269,'2025-12-06 22:35:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (446,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',344,'2025-12-06 22:39:06','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (447,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',287,'2025-12-06 22:58:13','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (448,1,'RECOVER_YOUR_DATA','RECOVER_YOUR_DATA',1,'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0','','',1,'执行成功',290,'2025-12-06 22:58:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (449,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-12-06 22:58:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (450,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',463,'2025-12-06 22:59:49','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (451,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',317,'2025-12-06 23:01:23','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (452,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',313,'2025-12-06 23:04:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (453,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',518,'2025-12-06 23:06:58','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (454,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',512,'2025-12-06 23:07:56','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (455,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',349,'2025-12-06 23:10:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (456,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',529,'2025-12-06 23:11:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (457,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',544,'2025-12-06 23:11:46','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (458,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-06 23:18:44','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (459,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',573,'2025-12-06 23:19:02','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (460,1,'test1','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',543,'2025-12-06 23:23:48','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (461,1,'test1','tool_service_deploy',1,'SELECT * FROM tool_service_deploy LIMIT 20 OFFSET 0','','',1,'执行成功',317,'2025-12-06 23:45:11','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (462,1,'test1','tool_link',1,'SELECT * FROM tool_link LIMIT 20 OFFSET 0','','',1,'执行成功',301,'2025-12-06 23:45:22','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (463,1,'test1','task_work',1,'SELECT * FROM task_work LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-12-06 23:45:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (464,1,'test1','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',316,'2025-12-07 00:02:12','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (465,1,'test1','RECOVER_YOUR_DATA',1,'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0','','',1,'执行成功',281,'2025-12-07 00:09:10','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (466,1,'test1','',5,'DROP TABLE RECOVER_YOUR_DATA','','',1,'执行成功',279,'2025-12-07 00:09:21','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (467,1,'gin-api','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',312,'2025-12-07 00:09:45','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (468,1,'devops','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-12-07 00:27:44','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (469,1,'devops','db_es_instance',1,'SELECT * FROM db_es_instance LIMIT 20 OFFSET 0','','',1,'执行成功',337,'2025-12-07 00:29:33','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (470,1,'test2','',5,'CREATE DATABASE IF NOT EXISTS test3
  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',512,'2025-12-07 16:43:30','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (471,1,'test2','',5,'CREATE DATABASE IF NOT EXISTS test3
  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',228,'2025-12-07 16:44:03','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (472,1,'test3','',5,'CREATE DATABASE IF NOT EXISTS test4
  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',221,'2025-12-07 16:51:25','',0);
INSERT INTO `db_sql_exec` (`id`,`db_id`,`db_name`,`table_name`,`type`,`sql`,`old_value`,`remark`,`status`,`res`,`exec_time`,`create_time`,`creator`,`creator_id`) VALUES (473,1,'test3','',5,'CREATE DATABASE IF NOT EXISTS test4
  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',781,'2025-12-07 17:01:52','',0);

-- Table structure for k8s_cluster
CREATE TABLE `k8s_cluster` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
 `name` varchar(100) NOT NULL COMMENT '''集群名称''',
 `version` varchar(50) NOT NULL COMMENT '''集群版本''',
 `status` bigint(19) NOT NULL DEFAULT 1 COMMENT '''集群状态:1-创建中,2-运行中,3-离线''',
 `credential` text COMMENT '''集群凭证(kubeconfig)''',
 `description` text COMMENT '''集群描述''',
 `cluster_type` bigint(19) NOT NULL DEFAULT 1 COMMENT '''集群类型:1-自建,2-导入''',
 `created_at` datetime,
 `updated_at` datetime,
 `node_count` bigint(19) DEFAULT 0 COMMENT '''节点数量''',
 `ready_nodes` bigint(19) DEFAULT 0 COMMENT '''就绪节点数''',
 `master_nodes` bigint(19) DEFAULT 0 COMMENT '''Master节点数''',
 `worker_nodes` bigint(19) DEFAULT 0 COMMENT '''Worker节点数''',
 `last_sync_at` datetime COMMENT '''最后同步时间''', 
PRIMARY KEY (id)
);
ALTER TABLE `k8s_cluster` ADD unique INDEX `idx_k8s_cluster_name`(`name`) USING BTREE;

-- Data for table k8s_cluster
INSERT INTO `k8s_cluster` (`id`,`name`,`version`,`status`,`credential`,`description`,`cluster_type`,`created_at`,`updated_at`,`node_count`,`ready_nodes`,`master_nodes`,`worker_nodes`,`last_sync_at`) VALUES (33,'腾讯云测试环境k8s集群','v1.22.5-tke.40',2,'apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeU1URXhOekV4TlRnd01sb1hEVE15TVRFeE5ERXhOVGd3TWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS1NyCjNwbDg2VExCT2c5Z1NyWFhLbC9PZWtTMDFtQllxM21GNFZmMnVOSmcxbFAybEQyZFJEN0JjTThoNVJYUUFwa2cKcDVmUU9vV1cxcHdsVHZpdVY5L0U2YjViSVJ3c1JwVThiVmthN283aVlxekZIS3ZLMXh5ay9WM2g0N2k1R0xkeQp5U2RqUW5WWTlOeVBRUDQ3ZGhKUGZGb21NbTMyQUkwMGIvZUNIVEdzblp5YVZSYW1GSG51ZXNkQTZIY05SbEJ2CkdIaS92TndtR3F5bldoR29jMXFMZEwxVnVDV1BPV1RwZEdCNGpydHByTmw0SnVwS1RCUjRxc0FVVHFCWEN5a3QKV3c0VGFCT0ZkeEJidDZ0d2RYVnUzbTc5dGNHd0cyR0V1UDVETzN3QVlZdVRJaXk4ZDJxLzR6dWJhbVA3T1hxMwppMVN0aktOUk9EUDBobnlubElNQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0tVTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFHL3BkV0hQdGdHTUxPRE9DNDRQd21IRndvVjAKbTNNWHIxN2JyT1ZFdFNZOVNEMVNlbUFYVUEwRU5wS2kxdHFrQVU3TzkzL0ZLQmVPekpkRnhaZzN4VDdzUGhhKwpCVlJ3MjYyV1VHamlFSERHNys3N3AwaWNxQlZidEUvMS9oM3EveU8wOUIxSmg5TTRVTURveVlEZXlKNTFsUFBXCnUyRjRmOTUzeGZUa3hKZmJMNVlSU3d5R3lEbW9LSy9ZWUdkQ25LeHlCc3lDTTl5RG1FbnZLcW5HOHh6UlF6QVEKRzcxckJhNHpQOTZOcXVBZVZwSVBvdnFxQ0ZIay9ydHg4SVNna2R4MXVtb1Vhb3hWREpkRU94TU1iSkt6dzFGTwp3MUIvWkpCS1llM2hNS3VIdVJ2NTRLeG5NeEIvbW5iVm8wWnhrTThERDlUOEhGYVJIR1p3TVFCbmF6ST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    server: https://10.9.0.13
  name: cls-301g0fi0
contexts:
- context:
    cluster: cls-301g0fi0
    user: "100028461074"
  name: cls-301g0fi0-100028461074-context-default
current-context: cls-301g0fi0-100028461074-context-default
kind: Config
preferences: {}
users:
- name: "100028461074"
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURERENDQWZTZ0F3SUJBZ0lJTDRsK3Z2MW9Gc2N3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TWpFeE1UY3hNakF3TlRGYUZ3MDBNakV4TVRjeE1qQXdOVEZhTURZeApFakFRQmdOVkJBb1RDWFJyWlRwMWMyVnljekVnTUI0R0ExVUVBeE1YTVRBd01ESTRORFl4TURjMExURTJOamcyCk9EWTBOVEV3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQzlEY0tPRlZBdGRhVEoKS0V6YXFvOWt1WWFRSU5GZnYrS3JoWDU2UzY1Qm1sVVp2UEhlZzRxZWo1Qy85OTg2K1luZVRHRnpTTExUMEp5WApwcGhlVjVMNE5KMFlsZXI3U3pNMzZtSUJIODZ6RnI3QUVTSndHUGxtdXg3RzBnakZjYkQzL3Z1REFYV082aVBaCkF6L3p2NGxGOUlQUlpTViswcW04alg0ZE9DTW5aMXhXV29UZW1Vc2czYktiUys3YWNOR2xqNnJNUzdWSmJRYXIKTXdiSmZaaUk1UDhGblZkUjNIbzQwSXYwNWFOZnBvSHVWRjcwUS9KYXJESmhMYTBTdWI4dU5nbjhDRWJ0TkJINQplMHlMdlF4ZTF2Yks4eWI0UGJDaGJHV3ZONENQWVVhVnB1Z05mK3FHNzR5cWo5cklNbWVMTkcycDVlU0wzT2N5CnU3NnhidTVGQWdNQkFBR2pQekE5TUE0R0ExVWREd0VCL3dRRUF3SUNoREFkQmdOVkhTVUVGakFVQmdnckJnRUYKQlFjREFnWUlLd1lCQlFVSEF3RXdEQVlEVlIwVEFRSC9CQUl3QURBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQQpUNHZlNnVyb2p1VWxZbmluSkI5Nkc0OFI5RFA2TGo5M0Y3TEdJMVlvbnZpTmVlaUxsUDhTYjNhR09jT215ME55CnNCNFhQOS9vdGY1T2Q2dnJ5QTJmT25UbDdFTFdqTGMyY3UvbGJxQ0lCQytuMjV2YzM1ZjIxRm5UOWpwQjZ5ZmcKR0lCWUh1eW9kc1RNbEhPaG01d0FPNkxNSENiUTh2cGdmVXlOUEZrc2ZrcFhPbDk5VGdwNGtGLzEveStLL0NnRwpHWmNaZUxkNUFkWVFpbXBtQS9QQzBmZ3ptUW85eVJZWktJYmNUMXBjcTI3cWc3R1RYb3N5MjYwdmU3eWgwSmxRCm8vUStmdEFEVDkvZXNTYnNLRWJrLzVkNC91Q0FIQjZJWkhaSmsvSkEyeGtQdGo4eVRySlUxQkxQdmg5L0NWM0gKOFJkSGdQRnZ0U0dTaWhXTDFXOXgyQT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcGdJQkFBS0NBUUVBdlEzQ2poVlFMWFdreVNoTTJxcVBaTG1Ha0NEUlg3L2lxNFYrZWt1dVFacFZHYnp4CjNvT0tubytRdi9mZk92bUoza3hoYzBpeTA5Q2NsNmFZWGxlUytEU2RHSlhxKzBzek4rcGlBUi9Pc3hhK3dCRWkKY0JqNVpyc2V4dElJeFhHdzkvNzdnd0YxanVvajJRTS84NytKUmZTRDBXVWxmdEtwdkkxK0hUZ2pKMmRjVmxxRQozcGxMSU4yeW0wdnUybkRScFkrcXpFdTFTVzBHcXpNR3lYMllpT1QvQloxWFVkeDZPTkNMOU9Xalg2YUI3bFJlCjlFUHlXcXd5WVMydEVybS9MallKL0FoRzdUUVIrWHRNaTcwTVh0YjJ5dk1tK0Qyd29XeGxyemVBajJGR2xhYm8KRFgvcWh1K01xby9heURKbml6UnRxZVhraTl6bk1ydStzVzd1UlFJREFRQUJBb0lCQVFDR2hNdnluV2UzY1FZWAp1LzB6ZWprZHlvVGMrUDBCekY5VS9GL0tWUjlndnBmZkFIcDVnYytZa3o3RDJYT1ZXd21taTJLcnpRYXhnTFpyCmtiYkpLamNCTlM2eTE3bWpwUmFFcXRtaFF5aW9kWDYwMjN4Tmx2ZWZDWFlVc2s1YmIxSWtvK0lDalg5QWVZbVkKV0orZ3FLRGNqZWRkajBMbXdxN1Q0em5oUUppQzUrZkx5SUdrWFhsclcwRC9Zem53VXE5ZU1abnF2Yk56MEZ2TgpOZm96R0tTemRFVjMwZXhVYUEyRElmSDFNSEF4WnZyanlWQkJuOVgwbWhOMFo4KzFmQ1Y0TVFkcy9VVFFndmIyCmlheGhyWnF3MWtxTHEvdjBZd212UDRqaGdmb3BzeG1PaDNkVFJxL2l6VVNtWlJCd2pDUk5SN1J2M3FKY1FyWEoKRVJhZ1Nic0JBb0dCQU42R0lZdjNNQ1FORjYzcFdqSzhSdkd5TjcyR1dqM2JGa1hkZTk2bDNXOExWVHlUVGg5QgpsUXlyS0IxYUdua0hQMEN1bjA4VWE4NkVZQ3dGSEhpK0F2RWZsbVlZN1NFTnZyR3BKQXJBUEFVbnJzL25TUFhvCkRMV2plZVJwSHVqNkFqY3M1d1puUUpaaERiY21XNkZFSGdyNnJYWFFFekR0YlJXN2ErL2JPWGdkQW9HQkFObCsKbnVDQXBZMXlRcGM4Y0lVUms5bjVEQ3JIYXFXck1GT2l6Q2NxeE1xRU8xZnNEZTcrSGQ3am0rNXBhbzgzTktNcgo2OE1XeVI5S0JZa2M1UXM4K2gyUnhhR3p2SDUzNW5nRmdDRURjZkN4bzhOeWlndWtqWklsZDVPcmZPR25UaGIvCmcvMFZVK2kzTGYzUU41emZVUkx6R2hYT05qVldlSUM1TWdoek53WkpBb0dCQU4zTUhpWHZsWWxwajErbVIyUm8KcStIL1BIWHBNNHVTRTJwQjhGU2xTTWlGS0VsL25nK2pnOWdqeFpFMEhVempQMVhPZURZUGFjRm55cUExM0FJeAppcWNpVC91bC9yM0NnSDkrMHREd2VTSytRZzB5UUZQekhxL3hxTzVWSlIwZERSWkg2R0pHaW1jWmg5MDFHNEhyCmJuS05HRlprdHF3L0MwcnREQklBbGhDMUFvR0JBTEdRaDZOZ1JFamJqUFF0K3VzWTg2MlVmMU41cFY4R0NSbHYKSXQrak1qRG9qZzUxQU5jcWpqYUhlcW5GVktNQi9pNXRVRGpoSkdnd0paN1FnVXFRaU42VlFWV0lDUTd6VXF5UQpTUkFxNE40U0dWSWE5QXNYV0NHR0J6Ny9WYUJnMkFjOHFvUjM0aDV2VnRRT0p5a2l4TzBFbmpwSk5jYU1CeXVJCmh3bVJEczVKQW9HQkFJMVY5cVM2UGI5eGRjYkcwUkhHYVAxczlwcE5rQmJrSGgxalN4cjRBWWlTUTNxalIxeCsKNVR6VXVNcC94UlQ0bGxBTkZXd3dVa0k0dnVLSDlwcHJSMzVrYlg2dkYvcmN2V0Y0TGNpZDdUdkdIU2tGR2Evagp1a1dMSTdKQlV3dFc3QWZqWlduR1NPcm95OFl5SmlGV0dnYisranEzUllrL2h1cWhBOEdWTXROcgotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=','',2,'2025-11-14 19:51:40','2025-11-26 16:35:00',8,8,0,8,'2025-11-14 19:51:41');

-- Table structure for monitor_agent
CREATE TABLE `monitor_agent` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `host_id` bigint(20) NOT NULL COMMENT '''主机ID''',
 `host_name` longtext COMMENT '''主机名称''',
 `version` varchar(191) DEFAULT 1.0.0 COMMENT '''Agent版本''',
 `status` bigint(19) COMMENT '''状态:1->部署中,2->部署失败,3->运行中,4->已停止''',
 `install_path` longtext COMMENT '''安装路径''',
 `port` bigint(19) DEFAULT 9100 COMMENT '''监听端口''',
 `pid` bigint(19) COMMENT '''进程ID''',
 `last_heartbeat` datetime COMMENT '''最后心跳时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `error_msg` text COMMENT '''错误信息''',
 `install_progress` bigint(19) DEFAULT 0 COMMENT '''安装进度(0-100)''', 
PRIMARY KEY (id)
);
ALTER TABLE `monitor_agent` ADD  INDEX `idx_monitor_agent_host_id`(`host_id`) USING BTREE;

-- Data for table monitor_agent
INSERT INTO `monitor_agent` (`id`,`host_id`,`host_name`,`version`,`status`,`install_path`,`port`,`pid`,`last_heartbeat`,`update_time`,`create_time`,`error_msg`,`install_progress`) VALUES (84,501,'go-ops','1.0.0',3,'/opt/agent',9100,0,NULL,'2025-12-03 16:21:09','2025-12-03 16:20:19','',100);
INSERT INTO `monitor_agent` (`id`,`host_id`,`host_name`,`version`,`status`,`install_path`,`port`,`pid`,`last_heartbeat`,`update_time`,`create_time`,`error_msg`,`install_progress`) VALUES (85,506,'hw-ops','1.0.0',3,'/opt/agent',9100,0,NULL,'2025-12-03 16:20:59','2025-12-03 16:20:19','',100);

-- Table structure for monitor_domain
CREATE TABLE `monitor_domain` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `domain` varchar(255) NOT NULL,
 `tags` varchar(500),
 `remark` text,
 `status` bigint(19) DEFAULT 1 COMMENT '状态(1:启用,0:禁用)',
 `is_alive` bigint(19) DEFAULT 0 COMMENT '存活状态(1:正常,0:异常)',
 `status_code` bigint(19) COMMENT 'HTTP状态码',
 `response_time` bigint(19) COMMENT '响应时间(ms)',
 `ssl_expire_at` datetime COMMENT 'SSL证书过期时间',
 `ssl_days_left` bigint(19) COMMENT 'SSL证书剩余天数',
 `ssl_issuer` varchar(255) COMMENT 'SSL证书颁发者',
 `last_check_at` datetime COMMENT '最后检查时间',
 `error_msg` text COMMENT '错误信息',
 `create_time` datetime,
 `update_time` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `monitor_domain` ADD unique INDEX `idx_monitor_domain_domain`(`domain`) USING BTREE;

-- Data for table monitor_domain
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (2,'ops-flashcat.dding.net','bline=ops','bline=ops',1,1,200,405,'2027-01-03 23:59:59',394,'Encryption Everywhere DV TLS CA - G2','2025-12-06 01:14:25','','2025-12-03 21:44:12','2025-12-06 01:14:24');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (4,'www.baidu.com','bline=saas3,env=prod','123',1,1,200,163,'2026-08-10 07:01:01',247,'GlobalSign RSA OV SSL CA 2018','2025-12-06 01:14:25','','2025-12-03 21:57:12','2025-12-06 01:14:24');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (5,'ai.feishu.cn','test','',1,1,403,637,'2026-03-25 23:59:59',110,'RapidSSL TLS RSA CA G1','2025-12-06 01:14:25','','2025-12-04 10:20:57','2025-12-06 01:14:25');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (6,'gitee.com','','',1,1,200,704,'2026-02-17 23:59:59',74,'TrustAsia DV TLS RSA CA 2025','2025-12-06 01:14:25','','2025-12-04 10:21:13','2025-12-06 01:14:25');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (7,'ops.dding.net','','',1,1,200,444,'2025-12-20 23:59:59',15,'Encryption Everywhere DV TLS CA - G1','2025-12-06 01:14:25','','2025-12-04 10:21:28','2025-12-06 01:14:24');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (8,'docker.aityp.com','','',1,1,200,952,'2025-12-28 23:58:52',23,'E8','2025-12-06 01:14:25','','2025-12-04 10:21:54','2025-12-06 01:14:25');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (9,'test-ops-jenkins-tc2.dding.net','','',1,1,200,229,NULL,0,'','2025-12-06 01:14:25','','2025-12-04 10:22:35','2025-12-06 01:14:24');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (10,'www.bejson.com','','',1,1,200,224,'2026-01-21 23:59:59',47,'ZeroSSL ECC Domain Secure Site CA','2025-12-06 01:14:25','','2025-12-04 10:23:17','2025-12-06 01:14:24');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (11,'www.deviops.com','prod','',1,1,200,4318,'2026-01-22 03:15:51',47,'WE1','2025-12-06 01:14:29','','2025-12-04 13:30:23','2025-12-06 01:14:28');
INSERT INTO `monitor_domain` (`id`,`domain`,`tags`,`remark`,`status`,`is_alive`,`status_code`,`response_time`,`ssl_expire_at`,`ssl_days_left`,`ssl_issuer`,`last_check_at`,`error_msg`,`create_time`,`update_time`) VALUES (12,'www.jd.com','ops','123',1,1,200,154,'2025-12-15 07:09:54',9,'GlobalSign RSA OV SSL CA 2018','2025-12-06 01:14:25','','2025-12-05 20:50:58','2025-12-06 01:14:24');

-- Table structure for monitor_domain_schedule
CREATE TABLE `monitor_domain_schedule` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `enabled` tinyint(3) DEFAULT 0 COMMENT '是否启用',
 `interval_mins` bigint(19) COMMENT '检查间隔(分钟)',
 `next_run_at` datetime COMMENT '下次执行时间',
 `last_run_at` datetime COMMENT '上次执行时间',
 `status` varchar(50) COMMENT '任务状态(running/paused/stopped)',
 `create_time` datetime,
 `update_time` datetime, 
PRIMARY KEY (id)
);

-- Data for table monitor_domain_schedule
INSERT INTO `monitor_domain_schedule` (`id`,`enabled`,`interval_mins`,`next_run_at`,`last_run_at`,`status`,`create_time`,`update_time`) VALUES (1,0,2,NULL,NULL,'stopped','2025-12-04 23:09:32','2025-12-05 23:12:27');
INSERT INTO `monitor_domain_schedule` (`id`,`enabled`,`interval_mins`,`next_run_at`,`last_run_at`,`status`,`create_time`,`update_time`) VALUES (2,0,120,NULL,NULL,'stopped','2025-12-04 23:09:32','2025-12-04 23:09:32');

-- Table structure for quick_deployment_tasks
CREATE TABLE `quick_deployment_tasks` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `deployment_id` bigint(20) NOT NULL,
 `app_id` bigint(20) NOT NULL,
 `app_name` varchar(255),
 `app_code` varchar(255),
 `environment` varchar(50),
 `jenkins_env_id` bigint(20) NOT NULL,
 `jenkins_job_url` varchar(500),
 `build_number` bigint(19),
 `status` tinyint(3) DEFAULT 1,
 `execute_order` bigint(19) NOT NULL,
 `start_time` datetime,
 `end_time` datetime,
 `duration` bigint(19),
 `error_message` text,
 `log_url` varchar(500),
 `created_at` datetime,
 `updated_at` datetime,
 `parameters` text, 
PRIMARY KEY (id)
);
ALTER TABLE `quick_deployment_tasks` ADD  INDEX `fk_quick_deployment_tasks_application`(`app_id`) USING BTREE;
ALTER TABLE `quick_deployment_tasks` ADD  INDEX `fk_quick_deployment_tasks_jenkins_env`(`jenkins_env_id`) USING BTREE;
ALTER TABLE `quick_deployment_tasks` ADD  INDEX `idx_quick_deployment_tasks_deployment_id`(`deployment_id`) USING BTREE;

-- Data for table quick_deployment_tasks
INSERT INTO `quick_deployment_tasks` (`id`,`deployment_id`,`app_id`,`app_name`,`app_code`,`environment`,`jenkins_env_id`,`jenkins_job_url`,`build_number`,`status`,`execute_order`,`start_time`,`end_time`,`duration`,`error_message`,`log_url`,`created_at`,`updated_at`,`parameters`) VALUES (30,21,14,'dev-lku-sass-pack','dev-lku-sass-pack','test2',38,'',407,3,1,'2025-11-14 19:36:56','2025-11-14 19:41:03',246,'','http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/407/console','2025-11-14 19:35:44','2025-11-14 19:41:04','{"commit_id":"3dbb6f7c83ab1d1bd5928183f8cf1e79bd3b2f3f","compile":"true"}');
INSERT INTO `quick_deployment_tasks` (`id`,`deployment_id`,`app_id`,`app_name`,`app_code`,`environment`,`jenkins_env_id`,`jenkins_job_url`,`build_number`,`status`,`execute_order`,`start_time`,`end_time`,`duration`,`error_message`,`log_url`,`created_at`,`updated_at`,`parameters`) VALUES (31,21,15,'dev-lku-sass-deploy','dev-lku-sass-deploy','dev',41,'',370,3,2,'2025-11-14 19:41:04','2025-11-14 19:42:00',56,'','http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-deploy/370/console','2025-11-14 19:35:44','2025-11-14 19:42:00','{"app_name":"swagger","commit_id":"3dbb6f7c83ab1d1bd5928183f8cf1e79bd3b2f3f"}');
INSERT INTO `quick_deployment_tasks` (`id`,`deployment_id`,`app_id`,`app_name`,`app_code`,`environment`,`jenkins_env_id`,`jenkins_job_url`,`build_number`,`status`,`execute_order`,`start_time`,`end_time`,`duration`,`error_message`,`log_url`,`created_at`,`updated_at`,`parameters`) VALUES (43,33,19,'prod_saas3.0_data-export','prod-saas30-data-export','test',52,'',9,3,1,'2025-12-05 20:44:14','2025-12-05 20:44:26',11,'','http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/9/console','2025-12-05 20:43:58','2025-12-05 20:44:26','{"commit_id":"123456789"}');

-- Table structure for quick_deployments
CREATE TABLE `quick_deployments` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `title` varchar(255) NOT NULL,
 `business_group_id` bigint(20) NOT NULL,
 `business_dept_id` bigint(20) NOT NULL,
 `description` text,
 `status` tinyint(3) DEFAULT 1,
 `task_count` bigint(19) NOT NULL DEFAULT 0,
 `creator_id` bigint(20) NOT NULL,
 `creator_name` varchar(100),
 `start_time` datetime,
 `end_time` datetime,
 `duration` bigint(19),
 `created_at` datetime,
 `updated_at` datetime,
 `execution_mode` tinyint(3) DEFAULT 1, 
PRIMARY KEY (id)
);

-- Data for table quick_deployments
INSERT INTO `quick_deployments` (`id`,`title`,`business_group_id`,`business_dept_id`,`description`,`status`,`task_count`,`creator_id`,`creator_name`,`start_time`,`end_time`,`duration`,`created_at`,`updated_at`,`execution_mode`) VALUES (19,'lockin-测试环境打包',18,2,'lockin-测试环境打包',3,1,1,'管理员','2025-11-11 11:18:31','2025-11-11 11:20:57',144,'2025-11-11 11:05:34','2025-11-11 11:20:58',2);
INSERT INTO `quick_deployments` (`id`,`title`,`business_group_id`,`business_dept_id`,`description`,`status`,`task_count`,`creator_id`,`creator_name`,`start_time`,`end_time`,`duration`,`created_at`,`updated_at`,`execution_mode`) VALUES (20,'lockin-研发环境打包',18,5,'lockin-研发环境打包',3,1,1,'管理员','2025-11-11 11:12:52','2025-11-11 11:15:33',159,'2025-11-11 11:10:07','2025-11-11 11:15:33',2);
INSERT INTO `quick_deployments` (`id`,`title`,`business_group_id`,`business_dept_id`,`description`,`status`,`task_count`,`creator_id`,`creator_name`,`start_time`,`end_time`,`duration`,`created_at`,`updated_at`,`execution_mode`) VALUES (21,'新加坡saas服务打包k8s发布',18,12,'快速发布',3,2,1,'管理员','2025-11-14 19:36:56','2025-11-14 19:42:01',302,'2025-11-14 19:35:44','2025-11-14 19:42:01',2);
INSERT INTO `quick_deployments` (`id`,`title`,`business_group_id`,`business_dept_id`,`description`,`status`,`task_count`,`creator_id`,`creator_name`,`start_time`,`end_time`,`duration`,`created_at`,`updated_at`,`execution_mode`) VALUES (33,'saas3.0服务快速上线',18,2,'123',3,1,1,'管理员','2025-12-05 20:44:13','2025-12-05 20:44:26',11,'2025-12-05 20:43:57','2025-12-05 20:44:26',2);

-- Table structure for redis_instance
CREATE TABLE `redis_instance` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `code` varchar(36) NOT NULL COMMENT '''实例编码''',
 `name` varchar(100) NOT NULL COMMENT '''实例名称''',
 `mode` varchar(20) NOT NULL COMMENT '''模式:standalone,cluster,sentinel''',
 `host` varchar(300) NOT NULL COMMENT '''主机: standalone为host:port, cluster为逗号分隔, sentinel为master=hosts''',
 `port` bigint(19) DEFAULT 0 COMMENT '''端口(standalone可用)''',
 `db` bigint(19) DEFAULT 0 COMMENT '''默认库号''',
 `username` varchar(100) COMMENT '''用户名(可选)''',
 `password` varchar(500) NOT NULL COMMENT '''密码(加密)''',
 `redis_node_password` varchar(500) COMMENT '''节点密码(仅sentinel)''',
 `remark` varchar(500) COMMENT '''备注''',
 `ssh_tunnel_machine_id` bigint(20) DEFAULT 0 COMMENT '''SSH隧道机器ID''',
 `status` bigint(19) DEFAULT 1 COMMENT '''状态:1->启用,2->禁用''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''',
 `creator` varchar(64) COMMENT '''创建人''',
 `creator_id` bigint(20) COMMENT '''创建人ID''',
 `modifier` varchar(64) COMMENT '''修改人''',
 `modifier_id` bigint(20) COMMENT '''修改人ID''', 
PRIMARY KEY (id)
);
ALTER TABLE `redis_instance` ADD unique INDEX `idx_redis_instance_code`(`code`) USING BTREE;

-- Data for table redis_instance
-- No data

-- Table structure for sys_admin
CREATE TABLE `sys_admin` (
 `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键',
 `post_id` int(10) COMMENT '岗位id',
 `dept_id` int(10) COMMENT '部门id',
 `username` varchar(64) NOT NULL COMMENT '账号',
 `password` varchar(64) NOT NULL COMMENT '密码',
 `nickname` varchar(64) COMMENT '昵称',
 `icon` varchar(500) COMMENT '头像',
 `email` varchar(64) COMMENT '邮箱',
 `phone` varchar(64) COMMENT '手机',
 `note` varchar(500) COMMENT '备注信息',
 `create_time` datetime NOT NULL COMMENT '创建时间',
 `status` int(10) NOT NULL DEFAULT 1 COMMENT '帐号启用状态：1->启用,2->禁用', 
PRIMARY KEY (id)
) COMMENT '后台管理员表';
ALTER TABLE `sys_admin` ADD  INDEX `create_time`(`create_time`) USING BTREE;
ALTER TABLE `sys_admin` ADD unique INDEX `username`(`username`) USING BTREE;

-- Data for table sys_admin
INSERT INTO `sys_admin` (`id`,`post_id`,`dept_id`,`username`,`password`,`nickname`,`icon`,`email`,`phone`,`note`,`create_time`,`status`) VALUES (89,1,15,'admin','e10adc3949ba59abbe56e057f20f883e','admin','http://10.7.16.22:8080/api/v1/upload/20251010/702910000.png','123456789@qq.com','13754354536','后端研发','2023-05-23 22:15:50',1);
INSERT INTO `sys_admin` (`id`,`post_id`,`dept_id`,`username`,`password`,`nickname`,`icon`,`email`,`phone`,`note`,`create_time`,`status`) VALUES (98,12,15,'lisi','e10adc3949ba59abbe56e057f20f883e','李四','http://127.0.0.1:8000/upload/20250629/500333800.png','123@qq.com','13826541566','ops','2025-06-20 11:54:02',1);
INSERT INTO `sys_admin` (`id`,`post_id`,`dept_id`,`username`,`password`,`nickname`,`icon`,`email`,`phone`,`note`,`create_time`,`status`) VALUES (102,12,14,'test','e10adc3949ba59abbe56e057f20f883e','test','http://10.7.16.22:8080/api/v1/upload/20251009/699174000.png','zfwh1024@163.com','13826541511','游客','2025-09-24 12:49:06',1);
INSERT INTO `sys_admin` (`id`,`post_id`,`dept_id`,`username`,`password`,`nickname`,`icon`,`email`,`phone`,`note`,`create_time`,`status`) VALUES (103,10,14,'zhangsan','e10adc3949ba59abbe56e057f20f883e','张三','','zfwh1024@163.com','11111111111','123','2025-11-22 19:26:31',1);
INSERT INTO `sys_admin` (`id`,`post_id`,`dept_id`,`username`,`password`,`nickname`,`icon`,`email`,`phone`,`note`,`create_time`,`status`) VALUES (104,10,16,'wangwu','e10adc3949ba59abbe56e057f20f883e','王五','','zfwh1024@163.com','11111111111','123','2025-11-22 19:27:30',1);

-- Table structure for sys_admin_role
CREATE TABLE `sys_admin_role` (
 `admin_id` int(10) NOT NULL COMMENT '管理员id',
 `role_id` int(10) NOT NULL COMMENT '角色id'
) COMMENT '管理员和角色关系表';

-- Data for table sys_admin_role
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (89,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (102,13);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (104,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (98,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (103,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (89,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (102,13);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (104,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (98,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (103,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (89,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (102,13);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (104,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (98,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (103,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (89,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (102,13);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (104,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (98,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (103,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (89,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (102,13);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (104,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (98,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (103,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (89,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (102,13);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (104,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (98,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (103,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (89,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (102,13);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (104,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (98,1);
INSERT INTO `sys_admin_role` (`admin_id`,`role_id`) VALUES (103,1);

-- Table structure for sys_dept
CREATE TABLE `sys_dept` (
 `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键',
 `parent_id` int(10) NOT NULL COMMENT '父id',
 `dept_type` int(10) NOT NULL COMMENT '部门类型（1->公司, 2->中心，3->部门）',
 `dept_name` varchar(30) NOT NULL COMMENT '部门名称',
 `dept_status` int(10) NOT NULL DEFAULT 1 COMMENT '部门状态（1->正常 2->停用）',
 `create_time` datetime NOT NULL COMMENT '创建时间', 
PRIMARY KEY (id)
) COMMENT '部门表';
ALTER TABLE `sys_dept` ADD unique INDEX `dept_name`(`dept_name`) USING BTREE;

-- Data for table sys_dept
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (1,0,1,'神舟科技有限公司',1,'2023-06-14 17:53:23');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (2,1,2,'深圳研发中心',1,'2023-06-14 17:53:55');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (3,2,3,'架构设计部门',1,'2023-06-14 17:54:15');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (5,2,3,'后端研发部门',1,'2023-06-14 17:55:25');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (6,2,3,'系统测试部门',1,'2023-06-14 17:55:31');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (12,1,2,'北京研发中心',1,'2025-06-28 23:42:46');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (13,1,2,'重庆研发中心',1,'2025-06-28 23:43:15');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (14,12,3,'运维1部',1,'2025-06-28 23:43:34');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (15,13,3,'运维2部',1,'2025-06-28 23:44:15');
INSERT INTO `sys_dept` (`id`,`parent_id`,`dept_type`,`dept_name`,`dept_status`,`create_time`) VALUES (16,13,3,'重庆研发部-001',1,'2025-07-04 13:10:58');

-- Table structure for sys_login_info
CREATE TABLE `sys_login_info` (
 `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'ID',
 `username` varchar(50) COMMENT '用户账号',
 `ip_address` varchar(128) COMMENT '登录IP地址',
 `login_location` varchar(255) COMMENT '登录地点',
 `browser` varchar(500) COMMENT '浏览器类型',
 `os` varchar(50) COMMENT '操作系统',
 `login_status` int(10) COMMENT '登录状态（1-成功 2-失败）',
 `message` varchar(255) COMMENT '提示消息',
 `login_time` datetime COMMENT '访问时间', 
PRIMARY KEY (id)
) COMMENT '登录日志记录';

-- Data for table sys_login_info
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (1,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-26 15:10:07');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (2,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-06-27 10:11:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (3,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 10:11:22');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (4,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 10:35:36');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (5,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 10:39:36');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (6,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-06-27 10:51:17');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (7,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 10:51:22');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (8,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 11:06:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (9,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 13:56:34');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (10,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 14:56:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (11,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-27 19:19:50');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (12,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',2,'密码不正确','2025-06-28 14:29:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (13,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-28 14:29:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (14,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-28 15:57:37');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (15,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-28 15:58:50');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (16,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-28 16:01:59');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (17,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-28 19:06:18');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (18,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-28 23:36:22');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (23,'lisi','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-29 15:23:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (24,'test','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',2,'密码不正确','2025-06-29 15:23:55');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (25,'test','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-29 15:24:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (26,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-29 15:24:42');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (27,'lisi','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-29 15:25:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (28,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',2,'密码不正确','2025-06-29 15:26:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (31,'zhangfan','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-29 16:16:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (32,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-30 00:06:13');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (33,'admin','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-30 00:53:23');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (34,'lisi','192.168.3.40','局域网','Chrome/137.0.0.0','Windows 10',1,'登录成功','2025-06-30 00:53:59');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (35,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-06-30 10:28:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (36,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-01 10:30:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (37,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-02 11:58:08');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (38,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 11:59:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (39,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 14:12:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (40,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 14:27:21');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (41,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 16:53:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (42,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-07-03 17:26:51');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (43,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-07-03 17:26:57');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (44,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-07-03 17:27:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (45,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-07-03 17:27:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (46,'zhangsan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-07-03 17:27:18');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (47,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 17:27:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (48,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 17:28:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (49,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 17:29:34');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (50,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 17:47:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (51,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 17:52:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (52,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 17:54:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (53,'lisi','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 18:54:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (54,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-03 18:55:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (55,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-04 13:01:03');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (56,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-04 13:08:12');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (57,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-04 16:46:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (58,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-04 17:14:45');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (59,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-04 17:54:44');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (60,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-04 17:55:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (61,'zhangfan','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-04 17:55:39');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (62,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-07 15:30:48');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (63,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-08 17:57:42');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (64,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-09 11:38:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (65,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-10 10:39:06');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (66,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-10 13:07:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (67,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-10 14:22:56');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (68,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-07-10 14:52:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (69,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-07-10 14:52:46');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (70,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-10 14:52:52');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (71,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-14 11:40:24');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (72,'admin','127.0.0.1','服务器登录','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-15 11:47:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (73,'admin','127.0.0.1','服务器登录','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-15 14:24:52');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (74,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-16 14:36:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (75,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-16 15:51:44');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (76,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-17 10:55:44');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (78,'admin','10.7.16.22','局域网','Chrome/137.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-17 20:45:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (79,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-18 14:46:12');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (80,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-18 14:49:31');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (81,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-18 15:07:07');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (82,'zhangsan','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-18 15:32:43');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (83,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-18 15:33:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (84,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-18 16:59:52');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (85,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-22 10:20:45');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (86,'zhangfan','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-07-22 10:29:25');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (87,'zhangfan','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-07-22 10:29:36');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (88,'zhangsan','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-22 10:29:45');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (89,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-22 10:30:24');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (90,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-22 14:22:55');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (91,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-22 16:09:14');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (92,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-07-22 22:45:37');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (93,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-23 10:26:42');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (94,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-24 10:28:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (95,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-07-24 14:39:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (96,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-07-24 14:40:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (97,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-07-24 14:42:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (98,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-24 15:21:36');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (99,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-24 19:42:47');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (130,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 12:41:48');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (131,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 12:45:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (132,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 12:46:34');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (133,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 13:06:53');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (134,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 13:20:24');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (135,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 13:21:43');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (136,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-07-31 19:45:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (137,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 19:45:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (138,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-07-31 21:19:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (139,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-01 21:53:34');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (140,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-02 00:16:46');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (141,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-02 16:41:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (142,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-03 16:51:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (143,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-04 10:13:37');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (144,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-04 19:35:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (145,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-08-04 22:34:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (146,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-05 09:59:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (147,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-05 14:40:59');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (148,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-05 14:41:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (149,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-05 14:41:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (150,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-08-05 23:29:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (151,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-06 13:37:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (152,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-06 13:39:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (153,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-06 14:37:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (154,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-06 16:59:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (155,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-08-06 22:34:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (156,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-08-07 14:54:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (157,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-07 14:54:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (158,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-08-07 22:48:24');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (159,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-08-07 23:13:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (160,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-08 15:04:01');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (161,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-08 23:48:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (162,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-09 12:59:01');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (163,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-10 13:13:50');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (164,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-11 11:14:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (165,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-12 11:16:21');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (166,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-13 10:33:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (167,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-14 10:36:39');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (168,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-14 21:41:10');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (169,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-15 10:57:16');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (170,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-17 18:12:31');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (171,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-21 11:33:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (172,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-21 11:39:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (173,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-22 10:41:13');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (174,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-23 12:35:49');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (175,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-23 13:47:27');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (176,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-23 18:36:52');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (177,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-23 22:43:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (178,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-23 22:43:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (179,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-24 15:27:16');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (180,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-25 10:15:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (182,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-08-25 15:19:48');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (183,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-25 15:19:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (184,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-26 19:59:49');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (185,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-28 10:17:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (186,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-29 10:28:08');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (187,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-29 11:11:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (188,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-08-29 11:30:08');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (189,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-09-04 20:10:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (190,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-04 20:10:06');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (191,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',2,'密码不正确','2025-09-04 23:21:57');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (192,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-09-04 23:22:56');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (193,'admin','127.0.0.1','服务器登录','Chrome/139.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-05 20:11:51');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (194,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-06 14:59:37');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (195,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-06 15:30:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (196,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-07 16:02:46');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (197,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-08 10:26:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (198,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-08 13:28:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (199,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-09 15:05:23');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (200,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-10 11:10:32');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (201,'admin','127.0.0.1','服务器登录','Chrome/138.0.0.0','Windows 10',1,'登录成功','2025-09-10 22:47:46');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (202,'zhangfan','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-09-11 10:43:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (203,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-11 10:44:08');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (204,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-11 12:29:57');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (205,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-11 15:06:41');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (206,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-11 15:09:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (207,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-11 20:11:48');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (208,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-11 20:49:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (209,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-11 21:03:03');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (210,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-12 10:23:06');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (211,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-12 10:33:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (212,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-12 11:20:47');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (213,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-09-12 17:52:39');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (214,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-12 17:52:46');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (215,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-12 23:13:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (216,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-13 15:21:53');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (217,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-14 15:26:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (218,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-14 23:19:13');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (219,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-15 10:47:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (220,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-15 15:19:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (221,'root','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-09-15 20:46:53');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (222,'root','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-09-15 20:47:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (223,'root','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-09-15 20:47:01');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (224,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-15 20:47:10');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (225,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-16 09:54:18');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (226,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-16 09:55:10');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (227,'root','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-09-16 17:13:49');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (228,'root','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-09-16 17:14:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (229,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-16 17:14:06');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (230,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-16 17:20:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (231,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-16 17:20:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (232,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-17 10:08:32');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (233,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:13:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (234,'zhangsan','127.0.0.1','服务器登录','Safari/605.1.15','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:14:56');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (235,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:23:53');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (237,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:30:10');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (238,'zhangsan','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:31:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (241,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:33:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (242,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:36:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (243,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 19:56:06');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (244,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-17 20:15:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (246,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 20:24:43');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (247,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 20:40:16');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (248,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 20:41:45');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (249,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-17 20:49:04');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (250,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-17 20:52:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (251,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-09-17 20:58:27');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (252,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-17 20:58:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (253,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 21:01:53');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (254,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-17 21:11:34');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (255,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-18 11:03:24');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (256,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-18 11:40:50');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (257,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-18 15:16:51');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (258,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-18 16:03:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (259,'zhangsan','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-18 18:30:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (260,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-19 21:29:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (261,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-19 22:46:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (262,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-19 22:50:23');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (263,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-20 02:01:52');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (264,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-21 02:16:14');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (265,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-21 15:13:08');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (266,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-21 15:59:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (267,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-22 11:44:03');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (268,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-23 11:59:48');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (269,'admin','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-23 15:22:41');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (270,'admin','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-24 12:50:17');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (271,'test','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-24 12:50:34');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (272,'test','127.0.0.1','服务器登录','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-09-24 12:51:31');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (273,'test','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-24 12:55:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (274,'test','127.0.0.1','服务器登录','Chrome/140.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-09-24 13:23:57');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (305,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-16 19:39:43');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (306,'admin','127.0.0.1','服务器登录','Chrome/141.0.0.0','Windows 10',2,'密码不正确','2025-10-16 22:12:52');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (307,'admin','127.0.0.1','服务器登录','Chrome/141.0.0.0','Windows 10',1,'登录成功','2025-10-16 22:14:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (308,'admin','192.168.3.7','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-18 15:58:49');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (309,'admin','192.168.3.7','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-18 17:53:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (310,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-21 17:04:41');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (311,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-21 17:10:32');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (312,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-10-22 17:29:11');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (313,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-22 17:29:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (314,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-23 10:24:16');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (315,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-23 15:41:50');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (316,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-24 20:05:14');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (317,'admin','192.168.3.7','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-25 15:49:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (318,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-28 17:14:31');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (319,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-29 17:16:37');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (320,'admin','10.7.16.22','局域网','Chrome/141.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-10-30 17:17:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (321,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-04 15:41:31');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (322,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-04 20:29:57');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (323,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-05 20:31:59');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (324,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-06 20:36:06');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (325,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-08 18:34:38');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (326,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-10 11:31:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (327,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-11 11:37:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (328,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-11 11:38:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (329,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-11 21:09:08');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (330,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-12 11:04:00');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (331,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-12 11:08:18');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (332,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-12 11:41:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (333,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-13 11:46:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (334,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-13 14:40:25');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (335,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-13 20:22:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (336,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-14 15:43:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (337,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-14 15:50:13');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (338,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-14 19:23:27');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (339,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-14 19:50:36');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (340,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-14 19:50:38');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (341,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-17 14:54:25');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (342,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-19 11:34:27');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (343,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-20 11:34:38');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (344,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-20 11:49:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (345,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-20 11:50:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (346,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-11-20 14:07:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (347,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-20 14:07:33');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (348,'admin','10.7.16.22','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-21 14:30:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (349,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-21 22:17:46');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (350,'admin','127.0.0.1','服务器登录','Chrome/135.0.0.0','Windows 10',2,'验证码已过期','2025-11-22 15:29:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (351,'admin','127.0.0.1','服务器登录','Chrome/135.0.0.0','Windows 10',1,'登录成功','2025-11-22 15:29:18');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (352,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-22 22:39:56');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (353,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-23 12:46:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (354,'lisi','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-11-23 12:47:08');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (355,'lisi','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-11-23 12:47:32');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (356,'lisi','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'密码不正确','2025-11-23 12:47:42');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (357,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-23 12:47:51');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (358,'lisi','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-23 12:48:32');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (359,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-11-24 00:13:30');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (360,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-24 00:14:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (361,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-24 00:14:56');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (362,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-24 00:34:12');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (363,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-24 00:43:33');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (364,'wangwu','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-24 00:44:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (365,'lisi','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-24 00:45:23');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (366,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-24 10:33:37');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (367,'admin','127.0.0.1','服务器登录','Chrome/135.0.0.0','Windows 10',1,'登录成功','2025-11-24 22:10:23');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (368,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-11-25 11:15:31');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (369,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-25 11:15:42');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (370,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-26 16:26:48');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (371,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-26 16:34:42');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (372,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-27 11:22:54');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (373,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-11-27 16:41:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (374,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-11-27 16:41:31');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (375,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-27 16:41:37');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (376,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-28 11:14:21');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (377,'lisi','10.7.16.60','局域网','Chrome/98.0.4758.139','Mac OS X 10_15_7',1,'登录成功','2025-11-28 14:18:34');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (378,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-28 16:27:03');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (379,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-28 22:57:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (380,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-29 19:00:27');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (381,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-30 11:37:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (382,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-30 13:27:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (383,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-30 23:22:27');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (384,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-30 23:23:01');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (385,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-11-30 23:47:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (386,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-01 00:41:42');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (387,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-01 00:57:22');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (388,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-01 10:18:24');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (389,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-01 15:29:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (390,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-01 17:20:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (391,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-01 20:02:35');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (392,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-01 21:06:17');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (393,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-02 20:20:46');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (394,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-03 10:24:52');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (395,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-03 19:50:02');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (396,'admin','10.7.16.60','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-03 21:22:45');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (397,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-12-04 22:53:58');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (398,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-04 22:54:05');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (399,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 00:24:40');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (400,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 00:32:45');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (401,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 00:45:07');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (402,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 00:53:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (403,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 01:03:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (404,'admin','10.7.16.100','局域网','Chrome/143.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 12:55:09');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (405,'admin','10.7.16.100','局域网','Chrome/143.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 13:26:29');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (406,'admin','10.7.16.100','局域网','Chrome/143.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 17:58:41');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (407,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 20:15:33');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (408,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码已过期','2025-12-05 23:10:15');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (409,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-05 23:10:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (410,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 00:42:41');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (411,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 00:46:43');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (412,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 00:55:19');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (413,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 01:02:48');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (414,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',2,'验证码不正确','2025-12-06 01:40:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (415,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 01:40:26');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (416,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 01:55:23');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (417,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 11:39:20');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (418,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 21:10:28');
INSERT INTO `sys_login_info` (`id`,`username`,`ip_address`,`login_location`,`browser`,`os`,`login_status`,`message`,`login_time`) VALUES (419,'admin','192.168.3.7','局域网','Chrome/142.0.0.0','Mac OS X 10_15_7',1,'登录成功','2025-12-06 21:52:49');

-- Table structure for sys_menu
CREATE TABLE `sys_menu` (
 `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键',
 `parent_id` int(10) COMMENT '父级菜单id',
 `menu_name` varchar(100) COMMENT '菜单名称',
 `icon` varchar(100) COMMENT '图标',
 `value` varchar(100) COMMENT '接口权限值',
 `menu_type` int(10) COMMENT '菜单类型：1->目录；2->菜单；3->按钮（接口绑定权限）',
 `url` varchar(200) COMMENT '菜单url',
 `menu_status` int(10) DEFAULT 2 COMMENT '启用状态；1->禁用；2->启用',
 `sort` int(10) COMMENT '排序',
 `create_time` datetime COMMENT '创建时间', 
PRIMARY KEY (id)
) COMMENT '菜单表';

-- Data for table sys_menu
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (4,0,'系统管理','StarFilled','',1,'',2,7,'2022-09-04 13:57:39');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (6,4,'用户信息','Avatar','base:admin:list',2,'system/admin',2,1,'2022-09-04 13:59:39');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (7,4,'角色信息','InfoFilled','base:role:list',2,'system/role',2,2,'2022-09-04 14:00:12');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (8,4,'菜单信息','Histogram','base:menu:list',2,'system/menu',2,3,'2022-09-04 14:00:17');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (9,4,'部门信息','Menu','base:dept:list',2,'system/dept',2,4,'2022-09-04 14:01:58');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (10,4,'岗位信息','Promotion','base:post:list',2,'system/post',2,5,'2022-09-04 14:02:06');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (16,6,'新增用户','','base:admin:add',3,'',2,1,'2022-09-04 18:32:55');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (17,6,'修改用户','','base:admin:edit',3,'',2,2,'2022-09-04 18:33:29');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (18,6,'删除用户','','base:admin:delete',3,'',2,3,'2022-09-04 18:33:51');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (21,7,'新增角色','','base:role:add',3,'',2,1,'2022-09-04 18:44:47');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (22,7,'修改角色','','base:role:edit',3,'',2,2,'2022-09-04 18:45:10');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (23,7,'删除角色','','base:role:delete',3,'',2,3,'2022-09-04 18:45:46');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (24,7,'分配权限','','base:role:assign',3,'',2,4,'2022-09-04 18:46:20');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (26,8,'新增菜单','','base:menu:add',3,'',2,1,'2022-09-04 18:49:51');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (27,8,'修改菜单','','base:menu:edit',3,'',2,2,'2022-09-04 18:50:24');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (28,8,'删除菜单','','base:menu:delete',3,'',2,3,'2022-09-04 18:50:53');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (29,9,'新增部门','','base:dept:add',3,'',2,1,'2022-09-04 18:52:16');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (30,9,'修改部门','','base:dept:edit',3,'',2,2,'2022-09-04 18:52:37');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (31,9,'删除部门','','base:dept:delete',3,'',2,3,'2022-09-04 18:52:50');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (32,10,'新增岗位','','base:post:add',3,'',2,1,'2022-09-04 18:53:28');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (33,10,'修改岗位','','base:post:edit',3,'',2,2,'2022-09-04 18:53:48');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (34,10,'删除岗位','','base:post:delete',3,'',2,3,'2022-09-04 18:54:00');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (44,0,'操作审计','BellFilled','',1,'',2,9,'2022-09-05 11:06:57');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (45,44,'操作日志','User','monitor:operator:list',2,'monitor/operator',2,1,'2022-09-05 11:10:54');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (46,44,'登录日志','DocumentRemove','monitor:loginLog:list',2,'monitor/loginlog',2,2,'2022-09-05 11:11:31');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (47,45,'清空操作日志','','monitor:operator:clean',3,'',2,1,'2022-09-05 11:12:36');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (49,46,'清空登录日志','','monitor:loginLog:clean',3,'',2,1,'2022-09-05 11:16:01');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (60,6,'重置密码',NULL,'base:admin:reset',3,NULL,2,6,'2022-12-01 16:33:34');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (62,46,'删除登录日志','','monitor:loginLog:delete',3,'',2,2,'2022-12-02 15:41:56');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (72,0,'仪表盘','HomeFilled','',1,'dashboard',2,1,'2023-05-24 22:11:13');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (73,45,'删除操作日志','','monitor:operator:delete',3,'',2,3,'2023-06-02 10:09:38');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (78,80,'主机管理','Platform','cmdb:ecs:list',2,'cmdb/ecs',2,1,'2025-06-29 00:30:35');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (80,0,'资产管理','TrendCharts','',1,'',2,2,'2025-07-03 11:47:07');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (81,0,'容器管理','UploadFilled','',1,'',2,3,'2025-07-03 11:50:47');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (82,81,'集群管理','Menu','cloud:k8s:list',2,'k8s/list',2,1,'2025-07-03 11:56:44');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (83,81,'节点管理','Help','cloud:k8s:node',2,'k8s/node',2,2,'2025-07-03 12:04:59');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (84,0,'配置中心','Tools','',1,'',2,8,'2025-07-04 17:00:01');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (85,84,'主机凭据','Setting','config:ecs:key',2,'config/ecskey',2,1,'2025-07-04 17:03:10');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (86,84,'通用凭据','User','config:accountauth:key',2,'config/accountauth',2,2,'2025-07-04 17:08:20');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (88,80,'业务分组','Shop','cmdb:group',2,'cmdb/group',2,2,'2025-07-16 15:17:14');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (89,88,'创建分组','','cmdb:group:add',3,'',2,1,'2025-07-18 15:24:31');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (90,88,'修改分组','','cmdb:group:update',3,'',2,2,'2025-07-18 15:25:49');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (91,88,'删除分组','','cmdb:group:delete',3,'',2,3,'2025-07-18 15:26:21');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (93,81,'工作负载','Star','cloud:k8s:workload',2,'k8s/workload',2,4,'2025-07-24 14:38:31');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (95,80,'DB数据管理','Coin','cmdb:db',2,'cmdb/db',2,3,'2025-07-29 15:27:50');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (96,44,'数据日志','Coin','monitor:dblog:list',2,'monitor/dblog',2,3,'2025-07-31 12:44:07');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (97,0,'任务中心','User','',1,'',2,5,'2025-08-06 13:33:47');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (98,97,'任务模版','connection','task:template',2,'task/template',2,2,'2025-08-06 13:35:19');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (99,97,'任务作业','key','task:job',2,'task/job',2,1,'2025-08-06 13:36:06');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (100,97,'Ansible任务','Eleme','task:ansible',2,'task/ansible',2,3,'2025-08-23 18:35:24');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (101,0,'运维工具','Search','',1,'',2,6,'2025-08-29 10:59:35');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (102,101,'agent列表','price-tag','ops:agent',2,'ops/agent',2,2,'2025-08-29 11:22:20');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (103,101,'工具列表','present','ops:tools',2,'ops/tools',2,1,'2025-08-29 11:29:02');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (104,84,'密钥管理','Phone','config:keymanage:key',2,'config/keymanage',2,3,'2025-09-08 13:24:40');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (105,81,'命名空间','discount','k8s:namespace',2,'k8s/namespace',2,3,'2025-09-11 15:02:14');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (106,81,'网络管理','guide','k8s:network',2,'k8s/network',2,5,'2025-09-11 15:04:14');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (107,81,'配置管理','connection','k8s:config',2,'k8s/config',2,7,'2025-09-11 15:04:52');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (108,81,'存储管理','Coin','k8s:storage',2,'k8s/storage',2,6,'2025-09-11 15:05:40');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (109,0,'服务管理','ElemeFilled','',1,'',2,4,'2025-09-16 09:49:55');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (110,109,'应用列表','Menu','app:application',2,'app/application',2,1,'2025-09-16 09:52:58');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (111,109,'快速发布','TrendCharts','app:quick-release',2,'app/quick-release',2,2,'2025-09-16 17:12:11');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (113,45,'批量删除','','monitor:operator:delete',3,'',2,2,'2025-09-17 20:55:13');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (114,104,'修改密钥','','config:keymanage:edit',3,'',2,1,'2025-09-18 10:45:57');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (115,104,'删除密钥','','config:keymanage:delete',3,'',2,2,'2025-09-18 10:53:44');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (117,104,'同步主机','','config:keymanage:rsync',3,'',2,3,'2025-09-18 10:57:25');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (118,104,'创建密钥','','config:keymanage:create',3,'',2,4,'2025-09-18 11:01:12');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (119,86,'修改账号','','config:common:edit',3,'',2,2,'2025-09-18 11:47:33');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (120,86,'删除账号','','config:common:delete',3,'',2,3,'2025-09-18 11:48:17');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (121,86,'解密账号','','config:common:decrypt',3,'',2,4,'2025-09-18 11:48:57');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (122,86,'创建账号','','config:common:add',3,'',2,1,'2025-09-18 11:49:30');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (123,85,'修改凭据','','config:ecs:edit',3,'',2,1,'2025-09-18 11:54:16');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (124,85,'删除凭据','','config:ecs:delete',3,'',2,2,'2025-09-18 11:54:51');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (125,85,'创建凭据','','config:ecs:create',3,'',2,3,'2025-09-18 11:55:21');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (126,102,'卸载agent','','ops:agent:delete',3,'',2,1,'2025-09-18 12:47:48');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (127,102,'查看agent','','ops:agent:get',3,'',2,2,'2025-09-18 12:49:02');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (128,102,'部署agent','','ops:agent:create',3,'',2,3,'2025-09-18 12:49:56');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (129,102,'批量卸载agent','','ops:agent:deleteall',3,'',2,4,'2025-09-18 12:50:52');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (130,100,'启动ansible任务流程','','task:ansible:start',3,'',2,1,'2025-09-18 12:59:30');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (131,100,'删除ansible任务','','task:ansible:delete',3,'',2,2,'2025-09-18 13:00:03');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (132,100,'新增ansible任务','','task:ansible:create',3,'',2,3,'2025-09-18 13:00:45');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (133,98,'新增模版','','task:template:add',3,'',2,1,'2025-09-18 13:16:38');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (134,98,'修改模版','','task:template:edit',3,'',2,2,'2025-09-18 13:17:04');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (135,98,'删除模版','','task:template:delete',3,'',2,3,'2025-09-18 13:18:25');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (136,99,'新增任务','','task:job:add',3,'',2,1,'2025-09-18 13:24:19');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (137,99,'启动任务','','task:job:start',3,'',2,2,'2025-09-18 13:24:59');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (138,99,'删除任务','','task:job:delete',3,'',2,3,'2025-09-18 13:25:41');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (139,111,'新建发布','','app:quick-release:add',3,'',2,1,'2025-09-18 13:30:53');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (140,111,'启动发布','','app:quick-release:start',3,'',2,2,'2025-09-18 13:32:11');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (141,111,'删除发布','','app:quick-release:delete',3,'',2,3,'2025-09-18 13:32:32');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (142,110,'创建应用','','app:application:add',3,'',2,1,'2025-09-18 14:28:07');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (143,110,'修改应用','','app:application:edit',3,'',2,2,'2025-09-18 14:28:59');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (144,110,'配置应用环境','','app:application:env',3,'',2,3,'2025-09-18 14:29:34');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (145,110,'删除应用','','app:application:delete',3,'',2,4,'2025-09-18 14:30:11');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (146,95,'创建数据库账号','','cmdb:db:add',3,'',2,1,'2025-09-18 14:41:32');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (147,95,'修改数据库配置','','cmdb:db:edit',3,'',2,2,'2025-09-18 14:42:47');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (148,95,'删除数据库账号','','cmdb:db:delete',3,'',2,3,'2025-09-18 14:43:57');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (149,78,'创建主机','','cmdb:ecs:add',3,'',2,1,'2025-09-18 14:47:42');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (150,78,'主机终端','','cmdb:ecs:terminal',3,'',2,2,'2025-09-18 14:48:36');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (151,78,'修改主机信息','','cmdb:ecs:edit',3,'',2,3,'2025-09-18 14:49:43');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (152,78,'上传文件到主机','','cmdb:ecs:upload',3,'',2,4,'2025-09-18 14:50:38');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (153,78,'执行主机命令','','cmdb:ecs:shell',3,'',2,5,'2025-09-18 14:51:10');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (154,78,'监控主机','','cmdb:ecs:monitor',3,'',2,6,'2025-09-18 14:51:52');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (155,78,'删除主机','','cmdb:ecs:delete',3,'',2,7,'2025-09-18 14:52:20');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (156,99,'启动脚本','','task:job:jobstart',3,'',2,4,'2025-09-18 18:36:38');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (157,99,'停止脚本','','task:job:jobstop',3,'',2,5,'2025-09-18 18:39:23');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (159,100,'开始ansible任务作业','','task:ansible:jobstart',3,'',2,4,'2025-09-18 18:43:40');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (160,111,'启动jenkins任务','','app:quick-release:jobstart',3,'',2,4,'2025-09-18 18:47:39');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (161,111,'停止jenkins任务','','app:quick-release:jobstop',3,'',2,5,'2025-09-18 18:48:16');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (162,110,'新增环境','','app:application:envadd',3,'',2,5,'2025-09-18 21:02:28');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (163,110,'修改环境','','app:application:envedit',3,'',2,6,'2025-09-18 21:03:08');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (164,110,'删除环境','','app:application:envdelete',3,'',2,7,'2025-09-18 21:04:43');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (165,78,'连接主机终端','','cmdb:ecs:connecthost',3,'',2,8,'2025-09-18 21:11:43');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (166,78,'同步主机信息','','cmdb:ecs:rsync',3,'',2,9,'2025-09-19 21:35:06');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (167,82,'注册集群','','cloud:k8s:register',3,'',2,1,'2025-09-19 21:57:54');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (168,82,'创建集群','','cloud:k8s:add',3,'',2,2,'2025-09-19 21:58:19');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (169,82,'修改集群','','cloud:k8s:edit',3,'',2,3,'2025-09-19 21:59:06');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (170,82,'同步集群','','cloud:k8s:rsync',3,'',2,4,'2025-09-19 21:59:31');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (171,82,'删除集群','','cloud:k8s:delete',3,'',2,5,'2025-09-19 21:59:56');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (172,83,'查看监控仪表盘','','k8s:node:monitor',3,'',2,1,'2025-09-20 00:19:49');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (173,83,'查看节点资源详情','','k8s:node:details',3,'',2,2,'2025-09-20 00:21:20');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (174,83,'配置节点污点','','k8s:node:stain',3,'',2,3,'2025-09-20 00:22:17');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (175,83,'增加标签','','k8s:node:label',3,'',2,4,'2025-09-20 00:23:15');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (176,83,'封锁节点','','k8s:node:close',3,'',2,5,'2025-09-20 00:24:13');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (177,83,'驱逐节点','','k8s:node:expel',3,'',2,6,'2025-09-20 00:25:04');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (178,105,'创建命名空间','','k8s:namespace:add',3,'',2,1,'2025-09-20 00:36:14');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (179,105,'查看命名空间详情','','k8s:namespace:details',3,'',2,2,'2025-09-20 00:37:22');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (180,105,'查看命名空间资源配置','','k8s:namespace:setup',3,'',2,3,'2025-09-20 00:39:05');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (181,105,'添加命名空间资源配置','','k8s:namespace:setupadd',3,'',2,4,'2025-09-20 00:40:06');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (182,105,'查看限制命名空间资源','','k8s:namespace:restriction',3,'',2,5,'2025-09-20 00:41:25');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (183,105,'添加限制命名空间资源','','k8s:namespace:restrictionadd',3,'',2,6,'2025-09-20 00:42:23');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (184,105,'删除命名空间','','k8s:namespace:delete',3,'',2,7,'2025-09-20 00:43:03');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (185,93,'新增工作负载','','k8s:workload:add',3,'',2,1,'2025-09-20 01:05:08');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (186,93,'伸缩pod','','k8s:workload:expandable',3,'',2,2,'2025-09-20 01:06:18');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (187,93,'重启pod','','k8s:workload:restart',3,'',2,3,'2025-09-20 01:07:13');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (188,93,'更新pod资源限制','','k8s:workload:resource',3,'',2,4,'2025-09-20 01:08:39');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (189,93,'更新pod资调度','','k8s:workload:dispatch',3,'',2,5,'2025-09-20 01:09:36');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (190,93,'更新yaml文件','','k8s:workload:edityaml',3,'',2,6,'2025-09-20 01:10:57');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (191,93,'删除工作负载','','k8s:workload:delete',3,'',2,7,'2025-09-20 01:11:52');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (192,93,'回滚工作负载版本','','k8s:workload:rollback_version',3,'',2,8,'2025-09-20 01:39:38');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (193,93,'查看pod日志','','k8s:workload:podlog',3,'',2,9,'2025-09-20 01:44:37');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (194,93,'删除pod','','k8s:workload:poddelete',3,'',2,10,'2025-09-20 01:45:16');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (195,93,'登陆pod终端','','k8s:workload:terminal',3,'',2,11,'2025-09-20 01:46:07');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (196,93,'编辑pod yaml文件','','k8s:workload:edityaml',3,'',2,12,'2025-09-20 01:47:22');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (197,106,'新增service','','k8s:network:addservice',3,'',2,1,'2025-09-20 02:14:21');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (198,106,'编辑 Service','','k8s:network:editservice',3,'',2,2,'2025-09-20 02:15:19');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (199,106,'编辑service YAML','','k8s:network:edit_service_yaml',3,'',2,3,'2025-09-20 02:16:15');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (200,106,'查看Service 事件','','k8s:network:service_event',3,'',2,4,'2025-09-20 02:18:11');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (201,106,'删除Service','','k8s:network:deleteservice',3,'',2,5,'2025-09-20 02:18:59');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (202,106,'编辑ingress','','k8s:network:editingress',3,'',2,6,'2025-09-20 02:26:59');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (203,106,'新增ingress','','k8s:network:addingress',3,'',2,7,'2025-09-20 02:27:29');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (204,106,'编辑ingress_yaml','','k8s:network:edit_ingress_yaml',3,'',2,8,'2025-09-20 02:28:23');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (205,106,'查看ingress 事件','','k8s:network:ingress_event',3,'',2,9,'2025-09-20 02:29:24');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (206,106,'删除ingress','','k8s:network:delete_ingress',3,'',2,10,'2025-09-20 02:30:04');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (207,109,'服务上线','Platform','app:upgrade',2,'app/upgrade',2,3,'2025-11-20 11:46:53');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (208,109,'脚本上线','Phone','app:sh_release',2,'app/sh-release',2,4,'2025-11-28 11:13:05');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (212,0,'监控中心','Shop','',1,'',2,4,'2025-12-03 21:21:04');
INSERT INTO `sys_menu` (`id`,`parent_id`,`menu_name`,`icon`,`value`,`menu_type`,`url`,`menu_status`,`sort`,`create_time`) VALUES (213,212,'域名监控','UploadFilled','monitor:domain',2,'monitor/domain',2,1,'2025-12-03 21:22:11');

-- Table structure for sys_operation_log
CREATE TABLE `sys_operation_log` (
 `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键',
 `admin_id` bigint(20) NOT NULL COMMENT '''管理员id''',
 `username` varchar(64) NOT NULL COMMENT '''管理员账号''',
 `method` varchar(64) NOT NULL COMMENT '''请求方式''',
 `ip` varchar(64) COMMENT '''IP''',
 `url` varchar(500) COMMENT '''URL''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `description` varchar(255) COMMENT '''操作描述''', 
PRIMARY KEY (id)
) COMMENT '操作日志记录表';
ALTER TABLE `sys_operation_log` ADD  INDEX `create_time`(`create_time`) USING BTREE;
ALTER TABLE `sys_operation_log` ADD  INDEX `username`(`username`) USING BTREE;

-- Data for table sys_operation_log
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (1,89,'admin','delete','10.7.16.22','/api/v1/sysOperationLog/clean','2025-10-16 10:59:42','清空操作日志');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (2,89,'admin','delete','10.7.16.22','/api/v1/sysLoginInfo/clean','2025-10-16 10:59:48','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (3,89,'admin','delete','10.7.16.22','/api/v1/sysLoginInfo/clean','2025-10-16 10:59:52','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (4,89,'admin','delete','10.7.16.22','/api/v1/sysLoginInfo/clean','2025-10-16 10:59:57','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (5,89,'admin','delete','10.7.16.22','/api/v1/sysLoginInfo/clean','2025-10-16 11:00:03','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (6,89,'admin','delete','10.7.16.22','/api/v1/sysLoginInfo/batch/delete','2025-10-16 11:00:08','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (7,89,'admin','delete','10.7.16.22','/api/v1/sysLoginInfo/batch/delete','2025-10-16 11:00:14','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (9,89,'admin','delete','192.168.3.7','/api/v1/monitor/agent/delete/75','2025-10-18 16:12:37','删除Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (10,89,'admin','post','192.168.3.7','/api/v1/config/ecsauthadd','2025-10-18 16:13:33','新增ECS认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (11,89,'admin','post','192.168.3.7','/api/v1/cmdb/hostcreate','2025-10-18 16:14:25','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (12,89,'admin','post','192.168.3.7','/api/v1/cmdb/hostsync','2025-10-18 16:14:47','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (13,89,'admin','post','192.168.3.7','/api/v1/monitor/agent/deploy','2025-10-18 16:17:39','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (14,89,'admin','post','192.168.3.7','/api/v1/cmdb/hostsync','2025-10-18 17:55:31','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (15,89,'admin','put','192.168.3.7','/api/v1/cmdb/hostupdate','2025-10-18 17:57:43','修改主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (16,89,'admin','post','192.168.3.7','/api/v1/cmdb/hostssh/upload/501','2025-10-18 17:58:06','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (17,89,'admin','post','192.168.3.7','/api/v1/cmdb/sql/databaselist','2025-10-18 18:00:26','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (18,89,'admin','post','192.168.3.7','/api/v1/config/accountauth/decrypt','2025-10-18 18:05:10','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (19,89,'admin','delete','10.7.16.22','/api/v1/monitor/agent/uninstall','2025-10-21 17:12:00','卸载Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (20,89,'admin','put','10.7.16.22','/api/v1/cmdb/hostupdate','2025-10-21 17:19:45','修改主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (21,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostsync','2025-10-21 17:20:47','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (22,89,'admin','post','10.7.16.22','/api/v1/monitor/agent/deploy','2025-10-21 17:28:02','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (23,89,'admin','delete','10.7.16.22','/api/v1/monitor/agent/uninstall','2025-10-21 17:34:17','卸载Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (24,89,'admin','post','10.7.16.22','/api/v1/monitor/agent/deploy','2025-10-21 17:34:36','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (25,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcreate','2025-10-21 17:39:30','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (26,89,'admin','post','10.7.16.22','/api/v1/monitor/agent/deploy','2025-10-21 19:47:07','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (27,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 10:36:04','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (28,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:47:37','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (29,89,'admin','post','10.7.16.22','/api/v1/tool','2025-10-23 15:48:09','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (30,89,'admin','delete','10.7.16.22','/api/v1/tool/1','2025-10-23 15:50:23','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (31,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:51:14','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (32,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:53:49','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (33,89,'admin','post','10.7.16.22','/api/v1/tool','2025-10-23 15:54:08','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (34,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:55:56','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (35,89,'admin','post','10.7.16.22','/api/v1/tool','2025-10-23 15:55:59','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (36,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:56:12','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (37,89,'admin','post','10.7.16.22','/api/v1/tool','2025-10-23 15:56:22','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (38,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:56:37','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (39,89,'admin','post','10.7.16.22','/api/v1/tool','2025-10-23 15:56:48','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (40,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:57:07','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (41,89,'admin','post','10.7.16.22','/api/v1/tool','2025-10-23 15:57:21','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (42,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-23 15:57:40','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (43,89,'admin','post','10.7.16.22','/api/v1/tool','2025-10-23 15:57:55','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (44,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostsync','2025-10-23 19:34:27','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (45,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 11:32:51','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (46,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 11:47:35','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (47,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 14:26:23','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (48,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/3','2025-10-30 15:13:50','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (49,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/1','2025-10-30 15:13:57','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (50,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 16:51:25','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (51,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/4','2025-10-30 17:03:49','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (52,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/2','2025-10-30 17:03:53','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (53,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 17:05:27','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (54,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/5','2025-10-30 17:18:17','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (55,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 17:18:51','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (56,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/6','2025-10-30 17:19:39','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (57,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 17:20:01','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (58,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-30 19:45:17','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (59,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-31 11:55:07','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (60,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/9','2025-10-31 12:42:51','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (61,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/8','2025-10-31 12:42:54','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (62,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-31 12:43:12','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (63,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-31 12:51:42','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (64,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-31 12:51:56','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (65,89,'admin','delete','10.7.16.22','/api/v1/tool/deploy/10','2025-10-31 12:52:43','部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (66,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-31 12:53:21','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (67,89,'admin','post','10.7.16.22','/api/v1/upload','2025-10-31 15:36:29','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (68,89,'admin','post','10.7.16.22','/api/v1/tool/deploy','2025-10-31 15:40:42','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (69,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster','2025-11-04 15:42:04','创建K8s集群');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (70,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/rental-common/pods/yaml','2025-11-04 17:22:05','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (71,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/rental-common/pods/yaml','2025-11-04 17:23:07','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (72,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/rental-common/pods/yaml','2025-11-04 17:29:16','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (73,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/rental-common/ingresses/test2-rental-ingress-rental-xxljob/yaml','2025-11-04 17:38:03','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (74,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/rental-common/ingresses/test2-rental-ingress-rental-xxljob/yaml','2025-11-04 17:50:53','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (75,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/secrets','2025-11-05 19:56:12','创建K8s密钥');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (76,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/configmaps','2025-11-05 19:56:52','创建K8s配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (77,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/configmaps','2025-11-05 19:58:04','创建K8s配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (78,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/secrets','2025-11-05 19:59:26','创建K8s密钥');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (79,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/configmaps/istio-ca-root-cert/yaml','2025-11-05 20:01:32','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (80,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/secrets/dding-net2025/yaml','2025-11-05 20:02:36','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (81,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/secrets/dding-net2025/yaml','2025-11-05 20:06:01','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (82,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/secrets/dding-net2025/yaml','2025-11-05 20:08:44','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (83,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/secrets/dding-net2025/yaml','2025-11-05 20:09:41','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (84,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/default/secrets/dding-net2025/yaml','2025-11-05 20:15:33','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (85,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/dev-iot-proxy/deployments/deploy-iot-proxy-factory/scale','2025-11-06 14:25:26','伸缩K8s部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (86,89,'admin','post','10.7.16.22','/api/v1/task/add','2025-11-06 15:50:35','新增任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (87,89,'admin','post','10.7.16.22','/api/v1/taskjob/start','2025-11-06 15:50:49','启动任务作业');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (88,89,'admin','post','10.7.16.22','/api/v1/task/add','2025-11-06 15:51:56','新增任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (89,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-06 16:01:52','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (90,89,'admin','post','10.7.16.22','/api/v1/config/accountauth/decrypt','2025-11-06 16:10:28','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (91,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-06 16:11:27','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (92,89,'admin','post','10.7.16.22','/api/v1/config/accountauth/decrypt','2025-11-06 16:12:23','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (93,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcloudcreatealiyun','2025-11-06 16:56:27','创建阿里云主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (94,89,'admin','post','10.7.16.22','/api/v1/config/accountauth','2025-11-06 20:37:34','新增账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (95,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-06 20:39:21','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (96,89,'admin','put','10.7.16.22','/api/v1/config/accountauth','2025-11-06 20:40:11','修改账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (97,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-06 20:40:34','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (98,89,'admin','delete','10.7.16.22','/api/v1/cmdb/hostdelete','2025-11-07 14:33:02','删除主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (99,89,'admin','delete','10.7.16.22','/api/v1/cmdb/hostdelete','2025-11-07 14:33:05','删除主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (100,89,'admin','delete','10.7.16.22','/api/v1/cmdb/hostdelete','2025-11-07 14:33:11','删除主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (101,89,'admin','delete','10.7.16.22','/api/v1/cmdb/hostdelete','2025-11-07 14:33:30','删除主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (102,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-lkadmin/pods/yaml','2025-11-07 14:41:13','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (103,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-lkadmin/pods/yaml','2025-11-07 14:43:50','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (104,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-lkadmin/ingresses/ingress-test-lkadmin-lockin-manager-boot','2025-11-07 15:24:16','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (105,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-lkadmin/pods/yaml','2025-11-07 15:29:00','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (106,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-lkadmin/pods/yaml','2025-11-07 15:29:39','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (107,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-lkadmin/ingresses/test-lkadmin-lockin-management-boot/yaml','2025-11-07 16:27:32','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (108,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:19:11','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (109,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:19:44','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (110,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:24:52','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (111,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:25:02','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (112,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:30:50','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (113,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:32:17','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (114,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:32:53','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (115,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:33:18','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (116,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:34:15','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (117,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:34:59','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (118,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:39:07','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (119,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test2-lkadmin/ingresses/test-backend','2025-11-07 17:41:17','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (120,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-saas5/ingresses/test-backend','2025-11-07 17:45:59','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (121,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-saas5/ingresses/test-backend','2025-11-07 17:50:14','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (122,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/32/namespaces/test-saas5/ingresses/test-backend','2025-11-07 17:51:10','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (123,89,'admin','post','10.7.16.22','/api/v1/config/accountauth','2025-11-07 20:16:23','新增账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (124,89,'admin','delete','10.7.16.22','/api/v1/config/accountauth','2025-11-07 20:16:27','删除账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (125,89,'admin','post','10.7.16.22','/api/v1/apps','2025-11-07 20:17:54','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (126,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-07 20:24:09','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (127,89,'admin','put','10.7.16.22','/api/v1/config/accountauth','2025-11-07 20:27:57','修改账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (128,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-07 20:28:39','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (129,89,'admin','post','192.168.3.7','/api/v1/config/accountauth/decrypt','2025-11-08 18:49:42','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (130,89,'admin','post','192.168.3.7','/api/v1/config/accountauth/decrypt','2025-11-08 18:49:46','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (131,89,'admin','delete','10.7.16.22','/api/v1/monitor/agent/uninstall','2025-11-10 11:32:27','卸载Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (132,89,'admin','post','10.7.16.22','/api/v1/monitor/agent/deploy','2025-11-10 11:32:52','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (133,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 11:48:34','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (134,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 11:51:26','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (135,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 12:04:07','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (136,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 12:22:56','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (137,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 12:23:45','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (138,89,'admin','delete','10.7.16.22','/api/v1/apps/4','2025-11-10 12:36:25','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (139,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 12:40:15','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (140,89,'admin','put','10.7.16.22','/api/v1/apps/11/jenkins-envs/27','2025-11-10 12:40:21','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (141,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 12:45:46','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (142,89,'admin','put','10.7.16.22','/api/v1/apps/11/jenkins-envs/28','2025-11-10 12:45:50','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (143,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-10 14:39:07','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (144,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-10 14:39:37','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (145,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-10 14:45:56','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (146,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-10 14:50:56','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (147,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-10 15:00:11','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (148,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-10 15:04:53','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (149,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-10 15:06:06','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (150,89,'admin','post','10.7.16.22','/api/v1/apps','2025-11-10 15:08:14','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (151,89,'admin','put','10.7.16.22','/api/v1/apps/11','2025-11-10 15:08:25','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (152,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 15:08:49','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (153,89,'admin','put','10.7.16.22','/api/v1/apps/12/jenkins-envs/30','2025-11-10 15:08:55','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (154,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-10 15:11:17','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (155,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-10 15:11:44','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (156,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-10 15:41:08','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (157,89,'admin','delete','10.7.16.22','/api/v1/apps/4','2025-11-11 10:35:08','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (158,89,'admin','delete','10.7.16.22','/api/v1/apps/11','2025-11-11 10:35:43','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (159,89,'admin','delete','10.7.16.22','/api/v1/apps/4/jenkins-envs/5','2025-11-11 10:37:37','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (160,89,'admin','delete','10.7.16.22','/api/v1/apps/4/jenkins-envs/4','2025-11-11 10:37:41','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (161,89,'admin','delete','10.7.16.22','/api/v1/apps/4/jenkins-envs/6','2025-11-11 10:37:45','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (162,89,'admin','delete','10.7.16.22','/api/v1/apps/4/jenkins-envs/5','2025-11-11 10:37:54','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (163,89,'admin','delete','10.7.16.22','/api/v1/apps/4/jenkins-envs/6','2025-11-11 10:38:01','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (164,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/9','2025-11-11 10:39:22','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (165,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/8','2025-11-11 10:39:26','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (166,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/7','2025-11-11 10:40:45','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (167,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/6','2025-11-11 10:40:48','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (168,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/10','2025-11-11 10:40:51','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (169,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/11','2025-11-11 10:40:54','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (170,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/12','2025-11-11 10:40:57','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (171,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/13','2025-11-11 10:41:01','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (172,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/14','2025-11-11 10:41:04','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (173,89,'admin','delete','10.7.16.22','/api/v1/apps/4','2025-11-11 10:45:47','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (174,89,'admin','delete','10.7.16.22','/api/v1/apps/5','2025-11-11 10:45:58','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (175,89,'admin','delete','10.7.16.22','/api/v1/apps/6','2025-11-11 10:46:05','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (176,89,'admin','delete','10.7.16.22','/api/v1/apps/7','2025-11-11 10:46:20','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (177,89,'admin','delete','10.7.16.22','/api/v1/apps/8','2025-11-11 10:46:28','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (178,89,'admin','delete','10.7.16.22','/api/v1/apps/9','2025-11-11 10:46:35','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (179,89,'admin','delete','10.7.16.22','/api/v1/apps/10','2025-11-11 10:47:02','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (180,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/17','2025-11-11 10:49:12','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (181,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/18','2025-11-11 10:49:15','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (182,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/15','2025-11-11 10:52:00','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (183,89,'admin','delete','10.7.16.22','/api/v1/apps/deployment/16','2025-11-11 10:52:03','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (184,89,'admin','post','10.7.16.22','/api/v1/apps','2025-11-11 10:54:10','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (185,89,'admin','delete','10.7.16.22','/api/v1/apps/11','2025-11-11 10:54:58','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (186,89,'admin','delete','10.7.16.22','/api/v1/apps/12','2025-11-11 10:55:02','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (187,89,'admin','delete','10.7.16.22','/api/v1/config/accountauth','2025-11-11 10:57:27','删除账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (188,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-11 11:02:33','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (189,89,'admin','put','10.7.16.22','/api/v1/apps/13/jenkins-envs/33','2025-11-11 11:02:39','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (190,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-11 11:04:02','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (191,89,'admin','put','10.7.16.22','/api/v1/apps/13/jenkins-envs/34','2025-11-11 11:04:09','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (192,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-11 11:05:35','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (193,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-11 11:10:08','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (194,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-11 11:12:52','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (195,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-11 11:18:32','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (196,89,'admin','post','10.7.16.22','/api/v1/config/ecsauthadd','2025-11-11 17:24:29','新增ECS认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (197,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcreate','2025-11-11 17:24:36','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (198,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostsync','2025-11-11 17:27:25','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (199,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcloudcreatehuawei','2025-11-14 15:44:35','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (200,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcloudcreatehuawei','2025-11-14 15:51:12','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (201,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcloudcreatehuawei','2025-11-14 15:53:52','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (202,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcreate','2025-11-14 17:25:42','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (203,89,'admin','delete','10.7.16.22','/api/v1/cmdb/hostdelete','2025-11-14 17:36:08','删除主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (204,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcreate','2025-11-14 17:36:40','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (205,89,'admin','delete','10.7.16.22','/api/v1/cmdb/hostdelete','2025-11-14 17:43:36','删除主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (206,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostcreate','2025-11-14 17:43:50','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (207,89,'admin','post','10.7.16.22','/api/v1/monitor/agent/deploy','2025-11-14 17:55:56','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (208,89,'admin','put','10.7.16.22','/api/v1/cmdb/hostupdate','2025-11-14 18:02:09','修改主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (209,89,'admin','delete','10.7.16.22','/api/v1/cmdb/hostdelete','2025-11-14 18:02:14','删除主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (210,89,'admin','put','10.7.16.22','/api/v1/cmdb/hostupdate','2025-11-14 18:02:20','修改主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (211,89,'admin','put','10.7.16.22','/api/v1/cmdb/hostupdate','2025-11-14 18:02:26','修改主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (212,89,'admin','post','10.7.16.22','/api/v1/config/accountauth','2025-11-14 19:25:19','新增账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (213,89,'admin','post','10.7.16.22','/api/v1/apps','2025-11-14 19:27:43','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (214,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-14 19:28:50','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (215,89,'admin','post','10.7.16.22','/api/v1/apps/14/jenkins-envs','2025-11-14 19:29:00','添加Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (216,89,'admin','post','10.7.16.22','/api/v1/apps','2025-11-14 19:30:30','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (217,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-14 19:30:57','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (218,89,'admin','put','10.7.16.22','/api/v1/apps/15/jenkins-envs/41','2025-11-14 19:31:01','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (219,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-14 19:35:44','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (220,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-14 19:36:56','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (221,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/32','2025-11-14 19:49:22','操作K8s集群');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (222,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster','2025-11-14 19:51:40','创建K8s集群');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (223,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces','2025-11-14 19:54:51','创建K8s命名空间');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (224,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/test123','2025-11-14 19:55:10','操作K8s命名空间');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (225,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-pingtai/deployments/deploy-pingtai-lockin-utils-service/scale','2025-11-14 19:56:51','伸缩K8s部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (226,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-pingtai/deployments/deploy-pingtai-lockin-utils-service','2025-11-14 19:58:48','操作K8s部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (227,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-pingtai/pods/deploy-pingtai-lockin-utils-service-7cbcc57bf4-h65vs','2025-11-14 19:59:39','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (228,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-pingtai/pods/deploy-pingtai-lockin-utils-service-7cbcc57bf4-9pzdg','2025-11-14 19:59:49','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (229,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-pingtai/deployments/deploy-pingtai-lockin-utils-service/scale','2025-11-14 20:00:11','伸缩K8s部署');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (230,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-pingtai/pods/deploy-pingtai-lockin-utils-service-7cbcc57bf4-dpspg','2025-11-14 20:00:27','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (231,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-pingtai/pods/deploy-pingtai-lockin-utils-service-7cbcc57bf4-jzrjk','2025-11-14 20:00:38','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (232,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/lockin-saas-test/ingresses/test-backend','2025-11-14 20:06:49','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (233,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/lockin-saas-dev/ingresses/test-backend','2025-11-14 20:07:44','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (234,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/lockin-saas-dev/ingresses/test-backend','2025-11-14 20:08:15','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (235,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/rental-test/ingresses/test-backend','2025-11-14 20:08:53','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (236,89,'admin','post','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/rental-test/ingresses/test-backend','2025-11-14 20:09:32','操作K8s Ingress');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (237,89,'admin','put','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/lockin-test/workload-yaml','2025-11-17 15:01:43','更新工作负载YAML');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (238,89,'admin','delete','10.7.16.22','/api/v1/k8s/cluster/33/namespaces/dev-lkadmin/pods/deploy-lkadmin-lockin-management-776fdc5557-rg6dl','2025-11-17 15:03:56','操作K8s Pod');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (239,89,'admin','post','10.7.16.22','/api/v1/config/accountauth/decrypt','2025-11-19 13:10:49','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (240,89,'admin','delete','10.7.16.22','/api/v1/monitor/agent/uninstall','2025-11-19 15:11:25','卸载Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (241,89,'admin','post','10.7.16.22','/api/v1/monitor/agent/deploy','2025-11-19 15:11:37','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (242,89,'admin','delete','10.7.16.22','/api/v1/monitor/agent/delete/79','2025-11-19 20:28:31','删除Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (243,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostsync','2025-11-19 20:28:55','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (244,89,'admin','delete','10.7.16.22','/api/v1/monitor/agent/uninstall','2025-11-19 20:29:20','卸载Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (245,89,'admin','post','10.7.16.22','/api/v1/monitor/agent/deploy','2025-11-19 20:29:44','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (246,89,'admin','post','10.7.16.22','/api/v1/menu/add','2025-11-20 11:46:53','新增菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (247,89,'admin','put','10.7.16.22','/api/v1/menu/update','2025-11-20 11:47:52','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (248,89,'admin','put','10.7.16.22','/api/v1/role/assignPermissions','2025-11-20 11:49:52','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (249,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostsync','2025-11-20 14:27:57','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (250,89,'admin','post','10.7.16.22','/api/v1/cmdb/hostsync','2025-11-20 14:28:11','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (251,89,'admin','post','10.7.16.22','/api/v1/task/ansible','2025-11-20 14:36:27','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (252,89,'admin','post','10.7.16.22','/api/v1/task/ansible/44/start','2025-11-20 14:36:38','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (253,89,'admin','delete','10.7.16.22','/api/v1/admin/delete','2025-11-20 14:46:15','删除管理员');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (254,89,'admin','delete','10.7.16.22','/api/v1/admin/delete','2025-11-20 14:46:21','删除管理员');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (255,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 12:32:01','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (256,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 12:33:13','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (257,89,'admin','put','10.7.16.22','/api/v1/apps/15','2025-11-21 12:49:32','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (258,89,'admin','post','10.7.16.22','/api/v1/config/accountauth','2025-11-21 14:31:47','新增账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (259,89,'admin','post','10.7.16.22','/api/v1/apps','2025-11-21 15:02:25','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (260,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-21 15:02:55','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (261,89,'admin','put','10.7.16.22','/api/v1/apps/16/jenkins-envs/43','2025-11-21 15:03:34','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (262,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:04:32','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (263,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:04:43','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (264,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:19:22','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (265,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:20:03','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (266,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:26:53','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (267,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:28:03','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (268,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:35:36','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (269,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:35:51','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (270,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:42:04','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (271,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:42:11','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (272,89,'admin','post','10.7.16.22','/api/v1/apps/jenkins-job/validate','2025-11-21 15:46:43','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (273,89,'admin','put','10.7.16.22','/api/v1/apps/16/jenkins-envs/44','2025-11-21 15:46:46','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (274,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:47:24','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (275,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:47:31','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (276,89,'admin','put','10.7.16.22','/api/v1/config/accountauth','2025-11-21 15:49:02','修改账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (277,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:50:40','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (278,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:50:47','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (279,89,'admin','put','10.7.16.22','/api/v1/config/accountauth','2025-11-21 15:51:28','修改账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (280,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:52:05','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (281,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:52:11','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (282,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 15:58:45','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (283,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 15:58:59','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (284,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/quick','2025-11-21 16:04:02','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (285,89,'admin','post','10.7.16.22','/api/v1/apps/deployment/execute','2025-11-21 16:04:11','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (286,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-21 22:23:56','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (287,89,'admin','post','192.168.3.7','/api/v1/apps/jenkins-job/validate','2025-11-21 22:35:23','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (288,89,'admin','put','192.168.3.7','/api/v1/apps/16/jenkins-envs/42','2025-11-21 22:35:26','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (289,89,'admin','delete','192.168.3.7','/api/v1/config/accountauth','2025-11-21 22:37:41','删除账号认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (290,89,'admin','post','192.168.3.7','/api/v1/apps/jenkins-job/validate','2025-11-21 23:06:26','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (291,89,'admin','put','192.168.3.7','/api/v1/apps/15/jenkins-envs/39','2025-11-21 23:06:30','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (292,89,'admin','put','192.168.3.7','/api/v1/apps/14','2025-11-21 23:07:48','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (293,89,'admin','put','192.168.3.7','/api/v1/apps/14','2025-11-21 23:08:39','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (294,89,'admin','put','192.168.3.7','/api/v1/apps/15','2025-11-21 23:08:49','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (295,89,'admin','delete','192.168.3.7','/api/v1/apps/13','2025-11-21 23:09:06','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (296,89,'admin','delete','192.168.3.7','/api/v1/apps/16','2025-11-21 23:09:15','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (297,89,'admin','put','192.168.3.7','/api/v1/apps/15','2025-11-21 23:09:29','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (298,89,'admin','post','192.168.3.7','/api/v1/apps','2025-11-21 23:10:43','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (299,89,'admin','post','192.168.3.7','/api/v1/apps','2025-11-21 23:11:32','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (300,89,'admin','post','192.168.3.7','/api/v1/apps/jenkins-job/validate','2025-11-21 23:12:04','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (301,89,'admin','put','192.168.3.7','/api/v1/apps/17/jenkins-envs/45','2025-11-21 23:12:07','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (302,89,'admin','post','192.168.3.7','/api/v1/apps/jenkins-job/validate','2025-11-21 23:12:26','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (303,89,'admin','put','192.168.3.7','/api/v1/apps/18/jenkins-envs/48','2025-11-21 23:12:29','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (304,89,'admin','put','192.168.3.7','/api/v1/apps/18','2025-11-21 23:27:44','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (305,89,'admin','put','192.168.3.7','/api/v1/apps/18','2025-11-21 23:28:28','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (306,89,'admin','put','192.168.3.7','/api/v1/apps/18','2025-11-21 23:29:32','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (307,89,'admin','put','192.168.3.7','/api/v1/apps/18','2025-11-21 23:29:55','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (308,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-22 17:47:51','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (309,89,'admin','post','192.168.3.7','/api/v1/apps/jenkins-job/validate','2025-11-22 18:11:43','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (310,89,'admin','put','192.168.3.7','/api/v1/apps/14/jenkins-envs/35','2025-11-22 18:11:54','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (311,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-22 18:13:38','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (312,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 18:16:13','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (313,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 18:16:27','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (314,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 18:16:36','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (315,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/deploy','2025-11-22 18:18:03','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (316,89,'admin','post','192.168.3.7','/api/v1/admin/add','2025-11-22 19:26:31','新增管理员');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (317,89,'admin','post','192.168.3.7','/api/v1/admin/add','2025-11-22 19:27:30','新增管理员');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (318,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-22 19:32:00','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (319,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 19:33:08','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (320,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 19:34:25','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (321,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 19:34:41','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (322,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/deploy','2025-11-22 19:36:09','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (323,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-22 20:00:20','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (324,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 20:28:36','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (325,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 20:28:48','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (326,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 20:28:55','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (327,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/deploy','2025-11-22 20:29:16','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (328,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/regression-test','2025-11-22 23:01:55','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (329,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/regression-test','2025-11-22 23:01:55','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (330,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-22 23:13:41','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (331,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 23:13:56','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (332,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 23:16:27','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (333,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-22 23:16:34','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (334,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/deploy','2025-11-22 23:52:17','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (335,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/regression-test','2025-11-23 00:50:11','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (336,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-23 11:32:21','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (337,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-23 11:33:16','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (338,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-23 11:33:23','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (339,89,'admin','delete','192.168.3.7','/api/v1/apps/deployment/32','2025-11-23 11:36:44','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (340,89,'admin','delete','192.168.3.7','/api/v1/apps/deployment/31','2025-11-23 11:36:47','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (341,89,'admin','delete','192.168.3.7','/api/v1/apps/deployment/28','2025-11-23 11:36:51','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (342,89,'admin','delete','192.168.3.7','/api/v1/apps/deployment/22','2025-11-23 11:36:54','操作快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (343,89,'admin','put','192.168.3.7','/api/v1/admin/update','2025-11-23 12:46:17','修改管理员');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (344,89,'admin','put','192.168.3.7','/api/v1/admin/updatePassword','2025-11-23 12:48:11','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (345,98,'lisi','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-23 12:57:50','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (346,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/21','2025-11-23 22:26:59','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (347,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/22','2025-11-23 22:27:01','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (348,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/23','2025-11-23 22:27:04','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (349,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/24','2025-11-23 22:27:07','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (350,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/25','2025-11-23 22:27:10','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (351,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/27','2025-11-23 22:27:12','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (352,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/36','2025-11-23 22:27:15','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (353,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/37','2025-11-23 22:27:17','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (354,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/38','2025-11-23 22:27:20','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (355,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/40','2025-11-23 22:27:24','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (356,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/42','2025-11-23 22:27:27','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (357,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/43','2025-11-23 22:27:37','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (358,98,'lisi','post','192.168.3.7','/api/v1/task/ansible','2025-11-23 23:21:06','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (359,98,'lisi','post','192.168.3.7','/api/v1/task/ansible/45/start','2025-11-23 23:21:54','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (360,98,'lisi','post','192.168.3.7','/api/v1/task/ansible','2025-11-23 23:31:39','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (361,98,'lisi','post','192.168.3.7','/api/v1/task/ansible','2025-11-23 23:34:41','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (362,98,'lisi','post','192.168.3.7','/api/v1/task/ansible/44/start','2025-11-23 23:35:48','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (363,98,'lisi','post','192.168.3.7','/api/v1/task/ansible/47/start','2025-11-23 23:36:10','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (364,98,'lisi','delete','192.168.3.7','/api/v1/task/ansible/46','2025-11-23 23:39:08','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (365,98,'lisi','post','192.168.3.7','/api/v1/cmdb/hostcreate','2025-11-23 23:44:35','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (366,98,'lisi','put','192.168.3.7','/api/v1/config/ecsauthupdate','2025-11-23 23:45:16','修改ECS认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (367,98,'lisi','put','192.168.3.7','/api/v1/cmdb/hostupdate','2025-11-23 23:45:26','修改主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (368,98,'lisi','post','192.168.3.7','/api/v1/cmdb/hostsync','2025-11-23 23:45:31','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (369,98,'lisi','post','192.168.3.7','/api/v1/cmdb/hostcreate','2025-11-23 23:46:04','创建主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (370,98,'lisi','post','192.168.3.7','/api/v1/task/ansible','2025-11-23 23:49:26','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (371,98,'lisi','post','192.168.3.7','/api/v1/task/ansible/48/start','2025-11-23 23:49:34','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (372,98,'lisi','post','192.168.3.7','/api/v1/task/ansible','2025-11-24 00:12:23','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (373,98,'lisi','post','192.168.3.7','/api/v1/task/ansible/49/start','2025-11-24 00:12:28','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (374,89,'admin','post','192.168.3.7','/api/v1/task/ansible','2025-11-24 00:19:30','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (375,89,'admin','post','192.168.3.7','/api/v1/task/ansible/50/start','2025-11-24 00:19:52','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (376,89,'admin','post','192.168.3.7','/api/v1/task/ansible','2025-11-24 00:22:11','创建Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (377,89,'admin','post','192.168.3.7','/api/v1/task/ansible/51/start','2025-11-24 00:24:21','启动Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (378,89,'admin','post','192.168.3.7','/api/v1/task/add','2025-11-24 00:28:34','新增任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (379,89,'admin','post','192.168.3.7','/api/v1/taskjob/start','2025-11-24 00:29:05','启动任务作业');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (380,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-11-24 00:42:25','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (381,89,'admin','put','192.168.3.7','/api/v1/admin/updatePassword','2025-11-24 00:43:19','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (382,89,'admin','put','192.168.3.7','/api/v1/admin/update','2025-11-24 00:43:57','修改管理员');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (383,104,'wangwu','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-24 00:44:47','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (384,98,'lisi','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-24 00:45:40','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (385,98,'lisi','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-11-24 00:45:46','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (386,98,'lisi','put','192.168.3.7','/api/v1/apps/14','2025-11-24 00:46:33','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (387,98,'lisi','post','192.168.3.7','/api/v1/apps/service-release/deploy','2025-11-24 00:46:52','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (388,98,'lisi','post','192.168.3.7','/api/v1/apps/service-release/regression-test','2025-11-24 00:53:04','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (389,89,'admin','put','10.7.16.60','/api/v1/apps/17','2025-11-24 14:59:30','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (390,89,'admin','put','10.7.16.60','/api/v1/apps/14','2025-11-24 14:59:42','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (391,89,'admin','put','10.7.16.60','/api/v1/apps/14','2025-11-24 14:59:49','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (392,89,'admin','put','10.7.16.60','/api/v1/apps/18','2025-11-24 14:59:57','操作应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (393,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostsync','2025-11-25 11:15:58','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (394,89,'admin','post','10.7.16.60','/api/v1/config/ecsauthadd','2025-11-25 16:39:37','新增ECS认证');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (395,89,'admin','put','10.7.16.60','/api/v1/cmdb/hostupdate','2025-11-25 16:39:53','修改主机');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (396,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostsync','2025-11-25 16:39:58','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (397,89,'admin','delete','10.7.16.60','/api/v1/task/ansible/44','2025-11-25 16:40:23','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (398,89,'admin','delete','10.7.16.60','/api/v1/task/ansible/45','2025-11-25 16:40:27','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (399,89,'admin','delete','10.7.16.60','/api/v1/task/ansible/47','2025-11-25 16:40:30','删除Ansible任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (400,89,'admin','put','10.7.16.60','/api/v1/role/assignPermissions','2025-11-26 16:34:08','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (401,89,'admin','put','10.7.16.60','/api/v1/k8s/cluster/33','2025-11-26 16:35:00','操作K8s集群');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (402,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostsync','2025-11-26 17:22:06','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (403,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostssh/upload/501','2025-11-27 11:09:59','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (404,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostsync','2025-11-27 11:10:12','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (405,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostssh/upload/501','2025-11-27 11:10:26','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (406,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostssh/upload/501','2025-11-27 11:13:16','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (407,89,'admin','post','10.7.16.60','/api/v1/cmdb/hostssh/upload/501','2025-11-27 11:19:24','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (408,89,'admin','put','10.7.16.60','/api/v1/role/assignPermissions','2025-11-27 11:22:26','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (409,89,'admin','post','10.7.16.60','/api/v1/menu/add','2025-11-28 11:13:05','新增菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (410,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-11-28 11:13:43','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (411,89,'admin','put','10.7.16.60','/api/v1/role/assignPermissions','2025-11-28 11:14:07','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (412,89,'admin','post','10.7.16.60','/api/v1/apps','2025-11-28 12:39:28','创建应用');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (413,89,'admin','post','10.7.16.60','/api/v1/apps/jenkins-job/validate','2025-11-28 12:39:52','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (414,89,'admin','put','10.7.16.60','/api/v1/apps/19/jenkins-envs/51','2025-11-28 12:39:55','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (415,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release','2025-11-28 12:41:29','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (416,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/approve','2025-11-28 12:42:48','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (417,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/execute','2025-11-28 12:43:09','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (418,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release','2025-11-28 13:37:28','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (419,98,'lisi','post','10.7.16.60','/api/v1/apps/sh-release/approve','2025-11-28 14:18:59','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (420,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/pull-code','2025-11-28 14:47:22','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (421,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/pull-code','2025-11-28 14:47:23','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (422,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release','2025-11-28 14:49:33','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (423,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release','2025-11-28 15:09:19','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (424,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/approve','2025-11-28 15:19:15','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (425,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/pull-code','2025-11-28 15:19:21','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (426,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/execute-script','2025-11-28 15:24:27','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (427,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release','2025-11-28 15:56:13','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (428,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/approve','2025-11-28 15:56:35','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (429,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/pull-code','2025-11-28 15:56:44','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (430,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/execute-script','2025-11-28 15:57:23','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (431,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release','2025-11-28 16:10:22','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (432,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/approve','2025-11-28 16:12:04','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (433,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/pull-code','2025-11-28 16:12:11','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (434,89,'admin','post','10.7.16.60','/api/v1/apps/sh-release/execute-script','2025-11-28 16:19:18','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (435,89,'admin','post','192.168.3.7','/api/v1/menu/add','2025-11-30 13:26:13','新增菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (436,89,'admin','put','192.168.3.7','/api/v1/menu/update','2025-11-30 13:26:29','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (437,89,'admin','put','192.168.3.7','/api/v1/menu/update','2025-11-30 13:26:40','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (438,89,'admin','put','192.168.3.7','/api/v1/role/assignPermissions','2025-11-30 13:27:00','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (439,89,'admin','post','192.168.3.7','/api/v1/tool/deploy','2025-11-30 22:33:06','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (440,89,'admin','post','192.168.3.7','/api/v1/menu/add','2025-11-30 23:20:38','新增菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (441,89,'admin','put','192.168.3.7','/api/v1/role/assignPermissions','2025-11-30 23:22:48','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (442,89,'admin','post','192.168.3.7','/api/v1/apps/sh-release','2025-12-01 01:01:09','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (443,89,'admin','post','192.168.3.7','/api/v1/apps/sh-release/approve','2025-12-01 01:01:48','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (444,89,'admin','post','192.168.3.7','/api/v1/apps/sh-release/pull-code','2025-12-01 01:02:00','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (445,89,'admin','post','192.168.3.7','/api/v1/apps/sh-release/execute-script','2025-12-01 01:02:37','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (446,89,'admin','delete','10.7.16.60','/api/v1/template/delete','2025-12-01 14:41:55','删除任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (447,89,'admin','delete','10.7.16.60','/api/v1/template/delete','2025-12-01 14:41:57','删除任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (448,89,'admin','delete','10.7.16.60','/api/v1/template/delete','2025-12-01 14:41:59','删除任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (449,89,'admin','delete','10.7.16.60','/api/v1/template/delete','2025-12-01 14:42:03','删除任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (450,89,'admin','post','10.7.16.60','/api/v1/template/add','2025-12-01 15:35:40','新增任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (451,89,'admin','post','10.7.16.60','/api/v1/template/add','2025-12-01 15:44:13','新增任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (452,89,'admin','post','10.7.16.60','/api/v1/template/add','2025-12-01 16:56:29','新增任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (453,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-12-01 17:19:34','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (454,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-12-01 17:20:22','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (455,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-12-01 19:59:12','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (456,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-12-01 19:59:31','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (457,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-12-01 19:59:53','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (458,89,'admin','post','10.7.16.60','/api/v1/menu/add','2025-12-01 21:05:27','新增菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (459,89,'admin','put','10.7.16.60','/api/v1/role/assignPermissions','2025-12-01 21:05:55','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (460,89,'admin','delete','10.7.16.60','/api/v1/menu/delete','2025-12-03 10:23:43','删除菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (461,89,'admin','put','10.7.16.60','/api/v1/role/assignPermissions','2025-12-03 10:23:59','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (462,89,'admin','delete','10.7.16.60','/api/v1/menu/delete','2025-12-03 10:24:26','删除菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (463,89,'admin','delete','10.7.16.60','/api/v1/menu/delete','2025-12-03 10:24:32','删除菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (464,89,'admin','delete','10.7.16.60','/api/v1/menu/delete','2025-12-03 10:24:37','删除菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (465,89,'admin','delete','10.7.16.60','/api/v1/monitor/agent/uninstall','2025-12-03 16:17:11','卸载Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (466,89,'admin','delete','10.7.16.60','/api/v1/monitor/agent/uninstall','2025-12-03 16:17:27','卸载Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (467,89,'admin','post','10.7.16.60','/api/v1/monitor/agent/deploy','2025-12-03 16:20:20','部署Agent');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (468,89,'admin','post','10.7.16.60','/api/v1/menu/add','2025-12-03 21:21:04','新增菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (469,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-12-03 21:21:12','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (470,89,'admin','put','10.7.16.60','/api/v1/menu/update','2025-12-03 21:21:25','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (471,89,'admin','post','10.7.16.60','/api/v1/menu/add','2025-12-03 21:22:11','新增菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (472,89,'admin','put','10.7.16.60','/api/v1/role/assignPermissions','2025-12-03 21:22:34','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (473,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-03 21:33:47','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (474,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-03 21:35:47','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (475,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-03 21:40:52','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (476,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-03 21:44:12','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (477,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-03 21:44:42','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (478,89,'admin','delete','10.7.16.60','/api/v1/monitor/domain/1','2025-12-03 21:51:46','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (479,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-03 21:53:16','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (480,89,'admin','delete','10.7.16.60','/api/v1/monitor/domain/3','2025-12-03 21:57:04','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (481,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-03 21:57:13','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (482,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check/4','2025-12-03 21:59:58','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (483,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check/4','2025-12-04 10:09:44','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (484,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 10:19:07','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (485,89,'admin','put','10.7.16.60','/api/v1/monitor/domain','2025-12-04 10:19:52','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (486,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-04 10:20:58','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (487,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-04 10:21:14','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (488,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-04 10:21:29','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (489,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-04 10:21:55','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (490,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-04 10:22:36','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (491,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-04 10:23:18','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (492,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 10:24:52','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (493,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 12:30:45','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (494,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 12:38:17','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (495,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 12:38:57','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (496,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 12:42:20','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (497,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 12:44:47','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (498,89,'admin','post','10.7.16.60','/api/v1/monitor/domain/check-all','2025-12-04 13:14:50','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (499,89,'admin','post','10.7.16.60','/api/v1/monitor/domain','2025-12-04 13:30:26','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (500,89,'admin','put','192.168.3.7','/api/v1/monitor/domain','2025-12-04 22:57:48','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (501,89,'admin','put','192.168.3.7','/api/v1/monitor/domain','2025-12-04 22:57:57','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (502,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/resume','2025-12-04 23:09:59','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (503,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/pause','2025-12-04 23:10:03','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (504,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:18:00','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (505,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:18:13','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (506,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:18:38','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (507,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:18:53','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (508,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:23:01','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (509,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:23:56','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (510,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/check-all','2025-12-04 23:27:10','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (511,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:37:58','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (512,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-04 23:51:38','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (513,89,'admin','delete','192.168.3.7','/api/v1/cmdb/database','2025-12-05 00:25:13','删除数据库');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (514,89,'admin','post','10.7.16.100','/api/v1/monitor/domain/check-all','2025-12-05 13:05:10','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (515,89,'admin','post','192.168.3.7','/api/v1/upload','2025-12-05 20:11:16','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (516,89,'admin','post','192.168.3.7','/api/v1/upload','2025-12-05 20:12:00','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (517,89,'admin','put','192.168.3.7','/api/v1/tool','2025-12-05 20:12:09','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (518,89,'admin','post','192.168.3.7','/api/v1/cmdb/groupadd','2025-12-05 20:17:08','新增资产分组');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (519,89,'admin','post','192.168.3.7','/api/v1/cmdb/groupadd','2025-12-05 20:17:32','新增资产分组');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (520,89,'admin','post','192.168.3.7','/api/v1/cmdb/groupadd','2025-12-05 20:17:44','新增资产分组');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (521,89,'admin','post','192.168.3.7','/api/v1/cmdb/hostsync','2025-12-05 20:18:55','同步主机信息');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (522,89,'admin','post','192.168.3.7','/api/v1/cmdb/hostssh/upload/501','2025-12-05 20:23:26','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (523,89,'admin','post','192.168.3.7','/api/v1/apps/jenkins-job/validate','2025-12-05 20:42:56','验证Jenkins任务');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (524,89,'admin','put','192.168.3.7','/api/v1/apps/19/jenkins-envs/52','2025-12-05 20:43:02','操作Jenkins环境配置');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (525,89,'admin','post','192.168.3.7','/api/v1/apps/deployment/quick','2025-12-05 20:43:58','创建快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (526,89,'admin','post','192.168.3.7','/api/v1/apps/deployment/execute','2025-12-05 20:44:14','执行快速发布');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (527,89,'admin','post','192.168.3.7','/api/v1/apps/service-release','2025-12-05 20:46:54','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (528,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-12-05 20:47:14','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (529,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-12-05 20:47:18','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (530,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/approve','2025-12-05 20:47:27','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (531,89,'admin','post','192.168.3.7','/api/v1/apps/service-release/deploy','2025-12-05 20:47:39','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (532,89,'admin','post','192.168.3.7','/api/v1/monitor/domain','2025-12-05 20:50:59','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (533,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/check/11','2025-12-05 20:52:42','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (534,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-05 20:53:29','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (535,89,'admin','post','192.168.3.7','/api/v1/template/add','2025-12-05 20:54:54','新增任务模板');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (536,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/schedule/config','2025-12-05 23:12:28','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (537,89,'admin','post','192.168.3.7','/api/v1/monitor/domain/check-all','2025-12-06 01:14:24','');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (538,89,'admin','put','192.168.3.7','/api/v1/menu/update','2025-12-06 01:37:15','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (539,89,'admin','put','192.168.3.7','/api/v1/menu/update','2025-12-06 01:37:52','修改菜单');
INSERT INTO `sys_operation_log` (`id`,`admin_id`,`username`,`method`,`ip`,`url`,`create_time`,`description`) VALUES (540,89,'admin','put','192.168.3.7','/api/v1/role/assignPermissions','2025-12-06 21:52:39','');

-- Table structure for sys_post
CREATE TABLE `sys_post` (
 `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
 `post_code` varchar(64) NOT NULL COMMENT '岗位编码',
 `post_name` varchar(50) NOT NULL COMMENT '岗位名称',
 `post_status` int(10) NOT NULL DEFAULT 1 COMMENT '状态（1->正常 2->停用）',
 `create_time` datetime COMMENT '创建时间',
 `remark` varchar(500) COMMENT '备注', 
PRIMARY KEY (id)
) COMMENT '岗位信息表';

-- Data for table sys_post
INSERT INTO `sys_post` (`id`,`post_code`,`post_name`,`post_status`,`create_time`,`remark`) VALUES (1,'AAA','研发总监',1,'2023-06-14 20:08:22','主管各个部门');
INSERT INTO `sys_post` (`id`,`post_code`,`post_name`,`post_status`,`create_time`,`remark`) VALUES (10,'ops','运维工程师',1,'2025-06-28 22:46:33','运维工程师');
INSERT INTO `sys_post` (`id`,`post_code`,`post_name`,`post_status`,`create_time`,`remark`) VALUES (11,'dev','研发工程师',1,'2025-06-28 22:50:29','研发工程师');
INSERT INTO `sys_post` (`id`,`post_code`,`post_name`,`post_status`,`create_time`,`remark`) VALUES (12,'test','测试工程师',2,'2025-06-28 22:52:57','测试工程师');

-- Table structure for sys_role
CREATE TABLE `sys_role` (
 `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键',
 `role_name` varchar(64) NOT NULL COMMENT '角色名称',
 `role_key` varchar(64) NOT NULL COMMENT '角色权限字符串',
 `status` int(10) NOT NULL DEFAULT 1 COMMENT '启用状态：1->启用；2->禁用',
 `description` varchar(500) COMMENT '描述',
 `create_time` datetime NOT NULL COMMENT '创建时间', 
PRIMARY KEY (id)
) COMMENT '后台角色表';
ALTER TABLE `sys_role` ADD  INDEX `create_time`(`create_time`) USING BTREE;
ALTER TABLE `sys_role` ADD  INDEX `name`(`role_name`) USING BTREE;

-- Data for table sys_role
INSERT INTO `sys_role` (`id`,`role_name`,`role_key`,`status`,`description`,`create_time`) VALUES (1,'超级管理员','admin',1,'最大权限','2023-06-12 20:04:53');
INSERT INTO `sys_role` (`id`,`role_name`,`role_key`,`status`,`description`,`create_time`) VALUES (13,'游客','test1',1,'test1','2025-07-03 18:47:25');

-- Table structure for sys_role_menu
CREATE TABLE `sys_role_menu` (
 `role_id` int(10) COMMENT '角色ID',
 `menu_id` int(10) COMMENT '菜单ID'
) COMMENT '角色和菜单关系表';

-- Data for table sys_role_menu
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,151);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,155);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,90);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,91);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,148);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,169);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,171);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,176);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,177);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,183);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,184);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,187);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,188);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,189);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,191);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,192);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,194);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,196);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,198);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,199);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,201);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,202);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,204);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,206);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,143);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,145);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,163);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,164);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,141);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,161);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,212);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,213);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,138);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,157);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,134);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,135);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,131);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,126);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,129);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,17);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,18);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,60);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,22);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,23);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,24);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,27);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,28);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,30);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,31);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,33);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,34);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,123);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,124);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,114);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,115);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,117);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,47);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,113);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,49);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,151);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,155);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,90);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,91);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,148);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,169);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,171);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,176);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,177);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,183);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,184);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,187);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,188);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,189);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,191);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,192);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,194);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,196);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,198);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,199);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,201);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,202);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,204);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,206);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,143);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,145);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,163);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,164);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,141);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,161);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,212);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,213);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,138);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,157);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,134);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,135);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,131);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,126);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,129);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,17);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,18);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,60);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,22);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,23);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,24);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,27);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,28);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,30);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,31);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,33);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,34);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,123);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,124);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,114);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,115);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,117);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,47);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,113);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,49);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,151);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,155);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,90);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,91);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,148);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,169);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,171);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,176);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,177);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,183);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,184);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,187);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,188);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,189);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,191);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,192);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,194);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,196);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,198);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,199);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,201);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,202);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,204);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,206);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,143);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,145);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,163);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,164);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,141);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,161);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,212);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,213);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,138);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,157);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,134);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,135);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,131);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,126);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,129);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,17);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,18);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,60);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,22);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,23);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,24);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,27);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,28);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,30);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,31);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,33);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,34);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,123);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,124);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,114);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,115);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,117);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,47);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,113);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,49);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,151);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,155);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,90);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,91);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,148);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,169);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,171);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,176);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,177);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,183);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,184);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,187);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,188);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,189);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,191);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,192);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,194);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,196);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,198);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,199);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,201);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,202);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,204);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,206);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,143);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,145);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,163);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,164);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,141);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,161);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,212);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,213);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,138);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,157);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,134);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,135);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,131);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,126);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,129);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,17);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,18);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,60);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,22);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,23);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,24);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,27);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,28);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,30);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,31);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,33);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,34);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,123);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,124);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,114);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,115);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,117);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,47);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,113);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,49);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,151);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,155);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,90);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,91);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,148);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,169);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,171);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,176);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,177);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,183);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,184);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,187);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,188);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,189);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,191);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,192);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,194);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,196);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,198);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,199);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,201);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,202);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,204);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,206);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,143);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,145);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,163);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,164);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,141);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,161);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,212);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,213);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,138);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,157);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,134);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,135);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,131);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,126);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,129);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,17);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,18);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,60);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,22);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,23);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,24);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,27);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,28);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,30);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,31);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,33);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,34);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,123);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,124);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,114);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,115);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,117);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,47);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,113);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,49);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,151);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,155);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,90);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,91);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,148);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,169);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,171);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,176);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,177);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,183);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,184);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,187);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,188);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,189);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,191);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,192);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,194);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,196);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,198);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,199);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,201);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,202);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,204);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,206);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,143);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,145);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,163);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,164);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,141);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,161);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,212);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,213);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,138);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,157);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,134);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,135);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,131);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,126);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,129);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,17);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,18);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,60);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,22);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,23);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,24);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,27);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,28);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,30);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,31);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,33);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,34);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,123);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,124);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,114);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,115);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,117);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,47);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,113);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,49);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,109);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (13,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,72);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,80);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,78);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,149);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,150);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,151);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,152);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,153);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,154);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,155);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,165);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,166);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,88);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,89);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,90);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,91);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,95);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,146);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,148);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,81);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,82);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,167);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,168);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,169);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,170);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,171);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,83);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,172);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,173);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,174);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,175);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,176);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,177);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,105);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,178);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,179);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,180);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,181);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,182);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,183);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,184);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,93);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,185);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,186);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,187);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,188);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,189);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,190);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,191);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,192);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,193);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,194);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,195);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,196);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,106);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,197);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,198);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,199);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,200);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,201);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,202);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,203);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,204);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,205);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,206);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,108);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,107);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,110);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,142);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,143);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,144);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,145);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,162);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,163);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,164);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,111);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,139);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,140);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,141);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,160);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,161);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,212);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,213);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,97);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,99);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,136);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,137);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,138);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,156);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,157);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,98);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,133);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,134);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,135);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,100);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,130);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,131);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,132);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,159);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,101);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,103);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,102);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,126);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,127);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,128);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,129);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,4);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,6);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,16);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,17);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,18);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,60);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,7);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,21);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,22);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,23);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,24);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,8);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,26);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,27);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,28);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,9);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,29);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,30);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,31);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,10);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,32);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,33);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,34);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,84);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,85);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,123);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,124);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,125);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,86);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,122);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,119);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,120);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,121);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,104);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,114);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,115);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,117);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,118);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,44);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,45);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,47);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,113);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,73);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,46);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,49);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,62);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,96);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) VALUES (1,109);

-- Table structure for task_ansible
CREATE TABLE `task_ansible` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
 `name` varchar(100) NOT NULL COMMENT '''任务名称''',
 `description` text COMMENT '''任务描述''',
 `type` bigint(19) NOT NULL DEFAULT 1 COMMENT '''任务类型:1-手动,2-Git,3-K8s''',
 `git_repo` varchar(255) COMMENT '''Git仓库地址''',
 `host_groups` text NOT NULL COMMENT '''主机分组JSON''',
 `all_host_ids` text NOT NULL COMMENT '''所有主机ID JSON数组''',
 `global_vars` text COMMENT '''全局变量JSON''',
 `status` bigint(19) NOT NULL DEFAULT 1 COMMENT '''任务状态:1-等待中,2-运行中,3-成功,4-异常''',
 `created_at` datetime NOT NULL COMMENT '''创建时间''',
 `updated_at` datetime NOT NULL COMMENT '''更新时间''',
 `error_msg` text COMMENT '''错误信息''',
 `task_count` bigint(19) NOT NULL DEFAULT 0 COMMENT '''任务数量(Type=1时为上传文件数,Type=2时为解析的playbook数,Type=3时固定为1)''',
 `total_duration` bigint(19) NOT NULL DEFAULT 0 COMMENT '''任务执行总耗时(秒,所有子任务耗时总和)''', 
PRIMARY KEY (id)
);
ALTER TABLE `task_ansible` ADD unique INDEX `idx_task_ansible_name`(`name`) USING BTREE;
ALTER TABLE `task_ansible` ADD  INDEX `idx_task_status`(`status`) USING BTREE;

-- Data for table task_ansible
INSERT INTO `task_ansible` (`id`,`name`,`description`,`type`,`git_repo`,`host_groups`,`all_host_ids`,`global_vars`,`status`,`created_at`,`updated_at`,`error_msg`,`task_count`,`total_duration`) VALUES (48,'本地多主机测试','',2,'git@gitee.com:zhang_fan1024/ansible-playbook.git','{"api":[512],"web":[511]}','[511,512]','',4,'2025-11-23 23:49:23','2025-11-23 23:50:42','',2,67);
INSERT INTO `task_ansible` (`id`,`name`,`description`,`type`,`git_repo`,`host_groups`,`all_host_ids`,`global_vars`,`status`,`created_at`,`updated_at`,`error_msg`,`task_count`,`total_duration`) VALUES (49,'test2','',1,'','{"web":[511]}','[511]','',3,'2025-11-24 00:12:22','2025-11-24 00:12:47','',1,17);
INSERT INTO `task_ansible` (`id`,`name`,`description`,`type`,`git_repo`,`host_groups`,`all_host_ids`,`global_vars`,`status`,`created_at`,`updated_at`,`error_msg`,`task_count`,`total_duration`) VALUES (50,'test3-手动任务','',1,'','{"web":[511]}','[511]','',3,'2025-11-24 00:19:29','2025-11-24 00:20:20','',1,27);
INSERT INTO `task_ansible` (`id`,`name`,`description`,`type`,`git_repo`,`host_groups`,`all_host_ids`,`global_vars`,`status`,`created_at`,`updated_at`,`error_msg`,`task_count`,`total_duration`) VALUES (51,'test4-自动任务','',2,'git@gitee.com:zhang_fan1024/ansible-playbook.git','{"api":[512],"web":[511]}','[511,512]','',4,'2025-11-24 00:22:08','2025-11-24 00:26:06','',2,103);

-- Table structure for task_ansiblework
CREATE TABLE `task_ansiblework` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
 `task_id` bigint(20) NOT NULL COMMENT '''父任务ID''',
 `entry_file_name` varchar(255) NOT NULL COMMENT '''入口文件名''',
 `entry_file_path` varchar(255) NOT NULL COMMENT '''入口文件路径''',
 `log_path` varchar(255) COMMENT '''日志路径''',
 `status` bigint(19) NOT NULL DEFAULT 1 COMMENT '''子任务状态:1-等待中,2-运行中,3-成功,4-异常''',
 `start_time` datetime COMMENT '''开始时间''',
 `end_time` datetime COMMENT '''结束时间''',
 `duration` bigint(19) COMMENT '''执行耗时(秒)''',
 `exit_code` bigint(19) COMMENT '''退出代码''',
 `error_msg` text COMMENT '''错误信息''',
 `log` text COMMENT '''日志内容''', 
PRIMARY KEY (id)
);
ALTER TABLE `task_ansiblework` ADD  INDEX `idx_task_ansiblework_task_id`(`task_id`) USING BTREE;
ALTER TABLE `task_ansiblework` ADD  INDEX `idx_task_id`(`task_id`) USING BTREE;
ALTER TABLE `task_ansiblework` ADD  INDEX `idx_task_work_composite`(`task_id`,`status`) USING BTREE;

-- Data for table task_ansiblework
INSERT INTO `task_ansiblework` (`id`,`task_id`,`entry_file_name`,`entry_file_path`,`log_path`,`status`,`start_time`,`end_time`,`duration`,`exit_code`,`error_msg`,`log`) VALUES (47,48,'01-linux-os.yaml','task/48/本地多主机测试/01-linux-os.yaml','logs/ansible/48/47/01-linux-os.yaml.log',3,'2025-11-23 23:49:34','2025-11-23 23:50:09',35,0,'',NULL);
INSERT INTO `task_ansiblework` (`id`,`task_id`,`entry_file_name`,`entry_file_path`,`log_path`,`status`,`start_time`,`end_time`,`duration`,`exit_code`,`error_msg`,`log`) VALUES (48,48,'02-os.yaml','task/48/本地多主机测试/02-os.yaml','logs/ansible/48/48/02-os.yaml.log',4,'2025-11-23 23:50:10','2025-11-23 23:50:42',32,2,'exit status 2',NULL);
INSERT INTO `task_ansiblework` (`id`,`task_id`,`entry_file_name`,`entry_file_path`,`log_path`,`status`,`start_time`,`end_time`,`duration`,`exit_code`,`error_msg`,`log`) VALUES (49,49,'01-linux-os.yaml.yml','task/49/test2/01-linux-os.yaml.yml','logs/ansible/49/49/01-linux-os.yaml.yml.log',3,'2025-11-24 00:12:28','2025-11-24 00:12:46',17,0,'',NULL);
INSERT INTO `task_ansiblework` (`id`,`task_id`,`entry_file_name`,`entry_file_path`,`log_path`,`status`,`start_time`,`end_time`,`duration`,`exit_code`,`error_msg`,`log`) VALUES (50,50,'01-linux-os.yaml.yml','task/50/test3-手动任务/01-linux-os.yaml.yml','logs/ansible/50/50/01-linux-os.yaml.yml.log',3,'2025-11-24 00:19:52','2025-11-24 00:20:19',27,0,'',NULL);
INSERT INTO `task_ansiblework` (`id`,`task_id`,`entry_file_name`,`entry_file_path`,`log_path`,`status`,`start_time`,`end_time`,`duration`,`exit_code`,`error_msg`,`log`) VALUES (51,51,'01-linux-os.yaml','task/51/test4-自动任务/01-linux-os.yaml','logs/ansible/51/51/01-linux-os.yaml.log',3,'2025-11-24 00:24:21','2025-11-24 00:25:02',40,0,'',NULL);
INSERT INTO `task_ansiblework` (`id`,`task_id`,`entry_file_name`,`entry_file_path`,`log_path`,`status`,`start_time`,`end_time`,`duration`,`exit_code`,`error_msg`,`log`) VALUES (52,51,'02-os.yaml','task/51/test4-自动任务/02-os.yaml','logs/ansible/51/52/02-os.yaml.log',4,'2025-11-24 00:25:02','2025-11-24 00:26:05',63,2,'exit status 2',NULL);

-- Table structure for task_job
CREATE TABLE `task_job` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(255) COMMENT '任务标题',
 `type` bigint(19) COMMENT '任务类型 1=普通任务,2=定时任务,3=ansible任务',
 `shell` text COMMENT '任务内容(任务模板ID,多个用逗号分隔)',
 `host_ids` text COMMENT '主机ID(多个用逗号分隔)',
 `cron_expr` varchar(255) COMMENT '定时表达式(* * * * *)',
 `tasklog` text COMMENT '任务执行日志',
 `status` bigint(19) COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常,5=已暂停',
 `duration` bigint(19) COMMENT '执行耗时(秒)',
 `remark` text COMMENT '任务备注',
 `start_time` datetime COMMENT '任务开始时间',
 `end_time` datetime COMMENT '任务结束时间',
 `created_at` datetime COMMENT '任务创建时间',
 `task_count` bigint(19) COMMENT '任务数量',
 `is_recurring` tinyint(3) COMMENT '是否周期性任务',
 `scheduled_time` datetime COMMENT '计划执行时间',
 `log_path` varchar(500) COMMENT '日志文件路径',
 `execute_count` bigint(19) DEFAULT 0 COMMENT '执行次数',
 `next_run_time` datetime COMMENT '下次执行时间', 
PRIMARY KEY (id)
);

-- Data for table task_job
INSERT INTO `task_job` (`id`,`name`,`type`,`shell`,`host_ids`,`cron_expr`,`tasklog`,`status`,`duration`,`remark`,`start_time`,`end_time`,`created_at`,`task_count`,`is_recurring`,`scheduled_time`,`log_path`,`execute_count`,`next_run_time`) VALUES (49,'test4',1,'12','1','','',3,5,'',NULL,'2025-09-29 21:15:32','2025-09-29 21:01:51',1,NULL,NULL,NULL,1,NULL);
INSERT INTO `task_job` (`id`,`name`,`type`,`shell`,`host_ids`,`cron_expr`,`tasklog`,`status`,`duration`,`remark`,`start_time`,`end_time`,`created_at`,`task_count`,`is_recurring`,`scheduled_time`,`log_path`,`execute_count`,`next_run_time`) VALUES (56,'定时任务01',2,'11','1','*/3 * * * *','',5,5,'',NULL,'2025-10-01 14:03:07','2025-09-30 17:15:17',1,NULL,NULL,NULL,26,'2025-10-01 14:06:00');
INSERT INTO `task_job` (`id`,`name`,`type`,`shell`,`host_ids`,`cron_expr`,`tasklog`,`status`,`duration`,`remark`,`start_time`,`end_time`,`created_at`,`task_count`,`is_recurring`,`scheduled_time`,`log_path`,`execute_count`,`next_run_time`) VALUES (57,'定时任务02',2,'12','1','*/5 * * * *','',5,6,'',NULL,'2025-10-01 14:00:08','2025-09-30 17:17:02',1,NULL,NULL,NULL,18,'2025-10-01 14:05:00');
INSERT INTO `task_job` (`id`,`name`,`type`,`shell`,`host_ids`,`cron_expr`,`tasklog`,`status`,`duration`,`remark`,`start_time`,`end_time`,`created_at`,`task_count`,`is_recurring`,`scheduled_time`,`log_path`,`execute_count`,`next_run_time`) VALUES (58,'test5',1,'10','1','','',3,6,'123',NULL,'2025-10-01 14:53:47','2025-10-01 14:53:31',1,NULL,NULL,NULL,1,NULL);
INSERT INTO `task_job` (`id`,`name`,`type`,`shell`,`host_ids`,`cron_expr`,`tasklog`,`status`,`duration`,`remark`,`start_time`,`end_time`,`created_at`,`task_count`,`is_recurring`,`scheduled_time`,`log_path`,`execute_count`,`next_run_time`) VALUES (59,'test123',1,'2','501','','',3,106,'123',NULL,'2025-11-06 15:52:36','2025-11-06 15:50:35',1,NULL,NULL,NULL,1,NULL);
INSERT INTO `task_job` (`id`,`name`,`type`,`shell`,`host_ids`,`cron_expr`,`tasklog`,`status`,`duration`,`remark`,`start_time`,`end_time`,`created_at`,`task_count`,`is_recurring`,`scheduled_time`,`log_path`,`execute_count`,`next_run_time`) VALUES (60,'test456',1,'3','501','','',1,0,'456',NULL,NULL,'2025-11-06 15:51:56',1,NULL,NULL,NULL,0,NULL);
INSERT INTO `task_job` (`id`,`name`,`type`,`shell`,`host_ids`,`cron_expr`,`tasklog`,`status`,`duration`,`remark`,`start_time`,`end_time`,`created_at`,`task_count`,`is_recurring`,`scheduled_time`,`log_path`,`execute_count`,`next_run_time`) VALUES (61,'打印数字',1,'2','511','','',3,95,'123',NULL,'2025-11-24 00:30:41','2025-11-24 00:28:34',1,NULL,NULL,NULL,1,NULL);

-- Table structure for task_template
CREATE TABLE `task_template` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(100) NOT NULL,
 `type` bigint(19) NOT NULL,
 `content` text NOT NULL,
 `created_at` datetime,
 `updated_at` datetime,
 `created_by` varchar(50),
 `updated_by` varchar(50),
 `remark` varchar(500), 
PRIMARY KEY (id)
);

-- Data for table task_template
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (2,'1-数字 1-100',1,'#!/bin/bash

# 从 1 到 100，每秒打印一个数字
for ((i = 1; i <= 100; i++)); do
    echo "[$(date +%H:%M:%S)] $i"
    sleep 1
done

echo "完成：所有数字 1-100 已打印完毕。"
','2025-08-06 12:47:57','2025-08-12 16:14:49','admin','admin','测试脚本');
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (3,'2-系统状态检查脚本',1,'#!/bin/bash

#==================================================
#       简化版 Linux 系统状态检查脚本
#       输出直接打印到终端，不保存日志
#==================================================

# 颜色定义（可选，让输出更清晰）
GREEN=''\033[0;32m''
YELLOW=''\033[1;33m''
RED=''\033[0;31m''
NC=''\033[0m''

log() { echo -e "${GREEN}[INFO] $1${NC}"; }
warn() { echo -e "${YELLOW}[WARN] $1${NC}"; }
error() { echo -e "${RED}[ERROR] $1${NC}"; }

separator() {
    echo "----------------------------------------"
}

# 1. 系统信息
separator
log "系统信息"
separator
echo "主机名: $(hostname)"
echo "操作系统: $(grep ''^PRETTY_NAME='' /etc/os-release | cut -d''"'' -f2 2>/dev/null || echo ''未知'')"
echo "内核: $(uname -r)"

# 2. 系统负载
separator
log "系统负载 (1/5/15分钟)"
separator
uptime | awk -F''load average:'' ''{print "负载: " $2}''

# 3. CPU 使用率
separator
log "CPU 使用率"
separator
top -bn1 | grep "Cpu(s)" | awk ''{print "CPU 使用: " $2 "%"}''

# 4. 内存使用
separator
log "内存使用"
separator
free -h | awk ''/^Mem:/ {printf "总内存: %s\t已用: %s\t空闲: %s\n", $2, $3, $4}''

# 5. 磁盘使用（仅根分区）
separator
log "磁盘使用 (/)"
separator
df -h / | awk ''NR==2 {printf "设备: %s\t总大小: %s\t已用: %s (%s)\n", $1, $2, $3, $5}''

# 6. 活跃进程数
separator
log "进程统计"
separator
PROC=$(ps aux | wc -l)
echo "当前运行进程数: $((PROC - 1))"

# 7. 网络监听端口（常用服务）
separator
log "监听端口 (常见服务)"
separator
ss -tuln | grep '':80\|:443\|:22\|:3306\|:5432'' || echo "无常见服务监听"

# 8. 最近登录与 SSH 尝试
separator
log "最近登录与安全"
separator
echo "最近3次登录:"
last | head -3 | awk ''{print $1, $3, $4, $5}''

FAILED=$(grep "Failed password" /var/log/auth.log 2>/dev/null | wc -l)
if [ $FAILED -gt 0 ]; then
    warn "发现 $FAILED 次 SSH 登录失败，请检查安全！"
else
    log "无 SSH 登录失败记录"
fi

separator
log "系统状态检查完成。"','2025-08-06 12:48:34','2025-08-12 16:14:21','root','admin','test');
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (4,'3-自动优化 Linux 内核参数',2,'#!/bin/bash

# 自动优化 Linux 内核参数
# 脚本针对 RHEL7

cat >> /usr/lib/sysctl.d/00-system.conf << ''EOF''
fs.file-max=65535
net.ipv4.tcp_timestamps = 0
net.ipv4.tcp_synack_retries = 5
net.ipv4.tcp_syn_retries = 5
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_fin_timeout = 30
#net.ipv4.tcp_keepalive_time = 120
net.ipv4.ip_local_port_range = 1024 65535
kernel.shmall = 2097152
kernel.shmmax = 2147483648
kernel.shmmni = 4096
kernel.sem = 5010 641280 5010 128
net.core.wmem_default=262144
net.core.wmem_max=262144
net.core.rmem_default=4194304
net.core.rmem_max=4194304
net.ipv4.tcp_fin_timeout = 10
net.ipv4.tcp_keepalive_time = 30
net.ipv4.tcp_window_scaling = 0
net.ipv4.tcp_sack = 0
EOF

sysctl -p
sleep 3
echo  "已经完成优化"
','2025-08-06 12:49:39','2025-09-06 16:27:02','lisi','admin','test');
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (6,'切割 Nginx 日志文件',2,'#mkdir  /data/scripts  
#vim   /data/scripts/nginx_log.sh    
#!/bin/bash  
  
# 切割 Nginx 日志文件(防止单个文件过大,后期处理很困难)   
logs_path="/usr/local/nginx/logs/"  
mv ${logs_path}access.log ${logs_path}access_$(date -d "yesterday" +"%Y%m%d").log  
kill -USR1  `cat /usr/local/nginx/logs/nginx.pid`  
  
# chmod +x  /data/scripts/nginx_log.sh  
# crontab  ‐e                    #脚本写完后,将脚本放入计划任务每天执行一次脚本  
0  1  *  *   *   /data/scripts/nginx_log.sh  ','2025-08-06 12:53:39','2025-09-06 16:27:07','zhangsan','admin','test');
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (9,'查找 Linux 系统中的僵尸进程  ',2,'#!/bin/bash  
  
# 查找 Linux 系统中的僵尸进程  
# awk 判断 ps 命令输出的第 8 列为 Z 是,显示该进程的 PID 和进程命令  
ps aux | awk ''{if($8 == "Z"){print $2,$11}}''  ','2025-08-06 14:45:58','2025-09-06 16:27:18','admin','admin','查找 Linux 系统中的僵尸进程  ');
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (14,'test123',1,'#!/bin/bash

# ==============================================
# 系统初始化配置脚本 - 完全非交互式版本
# 支持：CentOS/RHEL/龙溪/UOS/麒麟/欧拉/Ubuntu/Debian等
# 用法：curl -sSL https://raw.githubusercontent.com/xxx/init.sh | bash
#       或 wget -qO- https://raw.githubusercontent.com/xxx/init.sh | bash
# ==============================================

set -e  # 遇到错误退出

# ==============================================
# 配置变量（可根据需要修改）
# ==============================================

# 基本配置
DISABLE_SELINUX="yes"           # 是否禁用SELinux: yes/no
DISABLE_IPV6="yes"              # 是否禁用IPv6: yes/no
INSTALL_PACKAGES="yes"         # 是否安装常用工具: yes/no
ENABLE_SWAP="yes"              # 是否启用swap优化: yes/no
SSH_PORT=""                    # SSH端口，留空保持默认22

# 性能参数
FILE_LIMIT="65535"            # 文件描述符限制
TCP_MAX_TW_BUCKETS="8000"     # TIME_WAIT连接数
TCP_TW_REUSE="1"              # 重用TIME_WAIT连接
SWAPPINESS="10"               # swappiness值

# 日志配置
LOG_DIR="/var/log/system_init"
BACKUP_DIR="/backup/system_config_$(date +%Y%m%d)"

# ==============================================
# 初始化设置
# ==============================================

# 颜色输出
RED=''\033[0;31m''
GREEN=''\033[0;32m''
YELLOW=''\033[1;33m''
BLUE=''\033[0;34m''
NC=''\033[0m''

# 日志初始化
init_logging() {
    mkdir -p "$LOG_DIR"
    LOG_FILE="$LOG_DIR/init_$(date +%Y%m%d_%H%M%S).log"
    exec > >(tee -a "$LOG_FILE") 2>&1
    
    echo "========================================="
    echo "系统初始化脚本 - 非交互式版本"
    echo "开始时间: $(date ''+%Y-%m-%d %H:%M:%S'')"
    echo "日志文件: $LOG_FILE"
    echo "========================================="
    echo ""
}

# 检查root权限
check_root() {
    if [ "$EUID" -ne 0 ]; then 
        echo -e "${RED}错误: 请使用root权限运行此脚本${NC}"
        echo "用法: sudo $0 或 curl -sSL URL | sudo bash"
        exit 1
    fi
}

# 系统检测
detect_system() {
    echo -e "${BLUE}[1/10] 检测系统信息...${NC}"
    
    # 检测系统
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_NAME="$NAME"
        OS_ID="$ID"
        OS_VERSION="$VERSION_ID"
        OS_ID_LIKE="$ID_LIKE"
    elif [ -f /etc/redhat-release ]; then
        OS_NAME=$(cat /etc/redhat-release)
        OS_ID="rhel"
        OS_VERSION=$(cat /etc/redhat-release | grep -oE ''[0-9]+\.[0-9]+'' || echo "unknown")
    else
        OS_NAME="未知"
        OS_ID="unknown"
        OS_VERSION="unknown"
    fi
    
    # 架构和内核
    ARCH=$(uname -m)
    KERNEL=$(uname -r)
    
    # 内存和CPU
    MEM_TOTAL=$(free -g | awk ''/^Mem:/{print $2}'' || echo "0")
    CPU_CORES=$(nproc 2>/dev/null || grep -c ^processor /proc/cpuinfo || echo "1")
    
    # 存储
    DISK_TOTAL=$(df -h / | awk ''NR==2 {print $2}'' || echo "unknown")
    
    echo "✅ 系统检测完成:"
    echo "   - 系统: $OS_NAME"
    echo "   - 版本: $OS_VERSION"
    echo "   - 架构: $ARCH"
    echo "   - 内核: $KERNEL"
    echo "   - CPU核心: $CPU_CORES"
    echo "   - 内存: ${MEM_TOTAL}GB"
    echo "   - 根分区: $DISK_TOTAL"
    echo ""
    
    export OS_ID OS_VERSION ARCH MEM_TOTAL CPU_CORES
}

# ==============================================
# 备份功能
# ==============================================

backup_configs() {
    echo -e "${BLUE}[2/10] 备份系统配置...${NC}"
    
    mkdir -p "$BACKUP_DIR"
    
    # 要备份的配置文件列表
    local configs=(
        "/etc/sysctl.conf"
        "/etc/security/limits.conf"
        "/etc/ssh/sshd_config"
        "/etc/selinux/config"
        "/etc/hosts"
        "/etc/resolv.conf"
        "/etc/fstab"
        "/etc/profile"
        "/etc/bashrc"
    )
    
    local backup_count=0
    for config in "${configs[@]}"; do
        if [ -f "$config" ]; then
            cp -p "$config" "$BACKUP_DIR/" 2>/dev/null && {
                echo "   📁 已备份: $config"
                ((backup_count++))
            } || echo "   ⚠️  备份失败: $config"
        fi
    done
    
    # 备份已加载的内核参数
    sysctl -a 2>/dev/null > "$BACKUP_DIR/sysctl_current.conf" && \
        echo "   📁 已备份当前内核参数"
    
    echo "✅ 备份完成: 共备份 $backup_count 个文件到 $BACKUP_DIR"
    echo ""
}

# ==============================================
# 内核优化
# ==============================================

optimize_kernel() {
    echo -e "${BLUE}[3/10] 优化内核参数...${NC}"
    
    # 创建优化配置文件
    local sysctl_file="/etc/sysctl.d/99-system-optimize.conf"
    
    cat > "$sysctl_file" << EOF
# ==============================================
# 系统内核优化配置
# 生成时间: $(date)
# 系统: $OS_NAME $OS_VERSION
# 内存: ${MEM_TOTAL}GB
# ==============================================

# --- 内存优化 ---
# swappiness (0-100, 值越低越少使用swap)
vm.swappiness = $SWAPPINESS

# 缓存压力 (1-100, 值越低越倾向缓存)
vm.vfs_cache_pressure = 50

# 脏页写回比例
vm.dirty_ratio = 10
vm.dirty_background_ratio = 5

# 内存分配策略 (1允许超额分配)
vm.overcommit_memory = 1
vm.overcommit_ratio = 90

# 虚拟内存映射数量 (Elasticsearch等需要)
vm.max_map_count = 262144

# --- 网络优化 ---
# TCP连接队列
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 65535
net.ipv4.tcp_max_syn_backlog = 65536

# TIME_WAIT优化
net.ipv4.tcp_max_tw_buckets = $TCP_MAX_TW_BUCKETS
net.ipv4.tcp_tw_reuse = $TCP_TW_REUSE
net.ipv4.tcp_fin_timeout = 30

# 端口范围
net.ipv4.ip_local_port_range = 1024 65535

# Keepalive设置
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_keepalive_probes = 5
net.ipv4.tcp_keepalive_intvl = 15

# 缓冲区大小
net.core.rmem_default = 262144
net.core.rmem_max = 16777216
net.core.wmem_default = 262144
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216

# TCP拥塞控制
net.ipv4.tcp_congestion_control = cubic

# 其他优化
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_synack_retries = 2
net.ipv4.tcp_syn_retries = 3
net.ipv4.tcp_max_orphans = 65536
net.ipv4.tcp_timestamps = 1
net.ipv4.tcp_sack = 1
net.ipv4.tcp_window_scaling = 1
net.ipv4.tcp_fastopen = 3

# 根据内存大小调整TCP内存
EOF
    
    # 根据内存调整TCP内存参数
    if [ "$MEM_TOTAL" -lt 4 ]; then
        echo "net.ipv4.tcp_mem = 98304 131072 196608" >> "$sysctl_file"
    elif [ "$MEM_TOTAL" -lt 16 ]; then
        echo "net.ipv4.tcp_mem = 196608 262144 393216" >> "$sysctl_file"
    else
        echo "net.ipv4.tcp_mem = 393216 524288 786432" >> "$sysctl_file"
    fi
    
    # IPv6配置
    if [ "$DISABLE_IPV6" = "yes" ]; then
        cat >> "$sysctl_file" << EOF

# --- IPv6禁用 ---
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6 = 1
EOF
    fi
    
    # 应用配置
    sysctl -p "$sysctl_file" >/dev/null 2>&1
    echo "✅ 内核参数优化完成: $sysctl_file"
    echo ""
}

# ==============================================
# 资源限制设置
# ==============================================

set_resource_limits() {
    echo -e "${BLUE}[4/10] 设置系统资源限制...${NC}"
    
    # 备份原文件
    if [ -f /etc/security/limits.conf ]; then
        cp /etc/security/limits.conf /etc/security/limits.conf.bak.$(date +%Y%m%d)
    fi
    
    # 清理旧的自定义设置（如果有）
    sed -i ''/^# ===== 系统资源限制配置 =====/,/^# ===============================/d'' /etc/security/limits.conf 2>/dev/null
    
    # 添加新的限制
    cat >> /etc/security/limits.conf << EOF

# ===== 系统资源限制配置 =====
# 配置时间: $(date)
# 文件描述符限制: $FILE_LIMIT
# ===============================

# 全局文件描述符限制
* soft nofile $FILE_LIMIT
* hard nofile $FILE_LIMIT

# 进程数限制
* soft nproc $FILE_LIMIT
* hard nproc $FILE_LIMIT

# 核心文件大小
* soft core unlimited
* hard core unlimited

# 栈大小
* soft stack 10240
* hard stack 32768

# 实时优先级 (某些应用需要)
#@realtime soft rtprio 99
#@realtime hard rtprio 99
#@realtime soft memlock unlimited
#@realtime hard memlock unlimited

EOF
    
    # 对于systemd系统
    if [ -d /etc/systemd/system.conf.d ]; then
        mkdir -p /etc/systemd/system.conf.d
        cat > /etc/systemd/system.conf.d/limits.conf << EOF
[Manager]
DefaultLimitNOFILE=$FILE_LIMIT
DefaultLimitNPROC=$FILE_LIMIT
DefaultLimitMEMLOCK=infinity
EOF
    fi
    
    # 对于使用pam_limits的系统
    if [ -f /etc/pam.d/common-session ] && [ "$OS_ID" = "ubuntu" ] || [ "$OS_ID" = "debian" ]; then
        grep -q "pam_limits.so" /etc/pam.d/common-session || \
            echo "session required pam_limits.so" >> /etc/pam.d/common-session
    fi
    
    echo "✅ 资源限制设置完成: 文件限制=$FILE_LIMIT"
    echo ""
}

# ==============================================
# SELinux配置
# ==============================================

configure_selinux() {
    echo -e "${BLUE}[5/10] 配置SELinux...${NC}"
    
    # 检查SELinux是否安装
    if ! command -v getenforce >/dev/null 2>&1 && [ ! -f /etc/selinux/config ]; then
        echo "   ℹ️  SELinux未安装，跳过配置"
        return 0
    fi
    
    local current_state=$(getenforce 2>/dev/null || echo "Disabled")
    echo "   当前SELinux状态: $current_state"
    
    if [ "$DISABLE_SELINUX" = "yes" ]; then
        echo "   🔧 正在禁用SELinux..."
        
        # 临时设置为permissive
        setenforce 0 2>/dev/null || true
        
        # 永久禁用
        if [ -f /etc/selinux/config ]; then
            sed -i ''s/^SELINUX=.*/SELINUX=disabled/'' /etc/selinux/config
            sed -i ''s/^SELINUXTYPE=.*/#SELINUXTYPE=targeted/'' /etc/selinux/config
            echo "   ✅ SELinux已永久禁用，重启后生效"
        fi
    else
        echo "   ℹ️  保持SELinux当前状态 ($current_state)"
        echo "   提示: 如需禁用，设置 DISABLE_SELINUX=''yes''"
    fi
    
    echo ""
}

# ==============================================
# SSH优化
# ==============================================

optimize_ssh() {
    echo -e "${BLUE}[6/10] 优化SSH配置...${NC}"
    
    if [ ! -f /etc/ssh/sshd_config ]; then
        echo "   ℹ️  SSH未安装，跳过配置"
        return 0
    fi
    
    # 备份原配置
    cp /etc/ssh/sshd_config /etc/ssh/sshd_config.bak.$(date +%Y%m%d)
    
    # 优化配置参数
    local ssh_params=(
        "UseDNS no"
        "GSSAPIAuthentication no"
        "ClientAliveInterval 60"
        "ClientAliveCountMax 3"
        "MaxAuthTries 3"
        "LoginGraceTime 60"
        "PermitRootLogin yes"
        "PasswordAuthentication yes"
        "PubkeyAuthentication yes"
        "X11Forwarding no"
        "PrintMotd no"
        "PrintLastLog yes"
        "TCPKeepAlive yes"
        "AllowTcpForwarding yes"
        "Compression delayed"
    )
    
    # 应用配置
    for param in "${ssh_params[@]}"; do
        local key=$(echo "$param" | awk ''{print $1}'')
        local value=$(echo "$param" | awk ''{$1=""; print substr($0,2)}'')
        
        if grep -q "^$key" /etc/ssh/sshd_config; then
            sed -i "s/^$key.*/$key $value/" /etc/ssh/sshd_config
        elif grep -q "^#$key" /etc/ssh/sshd_config; then
            sed -i "s/^#$key.*/$key $value/" /etc/ssh/sshd_config
        else
            echo "$key $value" >> /etc/ssh/sshd_config
        fi
    done
    
    # 修改端口（如果指定）
    if [ -n "$SSH_PORT" ] && [ "$SSH_PORT" != "22" ]; then
        if grep -q "^Port" /etc/ssh/sshd_config; then
            sed -i "s/^Port.*/Port $SSH_PORT/" /etc/ssh/sshd_config
        elif grep -q "^#Port" /etc/ssh/sshd_config; then
            sed -i "s/^#Port.*/Port $SSH_PORT/" /etc/ssh/sshd_config
        else
            echo "Port $SSH_PORT" >> /etc/ssh/sshd_config
        fi
        echo "   🔧 SSH端口已修改为: $SSH_PORT"
    fi
    
    # 重启SSH服务
    if systemctl is-active sshd >/dev/null 2>&1; then
        systemctl restart sshd
        echo "   🔄 SSH服务已重启"
    elif systemctl is-active ssh >/dev/null 2>&1; then
        systemctl restart ssh
        echo "   🔄 SSH服务已重启"
    fi
    
    echo "✅ SSH优化完成"
    echo ""
}

# ==============================================
# Swap优化
# ==============================================

optimize_swap() {
    echo -e "${BLUE}[7/10] 优化Swap配置...${NC}"
    
    if [ "$ENABLE_SWAP" != "yes" ]; then
        echo "   ℹ️  Swap优化已跳过 (ENABLE_SWAP=no)"
        return 0
    fi
    
    # 检查现有swap
    local current_swap=$(free -g | awk ''/^Swap:/{print $2}'')
    
    if [ "$current_swap" -eq 0 ]; then
        echo "   ℹ️  系统未配置Swap，正在创建..."
        
        # 创建swap文件 (大小为内存的2倍，最大不超过32GB)
        local swap_size=$((MEM_TOTAL * 2))
        [ $swap_size -gt 32 ] && swap_size=32
        
        dd if=/dev/zero of=/swapfile bs=1G count=$swap_size status=progress
        chmod 600 /swapfile
        mkswap /swapfile
        swapon /swapfile
        
        # 添加到fstab
        echo "/swapfile none swap sw 0 0" >> /etc/fstab
        
        echo "   ✅ 已创建 ${swap_size}GB 的Swap文件"
    else
        echo "   ℹ️  系统已有 ${current_swap}GB Swap"
    fi
    
    # 设置swappiness
    if [ -f /proc/sys/vm/swappiness ]; then
        echo "$SWAPPINESS" > /proc/sys/vm/swappiness
    fi
    
    echo "✅ Swap优化完成"
    echo ""
}

# ==============================================
# 安装常用工具
# ==============================================

install_system_tools() {
    echo -e "${BLUE}[8/10] 安装系统工具...${NC}"
    
    if [ "$INSTALL_PACKAGES" != "yes" ]; then
        echo "   ℹ️  包安装已跳过 (INSTALL_PACKAGES=no)"
        return 0
    fi
    
    # 基本工具列表
    local base_tools="vim git wget curl net-tools lsof htop tree unzip tar ncdu psmisc jq"
    
    case $OS_ID in
        centos|rhel|fedora|alinux|dragon|longxi|openeuler|euleros)
            echo "   📦 使用yum安装..."
            yum update -y >/dev/null 2>&1
            
            # 安装EPEL仓库
            if ! yum repolist enabled | grep -q "epel"; then
                yum install -y epel-release >/dev/null 2>&1 || \
                    rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-$(rpm -E %rhel).noarch.rpm 2>/dev/null || true
            fi
            
            # 安装开发工具组
            yum groupinstall -y "Development Tools" >/dev/null 2>&1 || true
            
            # 安装基础工具
            yum install -y $base_tools sysstat iotop iftop glances tmux >/dev/null 2>&1 || \
                echo "   ⚠️  部分包安装失败，继续..."
            ;;
            
        ubuntu|debian|uos|uos-server)
            echo "   📦 使用apt安装..."
            export DEBIAN_FRONTEND=noninteractive
            
            apt update >/dev/null 2>&1
            
            # 安装基础工具
            apt install -y $base_tools sysstat iotop iftop glances tmux >/dev/null 2>&1 || \
                echo "   ⚠️  部分包安装失败，继续..."
                
            # 清理
            apt autoremove -y >/dev/null 2>&1
            ;;
            
        *)
            echo "   ⚠️  未知系统类型，尝试通用安装..."
            if command -v yum >/dev/null 2>&1; then
                yum install -y vim wget curl net-tools >/dev/null 2>&1 || true
            elif command -v apt-get >/dev/null 2>&1; then
                apt-get update && apt-get install -y vim wget curl net-tools >/dev/null 2>&1 || true
            else
                echo "   ❌ 无法识别包管理器"
            fi
            ;;
    esac
    
    # 安装docker-compose（可选）
    if command -v docker >/dev/null 2>&1 && [ ! -f /usr/local/bin/docker-compose ]; then
        echo "   📦 安装docker-compose..."
        curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
            -o /usr/local/bin/docker-compose 2>/dev/null && \
            chmod +x /usr/local/bin/docker-compose || \
            echo "   ⚠️  docker-compose安装失败"
    fi
    
    echo "✅ 系统工具安装完成"
    echo ""
}

# ==============================================
# 时区和时间同步
# ==============================================

configure_time() {
    echo -e "${BLUE}[9/10] 配置时区和时间同步...${NC}"
    
    # 设置时区为上海
    timedatectl set-timezone Asia/Shanghai 2>/dev/null || \
        ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 2>/dev/null || true
    
    # 安装配置chrony/NTP
    case $OS_ID in
        centos|rhel|fedora|alinux|dragon|longxi|openeuler|euleros)
            yum install -y chrony >/dev/null 2>&1 || true
            systemctl enable chronyd 2>/dev/null && systemctl start chronyd 2>/dev/null || true
            ;;
        ubuntu|debian|uos|uos-server)
            apt install -y chrony >/dev/null 2>&1 || true
            systemctl enable chrony 2>/dev/null && systemctl start chrony 2>/dev/null || true
            ;;
    esac
    
    # 使用阿里云NTP服务器
    if [ -f /etc/chrony.conf ]; then
        cp /etc/chrony.conf /etc/chrony.conf.bak
        cat > /etc/chrony.conf << EOF
server ntp.aliyun.com iburst
server ntp1.aliyun.com iburst
server ntp2.aliyun.com iburst
server ntp3.aliyun.com iburst
stratumweight 0
driftfile /var/lib/chrony/drift
rtcsync
makestep 10 3
bindcmdaddress 127.0.0.1
bindcmdaddress ::1
keyfile /etc/chrony.keys
commandkey 1
generatecommandkey
noclientlog
logchange 0.5
logdir /var/log/chrony
EOF
        systemctl restart chronyd 2>/dev/null || true
    fi
    
    # 强制同步时间
    hwclock --systohc 2>/dev/null || true
    
    echo "✅ 时区和时间配置完成: $(date)"
    echo ""
}

# ==============================================
# 完成和验证
# ==============================================

finalize_and_verify() {
    echo -e "${BLUE}[10/10] 验证配置...${NC}"
    
    echo "📊 配置验证报告:"
    echo "----------------------------------------"
    
    # 1. 验证内核参数
    echo "1. 内核参数:"
    local kernel_params=(
        "net.core.somaxconn"
        "net.ipv4.tcp_max_syn_backlog"
        "vm.swappiness"
        "vm.max_map_count"
    )
    for param in "${kernel_params[@]}"; do
        local value=$(sysctl -n $param 2>/dev/null || echo "N/A")
        echo "   - $param = $value"
    done
    
    # 2. 验证文件限制
    echo ""
    echo "2. 资源限制:"
    local ulimit_nofile=$(ulimit -n 2>/dev/null || echo "N/A")
    echo "   - 文件描述符限制: $ulimit_nofile"
    
    # 3. 验证服务状态
    echo ""
    echo "3. 服务状态:"
    if systemctl is-active sshd >/dev/null 2>&1; then
        echo "   - SSH服务: 运行中"
    elif systemctl is-active ssh >/dev/null 2>&1; then
        echo "   - SSH服务: 运行中"
    else
        echo "   - SSH服务: 未运行"
    fi
    
    # 4. 验证SELinux
    echo ""
    echo "4. SELinux状态:"
    if command -v getenforce >/dev/null 2>&1; then
        echo "   - 当前: $(getenforce 2>/dev/null || echo ''Disabled'')"
    else
        echo "   - 未安装"
    fi
    
    # 5. 验证时间
    echo ""
    echo "5. 时间配置:"
    echo "   - 当前时间: $(date)"
    echo "   - 时区: $(timedatectl show --property=Timezone --value 2>/dev/null || echo ''Asia/Shanghai'')"
    
    # 6. 验证Swap
    echo ""
    echo "6. 内存和Swap:"
    free -h | awk ''NR<=2 {print "   - "$0}''
    
    echo "----------------------------------------"
    echo ""
    
    # 生成摘要报告
    cat > "$BACKUP_DIR/init_summary.txt" << EOF
系统初始化完成报告
====================
完成时间: $(date)
系统: $OS_NAME $OS_VERSION
架构: $ARCH
内存: ${MEM_TOTAL}GB

应用配置:
- 文件描述符限制: $FILE_LIMIT
- SELinux状态: $( [ "$DISABLE_SELINUX" = "yes" ] && echo "已禁用" || echo "保持原样" )
- IPv6状态: $( [ "$DISABLE_IPV6" = "yes" ] && echo "已禁用" || echo "启用" )
- Swap优化: $( [ "$ENABLE_SWAP" = "yes" ] && echo "已启用" || echo "未启用" )
- 包安装: $( [ "$INSTALL_PACKAGES" = "yes" ] && echo "已执行" || echo "已跳过" )

重要文件:
- 备份目录: $BACKUP_DIR
- 内核配置: /etc/sysctl.d/99-system-optimize.conf
- SSH配置: /etc/ssh/sshd_config
- 日志文件: $LOG_FILE

注意: 部分配置需要重启系统才能完全生效
重启命令: reboot
EOF
    
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}         系统初始化完成！              ${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo -e "${YELLOW}📋 重要信息:${NC}"
    echo "1. 备份文件: $BACKUP_DIR"
    echo "2. 详细日志: $LOG_FILE"
    echo "3. 摘要报告: $BACKUP_DIR/init_summary.txt"
    echo ""
    echo -e "${YELLOW}⚠️  注意事项:${NC}"
    echo "- 部分内核参数和SELinux设置需要重启才能生效"
    echo "- 建议检查防火墙配置"
    echo "- 建议修改SSH端口和加强安全设置"
    echo ""
    echo -e "${RED}🔄 如需应用所有更改，请重启系统:${NC}"
    echo -e "${RED}   sudo reboot${NC}"
    echo ""
    echo "当前时间: $(date ''+%Y-%m-%d %H:%M:%S'')"
    echo "运行时长: $SECONDS 秒"
}

# ==============================================
# 主程序
# ==============================================

main() {
    # 记录开始时间
    local start_time=$(date +%s)
    
    # 初始化
    check_root
    init_logging
    detect_system
    
    # 执行所有配置步骤
    backup_configs
    optimize_kernel
    set_resource_limits
    configure_selinux
    optimize_ssh
    optimize_swap
    install_system_tools
    configure_time
    finalize_and_verify
    
    # 计算总耗时
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    echo -e "${GREEN}🎉 所有任务完成！总耗时: ${duration}秒${NC}"
}

# 设置错误处理
trap ''echo -e "${RED}❌ 脚本执行中断${NC}"; exit 1'' INT TERM

# 运行主程序
main "$@"','2025-12-01 15:44:12','2025-12-01 15:44:12','admin','','123');
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (15,'test456',1,'#!/bin/bash

# ==============================================
# 系统初始化配置脚本 - 完全非交互式版本
# 支持：CentOS/RHEL/龙溪/UOS/麒麟/欧拉/Ubuntu/Debian等
# 用法：curl -sSL https://raw.githubusercontent.com/xxx/init.sh | bash
#       或 wget -qO- https://raw.githubusercontent.com/xxx/init.sh | bash
# ==============================================

set -e  # 遇到错误退出

# ==============================================
# 配置变量（可根据需要修改）
# ==============================================

# 基本配置
DISABLE_SELINUX="yes"           # 是否禁用SELinux: yes/no
DISABLE_IPV6="yes"              # 是否禁用IPv6: yes/no
INSTALL_PACKAGES="yes"         # 是否安装常用工具: yes/no
ENABLE_SWAP="yes"              # 是否启用swap优化: yes/no
SSH_PORT=""                    # SSH端口，留空保持默认22

# 性能参数
FILE_LIMIT="65535"            # 文件描述符限制
TCP_MAX_TW_BUCKETS="8000"     # TIME_WAIT连接数
TCP_TW_REUSE="1"              # 重用TIME_WAIT连接
SWAPPINESS="10"               # swappiness值

# 日志配置
LOG_DIR="/var/log/system_init"
BACKUP_DIR="/backup/system_config_$(date +%Y%m%d)"

# ==============================================
# 初始化设置
# ==============================================

# 颜色输出
RED=''\033[0;31m''
GREEN=''\033[0;32m''
YELLOW=''\033[1;33m''
BLUE=''\033[0;34m''
NC=''\033[0m''

# 日志初始化
init_logging() {
    mkdir -p "$LOG_DIR"
    LOG_FILE="$LOG_DIR/init_$(date +%Y%m%d_%H%M%S).log"
    exec > >(tee -a "$LOG_FILE") 2>&1
    
    echo "========================================="
    echo "系统初始化脚本 - 非交互式版本"
    echo "开始时间: $(date ''+%Y-%m-%d %H:%M:%S'')"
    echo "日志文件: $LOG_FILE"
    echo "========================================="
    echo ""
}

# 检查root权限
check_root() {
    if [ "$EUID" -ne 0 ]; then 
        echo -e "${RED}错误: 请使用root权限运行此脚本${NC}"
        echo "用法: sudo $0 或 curl -sSL URL | sudo bash"
        exit 1
    fi
}

# 系统检测
detect_system() {
    echo -e "${BLUE}[1/10] 检测系统信息...${NC}"
    
    # 检测系统
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_NAME="$NAME"
        OS_ID="$ID"
        OS_VERSION="$VERSION_ID"
        OS_ID_LIKE="$ID_LIKE"
    elif [ -f /etc/redhat-release ]; then
        OS_NAME=$(cat /etc/redhat-release)
        OS_ID="rhel"
        OS_VERSION=$(cat /etc/redhat-release | grep -oE ''[0-9]+\.[0-9]+'' || echo "unknown")
    else
        OS_NAME="未知"
        OS_ID="unknown"
        OS_VERSION="unknown"
    fi
    
    # 架构和内核
    ARCH=$(uname -m)
    KERNEL=$(uname -r)
    
    # 内存和CPU
    MEM_TOTAL=$(free -g | awk ''/^Mem:/{print $2}'' || echo "0")
    CPU_CORES=$(nproc 2>/dev/null || grep -c ^processor /proc/cpuinfo || echo "1")
    
    # 存储
    DISK_TOTAL=$(df -h / | awk ''NR==2 {print $2}'' || echo "unknown")
    
    echo "✅ 系统检测完成:"
    echo "   - 系统: $OS_NAME"
    echo "   - 版本: $OS_VERSION"
    echo "   - 架构: $ARCH"
    echo "   - 内核: $KERNEL"
    echo "   - CPU核心: $CPU_CORES"
    echo "   - 内存: ${MEM_TOTAL}GB"
    echo "   - 根分区: $DISK_TOTAL"
    echo ""
    
    export OS_ID OS_VERSION ARCH MEM_TOTAL CPU_CORES
}

# ==============================================
# 备份功能
# ==============================================

backup_configs() {
    echo -e "${BLUE}[2/10] 备份系统配置...${NC}"
    
    mkdir -p "$BACKUP_DIR"
    
    # 要备份的配置文件列表
    local configs=(
        "/etc/sysctl.conf"
        "/etc/security/limits.conf"
        "/etc/ssh/sshd_config"
        "/etc/selinux/config"
        "/etc/hosts"
        "/etc/resolv.conf"
        "/etc/fstab"
        "/etc/profile"
        "/etc/bashrc"
    )
    
    local backup_count=0
    for config in "${configs[@]}"; do
        if [ -f "$config" ]; then
            cp -p "$config" "$BACKUP_DIR/" 2>/dev/null && {
                echo "   📁 已备份: $config"
                ((backup_count++))
            } || echo "   ⚠️  备份失败: $config"
        fi
    done
    
    # 备份已加载的内核参数
    sysctl -a 2>/dev/null > "$BACKUP_DIR/sysctl_current.conf" && \
        echo "   📁 已备份当前内核参数"
    
    echo "✅ 备份完成: 共备份 $backup_count 个文件到 $BACKUP_DIR"
    echo ""
}

# ==============================================
# 内核优化
# ==============================================

optimize_kernel() {
    echo -e "${BLUE}[3/10] 优化内核参数...${NC}"
    
    # 创建优化配置文件
    local sysctl_file="/etc/sysctl.d/99-system-optimize.conf"
    
    cat > "$sysctl_file" << EOF
# ==============================================
# 系统内核优化配置
# 生成时间: $(date)
# 系统: $OS_NAME $OS_VERSION
# 内存: ${MEM_TOTAL}GB
# ==============================================

# --- 内存优化 ---
# swappiness (0-100, 值越低越少使用swap)
vm.swappiness = $SWAPPINESS

# 缓存压力 (1-100, 值越低越倾向缓存)
vm.vfs_cache_pressure = 50

# 脏页写回比例
vm.dirty_ratio = 10
vm.dirty_background_ratio = 5

# 内存分配策略 (1允许超额分配)
vm.overcommit_memory = 1
vm.overcommit_ratio = 90

# 虚拟内存映射数量 (Elasticsearch等需要)
vm.max_map_count = 262144

# --- 网络优化 ---
# TCP连接队列
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 65535
net.ipv4.tcp_max_syn_backlog = 65536

# TIME_WAIT优化
net.ipv4.tcp_max_tw_buckets = $TCP_MAX_TW_BUCKETS
net.ipv4.tcp_tw_reuse = $TCP_TW_REUSE
net.ipv4.tcp_fin_timeout = 30

# 端口范围
net.ipv4.ip_local_port_range = 1024 65535

# Keepalive设置
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_keepalive_probes = 5
net.ipv4.tcp_keepalive_intvl = 15

# 缓冲区大小
net.core.rmem_default = 262144
net.core.rmem_max = 16777216
net.core.wmem_default = 262144
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216

# TCP拥塞控制
net.ipv4.tcp_congestion_control = cubic

# 其他优化
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_synack_retries = 2
net.ipv4.tcp_syn_retries = 3
net.ipv4.tcp_max_orphans = 65536
net.ipv4.tcp_timestamps = 1
net.ipv4.tcp_sack = 1
net.ipv4.tcp_window_scaling = 1
net.ipv4.tcp_fastopen = 3

# 根据内存大小调整TCP内存
EOF
    
    # 根据内存调整TCP内存参数
    if [ "$MEM_TOTAL" -lt 4 ]; then
        echo "net.ipv4.tcp_mem = 98304 131072 196608" >> "$sysctl_file"
    elif [ "$MEM_TOTAL" -lt 16 ]; then
        echo "net.ipv4.tcp_mem = 196608 262144 393216" >> "$sysctl_file"
    else
        echo "net.ipv4.tcp_mem = 393216 524288 786432" >> "$sysctl_file"
    fi
    
    # IPv6配置
    if [ "$DISABLE_IPV6" = "yes" ]; then
        cat >> "$sysctl_file" << EOF

# --- IPv6禁用 ---
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6 = 1
EOF
    fi
    
    # 应用配置
    sysctl -p "$sysctl_file" >/dev/null 2>&1
    echo "✅ 内核参数优化完成: $sysctl_file"
    echo ""
}

# ==============================================
# 资源限制设置
# ==============================================

set_resource_limits() {
    echo -e "${BLUE}[4/10] 设置系统资源限制...${NC}"
    
    # 备份原文件
    if [ -f /etc/security/limits.conf ]; then
        cp /etc/security/limits.conf /etc/security/limits.conf.bak.$(date +%Y%m%d)
    fi
    
    # 清理旧的自定义设置（如果有）
    sed -i ''/^# ===== 系统资源限制配置 =====/,/^# ===============================/d'' /etc/security/limits.conf 2>/dev/null
    
    # 添加新的限制
    cat >> /etc/security/limits.conf << EOF

# ===== 系统资源限制配置 =====
# 配置时间: $(date)
# 文件描述符限制: $FILE_LIMIT
# ===============================

# 全局文件描述符限制
* soft nofile $FILE_LIMIT
* hard nofile $FILE_LIMIT

# 进程数限制
* soft nproc $FILE_LIMIT
* hard nproc $FILE_LIMIT

# 核心文件大小
* soft core unlimited
* hard core unlimited

# 栈大小
* soft stack 10240
* hard stack 32768

# 实时优先级 (某些应用需要)
#@realtime soft rtprio 99
#@realtime hard rtprio 99
#@realtime soft memlock unlimited
#@realtime hard memlock unlimited

EOF
    
    # 对于systemd系统
    if [ -d /etc/systemd/system.conf.d ]; then
        mkdir -p /etc/systemd/system.conf.d
        cat > /etc/systemd/system.conf.d/limits.conf << EOF
[Manager]
DefaultLimitNOFILE=$FILE_LIMIT
DefaultLimitNPROC=$FILE_LIMIT
DefaultLimitMEMLOCK=infinity
EOF
    fi
    
    # 对于使用pam_limits的系统
    if [ -f /etc/pam.d/common-session ] && [ "$OS_ID" = "ubuntu" ] || [ "$OS_ID" = "debian" ]; then
        grep -q "pam_limits.so" /etc/pam.d/common-session || \
            echo "session required pam_limits.so" >> /etc/pam.d/common-session
    fi
    
    echo "✅ 资源限制设置完成: 文件限制=$FILE_LIMIT"
    echo ""
}

# ==============================================
# SELinux配置
# ==============================================

configure_selinux() {
    echo -e "${BLUE}[5/10] 配置SELinux...${NC}"
    
    # 检查SELinux是否安装
    if ! command -v getenforce >/dev/null 2>&1 && [ ! -f /etc/selinux/config ]; then
        echo "   ℹ️  SELinux未安装，跳过配置"
        return 0
    fi
    
    local current_state=$(getenforce 2>/dev/null || echo "Disabled")
    echo "   当前SELinux状态: $current_state"
    
    if [ "$DISABLE_SELINUX" = "yes" ]; then
        echo "   🔧 正在禁用SELinux..."
        
        # 临时设置为permissive
        setenforce 0 2>/dev/null || true
        
        # 永久禁用
        if [ -f /etc/selinux/config ]; then
            sed -i ''s/^SELINUX=.*/SELINUX=disabled/'' /etc/selinux/config
            sed -i ''s/^SELINUXTYPE=.*/#SELINUXTYPE=targeted/'' /etc/selinux/config
            echo "   ✅ SELinux已永久禁用，重启后生效"
        fi
    else
        echo "   ℹ️  保持SELinux当前状态 ($current_state)"
        echo "   提示: 如需禁用，设置 DISABLE_SELINUX=''yes''"
    fi
    
    echo ""
}

# ==============================================
# SSH优化
# ==============================================

optimize_ssh() {
    echo -e "${BLUE}[6/10] 优化SSH配置...${NC}"
    
    if [ ! -f /etc/ssh/sshd_config ]; then
        echo "   ℹ️  SSH未安装，跳过配置"
        return 0
    fi
    
    # 备份原配置
    cp /etc/ssh/sshd_config /etc/ssh/sshd_config.bak.$(date +%Y%m%d)
    
    # 优化配置参数
    local ssh_params=(
        "UseDNS no"
        "GSSAPIAuthentication no"
        "ClientAliveInterval 60"
        "ClientAliveCountMax 3"
        "MaxAuthTries 3"
        "LoginGraceTime 60"
        "PermitRootLogin yes"
        "PasswordAuthentication yes"
        "PubkeyAuthentication yes"
        "X11Forwarding no"
        "PrintMotd no"
        "PrintLastLog yes"
        "TCPKeepAlive yes"
        "AllowTcpForwarding yes"
        "Compression delayed"
    )
    
    # 应用配置
    for param in "${ssh_params[@]}"; do
        local key=$(echo "$param" | awk ''{print $1}'')
        local value=$(echo "$param" | awk ''{$1=""; print substr($0,2)}'')
        
        if grep -q "^$key" /etc/ssh/sshd_config; then
            sed -i "s/^$key.*/$key $value/" /etc/ssh/sshd_config
        elif grep -q "^#$key" /etc/ssh/sshd_config; then
            sed -i "s/^#$key.*/$key $value/" /etc/ssh/sshd_config
        else
            echo "$key $value" >> /etc/ssh/sshd_config
        fi
    done
    
    # 修改端口（如果指定）
    if [ -n "$SSH_PORT" ] && [ "$SSH_PORT" != "22" ]; then
        if grep -q "^Port" /etc/ssh/sshd_config; then
            sed -i "s/^Port.*/Port $SSH_PORT/" /etc/ssh/sshd_config
        elif grep -q "^#Port" /etc/ssh/sshd_config; then
            sed -i "s/^#Port.*/Port $SSH_PORT/" /etc/ssh/sshd_config
        else
            echo "Port $SSH_PORT" >> /etc/ssh/sshd_config
        fi
        echo "   🔧 SSH端口已修改为: $SSH_PORT"
    fi
    
    # 重启SSH服务
    if systemctl is-active sshd >/dev/null 2>&1; then
        systemctl restart sshd
        echo "   🔄 SSH服务已重启"
    elif systemctl is-active ssh >/dev/null 2>&1; then
        systemctl restart ssh
        echo "   🔄 SSH服务已重启"
    fi
    
    echo "✅ SSH优化完成"
    echo ""
}

# ==============================================
# Swap优化
# ==============================================

optimize_swap() {
    echo -e "${BLUE}[7/10] 优化Swap配置...${NC}"
    
    if [ "$ENABLE_SWAP" != "yes" ]; then
        echo "   ℹ️  Swap优化已跳过 (ENABLE_SWAP=no)"
        return 0
    fi
    
    # 检查现有swap
    local current_swap=$(free -g | awk ''/^Swap:/{print $2}'')
    
    if [ "$current_swap" -eq 0 ]; then
        echo "   ℹ️  系统未配置Swap，正在创建..."
        
        # 创建swap文件 (大小为内存的2倍，最大不超过32GB)
        local swap_size=$((MEM_TOTAL * 2))
        [ $swap_size -gt 32 ] && swap_size=32
        
        dd if=/dev/zero of=/swapfile bs=1G count=$swap_size status=progress
        chmod 600 /swapfile
        mkswap /swapfile
        swapon /swapfile
        
        # 添加到fstab
        echo "/swapfile none swap sw 0 0" >> /etc/fstab
        
        echo "   ✅ 已创建 ${swap_size}GB 的Swap文件"
    else
        echo "   ℹ️  系统已有 ${current_swap}GB Swap"
    fi
    
    # 设置swappiness
    if [ -f /proc/sys/vm/swappiness ]; then
        echo "$SWAPPINESS" > /proc/sys/vm/swappiness
    fi
    
    echo "✅ Swap优化完成"
    echo ""
}

# ==============================================
# 安装常用工具
# ==============================================

install_system_tools() {
    echo -e "${BLUE}[8/10] 安装系统工具...${NC}"
    
    if [ "$INSTALL_PACKAGES" != "yes" ]; then
        echo "   ℹ️  包安装已跳过 (INSTALL_PACKAGES=no)"
        return 0
    fi
    
    # 基本工具列表
    local base_tools="vim git wget curl net-tools lsof htop tree unzip tar ncdu psmisc jq"
    
    case $OS_ID in
        centos|rhel|fedora|alinux|dragon|longxi|openeuler|euleros)
            echo "   📦 使用yum安装..."
            yum update -y >/dev/null 2>&1
            
            # 安装EPEL仓库
            if ! yum repolist enabled | grep -q "epel"; then
                yum install -y epel-release >/dev/null 2>&1 || \
                    rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-$(rpm -E %rhel).noarch.rpm 2>/dev/null || true
            fi
            
            # 安装开发工具组
            yum groupinstall -y "Development Tools" >/dev/null 2>&1 || true
            
            # 安装基础工具
            yum install -y $base_tools sysstat iotop iftop glances tmux >/dev/null 2>&1 || \
                echo "   ⚠️  部分包安装失败，继续..."
            ;;
            
        ubuntu|debian|uos|uos-server)
            echo "   📦 使用apt安装..."
            export DEBIAN_FRONTEND=noninteractive
            
            apt update >/dev/null 2>&1
            
            # 安装基础工具
            apt install -y $base_tools sysstat iotop iftop glances tmux >/dev/null 2>&1 || \
                echo "   ⚠️  部分包安装失败，继续..."
                
            # 清理
            apt autoremove -y >/dev/null 2>&1
            ;;
            
        *)
            echo "   ⚠️  未知系统类型，尝试通用安装..."
            if command -v yum >/dev/null 2>&1; then
                yum install -y vim wget curl net-tools >/dev/null 2>&1 || true
            elif command -v apt-get >/dev/null 2>&1; then
                apt-get update && apt-get install -y vim wget curl net-tools >/dev/null 2>&1 || true
            else
                echo "   ❌ 无法识别包管理器"
            fi
            ;;
    esac
    
    # 安装docker-compose（可选）
    if command -v docker >/dev/null 2>&1 && [ ! -f /usr/local/bin/docker-compose ]; then
        echo "   📦 安装docker-compose..."
        curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
            -o /usr/local/bin/docker-compose 2>/dev/null && \
            chmod +x /usr/local/bin/docker-compose || \
            echo "   ⚠️  docker-compose安装失败"
    fi
    
    echo "✅ 系统工具安装完成"
    echo ""
}

# ==============================================
# 时区和时间同步
# ==============================================

configure_time() {
    echo -e "${BLUE}[9/10] 配置时区和时间同步...${NC}"
    
    # 设置时区为上海
    timedatectl set-timezone Asia/Shanghai 2>/dev/null || \
        ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 2>/dev/null || true
    
    # 安装配置chrony/NTP
    case $OS_ID in
        centos|rhel|fedora|alinux|dragon|longxi|openeuler|euleros)
            yum install -y chrony >/dev/null 2>&1 || true
            systemctl enable chronyd 2>/dev/null && systemctl start chronyd 2>/dev/null || true
            ;;
        ubuntu|debian|uos|uos-server)
            apt install -y chrony >/dev/null 2>&1 || true
            systemctl enable chrony 2>/dev/null && systemctl start chrony 2>/dev/null || true
            ;;
    esac
    
    # 使用阿里云NTP服务器
    if [ -f /etc/chrony.conf ]; then
        cp /etc/chrony.conf /etc/chrony.conf.bak
        cat > /etc/chrony.conf << EOF
server ntp.aliyun.com iburst
server ntp1.aliyun.com iburst
server ntp2.aliyun.com iburst
server ntp3.aliyun.com iburst
stratumweight 0
driftfile /var/lib/chrony/drift
rtcsync
makestep 10 3
bindcmdaddress 127.0.0.1
bindcmdaddress ::1
keyfile /etc/chrony.keys
commandkey 1
generatecommandkey
noclientlog
logchange 0.5
logdir /var/log/chrony
EOF
        systemctl restart chronyd 2>/dev/null || true
    fi
    
    # 强制同步时间
    hwclock --systohc 2>/dev/null || true
    
    echo "✅ 时区和时间配置完成: $(date)"
    echo ""
}

# ==============================================
# 完成和验证
# ==============================================

finalize_and_verify() {
    echo -e "${BLUE}[10/10] 验证配置...${NC}"
    
    echo "📊 配置验证报告:"
    echo "----------------------------------------"
    
    # 1. 验证内核参数
    echo "1. 内核参数:"
    local kernel_params=(
        "net.core.somaxconn"
        "net.ipv4.tcp_max_syn_backlog"
        "vm.swappiness"
        "vm.max_map_count"
    )
    for param in "${kernel_params[@]}"; do
        local value=$(sysctl -n $param 2>/dev/null || echo "N/A")
        echo "   - $param = $value"
    done
    
    # 2. 验证文件限制
    echo ""
    echo "2. 资源限制:"
    local ulimit_nofile=$(ulimit -n 2>/dev/null || echo "N/A")
    echo "   - 文件描述符限制: $ulimit_nofile"
    
    # 3. 验证服务状态
    echo ""
    echo "3. 服务状态:"
    if systemctl is-active sshd >/dev/null 2>&1; then
        echo "   - SSH服务: 运行中"
    elif systemctl is-active ssh >/dev/null 2>&1; then
        echo "   - SSH服务: 运行中"
    else
        echo "   - SSH服务: 未运行"
    fi
    
    # 4. 验证SELinux
    echo ""
    echo "4. SELinux状态:"
    if command -v getenforce >/dev/null 2>&1; then
        echo "   - 当前: $(getenforce 2>/dev/null || echo ''Disabled'')"
    else
        echo "   - 未安装"
    fi
    
    # 5. 验证时间
    echo ""
    echo "5. 时间配置:"
    echo "   - 当前时间: $(date)"
    echo "   - 时区: $(timedatectl show --property=Timezone --value 2>/dev/null || echo ''Asia/Shanghai'')"
    
    # 6. 验证Swap
    echo ""
    echo "6. 内存和Swap:"
    free -h | awk ''NR<=2 {print "   - "$0}''
    
    echo "----------------------------------------"
    echo ""
    
    # 生成摘要报告
    cat > "$BACKUP_DIR/init_summary.txt" << EOF
系统初始化完成报告
====================
完成时间: $(date)
系统: $OS_NAME $OS_VERSION
架构: $ARCH
内存: ${MEM_TOTAL}GB

应用配置:
- 文件描述符限制: $FILE_LIMIT
- SELinux状态: $( [ "$DISABLE_SELINUX" = "yes" ] && echo "已禁用" || echo "保持原样" )
- IPv6状态: $( [ "$DISABLE_IPV6" = "yes" ] && echo "已禁用" || echo "启用" )
- Swap优化: $( [ "$ENABLE_SWAP" = "yes" ] && echo "已启用" || echo "未启用" )
- 包安装: $( [ "$INSTALL_PACKAGES" = "yes" ] && echo "已执行" || echo "已跳过" )

重要文件:
- 备份目录: $BACKUP_DIR
- 内核配置: /etc/sysctl.d/99-system-optimize.conf
- SSH配置: /etc/ssh/sshd_config
- 日志文件: $LOG_FILE

注意: 部分配置需要重启系统才能完全生效
重启命令: reboot
EOF
    
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}         系统初始化完成！              ${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo -e "${YELLOW}📋 重要信息:${NC}"
    echo "1. 备份文件: $BACKUP_DIR"
    echo "2. 详细日志: $LOG_FILE"
    echo "3. 摘要报告: $BACKUP_DIR/init_summary.txt"
    echo ""
    echo -e "${YELLOW}⚠️  注意事项:${NC}"
    echo "- 部分内核参数和SELinux设置需要重启才能生效"
    echo "- 建议检查防火墙配置"
    echo "- 建议修改SSH端口和加强安全设置"
    echo ""
    echo -e "${RED}🔄 如需应用所有更改，请重启系统:${NC}"
    echo -e "${RED}   sudo reboot${NC}"
    echo ""
    echo "当前时间: $(date ''+%Y-%m-%d %H:%M:%S'')"
    echo "运行时长: $SECONDS 秒"
}

# ==============================================
# 主程序
# ==============================================

main() {
    # 记录开始时间
    local start_time=$(date +%s)
    
    # 初始化
    check_root
    init_logging
    detect_system
    
    # 执行所有配置步骤
    backup_configs
    optimize_kernel
    set_resource_limits
    configure_selinux
    optimize_ssh
    optimize_swap
    install_system_tools
    configure_time
    finalize_and_verify
    
    # 计算总耗时
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    echo -e "${GREEN}🎉 所有任务完成！总耗时: ${duration}秒${NC}"
}

# 设置错误处理
trap ''echo -e "${RED}❌ 脚本执行中断${NC}"; exit 1'' INT TERM

# 运行主程序
main "$@"','2025-12-01 16:56:29','2025-12-01 16:56:29','admin','','123');
INSERT INTO `task_template` (`id`,`name`,`type`,`content`,`created_at`,`updated_at`,`created_by`,`updated_by`,`remark`) VALUES (16,'ls',1,'ls','2025-12-05 20:54:54','2025-12-05 20:54:54','admin','','123');

-- Table structure for task_work
CREATE TABLE `task_work` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `task_id` bigint(20) COMMENT '关联的任务ID',
 `template_id` bigint(20) COMMENT '任务模板ID',
 `host_id` bigint(20) COMMENT '执行主机ID',
 `status` bigint(19) COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常',
 `log` text COMMENT '任务日志',
 `log_path` text COMMENT '日志文件路径',
 `start_time` datetime COMMENT '任务开始时间',
 `end_time` datetime COMMENT '任务结束时间',
 `duration` bigint(19) COMMENT '执行耗时(秒)',
 `created_at` datetime COMMENT '创建时间',
 `type` bigint(19) COMMENT '任务类型 1=普通任务,2=定时任务',
 `scheduled_time` datetime COMMENT '定时任务执行时间',
 `cron_expr` longtext COMMENT 'cron表达式',
 `is_recurring` tinyint(3) COMMENT '是否周期性任务', 
PRIMARY KEY (id)
);
ALTER TABLE `task_work` ADD  INDEX `idx_task_work_task_id`(`task_id`) USING BTREE;
ALTER TABLE `task_work` ADD  INDEX `idx_task_work_template_id`(`template_id`) USING BTREE;

-- Data for table task_work
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (90,48,11,1,2,'[2025-09-29-21:42:02] 任务开始
进程统计信息如下
总进程数量为:79
Running 进程数为:1
Stoped 进程数为:0
Sleeping 进程数为:55
Zombie 进程数为:0
[2025-09-29-21:42:02] 任务完成
','logs/task_48/task_48_template_11.log','2025-09-29 21:42:00','2025-09-29 21:42:06',5,'2025-09-29 20:58:17',0,'2025-09-29 21:45:00',NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (91,48,0,0,4,'','',NULL,NULL,0,'2025-09-29 20:58:17',2,'2025-09-29 21:00:00',NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (92,49,12,1,3,'[2025-09-29-21:15:27] 任务开始
协议:tcp 端口号:6379
协议:tcp 端口号:4330
协议:tcp 端口号:3306
协议:tcp 端口号:80
协议:tcp 端口号:33060
协议:tcp 端口号:8088
协议:tcp 端口号:8080
协议:tcp 端口号:4330
协议:tcp 端口号:80
协议:tcp 端口号:8080
协议:tcp 端口号:8086
协议:tcp 端口号:9091
协议:tcp 端口号:9090
[2025-09-29-21:15:27] 任务完成
','logs/task_49/task_49_template_12.log','2025-09-29 21:15:25','2025-09-29 21:15:31',5,'2025-09-29 21:01:51',0,NULL,NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (104,56,11,1,1,'[2025-10-01-14:03:02] 任务开始
进程统计信息如下
总进程数量为:80
Running 进程数为:0
Stoped 进程数为:0
Sleeping 进程数为:56
Zombie 进程数为:0
[2025-10-01-14:03:03] 任务完成
','logs/task_56/task_56_template_11.log','2025-10-01 14:03:00','2025-10-01 14:03:06',5,'2025-09-30 17:15:17',0,'2025-09-30 17:18:00',NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (105,57,12,1,1,'[2025-10-01-14:00:03] 任务开始
协议:tcp 端口号:6379
协议:tcp 端口号:4330
协议:tcp 端口号:3306
协议:tcp 端口号:80
协议:tcp 端口号:33060
协议:tcp 端口号:8088
协议:tcp 端口号:8080
协议:tcp 端口号:4330
协议:tcp 端口号:80
协议:tcp 端口号:8080
协议:tcp 端口号:8086
协议:tcp 端口号:9091
协议:tcp 端口号:9090
[2025-10-01-14:00:03] 任务完成
','logs/task_57/task_57_template_12.log','2025-10-01 14:00:00','2025-10-01 14:00:07',6,'2025-09-30 17:17:02',0,'2025-09-30 17:20:00',NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (106,58,10,1,3,'[2025-10-01-14:53:41] 任务开始
 本 机 MAC 地 址 信 息 如 下 :
veth47f9cf4@if2:
86:22:92:32:6b:13
[2025-10-01-14:53:41] 任务完成
','logs/task_58/task_58_template_10.log','2025-10-01 14:53:39','2025-10-01 14:53:45',6,'2025-10-01 14:53:31',0,NULL,NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (107,59,2,501,3,'[2025-11-06-15:50:52] 任务开始
[15:50:52] 1
[15:50:53] 2
[15:50:54] 3
[15:50:55] 4
[15:50:56] 5
[15:50:57] 6
[15:50:58] 7
[15:50:59] 8
[15:51:00] 9
[15:51:01] 10
[15:51:02] 11
[15:51:03] 12
[15:51:04] 13
[15:51:05] 14
[15:51:06] 15
[15:51:07] 16
[15:51:08] 17
[15:51:09] 18
[15:51:10] 19
[15:51:11] 20
[15:51:12] 21
[15:51:13] 22
[15:51:14] 23
[15:51:15] 24
[15:51:16] 25
[15:51:17] 26
[15:51:18] 27
[15:51:19] 28
[15:51:20] 29
[15:51:21] 30
[15:51:22] 31
[15:51:23] 32
[15:51:24] 33
[15:51:25] 34
[15:51:26] 35
[15:51:27] 36
[15:51:28] 37
[15:51:29] 38
[15:51:30] 39
[15:51:31] 40
[15:51:32] 41
[15:51:33] 42
[15:51:34] 43
[15:51:35] 44
[15:51:36] 45
[15:51:37] 46
[15:51:38] 47
[15:51:39] 48
[15:51:40] 49
[15:51:41] 50
[15:51:42] 51
[15:51:43] 52
[15:51:44] 53
[15:51:45] 54
[15:51:46] 55
[15:51:47] 56
[15:51:48] 57
[15:51:49] 58
[15:51:50] 59
[15:51:51] 60
[15:51:52] 61
[15:51:53] 62
[15:51:54] 63
[15:51:55] 64
[15:51:56] 65
[15:51:57] 66
[15:51:58] 67
[15:51:59] 68
[15:52:00] 69
[15:52:01] 70
[15:52:02] 71
[15:52:03] 72
[15:52:04] 73
[15:52:05] 74
[15:52:06] 75
[15:52:07] 76
[15:52:08] 77
[15:52:09] 78
[15:52:10] 79
[15:52:11] 80
[15:52:12] 81
[15:52:13] 82
[15:52:14] 83
[15:52:15] 84
[15:52:16] 85
[15:52:17] 86
[15:52:18] 87
[15:52:19] 88
[15:52:20] 89
[15:52:21] 90
[15:52:22] 91
[15:52:23] 92
[15:52:24] 93
[15:52:25] 94
[15:52:26] 95
[15:52:27] 96
[15:52:28] 97
[15:52:29] 98
[15:52:30] 99
[15:52:31] 100
完成：所有数字 1-100 已打印完毕。
[2025-11-06-15:52:32] 任务完成
','logs/task_59/task_59_template_2.log','2025-11-06 15:50:49','2025-11-06 15:52:35',106,'2025-11-06 15:50:35',0,NULL,NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (108,60,3,501,1,'','',NULL,NULL,0,'2025-11-06 15:51:56',0,NULL,NULL,NULL);
INSERT INTO `task_work` (`id`,`task_id`,`template_id`,`host_id`,`status`,`log`,`log_path`,`start_time`,`end_time`,`duration`,`created_at`,`type`,`scheduled_time`,`cron_expr`,`is_recurring`) VALUES (109,61,2,511,3,'[2025-09-26-22:13:03] 任务开始
[22:13:03] 1
[22:13:04] 2
[22:13:05] 3
[22:13:06] 4
[22:13:07] 5
[22:13:08] 6
[22:13:09] 7
[22:13:10] 8
[22:13:11] 9
[22:13:12] 10
[22:13:13] 11
[22:13:14] 12
[22:13:15] 13
[22:13:16] 14
[22:13:17] 15
[22:13:18] 16
[22:13:19] 17
[22:13:20] 18
[22:13:21] 19
[22:13:22] 20
[22:13:23] 21
[22:13:24] 22
[22:13:25] 23
[22:13:26] 24
[22:13:27] 25
[22:13:28] 26
[22:13:29] 27
[22:13:30] 28
[22:13:31] 29
[22:13:32] 30
[22:13:33] 31
[22:13:34] 32
[22:13:35] 33
[22:13:36] 34
[22:13:37] 35
[22:13:38] 36
[22:13:39] 37
[22:13:40] 38
[22:13:41] 39
[22:13:42] 40
[22:13:43] 41
[22:13:44] 42
[22:13:45] 43
[22:13:46] 44
[22:13:47] 45
[22:13:48] 46
[22:13:49] 47
[22:13:50] 48
[22:13:51] 49
[22:13:52] 50
[22:13:53] 51
[22:13:54] 52
[22:13:55] 53
[22:13:56] 54
[22:13:57] 55
[22:13:58] 56
[22:13:59] 57
[22:14:00] 58
[22:14:01] 59
[22:14:02] 60
[22:14:03] 61
[22:14:04] 62
[22:14:05] 63
[22:14:06] 64
[22:14:07] 65
[22:14:09] 66
[22:14:10] 67
[22:14:11] 68
[22:14:12] 69
[22:14:13] 70
[22:14:14] 71
[22:14:15] 72
[22:14:16] 73
[22:14:17] 74
[22:14:18] 75
[22:14:19] 76
[22:14:20] 77
[22:14:21] 78
[22:14:22] 79
[22:14:23] 80
[22:14:24] 81
[22:14:25] 82
[22:14:26] 83
[22:14:27] 84
[22:14:28] 85
[22:14:29] 86
[22:14:30] 87
[22:14:31] 88
[22:14:32] 89
[22:14:33] 90
[22:14:34] 91
[22:14:35] 92
[22:14:36] 93
[22:14:37] 94
[22:14:38] 95
[22:14:39] 96
[22:14:40] 97
[22:14:41] 98
[22:14:42] 99
[22:14:43] 100
完成：所有数字 1-100 已打印完毕。
[2025-09-26-22:14:44] 任务完成
','logs/task_61/task_61_template_2.log','2025-11-24 00:29:05','2025-11-24 00:30:40',95,'2025-11-24 00:28:34',0,NULL,NULL,NULL);

-- Table structure for tool_link
CREATE TABLE `tool_link` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `title` longtext NOT NULL COMMENT '''导航标题''',
 `icon` longtext COMMENT '''导航图标''',
 `link` longtext NOT NULL COMMENT '''链接地址''',
 `sort` bigint(19) DEFAULT 0 COMMENT '''排序''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''', 
PRIMARY KEY (id)
);

-- Data for table tool_link
INSERT INTO `tool_link` (`id`,`title`,`icon`,`link`,`sort`,`create_time`,`update_time`) VALUES (2,'百度','http://10.7.16.22:8080/api/v1/upload/20251023/625775000.svg','https://www.baidu.com/',0,'2025-10-23 15:54:08','2025-10-23 15:54:08');
INSERT INTO `tool_link` (`id`,`title`,`icon`,`link`,`sort`,`create_time`,`update_time`) VALUES (3,'aws','http://10.7.16.22:8080/api/v1/upload/20251023/985806000.svg','https://us-west-2.console.aws.amazon.com/eks/clusters?region=us-west-2#',0,'2025-10-23 15:55:59','2025-10-23 15:55:59');
INSERT INTO `tool_link` (`id`,`title`,`icon`,`link`,`sort`,`create_time`,`update_time`) VALUES (4,'美女一号','http://10.7.16.22:8080/api/v1/upload/20251023/646700000.png','https://gitee.com/zhang_fan1024',0,'2025-10-23 15:56:22','2025-10-23 15:56:22');
INSERT INTO `tool_link` (`id`,`title`,`icon`,`link`,`sort`,`create_time`,`update_time`) VALUES (5,'美女二号','http://10.7.16.22:8080/api/v1/upload/20251023/733520000.png','https://demo.spug.cc/',0,'2025-10-23 15:56:48','2025-10-23 15:56:48');
INSERT INTO `tool_link` (`id`,`title`,`icon`,`link`,`sort`,`create_time`,`update_time`) VALUES (6,'凡人修仙传','http://10.7.16.22:8080/api/v1/upload/20251023/771236000.png','http://115.190.10.126/#/dashboard',0,'2025-10-23 15:57:21','2025-10-23 15:57:21');
INSERT INTO `tool_link` (`id`,`title`,`icon`,`link`,`sort`,`create_time`,`update_time`) VALUES (7,'腾讯','http://10.7.16.22:8080/api/v1/upload/20251205/92506000.png','https://cloud.tencent.com/login?s_url=https%3A%2F%2Fconsole.cloud.tencent.com%2Ftke2%2Fcluster%2Fsub%2Flist%2Fbasic%2Finfo%3Frid%3D1%26clusterId%3Dcls-301g0fi0',0,'2025-10-23 15:57:55','2025-12-05 20:12:08');

-- Table structure for tool_service_deploy
CREATE TABLE `tool_service_deploy` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '''主键''',
 `service_name` longtext NOT NULL COMMENT '''服务名称''',
 `service_id` longtext NOT NULL COMMENT '''服务ID''',
 `version` longtext NOT NULL COMMENT '''服务版本''',
 `host_id` bigint(20) NOT NULL COMMENT '''主机ID''',
 `host_ip` longtext NOT NULL COMMENT '''主机IP''',
 `install_dir` longtext NOT NULL COMMENT '''安装目录''',
 `container_name` longtext COMMENT '''容器名称''',
 `ports` longtext COMMENT '''端口映射(JSON)''',
 `env_vars` longtext COMMENT '''环境变量(JSON)''',
 `status` bigint(19) DEFAULT 0 COMMENT '''状态:0->部署中,1->运行中,2->已停止,3->部署失败''',
 `deploy_log` longtext COMMENT '''部署日志''',
 `create_time` datetime NOT NULL COMMENT '''创建时间''',
 `update_time` datetime COMMENT '''更新时间''', 
PRIMARY KEY (id)
);

-- Data for table tool_service_deploy
INSERT INTO `tool_service_deploy` (`id`,`service_name`,`service_id`,`version`,`host_id`,`host_ip`,`install_dir`,`container_name`,`ports`,`env_vars`,`status`,`deploy_log`,`create_time`,`update_time`) VALUES (7,'Redis','redis','redis-7.2',501,'8.130.14.34','/opt/data/redis','','','{"REDIS_MAXMEMORY":"2gb","REDIS_PASSWORD":"redis123456","REDIS_PORT":"6370"}',1,'[2025-10-30 17:20:01] 开始部署 Redis Redis 7.2
[2025-10-30 17:20:01] 连接主机 8.130.14.34...
[2025-10-30 17:20:02] SSH连接成功
[2025-10-30 17:20:03] 创建安装目录 /opt/data/redis...
[2025-10-30 17:20:03] 读取模板文件 common/templates/05-redis/versions/redis-7.2-docker-compose.yml...
[2025-10-30 17:20:03] 生成环境变量配置...
[2025-10-30 17:20:03] 上传 docker-compose.yml...
[2025-10-30 17:20:04] 上传 .env...
[2025-10-30 17:20:04] 检查Docker环境...
[2025-10-30 17:20:04] 启动服务容器...
[2025-10-30 17:20:05] 容器启动输出:

[2025-10-30 17:20:05] 验证容器状态...
[2025-10-30 17:20:08] 容器状态:
NAME      IMAGE                                                                                      COMMAND                  SERVICE   CREATED         STATUS                            PORTS
redis72   crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/redis:7.2-alpine   "docker-entrypoint.s…"   redis     4 seconds ago   Up 3 seconds (health: starting)   0.0.0.0:6370->6379/tcp, [::]:6370->6379/tcp

[2025-10-30 17:20:08] 部署完成！
','2025-10-30 17:20:01','2025-10-30 17:20:01');
INSERT INTO `tool_service_deploy` (`id`,`service_name`,`service_id`,`version`,`host_id`,`host_ip`,`install_dir`,`container_name`,`ports`,`env_vars`,`status`,`deploy_log`,`create_time`,`update_time`) VALUES (11,'Java','java','java-17',501,'8.130.14.34','/opt/data/java','','','{"APP_PORT":"8080","JAVA_OPTS":"-Xmx512m -Xms256m"}',1,'[2025-10-31 12:53:21] 开始部署 Java Java 17 LTS
[2025-10-31 12:53:21] 连接主机 8.130.14.34...
[2025-10-31 12:53:22] SSH连接成功
[2025-10-31 12:53:22] 检查Docker环境...
[2025-10-31 12:53:22] 使用镜像: crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/openjdk:17-jdk
[2025-10-31 12:53:22] 拉取镜像...
[2025-10-31 12:53:23] 镜像拉取成功
[2025-10-31 12:53:23] 创建临时容器...
[2025-10-31 12:53:23] 提取文件 /usr/local/openjdk-17 -> /usr/local/java17...
[2025-10-31 12:53:29] 清理临时容器...
[2025-10-31 12:53:29] 读取安装脚本 common/templates/02-java/versions/java-17-install.sh...
[2025-10-31 12:53:29] 上传安装脚本...
[2025-10-31 12:53:29] 执行安装脚本...
[2025-10-31 12:53:30] 安装脚本输出:
===== Java 17 LTS 安装配置 =====
安装路径: /usr/local/java17
环境变量文件: /etc/profile.d/java17.sh
环境变量已配置: /etc/profile.d/java17.sh
===== 安装成功 =====

[2025-10-31 12:53:30] 验证安装...
[2025-10-31 12:53:30] 验证结果:
openjdk version "17.0.0.1" 2024-07-02
OpenJDK Runtime Environment (build 17.0.0.1+2-3)
OpenJDK 64-Bit Server VM (build 17.0.0.1+2-3, mixed mode, sharing)

[2025-10-31 12:53:30] 部署完成！
','2025-10-31 12:53:21','2025-10-31 12:53:21');
INSERT INTO `tool_service_deploy` (`id`,`service_name`,`service_id`,`version`,`host_id`,`host_ip`,`install_dir`,`container_name`,`ports`,`env_vars`,`status`,`deploy_log`,`create_time`,`update_time`) VALUES (12,'Elasticsearch','elasticsearch','elasticsearch-8.x',506,'139.9.205.38','/opt/data/elasticsearch','','','{"ES_HEAP_SIZE":"1g","ES_HTTP_PORT":"9200"}',3,'[2025-11-30 22:33:06] 开始部署 Elasticsearch Elasticsearch 8.x
[2025-11-30 22:33:06] 连接主机 139.9.205.38...
[2025-11-30 22:33:07] SSH连接成功
[2025-11-30 22:33:07] 创建安装目录 /opt/data/elasticsearch...
[2025-11-30 22:33:07] 读取模板文件 common/templates/06-elasticsearch/versions/elasticsearch-8.x-docker-compose.yml...
[2025-11-30 22:33:07] 生成环境变量配置...
[2025-11-30 22:33:07] 上传 docker-compose.yml...
[2025-11-30 22:33:07] 上传 .env...
[2025-11-30 22:33:07] 检查Docker环境...
[2025-11-30 22:33:08] 启动服务容器...
[2025-11-30 22:33:08] 启动失败: Process exited with status 127
输出: 
STDERR:
bash: line 1: docker-compose: command not found

','2025-11-30 22:33:06','2025-11-30 22:33:06');

