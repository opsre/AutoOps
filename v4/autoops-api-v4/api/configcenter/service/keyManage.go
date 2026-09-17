package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	cmdbService "dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	aliyuncloud "dodevops-api/common/util"
	tengxuncloud "dodevops-api/common/util"
	baiducloud "dodevops-api/common/util"
	huaweicloud "dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

// SyncResult 同步结果统计
type SyncResult struct {
	KeyID        uint           // 密钥ID
	KeyName      string         // 密钥名称
	CloudType    string         // 云厂商类型
	TotalFetched int            // 从云端获取的总数
	NewAdded     int            // 新增主机数
	FailedCount  int            // 失败数量
	SkippedCount int            // 跳过数量（已存在或非运行状态）
	RegionStats  map[string]int // 区域分布统计
	Error        string         // 错误信息（如果有）
}

type KeyManageService struct {
	dao *dao.KeyManageDao
}

func NewKeyManageService() *KeyManageService {
	return &KeyManageService{
		dao: dao.NewKeyManageDao(),
	}
}

// Create 创建密钥
func (s *KeyManageService) Create(keyManage *model.KeyManage) error {
	return s.dao.Create(keyManage)
}

// Update 更新密钥
func (s *KeyManageService) Update(keyManage *model.KeyManage) error {
	return s.dao.Update(keyManage)
}

// Delete 删除密钥
func (s *KeyManageService) Delete(id uint) error {
	return s.dao.Delete(id)
}

// GetByID 根据ID查询密钥
func (s *KeyManageService) GetByID(id uint) (*model.KeyManage, error) {
	return s.dao.GetByID(id)
}

// List 获取密钥列表
func (s *KeyManageService) List() ([]model.KeyManage, error) {
	return s.dao.List()
}

// ListWithPage 获取密钥列表（分页）
func (s *KeyManageService) ListWithPage(page, pageSize int) ([]model.KeyManage, int64, error) {
	return s.dao.ListWithPage(page, pageSize)
}

// DecryptKeys 解密密钥信息
func (s *KeyManageService) DecryptKeys(id uint) (string, string, error) {
	keyManage, err := s.dao.GetByID(id)
	if err != nil {
		return "", "", err
	}
	return keyManage.DecryptKeys()
}

// GetByType 根据云厂商类型查询密钥
func (s *KeyManageService) GetByType(keyType int) ([]model.KeyManage, error) {
	return s.dao.GetByType(keyType)
}

// GetDecryptedKeyForCloudAPI 获取解密后的密钥信息用于调用云API
// 返回解密后的KeyID和KeySecret，专用于云API调用
func (s *KeyManageService) GetDecryptedKeyForCloudAPI(id uint) (keyID string, keySecret string, keyType int, err error) {
	keyManage, err := s.dao.GetByID(id)
	if err != nil {
		return "", "", 0, err
	}
	
	decryptedKeyID, decryptedKeySecret, err := keyManage.DecryptKeys()
	if err != nil {
		return "", "", 0, err
	}
	
	return decryptedKeyID, decryptedKeySecret, keyManage.KeyType, nil
}

// GetDecryptedKeyByType 根据云厂商类型获取第一个可用的解密密钥
// 用于自动选择对应云厂商的密钥调用API
func (s *KeyManageService) GetDecryptedKeyByType(keyType int) (keyID string, keySecret string, err error) {
	keyManages, err := s.dao.GetByType(keyType)
	if err != nil {
		return "", "", err
	}
	
	if len(keyManages) == 0 {
		return "", "", fmt.Errorf("未找到类型为 %d 的密钥配置", keyType)
	}
	
	// 使用第一个找到的密钥
	decryptedKeyID, decryptedKeySecret, err := keyManages[0].DecryptKeys()
	if err != nil {
		return "", "", err
	}
	
	return decryptedKeyID, decryptedKeySecret, nil
}

