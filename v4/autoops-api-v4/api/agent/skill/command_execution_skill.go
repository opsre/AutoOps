// 命令执行 Skill
package skill

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	monitorService "dodevops-api/api/monitor/service"
	"dodevops-api/common/util"
	baseSkill "dodevops-api/pkg/agent/skill"
)

// CommandExecutionSkill 命令执行技能
type CommandExecutionSkill struct {
	baseSkill.BaseSkill
	hostDao    cmdbDao.CmdbHostDao
	ecsAuthDao monitorService.EcsAuthDao
	sshUtil    *util.SSHUtil
}

// NewCommandExecutionSkill 创建命令执行技能实例
func NewCommandExecutionSkill() *CommandExecutionSkill {
	return &CommandExecutionSkill{
		hostDao:    cmdbDao.NewCmdbHostDao(),
		ecsAuthDao: monitorService.NewEcsAuthDao(),
		sshUtil:    util.NewSSHUtil(),
	}
}

// Name 技能名称
func (s *CommandExecutionSkill) Name() string {
	return "command_execution"
}

// Description 技能描述
func (s *CommandExecutionSkill) Description() string {
	return "在指定服务器上执行Shell命令"
}

// Category 技能分类
func (s *CommandExecutionSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryOperation
}

// CanHandle 判断是否能处理该意图
func (s *CommandExecutionSkill) CanHandle(intent *baseSkill.Intent) bool {
	// 检测"执行"、"运行"等关键词，并且有"命令"关键词
	actionType := intent.Entities["actionType"]
	return (intent.Type == "incident" || intent.Type == "query") &&
		(actionType == "execute_command" || actionType == "run_command")
}

// Validate 验证参数
func (s *CommandExecutionSkill) Validate(req *baseSkill.SkillRequest) error {
	// 检查目标列表
	if len(req.TargetList) == 0 {
		targetName := req.Intent.Entities["targetName"]
		if targetName == "" {
			return fmt.Errorf("请指定要执行命令的主机IP或名称")
		}
	}

	// 检查命令
	command := req.Intent.Entities["command"]
	if fileOp := req.Intent.Entities["fileOp"]; fileOp != "" {
		filePath := req.Intent.Entities["filePath"]
		if err := validateManagedFileOperation(fileOp, filePath); err != nil {
			return err
		}
		req.Intent.Entities["command"] = buildManagedFileOperationCommand(fileOp, filePath)
		return nil
	}
	if command == "" {
		// 尝试从消息中提取命令
		command = s.extractCommand(req.Message)
		if command == "" {
			return fmt.Errorf("请指定要执行的命令")
		}
		req.Intent.Entities["command"] = command
	}

	return nil
}

