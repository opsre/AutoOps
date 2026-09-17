package service

import (
	"fmt"
	cmdbDao "dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	configDao "dodevops-api/api/configcenter/dao"
	systemModel "dodevops-api/api/system/model"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"time"

	"github.com/gin-gonic/gin"
)

type CmdbHostServiceInterface interface {
	GetCmdbHostList(c *gin.Context)                                                                           // 获取主机列表
	GetCmdbHostListWithPage(c *gin.Context, page, pageSize int, sortBy, sortOrder string)                     // 获取主机列表(分页)
	GetCmdbHostById(c *gin.Context, id uint)                                                                  // 根据ID获取主机
	GetCmdbHostByName(c *gin.Context, name string)                                                            // 根据名称获取主机
	CreateCmdbHost(c *gin.Context, dto *model.CreateCmdbHostDto)                                              // 创建主机
	UpdateCmdbHost(c *gin.Context, id uint, dto *model.UpdateCmdbHostDto)                                     // 更新主机
	DeleteCmdbHost(c *gin.Context, id uint)                                                                   // 删除主机
	GetCmdbHostsByGroupId(c *gin.Context, groupId uint, page, pageSize int)                                  // 根据分组ID获取主机列表(分页)
	GetCmdbHostsByHostNameLike(c *gin.Context, name string)                                                   // 根据主机名称模糊查询
	GetCmdbHostsByIP(c *gin.Context, ip string)                                                               // 根据IP查询(内网/公网/SSH)
	GetCmdbHostsByStatus(c *gin.Context, status int)                                                          // 根据状态查询
	ImportHostsFromExcel(c *gin.Context, dto *model.ImportHostsFromExcelDto, hosts []model.ExcelHostTemplate) // 从Excel导入主机
	SyncHostInfo(c *gin.Context, id uint)                                                                     // 同步主机基本信息
	GetCloudHostStatistics(c *gin.Context)                                                                    // 获取云主机统计信息
	BatchUpdateHostGroup(c *gin.Context, dto *model.BatchUpdateHostGroupDto)                                  // 批量更新主机分组
	BatchDeleteHost(c *gin.Context, dto *model.BatchDeleteHostDto)                                            // 批量删除主机
}

type CmdbHostServiceImpl struct {
	dao      cmdbDao.CmdbHostDao
	groupDao cmdbDao.CmdbGroupDao
}