// SyncAliyunHosts 同步阿里云主机（使用密钥管理中的凭据）
func (s *KeyManageService) SyncAliyunHosts(c *gin.Context, keyID uint, groupID uint, region string) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "阿里云主机同步请求已接收，正在后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理同步
	go func() {
		hostDao := cmdbDao.NewCmdbHostDao()
		
		hosts, err := s.syncAliyunHosts(keyID, groupID, region)
		if err != nil {
			fmt.Printf("[ERROR] 阿里云主机同步失败: %v\n", err)
			return
		}

		// 过滤掉已存在的主机
		newHosts := s.filterExistingHosts(hosts)
		fmt.Printf("[DEBUG] 准备同步 %d 台新的阿里云主机（总共获取 %d 台）\n", len(newHosts), len(hosts))

		// 批量导入新主机
		successCount := 0
		for _, host := range newHosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 阿里云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 阿里云主机同步成功: %s\n", host.InstanceID)
		}

		fmt.Printf("[INFO] 阿里云主机同步完成: 成功 %d 台, 失败 %d 台, 跳过已存在 %d 台\n",
			successCount, len(newHosts)-successCount, len(hosts)-len(newHosts))
	}()
}

// SyncAliyunHostsBackground 后台同步阿里云主机（用于定时任务，不依赖gin.Context）
// 使用精细化逻辑：批量导入、标签分组、去重检测、区域统计
func (s *KeyManageService) SyncAliyunHostsBackground(keyID uint, groupID uint, region string) (*SyncResult, error) {
	log.Printf("开始后台同步阿里云主机: keyID=%d, groupID=%d, region=%s", keyID, groupID, region)

	result := &SyncResult{
		KeyID:       keyID,
		CloudType:   "阿里云",
		RegionStats: make(map[string]int),
	}

	// 获取密钥名称和凭证
	keyManage, err := s.GetByID(keyID)
	if err != nil {
		result.Error = fmt.Sprintf("获取密钥信息失败: %v", err)
		log.Printf("阿里云主机同步失败: %v", err)
		return result, err
	}
	result.KeyName = keyManage.KeyName

	// 解密密钥
	accessKey, secretKey, err := keyManage.DecryptKeys()
	if err != nil {
		result.Error = fmt.Sprintf("解密密钥失败: %v", err)
		log.Printf("阿里云密钥解密失败: %v", err)
		return result, err
	}

	// 构造DTO，调用精细化同步逻辑
	dto := &cmdbModel.CreateCmdbHostCloudDto{
		AccessKey: accessKey,
		SecretKey: secretKey,
		GroupID:   groupID,
		Region:    region,
	}

	// 调用 CmdbHostcloud 的精细化同步逻辑
	cloudService := cmdbService.GetCmdbHostCloudService()
	syncResult, err := cloudService.SyncAliyunHostBackground(dto)
	if err != nil {
		result.Error = fmt.Sprintf("同步失败: %v", err)
		log.Printf("阿里云主机同步失败: %v", err)
		return result, err
	}

	// 映射结果
	result.TotalFetched = syncResult.TotalFetched
	result.NewAdded = syncResult.NewAdded
	result.FailedCount = syncResult.FailedCount
	result.SkippedCount = syncResult.SkippedCount
	result.RegionStats = syncResult.RegionStats

	log.Printf("阿里云主机同步完成: 成功 %d 台，失败 %d 台，跳过 %d 台",
		result.NewAdded, result.FailedCount, result.SkippedCount)
	return result, nil
}

