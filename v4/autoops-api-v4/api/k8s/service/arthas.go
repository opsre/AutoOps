package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// IK8sArthasService K8s Arthas诊断服务接口
type IK8sArthasService interface {
	ListJavaProcesses(clusterId uint, namespace, pod, container string) ([]model.ProcessInfo, error)
	CheckArthasInstalled(clusterId uint, namespace, pod, container string) (*model.ArthasCheckResult, error)
	InstallArthas(clusterId uint, namespace, pod, container string) error
	DownloadArthas() (bool, error) // bool=true 表示文件已存在，false 表示本次下载
	UploadArthas(fileName string, fileContent []byte) error
	ExecuteArthasCommand(clusterId uint, namespace, pod, container, processID, command string) (interface{}, error)
	GetDashboard(clusterId uint, namespace, pod, container, processID string) (interface{}, error)
	GetThreadList(clusterId uint, namespace, pod, container, processID string) (interface{}, error)
	GetThreadStack(clusterId uint, namespace, pod, container, processID, threadID string) (interface{}, error)
	GetJvmInfo(clusterId uint, namespace, pod, container, processID string) (interface{}, error)
	GetSysEnv(clusterId uint, namespace, pod, container, processID string) (interface{}, error)
	GetSysProp(clusterId uint, namespace, pod, container, processID string) (interface{}, error)
	GetMemory(clusterId uint, namespace, pod, container, processID string) (interface{}, error)
	GetPerfCounter(clusterId uint, namespace, pod, container, processID string) (interface{}, error)
	DecompileClass(clusterId uint, namespace, pod, container, processID, className string) (interface{}, error)
	GetStaticField(clusterId uint, namespace, pod, container, processID, className, fieldName string) (interface{}, error)
	SearchClass(clusterId uint, namespace, pod, container, processID, pattern string) (interface{}, error)
	SearchMethod(clusterId uint, namespace, pod, container, processID, className, methodName string) (interface{}, error)
	StartFlameGraph(clusterId uint, namespace, pod, container, processID string, duration int) (*model.FlameGraphTask, error)
	GetFlameGraphTask(taskId string) (*model.FlameGraphTask, error)
	ArthasWebSocket(clusterId uint, namespace, pod, container, processID string, conn *websocket.Conn)
}

// K8sArthasService K8s Arthas诊断服务实现
type K8sArthasService struct {
	dao            *dao.KubeClusterDao
	arthasLocalDir string // 本地 arthas-boot.jar 目录（如：/opt/arthas）
	arthasURL      string // 远程 arthas-boot.jar 下载 URL（如：http://内部镜像/arthas-boot.jar）
	tasks          sync.Map // 异步任务存储 taskId -> *model.FlameGraphTask
}

// NewK8sArthasService 创建K8s Arthas诊断服务
func NewK8sArthasService(db *gorm.DB) IK8sArthasService {
	svc := &K8sArthasService{
		dao: dao.NewKubeClusterDao(db),
	}

	// 从环境变量读取配置
	if dir := getEnv("ARTHAS_LOCAL_DIR", ""); dir != "" {
		svc.arthasLocalDir = dir
	} else {
		svc.arthasLocalDir = "upload/arthas"
	}

	// 转换为绝对路径，避免因服务器启动目录不同导致路径解析错误
	if absDir, err := filepath.Abs(svc.arthasLocalDir); err == nil {
		svc.arthasLocalDir = absDir
	}

	if url := getEnv("ARTHAS_REMOTE_URL", ""); url != "" {
		svc.arthasURL = url
	} else {
		// 默认使用官方镜像
		svc.arthasURL = "https://arthas.aliyun.com/arthas-boot.jar"
	}

	// 初始化：创建目录（如果不存在）
	os.MkdirAll(svc.arthasLocalDir, 0755)

	return svc
}

// getEnv 获取环境变量，返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getK8sClientWithConfig 从数据库获取集群凭证，创建 clientset 和 rest.Config
func (s *K8sArthasService) getK8sClientWithConfig(clusterId uint) (*kubernetes.Clientset, *rest.Config, error) {
	cluster, err := s.dao.GetByID(clusterId)
	if err != nil {
		return nil, nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, nil, fmt.Errorf("创建Kubernetes配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	return clientset, config, nil
}

// execInPod 在 Pod 内执行命令，返回 stdout
// 使用 remotecommand.NewSPDYExecutor + StreamWithContext（无TTY，仅捕获stdout/stderr）
func (s *K8sArthasService) execInPod(clusterId uint, namespace, pod, container string, commands []string) (string, error) {
	clientset, config, err := s.getK8sClientWithConfig(clusterId)
	if err != nil {
		return "", err
	}

	// 创建exec请求
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container)

	// 添加每个命令参数
	for _, cmd := range commands {
		req = req.Param("command", cmd)
	}

	// 添加标准输出/错误，无TTY
	req = req.
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	// 创建executor
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("创建executor失败: %v", err)
	}

	// 创建缓冲区捕获输出
	var stdout, stderr bytes.Buffer

	// 执行命令
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	// 无论成功与否，先获取输出内容
	output := stdout.String()
	if output == "" {
		output = stderr.String()
	}

	if err != nil {
		// 将 Pod 内的实际输出附加到错误信息，方便定位原因
		if output != "" {
			return "", fmt.Errorf("执行命令失败: %v\n详情: %s", err, output)
		}
		return "", fmt.Errorf("执行命令失败: %v", err)
	}

	return output, nil
}

