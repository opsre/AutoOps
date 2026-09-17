-- Database: mayfly-go
-- Export Time: 2025-12-10 10:28:17
-- Tables: 43
-- Export Type: full

-- Table structure for t_db
CREATE TABLE `t_db` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `code` varchar(32),
 `name` varchar(191),
 `get_database_mode` tinyint(3) COMMENT '库名获取方式（-1.实时获取、1.指定库名）',
 `database` varchar(1000),
 `remark` varchar(191),
 `instance_id` bigint(20) NOT NULL,
 `auth_cert_name` varchar(36) COMMENT '授权凭证名',
 `create_time` datetime,
 `creator_id` bigint(19),
 `creator` varchar(191),
 `update_time` datetime,
 `modifier_id` bigint(19),
 `modifier` varchar(191),
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '数据库资源信息表';
ALTER TABLE `t_db` ADD  INDEX `idx_code`(`code`) USING BTREE;

-- Data for table t_db
-- No data

-- Table structure for t_db_backup
CREATE TABLE `t_db_backup` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(32) NOT NULL COMMENT '备份名称',
 `db_instance_id` bigint(20) NOT NULL COMMENT '数据库实例ID',
 `db_name` varchar(64) NOT NULL COMMENT '数据库名称',
 `repeated` tinyint(3) COMMENT '是否重复执行',
 `interval` bigint(19) COMMENT '备份周期',
 `max_save_days` int(10) NOT NULL DEFAULT 0 COMMENT '最大保留天数',
 `start_time` datetime COMMENT '首次备份时间',
 `enabled` tinyint(3) COMMENT '是否启用',
 `enabled_desc` varchar(64) COMMENT '任务启用描述',
 `last_status` tinyint(3) COMMENT '上次备份状态',
 `last_result` varchar(256) COMMENT '上次备份结果',
 `last_time` datetime COMMENT '上次备份时间',
 `create_time` datetime,
 `creator_id` bigint(20),
 `creator` varchar(32),
 `update_time` datetime,
 `modifier_id` bigint(20),
 `modifier` varchar(32),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `t_db_backup` ADD  INDEX `idx_db_instance_id`(`db_instance_id`) USING BTREE;
ALTER TABLE `t_db_backup` ADD  INDEX `idx_db_name`(`db_name`) USING BTREE;

-- Data for table t_db_backup
-- No data

-- Table structure for t_db_backup_history
CREATE TABLE `t_db_backup_history` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(64) NOT NULL COMMENT '历史备份名称',
 `db_backup_id` bigint(20) NOT NULL COMMENT '数据库备份ID',
 `db_instance_id` bigint(20) NOT NULL COMMENT '数据库实例ID',
 `db_name` varchar(64) NOT NULL COMMENT '数据库名称',
 `uuid` varchar(36) NOT NULL COMMENT '历史备份uuid',
 `binlog_file_name` varchar(32) COMMENT 'BINLOG文件名',
 `binlog_sequence` bigint(19) COMMENT 'BINLOG序列号',
 `binlog_position` bigint(19) COMMENT 'BINLOG位置',
 `create_time` datetime COMMENT '历史备份创建时间',
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime,
 `restoring` tinyint(3) NOT NULL DEFAULT 0 COMMENT '备份历史恢复标识',
 `deleting` tinyint(3) NOT NULL DEFAULT 0 COMMENT '备份历史删除标识', 
PRIMARY KEY (id)
);
ALTER TABLE `t_db_backup_history` ADD  INDEX `idx_db_backup_id`(`db_backup_id`) USING BTREE;
ALTER TABLE `t_db_backup_history` ADD  INDEX `idx_db_instance_id`(`db_instance_id`) USING BTREE;
ALTER TABLE `t_db_backup_history` ADD  INDEX `idx_db_name`(`db_name`) USING BTREE;

-- Data for table t_db_backup_history
-- No data

-- Table structure for t_db_binlog
CREATE TABLE `t_db_binlog` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `db_instance_id` bigint(20) NOT NULL COMMENT '数据库实例ID',
 `last_status` bigint(19) COMMENT '上次下载状态',
 `last_result` varchar(256) COMMENT '上次下载结果',
 `last_time` datetime COMMENT '上次下载时间',
 `create_time` datetime,
 `creator_id` bigint(20),
 `creator` varchar(32),
 `update_time` datetime,
 `modifier_id` bigint(20),
 `modifier` varchar(32),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `t_db_binlog` ADD  INDEX `idx_db_instance_id`(`db_instance_id`) USING BTREE;

-- Data for table t_db_binlog
-- No data

-- Table structure for t_db_binlog_history
CREATE TABLE `t_db_binlog_history` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `db_instance_id` bigint(20) NOT NULL COMMENT '数据库实例ID',
 `file_name` varchar(32) COMMENT 'BINLOG文件名称',
 `file_size` bigint(19) COMMENT 'BINLOG文件大小',
 `sequence` bigint(19) COMMENT 'BINLOG序列号',
 `first_event_time` datetime COMMENT '首次事件时间',
 `last_event_time` datetime COMMENT '最新事件时间',
 `create_time` datetime,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `t_db_binlog_history` ADD  INDEX `idx_db_instance_id`(`db_instance_id`) USING BTREE;

-- Data for table t_db_binlog_history
-- No data

-- Table structure for t_db_data_sync_log
CREATE TABLE `t_db_data_sync_log` (
 `id` bigint(19) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
 `task_id` bigint(19) NOT NULL COMMENT '同步任务表id',
 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `data_sql_full` text NOT NULL COMMENT '执行的完整sql',
 `res_num` int(10) COMMENT '收到数据条数',
 `err_text` text COMMENT '错误日志',
 `status` tinyint(3) NOT NULL DEFAULT 1 COMMENT '状态:1.成功  0.失败',
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0 COMMENT '是否删除 1是 0 否', 
PRIMARY KEY (id)
) COMMENT '数据同步日志';
ALTER TABLE `t_db_data_sync_log` ADD  INDEX `t_db_data_sync_log_taskid_idx`(`task_id`) USING BTREE COMMENT 't_db_data_sync_log表(taskid)普通索引';

-- Data for table t_db_data_sync_log
-- No data

-- Table structure for t_db_data_sync_task
CREATE TABLE `t_db_data_sync_task` (
 `id` bigint(19) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
 `creator_id` bigint(19) NOT NULL COMMENT '创建人id',
 `creator` varchar(100) NOT NULL COMMENT '创建人姓名',
 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
 `modifier` varchar(100) NOT NULL COMMENT '修改人姓名',
 `modifier_id` bigint(19) NOT NULL COMMENT '修改人id',
 `task_name` varchar(500) NOT NULL COMMENT '任务名',
 `task_cron` varchar(50) NOT NULL COMMENT '任务Cron表达式',
 `src_db_id` bigint(19) NOT NULL COMMENT '源数据库ID',
 `src_db_name` varchar(100) COMMENT '源数据库名',
 `src_tag_path` varchar(200) COMMENT '源数据库tag路径',
 `target_db_id` bigint(19) NOT NULL COMMENT '目标数据库ID',
 `target_db_name` varchar(100) COMMENT '目标数据库名',
 `target_tag_path` varchar(200) COMMENT '目标数据库tag路径',
 `target_table_name` varchar(100) COMMENT '目标数据库表名',
 `data_sql` text NOT NULL COMMENT '数据查询sql',
 `page_size` int(10) NOT NULL COMMENT '数据同步分页大小',
 `upd_field` varchar(100) NOT NULL DEFAULT id COMMENT '更新字段，默认"id"',
 `upd_field_val` varchar(100) COMMENT '当前更新值',
 `id_rule` tinyint(3) NOT NULL DEFAULT 1 COMMENT 'id生成规则：1、MD5(时间戳+更新字段的值)。2、无(不自动生成id，选择无的时候需要指定主键ID字段是数据源哪个字段)',
 `pk_field` varchar(100) DEFAULT id COMMENT '主键id字段名，默认"id"',
 `field_map` text COMMENT '字段映射json',
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime,
 `status` tinyint(3) NOT NULL DEFAULT 1 COMMENT '状态 1启用 2停用',
 `recent_state` tinyint(3) NOT NULL DEFAULT 0 COMMENT '最近一次状态 0未执行 1成功 2失败',
 `task_key` varchar(100) COMMENT '定时任务唯一uuid key',
 `running_state` tinyint(3) DEFAULT 2 COMMENT '运行时状态 1运行中、2待运行、3已停止',
 `duplicate_strategy` tinyint(3) DEFAULT -1 COMMENT '唯一键冲突策略 -1：无，1：忽略，2：覆盖', 
PRIMARY KEY (id)
) COMMENT '数据同步';

-- Data for table t_db_data_sync_task
-- No data

-- Table structure for t_db_instance
CREATE TABLE `t_db_instance` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `code` varchar(36) COMMENT '唯一编号',
 `name` varchar(32) COMMENT '数据库实例名称',
 `host` varchar(100) NOT NULL,
 `port` int(10),
 `extra` varchar(255) COMMENT '连接需要的额外参数，如oracle数据库需要sid等',
 `type` varchar(20) NOT NULL COMMENT '数据库实例类型(mysql...)',
 `params` varchar(125) COMMENT '其他连接参数',
 `network` varchar(20),
 `ssh_tunnel_machine_id` bigint(19) COMMENT 'ssh隧道的机器id',
 `remark` varchar(125) COMMENT '备注，描述等',
 `create_time` datetime,
 `creator_id` bigint(19),
 `creator` varchar(32),
 `update_time` datetime,
 `modifier_id` bigint(19),
 `modifier` varchar(32),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '数据库实例信息表';

-- Data for table t_db_instance
-- No data

-- Table structure for t_db_restore
CREATE TABLE `t_db_restore` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `db_instance_id` bigint(20) NOT NULL COMMENT '数据库实例ID',
 `db_name` varchar(64) NOT NULL COMMENT '数据库名称',
 `repeated` tinyint(3) COMMENT '是否重复执行',
 `interval` bigint(19) COMMENT '恢复周期',
 `start_time` datetime COMMENT '首次恢复时间',
 `enabled` tinyint(3) COMMENT '是否启用',
 `enabled_desc` varchar(64) COMMENT '任务启用描述',
 `last_status` tinyint(3) COMMENT '上次恢复状态',
 `last_result` varchar(256) COMMENT '上次恢复结果',
 `last_time` datetime COMMENT '上次恢复时间',
 `point_in_time` datetime COMMENT '恢复时间点',
 `db_backup_id` bigint(20) COMMENT '备份ID',
 `db_backup_history_id` bigint(20) COMMENT '历史备份ID',
 `db_backup_history_name` varchar(64) COMMENT '历史备份名称',
 `create_time` datetime,
 `creator_id` bigint(20),
 `creator` varchar(32),
 `update_time` datetime,
 `modifier_id` bigint(20),
 `modifier` varchar(32),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `t_db_restore` ADD  INDEX `idx_db_instane_id`(`db_instance_id`) USING BTREE;
