package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-image/internal/service"
)

// HomeHandler 处理首页请求
func HomeHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Go-Image 个人图床",
		"user":  user,
	})
}

// LoginPageHandler 显示登录页面
func LoginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录 - Go-Image",
	})
}

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 创建认证服务实例
	authService := service.NewAuthService()

	// 验证用户凭据
	user, err := authService.Authenticate(username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "登录 - Go-Image",
			"error": "用户名或密码不正确",
		})
		return
	}

	// 设置会话
	session := sessions.Default(c)
	session.Set("user", user.Username)
	session.Save()

	// 重定向到首页
	c.Redirect(http.StatusFound, "/")
}

// LogoutHandler 处理登出请求
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}

// UploadPageHandler 显示上传页面
func UploadPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "上传图片 - Go-Image",
	})
}

// UploadHandler 处理图片上传
func UploadHandler(imageService *service.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取上传的文件
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的图片"})
			return
		}

		// 检查文件大小
		const maxSize = 10 * 1024 * 1024 // 10MB
		if file.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过10MB"})
			return
		}

		// 上传图片
		imageInfo, err := imageService.UploadImage(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("上传失败: %v", err)})
			return
		}

		// 构建完整的URL
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
		imageURL := fmt.Sprintf("%s/i/%s", baseURL, imageInfo.ID)

		// 返回上传成功的信息
		c.JSON(http.StatusOK, gin.H{
			"message":  "上传成功",
			"id":       imageInfo.ID,
			"filename": imageInfo.Filename,
			"url":      imageURL,
		})
	}
}

// ListImagesHandler 列出所有图片
func ListImagesHandler(imageService *service.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		images, err := imageService.ListImages()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片列表失败"})
			return
		}

		// 获取配置服务
		configService, exists := c.Get("configService")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "配置服务未初始化"})
			return
		}

		// 获取存储空间信息
		config := configService.(*service.ConfigService).GetConfig()
		usedStorage := float64(config.CurrentStorage) / float64(1024*1024*1024) // 转换为GB
		totalStorage := float64(config.StorageLimit) / float64(1024*1024*1024)  // 转换为GB

		c.HTML(http.StatusOK, "images.html", gin.H{
			"title":        "我的图片 - Go-Image",
			"images":       images,
			"usedStorage":  fmt.Sprintf("%.2f GB", usedStorage),
			"totalStorage": fmt.Sprintf("%.2f GB", totalStorage),
		})
	}
}

// GetImageHandler 获取单个图片信息
func GetImageHandler(imageService *service.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		image, err := imageService.GetImage(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
			return
		}

		c.JSON(http.StatusOK, image)
	}
}

// DeleteImageHandler 删除图片
func DeleteImageHandler(imageService *service.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := imageService.DeleteImage(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除失败: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	}
}

// ServeImageHandler 提供图片访问
func ServeImageHandler(imageService *service.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		image, err := imageService.GetImage(id)
		if err != nil {
			c.String(http.StatusNotFound, "图片不存在")
			return
		}

		// 构建文件路径
		filePath := filepath.Join("./static/uploads", image.Path)

		// 检查是否请求缩略图
		if c.Query("thumb") == "1" && image.ThumbPath != "" {
			filePath = filepath.Join("./static/uploads/thumbnails", image.ThumbPath)
		}

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "图片文件不存在")
			return
		}

		// 提供文件
		c.File(filePath)
	}
}
