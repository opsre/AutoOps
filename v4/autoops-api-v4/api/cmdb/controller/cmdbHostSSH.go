package controller

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util/recording"
	wsUtil "dodevops-api/common/util/websocket"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type noEarlyResponseWriter struct {
	gin.ResponseWriter
	headerWritten bool
}

type hijackedResponseWriter struct {
	conn net.Conn
}

func (w *noEarlyResponseWriter) WriteHeader(code int) {
	if !w.headerWritten {
		w.headerWritten = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *noEarlyResponseWriter) WriteHeaderNow() {
	if !w.headerWritten {
		w.headerWritten = true
		w.ResponseWriter.WriteHeaderNow()
	}
}

func (w *hijackedResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *hijackedResponseWriter) Write(data []byte) (int, error) {
	return w.conn.Write(data)
}

func (w *hijackedResponseWriter) WriteHeader(statusCode int) {
	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "status code " + strconv.Itoa(statusCode)
	}
	w.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)))
}

type CmdbHostSSHController struct {
	hostSSHService   service.CmdbHostSSHServiceInterface
	recordingHelper  *RecordingHelper
	recordingEnabled bool
}

func NewCmdbHostSSHController(hostSSHService service.CmdbHostSSHServiceInterface) *CmdbHostSSHController {
	// 初始化录制功能(默认存储路径)
	recordingHelper, err := NewRecordingHelper("./upload")
	if err != nil {
		log.Printf("初始化录制功能失败: %v (录制功能将被禁用)", err)
		return &CmdbHostSSHController{
			hostSSHService:   hostSSHService,
			recordingEnabled: false,
		}
	}

	return &CmdbHostSSHController{
		hostSSHService:   hostSSHService,
		recordingHelper:  recordingHelper,
		recordingEnabled: true,
	}
}

// ConnectTerminal 连接SSH终端
func (c *CmdbHostSSHController) ConnectTerminal(ctx *gin.Context) {
	hostID := ctx.Param("id")
	log.Printf("开始处理WebSocket连接请求, hostID: %s", hostID)

	// 解析主机ID
	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		log.Printf("无效的主机ID: %s, 错误: %v", hostID, err)
		ctx.String(http.StatusBadRequest, "无效的主机ID")
		return
	}

	// 使用标准WebSocket升级
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 优先从URL参数获取token
	token := ctx.Query("token")
	if token == "" {
		// 如果URL中没有token,检查Authorization头
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("缺少token或Authorization头")
			ctx.String(http.StatusUnauthorized, "缺少token或Authorization头")
			return
		}

		// 验证Bearer令牌格式
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			log.Println("无效的Authorization头格式")
			ctx.String(http.StatusUnauthorized, "无效的Authorization头格式")
			return
		}
		token = authHeader[7:]
	}
	claims, err := jwt.ValidateToken(token)
	if err != nil {
		log.Printf("令牌验证失败: %v", err)
		ctx.String(http.StatusUnauthorized, "令牌验证失败")
		return
	}
	// WebSocket接口未经过常规JWT中间件,这里手动注入用户信息,供录制等逻辑使用
	ctx.Set(constant.ContextKeyUserObj, claims)

	if !websocket.IsWebSocketUpgrade(ctx.Request) {
		log.Printf("无效的WebSocket升级请求: method=%s, connection=%q, upgrade=%q, uri=%s",
			ctx.Request.Method,
			ctx.GetHeader("Connection"),
			ctx.GetHeader("Upgrade"),
			ctx.Request.RequestURI,
		)
		ctx.String(http.StatusBadRequest, "请求不是有效的WebSocket升级请求")
		return
	}

	// 升级WebSocket连接
	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 连接SSH终端
	log.Println("尝试连接SSH终端")
	webSSH, err := c.hostSSHService.ConnectTerminal(ctx, uint(id))
	if err != nil {
		log.Printf("SSH连接失败: %v", err)
		wsConn.Close()
		ctx.String(http.StatusInternalServerError, "SSH连接失败")
		return
	}

	// 包装为带录制的WebSSH
	var recordingWS *wsUtil.RecordingWebSSH
	if c.recordingEnabled {
		recordingWS, err = c.recordingHelper.WrapWithRecording(ctx, webSSH, uint(id))
		if err != nil {
			log.Printf("创建录制WebSSH失败: %v (继续使用普通WebSSH)", err)
			// 录制失败不影响正常使用,降级为普通WebSSH
			c.recordingEnabled = false
		}
	}

	// 连接WebSocket
	if c.recordingEnabled && recordingWS != nil {
		// 使用带录制的WebSSH
		if err := recordingWS.Connect(wsConn); err != nil {
			log.Printf("RecordingWebSSH连接失败: %v", err)
			var sessionID string
			var recorder recording.Recorder
			if recorder = recordingWS.GetRecorder(); recorder != nil {
				sessionID = recorder.GetSessionID()
			}
			recordingWS.Close()
			wsConn.Close()
			if sessionID != "" {
				c.recordingHelper.OnRecordingComplete(sessionID, nil, recorder, err)
			}
			ctx.String(http.StatusInternalServerError, "WebSSH连接失败")
			return
		}
		log.Println("WebSSH连接成功 (带录制)")

		// 等待连接结束
		recordingWS.Wait()
		log.Println("RecordingWebSSH连接已结束")

		// 停止录制并获取信息
		var info *recording.RecordingInfo
		var stopErr error
		var sessionID string
		var recorder recording.Recorder
		if recorder = recordingWS.GetRecorder(); recorder != nil {
			sessionID = recorder.GetSessionID()
		}
		info, stopErr = recordingWS.FinalizeRecording()
		if sessionID != "" && info == nil && stopErr == nil {
			stopErr = fmt.Errorf("recording stopped without info")
		}

		// 关闭连接
		recordingWS.Close()
		wsConn.Close()

		// 处理录制完成回调(更新数据库)
		if sessionID != "" {
			c.recordingHelper.OnRecordingComplete(sessionID, info, recorder, stopErr)
		}
	} else {
		// 使用普通WebSSH
		if err := webSSH.Connect(wsConn); err != nil {
			log.Printf("WebSSH连接失败: %v", err)
			webSSH.Close()
			wsConn.Close()
			ctx.String(http.StatusInternalServerError, "WebSSH连接失败")
			return
		}
		log.Println("WebSSH连接成功 (无录制)")

		// 等待SSH会话结束
		if session := webSSH.GetSession(); session != nil {
			err := session.Wait()
			if err != nil {
				log.Printf("SSH会话结束: %v", err)
			} else {
				log.Println("SSH会话正常结束")
			}
		}
		webSSH.Close()
		wsConn.Close()
	}
}

