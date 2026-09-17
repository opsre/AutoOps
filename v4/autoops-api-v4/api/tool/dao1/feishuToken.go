// 飞书用户 OAuth Token DAO
package dao

import (
	"dodevops-api/api/tool/model"
	. "dodevops-api/pkg/db"
)

// UpsertFeishuUserToken 创建或更新飞书用户 token（按 key_manage_id + open_id 唯一）
func UpsertFeishuUserToken(token *model.FeishuUserToken) error {
	return Db.Where("key_manage_id = ? AND open_id = ?", token.KeyManageID, token.OpenID).
		Assign(model.FeishuUserToken{
			UnionID:          token.UnionID,
			UserName:         token.UserName,
			AvatarURL:        token.AvatarURL,
			AccessToken:      token.AccessToken,
			RefreshToken:     token.RefreshToken,
			ExpiresAt:        token.ExpiresAt,
			RefreshExpiresAt: token.RefreshExpiresAt,
		}).
		FirstOrCreate(token).Error
}

// GetFeishuUserToken 根据 openId 获取 token 记录
func GetFeishuUserToken(keyManageID uint, openID string) (*model.FeishuUserToken, error) {
	var token model.FeishuUserToken
	err := Db.Where("key_manage_id = ? AND open_id = ?", keyManageID, openID).First(&token).Error
	return &token, err
}

// UpdateFeishuUserToken 更新 token（全量保存）
func UpdateFeishuUserToken(token *model.FeishuUserToken) error {
	return Db.Save(token).Error
}

// ListFeishuUserTokens 列出某应用下已授权的用户（脱敏，不含 token 字段）
func ListFeishuUserTokens(keyManageID uint) ([]model.FeishuUserTokenVO, error) {
	var list []model.FeishuUserTokenVO
	err := Db.Model(&model.FeishuUserToken{}).
		Where("key_manage_id = ?", keyManageID).
		Select("id, key_manage_id, open_id, user_name, avatar_url, expires_at, refresh_expires_at, created_at").
		Scan(&list).Error
	return list, err
}

// DeleteFeishuUserToken 撤销用户授权
func DeleteFeishuUserToken(keyManageID uint, openID string) error {
	return Db.Where("key_manage_id = ? AND open_id = ?", keyManageID, openID).
		Delete(&model.FeishuUserToken{}).Error
}
