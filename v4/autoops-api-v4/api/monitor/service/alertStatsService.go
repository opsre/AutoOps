package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type AlertStatsServiceInterface interface {
	// 告警源管理
	CreateAlertSource(c *gin.Context, dto *model.CreateAlertSourceDto)
	UpdateAlertSource(c *gin.Context, dto *model.UpdateAlertSourceDto)
	DeleteAlertSource(c *gin.Context, id uint)
	GetAlertSourceById(c *gin.Context, id uint)
	GetAlertSourceList(c *gin.Context, query *model.AlertSourceQueryDto)
	BatchDeleteAlertSources(c *gin.Context, ids []uint)
	UpdateAlertSourceStatus(c *gin.Context, id uint, status int)

	// 告警统计
	GetAlertStatsToday(c *gin.Context, sourceID uint, label string)
	GetAlertStatsYesterday(c *gin.Context, sourceID uint, label string)
	GetAlertStatsLast7Days(c *gin.Context, sourceID uint, label string)
	GetAlertStatsLast30Days(c *gin.Context, sourceID uint, label string)
	GetAlertStatsTodayVsYesterday(c *gin.Context, sourceID uint, label string)
	GetAlertStatsThisWeekVsLastWeek(c *gin.Context, sourceID uint, label string)

	// 告警通知统计
	GetAlertNotificationsToday(c *gin.Context, sourceID uint)
	GetAlertNotificationsYesterday(c *gin.Context, sourceID uint)
	GetAlertNotificationsLast7Days(c *gin.Context, sourceID uint)
	GetAlertNotificationsLast30Days(c *gin.Context, sourceID uint)

	// 协作空间统计
	GetChannelStatsToday(c *gin.Context, sourceID uint)
	GetChannelStatsYesterday(c *gin.Context, sourceID uint)
	GetChannelStatsLast7Days(c *gin.Context, sourceID uint)
	GetChannelStatsLast30Days(c *gin.Context, sourceID uint)
}

type AlertStatsServiceImpl struct{}

const (
	// Redis缓存key前缀
	cacheKeyPrefixAlertStats        = "alert:stats:"
	cacheKeyPrefixAlertNotification = "alert:notification:"
	cacheKeyPrefixChannelStats      = "alert:channel:"

	// 缓存过期时间：10分钟
	cacheTTL = 10 * time.Minute
)

// generateCacheKey 生成缓存key
func generateCacheKey(prefix string, params ...interface{}) string {
	data := fmt.Sprintf("%v", params)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%s%x", prefix, hash)
}

// getFromCache 从缓存获取数据
func (s *AlertStatsServiceImpl) getFromCache(key string, result interface{}) (bool, error) {
	redisClient := common.GetRedisClient()
	if redisClient == nil {
		return false, nil
	}

	ctx := context.Background()
	data, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		// 缓存未命中或其他错误
		return false, nil
	}

	if err := json.Unmarshal([]byte(data), result); err != nil {
		return false, err
	}

	return true, nil
}

// setToCache 设置数据到缓存
func (s *AlertStatsServiceImpl) setToCache(key string, data interface{}) error {
	redisClient := common.GetRedisClient()
	if redisClient == nil {
		return nil
	}

	ctx := context.Background()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, jsonData, cacheTTL).Err()
}

// CreateAlertSource 创建告警源
func (s *AlertStatsServiceImpl) CreateAlertSource(c *gin.Context, dto *model.CreateAlertSourceDto) {
	alertStatsDao := dao.NewAlertStatsDao()

	// 检查名称是否已存在
	exists, err := alertStatsDao.CheckAlertSourceExists(dto.Name, 0)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "检查告警源名称失败")
		return
	}
	if exists {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源名称已存在")
		return
	}

	// 设置默认值
	if dto.Status == 0 {
		dto.Status = 1 // 默认启用
	}
	if dto.ApiBaseURL == "" {
		// 根据类型设置默认API地址
		switch dto.Type {
		case model.AlertSourceTypeFlashDuty:
			dto.ApiBaseURL = "https://api.flashcat.cloud"
		}
	}

	source := &model.AlertSource{
		Name:       dto.Name,
		Type:       dto.Type,
		KeyID:      dto.KeyID,
		HostID:     dto.HostID,
		AppKey:     dto.AppKey,
		ApiBaseURL: dto.ApiBaseURL,
		Status:     dto.Status,
		Remark:     dto.Remark,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	err = alertStatsDao.CreateAlertSource(source)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建告警源失败")
		return
	}

	result.Success(c, source)
}