// ExecuteCommand 执行命令
// @Summary 执行SSH命令
// @Description 在SSH终端执行命令
// @Tags CMDB主机SSH
// @Accept  json
// @Produce  json
// @Param id path uint true "主机ID"
// @Param command query string true "命令"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostssh/command/{id} [get]
// @Security ApiKeyAuth
func (c *CmdbHostSSHController) ExecuteCommand(ctx *gin.Context) {
	hostID := ctx.Param("id")
	command := ctx.Query("command")

	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	output, err := c.hostSSHService.ExecuteCommand(ctx, uint(id), command)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "执行命令失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"hostId":  hostID,
		"command": command,
		"output":  output,
	})
}

// UploadFile 上传文件到SSH服务器
// @Summary 上传文件到SSH服务器
// @Description 上传本地文件到远程SSH服务器
// @Tags CMDB主机SSH
// @Accept multipart/form-data
// @Produce json
// @Param id path uint true "主机ID"
// @Param file formData file true "要上传的文件"
// @Param destPath formData string true "远程服务器目标路径"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostssh/upload/{id} [post]
// @Security ApiKeyAuth
func (c *CmdbHostSSHController) UploadFile(ctx *gin.Context) {
	hostID := ctx.Param("id")
	file, err := ctx.FormFile("file")
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "获取上传文件失败: "+err.Error())
		return
	}

	destPath := ctx.PostForm("destPath")
	if destPath == "" {
		result.Failed(ctx, http.StatusBadRequest, "目标路径不能为空")
		return
	}

	// 保存临时文件（使用绝对路径并处理中文文件名）
	tempDir := "/tmp/ssh_uploads"
	log.Printf("创建临时目录: %s", tempDir)

	// 确保目录存在且有正确权限
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.Printf("创建临时目录失败: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "创建临时目录失败: "+err.Error())
		return
	}

	// 检查目录权限
	if fi, err := os.Stat(tempDir); err == nil {
		log.Printf("临时目录权限: %v", fi.Mode())
	} else {
		log.Printf("无法获取临时目录状态: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "无法访问临时目录: "+err.Error())
		return
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	log.Printf("尝试保存文件到: %s", tempFilePath)

	// 检查文件是否已存在
	if _, err := os.Stat(tempFilePath); err == nil {
		log.Printf("文件已存在，将被覆盖: %s", tempFilePath)
		if err := os.Remove(tempFilePath); err != nil {
			log.Printf("删除已存在文件失败: %v", err)
			result.Failed(ctx, http.StatusInternalServerError, "删除已存在文件失败: "+err.Error())
			return
		}
	}

	// 保存文件并设置权限
	log.Printf("保存文件: %s (大小: %d bytes)", file.Filename, file.Size)
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		log.Printf("保存文件失败: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "保存临时文件失败: "+err.Error())
		return
	}

	// 设置文件权限为0644
	if err := os.Chmod(tempFilePath, 0644); err != nil {
		log.Printf("设置文件权限失败: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "设置文件权限失败: "+err.Error())
		return
	}

	// 严格验证文件是否保存成功
	if fi, err := os.Stat(tempFilePath); err != nil {
		log.Printf("文件保存验证失败: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "文件保存验证失败: "+err.Error())
		return
	} else {
		log.Printf("文件保存成功: %s (大小: %d bytes, 权限: %v)",
			tempFilePath, fi.Size(), fi.Mode())
	}

	// 强制同步文件系统
	if file, err := os.Open(tempFilePath); err == nil {
		file.Sync()
		file.Close()
		log.Printf("已同步文件系统: %s", tempFilePath)
	} else {
		log.Printf("无法打开文件进行同步: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "无法验证文件状态: "+err.Error())
		return
	}

	// 再次验证文件存在
	if _, err := os.Stat(tempFilePath); err != nil {
		log.Printf("文件在同步后丢失: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "文件同步后丢失: "+err.Error())
		return
	}

	// 转换为uint
	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	// 验证文件在调用服务前仍然存在
	if _, err := os.Stat(tempFilePath); err != nil {
		log.Printf("文件在调用服务前已不存在: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "临时文件已丢失: "+err.Error())
		return
	}

	// 创建后台上传任务。浏览器到后端的 multipart 上传到这里已经完成，
	// 后续后端到远程主机的传输改为后台执行，前端通过任务状态接口轮询进度。
	log.Printf("开始创建SSH上传后台任务: %s -> %s", tempFilePath, destPath)
	task, err := c.hostSSHService.CreateUploadTask(ctx, uint(id), tempFilePath, file.Filename, destPath)
	if err != nil {
		log.Printf("创建上传任务失败: %v", err)
		if removeErr := os.Remove(tempFilePath); removeErr != nil && !os.IsNotExist(removeErr) {
			log.Printf("清理失败的上传临时文件失败: %v", removeErr)
		}
		result.Failed(ctx, http.StatusInternalServerError, "创建上传任务失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"hostId":   hostID,
		"fileName": file.Filename,
		"destPath": destPath,
		"task":     task,
		"message":  "文件已接收，后台上传任务已启动",
	})
}

// GetUploadTaskStatus 查询SSH上传任务状态
// @Summary 查询SSH上传任务状态
// @Description 查询后台SSH文件上传任务的实时进度和结果
// @Tags CMDB主机SSH
// @Accept json
// @Produce json
// @Param taskId path string true "上传任务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostssh/upload/task/{taskId} [get]
// @Security ApiKeyAuth
func (c *CmdbHostSSHController) GetUploadTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	if taskID == "" {
		result.Failed(ctx, http.StatusBadRequest, "任务ID不能为空")
		return
	}

	task, err := c.hostSSHService.GetUploadTask(taskID)
	if err != nil {
		result.Failed(ctx, http.StatusNotFound, err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"task": task,
	})
}

// DownloadFile 从SSH服务器下载文件
// @Summary 从SSH服务器下载文件
// @Description 从远程SSH服务器流式下载文件
// @Tags CMDB主机SSH
// @Accept json
// @Produce octet-stream
// @Param id path uint true "主机ID"
// @Param remotePath query string true "远程文件路径"
// @Success 200 {file} file "下载的文件"
// @Router /api/v1/cmdb/hostssh/download/{id} [get]
// @Security ApiKeyAuth
func (c *CmdbHostSSHController) DownloadFile(ctx *gin.Context) {
	hostID := ctx.Param("id")
	remotePath := ctx.Query("remotePath")

	if remotePath == "" {
		result.Failed(ctx, http.StatusBadRequest, "远程文件路径不能为空")
		return
	}

	// 转换主机ID
	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	log.Printf("开始流式下载文件: hostID=%d, remotePath=%s", id, remotePath)

	// 直接调用服务层的流式下载（不再需要本地临时文件）
	// Service 层会直接将数据流式写入 HTTP 响应
	err = c.hostSSHService.DownloadFile(ctx, uint(id), remotePath, "")
	if err != nil {
		log.Printf("文件下载失败: %v", err)
		// 注意：如果已经开始写入响应，这里的错误处理可能不会生效
		// 因为HTTP头已经发送
		result.Failed(ctx, http.StatusInternalServerError, "文件下载失败: "+err.Error())
		return
	}

	log.Printf("文件流式下载完成: %s", remotePath)
}

// ListFiles 列出SSH服务器目录下的文件
// @Summary 列出SSH服务器目录下的文件
// @Description 列出远程SSH服务器指定目录下的所有文件和文件夹
// @Tags CMDB主机SSH
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Param remotePath query string false "远程目录路径(默认为根目录/)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostssh/listfiles/{id} [get]
// @Security ApiKeyAuth
func (c *CmdbHostSSHController) ListFiles(ctx *gin.Context) {
	hostID := ctx.Param("id")
	remotePath := ctx.Query("remotePath")

	// 如果未指定路径，默认为根目录
	if remotePath == "" {
		remotePath = "/"
	}

	// 转换主机ID
	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	// 调用服务列出文件
	fileInfos, err := c.hostSSHService.ListFiles(ctx, uint(id), remotePath)
	if err != nil {
		log.Printf("列出文件失败: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "列出文件失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"hostId":     hostID,
		"remotePath": remotePath,
		"files":      fileInfos,
		"total":      len(fileInfos),
	})
}
