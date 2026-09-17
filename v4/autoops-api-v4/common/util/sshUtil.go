package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHUtil struct{}

// NewSSHUtil 创建SSHUtil实例
func NewSSHUtil() *SSHUtil {
	return &SSHUtil{}
}

// ExecuteRemoteCommand 执行远程命令并返回输出
func (s *SSHUtil) ExecuteRemoteCommand(auth *SSHConfig, command string) (string, error) {
	var conn *ssh.Client
	var err error

	// 根据认证类型建立SSH连接
	if auth.Type == 4 {
		// 跳板机认证
		bastionAuth := &SSHConfig{
			IP:        auth.BastionHost,
			Port:      auth.BastionPort,
			Username:  auth.BastionUsername,
			Type:      auth.BastionAuthType,
			Password:  auth.BastionPassword,
			PublicKey: auth.BastionPrivateKey,
		}
		targetAuth := &SSHConfig{
			IP:        auth.IP,
			Port:      auth.Port,
			Username:  auth.Username,
			Type:      auth.TargetAuthType,
			Password:  auth.TargetPassword,
			PublicKey: auth.TargetPrivateKey,
		}
		conn, err = s.CreateBastionSSHClient(bastionAuth, targetAuth)
		if err != nil {
			return "", fmt.Errorf("failed to create bastion SSH client: %v", err)
		}
	} else {
		// 普通认证（type=1,2,3）
		config, err := s.GetSSHConfig(auth)
		if err != nil {
			return "", fmt.Errorf("failed to create SSH config: %v", err)
		}

		// 建立SSH连接
		conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", auth.IP, auth.Port), config)
		if err != nil {
			fmt.Printf("SSH连接失败: %v\n", err)
			return "", fmt.Errorf("failed to dial: %v", err)
		}
	}
	defer conn.Close()

	// 创建会话
	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// 设置输出缓冲区
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	// 执行命令
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return stdoutBuf.String(), nil
}

// GetSSHConfig 根据认证类型创建SSH配置
func (s *SSHUtil) GetSSHConfig(auth *SSHConfig) (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod

	fmt.Printf("[getSSHConfig] 创建SSH配置: type=%d, username=%s, ip=%s:%d\n",
		auth.Type, auth.Username, auth.IP, auth.Port)

	switch auth.Type {
	case 1: // 密码认证
		fmt.Printf("[getSSHConfig] 使用密码认证\n")
		authMethods = []ssh.AuthMethod{ssh.Password(auth.Password)}
	case 2: // 私钥认证 - 注意：PublicKey字段实际存储的是私钥内容
		fmt.Printf("[getSSHConfig] 使用私钥认证, 私钥长度=%d\n", len(auth.PublicKey))
		// 显示私钥的前后几个字符用于调试
		keyPreview := auth.PublicKey
		if len(keyPreview) > 100 {
			fmt.Printf("[getSSHConfig] 私钥开头: %s...\n", keyPreview[:80])
			fmt.Printf("[getSSHConfig] 私钥结尾: ...%s\n", keyPreview[len(keyPreview)-80:])
		}
		authMethod, err := s.PrivateKeyAuth(auth.PublicKey) // 这里的PublicKey实际是私钥
		if err != nil {
			fmt.Printf("[getSSHConfig] 私钥认证创建失败: %v\n", err)
			return nil, err
		}
		fmt.Printf("[getSSHConfig] 私钥认证方法创建成功\n")
		authMethods = []ssh.AuthMethod{authMethod}
	case 3: // 公钥免认证 - 公钥已提前拷贝到服务器，自动查找本地私钥
		// 按顺序尝试找到可用的私钥
		// 1. 先尝试用户主目录
		if userKeyAuth, err := s.UserKeyAuth(); err == nil {
			authMethods = append(authMethods, userKeyAuth)
		} else {
			// 2. 再尝试系统默认路径
			if defaultKeyAuth, err := s.DefaultKeyAuth(); err == nil {
				authMethods = append(authMethods, defaultKeyAuth)
			}
		}

		// 如果没有找到任何私钥，返回清晰的错误信息
		if len(authMethods) == 0 {
			return nil, errors.New("type=3 公钥免认证失败: 本地未找到私钥文件，请确保以下路径之一存在私钥:\n" +
				"用户目录: ~/.ssh/id_rsa, ~/.ssh/id_ed25519, ~/.ssh/id_ecdsa\n" +
				"系统目录: /root/.ssh/id_rsa, /root/.ssh/id_ed25519, /root/.ssh/id_ecdsa")
		}
	default:
		return nil, errors.New("unsupported authentication type")
	}

	return &ssh.ClientConfig{
		User: auth.Username,
		Auth: authMethods,
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil // 忽略主机密钥检查
		}),
		Timeout: 30 * time.Second,
	}, nil
}

