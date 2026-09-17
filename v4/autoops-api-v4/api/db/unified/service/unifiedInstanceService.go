package service

import (
	"dodevops-api/api/db/unified/dao"
	"dodevops-api/api/db/unified/model"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

// IUnifiedInstanceService 统一实例服务接口
type IUnifiedInstanceService interface {
	GetUnifiedInstanceList(c *gin.Context, query *model.UnifiedInstanceQuery)
	GetUnifiedInstanceById(c *gin.Context, dbType string, id uint)
	GetUnifiedInstanceStats(c *gin.Context)
	GetInstancePassword(c *gin.Context, dbType string, id uint)
}

type UnifiedInstanceServiceImpl struct{}

func UnifiedInstanceService() IUnifiedInstanceService {
	return UnifiedInstanceServiceImpl{}
}

// GetUnifiedInstanceList 获取统一实例列表
func (UnifiedInstanceServiceImpl) GetUnifiedInstanceList(c *gin.Context, query *model.UnifiedInstanceQuery) {
	// 设置默认分页参数
	if query.PageNum == 0 {
		query.PageNum = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	list, total, err := dao.GetUnifiedInstanceList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 转换为 VO
	voList := make([]model.UnifiedInstanceVO, len(list))
	for i, inst := range list {
		voList[i] = inst.ToVO()
	}

	result.Success(c, map[string]interface{}{
		"list":     voList,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	})
}

// GetUnifiedInstanceById 根据类型和ID获取实例详情
func (UnifiedInstanceServiceImpl) GetUnifiedInstanceById(c *gin.Context, dbType string, id uint) {
	instance, err := dao.GetUnifiedInstanceById(dbType, id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	vo := instance.ToVO()
	result.Success(c, vo)
}

// GetUnifiedInstanceStats 获取统计信息
func (UnifiedInstanceServiceImpl) GetUnifiedInstanceStats(c *gin.Context) {
	stats, err := dao.GetUnifiedInstanceStats()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 计算总数
	var total int64 = 0
	for _, count := range stats {
		total += count
	}

	result.Success(c, map[string]interface{}{
		"stats": stats,
		"total": total,
	})
}

// GetInstancePassword 获取数据库实例的真实密码
func (UnifiedInstanceServiceImpl) GetInstancePassword(c *gin.Context, dbType string, id uint) {
	password, err := dao.GetInstancePassword(dbType, id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"password": password,
		"dbType":   dbType,
		"id":       id,
	})
}
