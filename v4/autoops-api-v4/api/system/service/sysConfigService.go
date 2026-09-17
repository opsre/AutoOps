// 系统配置服务层
package service

import (
	"encoding/json"
	"fmt"

	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SysConfigService 系统配置服务接口
type SysConfigService interface {
	CreateConfig(c *gin.Context, dto *model.SysConfigDto)
	GetConfigByID(c *gin.Context, id uint)
	GetConfigByKey(c *gin.Context, configKey string)
	UpdateConfig(c *gin.Context, dto *model.SysConfigDto)
	DeleteConfig(c *gin.Context, id uint)
	GetConfigList(c *gin.Context, query *model.SysConfigQuery)
	UpdateConfigStatus(c *gin.Context, id uint, status int)

	// LDAP相关
	SaveLDAPConfig(c *gin.Context, dto *model.LDAPConfigDto)
	GetLDAPConfig(c *gin.Context)
	TestLDAPConnection(c *gin.Context)
	LDAPLogin(c *gin.Context, username, password string)
}

// SysConfigServiceImpl 系统配置服务实现
type SysConfigServiceImpl struct {
	configDao dao.SysConfigDao
}

// NewSysConfigService 创建系统配置服务实例
func NewSysConfigService() SysConfigService {
	return &SysConfigServiceImpl{
		configDao: dao.NewSysConfigDao(),
	}
}

// CreateConfig 创建系统配置
func (s *SysConfigServiceImpl) CreateConfig(c *gin.Context, dto *model.SysConfigDto) {
	// 检查配置键是否已存在
	existing, _ := s.configDao.GetByKey(dto.ConfigKey)
	if existing != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "配置键已存在")
		return
	}

	config := &model.SysConfig{
		ConfigKey:  dto.ConfigKey,
		ConfigType: dto.ConfigType,
		ConfigData: dto.ConfigData,
		Status:     1, // 默认启用
		Remark:     dto.Remark,
	}

	if err := s.configDao.Create(config); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建失败: %v", err))
		return
	}

	result.Success(c, config)
}

// GetConfigByID 根据ID获取系统配置
func (s *SysConfigServiceImpl) GetConfigByID(c *gin.Context, id uint) {
	config, err := s.configDao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "配置不存在")
		} else {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		}
		return
	}

	result.Success(c, config)
}

// GetConfigByKey 根据配置键获取系统配置
func (s *SysConfigServiceImpl) GetConfigByKey(c *gin.Context, configKey string) {
	config, err := s.configDao.GetByKey(configKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "配置不存在")
		} else {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		}
		return
	}

	result.Success(c, config)
}

// UpdateConfig 更新系统配置
func (s *SysConfigServiceImpl) UpdateConfig(c *gin.Context, dto *model.SysConfigDto) {
	// 检查配置是否存在
	config, err := s.configDao.GetByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "配置不存在")
		return
	}

	// 如果修改了配置键，检查新配置键是否已存在
	if dto.ConfigKey != config.ConfigKey {
		existing, _ := s.configDao.GetByKey(dto.ConfigKey)
		if existing != nil && existing.ID != dto.ID {
			result.Failed(c, int(result.ApiCode.FAILED), "配置键已存在")
			return
		}
	}

	// 更新字段
	config.ConfigKey = dto.ConfigKey
	config.ConfigType = dto.ConfigType
	config.ConfigData = dto.ConfigData
	config.Status = dto.Status
	config.Remark = dto.Remark

	if err := s.configDao.Update(config); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新失败: %v", err))
		return
	}

	result.Success(c, config)
}

// DeleteConfig 删除系统配置
func (s *SysConfigServiceImpl) DeleteConfig(c *gin.Context, id uint) {
	// 检查是否存在
	_, err := s.configDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "配置不存在")
		return
	}

	if err := s.configDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除失败: %v", err))
		return
	}

	result.Success(c, nil)
}

// GetConfigList 获取系统配置列表
func (s *SysConfigServiceImpl) GetConfigList(c *gin.Context, query *model.SysConfigQuery) {
	configs, total, err := s.configDao.GetList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	data := map[string]interface{}{
		"list":     configs,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}

	result.Success(c, data)
}

// UpdateConfigStatus 更新配置状态
func (s *SysConfigServiceImpl) UpdateConfigStatus(c *gin.Context, id uint, status int) {
	// 检查是否存在
	_, err := s.configDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "配置不存在")
		return
	}

	if err := s.configDao.UpdateStatus(id, status); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新状态失败: %v", err))
		return
	}

	result.Success(c, nil)
}

