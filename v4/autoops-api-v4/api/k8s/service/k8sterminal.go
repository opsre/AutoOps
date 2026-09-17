package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// IK8sTerminalService 容器终端服务接口
type IK8sTerminalService interface {
	CreateK8sWebSocketStream(clusterId uint, namespaceName, podName, containerName, command string, conn *websocket.Conn) (*K8sWebSocketStream, error)
	GetPodContainers(clusterId uint, namespaceName, podName string) ([]string, error)
	CreatePodDebugContainer(clusterId uint, namespaceName, podName string, req *model.CreatePodDebugContainerRequest) (*model.CreatePodDebugContainerResponse, error)
	DeletePodDebugContainer(clusterId uint, namespaceName, podName, containerName string) (*model.DeletePodDebugContainerResponse, error)
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

const (
	debugManagedLabel              = "dodevops.io/debug-managed"
	debugModeLabel                 = "dodevops.io/debug-mode"
	debugModeEphemeralContainer    = "ephemeral-container"
	debugModeStandalonePod         = "standalone-pod"
	debugTargetPodAnnotation       = "dodevops.io/debug-target-pod"
	debugTargetPodUIDAnnotation    = "dodevops.io/debug-target-pod-uid"
	debugTargetContainerAnnotation = "dodevops.io/debug-target-container"
	debugImageAnnotation           = "dodevops.io/debug-image"
	debugContainerAnnotation       = "dodevops.io/debug-container"
)

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

func normalizeExecCommand(command string) (string, error) {
	switch strings.TrimSpace(strings.ToLower(command)) {
	case "", "bash", "/bin/bash":
		return "/bin/bash", nil
	case "sh", "/bin/sh":
		return "/bin/sh", nil
	default:
		return "", fmt.Errorf("仅支持 /bin/bash 或 /bin/sh")
	}
}

func (s *K8sTerminalService) getK8sClientWithConfig(clusterId uint) (*kubernetes.Clientset, *rest.Config, error) {
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

func trimContainerDisplaySuffix(name string) string {
	if idx := strings.Index(name, " "); idx > 0 {
		return name[:idx]
	}
	return name
}

func sanitizeDebugContainerName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "_", "-")
	name = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if name == "" {
		name = "debugger"
	}
	if len(name) > 63 {
		name = name[:63]
		name = strings.Trim(name, "-")
	}
	if name == "" {
		name = "debugger"
	}
	return name
}

func buildDebugContainerName(targetContainerName string) string {
	name := sanitizeDebugContainerName("debugger-" + targetContainerName)
	if len(name) > 63 {
		name = name[:63]
		name = strings.Trim(name, "-")
	}
	if name == "" {
		return "debugger"
	}
	return name
}

func buildStandaloneDebugPodGenerateName(sourcePodName, targetContainerName string) string {
	base := sanitizeDebugContainerName(fmt.Sprintf("debugger-%s-%s", sourcePodName, targetContainerName))
	if base == "" {
		base = "debugger"
	}
	return base + "-"
}

func debugContainerShellCommand() []string {
	return []string{
		"sh",
		"-c",
		`if command -v sleep >/dev/null 2>&1; then exec sleep 7d; elif command -v tail >/dev/null 2>&1; then exec tail -f /dev/null; else while true; do /bin/sh -c 'sleep 3600' 2>/dev/null || true; done; fi`,
	}
}

func boolPtr(value bool) *bool {
	return &value
}

func copyLocalObjectReferences(items []corev1.LocalObjectReference) []corev1.LocalObjectReference {
	if len(items) == 0 {
		return nil
	}
	result := make([]corev1.LocalObjectReference, len(items))
	copy(result, items)
	return result
}

func copyTolerations(items []corev1.Toleration) []corev1.Toleration {
	if len(items) == 0 {
		return nil
	}
	result := make([]corev1.Toleration, len(items))
	copy(result, items)
	return result
}

func shouldFallbackToStandaloneDebugPod(err error) bool {
	if err == nil {
		return false
	}
	if apierrors.IsNotFound(err) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "the server could not find the requested resource")
}

