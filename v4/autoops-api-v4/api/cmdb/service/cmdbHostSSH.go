package service

import (
	"dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	configModel "dodevops-api/api/configcenter/model"
	"dodevops-api/common"
	"dodevops-api/common/util"
	"dodevops-api/common/util/websocket"
	"dodevops-api/pkg/jwt"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type CmdbHostSSHServiceInterface interface {
	ConnectTerminal(c *gin.Context, hostID uint) (*websocket.WebSSH, error)
	ExecuteCommand(c *gin.Context, hostID uint, command string) (*CommandResponse, error)
	UploadFile(c *gin.Context, hostID uint, filePath string, destPath string) error
	CreateUploadTask(c *gin.Context, hostID uint, filePath string, originalFileName string, destPath string) (*cmdbModel.HostSSHUploadTask, error)
	GetUploadTask(taskID string) (*cmdbModel.HostSSHUploadTask, error)
	DownloadFile(c *gin.Context, hostID uint, remotePath string, localPath string) error
	ListFiles(c *gin.Context, hostID uint, remotePath string) ([]util.FileInfo, error)
}

type CmdbHostSSHServiceImpl struct{}

var (
	hostSSHUploadTasks   sync.Map
	hostSSHUploadTaskSeq uint64
)

func (s *CmdbHostSSHServiceImpl) ConnectTerminal(c *gin.Context, hostID uint) (*websocket.WebSSH, error) {
	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		log.Printf("获取主机ID=%d信息失败: %v", hostID, err)
		return nil, fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		log.Printf("主机ID=%d SSH信息不完整", hostID)
		return nil, errors.New("主机SSH信息不完整")
	}

	// 查询认证类型
	var auth configModel.EcsAuth
	if host.SSHKeyID != 0 {
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			log.Printf("获取认证类型失败: %v", err)
			return nil, fmt.Errorf("获取认证类型失败: %v", err)
		}
	} else {
		return nil, errors.New("主机未配置SSH凭据")
	}

	// 如果是跳板机认证(type=4)，使用跳板机连接
	if auth.Type == 4 {
		log.Printf("使用跳板机认证方式连接WebSSH终端")

		sshUtil := util.NewSSHUtil()

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

		log.Printf("跳板机: %s:%d, 目标主机: %s:%d",
			bastionAuth.IP, bastionAuth.Port, targetAuth.IP, targetAuth.Port)

		// 通过跳板机建立SSH客户端
		client, err := sshUtil.CreateBastionSSHClient(bastionAuth, targetAuth)
		if err != nil {
			log.Printf("创建跳板机SSH连接失败: %v", err)
			return nil, fmt.Errorf("创建跳板机SSH连接失败: %v", err)
		}

		// 使用已建立的SSH客户端创建WebSSH
		webSSH, err := websocket.NewWebSSHFromClient(client)
		if err != nil {
			client.Close()
			log.Printf("创建WebSSH失败: %v", err)
			return nil, fmt.Errorf("创建WebSSH失败: %v", err)
		}

		log.Printf("跳板机WebSSH连接创建成功")
		return webSSH, nil
	}

	credential, err := dao.NewCmdbHostSSHDao().GetSSHCredentials(host.SSHKeyID)
	if err != nil {
		log.Printf("获取SSH凭据失败: %v", err)
		return nil, fmt.Errorf("获取SSH凭据失败: %v", err)
	}

	var authMethod ssh.AuthMethod
	// 根据认证类型创建认证方法
	switch auth.Type {
	case 1: // 密码认证
		authMethod = ssh.Password(credential)
	case 2: // 密钥认证
		authMethod, err = util.NewSSHUtil().PublicKeyAuth(credential)
		if err != nil {
			log.Printf("创建密钥认证失败: %v", err)
			return nil, fmt.Errorf("创建密钥认证失败: %v", err)
		}
	case 3: // 公钥免认证 - 自动查找本地私钥
		// 先尝试用户主目录
		if userKeyAuth, err := util.NewSSHUtil().UserKeyAuth(); err == nil {
			authMethod = userKeyAuth
		} else {
			// 再尝试系统默认路径
			if defaultKeyAuth, err := util.NewSSHUtil().DefaultKeyAuth(); err == nil {
				authMethod = defaultKeyAuth
			} else {
				log.Printf("type=3 公钥免认证失败: 本地未找到私钥文件")
				return nil, errors.New("type=3 公钥免认证失败: 本地未找到私钥文件")
			}
		}
	default:
		return nil, errors.New("不支持的认证方式")
	}

	webSSH, err := websocket.NewWebSSHWithAuth(host.SSHIP, strconv.Itoa(host.SSHPort), host.SSHName, authMethod)
	if err != nil {
		log.Printf("创建WebSSH连接失败: %v", err)
		return nil, fmt.Errorf("创建WebSSH连接失败: %v", err)
	}

	return webSSH, nil
}

