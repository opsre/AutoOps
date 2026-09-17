package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	cmdbDao "dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

// SyncResult 同步结果统计
type SyncResult struct {
	TotalFetched int            // 从云端获取的总数
	NewAdded     int            // 新增主机数
	FailedCount  int            // 失败数量
	SkippedCount int            // 跳过数量（已存在或非运行状态）
	RegionStats  map[string]int // 区域分布统计
}

type CmdbHostCloudServiceInterface interface {
    CreateAliyunHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
    CreateTencentHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
    CreateBaiduHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
    CreateHuaweiHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
    CreateAWSHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
    // 后台同步版本（不依赖gin.Context）
    SyncAliyunHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error)
    SyncTencentHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error)
    SyncBaiduHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error)
    SyncHuaweiHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error)
    SyncAWSHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error)
}

type CmdbHostCloudServiceImpl struct {
	hostDao cmdbDao.CmdbHostDao
}

// 创建阿里云主机
func (s *CmdbHostCloudServiceImpl) CreateAliyunHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "阿里云主机导入请求已接收，正在扫描所有区域",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理导入
	go func() {
		// 初始化DAO
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.createAliyunHosts(dto)
		if err != nil {
			fmt.Printf("[ERROR] 阿里云主机导入失败: %v\n", err)
			return
		}

		// 打印所有要导入的主机信息
		fmt.Printf("[DEBUG] 准备导入 %d 台阿里云主机\n", len(hosts))
		for i, host := range hosts {
			fmt.Printf("[DEBUG] 主机 %d: %+v\n", i+1, host)
		}

		// 批量导入 - 每批50台
		const batchSize = 50
		successCount := 0
		totalBatches := (len(hosts) + batchSize - 1) / batchSize

		for i := 0; i < len(hosts); i += batchSize {
			end := i + batchSize
			if end > len(hosts) {
				end = len(hosts)
			}
			batch := hosts[i:end]
			batchNum := i/batchSize + 1

			fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

			// 使用批量插入方法
			if err := hostDao.BatchCreateCmdbHosts(batch); err != nil {
				fmt.Printf("[ERROR] 第 %d 批主机保存失败: %v\n", batchNum, err)
				// 批量插入失败，尝试逐条插入以确定具体失败的记录
				for _, host := range batch {
					if err := hostDao.CreateCmdbHost(&host); err != nil {
						fmt.Printf("[ERROR] 阿里云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
					} else {
						successCount++
						fmt.Printf("[INFO] 阿里云主机导入成功: %s\n", host.InstanceID)
					}
				}
			} else {
				successCount += len(batch)
				fmt.Printf("[INFO] 第 %d 批导入成功: %d 台\n", batchNum, len(batch))
			}
		}

		// 统计导入成功的主机按区域分布
		successRegionStats := make(map[string]int)
		successHosts := 0
		for _, host := range hosts {
			// 检查主机是否导入成功（通过查询数据库）
			if existing, err := hostDao.GetCmdbHostByInstanceID(host.InstanceID); err == nil && existing != nil {
				successRegionStats[host.Region]++
				successHosts++
			}
		}

		// 记录最终结果
		fmt.Printf("[INFO] 阿里云主机导入完成: 成功 %d 台, 失败 %d 台\n",
			successCount, len(hosts)-successCount)
		if successHosts > 0 {
			fmt.Println("[INFO] 导入成功的主机区域分布:")
			for region, count := range successRegionStats {
				fmt.Printf("  - %s: %d 台\n", region, count)
			}
		}

		// 这里可以添加将结果写入Redis或数据库的逻辑，供前端查询
	}()
}

func (s *CmdbHostCloudServiceImpl) createAliyunHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
	// 注意: dto.Region 参数在此处被忽略，自动扫描所有区域
	// 如果需要限制区域，可以在调用时传入空字符串表示"所有区域"
	aliyunService := util.NewAliyunCloudService(dto.AccessKey, dto.SecretKey, "")

	instances, err := aliyunService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取阿里云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何阿里云实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台阿里云实例\n", len(instances))

	// 统计每个区域的实例数量
	regionStats := make(map[string]int)
	for _, instance := range instances {
		regionStats[instance.Region]++
	}
	fmt.Println("[INFO] 区域分布统计:")
	for region, count := range regionStats {
		fmt.Printf("  - %s: %d 台\n", region, count)
	}

	// 初始化分组DAO
	groupDao := cmdbDao.NewCmdbGroupDao()

	// 初始化主机DAO用于去重检查
	hostDao := cmdbDao.NewCmdbHostDao()

	var hosts []model.CmdbHost
	skippedCount := 0
	duplicateCount := 0
	updatedCount := 0
	for _, instance := range instances {
		// 过滤非Running状态的实例(Expired/Stopped等)
		if instance.Status != "Running" {
			fmt.Printf("[INFO] 跳过非运行状态实例: %s (状态: %s)\n", instance.InstanceID, instance.Status)
			skippedCount++
			continue
		}

		// 检查实例是否已存在（根据InstanceID去重）
		if existing, _ := hostDao.GetCmdbHostByInstanceID(instance.InstanceID); existing != nil {
			// 检查配置是否发生变化
			needUpdate := false

			// 比对CPU
			newCPU := strconv.Itoa(instance.CPU)
			if existing.CPU != newCPU {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机CPU变更: %s (%s -> %s)\n", instance.InstanceID, existing.CPU, newCPU)
			}

			// 比对内存
			newMemory := strconv.Itoa(instance.MemoryGB)
			if existing.Memory != newMemory {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机内存变更: %s (%s -> %s)\n", instance.InstanceID, existing.Memory, newMemory)
			}

			// 比对磁盘
			newDisk := strings.TrimSuffix(instance.DiskGB, "GB")
			if existing.Disk != newDisk {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机磁盘变更: %s (%s -> %s)\n", instance.InstanceID, existing.Disk, newDisk)
			}

			// 比对IP地址
			newPrivateIP := strings.Join(instance.PrivateIPs, ",")
			if existing.PrivateIP != newPrivateIP {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机内网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PrivateIP, newPrivateIP)
			}

			newPublicIP := strings.Join(instance.PublicIPs, ",")
			if existing.PublicIP != newPublicIP {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机公网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PublicIP, newPublicIP)
			}

			// 比对SSH IP
			sshIP := ""
			if len(instance.PublicIPs) > 0 {
				sshIP = instance.PublicIPs[0]
			} else if len(instance.PrivateIPs) > 0 {
				sshIP = instance.PrivateIPs[0]
			}
			if sshIP != "" && existing.SSHIP != sshIP {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机SSH IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.SSHIP, sshIP)
			}

			// 比对操作系统
			if existing.OS != instance.OSName {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机OS变更: %s (%s -> %s)\n", instance.InstanceID, existing.OS, instance.OSName)
			}

			// 比对区域
			newRegion := strings.TrimSpace(instance.Region)
			if newRegion == "" {
				newRegion = "all"
			}
			if existing.Region != newRegion {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机区域变更: %s (%s -> %s)\n", instance.InstanceID, existing.Region, newRegion)
			}

			// 比对计费类型
			if existing.BillingType != instance.BillingType {
				needUpdate = true
				fmt.Printf("[INFO] 阿里云主机计费类型变更: %s (%s -> %s)\n", instance.InstanceID, existing.BillingType, instance.BillingType)
			}

			// 如果有变更，执行更新
			if needUpdate {
				if err := hostDao.UpdateCmdbHost(existing.ID, &model.CmdbHost{
					CPU:         newCPU,
					Memory:      newMemory,
					Disk:        newDisk,
					PrivateIP:   newPrivateIP,
					PublicIP:    newPublicIP,
					SSHIP:       sshIP,
					OS:          instance.OSName,
					Region:      newRegion,
					BillingType: instance.BillingType,
				}); err != nil {
					fmt.Printf("[ERROR] 阿里云主机配置更新失败: %s - %v\n", instance.InstanceID, err)
				} else {
					updatedCount++
					fmt.Printf("[SUCCESS] 阿里云主机配置已更新: %s\n", instance.InstanceID)
				}
			} else {
				fmt.Printf("[INFO] 阿里云主机配置无变化，跳过: %s (名称: %s, 区域: %s)\n", instance.InstanceID, instance.InstanceName, instance.Region)
			}

			duplicateCount++
			continue
		}

		if strings.TrimSpace(instance.Region) == "" {
			instance.Region = "all"
		}

		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			fmt.Printf("[WARN] 实例无有效IP，跳过: %s\n", instance.InstanceID)
			skippedCount++
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

		// 解析标签并匹配分组
		groupID, tagString := s.matchGroupByTags(instance.Tags, &groupDao)

		// 使用实例名称作为主机名,如果实例名称为空则使用实例ID
		hostName := instance.InstanceName
		if hostName == "" {
			hostName = instance.InstanceID
		}

		// 确定SSH端口，优先使用用户配置的端口，否则默认22
		sshPort := 22
		if dto.SSHPort > 0 && dto.SSHPort <= 65535 {
			sshPort = dto.SSHPort
		}

		host := model.CmdbHost{
			HostName:    hostName,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     sshPort,
			SSHName:     "root",
			Vendor:      2,
			Region:      strings.TrimSpace(instance.Region),
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2,
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.MemoryGB),
			Disk:        strings.TrimSuffix(instance.DiskGB, "GB"),
			BillingType: instance.BillingType,
			Tag:         tagString,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if skippedCount > 0 || duplicateCount > 0 || updatedCount > 0 {
		fmt.Printf("[INFO] 共跳过 %d 台非运行状态或无效实例, %d 台重复实例, %d 台配置已更新\n", skippedCount, duplicateCount, updatedCount)
	}

	if len(hosts) == 0 {
		if duplicateCount > 0 {
			if updatedCount > 0 {
				return nil, fmt.Errorf("所有实例均已导入，无新增实例（已更新 %d 台主机配置）", updatedCount)
			}
			return nil, errors.New("所有实例均已导入，无新增实例")
		}
		return nil, errors.New("没有处于运行状态的有效实例可导入")
	}
	return hosts, nil
}

