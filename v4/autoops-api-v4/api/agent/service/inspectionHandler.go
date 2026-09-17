// 巡检处理器 - 实现AgentHandler接口
package service

import (
	"context"
	"encoding/json"
	"fmt"

	"dodevops-api/api/agent/model"
)

// InspectionHandler 巡检处理器
type InspectionHandler struct {
	inspectionService InspectionService
}

// NewInspectionHandler 创建巡检处理器
func NewInspectionHandler() *InspectionHandler {
	return &InspectionHandler{
		inspectionService: NewInspectionService(),
	}
}

// Name 处理器名称
func (h *InspectionHandler) Name() string {
	return "inspection"
}

// Execute 执行任务
func (h *InspectionHandler) Execute(ctx context.Context, task *model.AgentTask) error {
	// 解析任务参数
	var params struct {
		InspectionID uint `json:"inspectionId"`
	}

	if err := json.Unmarshal([]byte(task.Params), &params); err != nil {
		return fmt.Errorf("解析任务参数失败: %v", err)
	}

	if params.InspectionID == 0 {
		return fmt.Errorf("巡检任务ID不能为空")
	}

	// 执行巡检
	err := h.inspectionService.ExecuteInspection(ctx, params.InspectionID)
	if err != nil {
		return err
	}

	// 保存执行结果
	result := map[string]interface{}{
		"inspectionId": params.InspectionID,
		"status":       "success",
		"message":      "巡检任务执行完成",
	}

	resultJSON, _ := json.Marshal(result)
	task.Result = string(resultJSON)

	return nil
}

// Validate 验证任务
func (h *InspectionHandler) Validate(task *model.AgentTask) error {
	if task.Type != "inspection" {
		return fmt.Errorf("任务类型错误")
	}

	// 解析参数
	var params struct {
		InspectionID uint `json:"inspectionId"`
	}

	if err := json.Unmarshal([]byte(task.Params), &params); err != nil {
		return fmt.Errorf("参数格式错误: %v", err)
	}

	if params.InspectionID == 0 {
		return fmt.Errorf("巡检任务ID不能为空")
	}

	return nil
}

// CanRetry 是否可以重试
func (h *InspectionHandler) CanRetry(task *model.AgentTask, err error) bool {
	// 网络错误、超时等可以重试
	if err != nil {
		errMsg := err.Error()
		// 可重试的错误类型
		retryableErrors := []string{
			"timeout",
			"connection refused",
			"network",
			"temporarily unavailable",
		}

		for _, retryErr := range retryableErrors {
			if contains(errMsg, retryErr) {
				return true
			}
		}
	}

	// 其他错误不重试
	return false
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
