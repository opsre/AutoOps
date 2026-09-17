// 录制辅助工具
package recording

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateSessionID 生成唯一的会话ID
// 格式: timestamp_randomhex
func GenerateSessionID() string {
	timestamp := time.Now().Unix()

	// 生成8字节随机数
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果随机数生成失败，使用纳秒时间戳
		return fmt.Sprintf("%d_%d", timestamp, time.Now().UnixNano())
	}

	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%d_%s", timestamp, randomHex)
}