// Execute 执行命令
func (s *CommandExecutionSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	command := req.Intent.Entities["command"]
	if fileOp := req.Intent.Entities["fileOp"]; fileOp != "" {
		filePath := req.Intent.Entities["filePath"]
		if err := validateManagedFileOperation(fileOp, filePath); err != nil {
			return nil, err
		}
		command = buildManagedFileOperationCommand(fileOp, filePath)
		req.Intent.Entities["command"] = command
	}

	// 安全检查：检测高风险命令
	if err := s.checkCommandSafety(command); err != nil {
		return nil, err
	}

	// 支持批量处理：如果 TargetList 有多个主机，批量处理
	targetList := req.TargetList
	if len(targetList) == 0 {
		// 兼容旧架构：从 entities 中提取
		targetName := req.Intent.Entities["targetName"]
		if targetName != "" {
			targetList = []string{targetName}
		}
	}

	if len(targetList) == 0 {
		return nil, fmt.Errorf("请指定要执行命令的主机IP或名称")
	}

	// 单主机处理（快速路径）
	if len(targetList) == 1 {
		host, err := s.getHostInfo(targetList[0])
		if err != nil {
			return nil, err
		}

		// 执行命令
		result, err := s.executeCommand(host, command)
		if err != nil {
			return nil, err
		}

		return &baseSkill.SkillResponse{
			Message: fmt.Sprintf("✅ 命令在主机 %s (%s) 上执行完成\n\n**命令**: `%s`\n\n**输出**:\n```\n%s\n```",
				host.HostName, host.SSHIP, command, result),
			Type: "text",
			TaskInfo: map[string]interface{}{
				"taskType": "command_execution",
				"hostId":   host.ID,
				"hostName": host.HostName,
				"hostIP":   host.SSHIP,
				"command":  command,
				"result":   result,
				"fileOp":   req.Intent.Entities["fileOp"],
				"filePath": req.Intent.Entities["filePath"],
			},
		}, nil
	}

	// 多主机批量处理
	var batchResults []map[string]interface{}
	var hostNames []string

	for _, target := range targetList {
		host, err := s.getHostInfo(target)
		if err != nil {
			// 记录错误但不中断
			batchResults = append(batchResults, map[string]interface{}{
				"target": target,
				"error":  err.Error(),
			})
			continue
		}

		// 执行命令
		result, err := s.executeCommand(host, command)
		if err != nil {
			batchResults = append(batchResults, map[string]interface{}{
				"hostName": host.HostName,
				"hostIP":   host.SSHIP,
				"error":    err.Error(),
			})
			continue
		}

		batchResults = append(batchResults, map[string]interface{}{
			"taskType": "command_execution",
			"hostId":   host.ID,
			"hostName": host.HostName,
			"hostIP":   host.SSHIP,
			"command":  command,
			"result":   result,
			"fileOp":   req.Intent.Entities["fileOp"],
			"filePath": req.Intent.Entities["filePath"],
		})
		hostNames = append(hostNames, fmt.Sprintf("%s(%s)", host.HostName, host.SSHIP))
	}

	return &baseSkill.SkillResponse{
		Message:     fmt.Sprintf("正在在 %d 台主机上执行命令：%s\n\n命令: `%s`", len(batchResults), strings.Join(hostNames, ", "), command),
		Type:        "batch",
		BatchResult: batchResults,
		Metadata: map[string]interface{}{
			"totalHosts": len(targetList),
			"command":    command,
		},
	}, nil
}

// getHostInfo 根据IP或名称获取主机信息
func (s *CommandExecutionSkill) getHostInfo(target string) (*cmdbModel.CmdbHost, error) {
	// 先尝试按IP查询
	hostByIP, err := s.hostDao.GetCmdbHostByIP(target)
	if err == nil {
		return &hostByIP, nil
	}

	// 再尝试按名称查询
	hostByName, err := s.hostDao.GetCmdbHostByName(target)
	if err == nil {
		return &hostByName, nil
	}

	// 都失败，返回友好错误
	if strings.Contains(err.Error(), "record not found") {
		return nil, fmt.Errorf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)
	}
	return nil, fmt.Errorf("查找主机 %s 失败: %v", target, err)
}

// executeCommand 在主机上执行命令
func (s *CommandExecutionSkill) executeCommand(host *cmdbModel.CmdbHost, command string) (string, error) {
	// 获取SSH认证信息
	ecsAuth, err := s.ecsAuthDao.GetEcsAuthById(host.SSHKeyID)
	if err != nil {
		return "", fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	// 构建SSH配置
	sshConfig := &util.SSHConfig{
		IP:       host.SSHIP,
		Port:     host.SSHPort,
		Username: host.SSHName,
		Type:     ecsAuth.Type,
	}

	// 根据认证类型设置认证信息
	switch ecsAuth.Type {
	case 1: // 密码认证
		sshConfig.Password = ecsAuth.Password
	case 2: // 私钥认证
		sshConfig.PublicKey = ecsAuth.PublicKey
	case 3: // 公钥免密
		// 无需额外字段，底层会自动查找本地私钥
	case 4: // 跳板机认证
		sshConfig.BastionHost = ecsAuth.BastionHost
		sshConfig.BastionPort = ecsAuth.BastionPort
		sshConfig.BastionUsername = ecsAuth.BastionUsername
		sshConfig.BastionAuthType = ecsAuth.BastionAuthType
		sshConfig.BastionPassword = ecsAuth.BastionPassword
		sshConfig.BastionPrivateKey = ecsAuth.BastionPrivateKey
		sshConfig.TargetAuthType = ecsAuth.TargetAuthType
		sshConfig.TargetPassword = ecsAuth.TargetPassword
		sshConfig.TargetPrivateKey = ecsAuth.TargetPrivateKey
	default:
		return "", fmt.Errorf("不支持的认证类型: %d", ecsAuth.Type)
	}

	// 执行命令
	result, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, command)
	if err != nil {
		return "", fmt.Errorf("命令执行失败: %v", err)
	}

	return result, nil
}

