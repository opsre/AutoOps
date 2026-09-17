package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/config"
	"dodevops-api/common/winrm"
	"fmt"
	"log"
	"strings"
)

type WindowsExporterService struct {
	snmpDeviceDao *dao.CmdbSnmpDeviceDao
}

func NewWindowsExporterService() *WindowsExporterService {
	return &WindowsExporterService{
		snmpDeviceDao: dao.NewCmdbSnmpDeviceDao(),
	}
}

// InstallExporterRequest 安装请求
type InstallExporterRequest struct {
	DeviceID     uint   `json:"deviceId" validate:"required"`
	InstallerURL string `json:"installerUrl"` // MSI下载地址（可选，默认从配置读取）
	Version      string `json:"version"`      // 版本号（可选，默认从配置读取）
	ExporterPort int    `json:"exporterPort"` // 端口（可选，默认从配置读取）
}

// ExporterStatusResponse 状态响应
type ExporterStatusResponse struct {
	DeviceID    uint   `json:"deviceId"`
	Installed   bool   `json:"installed"`
	Running     bool   `json:"running"`
	Version     string `json:"version"`
	Port        int    `json:"port"`
	Status      int    `json:"status"`           // 状态码
	StatusText  string `json:"statusText"`       // 状态描述
	StatusColor string `json:"statusColor"`      // 前端显示颜色
}