ALTER TABLE `t_db_restore` ADD  INDEX `idx_db_name`(`db_name`) USING BTREE;

-- Data for table t_db_restore
-- No data

-- Table structure for t_db_restore_history
CREATE TABLE `t_db_restore_history` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `db_restore_id` bigint(20) NOT NULL COMMENT '恢复ID',
 `create_time` datetime COMMENT '历史恢复创建时间',
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);
ALTER TABLE `t_db_restore_history` ADD  INDEX `idx_db_restore_id`(`db_restore_id`) USING BTREE;

-- Data for table t_db_restore_history
-- No data

-- Table structure for t_db_sql
CREATE TABLE `t_db_sql` (
 `id` bigint(19) NOT NULL AUTO_INCREMENT,
 `db_id` bigint(19) NOT NULL COMMENT '数据库实例id',
 `db` varchar(125) NOT NULL COMMENT '数据库',
 `name` varchar(60) COMMENT 'sql模板名',
 `sql` text,
 `type` tinyint(3) NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(32),
 `create_time` datetime NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19),
 `modifier` varchar(255),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '数据库sql信息';

-- Data for table t_db_sql
-- No data

-- Table structure for t_db_sql_exec
CREATE TABLE `t_db_sql_exec` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `db_id` bigint(19) NOT NULL COMMENT '数据库id',
 `db` varchar(128) NOT NULL COMMENT '数据库',
 `table` varchar(128) NOT NULL COMMENT '表名',
 `type` varchar(255) NOT NULL COMMENT 'sql类型',
 `sql` varchar(5000) NOT NULL COMMENT '执行sql',
 `old_value` varchar(5000) COMMENT '操作前旧值',
 `remark` varchar(128) COMMENT '备注',
 `status` tinyint(3) COMMENT '执行状态',
 `flow_biz_key` varchar(64) COMMENT '流程关联的业务key',
 `res` varchar(1000) COMMENT '执行结果',
 `create_time` datetime NOT NULL,
 `creator` varchar(36) NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier` varchar(36) NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '数据库sql执行记录表';
ALTER TABLE `t_db_sql_exec` ADD  INDEX `idx_flow_biz_key`(`flow_biz_key`) USING BTREE;

-- Data for table t_db_sql_exec
-- No data

-- Table structure for t_db_transfer_files
CREATE TABLE `t_db_transfer_files` (
 `id` bigint(19) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0 COMMENT '是否删除',
 `delete_time` datetime COMMENT '删除时间',
 `status` tinyint(3) NOT NULL DEFAULT 1 COMMENT '状态，1、执行中 2、执行失败 3、 执行成功',
 `task_id` bigint(19) COMMENT '迁移任务ID',
 `log_id` bigint(19) COMMENT '日志ID',
 `file_db_type` varchar(200) COMMENT 'sql文件数据库类型',
 `file_key` varchar(50) COMMENT '文件', 
PRIMARY KEY (id)
) COMMENT '数据库迁移文件管理';

-- Data for table t_db_transfer_files
-- No data

-- Table structure for t_db_transfer_task
CREATE TABLE `t_db_transfer_task` (
 `id` bigint(19) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
 `creator_id` bigint(19) NOT NULL COMMENT '创建人id',
 `creator` varchar(100) NOT NULL COMMENT '创建人姓名',
 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `modifier_id` bigint(19) NOT NULL COMMENT '修改人id',
 `modifier` varchar(100) NOT NULL COMMENT '修改人姓名',
 `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
 `is_deleted` tinyint(3) DEFAULT 0 COMMENT '是否删除',
 `task_name` varchar(100) COMMENT '任务名',
 `cron_able` tinyint(3) NOT NULL DEFAULT 0 COMMENT '是否定时  1是 -1否',
 `cron` varchar(20) COMMENT '定时任务cron表达式',
 `task_key` varchar(100) COMMENT '定时任务唯一uuid key',
 `mode` tinyint(3) NOT NULL DEFAULT 1 COMMENT '数据迁移方式，1、迁移到数据库  2、迁移到文件',
 `target_file_db_type` varchar(200) COMMENT '目标文件语言类型，类型枚举同target_db_type',
 `file_save_days` int(10) COMMENT '文件保存天数',
 `status` tinyint(3) NOT NULL DEFAULT 1 COMMENT '启用状态 1启用 -1禁用',
 `upd_field_src` varchar(100) COMMENT '更新值来源字段，默认同更新字段，如果查询结果指定了字段别名且与原更新字段不一致，则取这个字段值为当前更新值',
 `delete_time` datetime COMMENT '删除时间',
 `checked_keys` text NOT NULL COMMENT '选中需要迁移的表',
 `delete_table` tinyint(3) NOT NULL COMMENT '创建表前是否删除表  1是  -1否',
 `name_case` tinyint(3) NOT NULL COMMENT '表名、字段大小写转换  1无  2大写  3小写',
 `strategy` tinyint(3) NOT NULL COMMENT '迁移策略  1全量  2增量',
 `running_state` tinyint(3) DEFAULT 2 COMMENT '运行状态 1运行中  2待运行',
 `src_db_id` bigint(19) NOT NULL COMMENT '源库id',
 `src_db_name` varchar(200) NOT NULL COMMENT '源库名',
 `src_tag_path` varchar(200) NOT NULL COMMENT '源库tagPath',
 `src_db_type` varchar(200) NOT NULL COMMENT '源库类型',
 `src_inst_name` varchar(200) NOT NULL COMMENT '源库实例名',
 `target_db_id` bigint(19) NOT NULL COMMENT '目标库id',
 `target_db_name` varchar(200) NOT NULL COMMENT '目标库名',
 `target_tag_path` varchar(200) NOT NULL COMMENT '目标库类型',
 `target_db_type` varchar(200) NOT NULL COMMENT '目标库实例名',
 `target_inst_name` varchar(200) NOT NULL COMMENT '目标库tagPath',
 `log_id` bigint(19) NOT NULL COMMENT '日志id', 
PRIMARY KEY (id)
) COMMENT '数据库迁移任务表';

-- Data for table t_db_transfer_task
-- No data

-- Table structure for t_flow_procdef
CREATE TABLE `t_flow_procdef` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `def_key` varchar(100) NOT NULL COMMENT '流程定义key',
 `name` varchar(191) COMMENT '流程名称',
 `status` tinyint(3) COMMENT '状态',
 `tasks` varchar(3000) NOT NULL COMMENT '审批节点任务信息',
 `condition` text COMMENT '触发审批的条件（计算结果返回1则需要启用该流程）',
 `remark` varchar(191),
 `create_time` datetime NOT NULL,
 `creator` varchar(191),
 `creator_id` bigint(19) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier` varchar(191),
 `modifier_id` bigint(19) NOT NULL,
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '流程-流程定义';

-- Data for table t_flow_procdef
-- No data

-- Table structure for t_flow_procinst
CREATE TABLE `t_flow_procinst` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `procdef_id` bigint(19) NOT NULL COMMENT '流程定义id',
 `procdef_name` varchar(64) NOT NULL COMMENT '流程定义名称',
 `task_key` varchar(100) COMMENT '当前任务key',
 `status` tinyint(3) COMMENT '状态',
 `biz_type` varchar(64) NOT NULL COMMENT '关联业务类型',
 `biz_key` varchar(64) NOT NULL COMMENT '关联业务key',
 `biz_form` text COMMENT '业务form',
 `biz_status` tinyint(3) COMMENT '业务状态',
 `biz_handle_res` varchar(1000) COMMENT '关联的业务处理结果',
 `remark` varchar(191),
 `end_time` datetime COMMENT '结束时间',
 `duration` bigint(19) COMMENT '流程持续时间（开始到结束）',
 `create_time` datetime NOT NULL COMMENT '流程发起时间',
 `creator` varchar(191) COMMENT '流程发起人',
 `creator_id` bigint(19) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier` varchar(191),
 `modifier_id` bigint(19) NOT NULL,
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '流程-流程实例(根据流程定义开启一个流程)';
ALTER TABLE `t_flow_procinst` ADD  INDEX `idx_procdef_id`(`procdef_id`) USING BTREE;

-- Data for table t_flow_procinst
-- No data