func (s *K8sTerminalService) findPodContainer(pod *corev1.Pod, containerName string) bool {
	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			return true
		}
	}
	for _, container := range pod.Spec.EphemeralContainers {
		if container.Name == containerName {
			return true
		}
	}
	return false
}

func (s *K8sTerminalService) resolveTargetContainer(pod *corev1.Pod, containerName string) (string, error) {
	containerName = trimContainerDisplaySuffix(containerName)
	if containerName == "" {
		if len(pod.Spec.Containers) == 0 {
			return "", fmt.Errorf("Pod中没有业务容器")
		}
		return pod.Spec.Containers[0].Name, nil
	}

	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			return container.Name, nil
		}
	}

	return "", fmt.Errorf("目标容器不存在: %s", containerName)
}

func (s *K8sTerminalService) findReusableDebugContainer(pod *corev1.Pod, targetContainerName, debugImage string) string {
	for _, container := range pod.Spec.EphemeralContainers {
		if container.TargetContainerName == targetContainerName &&
			container.Image == debugImage &&
			strings.HasPrefix(container.Name, "debugger-") &&
			!s.isEphemeralContainerTerminated(pod, container.Name) {
			return container.Name
		}
	}
	return ""
}

func (s *K8sTerminalService) debugContainerExists(pod *corev1.Pod, containerName string) bool {
	for _, container := range pod.Spec.EphemeralContainers {
		if container.Name == containerName {
			return true
		}
	}
	return false
}

func (s *K8sTerminalService) buildTerminalPath(clusterId uint, namespaceName, podName, containerName string) string {
	values := url.Values{}
	values.Set("containerName", containerName)
	values.Set("shell", "bash")
	return fmt.Sprintf("/api/v1/k8s/cluster/%d/namespaces/%s/pods/%s/terminal?%s",
		clusterId,
		pathEscapeSegment(namespaceName),
		pathEscapeSegment(podName),
		values.Encode(),
	)
}

func pathEscapeSegment(value string) string {
	return strings.ReplaceAll(url.PathEscape(value), "+", "%20")
}

func isStandaloneDebugPod(pod *corev1.Pod) bool {
	if pod == nil {
		return false
	}
	return pod.Labels[debugManagedLabel] == "true" && pod.Labels[debugModeLabel] == debugModeStandalonePod
}

func (s *K8sTerminalService) getPodPrimaryContainerName(pod *corev1.Pod) string {
	if pod == nil || len(pod.Spec.Containers) == 0 {
		return ""
	}
	return pod.Spec.Containers[0].Name
}

func (s *K8sTerminalService) standaloneDebugPodMatches(pod *corev1.Pod, sourcePodName, sourcePodUID, targetContainerName, debugImage, debugContainerName string) bool {
	if !isStandaloneDebugPod(pod) || pod.DeletionTimestamp != nil {
		return false
	}
	if pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == corev1.PodSucceeded {
		return false
	}

	annotations := pod.Annotations
	if annotations[debugTargetPodAnnotation] != sourcePodName {
		return false
	}
	if sourcePodUID != "" {
		if uid := annotations[debugTargetPodUIDAnnotation]; uid != "" && uid != sourcePodUID {
			return false
		}
	}
	if targetContainerName != "" && annotations[debugTargetContainerAnnotation] != targetContainerName {
		return false
	}
	if debugImage != "" && annotations[debugImageAnnotation] != debugImage {
		return false
	}
	if debugContainerName != "" {
		actualContainerName := annotations[debugContainerAnnotation]
		if actualContainerName == "" {
			actualContainerName = s.getPodPrimaryContainerName(pod)
		}
		if actualContainerName != debugContainerName {
			return false
		}
	}

	return true
}

