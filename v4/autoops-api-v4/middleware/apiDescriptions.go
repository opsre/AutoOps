// API描述映射
// 自动生成，请勿手动修改
// 使用 sync_api_descriptions.py 生成

package middleware

import (
	"regexp"
	"strings"
)

// GetAPIDescription 根据 URL 路径和 HTTP 方法获取 API 的中文描述
func GetAPIDescription(url, method string) string {
	// 精确匹配
	if desc, ok := getDescriptionMap()[url]; ok {
		return desc
	}

	// 方法+路径匹配
	key := method + ":" + url
	if desc, ok := getMethodBasedDescriptionMap()[key]; ok {
		return desc
	}

	// 动态路由匹配
	return matchDynamicRoute(url)
}

// getDescriptionMap 返回 API 路径到描述的映射
func getDescriptionMap() map[string]string {
	descriptions := map[string]string{
		"/api/v1/ai/chat": "使用AI模型进行对话（支持智能运维Agent）",
		"/api/v1/ai/model/add": "创建新的AI模型配置",
		"/api/v1/ai/model/delete": "根据ID删除AI模型",
		"/api/v1/ai/model/enabled": "获取所有启用状态的AI模型",
		"/api/v1/ai/model/info": "根据ID获取AI模型详情",
		"/api/v1/ai/model/list": "获取AI模型列表（支持分页和筛选）",
		"/api/v1/ai/model/update": "更新AI模型信息",
		"/api/v1/ai/model/updateStatus": "更新AI模型状态",
		"/api/v1/config/sync-schedule/trigger": "手动触发同步任务",
		"/api/v1/monitor/aliyun-config/list": "查询所有阿里云配置列表",
		"/api/v1/monitor/billing/aliyun-product": "根据密钥ID获取阿里云账单（按产品统计）",
		"/api/v1/monitor/billing/aliyun-tag": "根据密钥ID获取阿里云账单（按标签bline统计）",
		"/api/v1/monitor/billing/aws": "根据密钥ID获取AWS账单（美元转换为人民币，按SERVICE统计）",
		"/api/v1/monitor/billing/last-12-months-trend": "获取各云厂商最近12个月的费用趋势数据（折线图）",
		"/api/v1/monitor/billing/last-month-summary": "获取各云厂商上月费用汇总以及与上上月的环比",
		"/api/v1/monitor/billing/tencent": "根据密钥ID获取腾讯云账单（按业务线统计）",
		"/api/v1/monitor/domain/batch-upload": "上传包含域名的文本文件，每行一个域名，自动去重并异步创建",
		"/api/v1/monitor/domain/check-all": "立即检查所有启用的域名",
		"/api/v1/monitor/domain/check/{id}": "立即检查指定域名的存活状态和SSL证书",
		"/api/v1/monitor/domain/import-aliyun-dns": "通过阿里云密钥ID从阿里云解析自动导入所有域名,支持导入子域名",
		"/api/v1/monitor/domain/list": "获取域名监控列表（支持分页和筛选）",
		"/api/v1/monitor/domain/schedule/pause": "暂停域名监控定时任务",
		"/api/v1/monitor/domain/schedule/resume": "恢复域名监控定时任务",
		"/api/v1/monitor/domain/schedule/status": "获取域名监控定时任务运行状态",
		"/api/v1/monitor/domain/ssl-expiring": "获取SSL证书即将过期的域名列表",
		"/api/v1/monitor/domain/statistics": "获取域名监控统计信息（总数、正常、异常、SSL过期等）",
		"/api/v1/monitor/incident/batch-delete": "批量删除故障记录",
		"/api/v1/monitor/incident/list": "获取故障列表",
		"/api/v1/monitor/incident/statistics": "获取故障统计信息",
		"/api/v1/monitor/incident/{id}/status": "更新故障状态",
		"/api/v1/monitor/robot/list": "分页获取通知机器人列表",
		"/api/v1/monitor/robot/test/{id}": "测试发送消息到指定机器人",
		"/api/v1/monitor/ssl-cert/apply": "申请阿里云免费SSL证书",
		"/api/v1/monitor/ssl-cert/deploy": "异步部署SSL证书到目标主机(备份、拷贝、重载Nginx)，立即返回部署日志ID，前端可通过查询部署日志接口获取实时进度",
		"/api/v1/monitor/ssl-cert/deploy-logs": "获取SSL证书部署日志列表，支持多条件筛选和分页",
		"/api/v1/monitor/ssl-cert/list": "查询SSL证书列表（分页）",
		"/api/v1/monitor/ssl-cert/upload": "上传SSL证书压缩包(仅支持zip格式),自动解析证书和私钥并录入数据库。压缩包中应包含xxx.key和xxx.pem/crt/cer文件,域名从文件名自动提取",
		"/api/v1/monitor/ssl-cert/{id}": "删除SSL证书记录",
		"/api/v1/monitor/ssl-cert/{id}/check-validation": "检查DNS验证是否通过",
		"/api/v1/monitor/ssl-cert/{id}/dns-validation": "自动添加DNS TXT验证记录到阿里云DNS",
		"/api/v1/monitor/ssl-cert/{id}/download": "下载已签发的SSL证书",
		"/api/v1/monitor/ssl-cert/{id}/state": "查询证书申请状态和验证信息",
		"/api/v1/monitor/webhook/log/list": "分页获取Webhook日志列表",
		"/api/v1/monitor/webhook/log/{id}": "根据ID获取Webhook日志详情",
		"/api/v1/monitor/webhook/log/{webhookLogId}/notify-logs": "获取指定Webhook日志的通知详情列表",
		"/api/v1/monitor/webhook/receive": "接收外部系统的Webhook通知并转发到配置的通知机器人",
		"/api/v1/system/config/add": "创建新的系统配置",
		"/api/v1/system/config/delete": "根据ID删除系统配置",
		"/api/v1/system/config/info": "根据ID获取系统配置详情",
		"/api/v1/system/config/key": "根据配置键获取系统配置",
		"/api/v1/system/config/list": "获取系统配置列表（支持分页和筛选）",
		"/api/v1/system/config/update": "更新系统配置信息",
		"/api/v1/system/config/updateStatus": "更新系统配置状态",
		"/api/v1/system/ldap/login": "使用LDAP进行用户认证登录",
		"/api/v1/system/ldap/test": "测试LDAP服务器连接",
		"/api/v1/tool/all": "获取所有导航工具（不分页）",
		"/api/v1/tool/deploy": "创建一个服务部署任务，异步执行部署",
		"/api/v1/tool/deploy/list": "获取服务部署历史记录列表（分页）",
		"/api/v1/tool/deploy/{id}": "停止并删除已部署的服务",
		"/api/v1/tool/deploy/{id}/status": "根据部署ID获取部署状态和日志",
		"/api/v1/tool/list": "获取导航工具列表（分页）",
		"/api/v1/tool/services": "获取所有可部署的服务列表及分类",
		"/api/v1/tool/services/{serviceId}": "根据服务ID获取服务的详细信息",
		"/api/v1/upload": "单图片上传接口",
	}

	return descriptions
}