// buildArthasScript 构建在 Pod 内执行 Arthas 命令的完整 shell 脚本
func (s *K8sArthasService) buildArthasScript(pid, command string) string {
	return fmt.Sprintf(`
TARGET_PID=%s
COMMAND='%s'
echo "[INFO] Executing Arthas command on process $TARGET_PID: $COMMAND"

# 查找可用的 Java 二进制（优先 JDK，JRE 也可尝试）
find_java() {
    # 1. 优先使用 JAVA_HOME
    if [ -n "$JAVA_HOME" ] && [ -x "$JAVA_HOME/bin/java" ]; then
        echo "$JAVA_HOME/bin/java"
        return 0
    fi

    # 2. 通过 which java 找到 java 并解析软链（适配容器环境）
    WHICH_JAVA=$(which java 2>/dev/null)
    if [ -n "$WHICH_JAVA" ]; then
        # 解析软链得到真实路径
        REAL_JAVA=$(readlink -f "$WHICH_JAVA" 2>/dev/null || realpath "$WHICH_JAVA" 2>/dev/null || echo "$WHICH_JAVA")
        echo "$REAL_JAVA"
        return 0
    fi

    # 3. 常见容器路径搜索（包括 openjdk/temurin/dragonwell 等）
    for java_bin in \
        /usr/local/openjdk-*/bin/java \
        /usr/lib/jvm/*/bin/java \
        /opt/java/*/bin/java \
        /opt/jdk/*/bin/java \
        /usr/java/*/bin/java; do
        if [ -x "$java_bin" ]; then
            echo "$java_bin"
            return 0
        fi
    done

    return 1
}

JAVA_BIN=$(find_java)
if [ -z "$JAVA_BIN" ]; then
    echo "[ERROR] 未找到 Java 环境，请确认容器内安装了 JDK"
    exit 1
fi
echo "[INFO] Using Java: $JAVA_BIN ($($JAVA_BIN -version 2>&1 | head -1))"

# 端口基于进程ID计算，避免冲突
BASE_PORT=$((3658 + (TARGET_PID %% 100)))
PORT=$BASE_PORT
MAX_ATTEMPTS=3
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    ATTEMPT=$((ATTEMPT + 1))
    echo "[INFO] Attempt $ATTEMPT/$MAX_ATTEMPTS on port $PORT"

    # 清空 JAVA_TOOL_OPTIONS，避免 Arthas 启动器被注入 javaagent（如 SkyWalking）
    OUTPUT=$(unset JAVA_TOOL_OPTIONS; $JAVA_BIN -jar /tmp/arthas-boot.jar $TARGET_PID --telnet-port $PORT --http-port -1 -c "$COMMAND" 2>&1)

    # 判断是否获取到有效数据
    if echo "$OUTPUT" | grep -qE "ID.*NAME|heap|nonheap|Memory|THREAD|RUNNABLE|WAITING|BLOCKED|TIMED_|gc\.|eden|os\.|java\.|sun\."; then
        echo "$OUTPUT"
        exit 0
    fi

    # 端口已被占用，换端口重试
    if echo "$OUTPUT" | grep -q "telnet port.*is used"; then
        echo "[WARN] Port $PORT in use, trying next..."
        PORT=$((PORT + 1))
        sleep 1
        continue
    fi

    # 进程不存在
    if echo "$OUTPUT" | grep -qE "No process|Can not find"; then
        echo "[ERROR] Process $TARGET_PID not found"
        echo "$OUTPUT"
        exit 1
    fi

    # 非 JDK 环境（缺少 attach API）
    if echo "$OUTPUT" | grep -qE "NoClassDefFoundError.*attach|Unable to attach|attach not supported|tools.jar"; then
        echo "[ERROR] ❌ Java 环境缺少 JDK attach API，Arthas 无法运行"
        echo "[ERROR] 当前 Java: $JAVA_BIN"
        echo "[ERROR] 容器使用的是 JRE（运行时），需要完整 JDK（开发工具包）"
        echo "[HINT] 解决方案："
        echo "[HINT]   1. 修改 Dockerfile 基础镜像为 JDK 版本："
        echo "[HINT]      - openjdk:11-jre  → openjdk:11-jdk"
        echo "[HINT]      - eclipse-temurin:11-jre → eclipse-temurin:11-jdk"
        echo "[HINT]   2. 或在 Dockerfile 中安装 JDK: apt-get install openjdk-11-jdk"
        echo "$OUTPUT"
        exit 1
    fi

    echo "$OUTPUT"

    if [ $ATTEMPT -lt $MAX_ATTEMPTS ]; then
        sleep 5
    fi
done

# 最后一次尝试（换端口）
FINAL_PORT=$((BASE_PORT + 100))
echo "[INFO] Final attempt with port $FINAL_PORT"
OUTPUT=$(unset JAVA_TOOL_OPTIONS; $JAVA_BIN -jar /tmp/arthas-boot.jar $TARGET_PID --telnet-port $FINAL_PORT --http-port -1 -c "$COMMAND" 2>&1)
echo "$OUTPUT"
if echo "$OUTPUT" | grep -qE "ID.*NAME|heap|nonheap|Memory|THREAD|RUNNABLE|WAITING|BLOCKED|TIMED_|gc\.|eden|os\.|java\.|sun\."; then
    exit 0
fi

echo "[ERROR] Arthas command failed after all attempts"
exit 1
`, pid, command)
}

