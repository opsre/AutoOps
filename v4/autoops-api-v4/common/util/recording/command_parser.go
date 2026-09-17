// 命令解析器 - 从输出流提取完整命令(包含Tab补全)
// 通过解析终端回显获取实际执行的命令,而不是原始键盘输入
package recording

import (
	"bytes"
	"regexp"
	"strings"
)

// CommandParser 命令解析器
type CommandParser struct {
	inputBuffer        []byte          // 输入缓冲区
	outputBuffer       []byte          // 输出缓冲区
	commands           []ParsedCommand // 已解析的命令列表
	commandCount       int             // 命令计数器
	currentCommand     []rune          // 当前正在构建的命令(包含Tab补全)
	pendingTabComplete bool            // 是否正在等待Tab补全响应
	inFullScreenApp    bool            // 是否在全屏应用中(vim/less等)
}

// ParsedCommand 解析后的命令
type ParsedCommand struct {
	Command   string  // 命令文本
	Timestamp float64 // 相对时间戳
	Sequence  int     // 序号
}

// NewCommandParser 创建命令解析器
func NewCommandParser() *CommandParser {
	return &CommandParser{
		inputBuffer:        make([]byte, 0, 4096),
		outputBuffer:       make([]byte, 0, 8192),
		commands:           make([]ParsedCommand, 0),
		commandCount:       0,
		currentCommand:     make([]rune, 0, 256),
		pendingTabComplete: false,
		inFullScreenApp:    false,
	}
}

// ParseInput 解析输入数据,提取命令(支持Tab补全和退格键)
// 返回: 是否解析到完整命令
func (p *CommandParser) ParseInput(data []byte, timestamp float64) (bool, *ParsedCommand) {
	// 如果在全屏应用中(vim/less等),不记录命令
	if p.inFullScreenApp {
		return false, nil
	}

	// 检查是否包含退格键
	hasBackspace := bytes.Contains(data, []byte("\x08")) || bytes.Contains(data, []byte("\x7f"))

	// 检查是否包含Tab键
	hasTab := bytes.Contains(data, []byte("\t"))

	// 检测命令结束标记(回车或换行)
	hasNewline := bytes.Contains(data, []byte("\n")) ||
		bytes.Contains(data, []byte("\r"))

	// 处理退格键
	if hasBackspace {
		for _, b := range data {
			if b == '\x08' || b == '\x7f' { // 退格符
				if len(p.currentCommand) > 0 {
					p.currentCommand = p.currentCommand[:len(p.currentCommand)-1] // 删除最后一个字符
				}
			} else if b >= 32 && b <= 126 {
				// 普通字符,添加到命令
				p.currentCommand = append(p.currentCommand, rune(b))
			}
		}
		return false, nil
	}

	// 如果包含Tab键,设置等待标记,稍后从输出流获取补全内容
	if hasTab {
		p.pendingTabComplete = true
		// 将Tab之前的内容添加到当前命令
		for _, b := range data {
			if b == '\t' {
				continue // 跳过Tab键本身
			}
			if b == '\n' || b == '\r' {
				break
			}
			// 过滤控制字符
			if b >= 32 && b <= 126 {
				p.currentCommand = append(p.currentCommand, rune(b))
			}
		}
		return false, nil
	}

	// 如果是命令结束,提取完整命令
	if hasNewline {
		// 将剩余输入添加到当前命令
		for _, b := range data {
			if b == '\n' || b == '\r' {
				break
			}
			// 过滤控制字符
			if b >= 32 && b <= 126 {
				p.currentCommand = append(p.currentCommand, rune(b))
			}
		}

		// 构建完整命令文本
		cmdText := string(p.currentCommand)
		cmdText = strings.TrimSpace(cmdText)

		// 清空当前命令缓冲区
		p.currentCommand = p.currentCommand[:0]
		p.pendingTabComplete = false

		// 如果命令为空,不记录
		if cmdText == "" {
			return false, nil
		}

		// 检测是否启动全屏应用
		if isFullScreenApp(cmdText) {
			p.inFullScreenApp = true
		}

		// 创建命令记录
		p.commandCount++
		cmd := ParsedCommand{
			Command:   cmdText,
			Timestamp: timestamp,
			Sequence:  p.commandCount,
		}

		p.commands = append(p.commands, cmd)
		return true, &cmd
	}

	// 普通输入,添加到当前命令
	for _, b := range data {
		// 过滤控制字符
		if b >= 32 && b <= 126 {
			p.currentCommand = append(p.currentCommand, rune(b))
		}
	}

	return false, nil
}

