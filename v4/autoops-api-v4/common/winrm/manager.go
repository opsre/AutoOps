package winrm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/masterzen/winrm"
)

// Manager WinRM管理器
type Manager struct {
	Host     string
	Port     int
	Username string
	Password string
	UseHTTPS bool
	Timeout  time.Duration
	client   *winrm.Client
}

// Config WinRM配置
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	UseHTTPS bool
}

// NewManager 创建WinRM管理器
func NewManager(config Config) *Manager {
	port := config.Port
	if port == 0 {
		if config.UseHTTPS {
			port = 5986
		} else {
			port = 5985
		}
	}

	// 构建 WinRM endpoint
	endpoint := winrm.NewEndpoint(config.Host, port, config.UseHTTPS, false, nil, nil, nil, 60*time.Second)

	// 创建客户端
	client, err := winrm.NewClient(endpoint, config.Username, config.Password)
	if err != nil {
		log.Printf("创建 WinRM 客户端失败: %v", err)
		return nil
	}

	return &Manager{
		Host:     config.Host,
		Port:     port,
		Username: config.Username,
		Password: config.Password,
		UseHTTPS: config.UseHTTPS,
		Timeout:  60 * time.Second,
		client:   client,
	}
}

// InstallWindowsExporter 安装Windows Exporter（支持主动推送和被动采集两种模式）
func (m *Manager) InstallWindowsExporter(installerURL, version string, exporterPort int, pushgatewayURL string, pushInterval int, deviceName string) error {
	// 1. 下载安装包到临时目录
	log.Printf("   📥 正在下载安装包...")
	log.Printf("      URL: %s", installerURL)

	// 使用固定路径，避免路径解析问题
	installerPath := fmt.Sprintf("C:\\Temp\\windows_exporter-%s.msi", version)

	downloadCmd := fmt.Sprintf(`
# 创建临时目录
if (-Not (Test-Path "C:\Temp")) {
    New-Item -ItemType Directory -Path "C:\Temp" | Out-Null
}

# 下载文件
$installerPath = "%s"
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Write-Output "开始下载到: $installerPath"
Invoke-WebRequest -Uri "%s" -OutFile $installerPath

# 验证文件
if (Test-Path $installerPath) {
    $fileSize = (Get-Item $installerPath).Length / 1MB
    Write-Output "下载成功: $($fileSize.ToString('0.00')) MB"
    Write-Output "文件路径: $installerPath"
} else {
    Write-Output "下载失败: 文件不存在"
    exit 1
}
`, installerPath, installerURL)

	output, err := m.RunCommand(downloadCmd)
	if err != nil {
		return fmt.Errorf("下载安装包失败: %v, 输出: %s", err, output)
	}
	log.Printf("   ✅ 下载完成")
	log.Printf("      %s", strings.TrimSpace(output))

	// 2. 创建安装脚本文件（避免命令行长度限制）
	log.Printf("   📦 正在安装 Windows Exporter...")
	log.Printf("      版本: %s, 端口: %d", version, exporterPort)

	// 先创建安装脚本文件（分段写入避免命令行长度限制）
	scriptPath := "C:\\Temp\\install_exporter.ps1"

	// 第1部分：基础安装逻辑
	part1Cmd := fmt.Sprintf(`
$installerPath = "%s"
$exporterPort = %d
"开始安装: $installerPath" | Out-File -FilePath "%s" -Encoding UTF8
"检查安装包..." | Out-File -FilePath "%s" -Append
if (-Not (Test-Path $installerPath)) { "错误: 安装包不存在" | Out-File -FilePath "%s" -Append; exit 1 }
$fileSize = (Get-Item $installerPath).Length / 1MB
"安装包大小: $($fileSize.ToString('0.00')) MB" | Out-File -FilePath "%s" -Append
$logPath = "C:\Temp\windows_exporter_install.log"
    $arguments = @(
        "/i",
        $installerPath,
        "/qn", "/norestart",
        "/l*v", $logPath,
        "ENABLED_COLLECTORS=cpu,logical_disk,net,os,system,service,memory,process",
        "LISTEN_PORT=$exporterPort",
        'EXTRA_FLAGS="--collector.process.include=.*"'
    )
"执行MSI安装..." | Out-File -FilePath "%s" -Append
"MSI Arguments:" | Out-File -FilePath "%s" -Append
$arguments | ForEach-Object { "  ARG: $_" | Out-File -FilePath "%s" -Append }
$process = Start-Process "msiexec.exe" -ArgumentList $arguments -Wait -PassThru -NoNewWindow
"MSI退出码: $($process.ExitCode)" | Out-File -FilePath "%s" -Append
Start-Sleep -Seconds 3
`, installerPath, exporterPort, scriptPath, scriptPath, scriptPath, scriptPath, scriptPath, scriptPath, scriptPath, scriptPath)

	// 执行第1部分（创建脚本并执行安装）
	output, err = m.RunCommand(part1Cmd)
	if err != nil {
		return fmt.Errorf("MSI安装失败: %v, 输出: %s", err, output)
	}
	log.Printf("      MSI安装完成")

	// 第2部分：修复服务配置并启动
	part2Cmd := fmt.Sprintf(`
$exporterPort = %d
$service = Get-Service -Name "windows_exporter" -ErrorAction SilentlyContinue
if (-not $service) { Write-Output "错误: 服务未创建"; exit 1 }
Write-Output "服务已创建: $($service.Status)"

# 修复服务配置（移除 --config.file 参数，避免与 MSI 安装配置冲突）
$exePath = "C:\Program Files\windows_exporter\windows_exporter.exe"
    $serviceArgs = "--collectors.enabled cpu,logical_disk,net,os,system,service,memory,process --collector.process.include=.* --web.listen-address=:$exporterPort"
Stop-Service -Name "windows_exporter" -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 2

# 使用正确的 sc.exe 语法（注意转义引号）
$binPathWithQuotes = '\"' + $exePath + '\" ' + $serviceArgs
$result = sc.exe config windows_exporter binPath= $binPathWithQuotes
Write-Output "sc.exe 执行结果: $result"

# 验证服务配置
$serviceInfo = sc.exe qc windows_exporter
Write-Output "当前服务配置:"
Write-Output $serviceInfo

# 启动服务
try {
    Start-Service -Name "windows_exporter" -ErrorAction Stop
    Write-Output "启动命令已执行"
} catch {
    Write-Output "启动服务异常: $($_.Exception.Message)"
}

# 等待启动
for ($i = 0; $i -lt 15; $i++) {
    Start-Sleep -Seconds 1
    $service = Get-Service -Name "windows_exporter"
    if ($service.Status -eq 'Running') {
        Write-Output "服务已启动"
        exit 0
    }
}

# 服务未能启动，查看详细错误
$service = Get-Service -Name "windows_exporter"
Write-Output "服务未能启动: $($service.Status)"

# 查看事件日志错误（包含更详细的错误信息）
Write-Output "查看事件日志错误信息:"
Get-EventLog -LogName Application -Source "windows_exporter" -Newest 5 -ErrorAction SilentlyContinue |
    ForEach-Object { Write-Output "$($_.TimeGenerated) - $($_.EntryType) - $($_.Message)" }
Get-EventLog -LogName System -Source "Service Control Manager" -Newest 10 -ErrorAction SilentlyContinue |
    Where-Object {$_.Message -like "*windows_exporter*"} |
    ForEach-Object { Write-Output "$($_.TimeGenerated) - $($_.EntryType) - $($_.Message)" }

# 尝试手动执行看看错误（加上完整的启动参数）
Write-Output "尝试手动执行 exporter:"
$exePath = "C:\Program Files\windows_exporter\windows_exporter.exe"
if (Test-Path $exePath) {
    $testOutput = & $exePath --version 2>&1
    Write-Output "Exporter版本检测: $testOutput"
    Write-Output "尝试运行5秒测试启动:"
    $testRun = Start-Process -FilePath $exePath -ArgumentList $serviceArgs -NoNewWindow -PassThru
    Start-Sleep -Seconds 5
    if (-not $testRun.HasExited) {
        Write-Output "测试启动成功,进程正常运行"
        Stop-Process -Id $testRun.Id -Force
    } else {
        Write-Output "测试启动失败,进程退出码: $($testRun.ExitCode)"
    }
}

exit 1
`, exporterPort)

	//执行第2部分（修复服务并启动）
	output, err = m.RunCommand(part2Cmd)
	if err != nil {
		return fmt.Errorf("服务启动失败: %v, 输出: %s", err, output)
	}
	log.Printf("      %s", strings.TrimSpace(output))
	log.Printf("   ✅ Windows Exporter 安装并启动成功")

	// 3. 如果启用了Pushgateway推送模式，创建定时任务
	if pushgatewayURL != "" && pushInterval > 0 {
		log.Printf("   ⚙️  配置 Pushgateway 推送任务...")
		log.Printf("      目标: %s", pushgatewayURL)
		log.Printf("      间隔: %d 秒", pushInterval)

		pushScriptCmd := fmt.Sprintf(`
# 创建推送脚本
$scriptPath = "C:\Windows\System32\push_metrics_to_gateway.ps1"
$scriptContent = @'
$exporterUrl = "http://localhost:%d/metrics"
$deviceName = "%s"
$pushgatewayUrl = "%s/metrics/job/$deviceName/instance/$env:COMPUTERNAME"
try {
    $metrics = Invoke-WebRequest -Uri $exporterUrl -UseBasicParsing
    Invoke-WebRequest -Uri $pushgatewayUrl -Method Post -Body $metrics.Content -ContentType "text/plain" -UseBasicParsing | Out-Null
} catch {
    Write-Error "推送失败: $_"
}
'@
Set-Content -Path $scriptPath -Value $scriptContent -Force

# 创建计划任务（每%d秒执行一次）
$action = New-ScheduledTaskAction -Execute "powershell.exe" -Argument "-NoProfile -ExecutionPolicy Bypass -File $scriptPath"
$trigger = New-ScheduledTaskTrigger -Once -At (Get-Date) -RepetitionInterval (New-TimeSpan -Seconds %d)
$principal = New-ScheduledTaskPrincipal -UserId "SYSTEM" -LogonType ServiceAccount -RunLevel Highest
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable
Register-ScheduledTask -TaskName "PushMetricsToGateway" -Action $action -Trigger $trigger -Principal $principal -Settings $settings -Force
Write-Output "已配置Pushgateway推送任务"
`, exporterPort, deviceName, pushgatewayURL, pushInterval, pushInterval)

		_, err = m.RunCommand(pushScriptCmd)
		if err != nil {
			return fmt.Errorf("配置推送任务失败: %v", err)
		}
		log.Printf("   ✅ 推送任务配置完成")
	}

	// 4. 清理临时文件
	log.Printf("   🧹 清理临时文件...")
	cleanCmd := fmt.Sprintf(`Remove-Item -Path "%s" -Force -ErrorAction SilentlyContinue`, installerPath)
	m.RunCommand(cleanCmd)
	log.Printf("   ✅ 清理完成")

	return nil
}

