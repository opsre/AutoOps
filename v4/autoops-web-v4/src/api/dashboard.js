/**
 * Dashboard API接口管理
 *
 * @author Generated
 */

import request from "@/utils/request"

// Dashboard统计数据API
export function getDashboardStats() {
    return request({
        url: 'dashboard/stats',
        method: 'get'
    })
}

// 业务分布统计API
export function getBusinessDistribution() {
    return request({
        url: 'dashboard/business-distribution',
        method: 'get'
    })
}

// 资产数量统计API
export function getAssetsStatistics() {
    return request({
        url: 'dashboard/assets',
        method: 'get'
    })
}

// ========== 告警源管理 ==========

// 创建告警源
export function createAlertSource(data) {
    return request({
        url: '/api/v1/monitor/alert-sources',
        method: 'post',
        data
    })
}

// 更新告警源
export function updateAlertSource(data) {
    return request({
        url: '/api/v1/monitor/alert-sources',
        method: 'put',
        data
    })
}

// 删除告警源
export function deleteAlertSource(id) {
    return request({
        url: `/api/v1/monitor/alert-sources/${id}`,
        method: 'delete'
    })
}

// 获取告警源详情
export function getAlertSource(id) {
    return request({
        url: `/api/v1/monitor/alert-sources/${id}`,
        method: 'get'
    })
}

// 获取告警源列表
export function getAlertSourceList(data) {
    return request({
        url: '/api/v1/monitor/alert-sources/list',
        method: 'get',
        data
    })
}

// 批量删除告警源
export function batchDeleteAlertSource(ids) {
    return request({
        url: '/api/v1/monitor/alert-sources/batch-delete',
        method: 'post',
        data: ids
    })
}

// 更新告警源状态
export function updateAlertSourceStatus(id, status) {
    return request({
        url: `/api/v1/monitor/alert-sources/${id}/status`,
        method: 'put',
        data: { status }
    })
}

// ========== 告警标签统计 ==========

// 获取今天的告警统计
export function getAlertStatsToday(sourceId, label) {
    return request({
        url: '/api/v1/monitor/alert-stats/today',
        method: 'get',
        params: { sourceId, label }
    })
}

// 获取昨天的告警统计
export function getAlertStatsYesterday(sourceId, label) {
    return request({
        url: '/api/v1/monitor/alert-stats/yesterday',
        method: 'get',
        params: { sourceId, label }
    })
}

// 获取最近7天的告警统计
export function getAlertStatsLast7Days(sourceId, label) {
    return request({
        url: '/api/v1/monitor/alert-stats/last7days',
        method: 'get',
        params: { sourceId, label }
    })
}

// 获取最近30天的告警统计
export function getAlertStatsLast30Days(sourceId, label) {
    return request({
        url: '/api/v1/monitor/alert-stats/last30days',
        method: 'get',
        params: { sourceId, label }
    })
}

// 今天vs昨天同环比
export function getAlertStatsTodayVsYesterday(sourceId, label) {
    return request({
        url: '/api/v1/monitor/alert-stats/today-vs-yesterday',
        method: 'get',
        params: { sourceId, label }
    })
}

// 本周vs上周同环比
export function getAlertStatsThisWeekVsLastWeek(sourceId, label) {
    return request({
        url: '/api/v1/monitor/alert-stats/thisweek-vs-lastweek',
        method: 'get',
        params: { sourceId, label }
    })
}

// ========== 告警通知统计 ==========

// 获取今天的告警通知统计
export function getAlertNotificationsToday(sourceId) {
    return request({
        url: '/api/v1/monitor/alert-notifications/today',
        method: 'get',
        params: { sourceId }
    })
}

// 获取昨天的告警通知统计
export function getAlertNotificationsYesterday(sourceId) {
    return request({
        url: '/api/v1/monitor/alert-notifications/yesterday',
        method: 'get',
        params: { sourceId }
    })
}

// 获取最近7天的告警通知统计
export function getAlertNotificationsLast7Days(sourceId) {
    return request({
        url: '/api/v1/monitor/alert-notifications/last7days',
        method: 'get',
        params: { sourceId }
    })
}

// 获取最近30天的告警通知统计
export function getAlertNotificationsLast30Days(sourceId) {
    return request({
        url: '/api/v1/monitor/alert-notifications/last30days',
        method: 'get',
        params: { sourceId }
    })
}

// ========== 云账单管理 ==========

// 获取腾讯云账单
export function getTencentBilling(params) {
    return request({
        url: '/api/v1/monitor/billing/tencent',
        method: 'get',
        params
    })
}

// 获取阿里云账单（按产品）
export function getAliyunProductBilling(params) {
    return request({
        url: '/api/v1/monitor/billing/aliyun-product',
        method: 'get',
        params
    })
}