// isFullScreenApp 检测是否为全屏应用命令
func isFullScreenApp(cmd string) bool {
	cmd = strings.TrimSpace(strings.ToLower(cmd))

	// 检测常见的全屏应用
	fullScreenApps := []string{
		"vim", "vi", "nano", "emacs",
		"less", "more",
		"top", "htop", "iotop",
		"man",
	}

	// 检查命令是否以这些应用开头
	for _, app := range fullScreenApps {
		if strings.HasPrefix(cmd, app+" ") || cmd == app {
			return true
		}
	}

	return false
}

// ParseOutput 解析SSH输出流,捕获Tab补全内容和检测全屏应用退出
// 返回: 新解析的命令列表
func (p *CommandParser) ParseOutput(data []byte, timestamp float64) []ParsedCommand {
	// 检测全屏应用退出的特征序列
	// \x1b[?1049l 是退出备用屏幕缓冲区的ANSI序列(vim/less等使用)
	if bytes.Contains(data, []byte("\x1b[?1049l")) {
		p.inFullScreenApp = false
		return []ParsedCommand{}
	}

	// 如果正在等待Tab补全响应
	if p.pendingTabComplete && !p.inFullScreenApp {
		// 清理ANSI转义序列
		cleaned := cleanOutputData(data)

		// 将补全内容添加到当前命令
		// Tab补全通常返回的是补全的后缀部分
		for _, ch := range cleaned {
			// 只接受可打印字符(避免控制字符)
			if ch >= 32 && ch <= 126 {
				p.currentCommand = append(p.currentCommand, ch)
			}
		}

		// 补全完成后立即清除等待标记
		// 避免后续的输出被误认为补全内容
		p.pendingTabComplete = false
	}

	// 目前不从输出流独立生成命令,只作为输入流的补充
	return []ParsedCommand{}
}

// cleanOutputData 清理输出数据中的ANSI转义序列和控制字符
func cleanOutputData(data []byte) []rune {
	// 移除ANSI转义序列
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cleaned := ansiRegex.ReplaceAll(data, []byte{})

	// 移除其他常见的控制序列
	// \x1b]... 序列
	oscRegex := regexp.MustCompile(`\x1b\][^\x07]*\x07`)
	cleaned = oscRegex.ReplaceAll(cleaned, []byte{})

	// 转换为rune数组,过滤控制字符
	result := []rune{}
	for _, b := range cleaned {
		// 跳过特定控制字符,但保留空格
		if b == '\r' || b == '\n' || b == '\t' || b == '\b' || b == '\x07' {
			continue
		}
		if b >= 32 && b <= 126 {
			result = append(result, rune(b))
		}
	}

	return result
}

// GetCommands 获取所有已解析的命令
func (p *CommandParser) GetCommands() []ParsedCommand {
	return p.commands
}

// GetCommandCount 获取命令数量
func (p *CommandParser) GetCommandCount() int {
	return p.commandCount
}

// Reset 重置解析器
func (p *CommandParser) Reset() {
	p.inputBuffer = p.inputBuffer[:0]
	p.outputBuffer = p.outputBuffer[:0]
	p.commands = p.commands[:0]
	p.commandCount = 0
}

