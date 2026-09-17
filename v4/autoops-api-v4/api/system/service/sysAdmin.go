// 用户服务层
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
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// 定义接口
type ISysAminService interface {
	Login(c *gin.Context, dto model.LoginDto)                                                           // 登录
	CreateSysAdmin(c *gin.Context, dto model.AddSysAdminDto)                                            // 新增用户
	GetSysAdminInfo(c *gin.Context, Id int)                                                             //根据id查询用户信息
	UpdateSysAdmin(c *gin.Context, dto model.UpdateSysAdminDto)                                         // 修改用户
	DeleteSysAdminById(c *gin.Context, dto model.SysAdminIdDto)                                         // 删除用户
	UpdateSysAdminStatus(c *gin.Context, dto model.UpdateSysAdminStatusDto)                             // 修改用户状态
	ResetSysAdminPassword(c *gin.Context, dto model.ResetSysAdminPasswordDto)                           // 重置密码
	GetSysAdminList(c *gin.Context, PageSize, PageNum int, Username, Status, BeginTime, EndTime string) // 分页查询用户列表
	UpdatePersonal(c *gin.Context, dto model.UpdatePersonalDto)                                         // 修改个人信息
	UpdatePersonalPassword(c *gin.Context, dto model.UpdatePersonalPasswordDto)                         // 修改个人密码
}
type SysAdminServiceImpl struct{}

// 用户登录
func (s SysAdminServiceImpl) Login(c *gin.Context, dto model.LoginDto) {
	// 登录参数校验
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingLoginParameter), result.ApiCode.GetMessage(result.ApiCode.MissingLoginParameter))
		return
	}
	ip := c.ClientIP()
	// 验证码是否过期
	code := util.RedisStore{}.Get(dto.IdKey, true)
	if len(code) == 0 {
		dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "验证码已过期", 2)
		result.Failed(c, int(result.ApiCode.VerificationCodeHasExpired), result.ApiCode.GetMessage(result.ApiCode.VerificationCodeHasExpired))
		return
	}
	// 校验验证码
	verifyRes := CaptVerify(dto.IdKey, dto.Image)
	if !verifyRes {
		dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "验证码不正确", 2)
		result.Failed(c, int(result.ApiCode.CAPTCHANOTTRUE), result.ApiCode.GetMessage(result.ApiCode.CAPTCHANOTTRUE))
		return
	}

	// 先查询本地账号是否存在
	sysAdmin := dao.SysAdminDetail(dto)
	var loginType string = "local" // 登录类型

	// 如果本地账号存在,使用本地密码验证
	if sysAdmin.ID > 0 {
		if sysAdmin.Password != util.EncryptionMd5(dto.Password) {
			// 本地密码验证失败,尝试LDAP认证（支持LDAP用户使用统一登录接口）
			ldapUser, ldapErr := s.tryLDAPLogin(dto.Username, dto.Password)
			if ldapErr != nil {
				// LDAP认证也失败,返回密码错误
				dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "密码不正确", 2)
				result.Failed(c, int(result.ApiCode.PASSWORDNOTTRUE), result.ApiCode.GetMessage(result.ApiCode.PASSWORDNOTTRUE))
				return
			}
			// LDAP认证成功,同步用户信息
			sysAdmin = s.syncLDAPUser(dto.Username, ldapUser)
			loginType = "ldap"
		}
	} else {
		// 本地账号不存在,尝试LDAP认证
		ldapUser, ldapErr := s.tryLDAPLogin(dto.Username, dto.Password)
		if ldapErr != nil {
			// LDAP认证也失败,返回用户不存在
			dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), fmt.Sprintf("用户不存在或密码错误: %v", ldapErr), 2)
			result.Failed(c, int(result.ApiCode.PASSWORDNOTTRUE), "用户不存在或密码错误")
			return
		}

		// LDAP认证成功,创建或更新本地用户
		sysAdmin = s.syncLDAPUser(dto.Username, ldapUser)
		loginType = "ldap"
	}

	// 检查账号状态
	const status int = 2
	if sysAdmin.Status == status {
		dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "账号已停用", 2)
		result.Failed(c, int(result.ApiCode.STATUSISENABLE), result.ApiCode.GetMessage(result.ApiCode.STATUSISENABLE))
		return
	}

	// 生成token
	tokenString, _ := jwt.GenerateTokenByAdmin(sysAdmin)
	dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "登录成功", 1)

	// 左侧菜单列表
	var leftMenuVo []model.LeftMenuVo
	leftMenuList := dao.QueryLeftMenuList(sysAdmin.ID)
	for _, value := range leftMenuList {
		menuSvoList := dao.QueryMenuVoList(sysAdmin.ID, value.Id)
		item := model.LeftMenuVo{}
		item.MenuSvoList = menuSvoList
		item.Id = value.Id
		item.MenuName = value.MenuName
		item.Icon = value.Icon
		item.Url = value.Url
		leftMenuVo = append(leftMenuVo, item)
	}

	// 权限列表
	permissionList := dao.QueryPermissionList(sysAdmin.ID)
	var stringList = make([]string, 0)
	for _, value := range permissionList {
		stringList = append(stringList, value.Value)
	}

	licenseStatus := map[string]interface{}{
		"activated": false,
		"expired":   false,
	}
	if availability, err := NewLicenseService().GetAvailability(); err == nil {
		licenseStatus = map[string]interface{}{
			"activated":     availability.Activated,
			"expired":       availability.Expired,
			"machineCode":   availability.MachineCode,
			"customerName":  availability.CustomerName,
			"licenseType":   availability.LicenseType,
			"issuedAt":      availability.IssuedAt,
			"activatedAt":   availability.ActivatedAt,
			"expireAt":      availability.ExpireAt,
			"daysRemaining": availability.DaysRemaining,
			"validityText":  availability.ValidityText,
		}
	}

	result.Success(c, map[string]interface{}{
		"token":          tokenString,
		"sysAdmin":       sysAdmin,
		"leftMenuList":   leftMenuVo,
		"permissionList": stringList,
		"loginType":      loginType,
		"licenseStatus":  licenseStatus,
	})
}