// UpdateAlertSource 更新告警源
func (s *AlertStatsServiceImpl) UpdateAlertSource(c *gin.Context, dto *model.UpdateAlertSourceDto) {
	alertStatsDao := dao.NewAlertStatsDao()

	// 检查记录是否存在
	_, err := alertStatsDao.GetAlertSourceById(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	// 如果修改了名称，检查名称是否已被其他记录使用
	if dto.Name != "" {
		exists, err := alertStatsDao.CheckAlertSourceExists(dto.Name, dto.ID)
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), "检查告警源名称失败")
			return
		}
		if exists {
			result.Failed(c, int(result.ApiCode.FAILED), "告警源名称已存在")
			return
		}
	}

	source := &model.AlertSource{
		ID:         dto.ID,
		Name:       dto.Name,
		Type:       dto.Type,
		KeyID:      dto.KeyID,
		HostID:     dto.HostID,
		AppKey:     dto.AppKey,
		ApiBaseURL: dto.ApiBaseURL,
		Status:     dto.Status,
		Remark:     dto.Remark,
		UpdateTime: time.Now().Unix(),
	}

	err = alertStatsDao.UpdateAlertSource(dto.ID, source)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新告警源失败")
		return
	}

	result.Success(c, true)
}

// DeleteAlertSource 删除告警源
func (s *AlertStatsServiceImpl) DeleteAlertSource(c *gin.Context, id uint) {
	alertStatsDao := dao.NewAlertStatsDao()

	err := alertStatsDao.DeleteAlertSource(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除告警源失败")
		return
	}

	result.Success(c, true)
}

// GetAlertSourceById 根据ID获取告警源
func (s *AlertStatsServiceImpl) GetAlertSourceById(c *gin.Context, id uint) {
	alertStatsDao := dao.NewAlertStatsDao()

	source, err := alertStatsDao.GetAlertSourceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	result.Success(c, source)
}

// GetAlertSourceList 获取告警源列表
func (s *AlertStatsServiceImpl) GetAlertSourceList(c *gin.Context, query *model.AlertSourceQueryDto) {
	alertStatsDao := dao.NewAlertStatsDao()

	// 设置默认分页参数
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	list, total, err := alertStatsDao.GetAlertSourceList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询告警源列表失败")
		return
	}

	result.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	})
}

// BatchDeleteAlertSources 批量删除告警源
func (s *AlertStatsServiceImpl) BatchDeleteAlertSources(c *gin.Context, ids []uint) {
	if len(ids) == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "请选择要删除的记录")
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()
	err := alertStatsDao.BatchDeleteAlertSources(ids)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "批量删除失败")
		return
	}

	result.Success(c, true)
}

// UpdateAlertSourceStatus 更新告警源状态
func (s *AlertStatsServiceImpl) UpdateAlertSourceStatus(c *gin.Context, id uint, status int) {
	if status < 0 || status > 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "无效的状态值")
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()
	err := alertStatsDao.UpdateAlertSourceStatus(id, status)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新状态失败")
		return
	}

	result.Success(c, true)
}

// GetAlertStats 获取告警统计数据（基础版）
func (s *AlertStatsServiceImpl) GetAlertStats(c *gin.Context, query *model.AlertStatsQueryDto) {
	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixAlertStats, query.SourceID, query.StartTime, query.EndTime, query.Label)

	// 尝试从缓存获取
	var cachedResult model.AlertStatsResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(query.SourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	// 设置默认值
	if query.TimeZone == "" {
		query.TimeZone = "Asia/Shanghai"
	}
	if query.K == 0 {
		query.K = 100
	}
	if query.OrderBy == "" {
		query.OrderBy = "total_alert_event_cnt"
	}

	// 根据告警源类型调用不同的API
	var statsResult *model.AlertStatsResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		statsResult, err = s.fetchFlashDutyStats(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取告警统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, statsResult)

	result.Success(c, statsResult)
}

