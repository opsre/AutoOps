package middleware

import (
	"dodevops-api/api/system/service"
	"dodevops-api/common/result"
	"strings"

	"github.com/gin-gonic/gin"
)

func LicenseMiddleware() gin.HandlerFunc {
	whitelist := map[string]struct{}{
		"/api/v1/license/status":       {},
		"/api/v1/license/machine-code": {},
		"/api/v1/license/activate":     {},
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if _, ok := whitelist[path]; ok {
			c.Next()
			return
		}

		// 允许激活相关路径的末尾斜杠形式
		trimmed := strings.TrimRight(path, "/")
		if _, ok := whitelist[trimmed]; ok {
			c.Next()
			return
		}

		licenseService := service.NewLicenseService()
		availability, err := licenseService.GetAvailability()
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), "读取激活状态失败: "+err.Error())
			c.Abort()
			return
		}

		if !availability.Activated {
			result.Failed(c, int(result.ApiCode.LICENSENOTACTIVATED), result.ApiCode.GetMessage(result.ApiCode.LICENSENOTACTIVATED))
			c.Abort()
			return
		}

		if availability.Expired {
			result.Failed(c, int(result.ApiCode.LICENSEEXPIRED), result.ApiCode.GetMessage(result.ApiCode.LICENSEEXPIRED))
			c.Abort()
			return
		}

		c.Next()
	}
}