// runArthasCommand 执行 Arthas 一次性命令
func (s *K8sArthasService) runArthasCommand(clusterId uint, namespace, pod, container, pid, command string) (interface{}, error) {
	// 每次直接从本地复制 arthas-boot.jar 到 Pod（保证版本一致）
	// 若本地文件不存在，自动尝试下载（不要求用户手动调用下载接口）
	localJarPath := s.arthasLocalDir + "/arthas-boot.jar"
	if _, err := os.Stat(localJarPath); err != nil {
		os.MkdirAll(s.arthasLocalDir, 0755)
		if dlErr := downloadFile(s.arthasURL, localJarPath); dlErr != nil {
			return nil, fmt.Errorf("本地 arthas-boot.jar 不存在（%s），自动下载也失败: %v", localJarPath, dlErr)
		}
	}
	if err := s.copyArthasToPod(clusterId, namespace, pod, container, localJarPath); err != nil {
		return nil, fmt.Errorf("复制 arthas-boot.jar 到 Pod 失败: %v", err)
	}

	// 执行诊断命令
	script := s.buildArthasScript(pid, command)
	output, err := s.execInPod(clusterId, namespace, pod, container, []string{"sh", "-c", script})
	if err != nil {
		return nil, fmt.Errorf("执行Arthas命令失败: %v", err)
	}

	var result interface{}
	if jsonErr := json.Unmarshal([]byte(output), &result); jsonErr != nil {
		return output, nil
	}
	return result, nil
}

// ensureArthasInPod 确保 Pod 内存在 arthas-boot.jar
// 优先级：1. upload/arthas/arthas-boot.jar 2. 从远程 URL 下载
func (s *K8sArthasService) ensureArthasInPod(clusterId uint, namespace, pod, container string) error {
	// 先检查 Pod 内是否已有 arthas-boot.jar
	checkCmd := "test -f /tmp/arthas-boot.jar && echo 'jar_found' || echo 'jar_missing'"
	output, _ := s.execInPod(clusterId, namespace, pod, container, []string{"sh", "-c", checkCmd})

	if strings.Contains(output, "jar_found") {
		return nil // 已经存在，无需重复复制
	}

	// 确保本地有 arthas-boot.jar 文件
	localJarPath := s.arthasLocalDir + "/arthas-boot.jar"
	if _, err := os.Stat(localJarPath); err != nil {
		// 本地文件不存在，从远程下载到 upload/arthas/
		if err := downloadFile(s.arthasURL, localJarPath); err != nil {
			return fmt.Errorf("下载 Arthas 失败: %v", err)
		}
	}

	// 通过 kubectl cp 复制到 Pod
	return s.copyArthasToPod(clusterId, namespace, pod, container, localJarPath)
}

// copyArthasToPod 通过 Kubernetes API 将本地 arthas-boot.jar 复制到 Pod
// 使用 cat > /tmp/arthas-boot.jar 方式，不依赖 tar 命令，兼容所有容器
func (s *K8sArthasService) copyArthasToPod(clusterId uint, namespace, pod, container, localJarPath string) error {
	fileContent, err := os.ReadFile(localJarPath)
	if err != nil {
		return fmt.Errorf("读取本地文件失败: %v", err)
	}

	clientset, config, err := s.getK8sClientWithConfig(clusterId)
	if err != nil {
		return err
	}

	// 使用 cat > /tmp/arthas-boot.jar，将文件内容通过 stdin 写入 Pod
	// 不依赖 tar，只需要容器有 sh + cat
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container).
		Param("command", "sh").
		Param("command", "-c").
		Param("command", "cat > /tmp/arthas-boot.jar").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("创建executor失败: %v", err)
	}

	var stdout, stderr bytes.Buffer
	if err := executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  bytes.NewReader(fileContent),
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		return fmt.Errorf("文件写入Pod失败: %v (stderr: %s)", err, stderr.String())
	}

	// 验证文件写入成功
	verifyOut, _ := s.execInPod(clusterId, namespace, pod, container, []string{"sh", "-c", "test -f /tmp/arthas-boot.jar && wc -c /tmp/arthas-boot.jar || echo fail"})
	if strings.Contains(verifyOut, "fail") || verifyOut == "" {
		return fmt.Errorf("复制后验证失败，/tmp/arthas-boot.jar 不存在")
	}

	return nil
}