// GetAlertStatsWithCompare 获取告警统计数据（含同环比）
func (s *AlertStatsServiceImpl) GetAlertStatsWithCompare(c *gin.Context, query *model.AlertStatsQueryDto) {
	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixAlertStats, "compare", query.SourceID, query.StartTime, query.EndTime, query.CompareStartTime, query.CompareEndTime, query.Label)

	// 尝试从缓存获取
	var cachedResult model.AlertStatsCompareResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(query.SourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	// 必须提供对比时间
	if query.CompareStartTime == 0 || query.CompareEndTime == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "请提供对比时间范围")
		return
	}

	// 设置默认值
	if query.TimeZone == "" {
		query.TimeZone = "Asia/Shanghai"
	}
	if query.K == 0 {
		query.K = 100
	}
	if query.OrderBy == "" {
		query.OrderBy = "total_alert_event_cnt"
	}

	// 根据告警源类型调用不同的API
	var compareResult *model.AlertStatsCompareResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		compareResult, err = s.fetchFlashDutyStatsWithCompare(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取告警统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, compareResult)

	result.Success(c, compareResult)
}

// fetchFlashDutyStats 调用FlashDuty API获取统计数据
func (s *AlertStatsServiceImpl) fetchFlashDutyStats(source *model.AlertSource, query *model.AlertStatsQueryDto) (*model.AlertStatsResult, error) {
	// 构建请求
	req := &model.FlashDutyTopkRequest{
		StartTime:          query.StartTime,
		EndTime:            query.EndTime,
		TeamIDs:            query.TeamIDs,
		IsMyTeam:           query.IsMyTeam,
		ChannelIDs:         query.ChannelIDs,
		Severities:         query.Severities,
		TimeZone:           query.TimeZone,
		Query:              query.Query,
		Labels:             query.Labels,
		Fields:             query.Fields,
		SecondsToCloseFrom: query.SecondsToCloseFrom,
		SecondsToCloseTo:   query.SecondsToCloseTo,
		SecondsToAckFrom:   query.SecondsToAckFrom,
		SecondsToAckTo:     query.SecondsToAckTo,
		OrderBy:            query.OrderBy,
		Asc:                query.Asc,
		K:                  query.K,
		Label:              query.Label,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 发送HTTP请求
	url := fmt.Sprintf("%s/insight/alert/topk-by-label?app_key=%s", source.ApiBaseURL, source.AppKey)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "*/*")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用FlashDuty API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FlashDuty API返回错误: %s", string(body))
	}

	// 解析响应
	var flashResp model.FlashDutyTopkResponse
	if err := json.Unmarshal(body, &flashResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 转换为统计结果
	statsResult := &model.AlertStatsResult{
		Items:      make([]model.AlertStatsItem, 0),
		StartTime:  query.StartTime,
		EndTime:    query.EndTime,
		QueryLabel: query.Label,
	}

	total := 0
	for _, item := range flashResp.Data.Items {
		total += item.TotalAlertEventCnt
	}

	for _, item := range flashResp.Data.Items {
		percentage := 0.0
		if total > 0 {
			percentage = float64(item.TotalAlertEventCnt) / float64(total) * 100
		}

		statsResult.Items = append(statsResult.Items, model.AlertStatsItem{
			Label:              item.Label,
			TotalAlertCnt:      item.TotalAlertCnt,
			TotalAlertEventCnt: item.TotalAlertEventCnt,
			Percentage:         percentage,
		})
	}

	statsResult.Total = total

	return statsResult, nil
}