type CommandResponse struct {
	Command string `json:"command"`
	HostID  string `json:"hostId"`
	Output  string `json:"output"`
}

func nextHostSSHUploadTaskID() string {
	seq := atomic.AddUint64(&hostSSHUploadTaskSeq, 1)
	return fmt.Sprintf("ssh-upload-%d-%d", time.Now().UnixNano(), seq)
}

func cloneHostSSHUploadTask(task *cmdbModel.HostSSHUploadTask) *cmdbModel.HostSSHUploadTask {
	if task == nil {
		return nil
	}
	copy := *task
	return &copy
}

func (s *CmdbHostSSHServiceImpl) ExecuteCommand(c *gin.Context, hostID uint, command string) (*CommandResponse, error) {
	// 创建独立的SSH连接，不使用WebSSH
	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		return nil, fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		return nil, errors.New("主机SSH信息不完整")
	}

	var auth configModel.EcsAuth
	if host.SSHKeyID != 0 {
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			return nil, fmt.Errorf("获取认证类型失败: %v", err)
		}
	} else {
		return nil, errors.New("主机未配置SSH凭据")
	}

	var client *ssh.Client

	// 根据认证类型创建SSH客户端
	switch auth.Type {
	case 4: // 跳板机认证
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

		// 通过跳板机连接
		client, err = util.NewSSHUtil().CreateBastionSSHClient(bastionAuth, targetAuth)
		if err != nil {
			return nil, fmt.Errorf("通过跳板机连接失败: %v", err)
		}

	default: // 其他认证方式 (1密码, 2私钥, 3免密)
		credential, err := dao.NewCmdbHostSSHDao().GetSSHCredentials(host.SSHKeyID)
		if err != nil {
			return nil, fmt.Errorf("获取SSH凭据失败: %v", err)
		}

		var authMethod ssh.AuthMethod
		switch auth.Type {
		case 1: // 密码认证
			authMethod = ssh.Password(credential)
		case 2: // 密钥认证
			authMethod, err = util.NewSSHUtil().PublicKeyAuth(credential)
			if err != nil {
				return nil, fmt.Errorf("创建密钥认证失败: %v", err)
			}
		case 3: // 公钥免认证 - 自动查找本地私钥
			// 先尝试用户主目录
			if userKeyAuth, err := util.NewSSHUtil().UserKeyAuth(); err == nil {
				authMethod = userKeyAuth
			} else {
				// 再尝试系统默认路径
				if defaultKeyAuth, err := util.NewSSHUtil().DefaultKeyAuth(); err == nil {
					authMethod = defaultKeyAuth
				} else {
					return nil, errors.New("type=3 公钥免认证失败: 本地未找到私钥文件")
				}
			}
		default:
			return nil, errors.New("不支持的认证方式")
		}

		config := &ssh.ClientConfig{
			User:            host.SSHName,
			Auth:            []ssh.AuthMethod{authMethod},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         15 * time.Second,
		}

		client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort), config)
		if err != nil {
			return nil, fmt.Errorf("创建SSH连接失败: %v", err)
		}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 执行命令并获取输出
	output, err := session.CombinedOutput(command)
	if err != nil {
		return nil, fmt.Errorf("命令执行失败: %v", err)
	}

	// 清理输出
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	cleanOutput := ansiRegex.ReplaceAllString(string(output), "")

	promptRegex := regexp.MustCompile(`.*@.*:[~#]\s*`)
	parts := promptRegex.Split(cleanOutput, -1)

	var finalOutput string
	if len(parts) > 2 {
		finalOutput = strings.TrimSpace(parts[1])
	} else {
		finalOutput = strings.TrimSpace(cleanOutput)
	}

	return &CommandResponse{
		Command: command,
		HostID:  strconv.Itoa(int(hostID)),
		Output:  finalOutput,
	}, nil
}

