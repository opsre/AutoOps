package controller

import (
	"encoding/json"
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// NotifyRobotController 通知机器人控制器（统一管理三种类型机器人）
type NotifyRobotController struct {
	robotService service.NotifyRobotService
}

// NewNotifyRobotController 创建通知机器人控制器实例
func NewNotifyRobotController() *NotifyRobotController {
	return &NotifyRobotController{
		robotService: service.NewNotifyRobotService(),
	}
}

// @Tags 通知机器人
// @Summary 创建通知机器人
// @Description 创建通知机器人
// @Accept json
// @Produce json
// @Description 创建通知机器人（支持钉钉/飞书/企业微信/邮箱/通用Webhook/Teams）
// @Param data body model.NotifyRobot true "机器人信息（type字段: dingtalk/feishu/wechat/email/webhook/teams）"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/robot [post]
// @Security ApiKeyAuth
func (ctrl *NotifyRobotController) CreateRobot(c *gin.Context) {
	// 先绑定到 map 进行预处理
	var rawData map[string]interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	// 检查类型字段
	robotType, ok := rawData["type"].(string)
	if !ok {
		result.Failed(c, int(result.ApiCode.FAILED), "缺少机器人类型")
		return
	}

	// 验证机器人类型
	validTypes := []string{"dingtalk", "feishu", "wechat", "email", "webhook", "teams"}
	isValidType := false
	for _, t := range validTypes {
		if robotType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		result.Failed(c, int(result.ApiCode.FAILED), "机器人类型错误，仅支持: dingtalk/feishu/wechat/email/webhook/teams")
		return
	}

	// 如果不是邮箱类型，删除 port 字段（避免空字符串导致的 JSON 反序列化错误）
	if robotType != "email" {
		delete(rawData, "port")
	}

	// 将 map 转换为 JSON 再转换为结构体
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "数据格式错误")
		return
	}

	var robot model.NotifyRobot
	if err := json.Unmarshal(jsonData, &robot); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.robotService.Create(c, &robot, robot.Type)
}

// @Tags 通知机器人
// @Summary 获取通知机器人详情
// @Description 获取通知机器人详情
// @Produce json
// @Description 根据ID获取通知机器人详情
// @Param id path int true "机器人ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/robot/{id} [get]
// @Security ApiKeyAuth
func (ctrl *NotifyRobotController) GetRobot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	// 不指定类型，返回任意类型的机器人
	ctrl.robotService.GetByID(c, uint(id), "")
}

// @Tags 通知机器人
// @Summary 更新通知机器人
// @Description 更新通知机器人
// @Accept json
// @Produce json
// @Description 更新通知机器人信息
// @Param data body model.NotifyRobot true "机器人信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/robot [put]
// @Security ApiKeyAuth
func (ctrl *NotifyRobotController) UpdateRobot(c *gin.Context) {
	// 先绑定到 map 进行预处理
	var rawData map[string]interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	// 检查 ID
	id, ok := rawData["id"].(float64)
	if !ok || id == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "ID不能为空")
		return
	}

	// 检查类型字段
	robotType, ok := rawData["type"].(string)
	if !ok {
		result.Failed(c, int(result.ApiCode.FAILED), "缺少机器人类型")
		return
	}

	// 验证机器人类型
	validTypes := []string{"dingtalk", "feishu", "wechat", "email", "webhook", "teams"}
	isValidType := false
	for _, t := range validTypes {
		if robotType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		result.Failed(c, int(result.ApiCode.FAILED), "机器人类型错误，仅支持: dingtalk/feishu/wechat/email/webhook/teams")
		return
	}

	// 如果不是邮箱类型，删除 port 字段（避免空字符串导致的 JSON 反序列化错误）
	if robotType != "email" {
		delete(rawData, "port")
	}

	// 将 map 转换为 JSON 再转换为结构体
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "数据格式错误")
		return
	}

	var robot model.NotifyRobot
	if err := json.Unmarshal(jsonData, &robot); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.robotService.Update(c, &robot, robot.Type)
}

// @Tags 通知机器人
// @Summary 删除通知机器人
// @Description 删除通知机器人
// @Produce json
// @Description 根据ID删除通知机器人
// @Param id path int true "机器人ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/robot/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *NotifyRobotController) DeleteRobot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	// 不指定类型，删除任意类型的机器人
	ctrl.robotService.Delete(c, uint(id), "")
}

// @Tags 通知机器人
// @Summary 获取通知机器人列表
// @Description 获取通知机器人列表
// @Produce json
// @Description 分页获取通知机器人列表
// @Param type query string false "机器人类型(dingtalk/feishu/wechat/email/webhook/teams)，不传则查询所有类型"
// @Param status query int false "状态(1:启用,0:禁用)"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/robot/list [get]
// @Security ApiKeyAuth
func (ctrl *NotifyRobotController) GetRobotList(c *gin.Context) {
	// 获取机器人类型（可选）
	robotType := c.Query("type")
	validTypes := []string{"dingtalk", "feishu", "wechat", "email", "webhook", "teams"}
	if robotType != "" {
		isValidType := false
		for _, t := range validTypes {
			if robotType == t {
				isValidType = true
				break
			}
		}
		if !isValidType {
			result.Failed(c, int(result.ApiCode.FAILED), "机器人类型错误，仅支持: dingtalk/feishu/wechat/email/webhook/teams")
			return
		}
	}

	// 获取状态参数
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	ctrl.robotService.GetList(c, robotType, status, page, pageSize)
}

// @Tags 通知机器人
// @Summary 测试发送消息
// @Description 测试发送消息
// @Accept json
// @Produce json
// @Description 测试发送消息到指定机器人（可选：传入临时模板进行测试）
// @Param id path int true "机器人ID"
// @Param data body map[string]string false "可选参数：template-临时模板（不影响数据库中的模板）"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/robot/test/{id} [post]
// @Security ApiKeyAuth
func (ctrl *NotifyRobotController) TestSendMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	// 尝试解析请求体，获取临时模板（可选）
	var requestBody map[string]string
	var tempTemplate string
	if err := c.ShouldBindJSON(&requestBody); err == nil {
		tempTemplate = requestBody["template"]
	}

	ctrl.robotService.TestSend(c, uint(id), tempTemplate)
}
