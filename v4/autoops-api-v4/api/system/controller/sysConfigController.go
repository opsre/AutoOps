// 系统配置控制器
package controller

import (
	"strconv"

	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// @Tags 系统配置
// @Summary 创建系统配置
// @Description 创建系统配置
// @Accept json
// @Produce json
// @Description 创建新的系统配置
// @Param data body model.SysConfigDto true "配置信息"
// @Success 200 {object} result.Result
// @Router /api/v1/system/config/add [post]
// @Security ApiKeyAuth
func CreateSysConfig(c *gin.Context) {
	var dto model.SysConfigDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	configService := service.NewSysConfigService()
	configService.CreateConfig(c, &dto)
}

// @Tags 系统配置
// @Summary 获取系统配置详情
// @Description 获取系统配置详情
// @Produce json
// @Description 根据ID获取系统配置详情
// @Param id query int true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/system/config/info [get]
// @Security ApiKeyAuth
func GetSysConfigInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	configService := service.NewSysConfigService()
	configService.GetConfigByID(c, uint(id))
}

// @Tags 系统配置
// @Summary 根据配置键获取配置
// @Description 根据配置键获取配置
// @Produce json
// @Description 根据配置键获取系统配置
// @Param key query string true "配置键"
// @Success 200 {object} result.Result
// @Router /api/v1/system/config/key [get]
// @Security ApiKeyAuth
func GetSysConfigByKey(c *gin.Context) {
	configKey := c.Query("key")
	if configKey == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "配置键不能为空")
		return
	}

	configService := service.NewSysConfigService()
	configService.GetConfigByKey(c, configKey)
}

// @Tags 系统配置
// @Summary 更新系统配置
// @Description 更新系统配置
// @Accept json
// @Produce json
// @Description 更新系统配置信息
// @Param data body model.SysConfigDto true "配置信息"
// @Success 200 {object} result.Result
// @Router /api/v1/system/config/update [put]
// @Security ApiKeyAuth
func UpdateSysConfig(c *gin.Context) {
	var dto model.SysConfigDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	if dto.ID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "ID不能为空")
		return
	}

	configService := service.NewSysConfigService()
	configService.UpdateConfig(c, &dto)
}

// @Tags 系统配置
// @Summary 删除系统配置
// @Description 删除系统配置
// @Produce json
// @Description 根据ID删除系统配置
// @Param id query int true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/system/config/delete [delete]
// @Security ApiKeyAuth
func DeleteSysConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	configService := service.NewSysConfigService()
	configService.DeleteConfig(c, uint(id))
}

// @Tags 系统配置
// @Summary 获取系统配置列表
// @Description 获取系统配置列表
// @Produce json
// @Description 获取系统配置列表（支持分页和筛选）
// @Param configKey query string false "配置键(模糊查询)"
// @Param configType query string false "配置类型"
// @Param status query int false "状态(-1:全部,1:启用,2:禁用)"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/system/config/list [get]
// @Security ApiKeyAuth
func GetSysConfigList(c *gin.Context) {
	query := &model.SysConfigQuery{
		ConfigKey:  c.Query("configKey"),
		ConfigType: c.Query("configType"),
		Status:     -1, // 默认查询全部
		PageNum:    1,
		PageSize:   10,
	}

	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if status, err := strconv.Atoi(c.Query("status")); err == nil {
		query.Status = status
	}

	configService := service.NewSysConfigService()
	configService.GetConfigList(c, query)
}

// @Tags 系统配置
// @Summary 更新配置状态
// @Description 更新配置状态
// @Accept json
// @Produce json
// @Description 更新系统配置状态
// @Param data body map[string]interface{} true "状态信息 {id:1, status:2}"
// @Success 200 {object} result.Result
// @Router /api/v1/system/config/updateStatus [put]
// @Security ApiKeyAuth
func UpdateSysConfigStatus(c *gin.Context) {
	var req struct {
		ID     uint `json:"id" validate:"required"`
		Status int  `json:"status" validate:"required,oneof=1 2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	configService := service.NewSysConfigService()
	configService.UpdateConfigStatus(c, req.ID, req.Status)
}

// ============== LDAP 相关接口 ==============

// @Tags LDAP配置
// @Summary 保存LDAP配置
// @Description 保存LDAP配置
// @Accept json
// @Produce json
// @Description 保存或更新LDAP配置
// @Param data body model.LDAPConfigDto true "LDAP配置信息"
// @Success 200 {object} result.Result
// @Router /api/v1/system/ldap/config [post]
// @Security ApiKeyAuth
func SaveLDAPConfig(c *gin.Context) {
	var dto model.LDAPConfigDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	configService := service.NewSysConfigService()
	configService.SaveLDAPConfig(c, &dto)
}

// @Tags LDAP配置
// @Summary 获取LDAP配置
// @Description 获取LDAP配置
// @Produce json
// @Description 获取LDAP配置信息
// @Success 200 {object} result.Result
// @Router /api/v1/system/ldap/config [get]
// @Security ApiKeyAuth
func GetLDAPConfig(c *gin.Context) {
	configService := service.NewSysConfigService()
	configService.GetLDAPConfig(c)
}

// @Tags LDAP配置
// @Summary 测试LDAP连接
// @Description 测试LDAP连接
// @Produce json
// @Description 测试LDAP服务器连接
// @Success 200 {object} result.Result
// @Router /api/v1/system/ldap/test [post]
// @Security ApiKeyAuth
func TestLDAPConnection(c *gin.Context) {
	configService := service.NewSysConfigService()
	configService.TestLDAPConnection(c)
}

// @Tags LDAP配置
// @Summary LDAP认证登录
// @Description LDAP认证登录
// @Accept json
// @Produce json
// @Description 使用LDAP进行用户认证登录
// @Param data body map[string]string true "登录信息 {username:用户名, password:密码}"
// @Success 200 {object} result.Result
// @Router /api/v1/system/ldap/login [post]
func LDAPLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	configService := service.NewSysConfigService()
	configService.LDAPLogin(c, req.Username, req.Password)
}