// SaveLDAPConfig 保存LDAP配置
func (s *SysConfigServiceImpl) SaveLDAPConfig(c *gin.Context, dto *model.LDAPConfigDto) {
	// 将LDAP配置转换为JSON
	configData, err := json.Marshal(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("配置数据序列化失败: %v", err))
		return
	}

	// 检查LDAP配置是否已存在
	existing, err := s.configDao.GetByKey("ldap")

	if err != nil && err != gorm.ErrRecordNotFound {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询配置失败: %v", err))
		return
	}

	if existing != nil {
		// 更新现有配置
		existing.ConfigData = string(configData)
		existing.Status = 1
		if dto.Remark != "" {
			existing.Remark = dto.Remark
		}

		if err := s.configDao.Update(existing); err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新LDAP配置失败: %v", err))
			return
		}

		result.Success(c, existing)
	} else {
		// 创建新配置
		config := &model.SysConfig{
			ConfigKey:  "ldap",
			ConfigType: "ldap",
			ConfigData: string(configData),
			Status:     1,
			Remark:     dto.Remark,
		}

		if err := s.configDao.Create(config); err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建LDAP配置失败: %v", err))
			return
		}

		result.Success(c, config)
	}
}

// GetLDAPConfig 获取LDAP配置
func (s *SysConfigServiceImpl) GetLDAPConfig(c *gin.Context) {
	config, err := s.configDao.GetByKey("ldap")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认配置
			defaultConfig := &model.LDAPConfig{
				Enable:          false,
				Host:            "",
				Port:            389,
				BaseDn:          "",
				BindUser:        "",
				BindPass:        "",
				AuthFilter:      "(&(cn=%s))",
				CoverAttributes: true,
				TLS:             false,
				StartTLS:        false,
				DefaultRoles:    []string{"Standard"},
				Attributes: model.LDAPAttributes{
					Nickname: "cn",
					Phone:    "mobile",
					Email:    "mail",
				},
			}
			result.Success(c, defaultConfig)
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	// 解析JSON配置
	var ldapConfig model.LDAPConfig
	if err := json.Unmarshal([]byte(config.ConfigData), &ldapConfig); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("解析配置失败: %v", err))
		return
	}

	result.Success(c, ldapConfig)
}

// TestLDAPConnection 测试LDAP连接
func (s *SysConfigServiceImpl) TestLDAPConnection(c *gin.Context) {
	config, err := s.configDao.GetByKey("ldap")
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "LDAP配置不存在，请先配置")
		return
	}

	// 解析JSON配置
	var ldapConfig model.LDAPConfig
	if err := json.Unmarshal([]byte(config.ConfigData), &ldapConfig); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("解析配置失败: %v", err))
		return
	}

	// 测试LDAP连接
	authConfig := &util.LDAPAuthConfig{
		Host:       ldapConfig.Host,
		Port:       ldapConfig.Port,
		BaseDn:     ldapConfig.BaseDn,
		BindUser:   ldapConfig.BindUser,
		BindPass:   ldapConfig.BindPass,
		AuthFilter: ldapConfig.AuthFilter,
		TLS:        ldapConfig.TLS,
		StartTLS:   ldapConfig.StartTLS,
	}

	if err := util.TestLDAPConnection(authConfig); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("LDAP连接测试失败: %v", err))
		return
	}

	result.Success(c, map[string]interface{}{
		"status":  "success",
		"message": "LDAP连接测试成功",
		"host":    ldapConfig.Host,
		"port":    ldapConfig.Port,
	})
}