// syncAliyunHosts 获取阿里云主机信息
func (s *KeyManageService) syncAliyunHosts(keyID uint, groupID uint, region string) ([]cmdbModel.CmdbHost, error) {
	// 获取解密后的密钥
	accessKey, secretKey, keyType, err := s.GetDecryptedKeyForCloudAPI(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取密钥失败: %v", err)
	}
	
	// 检查密钥类型是否为阿里云
	if keyType != 1 {
		return nil, fmt.Errorf("密钥类型不正确，期望阿里云(1)，实际为: %d", keyType)
	}

	aliyunService := aliyuncloud.NewAliyunCloudService(accessKey, secretKey, region)
	
	instances, err := aliyunService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取阿里云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到指定实例")
	}

	var hosts []cmdbModel.CmdbHost
	for _, instance := range instances {
		if strings.TrimSpace(instance.Region) == "" {
			instance.Region = "all"
		}

		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			continue
		}

		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}
		if instance.ExpireTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.ExpireTime); err == nil {
				expireTime = t
			}
		}

		// 过滤已过期的实例
		if !expireTime.IsZero() && expireTime.Before(time.Now()) {
			fmt.Printf("[INFO] 跳过已过期实例: %s (过期时间: %s)\n", instance.InstanceID, expireTime.Format("2006-01-02 15:04:05"))
			continue
		}

		host := cmdbModel.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      2, // 阿里云
			Region:      strings.TrimSpace(instance.Region),
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2, // 未认证
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.MemoryGB),
			Disk:        strings.TrimSuffix(instance.DiskGB, "GB"),
			BillingType: instance.BillingType,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return nil, errors.New("没有有效的实例可同步")
	}
	return hosts, nil
}

// SyncTencentHosts 同步腾讯云主机（使用密钥管理中的凭据）
func (s *KeyManageService) SyncTencentHosts(c *gin.Context, keyID uint, groupID uint) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "腾讯云主机同步请求已接收，正在后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理同步
	go func() {
		hostDao := cmdbDao.NewCmdbHostDao()
		
		hosts, err := s.syncTencentHosts(keyID, groupID)
		if err != nil {
			fmt.Printf("[ERROR] 腾讯云主机同步失败: %v\n", err)
			return
		}

		// 过滤掉已存在的主机
		newHosts := s.filterExistingHosts(hosts)
		fmt.Printf("[DEBUG] 准备同步 %d 台新的腾讯云主机（总共获取 %d 台）\n", len(newHosts), len(hosts))

		// 批量导入新主机
		successCount := 0
		for _, host := range newHosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 腾讯云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 腾讯云主机同步成功: %s\n", host.InstanceID)
		}

		fmt.Printf("[INFO] 腾讯云主机同步完成: 成功 %d 台, 失败 %d 台, 跳过已存在 %d 台\n",
			successCount, len(newHosts)-successCount, len(hosts)-len(newHosts))
	}()
}

// SyncTencentHostsBackground 后台同步腾讯云主机（用于定时任务，不依赖gin.Context）
// 使用精细化逻辑：批量导入、标签分组、去重检测、区域统计
func (s *KeyManageService) SyncTencentHostsBackground(keyID uint, groupID uint) (*SyncResult, error) {
	log.Printf("开始后台同步腾讯云主机: keyID=%d, groupID=%d", keyID, groupID)

	result := &SyncResult{
		KeyID:       keyID,
		CloudType:   "腾讯云",
		RegionStats: make(map[string]int),
	}

	// 获取密钥名称和凭证
	keyManage, err := s.GetByID(keyID)
	if err != nil {
		result.Error = fmt.Sprintf("获取密钥信息失败: %v", err)
		log.Printf("腾讯云主机同步失败: %v", err)
		return result, err
	}
	result.KeyName = keyManage.KeyName

	// 解密密钥
	accessKey, secretKey, err := keyManage.DecryptKeys()
	if err != nil {
		result.Error = fmt.Sprintf("解密密钥失败: %v", err)
		log.Printf("腾讯云密钥解密失败: %v", err)
		return result, err
	}

	// 构造DTO，调用精细化同步逻辑
	dto := &cmdbModel.CreateCmdbHostCloudDto{
		AccessKey: accessKey,
		SecretKey: secretKey,
		GroupID:   groupID,
	}

	// 调用 CmdbHostcloud 的精细化同步逻辑
	cloudService := cmdbService.GetCmdbHostCloudService()
	syncResult, err := cloudService.SyncTencentHostBackground(dto)
	if err != nil {
		result.Error = fmt.Sprintf("同步失败: %v", err)
		log.Printf("腾讯云主机同步失败: %v", err)
		return result, err
	}

	// 映射结果
	result.TotalFetched = syncResult.TotalFetched
	result.NewAdded = syncResult.NewAdded
	result.FailedCount = syncResult.FailedCount
	result.SkippedCount = syncResult.SkippedCount
	result.RegionStats = syncResult.RegionStats

	log.Printf("腾讯云主机同步完成: 成功 %d 台，失败 %d 台，跳过 %d 台",
		result.NewAdded, result.FailedCount, result.SkippedCount)
	return result, nil
}

