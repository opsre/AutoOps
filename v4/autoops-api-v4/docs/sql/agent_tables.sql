-- ================================
-- 智能运维Agent模块数据库表结构
-- 创建时间: 2026-01-09
-- ================================

-- ============== 运维巡检相关表 ==============

-- 巡检任务表
CREATE TABLE IF NOT EXISTS `agent_inspection` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(200) NOT NULL COMMENT '巡检任务名称',
  `type` varchar(50) NOT NULL COMMENT '巡检类型:system/service/security/performance/backup',
  `target_type` varchar(20) NOT NULL COMMENT '目标类型:host/group/tag',
  `target_ids` text NOT NULL COMMENT '目标ID列表(JSON数组)',
  `schedule_type` varchar(20) NOT NULL DEFAULT 'manual' COMMENT '调度类型:manual/cron',
  `cron_expr` varchar(100) DEFAULT NULL COMMENT 'Cron表达式',
  `enabled` tinyint NOT NULL DEFAULT 1 COMMENT '是否启用:1=启用,2=禁用',
  `template` text NOT NULL COMMENT '巡检模板(JSON配置)',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='巡检任务表';

-- 巡检记录表
CREATE TABLE IF NOT EXISTS `inspection_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `inspection_id` int unsigned NOT NULL COMMENT '关联巡检任务',
  `host_id` int unsigned NOT NULL COMMENT '主机ID',
  `host_name` varchar(200) NOT NULL COMMENT '主机名',
  `host_ip` varchar(50) DEFAULT NULL COMMENT '主机IP',
  `status` tinyint NOT NULL COMMENT '状态:1=执行中,2=成功,3=失败,4=部分成功',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `duration` int DEFAULT 0 COMMENT '耗时(秒)',
  `total_checks` int DEFAULT 0 COMMENT '总检查项',
  `passed_checks` int DEFAULT 0 COMMENT '通过数',
  `failed_checks` int DEFAULT 0 COMMENT '失败数',
  `warn_checks` int DEFAULT 0 COMMENT '警告数',
  `risk_level` tinyint DEFAULT 0 COMMENT '风险等级:0=正常,1=低危,2=中危,3=高危',
  `result` text DEFAULT NULL COMMENT '详细结果(JSON)',
  `summary` text DEFAULT NULL COMMENT 'AI生成的总结',
  `suggestions` text DEFAULT NULL COMMENT 'AI生成的建议',
  `report_url` varchar(500) DEFAULT NULL COMMENT '报告文件URL',
  `error_msg` text DEFAULT NULL COMMENT '错误信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_inspection_id` (`inspection_id`),
  KEY `idx_host_id` (`host_id`),
  KEY `idx_status` (`status`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='巡检记录表';

-- 巡检项目表
CREATE TABLE IF NOT EXISTS `inspection_item` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `record_id` int unsigned NOT NULL COMMENT '关联巡检记录',
  `category` varchar(50) NOT NULL COMMENT '分类:system/service/security/performance',
  `name` varchar(200) NOT NULL COMMENT '检查项名称',
  `status` tinyint NOT NULL COMMENT '状态:1=通过,2=失败,3=警告',
  `actual` text DEFAULT NULL COMMENT '实际值',
  `expected` varchar(500) DEFAULT NULL COMMENT '期望值',
  `message` text DEFAULT NULL COMMENT '详细信息',
  `suggestion` text DEFAULT NULL COMMENT '修复建议',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_record_id` (`record_id`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='巡检项目表';

-- ============== 故障分析相关表 ==============

-- 故障分析表
CREATE TABLE IF NOT EXISTS `incident_analysis` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `incident_id` varchar(100) NOT NULL COMMENT '故障ID(唯一标识)',
  `title` varchar(500) NOT NULL COMMENT '故障标题',
  `host_id` int unsigned NOT NULL COMMENT '故障主机ID',
  `host_name` varchar(200) NOT NULL COMMENT '主机名',
  `host_ip` varchar(50) DEFAULT NULL COMMENT '主机IP',
  `category` varchar(50) NOT NULL COMMENT '故障类别:cpu/memory/disk/network/service/custom',
  `severity` tinyint NOT NULL COMMENT '严重程度:1=低,2=中,3=高,4=紧急',
  `status` tinyint NOT NULL COMMENT '状态:1=分析中,2=已定位,3=修复中,4=已修复,5=已关闭',
  `detect_time` datetime NOT NULL COMMENT '检测时间',
  `start_time` datetime DEFAULT NULL COMMENT '分析开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '分析结束时间',
  `duration` int DEFAULT 0 COMMENT '耗时(秒)',
  `root_cause` text DEFAULT NULL COMMENT '根因分析结果',
  `symptoms` text DEFAULT NULL COMMENT '故障症状(JSON数组)',
  `evidence` text DEFAULT NULL COMMENT '证据数据(JSON)',
  `solution` text DEFAULT NULL COMMENT '解决方案',
  `fix_script` text DEFAULT NULL COMMENT '修复脚本',
  `risk_level` tinyint DEFAULT 0 COMMENT '自动修复风险等级:0=安全,1=低危,2=中危,3=高危',
  `is_auto_fixed` tinyint(1) DEFAULT 0 COMMENT '是否自动修复',
  `fix_result` text DEFAULT NULL COMMENT '修复结果',
  `error_msg` text DEFAULT NULL COMMENT '错误信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_incident_id` (`incident_id`),
  KEY `idx_host_id` (`host_id`),
  KEY `idx_category` (`category`),
  KEY `idx_severity` (`severity`),
  KEY `idx_status` (`status`),
  KEY `idx_detect_time` (`detect_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='故障分析表';

-- 故障证据表
CREATE TABLE IF NOT EXISTS `incident_evidence` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `incident_id` varchar(100) NOT NULL COMMENT '故障ID',
  `type` varchar(50) NOT NULL COMMENT '类型:log/metric/process/network/file',
  `source` varchar(500) NOT NULL COMMENT '来源:系统日志路径、监控指标名等',
  `content` text NOT NULL COMMENT '证据内容',
  `collect_time` datetime NOT NULL COMMENT '采集时间',
  `relevance` decimal(3,2) DEFAULT 0.5 COMMENT '相关度评分(0-1)',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_incident_id` (`incident_id`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='故障证据表';

-- 故障知识库表
CREATE TABLE IF NOT EXISTS `incident_knowledge` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(500) NOT NULL COMMENT '故障标题',
  `category` varchar(50) NOT NULL COMMENT '故障类别',
  `symptoms` text NOT NULL COMMENT '症状描述(JSON数组)',
  `root_cause` text NOT NULL COMMENT '根因',
  `solution` text NOT NULL COMMENT '解决方案',
  `fix_script` text DEFAULT NULL COMMENT '修复脚本',
  `keywords` varchar(1000) DEFAULT NULL COMMENT '关键词(用于匹配)',
  `match_count` int DEFAULT 0 COMMENT '匹配次数(用于排序)',
  `success_rate` decimal(5,2) DEFAULT 0.00 COMMENT '成功率',
  `enabled` tinyint DEFAULT 1 COMMENT '是否启用:1=启用,2=禁用',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category`),
  KEY `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='故障知识库表';

-- ============== 容量预测和告警相关表 ==============

-- 容量预测表
CREATE TABLE IF NOT EXISTS `capacity_prediction` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(200) NOT NULL COMMENT '预测任务名称',
  `target_type` varchar(20) NOT NULL COMMENT '目标类型:host/group/cluster',
  `target_id` int unsigned NOT NULL COMMENT '目标ID',
  `target_name` varchar(200) NOT NULL COMMENT '目标名称',
  `metric_type` varchar(50) NOT NULL COMMENT '指标类型:cpu/memory/disk/network/custom',
  `metric_name` varchar(100) NOT NULL COMMENT '指标名称',
  `current_value` decimal(10,2) DEFAULT 0.00 COMMENT '当前值',
  `threshold` decimal(10,2) NOT NULL COMMENT '阈值(百分比或绝对值)',
  `predict_days` int NOT NULL DEFAULT 30 COMMENT '预测天数',
  `algorithm` varchar(20) DEFAULT 'linear' COMMENT '算法:linear/polynomial/arima/lstm',
  `confidence` decimal(3,2) DEFAULT 0.00 COMMENT '置信度(0-1)',
  `predict_time` datetime DEFAULT NULL COMMENT '预测到达阈值时间',
  `days_left` int DEFAULT 0 COMMENT '剩余天数',
  `trend_type` varchar(20) DEFAULT NULL COMMENT '趋势类型:increasing/decreasing/stable/fluctuating',
  `risk_level` tinyint DEFAULT 0 COMMENT '风险等级:0=安全,1=关注,2=预警,3=紧急',
  `suggestion` text DEFAULT NULL COMMENT 'AI建议',
  `status` tinyint DEFAULT 1 COMMENT '状态:1=正常,2=预警,3=告警',
  `enabled` tinyint DEFAULT 1 COMMENT '是否启用:1=启用,2=禁用',
  `last_calc_time` datetime DEFAULT NULL COMMENT '最后计算时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_target_id` (`target_id`),
  KEY `idx_metric_type` (`metric_type`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='容量预测表';

-- 预测历史记录表
CREATE TABLE IF NOT EXISTS `prediction_history` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `prediction_id` int unsigned NOT NULL COMMENT '预测任务ID',
  `predict_time` datetime NOT NULL COMMENT '预测时间点',
  `predict_value` decimal(10,2) DEFAULT 0.00 COMMENT '预测值',
  `actual_value` decimal(10,2) DEFAULT 0.00 COMMENT '实际值',
  `accuracy` decimal(5,2) DEFAULT 0.00 COMMENT '准确度',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_prediction_id` (`prediction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预测历史记录表';

-- 智能告警规则表
CREATE TABLE IF NOT EXISTS `alert_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(200) NOT NULL COMMENT '规则名称',
  `type` varchar(50) NOT NULL COMMENT '类型:threshold/trend/anomaly/prediction',
  `target_type` varchar(20) NOT NULL COMMENT '目标类型:host/group/cluster',
  `target_ids` text NOT NULL COMMENT '目标ID列表(JSON)',
  `metric_type` varchar(50) NOT NULL COMMENT '指标类型',
  `condition` text NOT NULL COMMENT '告警条件(JSON配置)',
  `severity` tinyint NOT NULL COMMENT '严重程度:1=低,2=中,3=高,4=紧急',
  `notify_channels` varchar(500) DEFAULT NULL COMMENT '通知渠道(JSON数组)',
  `notify_users` text DEFAULT NULL COMMENT '通知用户(JSON数组)',
  `silence_period` int DEFAULT 60 COMMENT '静默期(分钟)',
  `enabled` tinyint DEFAULT 1 COMMENT '是否启用:1=启用,2=禁用',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_metric_type` (`metric_type`),
  KEY `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='智能告警规则表';

-- 告警事件表
CREATE TABLE IF NOT EXISTS `alert_event` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `rule_id` int unsigned NOT NULL COMMENT '关联告警规则',
  `rule_name` varchar(200) NOT NULL COMMENT '规则名称',
  `host_id` int unsigned NOT NULL COMMENT '主机ID',
  `host_name` varchar(200) NOT NULL COMMENT '主机名',
  `host_ip` varchar(50) DEFAULT NULL COMMENT '主机IP',
  `metric_type` varchar(50) NOT NULL COMMENT '指标类型',
  `metric_name` varchar(100) NOT NULL COMMENT '指标名称',
  `current_value` decimal(10,2) DEFAULT 0.00 COMMENT '当前值',
  `threshold` decimal(10,2) DEFAULT 0.00 COMMENT '阈值',
  `severity` tinyint NOT NULL COMMENT '严重程度:1=低,2=中,3=高,4=紧急',
  `status` tinyint NOT NULL COMMENT '状态:1=触发,2=恢复,3=已确认,4=已解决',
  `trigger_time` datetime NOT NULL COMMENT '触发时间',
  `recover_time` datetime DEFAULT NULL COMMENT '恢复时间',
  `duration` int DEFAULT 0 COMMENT '持续时间(秒)',
  `ai_summary` text DEFAULT NULL COMMENT 'AI分析摘要',
  `ai_suggestion` text DEFAULT NULL COMMENT 'AI处理建议',
  `notify_status` tinyint DEFAULT 1 COMMENT '通知状态:1=待发送,2=已发送,3=发送失败',
  `notify_error` text DEFAULT NULL COMMENT '通知错误信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_host_id` (`host_id`),
  KEY `idx_metric_type` (`metric_type`),
  KEY `idx_severity` (`severity`),
  KEY `idx_status` (`status`),
  KEY `idx_trigger_time` (`trigger_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警事件表';

-- ============== Agent通用表 ==============

-- Agent任务表
CREATE TABLE IF NOT EXISTS `agent_task` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `type` varchar(50) NOT NULL COMMENT '任务类型:inspection/incident/prediction/custom',
  `name` varchar(200) NOT NULL COMMENT '任务名称',
  `description` varchar(500) DEFAULT NULL COMMENT '任务描述',
  `priority` tinyint DEFAULT 5 COMMENT '优先级:1-10',
  `status` tinyint NOT NULL COMMENT '状态:1=待执行,2=执行中,3=成功,4=失败,5=取消',
  `params` text DEFAULT NULL COMMENT '任务参数(JSON)',
  `result` text DEFAULT NULL COMMENT '执行结果(JSON)',
  `error_msg` text DEFAULT NULL COMMENT '错误信息',
  `retry` int DEFAULT 0 COMMENT '重试次数',
  `max_retry` int DEFAULT 3 COMMENT '最大重试次数',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Agent任务表';

-- 知识库表
CREATE TABLE IF NOT EXISTS `knowledge_base` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `type` varchar(50) NOT NULL COMMENT '类型:command/solution/best-practice/faq',
  `category` varchar(50) NOT NULL COMMENT '分类',
  `title` varchar(500) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '内容',
  `keywords` varchar(1000) DEFAULT NULL COMMENT '关键词(JSON数组)',
  `tags` varchar(500) DEFAULT NULL COMMENT '标签(JSON数组)',
  `score` decimal(3,2) DEFAULT 0.50 COMMENT '评分',
  `use_count` int DEFAULT 0 COMMENT '使用次数',
  `enabled` tinyint DEFAULT 1 COMMENT '是否启用:1=启用,2=禁用',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_category` (`category`),
  KEY `idx_enabled` (`enabled`),
  FULLTEXT KEY `idx_fulltext` (`title`,`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='知识库表';

-- Prompt模板表
CREATE TABLE IF NOT EXISTS `prompt_template` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(200) NOT NULL COMMENT '模板名称',
  `category` varchar(50) NOT NULL COMMENT '分类:inspection/incident/prediction/general',
  `scene` varchar(50) NOT NULL COMMENT '场景:summary/analysis/suggestion/fix',
  `template` text NOT NULL COMMENT '模板内容(支持变量替换)',
  `variables` varchar(1000) DEFAULT NULL COMMENT '变量列表(JSON数组)',
  `system_prompt` text DEFAULT NULL COMMENT '系统提示词',
  `temperature` decimal(3,2) DEFAULT 0.70 COMMENT '温度参数',
  `max_tokens` int DEFAULT 2000 COMMENT '最大Token数',
  `model_id` int unsigned DEFAULT NULL COMMENT '默认模型ID',
  `enabled` tinyint DEFAULT 1 COMMENT '是否启用:1=启用,2=禁用',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_category` (`category`),
  KEY `idx_scene` (`scene`),
  KEY `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Prompt模板表';

-- Agent对话历史表
CREATE TABLE IF NOT EXISTS `ai_agent_chat_history` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `session_id` varchar(64) NOT NULL COMMENT '会话ID',
  `user_id` int unsigned NOT NULL COMMENT '用户ID',
  `role` varchar(20) NOT NULL COMMENT '角色:user/assistant/system',
  `message` text NOT NULL COMMENT '消息内容',
  `intent` varchar(50) DEFAULT NULL COMMENT '识别的意图类型',
  `intent_conf` decimal(3,2) DEFAULT NULL COMMENT '意图识别置信度',
  `entities` json DEFAULT NULL COMMENT '提取的实体(JSON)',
  `task_id` int unsigned DEFAULT NULL COMMENT '关联的任务ID',
  `task_type` varchar(50) DEFAULT NULL COMMENT '任务类型',
  `status` tinyint DEFAULT 1 COMMENT '状态:1=处理中,2=成功,3=失败',
  `error_msg` text DEFAULT NULL COMMENT '错误信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI对话历史表';