func (s *K8sTerminalService) findReusableStandaloneDebugPod(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, sourcePod *corev1.Pod, targetContainerName, debugImage string) (*corev1.Pod, error) {
	pods, err := clientset.CoreV1().Pods(namespaceName).List(ctx, metav1.ListOptions{
		LabelSelector: debugManagedLabel + "=true," + debugModeLabel + "=" + debugModeStandalonePod,
	})
	if err != nil {
		return nil, fmt.Errorf("查询独立调试Pod失败: %v", err)
	}

	var pending *corev1.Pod
	for i := range pods.Items {
		pod := &pods.Items[i]
		if !s.standaloneDebugPodMatches(pod, sourcePod.Name, string(sourcePod.UID), targetContainerName, debugImage, "") {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning {
			return pod, nil
		}
		if pending == nil {
			pending = pod
		}
	}

	return pending, nil
}

func (s *K8sTerminalService) findStandaloneDebugPodForDelete(ctx context.Context, clientset *kubernetes.Clientset, namespaceName, sourcePodName, sourcePodUID, debugContainerName string) (*corev1.Pod, error) {
	pods, err := clientset.CoreV1().Pods(namespaceName).List(ctx, metav1.ListOptions{
		LabelSelector: debugManagedLabel + "=true," + debugModeLabel + "=" + debugModeStandalonePod,
	})
	if err != nil {
		return nil, fmt.Errorf("查询独立调试Pod失败: %v", err)
	}

	for i := range pods.Items {
		pod := &pods.Items[i]
		if s.standaloneDebugPodMatches(pod, sourcePodName, sourcePodUID, "", "", debugContainerName) {
			return pod, nil
		}
	}

	return nil, nil
}

func (s *K8sTerminalService) buildStandaloneDebugPod(namespaceName string, sourcePod *corev1.Pod, targetContainerName, debugContainerName, debugImage string, privileged bool) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: buildStandaloneDebugPodGenerateName(sourcePod.Name, targetContainerName),
			Namespace:    namespaceName,
			Labels: map[string]string{
				debugManagedLabel: "true",
				debugModeLabel:    debugModeStandalonePod,
			},
			Annotations: map[string]string{
				debugTargetPodAnnotation:       sourcePod.Name,
				debugTargetPodUIDAnnotation:    string(sourcePod.UID),
				debugTargetContainerAnnotation: targetContainerName,
				debugImageAnnotation:           debugImage,
				debugContainerAnnotation:       debugContainerName,
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy:                corev1.RestartPolicyNever,
			ServiceAccountName:           sourcePod.Spec.ServiceAccountName,
			AutomountServiceAccountToken: sourcePod.Spec.AutomountServiceAccountToken,
			ImagePullSecrets:             copyLocalObjectReferences(sourcePod.Spec.ImagePullSecrets),
			NodeName:                     sourcePod.Spec.NodeName,
			Tolerations:                  copyTolerations(sourcePod.Spec.Tolerations),
			Containers: []corev1.Container{
				{
					Name:                     debugContainerName,
					Image:                    debugImage,
					ImagePullPolicy:          corev1.PullIfNotPresent,
					Command:                  debugContainerShellCommand(),
					Stdin:                    true,
					TTY:                      true,
					TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
					SecurityContext: &corev1.SecurityContext{
						Privileged:               boolPtr(privileged),
						AllowPrivilegeEscalation: boolPtr(privileged),
					},
				},
			},
		},
	}
}

func (s *K8sTerminalService) buildStandaloneDebugResponse(clusterId uint, namespaceName string, sourcePod *corev1.Pod, targetContainerName, debugContainerName, debugPodName, debugImage, message string, created bool) *model.CreatePodDebugContainerResponse {
	return &model.CreatePodDebugContainerResponse{
		Namespace:           namespaceName,
		PodName:             debugPodName,
		SourcePodName:       sourcePod.Name,
		TargetContainerName: targetContainerName,
		DebugContainerName:  debugContainerName,
		DebugPodName:        debugPodName,
		DebugImage:          debugImage,
		Created:             created,
		Mode:                debugModeStandalonePod,
		Message:             message,
		TerminalPath:        s.buildTerminalPath(clusterId, namespaceName, debugPodName, debugContainerName),
	}
}