// matchGroupByTags 根据标签匹配分组
// 返回: (分组ID, 标签字符串)
// 匹配规则:
// 1. 从标签中提取 bline 值
// 2. 先尝试精确匹配子分组名称(任意层级)
// 3. 如果匹配失败,返回默认分组ID(1)
func (s *CmdbHostCloudServiceImpl) matchGroupByTags(tags map[string]string, groupDao *cmdbDao.CmdbGroupDao) (uint, string) {
	// 默认分组ID
	defaultGroupID := uint(1)

	// 如果没有标签,返回默认分组
	if len(tags) == 0 {
		return defaultGroupID, ""
	}

	// 将标签转换为字符串格式: key=value,key=value
	var tagPairs []string
	for k, v := range tags {
		tagPairs = append(tagPairs, fmt.Sprintf("%s=%s", k, v))
	}
	tagString := strings.Join(tagPairs, ",")

	// 获取 bline 标签值
	bline, hasBline := tags["bline"]
	if !hasBline || bline == "" {
		fmt.Printf("[INFO] 标签中没有 bline 字段,使用默认分组 (标签: %s)\n", tagString)
		return defaultGroupID, tagString
	}

	// 直接使用bline值匹配分组名称(任意层级)
	// 例如: bline=saas3 -> 查找名称为"saas3"的分组
	if group, err := groupDao.GetCmdbGroupByExactName(bline); err == nil && group != nil {
		// 获取父分组信息用于日志输出
		var parentName string
		if group.ParentID != 0 {
			if parent, err := groupDao.GetCmdbGroupById(group.ParentID); err == nil {
				parentName = parent.Name
				fmt.Printf("[INFO] 标签匹配成功: %s -> 分组 %s/%s (ID: %d)\n", tagString, parentName, group.Name, group.ID)
			} else {
				fmt.Printf("[INFO] 标签匹配成功: %s -> 分组 %s (ID: %d)\n", tagString, group.Name, group.ID)
			}
		} else {
			fmt.Printf("[INFO] 标签匹配成功: %s -> 根分组 %s (ID: %d)\n", tagString, group.Name, group.ID)
		}
		return group.ID, tagString
	}

	// 匹配失败,返回默认分组
	fmt.Printf("[WARN] 标签匹配失败: 未找到名为 '%s' 的分组,使用默认分组 (标签: %s)\n", bline, tagString)
	return defaultGroupID, tagString
}


// 创建腾讯云主机
func (s *CmdbHostCloudServiceImpl) CreateTencentHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "腾讯云主机导入请求已接收，正在扫描所有区域",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理导入
	go func() {
		// 初始化DAO
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.createTencentHosts(dto)
		if err != nil {
			fmt.Printf("[ERROR] 腾讯云主机导入失败: %v\n", err)
			return
		}

		// 打印所有要导入的主机信息
		fmt.Printf("[DEBUG] 准备导入 %d 台腾讯云主机\n", len(hosts))
		for i, host := range hosts {
			fmt.Printf("[DEBUG] 主机 %d: %+v\n", i+1, host)
		}

		// 批量导入 - 每批50台
		const batchSize = 50
		successCount := 0
		totalBatches := (len(hosts) + batchSize - 1) / batchSize

		for i := 0; i < len(hosts); i += batchSize {
			end := i + batchSize
			if end > len(hosts) {
				end = len(hosts)
			}
			batch := hosts[i:end]
			batchNum := i/batchSize + 1

			fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

			// 使用批量插入方法
			if err := hostDao.BatchCreateCmdbHosts(batch); err != nil {
				fmt.Printf("[ERROR] 第 %d 批主机保存失败: %v\n", batchNum, err)
				// 批量插入失败，尝试逐条插入以确定具体失败的记录
				for _, host := range batch {
					if err := hostDao.CreateCmdbHost(&host); err != nil {
						fmt.Printf("[ERROR] 腾讯云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
					} else {
						successCount++
						fmt.Printf("[INFO] 腾讯云主机导入成功: %s\n", host.InstanceID)
					}
				}
			} else {
				successCount += len(batch)
				fmt.Printf("[INFO] 第 %d 批导入成功: %d 台\n", batchNum, len(batch))
			}
		}

		// 统计导入成功的主机按区域分布
		successRegionStats := make(map[string]int)
		successHosts := 0
		for _, host := range hosts {
			// 检查主机是否导入成功（通过查询数据库）
			if existing, err := hostDao.GetCmdbHostByInstanceID(host.InstanceID); err == nil && existing != nil {
				successRegionStats[host.Region]++
				successHosts++
			}
		}

		// 记录最终结果
		fmt.Printf("[INFO] 腾讯云主机导入完成: 成功 %d 台, 失败 %d 台\n",
			successCount, len(hosts)-successCount)
		if successHosts > 0 {
			fmt.Println("[INFO] 导入成功的主机区域分布:")
			for region, count := range successRegionStats {
				fmt.Printf("  - %s: %d 台\n", region, count)
			}
		}

		// 这里可以添加将结果写入Redis或数据库的逻辑，供前端查询
	}()
}