func (s *CmdbHostSSHServiceImpl) CreateUploadTask(c *gin.Context, hostID uint, filePath string, originalFileName string, destPath string) (*cmdbModel.HostSSHUploadTask, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取临时文件信息失败: %v", err)
	}

	var createdBy uint
	var createdByName string
	if admin, err := jwt.GetAdmin(c); err == nil && admin != nil {
		createdBy = admin.ID
		createdByName = admin.Username
	}

	now := time.Now()
	task := &cmdbModel.HostSSHUploadTask{
		TaskID:        nextHostSSHUploadTaskID(),
		HostID:        hostID,
		FileName:      originalFileName,
		FileSize:      fileInfo.Size(),
		DestPath:      destPath,
		Status:        cmdbModel.HostSSHUploadTaskPending,
		Progress:      0,
		UploadedBytes: 0,
		TotalBytes:    fileInfo.Size(),
		Message:       "任务已创建，等待上传",
		CreatedBy:     createdBy,
		CreatedByName: createdByName,
		CreatedAt:     now,
	}

	hostSSHUploadTasks.Store(task.TaskID, task)

	go func() {
		startedAt := time.Now()
		task.Status = cmdbModel.HostSSHUploadTaskUploading
		task.StartedAt = &startedAt
		task.Message = "正在上传到远程主机"

		err := s.uploadFileWithProgress(hostID, filePath, destPath, func(uploaded, total int64) {
			task.UploadedBytes = uploaded
			task.TotalBytes = total
			if total > 0 {
				task.Progress = float64(uploaded) / float64(total) * 100
			}
			task.Message = fmt.Sprintf("正在上传到远程主机: %.2f%%", task.Progress)
		})

		completedAt := time.Now()
		task.CompletedAt = &completedAt

		if err != nil {
			task.Status = cmdbModel.HostSSHUploadTaskFailed
			task.Error = err.Error()
			task.Message = "上传失败"
			log.Printf("SSH上传后台任务失败: taskID=%s err=%v", task.TaskID, err)
		} else {
			task.Status = cmdbModel.HostSSHUploadTaskCompleted
			task.Progress = 100
			task.UploadedBytes = task.TotalBytes
			task.Message = "上传完成"
		}

		if removeErr := os.Remove(filePath); removeErr != nil && !os.IsNotExist(removeErr) {
			log.Printf("删除上传临时文件失败: taskID=%s path=%s err=%v", task.TaskID, filePath, removeErr)
		}
	}()

	return cloneHostSSHUploadTask(task), nil
}

func (s *CmdbHostSSHServiceImpl) GetUploadTask(taskID string) (*cmdbModel.HostSSHUploadTask, error) {
	taskValue, ok := hostSSHUploadTasks.Load(taskID)
	if !ok {
		return nil, fmt.Errorf("上传任务不存在: %s", taskID)
	}

	task, ok := taskValue.(*cmdbModel.HostSSHUploadTask)
	if !ok {
		return nil, fmt.Errorf("上传任务数据异常: %s", taskID)
	}

	return cloneHostSSHUploadTask(task), nil
}

func (s *CmdbHostSSHServiceImpl) UploadFile(c *gin.Context, hostID uint, filePath string, destPath string) error {
	return s.uploadFileWithProgress(hostID, filePath, destPath, nil)
}

