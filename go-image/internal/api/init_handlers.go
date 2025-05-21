package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-image/internal/service"
)

// InitPageHandler 显示初始化页面
func InitPageHandler(configService *service.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果系统已经初始化，重定向到首页
		if configService.IsInitialized() {
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.HTML(http.StatusOK, "init.html", gin.H{
			"title": "系统初始化 - Go-Image",
		})
	}
}

// InitHandler 处理系统初始化请求
func InitHandler(configService *service.ConfigService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果系统已经初始化，返回错误
		if configService.IsInitialized() {
			c.HTML(http.StatusBadRequest, "init.html", gin.H{
				"title": "系统初始化 - Go-Image",
				"error": "系统已经初始化",
			})
			return
		}

		// 获取表单数据
		username := c.PostForm("username")
		password := c.PostForm("password")
		confirmPassword := c.PostForm("confirm_password")
		storageLimit := c.PostForm("storage_limit")

		// 验证表单数据
		if username == "" || password == "" || confirmPassword == "" || storageLimit == "" {
			c.HTML(http.StatusBadRequest, "init.html", gin.H{
				"title": "系统初始化 - Go-Image",
				"error": "请填写所有必填字段",
			})
			return
		}

		// 验证密码一致性
		if password != confirmPassword {
			c.HTML(http.StatusBadRequest, "init.html", gin.H{
				"title": "系统初始化 - Go-Image",
				"error": "两次输入的密码不一致",
			})
			return
		}

		// 转换存储空间限制为字节
		limitGB, err := strconv.ParseInt(storageLimit, 10, 64)
		if err != nil || limitGB < 1 || limitGB > 1000 {
			c.HTML(http.StatusBadRequest, "init.html", gin.H{
				"title": "系统初始化 - Go-Image",
				"error": "存储空间限制必须在1-1000GB之间",
			})
			return
		}
		limitBytes := limitGB * 1024 * 1024 * 1024 // 转换为字节

		// 初始化系统配置
		if err := configService.InitializeSystem(username, password, limitBytes); err != nil {
			c.HTML(http.StatusInternalServerError, "init.html", gin.H{
				"title": "系统初始化 - Go-Image",
				"error": "初始化系统失败: " + err.Error(),
			})
			return
		}

		// 创建管理员账户
		if err := authService.AddUser(username, password, true); err != nil {
			c.HTML(http.StatusInternalServerError, "init.html", gin.H{
				"title": "系统初始化 - Go-Image",
				"error": "创建管理员账户失败: " + err.Error(),
			})
			return
		}

		// 重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
	}
}