// fetchFlashDutyStatsWithCompare 调用FlashDuty API获取对比统计数据
func (s *AlertStatsServiceImpl) fetchFlashDutyStatsWithCompare(source *model.AlertSource, query *model.AlertStatsQueryDto) (*model.AlertStatsCompareResult, error) {
	// 获取当前周期数据
	currentQuery := *query
	currentStats, err := s.fetchFlashDutyStats(source, &currentQuery)
	if err != nil {
		return nil, fmt.Errorf("获取当前周期数据失败: %v", err)
	}

	// 获取对比周期数据
	compareQuery := *query
	compareQuery.StartTime = query.CompareStartTime
	compareQuery.EndTime = query.CompareEndTime
	compareStats, err := s.fetchFlashDutyStats(source, &compareQuery)
	if err != nil {
		return nil, fmt.Errorf("获取对比周期数据失败: %v", err)
	}

	// 构建对比结果
	compareResult := &model.AlertStatsCompareResult{
		Items:            make([]model.AlertStatsCompareItem, 0),
		Total:            currentStats.Total,
		CompareTotal:     compareStats.Total,
		StartTime:        query.StartTime,
		EndTime:          query.EndTime,
		CompareStartTime: query.CompareStartTime,
		CompareEndTime:   query.CompareEndTime,
		QueryLabel:       query.Label,
	}

	// 计算总数变化
	compareResult.TotalChange = compareResult.Total - compareResult.CompareTotal
	if compareResult.CompareTotal > 0 {
		compareResult.TotalChangeRate = float64(compareResult.TotalChange) / float64(compareResult.CompareTotal) * 100
	}

	// 构建标签映射
	compareMap := make(map[string]model.AlertStatsItem)
	for _, item := range compareStats.Items {
		compareMap[item.Label] = item
	}

	// 合并数据
	for _, currentItem := range currentStats.Items {
		compareItem := model.AlertStatsCompareItem{
			Label:              currentItem.Label,
			TotalAlertCnt:      currentItem.TotalAlertCnt,
			TotalAlertEventCnt: currentItem.TotalAlertEventCnt,
			Percentage:         currentItem.Percentage,
		}

		// 查找对比数据
		if prevItem, exists := compareMap[currentItem.Label]; exists {
			compareItem.CompareTotalAlertCnt = prevItem.TotalAlertCnt
			compareItem.CompareTotalAlertEventCnt = prevItem.TotalAlertEventCnt
			compareItem.ComparePercentage = prevItem.Percentage

			// 计算变化量
			compareItem.AlertCntChange = compareItem.TotalAlertCnt - compareItem.CompareTotalAlertCnt
			compareItem.AlertEventCntChange = compareItem.TotalAlertEventCnt - compareItem.CompareTotalAlertEventCnt

			// 计算变化率
			if compareItem.CompareTotalAlertCnt > 0 {
				compareItem.AlertCntChangeRate = float64(compareItem.AlertCntChange) / float64(compareItem.CompareTotalAlertCnt) * 100
			}
			if compareItem.CompareTotalAlertEventCnt > 0 {
				compareItem.AlertEventCntChangeRate = float64(compareItem.AlertEventCntChange) / float64(compareItem.CompareTotalAlertEventCnt) * 100
			}
		}

		compareResult.Items = append(compareResult.Items, compareItem)
	}

	return compareResult, nil
}

// GetAlertStatsToday 获取今天的告警统计
func (s *AlertStatsServiceImpl) GetAlertStatsToday(c *gin.Context, sourceID uint, label string) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	query := &model.AlertStatsQueryDto{
		SourceID:  sourceID,
		StartTime: startOfDay.Unix(),
		EndTime:   endOfDay.Unix(),
		TimeZone:  "Asia/Shanghai",
		Label:     label,
		K:         5, // 改为TOP 5
	}

	s.GetAlertStats(c, query)
}

// GetAlertStatsYesterday 获取昨天的告警统计
func (s *AlertStatsServiceImpl) GetAlertStatsYesterday(c *gin.Context, sourceID uint, label string) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, loc)

	query := &model.AlertStatsQueryDto{
		SourceID:  sourceID,
		StartTime: startOfDay.Unix(),
		EndTime:   endOfDay.Unix(),
		TimeZone:  "Asia/Shanghai",
		Label:     label,
		K:         5, // 改为TOP 5
	}

	s.GetAlertStats(c, query)
}

// GetAlertStatsLast7Days 获取最近7天的告警统计
func (s *AlertStatsServiceImpl) GetAlertStatsLast7Days(c *gin.Context, sourceID uint, label string) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	startTime := endTime.AddDate(0, 0, -6).Add(-23*time.Hour - 59*time.Minute - 59*time.Second)

	query := &model.AlertStatsQueryDto{
		SourceID:  sourceID,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
		TimeZone:  "Asia/Shanghai",
		Label:     label,
		K:         5, // 改为TOP 5
	}

	s.GetAlertStats(c, query)
}

