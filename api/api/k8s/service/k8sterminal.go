package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync"

	"dodevops-api/api/k8s/dao"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// IK8sTerminalService 容器终端服务接口
type IK8sTerminalService interface {
	CreateK8sWebSocketStream(clusterId uint, namespaceName, podName, containerName, command string, conn *websocket.Conn) (*K8sWebSocketStream, error)
	GetPodContainers(clusterId uint, namespaceName, podName string) ([]string, error)
	GetPodFileList(clusterId uint, namespaceName, podName, containerName, path string) ([]map[string]interface{}, error)
	GetPodFileContent(clusterId uint, namespaceName, podName, containerName, path string) (string, error)
	UpdatePodFileContent(clusterId uint, namespaceName, podName, containerName, path, content string) error
	UploadPodFile(clusterId uint, namespaceName, podName, containerName, path string, content []byte) error
	DownloadPodFile(clusterId uint, namespaceName, podName, containerName, path string) ([]byte, error)
	CreatePodDirectory(clusterId uint, namespaceName, podName, containerName, path string) error
	DeletePodFile(clusterId uint, namespaceName, podName, containerName, path string) error
	HotReloadPod(clusterId uint, namespaceName, podName, containerName string) error
}

// K8sMessage WebSocket消息结构
type K8sMessage struct {
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
	Cols      int         `json:"cols,omitempty"`
	Rows      int         `json:"rows,omitempty"`
}

// K8sWebSocketStream K8s WebSocket流处理
type K8sWebSocketStream struct {
	sync.RWMutex
	Conn     *websocket.Conn
	executor remotecommand.Executor
	Ctx      context.Context
	cancel   context.CancelFunc
	closed   bool
	reader   *io.PipeReader
	writer   *io.PipeWriter
}

// terminalConn 实现io.ReadWriter接口用于K8s executor
type terminalConn struct {
	stream *K8sWebSocketStream
}

func (tc *terminalConn) Read(p []byte) (n int, err error) {
	return tc.stream.reader.Read(p)
}

func (tc *terminalConn) Write(p []byte) (n int, err error) {
	return tc.stream.WriteToWebSocket(p)
}

// Close 关闭WebSocket流
func (kws *K8sWebSocketStream) Close() error {
	kws.Lock()
	defer kws.Unlock()

	if kws.closed {
		return nil
	}

	kws.closed = true
	if kws.cancel != nil {
		kws.cancel()
	}
	// 不要立即关闭管道，让K8s executor自然结束
	return nil
}

// IsClosed 检查是否已关闭
func (kws *K8sWebSocketStream) IsClosed() bool {
	kws.RLock()
	defer kws.RUnlock()
	return kws.closed
}

// WriteToWebSocket 写入数据到WebSocket
func (kws *K8sWebSocketStream) WriteToWebSocket(p []byte) (n int, err error) {
	if kws.IsClosed() {
		return 0, io.EOF
	}

	message := K8sMessage{
		Operation: "stdout",
		Data:      string(p),
	}

	if err = kws.Conn.WriteJSON(message); err != nil {
		go kws.Close() // 异步关闭，避免死锁
		return 0, err
	}

	return len(p), nil
}

// ReadFromWebSocket 从WebSocket读取数据
func (kws *K8sWebSocketStream) ReadFromWebSocket() {
	defer func() {
		kws.Close()
		// 关闭写入管道，通知K8s executor结束
		if kws.writer != nil {
			kws.writer.Close()
		}
	}()

	for {
		if kws.IsClosed() {
			return
		}

		var message K8sMessage
		err := kws.Conn.ReadJSON(&message)
		if err != nil {
			return
		}

		switch message.Operation {
		case "stdin":
			if kws.writer != nil && !kws.IsClosed() {
				var data string
				if str, ok := message.Data.(string); ok {
					data = str
				} else {
					continue
				}

				_, err := kws.writer.Write([]byte(data))
				if err != nil {
					return
				}
			}
		case "resize":
			// 处理终端大小调整，暂时忽略以避免问题
			continue
		default:
			// 忽略未知操作
			continue
		}
	}
}

