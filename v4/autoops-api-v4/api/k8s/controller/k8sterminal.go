package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// K8sTerminalController 容器终端控制器
type K8sTerminalController struct {
	service service.IK8sTerminalService
}

func NewK8sTerminalController(db *gorm.DB) *K8sTerminalController {
	return &K8sTerminalController{
		service: service.NewK8sTerminalService(db),
	}
}

func normalizeTerminalCommand(command, shell string) (string, error) {
	if command != "" {
		switch strings.TrimSpace(strings.ToLower(command)) {
		case "bash", "/bin/bash":
			return "/bin/bash", nil
		case "sh", "/bin/sh":
			return "/bin/sh", nil
		default:
			return "", fmt.Errorf("仅支持 /bin/bash 或 /bin/sh")
		}
	}

	switch strings.TrimSpace(strings.ToLower(shell)) {
	case "", "bash", "/bin/bash":
		return "/bin/bash", nil
	case "sh", "/bin/sh":
		return "/bin/sh", nil
	default:
		return "", fmt.Errorf("shell 仅支持 bash 或 sh")
	}
}

// ConnectPodTerminal 连接到Pod终端
// @Summary 连接到Pod终端
// @Description 通过WebSocket连接到指定Pod的终端
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName query string false "容器名称（默认为Pod中第一个容器）"
// @Param shell query string false "终端类型，可选 bash/sh，默认 bash"
// @Param command query string false "执行命令，兼容旧参数，仅支持 /bin/bash 或 /bin/sh"
// @Success 101 "Switching Protocols"
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/terminal [get]
func (ctrl *K8sTerminalController) ConnectPodTerminal(c *gin.Context) {
	// 参数验证
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	containerName := c.Query("containerName")
	command, err := normalizeTerminalCommand(c.Query("command"), c.DefaultQuery("shell", "bash"))
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// 创建K8s WebSocket流
	stream, err := ctrl.service.CreateK8sWebSocketStream(uint(clusterId), namespaceName, podName, containerName, command, conn)
	if err != nil {
		log.Printf("Failed to create K8s WebSocket stream: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	defer stream.Close()

	// 等待连接关闭
	select {
	case <-stream.Ctx.Done():
		log.Printf("K8s terminal connection closed")
	}
}

// GetPodContainers 获取Pod中的容器列表
// @Summary 获取Pod中的容器列表
// @Description 获取指定Pod中所有容器的名称列表，用于终端连接时选择容器
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Success 200 {object} result.Result{data=[]string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/containers [get]
func (ctrl *K8sTerminalController) GetPodContainers(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	containers, err := ctrl.service.GetPodContainers(uint(clusterId), namespaceName, podName)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("获取容器列表失败: %v", err))
		return
	}

	result.Success(c, containers)
}

// CreatePodDebugContainer 创建Pod临时调试容器
// @Summary 创建Pod临时调试容器
// @Description 在指定Running Pod中注入一个临时调试容器，并返回终端连接路径
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param request body model.CreatePodDebugContainerRequest false "调试容器创建参数"
// @Success 200 {object} result.Result{data=model.CreatePodDebugContainerResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/debug [post]
func (ctrl *K8sTerminalController) CreatePodDebugContainer(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	var req model.CreatePodDebugContainerRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			result.Failed(c, 400, "请求参数错误: "+err.Error())
			return
		}
	}

	resp, err := ctrl.service.CreatePodDebugContainer(uint(clusterId), namespaceName, podName, &req)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("创建临时调试容器失败: %v", err))
		return
	}

	result.Success(c, resp)
}

// DeletePodDebugContainer 结束Pod临时调试容器
// @Summary 结束Pod临时调试容器
// @Description 结束指定的临时调试容器。注意：ephemeral container 不会从 Pod 定义中物理移除，只会变为 terminated
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName path string true "调试容器名称"
// @Success 200 {object} result.Result{data=model.DeletePodDebugContainerResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/debug/{containerName} [delete]
func (ctrl *K8sTerminalController) DeletePodDebugContainer(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	containerName := c.Param("containerName")
	if containerName == "" {
		result.Failed(c, 400, "调试容器名称不能为空")
		return
	}

	resp, err := ctrl.service.DeletePodDebugContainer(uint(clusterId), namespaceName, podName, containerName)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("结束临时调试容器失败: %v", err))
		return
	}

	result.Success(c, resp)
}