// tryLDAPLogin 尝试LDAP登录
func (s SysAdminServiceImpl) tryLDAPLogin(username, password string) (*util.LDAPUserAttributes, error) {
	// 获取LDAP配置
	configDao := dao.NewSysConfigDao()
	config, err := configDao.GetByKey("ldap")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("LDAP配置不存在")
		}
		return nil, fmt.Errorf("获取LDAP配置失败: %v", err)
	}

	// 解析LDAP配置
	var ldapConfig model.LDAPConfig
	if err := json.Unmarshal([]byte(config.ConfigData), &ldapConfig); err != nil {
		return nil, fmt.Errorf("LDAP配置解析失败: %v", err)
	}

	// 检查LDAP是否启用
	if !ldapConfig.Enable {
		return nil, fmt.Errorf("LDAP认证未启用")
	}

	// 执行LDAP认证
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

	return util.LDAPAuth(authConfig, username, password)
}

// syncLDAPUser 同步LDAP用户到本地数据库
func (s SysAdminServiceImpl) syncLDAPUser(username string, ldapUser *util.LDAPUserAttributes) model.SysAdmin {
	// 获取LDAP配置用于属性映射
	configDao := dao.NewSysConfigDao()
	config, _ := configDao.GetByKey("ldap")
	var ldapConfig model.LDAPConfig
	json.Unmarshal([]byte(config.ConfigData), &ldapConfig)

	// 检查用户是否已存在
	existingUser := (dao.SysAdminDao{}).GetSysAdminByUsername(username)
	if existingUser.ID > 0 {
		// 如果配置了覆盖属性,更新用户信息
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
		return existingUser
	}

	// 创建新用户
	newUser := model.SysAdmin{
		Username: username,
		Password: util.EncryptionMd5("ldap_user_" + username),
		Nickname: ldapUser.CN,
		Status:   1,
		PostId:   1,
		DeptId:   1,
	}

	// 映射LDAP属性
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

	// 创建用户
	(dao.SysAdminDao{}).CreateSysAdminForLdap(&newUser)

	// 分配默认角色
	if ldapConfig.DefaultRoleID > 0 {
		(dao.SysAdminDao{}).AssignRoleToUser(newUser.ID, ldapConfig.DefaultRoleID)
	}

	return newUser
}

// 新增用户
func (s SysAdminServiceImpl) CreateSysAdmin(c *gin.Context, dto model.AddSysAdminDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingNewAdminParameter), result.ApiCode.GetMessage(result.ApiCode.MissingNewAdminParameter))
		return
	}
	bool := dao.CreateSysAdmin(dto)
	if !bool {
		result.Failed(c, int(result.ApiCode.USERNAMEALREADYEXISTS), result.ApiCode.GetMessage(result.ApiCode.USERNAMEALREADYEXISTS))
		return
	}
	result.Success(c, bool)
}

// 根据id查询用户信息
func (s SysAdminServiceImpl) GetSysAdminInfo(c *gin.Context, Id int) {
	result.Success(c, dao.GetSysAdminInfo(Id))
}