// 从Excel导入主机
func (s *CmdbHostServiceImpl) ImportHostsFromExcel(c *gin.Context, dto *model.ImportHostsFromExcelDto, hosts []model.ExcelHostTemplate) {
	// 检查分组是否存在
	group, err := s.groupDao.GetCmdbGroupById(dto.GroupID)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "分组不存在")
		return
	}
	fmt.Printf("导入主机到分组: ID=%d, Name=%s\n", group.ID, group.Name)

	// 批量创建主机
	var successCount, failCount int
	var failedHosts []string
	type hostResult struct {
		hostID uint
		err    error
	}
	resultChan := make(chan hostResult, len(hosts))

	for i, excelHost := range hosts {
		fmt.Printf("处理第%d条主机记录: %+v\n", i+1, excelHost)

		// 检查SSH IP是否已存在
		if existHost := s.dao.GetCmdbHostBySSHIP(excelHost.SSHIP); existHost != nil {
			fmt.Printf("警告: SSH地址 %s 已存在(主机别名: %s)，跳过导入\n", excelHost.SSHIP, existHost.HostName)
			failedHosts = append(failedHosts, fmt.Sprintf("%s(IP重复)", excelHost.HostAlias))
			failCount++
			continue
		}

		// 创建主机
		host := model.CmdbHost{
			HostName:   excelHost.HostAlias, // 使用主机别名作为主机名称
			Name:       excelHost.HostAlias, // 默认name也使用主机别名，后续可通过同步更新
			GroupID:    dto.GroupID,
			SSHIP:      excelHost.SSHIP,
			SSHName:    excelHost.SSHName,
			SSHPort:    excelHost.SSHPort,
			Remark:     excelHost.Remark,
			CreateTime: util.HTime{Time: time.Now()},
			Vendor:     1, // 默认创建主机都是为自建主机
			Status:     2, // 初始状态设为未认证
		}

		// 保存主机信息
		fmt.Printf("尝试创建主机: %+v\n", host)
		if err := s.dao.CreateCmdbHost(&host); err != nil {
			fmt.Printf("创建主机失败: %v\n主机信息: %+v\n", err, host)
			failedHosts = append(failedHosts, excelHost.HostAlias)
			failCount++
			continue
		}
		fmt.Printf("主机创建成功, ID=%d\n", host.ID)
		successCount++

		// 异步执行SSH采集（使用Type=3公钥免认证）
		go func(h model.CmdbHost) {
			fmt.Printf("开始SSH采集(公钥免认证): %s@%s:%d\n", h.SSHName, h.SSHIP, h.SSHPort)

			// 准备SSH配置 - 使用Type=3公钥免认证
			sshConfig := util.SSHConfig{
				IP:       h.SSHIP,
				Port:     h.SSHPort,
				Type:     3, // 公钥免认证
				Username: h.SSHName,
			}

			// 获取系统信息
			sshUtil := util.NewSSHUtil()
			systemInfo, err := sshUtil.GetSystemInfo(&sshConfig)
			if err != nil {
				// 更新状态为认证失败
				fmt.Printf("主机ID %d SSH连接失败(公钥免认证): %v\n", h.ID, err)
				s.dao.UpdateCmdbHost(h.ID, &model.CmdbHost{Status: 3})
				resultChan <- hostResult{hostID: h.ID, err: err}
				return
			}

			// 更新主机信息
			updateData := model.CmdbHost{
				PrivateIP:  systemInfo["privateIp"],
				PublicIP:   systemInfo["publicIp"],
				Name:       systemInfo["name"],
				OS:         systemInfo["os"],
				CPU:        systemInfo["cpu"],
				Memory:     systemInfo["memory"],
				Disk:       systemInfo["disk"],
				Status:     1, // 认证成功
				UpdateTime: util.HTime{Time: time.Now()},
			}
			if err := s.dao.UpdateCmdbHost(h.ID, &updateData); err != nil {
				fmt.Printf("主机ID %d 更新信息失败: %v\n", h.ID, err)
				resultChan <- hostResult{hostID: h.ID, err: err}
				return
			}

			fmt.Printf("主机ID %d 信息采集成功\n", h.ID)
			resultChan <- hostResult{hostID: h.ID, err: nil}
		}(host)
	}

	// 等待所有goroutine完成
	go func() {
		for i := 0; i < successCount; i++ {
			res := <-resultChan
			if res.err != nil {
				fmt.Printf("主机ID %d 采集失败: %v\n", res.hostID, res.err)
			}
		}
		close(resultChan)
	}()

	responseData := gin.H{
		"success": successCount,
		"fail":    failCount,
		"total":   len(hosts),
		"message": "主机导入成功，正在使用公钥免认证方式采集主机信息...",
	}

	if len(failedHosts) > 0 {
		responseData["failedHosts"] = failedHosts
	}

	result.Success(c, responseData)
}

