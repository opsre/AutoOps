// 图片上传 服务层
// author xiaoRui

package service

import (
	"fmt"
	"dodevops-api/common/config"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/log"
	"github.com/gin-gonic/gin"
	"path"
	"strconv"
	"time"
)

type IUploadService interface {
	Upload(c *gin.Context)
}

type UploadServiceImpl struct{}

// 图片上传
func (u UploadServiceImpl) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		result.Failed(c, int(result.ApiCode.FILEUPLOADERROR), result.ApiCode.GetMessage(result.ApiCode.FILEUPLOADERROR))
		return
	}

	now := time.Now()
	ext := path.Ext(file.Filename)
	fileName := strconv.Itoa(now.Nanosecond()) + ext

	// 文件保存路径（增加 avatar 一级目录）: upload/avatar/20060102/xxx.ext
	filePath := fmt.Sprintf("%savatar/%s/%s",
		config.Config.ImageSettings.UploadDir,
		now.Format("20060102"),
		fileName)

	// 创建目录（包含 avatar 子目录）
	dir := fmt.Sprintf("%savatar/%s",
		config.Config.ImageSettings.UploadDir,
		now.Format("20060102"))

	if err := util.CreateDir(dir); err != nil {
		result.Failed(c, int(result.ApiCode.FILEUPLOADERROR), "创建目录失败: "+err.Error())
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		result.Failed(c, int(result.ApiCode.FILEUPLOADERROR), "保存文件失败: "+err.Error())
		return
	}

	// 构造浏览器可访问的 URL（增加 avatar 路径）
	relativeURL := "/api/v1/upload/avatar/" + now.Format("20060102") + "/" + fileName

	// 动态获取当前请求的协议和Host
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	// 优先使用 X-Forwarded-Host，否则使用 Host header
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	// 注意：X-Forwarded-Host 可能已经包含端口，不需要再添加

	log.Log().Infof("Upload - Request Headers: Host=%s, X-Forwarded-Host=%s, X-Forwarded-Port=%s, X-Forwarded-Proto=%s",
		c.Request.Host, c.GetHeader("X-Forwarded-Host"), c.GetHeader("X-Forwarded-Port"), c.GetHeader("X-Forwarded-Proto"))

	imageHost := scheme + "://" + host
	finalURL := imageHost + relativeURL
	log.Log().Infof("Upload - Final URL: %s", finalURL)

	// 返回结果
	result.Success(c, finalURL)
}

var uploadService = UploadServiceImpl{}

func UploadService() IUploadService {
	return &uploadService
}