// downloadFile 从 URL 下载文件到本地路径
func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// ListJavaProcesses 列表Java进程
// 使用 /proc 文件系统，不依赖 jps 或 ps 命令（兼容最小化容器）
func (s *K8sArthasService) ListJavaProcesses(clusterId uint, namespace, pod, container string) ([]model.ProcessInfo, error) {
	// 方案：使用 /proc 文件系统列表 Java 进程
	// 先找出所有 Java 进程的 PID，然后分别获取 cmdline
	getPidsCmd := `for pid in $(ls -d /proc/[0-9]* 2>/dev/null); do
  if grep -q 'java' "$pid/cmdline" 2>/dev/null; then
    basename "$pid"
  fi
done`

	pidsOutput, err := s.execInPod(clusterId, namespace, pod, container, []string{"sh", "-c", getPidsCmd})
	if err != nil || strings.TrimSpace(pidsOutput) == "" {
		return nil, fmt.Errorf("未找到 Java 进程")
	}

	var processes []model.ProcessInfo
	pids := strings.Fields(strings.TrimSpace(pidsOutput))

	for _, pid := range pids {
		// 获取每个 PID 的完整 cmdline
		getCmdlineCmd := fmt.Sprintf("cat /proc/%s/cmdline 2>/dev/null | tr '\\0' ' ' || echo ''", pid)
		cmdlineOutput, _ := s.execInPod(clusterId, namespace, pod, container, []string{"sh", "-c", getCmdlineCmd})

		cmdline := strings.TrimSpace(cmdlineOutput)
		if cmdline != "" {
			processInfo := model.ProcessInfo{
				ProcessID: pid,
				MainClass: cmdline,
			}
			processes = append(processes, processInfo)
		}
	}

	if len(processes) == 0 {
		return nil, fmt.Errorf("未找到 Java 进程")
	}

	return processes, nil
}

// CheckArthasInstalled 检查是否可以使用Arthas（检查Java和后端是否能获取Arthas）
func (s *K8sArthasService) CheckArthasInstalled(clusterId uint, namespace, pod, container string) (*model.ArthasCheckResult, error) {
	// 检查 Java 是否可用
	output, err := s.execInPod(clusterId, namespace, pod, container, []string{"sh", "-c", "which java"})
	if err != nil || strings.TrimSpace(output) == "" {
		return &model.ArthasCheckResult{
			Installed: false,
			Message:   "容器内未找到 Java，无法运行 Arthas",
		}, nil
	}

	javaPath := strings.TrimSpace(output)

	// 尝试确保 Arthas 可用（不需要 curl/wget，由后端负责）
	err = s.ensureArthasInPod(clusterId, namespace, pod, container)

	return &model.ArthasCheckResult{
		Installed: err == nil,
		JavaPath:  javaPath,
		Message: func() string {
			if err == nil {
				return "✓ 环境就绪，可以进行诊断"
			}
			return fmt.Sprintf("⚠ 无法获取 Arthas: %v", err)
		}(),
	}, nil
}

// DownloadArthas 下载 arthas-boot.jar 到本地 upload/arthas 目录
// 返回 (alreadyExists bool, err error)
// - alreadyExists=true 表示文件已存在，跳过下载
// - 支持直接下载 .jar 或下载 .zip 包并自动解压
// - 可通过 ARTHAS_REMOTE_URL 指定完整包 URL（如 arthas-packaging-bin.zip）
func (s *K8sArthasService) DownloadArthas() (bool, error) {
	// 检查文件是否已存在
	localJarPath := s.arthasLocalDir + "/arthas-boot.jar"
	if _, err := os.Stat(localJarPath); err == nil {
		return true, nil // 文件已存在，无需重复下载
	}

	// 确保目录存在
	os.MkdirAll(s.arthasLocalDir, 0755)

	// 根据 URL 后缀判断是 zip 包还是 jar
	urlLower := strings.ToLower(s.arthasURL)
	if strings.HasSuffix(urlLower, ".zip") {
		if err := s.downloadAndExtractArthas(); err != nil {
			return false, err
		}
		// 验证解压后 arthas-boot.jar 是否存在
		if _, err := os.Stat(localJarPath); err != nil {
			return false, fmt.Errorf("解压完成但未找到 arthas-boot.jar，请检查压缩包格式")
		}
	} else {
		// 直接下载 jar
		if err := downloadFile(s.arthasURL, localJarPath); err != nil {
			return false, fmt.Errorf("下载 Arthas 失败: %v", err)
		}

		fileInfo, err := os.Stat(localJarPath)
		if err != nil {
			return false, fmt.Errorf("验证文件失败: %v", err)
		}
		if fileInfo.Size() < 50*1024 {
			os.Remove(localJarPath)
			return false, fmt.Errorf("下载的文件过小 (%d bytes)，可能下载不完整", fileInfo.Size())
		}
	}

	return false, nil
}

