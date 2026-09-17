-- ============================================
-- 数据库实例统一管理视图
-- 创建时间: 2025-12-02
-- 说明: 联合查询 MySQL/PostgreSQL、Redis、MongoDB、Elasticsearch 四张表
-- ============================================

-- 删除已存在的视图
DROP VIEW IF EXISTS `view_db_instance_all`;

-- 创建统一视图
CREATE VIEW `view_db_instance_all` AS

-- MySQL/PostgreSQL 实例
SELECT
    id,
    code,
    name,
    CASE
        WHEN type = 'postgres' THEN 'postgres'
        ELSE 'mysql'
    END AS db_type,
    type AS sub_type,
    host,
    port,
    username,
    password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'network', network,
        'params', params
    ) AS connection_config
FROM db_instance

UNION ALL

-- Redis 实例
SELECT
    id,
    code,
    name,
    'redis' AS db_type,
    mode AS sub_type,
    host,
    port,
    username,
    password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'mode', mode,
        'db', db,
        'redisNodePassword', redis_node_password
    ) AS connection_config
FROM db_redis_instance

UNION ALL

-- MongoDB 实例
SELECT
    id,
    code,
    name,
    'mongodb' AS db_type,
    type AS sub_type,
    NULL AS host,
    NULL AS port,
    NULL AS username,
    NULL AS password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'uri', uri
    ) AS connection_config
FROM db_mongo_instance

UNION ALL

-- Elasticsearch 实例
SELECT
    id,
    code,
    name,
    'elasticsearch' AS db_type,
    protocol AS sub_type,
    host,
    port,
    username,
    password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'protocol', protocol
    ) AS connection_config
FROM db_es_instance;

-- 验证视图
SELECT
    db_type,
    COUNT(*) as count
FROM view_db_instance_all
GROUP BY db_type;

-- 查询示例
-- SELECT * FROM view_db_instance_all WHERE db_type = 'mysql';
-- SELECT * FROM view_db_instance_all WHERE name LIKE '%生产%';
-- SELECT * FROM view_db_instance_all WHERE status = 1 ORDER BY create_time DESC;
