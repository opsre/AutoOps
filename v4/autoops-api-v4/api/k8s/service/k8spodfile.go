package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type podFileFlushWriter struct {
	writer  io.Writer
	flusher http.Flusher
}

const (
	podFileRecordSeparator = "\x1e"
	podFileFieldSeparator  = "\x1f"
)

func (w *podFileFlushWriter) Write(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err == nil && w.flusher != nil {
		w.flusher.Flush()
	}
	return n, err
}

func (s *K8sWorkloadServiceImpl) getPodExecClient(clusterId uint) (*kubernetes.Clientset, *rest.Config, error) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		return nil, nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, nil, fmt.Errorf("解析集群配置失败: %v", err)
	}

	// 文件传输为长连接流式操作，不使用固定请求超时。
	config.Timeout = 0

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	return clientset, config, nil
}

func (s *K8sWorkloadServiceImpl) resolvePodExecContainer(ctx context.Context, clientset *kubernetes.Clientset, namespaceName, podName, containerName string) (string, error) {
	pod, err := clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("获取Pod信息失败: %v", err)
	}

	if len(pod.Spec.Containers) == 0 {
		return "", fmt.Errorf("Pod中没有可执行命令的业务容器")
	}

	if containerName == "" {
		return pod.Spec.Containers[0].Name, nil
	}

	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			return containerName, nil
		}
	}
	for _, container := range pod.Spec.EphemeralContainers {
		if container.Name == containerName {
			return containerName, nil
		}
	}

	return "", fmt.Errorf("容器不存在: %s", containerName)
}

func (s *K8sWorkloadServiceImpl) streamPodExec(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	config *rest.Config,
	namespaceName, podName, containerName string,
	command []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName)

	for _, arg := range command {
		req = req.Param("command", arg)
	}

	if stdin != nil {
		req = req.Param("stdin", "true")
	} else {
		req = req.Param("stdin", "false")
	}
	if stdout != nil {
		req = req.Param("stdout", "true")
	} else {
		req = req.Param("stdout", "false")
	}
	if stderr != nil {
		req = req.Param("stderr", "true")
	} else {
		req = req.Param("stderr", "false")
	}
	req = req.Param("tty", "false")

	executor, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		return fmt.Errorf("创建容器执行器失败: %v", err)
	}

	return executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
	})
}

func (s *K8sWorkloadServiceImpl) execPodCommand(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	config *rest.Config,
	namespaceName, podName, containerName string,
	command []string,
) (string, string, error) {
	var stdout, stderr bytes.Buffer
	err := s.streamPodExec(ctx, clientset, config, namespaceName, podName, containerName, command, nil, &stdout, &stderr)
	return stdout.String(), stderr.String(), err
}

func normalizeContainerAbsolutePath(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if !strings.HasPrefix(value, "/") {
		return "", fmt.Errorf("路径必须为绝对路径")
	}

	cleaned := path.Clean(value)
	if cleaned == "." {
		return "", fmt.Errorf("路径不能为空")
	}

	return cleaned, nil
}

func sanitizeUploadedFileName(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	fileName = strings.ReplaceAll(fileName, "\\", "/")
	fileName = path.Base(fileName)
	if fileName == "" || fileName == "." || fileName == "/" {
		return ""
	}
	return fileName
}

func buildAttachmentDisposition(fileName string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	return fmt.Sprintf(
		"attachment; filename=\"%s\"; filename*=UTF-8''%s",
		replacer.Replace(fileName),
		url.PathEscape(fileName),
	)
}

func parsePodFileMetadata(output string) (status string, size int64, err error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return "", 0, fmt.Errorf("容器文件元数据为空")
	}

	status = strings.TrimSpace(lines[0])
	if len(lines) > 1 {
		sizeLine := strings.TrimSpace(lines[1])
		if sizeLine != "" {
			size, err = strconv.ParseInt(sizeLine, 10, 64)
			if err != nil {
				return "", 0, fmt.Errorf("解析文件大小失败: %v", err)
			}
		}
	}

	return status, size, nil
}

func fileTransferErrorMessage(prefix string, execErr error, stdout, stderr string) string {
	detail := strings.TrimSpace(stderr)
	if detail == "" {
		detail = strings.TrimSpace(stdout)
	}
	if detail == "" && execErr != nil {
		detail = execErr.Error()
	}
	if detail == "" {
		detail = "未知错误"
	}
	return prefix + detail
}

