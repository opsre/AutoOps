package cmdb

import (
	"dodevops-api/api/cmdb/controller"
	"dodevops-api/api/cmdb/service"
	"github.com/gin-gonic/gin"
)

// RegisterCmdbRoutes 注册系统相关路由
func RegisterCmdbRoutes(router *gin.RouterGroup) {
	// 资产分组
	router.POST("/cmdb/groupadd", controller.CreateCmdbGroup)                    // 添加资产分组
	router.GET("/cmdb/grouplist", controller.GetAllCmdbGroups)                   // 获取所有资产分组树
	router.GET("/cmdb/grouplistwithhosts", controller.GetAllCmdbGroupsWithHosts) // 获取所有资产分组树及关联主机
	router.PUT("/cmdb/groupupdate", controller.UpdateCmdbGroup)                  // 更新资产分组
	router.DELETE("/cmdb/groupdelete", controller.DeleteCmdbGroup)               // 删除资产分组
	router.GET("/cmdb/groupbyname", controller.GetCmdbGroupByName)               // 根据名称查询分组
	// 主机管理
	router.POST("/cmdb/hostcreate", controller.NewCmdbHostController().CreateCmdbHost)                 // 创建主机
	router.PUT("/cmdb/hostupdate", controller.NewCmdbHostController().UpdateCmdbHost)                  // 更新主机
	router.DELETE("/cmdb/hostdelete", controller.NewCmdbHostController().DeleteCmdbHost)               // 删除主机
	router.GET("/cmdb/hostlist", controller.NewCmdbHostController().GetCmdbHostListWithPage)           // 获取主机列表(分页)
	router.GET("/cmdb/hostinfo", controller.NewCmdbHostController().GetCmdbHostById)                   // 根据ID获取主机
	router.GET("/cmdb/hostgroup", controller.NewCmdbHostController().GetCmdbHostsByGroupId)            // 根据分组ID获取主机列表
	router.GET("/cmdb/hostbyname", controller.NewCmdbHostController().GetCmdbHostsByHostNameLike)      // 根据主机名称+ip模糊查询
	router.GET("/cmdb/hostbyip", controller.NewCmdbHostController().GetCmdbHostsByIP)                  // 根据IP查询主机
	router.GET("/cmdb/hostbystatus", controller.NewCmdbHostController().GetCmdbHostsByStatus)          // 根据状态查询主机
	router.POST("/cmdb/hostimport", controller.NewCmdbHostController().ImportHostsFromExcel)           // 从Excel导入主机
	router.GET("/cmdb/hosttemplate", controller.NewCmdbHostController().DownloadHostTemplate)          // 下载主机导入模板
	router.POST("/cmdb/hostsync", controller.NewCmdbHostController().SyncHostInfo)                     // 同步主机基本信息
	router.GET("/cmdb/hostcloudstatistics", controller.NewCmdbHostController().GetCloudHostStatistics) // 获取云主机统计信息
	router.PUT("/cmdb/hostbatchupdategroup", controller.NewCmdbHostController().BatchUpdateHostGroup)  // 批量更新主机分组
	router.DELETE("/cmdb/hostbatchdelete", controller.NewCmdbHostController().BatchDeleteHost)         // 批量删除主机
	// 云主机管理
	router.POST("/cmdb/hostcloudcreatealiyun", controller.NewCmdbHostCloudController().CreateAliyunHost)   // 创建阿里云主机
	router.POST("/cmdb/hostcloudcreatetencent", controller.NewCmdbHostCloudController().CreateTencentHost) // 创建腾讯云主机
	router.POST("/cmdb/hostcloudcreatebaidu", controller.NewCmdbHostCloudController().CreateBaiduHost)     // 创建百度云主机
	router.POST("/cmdb/hostcloudcreatehuawei", controller.NewCmdbHostCloudController().CreateHuaweiHost)   // 创建华为云主机
	router.POST("/cmdb/hostcloudcreateaws", controller.NewCmdbHostCloudController().CreateAWSHost)         // 创建AWS主机
	// JumpServer主机同步
	router.POST("/cmdb/hostjumpserversync", controller.NewCmdbHostJumpServerController().SyncJumpServerHosts)                                 // 同步JumpServer主机
	router.GET("/cmdb/hostssh/connect/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).ConnectTerminal)             // SSH终端连接
	router.GET("/cmdb/hostssh/command/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).ExecuteCommand)              // SSH执行命令
	router.POST("/cmdb/hostssh/upload/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).UploadFile)                  // SSH文件上传
	router.GET("/cmdb/hostssh/upload/task/:taskId", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).GetUploadTaskStatus) // SSH上传任务状态
	router.GET("/cmdb/hostssh/download/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).DownloadFile)               // SSH文件下载
	router.GET("/cmdb/hostssh/listfiles/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).ListFiles)                 // SSH列出文件
	// SQL执行
	router.POST("/cmdb/sql/select", controller.GetCmdbSQLRecordController().ExecuteSelect)       // 执行查询语句(通过数据库ID/名称)
	router.POST("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteInsert)              // 执行插入语句(通过数据库ID/名称)
	router.PUT("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteUpdate)               // 执行更新语句(通过数据库ID/名称)
	router.DELETE("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteDelete)            // 执行删除语句(通过数据库ID/名称)
	router.POST("/cmdb/sql/execute", controller.GetCmdbSQLRecordController().ExecuteSQL)         // 执行原生SQL语句(通过数据库ID/名称)
	router.POST("/cmdb/sql/databaselist", controller.GetCmdbSQLRecordController().ListDatabases) // 获取数据库列表(通过数据库ID)
	// SQL日志管理
	router.GET("/cmdb/sqlLog/list", controller.GetCmdbSqlLogList)         // 分页获取SQL操作日志列表
	router.DELETE("/cmdb/sqlLog/delete", controller.DeleteCmdbSqlLogById) // 根据id删除SQL操作日志
	router.DELETE("/cmdb/sqlLog/clean", controller.CleanCmdbSqlLog)       // 清空SQL操作日志
	// 数据库管理
	router.POST("/cmdb/database", controller.NewCmdbSQLController().CreateDatabase)           // 创建数据库
	router.PUT("/cmdb/database", controller.NewCmdbSQLController().UpdateDatabase)            // 更新数据库
	router.DELETE("/cmdb/database", controller.NewCmdbSQLController().DeleteDatabase)         // 删除数据库
	router.GET("/cmdb/database/info", controller.NewCmdbSQLController().GetDatabase)          // 根据ID获取数据库详情
	router.GET("/cmdb/databaselist", controller.NewCmdbSQLController().ListDatabases)         // 获取数据库列表
	router.GET("/cmdb/database/byname", controller.NewCmdbSQLController().GetDatabasesByName) // 根据名称查询数据库
	router.GET("/cmdb/database/bytype", controller.NewCmdbSQLController().GetDatabasesByType) // 根据类型查询数据库

	// SNMP设备资产管理(简化版)
	snmpDeviceCtrl := controller.NewCmdbSnmpDeviceController()
	router.POST("/cmdb/snmp-devices", snmpDeviceCtrl.CreateSnmpDevice)                            // 创建SNMP设备
	router.PUT("/cmdb/snmp-devices/:id", snmpDeviceCtrl.UpdateSnmpDevice)                         // 更新SNMP设备
	router.DELETE("/cmdb/snmp-devices/:id", snmpDeviceCtrl.DeleteSnmpDevice)                      // 删除SNMP设备
	router.GET("/cmdb/snmp-devices/:id", snmpDeviceCtrl.GetSnmpDeviceById)                        // 根据ID获取SNMP设备
	router.GET("/cmdb/snmp-devices", snmpDeviceCtrl.GetSnmpDeviceList)                            // 获取SNMP设备列表
	router.POST("/cmdb/snmp-devices/:id/collect", snmpDeviceCtrl.CollectNow)                      // 立即采集(手动触发)
	router.POST("/cmdb/snmp-devices/:id/enable", snmpDeviceCtrl.EnableMonitor)                    // 启用监控
	router.POST("/cmdb/snmp-devices/:id/disable", snmpDeviceCtrl.DisableMonitor)                  // 停止监控
	router.GET("/cmdb/snmp-devices/:id/monitor", snmpDeviceCtrl.GetDeviceMonitorMetrics)          // 获取设备监控指标
	router.GET("/cmdb/snmp-devices/:id/resource", snmpDeviceCtrl.GetDeviceResourceMetrics)        // 获取设备资源使用率(CPU+内存)
	router.GET("/cmdb/snmp-devices/batch-resource", snmpDeviceCtrl.GetBatchDeviceResourceMetrics) // 批量获取设备资源使用率

	// SNMP设备SSH管理
	snmpDeviceSSHCtrl := controller.NewCmdbSnmpDeviceSSHController()
	router.GET("/cmdb/snmp-devices/:id/ssh/connect", snmpDeviceSSHCtrl.ConnectTerminal) // SSH终端连接
	router.GET("/cmdb/snmp-devices/:id/ssh/command", snmpDeviceSSHCtrl.ExecuteCommand)  // SSH执行命令

	// Windows Exporter 管理
	exporterCtrl := controller.NewWindowsExporterController()
	router.POST("/cmdb/windows-exporter/install", exporterCtrl.InstallExporter)                 // 安装Windows Exporter
	router.GET("/cmdb/windows-exporter/status/:deviceId", exporterCtrl.GetExporterStatus)       // 获取安装状态
	router.DELETE("/cmdb/windows-exporter/uninstall/:deviceId", exporterCtrl.UninstallExporter) // 卸载Windows Exporter
	router.POST("/cmdb/windows-exporter/batch-install", exporterCtrl.BatchInstallExporter)      // 批量安装
	router.GET("/cmdb/windows-exporter/status-all", exporterCtrl.GetAllWindowsDevicesStatus)    // 获取所有Windows设备状态

	// Windows 主机监控
	windowsMonitorCtrl := controller.NewWindowsMonitorController()
	router.GET("/cmdb/windows-devices/:id/monitor", windowsMonitorCtrl.GetWindowsMonitorMetrics)          // 获取Windows主机监控指标(30m/1d/7d)
	router.GET("/cmdb/windows-devices/:id/realtime", windowsMonitorCtrl.GetWindowsRealtimeMetrics)        // 获取Windows主机实时CPU/内存
	router.GET("/cmdb/windows-devices/batch-realtime", windowsMonitorCtrl.GetBatchWindowsRealtimeMetrics) // 批量获取Windows主机实时CPU/内存

	// 堡垒机用户组管理
	bastionGroupCtrl := controller.NewBastionUserGroupController()
	router.POST("/bastion/group/create", bastionGroupCtrl.CreateUserGroup)           // 创建用户组
	router.PUT("/bastion/group/update", bastionGroupCtrl.UpdateUserGroup)            // 更新用户组
	router.DELETE("/bastion/group/delete", bastionGroupCtrl.DeleteUserGroup)         // 删除用户组
	router.GET("/bastion/group/info", bastionGroupCtrl.GetUserGroupById)             // 根据ID获取用户组
	router.GET("/bastion/group/list", bastionGroupCtrl.GetUserGroupList)             // 获取用户组列表(分页)
	router.POST("/bastion/group/member/add", bastionGroupCtrl.AddGroupMembers)       // 添加用户组成员
	router.POST("/bastion/group/member/remove", bastionGroupCtrl.RemoveGroupMembers) // 删除用户组成员
	router.GET("/bastion/group/member/list", bastionGroupCtrl.GetGroupMembers)       // 获取用户组成员列表
	router.GET("/bastion/user/groups", bastionGroupCtrl.GetUserGroups)               // 获取用户所属的用户组

	// 堡垒机授权管理
	bastionAuthCtrl := controller.NewBastionHostAuthorizationController()
	router.POST("/bastion/auth/user/authorize", bastionAuthCtrl.AuthorizeHostsToUser)     // 授权主机给用户
	router.POST("/bastion/auth/group/authorize", bastionAuthCtrl.AuthorizeHostsToGroup)   // 授权主机给用户组
	router.POST("/bastion/auth/user/revoke", bastionAuthCtrl.RevokeHostsFromUser)         // 取消用户的主机授权
	router.POST("/bastion/auth/group/revoke", bastionAuthCtrl.RevokeHostsFromGroup)       // 取消用户组的主机授权
	router.GET("/bastion/auth/hosts", bastionAuthCtrl.GetAuthorizedHosts)                 // 获取授权主机列表
	router.GET("/bastion/user/authorized-hosts", bastionAuthCtrl.GetUserAuthorizedHosts)  // 获取用户有权限的主机列表
	router.GET("/bastion/user/check-permission", bastionAuthCtrl.CheckUserHostPermission) // 检查用户是否有主机权限
}
