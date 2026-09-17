package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/k8s/service"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type K8sNodesController struct {
	service service.IK8sNodesService
}

func NewK8sNodesController(db *gorm.DB) *K8sNodesController {
	return &K8sNodesController{
		service: service.NewK8sNodesService(db),
	}
}

// GetNodes 获取K8s节点信息
// @Summary 获取K8s节点信息
// @Description 获取指定集群的所有节点详细信息，包括名称/IP地址、状态、配置、资源使用情况等
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Success 200 {object} result.Result{data=[]model.K8sNode} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) GetNodes(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	ctrl.service.GetNodes(c, uint(clusterId))
}

// GetNodeDetail 获取单个节点详细信息
// @Summary 获取单个节点详细信息
// @Description 获取指定节点的详细信息，包括容器组、资源使用详情等
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Success 200 {object} result.Result{data=model.K8sNodeDetail} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName} [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) GetNodeDetail(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	ctrl.service.GetNodeDetail(c, uint(clusterId), nodeName)
}

// GetNodeDetailEnhanced 获取增强的节点详细信息
// @Summary 获取增强的节点详细信息
// @Description 获取节点的完整详细信息，包括基本信息、系统信息、K8s组件版本、资源使用情况、监控信息、Pod列表等
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Success 200 {object} result.Result{data=model.NodeDetailResponse} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/enhanced [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) GetNodeDetailEnhanced(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	ctrl.service.GetNodeDetailEnhanced(c, uint(clusterId), nodeName)
}

// AddTaint 为节点添加污点
// @Summary 为节点添加污点
// @Description 为指定节点添加污点，控制Pod调度
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Param taint body model.AddTaintRequest true "污点信息"
// @Success 200 {object} result.Result "添加成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/taints [post]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) AddTaint(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	var req model.AddTaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.AddTaint(c, uint(clusterId), nodeName, &req)
}

// RemoveTaint 移除节点污点
// @Summary 移除节点污点
// @Description 移除指定节点的污点
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Param taint body model.RemoveTaintRequest true "污点信息"
// @Success 200 {object} result.Result "移除成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/taints [delete]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) RemoveTaint(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	var req model.RemoveTaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.RemoveTaint(c, uint(clusterId), nodeName, &req)
}

// AddLabel 为节点添加标签
// @Summary 为节点添加标签
// @Description 为指定节点添加标签
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Param label body model.AddLabelRequest true "标签信息"
// @Success 200 {object} result.Result "添加成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/labels [post]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) AddLabel(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	var req model.AddLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.AddLabel(c, uint(clusterId), nodeName, &req)
}

// RemoveLabel 移除节点标签
// @Summary 移除节点标签
// @Description 移除指定节点的标签
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Param label body model.RemoveLabelRequest true "标签信息"
// @Success 200 {object} result.Result "移除成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/labels [delete]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) RemoveLabel(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	var req model.RemoveLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.RemoveLabel(c, uint(clusterId), nodeName, &req)
}

// CordonNode 封锁/解封节点
// @Summary 封锁/解封节点
// @Description 设置节点的可调度状态
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Param cordon body model.CordonNodeRequest true "封锁信息"
// @Success 200 {object} result.Result "操作成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/cordon [post]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) CordonNode(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	var req model.CordonNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CordonNode(c, uint(clusterId), nodeName, &req)
}

// DrainNode 驱逐节点
// @Summary 驱逐节点
// @Description 驱逐节点上的所有Pod
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Param drain body model.DrainNodeRequest true "驱逐配置"
// @Success 200 {object} result.Result "驱逐成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/drain [post]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) DrainNode(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	var req model.DrainNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.DrainNode(c, uint(clusterId), nodeName, &req)
}

// GetNodeResourceAllocation 获取节点资源分配详情
// @Summary 获取节点资源分配详情
// @Description 获取节点的资源容量、分配情况和Pod资源使用详情
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Success 200 {object} result.Result{data=model.NodeResourceAllocation} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/resources [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) GetNodeResourceAllocation(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	ctrl.service.GetNodeResourceAllocation(c, uint(clusterId), nodeName)
}

// GetNodeYaml 获取节点YAML
// @Summary 获取节点YAML配置
// @Description 获取指定节点的YAML配置内容
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Success 200 {object} result.Result{data=map[string]string} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/yaml [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) GetNodeYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	ctrl.service.GetNodeYaml(c, uint(clusterId), nodeName)
}

// UpdateNodeYaml 更新节点YAML
// @Summary 更新节点YAML配置
// @Description 通过YAML内容更新指定节点的配置
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param nodeName path string true "节点名称"
// @Param yaml body model.UpdateNodeYAMLRequest true "YAML内容"
// @Success 200 {object} result.Result{data=model.UpdateNodeYAMLResponse} "更新成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "节点不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/{nodeName}/yaml [put]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) UpdateNodeYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	var req model.UpdateNodeYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateNodeYaml(c, uint(clusterId), nodeName, &req)
}

// SearchByName 根据名称模糊查询节点
// @Summary 根据名称模糊查询节点
// @Description 按节点名称进行模糊搜索（不区分大小写）
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param name query string true "节点名称（支持模糊匹配）"
// @Success 200 {object} result.Result{data=[]model.K8sNode} "查询成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/search [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) SearchByName(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	name := c.Query("name")
	if name == "" {
		result.Failed(c, http.StatusBadRequest, "节点名称不能为空")
		return
	}

	ctrl.service.SearchNodesByName(c, uint(clusterId), name)
}

// SearchByRole 根据角色查询节点
// @Summary 根据角色查询节点
// @Description 按节点角色进行搜索（如 master, worker, control-plane）
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param role query string true "节点角色（如 master, worker, control-plane）"
// @Success 200 {object} result.Result{data=[]model.K8sNode} "查询成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/byRole [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) SearchByRole(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	role := c.Query("role")
	if role == "" {
		result.Failed(c, http.StatusBadRequest, "节点角色不能为空")
		return
	}

	ctrl.service.SearchNodesByRole(c, uint(clusterId), role)
}

// SearchByStatus 根据状态查询节点
// @Summary 根据状态查询节点
// @Description 按节点状态进行搜索（如 Ready, NotReady）
// @Tags K8s节点管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param status query string true "节点状态（如 Ready, NotReady）"
// @Success 200 {object} result.Result{data=[]model.K8sNode} "查询成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/nodes/byStatus [get]
// @Security ApiKeyAuth
func (ctrl *K8sNodesController) SearchByStatus(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	status := c.Query("status")
	if status == "" {
		result.Failed(c, http.StatusBadRequest, "节点状态不能为空")
		return
	}

	ctrl.service.SearchNodesByStatus(c, uint(clusterId), status)
}