// 修改用户
func (s SysAdminServiceImpl) UpdateSysAdmin(c *gin.Context, dto model.UpdateSysAdminDto) {
	result.Success(c, dao.UpdateSysAdmin(dto))
}

// 根据id删除用户
func (s SysAdminServiceImpl) DeleteSysAdminById(c *gin.Context, dto model.SysAdminIdDto) {
	// 获取要删除的用户信息
	userInfo := dao.GetSysAdminInfo(int(dto.Id))

	// 禁止删除 admin 和 root 用户
	if userInfo.Username == "admin" || userInfo.Username == "root" {
		result.Failed(c, int(result.ApiCode.FAILED), "禁止删除系统管理员账号(admin/root)")
		return
	}

	dao.DeleteSysAdminById(dto)
	result.Success(c, true)
}

// 修改用户状态
func (s SysAdminServiceImpl) UpdateSysAdminStatus(c *gin.Context, dto model.UpdateSysAdminStatusDto) {
	dao.UpdateSysAdminStatus(dto)
	result.Success(c, true)
}

// 重置密码
func (s SysAdminServiceImpl) ResetSysAdminPassword(c *gin.Context, dto model.ResetSysAdminPasswordDto) {
	dao.ResetSysAdminPassword(dto)
	result.Success(c, true)
}

// 分页查询用户列表
func (s SysAdminServiceImpl) GetSysAdminList(c *gin.Context, PageSize, PageNum int, Username, Status, BeginTime, EndTime string) {
	if PageSize < 1 {
		PageSize = 10
	}
	if PageNum < 1 {
		PageNum = 1
	}
	sysAdmin, count := dao.GetSysAdminList(PageSize, PageNum, Username, Status, BeginTime, EndTime)
	result.Success(c, map[string]interface{}{"total": count, "pageSize": PageSize, "pageNum": PageNum, "list": sysAdmin})
}

// 修改个人信息
func (s SysAdminServiceImpl) UpdatePersonal(c *gin.Context, dto model.UpdatePersonalDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingModificationOfPersonalParameters), result.ApiCode.GetMessage(result.ApiCode.MissingModificationOfPersonalParameters))
		return
	}
	// 获取当前用户信息
	sysAdmin, _ := jwt.GetAdmin(c)
	// 禁止 test 用户修改个人信息
	if sysAdmin.Username == "test" {
		result.Failed(c, int(result.ApiCode.TESTUSERFORBIDDEN), result.ApiCode.GetMessage(result.ApiCode.TESTUSERFORBIDDEN))
		return
	}
	id, _ := jwt.GetAdminId(c)
	dto.Id = id
	result.Success(c, dao.UpdatePersonal(dto))
}

// 修改个人密码
func (s SysAdminServiceImpl) UpdatePersonalPassword(c *gin.Context, dto model.UpdatePersonalPasswordDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingChangePasswordParameter), result.ApiCode.GetMessage(result.ApiCode.MissingChangePasswordParameter))
		return
	}
	sysAdmin, _ := jwt.GetAdmin(c)
	// 禁止 test 用户修改密码
	if sysAdmin.Username == "test" {
		result.Failed(c, int(result.ApiCode.TESTUSERFORBIDDEN), result.ApiCode.GetMessage(result.ApiCode.TESTUSERFORBIDDEN))
		return
	}
	dto.Id = sysAdmin.ID
	sysAdminExist := dao.GetSysAdminByUsername(sysAdmin.Username)
	if sysAdminExist.Password != util.EncryptionMd5(dto.Password) {
		result.Failed(c, int(result.ApiCode.PASSWORDNOTTRUE), result.ApiCode.GetMessage(result.ApiCode.PASSWORDNOTTRUE))
		return
	}
	if dto.NewPassword != dto.ResetPassword {
		result.Failed(c, int(result.ApiCode.RESETPASSWORD), result.ApiCode.GetMessage(result.ApiCode.RESETPASSWORD))
		return
	}
	dto.NewPassword = util.EncryptionMd5(dto.NewPassword)
	sysAdminUpdatePwd := dao.UpdatePersonalPassword(dto)
	tokenString, _ := jwt.GenerateTokenByAdmin(sysAdminUpdatePwd)
	result.Success(c, map[string]interface{}{"token": tokenString, "sysAdmin": sysAdminUpdatePwd})
}

var sysAdminService = SysAdminServiceImpl{}

func SysAdminService() ISysAminService {
	return &sysAdminService
}