// downloadAndExtractArthas 下载 Arthas 完整 zip 包并解压到本地目录
func (s *K8sArthasService) downloadAndExtractArthas() error {
	// 下载到临时文件
	tmpFile, err := os.CreateTemp("", "arthas-download-*.zip")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	resp, err := http.Get(s.arthasURL)
	if err != nil {
		tmpFile.Close()
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tmpFile.Close()
		return fmt.Errorf("下载失败，HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("写入临时文件失败: %v", err)
	}
	tmpFile.Close()

	// 解压 zip
	r, err := zip.OpenReader(tmpPath)
	if err != nil {
		return fmt.Errorf("解压失败: %v", err)
	}
	defer r.Close()

	baseDir := filepath.Clean(s.arthasLocalDir) + string(os.PathSeparator)
	for _, f := range r.File {
		// 去掉顶层目录前缀（如 arthas-bin/）
		name := f.Name
		if idx := strings.Index(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		if name == "" {
			continue
		}

		destPath := filepath.Join(s.arthasLocalDir, name)
		// 防止路径穿越
		if !strings.HasPrefix(destPath, baseDir) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(destPath), 0755)

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("读取压缩文件 %s 失败: %v", f.Name, err)
		}

		destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("创建文件 %s 失败: %v", destPath, err)
		}

		_, err = io.Copy(destFile, rc)
		rc.Close()
		destFile.Close()
		if err != nil {
			return fmt.Errorf("解压文件 %s 失败: %v", f.Name, err)
		}
	}

	return nil
}

// InstallArthas 将 arthas-boot.jar 复制到 Pod（如果不存在则先下载）
func (s *K8sArthasService) InstallArthas(clusterId uint, namespace, pod, container string) error {
	return s.ensureArthasInPod(clusterId, namespace, pod, container)
}

// UploadArthas 验证上传的 arthas-boot.jar 文件
// arthas-boot.jar 是 bootstrap 启动器，大小通常在 100KB-200KB，可自动下载完整诊断工具
func (s *K8sArthasService) UploadArthas(fileName string, fileContent []byte) error {
	// 验证文件名
	if fileName != "arthas-boot.jar" {
		return fmt.Errorf("文件名必须是 arthas-boot.jar")
	}

	// 验证文件大小（bootstrap jar 通常 100KB-200KB）
	if len(fileContent) < 50*1024 { // 小于 50KB
		return fmt.Errorf("文件过小，可能不是有效的 arthas-boot.jar")
	}

	if len(fileContent) > 500*1024*1024 { // 大于 500MB
		return fmt.Errorf("文件过大，请检查是否正确")
	}

	// 确保目录存在
	os.MkdirAll(s.arthasLocalDir, 0755)

	// 保存文件
	filePath := s.arthasLocalDir + "/arthas-boot.jar"
	if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// ExecuteArthasCommand 执行Arthas自定义命令
func (s *K8sArthasService) ExecuteArthasCommand(clusterId uint, namespace, pod, container, processID, command string) (interface{}, error) {
	return s.runArthasCommand(clusterId, namespace, pod, container, processID, command)
}

// GetDashboard 获取仪表板（-n 1 表示只采集一次快照，避免无限循环）
func (s *K8sArthasService) GetDashboard(clusterId uint, namespace, pod, container, processID string) (interface{}, error) {
	rawOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, "dashboard -n 1")
	if err != nil {
		return nil, err
	}

	// 如果是字符串，解析为结构化数据
	if strOutput, ok := rawOutput.(string); ok {
		parsed := parseDashboard(strOutput)
		// 如果解析结果为空，返回原始文本用于调试
		if parsed.Threads.Live == 0 && parsed.Memory.Heap.Used == 0 {
			return map[string]interface{}{
				"parsed": parsed,
				"raw":    strOutput[:min(500, len(strOutput))], // 返回前500字符用于调试
			}, nil
		}
		return parsed, nil
	}
	return rawOutput, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetThreadList 获取线程列表
func (s *K8sArthasService) GetThreadList(clusterId uint, namespace, pod, container, processID string) (interface{}, error) {
	rawOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, "thread")
	if err != nil {
		return nil, err
	}

	if strOutput, ok := rawOutput.(string); ok {
		return parseThreadList(strOutput), nil
	}
	return rawOutput, nil
}