// K8sTerminalService 容器终端服务实现
type K8sTerminalService struct {
	dao *dao.KubeClusterDao
}

func NewK8sTerminalService(db *gorm.DB) IK8sTerminalService {
	return &K8sTerminalService{
		dao: dao.NewKubeClusterDao(db),
	}
}

// CreateK8sWebSocketStream 创建K8s WebSocket流
func (s *K8sTerminalService) CreateK8sWebSocketStream(clusterId uint, namespaceName, podName, containerName, command string, conn *websocket.Conn) (*K8sWebSocketStream, error) {
	// 获取集群信息
	cluster, err := s.dao.GetByID(clusterId)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 创建Kubernetes客户端
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	// 获取Pod信息，如果没有指定容器名称，使用第一个容器
	if containerName == "" {
		pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("获取Pod信息失败: %v", err)
		}
		if len(pod.Spec.Containers) == 0 {
			return nil, fmt.Errorf("Pod中没有找到容器")
		}
		containerName = pod.Spec.Containers[0].Name
	}

	// 创建exec请求
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName).
		Param("command", command).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true")

	// 创建executor
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("创建executor失败: %v", err)
	}

	// 创建长期运行的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 创建管道
	reader, writer := io.Pipe()

	// 创建流对象
	stream := &K8sWebSocketStream{
		Conn:     conn,
		executor: executor,
		Ctx:      ctx,
		cancel:   cancel,
		reader:   reader,
		writer:   writer,
	}

	// 创建终端连接
	termConn := &terminalConn{stream: stream}

	// 启动WebSocket读取goroutine
	go stream.ReadFromWebSocket()

	// 启动K8s executor goroutine
	go func() {
		defer func() {
			cancel()
			// 关闭读取管道
			if reader != nil {
				reader.Close()
			}
		}()

		executor.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:  termConn,
			Stdout: termConn,
			Stderr: termConn,
			Tty:    true,
		})
	}()

	return stream, nil
}

// GetPodContainers 获取Pod中的容器列表
func (s *K8sTerminalService) GetPodContainers(clusterId uint, namespaceName, podName string) ([]string, error) {
	// 获取集群信息
	cluster, err := s.dao.GetByID(clusterId)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 创建Kubernetes客户端
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	// 获取Pod信息
	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取Pod信息失败: %v", err)
	}

	// 提取容器名称
	var containers []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	for _, container := range pod.Spec.InitContainers {
		containers = append(containers, container.Name+" (init)")
	}

	return containers, nil
}

// ExecCommand 在Pod容器中执行命令（内部辅助方法）
func (s *K8sTerminalService) ExecCommand(clusterId uint, namespaceName, podName, containerName string, command []string) (string, error) {
	cluster, err := s.dao.GetByID(clusterId)
	if err != nil {
		return "", err
	}
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	if containerName == "" {
		pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		if len(pod.Spec.Containers) > 0 {
			containerName = pod.Spec.Containers[0].Name
		}
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName).
		Param("stdout", "true").
		Param("stderr", "true")

	for _, cmd := range command {
		req.Param("command", cmd)
	}

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %v, stderr: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

// GetPodFileList 获取Pod文件列表
func (s *K8sTerminalService) GetPodFileList(clusterId uint, namespaceName, podName, containerName, path string) ([]map[string]interface{}, error) {
	if path == "" {
		path = "/"
	}

	// 使用 ls -al -p 命令，可以区分出目录并显示所有文件包括隐藏文件
	cmd := []string{"ls", "-al", "-p", path}
	output, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, cmd)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var files []map[string]interface{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total ") {
			continue
		}

		// 简单的解析 ls 输出，适用于大多 bash 环境
		parts := strings.Fields(line)
		if len(parts) >= 9 {
			name := strings.Join(parts[8:], " ")
			// 跳过当前和上级目录表示
			if name == "./" || name == "../" {
				continue
			}

			isDir := strings.HasSuffix(name, "/")
			if isDir {
				name = strings.TrimSuffix(name, "/")
			} else {
				isDir = strings.HasPrefix(parts[0], "d")
			}

			files = append(files, map[string]interface{}{
				"name":    name,
				"isDir":   isDir,
				"size":    parts[4],
				"modTime": parts[5] + " " + parts[6] + " " + parts[7],
				"mode":    parts[0],
				"owner":   parts[2],
				"group":   parts[3],
			})
		}
	}
	return files, nil
}