// InstallWindowsExporter 安装Windows Exporter
func (s *WindowsExporterService) InstallWindowsExporter(req *InstallExporterRequest) error {
	log.Printf("======== 开始安装 Windows Exporter ========")
	log.Printf("设备ID: %d", req.DeviceID)

	// 1. 查询设备信息
	log.Printf("[1/9] 查询设备信息...")
	device, err := s.snmpDeviceDao.GetSnmpDeviceById(req.DeviceID)
	if err != nil {
		log.Printf("❌ 设备不存在: %v", err)
		return fmt.Errorf("设备不存在: %v", err)
	}
	log.Printf("✅ 设备信息: %s (%s)", device.DeviceName, device.IPAddress)

	// 2. 验证设备类型
	log.Printf("[2/9] 验证设备类型...")
	if device.DeviceType != "windows" {
		log.Printf("❌ 设备类型错误: %s", device.DeviceType)
		return fmt.Errorf("设备类型必须为 windows，当前为 %s", device.DeviceType)
	}
	log.Printf("✅ 设备类型正确: windows")

	// 3. 验证凭据
	log.Printf("[3/9] 验证远程凭据...")
	if device.RemoteUsername == "" || device.RemotePassword == "" {
		log.Printf("❌ 未配置远程凭据")
		return fmt.Errorf("设备未配置远程凭据，无法安装")
	}
	log.Printf("✅ 凭据验证通过: %s", device.RemoteUsername)

	// 4. 从配置文件读取默认值
	log.Printf("[4/9] 读取配置...")
	cfg := config.Config.WindowsExporter

	version := req.Version
	if version == "" {
		version = cfg.DefaultVersion
	}

	exporterPort := req.ExporterPort
	if exporterPort == 0 {
		exporterPort = cfg.DefaultPort
	}

	installerURL := req.InstallerURL
	if installerURL == "" {
		// 使用本地上传的MSI文件
		installerURL = fmt.Sprintf("%swindows_exporter-%s-amd64.msi", cfg.InstallerBaseURL, version)
	}
	log.Printf("✅ 安装配置: 版本=%s, 端口=%d", version, exporterPort)
	log.Printf("   安装包地址: %s", installerURL)

	// 5. 更新状态为安装中
	log.Printf("[5/9] 更新设备状态为\"安装中\"...")
	device.ExporterInstalled = 2 // 安装中
	if err := s.snmpDeviceDao.UpdateSnmpDevice(device.ID, device); err != nil {
		log.Printf("❌ 更新设备状态失败: %v", err)
		return fmt.Errorf("更新设备状态失败: %v", err)
	}
	log.Printf("✅ 状态已更新: exporter_installed=2 (安装中)")

	// 6. 创建WinRM管理器
	log.Printf("[6/9] 创建 WinRM 连接...")
	winrmMgr := winrm.NewManager(winrm.Config{
		Host:     device.IPAddress,
		Port:     5985, // HTTP WinRM端口
		Username: device.RemoteUsername,
		Password: device.RemotePassword,
		UseHTTPS: false,
	})
	log.Printf("✅ WinRM 管理器已创建: %s:5985", device.IPAddress)

	// 7. 测试连接
	log.Printf("[7/9] 测试 WinRM 连接...")
	if err := winrmMgr.TestConnection(); err != nil {
		log.Printf("❌ WinRM连接失败: %v", err)
		device.ExporterInstalled = 4 // 安装失败
		s.snmpDeviceDao.UpdateSnmpDevice(device.ID, device)
		return fmt.Errorf("WinRM连接失败: %v，请确保目标机器已启用WinRM", err)
	}
	log.Printf("✅ WinRM 连接成功")

	// 8. 执行安装（传递Pushgateway配置）
	log.Printf("[8/9] 开始安装 Windows Exporter...")
	var pushgatewayURL string
	var pushInterval int
	if cfg.PushMode.Enabled {
		pushgatewayURL = cfg.PushMode.PushgatewayURL
		pushInterval = cfg.PushMode.Interval
		log.Printf("   推送模式: 启用 (Pushgateway=%s, 间隔=%d秒)", pushgatewayURL, pushInterval)
	} else {
		log.Printf("   推送模式: 禁用 (被动采集)")
	}

	if err := winrmMgr.InstallWindowsExporter(installerURL, version, exporterPort, pushgatewayURL, pushInterval, device.DeviceName); err != nil {
		log.Printf("❌ 安装失败: %v", err)
		// 区分是安装失败还是启动失败
		if strings.Contains(err.Error(), "服务启动失败") || strings.Contains(err.Error(), "Start-Service") {
			device.ExporterInstalled = 3 // 已安装未运行
		} else {
			device.ExporterInstalled = 4 // 安装失败
		}
		s.snmpDeviceDao.UpdateSnmpDevice(device.ID, device)
		return fmt.Errorf("安装失败: %v", err)
	}
	log.Printf("✅ Windows Exporter 安装成功")

	// 9. 更新设备状态为已安装运行中
	log.Printf("[9/9] 更新设备状态为\"已安装运行中\"...")
	device.ExporterInstalled = 1 // 已安装运行中
	device.ExporterPort = exporterPort
	device.ExporterVersion = version
	if err := s.snmpDeviceDao.UpdateSnmpDevice(device.ID, device); err != nil {
		log.Printf("❌ 更新设备状态失败: %v", err)
		return fmt.Errorf("更新设备状态失败: %v", err)
	}
	log.Printf("✅ 状态已更新: exporter_installed=1 (已安装)")
	log.Printf("======== 安装完成 ========\n")

	return nil
}