// PrivateKeyAuth 使用私钥创建认证方法
// 支持 PEM 格式（BEGIN RSA PRIVATE KEY）和 OpenSSH 格式（BEGIN OPENSSH PRIVATE KEY）
func (s *SSHUtil) PrivateKeyAuth(privateKey string) (ssh.AuthMethod, error) {
	// 清理私钥字符串（去除可能的前后空白）
	privateKey = strings.TrimSpace(privateKey)

	// 检测私钥格式
	keyFormat := "unknown"
	if strings.Contains(privateKey, "BEGIN RSA PRIVATE KEY") {
		keyFormat = "PEM RSA"
	} else if strings.Contains(privateKey, "BEGIN OPENSSH PRIVATE KEY") {
		keyFormat = "OpenSSH"
	} else if strings.Contains(privateKey, "BEGIN EC PRIVATE KEY") {
		keyFormat = "PEM EC"
	} else if strings.Contains(privateKey, "BEGIN DSA PRIVATE KEY") {
		keyFormat = "PEM DSA"
	}
	fmt.Printf("检测到私钥格式: %s, 长度: %d\n", keyFormat, len(privateKey))

	// 方法1：尝试直接解析（支持无密码的 RSA/DSA/ECDSA/ED25519 和 OpenSSH 格式）
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err == nil {
		fmt.Printf("私钥解析成功（无密码）\n")
		return ssh.PublicKeys(signer), nil
	}
	fmt.Printf("私钥解析失败（无密码）: %v\n", err)

	// 方法2：尝试使用空密码解析（支持有密码保护但密码为空的情况）
	signer, err2 := ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(""))
	if err2 == nil {
		fmt.Printf("私钥解析成功（空密码）\n")
		return ssh.PublicKeys(signer), nil
	}
	fmt.Printf("私钥解析失败（空密码）: %v\n", err2)

	// 如果都失败，返回详细错误信息
	return nil, fmt.Errorf("failed to parse private key (format: %s, length: %d): %v (also tried with empty passphrase: %v)",
		keyFormat, len(privateKey), err, err2)
}

// PublicKeyAuth 保持向后兼容，实际调用PrivateKeyAuth
func (s *SSHUtil) PublicKeyAuth(privateKey string) (ssh.AuthMethod, error) {
	return s.PrivateKeyAuth(privateKey)
}

// DefaultKeyAuth 尝试使用默认路径的私钥文件
func (s *SSHUtil) DefaultKeyAuth() (ssh.AuthMethod, error) {
	// 常见的默认私钥路径
	defaultKeyPaths := []string{
		"/root/.ssh/id_rsa",
		"/root/.ssh/id_ed25519",
		"/root/.ssh/id_ecdsa",
	}

	for _, keyPath := range defaultKeyPaths {
		if keyData, err := os.ReadFile(keyPath); err == nil {
			// 尝试解析私钥
			if signer, err := ssh.ParsePrivateKey(keyData); err == nil {
				return ssh.PublicKeys(signer), nil
			}
			// 尝试无密码解析
			if signer, err := ssh.ParsePrivateKeyWithPassphrase(keyData, []byte("")); err == nil {
				return ssh.PublicKeys(signer), nil
			}
		}
	}

	return nil, errors.New("no default SSH key found")
}