func (s *CmdbHostCloudServiceImpl) createTencentHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
	// 腾讯云服务自动扫描所有区域
	tencentService := util.NewTencentCloudService(dto.AccessKey, dto.SecretKey)

	instances, err := tencentService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取腾讯云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何腾讯云实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台腾讯云实例\n", len(instances))

	// 统计每个区域的实例数量
	regionStats := make(map[string]int)
	for _, instance := range instances {
		regionStats[instance.Region]++
	}
	fmt.Println("[INFO] 区域分布统计:")
	for region, count := range regionStats {
		fmt.Printf("  - %s: %d 台\n", region, count)
	}

	// 初始化分组DAO
	groupDao := cmdbDao.NewCmdbGroupDao()

	// 初始化主机DAO用于去重检查
	hostDao := cmdbDao.NewCmdbHostDao()

	// 处理所有实例
	var hosts []model.CmdbHost
	skippedCount := 0
	duplicateCount := 0
	updatedCount := 0
	for _, instance := range instances {
		// 过滤非Running状态的实例(STOPPED/EXPIRED等)
		if instance.Status != "RUNNING" {
			fmt.Printf("[INFO] 跳过非运行状态实例: %s (状态: %s)\\n", instance.InstanceID, instance.Status)
			skippedCount++
			continue
		}

		// 检查实例是否已存在（根据InstanceID去重）
		if existing, _ := hostDao.GetCmdbHostByInstanceID(instance.InstanceID); existing != nil {
			// 检查配置是否发生变化
			needUpdate := false

			// 计算总磁盘
			totalDisk := instance.SystemDisk
			for _, disk := range instance.DataDisks {
				totalDisk += disk
			}

			// 比对CPU
			newCPU := strconv.Itoa(instance.CPU)
			if existing.CPU != newCPU {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机CPU变更: %s (%s -> %s)\n", instance.InstanceID, existing.CPU, newCPU)
			}

			// 比对内存
			newMemory := strconv.Itoa(instance.Memory)
			if existing.Memory != newMemory {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机内存变更: %s (%s -> %s)\n", instance.InstanceID, existing.Memory, newMemory)
			}

			// 比对磁盘
			newDisk := strconv.Itoa(totalDisk)
			if existing.Disk != newDisk {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机磁盘变更: %s (%s -> %s)\n", instance.InstanceID, existing.Disk, newDisk)
			}

			// 比对IP地址
			newPrivateIP := strings.Join(instance.PrivateIPs, ",")
			if existing.PrivateIP != newPrivateIP {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机内网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PrivateIP, newPrivateIP)
			}

			newPublicIP := strings.Join(instance.PublicIPs, ",")
			if existing.PublicIP != newPublicIP {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机公网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PublicIP, newPublicIP)
			}

			// 比对SSH IP
			sshIP := ""
			if len(instance.PublicIPs) > 0 {
				sshIP = instance.PublicIPs[0]
			} else if len(instance.PrivateIPs) > 0 {
				sshIP = instance.PrivateIPs[0]
			}
			if sshIP != "" && existing.SSHIP != sshIP {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机SSH IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.SSHIP, sshIP)
			}

			// 比对操作系统
			if existing.OS != instance.OSName {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机OS变更: %s (%s -> %s)\n", instance.InstanceID, existing.OS, instance.OSName)
			}

			// 比对区域
			if existing.Region != instance.Region {
				needUpdate = true
				fmt.Printf("[INFO] 腾讯云主机区域变更: %s (%s -> %s)\n", instance.InstanceID, existing.Region, instance.Region)
			}

			// 如果有变更，执行更新
			if needUpdate {
				if err := hostDao.UpdateCmdbHost(existing.ID, &model.CmdbHost{
					CPU:       newCPU,
					Memory:    newMemory,
					Disk:      newDisk,
					PrivateIP: newPrivateIP,
					PublicIP:  newPublicIP,
					SSHIP:     sshIP,
					OS:        instance.OSName,
					Region:    instance.Region,
				}); err != nil {
					fmt.Printf("[ERROR] 腾讯云主机配置更新失败: %s - %v\n", instance.InstanceID, err)
				} else {
					updatedCount++
					fmt.Printf("[SUCCESS] 腾讯云主机配置已更新: %s\n", instance.InstanceID)
				}
			} else {
				fmt.Printf("[INFO] 腾讯云主机配置无变化，跳过: %s (名称: %s, 区域: %s)\\n", instance.InstanceID, instance.InstanceName, instance.Region)
			}

			duplicateCount++
			continue
		}

		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			fmt.Printf("[WARN] 实例无有效IP，跳过: %s\\n", instance.InstanceID)
			skippedCount++
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

		// 解析标签并匹配分组
		groupID, tagString := s.matchGroupByTags(instance.Tags, &groupDao)

		// 使用实例名称作为主机名,如果实例名称为空则使用实例ID
		hostName := instance.InstanceName
		if hostName == "" {
			hostName = instance.InstanceID
		}

		// 确定SSH端口，优先使用用户配置的端口，否则默认22
		sshPort := 22
		if dto.SSHPort > 0 && dto.SSHPort <= 65535 {
			sshPort = dto.SSHPort
		}

		host := model.CmdbHost{
			HostName:    hostName,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     sshPort,
			SSHName:     "root",
			Vendor:      3,
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2,
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.Memory),
			Disk:        strconv.Itoa(totalDisk),
			BillingType: "包年包月",
			Tag:         tagString,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if skippedCount > 0 || duplicateCount > 0 || updatedCount > 0 {
		fmt.Printf("[INFO] 共跳过 %d 台非运行状态或无效实例, %d 台重复实例, %d 台配置已更新\\n", skippedCount, duplicateCount, updatedCount)
	}

	if len(hosts) == 0 {
		if duplicateCount > 0 {
			if updatedCount > 0 {
				return nil, fmt.Errorf("所有实例均已导入，无新增实例（已更新 %d 台主机配置）", updatedCount)
			}
			return nil, errors.New("所有实例均已导入，无新增实例")
		}
		return nil, errors.New("没有处于运行状态的有效实例可导入")
	}
	return hosts, nil
}