// getMethodBasedDescriptionMap 返回基于 HTTP Method + URL 的映射
func getMethodBasedDescriptionMap() map[string]string {
	return map[string]string{
		"delete:/api/v1/monitor/aliyun-config/{id}": "删除阿里云配置记录",
		"delete:/api/v1/monitor/domain/{id}": "根据ID删除域名监控",
		"delete:/api/v1/monitor/incident/{id}": "删除故障记录",
		"delete:/api/v1/monitor/robot/{id}": "根据ID删除通知机器人",
		"delete:/api/v1/tool/{id}": "删除导航工具接口",
		"get:/api/v1/monitor/aliyun-config/{id}": "根据ID获取阿里云配置详情",
		"get:/api/v1/monitor/domain/schedule/config": "获取域名监控定时任务配置信息",
		"get:/api/v1/monitor/domain/{id}": "根据ID获取域名监控详情",
		"get:/api/v1/monitor/incident/{id}": "根据ID获取故障记录",
		"get:/api/v1/monitor/robot/{id}": "根据ID获取通知机器人详情",
		"get:/api/v1/system/ldap/config": "获取LDAP配置信息",
		"get:/api/v1/tool/{id}": "根据ID获取导航工具详情",
		"post:/api/v1/monitor/aliyun-config": "添加新的阿里云AK/SK配置",
		"post:/api/v1/monitor/domain": "添加新的域名监控",
		"post:/api/v1/monitor/domain/schedule/config": "更新域名监控定时任务配置（开启/关闭，设置检查间隔）",
		"post:/api/v1/monitor/incident": "创建故障记录",
		"post:/api/v1/monitor/robot": "创建通知机器人（支持钉钉/飞书/企业微信/邮箱/通用Webhook/Teams）",
		"post:/api/v1/system/ldap/config": "保存或更新LDAP配置",
		"post:/api/v1/tool": "创建导航工具接口",
		"put:/api/v1/monitor/aliyun-config": "更新阿里云配置信息",
		"put:/api/v1/monitor/domain": "更新域名监控信息",
		"put:/api/v1/monitor/incident": "更新故障记录",
		"put:/api/v1/monitor/robot": "更新通知机器人信息",
		"put:/api/v1/tool": "更新导航工具接口",
	}
}

