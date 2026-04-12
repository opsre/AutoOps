// 鉴权中间件
// author xiaoRui

package middleware

import (
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 检查是否为websocket升级请求
		if c.GetHeader("Upgrade") == "websocket" {
			// 从URL参数获取token
			token := c.Query("token")

			// 如果URL没有token，尝试从Header获取
			if token == "" {
				authHeader := c.Request.Header.Get("Authorization")
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
				// 有些客户端会将token放在 Sec-WebSocket-Protocol 中
				if token == "" {
					protocols := c.Request.Header.Get("Sec-WebSocket-Protocol")
					if protocols != "" {
						parts := strings.Split(protocols, ",")
						for _, p := range parts {
							p = strings.TrimSpace(p)
							if p != "" {
								token = p
								break
							}
						}
					}
				}
			}

			if token != "" {
				mc, err := jwt.ValidateToken(token)
				if err == nil {
					c.Set(constant.ContextKeyUserObj, mc)
					c.Next()
					return
				}
			}

			// 直接返回HTTP 401未授权，终止握手
			result.Failed(c, int(result.ApiCode.NOAUTH), result.ApiCode.GetMessage(result.ApiCode.NOAUTH))
			c.Abort()
			return
		}

		// 常规HTTP请求鉴权 - 支持SSE的token query参数
		authHeader := c.Request.Header.Get("Authorization")

		// 如果没有Authorization头，检查是否为SSE连接并从query参数获取token
		if authHeader == "" {
			token := c.Query("token")
			if token != "" {
				// 验证token
				mc, err := jwt.ValidateToken(token)
				if err == nil {
					c.Set(constant.ContextKeyUserObj, mc)
					c.Next()
					return
				} else {
					// 添加调试信息
					fmt.Printf("[DEBUG] Token验证失败: %v, token长度: %d, Secret: %s\n", err, len(token), string(jwt.Secret))
				}
			}

			result.Failed(c, int(result.ApiCode.NOAUTH), result.ApiCode.GetMessage(result.ApiCode.NOAUTH))
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			result.Failed(c, int(result.ApiCode.AUTHFORMATERROR), result.ApiCode.GetMessage(result.ApiCode.AUTHFORMATERROR))
			c.Abort()
			return
		}
		mc, err := jwt.ValidateToken(parts[1])
		if err != nil {
			result.Failed(c, int(result.ApiCode.INVALIDTOKEN), result.ApiCode.GetMessage(result.ApiCode.INVALIDTOKEN))
			c.Abort()
			return
		}
		c.Set(constant.ContextKeyUserObj, mc)
		c.Next()
	}
}
