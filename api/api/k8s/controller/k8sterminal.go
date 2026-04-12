package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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
// @Param command query string false "执行命令（默认为/bin/bash）"
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
	command := c.Query("command")
	if command == "" {
		command = "/bin/bash"
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

// GetPodFileList 获取Pod容器中的文件列表
// @Summary 获取Pod容器中的文件列表
// @Description 通过执行ls命令获取Pod容器指定路径下的文件列表
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName query string false "容器名称（默认为Pod中第一个容器）"
// @Param path query string false "目标路径（默认为/）"
// @Success 200 {object} result.Result{data=[]map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/list [get]
func (ctrl *K8sTerminalController) GetPodFileList(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.Query("containerName")
	path := c.Query("path")

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}

	files, err := ctrl.service.GetPodFileList(uint(clusterId), namespaceName, podName, containerName, path)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("获取文件列表失败: %v", err))
		return
	}
	result.Success(c, files)
}

// DeletePodFile 删除Pod容器中的文件
// @Summary 删除Pod容器中的文件
// @Description 通过执行rm命令删除Pod容器指定路径下的文件
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName query string false "容器名称（默认为Pod中第一个容器）"
// @Param path query string true "要删除的目标路径"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files [delete]
func (ctrl *K8sTerminalController) DeletePodFile(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.Query("containerName")
	path := c.Query("path")

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}
	if path == "" {
		result.Failed(c, 400, "删除路径不能为空")
		return
	}

	err = ctrl.service.DeletePodFile(uint(clusterId), namespaceName, podName, containerName, path)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("删除文件失败: %v", err))
		return
	}
	result.Success(c, "文件删除成功")
}

// GetPodFileContent 获取Pod容器中的文件内容
// @Summary 获取Pod容器中的文件内容
// @Description 通过执行cat命令获取Pod容器指定路径下的文件内容
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName query string false "容器名称（默认为Pod中第一个容器）"
// @Param path query string true "目标文件路径"
// @Success 200 {object} result.Result{data=string}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/content [get]
func (ctrl *K8sTerminalController) GetPodFileContent(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.Query("containerName")
	path := c.Query("path")

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}
	if path == "" {
		result.Failed(c, 400, "文件路径不能为空")
		return
	}

	content, err := ctrl.service.GetPodFileContent(uint(clusterId), namespaceName, podName, containerName, path)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("获取文件内容失败: %v", err))
		return
	}
	result.Success(c, content)
}

// UpdatePodFileContent 更新Pod容器中的文件内容
// @Summary 更新Pod容器中的文件内容
// @Description 将新内容覆盖写入到容器的指定文件中
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/content [put]
func (ctrl *K8sTerminalController) UpdatePodFileContent(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")

	var req struct {
		ContainerName string `json:"containerName"`
		Path          string `json:"path"`
		Content       string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数解析失败")
		return
	}

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}
	if req.Path == "" {
		result.Failed(c, 400, "文件路径不能为空")
		return
	}

	err = ctrl.service.UpdatePodFileContent(uint(clusterId), namespaceName, podName, req.ContainerName, req.Path, req.Content)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("更新文件内容失败: %v", err))
		return
	}
	result.Success(c, "文件更新成功")
}

// UploadPodFile 上传文件到Pod容器
// @Summary 上传文件到Pod容器
// @Description 通过上传文件将其覆盖写入到容器的指定路径中
// @Tags K8s容器终端
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName formData string false "容器名称（默认为Pod中第一个容器）"
// @Param path formData string true "目标文件路径"
// @Param file formData file true "欲上传的文件"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/upload [post]
func (ctrl *K8sTerminalController) UploadPodFile(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.PostForm("containerName")
	path := c.PostForm("path")

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}
	if path == "" {
		result.Failed(c, 400, "文件路径不能为空")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		result.Failed(c, 400, "未找到文件，请检查表单是否包含名为`file`的文件对象")
		return
	}
	f, err := file.Open()
	if err != nil {
		result.Failed(c, 500, "无法读取文件数据")
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		result.Failed(c, 500, "读取文件字节失败")
		return
	}

	err = ctrl.service.UploadPodFile(uint(clusterId), namespaceName, podName, containerName, path, content)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("上传文件失败: %v", err))
		return
	}
	result.Success(c, "文件上传成功")
}

// DownloadPodFile 下载Pod容器中的文件
// @Summary 下载Pod容器中的文件
// @Description 下载Pod容器指定路径下的文件并返回二进制流
// @Tags K8s容器终端
// @Produce application/octet-stream
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName query string false "容器名称（默认为Pod中第一个容器）"
// @Param path query string true "目标文件路径"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/download [get]
func (ctrl *K8sTerminalController) DownloadPodFile(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.Query("containerName")
	path := c.Query("path")

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}
	if path == "" {
		result.Failed(c, 400, "文件路径不能为空")
		return
	}

	contentBytes, err := ctrl.service.DownloadPodFile(uint(clusterId), namespaceName, podName, containerName, path)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("下载文件失败: %v", err))
		return
	}

	// 从路径截取简单文件名
	fileName := "download_file"
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			if i < len(path)-1 {
				fileName = path[i+1:]
			}
			break
		}
	}
	if fileName == "download_file" && len(path) > 0 && path[0] != '/' {
		fileName = path
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "application/octet-stream", contentBytes)
}

// CreatePodDirectory 创建Pod容器中的目录
// @Summary 创建Pod容器中的目录
// @Description 通过执行mkdir -p命令在容器的指定路径下创建目录
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/directory [post]
func (ctrl *K8sTerminalController) CreatePodDirectory(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")

	var req struct {
		ContainerName string `json:"containerName"`
		Path          string `json:"path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数解析失败")
		return
	}

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}
	if req.Path == "" {
		result.Failed(c, 400, "目录路径不能为空")
		return
	}

	err = ctrl.service.CreatePodDirectory(uint(clusterId), namespaceName, podName, req.ContainerName, req.Path)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("创建目录失败: %v", err))
		return
	}
	result.Success(c, "目录创建成功")
}

// HotReloadPod 热加载Pod容器
// @Summary 热加载Pod容器
// @Description 类似于执行nginx -s reload效果，通常面向前台1号进程发送SIGHUP信号进行配置重载
// @Tags K8s容器终端
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/hot-reload [post]
func (ctrl *K8sTerminalController) HotReloadPod(c *gin.Context) {
	clusterId, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")

	var req struct {
		ContainerName string `json:"containerName"`
	}
	_ = c.ShouldBindJSON(&req)

	if namespaceName == "" || podName == "" {
		result.Failed(c, 400, "命名空间和Pod名称不能为空")
		return
	}

	err = ctrl.service.HotReloadPod(uint(clusterId), namespaceName, podName, req.ContainerName)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("热加载指令执行失败: %v", err))
		return
	}
	result.Success(c, "热加载指令已成功发送")
}