// 获取阿里云账单（按标签）
export function getAliyunTagBilling(params) {
    return request({
        url: '/api/v1/monitor/billing/aliyun-tag',
        method: 'get',
        params
    })
}

// 获取AWS账单
export function getAWSBilling(params) {
    return request({
        url: '/api/v1/monitor/billing/aws',
        method: 'get',
        params
    })
}

// 获取最近12个月费用趋势
export function getLast12MonthsTrend(data) {
    return request({
        url: '/api/v1/monitor/billing/last-12-months-trend',
        method: 'get',
        data
    })
}

// 获取上月费用汇总及同环比
export function getLastMonthSummary(params) {
    return request({
        url: '/api/v1/monitor/billing/last-month-summary',
        method: 'get',
        params
    })
}

// ========== 协作空间统计 ==========

// 获取今天的协作空间统计
export function getChannelStatsToday(sourceId) {
    return request({
        url: '/api/v1/monitor/channel-stats/today',
        method: 'get',
        params: { sourceId }
    })
}

// 获取昨天的协作空间统计
export function getChannelStatsYesterday(sourceId) {
    return request({
        url: '/api/v1/monitor/channel-stats/yesterday',
        method: 'get',
        params: { sourceId }
    })
}

// 获取最近7天的协作空间统计
export function getChannelStatsLast7Days(sourceId) {
    return request({
        url: '/api/v1/monitor/channel-stats/last7days',
        method: 'get',
        params: { sourceId }
    })
}

// 获取最近30天的协作空间统计
export function getChannelStatsLast30Days(sourceId) {
    return request({
        url: '/api/v1/monitor/channel-stats/last30days',
        method: 'get',
        params: { sourceId }
    })
}

// ========== 上线趋势统计 ==========

// 获取上线趋势统计
export function getDeploymentTrend(period = 'week') {
    return request({
        url: '/api/v1/apps/deployment-trend',
        method: 'get',
        params: { period }
    })
}

// 获取主机监控数据
export function getHostMetrics(hostId) {
    return request({
        url: `/api/v1/monitor/host/${hostId}`,
        method: 'get'
    })
}

// ========== Zabbix告警统计 ==========

// 获取告警TOP(按告警名称)
export function getZabbixAlertTopByName(accountId, range, topN = 5) {
    return request({
        url: '/api/v1/monitor/zabbix/alert-top-by-name',
        method: 'get',
        params: { accountId, range, topN }
    })
}

// 获取告警TOP(按主机)
export function getZabbixAlertTopByHost(accountId, range, topN = 5) {
    return request({
        url: '/api/v1/monitor/zabbix/alert-top-by-host',
        method: 'get',
        params: { accountId, range, topN }
    })
}

// 获取告警TOP(按主机群组)
export function getZabbixAlertTopByHostgroup(accountId, range, topN = 5) {
    return request({
        url: '/api/v1/monitor/zabbix/alert-top-by-hostgroup',
        method: 'get',
        params: { accountId, range, topN }
    })
}

// 获取告警通知统计
export function getZabbixNotificationStats(accountId, range) {
    return request({
        url: '/api/v1/monitor/zabbix/notification-stats',
        method: 'get',
        params: { accountId, range }
    })
}

export default {
    // Dashboard统计数据
    getDashboardStats,
    // 业务分布统计
    getBusinessDistribution,
    // 资产数量统计
    getAssetsStatistics,
    // 告警源管理
    createAlertSource,
    updateAlertSource,
    deleteAlertSource,
    getAlertSource,
    getAlertSourceList,
    batchDeleteAlertSource,
    updateAlertSourceStatus,
    // 告警标签统计
    getAlertStatsToday,
    getAlertStatsYesterday,
    getAlertStatsLast7Days,
    getAlertStatsLast30Days,
    getAlertStatsTodayVsYesterday,
    getAlertStatsThisWeekVsLastWeek,
    // 告警通知统计
    getAlertNotificationsToday,
    getAlertNotificationsYesterday,
    getAlertNotificationsLast7Days,
    getAlertNotificationsLast30Days,
    // 云账单管理
    getTencentBilling,
    getAliyunProductBilling,
    getAliyunTagBilling,
    getAWSBilling,
    getLast12MonthsTrend,
    getLastMonthSummary,
    // 协作空间统计
    getChannelStatsToday,
    getChannelStatsYesterday,
    getChannelStatsLast7Days,
    getChannelStatsLast30Days,
    // 上线趋势统计
    getDeploymentTrend,
    // 主机监控数据
    getHostMetrics,
    // Zabbix告警统计
    getZabbixAlertTopByName,
    getZabbixAlertTopByHost,
    getZabbixAlertTopByHostgroup,
    getZabbixNotificationStats
}