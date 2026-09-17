package model

import "time"

// FeishuUserToken 飞书 OAuth 用户令牌（方案A扩展：支持个人知识库访问）
// 用户通过 OAuth 授权后，后端代理调用其个人 Wiki 空间
type FeishuUserToken struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	KeyManageID      uint      `gorm:"not null;index" json:"keyManageId"`                             // 关联的飞书应用（KeyManage.KeyType=7）
	OpenID           string    `gorm:"type:varchar(100);uniqueIndex:uk_feishu_user" json:"openId"`   // 飞书用户唯一ID
	UnionID          string    `gorm:"type:varchar(100)" json:"unionId"`
	UserName         string    `gorm:"type:varchar(100)" json:"userName"`
	AvatarURL        string    `gorm:"type:varchar(500)" json:"avatarUrl"`
	AccessToken      string    `gorm:"type:text" json:"-"`        // AES 加密存储
	RefreshToken     string    `gorm:"type:text" json:"-"`        // AES 加密存储
	ExpiresAt        time.Time `json:"expiresAt"`                 // access_token 过期时间
	RefreshExpiresAt time.Time `json:"refreshExpiresAt"`          // refresh_token 过期时间（默认30天）
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (FeishuUserToken) TableName() string {
	return "tool_feishu_user_token"
}

// FeishuUserTokenVO 脱敏后对外展示（不含 token 字段）
type FeishuUserTokenVO struct {
	ID               uint      `json:"id"`
	KeyManageID      uint      `json:"keyManageId"`
	OpenID           string    `json:"openId"`
	UserName         string    `json:"userName"`
	AvatarURL        string    `json:"avatarUrl"`
	ExpiresAt        time.Time `json:"expiresAt"`
	RefreshExpiresAt time.Time `json:"refreshExpiresAt"`
	CreatedAt        time.Time `json:"createdAt"`
}
