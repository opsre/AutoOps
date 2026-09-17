package service

import (
	"errors"
	"fmt"
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/common/util"
	"dodevops-api/common/util/websocket"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

type CmdbSnmpDeviceSSHService struct {
	deviceDao *dao.CmdbSnmpDeviceDao
}

func NewCmdbSnmpDeviceSSHService() *CmdbSnmpDeviceSSHService {
	return &CmdbSnmpDeviceSSHService{
		deviceDao: dao.NewCmdbSnmpDeviceDao(),
	}
}

// ConnectTerminal 连接SNMP设备的SSH终端
func (s *CmdbSnmpDeviceSSHService) ConnectTerminal(c *gin.Context, deviceID uint) (*websocket.WebSSH, error) {
	// 获取设备信息
	device, err := s.deviceDao.GetSnmpDeviceById(deviceID)
	if err != nil {
		log.Printf("获取设备ID=%d信息失败: %v", deviceID, err)
		return nil, fmt.Errorf("获取设备信息失败: %v", err)
	}

	// 检查SSH配置是否完整
	if device.IPAddress == "" || device.RemoteUsername == "" || device.RemotePassword == "" {
		log.Printf("设备ID=%d SSH信息不完整", deviceID)
		return nil, errors.New("设备SSH信息不完整")
	}

	// 使用IP地址（不包含端口）
	sshHost := device.IPAddress

	// 设置SSH端口
	sshPort := device.RemotePort
	if sshPort == 0 {
		sshPort = 22 // 默认SSH端口
	}

	// 创建密码认证方式
	authMethod := ssh.Password(device.RemotePassword)

	// 创建WebSSH连接
	webSSH, err := websocket.NewWebSSHWithAuth(
		sshHost,
		strconv.Itoa(sshPort),
		device.RemoteUsername,
		authMethod,
	)
	if err != nil {
		log.Printf("创建WebSSH连接失败: %v", err)
		return nil, fmt.Errorf("创建WebSSH连接失败: %v", err)
	}

	return webSSH, nil
}

// ExecuteCommand 在SNMP设备上执行SSH命令
func (s *CmdbSnmpDeviceSSHService) ExecuteCommand(c *gin.Context, deviceID uint, command string) (string, error) {
	// 获取设备信息
	device, err := s.deviceDao.GetSnmpDeviceById(deviceID)
	if err != nil {
		return "", fmt.Errorf("获取设备信息失败: %v", err)
	}

	// 检查SSH配置
	if device.IPAddress == "" || device.RemoteUsername == "" || device.RemotePassword == "" {
		return "", errors.New("设备SSH信息不完整")
	}

	// 使用IP地址
	sshHost := device.IPAddress

	// 设置SSH端口
	sshPort := device.RemotePort
	if sshPort == 0 {
		sshPort = 22
	}

	// 创建SSH连接
	sshConfig := &util.SSHConfig{
		IP:       sshHost,
		Port:     sshPort,
		Username: device.RemoteUsername,
		Password: device.RemotePassword,
		Type:     1, // 密码认证
	}

	client, err := util.GetSSHClientByConfig(sshConfig)
	if err != nil {
		return "", fmt.Errorf("创建SSH连接失败: %v", err)
	}
	defer client.Close()

	// 创建会话并执行命令
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 执行命令
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("命令执行失败: %v", err)
	}

	return string(output), nil
}
