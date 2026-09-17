package middleware

import (
	"dodevops-api/common/config"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// WebhookTokenAuth Webhook Token认证中间件
// 用于验证外部系统调用webhook接口的权限
func WebhookTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取配置的Token
		configuredToken := config.Config.Monitor.Webhook.Token

		// 如果配置中没有设置Token,则跳过认证(不安全,仅用于测试环境)
		if configuredToken == "" {
			c.Next()
			return
		}

		// 从请求头中获取Token
		// 支持两种方式:
		// 1. Authorization: Bearer <token>
		// 2. X-Webhook-Token: <token>
		token := c.GetHeader("X-Webhook-Token")
		if token == "" {
			// 尝试从 Authorization header 获取
			authHeader := c.GetHeader("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			}
		}

		// 验证Token
		if token != configuredToken {
			result.Failed(c, int(result.ApiCode.NOAUTH), "无效的Webhook认证Token")
			c.Abort()
			return
		}

		c.Next()
	}
}
