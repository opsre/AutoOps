// 存储层抽象接口
// 定义录制文件的存储、读取、删除等操作
package recording

import "io"

// Storage 存储接口
type Storage interface {
	// Save 保存录制文件
	// sessionID: 会话唯一标识
	// data: 文件内容读取器
	// 返回: 文件路径, 文件大小(字节), 错误
	Save(sessionID string, data io.Reader) (path string, size int64, err error)

	// Load 加载录制文件
	// sessionID: 会话唯一标识
	// 返回: 文件内容读取器(需要调用者关闭), 错误
	Load(sessionID string) (io.ReadCloser, error)

	// Delete 删除录制文件
	// sessionID: 会话唯一标识
	Delete(sessionID string) error

	// Exists 检查文件是否存在
	// sessionID: 会话唯一标识
	Exists(sessionID string) (bool, error)

	// GetPath 获取文件路径
	// sessionID: 会话唯一标识
	GetPath(sessionID string) string
}