// extractCommand 从消息中提取命令
func (s *CommandExecutionSkill) extractCommand(message string) string {
	// 尝试提取命令的几种模式：
	// 1. "执行命令 xxx"
	// 2. "运行 xxx"
	// 3. "command=xxx"
	// 4. ` 符号包裹的命令

	messageLower := strings.ToLower(message)

	// 模式1: 执行命令 xxx
	if strings.Contains(messageLower, "执行命令") {
		parts := strings.SplitN(message, "执行命令", 2)
		if len(parts) > 1 {
			cmd := strings.TrimSpace(parts[1])
			// 去除末尾的 ip= 或 name= 参数
			if idx := strings.Index(cmd, " ip="); idx != -1 {
				cmd = strings.TrimSpace(cmd[:idx])
			}
			if idx := strings.Index(cmd, " name="); idx != -1 {
				cmd = strings.TrimSpace(cmd[:idx])
			}
			return cmd
		}
	}

	// 模式2: 运行 xxx
	if strings.Contains(messageLower, "运行") && !strings.Contains(messageLower, "运行的") {
		parts := strings.SplitN(message, "运行", 2)
		if len(parts) > 1 {
			cmd := strings.TrimSpace(parts[1])
			// 去除"命令"等词
			cmd = strings.TrimPrefix(cmd, "命令")
			cmd = strings.TrimSpace(cmd)
			// 去除末尾的 ip= 或 name= 参数
			if idx := strings.Index(cmd, " ip="); idx != -1 {
				cmd = strings.TrimSpace(cmd[:idx])
			}
			if idx := strings.Index(cmd, " name="); idx != -1 {
				cmd = strings.TrimSpace(cmd[:idx])
			}
			return cmd
		}
	}

	// 模式3: command=xxx
	if strings.Contains(message, "command=") {
		parts := strings.SplitN(message, "command=", 2)
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}

	// 模式4: ` 符号包裹
	if strings.Contains(message, "`") {
		start := strings.Index(message, "`")
		end := strings.LastIndex(message, "`")
		if start != -1 && end != -1 && start != end {
			return strings.TrimSpace(message[start+1 : end])
		}
	}

	return ""
}