// 获取主机列表(分页) - 增加权限过滤
func (s *CmdbHostServiceImpl) GetCmdbHostListWithPage(c *gin.Context, page, pageSize int, sortBy, sortOrder string) {
	// 获取当前登录用户ID
	userObj, exists := c.Get(constant.ContextKeyUserObj)
	if !exists {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "未登录或登录已过期")
		return
	}
	admin, ok := userObj.(*systemModel.JwtAdmin)
	if !ok {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "无法获取用户信息")
		return
	}

	// admin用户跳过权限检查，可以查看所有主机
	if admin.Username == "admin" || admin.ID == 1 {
		list, total := s.dao.GetCmdbHostListWithPage(page, pageSize, sortBy, sortOrder)
		var vos []model.CmdbHostVo
		for _, host := range list {
			vos = append(vos, model.CmdbHostVo{
				ID:          host.ID,
				HostName:    host.HostName,
				Name:        host.Name,
				GroupID:     host.GroupID,
				GroupName:   host.Group.Name,
				PrivateIP:   host.PrivateIP,
				PublicIP:    host.PublicIP,
				SSHIP:       host.SSHIP,
				SSHName:     host.SSHName,
				SSHKeyID:    host.SSHKeyID,
				SSHPort:     host.SSHPort,
				Remark:      host.Remark,
				Tag:         host.Tag,
				Vendor:      host.Vendor,
				Region:      host.Region,
				InstanceID:  host.InstanceID,
				OS:          host.OS,
				Status:      host.Status,
				CPU:         host.CPU,
				Memory:      host.Memory,
				Disk:        host.Disk,
				BillingType: host.BillingType,
				CreateTime:  host.CreateTime,
				ExpireTime:  host.ExpireTime,
				UpdateTime:  host.UpdateTime,
			})
		}
		result.SuccessWithPage(c, vos, total, page, pageSize)
		return
	}

	// 获取用户有权限的主机ID列表
	authDao := cmdbDao.NewBastionHostAuthorizationDao()
	authorizedHostIds, err := authDao.GetUserAuthorizedHostIds(admin.ID)
	if err != nil {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "获取授权信息失败")
		return
	}

	// 如果用户没有任何授权，返回空列表
	if len(authorizedHostIds) == 0 {
		result.SuccessWithPage(c, []model.CmdbHostVo{}, 0, page, pageSize)
		return
	}

	// 使用授权的主机ID过滤主机列表
	list, total := s.dao.GetCmdbHostListWithPageByIds(page, pageSize, sortBy, sortOrder, authorizedHostIds)

	var vos []model.CmdbHostVo
	for _, host := range list {
		// 使用预加载的Group，避免N+1查询问题
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.SuccessWithPage(c, vos, total, page, pageSize)
}

