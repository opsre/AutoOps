package service

import (
	"dodevops-api/api/db/mysql/dao"
	"dodevops-api/api/db/mysql/model"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

// IDbSqlExecService SQL执行历史服务接口
type IDbSqlExecService interface {
	GetDbSqlExecList(c *gin.Context, query *model.DbSqlExecQuery)
	GetDbSqlExecById(c *gin.Context, id uint)
	GetDbSqlExecStatistics(c *gin.Context, dbId uint, startTime, endTime string)
	DeleteOldRecords(c *gin.Context, days int)
}

type DbSqlExecServiceImpl struct{}

// GetDbSqlExecList 获取SQL执行记录列表
func (s DbSqlExecServiceImpl) GetDbSqlExecList(c *gin.Context, query *model.DbSqlExecQuery) {
	execs, total, err := dao.GetDbSqlExecList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	data := map[string]interface{}{
		"list":     execs,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}
	result.Success(c, data)
}

// GetDbSqlExecById 根据ID获取SQL执行记录
func (s DbSqlExecServiceImpl) GetDbSqlExecById(c *gin.Context, id uint) {
	exec, err := dao.GetDbSqlExecById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, exec)
}

// GetDbSqlExecStatistics 获取SQL执行统计
func (s DbSqlExecServiceImpl) GetDbSqlExecStatistics(c *gin.Context, dbId uint, startTime, endTime string) {
	stats, err := dao.GetDbSqlExecStatistics(dbId, startTime, endTime)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, stats)
}

// DeleteOldRecords 删除旧的执行记录
func (s DbSqlExecServiceImpl) DeleteOldRecords(c *gin.Context, days int) {
	if days <= 0 {
		days = 30 // 默认30天
	}

	if err := dao.DeleteOldDbSqlExecs(days); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}

// DbSqlExecService 服务实例
func DbSqlExecService() IDbSqlExecService {
	return DbSqlExecServiceImpl{}
}
