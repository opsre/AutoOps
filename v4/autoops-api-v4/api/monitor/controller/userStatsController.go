package controller

import (
	adminDao "dodevops-api/api/system/dao"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type UserStatsController struct{}

func NewUserStatsController() *UserStatsController {
	return &UserStatsController{}
}

// GetUserLoginStats 获取用户登录统计
// @Summary 获取用户登录统计
// @Description 获取今日用户登录次数、活跃用户数、系统总用户数
// @Tags 用户统计
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=UserLoginStatsVo}
// @Router /api/v1/monitor/user-stats/login [get]
// @Security ApiKeyAuth
func (ctrl *UserStatsController) GetUserLoginStats(c *gin.Context) {
	// 获取今日登录次数
	todayLoginCount := adminDao.GetTodayLoginCount()

	// 获取今日活跃用户数
	todayActiveUserCount := adminDao.GetTodayActiveUserCount()

	// 获取系统用户总数
	totalUserCount := adminDao.GetTotalUserCount()

	vo := UserLoginStatsVo{
		TodayLoginCount:      todayLoginCount,
		TodayActiveUserCount: todayActiveUserCount,
		TotalUserCount:       totalUserCount,
	}

	result.Success(c, vo)
}

// UserLoginStatsVo 用户登录统计视图对象
type UserLoginStatsVo struct {
	TodayLoginCount      int64 `json:"todayLoginCount"`      // 今日登录总次数
	TodayActiveUserCount int64 `json:"todayActiveUserCount"` // 今日活跃用户数
	TotalUserCount       int64 `json:"totalUserCount"`       // 系统用户总数
}
