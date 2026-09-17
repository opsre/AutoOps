package service

import (
	"dodevops-api/api/app/model"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
	"time"
)

// GetDeploymentTrend 获取上线趋势统计
func GetDeploymentTrend(c *gin.Context) {
	// 1. 绑定请求参数
	var req model.DeploymentTrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "参数错误: "+err.Error())
		return
	}

	// 2. 计算时间范围
	now := time.Now()
	var startDate, endDate time.Time
	var dateFormat string

	switch req.Period {
	case "week":
		// 近7天
		startDate = now.AddDate(0, 0, -6).Truncate(24 * time.Hour)
		endDate = now.Truncate(24 * time.Hour).Add(24*time.Hour - time.Second)
		dateFormat = "2006-01-02"
	case "month":
		// 近30天
		startDate = now.AddDate(0, 0, -29).Truncate(24 * time.Hour)
		endDate = now.Truncate(24 * time.Hour).Add(24*time.Hour - time.Second)
		dateFormat = "2006-01-02"
	case "year":
		// 近12个月
		startDate = now.AddDate(0, -11, 0)
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		endDate = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
		dateFormat = "2006-01"
	default:
		result.Failed(c, 400, "无效的时间周期")
		return
	}

	// 3. 生成日期列表
	dateList := generateDateList(startDate, endDate, req.Period)

	// 4. 查询服务上线数据
	serviceData := make(map[string]int64)
	if req.Period == "year" {
		// 按月统计
		var serviceResults []struct {
			Month string
			Count int64
		}
		err := common.GetDB().Table("app_service_release").
			Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count").
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("deleted_at IS NULL").
			Where("status = 2"). // 只统计发布成功的
			Group("month").
			Scan(&serviceResults).Error
		if err != nil {
			result.Failed(c, 500, "查询服务上线数据失败: "+err.Error())
			return
		}
		for _, r := range serviceResults {
			serviceData[r.Month] = r.Count
		}
	} else {
		// 按天统计
		var serviceResults []struct {
			Date  string
			Count int64
		}
		err := common.GetDB().Table("app_service_release").
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as date, COUNT(*) as count").
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("deleted_at IS NULL").
			Where("status = 2"). // 只统计发布成功的
			Group("date").
			Scan(&serviceResults).Error
		if err != nil {
			result.Failed(c, 500, "查询服务上线数据失败: "+err.Error())
			return
		}
		for _, r := range serviceResults {
			serviceData[r.Date] = r.Count
		}
	}

	// 5. 查询脚本上线数据
	scriptData := make(map[string]int64)
	if req.Period == "year" {
		// 按月统计
		var scriptResults []struct {
			Month string
			Count int64
		}
		err := common.GetDB().Table("app_sh_release").
			Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count").
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("deleted_at IS NULL").
			Where("status = 2"). // 只统计执行成功的
			Group("month").
			Scan(&scriptResults).Error
		if err != nil {
			result.Failed(c, 500, "查询脚本上线数据失败: "+err.Error())
			return
		}
		for _, r := range scriptResults {
			scriptData[r.Month] = r.Count
		}
	} else {
		// 按天统计
		var scriptResults []struct {
			Date  string
			Count int64
		}
		err := common.GetDB().Table("app_sh_release").
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as date, COUNT(*) as count").
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("deleted_at IS NULL").
			Where("status = 2"). // 只统计执行成功的
			Group("date").
			Scan(&scriptResults).Error
		if err != nil {
			result.Failed(c, 500, "查询脚本上线数据失败: "+err.Error())
			return
		}
		for _, r := range scriptResults {
			scriptData[r.Date] = r.Count
		}
	}

	// 6. 查询快速发布数据
	quickDeployData := make(map[string]int64)
	if req.Period == "year" {
		// 按月统计
		var quickResults []struct {
			Month string
			Count int64
		}
		err := common.GetDB().Table("quick_deployments").
			Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count").
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("status = 3"). // 只统计发布成功的
			Group("month").
			Scan(&quickResults).Error
		if err != nil {
			result.Failed(c, 500, "查询快速发布数据失败: "+err.Error())
			return
		}
		for _, r := range quickResults {
			quickDeployData[r.Month] = r.Count
		}
	} else {
		// 按天统计
		var quickResults []struct {
			Date  string
			Count int64
		}
		err := common.GetDB().Table("quick_deployments").
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as date, COUNT(*) as count").
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("status = 3"). // 只统计发布成功的
			Group("date").
			Scan(&quickResults).Error
		if err != nil {
			result.Failed(c, 500, "查询快速发布数据失败: "+err.Error())
			return
		}
		for _, r := range quickResults {
			quickDeployData[r.Date] = r.Count
		}
	}

	// 7. 组合数据
	var totalCount int64
	trendData := make([]model.DeploymentTrendData, 0, len(dateList))
	for _, date := range dateList {
		serviceCount := serviceData[date]
		scriptCount := scriptData[date]
		quickCount := quickDeployData[date]
		total := serviceCount + scriptCount + quickCount

		trendData = append(trendData, model.DeploymentTrendData{
			Date:             date,
			ServiceCount:     serviceCount,
			ScriptCount:      scriptCount,
			QuickDeployCount: quickCount,
			TotalCount:       total,
		})
		totalCount += total
	}

	// 8. 构建响应
	response := model.DeploymentTrendResponse{
		Period:     req.Period,
		StartDate:  startDate.Format(dateFormat),
		EndDate:    endDate.Format(dateFormat),
		TotalCount: totalCount,
		Data:       trendData,
	}

	result.Success(c, response)
}

// generateDateList 生成日期列表
func generateDateList(startDate, endDate time.Time, period string) []string {
	var dates []string

	if period == "year" {
		// 按月生成
		current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		end := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())

		for current.Before(end) || current.Equal(end) {
			dates = append(dates, current.Format("2006-01"))
			current = current.AddDate(0, 1, 0)
		}
	} else {
		// 按天生成
		current := startDate
		for current.Before(endDate) || current.Equal(endDate) {
			dates = append(dates, current.Format("2006-01-02"))
			current = current.AddDate(0, 0, 1)
		}
	}

	return dates
}
