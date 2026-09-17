package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/util"
	"errors"
	"time"
)

type CmdbSnmpDeviceService struct {
	deviceDao *dao.CmdbSnmpDeviceDao
}

func NewCmdbSnmpDeviceService() *CmdbSnmpDeviceService {
	return &CmdbSnmpDeviceService{
		deviceDao: dao.NewCmdbSnmpDeviceDao(),
	}
}

// CreateSnmpDevice 创建SNMP设备(简化版)
func (s *CmdbSnmpDeviceService) CreateSnmpDevice(dto *model.CreateSnmpDeviceDto, username string) (*model.CmdbSnmpDevice, error) {
	// 检查IP是否已存在
	exists, err := s.deviceDao.CheckIPExists(dto.IPAddress, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("IP地址已存在")
	}

	// 设置默认值
	snmpVersion := dto.SnmpVersion
	if snmpVersion == "" {
		snmpVersion = "v2c"
	}

	snmpCommunity := dto.SnmpCommunity
	if snmpCommunity == "" {
		snmpCommunity = "public"
	}

	// 构建设备对象(仅包含简化字段)
	device := &model.CmdbSnmpDevice{
		DeviceName:    dto.DeviceName,
		DeviceType:    dto.DeviceType,
		IPAddress:     dto.IPAddress,
		SnmpVersion:   snmpVersion,
		SnmpCommunity: snmpCommunity,
		MonitorEnabled: 1, // 默认启用监控
		Location:      dto.Location,
		Owner:         dto.Owner,
		Remark:        dto.Remark,
		Status:        "offline", // 默认离线状态
		CreateTime:    util.HTime{Time: time.Now()},
		UpdateTime:    util.HTime{Time: time.Now()},
	}

	// 创建设备
	err = s.deviceDao.CreateSnmpDevice(device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

// UpdateSnmpDevice 更新SNMP设备(简化版)
func (s *CmdbSnmpDeviceService) UpdateSnmpDevice(dto *model.UpdateSnmpDeviceDto, username string) error {
	// 检查设备是否存在
	_, err := s.deviceDao.GetSnmpDeviceById(dto.ID)
	if err != nil {
		return errors.New("设备不存在")
	}

	// 如果修改了IP地址,检查新IP是否已被其他设备使用
	if dto.IPAddress != "" {
		exists, err := s.deviceDao.CheckIPExists(dto.IPAddress, dto.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("IP地址已被其他设备使用")
		}
	}

	// 构建更新数据映射 - 使用map以支持更新为空值
	updateFields := map[string]interface{}{
		"update_time": time.Now(),
	}

	// 基础字段(非空才更新)
	if dto.DeviceName != "" {
		updateFields["device_name"] = dto.DeviceName
	}
	if dto.DeviceType != "" {
		updateFields["device_type"] = dto.DeviceType
	}
	if dto.IPAddress != "" {
		updateFields["ip_address"] = dto.IPAddress
	}
	if dto.SnmpVersion != "" {
		updateFields["snmp_version"] = dto.SnmpVersion
	}
	if dto.SnmpCommunity != "" {
		updateFields["snmp_community"] = dto.SnmpCommunity
	}
	if dto.MonitorEnabled >= 0 {
		updateFields["monitor_enabled"] = dto.MonitorEnabled
	}
	if dto.Location != "" {
		updateFields["location"] = dto.Location
	}
	if dto.Owner != "" {
		updateFields["owner"] = dto.Owner
	}

	// 远程连接字段(允许更新为空值或0)
	updateFields["remote_port"] = dto.RemotePort
	updateFields["remote_username"] = dto.RemoteUsername
	updateFields["remote_password"] = dto.RemotePassword
	updateFields["remote_domain"] = dto.RemoteDomain
	updateFields["remark"] = dto.Remark

	return s.deviceDao.UpdateSnmpDeviceFields(dto.ID, updateFields)
}

// DeleteSnmpDevice 删除SNMP设备
func (s *CmdbSnmpDeviceService) DeleteSnmpDevice(id uint) error {
	// 检查设备是否存在
	_, err := s.deviceDao.GetSnmpDeviceById(id)
	if err != nil {
		return errors.New("设备不存在")
	}

	return s.deviceDao.DeleteSnmpDevice(id)
}

// GetSnmpDeviceById 根据ID获取SNMP设备
func (s *CmdbSnmpDeviceService) GetSnmpDeviceById(id uint) (*model.CmdbSnmpDevice, error) {
	return s.deviceDao.GetSnmpDeviceById(id)
}

// GetSnmpDeviceList 获取SNMP设备列表(带分页和筛选)
func (s *CmdbSnmpDeviceService) GetSnmpDeviceList(query *model.SnmpDeviceQueryDto) (*model.SnmpDeviceListResult, error) {
	devices, total, err := s.deviceDao.GetSnmpDeviceList(query)
	if err != nil {
		return nil, err
	}

	return &model.SnmpDeviceListResult{
		Items:    devices,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// GetSnmpDeviceStats 获取SNMP设备统计信息
func (s *CmdbSnmpDeviceService) GetSnmpDeviceStats() (*model.SnmpDeviceStatsDto, error) {
	return s.deviceDao.GetSnmpDeviceStats()
}

// GetSnmpDevicesByType 根据设备类型获取设备列表
func (s *CmdbSnmpDeviceService) GetSnmpDevicesByType(deviceType string) ([]model.CmdbSnmpDevice, error) {
	return s.deviceDao.GetSnmpDevicesByType(deviceType)
}

// BatchDeleteSnmpDevices 批量删除SNMP设备
func (s *CmdbSnmpDeviceService) BatchDeleteSnmpDevices(ids []uint) error {
	for _, id := range ids {
		err := s.DeleteSnmpDevice(id)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchUpdateDeviceStatus 批量更新设备状态
func (s *CmdbSnmpDeviceService) BatchUpdateDeviceStatus(ids []uint, status string) error {
	return s.deviceDao.BatchUpdateDeviceStatus(ids, status)
}

// GetMonitorEnabledDevices 获取所有启用监控的设备
func (s *CmdbSnmpDeviceService) GetMonitorEnabledDevices() ([]model.CmdbSnmpDevice, error) {
	return s.deviceDao.GetMonitorEnabledDevices()
}

// SetMonitorEnabled 设置设备监控开关
func (s *CmdbSnmpDeviceService) SetMonitorEnabled(id uint, enabled int) error {
	// 检查设备是否存在
	_, err := s.deviceDao.GetSnmpDeviceById(id)
	if err != nil {
		return errors.New("设备不存在")
	}

	// 更新监控开关
	updateData := &model.CmdbSnmpDevice{
		MonitorEnabled: enabled,
		UpdateTime:     util.HTime{Time: time.Now()},
	}

	return s.deviceDao.UpdateSnmpDevice(id, updateData)
}