// syncTencentHosts 获取腾讯云主机信息
func (s *KeyManageService) syncTencentHosts(keyID uint, groupID uint) ([]cmdbModel.CmdbHost, error) {
	// 获取解密后的密钥
	accessKey, secretKey, keyType, err := s.GetDecryptedKeyForCloudAPI(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取密钥失败: %v", err)
	}
	
	// 检查密钥类型是否为腾讯云
	if keyType != 2 {
		return nil, fmt.Errorf("密钥类型不正确，期望腾讯云(2)，实际为: %d", keyType)
	}

	tencentService := tengxuncloud.NewTencentCloudService(accessKey, secretKey)
	
	instances, err := tencentService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取腾讯云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到指定实例")
	}

	var hosts []cmdbModel.CmdbHost
	for _, instance := range instances {
		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			continue
		}

		totalDisk := instance.SystemDisk
		for _, disk := range instance.DataDisks {
			totalDisk += disk
		}

		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}
		if instance.ExpireTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.ExpireTime); err == nil {
				expireTime = t
			}
		}

		// 过滤已过期的实例
		if !expireTime.IsZero() && expireTime.Before(time.Now()) {
			fmt.Printf("[INFO] 跳过已过期实例: %s (过期时间: %s)\n", instance.InstanceID, expireTime.Format("2006-01-02 15:04:05"))
			continue
		}

		host := cmdbModel.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      3, // 腾讯云
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2, // 未认证
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.Memory),
			Disk:        strconv.Itoa(totalDisk),
			BillingType: "包年包月",
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return nil, errors.New("没有有效的实例可同步")
	}
	return hosts, nil
}

// filterExistingHosts 过滤掉已存在的主机，避免重复导入（基于instance_id去重）
func (s *KeyManageService) filterExistingHosts(hosts []cmdbModel.CmdbHost) []cmdbModel.CmdbHost {
	var newHosts []cmdbModel.CmdbHost
	hostDao := cmdbDao.NewCmdbHostDao()

	for _, host := range hosts {
		// 通过instance_id来判断是否存在（更可靠的去重方式）
		existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
		if existing == nil {
			// 主机不存在，可以添加
			newHosts = append(newHosts, host)
		} else {
			fmt.Printf("[INFO] 主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
		}
	}

	return newHosts
}



// SyncBaiduHosts 同步百度云主机（使用密钥管理中的凭据）
func (s *KeyManageService) SyncBaiduHosts(c *gin.Context, keyID uint, groupID uint) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "百度云主机同步请求已接收，正在扫描所有区域并后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理同步
	go func() {
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.syncBaiduHosts(keyID, groupID)
		if err != nil {
			fmt.Printf("[ERROR] 百度云主机同步失败: %v\n", err)
			return
		}

		// 过滤掉已存在的主机
		newHosts := s.filterExistingHosts(hosts)
		fmt.Printf("[DEBUG] 准备同步 %d 台新的百度云主机（总共获取 %d 台）\n", len(newHosts), len(hosts))

		// 批量导入新主机
		successCount := 0
		for _, host := range newHosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 百度云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 百度云主机同步成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		fmt.Printf("[INFO] 百度云主机同步完成: 成功 %d 台, 失败 %d 台, 跳过已存在 %d 台\n",
			successCount, len(newHosts)-successCount, len(hosts)-len(newHosts))
	}()
}