-- Table structure for t_flow_procinst_task
CREATE TABLE `t_flow_procinst_task` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `procinst_id` bigint(19) NOT NULL COMMENT '流程实例id',
 `task_key` varchar(64) NOT NULL COMMENT '任务key',
 `task_name` varchar(64) COMMENT '任务名称',
 `assignee` varchar(64) COMMENT '分配到该任务的用户',
 `status` tinyint(3) COMMENT '状态',
 `remark` varchar(191),
 `end_time` datetime COMMENT '结束时间',
 `duration` bigint(19) COMMENT '任务持续时间（开始到结束）',
 `create_time` datetime NOT NULL COMMENT '任务开始时间',
 `creator` varchar(191),
 `creator_id` bigint(19) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier` varchar(191),
 `modifier_id` bigint(19) NOT NULL,
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '流程-流程实例任务';
ALTER TABLE `t_flow_procinst_task` ADD  INDEX `idx_procinst_id`(`procinst_id`) USING BTREE;

-- Data for table t_flow_procinst_task
-- No data

-- Table structure for t_machine
CREATE TABLE `t_machine` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `code` varchar(32) COMMENT 'code',
 `name` varchar(32),
 `ip` varchar(50) NOT NULL,
 `port` int(10) NOT NULL,
 `protocol` tinyint(3) DEFAULT 1 COMMENT '协议  1、SSH  2、RDP',
 `ssh_tunnel_machine_id` bigint(19) COMMENT 'ssh隧道的机器id',
 `enable_recorder` tinyint(3) COMMENT '是否启用终端回放记录',
 `status` tinyint(3) NOT NULL COMMENT '状态: 1:启用; -1:禁用',
 `remark` varchar(255),
 `extra` varchar(200) COMMENT '额外信息',
 `need_monitor` tinyint(3),
 `create_time` datetime NOT NULL,
 `creator` varchar(16),
 `creator_id` bigint(19),
 `update_time` datetime NOT NULL,
 `modifier` varchar(12),
 `modifier_id` bigint(19),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '机器信息';

-- Data for table t_machine
-- No data

-- Table structure for t_machine_cmd_conf
CREATE TABLE `t_machine_cmd_conf` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(100) COMMENT '名称',
 `cmds` varchar(500) COMMENT '命令配置',
 `status` tinyint(3) COMMENT '状态',
 `stratege` varchar(100) COMMENT '策略',
 `remark` varchar(50) COMMENT '备注',
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(36) NOT NULL,
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '机器命令配置';

-- Data for table t_machine_cmd_conf
-- No data

-- Table structure for t_machine_cron_job
CREATE TABLE `t_machine_cron_job` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `key` varchar(32) NOT NULL COMMENT 'key',
 `name` varchar(255) NOT NULL COMMENT '名称',
 `cron` varchar(255) NOT NULL COMMENT 'cron表达式',
 `script` text COMMENT '脚本内容',
 `remark` varchar(255) COMMENT '备注',
 `status` tinyint(3) COMMENT '状态',
 `save_exec_res_type` tinyint(3) COMMENT '保存执行记录类型',
 `last_exec_time` datetime COMMENT '最后执行时间',
 `creator_id` bigint(19),
 `creator` varchar(32),
 `modifier_id` bigint(19),
 `modifier` varchar(255),
 `create_time` datetime,
 `update_time` datetime,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '机器计划任务';

-- Data for table t_machine_cron_job
-- No data

-- Table structure for t_machine_cron_job_exec
CREATE TABLE `t_machine_cron_job_exec` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `cron_job_id` bigint(19),
 `machine_code` varchar(36),
 `status` tinyint(3) COMMENT '状态',
 `res` varchar(4000) COMMENT '执行结果',
 `exec_time` datetime COMMENT '执行时间',
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '机器计划任务执行记录';

-- Data for table t_machine_cron_job_exec
-- No data

-- Table structure for t_machine_file
CREATE TABLE `t_machine_file` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '机器文件配置（linux一切皆文件，故也可以表示目录）',
 `machine_id` bigint(19) NOT NULL,
 `name` varchar(45) NOT NULL,
 `path` varchar(45) NOT NULL,
 `type` varchar(45) NOT NULL COMMENT '1：目录；2：文件',
 `creator_id` bigint(20),
 `creator` varchar(45),
 `modifier_id` bigint(20),
 `modifier` varchar(45),
 `create_time` datetime NOT NULL,
 `update_time` datetime,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '机器文件';

-- Data for table t_machine_file
-- No data

-- Table structure for t_machine_monitor
CREATE TABLE `t_machine_monitor` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `machine_id` bigint(20) NOT NULL COMMENT '机器id',
 `cpu_rate` float(255,2),
 `mem_rate` float(255,2),
 `sys_load` varchar(32),
 `create_time` datetime NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);

-- Data for table t_machine_monitor
-- No data

-- Table structure for t_machine_script
CREATE TABLE `t_machine_script` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(255) NOT NULL COMMENT '脚本名',
 `machine_id` bigint(19) NOT NULL COMMENT '机器id[0:公共]',
 `script` text COMMENT '脚本内容',
 `params` varchar(512) COMMENT '脚本入参',
 `description` varchar(255) COMMENT '脚本描述',
 `type` tinyint(3) COMMENT '脚本类型[1: 有结果；2：无结果；3：实时交互]',
 `creator_id` bigint(19),
 `creator` varchar(32),
 `modifier_id` bigint(19),
 `modifier` varchar(255),
 `create_time` datetime,
 `update_time` datetime,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '机器脚本';

-- Data for table t_machine_script
INSERT INTO `t_machine_script` (`id`,`name`,`machine_id`,`script`,`params`,`description`,`type`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1,'sys_info',9999999,'# 获取系统cpu信息
function get_cpu_info() {
  Physical_CPUs=$(grep "physical id" /proc/cpuinfo | sort | uniq | wc -l)
  Virt_CPUs=$(grep "processor" /proc/cpuinfo | wc -l)
  CPU_Kernels=$(grep "cores" /proc/cpuinfo | uniq | awk -F '': '' ''{print $2}'')
  CPU_Type=$(grep "model name" /proc/cpuinfo | awk -F '': '' ''{print $2}'' | sort | uniq)
  CPU_Arch=$(uname -m)
  echo -e ''\n-------------------------- CPU信息 --------------------------''
  cat <<EOF | column -t
物理CPU个数: $Physical_CPUs
逻辑CPU个数: $Virt_CPUs
每CPU核心数: $CPU_Kernels
CPU型号: $CPU_Type
CPU架构: $CPU_Arch
EOF
}

# 获取系统信息
function get_systatus_info() {
  sys_os=$(uname -o)
  sys_release=$(cat /etc/redhat-release)
  sys_kernel=$(uname -r)
  sys_hostname=$(hostname)
  sys_selinux=$(getenforce)
  sys_lang=$(echo $LANG)
  sys_lastreboot=$(who -b | awk ''{print $3,$4}'')
  echo -e ''-------------------------- 系统信息 --------------------------''
  cat <<EOF | column -t
系统: ${sys_os}
发行版本:   ${sys_release}
系统内核:   ${sys_kernel}
主机名:    ${sys_hostname}
selinux状态:  ${sys_selinux}
系统语言:   ${sys_lang}
系统最后重启时间:   ${sys_lastreboot}
EOF
}

get_systatus_info
#echo -e "\n"
get_cpu_info',NULL,'获取系统信息',1,NULL,NULL,NULL,NULL,NULL,NULL,0,NULL);
INSERT INTO `t_machine_script` (`id`,`name`,`machine_id`,`script`,`params`,`description`,`type`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (2,'get_process_by_name',9999999,'#! /bin/bash
# Function: 根据输入的程序的名字过滤出所对应的PID，并显示出详细信息，如果有几个PID，则全部显示
NAME={{.processName}}
N=`ps -aux | grep $NAME | grep -v grep | wc -l`    ##统计进程总数
if [ $N -le 0 ];then
  echo "无该进程！"
fi
i=1
while [ $N -gt 0 ]
do
  echo "进程PID: `ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $2}''`"
  echo "进程命令：`ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $11}''`"
  echo "进程所属用户: `ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $1}''`"
  echo "CPU占用率：`ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $3}''`%"
  echo "内存占用率：`ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $4}''`%"
  echo "进程开始运行的时刻：`ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $9}''`"
  echo "进程运行的时间：`  ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $11}''`"
  echo "进程状态：`ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $8}''`"
  echo "进程虚拟内存：`ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $5}''`"
  echo "进程共享内存：`ps -aux | grep $NAME | grep -v grep | awk ''NR==''$i''{print $0}''| awk ''{print $6}''`"
  echo "***************************************************************"
  let N-- i++
done','[{"name": "进程名","model": "processName", "placeholder": "请输入进程名"}]','获取进程运行状态',1,NULL,NULL,1,'admin',NULL,'2021-07-12 15:33:41',0,NULL);
INSERT INTO `t_machine_script` (`id`,`name`,`machine_id`,`script`,`params`,`description`,`type`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (3,'sys_run_info',9999999,'#!/bin/bash
# 获取要监控的本地服务器IP地址
IP=`ifconfig | grep inet | grep -vE ''inet6|127.0.0.1'' | awk ''{print $2}''`
echo "IP地址："$IP
 
# 获取cpu总核数
cpu_num=`grep -c "model name" /proc/cpuinfo`
echo "cpu总核数："$cpu_num
 
# 1、获取CPU利用率
################################################
#us 用户空间占用CPU百分比
#sy 内核空间占用CPU百分比
#ni 用户进程空间内改变过优先级的进程占用CPU百分比
#id 空闲CPU百分比
#wa 等待输入输出的CPU时间百分比
#hi 硬件中断
#si 软件中断
#################################################
# 获取用户空间占用CPU百分比
cpu_user=`top -b -n 1 | grep Cpu | awk ''{print $2}'' | cut -f 1 -d "%"`
echo "用户空间占用CPU百分比："$cpu_user
 
# 获取内核空间占用CPU百分比
cpu_system=`top -b -n 1 | grep Cpu | awk ''{print $4}'' | cut -f 1 -d "%"`
echo "内核空间占用CPU百分比："$cpu_system
 
# 获取空闲CPU百分比
cpu_idle=`top -b -n 1 | grep Cpu | awk ''{print $8}'' | cut -f 1 -d "%"`
echo "空闲CPU百分比："$cpu_idle
 
# 获取等待输入输出占CPU百分比
cpu_iowait=`top -b -n 1 | grep Cpu | awk ''{print $10}'' | cut -f 1 -d "%"`
echo "等待输入输出占CPU百分比："$cpu_iowait
 
#2、获取CPU上下文切换和中断次数
# 获取CPU中断次数
cpu_interrupt=`vmstat -n 1 1 | sed -n 3p | awk ''{print $11}''`
echo "CPU中断次数："$cpu_interrupt
 
# 获取CPU上下文切换次数
cpu_context_switch=`vmstat -n 1 1 | sed -n 3p | awk ''{print $12}''`
echo "CPU上下文切换次数："$cpu_context_switch
 
#3、获取CPU负载信息
# 获取CPU15分钟前到现在的负载平均值
cpu_load_15min=`uptime | awk ''{print $11}'' | cut -f 1 -d '',''`
echo "CPU 15分钟前到现在的负载平均值："$cpu_load_15min
 
# 获取CPU5分钟前到现在的负载平均值
cpu_load_5min=`uptime | awk ''{print $10}'' | cut -f 1 -d '',''`
echo "CPU 5分钟前到现在的负载平均值："$cpu_load_5min
 
# 获取CPU1分钟前到现在的负载平均值
cpu_load_1min=`uptime | awk ''{print $9}'' | cut -f 1 -d '',''`
echo "CPU 1分钟前到现在的负载平均值："$cpu_load_1min
 
# 获取任务队列(就绪状态等待的进程数)
cpu_task_length=`vmstat -n 1 1 | sed -n 3p | awk ''{print $1}''`
echo "CPU任务队列长度："$cpu_task_length
 
#4、获取内存信息
# 获取物理内存总量
mem_total=`free -h | grep Mem | awk ''{print $2}''`
echo "物理内存总量："$mem_total
 
# 获取操作系统已使用内存总量
mem_sys_used=`free -h | grep Mem | awk ''{print $3}''`
echo "已使用内存总量(操作系统)："$mem_sys_used
 
# 获取操作系统未使用内存总量
mem_sys_free=`free -h | grep Mem | awk ''{print $4}''`
echo "剩余内存总量(操作系统)："$mem_sys_free
 
# 获取应用程序已使用的内存总量
mem_user_used=`free | sed -n 3p | awk ''{print $3}''`
echo "已使用内存总量(应用程序)："$mem_user_used
 
# 获取应用程序未使用内存总量
mem_user_free=`free | sed -n 3p | awk ''{print $4}''`
echo "剩余内存总量(应用程序)："$mem_user_free
 
# 获取交换分区总大小
mem_swap_total=`free | grep Swap | awk ''{print $2}''`
echo "交换分区总大小："$mem_swap_total
 
# 获取已使用交换分区大小
mem_swap_used=`free | grep Swap | awk ''{print $3}''`
echo "已使用交换分区大小："$mem_swap_used
 
# 获取剩余交换分区大小
mem_swap_free=`free | grep Swap | awk ''{print $4}''`
echo "剩余交换分区大小："$mem_swap_free',NULL,'获取cpu、内存等系统运行状态',1,NULL,NULL,NULL,NULL,NULL,'2021-04-25 15:07:16',0,NULL);
INSERT INTO `t_machine_script` (`id`,`name`,`machine_id`,`script`,`params`,`description`,`type`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (4,'top',9999999,'top',NULL,'实时获取系统运行状态',3,NULL,NULL,1,'admin',NULL,'2021-05-24 15:58:20',0,NULL);
INSERT INTO `t_machine_script` (`id`,`name`,`machine_id`,`script`,`params`,`description`,`type`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (5,'disk-mem',9999999,'df -h','','磁盘空间查看',1,1,'admin',1,'admin','2021-07-16 10:49:53','2021-07-16 10:49:53',0,NULL);

-- Table structure for t_machine_term_op
CREATE TABLE `t_machine_term_op` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
 `machine_id` bigint(19) NOT NULL COMMENT '机器id',
 `username` varchar(60) COMMENT '登录用户名',
 `file_key` varchar(36) COMMENT '文件',
 `exec_cmds` text COMMENT '执行的命令记录',
 `creator_id` bigint(20),
 `creator` varchar(191),
 `create_time` datetime NOT NULL,
 `end_time` datetime,
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '机器终端操作记录表';

-- Data for table t_machine_term_op
-- No data

-- Table structure for t_mongo
CREATE TABLE `t_mongo` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `code` varchar(32) COMMENT 'code',
 `name` varchar(36) NOT NULL COMMENT '名称',
 `uri` varchar(255) NOT NULL COMMENT '连接uri',
 `ssh_tunnel_machine_id` bigint(19) COMMENT 'ssh隧道的机器id',
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19),
 `creator` varchar(36),
 `update_time` datetime,
 `modifier_id` bigint(19),
 `modifier` varchar(36),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);

-- Data for table t_mongo
-- No data

-- Table structure for t_oauth2_account
CREATE TABLE `t_oauth2_account` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `account_id` bigint(19) NOT NULL,
 `identity` varchar(64),
 `create_time` datetime NOT NULL,
 `update_time` datetime NOT NULL,
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT 'oauth2关联账号';

-- Data for table t_oauth2_account
-- No data

-- Table structure for t_redis
CREATE TABLE `t_redis` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `code` varchar(32) COMMENT 'code',
 `name` varchar(255) COMMENT '名称',
 `host` varchar(255) NOT NULL,
 `db` varchar(64) COMMENT '库号: 多个库用,分割',
 `mode` varchar(32),
 `ssh_tunnel_machine_id` bigint(19) COMMENT 'ssh隧道的机器id',
 `remark` varchar(125),
 `creator` varchar(32),
 `creator_id` bigint(19),
 `create_time` datetime,
 `modifier` varchar(32),
 `modifier_id` bigint(19),
 `update_time` datetime,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT 'redis信息';

-- Data for table t_redis
-- No data

-- Table structure for t_resource_auth_cert
CREATE TABLE `t_resource_auth_cert` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(100) COMMENT '账号名称',
 `resource_code` varchar(36) COMMENT '资源编码',
 `resource_type` tinyint(3) NOT NULL COMMENT '资源类型',
 `type` tinyint(3) COMMENT '凭证类型',
 `username` varchar(100) COMMENT '用户名',
 `ciphertext` varchar(5000) COMMENT '密文内容',
 `ciphertext_type` tinyint(3) NOT NULL COMMENT '密文类型（-1.公共授权凭证 1.密码 2.秘钥）',
 `extra` varchar(200) COMMENT '账号需要的其他额外信息（如秘钥口令等）',
 `remark` varchar(50) COMMENT '备注',
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(36) NOT NULL,
 `is_deleted` tinyint(3) DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '资源授权凭证表';
ALTER TABLE `t_resource_auth_cert` ADD  INDEX `idx_name`(`name`) USING BTREE;
ALTER TABLE `t_resource_auth_cert` ADD  INDEX `idx_resource_code`(`resource_code`) USING BTREE;

-- Data for table t_resource_auth_cert
-- No data

-- Table structure for t_resource_op_log
CREATE TABLE `t_resource_op_log` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `code_path` varchar(600) NOT NULL COMMENT '资源标签路径',
 `resource_code` varchar(32) NOT NULL COMMENT '资源编号',
 `resource_type` tinyint(3) NOT NULL COMMENT '资源类型',
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36) NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '资源操作记录';
ALTER TABLE `t_resource_op_log` ADD  INDEX `idx_resource_code`(`resource_code`) USING BTREE;

-- Data for table t_resource_op_log
-- No data

-- Table structure for t_sys_account
CREATE TABLE `t_sys_account` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(30) NOT NULL,
 `username` varchar(30) NOT NULL,
 `password` varchar(64) NOT NULL,
 `status` tinyint(3),
 `otp_secret` varchar(100) COMMENT 'otp秘钥',
 `last_login_time` datetime COMMENT '最后登录时间',
 `last_login_ip` varchar(50),
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(12) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(12) NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '账号信息表';

-- Data for table t_sys_account
INSERT INTO `t_sys_account` (`id`,`name`,`username`,`password`,`status`,`otp_secret`,`last_login_time`,`last_login_ip`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (1,'管理员','admin','$2a$10$w3Wky2U.tinvR7c/s0aKPuwZsIu6pM1/DMJalwBDMbE6niHIxVrrm',1,'','2022-10-26 20:03:48','::1','2020-01-01 19:00:00',1,'admin','2020-01-01 19:00:00',1,'admin',0,NULL);

-- Table structure for t_sys_account_role
CREATE TABLE `t_sys_account_role` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Id',
 `account_id` bigint(19) NOT NULL COMMENT '账号id',
 `role_id` bigint(19) NOT NULL COMMENT '角色id',
 `creator` varchar(45),
 `creator_id` bigint(20),
 `create_time` datetime NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '账号角色关联表';

-- Data for table t_sys_account_role
-- No data

-- Table structure for t_sys_config
CREATE TABLE `t_sys_config` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(60) NOT NULL COMMENT '配置名',
 `key` varchar(120) NOT NULL COMMENT '配置key',
 `params` varchar(1500),
 `value` varchar(1500) COMMENT '配置value',
 `remark` varchar(255) COMMENT '备注',
 `permission` varchar(255) DEFAULT all COMMENT '操作权限',
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(36) NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
);

-- Data for table t_sys_config
INSERT INTO `t_sys_config` (`id`,`name`,`key`,`params`,`value`,`remark`,`permission`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (1,'system.sysconf.accountLoginConf','AccountLoginSecurity','[{"name":"system.sysconf.useCaptcha","model":"useCaptcha","placeholder":"system.sysconf.useCaptchaPlaceholder","options":"true,false"},{"name":"system.sysconf.useOtp","model":"useOtp","placeholder":"system.sysconf.useOtpPlaceholder","options":"true,false"},{"name":"system.sysconf.otpIssuer","model":"otpIssuer","placeholder":""},{"name":"system.sysconf.loginFailCount","model":"loginFailCount","placeholder":"system.sysconf.loginFailCountPlaceholder"},{"name":"system.sysconf.loginFainMin","model":"loginFailMin","placeholder":"system.sysconf.loginFailMinPlaceholder"}]','{"useCaptcha":"true","useOtp":"false","loginFailCount":"5","loginFailMin":"10","otpIssuer":"mayfly-go"}','system.sysconf.accountLoginConfRemark','all','2023-06-17 11:02:11',1,'admin','2023-06-17 14:18:07',1,'admin',0,NULL);
INSERT INTO `t_sys_config` (`id`,`name`,`key`,`params`,`value`,`remark`,`permission`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (2,'system.sysconf.oauth2LoginConf','Oauth2Login','[{"name":"system.sysconf.oauth2Enable","model":"enable","placeholder":"system.sysconf.oauth2EnablePlaceholder","options":"true,false"},{"name":"system.sysconf.name","model":"name","placeholder":"system.sysconf.namePlaceholder"},{"name":"system.sysconf.clientId","model":"clientId","placeholder":"system.sysconf.clientIdPlaceholder"},{"name":"system.sysconf.clientSecret","model":"clientSecret","placeholder":"system.sysconf.clientSecretPlaceholder"},{"name":"system.sysconf.authorizationUrl","model":"authorizationURL","placeholder":"system.sysconf.authorizationUrlPlaceholder"},{"name":"system.sysconf.accessTokenUrl","model":"accessTokenURL","placeholder":"system.sysconf.accessTokenUrlPlaceholder"},{"name":"system.sysconf.redirectUrl","model":"redirectURL","placeholder":"system.sysconf.redirectUrlPlaceholder"},{"name":"system.sysconf.scope","model":"scopes","placeholder":"system.sysconf.scopePlaceholder"},{"name":"system.sysconf.resourceUrl","model":"resourceURL","placeholder":"system.sysconf.resourceUrlPlaceholder"},{"name":"system.sysconf.userId","model":"userIdentifier","placeholder":"system.sysconf.userIdPlaceholder"},{"name":"system.sysconf.autoRegister","model":"autoRegister","placeholder":"","options":"true,false"}]','','system.sysconf.oauth2LoginConfRemark','admin,','2023-07-22 13:58:51',1,'admin','2023-07-22 19:34:37',1,'admin',0,NULL);
INSERT INTO `t_sys_config` (`id`,`name`,`key`,`params`,`value`,`remark`,`permission`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (3,'system.sysconf.ldapLoginConf','LdapLogin','[{"name":"system.sysconf.ldapEnable","model":"enable","placeholder":"system.sysconf.dapEnablePlaceholder","options":"true,false"},{"name":"system.sysconf.host","model":"host","placeholder":"system.sysconf.host"},{"name":"system.sysconf.port","model":"port","placeholder":"system.sysconf.port"},{"name":"system.sysconf.bindDN","model":"bindDN","placeholder":"system.sysconf.bindDnPlaceholder"},{"name":"system.sysconf.bindPwd","model":"bindPwd","placeholder":"system.sysconf.bindPwdPlaceholder"},{"name":"system.sysconf.baseDN","model":"baseDN","placeholder":"system.sysconf.baseDnPlaceholder"},{"name":"system.sysconf.userFilter","model":"userFilter","placeholder":"system.sysconf.userFilerPlaceholder"},{"name":"system.sysconf.uidMap","model":"uidMap","placeholder":"system.sysconf.uidMapPlaceholder"},{"name":"system.sysconf.udnMap","model":"udnMap","placeholder":"system.sysconf.udnMapPlaceholder"},{"name":"system.sysconf.emailMap","model":"emailMap","placeholder":"system.sysconf.emailMapPlaceholder"},{"name":"system.sysconf.skipTlsVerfify","model":"skipTLSVerify","placeholder":"system.sysconf.skipTlsVerfifyPlaceholder","options":"true,false"},{"name":"system.sysconf.securityProtocol","model":"securityProtocol","placeholder":"system.sysconf.securityProtocolPlaceholder","options":"Null,StartTLS,LDAPS"}]','','system.sysconf.ldapLoginConfRemark','admin,','2023-08-25 21:47:20',1,'admin','2023-08-25 22:56:07',1,'admin',0,NULL);
INSERT INTO `t_sys_config` (`id`,`name`,`key`,`params`,`value`,`remark`,`permission`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (4,'system.sysconf.systemConf','SysStyleConfig','[{"model":"logoIcon","name":"system.sysconf.logoIcon","placeholder":"system.sysconf.logoIconPlaceholder","required":false},{"model":"title","name":"system.sysconf.title","placeholder":"system.sysconf.titlePlaceholder","required":false},{"model":"viceTitle","name":"system.sysconf.viceTitle","placeholder":"system.sysconf.viceTitlePlaceholder","required":false},{"model":"useWatermark","name":"system.sysconf.useWatermark","placeholder":"system.sysconf.useWatermarkPlaceholder","options":"true,false","required":false},{"model":"watermarkContent","name":"system.sysconf.watermarkContent","placeholder":"system.sysconf.watermarkContentPlaceholder","required":false}]','{"title":"mayfly-go","viceTitle":"mayfly-go","logoIcon":"","useWatermark":"true","watermarkContent":""}','system.sysconf.systemConfRemark','all','2024-01-04 15:17:18',1,'admin','2024-01-05 09:40:44',1,'admin',0,NULL);
INSERT INTO `t_sys_config` (`id`,`name`,`key`,`params`,`value`,`remark`,`permission`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (5,'system.sysconf.machineConf','MachineConfig','[{"name":"system.sysconf.uploadMaxFileSize","model":"uploadMaxFileSize","placeholder":"system.sysconf.uploadMaxFileSizePlaceholder"},{"model":"termOpSaveDays","name":"system.sysconf.termOpSaveDays","placeholder":"system.sysconf.termOpSaveDaysPlaceholder"},{"model":"guacdHost","name":"system.sysconf.guacdHost","placeholder":"system.sysconf.guacdHostPlaceholder","required":false},{"name":"system.sysconf.guacdPort","model":"guacdPort","placeholder":"system.sysconf.guacdPortPlaceholder","required":false},{"model":"guacdFilePath","name":"system.sysconf.guacdFilePath","placeholder":"system.sysconf.guacdFilePathPlaceholder"}]','{"uploadMaxFileSize":"1000MB","termOpSaveDays":"30","guacdHost":"","guacdPort":"","guacdFilePath":"./guacd/rdp-file"}','system.sysconf.machineConfRemark','all','2023-07-13 16:26:44',1,'admin','2024-10-21 17:02:55',1,'admin',0,NULL);
INSERT INTO `t_sys_config` (`id`,`name`,`key`,`params`,`value`,`remark`,`permission`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (8,'system.sysconf.dbmsConf','DbmsConfig','[{"model":"querySqlSave","name":"system.sysconf.recordQuerySql","placeholder":"system.sysconf.recordQuerySqlPlaceholder","options":"true,false"},{"model":"maxResultSet","name":"system.sysconf.maxResultSet","placeholder":"system.sysconf.maxResultSetPlaceholder","options":""},{"model":"sqlExecTl","name":"system.sysconf.sqlExecLimt","placeholder":"system.sysconf.sqlExecLimtPlaceholder"}]','{"querySqlSave":"false","maxResultSet":"0","sqlExecTl":"60"}','system.sysconf.dbmsConfRemark','admin,','2024-03-06 13:30:51',1,'admin','2024-03-06 14:07:16',1,'admin',0,NULL);
INSERT INTO `t_sys_config` (`id`,`name`,`key`,`params`,`value`,`remark`,`permission`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (9,'system.sysconf.fileConf','FileConfig','[{"model":"basePath","name":"system.sysconf.basePath","placeholder":"system.sysconf.baesPathPlaceholder"}]','{"basePath":"./file"}','system.sysconf.fileConfRemark','admin,','2024-10-20 22:30:01',1,'admin','2024-10-21 13:51:17',1,'admin',0,NULL);

-- Table structure for t_sys_file
CREATE TABLE `t_sys_file` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `file_key` varchar(32) NOT NULL COMMENT 'key',
 `filename` varchar(255) NOT NULL COMMENT '文件名',
 `path` varchar(555) NOT NULL COMMENT '文件路径',
 `size` int(10) COMMENT '文件大小',
 `creator_id` bigint(19),
 `creator` varchar(32),
 `modifier_id` bigint(19),
 `modifier` varchar(255),
 `create_time` datetime,
 `update_time` datetime,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '系统文件表';
ALTER TABLE `t_sys_file` ADD  INDEX `idx_file_key`(`file_key`) USING BTREE;

-- Data for table t_sys_file
-- No data

-- Table structure for t_sys_log
CREATE TABLE `t_sys_log` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `type` tinyint(3) NOT NULL COMMENT '类型',
 `description` varchar(255) COMMENT '描述',
 `req_param` text COMMENT '请求信息',
 `resp` text COMMENT '响应信息',
 `creator` varchar(36) NOT NULL COMMENT '调用者',
 `creator_id` bigint(19) NOT NULL COMMENT '调用者id',
 `create_time` datetime NOT NULL COMMENT '操作时间',
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime,
 `extra` text, 
PRIMARY KEY (id)
) COMMENT '系统操作日志';
ALTER TABLE `t_sys_log` ADD  INDEX `idx_creator_id`(`creator_id`) USING BTREE;

-- Data for table t_sys_log
-- No data

-- Table structure for t_sys_msg
CREATE TABLE `t_sys_msg` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `type` int(10),
 `msg` varchar(2000) NOT NULL,
 `recipient_id` bigint(19) COMMENT '接收人id，-1为所有接收',
 `creator_id` bigint(19),
 `creator` varchar(36),
 `create_time` datetime NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '系统消息表';

-- Data for table t_sys_msg
-- No data

-- Table structure for t_sys_resource
CREATE TABLE `t_sys_resource` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `pid` int(10) NOT NULL COMMENT '父节点id',
 `ui_path` varchar(200) COMMENT '唯一标识路径',
 `type` tinyint(3) NOT NULL COMMENT '1：菜单路由；2：资源（按钮等）',
 `status` int(10) NOT NULL COMMENT '状态；1:可用，-1:禁用',
 `name` varchar(255) NOT NULL COMMENT '名称',
 `code` varchar(255) COMMENT '菜单路由为path，其他为唯一标识',
 `weight` int(10) COMMENT '权重顺序',
 `meta` varchar(455) COMMENT '元数据',
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(255) NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(255) NOT NULL,
 `create_time` datetime NOT NULL,
 `update_time` datetime NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '资源表';
ALTER TABLE `t_sys_resource` ADD  INDEX `t_sys_resource_ui_path_IDX`(`ui_path`) USING BTREE;

-- Data for table t_sys_resource
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1,0,'Aexqq77l/',1,1,'menu.index','/home',10000000,'{"component":"home/Home","icon":"HomeFilled","isAffix":true,"routeName":"Home"}',1,'admin',1,'admin','2021-05-25 16:44:41','2024-11-06 16:18:09',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (2,0,'12sSjal1/',1,1,'menu.machine','/machine',49999998,'{"icon":"Monitor","isKeepAlive":true,"redirect":"machine/list","routeName":"Machine"}',1,'admin',1,'admin','2021-05-25 16:48:16','2024-11-07 11:29:10',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (3,2,'12sSjal1/lskeiql1/',1,1,'menu.machineList','machines',20000000,'{"component":"ops/machine/MachineList","icon":"Monitor","isKeepAlive":true,"routeName":"MachineList"}',2,'admin',1,'admin','2021-05-25 16:50:04','2024-11-07 11:30:53',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (4,0,'Xlqig32x/',1,1,'menu.system','/sys',60000001,'{"icon":"Setting","isKeepAlive":true,"redirect":"/sys/resources","routeName":"sys"}',1,'admin',1,'admin','2021-05-26 15:20:20','2024-11-07 14:03:22',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (5,4,'Xlqig32x/UGxla231/',1,1,'menu.menuPermission','resources',9999998,'{"component":"system/resource/ResourceList","icon":"Menu","isKeepAlive":true,"routeName":"ResourceList"}',1,'admin',1,'admin','2021-05-26 15:23:07','2024-11-07 14:03:31',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (11,4,'Xlqig32x/lxqSiae1/',1,1,'menu.role','roles',10000001,'{"component":"system/role/RoleList","icon":"icon menu/role","isKeepAlive":true,"routeName":"RoleList"}',1,'admin',1,'admin','2021-05-27 11:15:35','2024-11-07 14:12:40',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (12,3,'12sSjal1/lskeiql1/Alw1Xkq3/',2,1,'menu.machineTerminal','machine:terminal',40000000,'null',1,'admin',1,'admin','2021-05-28 14:06:02','2024-11-07 11:54:29',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (14,4,'Xlqig32x/sfslfel/',1,1,'menu.account','accounts',9999999,'{"component":"system/account/AccountList","icon":"User","isKeepAlive":true,"routeName":"AccountList"}',1,'admin',1,'admin','2021-05-28 14:56:25','2024-11-07 14:11:08',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (15,3,'12sSjal1/lskeiql1/Lsew24Kx/',2,1,'menu.machineFileConf','machine:file',50000000,'null',1,'admin',1,'admin','2021-05-31 17:44:37','2024-11-07 11:54:40',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (16,3,'12sSjal1/lskeiql1/exIsqL31/',2,1,'menu.machineCreate','machine:add',10000000,'null',1,'admin',1,'admin','2021-05-31 17:46:11','2024-11-07 11:53:52',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (17,3,'12sSjal1/lskeiql1/Liwakg2x/',2,1,'menu.machineEdit','machine:update',20000000,'null',1,'admin',1,'admin','2021-05-31 17:46:23','2024-11-07 11:54:07',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (18,3,'12sSjal1/lskeiql1/Lieakenx/',2,1,'menu.machineDelete','machine:del',30000000,'null',1,'admin',1,'admin','2021-05-31 17:46:36','2024-11-07 11:54:17',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (19,14,'Xlqig32x/sfslfel/UUiex2xA/',2,1,'menu.accountRoleAllocation','account:saveRoles',50000001,'null',1,'admin',1,'admin','2021-05-31 17:50:51','2024-11-07 14:12:26',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (20,11,'Xlqig32x/lxqSiae1/EMq2Kxq3/',2,1,'menu.roleMenuPermissionAllocation','role:saveResources',40000002,'null',1,'admin',1,'admin','2021-05-31 17:51:41','2024-11-07 14:13:47',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (21,14,'Xlqig32x/sfslfel/Uexax2xA/',2,1,'menu.accountDelete','account:del',20000000,'null',1,'admin',1,'admin','2021-05-31 18:02:01','2024-11-07 14:12:05',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (22,11,'Xlqig32x/lxqSiae1/Elxq2Kxq3/',2,1,'menu.roleDelete','role:del',40000001,'null',1,'admin',1,'admin','2021-05-31 18:02:29','2024-11-07 14:13:38',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (23,11,'Xlqig32x/lxqSiae1/342xKxq3/',2,1,'menu.roleAdd','role:add',19999999,'null',1,'admin',1,'admin','2021-05-31 18:02:44','2024-11-07 14:13:14',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (24,11,'Xlqig32x/lxqSiae1/LexqKxq3/',2,1,'menu.roleEdit','role:update',40000000,'null',1,'admin',1,'admin','2021-05-31 18:02:57','2024-11-07 14:13:26',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (25,5,'Xlqig32x/UGxla231/Elxq23XK/',2,1,'menu.menuPermissionAdd','resource:add',10000000,'null',1,'admin',1,'admin','2021-05-31 18:03:33','2024-11-07 14:03:44',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (26,5,'Xlqig32x/UGxla231/eloq23XK/',2,1,'menu.menuPermissionDelete','resource:delete',30000001,'null',1,'admin',1,'admin','2021-05-31 18:03:47','2024-11-07 14:04:19',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (27,5,'Xlqig32x/UGxla231/JExq23XK/',2,1,'menu.menuPermissionEdit','resource:update',30000000,'null',1,'admin',1,'admin','2021-05-31 18:04:03','2024-11-07 14:04:10',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (28,5,'Xlqig32x/UGxla231/Elex13XK/',2,1,'menu.menuPermissionEnableDisable','resource:changeStatus',40000000,'null',1,'admin',1,'admin','2021-05-31 18:04:33','2024-11-07 14:04:53',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (29,14,'Xlqig32x/sfslfel/xlawx2xA/',2,1,'menu.accountAdd','account:add',19999999,'null',1,'admin',1,'admin','2021-05-31 19:23:42','2024-11-07 14:11:46',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (31,14,'Xlqig32x/sfslfel/eubale13/',2,1,'menu.accountBase','account',9999999,'null',1,'admin',1,'admin','2021-05-31 21:25:06','2024-11-07 14:11:27',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (32,5,'Xlqig32x/UGxla231/321q23XK/',2,1,'menu.menuPermissionBase','resource',9999999,'null',1,'admin',1,'admin','2021-05-31 21:25:25','2024-11-07 14:03:59',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (33,11,'Xlqig32x/lxqSiae1/908xKxq3/',2,1,'menu.roleBase','role',9999999,'null',1,'admin',1,'admin','2021-05-31 21:25:40','2024-11-07 14:13:03',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (34,14,'Xlqig32x/sfslfel/32alx2xA/',2,1,'menu.accountEnableDisable','account:changeStatus',50000000,'null',1,'admin',1,'admin','2021-05-31 21:29:48','2024-11-07 14:12:17',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (36,0,'dbms23ax/',1,1,'menu.dbms','/dbms',49999999,'{"icon":"Coin","isKeepAlive":true,"routeName":"DBMS"}',1,'admin',1,'admin','2021-06-01 14:01:33','2024-11-07 13:43:14',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (37,3,'12sSjal1/lskeiql1/Keiqkx4L/',2,1,'menu.machineFileConfCreate','machine:addFile',60000000,'null',1,'admin',1,'admin','2021-06-01 19:54:23','2024-11-07 12:02:41',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (38,36,'dbms23ax/exaeca2x/',1,1,'menu.dbDataOp','sql-exec',10000000,'{"component":"ops/db/SqlExec","icon":"Coin","isKeepAlive":true,"routeName":"SqlExec"}',1,'admin',1,'admin','2021-06-03 09:09:29','2024-11-07 13:43:23',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (39,0,'sl3as23x/',1,1,'menu.personalCenter','/personal',19999999,'{"component":"personal/index","icon":"UserFilled","isHide":true,"isKeepAlive":true,"routeName":"Personal"}',1,'admin',1,'admin','2021-06-03 14:25:35','2024-11-06 20:55:37',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (40,3,'12sSjal1/lskeiql1/Keal2Xke/',2,1,'menu.machineFileCreate','machine:file:add',70000000,'null',1,'admin',1,'admin','2021-06-08 11:06:26','2024-11-07 12:07:17',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (41,3,'12sSjal1/lskeiql1/Ihfs2xaw/',2,1,'menu.machineFileDelete','machine:file:del',80000000,'null',1,'admin',1,'admin','2021-06-08 11:06:49','2024-11-07 12:07:38',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (42,3,'12sSjal1/lskeiql1/3ldkxJDx/',2,1,'menu.machineFileWrite','machine:file:write',90000000,'null',1,'admin',1,'admin','2021-06-08 11:07:27','2024-11-07 12:07:48',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (43,3,'12sSjal1/lskeiql1/Ljewix43/',2,1,'menu.machineFileUpload','machine:file:upload',100000000,'null',1,'admin',1,'admin','2021-06-08 11:07:42','2024-11-07 12:08:00',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (44,3,'12sSjal1/lskeiql1/L12wix43/',2,1,'menu.machineFileConfDelete','machine:file:rm',69999999,'null',1,'admin',1,'admin','2021-06-08 11:08:12','2024-11-07 12:08:17',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (45,3,'12sSjal1/lskeiql1/Ljewisd3/',2,1,'menu.machineScriptSave','machine:script:save',120000000,'null',1,'admin',1,'admin','2021-06-08 11:09:01','2024-11-07 12:17:03',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (46,3,'12sSjal1/lskeiql1/Ljeew43/',2,1,'menu.machineScriptDelete','machine:script:del',130000000,'null',1,'admin',1,'admin','2021-06-08 11:09:27','2024-11-07 12:17:13',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (47,3,'12sSjal1/lskeiql1/ODewix43/',2,1,'menu.machineScriptRun','machine:script:run',140000000,'null',1,'admin',1,'admin','2021-06-08 11:09:50','2024-11-07 12:17:24',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (54,135,'dbms23ax/X0f4BxT0/leix3Axl/',2,1,'menu.dbSave','db:save',1693041086,'null',1,'admin',1,'admin','2021-07-08 17:30:36','2024-11-07 13:45:53',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (55,135,'dbms23ax/X0f4BxT0/ygjL3sxA/',2,1,'menu.dbDelete','db:del',1693041086,'null',1,'admin',1,'admin','2021-07-08 17:30:48','2024-11-07 13:46:06',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (57,3,'12sSjal1/lskeiql1/OJewex43/',2,1,'menu.machineBase','machine',9999999,'null',1,'admin',1,'admin','2021-07-09 10:48:02','2024-11-07 11:53:55',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (58,135,'dbms23ax/X0f4BxT0/AceXe321/',2,1,'menu.dbBase','db',1693041085,'null',1,'admin',1,'admin','2021-07-09 10:48:22','2024-11-07 13:45:25',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (59,38,'dbms23ax/exaeca2x/ealcia23/',2,1,'menu.dbDataOpBase','db:exec',10000000,'null',1,'admin',1,'admin','2021-07-09 10:50:13','2024-11-07 13:43:31',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (60,0,'RedisXq4/',1,1,'menu.redis','/redis',50000001,'{"icon":"icon redis/redis","isKeepAlive":true,"routeName":"RDS"}',1,'admin',1,'admin','2021-07-19 20:15:41','2024-11-07 13:57:51',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (61,60,'RedisXq4/Exitx4al/',1,1,'menu.redisDataOp','data-operation',10000000,'{"component":"ops/redis/DataOperation","icon":"icon icon-redis","isKeepAlive":true,"routeName":"DataOperation"}',1,'admin',1,'admin','2021-07-19 20:17:29','2024-11-07 13:58:03',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (62,61,'RedisXq4/Exitx4al/LSjie321/',2,1,'menu.redisDataOpBase','redis:data',10000000,'null',1,'admin',1,'admin','2021-07-19 20:18:54','2024-11-07 13:58:12',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (63,60,'RedisXq4/Eoaljc12/',1,1,'menu.redisManage','manage',20000000,'{"component":"ops/redis/RedisList","icon":"icon icon-redis","isKeepAlive":true,"routeName":"RedisList"}',1,'admin',1,'admin','2021-07-20 10:48:04','2024-11-07 13:58:52',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (64,63,'RedisXq4/Eoaljc12/IoxqAd31/',2,1,'menu.redisManageBase','redis:manage',10000000,'null',1,'admin',1,'admin','2021-07-20 10:48:26','2024-11-07 13:59:03',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (71,61,'RedisXq4/Exitx4al/IUlxia23/',2,1,'menu.redisDataOpSave','redis:data:save',29999999,'null',1,'admin',1,'admin','2021-08-17 11:20:37','2024-11-07 13:58:24',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (72,3,'12sSjal1/lskeiql1/LIEwix43/',2,1,'menu.machineKillprocess','machine:killprocess',49999999,'null',1,'admin',1,'admin','2021-08-17 11:20:37','2024-11-07 12:03:00',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (79,0,'Mongo452/',1,1,'menu.mongo','/mongo',50000002,'{"icon":"icon mongo/mongo","isKeepAlive":true,"routeName":"Mongo"}',1,'admin',1,'admin','2022-05-13 14:00:41','2024-11-07 13:59:12',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (80,79,'Mongo452/eggago31/',1,1,'menu.mongoDataOp','mongo-data-operation',10000000,'{"component":"ops/mongo/MongoDataOp","icon":"icon icon-mongo","isKeepAlive":true,"routeName":"MongoDataOp"}',1,'admin',1,'admin','2022-05-13 14:03:58','2024-11-07 13:59:23',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (81,80,'Mongo452/eggago31/egjglal3/',2,1,'menu.mongoDataOpBase','mongo:base',10000000,'null',1,'admin',1,'admin','2022-05-13 14:04:16','2024-11-07 13:59:32',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (82,79,'Mongo452/ghxagl43/',1,1,'menu.mongoManage','mongo-manage',20000000,'{"component":"ops/mongo/MongoList","icon":"icon icon-mongo","isKeepAlive":true,"routeName":"MongoList"}',1,'admin',1,'admin','2022-05-16 18:13:06','2024-11-07 14:00:23',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (83,82,'Mongo452/ghxagl43/egljbla3/',2,1,'menu.mongoManageBase','mongo:manage:base',10000000,'null',1,'admin',1,'admin','2022-05-16 18:13:25','2024-11-07 14:00:35',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (84,4,'Xlqig32x/exlaeAlx/',1,1,'menu.opLog','syslogs',20000000,'{"component":"system/syslog/SyslogList","icon":"Tickets","routeName":"SyslogList"}',1,'admin',1,'admin','2022-07-13 19:57:07','2024-11-07 14:15:09',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (85,84,'Xlqig32x/exlaeAlx/3xlqeXql/',2,1,'menu.opLogBase','syslog',10000000,'null',1,'admin',1,'admin','2022-07-13 19:57:55','2024-11-07 14:15:19',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (87,4,'Xlqig32x/Ulxaee23/',1,1,'menu.sysConf','configs',10000002,'{"component":"system/config/ConfigList","icon":"Setting","isKeepAlive":true,"routeName":"ConfigList"}',1,'admin',1,'admin','2022-08-25 22:18:55','2024-11-07 14:14:27',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (88,87,'Xlqig32x/Ulxaee23/exlqguA3/',2,1,'menu.sysConfBase','config:base',10000000,'null',1,'admin',1,'admin','2022-08-25 22:19:35','2024-11-07 14:14:48',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (93,0,'Tag3fhad/',1,1,'menu.tag','/tag',20000001,'{"icon":"CollectionTag","isKeepAlive":true,"routeName":"Tag"}',1,'admin',1,'admin','2022-10-24 15:18:40','2024-11-07 08:45:41',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (94,93,'Tag3fhad/glxajg23/',1,1,'menu.tagTree','tag-trees',10000000,'{"component":"ops/tag/TagTreeList","icon":"CollectionTag","isKeepAlive":true,"routeName":"TagTreeList"}',1,'admin',1,'admin','2022-10-24 15:19:40','2024-11-07 09:46:57',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (95,93,'Tag3fhad/Bjlag32x/',1,1,'menu.team','teams',20000000,'{"component":"ops/tag/TeamList","icon":"UserFilled","isKeepAlive":true,"routeName":"TeamList"}',1,'admin',1,'admin','2022-10-24 15:20:09','2024-11-07 11:24:34',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (96,94,'Tag3fhad/glxajg23/gkxagt23/',2,1,'menu.tagSave','tag:save',10000000,'null',1,'admin',1,'admin','2022-10-24 15:20:40','2024-11-07 11:22:08',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (97,95,'Tag3fhad/Bjlag32x/GJslag32/',2,1,'menu.teamSave','team:save',10000000,'null',1,'admin',1,'admin','2022-10-24 15:20:57','2024-11-07 11:24:43',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (98,94,'Tag3fhad/glxajg23/xjgalte2/',2,1,'menu.tagDelete','tag:del',20000000,'null',1,'admin',1,'admin','2022-10-26 13:58:47','2024-11-07 11:22:16',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (99,95,'Tag3fhad/Bjlag32x/Gguca23x/',2,1,'menu.teamDelete','team:del',20000000,'null',1,'admin',1,'admin','2022-10-26 13:59:06','2024-11-07 11:24:54',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (100,95,'Tag3fhad/Bjlag32x/Lgidsq32/',2,1,'menu.teamMemberAdd','team:member:save',30000000,'null',1,'admin',1,'admin','2022-10-26 13:59:27','2024-11-07 11:25:43',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (101,95,'Tag3fhad/Bjlag32x/Lixaue3G/',2,1,'menu.teamMemberDelete','team:member:del',40000000,'null',1,'admin',1,'admin','2022-10-26 13:59:43','2024-11-07 11:25:54',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (102,95,'Tag3fhad/Bjlag32x/Oygsq3xg/',2,1,'menu.teamTagSave','team:tag:save',50000000,'null',1,'admin',1,'admin','2022-10-26 13:59:57','2024-11-07 11:26:09',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (103,93,'Tag3fhad/exahgl32/',1,1,'menu.authorization','authcerts',19999999,'{"component":"ops/tag/AuthCertList","icon":"Ticket","isKeepAlive":true,"routeName":"AuthCertList"}',1,'admin',1,'admin','2023-02-23 11:36:26','2024-11-07 09:47:42',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (104,103,'Tag3fhad/exahgl32/egxahg24/',2,1,'menu.authorizationBase','authcert',10000000,'null',1,'admin',1,'admin','2023-02-23 11:37:24','2024-11-07 11:23:58',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (105,103,'Tag3fhad/exahgl32/yglxahg2/',2,1,'menu.authorizationSave','authcert:save',20000000,'null',1,'admin',1,'admin','2023-02-23 11:37:54','2024-11-07 11:24:11',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (106,103,'Tag3fhad/exahgl32/Glxag234/',2,1,'menu.authorizationDelete','authcert:del',30000000,'null',1,'admin',1,'admin','2023-02-23 11:38:09','2024-11-07 11:24:21',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (108,61,'RedisXq4/Exitx4al/Gxlagheg/',2,1,'menu.redisDataOpDelete','redis:data:del',30000000,'null',1,'admin',1,'admin','2023-03-14 17:20:00','2024-11-07 13:58:32',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (128,87,'Xlqig32x/Ulxaee23/MoOWr2N0/',2,1,'menu.sysConfSave','config:save',1687315135,'null',1,'admin',1,'admin','2023-06-21 10:38:55','2024-11-07 14:14:59',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (130,2,'12sSjal1/W9XKiabq/',1,1,'menu.machineCronJob','/machine/cron-job',1689646396,'{"component":"ops/machine/cronjob/CronJobList","icon":"AlarmClock","isKeepAlive":true,"routeName":"CronJobList"}',1,'admin',1,'admin','2023-07-18 10:13:16','2024-11-07 12:17:39',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (131,130,'12sSjal1/W9XKiabq/gEOqr2pD/',2,1,'menu.machineCronJobSvae','machine:cronjob:save',1689860087,'null',1,'admin',1,'admin','2023-07-20 21:34:47','2024-11-07 12:17:48',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (132,130,'12sSjal1/W9XKiabq/zxXM23i0/',2,1,'menu.machineCronJobDelete','machine:cronjob:del',1689860102,'null',1,'admin',1,'admin','2023-07-20 21:35:02','2024-11-07 12:18:01',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (133,80,'Mongo452/eggago31/xvpKk36u/',2,1,'menu.mongoDataOpSave','mongo:data:save',1692674943,'null',1,'admin',1,'admin','2023-08-22 11:29:04','2024-11-07 13:59:41',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (134,80,'Mongo452/eggago31/3sblw1Wb/',2,1,'menu.mongoDataOpDelete','mongo:data:del',1692674964,'null',1,'admin',1,'admin','2023-08-22 11:29:24','2024-11-07 14:00:00',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (135,36,'dbms23ax/X0f4BxT0/',1,1,'menu.dbInstance','instances',1693040706,'{"component":"ops/db/InstanceList","icon":"Coin","isKeepAlive":true,"routeName":"InstanceList"}',1,'admin',1,'admin','2023-08-26 09:05:07','2024-11-07 13:43:59',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (136,135,'dbms23ax/X0f4BxT0/D23fUiBr/',2,1,'menu.dbInstanceSave','db:instance:save',1693041001,'null',1,'admin',1,'admin','2023-08-26 09:10:02','2024-11-07 13:44:47',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (137,135,'dbms23ax/X0f4BxT0/mJlBeTCs/',2,1,'menu.dbInstanceBase','db:instance',1693041000,'null',1,'admin',1,'admin','2023-08-26 09:10:55','2024-11-07 13:44:29',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (138,135,'dbms23ax/X0f4BxT0/Sgg8uPwz/',2,1,'menu.dbInstanceDelete','db:instance:del',1693041084,'null',1,'admin',1,'admin','2023-08-26 09:11:24','2024-11-07 13:44:56',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (150,36,'Jra0n7De/',1,1,'menu.dbDataSync','sync',1693040707,'{"component":"ops/db/SyncTaskList","icon":"Refresh","isKeepAlive":true,"routeName":"SyncTaskList"}',12,'liuzongyang',1,'admin','2023-12-22 09:51:34','2024-11-07 13:46:26',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (151,150,'Jra0n7De/uAnHZxEV/',2,1,'menu.dbDataSync','db:sync',1703641202,'null',12,'liuzongyang',1,'admin','2023-12-27 09:40:02','2024-11-07 13:46:37',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (152,150,'Jra0n7De/zvAMo2vk/',2,1,'menu.dbDataSyncSave','db:sync:save',1703641320,'null',12,'liuzongyang',1,'admin','2023-12-27 09:42:00','2024-11-07 13:46:54',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (153,150,'Jra0n7De/pLOA2UYz/',2,1,'menu.dbDataSyncDelete','db:sync:del',1703641342,'null',12,'liuzongyang',1,'admin','2023-12-27 09:42:22','2024-11-07 13:47:06',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (154,150,'Jra0n7De/VBt68CDx/',2,1,'menu.dbDataSyncChangeStatus','db:sync:status',1703641364,'null',12,'liuzongyang',1,'admin','2023-12-27 09:42:45','2024-11-07 13:47:17',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (155,150,'Jra0n7De/PigmSGVg/',2,1,'menu.dbDataSyncLog','db:sync:log',1704266866,'null',12,'liuzongyang',1,'admin','2024-01-03 15:27:47','2024-11-07 13:47:28',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1707206386,2,'PDPt6217/',1,1,'menu.machineOp','machines-op',1,'{"component":"ops/machine/MachineOp","icon":"Monitor","isKeepAlive":true,"routeName":"MachineOp"}',12,'liuzongyang',1,'admin','2024-02-06 15:59:46','2024-11-07 11:29:21',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1707206421,1707206386,'PDPt6217/kQXTYvuM/',2,1,'menu.machineOpBase','machine-op',1707206421,'null',12,'liuzongyang',1,'admin','2024-02-06 16:00:22','2024-11-07 11:30:41',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1708910975,0,'6egfEVYr/',1,1,'menu.flow','/flow',60000000,'{"icon":"List","isKeepAlive":true,"routeName":"flow"}',1,'admin',1,'admin','2024-02-26 09:29:36','2024-11-07 14:01:27',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1708911264,1708910975,'6egfEVYr/fw0Hhvye/',1,1,'menu.flowProcDef','procdefs',1708911264,'{"component":"flow/ProcdefList","icon":"List","isKeepAlive":true,"routeName":"ProcdefList"}',1,'admin',1,'admin','2024-02-26 09:34:24','2024-11-07 14:01:59',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709045735,1708910975,'6egfEVYr/3r3hHEub/',1,1,'menu.myTask','procinst-tasks',1708911263,'{"component":"flow/ProcinstTaskList","icon":"Tickets","isKeepAlive":true,"routeName":"ProcinstTaskList"}',1,'admin',1,'admin','2024-02-27 22:55:35','2024-11-07 14:01:39',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709103180,1708910975,'6egfEVYr/oNCIbynR/',1,1,'menu.myFlow','procinsts',1708911263,'{"component":"flow/ProcinstList","icon":"Tickets","isKeepAlive":true,"routeName":"ProcinstList"}',1,'admin',1,'admin','2024-02-28 14:53:00','2024-11-07 14:01:48',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709194669,36,'SmLcpu6c/',1,1,'menu.dbTransfer','transfer',1709194669,'{"component":"ops/db/DbTransferList","icon":"Switch","isKeepAlive":true,"routeName":"DbTransferList"}',12,'liuzongyang',1,'admin','2024-02-29 16:17:50','2024-11-07 13:47:44',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709194694,1709194669,'SmLcpu6c/A9vAm4J8/',2,1,'menu.dbTransferBase','db:transfer',1709194694,'null',12,'liuzongyang',1,'admin','2024-02-29 16:18:14','2024-11-07 13:47:55',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709196697,1709194669,'SmLcpu6c/5oJwPzNb/',2,1,'menu.dbTransferSave','db:transfer:save',1709196697,'null',12,'liuzongyang',1,'admin','2024-02-29 16:51:37','2024-11-07 13:48:06',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709196707,1709194669,'SmLcpu6c/L3ybnAEW/',2,1,'menu.dbTransferDelete','db:transfer:del',1709196707,'null',12,'liuzongyang',1,'admin','2024-02-29 16:51:47','2024-11-07 13:48:21',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709196723,1709194669,'SmLcpu6c/hGiLN1VT/',2,1,'menu.dbTransferChangeStatus','db:transfer:status',1709196723,'null',12,'liuzongyang',1,'admin','2024-02-29 16:52:04','2024-11-07 13:49:01',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709196737,1709194669,'SmLcpu6c/CZhNIbWg/',2,1,'menu.dbTransferRunLog','db:transfer:log',1709196737,'null',12,'liuzongyang',1,'admin','2024-02-29 16:52:17','2024-11-07 13:49:26',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709196755,1709194669,'SmLcpu6c/b6yHt6V2/',2,1,'menu.dbTransferRun','db:transfer:run',1709196736,'null',12,'liuzongyang',1,'admin','2024-02-29 16:52:36','2024-11-07 13:49:15',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709208339,1708911264,'6egfEVYr/fw0Hhvye/r9ZMTHqC/',2,1,'menu.flowProcDefSave','flow:procdef:save',1709208339,'null',1,'admin',1,'admin','2024-02-29 20:05:40','2024-11-07 14:02:10',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1709208354,1708911264,'6egfEVYr/fw0Hhvye/b4cNf3iq/',2,1,'menu.flowProcDefDelete','flow:procdef:del',1709208354,'null',1,'admin',1,'admin','2024-02-29 20:05:54','2024-11-07 14:03:07',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1712717290,0,'tLb8TKLB/',1,1,'menu.noPagePermission','empty',60000002,'{"component":"empty","icon":"Menu","isHide":true,"isKeepAlive":true,"routeName":"empty"}',1,'admin',1,'admin','2024-04-10 10:48:10','2024-11-07 14:15:31',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1712717337,1712717290,'tLb8TKLB/m2abQkA8/',2,1,'menu.authcertShowciphertext','authcert:showciphertext',1712717337,'null',1,'admin',1,'admin','2024-04-10 10:48:58','2024-11-07 14:15:40',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1713875842,2,'12sSjal1/UnWIUhW0/',1,1,'menu.machineSecurityConfig','security',1713875842,'{"component":"ops/machine/security/SecurityConfList","icon":"Setting","isKeepAlive":true,"routeName":"SecurityConfList"}',1,'admin',1,'admin','2024-04-23 20:37:22','2024-11-07 12:18:13',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1714031981,1713875842,'12sSjal1/UnWIUhW0/tEzIKecl/',2,1,'menu.machineSecurityCmdSvae','cmdconf:save',1714031981,'null',1,'admin',1,'admin','2024-04-25 15:59:41','2024-11-07 12:18:22',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1714032002,1713875842,'12sSjal1/UnWIUhW0/0tJwC3Gf/',2,1,'menu.machineSecurityCmdDelete','cmdconf:del',1714032002,'null',1,'admin',1,'admin','2024-04-25 16:00:02','2024-11-07 12:18:33',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1724376022,1709194669,'SmLcpu6c/HIURtJJA/',2,1,'menu.dbTransferFileDelete','db:transfer:files:del',1724376022,'null',12,'liuzongyang',1,'admin','2024-08-23 09:20:23','2024-11-07 13:49:39',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1724395850,1709194669,'SmLcpu6c/FmqK4azt/',2,1,'menu.dbTransferFileDownload','db:transfer:files:down',1724395850,'null',12,'liuzongyang',1,'admin','2024-08-23 14:50:51','2024-11-07 13:49:58',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1724398262,1709194669,'SmLcpu6c/btVtrbhk/',2,1,'menu.dbTransferFileShow','db:transfer:files',1724376021,'null',12,'liuzongyang',1,'admin','2024-08-23 15:31:02','2024-11-07 13:51:32',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1724998419,1709194669,'SmLcpu6c/qINungml/',2,1,'menu.dbTransferFileRun','db:transfer:files:run',1724998419,'null',12,'liuzongyang',1,'admin','2024-08-30 14:13:39','2024-11-07 13:50:07',0,NULL);
INSERT INTO `t_sys_resource` (`id`,`pid`,`ui_path`,`type`,`status`,`name`,`code`,`weight`,`meta`,`creator_id`,`creator`,`modifier_id`,`modifier`,`create_time`,`update_time`,`is_deleted`,`delete_time`) VALUES (1729668131,38,'dbms23ax/exaeca2x/TGFPA3Ez/',2,1,'menu.dbDataOpSqlScriptRun','db:sqlscript:run',1729668131,'null',1,'admin',1,'admin','2024-10-23 15:22:12','2024-11-07 13:43:46',0,NULL);

-- Table structure for t_sys_role
CREATE TABLE `t_sys_role` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(16) NOT NULL,
 `code` varchar(64) NOT NULL COMMENT '角色code',
 `status` tinyint(3),
 `remark` varchar(255),
 `type` tinyint(3) NOT NULL COMMENT '类型：1:公共角色；2:特殊角色',
 `create_time` datetime,
 `creator_id` bigint(19),
 `creator` varchar(16),
 `update_time` datetime,
 `modifier_id` bigint(19),
 `modifier` varchar(16),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '角色表';

-- Data for table t_sys_role
INSERT INTO `t_sys_role` (`id`,`name`,`code`,`status`,`remark`,`type`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (7,'公共角色','COMMON',1,'所有账号基础角色',1,'2021-07-06 15:05:47',1,'admin','2021-07-06 15:05:47',1,'admin',0,NULL);

-- Table structure for t_sys_role_resource
CREATE TABLE `t_sys_role_resource` (
 `id` bigint(19) NOT NULL AUTO_INCREMENT,
 `role_id` bigint(19) NOT NULL,
 `resource_id` bigint(19) NOT NULL,
 `creator_id` bigint(20),
 `creator` varchar(45),
 `create_time` datetime,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '角色资源关联表';

-- Data for table t_sys_role_resource
-- No data

-- Table structure for t_tag_tree
CREATE TABLE `t_tag_tree` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `type` tinyint(3) NOT NULL DEFAULT -1 COMMENT '类型： -1.普通标签； 1机器  2db 3redis 4mongo',
 `code` varchar(36) NOT NULL COMMENT '标识符',
 `code_path` varchar(800) NOT NULL COMMENT '标识符路径',
 `name` varchar(36) COMMENT '名称',
 `remark` varchar(255),
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(36) NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '标签树';
ALTER TABLE `t_tag_tree` ADD  INDEX `idx_code`(`code`(32)) USING BTREE;
ALTER TABLE `t_tag_tree` ADD  INDEX `idx_code_path`(`code_path`(200)) USING BTREE;

-- Data for table t_tag_tree
INSERT INTO `t_tag_tree` (`id`,`type`,`code`,`code_path`,`name`,`remark`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (1,-1,'default','default/','默认','默认标签','2022-10-26 20:04:19',1,'admin','2022-10-26 20:04:19',1,'admin',0,NULL);

-- Table structure for t_tag_tree_relate
CREATE TABLE `t_tag_tree_relate` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `tag_id` bigint(19) NOT NULL COMMENT '标签树id',
 `relate_id` bigint(19) NOT NULL COMMENT '关联',
 `relate_type` tinyint(3) NOT NULL COMMENT '关联类型',
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(36) NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '与标签树有关联关系的表';
ALTER TABLE `t_tag_tree_relate` ADD  INDEX `idx_tag_id`(`tag_id`) USING BTREE;

-- Data for table t_tag_tree_relate
INSERT INTO `t_tag_tree_relate` (`id`,`tag_id`,`relate_id`,`relate_type`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (1,1,1,1,'2022-10-26 20:04:45',1,'admin','2022-10-26 20:04:45',1,'admin',0,NULL);

-- Table structure for t_team
CREATE TABLE `t_team` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `name` varchar(36) NOT NULL COMMENT '名称',
 `validity_start_date` datetime COMMENT '生效开始时间',
 `validity_end_date` datetime COMMENT '生效结束时间',
 `remark` varchar(255) COMMENT '备注',
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36),
 `update_time` datetime,
 `modifier_id` bigint(19),
 `modifier` varchar(36),
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '团队信息';

-- Data for table t_team
INSERT INTO `t_team` (`id`,`name`,`validity_start_date`,`validity_end_date`,`remark`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (1,'default_team','2024-05-01 00:00:00','2050-05-01 00:00:00','默认团队','2022-10-26 20:04:36',1,'admin','2022-10-26 20:04:36',1,'admin',0,NULL);

-- Table structure for t_team_member
CREATE TABLE `t_team_member` (
 `id` bigint(19) NOT NULL AUTO_INCREMENT,
 `team_id` bigint(19) NOT NULL COMMENT '项目团队id',
 `account_id` bigint(19) NOT NULL COMMENT '成员id',
 `username` varchar(36) NOT NULL,
 `create_time` datetime NOT NULL,
 `creator_id` bigint(19) NOT NULL,
 `creator` varchar(36) NOT NULL,
 `update_time` datetime NOT NULL,
 `modifier_id` bigint(19) NOT NULL,
 `modifier` varchar(36) NOT NULL,
 `is_deleted` tinyint(3) NOT NULL DEFAULT 0,
 `delete_time` datetime, 
PRIMARY KEY (id)
) COMMENT '团队成员表';

-- Data for table t_team_member
INSERT INTO `t_team_member` (`id`,`team_id`,`account_id`,`username`,`create_time`,`creator_id`,`creator`,`update_time`,`modifier_id`,`modifier`,`is_deleted`,`delete_time`) VALUES (1,1,1,'admin','2022-10-26 20:04:36',1,'admin','2022-10-26 20:04:36',1,'admin',0,NULL);