func (s *CmdbHostSSHServiceImpl) uploadFileWithProgress(hostID uint, filePath string, destPath string, progress func(uploaded, total int64)) error {
	// 验证文件路径是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("文件不存在: %s", filePath)
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		log.Printf("获取主机信息失败: %v", err)
		return fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		return errors.New("主机SSH信息不完整")
	}

	// 获取SSH认证信息
	var auth configModel.EcsAuth
	if host.SSHKeyID != 0 {
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			log.Printf("获取认证类型失败: %v", err)
			return fmt.Errorf("获取认证类型失败: %v", err)
		}
	} else {
		return errors.New("主机未配置SSH凭据")
	}

	// 使用绝对路径
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Printf("获取绝对路径失败: %v", err)
		return fmt.Errorf("获取文件绝对路径失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("文件不存在(绝对路径): %s", absFilePath)
		return fmt.Errorf("文件不存在: %s", absFilePath)
	}

	// 生成唯一临时目录路径
	tmpDir := fmt.Sprintf("/tmp/ssh_uploads_%d_%d", time.Now().UnixNano(), os.Getpid())
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Printf("创建临时目录失败: %v", err)
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir) // 完成后自动清理

	// 实现文件复制
	tmpFilePath := filepath.Join(tmpDir, filepath.Base(filePath))
	srcFile, err := os.Open(absFilePath)
	if err != nil {
		log.Printf("打开源文件失败: %v", err)
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(tmpFilePath)
	if err != nil {
		log.Printf("创建目标文件失败: %v", err)
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		log.Printf("复制文件内容失败: %v", err)
		return fmt.Errorf("复制文件内容失败: %v", err)
	}
	absFilePath = tmpFilePath // 使用临时文件路径

	// 构建SCP命令并添加详细调试
	log.Printf("文件上传调试信息:")
	log.Printf(" - 本地文件路径: %s", absFilePath)
	log.Printf(" - 文件存在: %t", func() bool {
		_, err := os.Stat(absFilePath)
		return err == nil
	}())
	log.Printf(" - 文件权限: %v", func() string {
		fi, err := os.Stat(absFilePath)
		if err != nil {
			return err.Error()
		}
		return fi.Mode().String()
	}())
	log.Printf(" - 当前工作目录: %s", func() string {
		wd, _ := os.Getwd()
		return wd
	}())
	log.Printf(" - 文件所在目录: %s", filepath.Dir(absFilePath))
	log.Printf(" - 文件名: %s", filepath.Base(absFilePath))

	var client *ssh.Client

	// 根据认证类型创建SSH客户端
	switch auth.Type {
	case 4: // 跳板机认证
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

		// 通过跳板机连接
		client, err = util.NewSSHUtil().CreateBastionSSHClient(bastionAuth, targetAuth)
		if err != nil {
			log.Printf("通过跳板机连接失败: %v", err)
			return fmt.Errorf("通过跳板机连接失败: %v", err)
		}

	default: // 其他认证方式 (1密码, 2私钥, 3免密)
		var authMethod ssh.AuthMethod
		switch auth.Type {
		case 1: // 密码认证
			authMethod = ssh.Password(auth.Password)
		case 2: // 密钥认证
			authMethod, err = util.NewSSHUtil().PublicKeyAuth(auth.PublicKey)
			if err != nil {
				log.Printf("创建密钥认证失败: %v", err)
				return fmt.Errorf("创建密钥认证失败: %v", err)
			}
		case 3: // 公钥免认证 - 自动查找本地私钥
			if userKeyAuth, err := util.NewSSHUtil().UserKeyAuth(); err == nil {
				authMethod = userKeyAuth
			} else if defaultKeyAuth, err := util.NewSSHUtil().DefaultKeyAuth(); err == nil {
				authMethod = defaultKeyAuth
			} else {
				log.Printf("type=3 公钥免认证失败: 本地未找到私钥文件")
				return errors.New("type=3 公钥免认证失败: 本地未找到私钥文件")
			}
		default:
			return errors.New("不支持的认证方式")
		}

		config := &ssh.ClientConfig{
			User: host.SSHName,
			Auth: []ssh.AuthMethod{
				authMethod,
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second, // 增加SSH连接超时时间
		}

		addr := fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort)

		// 先测试基本SSH连接
		testClient, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			log.Printf("SSH连接测试失败: %v", err)
			return fmt.Errorf("无法连接到SSH服务器(%s): %v (请检查网络连接和服务器状态)", addr, err)
		}
		defer testClient.Close()

		// 测试执行简单命令
		testSession, err := testClient.NewSession()
		if err != nil {
			log.Printf("创建测试会话失败: %v", err)
			return fmt.Errorf("SSH会话创建失败: %v", err)
		}
		defer testSession.Close()

		if _, err := testSession.CombinedOutput("echo connection-test"); err != nil {
			log.Printf("SSH命令测试失败: %v", err)
			return fmt.Errorf("SSH命令执行失败: %v (请检查远程服务器配置)", err)
		}

		// 创建正式连接
		client, err = ssh.Dial("tcp", addr, config)
		if err != nil {
			log.Printf("创建SSH连接失败: %v", err)
			return fmt.Errorf("创建SSH连接失败: %v", err)
		}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Printf("创建SSH会话失败: %v", err)
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 获取当前工作目录
	wd, _ := os.Getwd()
	log.Printf("当前工作目录: %s", wd)

	// 直接上传文件，依赖前端并发控制
	// 设置60秒超时(与前端一致)
	log.Printf("开始文件上传流程(超时60秒)...")

	// 使用Go的SCP实现上传文件
	file, err := os.Open(absFilePath)
	if err != nil {
		log.Printf("打开本地文件失败: %v", err)
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 创建远程目录 (使用单独的session)
	log.Printf("开始创建远程目录...")
	mkdirSession, err := client.NewSession()
	if err != nil {
		log.Printf("创建mkdir会话失败: %v", err)
		return fmt.Errorf("创建mkdir会话失败: %v", err)
	}
	defer mkdirSession.Close()

	mkdirCmd := fmt.Sprintf("mkdir -p %q", filepath.Dir(destPath))
	if output, err := mkdirSession.CombinedOutput(mkdirCmd); err != nil {
		log.Printf("创建远程目录失败: %v", err)
		log.Printf("命令输出: %s", string(output))
		return fmt.Errorf("创建远程目录失败: %v", err)
	}
	log.Printf("远程目录创建成功")

	// 尝试使用SCP协议上传文件（带进度反馈）
	scpSession, err := client.NewSession()
	if err != nil {
		log.Printf("创建SCP会话失败: %v", err)
		return fmt.Errorf("创建SCP会话失败: %v", err)
	}
	defer scpSession.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
		return fmt.Errorf("获取文件信息失败: %v", err)
	}
	fileSize := fileInfo.Size()

	log.Printf("开始上传文件: %s -> %s@%s:%s, 大小: %.2f MB",
		absFilePath, host.SSHName, host.SSHIP, destPath, float64(fileSize)/(1024*1024))

	// 准备SCP命令并实现带进度的上传
	uploadErr := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				uploadErr <- fmt.Errorf("上传goroutine panic: %v", r)
			}
		}()

		w, err := scpSession.StdinPipe()
		if err != nil {
			uploadErr <- fmt.Errorf("获取stdin管道失败: %v", err)
			return
		}
		defer w.Close()

		file, err := os.Open(absFilePath)
		if err != nil {
			uploadErr <- fmt.Errorf("打开本地文件失败: %v", err)
			return
		}
		defer file.Close()

		// 发送SCP头信息
		fmt.Fprintf(w, "C%04o %d %s\n", fileInfo.Mode().Perm(), fileSize, filepath.Base(destPath))

		// 流式上传并显示进度
		buffer := make([]byte, 1024*1024) // 1MB缓冲区
		var uploaded int64
		lastLogTime := time.Now()
		logInterval := 2 * time.Second // 每2秒记录一次进度

		for {
			n, err := file.Read(buffer)
			if n > 0 {
				_, writeErr := w.Write(buffer[:n])
				if writeErr != nil {
					uploadErr <- fmt.Errorf("写入远程文件失败: %v", writeErr)
					return
				}
				uploaded += int64(n)

				// 定期记录上传进度
				if time.Since(lastLogTime) >= logInterval {
					progress := float64(uploaded) / float64(fileSize) * 100
					uploadedMB := float64(uploaded) / (1024 * 1024)
					totalMB := float64(fileSize) / (1024 * 1024)
					log.Printf("上传进度: %.2f%% (%.2f MB / %.2f MB)",
						progress, uploadedMB, totalMB)
					lastLogTime = time.Now()
				}
				if progress != nil {
					progress(uploaded, fileSize)
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				uploadErr <- fmt.Errorf("读取本地文件失败: %v", err)
				return
			}
		}

		// 发送结束标记
		fmt.Fprint(w, "\x00")
		uploadErr <- nil
	}()

	// 确保目标路径包含文件名
	finalDestPath := filepath.Join(destPath, filepath.Base(absFilePath))

	// 执行SCP命令
	scpErr := scpSession.Run(fmt.Sprintf("scp -t %s", finalDestPath))

	// 等待上传goroutine完成
	uploadGoroutineErr := <-uploadErr

	if scpErr != nil {
		log.Printf("SCP命令执行失败: %v", scpErr)
		return fmt.Errorf("SCP上传失败: %v", scpErr)
	}

	if uploadGoroutineErr != nil {
		log.Printf("上传数据失败: %v", uploadGoroutineErr)
		return fmt.Errorf("上传数据失败: %v", uploadGoroutineErr)
	}

	log.Printf("文件上传成功: %s -> %s@%s:%s (%.2f MB)",
		absFilePath, host.SSHName, host.SSHIP, destPath, float64(fileSize)/(1024*1024))
	if progress != nil {
		progress(fileSize, fileSize)
	}

	// 验证远程文件
	checkCmd := fmt.Sprintf("ls -la %q", destPath)
	if output, err := session.CombinedOutput(checkCmd); err != nil {
		log.Printf("远程文件验证失败: %v", err)
	} else {
		log.Printf("远程文件验证成功: %s", string(output))
	}

	return nil
}

