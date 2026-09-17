package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// K8sArthasController K8s Arthas诊断控制器
type K8sArthasController struct {
	service service.IK8sArthasService
}

// NewK8sArthasController 创建K8s Arthas诊断控制器
func NewK8sArthasController(db *gorm.DB) *K8sArthasController {
	return &K8sArthasController{
		service: service.NewK8sArthasService(db),
	}
}

// ListJavaProcesses 列表Java进程
// @Summary 列表Java进程
// @Description 列表指定Pod中的Java进程
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/java-processes [get]
func (ctrl *K8sArthasController) ListJavaProcesses(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasCheckQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	processes, err := ctrl.service.ListJavaProcesses(uint(clusterID), query.Namespace, query.Pod, query.Container)
	if err != nil {
		result.Failed(c, 500, "列表Java进程失败: "+err.Error())
		return
	}

	result.Success(c, processes)
}

// CheckArthasInstalled 检查Arthas是否已安装
// @Summary 检查Arthas是否已安装
// @Description 检查指定Pod中Arthas是否已安装
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/check [get]
func (ctrl *K8sArthasController) CheckArthasInstalled(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasCheckQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	checkResult, err := ctrl.service.CheckArthasInstalled(uint(clusterID), query.Namespace, query.Pod, query.Container)
	if err != nil {
		result.Failed(c, 500, "检查Arthas失败: "+err.Error())
		return
	}

	result.Success(c, checkResult)
}

