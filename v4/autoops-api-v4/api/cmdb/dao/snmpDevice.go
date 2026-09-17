package dao

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common"
	"gorm.io/gorm"
)

type CmdbSnmpDeviceDao struct {
	db *gorm.DB
}

func NewCmdbSnmpDeviceDao() *CmdbSnmpDeviceDao {
	return &CmdbSnmpDeviceDao{
		db: common.GetDB(),
	}
}

// CreateSnmpDevice 创建SNMP设备
func (d *CmdbSnmpDeviceDao) CreateSnmpDevice(device *model.CmdbSnmpDevice) error {
	return d.db.Create(device).Error
}

// GetSnmpDeviceById 根据ID获取SNMP设备
func (d *CmdbSnmpDeviceDao) GetSnmpDeviceById(id uint) (*model.CmdbSnmpDevice, error) {
	var device model.CmdbSnmpDevice
	err := d.db.Where("id = ?", id).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// GetSnmpDeviceByIP 根据IP获取SNMP设备
func (d *CmdbSnmpDeviceDao) GetSnmpDeviceByIP(ip string) (*model.CmdbSnmpDevice, error) {
	var device model.CmdbSnmpDevice
	err := d.db.Where("ip_address = ?", ip).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// UpdateSnmpDevice 更新SNMP设备
func (d *CmdbSnmpDeviceDao) UpdateSnmpDevice(id uint, device *model.CmdbSnmpDevice) error {
	// 使用map更新,支持零值字段
	updateMap := make(map[string]interface{})

	if device.DeviceName != "" {
		updateMap["device_name"] = device.DeviceName
	}
	if device.DeviceType != "" {
		updateMap["device_type"] = device.DeviceType
	}
	if device.IPAddress != "" {
		updateMap["ip_address"] = device.IPAddress
	}
	if device.SnmpVersion != "" {
		updateMap["snmp_version"] = device.SnmpVersion
	}
	if device.SnmpCommunity != "" {
		updateMap["snmp_community"] = device.SnmpCommunity
	}
	// MonitorEnabled 可以为0,需要特殊处理
	updateMap["monitor_enabled"] = device.MonitorEnabled

	if device.Location != "" {
		updateMap["location"] = device.Location
	}
	if device.Owner != "" {
		updateMap["owner"] = device.Owner
	}
	if device.Remark != "" {
		updateMap["remark"] = device.Remark
	}
	if device.Status != "" {
		updateMap["status"] = device.Status
	}

	// Exporter 相关字段（支持零值）
	updateMap["exporter_installed"] = device.ExporterInstalled
	if device.ExporterPort > 0 {
		updateMap["exporter_port"] = device.ExporterPort
	}
	if device.ExporterVersion != "" {
		updateMap["exporter_version"] = device.ExporterVersion
	}

	// 总是更新 update_time
	updateMap["update_time"] = device.UpdateTime.Time

	return d.db.Model(&model.CmdbSnmpDevice{}).Where("id = ?", id).Updates(updateMap).Error
}

// UpdateSnmpDeviceFields 使用map更新设备字段(支持更新为空值)
func (d *CmdbSnmpDeviceDao) UpdateSnmpDeviceFields(id uint, fields map[string]interface{}) error {
	return d.db.Model(&model.CmdbSnmpDevice{}).Where("id = ?", id).Updates(fields).Error
}

// DeleteSnmpDevice 删除SNMP设备
func (d *CmdbSnmpDeviceDao) DeleteSnmpDevice(id uint) error {
	return d.db.Delete(&model.CmdbSnmpDevice{}, id).Error
}

// GetSnmpDeviceList 获取SNMP设备列表(带分页和筛选)
func (d *CmdbSnmpDeviceDao) GetSnmpDeviceList(query *model.SnmpDeviceQueryDto) ([]model.CmdbSnmpDevice, int64, error) {
	var devices []model.CmdbSnmpDevice
	var total int64

	db := d.db.Model(&model.CmdbSnmpDevice{})

	// 筛选条件
	if query.DeviceType != "" {
		db = db.Where("device_type = ?", query.DeviceType)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.MonitorEnabled >= 0 {
		db = db.Where("monitor_enabled = ?", query.MonitorEnabled)
	}
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("device_name LIKE ? OR ip_address LIKE ? OR location LIKE ?", keyword, keyword, keyword)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	err = db.Order("id DESC").Offset(offset).Limit(query.PageSize).Find(&devices).Error
	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

// GetAllSnmpDevices 获取所有SNMP设备(不分页)
func (d *CmdbSnmpDeviceDao) GetAllSnmpDevices() ([]model.CmdbSnmpDevice, error) {
	var devices []model.CmdbSnmpDevice
	err := d.db.Order("id DESC").Find(&devices).Error
	return devices, err
}

// GetSnmpDevicesByType 根据设备类型获取设备列表
func (d *CmdbSnmpDeviceDao) GetSnmpDevicesByType(deviceType string) ([]model.CmdbSnmpDevice, error) {
	var devices []model.CmdbSnmpDevice
	err := d.db.Where("device_type = ?", deviceType).Find(&devices).Error
	return devices, err
}

// GetMonitorEnabledDevices 获取所有启用监控的设备
func (d *CmdbSnmpDeviceDao) GetMonitorEnabledDevices() ([]model.CmdbSnmpDevice, error) {
	var devices []model.CmdbSnmpDevice
	err := d.db.Where("monitor_enabled = 1").Find(&devices).Error
	return devices, err
}

// UpdateDeviceStatus 更新设备状态
func (d *CmdbSnmpDeviceDao) UpdateDeviceStatus(id uint, status string) error {
	return d.db.Model(&model.CmdbSnmpDevice{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateLastSeen 更新最后在线时间
func (d *CmdbSnmpDeviceDao) UpdateLastSeen(id uint) error {
	return d.db.Model(&model.CmdbSnmpDevice{}).Where("id = ?", id).Update("last_seen", gorm.Expr("NOW()")).Error
}

// GetSnmpDeviceStats 获取SNMP设备统计信息
func (d *CmdbSnmpDeviceDao) GetSnmpDeviceStats() (*model.SnmpDeviceStatsDto, error) {
	stats := &model.SnmpDeviceStatsDto{}

	// 总数
	var totalCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Count(&totalCount)
	stats.TotalCount = int(totalCount)

	// 在线数
	var onlineCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Where("status = 'online'").Count(&onlineCount)
	stats.OnlineCount = int(onlineCount)

	// 离线数
	var offlineCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Where("status = 'offline'").Count(&offlineCount)
	stats.OfflineCount = int(offlineCount)

	// 维护中数
	var maintenanceCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Where("status = 'maintenance'").Count(&maintenanceCount)
	stats.MaintenanceCount = int(maintenanceCount)

	// 按类型统计
	var switchCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Where("device_type = 'switch'").Count(&switchCount)
	stats.SwitchCount = int(switchCount)

	var windowsCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Where("device_type = 'windows'").Count(&windowsCount)
	stats.WindowsCount = int(windowsCount)

	var routerCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Where("device_type = 'router'").Count(&routerCount)
	stats.RouterCount = int(routerCount)

	var otherCount int64
	d.db.Model(&model.CmdbSnmpDevice{}).Where("device_type NOT IN ('switch', 'windows', 'router')").Count(&otherCount)
	stats.OtherCount = int(otherCount)

	return stats, nil
}

// BatchUpdateDeviceStatus 批量更新设备状态
func (d *CmdbSnmpDeviceDao) BatchUpdateDeviceStatus(ids []uint, status string) error {
	return d.db.Model(&model.CmdbSnmpDevice{}).Where("id IN ?", ids).Update("status", status).Error
}

// CheckIPExists 检查IP是否已存在
func (d *CmdbSnmpDeviceDao) CheckIPExists(ip string, excludeID uint) (bool, error) {
	var count int64
	query := d.db.Model(&model.CmdbSnmpDevice{}).Where("ip_address = ?", ip)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