// 创建百度云主机
func (s *CmdbHostCloudServiceImpl) CreateBaiduHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "百度云主机导入请求已接收，正在扫描所有区域",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理导入
	go func() {
		// 初始化DAO
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.createBaiduHosts(dto)
		if err != nil {
			fmt.Printf("[ERROR] 百度云主机导入失败: %v\n", err)
			return
		}

		// 打印所有要导入的主机信息
		fmt.Printf("[DEBUG] 准备导入 %d 台百度云主机\n", len(hosts))
		for i, host := range hosts {
			fmt.Printf("[DEBUG] 主机 %d: %+v\n", i+1, host)
		}

		// 批量导入（自动去重）
		successCount := 0
		duplicateCount := 0
		for _, host := range hosts {
			// 检查是否已存在（根据InstanceID去重）
			existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
			if existing != nil {
				fmt.Printf("[INFO] 百度云主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
				duplicateCount++
				continue
			}

			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 百度云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 百度云主机导入成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		// 记录最终结果
		fmt.Printf("[INFO] 百度云主机导入完成: 成功 %d 台, 重复 %d 台, 失败 %d 台\n",
			successCount, duplicateCount, len(hosts)-successCount-duplicateCount)
	}()
}

func (s *CmdbHostCloudServiceImpl) createBaiduHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
	baiduService := util.NewBaiduCloudService(dto.AccessKey, dto.SecretKey)

	// 获取所有区域的实例（自动扫描所有区域）
	instances, err := baiduService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取百度云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何百度云实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台百度云实例\n", len(instances))

	// 统计每个区域的实例数量
	regionStats := make(map[string]int)
	for _, instance := range instances {
		regionStats[instance.Region]++
	}
	fmt.Println("[INFO] 区域分布统计:")
	for region, count := range regionStats {
		fmt.Printf("  - %s: %d 台\n", region, count)
	}

	// 初始化主机DAO用于去重检查
	hostDao := cmdbDao.NewCmdbHostDao()

	var hosts []model.CmdbHost
	skippedCount := 0
	duplicateCount := 0
	updatedCount := 0
	for _, instance := range instances {
		// 过滤非Running状态的实例
		if instance.Status != "Running" && instance.Status != "running" {
			fmt.Printf("[INFO] 跳过非运行状态实例: %s (状态: %s)\n", instance.InstanceID, instance.Status)
			skippedCount++
			continue
		}

		// 检查实例是否已存在（根据InstanceID去重）
		if existing, _ := hostDao.GetCmdbHostByInstanceID(instance.InstanceID); existing != nil {
			// 检查配置是否发生变化
			needUpdate := false

			// 计算总磁盘
			totalDisk := instance.SystemDisk
			for _, disk := range instance.DataDisks {
				totalDisk += disk
			}

			// 比对CPU
			newCPU := strconv.Itoa(instance.CPU)
			if existing.CPU != newCPU {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机CPU变更: %s (%s -> %s)\n", instance.InstanceID, existing.CPU, newCPU)
			}

			// 比对内存
			newMemory := strconv.Itoa(instance.Memory)
			if existing.Memory != newMemory {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机内存变更: %s (%s -> %s)\n", instance.InstanceID, existing.Memory, newMemory)
			}

			// 比对磁盘
			newDisk := strconv.Itoa(totalDisk)
			if existing.Disk != newDisk {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机磁盘变更: %s (%s -> %s)\n", instance.InstanceID, existing.Disk, newDisk)
			}

			// 比对IP地址
			newPrivateIP := strings.Join(instance.PrivateIPs, ",")
			if existing.PrivateIP != newPrivateIP {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机内网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PrivateIP, newPrivateIP)
			}

			newPublicIP := strings.Join(instance.PublicIPs, ",")
			if existing.PublicIP != newPublicIP {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机公网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PublicIP, newPublicIP)
			}

			// 比对SSH IP
			sshIP := ""
			if len(instance.PublicIPs) > 0 {
				sshIP = instance.PublicIPs[0]
			} else if len(instance.PrivateIPs) > 0 {
				sshIP = instance.PrivateIPs[0]
			}
			if sshIP != "" && existing.SSHIP != sshIP {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机SSH IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.SSHIP, sshIP)
			}

			// 比对操作系统
			if existing.OS != instance.OSName {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机OS变更: %s (%s -> %s)\n", instance.InstanceID, existing.OS, instance.OSName)
			}

			// 比对区域
			if existing.Region != instance.Region {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机区域变更: %s (%s -> %s)\n", instance.InstanceID, existing.Region, instance.Region)
			}

			// 比对计费类型
			if existing.BillingType != instance.BillingType {
				needUpdate = true
				fmt.Printf("[INFO] 百度云主机计费类型变更: %s (%s -> %s)\n", instance.InstanceID, existing.BillingType, instance.BillingType)
			}

			// 如果有变更，执行更新
			if needUpdate {
				if err := hostDao.UpdateCmdbHost(existing.ID, &model.CmdbHost{
					CPU:         newCPU,
					Memory:      newMemory,
					Disk:        newDisk,
					PrivateIP:   newPrivateIP,
					PublicIP:    newPublicIP,
					SSHIP:       sshIP,
					OS:          instance.OSName,
					Region:      instance.Region,
					BillingType: instance.BillingType,
				}); err != nil {
					fmt.Printf("[ERROR] 百度云主机配置更新失败: %s - %v\n", instance.InstanceID, err)
				} else {
					updatedCount++
					fmt.Printf("[SUCCESS] 百度云主机配置已更新: %s\n", instance.InstanceID)
				}
			} else {
				fmt.Printf("[INFO] 百度云主机配置无变化，跳过: %s (名称: %s, 区域: %s)\n", instance.InstanceID, instance.InstanceName, instance.Region)
			}

			duplicateCount++
			continue
		}

		// 确定SSH连接IP
		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			fmt.Printf("[WARN] 百度云实例无有效IP，跳过: %s\n", instance.InstanceID)
			skippedCount++
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

		// 百度云暂不支持标签，使用默认分组或用户指定分组
		groupID := dto.GroupID
		if groupID == 0 {
			groupID = 1 // 默认分组
		}

		// 百度云使用实例名称作为主机名,如果实例名称为空则使用实例ID
		hostName := instance.InstanceName
		if hostName == "" {
			hostName = instance.InstanceID
		}

		// 确定SSH端口，优先使用用户配置的端口，否则默认22
		sshPort := 22
		if dto.SSHPort > 0 && dto.SSHPort <= 65535 {
			sshPort = dto.SSHPort
		}

		host := model.CmdbHost{
			HostName:    hostName,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     sshPort,
			SSHName:     "root",
			Vendor:      4, // 百度云厂商代码为4
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2, // 默认未认证状态
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.Memory),
			Disk:        strconv.Itoa(totalDisk),
			BillingType: instance.BillingType,
			Tag:         "", // 百度云暂不支持标签
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if skippedCount > 0 || duplicateCount > 0 || updatedCount > 0 {
		fmt.Printf("[INFO] 共跳过 %d 台非运行状态或无效实例, %d 台重复实例, %d 台配置已更新\n", skippedCount, duplicateCount, updatedCount)
	}

	if len(hosts) == 0 {
		if duplicateCount > 0 {
			if updatedCount > 0 {
				return nil, fmt.Errorf("所有实例均已导入，无新增实例（已更新 %d 台主机配置）", updatedCount)
			}
			return nil, errors.New("所有实例均已导入，无新增实例")
		}
		return nil, errors.New("没有处于运行状态的有效实例可导入")
	}

	fmt.Printf("[INFO] 成功解析 %d 台百度云主机\n", len(hosts))
	return hosts, nil
}