// SyncBaiduHostsBackground 后台同步百度云主机（用于定时任务，不依赖gin.Context）
// 使用精细化逻辑：批量导入、标签分组、去重检测、区域统计
func (s *KeyManageService) SyncBaiduHostsBackground(keyID uint, groupID uint) (*SyncResult, error) {
	log.Printf("开始后台同步百度云主机: keyID=%d, groupID=%d", keyID, groupID)

	result := &SyncResult{
		KeyID:       keyID,
		CloudType:   "百度云",
		RegionStats: make(map[string]int),
	}

	// 获取密钥名称和凭证
	keyManage, err := s.GetByID(keyID)
	if err != nil {
		result.Error = fmt.Sprintf("获取密钥信息失败: %v", err)
		log.Printf("百度云主机同步失败: %v", err)
		return result, err
	}
	result.KeyName = keyManage.KeyName

	// 解密密钥
	accessKey, secretKey, err := keyManage.DecryptKeys()
	if err != nil {
		result.Error = fmt.Sprintf("解密密钥失败: %v", err)
		log.Printf("百度云密钥解密失败: %v", err)
		return result, err
	}

	// 构造DTO，调用精细化同步逻辑
	dto := &cmdbModel.CreateCmdbHostCloudDto{
		AccessKey: accessKey,
		SecretKey: secretKey,
		GroupID:   groupID,
	}

	// 调用 CmdbHostcloud 的精细化同步逻辑
	cloudService := cmdbService.GetCmdbHostCloudService()
	syncResult, err := cloudService.SyncBaiduHostBackground(dto)
	if err != nil {
		result.Error = fmt.Sprintf("同步失败: %v", err)
		log.Printf("百度云主机同步失败: %v", err)
		return result, err
	}

	// 映射结果
	result.TotalFetched = syncResult.TotalFetched
	result.NewAdded = syncResult.NewAdded
	result.FailedCount = syncResult.FailedCount
	result.SkippedCount = syncResult.SkippedCount
	result.RegionStats = syncResult.RegionStats

	log.Printf("百度云主机同步完成: 成功 %d 台，失败 %d 台，跳过 %d 台",
		result.NewAdded, result.FailedCount, result.SkippedCount)
	return result, nil
}

// syncBaiduHosts 获取百度云主机信息（自动扫描所有区域）
func (s *KeyManageService) syncBaiduHosts(keyID uint, groupID uint) ([]cmdbModel.CmdbHost, error) {
	// 获取解密后的密钥
	accessKey, secretKey, keyType, err := s.GetDecryptedKeyForCloudAPI(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取密钥失败: %v", err)
	}

	// 检查密钥类型是否为百度云
	if keyType != 3 {
		return nil, fmt.Errorf("密钥类型不正确，期望百度云(3)，实际为: %d", keyType)
	}

	baiduService := baiducloud.NewBaiduCloudService(accessKey, secretKey)

	// 获取所有区域的实例（自动扫描所有区域）
	instances, err := baiduService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取百度云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何百度云实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台百度云实例\n", len(instances))

	var hosts []cmdbModel.CmdbHost
	for _, instance := range instances {
		// 确定SSH连接IP
		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			fmt.Printf("[WARN] 百度云实例无有效IP，跳过: %s\n", instance.InstanceID)
			continue
		}

		// 计算总磁盘容量
		totalDisk := instance.SystemDisk
		for _, disk := range instance.DataDisks {
			totalDisk += disk
		}

		// 解析时间
		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}
		if instance.ExpireTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.ExpireTime); err == nil {
				expireTime = t
			}
		}

		// 确定主机状态
		status := 2 // 默认在线
		if instance.Status == "Stopped" {
			status = 3 // 离线
		}

		host := cmdbModel.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      3, // 百度云厂商代码为3
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      status,
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.Memory),
			Disk:        strconv.Itoa(totalDisk),
			BillingType: instance.BillingType,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return nil, errors.New("没有有效的百度云实例可同步")
	}

	fmt.Printf("[INFO] 成功解析 %d 台百度云主机\n", len(hosts))
	return hosts, nil
}

