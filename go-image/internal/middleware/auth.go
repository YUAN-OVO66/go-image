package middleware

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-image/internal/service"
)

// APIAuthRequired API认证中间件
func APIAuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取认证信息
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
			c.Abort()
			return
		}

		// 验证用户凭据
		user, err := authService.Authenticate(username, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败"})
			c.Abort()
			return
		}

		// 将用户信息和ID存储在上下文中
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Next()
	}
}

// AuthRequired 验证用户是否已登录的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取会话
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
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