func (s *CmdbHostSSHServiceImpl) DownloadFile(c *gin.Context, hostID uint, remotePath string, localPath string) error {
	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		log.Printf("获取主机信息失败: %v", err)
		return fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		return errors.New("主机SSH信息不完整")
	}

	// 获取SSH认证信息
	var auth configModel.EcsAuth
	if host.SSHKeyID != 0 {
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			log.Printf("获取认证类型失败: %v", err)
			return fmt.Errorf("获取认证类型失败: %v", err)
		}
	} else {
		return errors.New("主机未配置SSH凭据")
	}

	var client *ssh.Client

	// 根据认证类型创建SSH客户端
	if auth.Type == 4 { // 跳板机认证
		bastionAuth := &util.SSHConfig{
			IP:        auth.BastionHost,
			Port:      auth.BastionPort,
			Username:  auth.BastionUsername,
			Type:      auth.BastionAuthType,
			Password:  auth.BastionPassword,
			PublicKey: auth.BastionPrivateKey,
		}

		targetAuth := &util.SSHConfig{
			IP:        host.SSHIP,
			Port:      host.SSHPort,
			Username:  host.SSHName,
			Type:      auth.TargetAuthType,
			Password:  auth.TargetPassword,
			PublicKey: auth.TargetPrivateKey,
		}

		client, err = util.NewSSHUtil().CreateBastionSSHClient(bastionAuth, targetAuth)
		if err != nil {
			log.Printf("通过跳板机连接失败: %v", err)
			return fmt.Errorf("通过跳板机连接失败: %v", err)
		}
	} else { // 其他认证方式
		sshConfig := &util.SSHConfig{
			IP:        host.SSHIP,
			Port:      host.SSHPort,
			Type:      auth.Type,
			Username:  host.SSHName,
			Password:  auth.Password,
			PublicKey: auth.PublicKey,
		}

		client, err = util.GetSSHClientByConfig(sshConfig)
		if err != nil {
			log.Printf("SSH连接失败: %v", err)
			return fmt.Errorf("SSH连接失败: %v", err)
		}
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Printf("创建SFTP客户端失败: %v", err)
		return fmt.Errorf("创建SFTP客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 打开远程文件
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		log.Printf("打开远程文件失败: %v", err)
		return fmt.Errorf("打开远程文件失败: %v", err)
	}
	defer remoteFile.Close()

	// 获取文件大小
	fileInfo, err := remoteFile.Stat()
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
		return fmt.Errorf("获取文件信息失败: %v", err)
	}
	fileSize := fileInfo.Size()

	// 设置响应头 - 让前端知道文件大小
	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(remotePath)))
	c.Header("X-Accel-Buffering", "no") // 禁用Nginx缓冲

	// 直接流式传输到HTTP响应 - 关键改动！
	buffer := make([]byte, 32*1024) // 32KB 缓冲区（更小的缓冲区让前端更快感知进度）
	var transferred int64
	lastLogTime := time.Now()
	logInterval := 2 * time.Second

	log.Printf("开始流式下载文件: %s@%s:%s, 大小: %.2f MB",
		host.SSHName, host.SSHIP, remotePath, float64(fileSize)/(1024*1024))

	// 直接写入HTTP响应流
	writer := c.Writer
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return errors.New("响应流不支持Flush操作")
	}

	for {
		n, err := remoteFile.Read(buffer)
		if n > 0 {
			// 直接写入HTTP响应
			_, writeErr := writer.Write(buffer[:n])
			if writeErr != nil {
				log.Printf("写入HTTP响应失败: %v", writeErr)
				return fmt.Errorf("写入HTTP响应失败: %v", writeErr)
			}

			// 立即刷新缓冲区，让前端能接收到数据
			flusher.Flush()

			transferred += int64(n)

			// 定期记录进度
			if time.Since(lastLogTime) >= logInterval {
				progress := float64(transferred) / float64(fileSize) * 100
				transferredMB := float64(transferred) / (1024 * 1024)
				totalMB := float64(fileSize) / (1024 * 1024)
				log.Printf("流式传输进度: %.2f%% (%.2f MB / %.2f MB)",
					progress, transferredMB, totalMB)
				lastLogTime = time.Now()
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("读取远程文件失败: %v", err)
			return fmt.Errorf("读取远程文件失败: %v", err)
		}
	}

	log.Printf("文件流式传输完成: %s@%s:%s (%.2f MB)",
		host.SSHName, host.SSHIP, remotePath, float64(fileSize)/(1024*1024))
	return nil
}