// SyncHuaweiHosts 同步华为云主机（使用密钥管理中的凭据）
func (s *KeyManageService) SyncHuaweiHosts(c *gin.Context, keyID uint, groupID uint) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "华为云主机同步请求已接收，正在扫描所有区域并后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理同步
	go func() {
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.syncHuaweiHosts(keyID, groupID)
		if err != nil {
			fmt.Printf("[ERROR] 华为云主机同步失败: %v\n", err)
			return
		}

		// 过滤掉已存在的主机
		newHosts := s.filterExistingHosts(hosts)
		fmt.Printf("[DEBUG] 准备同步 %d 台新的华为云主机（总共获取 %d 台）\n", len(newHosts), len(hosts))

		// 批量导入新主机
		successCount := 0
		for _, host := range newHosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 华为云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 华为云主机同步成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		fmt.Printf("[INFO] 华为云主机同步完成: 成功 %d 台, 失败 %d 台, 跳过已存在 %d 台\n",
			successCount, len(newHosts)-successCount, len(hosts)-len(newHosts))
	}()
}

// SyncHuaweiHostsBackground 后台同步华为云主机（用于定时任务，不依赖gin.Context）
// 使用精细化逻辑：批量导入、标签分组、去重检测、区域统计
func (s *KeyManageService) SyncHuaweiHostsBackground(keyID uint, groupID uint) (*SyncResult, error) {
	log.Printf("开始后台同步华为云主机: keyID=%d, groupID=%d", keyID, groupID)

	result := &SyncResult{
		KeyID:       keyID,
		CloudType:   "华为云",
		RegionStats: make(map[string]int),
	}

	// 获取密钥名称和凭证
	keyManage, err := s.GetByID(keyID)
	if err != nil {
		result.Error = fmt.Sprintf("获取密钥信息失败: %v", err)
		log.Printf("华为云主机同步失败: %v", err)
		return result, err
	}
	result.KeyName = keyManage.KeyName

	// 解密密钥
	accessKey, secretKey, err := keyManage.DecryptKeys()
	if err != nil {
		result.Error = fmt.Sprintf("解密密钥失败: %v", err)
		log.Printf("华为云密钥解密失败: %v", err)
		return result, err
	}

	// 构造DTO，调用精细化同步逻辑
	dto := &cmdbModel.CreateCmdbHostCloudDto{
		AccessKey: accessKey,
		SecretKey: secretKey,
		GroupID:   groupID,
	}

	// 调用 CmdbHostcloud 的精细化同步逻辑
	cloudService := cmdbService.GetCmdbHostCloudService()
	syncResult, err := cloudService.SyncHuaweiHostBackground(dto)
	if err != nil {
		result.Error = fmt.Sprintf("同步失败: %v", err)
		log.Printf("华为云主机同步失败: %v", err)
		return result, err
	}

	// 映射结果
	result.TotalFetched = syncResult.TotalFetched
	result.NewAdded = syncResult.NewAdded
	result.FailedCount = syncResult.FailedCount
	result.SkippedCount = syncResult.SkippedCount
	result.RegionStats = syncResult.RegionStats

	log.Printf("华为云主机同步完成: 成功 %d 台，失败 %d 台，跳过 %d 台",
		result.NewAdded, result.FailedCount, result.SkippedCount)
	return result, nil
}

