// 本地文件存储实现
package recording

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalStorage 本地文件存储实现
type LocalStorage struct {
	basePath string // 基础存储路径，如: ./upload
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// 确保基础目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create base directory failed: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// generatePath 生成文件路径
// 路径格式: {basePath}/{YYYYMMDD}/session_{sessionID}.cast
// 例如: ./upload/20251230/session_abc123.cast
func (s *LocalStorage) generatePath(sessionID string) string {
	now := time.Now()
	return filepath.Join(
		s.basePath,
		now.Format("20060102"),  // 日期格式 YYYYMMDD
		fmt.Sprintf("session_%s.cast", sessionID),
	)
}

// GetPath 获取文件路径(实现Storage接口)
func (s *LocalStorage) GetPath(sessionID string) string {
	return s.generatePath(sessionID)
}

// Save 保存录制文件(实现Storage接口)
func (s *LocalStorage) Save(sessionID string, data io.Reader) (string, int64, error) {
	filePath := s.generatePath(sessionID)

	// 确保父目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", 0, fmt.Errorf("create directory failed: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("create file failed: %w", err)
	}
	defer file.Close()

	// 写入数据
	size, err := io.Copy(file, data)
	if err != nil {
		// 写入失败，删除已创建的文件
		os.Remove(filePath)
		return "", 0, fmt.Errorf("write file failed: %w", err)
	}

	return filePath, size, nil
}

// Load 加载录制文件(实现Storage接口)
func (s *LocalStorage) Load(sessionID string) (io.ReadCloser, error) {
	filePath := s.generatePath(sessionID)

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("recording file not found: %s", sessionID)
		}
		return nil, fmt.Errorf("open file failed: %w", err)
	}

	return file, nil
}

// Delete 删除录制文件(实现Storage接口)
func (s *LocalStorage) Delete(sessionID string) error {
	filePath := s.generatePath(sessionID)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，视为删除成功
		}
		return fmt.Errorf("delete file failed: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在(实现Storage接口)
func (s *LocalStorage) Exists(sessionID string) (bool, error) {
	filePath := s.generatePath(sessionID)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat file failed: %w", err)
	}

	return true, nil
}
