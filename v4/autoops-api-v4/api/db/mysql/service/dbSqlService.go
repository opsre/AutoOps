package service

import (
	"dodevops-api/api/db/mysql/dao"
	"dodevops-api/api/db/mysql/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"github.com/gin-gonic/gin"
	"time"
)

// IDbSqlService SQL模板服务接口
type IDbSqlService interface {
	GetDbSqlList(c *gin.Context, query *model.DbSqlQuery)
	GetDbSqlById(c *gin.Context, id uint)
	GetDbSqlNames(c *gin.Context, dbId uint, db string)
	CreateDbSql(c *gin.Context, dto *model.DbSqlDto)
	UpdateDbSql(c *gin.Context, dto *model.DbSqlDto)
	DeleteDbSql(c *gin.Context, id uint)
}

type DbSqlServiceImpl struct{}

// GetDbSqlList 获取SQL模板列表
func (s DbSqlServiceImpl) GetDbSqlList(c *gin.Context, query *model.DbSqlQuery) {
	sqls, total, err := dao.GetDbSqlList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	data := map[string]interface{}{
		"list":     sqls,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}
	result.Success(c, data)
}

// GetDbSqlById 根据ID获取SQL模板
func (s DbSqlServiceImpl) GetDbSqlById(c *gin.Context, id uint) {
	sql, err := dao.GetDbSqlById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, sql)
}

// GetDbSqlNames 获取SQL模板名称列表
func (s DbSqlServiceImpl) GetDbSqlNames(c *gin.Context, dbId uint, db string) {
	names, err := dao.GetDbSqlNames(dbId, db)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, names)
}

// CreateDbSql 创建SQL模板
func (s DbSqlServiceImpl) CreateDbSql(c *gin.Context, dto *model.DbSqlDto) {
	// 检查名称是否已存在
	if dao.CheckDbSqlExists(dto.DbId, dto.Name, 0) {
		result.Failed(c, int(result.ApiCode.FAILED), "SQL模板名称已存在")
		return
	}

	now := time.Now()
	sql := &model.DbSql{
		DbId:       dto.DbId,
		Db:         dto.Db,
		Name:       dto.Name,
		Type:       dto.Type,
		Sql:        dto.Sql,
		Remark:     dto.Remark,
		CreateTime: util.HTime{Time: now},
		UpdateTime: util.HTime{Time: now},
		Creator:    c.GetString("username"),
		CreatorId:  uint(c.GetInt("userId")),
	}

	if err := dao.CreateDbSql(sql); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, sql.ID)
}

// UpdateDbSql 更新SQL模板
func (s DbSqlServiceImpl) UpdateDbSql(c *gin.Context, dto *model.DbSqlDto) {
	// 检查SQL模板是否存在
	sql, err := dao.GetDbSqlById(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "SQL模板不存在")
		return
	}

	// 检查名称是否重复
	if dto.Name != sql.Name && dao.CheckDbSqlExists(dto.DbId, dto.Name, dto.ID) {
		result.Failed(c, int(result.ApiCode.FAILED), "SQL模板名称已存在")
		return
	}

	// 更新字段
	sql.Name = dto.Name
	sql.Type = dto.Type
	sql.Sql = dto.Sql
	sql.Remark = dto.Remark
	sql.UpdateTime = util.HTime{Time: time.Now()}
	sql.Modifier = c.GetString("username")
	sql.ModifierId = uint(c.GetInt("userId"))

	if err := dao.UpdateDbSql(&sql); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}

// DeleteDbSql 删除SQL模板
func (s DbSqlServiceImpl) DeleteDbSql(c *gin.Context, id uint) {
	if err := dao.DeleteDbSql(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}

// DbSqlService 服务实例
func DbSqlService() IDbSqlService {
	return DbSqlServiceImpl{}
}
