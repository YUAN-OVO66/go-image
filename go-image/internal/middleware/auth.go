package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-image/internal/service"
)

// AuthRequired 验证用户是否已登录的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取配置服务
		configService, exists := c.Get("configService")
		if !exists {
			c.JSON(500, gin.H{"error": "系统配置服务未初始化"})
			c.Abort()
			return
		}

		// 检查系统是否已初始化
		if !configService.(*service.ConfigService).IsInitialized() {
			c.Redirect(302, "/init")
			c.Abort()
			return
		}

		// 检查用户是否已登录
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			// 用户未登录，重定向到登录页面
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中，以便后续处理器使用
		c.Set("user", user)
		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		c.Next()
		// 请求后
	}
}

// ConfigMiddleware 配置服务中间件
func ConfigMiddleware(configService *service.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将配置服务添加到上下文中
		c.Set("configService", configService)
		c.Next()
	}
}