// UserKeyAuth 尝试使用用户主目录的私钥文件
func (s *SSHUtil) UserKeyAuth() (ssh.AuthMethod, error) {
	// 获取当前用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	// 用户主目录下的私钥路径
	userKeyPaths := []string{
		homeDir + "/.ssh/id_rsa",
		homeDir + "/.ssh/id_ed25519",
		homeDir + "/.ssh/id_ecdsa",
	}

	for _, keyPath := range userKeyPaths {
		if keyData, err := os.ReadFile(keyPath); err == nil {
			// 尝试解析私钥
			if signer, err := ssh.ParsePrivateKey(keyData); err == nil {
				return ssh.PublicKeys(signer), nil
			}
			// 尝试无密码解析
			if signer, err := ssh.ParsePrivateKeyWithPassphrase(keyData, []byte("")); err == nil {
				return ssh.PublicKeys(signer), nil
			}
		}
	}

	return nil, errors.New("no user SSH key found")
}

// ExecuteScript 执行远程脚本
func (s *SSHUtil) ExecuteScript(auth *SSHConfig, script string) (string, error) {
	return s.ExecuteRemoteCommand(auth, script)
}

// TerminalLogin 建立SSH终端连接
func (s *SSHUtil) TerminalLogin(auth *SSHConfig) (*ssh.Client, error) {
	config, err := s.GetSSHConfig(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH config: %v", err)
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", auth.IP, auth.Port), config)
}

// UploadFile 上传文件到远程主机
func (s *SSHUtil) UploadFile(auth *SSHConfig, localPath, remotePath string) error {
	client, err := s.TerminalLogin(auth)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// 使用scp命令上传文件
	cmd := fmt.Sprintf("scp %s %s@%s:%s", localPath, auth.Username, auth.IP, remotePath)
	_, err = s.ExecuteRemoteCommand(auth, cmd)
	return err
}

// GetSystemInfo 获取远程主机系统信息
func (s *SSHUtil) GetSystemInfo(auth *SSHConfig) (map[string]string, error) {
	// 1. 测试基本连接
	if _, err := s.ExecuteRemoteCommand(auth, "echo 'Testing SSH connection...'"); err != nil {
		return nil, fmt.Errorf("SSH connection test failed: %v", err)
	}

	// 2. 创建并执行脚本
	output, err := s.ExecuteRemoteCommand(auth, getScriptContent())
	if err != nil {
		return nil, fmt.Errorf("failed to execute script: %v", err)
	}

	// 3. 解析结果
	var info map[string]string
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		return nil, fmt.Errorf("invalid script output: %v (output: %s)", err, output)
	}

	// 4. 验证必要字段
	requiredFields := []string{"privateIp", "os", "cpu", "memory", "disk", "name"}
	var missingFields []string
	for _, field := range requiredFields {
		if _, ok := info[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) > 0 {
		return info, fmt.Errorf("missing required fields: %v (partial info: %+v)", missingFields, info)
	}

	return info, nil
}

// GetScriptContent 返回可直接执行的脚本命令（导出供外部使用）
func (s *SSHUtil) GetScriptContent() string {
	return getScriptContent()
}

// ParseSystemInfo 解析系统信息输出为map
func ParseSystemInfo(output string) (map[string]string, error) {
	var info map[string]string
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		return nil, fmt.Errorf("invalid script output: %v (output: %s)", err, output)
	}

	// 验证必要字段
	requiredFields := []string{"privateIp", "os", "cpu", "memory", "disk", "name"}
	var missingFields []string
	for _, field := range requiredFields {
		if _, ok := info[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) > 0 {
		return info, fmt.Errorf("missing required fields: %v (partial info: %+v)", missingFields, info)
	}

	return info, nil
}

// getScriptContent 返回可直接执行的脚本命令
func getScriptContent() string {
	return `#!/bin/bash

# 获取系统信息
privateIp=$(hostname -I | awk '{print $1}' || echo "unknown")
publicIp=$(curl -s ipinfo.io/ip 2>/dev/null || echo "")
os=$(cat /etc/os-release 2>/dev/null | grep PRETTY_NAME | cut -d= -f2 | tr -d '"' | sed 's/ LTS//;s/ //g' || echo "unknown")

# 获取CPU核心数(去掉单位)
cpu=$(nproc 2>/dev/null || echo "unknown")

# 获取内存大小(去掉单位)
memory=$(free -m 2>/dev/null | awk '/^Mem:/{printf "%.0f\n",$2/1024}' || echo "unknown")

# 获取磁盘总容量(去掉单位)
disk=$(df -h 2>/dev/null | awk '/\/$/ {print $2}' | sed 's/G//' || echo "unknown")

# 确保获取主机名称
name=$(hostname 2>/dev/null)
if [ -z "$name" ]; then
    name=$(uname -n 2>/dev/null || echo "unknown")
fi

# 输出JSON结果
echo '{"privateIp":"'$privateIp'","publicIp":"'$publicIp'","os":"'$os'","cpu":"'$cpu'","memory":"'$memory'","disk":"'$disk'","name":"'$name'"}'
`
}

// SSHExec 执行SSH命令
func SSHExec(ip string, port int, user string, password string, command string) (string, error) {
	sshUtil := NewSSHUtil()
	auth := &SSHConfig{
		IP:       ip,
		Port:     port,
		Username: user,
		Password: password,
		Type:     1, // 密码认证
	}
	return sshUtil.ExecuteRemoteCommand(auth, command)
}

// GetSSHClientByConfig 通过 SSHConfig 获取SSH客户端连接
func GetSSHClientByConfig(auth *SSHConfig) (*ssh.Client, error) {
	sshUtil := NewSSHUtil()
	config, err := sshUtil.GetSSHConfig(auth)
	if err != nil {
		return nil, fmt.Errorf("创建SSH配置失败: %v", err)
	}

	// 建立SSH连接
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", auth.IP, auth.Port), config)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %v", err)
	}

	return client, nil
}