func buildListPodFilesCommand(remotePath string) []string {
	script := `target="$1"; rs=$(printf '\036'); us=$(printf '\037');
if [ ! -e "$target" ]; then printf '__status__%smissing%s' "$us" "$rs"; exit 0; fi;
if [ ! -d "$target" ]; then printf '__status__%snot_directory%s' "$us" "$rs"; exit 0; fi;
if [ ! -r "$target" ]; then printf '__status__%sunreadable%s' "$us" "$rs"; exit 0; fi;
printf '__status__%sok%s' "$us" "$rs";
for entry in "$target"/* "$target"/.[!.]* "$target"/..?*; do
  [ -e "$entry" ] || continue;
  name=${entry##*/};
  [ "$name" = "." ] && continue;
  [ "$name" = ".." ] && continue;
  mode="";
  size="0";
  mod_time="";
  set -- $(ls -ldn "$entry" 2>/dev/null || printf '');
  if [ "$#" -ge 8 ]; then
    mode="$1";
    size="$5";
    mod_time="$6 $7 $8";
  fi;
  is_dir="false";
  if [ -d "$entry" ]; then is_dir="true"; fi;
  printf '%s%s%s%s%s%s%s%s%s%s' "$name" "$us" "$entry" "$us" "$size" "$us" "$mode" "$us" "$mod_time" "$us";
  printf '%s%s' "$is_dir" "$rs";
done`

	return []string{"sh", "-c", script, "list", remotePath}
}

func parsePodFilesOutput(output string) (string, []model.K8sPodFileInfo, error) {
	records := strings.Split(output, podFileRecordSeparator)
	files := make([]model.K8sPodFileInfo, 0, len(records))
	status := ""

	for _, record := range records {
		record = strings.TrimSuffix(record, "\n")
		if record == "" {
			continue
		}

		fields := strings.Split(record, podFileFieldSeparator)
		if len(fields) == 0 {
			continue
		}

		if fields[0] == "__status__" {
			if len(fields) < 2 {
				return "", nil, fmt.Errorf("目录状态输出格式无效")
			}
			status = fields[1]
			continue
		}

		if len(fields) < 6 {
			return "", nil, fmt.Errorf("目录项输出格式无效")
		}

		size, err := strconv.ParseInt(strings.TrimSpace(fields[2]), 10, 64)
		if err != nil {
			size = 0
		}

		files = append(files, model.K8sPodFileInfo{
			Name:    fields[0],
			Path:    fields[1],
			Size:    size,
			Mode:    fields[3],
			ModTime: fields[4],
			IsDir:   strings.EqualFold(fields[5], "true"),
		})
	}

	if status == "" {
		return "", nil, fmt.Errorf("目录状态缺失")
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return status, files, nil
}

// ListPodFiles 列出Pod容器目录文件
func (s *K8sWorkloadServiceImpl) ListPodFiles(c *gin.Context, clusterId uint, namespaceName, podName, containerName, remotePath string) {
	remotePath, err := normalizeContainerAbsolutePath(remotePath)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "远程目录路径无效: "+err.Error())
		return
	}

	clientset, config, err := s.getPodExecClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	lookupCtx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	containerName, err = s.resolvePodExecContainer(lookupCtx, clientset, namespaceName, podName, containerName)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	command := buildListPodFilesCommand(remotePath)
	stdout, stderr, err := s.execPodCommand(lookupCtx, clientset, config, namespaceName, podName, containerName, command)
	if err != nil {
		log.Printf("列出容器文件失败: cluster=%d namespace=%s pod=%s container=%s remotePath=%s err=%v stderr=%s", clusterId, namespaceName, podName, containerName, remotePath, err, strings.TrimSpace(stderr))
		result.Failed(c, http.StatusInternalServerError, fileTransferErrorMessage("列出容器文件失败: ", err, stdout, stderr))
		return
	}

	status, files, err := parsePodFilesOutput(stdout)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "解析容器目录失败: "+err.Error())
		return
	}

	switch status {
	case "ok":
	case "missing":
		result.Failed(c, http.StatusNotFound, "容器内目录不存在")
		return
	case "not_directory":
		result.Failed(c, http.StatusBadRequest, "指定路径不是目录")
		return
	case "unreadable":
		result.Failed(c, http.StatusForbidden, "容器内目录不可读")
		return
	default:
		result.Failed(c, http.StatusInternalServerError, "未知的容器目录状态: "+status)
		return
	}

	result.Success(c, gin.H{
		"namespace":     namespaceName,
		"podName":       podName,
		"containerName": containerName,
		"remotePath":    remotePath,
		"files":         files,
		"total":         len(files),
	})
}