// cleanCommand 清理命令文本
// 移除控制字符、ANSI转义序列和多余的空白
func cleanCommand(raw string) string {
	// 移除ANSI转义序列 (如颜色代码)
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cleaned := ansiRegex.ReplaceAllString(raw, "")

	// 移除退格符及其前面的字符
	cleaned = processBackspace(cleaned)

	// 移除其他控制字符(保留可打印字符)
	cleaned = removeControlChars(cleaned)

	// 移除首尾空白
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

// processBackspace 处理退格符(\x08或\x7f)
func processBackspace(s string) string {
	result := []rune{}

	for _, ch := range s {
		if ch == '\x08' || ch == '\x7f' { // 退格符
			if len(result) > 0 {
				result = result[:len(result)-1] // 删除前一个字符
			}
		} else {
			result = append(result, ch)
		}
	}

	return string(result)
}

// removeControlChars 移除控制字符(保留可打印字符和空白)
func removeControlChars(s string) string {
	result := []rune{}

	for _, ch := range s {
		// 保留可打印字符 (32-126) 和制表符(9)、换行(10)、回车(13)
		if (ch >= 32 && ch <= 126) || ch == 9 || ch == 10 || ch == 13 {
			result = append(result, ch)
		}
	}

	return string(result)
}

// IsSensitiveCommand 判断是否为敏感命令
func IsSensitiveCommand(cmd string) bool {
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	// 敏感命令关键词列表
	sensitiveKeywords := []string{
		"rm", "mkfs", "dd", "fdisk", "parted",
		"userdel", "groupdel", "passwd",
		"chmod 777", "chown",
		"iptables -f", "iptables -x",
		"systemctl stop", "systemctl disable",
		"kill -9", "pkill", "killall",
		"reboot", "shutdown", "halt", "poweroff",
		"wget", "curl", "nc", "netcat",
		"su ", "sudo ",
		"history -c", "history -w",
		"/etc/shadow", "/etc/passwd",
	}

	for _, keyword := range sensitiveKeywords {
		if strings.Contains(cmd, keyword) {
			return true
		}
	}

	return false
}

// AssessCommandRisk 评估命令风险等级
// 返回: (风险等级, 风险原因)
// 风险等级: 0-正常 1-可疑 2-高危
func AssessCommandRisk(cmd string) (int, string) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	// 高危命令模式
	criticalPatterns := map[string]string{
		`rm\s+-rf\s+/`:                      "尝试删除根目录",
		`rm\s+-rf\s+/\*`:                    "尝试删除根目录下所有文件",
		`rm\s+-rf\s+--no-preserve-root\s+/`: "尝试强制删除根目录",
		`mkfs`:                              "格式化磁盘操作",
		`dd.*if=/dev/zero.*of=/dev/`:        "覆写磁盘数据",
		`chmod\s+-r\s+777`:                  "递归赋予所有权限",
		`iptables\s+-f`:                     "清空防火墙规则",
		`history\s+-c`:                      "清除命令历史",
		`/etc/shadow`:                       "访问密码文件",
	}

	for pattern, reason := range criticalPatterns {
		if matched, _ := regexp.MatchString(pattern, cmd); matched {
			return 2, reason // 高危
		}
	}

	// 可疑命令模式
	suspiciousPatterns := map[string]string{
		`(wget|curl).*\|.*sh`:       "下载并执行未知脚本",
		`(wget|curl).*\|.*bash`:     "下载并执行未知脚本",
		`nc\s+-[el]`:                "使用netcat监听端口",
		`python.*-m\s+http\.server`: "启动HTTP服务器",
		`sudo\s+su`:                 "切换到root用户",
		`base64\s+-d`:               "解码Base64内容",
		`eval\s+`:                   "执行动态代码",
	}

	for pattern, reason := range suspiciousPatterns {
		if matched, _ := regexp.MatchString(pattern, cmd); matched {
			return 1, reason // 可疑
		}
	}

	// 检查是否包含敏感关键词(但不是高危或可疑)
	if IsSensitiveCommand(cmd) {
		return 1, "包含敏感操作关键词" // 可疑
	}

	return 0, "" // 正常
}