// CheckExporterStatus 检查Windows Exporter状态
func (m *Manager) CheckExporterStatus() (bool, string, error) {
	cmd := `
$service = Get-Service -Name "windows_exporter" -ErrorAction SilentlyContinue
if ($service) {
    $version = (Get-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*" | Where-Object { $_.DisplayName -like "*windows_exporter*" }).DisplayVersion
    Write-Output "installed|$($service.Status)|$version"
} else {
    Write-Output "not_installed"
}
`
	output, err := m.RunCommand(cmd)
	if err != nil {
		return false, "", err
	}

	output = strings.TrimSpace(output)
	if output == "not_installed" {
		return false, "", nil
	}

	parts := strings.Split(output, "|")
	if len(parts) >= 3 {
		isRunning := parts[1] == "Running"
		version := parts[2]
		return isRunning, version, nil
	}

	return false, "", fmt.Errorf("无法解析服务状态")
}

// UninstallExporter 卸载Windows Exporter
func (m *Manager) UninstallExporter() error {
	cmd := `
$app = Get-WmiObject -Class Win32_Product | Where-Object { $_.Name -like "*windows_exporter*" }
if ($app) {
    $app.Uninstall()
    Write-Output "卸载成功"
} else {
    Write-Output "未找到安装程序"
}
`
	_, err := m.RunCommand(cmd)
	return err
}

// RunCommand 执行PowerShell命令
func (m *Manager) RunCommand(command string) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("WinRM 客户端未初始化")
	}

	// 使用 winrm 库执行 PowerShell 命令
	// 注意：使用 RunPSWithString 而不是 RunWithString
	stdout, stderr, exitCode, err := m.client.RunPSWithString(command, "")
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %v, stderr: %s", err, stderr)
	}

	if exitCode != 0 {
		return stdout, fmt.Errorf("命令返回非零退出码 %d: %s", exitCode, stderr)
	}

	return stdout, nil
}

// TestConnection 测试WinRM连接
func (m *Manager) TestConnection() error {
	output, err := m.RunCommand("Write-Output 'OK'")
	if err != nil {
		return err
	}

	if !strings.Contains(output, "OK") {
		return fmt.Errorf("连接测试失败: 未收到预期响应")
	}

	return nil
}