// matchDynamicRoute 匹配带参数的动态路由
func matchDynamicRoute(url string) string {
	dynamicRoutes := []struct {
		pattern *regexp.Regexp
		desc    string
	}{
		{regexp.MustCompile(`^/api/v1/monitor/aliyun-config/\d+$`), "根据ID获取阿里云配置详情"},
		{regexp.MustCompile(`^/api/v1/monitor/domain/\d+$`), "根据ID获取域名监控详情"},
		{regexp.MustCompile(`^/api/v1/monitor/domain/check/\d+$`), "立即检查指定域名的存活状态和SSL证书"},
		{regexp.MustCompile(`^/api/v1/monitor/incident/\d+$`), "删除故障记录"},
		{regexp.MustCompile(`^/api/v1/monitor/incident/\d+/status$`), "更新故障状态"},
		{regexp.MustCompile(`^/api/v1/monitor/robot/\d+$`), "根据ID获取通知机器人详情"},
		{regexp.MustCompile(`^/api/v1/monitor/robot/test/\d+$`), "测试发送消息到指定机器人"},
		{regexp.MustCompile(`^/api/v1/monitor/ssl-cert/\d+$`), "删除SSL证书记录"},
		{regexp.MustCompile(`^/api/v1/monitor/ssl-cert/\d+/check-validation$`), "检查DNS验证是否通过"},
		{regexp.MustCompile(`^/api/v1/monitor/ssl-cert/\d+/dns-validation$`), "自动添加DNS TXT验证记录到阿里云DNS"},
		{regexp.MustCompile(`^/api/v1/monitor/ssl-cert/\d+/download$`), "下载已签发的SSL证书"},
		{regexp.MustCompile(`^/api/v1/monitor/ssl-cert/\d+/state$`), "查询证书申请状态和验证信息"},
		{regexp.MustCompile(`^/api/v1/monitor/webhook/log/[^/]+/notify-logs$`), "获取指定Webhook日志的通知详情列表"},
		{regexp.MustCompile(`^/api/v1/monitor/webhook/log/\d+$`), "根据ID获取Webhook日志详情"},
		{regexp.MustCompile(`^/api/v1/tool/\d+$`), "根据ID获取导航工具详情"},
		{regexp.MustCompile(`^/api/v1/tool/deploy/\d+$`), "停止并删除已部署的服务"},
		{regexp.MustCompile(`^/api/v1/tool/deploy/\d+/status$`), "根据部署ID获取部署状态和日志"},
		{regexp.MustCompile(`^/api/v1/tool/services/[^/]+$`), "根据服务ID获取服务的详细信息"},
	}

	for _, route := range dynamicRoutes {
		if route.pattern.MatchString(url) {
			return route.desc
		}
	}

	return extractActionFromURL(url)
}

// extractActionFromURL 从 URL 中提取操作类型
func extractActionFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		action := parts[len(parts)-1]
		switch action {
		case "add", "create":
			return "新增"
		case "update":
			return "修改"
		case "delete":
			return "删除"
		case "batchDelete", "batch-delete":
			return "批量删除"
		case "clean":
			return "清空"
		case "execute", "start":
			return "执行"
		case "stop":
			return "停止"
		case "deploy":
			return "部署"
		case "sync":
			return "同步"
		}
	}
	return ""
}