// checkCommandSafety 检查命令安全性，拦截高风险命令
func (s *CommandExecutionSkill) checkCommandSafety(command string) error {
	if command == "" {
		return fmt.Errorf("命令不能为空")
	}
	if isManagedFileOperationCommand(command) {
		return nil
	}

	// 转换为小写以进行不区分大小写的匹配
	cmdLower := strings.ToLower(command)

	if err := validateDeleteCommand(command); err != nil {
		return err
	}

	// 定义高风险命令类别及其匹配规则
	type RiskPattern struct {
		Pattern     string
		Category    string
		Description string
	}

	riskPatterns := []RiskPattern{
		// 1. 系统破坏类
		{"rm -rf /", "系统破坏", "删除根目录"},
		{"rm -rf /*", "系统破坏", "删除根目录所有内容"},
		{"rm -rf --no-preserve-root /", "系统破坏", "强制删除根目录"},
		{"dd if=/dev/zero of=/dev/sd", "系统破坏", "清除磁盘数据"},
		{"dd if=/dev/zero of=/dev/hd", "系统破坏", "清除磁盘数据"},
		{"dd if=/dev/zero of=/dev/nvme", "系统破坏", "清除磁盘数据"},
		{"mkfs.", "系统破坏", "格式化文件系统"},
		{"fdisk", "系统破坏", "磁盘分区工具"},
		{"parted", "系统破坏", "磁盘分区工具"},
		{"cfdisk", "系统破坏", "磁盘分区工具"},
		{"wipefs", "系统破坏", "擦除文件系统签名"},

		// 2. 权限与提权类
		{"sudo", "权限提升", "超级用户权限"},
		{"su ", "权限提升", "切换用户"},
		{"su -", "权限提升", "切换到root用户"},
		{"passwd", "权限提升", "修改密码"},
		{"useradd", "权限提升", "添加用户"},
		{"userdel", "权限提升", "删除用户"},
		{"usermod", "权限提升", "修改用户"},
		{"groupadd", "权限提升", "添加用户组"},
		{"chmod 777", "权限提升", "设置危险的文件权限"},
		{"chmod -r 777", "权限提升", "递归设置危险权限"},
		{"chown root", "权限提升", "修改文件所有者为root"},
		{"chgrp root", "权限提升", "修改文件组为root"},
		{"visudo", "权限提升", "编辑sudo配置"},
		{"pkexec", "权限提升", "以root权限执行"},

		// 3. 敏感信息读取类
		{"cat /etc/shadow", "敏感信息", "读取密码哈希"},
		{"/etc/shadow", "敏感信息", "访问密码文件"},
		{"cat ~/.ssh/id_rsa", "敏感信息", "读取SSH私钥"},
		{"~/.ssh/id_rsa", "敏感信息", "访问SSH私钥"},
		{"/root/.ssh", "敏感信息", "访问root SSH密钥"},
		{"printenv", "敏感信息", "打印环境变量（可能包含密钥）"},

		// 4. 网络与远程控制类
		{"nc ", "网络控制", "netcat反向shell"},
		{"netcat", "网络控制", "netcat反向shell"},
		{"socat", "网络控制", "网络数据传输"},
		{"ncat", "网络控制", "nmap netcat"},
		{"-e /bin/sh", "网络控制", "执行shell"},
		{"-e /bin/bash", "网络控制", "执行bash"},
		{"mkfifo", "网络控制", "创建命名管道（反向shell）"},

		// 5. 进程与系统控制类
		{"kill -9 1", "系统控制", "杀死init进程"},
		{"kill -9 -1", "系统控制", "杀死所有进程"},
		{"killall -9", "系统控制", "强制杀死所有进程"},
		{"reboot", "系统控制", "重启系统"},
		{"shutdown", "系统控制", "关闭系统"},
		{"halt", "系统控制", "停止系统"},
		{"poweroff", "系统控制", "关闭电源"},
		{"init 0", "系统控制", "关闭系统"},
		{"init 6", "系统控制", "重启系统"},
		{"telinit", "系统控制", "更改运行级别"},
		{"sysctl -w", "系统控制", "修改内核参数"},
		{"modprobe", "系统控制", "加载内核模块"},
		{"insmod", "系统控制", "插入内核模块"},
		{"rmmod", "系统控制", "删除内核模块"},

		// 6. 资源耗尽类（DoS）
		{":(){ :|:& };:", "资源耗尽", "Fork炸弹"},
		{":(){", "资源耗尽", "Fork炸弹变体"},
		{"yes >", "资源耗尽", "CPU耗尽"},
		{"yes |", "资源耗尽", "CPU耗尽"},
		{"while true; do", "资源耗尽", "无限循环"},
		{"for((;;)); do", "资源耗尽", "无限循环"},
		{"stress-ng", "资源耗尽", "压力测试工具"},
		{"stress ", "资源耗尽", "压力测试工具"},
		{"/dev/urandom", "资源耗尽", "随机数耗尽资源"},

		// 7. 文件系统挂载/危险操作类
		{"mount ", "文件系统", "挂载文件系统"},
		{"umount ", "文件系统", "卸载文件系统"},
		{"> /dev/sd", "文件系统", "直接写入磁盘设备"},
		{"> /dev/hd", "文件系统", "直接写入磁盘设备"},
		{"> /dev/nvme", "文件系统", "直接写入磁盘设备"},
		{"dd of=/dev", "文件系统", "直接写入设备"},

		// 8. 危险的重定向和管道
		{"> /etc/", "文件覆盖", "覆盖系统配置文件"},
		{"> /bin/", "文件覆盖", "覆盖系统二进制文件"},
		{"> /usr/bin/", "文件覆盖", "覆盖系统二进制文件"},
		{"> /sbin/", "文件覆盖", "覆盖系统二进制文件"},
		{"> /boot/", "文件覆盖", "覆盖启动文件"},
	}

	// 检查命令是否匹配高风险模式
	for _, risk := range riskPatterns {
		if strings.Contains(cmdLower, risk.Pattern) {
			return fmt.Errorf("⛔ 安全拦截：该命令属于【%s】类高风险操作，已被拒绝执行\n"+
				"原因: %s\n"+
				"检测到的危险模式: %s\n"+
				"如确需执行，请手动登录服务器进行操作",
				risk.Category, risk.Description, risk.Pattern)
		}
	}

	// 额外检查：命令链中的危险组合
	if s.containsDangerousChaining(cmdLower) {
		return fmt.Errorf("⛔ 安全拦截：检测到危险的命令链或管道操作，已被拒绝执行\n" +
			"该命令可能会绕过安全检查或造成系统损害\n" +
			"如确需执行，请手动登录服务器进行操作")
	}

	// 通过安全检查
	return nil
}

