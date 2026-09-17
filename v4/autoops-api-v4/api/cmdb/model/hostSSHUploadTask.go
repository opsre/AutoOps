package model

import "time"

type HostSSHUploadTaskStatus string

const (
	HostSSHUploadTaskPending   HostSSHUploadTaskStatus = "pending"
	HostSSHUploadTaskUploading HostSSHUploadTaskStatus = "uploading"
	HostSSHUploadTaskCompleted HostSSHUploadTaskStatus = "completed"
	HostSSHUploadTaskFailed    HostSSHUploadTaskStatus = "failed"
)

type HostSSHUploadTask struct {
	TaskID        string                  `json:"taskId"`
	HostID        uint                    `json:"hostId"`
	FileName      string                  `json:"fileName"`
	FileSize      int64                   `json:"fileSize"`
	DestPath      string                  `json:"destPath"`
	Status        HostSSHUploadTaskStatus `json:"status"`
	Progress      float64                 `json:"progress"`
	UploadedBytes int64                   `json:"uploadedBytes"`
	TotalBytes    int64                   `json:"totalBytes"`
	Message       string                  `json:"message"`
	Error         string                  `json:"error"`
	CreatedBy     uint                    `json:"createdBy"`
	CreatedByName string                  `json:"createdByName"`
	CreatedAt     time.Time               `json:"createdAt"`
	StartedAt     *time.Time              `json:"startedAt,omitempty"`
	CompletedAt   *time.Time              `json:"completedAt,omitempty"`
}