// LDAPLogin handles LDAP authentication and user provisioning
func (s *SysConfigServiceImpl) LDAPLogin(c *gin.Context, username, password string) {
	config, err := s.configDao.GetByKey("ldap")
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "LDAP配置不存在，请联系管理员配置")
		return
	}

	var ldapConfig model.LDAPConfig
	if err := json.Unmarshal([]byte(config.ConfigData), &ldapConfig); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "LDAP配置解析失败")
		return
	}

	if !ldapConfig.Enable {
		result.Failed(c, int(result.ApiCode.FAILED), "LDAP认证未启用")
		return
	}

	authConfig := &util.LDAPAuthConfig{
		Host:       ldapConfig.Host,
		Port:       ldapConfig.Port,
		BaseDn:     ldapConfig.BaseDn,
		BindUser:   ldapConfig.BindUser,
		BindPass:   ldapConfig.BindPass,
		AuthFilter: ldapConfig.AuthFilter,
		TLS:        ldapConfig.TLS,
		StartTLS:   ldapConfig.StartTLS,
	}

	ip := c.ClientIP()
	ldapUser, err := util.LDAPAuth(authConfig, username, password)
	if err != nil {
		dao.CreateSysLoginInfo(username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), fmt.Sprintf("LDAP认证失败: %v", err), 2)
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("LDAP认证失败: %v", err))
		return
	}

	existingUser := (dao.SysAdminDao{}).GetSysAdminByUsername(username)
	var sysAdmin model.SysAdmin

	if existingUser.ID > 0 {
		if ldapConfig.CoverAttributes {
			if key := ldapConfig.Attributes.Nickname; key != "" {
				if val, ok := ldapUser.Attrs[key]; ok && val != "" {
					existingUser.Nickname = val
				}
			}
			if key := ldapConfig.Attributes.Email; key != "" {
				if val, ok := ldapUser.Attrs[key]; ok && val != "" {
					existingUser.Email = val
				}
			}
			if key := ldapConfig.Attributes.Phone; key != "" {
				if val, ok := ldapUser.Attrs[key]; ok && val != "" {
					existingUser.Phone = val
				}
			}
			(dao.SysAdminDao{}).UpdateSysAdminLdapInfo(existingUser)
		}
		sysAdmin = existingUser
	} else {
		newUser := model.SysAdmin{
			Username: username,
			Password: util.EncryptionMd5("ldap_user_" + username),
			Nickname: ldapUser.CN,
			Status:   1,
			PostId:   1,
			DeptId:   1,
		}

		if key := ldapConfig.Attributes.Nickname; key != "" {
			if val, ok := ldapUser.Attrs[key]; ok && val != "" {
				newUser.Nickname = val
			}
		}
		if key := ldapConfig.Attributes.Email; key != "" {
			if val, ok := ldapUser.Attrs[key]; ok && val != "" {
				newUser.Email = val
			}
		}
		if key := ldapConfig.Attributes.Phone; key != "" {
			if val, ok := ldapUser.Attrs[key]; ok && val != "" {
				newUser.Phone = val
			}
		}

		if err := (dao.SysAdminDao{}).CreateSysAdminForLdap(&newUser); err != nil {
			dao.CreateSysLoginInfo(username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "LDAP用户创建失败", 2)
			result.Failed(c, int(result.ApiCode.FAILED), "用户创建失败")
			return
		}

		if ldapConfig.DefaultRoleID > 0 {
			(dao.SysAdminDao{}).AssignRoleToUser(newUser.ID, ldapConfig.DefaultRoleID)
		} else if len(ldapConfig.DefaultRoles) > 0 {
			s.assignDefaultRolesToUser(newUser.ID, ldapConfig.DefaultRoles)
		}

		sysAdmin = newUser
	}

	if sysAdmin.Status == 2 {
		dao.CreateSysLoginInfo(username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "账号已停用", 2)
		result.Failed(c, int(result.ApiCode.STATUSISENABLE), "账号已停用")
		return
	}

	tokenString, _ := jwt.GenerateTokenByAdmin(sysAdmin)
	dao.CreateSysLoginInfo(username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "LDAP登录成功", 1)

	leftMenuList := dao.QueryLeftMenuList(sysAdmin.ID)
	leftMenuVo := make([]model.LeftMenuVo, 0, len(leftMenuList))
	for _, value := range leftMenuList {
		menuSvoList := dao.QueryMenuVoList(sysAdmin.ID, value.Id)
		item := model.LeftMenuVo{
			MenuSvoList: menuSvoList,
			Id:          value.Id,
			MenuName:    value.MenuName,
			Icon:        value.Icon,
			Url:         value.Url,
		}
		leftMenuVo = append(leftMenuVo, item)
	}

	permissionList := dao.QueryPermissionList(sysAdmin.ID)
	stringList := make([]string, 0, len(permissionList))
	for _, value := range permissionList {
		stringList = append(stringList, value.Value)
	}

	result.Success(c, map[string]interface{}{
		"token":          tokenString,
		"sysAdmin":       sysAdmin,
		"leftMenuList":   leftMenuVo,
		"permissionList": stringList,
		"loginType":      "ldap",
	})
}

func (s *SysConfigServiceImpl) assignDefaultRolesToUser(userID uint, roleNames []string) {
	for _, roleName := range roleNames {
		role := dao.SysRoleDao{}.GetRoleByName(roleName)
		if role.ID > 0 {
			(dao.SysAdminDao{}).AssignRoleToUser(userID, role.ID)
		}
	}
}