func (s *K8sTerminalService) createStandaloneDebugPod(ctx context.Context, clientset *kubernetes.Clientset, clusterId uint, namespaceName string, sourcePod *corev1.Pod, targetContainerName, debugContainerName, debugImage string, privileged bool) (*model.CreatePodDebugContainerResponse, error) {
	existingPod, err := s.findReusableStandaloneDebugPod(ctx, clientset, namespaceName, sourcePod, targetContainerName, debugImage)
	if err != nil {
		return nil, err
	}
	if existingPod != nil {
		actualContainerName := existingPod.Annotations[debugContainerAnnotation]
		if actualContainerName == "" {
			actualContainerName = s.getPodPrimaryContainerName(existingPod)
		}
		return s.buildStandaloneDebugResponse(
			clusterId,
			namespaceName,
			sourcePod,
			targetContainerName,
			actualContainerName,
			existingPod.Name,
			debugImage,
			"集群不支持临时调试容器，已复用独立调试Pod",
			false,
		), nil
	}

	debugPod := s.buildStandaloneDebugPod(namespaceName, sourcePod, targetContainerName, debugContainerName, debugImage, privileged)
	createdPod, err := clientset.CoreV1().Pods(namespaceName).Create(ctx, debugPod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("创建独立调试Pod失败: %v", err)
	}

	return s.buildStandaloneDebugResponse(
		clusterId,
		namespaceName,
		sourcePod,
		targetContainerName,
		debugContainerName,
		createdPod.Name,
		debugImage,
		"集群不支持临时调试容器，已创建独立调试Pod",
		true,
	), nil
}

func (s *K8sTerminalService) deleteStandaloneDebugPod(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, debugPod *corev1.Pod, containerName string) (*model.DeletePodDebugContainerResponse, error) {
	actualContainerName := debugPod.Annotations[debugContainerAnnotation]
	if actualContainerName == "" {
		actualContainerName = s.getPodPrimaryContainerName(debugPod)
	}
	if actualContainerName != containerName {
		return nil, fmt.Errorf("独立调试Pod中不存在调试容器: %s", containerName)
	}

	gracePeriodSeconds := int64(0)
	if err := clientset.CoreV1().Pods(namespaceName).Delete(ctx, debugPod.Name, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
	}); err != nil {
		return nil, fmt.Errorf("删除独立调试Pod失败: %v", err)
	}

	targetPodName := debugPod.Annotations[debugTargetPodAnnotation]
	if targetPodName == "" {
		targetPodName = debugPod.Name
	}

	return &model.DeletePodDebugContainerResponse{
		Namespace:          namespaceName,
		PodName:            targetPodName,
		DebugContainerName: containerName,
		DebugPodName:       debugPod.Name,
		Stopped:            true,
		Mode:               debugModeStandalonePod,
		Message:            "独立调试Pod已删除",
	}, nil
}

func (s *K8sTerminalService) isEphemeralContainerTerminated(pod *corev1.Pod, containerName string) bool {
	for _, status := range pod.Status.EphemeralContainerStatuses {
		if status.Name == containerName {
			return status.State.Terminated != nil
		}
	}
	return false
}

func (s *K8sTerminalService) isEphemeralContainerRunning(pod *corev1.Pod, containerName string) bool {
	for _, status := range pod.Status.EphemeralContainerStatuses {
		if status.Name == containerName {
			return status.State.Running != nil
		}
	}
	return false
}

func (s *K8sTerminalService) getEphemeralContainer(pod *corev1.Pod, containerName string) (*corev1.EphemeralContainer, bool) {
	for i := range pod.Spec.EphemeralContainers {
		if pod.Spec.EphemeralContainers[i].Name == containerName {
			return &pod.Spec.EphemeralContainers[i], true
		}
	}
	return nil, false
}

func (s *K8sTerminalService) execCommandInContainer(config *rest.Config, clientset *kubernetes.Clientset, namespaceName, podName, containerName string, command []string) (string, string, error) {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName)

	for _, arg := range command {
		req.Param("command", arg)
	}

	req.Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("创建executor失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	return stdout.String(), stderr.String(), err
}

func (s *K8sTerminalService) probeTerminalCommand(config *rest.Config, clientset *kubernetes.Clientset, namespaceName, podName, containerName, command string) (string, error) {
	candidates := []string{command}
	if command == "/bin/bash" {
		candidates = append(candidates, "/bin/sh")
	}

	var lastErr error
	for _, candidate := range candidates {
		if err := s.runTerminalProbe(config, clientset, namespaceName, podName, containerName, candidate); err == nil {
			return candidate, nil
		} else {
			lastErr = err
		}
	}

	return "", fmt.Errorf("容器终端不可用，已尝试 /bin/bash 和 /bin/sh: %w", lastErr)
}