// GetThreadStack 获取线程堆栈
func (s *K8sArthasService) GetThreadStack(clusterId uint, namespace, pod, container, processID, threadID string) (interface{}, error) {
	command := fmt.Sprintf("thread %s", threadID)
	return s.runArthasCommand(clusterId, namespace, pod, container, processID, command)
}

// GetJvmInfo 获取JVM信息
func (s *K8sArthasService) GetJvmInfo(clusterId uint, namespace, pod, container, processID string) (interface{}, error) {
	rawOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, "jvm")
	if err != nil {
		return nil, err
	}

	if strOutput, ok := rawOutput.(string); ok {
		return parseJvmInfo(strOutput), nil
	}
	return rawOutput, nil
}

// GetSysEnv 获取系统环境变量
func (s *K8sArthasService) GetSysEnv(clusterId uint, namespace, pod, container, processID string) (interface{}, error) {
	rawOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, "sysenv")
	if err != nil {
		return nil, err
	}

	if strOutput, ok := rawOutput.(string); ok {
		return parseSysprop(strOutput), nil // sysenv 和 sysprop 格式相同
	}
	return rawOutput, nil
}

// GetSysProp 获取系统属性
func (s *K8sArthasService) GetSysProp(clusterId uint, namespace, pod, container, processID string) (interface{}, error) {
	rawOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, "sysprop")
	if err != nil {
		return nil, err
	}

	if strOutput, ok := rawOutput.(string); ok {
		return parseSysprop(strOutput), nil
	}
	return rawOutput, nil
}

// GetMemory 获取内存信息
func (s *K8sArthasService) GetMemory(clusterId uint, namespace, pod, container, processID string) (interface{}, error) {
	rawOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, "memory")
	if err != nil {
		return nil, err
	}

	if strOutput, ok := rawOutput.(string); ok {
		return parseMemoryInfo(strOutput), nil
	}
	return rawOutput, nil
}

// GetPerfCounter 获取性能计数器
func (s *K8sArthasService) GetPerfCounter(clusterId uint, namespace, pod, container, processID string) (interface{}, error) {
	return s.runArthasCommand(clusterId, namespace, pod, container, processID, "perfcounter")
}

// DecompileClass 反编译类
func (s *K8sArthasService) DecompileClass(clusterId uint, namespace, pod, container, processID, className string) (interface{}, error) {
	command := fmt.Sprintf("jad %s", className)
	return s.runArthasCommand(clusterId, namespace, pod, container, processID, command)
}

// GetStaticField 获取静态字段
func (s *K8sArthasService) GetStaticField(clusterId uint, namespace, pod, container, processID, className, fieldName string) (interface{}, error) {
	command := fmt.Sprintf("getstatic %s %s", className, fieldName)
	return s.runArthasCommand(clusterId, namespace, pod, container, processID, command)
}

// SearchClass 搜索类
func (s *K8sArthasService) SearchClass(clusterId uint, namespace, pod, container, processID, pattern string) (interface{}, error) {
	command := fmt.Sprintf("sc %s", pattern)
	return s.runArthasCommand(clusterId, namespace, pod, container, processID, command)
}

// SearchMethod 搜索方法
func (s *K8sArthasService) SearchMethod(clusterId uint, namespace, pod, container, processID, className, methodName string) (interface{}, error) {
	command := fmt.Sprintf("sm %s %s", className, methodName)
	return s.runArthasCommand(clusterId, namespace, pod, container, processID, command)
}

// StartFlameGraph 异步启动火焰图采样，立即返回任务ID
func (s *K8sArthasService) StartFlameGraph(clusterId uint, namespace, pod, container, processID string, duration int) (*model.FlameGraphTask, error) {
	// 生成唯一任务ID
	taskId := fmt.Sprintf("%s_%s_%s_%d", namespace, pod, processID, time.Now().UnixMilli())

	task := &model.FlameGraphTask{
		TaskID:    taskId,
		Status:    model.FlameGraphTaskPending,
		Duration:  duration,
		StartedAt: time.Now().Unix(),
		Message:   "任务已提交，等待采样开始",
	}
	s.tasks.Store(taskId, task)

	// 后台 goroutine 执行采样
	go func() {
		// 更新状态为 running
		task.Status = model.FlameGraphTaskRunning
		task.Message = fmt.Sprintf("正在采样，持续 %d 秒...", duration)
		s.tasks.Store(taskId, task)

		result, err := s.doGenerateFlameGraph(clusterId, namespace, pod, container, processID, duration)

		task.FinishAt = time.Now().Unix()
		if err != nil {
			task.Status = model.FlameGraphTaskFailed
			task.Message = err.Error()
		} else {
			task.Status = model.FlameGraphTaskDone
			task.Message = result.Message
			task.Result = result
		}
		s.tasks.Store(taskId, task)
	}()

	return task, nil
}

// GetFlameGraphTask 查询火焰图任务状态
func (s *K8sArthasService) GetFlameGraphTask(taskId string) (*model.FlameGraphTask, error) {
	val, ok := s.tasks.Load(taskId)
	if !ok {
		return nil, fmt.Errorf("任务不存在: %s", taskId)
	}
	return val.(*model.FlameGraphTask), nil
}