// GetAlertStatsLast30Days 获取最近30天的告警统计
func (s *AlertStatsServiceImpl) GetAlertStatsLast30Days(c *gin.Context, sourceID uint, label string) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	startTime := endTime.AddDate(0, 0, -29).Add(-23*time.Hour - 59*time.Minute - 59*time.Second)

	query := &model.AlertStatsQueryDto{
		SourceID:  sourceID,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
		TimeZone:  "Asia/Shanghai",
		Label:     label,
		K:         5, // 改为TOP 5
	}

	s.GetAlertStats(c, query)
}

// GetAlertStatsTodayVsYesterday 今天 vs 昨天同环比
func (s *AlertStatsServiceImpl) GetAlertStatsTodayVsYesterday(c *gin.Context, sourceID uint, label string) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 今天
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	// 昨天
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, loc)

	query := &model.AlertStatsQueryDto{
		SourceID:         sourceID,
		StartTime:        todayStart.Unix(),
		EndTime:          todayEnd.Unix(),
		CompareStartTime: yesterdayStart.Unix(),
		CompareEndTime:   yesterdayEnd.Unix(),
		TimeZone:         "Asia/Shanghai",
		Label:            label,
		K:                5, // 改为TOP 5
	}

	s.GetAlertStatsWithCompare(c, query)
}

// GetAlertStatsThisWeekVsLastWeek 本周 vs 上周同环比
func (s *AlertStatsServiceImpl) GetAlertStatsThisWeekVsLastWeek(c *gin.Context, sourceID uint, label string) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 本周（周一到今天）
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 将周日从0改为7
	}
	thisWeekStart := now.AddDate(0, 0, -(weekday - 1))
	thisWeekStart = time.Date(thisWeekStart.Year(), thisWeekStart.Month(), thisWeekStart.Day(), 0, 0, 0, 0, loc)
	thisWeekEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	// 上周同期（上周一到上周同一天）
	lastWeekStart := thisWeekStart.AddDate(0, 0, -7)
	lastWeekEnd := thisWeekEnd.AddDate(0, 0, -7)

	query := &model.AlertStatsQueryDto{
		SourceID:         sourceID,
		StartTime:        thisWeekStart.Unix(),
		EndTime:          thisWeekEnd.Unix(),
		CompareStartTime: lastWeekStart.Unix(),
		CompareEndTime:   lastWeekEnd.Unix(),
		TimeZone:         "Asia/Shanghai",
		Label:            label,
		K:                5, // 改为TOP 5
	}

	s.GetAlertStatsWithCompare(c, query)
}

// fetchAlertNotifications 调用FlashDuty API获取告警通知统计
func (s *AlertStatsServiceImpl) fetchAlertNotifications(source *model.AlertSource, query *model.AlertNotificationQueryDto) (*model.AlertNotificationResult, error) {
	// 构建请求
	req := &model.FlashDutyNotificationRequest{
		ChannelIDs: query.ChannelIDs,
		Step:       query.Step,
		StartTime:  query.StartTime,
		EndTime:    query.EndTime,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 发送HTTP请求
	url := fmt.Sprintf("%s/report/oncall/notifications?app_key=%s", source.ApiBaseURL, source.AppKey)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "*/*")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用FlashDuty API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FlashDuty API返回错误: %s", string(body))
	}

	// 解析响应
	var flashResp model.FlashDutyNotificationResponse
	if err := json.Unmarshal(body, &flashResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 转换为统计结果
	notificationResult := &model.AlertNotificationResult{
		Items:      make([]model.AlertNotificationItem, 0),
		Step:       query.Step,
		StartTime:  query.StartTime,
		EndTime:    query.EndTime,
	}

	totalSms := 0
	totalEmail := 0
	totalVoice := 0

	for _, item := range flashResp.Data.Items {
		totalSms += item.SmsCount
		totalEmail += item.EmailCount
		totalVoice += item.VoiceCount

		notificationResult.Items = append(notificationResult.Items, model.AlertNotificationItem{
			Timestamp:  item.Timestamp,
			SmsCount:   item.SmsCount,
			EmailCount: item.EmailCount,
			VoiceCount: item.VoiceCount,
			TotalCount: item.SmsCount + item.EmailCount + item.VoiceCount,
		})
	}

	notificationResult.TotalSms = totalSms
	notificationResult.TotalEmail = totalEmail
	notificationResult.TotalVoice = totalVoice
	notificationResult.GrandTotal = totalSms + totalEmail + totalVoice

	return notificationResult, nil
}