func (s *K8sTerminalService) runTerminalProbe(config *rest.Config, clientset *kubernetes.Clientset, namespaceName, podName, containerName, command string) error {
	var probeArgs []string
	switch command {
	case "/bin/bash":
		probeArgs = []string{"/bin/bash", "-lc", "true"}
	case "/bin/sh":
		probeArgs = []string{"/bin/sh", "-c", "true"}
	default:
		return fmt.Errorf("unsupported shell: %s", command)
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName)
	for _, arg := range probeArgs {
		req.Param("command", arg)
	}
	req.Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stdout, stderr bytes.Buffer
	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return fmt.Errorf("%w: %s", err, errMsg)
		}
		return err
	}

	return nil
}

// CreateK8sWebSocketStream 创建K8s WebSocket流
func (s *K8sTerminalService) CreateK8sWebSocketStream(clusterId uint, namespaceName, podName, containerName, command string, conn *websocket.Conn) (*K8sWebSocketStream, error) {
	command, err := normalizeExecCommand(command)
	if err != nil {
		return nil, err
	}

	containerName = trimContainerDisplaySuffix(containerName)

	clientset, config, err := s.getK8sClientWithConfig(clusterId)
	if err != nil {
		return nil, err
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

	command, err = s.probeTerminalCommand(config, clientset, namespaceName, podName, containerName, command)
	if err != nil {
		return nil, err
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
	clientset, _, err := s.getK8sClientWithConfig(clusterId)
	if err != nil {
		return nil, err
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
	for _, container := range pod.Spec.EphemeralContainers {
		containers = append(containers, container.Name+" (ephemeral)")
	}

	return containers, nil
}

func (s *K8sTerminalService) CreatePodDebugContainer(clusterId uint, namespaceName, podName string, req *model.CreatePodDebugContainerRequest) (*model.CreatePodDebugContainerResponse, error) {
	if req == nil {
		req = &model.CreatePodDebugContainerRequest{}
	}

	clientset, _, err := s.getK8sClientWithConfig(clusterId)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pod, err := clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取Pod信息失败: %v", err)
	}

	if pod.DeletionTimestamp != nil {
		return nil, fmt.Errorf("Pod正在删除中，无法创建调试容器")
	}
	if pod.Status.Phase != corev1.PodRunning {
		return nil, fmt.Errorf("仅支持对 Running 状态的Pod创建调试容器，当前状态: %s", pod.Status.Phase)
	}

	targetContainerName, err := s.resolveTargetContainer(pod, req.TargetContainerName)
	if err != nil {
		return nil, err
	}

	debugImage := strings.TrimSpace(req.DebugImage)
	if debugImage == "" {
		debugImage = "crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/busybox:latest"
	}

	if existing := s.findReusableDebugContainer(pod, targetContainerName, debugImage); existing != "" {
		return &model.CreatePodDebugContainerResponse{
			Namespace:           namespaceName,
			PodName:             podName,
			TargetContainerName: targetContainerName,
			DebugContainerName:  existing,
			DebugImage:          debugImage,
			Created:             false,
			Mode:                debugModeEphemeralContainer,
			Message:             "调试容器已存在，直接复用",
			TerminalPath:        s.buildTerminalPath(clusterId, namespaceName, podName, existing),
		}, nil
	}

	debugContainerName := strings.TrimSpace(req.ContainerName)
	if debugContainerName == "" {
		debugContainerName = buildDebugContainerName(targetContainerName)
	} else {
		debugContainerName = sanitizeDebugContainerName(debugContainerName)
	}

	if s.findPodContainer(pod, debugContainerName) {
		debugContainerName = sanitizeDebugContainerName(fmt.Sprintf("%s-%d", debugContainerName, time.Now().Unix()%100000))
	}

	privileged := true
	if req.Privileged != nil {
		privileged = *req.Privileged
	}

	podCopy := pod.DeepCopy()
	podCopy.Spec.EphemeralContainers = append(podCopy.Spec.EphemeralContainers, corev1.EphemeralContainer{
		EphemeralContainerCommon: corev1.EphemeralContainerCommon{
			Name:                     debugContainerName,
			Image:                    debugImage,
			ImagePullPolicy:          corev1.PullIfNotPresent,
			Command:                  debugContainerShellCommand(),
			Stdin:                    true,
			TTY:                      true,
			TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
			SecurityContext: &corev1.SecurityContext{
				Privileged:               boolPtr(privileged),
				AllowPrivilegeEscalation: boolPtr(privileged),
			},
		},
		TargetContainerName: targetContainerName,
	})

	if _, err := clientset.CoreV1().Pods(namespaceName).UpdateEphemeralContainers(ctx, podName, podCopy, metav1.UpdateOptions{}); err != nil {
		if shouldFallbackToStandaloneDebugPod(err) {
			return s.createStandaloneDebugPod(ctx, clientset, clusterId, namespaceName, pod, targetContainerName, debugContainerName, debugImage, privileged)
		}
		return nil, fmt.Errorf("创建临时调试容器失败: %v", err)
	}

	return &model.CreatePodDebugContainerResponse{
		Namespace:           namespaceName,
		PodName:             podName,
		TargetContainerName: targetContainerName,
		DebugContainerName:  debugContainerName,
		DebugImage:          debugImage,
		Created:             true,
		Mode:                debugModeEphemeralContainer,
		Message:             "临时调试容器创建成功",
		TerminalPath:        s.buildTerminalPath(clusterId, namespaceName, podName, debugContainerName),
	}, nil
}

func (s *K8sTerminalService) DeletePodDebugContainer(clusterId uint, namespaceName, podName, containerName string) (*model.DeletePodDebugContainerResponse, error) {
	containerName = trimContainerDisplaySuffix(containerName)
	if containerName == "" {
		return nil, fmt.Errorf("调试容器名称不能为空")
	}

	clientset, config, err := s.getK8sClientWithConfig(clusterId)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pod, err := clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if err == nil && isStandaloneDebugPod(pod) {
		return s.deleteStandaloneDebugPod(ctx, clientset, namespaceName, pod, containerName)
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, fmt.Errorf("获取Pod信息失败: %v", err)
	}

	if err == nil {
		if _, ok := s.getEphemeralContainer(pod, containerName); ok {
			if s.isEphemeralContainerTerminated(pod, containerName) {
				return &model.DeletePodDebugContainerResponse{
					Namespace:          namespaceName,
					PodName:            podName,
					DebugContainerName: containerName,
					Stopped:            false,
					Mode:               debugModeEphemeralContainer,
					Message:            "调试容器已经是停止状态",
				}, nil
			}

			if !s.isEphemeralContainerRunning(pod, containerName) {
				return &model.DeletePodDebugContainerResponse{
					Namespace:          namespaceName,
					PodName:            podName,
					DebugContainerName: containerName,
					Stopped:            false,
					Mode:               debugModeEphemeralContainer,
					Message:            "调试容器当前未运行，暂时无法主动结束",
				}, nil
			}

			stdout, stderr, err := s.execCommandInContainer(config, clientset, namespaceName, podName, containerName, []string{"sh", "-c", "kill -TERM 1"})
			if err != nil {
				detail := strings.TrimSpace(stderr)
				if detail == "" {
					detail = strings.TrimSpace(stdout)
				}
				if detail == "" {
					detail = err.Error()
				}
				return nil, fmt.Errorf("结束调试容器失败: %s", detail)
			}

			return &model.DeletePodDebugContainerResponse{
				Namespace:          namespaceName,
				PodName:            podName,
				DebugContainerName: containerName,
				Stopped:            true,
				Mode:               debugModeEphemeralContainer,
				Message:            "调试容器已结束。注意：ephemeral container 不会从 Pod 定义中物理删除，只会变为 terminated",
			}, nil
		}
	}

	sourcePodUID := ""
	if err == nil {
		sourcePodUID = string(pod.UID)
	}

	standalonePod, err := s.findStandaloneDebugPodForDelete(ctx, clientset, namespaceName, podName, sourcePodUID, containerName)
	if err != nil {
		return nil, err
	}
	if standalonePod != nil {
		return s.deleteStandaloneDebugPod(ctx, clientset, namespaceName, standalonePod, containerName)
	}

	return nil, fmt.Errorf("临时调试容器不存在: %s", containerName)
}