// 创建华为云主机
func (s *CmdbHostCloudServiceImpl) CreateHuaweiHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
    // 立即返回成功响应
    result.Success(c, gin.H{
        "status": "processing",
        "msg":    "华为云主机导入请求已接收，正在后台处理",
        "data": gin.H{
            "task_id":    time.Now().UnixNano(),
            "start_time": time.Now().Format(time.RFC3339),
        },
    })

    // 异步处理导入
    go func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Printf("[ERROR] 华为云主机导入发生panic: %v\n", r)
            }
        }()
        hostDao := cmdbDao.NewCmdbHostDao()

        hosts, err := s.createHuaweiHosts(dto)
        if err != nil {
            fmt.Printf("[ERROR] 华为云主机导入失败: %v\n", err)
            return
        }

        successCount := 0
        duplicateCount := 0
        for _, host := range hosts {
            // 根据实例ID去重
            existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
            if existing != nil {
                fmt.Printf("[INFO] 华为云主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
                duplicateCount++
                continue
            }
            if err := hostDao.CreateCmdbHost(&host); err != nil {
                fmt.Printf("[ERROR] 华为云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
                continue
            }
            successCount++
            fmt.Printf("[INFO] 华为云主机导入成功: %s (区域: %s)\n", host.InstanceID, host.Region)
        }

        fmt.Printf("[INFO] 华为云主机导入完成: 成功 %d 台, 重复 %d 台, 失败 %d 台\n",
            successCount, duplicateCount, len(hosts)-successCount-duplicateCount)
    }()
}

func (s *CmdbHostCloudServiceImpl) createHuaweiHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
    hService := util.NewHuaweiCloudService(dto.AccessKey, dto.SecretKey)

    instances, err := hService.GetInstances()
    if err != nil {
        return nil, fmt.Errorf("获取华为云实例失败: %v", err)
    }
    if len(instances) == 0 {
        return nil, errors.New("未找到任何华为云实例")
    }

    fmt.Printf("[INFO] 从所有区域共获取到 %d 台华为云实例\n", len(instances))

    // 统计每个区域的实例数量
    regionStats := make(map[string]int)
    for _, instance := range instances {
        regionStats[instance.Region]++
    }
    fmt.Println("[INFO] 区域分布统计:")
    for region, count := range regionStats {
        fmt.Printf("  - %s: %d 台\n", region, count)
    }

    // 初始化主机DAO用于去重检查
    hostDao := cmdbDao.NewCmdbHostDao()

    var hosts []model.CmdbHost
    skippedCount := 0
    duplicateCount := 0
    updatedCount := 0
    for _, instance := range instances {
        // 华为云API返回的实例都是有效的，暂不进行状态过滤（因为结构体中没有Status字段）

        // 检查实例是否已存在（根据InstanceID去重）
        if existing, _ := hostDao.GetCmdbHostByInstanceID(instance.InstanceID); existing != nil {
            // 检查配置是否发生变化
            needUpdate := false

            // 统计磁盘大小
            totalDisk := instance.SystemDisk
            for _, d := range instance.DataDisks {
                totalDisk += d
            }

            // 比对CPU
            newCPU := strconv.Itoa(instance.CPU)
            if existing.CPU != newCPU {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机CPU变更: %s (%s -> %s)\n", instance.InstanceID, existing.CPU, newCPU)
            }

            // 比对内存
            newMemory := strconv.Itoa(instance.MemoryGB)
            if existing.Memory != newMemory {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机内存变更: %s (%s -> %s)\n", instance.InstanceID, existing.Memory, newMemory)
            }

            // 比对磁盘
            newDisk := strconv.Itoa(totalDisk)
            if existing.Disk != newDisk {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机磁盘变更: %s (%s -> %s)\n", instance.InstanceID, existing.Disk, newDisk)
            }

            // 比对IP地址
            newPrivateIP := strings.Join(instance.PrivateIPs, ",")
            if existing.PrivateIP != newPrivateIP {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机内网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PrivateIP, newPrivateIP)
            }

            newPublicIP := strings.Join(instance.PublicIPs, ",")
            if existing.PublicIP != newPublicIP {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机公网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PublicIP, newPublicIP)
            }

            // 比对SSH IP
            sshIP := ""
            if len(instance.PublicIPs) > 0 {
                sshIP = instance.PublicIPs[0]
            } else if len(instance.PrivateIPs) > 0 {
                sshIP = instance.PrivateIPs[0]
            }
            if sshIP != "" && existing.SSHIP != sshIP {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机SSH IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.SSHIP, sshIP)
            }

            // 比对操作系统
            if existing.OS != instance.OSName {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机OS变更: %s (%s -> %s)\n", instance.InstanceID, existing.OS, instance.OSName)
            }

            // 比对区域
            if existing.Region != instance.Region {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机区域变更: %s (%s -> %s)\n", instance.InstanceID, existing.Region, instance.Region)
            }

            // 比对计费类型
            if existing.BillingType != instance.BillingType {
                needUpdate = true
                fmt.Printf("[INFO] 华为云主机计费类型变更: %s (%s -> %s)\n", instance.InstanceID, existing.BillingType, instance.BillingType)
            }

            // 如果有变更，执行更新
            if needUpdate {
                if err := hostDao.UpdateCmdbHost(existing.ID, &model.CmdbHost{
                    CPU:         newCPU,
                    Memory:      newMemory,
                    Disk:        newDisk,
                    PrivateIP:   newPrivateIP,
                    PublicIP:    newPublicIP,
                    SSHIP:       sshIP,
                    OS:          instance.OSName,
                    Region:      instance.Region,
                    BillingType: instance.BillingType,
                }); err != nil {
                    fmt.Printf("[ERROR] 华为云主机配置更新失败: %s - %v\n", instance.InstanceID, err)
                } else {
                    updatedCount++
                    fmt.Printf("[SUCCESS] 华为云主机配置已更新: %s\n", instance.InstanceID)
                }
            } else {
                fmt.Printf("[INFO] 华为云主机配置无变化，跳过: %s (名称: %s, 区域: %s)\n", instance.InstanceID, instance.InstanceName, instance.Region)
            }

            duplicateCount++
            continue
        }

        // 选择SSH优先级：公网 -> 私网
        sshIP := ""
        if len(instance.PublicIPs) > 0 {
            sshIP = instance.PublicIPs[0]
        } else if len(instance.PrivateIPs) > 0 {
            sshIP = instance.PrivateIPs[0]
        } else {
            fmt.Printf("[WARN] 华为云实例无有效IP，跳过: %s\n", instance.InstanceID)
            skippedCount++
            continue
        }

        // 统计磁盘大小
        totalDisk := instance.SystemDisk
        for _, d := range instance.DataDisks {
            totalDisk += d
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

        // 华为云暂不支持标签，使用默认分组或用户指定分组
        groupID := dto.GroupID
        if groupID == 0 {
            groupID = 1 // 默认分组
        }

        // 华为云使用实例名称作为主机名,如果实例名称为空则使用实例ID
        hostName := instance.InstanceName
        if hostName == "" {
            hostName = instance.InstanceID
        }

        // 确定SSH端口，优先使用用户配置的端口，否则默认22
        sshPort := 22
        if dto.SSHPort > 0 && dto.SSHPort <= 65535 {
            sshPort = dto.SSHPort
        }

        host := model.CmdbHost{
            HostName:    hostName,
            GroupID:     groupID,
            PrivateIP:   strings.Join(instance.PrivateIPs, ","),
            PublicIP:    strings.Join(instance.PublicIPs, ","),
            SSHIP:       sshIP,
            SSHPort:     sshPort,
            SSHName:     "root",
            Vendor:      5, // 华为云
            Region:      instance.Region,
            InstanceID:  instance.InstanceID,
            OS:          instance.OSName,
            Status:      2, // 默认未认证状态
            CPU:         strconv.Itoa(instance.CPU),
            Memory:      strconv.Itoa(instance.MemoryGB),
            Disk:        strconv.Itoa(totalDisk),
            BillingType: instance.BillingType,
            Tag:         "", // 华为云暂不支持标签
            CreateTime:  util.HTime{Time: createTime.UTC()},
            ExpireTime:  util.HTime{Time: expireTime.UTC()},
        }
        hosts = append(hosts, host)
    }

    if skippedCount > 0 || duplicateCount > 0 || updatedCount > 0 {
        fmt.Printf("[INFO] 共跳过 %d 台非运行状态或无效实例, %d 台重复实例, %d 台配置已更新\n", skippedCount, duplicateCount, updatedCount)
    }

    if len(hosts) == 0 {
        if duplicateCount > 0 {
            if updatedCount > 0 {
                return nil, fmt.Errorf("所有实例均已导入，无新增实例（已更新 %d 台主机配置）", updatedCount)
            }
            return nil, errors.New("所有实例均已导入，无新增实例")
        }
        return nil, errors.New("没有处于运行状态的有效实例可导入")
    }

    fmt.Printf("[INFO] 成功解析 %d 台华为云主机\n", len(hosts))
    return hosts, nil
}