// GetAlertNotificationsToday 获取今天的告警通知统计
func (s *AlertStatsServiceImpl) GetAlertNotificationsToday(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixAlertNotification, "today", sourceID, startOfDay.Unix(), endOfDay.Unix())

	// 尝试从缓存获取
	var cachedResult model.AlertNotificationResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.AlertNotificationQueryDto{
		SourceID:   sourceID,
		ChannelIDs: []int64{},
		Step:       "day",
		StartTime:  startOfDay.Unix(),
		EndTime:    endOfDay.Unix(),
	}

	// 根据告警源类型调用不同的API
	var notificationResult *model.AlertNotificationResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		notificationResult, err = s.fetchAlertNotifications(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取告警通知统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, notificationResult)

	result.Success(c, notificationResult)
}

// GetAlertNotificationsYesterday 获取昨天的告警通知统计
func (s *AlertStatsServiceImpl) GetAlertNotificationsYesterday(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, loc)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixAlertNotification, "yesterday", sourceID, startOfDay.Unix(), endOfDay.Unix())

	// 尝试从缓存获取
	var cachedResult model.AlertNotificationResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.AlertNotificationQueryDto{
		SourceID:   sourceID,
		ChannelIDs: []int64{},
		Step:       "day",
		StartTime:  startOfDay.Unix(),
		EndTime:    endOfDay.Unix(),
	}

	// 根据告警源类型调用不同的API
	var notificationResult *model.AlertNotificationResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		notificationResult, err = s.fetchAlertNotifications(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取告警通知统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, notificationResult)

	result.Success(c, notificationResult)
}

// GetAlertNotificationsLast7Days 获取最近7天的告警通知统计
func (s *AlertStatsServiceImpl) GetAlertNotificationsLast7Days(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	startTime := endTime.AddDate(0, 0, -6).Add(-23*time.Hour - 59*time.Minute - 59*time.Second)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixAlertNotification, "last7days", sourceID, startTime.Unix(), endTime.Unix())

	// 尝试从缓存获取
	var cachedResult model.AlertNotificationResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.AlertNotificationQueryDto{
		SourceID:   sourceID,
		ChannelIDs: []int64{},
		Step:       "day",
		StartTime:  startTime.Unix(),
		EndTime:    endTime.Unix(),
	}

	// 根据告警源类型调用不同的API
	var notificationResult *model.AlertNotificationResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		notificationResult, err = s.fetchAlertNotifications(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取告警通知统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, notificationResult)

	result.Success(c, notificationResult)
}

// GetAlertNotificationsLast30Days 获取最近30天的告警通知统计
func (s *AlertStatsServiceImpl) GetAlertNotificationsLast30Days(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	startTime := endTime.AddDate(0, 0, -29).Add(-23*time.Hour - 59*time.Minute - 59*time.Second)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixAlertNotification, "last30days", sourceID, startTime.Unix(), endTime.Unix())

	// 尝试从缓存获取
	var cachedResult model.AlertNotificationResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.AlertNotificationQueryDto{
		SourceID:   sourceID,
		ChannelIDs: []int64{},
		Step:       "day",
		StartTime:  startTime.Unix(),
		EndTime:    endTime.Unix(),
	}

	// 根据告警源类型调用不同的API
	var notificationResult *model.AlertNotificationResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		notificationResult, err = s.fetchAlertNotifications(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取告警通知统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, notificationResult)

	result.Success(c, notificationResult)
}