// syncHuaweiHosts 获取华为云主机信息（自动扫描所有区域）
func (s *KeyManageService) syncHuaweiHosts(keyID uint, groupID uint) ([]cmdbModel.CmdbHost, error) {
	// 获取解密后的密钥
	accessKey, secretKey, keyType, err := s.GetDecryptedKeyForCloudAPI(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取密钥失败: %v", err)
	}

	// 检查密钥类型是否为华为云
	if keyType != 4 {
		return nil, fmt.Errorf("密钥类型不正确，期望华为云(4)，实际为: %d", keyType)
	}

	huaweiService := huaweicloud.NewHuaweiCloudService(accessKey, secretKey)

	// 获取所有区域的实例（自动扫描所有区域）
	instances, err := huaweiService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取华为云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何华为云实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台华为云实例\n", len(instances))

	var hosts []cmdbModel.CmdbHost
	for _, instance := range instances {
		// 确定SSH连接IP
		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			fmt.Printf("[WARN] 华为云实例无有效IP，跳过: %s\n", instance.InstanceID)
			continue
		}

		// 计算总磁盘容量
		totalDisk := instance.SystemDisk
		for _, disk := range instance.DataDisks {
			totalDisk += disk
		}

		// 解析时间
		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}
		if instance.ExpireTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.ExpireTime); err == nil {
				expireTime = t
			}
		}

		host := cmdbModel.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      4, // 华为云厂商代码为4
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2, // 默认未认证状态
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.MemoryGB),
			Disk:        strconv.Itoa(totalDisk),
			BillingType: instance.BillingType,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return nil, errors.New("没有有效的华为云实例可同步")
	}

	fmt.Printf("[INFO] 成功解析 %d 台华为云主机\n", len(hosts))
	return hosts, nil
}

// SyncAWSHostsBackground 后台同步AWS主机（用于定时任务，不依赖gin.Context）
func (s *KeyManageService) SyncAWSHostsBackground(keyID uint, groupID uint) (*SyncResult, error) {
	log.Printf("开始后台同步AWS主机: keyID=%d, groupID=%d", keyID, groupID)

	result := &SyncResult{
		KeyID:       keyID,
		CloudType:   "AWS云",
		RegionStats: make(map[string]int),
	}

	// 获取密钥名称和凭证
	keyManage, err := s.GetByID(keyID)
	if err != nil {
		result.Error = fmt.Sprintf("获取密钥信息失败: %v", err)
		log.Printf("AWS主机同步失败: %v", err)
		return result, err
	}
	result.KeyName = keyManage.KeyName

	// 解密密钥
	accessKey, secretKey, err := keyManage.DecryptKeys()
	if err != nil {
		result.Error = fmt.Sprintf("解密密钥失败: %v", err)
		log.Printf("AWS密钥解密失败: %v", err)
		return result, err
	}

	// 构造DTO，调用精细化同步逻辑
	dto := &cmdbModel.CreateCmdbHostCloudDto{
		AccessKey: accessKey,
		SecretKey: secretKey,
		GroupID:   groupID,
	}

	// 调用 CmdbHostcloud 的精细化同步逻辑
	cloudService := cmdbService.GetCmdbHostCloudService()
	syncResult, err := cloudService.SyncAWSHostBackground(dto)
	if err != nil {
		result.Error = fmt.Sprintf("同步失败: %v", err)
		log.Printf("AWS主机同步失败: %v", err)
		return result, err
	}

	// 映射结果
	result.TotalFetched = syncResult.TotalFetched
	result.NewAdded = syncResult.NewAdded
	result.FailedCount = syncResult.FailedCount
	result.SkippedCount = syncResult.SkippedCount
	result.RegionStats = syncResult.RegionStats

	log.Printf("AWS主机同步完成: 成功 %d 台，失败 %d 台，跳过 %d 台",
		result.NewAdded, result.FailedCount, result.SkippedCount)
	return result, nil
}