// CreateBastionSSHClient 创建跳板机SSH连接
// 先连接到跳板机，再从跳板机连接到目标主机
func (s *SSHUtil) CreateBastionSSHClient(bastionAuth, targetAuth *SSHConfig) (*ssh.Client, error) {
	fmt.Printf("[CreateBastionSSHClient] 开始创建跳板机连接\n")
	fmt.Printf("[CreateBastionSSHClient] 跳板机: %s:%d (type=%d, user=%s)\n",
		bastionAuth.IP, bastionAuth.Port, bastionAuth.Type, bastionAuth.Username)
	fmt.Printf("[CreateBastionSSHClient] 目标主机: %s:%d (type=%d, user=%s)\n",
		targetAuth.IP, targetAuth.Port, targetAuth.Type, targetAuth.Username)

	// 1. 先连接到跳板机
	fmt.Printf("[CreateBastionSSHClient] 步骤1: 创建跳板机SSH配置\n")
	bastionConfig, err := s.GetSSHConfig(bastionAuth)
	if err != nil {
		return nil, fmt.Errorf("创建跳板机SSH配置失败: %v", err)
	}

	fmt.Printf("[CreateBastionSSHClient] 步骤2: 连接到跳板机\n")
	bastionClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", bastionAuth.IP, bastionAuth.Port), bastionConfig)
	if err != nil {
		return nil, fmt.Errorf("连接跳板机失败(%s:%d): %v", bastionAuth.IP, bastionAuth.Port, err)
	}
	fmt.Printf("[CreateBastionSSHClient] 跳板机连接成功\n")

	// 2. 从跳板机连接到目标主机
	// 特殊处理：如果目标认证类型是4，表示使用跳板机的身份（Agent Forwarding）
	if targetAuth.Type == 4 {
		fmt.Printf("[CreateBastionSSHClient] 使用跳板机代理认证（目标机器使用跳板机的身份）\n")
		// 通过跳板机的SSH Agent进行认证
		// 这需要在跳板机上执行 ssh user@target 来建立连接
		return s.createBastionProxyClient(bastionClient, targetAuth)
	}

	fmt.Printf("[CreateBastionSSHClient] 步骤3: 创建目标主机SSH配置\n")
	targetConfig, err := s.GetSSHConfig(targetAuth)
	if err != nil {
		bastionClient.Close()
		return nil, fmt.Errorf("创建目标主机SSH配置失败: %v", err)
	}

	// 从跳板机建立到目标主机的连接
	targetAddr := fmt.Sprintf("%s:%d", targetAuth.IP, targetAuth.Port)
	fmt.Printf("[CreateBastionSSHClient] 步骤4: 从跳板机连接目标主机 %s\n", targetAddr)
	bastionConn, err := bastionClient.Dial("tcp", targetAddr)
	if err != nil {
		bastionClient.Close()
		return nil, fmt.Errorf("从跳板机连接目标主机失败(%s): %v", targetAddr, err)
	}
	fmt.Printf("[CreateBastionSSHClient] 目标主机TCP连接建立成功\n")

	// 3. 在跳板机连接的基础上建立SSH客户端
	fmt.Printf("[CreateBastionSSHClient] 步骤5: 进行SSH握手\n")
	ncc, chans, reqs, err := ssh.NewClientConn(bastionConn, targetAddr, targetConfig)
	if err != nil {
		bastionConn.Close()
		bastionClient.Close()
		return nil, fmt.Errorf("创建目标主机SSH客户端失败: %v", err)
	}
	fmt.Printf("[CreateBastionSSHClient] SSH握手成功\n")

	// 返回目标主机的SSH客户端
	targetClient := ssh.NewClient(ncc, chans, reqs)
	fmt.Printf("[CreateBastionSSHClient] 跳板机SSH客户端创建完成\n")

	// 注意：这里不关闭bastionClient，因为targetClient依赖它
	// 当targetClient关闭时，bastionClient也会被关闭
	return targetClient, nil
}

