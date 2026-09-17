package controller

import (
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ZabbixStatsController struct {
	service *service.ZabbixStatsService
}

func NewZabbixStatsController() *ZabbixStatsController {
	return &ZabbixStatsController{
		service: service.NewZabbixStatsService(),
	}
}

// GetZabbixAlertTopByName 获取Zabbix告警TOP(按告警名称)
// @Summary 获取Zabbix告警TOP(按告警名称)
// @Description 获取指定时间范围内的Zabbix告警TOP5(按告警名称统计)
// @Tags Zabbix监控统计
// @Accept json
// @Produce json
// @Param accountId query uint true "Zabbix账号ID"
// @Param range query string true "时间范围(today/yesterday/last7days/last30days)"
// @Param topN query int false "TOP数量,默认5"
// @Success 200 {object} result.Result{data=model.ZabbixAlertTopResult}
// @Router /api/v1/monitor/zabbix/alert-top-by-name [get]
// @Security ApiKeyAuth
func (ctrl *ZabbixStatsController) GetZabbixAlertTopByName(c *gin.Context) {
	accountIDStr := c.Query("accountId")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的账号ID")
		return
	}

	rangeType := c.DefaultQuery("range", "today")
	topN, _ := strconv.Atoi(c.DefaultQuery("topN", "5"))
	if topN <= 0 {
		topN = 5
	}

	startTime, endTime := getTimeRange(rangeType)

	topResult, err := ctrl.service.GetAlertTopByName(uint(accountID), startTime, endTime, topN, rangeType)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取告警TOP失败: "+err.Error())
		return
	}

	result.Success(c, topResult)
}

// GetZabbixAlertTopByHost 获取Zabbix告警TOP(按主机)
// @Summary 获取Zabbix告警TOP(按主机)
// @Description 获取指定时间范围内的Zabbix告警TOP5(按主机统计)
// @Tags Zabbix监控统计
// @Accept json
// @Produce json
// @Param accountId query uint true "Zabbix账号ID"
// @Param range query string true "时间范围(today/yesterday/last7days/last30days)"
// @Param topN query int false "TOP数量,默认5"
// @Success 200 {object} result.Result{data=model.ZabbixAlertTopResult}
// @Router /api/v1/monitor/zabbix/alert-top-by-host [get]
// @Security ApiKeyAuth
func (ctrl *ZabbixStatsController) GetZabbixAlertTopByHost(c *gin.Context) {
	accountIDStr := c.Query("accountId")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的账号ID")
		return
	}

	rangeType := c.DefaultQuery("range", "today")
	topN, _ := strconv.Atoi(c.DefaultQuery("topN", "5"))
	if topN <= 0 {
		topN = 5
	}

	startTime, endTime := getTimeRange(rangeType)

	topResult, err := ctrl.service.GetAlertTopByHost(uint(accountID), startTime, endTime, topN, rangeType)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取告警TOP失败: "+err.Error())
		return
	}

	result.Success(c, topResult)
}

// GetZabbixAlertTopByHostGroup 获取Zabbix告警TOP(按主机群组)
// @Summary 获取Zabbix告警TOP(按主机群组)
// @Description 获取指定时间范围内的Zabbix告警TOP5(按主机群组统计)
// @Tags Zabbix监控统计
// @Accept json
// @Produce json
// @Param accountId query uint true "Zabbix账号ID"
// @Param range query string true "时间范围(today/yesterday/last7days/last30days)"
// @Param topN query int false "TOP数量,默认5"
// @Success 200 {object} result.Result{data=model.ZabbixAlertTopResult}
// @Router /api/v1/monitor/zabbix/alert-top-by-hostgroup [get]
// @Security ApiKeyAuth
func (ctrl *ZabbixStatsController) GetZabbixAlertTopByHostGroup(c *gin.Context) {
	accountIDStr := c.Query("accountId")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的账号ID")
		return
	}

	rangeType := c.DefaultQuery("range", "today")
	topN, _ := strconv.Atoi(c.DefaultQuery("topN", "5"))
	if topN <= 0 {
		topN = 5
	}

	startTime, endTime := getTimeRange(rangeType)

	topResult, err := ctrl.service.GetAlertTopByHostGroup(uint(accountID), startTime, endTime, topN, rangeType)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取告警TOP失败: "+err.Error())
		return
	}

	result.Success(c, topResult)
}

// GetZabbixNotificationStats 获取Zabbix每日告警数量统计
// @Summary 获取Zabbix每日告警数量统计(按天)
// @Description 获取指定时间范围内的Zabbix告警数量统计,按天返回数据用于绘制折线图。使用problem.get接口获取告警问题数量
// @Tags Zabbix监控统计
// @Accept json
// @Produce json
// @Param accountId query uint true "Zabbix账号ID"
// @Param range query string true "时间范围(today/yesterday/last7days/last30days)"
// @Success 200 {object} result.Result{data=model.ZabbixNotificationStats}
// @Router /api/v1/monitor/zabbix/notification-stats [get]
// @Security ApiKeyAuth
func (ctrl *ZabbixStatsController) GetZabbixNotificationStats(c *gin.Context) {
	accountIDStr := c.Query("accountId")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的账号ID")
		return
	}

	rangeType := c.DefaultQuery("range", "today")
	startTime, endTime := getTimeRange(rangeType)

	stats, err := ctrl.service.GetNotificationStats(uint(accountID), startTime, endTime, rangeType)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取告警统计失败: "+err.Error())
		return
	}

	result.Success(c, stats)
}

// getTimeRange 获取时间范围
func getTimeRange(rangeType string) (startTime, endTime int64) {
	now := time.Now()
	endTime = now.Unix()

	switch rangeType {
	case "today":
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startTime = today.Unix()
	case "yesterday":
		yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		yesterdayEnd := yesterday.Add(24 * time.Hour)
		startTime = yesterday.Unix()
		endTime = yesterdayEnd.Unix()
	case "last7days":
		startTime = now.Add(-7 * 24 * time.Hour).Unix()
	case "last30days":
		startTime = now.Add(-30 * 24 * time.Hour).Unix()
	default:
		// 默认今天
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startTime = today.Unix()
	}

	return startTime, endTime
}