// 创建AWS主机
func (s *CmdbHostCloudServiceImpl) CreateAWSHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "AWS主机导入请求已接收，正在扫描所有区域",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理导入
	go func() {
		// 初始化DAO
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.createAWSHosts(dto)
		if err != nil {
			fmt.Printf("[ERROR] AWS主机导入失败: %v\n", err)
			return
		}

		// 打印所有要导入的主机信息
		fmt.Printf("[DEBUG] 准备导入 %d 台AWS主机\n", len(hosts))
		for i, host := range hosts {
			fmt.Printf("[DEBUG] 主机 %d: %+v\n", i+1, host)
		}

		// 批量导入 - 每批50台
		const batchSize = 50
		successCount := 0
		totalBatches := (len(hosts) + batchSize - 1) / batchSize

		for i := 0; i < len(hosts); i += batchSize {
			end := i + batchSize
			if end > len(hosts) {
				end = len(hosts)
			}
			batch := hosts[i:end]
			batchNum := i/batchSize + 1

			fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

			// 使用批量插入方法
			if err := hostDao.BatchCreateCmdbHosts(batch); err != nil {
				fmt.Printf("[ERROR] 第 %d 批主机保存失败: %v\n", batchNum, err)
				// 批量插入失败，尝试逐条插入以确定具体失败的记录
				for _, host := range batch {
					if err := hostDao.CreateCmdbHost(&host); err != nil {
						fmt.Printf("[ERROR] AWS主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
					} else {
						successCount++
						fmt.Printf("[INFO] AWS主机导入成功: %s\n", host.InstanceID)
					}
				}
			} else {
				successCount += len(batch)
				fmt.Printf("[INFO] 第 %d 批导入成功: %d 台\n", batchNum, len(batch))
			}
		}

		// 统计导入成功的主机按区域分布
		successRegionStats := make(map[string]int)
		successHosts := 0
		for _, host := range hosts {
			// 检查主机是否导入成功（通过查询数据库）
			if existing, err := hostDao.GetCmdbHostByInstanceID(host.InstanceID); err == nil && existing != nil {
				successRegionStats[host.Region]++
				successHosts++
			}
		}

		// 记录最终结果
		fmt.Printf("[INFO] AWS主机导入完成: 成功 %d 台, 失败 %d 台\n",
			successCount, len(hosts)-successCount)
		if successHosts > 0 {
			fmt.Println("[INFO] 导入成功的主机区域分布:")
			for region, count := range successRegionStats {
				fmt.Printf("  - %s: %d 台\n", region, count)
			}
		}

		// 这里可以添加将结果写入Redis或数据库的逻辑，供前端查询
	}()
}