// UploadPodFile 上传文件到Pod容器
func (s *K8sWorkloadServiceImpl) UploadPodFile(c *gin.Context, clusterId uint, namespaceName, podName, containerName, destPath, fileName string, file io.Reader) {
	destDir, err := normalizeContainerAbsolutePath(destPath)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "目标路径无效: "+err.Error())
		return
	}

	fileName = sanitizeUploadedFileName(fileName)
	if fileName == "" {
		result.Failed(c, http.StatusBadRequest, "上传文件名不能为空")
		return
	}

	clientset, config, err := s.getPodExecClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	lookupCtx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	containerName, err = s.resolvePodExecContainer(lookupCtx, clientset, namespaceName, podName, containerName)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	targetPath := path.Join(destDir, fileName)
	command := []string{
		"sh",
		"-c",
		`set -eu; target="$1"; dir="$2"; mkdir -p "$dir"; cat > "$target"`,
		"upload",
		targetPath,
		destDir,
	}

	var stdout, stderr bytes.Buffer
	log.Printf("开始上传容器文件: cluster=%d namespace=%s pod=%s container=%s target=%s", clusterId, namespaceName, podName, containerName, targetPath)

	if err := s.streamPodExec(c.Request.Context(), clientset, config, namespaceName, podName, containerName, command, file, &stdout, &stderr); err != nil {
		log.Printf("上传容器文件失败: cluster=%d namespace=%s pod=%s container=%s target=%s err=%v stderr=%s", clusterId, namespaceName, podName, containerName, targetPath, err, strings.TrimSpace(stderr.String()))
		result.Failed(c, http.StatusInternalServerError, fileTransferErrorMessage("上传文件到容器失败: ", err, stdout.String(), stderr.String()))
		return
	}

	log.Printf("上传容器文件成功: cluster=%d namespace=%s pod=%s container=%s target=%s", clusterId, namespaceName, podName, containerName, targetPath)
	result.Success(c, gin.H{
		"namespace":     namespaceName,
		"podName":       podName,
		"containerName": containerName,
		"fileName":      fileName,
		"destPath":      destDir,
		"targetPath":    targetPath,
		"message":       "文件上传成功",
	})
}

// DownloadPodFile 从Pod容器下载文件
func (s *K8sWorkloadServiceImpl) DownloadPodFile(c *gin.Context, clusterId uint, namespaceName, podName, containerName, remotePath string) {
	remotePath, err := normalizeContainerAbsolutePath(remotePath)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "远程文件路径无效: "+err.Error())
		return
	}

	clientset, config, err := s.getPodExecClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	lookupCtx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	containerName, err = s.resolvePodExecContainer(lookupCtx, clientset, namespaceName, podName, containerName)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	metadataCommand := []string{
		"sh",
		"-c",
		`target="$1"; if [ ! -e "$target" ]; then printf 'missing\n'; exit 0; fi; if [ -d "$target" ]; then printf 'directory\n'; exit 0; fi; if [ ! -r "$target" ]; then printf 'unreadable\n'; exit 0; fi; printf 'file\n'; wc -c < "$target" 2>/dev/null || true`,
		"download",
		remotePath,
	}

	stdout, stderr, err := s.execPodCommand(lookupCtx, clientset, config, namespaceName, podName, containerName, metadataCommand)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, fileTransferErrorMessage("获取容器文件信息失败: ", err, stdout, stderr))
		return
	}

	status, fileSize, err := parsePodFileMetadata(stdout)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "解析容器文件信息失败: "+err.Error())
		return
	}

	switch status {
	case "missing":
		result.Failed(c, http.StatusNotFound, "容器内文件不存在")
		return
	case "directory":
		result.Failed(c, http.StatusBadRequest, "暂不支持目录下载，请传入具体文件路径")
		return
	case "unreadable":
		result.Failed(c, http.StatusForbidden, "容器内文件不可读")
		return
	case "file":
	default:
		result.Failed(c, http.StatusInternalServerError, "未知的容器文件状态: "+status)
		return
	}

	downloadName := path.Base(remotePath)
	if downloadName == "" || downloadName == "." || downloadName == "/" {
		downloadName = "download.bin"
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", buildAttachmentDisposition(downloadName))
	c.Header("X-Accel-Buffering", "no")
	if fileSize > 0 {
		c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
	}

	var stdoutWriter io.Writer = c.Writer
	if flusher, ok := c.Writer.(http.Flusher); ok {
		stdoutWriter = &podFileFlushWriter{
			writer:  c.Writer,
			flusher: flusher,
		}
	}

	var streamStderr bytes.Buffer
	streamCommand := []string{
		"sh",
		"-c",
		`cat "$1"`,
		"download",
		remotePath,
	}

	log.Printf("开始下载容器文件: cluster=%d namespace=%s pod=%s container=%s remotePath=%s", clusterId, namespaceName, podName, containerName, remotePath)

	if err := s.streamPodExec(c.Request.Context(), clientset, config, namespaceName, podName, containerName, streamCommand, nil, stdoutWriter, &streamStderr); err != nil {
		log.Printf("下载容器文件失败: cluster=%d namespace=%s pod=%s container=%s remotePath=%s err=%v stderr=%s", clusterId, namespaceName, podName, containerName, remotePath, err, strings.TrimSpace(streamStderr.String()))
		if !c.Writer.Written() {
			result.Failed(c, http.StatusInternalServerError, fileTransferErrorMessage("下载容器文件失败: ", err, "", streamStderr.String()))
		}
		return
	}

	log.Printf("下载容器文件成功: cluster=%d namespace=%s pod=%s container=%s remotePath=%s", clusterId, namespaceName, podName, containerName, remotePath)
}