// fetchChannelStats 调用FlashDuty API获取协作空间统计
func (s *AlertStatsServiceImpl) fetchChannelStats(source *model.AlertSource, query *model.ChannelStatsQueryDto) (*model.ChannelStatsResult, error) {
	// 构建请求
	req := &model.FlashDutyChannelRequest{
		StartTime:          query.StartTime,
		EndTime:            query.EndTime,
		TeamIDs:            query.TeamIDs,
		IsMyTeam:           query.IsMyTeam,
		ChannelIDs:         query.ChannelIDs,
		Severities:         query.Severities,
		TimeZone:           query.TimeZone,
		Query:              query.Query,
		Labels:             query.Labels,
		Fields:             query.Fields,
		SecondsToCloseFrom: query.SecondsToCloseFrom,
		SecondsToCloseTo:   query.SecondsToCloseTo,
		SecondsToAckFrom:   query.SecondsToAckFrom,
		SecondsToAckTo:     query.SecondsToAckTo,
		AggregateUnit:      query.AggregateUnit,
		SplitHours:         query.SplitHours,
		OrderBy:            query.OrderBy,
		Asc:                query.Asc,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 发送HTTP请求
	url := fmt.Sprintf("%s/insight/channel?app_key=%s", source.ApiBaseURL, source.AppKey)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "*/*")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用FlashDuty API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FlashDuty API返回错误: %s", string(body))
	}

	// 解析响应
	var flashResp model.FlashDutyChannelResponse
	if err := json.Unmarshal(body, &flashResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 转换为统计结果
	channelStatsResult := &model.ChannelStatsResult{
		Items:         make([]model.ChannelStatsItem, 0),
		StartTime:     query.StartTime,
		EndTime:       query.EndTime,
		AggregateUnit: query.AggregateUnit,
	}

	// 先收集所有数据
	allItems := make([]model.ChannelStatsItem, 0)
	for _, item := range flashResp.Data.Items {
		// 计算自动关闭率
		autoClosePct := 0.0
		if item.TotalIncidentsClosed > 0 {
			autoClosePct = float64(item.TotalIncidentsAutoClosed) / float64(item.TotalIncidentsClosed) * 100
		}

		allItems = append(allItems, model.ChannelStatsItem{
			Timestamp:                       item.Timestamp,
			ChannelID:                       item.ChannelID,
			ChannelName:                     item.ChannelName,
			TotalIncidentCnt:                item.TotalIncidentCnt,
			TotalIncidentsAcknowledged:      item.TotalIncidentsAcknowledged,
			TotalIncidentsClosed:            item.TotalIncidentsClosed,
			TotalIncidentsAutoClosed:        item.TotalIncidentsAutoClosed,
			TotalIncidentsManuallyClosed:    item.TotalIncidentsManuallyClosed,
			TotalIncidentsTimeoutClosed:     item.TotalIncidentsTimeoutClosed,
			TotalIncidentsEscalated:         item.TotalIncidentsEscalated,
			TotalIncidentsManuallyEscalated: item.TotalIncidentsManuallyEscalated,
			TotalIncidentsTimeoutEscalated:  item.TotalIncidentsTimeoutEscalated,
			TotalIncidentsReassigned:        item.TotalIncidentsReassigned,
			TotalInterruptions:              item.TotalInterruptions,
			TotalNotifications:              item.TotalNotifications,
			TotalEngagedSeconds:             item.TotalEngagedSeconds,
			TotalSecondsToAck:               item.TotalSecondsToAck,
			TotalSecondsToClose:             item.TotalSecondsToClose,
			MeanSecondsToAck:                item.MeanSecondsToAck,
			MeanSecondsToClose:              item.MeanSecondsToClose,
			NoiseReductionPct:               item.NoiseReductionPct,
			AcknowledgementPct:              item.AcknowledgementPct,
			AutoClosePct:                    autoClosePct,
			TotalAlertCnt:                   item.TotalAlertCnt,
			TotalAlertEventCnt:              item.TotalAlertEventCnt,
		})
	}

	// 按 totalIncidentCnt 降序排序
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].TotalIncidentCnt > allItems[j].TotalIncidentCnt
	})

	// 只取前6条
	limit := 6
	if len(allItems) < limit {
		limit = len(allItems)
	}
	channelStatsResult.Items = allItems[:limit]

	channelStatsResult.Total = len(channelStatsResult.Items)

	return channelStatsResult, nil
}

// GetChannelStatsToday 获取今天的协作空间统计
func (s *AlertStatsServiceImpl) GetChannelStatsToday(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixChannelStats, "today", sourceID, startOfDay.Unix(), endOfDay.Unix())

	// 尝试从缓存获取
	var cachedResult model.ChannelStatsResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.ChannelStatsQueryDto{
		SourceID:   sourceID,
		StartTime:  startOfDay.Unix(),
		EndTime:    endOfDay.Unix(),
		TimeZone:   "Asia/Shanghai",
		OrderBy:    "total_incident_cnt",
		Asc:        false,
	}

	// 根据告警源类型调用不同的API
	var channelStatsResult *model.ChannelStatsResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		channelStatsResult, err = s.fetchChannelStats(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取协作空间统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, channelStatsResult)

	result.Success(c, channelStatsResult)
}