// 获取主机列表
func (s *CmdbHostServiceImpl) GetCmdbHostList(c *gin.Context) {
	list := s.dao.GetCmdbHostList()
	var vos []model.CmdbHostVo
	for _, host := range list {
		// 使用预加载的Group，避免N+1查询问题
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

// 创建主机
func (s *CmdbHostServiceImpl) CreateCmdbHost(c *gin.Context, dto *model.CreateCmdbHostDto) {
	// 检查名称是否已存在
	if s.dao.CheckNameExists(dto.HostName) {
		result.FailedWithCode(c, constant.CMDB_HOST_NAME_EXISTS, "主机名称已存在")
		return
	}

	// 获取SSH凭据 (前端已确保SSHKeyID有效)
	authDao := configDao.NewEcsAuthDao()
	auth, _ := authDao.GetEcsAuthById(dto.SSHKeyID)

	// 初始保存连接信息
	host := model.CmdbHost{
		HostName:   dto.HostName,
		GroupID:    dto.GroupID,
		SSHIP:      dto.SSHIP,
		SSHName:    dto.SSHName,
		SSHKeyID:   dto.SSHKeyID,
		SSHPort:    dto.SSHPort,
		Remark:     dto.Remark,
		Tag:        dto.Tag,
		CreateTime: util.HTime{Time: time.Now()},
		Vendor:     1, // 默认创建主机都是为自建主机
		Status:     2, // 初始状态设为未认证
	}

	// 先保存基本信息
	if err := s.dao.CreateCmdbHost(&host); err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_CREATE_FAILED, err.Error())
		return
	}

	// 立即返回成功响应，后台异步执行SSH操作
	go func() {
		fmt.Println("开始尝试SSH连接获取系统信息...")

		sshUtil := util.NewSSHUtil()
		var systemInfo map[string]string
		var err error

		// 根据认证类型处理SSH连接
		if auth.Type == 4 {
			// 跳板机认证：先连接跳板机，再连接目标主机
			fmt.Println("使用跳板机认证方式")

			// 创建跳板机认证配置
			bastionAuth := &util.SSHConfig{
				IP:        auth.BastionHost,
				Port:      auth.BastionPort,
				Username:  auth.BastionUsername,
				Type:      auth.BastionAuthType,
				Password:  auth.BastionPassword,
				PublicKey: auth.BastionPrivateKey,
			}

			// 创建目标主机认证配置
			targetAuth := &util.SSHConfig{
				IP:        dto.SSHIP,
				Port:      dto.SSHPort,
				Username:  dto.SSHName,
				Type:      auth.TargetAuthType,
				Password:  auth.TargetPassword,
				PublicKey: auth.TargetPrivateKey,
			}

			fmt.Printf("跳板机配置: %+v\n", bastionAuth)
			fmt.Printf("目标主机配置: %+v\n", targetAuth)

			// 通过跳板机建立SSH客户端
			client, err := sshUtil.CreateBastionSSHClient(bastionAuth, targetAuth)
			if err != nil {
				fmt.Printf("通过跳板机连接失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}
			defer client.Close()

			// 在已建立的SSH客户端上执行系统信息采集命令
			output, err := util.ExecuteSSHCommandWithOutput(client, sshUtil.GetScriptContent())
			if err != nil {
				fmt.Printf("执行系统信息采集失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}

			// 解析输出
			systemInfo, err = util.ParseSystemInfo(output)
			if err != nil {
				fmt.Printf("解析系统信息失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}
		} else {
			// 普通认证方式（密码/私钥/免密）
			fmt.Printf("使用普通认证方式: type=%d\n", auth.Type)

			sshConfig := &util.SSHConfig{
				IP:        dto.SSHIP,
				Port:      dto.SSHPort,
				Type:      auth.Type,
				Username:  dto.SSHName,
				Password:  auth.Password,
				PublicKey: auth.PublicKey,
			}

			fmt.Printf("SSH配置: %+v\n", sshConfig)

			systemInfo, err = sshUtil.GetSystemInfo(sshConfig)
			if err != nil {
				fmt.Printf("SSH获取系统信息失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}
		}

		fmt.Printf("成功获取系统信息: %+v\n", systemInfo)

		// 验证必要字段是否存在
		if systemInfo["privateIp"] == "" || systemInfo["os"] == "" {
			fmt.Println("警告: 获取的系统信息不完整")
		}

		// 更新主机信息
		updateData := model.CmdbHost{
			PrivateIP:  systemInfo["privateIp"],
			PublicIP:   systemInfo["publicIp"],
			Name:       systemInfo["name"], // 添加name字段
			OS:         systemInfo["os"],
			CPU:        systemInfo["cpu"],
			Memory:     systemInfo["memory"],
			Disk:       systemInfo["disk"],
			Status:     1, // 认证成功
			UpdateTime: util.HTime{Time: time.Now()},
		}
		s.dao.UpdateCmdbHost(host.ID, &updateData)
	}()

	// 返回成功响应，前端可以通过轮询获取最新状态
	result.Success(c, gin.H{
		"id":     host.ID,
		"status": host.Status,
		"msg":    "主机创建成功，系统信息收集中...",
	})
}

// 更新主机
func (s *CmdbHostServiceImpl) UpdateCmdbHost(c *gin.Context, id uint, dto *model.UpdateCmdbHostDto) {
	// 不再需要查询认证凭证信息，直接从dto获取SSHName和SSHPort

	host := model.CmdbHost{
		HostName:   dto.HostName,
		GroupID:    dto.GroupID,
		SSHIP:      dto.SSHIP,
		SSHName:    dto.SSHName,
		SSHKeyID:   dto.SSHKeyID,
		SSHPort:    dto.SSHPort,
		Vendor:     dto.Vendor,
		Remark:     dto.Remark,
		Tag:        dto.Tag,
		UpdateTime: util.HTime{Time: time.Now()}, // 设置更新时间
	}
	err := s.dao.UpdateCmdbHost(id, &host)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_UPDATE_FAILED, err.Error())
		return
	}
	result.Success(c, true)
}

// 根据ID获取主机
func (s *CmdbHostServiceImpl) GetCmdbHostById(c *gin.Context, id uint) {
	host, err := s.dao.GetCmdbHostById(id)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "主机不存在")
		return
	}

	group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
	vo := model.CmdbHostVo{
		ID:          host.ID,
		HostName:    host.HostName,
		Name:        host.Name,
		GroupID:     host.GroupID,
		GroupName:   group.Name,
		PrivateIP:   host.PrivateIP,
		PublicIP:    host.PublicIP,
		SSHIP:       host.SSHIP,
		SSHName:     host.SSHName,
		SSHKeyID:    host.SSHKeyID,
		SSHPort:     host.SSHPort,
		Remark:      host.Remark,
		Tag:         host.Tag,
		Vendor:      host.Vendor,
		Region:      host.Region,
		InstanceID:  host.InstanceID,
		OS:          host.OS,
		Status:      host.Status,
		CPU:         host.CPU,
		Memory:      host.Memory,
		Disk:        host.Disk,
		BillingType: host.BillingType,
		CreateTime:  host.CreateTime,
		ExpireTime:  host.ExpireTime,
		UpdateTime:  host.UpdateTime,
	}
	result.Success(c, vo)
}

// 根据名称获取主机
func (s *CmdbHostServiceImpl) GetCmdbHostByName(c *gin.Context, name string) {
	host, err := s.dao.GetCmdbHostByName(name)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "主机不存在")
		return
	}

	group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
	vo := model.CmdbHostVo{
		ID:          host.ID,
		HostName:    host.HostName,
		Name:        host.Name,
		GroupID:     host.GroupID,
		GroupName:   group.Name,
		PrivateIP:   host.PrivateIP,
		PublicIP:    host.PublicIP,
		SSHIP:       host.SSHIP,
		SSHName:     host.SSHName,
		SSHKeyID:    host.SSHKeyID,
		SSHPort:     host.SSHPort,
		Remark:      host.Remark,
		Tag:         host.Tag,
		Vendor:      host.Vendor,
		Region:      host.Region,
		InstanceID:  host.InstanceID,
		OS:          host.OS,
		Status:      host.Status,
		CPU:         host.CPU,
		Memory:      host.Memory,
		Disk:        host.Disk,
		BillingType: host.BillingType,
		CreateTime:  host.CreateTime,
		ExpireTime:  host.ExpireTime,
		UpdateTime:  host.UpdateTime,
	}
	result.Success(c, vo)
}

// 删除主机
func (s *CmdbHostServiceImpl) DeleteCmdbHost(c *gin.Context, id uint) {
	err := s.dao.DeleteCmdbHost(id)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_DELETE_FAILED, err.Error())
		return
	}
	result.Success(c, true)
}

