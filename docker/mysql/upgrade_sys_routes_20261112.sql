-- Upgrade initialization data from devops001.sql.bak to devops001.sql.
-- This script only inserts or updates sys_* route and permission data.
-- It does not create/drop tables or delete existing rows.

SET NAMES utf8mb4;
START TRANSACTION;

-- Insert new routes and update route labels/order where the latest seed changed.
INSERT INTO `sys_menu`
  (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`)
VALUES
  (293, 80, '凭据管理', 'Key', 'cmdb:credentials', 2, 'cmdb/credentials', 2, 4, '2026-04-03 18:01:20'),
  (358, 80, '数据库管理', 'Eleme', 'cmdb:db-management', 2, 'cmdb/db-management', 2, 7, '2026-05-19 11:43:47'),
  (370, 80, '主机巡检', 'Wallet', 'cmdb:inspection', 2, 'cmdb/inspection', 2, 6, '2026-06-08 17:27:05'),
  (374, 80, '中间件管理', 'Eleme', 'cmdb:middleware', 2, 'cmdb/middleware', 2, 8, '2026-07-11 00:10:18'),
  (394, 393, '告警中心', 'List', 'alert:view', 2, 'alert/view', 2, 1, '2026-07-11 00:31:32'),
  (399, 393, '告警数据源', 'List', 'alert:source', 2, 'alert/source', 2, 7, '2026-07-11 00:36:27'),
  (404, 393, '域名告警', 'Link', 'alert:domain:list', 2, 'alert/domain', 2, 6, '2026-07-30 23:13:51'),
  (413, 383, '数据慢日志', 'Menu', 'logs:mysql-slow', 2, 'logs/query/mysql-slow', 2, 6, '2026-11-10 20:46:03'),
  (414, 383, '业务访问日志', 'Menu', 'logs:nginx', 2, 'logs/nginx', 2, 7, '2026-11-10 20:46:20'),
  (415, 393, '监控中心', 'List', 'alert:monitoring', 2, 'alert/monitoring', 2, 0, '2026-11-10 20:47:00')
ON DUPLICATE KEY UPDATE
  `parent_id` = VALUES(`parent_id`),
  `menu_name` = VALUES(`menu_name`),
  `icon` = VALUES(`icon`),
  `value` = VALUES(`value`),
  `menu_type` = VALUES(`menu_type`),
  `url` = VALUES(`url`),
  `menu_status` = VALUES(`menu_status`),
  `sort` = VALUES(`sort`),
  `create_time` = VALUES(`create_time`);

-- Add the latest permission grants without removing existing grants.
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, 363 FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` = 363
);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, 413 FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` = 413
);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, 414 FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` = 414
);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, 415 FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` = 415
);

-- Seed the two navigation categories used by the navigation page.
INSERT INTO `sys_nav_category`
  (`id`, `name`, `code`, `icon`, `description`, `sort`, `status`, `create_by`, `create_time`, `update_time`)
VALUES
  (1, '公共服务', 'public_service', 'Grid', '公司内部公共系统、项目系统和运维平台入口', 10, 1, 0, '2026-11-10 20:06:35.387', '2026-11-10 20:06:35.387'),
  (2, '个人工具', 'personal_tools', 'Star', '个人常用工具网站和个人链接', 20, 1, 0, '2026-11-10 20:06:35.387', '2026-11-10 20:06:35.387')
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `icon` = VALUES(`icon`),
  `description` = VALUES(`description`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `update_time` = VALUES(`update_time`);

COMMIT;
