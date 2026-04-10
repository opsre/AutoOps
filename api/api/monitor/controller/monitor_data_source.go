package controller

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MonitorDataSourceDTO struct {
	ID           uint        `json:"id"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	DeployMethod string      `json:"deployMethod"`
	Url          string      `json:"url"`
	ApiUrl       string      `json:"apiUrl"`
	Config       interface{} `json:"config"`
	Status       int         `json:"status"`
}

func (dto *MonitorDataSourceDTO) ToModel() *model.MonitorDataSource {
	var configStr string
	switch v := dto.Config.(type) {
	case string:
		configStr = v
	case map[string]interface{}:
		b, _ := json.Marshal(v)
		configStr = string(b)
	default:
		configStr = "{}"
	}

	apiUrl := dto.ApiUrl
	if apiUrl == "" && dto.Url != "" {
		apiUrl = dto.Url
	}
	return &model.MonitorDataSource{
		ID:           dto.ID,
		Name:         dto.Name,
		Type:         dto.Type,
		DeployMethod: dto.DeployMethod,
		ApiUrl:       apiUrl,
		Config:       configStr,
		Status:       dto.Status,
	}
}
type MonitorDataSourceController struct {
	dataSourceService service.MonitorDataSourceService
}

func NewMonitorDataSourceController() *MonitorDataSourceController {
	return &MonitorDataSourceController{
		dataSourceService: service.NewMonitorDataSourceService(),
	}
}

// Create 创建监控数据源
func (c *MonitorDataSourceController) Create(ctx *gin.Context) {
	var dto MonitorDataSourceDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "参数绑定失败")
		return
	}

	data := dto.ToModel()
	if err := c.dataSourceService.Create(data); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "创建监控数据源失败")
		return
	}
	result.Success(ctx, data)
}

// Delete 删除监控数据源
func (c *MonitorDataSourceController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "无效的ID")
		return
	}

	if err := c.dataSourceService.Delete(uint(id)); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "删除监控数据源失败")
		return
	}
	result.Success(ctx, nil)
}

// Update 更新监控数据源
func (c *MonitorDataSourceController) Update(ctx *gin.Context) {
	var dto MonitorDataSourceDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "参数绑定失败")
		return
	}

	data := dto.ToModel()
	if data.ID == 0 {
		result.Failed(ctx, int(result.ApiCode.FAILED), "监控数据源ID不能为空")
		return
	}

	if err := c.dataSourceService.Update(data); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "更新监控数据源失败")
		return
	}
	result.Success(ctx, data)
}

// GetByID 获取单个监控数据源
func (c *MonitorDataSourceController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "无效的ID")
		return
	}

	data, err := c.dataSourceService.GetByID(uint(id))
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "获取监控数据源失败")
		return
	}
	result.Success(ctx, data)
}

// GetList 获取监控数据源列表
func (c *MonitorDataSourceController) GetList(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	list, total, err := c.dataSourceService.GetList(page, pageSize)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "获取监控数据源列表失败")
		return
	}
	result.SuccessWithPage(ctx, list, total, page, pageSize)
}
