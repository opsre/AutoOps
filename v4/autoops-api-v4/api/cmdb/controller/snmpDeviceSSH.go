package controller

import (
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CmdbSnmpDeviceSSHController struct {
	sshService *service.CmdbSnmpDeviceSSHService
}

func NewCmdbSnmpDeviceSSHController() *CmdbSnmpDeviceSSHController {
	return &CmdbSnmpDeviceSSHController{
		sshService: service.NewCmdbSnmpDeviceSSHService(),
	}
}

// ConnectTerminal 连接SNMP设备SSH终端
// @Summary 连接SNMP设备SSH终端
// @Description 连接SNMP设备SSH终端
// @Tags CMDB-SNMP设备SSH
// @Param id path uint true "设备ID"
// @Param token query string false "认证Token"
// @Success 101 "Switching Protocols"
// @Router /api/v1/cmdb/snmp-devices/{id}/ssh/connect [get]
// @Security ApiKeyAuth
func (ctrl *CmdbSnmpDeviceSSHController) ConnectTerminal(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	log.Printf("开始处理SNMP设备WebSocket连接请求, deviceID: %s", deviceID)

	// 解析设备ID
	id, err := strconv.ParseUint(deviceID, 10, 32)
	if err != nil {
		log.Printf("无效的设备ID: %s, 错误: %v", deviceID, err)
		ctx.String(http.StatusBadRequest, "无效的设备ID")
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
		// 如果URL中没有token，检查Authorization头
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

	_, err = jwt.ValidateToken(token)
	if err != nil {
		log.Printf("令牌验证失败: %v", err)
		ctx.String(http.StatusUnauthorized, "令牌验证失败")
		return
	}

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
	log.Println("尝试连接SNMP设备SSH终端")
	webSSH, err := ctrl.sshService.ConnectTerminal(ctx, uint(id))
	if err != nil {
		log.Printf("SSH连接失败: %v", err)
		wsConn.Close()
		ctx.String(http.StatusInternalServerError, "SSH连接失败: "+err.Error())
		return
	}

	// 连接WebSocket
	if err := webSSH.Connect(wsConn); err != nil {
		log.Printf("WebSSH连接失败: %v", err)
		webSSH.Close()
		wsConn.Close()
		ctx.String(http.StatusInternalServerError, "WebSSH连接失败")
		return
	}

	// SSH连接测试通过
	log.Println("SNMP设备SSH连接测试通过")
	log.Printf("SNMP设备SSH连接测试通过, WebSocket和SSH连接已建立, 远程地址: %s", wsConn.RemoteAddr().String())

	// 设置defer确保资源释放
	defer func() {
		webSSH.Close()
		wsConn.Close()
	}()

	// 保持连接
	select {}
}

// ExecuteCommand 在SNMP设备上执行SSH命令
// @Summary 在SNMP设备上执行SSH命令
// @Description 在SNMP设备上执行SSH命令
// @Tags CMDB-SNMP设备SSH
// @Param id path uint true "设备ID"
// @Param command query string true "命令"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/snmp-devices/{id}/ssh/command [get]
// @Security ApiKeyAuth
func (ctrl *CmdbSnmpDeviceSSHController) ExecuteCommand(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	command := ctx.Query("command")

	id, err := strconv.ParseUint(deviceID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的设备ID")
		return
	}

	output, err := ctrl.sshService.ExecuteCommand(ctx, uint(id), command)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "执行命令失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"deviceId": deviceID,
		"command":  command,
		"output":   output,
	})
}