// doGenerateFlameGraph 实际执行火焰图生成（内部调用）
func (s *K8sArthasService) doGenerateFlameGraph(clusterId uint, namespace, pod, container, processID string, duration int) (*model.FlameGraphResult, error) {
	// 0. 先尝试停止可能存在的旧 profiler 会话（避免冲突）
	s.runArthasCommand(clusterId, namespace, pod, container, processID, "profiler stop")

	// 1. 启动 profiler 采样
	startCmd := fmt.Sprintf("profiler start --duration %d", duration)
	startOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, startCmd)
	if err != nil {
		return nil, fmt.Errorf("启动 profiler 失败: %v", err)
	}

	startStr := fmt.Sprintf("%v", startOutput)
	if strings.Contains(startStr, "Error") || strings.Contains(startStr, "Failed") {
		return nil, fmt.Errorf("profiler 启动失败: %s", startStr)
	}

	// 2. 等待采样完成
	time.Sleep(time.Duration(duration+3) * time.Second)

	// 3. 停止并生成 collapsed 格式
	stopOutput, err := s.runArthasCommand(clusterId, namespace, pod, container, processID, "profiler stop --format collapsed")
	if err != nil {
		return nil, fmt.Errorf("生成火焰图失败: %v", err)
	}

	stopStr := fmt.Sprintf("%v", stopOutput)

	// 4. 提取文件路径
	filePath := extractProfilerOutputPath(stopStr)
	if filePath == "" {
		return nil, fmt.Errorf("未找到 profiler 输出文件路径")
	}

	// 5. 从 Pod 复制文件
	collapsedContent, err := s.copyFileFromPod(clusterId, namespace, pod, container, filePath)
	if err != nil {
		return nil, fmt.Errorf("复制 profiler 文件失败: %v", err)
	}

	// 6. 解析并优化数据结构
	flameGraphData, totalSamples := parseCollapsedFormat(string(collapsedContent))
	optimizedData := optimizeFlameGraphData(flameGraphData)

	return &model.FlameGraphResult{
		Data:    optimizedData,
		Message: fmt.Sprintf("火焰图生成成功，共采集 %d 个样本", totalSamples),
	}, nil
}

// optimizeFlameGraphData 优化火焰图数据结构
// 1. 合并只有单个子节点的链路（A->B->C 变成 "A;B;C"）
// 2. 限制最大深度
func optimizeFlameGraphData(data *model.FlameGraphData) *model.FlameGraphData {
	if data == nil {
		return nil
	}

	return optimizeNode(data, 0, 50) // 最大深度 50 层
}

// optimizeNode 递归优化节点
func optimizeNode(node *model.FlameGraphData, depth int, maxDepth int) *model.FlameGraphData {
	if node == nil || depth >= maxDepth {
		return node
	}

	// 如果只有一个子节点，合并到当前节点
	if len(node.Children) == 1 {
		fullPath := node.Name
		currentNode := node.Children[0]

		// 持续合并单链路
		for len(currentNode.Children) == 1 {
			fullPath += ";" + currentNode.Name
			currentNode = currentNode.Children[0]
		}
		fullPath += ";" + currentNode.Name

		// 取最后一段作为显示名
		parts := strings.Split(fullPath, ";")
		displayName := parts[len(parts)-1]

		// 创建合并后的节点
		optimized := &model.FlameGraphData{
			Name:     displayName,
			FullName: fullPath,
			Value:    node.Value,
			Children: make([]*model.FlameGraphData, 0),
		}

		// 递归处理子节点
		for _, grandChild := range currentNode.Children {
			optimized.Children = append(optimized.Children, optimizeNode(grandChild, depth+1, maxDepth))
		}

		return optimized
	}

	// 多个子节点，正常递归处理
	optimized := &model.FlameGraphData{
		Name:     shortName(node.Name),
		FullName: node.Name,
		Value:    node.Value,
		Children: make([]*model.FlameGraphData, 0, len(node.Children)),
	}

	for _, child := range node.Children {
		optimized.Children = append(optimized.Children, optimizeNode(child, depth+1, maxDepth))
	}

	return optimized
}

// shortName 取函数名最后一段（去掉包路径前缀）
func shortName(name string) string {
	// 取分号分隔的最后一段
	if idx := strings.LastIndex(name, ";"); idx >= 0 {
		name = name[idx+1:]
	}
	// 对 Java 类名取最后一段（去掉包名）
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	return name
}