// 根据分组ID获取主机列表(分页) - 增加权限过滤
func (s *CmdbHostServiceImpl) GetCmdbHostsByGroupId(c *gin.Context, groupId uint, page, pageSize int) {
	// 获取当前登录用户ID
	userObj, exists := c.Get(constant.ContextKeyUserObj)
	if !exists {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "未登录或登录已过期")
		return
	}
	admin, ok := userObj.(*systemModel.JwtAdmin)
	if !ok {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "无法获取用户信息")
		return
	}

	// admin用户跳过权限检查，可以查看所有主机
	if admin.Username == "admin" || admin.ID == 1 {
		list, total := s.dao.GetCmdbHostsByGroupIdWithPage(groupId, page, pageSize)
		var vos []model.CmdbHostVo
		for _, host := range list {
			vos = append(vos, model.CmdbHostVo{
				ID:          host.ID,
				HostName:    host.HostName,
				Name:        host.Name,
				GroupID:     host.GroupID,
				GroupName:   host.Group.Name,
				PrivateIP:   host.PrivateIP,
				PublicIP:    host.PublicIP,
				SSHIP:       host.SSHIP,
				SSHName:     host.SSHName,
				SSHKeyID:    host.SSHKeyID,
				SSHPort:     host.SSHPort,
				Remark:      host.Remark,
				Tag:         host.Tag,
				Vendor:      host.Vendor,
				Region:      host.Region,
				InstanceID:  host.InstanceID,
				OS:          host.OS,
				Status:      host.Status,
				CPU:         host.CPU,
				Memory:      host.Memory,
				Disk:        host.Disk,
				BillingType: host.BillingType,
				CreateTime:  host.CreateTime,
				ExpireTime:  host.ExpireTime,
				UpdateTime:  host.UpdateTime,
			})
		}
		result.SuccessWithPage(c, vos, total, page, pageSize)
		return
	}

	// 获取用户有权限的主机ID列表
	authDao := cmdbDao.NewBastionHostAuthorizationDao()
	authorizedHostIds, err := authDao.GetUserAuthorizedHostIds(admin.ID)
	if err != nil {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "获取授权信息失败")
		return
	}

	// 如果用户没有任何授权，返回空列表
	if len(authorizedHostIds) == 0 {
		result.SuccessWithPage(c, []model.CmdbHostVo{}, 0, page, pageSize)
		return
	}

	// 获取分组下的主机，并用授权列表过滤
	list, total := s.dao.GetCmdbHostsByGroupIdWithPageFiltered(groupId, page, pageSize, authorizedHostIds)

	var vos []model.CmdbHostVo
	for _, host := range list {
		// 使用预加载的Group，避免N+1查询问题
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.SuccessWithPage(c, vos, total, page, pageSize)
}