func (s *CmdbHostSSHServiceImpl) ListFiles(c *gin.Context, hostID uint, remotePath string) ([]util.FileInfo, error) {
	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		log.Printf("获取主机信息失败: %v", err)
		return nil, fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		return nil, errors.New("主机SSH信息不完整")
	}

	// 获取SSH认证信息
	var auth configModel.EcsAuth
	if host.SSHKeyID != 0 {
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			log.Printf("获取认证类型失败: %v", err)
			return nil, fmt.Errorf("获取认证类型失败: %v", err)
		}
	} else {
		return nil, errors.New("主机未配置SSH凭据")
	}

	var client *ssh.Client

	// 根据认证类型创建SSH客户端
	if auth.Type == 4 { // 跳板机认证
		bastionAuth := &util.SSHConfig{
			IP:        auth.BastionHost,
			Port:      auth.BastionPort,
			Username:  auth.BastionUsername,
			Type:      auth.BastionAuthType,
			Password:  auth.BastionPassword,
			PublicKey: auth.BastionPrivateKey,
		}

		targetAuth := &util.SSHConfig{
			IP:        host.SSHIP,
			Port:      host.SSHPort,
			Username:  host.SSHName,
			Type:      auth.TargetAuthType,
			Password:  auth.TargetPassword,
			PublicKey: auth.TargetPrivateKey,
		}

		client, err = util.NewSSHUtil().CreateBastionSSHClient(bastionAuth, targetAuth)
		if err != nil {
			log.Printf("通过跳板机连接失败: %v", err)
			return nil, fmt.Errorf("通过跳板机连接失败: %v", err)
		}
	} else { // 其他认证方式
		sshConfig := &util.SSHConfig{
			IP:        host.SSHIP,
			Port:      host.SSHPort,
			Type:      auth.Type,
			Username:  host.SSHName,
			Password:  auth.Password,
			PublicKey: auth.PublicKey,
		}

		client, err = util.GetSSHClientByConfig(sshConfig)
		if err != nil {
			log.Printf("SSH连接失败: %v", err)
			return nil, fmt.Errorf("SSH连接失败: %v", err)
		}
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Printf("创建SFTP客户端失败: %v", err)
		return nil, fmt.Errorf("创建SFTP客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 读取目录
	files, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		log.Printf("读取目录失败: %v", err)
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	// 转换为 FileInfo 结构
	var fileInfos []util.FileInfo
	for _, file := range files {
		fullPath := remotePath
		if !strings.HasSuffix(fullPath, "/") {
			fullPath += "/"
		}
		fullPath += file.Name()

		fileInfos = append(fileInfos, util.FileInfo{
			Name:    file.Name(),
			Path:    fullPath,
			Size:    file.Size(),
			Mode:    file.Mode().String(),
			ModTime: file.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:   file.IsDir(),
		})
	}

	log.Printf("列出文件成功: %s@%s:%s (%d个文件)", host.SSHName, host.SSHIP, remotePath, len(fileInfos))
	return fileInfos, nil
}

func GetCmdbHostSSHService() CmdbHostSSHServiceInterface {
	return &CmdbHostSSHServiceImpl{}
}
