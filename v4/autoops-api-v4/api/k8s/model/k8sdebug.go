package model

type CreatePodDebugContainerRequest struct {
	TargetContainerName string `json:"targetContainerName"` // 目标业务容器，不传默认第一个业务容器
	ContainerName       string `json:"containerName"`       // 调试容器名称，不传自动生成
	DebugImage          string `json:"debugImage"`          // 调试镜像，不传默认 crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/busybox:latest
	Privileged          *bool  `json:"privileged"`          // 是否启用特权模式，默认 true
}

type CreatePodDebugContainerResponse struct {
	Namespace           string `json:"namespace"`
	PodName             string `json:"podName"`
	SourcePodName       string `json:"sourcePodName,omitempty"`
	TargetContainerName string `json:"targetContainerName"`
	DebugContainerName  string `json:"debugContainerName"`
	DebugPodName        string `json:"debugPodName,omitempty"`
	DebugImage          string `json:"debugImage"`
	Created             bool   `json:"created"`
	Mode                string `json:"mode"`
	Message             string `json:"message"`
	TerminalPath        string `json:"terminalPath"`
}

type DeletePodDebugContainerResponse struct {
	Namespace          string `json:"namespace"`
	PodName            string `json:"podName"`
	DebugContainerName string `json:"debugContainerName"`
	DebugPodName       string `json:"debugPodName,omitempty"`
	Stopped            bool   `json:"stopped"`
	Mode               string `json:"mode"`
	Message            string `json:"message"`
}