// 根据主机名称模糊查询
func (s *CmdbHostServiceImpl) GetCmdbHostsByHostNameLike(c *gin.Context, name string) {
	list := s.dao.GetCmdbHostsByHostNameLike(name)
	var vos []model.CmdbHostVo
	for _, host := range list {
		// 使用预加载的Group，避免N+1查询问题
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

// 根据IP查询(内网/公网/SSH)
func (s *CmdbHostServiceImpl) GetCmdbHostsByIP(c *gin.Context, ip string) {
	list := s.dao.GetCmdbHostsByIP(ip)
	var vos []model.CmdbHostVo
	for _, host := range list {
		// 使用预加载的Group，避免N+1查询问题
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

// 根据状态查询
func (s *CmdbHostServiceImpl) GetCmdbHostsByStatus(c *gin.Context, status int) {
	list := s.dao.GetCmdbHostsByStatus(status)
	var vos []model.CmdbHostVo
	for _, host := range list {
		// 使用预加载的Group，避免N+1查询问题
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

func GetCmdbHostService() CmdbHostServiceInterface {
	return &CmdbHostServiceImpl{
		dao:      cmdbDao.NewCmdbHostDao(),
		groupDao: cmdbDao.NewCmdbGroupDao(),
	}
}

// 同步主机基本信息
func (s *CmdbHostServiceImpl) SyncHostInfo(c *gin.Context, id uint) {
	// 1. 验证主机是否存在
	host, err := s.dao.GetCmdbHostById(id)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "主机不存在")
		return
	}

	// 2. 检查SSH连接信息是否完整
	if host.SSHIP == "" || host.SSHName == "" || host.SSHKeyID == 0 {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "主机SSH连接信息不完整，无法同步")
		return
	}

	// 3. 获取SSH认证信息
	authDao := configDao.NewEcsAuthDao()
	auth, err := authDao.GetEcsAuthById(host.SSHKeyID)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "SSH认证凭据不存在")
		return
	}

	// 4. 立即返回成功响应，后台异步执行同步操作
	result.Success(c, gin.H{
		"id":      host.ID,
		"status":  host.Status,
		"message": "开始同步主机信息，请稍后查看结果",
	})

	// 5. 异步执行同步操作
	go func() {
		fmt.Printf("开始同步主机信息: ID=%d, Name=%s\n", host.ID, host.HostName)

		// 设置状态为同步中（使用状态2表示同步中）
		s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 2})

		sshUtil := util.NewSSHUtil()
		var systemInfo map[string]string
		var err error

		// 根据认证类型处理SSH连接
		if auth.Type == 4 {
			// 跳板机认证：先连接跳板机，再连接目标主机
			fmt.Printf("使用跳板机认证同步: 跳板机=%s:%d, 目标=%s:%d\n",
				auth.BastionHost, auth.BastionPort, host.SSHIP, host.SSHPort)

			// 创建跳板机认证配置
			bastionAuth := &util.SSHConfig{
				IP:        auth.BastionHost,
				Port:      auth.BastionPort,
				Username:  auth.BastionUsername,
				Type:      auth.BastionAuthType,
				Password:  auth.BastionPassword,
				PublicKey: auth.BastionPrivateKey,
			}

			// 创建目标主机认证配置
			targetAuth := &util.SSHConfig{
				IP:        host.SSHIP,
				Port:      host.SSHPort,
				Username:  host.SSHName,
				Type:      auth.TargetAuthType,
				Password:  auth.TargetPassword,
				PublicKey: auth.TargetPrivateKey,
			}

			fmt.Printf("跳板机认证: type=%d, username=%s, hasPassword=%v, hasKey=%v\n",
				bastionAuth.Type, bastionAuth.Username,
				bastionAuth.Password != "", bastionAuth.PublicKey != "")
			fmt.Printf("目标主机认证: type=%d, username=%s, hasPassword=%v, hasKeyLen=%d\n",
				targetAuth.Type, targetAuth.Username,
				targetAuth.Password != "", len(targetAuth.PublicKey))

			// 通过跳板机建立SSH客户端
			client, err := sshUtil.CreateBastionSSHClient(bastionAuth, targetAuth)
			if err != nil {
				fmt.Printf("通过跳板机连接失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}
			defer client.Close()

			// 在已建立的SSH客户端上执行系统信息采集命令
			output, err := util.ExecuteSSHCommandWithOutput(client, sshUtil.GetScriptContent())
			if err != nil {
				fmt.Printf("执行系统信息采集失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}

			// 解析输出
			systemInfo, err = util.ParseSystemInfo(output)
			if err != nil {
				fmt.Printf("解析系统信息失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}
		} else {
			// 普通认证方式（密码/私钥/免密）
			fmt.Printf("使用普通认证同步: %s@%s:%d (type=%d)\n",
				host.SSHName, host.SSHIP, host.SSHPort, auth.Type)

			sshConfig := &util.SSHConfig{
				IP:        host.SSHIP,
				Port:      host.SSHPort,
				Type:      auth.Type,
				Username:  host.SSHName,
				Password:  auth.Password,
				PublicKey: auth.PublicKey,
			}

			systemInfo, err = sshUtil.GetSystemInfo(sshConfig)
			if err != nil {
				fmt.Printf("SSH获取系统信息失败: %v\n", err)
				s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
				return
			}
		}

		fmt.Printf("成功获取系统信息: %+v\n", systemInfo)

		// 更新主机信息
		updateData := model.CmdbHost{
			Name:       systemInfo["name"],      // 主机名称
			OS:         systemInfo["os"],        // 操作系统
			CPU:        systemInfo["cpu"],       // CPU
			Memory:     systemInfo["memory"],    // 内存
			Disk:       systemInfo["disk"],      // 磁盘
			PrivateIP:  systemInfo["privateIp"], // 内网IP
			PublicIP:   systemInfo["publicIp"],  // 公网IP
			Status:     1,                       // 同步成功
			UpdateTime: util.HTime{Time: time.Now()},
		}

		if err := s.dao.UpdateCmdbHost(host.ID, &updateData); err != nil {
			fmt.Printf("更新主机信息失败: %v\n", err)
			// 更新状态为同步失败
			s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
			return
		}

		fmt.Printf("主机信息同步完成: ID=%d\n", host.ID)
	}()
}

// 获取云主机统计信息
func (s *CmdbHostServiceImpl) GetCloudHostStatistics(c *gin.Context) {
	// 获取主机按厂商分组的统计数据(包括自建)
	stats := s.dao.GetCloudHostStatistics()

	// 获取主机总数
	totalCount := s.dao.GetCloudHostTotalCount()

	// 构建返回数据
	vendorNames := map[int]string{
		1: "自建",
		2: "阿里云",
		3: "腾讯云",
		4: "百度云",
		5: "华为云",
		6: "AWS",
	}

	// 构建详细统计列表(按类型统计)
	var details []map[string]interface{}
	for vendor := 1; vendor <= 6; vendor++ {
		count := stats[vendor]
		details = append(details, map[string]interface{}{
			"vendor":     vendor,
			"vendorName": vendorNames[vendor],
			"count":      count,
		})
	}

	result.Success(c, map[string]interface{}{
		"total":   totalCount,
		"details": details,
		"summary": map[string]interface{}{
			"selfBuilt": stats[1], // 自建
			"aliyun":    stats[2], // 阿里云
			"tencent":   stats[3], // 腾讯云
			"baidu":     stats[4], // 百度云
			"huawei":    stats[5], // 华为云
			"aws":       stats[6], // AWS
		},
	})
}

// 批量更新主机分组
func (s *CmdbHostServiceImpl) BatchUpdateHostGroup(c *gin.Context, dto *model.BatchUpdateHostGroupDto) {
	// 1. 验证分组是否存在
	group, err := s.groupDao.GetCmdbGroupById(dto.GroupID)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "分组不存在")
		return
	}

	// 2. 验证是否为二级分组（ParentID != 0）
	isSecondLevel, err := s.groupDao.IsSecondLevelGroup(dto.GroupID)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "分组验证失败")
		return
	}
	if !isSecondLevel {
		result.Failed(c, constant.INVALID_PARAMS, "只能选择二级分组（子分组）")
		return
	}

	// 3. 验证主机ID列表不为空
	if len(dto.HostIDs) == 0 {
		result.Failed(c, constant.INVALID_PARAMS, "主机ID列表不能为空")
		return
	}

	// 4. 验证主机是否存在
	hosts, err := s.dao.GetCmdbHostsByIds(dto.HostIDs)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "查询主机失败")
		return
	}
	if len(hosts) != len(dto.HostIDs) {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "部分主机不存在")
		return
	}

	// 5. 批量更新分组
	if err := s.dao.BatchUpdateHostGroup(dto.HostIDs, dto.GroupID); err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_UPDATE_FAILED, fmt.Sprintf("批量更新失败: %v", err))
		return
	}

	result.Success(c, map[string]interface{}{
		"updatedCount": len(dto.HostIDs),
		"groupId":      dto.GroupID,
		"groupName":    group.Name,
		"message":      fmt.Sprintf("成功更新 %d 台主机到分组 %s", len(dto.HostIDs), group.Name),
	})
}

// 批量删除主机
func (s *CmdbHostServiceImpl) BatchDeleteHost(c *gin.Context, dto *model.BatchDeleteHostDto) {
	// 1. 验证主机ID列表不为空
	if len(dto.HostIDs) == 0 {
		result.Failed(c, constant.INVALID_PARAMS, "主机ID列表不能为空")
		return
	}

	// 2. 验证主机是否存在
	hosts, err := s.dao.GetCmdbHostsByIds(dto.HostIDs)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "查询主机失败")
		return
	}
	if len(hosts) != len(dto.HostIDs) {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "部分主机不存在")
		return
	}

	// 3. 批量删除主机
	if err := s.dao.BatchDeleteHost(dto.HostIDs); err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_DELETE_FAILED, fmt.Sprintf("批量删除失败: %v", err))
		return
	}

	result.Success(c, map[string]interface{}{
		"deletedCount": len(dto.HostIDs),
		"message":      fmt.Sprintf("成功删除 %d 台主机", len(dto.HostIDs)),
	})
}