// createBastionProxyClient 通过跳板机代理创建到目标主机的连接
// 这种方式让跳板机使用自己的密钥来连接目标机器（相当于在跳板机上执行 ssh target）
// 返回一个包装的客户端，实际上是通过跳板机转发命令到目标机器
func (s *SSHUtil) createBastionProxyClient(bastionClient *ssh.Client, targetAuth *SSHConfig) (*ssh.Client, error) {
	fmt.Printf("[createBastionProxyClient] 通过跳板机免密连接目标: %s@%s:%d\n",
		targetAuth.Username, targetAuth.IP, targetAuth.Port)

	// 先测试跳板机能否免密登录目标机器
	testCmd := fmt.Sprintf("ssh -o StrictHostKeyChecking=no -o BatchMode=yes -o ConnectTimeout=5 %s@%s 'echo ok' 2>&1",
		targetAuth.Username, targetAuth.IP)
	fmt.Printf("[createBastionProxyClient] 测试跳板机免密登录: %s\n", testCmd)

	session, err := bastionClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建测试会话失败: %v", err)
	}
	output, err := session.CombinedOutput(testCmd)
	session.Close()

	fmt.Printf("[createBastionProxyClient] 测试结果: %s\n", string(output))
	if err != nil || !strings.Contains(string(output), "ok") {
		return nil, fmt.Errorf("跳板机无法免密登录目标机器，请确保跳板机的公钥已添加到目标机器的 ~/.ssh/authorized_keys。错误: %v", err)
	}
	fmt.Printf("[createBastionProxyClient] 免密登录测试成功\n")

	// 使用跳板机直接执行SSH命令的方式
	// 通过跳板机建立到目标的实际SSH连接
	targetAddr := fmt.Sprintf("%s:%d", targetAuth.IP, targetAuth.Port)

	// 使用nc (netcat) 或者直接的TCP转发来建立隧道
	bastionConn, err := bastionClient.Dial("tcp", targetAddr)
	if err != nil {
		return nil, fmt.Errorf("从跳板机建立TCP连接到目标失败: %v", err)
	}
	fmt.Printf("[createBastionProxyClient] TCP隧道建立成功\n")

	// 从跳板机读取私钥内容
	// 尝试使用跳板机上的私钥文件
	var signer ssh.Signer
	for _, keyPath := range []string{"~/.ssh/id_rsa", "~/.ssh/id_ed25519", "~/.ssh/id_ecdsa"} {
		readKeyCmd := fmt.Sprintf("cat %s 2>/dev/null", keyPath)
		session, err := bastionClient.NewSession()
		if err != nil {
			continue
		}
		keyContent, err := session.CombinedOutput(readKeyCmd)
		session.Close()

		if err == nil && len(keyContent) > 0 {
			fmt.Printf("[createBastionProxyClient] 从跳板机读取私钥: %s (长度=%d)\n", keyPath, len(keyContent))
			// 尝试解析私钥
			signer, err = ssh.ParsePrivateKey(keyContent)
			if err == nil {
				fmt.Printf("[createBastionProxyClient] 私钥解析成功\n")
				break
			}
			fmt.Printf("[createBastionProxyClient] 私钥解析失败: %v，尝试下一个\n", err)
		}
	}

	if signer == nil {
		bastionConn.Close()
		return nil, fmt.Errorf("无法从跳板机读取私钥文件")
	}

	// 使用读取到的私钥建立SSH连接
	targetConfig := &ssh.ClientConfig{
		User: targetAuth.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	fmt.Printf("[createBastionProxyClient] 通过隧道进行SSH握手\n")
	ncc, chans, reqs, err := ssh.NewClientConn(bastionConn, targetAddr, targetConfig)
	if err != nil {
		bastionConn.Close()
		return nil, fmt.Errorf("SSH握手失败: %v", err)
	}
	fmt.Printf("[createBastionProxyClient] SSH连接建立成功\n")

	targetClient := ssh.NewClient(ncc, chans, reqs)
	return targetClient, nil
}