func validateManagedFileOperation(fileOp, filePath string) error {
	if fileOp != "clear" && fileOp != "delete" {
		return fmt.Errorf("不支持的文件操作: %s", fileOp)
	}
	if !isStrictManagedFilePath(filePath) {
		return fmt.Errorf("文件路径不合法或不允许操作: %s", filePath)
	}
	if !isManagedFileSuffixAllowed(filePath) {
		return fmt.Errorf("仅允许操作 *.log、*.txt、*.gz 文件")
	}
	return nil
}

func validateDeleteCommand(command string) error {
	fields := strings.Fields(strings.TrimSpace(command))
	if len(fields) == 0 || fields[0] != "rm" {
		return nil
	}

	var targets []string
	for _, field := range fields[1:] {
		if strings.HasPrefix(field, "-") {
			if strings.ContainsAny(field, "rRd") {
				return fmt.Errorf("⛔ 安全拦截：AI 仅允许删除 .txt、.log、.gz 文件，禁止递归删除、目录删除或强破坏删除")
			}
			continue
		}
		targets = append(targets, field)
	}

	if len(targets) == 0 {
		return fmt.Errorf("⛔ 安全拦截：删除命令缺少明确的文件路径")
	}

	for _, target := range targets {
		if !isStrictManagedFilePath(target) {
			return fmt.Errorf("⛔ 安全拦截：删除文件路径不合法或包含不允许的字符: %s", target)
		}
		if strings.ContainsAny(target, "*?[]{}") {
			return fmt.Errorf("⛔ 安全拦截：不允许使用通配符批量删除: %s", target)
		}
		if ext := strings.ToLower(filepath.Ext(target)); ext != ".txt" && ext != ".log" && ext != ".gz" {
			return fmt.Errorf("⛔ 安全拦截：AI 仅允许删除 .txt、.log、.gz 文件，当前文件需手动登录服务器删除: %s", target)
		}
	}

	return nil
}

func buildManagedFileOperationCommand(fileOp, filePath string) string {
	switch fileOp {
	case "clear":
		return fmt.Sprintf("if test -f %s; then : > %s && printf 'file=%s\\naction=clear\\nstatus=cleared\\n'; else printf 'file=%s\\naction=clear\\nstatus=not_found\\n'; fi", filePath, filePath, filePath, filePath)
	case "delete":
		dirPath := path.Dir(filePath)
		return fmt.Sprintf("if test -f %s; then rm -f %s && printf 'file=%s\\naction=delete\\nstatus=deleted\\nremaining_dir=%s\\n'; ls -1A %s 2>/dev/null | sed 's/^/remaining_file=/'; else printf 'file=%s\\naction=delete\\nstatus=not_found\\n'; fi", filePath, filePath, filePath, dirPath, dirPath, filePath)
	default:
		return ""
	}
}

func isManagedFileOperationCommand(command string) bool {
	return strings.Contains(command, "\naction=clear") ||
		strings.Contains(command, "\naction=delete") ||
		strings.Contains(command, "action=clear\\n") ||
		strings.Contains(command, "action=delete\\n")
}

func isStrictManagedFilePath(path string) bool {
	if path == "" || !strings.HasPrefix(path, "/") {
		return false
	}
	for _, r := range path {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '/', r == '-', r == '_', r == '.', r == '@':
		default:
			return false
		}
	}
	return true
}

func isManagedFileSuffixAllowed(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	return strings.HasSuffix(lower, ".log") ||
		strings.HasSuffix(lower, ".txt") ||
		strings.HasSuffix(lower, ".gz")
}

// containsDangerousChaining 检查命令是否包含危险的命令链
func (s *CommandExecutionSkill) containsDangerousChaining(command string) bool {
	// 检查是否试图绕过限制的模式
	dangerousPatterns := []string{
		"bash -c", // 可能用于执行复杂的危险命令
		"sh -c",
		"eval ",  // 可能用于动态执行危险代码
		"exec ",  // 可能用于替换当前进程
		"$(curl", // 从远程下载并执行
		"$(wget",
		"`curl", // 反引号执行
		"`wget",
		"&&rm", // 命令链中包含删除
		";rm",
		"|rm",
		"&&dd", // 命令链中包含dd
		";dd",
		"|dd",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(command, pattern) {
			return true
		}
	}

	return false
}