// InstallArthas 安装Arthas
// @Summary 安装Arthas
// @Description 在指定Pod中安装Arthas诊断工具
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param body body model.ArthasCheckQuery true "安装请求"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/install [post]
func (ctrl *K8sArthasController) InstallArthas(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var req model.ArthasCheckQuery
	if err = c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	err = ctrl.service.InstallArthas(uint(clusterID), req.Namespace, req.Pod, req.Container)
	if err != nil {
		result.Failed(c, 500, "安装Arthas失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{"message": "Arthas安装成功"})
}

// ExecuteArthasCommand 执行Arthas自定义命令
// @Summary 执行Arthas自定义命令
// @Description 在指定Pod的Java进程中执行Arthas命令
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param body body model.ArthasCommandRequest true "命令请求"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/command [post]
func (ctrl *K8sArthasController) ExecuteArthasCommand(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var req model.ArthasCommandRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.ExecuteArthasCommand(uint(clusterID), req.Namespace, req.Pod, req.Container, req.ProcessID, req.Command)
	if err != nil {
		result.Failed(c, 500, "执行命令失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetDashboard 获取仪表板
// @Summary 获取仪表板
// @Description 获取Java进程仪表板信息
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/dashboard [get]
func (ctrl *K8sArthasController) GetDashboard(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasBaseQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.GetDashboard(uint(clusterID), query.Namespace, query.Pod, query.Container, query.ProcessID)
	if err != nil {
		result.Failed(c, 500, "获取仪表板失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetThreadList 获取线程列表
// @Summary 获取线程列表
// @Description 获取Java进程的线程列表
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/thread [get]
func (ctrl *K8sArthasController) GetThreadList(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasBaseQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.GetThreadList(uint(clusterID), query.Namespace, query.Pod, query.Container, query.ProcessID)
	if err != nil {
		result.Failed(c, 500, "获取线程列表失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetThreadStack 获取线程堆栈
// @Summary 获取线程堆栈
// @Description 获取指定线程的堆栈跟踪
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Param threadId query string true "线程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/thread/stack [get]
func (ctrl *K8sArthasController) GetThreadStack(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("container")
	processID := c.Query("processId")
	threadID := c.Query("threadId")

	if namespace == "" || pod == "" || container == "" || processID == "" || threadID == "" {
		result.Failed(c, 400, "必填参数不能为空")
		return
	}

	output, err := ctrl.service.GetThreadStack(uint(clusterID), namespace, pod, container, processID, threadID)
	if err != nil {
		result.Failed(c, 500, "获取线程堆栈失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetJvmInfo 获取JVM信息
// @Summary 获取JVM信息
// @Description 获取Java进程的JVM信息
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/jvm [get]
func (ctrl *K8sArthasController) GetJvmInfo(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasBaseQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.GetJvmInfo(uint(clusterID), query.Namespace, query.Pod, query.Container, query.ProcessID)
	if err != nil {
		result.Failed(c, 500, "获取JVM信息失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetSysEnv 获取系统环境变量
// @Summary 获取系统环境变量
// @Description 获取Java进程的系统环境变量
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/sysenv [get]
func (ctrl *K8sArthasController) GetSysEnv(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasBaseQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.GetSysEnv(uint(clusterID), query.Namespace, query.Pod, query.Container, query.ProcessID)
	if err != nil {
		result.Failed(c, 500, "获取系统环境变量失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetSysProp 获取系统属性
// @Summary 获取系统属性
// @Description 获取Java进程的系统属性
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/sysprop [get]
func (ctrl *K8sArthasController) GetSysProp(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasBaseQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.GetSysProp(uint(clusterID), query.Namespace, query.Pod, query.Container, query.ProcessID)
	if err != nil {
		result.Failed(c, 500, "获取系统属性失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetMemory 获取内存信息
// @Summary 获取内存信息
// @Description 获取Java进程的内存信息
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/memory [get]
func (ctrl *K8sArthasController) GetMemory(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasBaseQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.GetMemory(uint(clusterID), query.Namespace, query.Pod, query.Container, query.ProcessID)
	if err != nil {
		result.Failed(c, 500, "获取内存信息失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetPerfCounter 获取性能计数器
// @Summary 获取性能计数器
// @Description 获取Java进程的性能计数器
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/perfcounter [get]
func (ctrl *K8sArthasController) GetPerfCounter(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var query model.ArthasBaseQuery
	if err = c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, 400, "参数绑定失败: "+err.Error())
		return
	}

	output, err := ctrl.service.GetPerfCounter(uint(clusterID), query.Namespace, query.Pod, query.Container, query.ProcessID)
	if err != nil {
		result.Failed(c, 500, "获取性能计数器失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// DecompileClass 反编译类
// @Summary 反编译类
// @Description 反编译指定的Java类
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Param className query string true "类名"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/jad [get]
func (ctrl *K8sArthasController) DecompileClass(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("container")
	processID := c.Query("processId")
	className := c.Query("className")

	if namespace == "" || pod == "" || container == "" || processID == "" || className == "" {
		result.Failed(c, 400, "必填参数不能为空")
		return
	}

	output, err := ctrl.service.DecompileClass(uint(clusterID), namespace, pod, container, processID, className)
	if err != nil {
		result.Failed(c, 500, "反编译类失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// GetStaticField 获取静态字段
// @Summary 获取静态字段
// @Description 获取指定类的静态字段值
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Param className query string true "类名"
// @Param fieldName query string true "字段名"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/getstatic [get]
func (ctrl *K8sArthasController) GetStaticField(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("container")
	processID := c.Query("processId")
	className := c.Query("className")
	fieldName := c.Query("fieldName")

	if namespace == "" || pod == "" || container == "" || processID == "" || className == "" || fieldName == "" {
		result.Failed(c, 400, "必填参数不能为空")
		return
	}

	output, err := ctrl.service.GetStaticField(uint(clusterID), namespace, pod, container, processID, className, fieldName)
	if err != nil {
		result.Failed(c, 500, "获取静态字段失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// SearchClass 搜索类
// @Summary 搜索类
// @Description 根据模式搜索类
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Param pattern query string true "搜索模式"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/sc [get]
func (ctrl *K8sArthasController) SearchClass(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("container")
	processID := c.Query("processId")
	pattern := c.Query("pattern")

	if namespace == "" || pod == "" || container == "" || processID == "" || pattern == "" {
		result.Failed(c, 400, "必填参数不能为空")
		return
	}

	output, err := ctrl.service.SearchClass(uint(clusterID), namespace, pod, container, processID, pattern)
	if err != nil {
		result.Failed(c, 500, "搜索类失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// SearchMethod 搜索方法
// @Summary 搜索方法
// @Description 搜索指定类中的方法
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Param className query string true "类名"
// @Param methodName query string true "方法名"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/sm [get]
func (ctrl *K8sArthasController) SearchMethod(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("container")
	processID := c.Query("processId")
	className := c.Query("className")
	methodName := c.Query("methodName")

	if namespace == "" || pod == "" || container == "" || processID == "" || className == "" || methodName == "" {
		result.Failed(c, 400, "必填参数不能为空")
		return
	}

	output, err := ctrl.service.SearchMethod(uint(clusterID), namespace, pod, container, processID, className, methodName)
	if err != nil {
		result.Failed(c, 500, "搜索方法失败: "+err.Error())
		return
	}

	result.Success(c, output)
}

// StartFlameGraph 异步启动火焰图
// @Summary 启动火焰图采样
// @Description 异步启动CPU火焰图采样，立即返回taskId，通过查询接口获取结果
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Param duration query int true "采样时间（秒），支持 10/30/60/120"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/profiler/start [post]
func (ctrl *K8sArthasController) StartFlameGraph(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("container")
	processID := c.Query("processId")
	durationStr := c.Query("duration")

	if namespace == "" || pod == "" || container == "" || processID == "" || durationStr == "" {
		result.Failed(c, 400, "必填参数不能为空")
		return
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		result.Failed(c, 400, "duration必须是正整数")
		return
	}

	task, err := ctrl.service.StartFlameGraph(uint(clusterID), namespace, pod, container, processID, duration)
	if err != nil {
		result.Failed(c, 500, "启动火焰图失败: "+err.Error())
		return
	}

	result.Success(c, task)
}

// GetFlameGraphTask 查询火焰图任务
// @Summary 查询火焰图任务状态
// @Description 轮询查询火焰图采样任务，status=done时返回火焰图数据
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param taskId path string true "任务ID"
// @Success 200 {object} result.Result
// @Failure 404 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/profiler/task/{taskId} [get]
func (ctrl *K8sArthasController) GetFlameGraphTask(c *gin.Context) {
	taskId := c.Param("taskId")
	if taskId == "" {
		result.Failed(c, 400, "taskId不能为空")
		return
	}

	task, err := ctrl.service.GetFlameGraphTask(taskId)
	if err != nil {
		result.Failed(c, 404, err.Error())
		return
	}

	result.Success(c, task)
}

// ArthasWebSocket WebSocket实时交互
// @Summary WebSocket实时交互
// @Description 通过WebSocket与Arthas进行实时交互（trace/watch/monitor等）
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespace query string true "命名空间"
// @Param pod query string true "Pod名称"
// @Param container query string true "容器名称"
// @Param processId query string true "进程ID"
// @Success 101 "Switching Protocols"
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/arthas/ws [get]
func (ctrl *K8sArthasController) ArthasWebSocket(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("container")
	processID := c.Query("processId")

	if namespace == "" || pod == "" || container == "" || processID == "" {
		result.Failed(c, 400, "必填参数不能为空")
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "WebSocket升级失败: "+err.Error())
		return
	}

	// 启动WebSocket会话处理
	ctrl.service.ArthasWebSocket(uint(clusterID), namespace, pod, container, processID, conn)
}

// DownloadArthas 下载 arthas-boot.jar
// @Summary 下载 arthas-boot.jar
// @Description 从阿里云下载 arthas-boot.jar 到 upload/arthas 目录
// @Tags K8s应用诊断
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Success 200 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/arthas/download [post]
func (ctrl *K8sArthasController) DownloadArthas(c *gin.Context) {
	alreadyExists, err := ctrl.service.DownloadArthas()
	if err != nil {
		result.Failed(c, 500, "下载失败: "+err.Error())
		return
	}

	jarPath := "upload/arthas/arthas-boot.jar"
	fileInfo, _ := os.Stat(jarPath)
	sizeStr := ""
	if fileInfo != nil {
		kb := float64(fileInfo.Size()) / 1024
		if kb >= 1024 {
			sizeStr = fmt.Sprintf("%.1f MB", kb/1024)
		} else {
			sizeStr = fmt.Sprintf("%.0f KB", kb)
		}
	}

	if alreadyExists {
		result.Success(c, gin.H{
			"message": "✓ arthas-boot.jar 已存在，无需重复下载",
			"path":    jarPath,
			"size":    sizeStr,
		})
		return
	}

	result.Success(c, gin.H{
		"message": "✓ Arthas 下载成功，已保存到 upload/arthas/arthas-boot.jar",
		"path":    jarPath,
		"size":    sizeStr,
	})
}

// UploadArthas 上传 arthas-boot.jar
// @Summary 上传 arthas-boot.jar
// @Description 上传 arthas-boot.jar 到 upload/arthas 目录
// @Tags K8s应用诊断
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param file formData file true "arthas-boot.jar 文件"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/arthas/upload [post]
func (ctrl *K8sArthasController) UploadArthas(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		result.Failed(c, 400, "获取文件失败: "+err.Error())
		return
	}

	// 检查文件名
	if file.Filename != "arthas-boot.jar" {
		result.Failed(c, 400, "文件名必须是 arthas-boot.jar")
		return
	}

	// 检查文件大小（arthas-boot.jar 是 bootstrap 启动器，大小 100KB-200KB）
	if file.Size < 50*1024 { // 小于 50KB
		result.Failed(c, 400, "文件过小，可能不是有效的 arthas-boot.jar")
		return
	}

	if file.Size > 500*1024*1024 { // 大于 500MB
		result.Failed(c, 400, "文件过大，请检查是否正确")
		return
	}

	// 保存文件到 upload/arthas 目录
	savePath := "upload/arthas/arthas-boot.jar"
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		result.Failed(c, 500, "保存文件失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"message": "✓ Arthas 上传成功",
		"path":    savePath,
		"size":    file.Size,
	})
}
