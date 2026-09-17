// utils/system/system.go

package system

import (
	"dodevops-api/api/system/controller"

	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 注册系统相关路由
func RegisterSystemRoutes(router *gin.RouterGroup) {
	// 使用传入的 router，不再创建子分组
	systemGroup := router

	// 岗位
	systemGroup.POST("/post/add", controller.CreateSysPost)
	systemGroup.GET("/post/list", controller.GetSysPostList)
	systemGroup.GET("/post/info", controller.GetSysPostById)
	systemGroup.PUT("/post/update", controller.UpdateSysPost)
	systemGroup.DELETE("/post/delete", controller.DeleteSysPostById)
	systemGroup.DELETE("/post/batch/delete", controller.BatchDeleteSysPost)
	systemGroup.PUT("/post/updateStatus", controller.UpdateSysPostStatus)
	systemGroup.GET("/post/vo/list", controller.QuerySysPostVoList)
	//  部门
	systemGroup.GET("/dept/list", controller.GetSysDeptList)         // 部门列表
	systemGroup.POST("/dept/add", controller.CreateSysDept)          // 添加部门
	systemGroup.GET("/dept/info", controller.GetSysDeptById)         // 根据id查询部门
	systemGroup.PUT("/dept/update", controller.UpdateSysDept)        // 修改部门
	systemGroup.DELETE("/dept/delete", controller.DeleteSysDeptById) // 删除部门
	systemGroup.GET("/dept/vo/list", controller.QuerySysDeptVoList)  // 查询部门树形结构
	systemGroup.GET("/dept/users", controller.GetDeptUsers)          // 查询部门下的用户
	// 菜单
	systemGroup.POST("/menu/add", controller.CreateSysMenu)
	systemGroup.GET("/menu/vo/list", controller.QuerySysMenuVoList)
	systemGroup.GET("/menu/info", controller.GetSysMenu)
	systemGroup.PUT("/menu/update", controller.UpdateSysMenu)
	systemGroup.DELETE("/menu/delete", controller.DeleteSysMenu)
	systemGroup.GET("/menu/list", controller.GetSysMenuList)
	// 角色
	systemGroup.POST("/role/add", controller.CreateSysRole)
	systemGroup.GET("/role/info", controller.GetSysRoleById)
	systemGroup.PUT("/role/update", controller.UpdateSysRole)
	systemGroup.DELETE("/role/delete", controller.DeleteSysRoleById)
	systemGroup.PUT("/role/updateStatus", controller.UpdateSysRoleStatus)
	systemGroup.GET("/role/list", controller.GetSysRoleList)
	systemGroup.GET("/role/vo/list", controller.QuerySysRoleVoList)
	systemGroup.GET("/role/vo/idList", controller.QueryRoleMenuIdList)
	systemGroup.PUT("/role/assignPermissions", controller.AssignPermissions)
	// 用户
	systemGroup.POST("/admin/add", controller.CreateSysAdmin)
	systemGroup.GET("/admin/info", controller.GetSysAdminInfo)
	systemGroup.PUT("/admin/update", controller.UpdateSysAdmin)
	systemGroup.DELETE("/admin/delete", controller.DeleteSysAdminById)
	systemGroup.PUT("/admin/updateStatus", controller.UpdateSysAdminStatus)
	systemGroup.PUT("/admin/updatePassword", controller.ResetSysAdminPassword)
	systemGroup.GET("/admin/list", controller.GetSysAdminList)
	systemGroup.POST("/upload", controller.Upload)
	systemGroup.PUT("/admin/updatePersonal", controller.UpdatePersonal)
	systemGroup.PUT("/admin/updatePersonalPassword", controller.UpdatePersonalPassword)
	// 日志
	systemGroup.GET("/sysLoginInfo/list", controller.GetSysLoginInfoList)
	systemGroup.DELETE("/sysLoginInfo/batch/delete", controller.BatchDeleteSysLoginInfo)
	systemGroup.DELETE("/sysLoginInfo/delete", controller.DeleteSysLoginInfoById)
	systemGroup.DELETE("/sysLoginInfo/clean", controller.CleanSysLoginInfo)
	systemGroup.GET("/sysOperationLog/list", controller.GetSysOperationLogList)
	systemGroup.DELETE("/sysOperationLog/delete", controller.DeleteSysOperationLogById)
	systemGroup.DELETE("/sysOperationLog/batch/delete", controller.BatchDeleteSysOperationLog)
	systemGroup.DELETE("/sysOperationLog/clean", controller.CleanSysOperationLog)
	// 系统动态日志
	activityLogController := controller.NewSysActivityLogController()
	systemGroup.POST("/activity-log", activityLogController.Create)
	systemGroup.GET("/activity-log/detail", activityLogController.GetByID)
	systemGroup.GET("/activity-log/list", activityLogController.List)
	systemGroup.GET("/activity-log/latest", activityLogController.GetLatest)
	systemGroup.GET("/activity-log/related", activityLogController.GetByRelatedID)
	systemGroup.DELETE("/activity-log", activityLogController.DeleteByID)
	systemGroup.DELETE("/activity-log/batch", activityLogController.DeleteBatch)
	systemGroup.DELETE("/activity-log/clean", activityLogController.Clean)
	// 系统配置
	systemGroup.POST("/config/add", controller.CreateSysConfig)
	systemGroup.GET("/config/info", controller.GetSysConfigInfo)
	systemGroup.GET("/config/key", controller.GetSysConfigByKey)
	systemGroup.PUT("/config/update", controller.UpdateSysConfig)
	systemGroup.DELETE("/config/delete", controller.DeleteSysConfig)
	systemGroup.GET("/config/list", controller.GetSysConfigList)
	systemGroup.PUT("/config/updateStatus", controller.UpdateSysConfigStatus)
	// LDAP配置
	systemGroup.POST("/ldap/config", controller.SaveLDAPConfig)
	systemGroup.GET("/ldap/config", controller.GetLDAPConfig)
	systemGroup.POST("/ldap/test", controller.TestLDAPConnection)
	systemGroup.POST("/ldap/login", controller.LDAPLogin) // LDAP登录（不需要认证）
	// 系统激活
	systemGroup.GET("/license/machine-code", controller.GetLicenseMachineCode)
	systemGroup.GET("/license/status", controller.GetLicenseStatus)
	systemGroup.POST("/license/activate", controller.ActivateLicense)

	// 会话录制管理
	recordingCtrl, err := controller.NewSysSessionRecordingController("./upload")
	if err == nil {
		systemGroup.GET("/recordings", recordingCtrl.GetRecordingList)
		systemGroup.GET("/recordings/:id", recordingCtrl.GetRecordingDetail)
		systemGroup.GET("/recordings/:id/cast", recordingCtrl.GetRecordingCast)
		systemGroup.GET("/recordings/:id/download", recordingCtrl.DownloadRecording)
		systemGroup.DELETE("/recordings/:id", recordingCtrl.DeleteRecording)
		systemGroup.DELETE("/recordings/batch", recordingCtrl.BatchDeleteRecording)
	}

	// 命令审计管理 - 注意：具体路径要放在参数路径之前
	commandAuditCtrl := controller.NewSysCommandAuditController()
	systemGroup.GET("/recordings/commands", commandAuditCtrl.GetCommandAuditList)
	systemGroup.GET("/recordings/commands/sensitive-count/:recordingId", commandAuditCtrl.GetSensitiveCommandCount)
	systemGroup.GET("/recordings/commands/highest-risk/:recordingId", commandAuditCtrl.GetHighestRiskLevel)
}