// GetChannelStatsYesterday 获取昨天的协作空间统计
func (s *AlertStatsServiceImpl) GetChannelStatsYesterday(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, loc)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixChannelStats, "yesterday", sourceID, startOfDay.Unix(), endOfDay.Unix())

	// 尝试从缓存获取
	var cachedResult model.ChannelStatsResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.ChannelStatsQueryDto{
		SourceID:   sourceID,
		StartTime:  startOfDay.Unix(),
		EndTime:    endOfDay.Unix(),
		TimeZone:   "Asia/Shanghai",
		OrderBy:    "total_incident_cnt",
		Asc:        false,
	}

	// 根据告警源类型调用不同的API
	var channelStatsResult *model.ChannelStatsResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		channelStatsResult, err = s.fetchChannelStats(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取协作空间统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, channelStatsResult)

	result.Success(c, channelStatsResult)
}

// GetChannelStatsLast7Days 获取最近7天的协作空间统计
func (s *AlertStatsServiceImpl) GetChannelStatsLast7Days(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	startTime := endTime.AddDate(0, 0, -6).Add(-23*time.Hour - 59*time.Minute - 59*time.Second)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixChannelStats, "last7days", sourceID, startTime.Unix(), endTime.Unix())

	// 尝试从缓存获取
	var cachedResult model.ChannelStatsResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.ChannelStatsQueryDto{
		SourceID:   sourceID,
		StartTime:  startTime.Unix(),
		EndTime:    endTime.Unix(),
		TimeZone:   "Asia/Shanghai",
		OrderBy:    "total_incident_cnt",
		Asc:        false,
	}

	// 根据告警源类型调用不同的API
	var channelStatsResult *model.ChannelStatsResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		channelStatsResult, err = s.fetchChannelStats(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取协作空间统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, channelStatsResult)

	result.Success(c, channelStatsResult)
}

// GetChannelStatsLast30Days 获取最近30天的协作空间统计
func (s *AlertStatsServiceImpl) GetChannelStatsLast30Days(c *gin.Context, sourceID uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	startTime := endTime.AddDate(0, 0, -29).Add(-23*time.Hour - 59*time.Minute - 59*time.Second)

	// 生成缓存key
	cacheKey := generateCacheKey(cacheKeyPrefixChannelStats, "last30days", sourceID, startTime.Unix(), endTime.Unix())

	// 尝试从缓存获取
	var cachedResult model.ChannelStatsResult
	if found, _ := s.getFromCache(cacheKey, &cachedResult); found {
		result.Success(c, &cachedResult)
		return
	}

	alertStatsDao := dao.NewAlertStatsDao()

	// 获取告警源配置
	source, err := alertStatsDao.GetAlertSourceById(sourceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源不存在")
		return
	}

	if source.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "告警源已禁用")
		return
	}

	query := &model.ChannelStatsQueryDto{
		SourceID:   sourceID,
		StartTime:  startTime.Unix(),
		EndTime:    endTime.Unix(),
		TimeZone:   "Asia/Shanghai",
		OrderBy:    "total_incident_cnt",
		Asc:        false,
	}

	// 根据告警源类型调用不同的API
	var channelStatsResult *model.ChannelStatsResult
	switch source.Type {
	case model.AlertSourceTypeFlashDuty:
		channelStatsResult, err = s.fetchChannelStats(&source, query)
	case model.AlertSourceTypeZabbix:
		result.Failed(c, int(result.ApiCode.FAILED), "Zabbix告警源暂未实现")
		return
	case model.AlertSourceTypeN9e:
		result.Failed(c, int(result.ApiCode.FAILED), "N9e告警源暂未实现")
		return
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的告警源类型")
		return
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取协作空间统计失败: %v", err))
		return
	}

	// 设置缓存
	s.setToCache(cacheKey, channelStatsResult)

	result.Success(c, channelStatsResult)
}

// GetAlertStatsService 获取告警统计服务实例
func GetAlertStatsService() AlertStatsServiceInterface {
	return &AlertStatsServiceImpl{}
}