// GetExporterStatus 获取Exporter状态
func (s *WindowsExporterService) GetExporterStatus(deviceID uint) (*ExporterStatusResponse, error) {
	// 1. 查询设备信息
	device, err := s.snmpDeviceDao.GetSnmpDeviceById(deviceID)
	if err != nil {
		return nil, fmt.Errorf("设备不存在: %v", err)
	}

	// 2. 验证设备类型
	if device.DeviceType != "windows" {
		return nil, fmt.Errorf("设备类型必须为 windows")
	}

	// 3. 根据数据库状态码构建响应
	response := &ExporterStatusResponse{
		DeviceID: device.ID,
		Version:  device.ExporterVersion,
		Port:     device.ExporterPort,
		Status:   device.ExporterInstalled,
	}

	// 4. 设置状态文本和颜色
	switch device.ExporterInstalled {
	case 0:
		response.Installed = false
		response.Running = false
		response.StatusText = "未安装"
		response.StatusColor = "default"
	case 1:
		response.Installed = true
		response.Running = true
		response.StatusText = "运行中"
		response.StatusColor = "success"
	case 2:
		response.Installed = false
		response.Running = false
		response.StatusText = "安装中"
		response.StatusColor = "processing"
	case 3:
		response.Installed = true
		response.Running = false
		response.StatusText = "已安装未运行"
		response.StatusColor = "warning"
	case 4:
		response.Installed = false
		response.Running = false
		response.StatusText = "安装失败"
		response.StatusColor = "error"
	case 5:
		response.Installed = true
		response.Running = false
		response.StatusText = "启动中"
		response.StatusColor = "processing"
	default:
		response.Installed = false
		response.Running = false
		response.StatusText = "未知状态"
		response.StatusColor = "default"
	}

	// 5. 如果数据库显示已安装，实时查询运行状态验证
	if (device.ExporterInstalled == 1 || device.ExporterInstalled == 3) && device.RemoteUsername != "" && device.RemotePassword != "" {
		winrmMgr := winrm.NewManager(winrm.Config{
			Host:     device.IPAddress,
			Port:     5985,
			Username: device.RemoteUsername,
			Password: device.RemotePassword,
			UseHTTPS: false,
		})

		isRunning, version, err := winrmMgr.CheckExporterStatus()
		if err == nil {
			response.Running = isRunning
			if version != "" {
				response.Version = version
			}
			// 实时更新状态
			if isRunning && device.ExporterInstalled != 1 {
				// 数据库状态不对，更新为运行中
				response.Status = 1
				response.StatusText = "运行中"
				response.StatusColor = "success"
			} else if !isRunning && device.ExporterInstalled == 1 {
				// 服务停止了，更新为已安装未运行
				response.Status = 3
				response.StatusText = "已安装未运行"
				response.StatusColor = "warning"
			}
		}
	}

	return response, nil
}

// UninstallExporter 卸载Exporter
func (s *WindowsExporterService) UninstallExporter(deviceID uint) error {
	// 1. 查询设备信息
	device, err := s.snmpDeviceDao.GetSnmpDeviceById(deviceID)
	if err != nil {
		return fmt.Errorf("设备不存在: %v", err)
	}

	// 2. 验证设备类型
	if device.DeviceType != "windows" {
		return fmt.Errorf("设备类型必须为 windows")
	}

	// 3. 验证凭据
	if device.RemoteUsername == "" || device.RemotePassword == "" {
		return fmt.Errorf("设备未配置远程凭据，无法卸载")
	}

	// 4. 创建WinRM管理器
	winrmMgr := winrm.NewManager(winrm.Config{
		Host:     device.IPAddress,
		Port:     5985,
		Username: device.RemoteUsername,
		Password: device.RemotePassword,
		UseHTTPS: false,
	})

	// 5. 执行卸载
	if err := winrmMgr.UninstallExporter(); err != nil {
		return fmt.Errorf("卸载失败: %v", err)
	}

	// 6. 更新设备状态
	device.ExporterInstalled = 0
	device.ExporterPort = 0
	device.ExporterVersion = ""
	if err := s.snmpDeviceDao.UpdateSnmpDevice(device.ID, device); err != nil {
		return fmt.Errorf("更新设备状态失败: %v", err)
	}

	return nil
}

// BatchInstallExporter 批量安装Exporter
func (s *WindowsExporterService) BatchInstallExporter(req *InstallExporterRequest, deviceIDs []uint) map[uint]error {
	results := make(map[uint]error)

	for _, deviceID := range deviceIDs {
		req.DeviceID = deviceID
		err := s.InstallWindowsExporter(req)
		results[deviceID] = err
	}

	return results
}

// GetAllWindowsDevicesStatus 获取所有Windows设备的Exporter状态
func (s *WindowsExporterService) GetAllWindowsDevicesStatus() ([]*ExporterStatusResponse, error) {
	// 查询所有Windows设备
	query := &model.SnmpDeviceQueryDto{
		DeviceType: "windows",
		Page:       1,
		PageSize:   1000,
	}

	devices, total, err := s.snmpDeviceDao.GetSnmpDeviceList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]*ExporterStatusResponse, 0, total)
	for _, device := range devices {
		status, err := s.GetExporterStatus(device.ID)
		if err == nil {
			responses = append(responses, status)
		}
	}

	return responses, nil
}