func (s *CmdbHostCloudServiceImpl) createAWSHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
	// AWS服务自动扫描所有区域
	awsService := util.NewAWSCloudService(dto.AccessKey, dto.SecretKey)

	instances, err := awsService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取AWS实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何AWS实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台AWS实例\n", len(instances))

	// 统计每个区域的实例数量
	regionStats := make(map[string]int)
	for _, instance := range instances {
		regionStats[instance.Region]++
	}
	fmt.Println("[INFO] 区域分布统计:")
	for region, count := range regionStats {
		fmt.Printf("  - %s: %d 台\n", region, count)
	}

	// 初始化分组DAO
	groupDao := cmdbDao.NewCmdbGroupDao()

	// 初始化主机DAO用于去重检查
	hostDao := cmdbDao.NewCmdbHostDao()

	// 处理所有实例
	var hosts []model.CmdbHost
	skippedCount := 0
	duplicateCount := 0
	updatedCount := 0
	for _, instance := range instances {
		// 过滤非running状态的实例(stopped/terminated等)
		if instance.Status != "running" {
			fmt.Printf("[INFO] 跳过非运行状态实例: %s (状态: %s)\n", instance.InstanceID, instance.Status)
			skippedCount++
			continue
		}

		// 检查实例是否已存在（根据InstanceID去重）
		if existing, _ := hostDao.GetCmdbHostByInstanceID(instance.InstanceID); existing != nil {
			// 检查配置是否发生变化
			needUpdate := false

			// 解析磁盘大小
			diskSize, _ := strconv.Atoi(instance.DiskGB)
			if diskSize == 0 {
				diskSize = 8 // 默认8GB
			}

			// 比对CPU
			newCPU := strconv.Itoa(instance.CPU)
			if existing.CPU != newCPU {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机CPU变更: %s (%s -> %s)\n", instance.InstanceID, existing.CPU, newCPU)
			}

			// 比对内存
			newMemory := strconv.Itoa(instance.Memory)
			if existing.Memory != newMemory {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机内存变更: %s (%s -> %s)\n", instance.InstanceID, existing.Memory, newMemory)
			}

			// 比对磁盘
			newDisk := strconv.Itoa(diskSize)
			if existing.Disk != newDisk {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机磁盘变更: %s (%s -> %s)\n", instance.InstanceID, existing.Disk, newDisk)
			}

			// 比对IP地址
			newPrivateIP := strings.Join(instance.PrivateIPs, ",")
			if existing.PrivateIP != newPrivateIP {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机内网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PrivateIP, newPrivateIP)
			}

			newPublicIP := strings.Join(instance.PublicIPs, ",")
			if existing.PublicIP != newPublicIP {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机公网IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.PublicIP, newPublicIP)
			}

			// 比对SSH IP
			sshIP := ""
			if len(instance.PublicIPs) > 0 {
				sshIP = instance.PublicIPs[0]
			} else if len(instance.PrivateIPs) > 0 {
				sshIP = instance.PrivateIPs[0]
			}
			if sshIP != "" && existing.SSHIP != sshIP {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机SSH IP变更: %s (%s -> %s)\n", instance.InstanceID, existing.SSHIP, sshIP)
			}

			// 比对操作系统
			if existing.OS != instance.OSName {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机OS变更: %s (%s -> %s)\n", instance.InstanceID, existing.OS, instance.OSName)
			}

			// 比对区域
			if existing.Region != instance.Region {
				needUpdate = true
				fmt.Printf("[INFO] AWS主机区域变更: %s (%s -> %s)\n", instance.InstanceID, existing.Region, instance.Region)
			}

			// 如果有变更，执行更新
			if needUpdate {
				if err := hostDao.UpdateCmdbHost(existing.ID, &model.CmdbHost{
					CPU:       newCPU,
					Memory:    newMemory,
					Disk:      newDisk,
					PrivateIP: newPrivateIP,
					PublicIP:  newPublicIP,
					SSHIP:     sshIP,
					OS:        instance.OSName,
					Region:    instance.Region,
				}); err != nil {
					fmt.Printf("[ERROR] AWS主机配置更新失败: %s - %v\n", instance.InstanceID, err)
				} else {
					updatedCount++
					fmt.Printf("[SUCCESS] AWS主机配置已更新: %s\n", instance.InstanceID)
				}
			} else {
				fmt.Printf("[INFO] AWS主机配置无变化，跳过: %s (名称: %s, 区域: %s)\n", instance.InstanceID, instance.InstanceName, instance.Region)
			}

			duplicateCount++
			continue
		}

		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			fmt.Printf("[WARN] 实例无有效IP，跳过: %s\n", instance.InstanceID)
			skippedCount++
			continue
		}

		// 解析磁盘大小
		diskSize, _ := strconv.Atoi(instance.DiskGB)
		if diskSize == 0 {
			diskSize = 8 // 默认8GB
		}

		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}

		// 解析标签并匹配分组
		groupID, tagString := s.matchGroupByTags(instance.Tags, &groupDao)

		// 使用实例名称作为主机名,如果实例名称为空则使用实例ID
		hostName := instance.InstanceName
		if hostName == "" {
			hostName = instance.InstanceID
		}

		// 确定SSH端口，优先使用用户配置的端口，否则默认22
		sshPort := 22
		if dto.SSHPort > 0 && dto.SSHPort <= 65535 {
			sshPort = dto.SSHPort
		}

		host := model.CmdbHost{
			HostName:    hostName,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     sshPort,
			SSHName:     "ec2-user", // AWS EC2默认用户名
			Vendor:      6,          // AWS厂商代码为6
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2, // 未认证状态
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.Memory),
			Disk:        strconv.Itoa(diskSize),
			BillingType: "按需付费",
			Tag:         tagString,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if skippedCount > 0 || duplicateCount > 0 || updatedCount > 0 {
		fmt.Printf("[INFO] 共跳过 %d 台非运行状态或无效实例, %d 台重复实例, %d 台配置已更新\n", skippedCount, duplicateCount, updatedCount)
	}

	if len(hosts) == 0 {
		if duplicateCount > 0 {
			if updatedCount > 0 {
				return nil, fmt.Errorf("所有实例均已导入，无新增实例（已更新 %d 台主机配置）", updatedCount)
			}
			return nil, errors.New("所有实例均已导入，无新增实例")
		}
		return nil, errors.New("没有处于运行状态的有效实例可导入")
	}
	return hosts, nil
}

func GetCmdbHostCloudService() CmdbHostCloudServiceInterface {
    return &CmdbHostCloudServiceImpl{
        hostDao: cmdbDao.NewCmdbHostDao(),
    }
}

// SyncAliyunHostBackground 后台同步阿里云主机（复用精细化逻辑，不依赖gin.Context）
func (s *CmdbHostCloudServiceImpl) SyncAliyunHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error) {
	result := &SyncResult{
		RegionStats: make(map[string]int),
	}

	// 初始化DAO
	hostDao := cmdbDao.NewCmdbHostDao()

	// 获取主机列表（复用精细化逻辑）
	hosts, err := s.createAliyunHosts(dto)
	if err != nil {
		return result, err
	}

	result.TotalFetched = len(hosts)
	fmt.Printf("[INFO] 准备导入 %d 台阿里云主机\n", len(hosts))

	// 批量导入 - 每批50台（精细化逻辑）
	const batchSize = 50
	successCount := 0
	totalBatches := (len(hosts) + batchSize - 1) / batchSize

	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}
		batch := hosts[i:end]
		batchNum := i/batchSize + 1

		fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

		// 使用批量插入方法
		if err := hostDao.BatchCreateCmdbHosts(batch); err != nil {
			fmt.Printf("[WARN] 第 %d 批主机保存失败: %v，尝试逐条插入\n", batchNum, err)
			// 批量插入失败，尝试逐条插入
			for _, host := range batch {
				if err := hostDao.CreateCmdbHost(&host); err != nil {
					fmt.Printf("[ERROR] 阿里云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
					result.FailedCount++
				} else {
					successCount++
					result.RegionStats[host.Region]++
					fmt.Printf("[INFO] 阿里云主机导入成功: %s\n", host.InstanceID)
				}
			}
		} else {
			successCount += len(batch)
			// 统计区域分布
			for _, host := range batch {
				result.RegionStats[host.Region]++
			}
			fmt.Printf("[INFO] 第 %d 批导入成功: %d 台\n", batchNum, len(batch))
		}
	}

	result.NewAdded = successCount
	result.SkippedCount = result.TotalFetched - result.NewAdded - result.FailedCount

	fmt.Printf("[INFO] 阿里云主机导入完成: 成功 %d 台, 失败 %d 台\n", successCount, result.FailedCount)
	if len(result.RegionStats) > 0 {
		fmt.Println("[INFO] 导入成功的主机区域分布:")
		for region, count := range result.RegionStats {
			fmt.Printf("  - %s: %d 台\n", region, count)
		}
	}

	return result, nil
}

// SyncTencentHostBackground 后台同步腾讯云主机（复用精细化逻辑，不依赖gin.Context）
func (s *CmdbHostCloudServiceImpl) SyncTencentHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error) {
	result := &SyncResult{
		RegionStats: make(map[string]int),
	}

	// 初始化DAO
	hostDao := cmdbDao.NewCmdbHostDao()

	// 获取主机列表（复用精细化逻辑）
	hosts, err := s.createTencentHosts(dto)
	if err != nil {
		return result, err
	}

	result.TotalFetched = len(hosts)
	fmt.Printf("[INFO] 准备导入 %d 台腾讯云主机\n", len(hosts))

	// 批量导入 - 每批50台（精细化逻辑）
	const batchSize = 50
	successCount := 0
	totalBatches := (len(hosts) + batchSize - 1) / batchSize

	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}
		batch := hosts[i:end]
		batchNum := i/batchSize + 1

		fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

		// 使用批量插入方法
		if err := hostDao.BatchCreateCmdbHosts(batch); err != nil {
			fmt.Printf("[WARN] 第 %d 批主机保存失败: %v，尝试逐条插入\n", batchNum, err)
			// 批量插入失败，尝试逐条插入
			for _, host := range batch {
				if err := hostDao.CreateCmdbHost(&host); err != nil {
					fmt.Printf("[ERROR] 腾讯云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
					result.FailedCount++
				} else {
					successCount++
					result.RegionStats[host.Region]++
					fmt.Printf("[INFO] 腾讯云主机导入成功: %s\n", host.InstanceID)
				}
			}
		} else {
			successCount += len(batch)
			// 统计区域分布
			for _, host := range batch {
				result.RegionStats[host.Region]++
			}
			fmt.Printf("[INFO] 第 %d 批导入成功: %d 台\n", batchNum, len(batch))
		}
	}

	result.NewAdded = successCount
	result.SkippedCount = result.TotalFetched - result.NewAdded - result.FailedCount

	fmt.Printf("[INFO] 腾讯云主机导入完成: 成功 %d 台, 失败 %d 台\n", successCount, result.FailedCount)
	if len(result.RegionStats) > 0 {
		fmt.Println("[INFO] 导入成功的主机区域分布:")
		for region, count := range result.RegionStats {
			fmt.Printf("  - %s: %d 台\n", region, count)
		}
	}

	return result, nil
}

