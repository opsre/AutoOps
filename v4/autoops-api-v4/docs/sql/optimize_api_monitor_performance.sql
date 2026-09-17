-- =====================================================
-- API监控性能优化SQL
-- =====================================================

-- 1. 为 monitor_api_monitor_result 表添加复合索引
-- 这是最关键的优化,可以将查询时间从 100ms+ 降低到 10ms 以内
DROP INDEX IF EXISTS idx_monitor_result_monitor_time ON monitor_api_monitor_result;
CREATE INDEX idx_monitor_result_monitor_time
ON monitor_api_monitor_result(monitor_id, check_time);

-- 2. 添加覆盖索引,包含 success 字段,避免回表查询
DROP INDEX IF EXISTS idx_monitor_result_query ON monitor_api_monitor_result;
CREATE INDEX idx_monitor_result_query
ON monitor_api_monitor_result(monitor_id, check_time, success);

-- 3. 为 monitor_api_monitor 表的 group_id 确保有索引(应该已存在,但确认一下)
CREATE INDEX IF NOT EXISTS idx_monitor_group_id
ON monitor_api_monitor(group_id);

-- 4. 为经常一起查询的字段添加复合索引
CREATE INDEX IF NOT EXISTS idx_monitor_group_status
ON monitor_api_monitor(group_id, status);

-- 5. 查看当前索引情况
SHOW INDEX FROM monitor_api_monitor;
SHOW INDEX FROM monitor_api_monitor_result;

-- 6. 分析表统计信息(帮助MySQL优化器选择正确的索引)
ANALYZE TABLE monitor_api_monitor;
ANALYZE TABLE monitor_api_monitor_result;