// ExecuteSSHCommand 在SSH客户端上执行命令（不返回输出）
func ExecuteSSHCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	return session.Run(command)
}

// ExecuteSSHCommandWithOutput 在SSH客户端上执行命令并返回输出
func ExecuteSSHCommandWithOutput(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(command)
	output := stdoutBuf.String()
	if err != nil {
		// 即使出错也返回输出，方便调试
		return output + "\nSTDERR:\n" + stderrBuf.String(), err
	}

	return output, nil
}

// DownloadFile 从远程主机下载文件到本地
func (s *SSHUtil) DownloadFile(auth *SSHConfig, remotePath, localPath string) error {
	// 导入sftp包
	client, err := s.TerminalLogin(auth)
	if err != nil {
		return fmt.Errorf("SSH连接失败: %v", err)
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := createSFTPClient(client)
	if err != nil {
		return fmt.Errorf("创建SFTP客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 打开远程文件
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("打开远程文件失败: %v", err)
	}
	defer remoteFile.Close()

	// 创建本地文件
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer localFile.Close()

	// 复制文件内容
	_, err = remoteFile.WriteTo(localFile)
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}

	return nil
}

// createSFTPClient 创建SFTP客户端（内部辅助函数）
func createSFTPClient(client *ssh.Client) (*sftp.Client, error) {
	return sftp.NewClient(client)
}

// FileInfo 文件信息结构体
type FileInfo struct {
	Name    string `json:"name"`    // 文件名
	Path    string `json:"path"`    // 完整路径
	Size    int64  `json:"size"`    // 文件大小（字节）
	Mode    string `json:"mode"`    // 权限模式
	ModTime string `json:"modTime"` // 修改时间
	IsDir   bool   `json:"isDir"`   // 是否是目录
}

// ListFiles 列出远程目录下的文件
func (s *SSHUtil) ListFiles(auth *SSHConfig, remotePath string) ([]FileInfo, error) {
	client, err := s.TerminalLogin(auth)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %v", err)
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := createSFTPClient(client)
	if err != nil {
		return nil, fmt.Errorf("创建SFTP客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 读取目录
	files, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	// 转换为 FileInfo 结构
	var fileInfos []FileInfo
	for _, file := range files {
		fullPath := remotePath
		if !strings.HasSuffix(fullPath, "/") {
			fullPath += "/"
		}
		fullPath += file.Name()

		fileInfos = append(fileInfos, FileInfo{
			Name:    file.Name(),
			Path:    fullPath,
			Size:    file.Size(),
			Mode:    file.Mode().String(),
			ModTime: file.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:   file.IsDir(),
		})
	}

	return fileInfos, nil
}