// SyncBaiduHostBackground 后台同步百度云主机（复用精细化逻辑，不依赖gin.Context）
func (s *CmdbHostCloudServiceImpl) SyncBaiduHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error) {
	result := &SyncResult{
		RegionStats: make(map[string]int),
	}

	// 初始化DAO
	hostDao := cmdbDao.NewCmdbHostDao()

	// 获取主机列表（复用精细化逻辑）
	hosts, err := s.createBaiduHosts(dto)
	if err != nil {
		return result, err
	}

	result.TotalFetched = len(hosts)
	fmt.Printf("[INFO] 准备导入 %d 台百度云主机\n", len(hosts))

	// 批量导入 - 每批50台（精细化逻辑）
	const batchSize = 50
	successCount := 0
	totalBatches := (len(hosts) + batchSize - 1) / batchSize

	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}
		batch := hosts[i:end]
		batchNum := i/batchSize + 1

		fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

		for _, host := range batch {
			// 检查实例ID是否存在
			existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
			if existing != nil {
				fmt.Printf("[INFO] 百度云主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
				result.SkippedCount++
				continue
			}

			// 创建主机
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 百度云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				result.FailedCount++
				continue
			}
			successCount++
			fmt.Printf("[INFO] 百度云主机导入成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		// 仅统计成功导入的区域分布
		if batchNum == 1 || successCount > 0 {
			for _, host := range batch {
				result.RegionStats[host.Region]++
			}
			fmt.Printf("[INFO] 第 %d 批导入成功: %d 台\n", batchNum, len(batch))
		}
	}

	result.NewAdded = successCount
	result.SkippedCount = result.TotalFetched - result.NewAdded - result.FailedCount

	fmt.Printf("[INFO] 百度云主机导入完成: 成功 %d 台, 失败 %d 台\n", successCount, result.FailedCount)
	if len(result.RegionStats) > 0 {
		fmt.Println("[INFO] 导入成功的主机区域分布:")
		for region, count := range result.RegionStats {
			fmt.Printf("  - %s: %d 台\n", region, count)
		}
	}

	return result, nil
}

// SyncHuaweiHostBackground 后台同步华为云主机（复用精细化逻辑，不依赖gin.Context）
func (s *CmdbHostCloudServiceImpl) SyncHuaweiHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error) {
	result := &SyncResult{
		RegionStats: make(map[string]int),
	}

	// 初始化DAO
	hostDao := cmdbDao.NewCmdbHostDao()

	// 获取主机列表（复用精细化逻辑）
	hosts, err := s.createHuaweiHosts(dto)
	if err != nil {
		return result, err
	}

	result.TotalFetched = len(hosts)
	fmt.Printf("[INFO] 准备导入 %d 台华为云主机\n", len(hosts))

	// 批量导入 - 每批50台（精细化逻辑）
	const batchSize = 50
	successCount := 0
	totalBatches := (len(hosts) + batchSize - 1) / batchSize

	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}
		batch := hosts[i:end]
		batchNum := i/batchSize + 1

		fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

		for _, host := range batch {
			// 检查实例ID是否存在
			existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
			if existing != nil {
				fmt.Printf("[INFO] 华为云主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
				result.SkippedCount++
				continue
			}

			// 创建主机
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 华为云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				result.FailedCount++
				continue
			}
			successCount++
			fmt.Printf("[INFO] 华为云主机导入成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		// 仅统计成功导入的区域分布
		if batchNum == 1 || successCount > 0 {
			for _, host := range batch {
				result.RegionStats[host.Region]++
			}
			fmt.Printf("[INFO] 第 %d 批导入成功: %d 台\n", batchNum, len(batch))
		}
	}

	result.NewAdded = successCount
	result.SkippedCount = result.TotalFetched - result.NewAdded - result.FailedCount

	fmt.Printf("[INFO] 华为云主机导入完成: 成功 %d 台, 失败 %d 台\n", successCount, result.FailedCount)
	if len(result.RegionStats) > 0 {
		fmt.Println("[INFO] 导入成功的主机区域分布:")
		for region, count := range result.RegionStats {
			fmt.Printf("  - %s: %d 台\n", region, count)
		}
	}

	return result, nil
}

// SyncAWSHostBackground 后台同步AWS主机（复用精细化逻辑，不依赖gin.Context）
func (s *CmdbHostCloudServiceImpl) SyncAWSHostBackground(dto *model.CreateCmdbHostCloudDto) (*SyncResult, error) {
	result := &SyncResult{
		RegionStats: make(map[string]int),
	}

	// 初始化DAO
	hostDao := cmdbDao.NewCmdbHostDao()

	// 获取主机列表（复用精细化逻辑）
	hosts, err := s.createAWSHosts(dto)
	if err != nil {
		return result, err
	}

	result.TotalFetched = len(hosts)
	fmt.Printf("[INFO] 准备导入 %d 台AWS主机\n", len(hosts))

	// 批量导入 - 每批50台（精细化逻辑）
	const batchSize = 50
	successCount := 0
	totalBatches := (len(hosts) + batchSize - 1) / batchSize

	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}
		batch := hosts[i:end]
		batchNum := i/batchSize + 1

		fmt.Printf("[INFO] 正在导入第 %d/%d 批 (本批 %d 台)...\n", batchNum, totalBatches, len(batch))

		for _, host := range batch {
			// 检查实例ID是否存在
			existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
			if existing != nil {
				fmt.Printf("[INFO] AWS主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
				result.SkippedCount++
				continue
			}

			// 创建主机
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] AWS主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				result.FailedCount++
				continue
			}
			successCount++
			result.RegionStats[host.Region]++
			fmt.Printf("[INFO] AWS主机导入成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		fmt.Printf("[INFO] 第 %d 批导入完成\n", batchNum)
	}

	result.NewAdded = successCount
	result.SkippedCount = result.TotalFetched - result.NewAdded - result.FailedCount

	fmt.Printf("[INFO] AWS主机导入完成: 成功 %d 台, 失败 %d 台, 跳过 %d 台\n",
		successCount, result.FailedCount, result.SkippedCount)
	if len(result.RegionStats) > 0 {
		fmt.Println("[INFO] 导入成功的主机区域分布:")
		for region, count := range result.RegionStats {
			fmt.Printf("  - %s: %d 台\n", region, count)
		}
	}

	return result, nil
}