// extractProfilerOutputPath 从 profiler 输出中提取文件路径
func extractProfilerOutputPath(output string) string {
	// Arthas profiler 输出格式示例：
	// profiler output file: /tmp/arthas-output/20240305-143022.txt
	re := regexp.MustCompile(`(?i)profiler output file:\s*([^\s\n]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}

	// 备用正则：直接匹配文件路径
	re2 := regexp.MustCompile(`(/[^\s]+\.(txt|collapsed))`)
	matches2 := re2.FindStringSubmatch(output)
	if len(matches2) > 1 {
		return matches2[1]
	}

	return ""
}

// parseCollapsedFormat 解析 collapsed 格式为树形结构
// collapsed 格式示例：
// java/lang/Thread.run;org/apache/catalina/core/StandardEngineValve.invoke 100
// java/lang/Thread.run;org/apache/catalina/core/StandardEngineValve.invoke;com/example/MyService.process 50
func parseCollapsedFormat(collapsed string) (*model.FlameGraphData, int64) {
	root := &model.FlameGraphNode{
		Name:     "all",
		Value:    0,
		Children: make(map[string]*model.FlameGraphNode),
	}

	var totalSamples int64

	lines := strings.Split(collapsed, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析每一行：stack_trace sample_count
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		stackTrace := parts[0]
		sampleCount := int64(0)
		fmt.Sscanf(parts[1], "%d", &sampleCount)

		if sampleCount == 0 {
			continue
		}

		totalSamples += sampleCount

		// 分割调用栈（用分号分隔）
		frames := strings.Split(stackTrace, ";")

		// 构建树形结构
		currentNode := root
		for _, frame := range frames {
			frame = strings.TrimSpace(frame)
			if frame == "" {
				continue
			}

			// 如果子节点不存在，创建新节点
			if currentNode.Children[frame] == nil {
				currentNode.Children[frame] = &model.FlameGraphNode{
					Name:     frame,
					Value:    0,
					Children: make(map[string]*model.FlameGraphNode),
				}
			}

			currentNode = currentNode.Children[frame]
		}

		// 在叶子节点累加样本数
		currentNode.Value += sampleCount
	}

	// 递归计算每个节点的总值（包括子节点）
	calculateNodeValues(root)

	// 转换为前端需要的格式
	return convertToFlameGraphData(root), totalSamples
}

// calculateNodeValues 递归计算节点值（自底向上累加）
func calculateNodeValues(node *model.FlameGraphNode) int64 {
	if len(node.Children) == 0 {
		return node.Value
	}

	total := node.Value
	for _, child := range node.Children {
		total += calculateNodeValues(child)
	}
	node.Value = total
	return total
}

// convertToFlameGraphData 转换为前端需要的数据格式
func convertToFlameGraphData(node *model.FlameGraphNode) *model.FlameGraphData {
	data := &model.FlameGraphData{
		Name:  node.Name,
		Value: node.Value,
	}

	if len(node.Children) > 0 {
		data.Children = make([]*model.FlameGraphData, 0, len(node.Children))
		for _, child := range node.Children {
			data.Children = append(data.Children, convertToFlameGraphData(child))
		}
	}

	return data
}

// copyFileFromPod 从 Pod 复制文件到本地（返回文件内容）
func (s *K8sArthasService) copyFileFromPod(clusterId uint, namespace, pod, container, remotePath string) ([]byte, error) {
	// 使用 cat 命令读取文件内容
	catCmd := fmt.Sprintf("cat %s", remotePath)
	output, err := s.execInPod(clusterId, namespace, pod, container, []string{"sh", "-c", catCmd})
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	return []byte(output), nil
}

// ArthasWebSocket WebSocket实时交互会话
func (s *K8sArthasService) ArthasWebSocket(clusterId uint, namespace, pod, container, processID string, conn *websocket.Conn) {
	defer conn.Close()

	clientset, config, err := s.getK8sClientWithConfig(clusterId)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": fmt.Sprintf("获取集群信息失败: %v", err)})
		return
	}

	// 第一步：确保 Pod 内有 arthas-boot.jar
	if err := s.ensureArthasInPod(clusterId, namespace, pod, container); err != nil {
		conn.WriteJSON(map[string]string{"error": fmt.Sprintf("无法获取 Arthas: %v", err)})
		return
	}

	// 第二步：启动交互式Arthas会话
	arthasCmd := fmt.Sprintf("java -jar /tmp/arthas-boot.jar --pid %s", processID)

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container).
		Param("command", arthasCmd).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true")

	// 创建executor
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		conn.WriteJSON(map[string]string{"error": fmt.Sprintf("创建executor失败: %v", err)})
		return
	}

	// 创建上下文和管道
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	// 创建流对象，复用Terminal的WebSocketStream实现
	stream := &K8sWebSocketStream{
		Conn:   conn,
		reader: reader,
		writer: writer,
		Ctx:    ctx,
		cancel: cancel,
	}

	// 创建终端连接
	termConn := &terminalConn{stream: stream}

	// 启动WebSocket读取goroutine
	go stream.ReadFromWebSocket()

	// 启动K8s executor goroutine
	executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  termConn,
		Stdout: termConn,
		Stderr: termConn,
		Tty:    true,
	})

	stream.Close()
}