// GetPodFileContent 获取Pod容器中的文件内容
func (s *K8sTerminalService) GetPodFileContent(clusterId uint, namespaceName, podName, containerName, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	output, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, []string{"cat", path})
	if err != nil {
		return "", err
	}
	return output, nil
}

// UpdatePodFileContent 更新Pod容器中的文件内容
func (s *K8sTerminalService) UpdatePodFileContent(clusterId uint, namespaceName, podName, containerName, path, content string) error {
	if path == "" {
		return fmt.Errorf("路径不能为空")
	}
	// 利用base64编码避免特殊字符(引号、换行等)破坏shell执行
	encoded := base64.StdEncoding.EncodeToString([]byte(content))
	cmd := []string{"sh", "-c", fmt.Sprintf("echo '%s' | base64 -d > '%s'", encoded, path)}
	_, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, cmd)
	return err
}

// UploadPodFile 上传文件到Pod容器
func (s *K8sTerminalService) UploadPodFile(clusterId uint, namespaceName, podName, containerName, path string, content []byte) error {
	if path == "" {
		return fmt.Errorf("路径不能为空")
	}
	// 利用base64编码避免特殊字符破坏shell执行
	encoded := base64.StdEncoding.EncodeToString(content)
	cmd := []string{"sh", "-c", fmt.Sprintf("echo '%s' | base64 -d > '%s'", encoded, path)}
	_, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, cmd)
	return err
}

// DownloadPodFile 下载Pod容器中的文件
func (s *K8sTerminalService) DownloadPodFile(clusterId uint, namespaceName, podName, containerName, path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("路径不能为空")
	}
	// 利用base64编码输出，避免二进制数据解析或截断错误
	cmd := []string{"sh", "-c", fmt.Sprintf("base64 '%s'", path)}
	output, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, cmd)
	if err != nil {
		return nil, err
	}
	// 去除可能存在的换行符等影响解码的字符
	output = strings.ReplaceAll(output, "\n", "")
	output = strings.ReplaceAll(output, "\r", "")
	output = strings.TrimSpace(output)
	return base64.StdEncoding.DecodeString(output)
}

// CreatePodDirectory 创建Pod容器中的目录
func (s *K8sTerminalService) CreatePodDirectory(clusterId uint, namespaceName, podName, containerName, path string) error {
	if path == "" {
		return fmt.Errorf("路径不能为空")
	}
	_, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, []string{"mkdir", "-p", path})
	return err
}

// DeletePodFile 删除Pod容器中的文件
func (s *K8sTerminalService) DeletePodFile(clusterId uint, namespaceName, podName, containerName, path string) error {
	if path == "" || path == "/" {
		return fmt.Errorf("不能删除根目录或空路径")
	}
	_, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, []string{"rm", "-rf", path})
	return err
}

// HotReloadPod 热加载Pod（通常向 1 号进程发送 HUP 信号）
func (s *K8sTerminalService) HotReloadPod(clusterId uint, namespaceName, podName, containerName string) error {
	// 向前台主进程发送 SIGHUP，常见如 Nginx 会平滑重载配置
	cmd := []string{"sh", "-c", "kill -HUP 1 || kill -1 1"}
	_, err := s.ExecCommand(clusterId, namespaceName, podName, containerName, cmd)
	return err
}
