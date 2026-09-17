package model

type K8sPodFileInfo struct {
	Name    string `json:"name"`    // 文件名
	Path    string `json:"path"`    // 完整路径
	Size    int64  `json:"size"`    // 文件大小（字节）
	Mode    string `json:"mode"`    // 权限模式
	ModTime string `json:"modTime"` // 修改时间
	IsDir   bool   `json:"isDir"`   // 是否是目录
}
