package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

func parsePodFileParams(c *gin.Context) (clusterID uint, namespaceName, podName string, ok bool) {
	clusterIDValue, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterIDValue <= 0 {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return 0, "", "", false
	}

	namespaceName = c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, http.StatusBadRequest, "命名空间名称不能为空")
		return 0, "", "", false
	}

	podName = c.Param("podName")
	if podName == "" {
		result.Failed(c, http.StatusBadRequest, "Pod名称不能为空")
		return 0, "", "", false
	}

	return uint(clusterIDValue), namespaceName, podName, true
}

// ListPodFiles 列出Pod容器目录文件
// @Summary 列出Pod容器目录文件
// @Description 列出指定Pod容器目录下的文件和文件夹，供前端浏览并选择下载
// @Tags Pod文件管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param remotePath query string false "容器内目录路径，默认 /"
// @Param containerName query string false "容器名称，不传默认第一个业务容器"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files [get]
func (ctrl *K8sWorkloadController) ListPodFiles(c *gin.Context) {
	clusterID, namespaceName, podName, ok := parsePodFileParams(c)
	if !ok {
		return
	}

	remotePath := c.Query("remotePath")
	if remotePath == "" {
		remotePath = "/"
	}

	containerName := c.Query("containerName")
	if containerName == "" {
		containerName = c.Query("container")
	}

	ctrl.service.ListPodFiles(c, clusterID, namespaceName, podName, containerName, remotePath)
}

// UploadPodFile 上传文件到Pod容器
// @Summary 上传文件到Pod容器
// @Description 通过表单文件上传，将文件写入指定Pod容器目录
// @Tags Pod文件管理
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param file formData file true "上传文件"
// @Param destPath formData string true "容器内目标目录，必须为绝对路径"
// @Param containerName formData string false "容器名称，不传默认第一个业务容器"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/upload [post]
func (ctrl *K8sWorkloadController) UploadPodFile(c *gin.Context) {
	clusterID, namespaceName, podName, ok := parsePodFileParams(c)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "获取上传文件失败: "+err.Error())
		return
	}

	destPath := c.PostForm("destPath")
	if destPath == "" {
		destPath = c.Query("destPath")
	}
	if destPath == "" {
		result.Failed(c, http.StatusBadRequest, "目标路径不能为空")
		return
	}

	containerName := c.PostForm("containerName")
	if containerName == "" {
		containerName = c.PostForm("container")
	}
	if containerName == "" {
		containerName = c.Query("containerName")
	}
	if containerName == "" {
		containerName = c.Query("container")
	}

	file, err := fileHeader.Open()
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "打开上传文件失败: "+err.Error())
		return
	}
	defer file.Close()

	ctrl.service.UploadPodFile(c, clusterID, namespaceName, podName, containerName, destPath, fileHeader.Filename, file)
}

// DownloadPodFile 从Pod容器下载文件
// @Summary 从Pod容器下载文件
// @Description 流式下载指定Pod容器中的文件
// @Tags Pod文件管理
// @Accept json
// @Produce octet-stream
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param remotePath query string true "容器内文件绝对路径"
// @Param containerName query string false "容器名称，不传默认第一个业务容器"
// @Success 200 {file} file "下载的文件"
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/files/download [get]
func (ctrl *K8sWorkloadController) DownloadPodFile(c *gin.Context) {
	clusterID, namespaceName, podName, ok := parsePodFileParams(c)
	if !ok {
		return
	}

	remotePath := c.Query("remotePath")
	if remotePath == "" {
		result.Failed(c, http.StatusBadRequest, "远程文件路径不能为空")
		return
	}

	containerName := c.Query("containerName")
	if containerName == "" {
		containerName = c.Query("container")
	}

	ctrl.service.DownloadPodFile(c, clusterID, namespaceName, podName, containerName, remotePath)
}
