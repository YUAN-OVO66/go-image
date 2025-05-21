package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-image/internal/api"
	"github.com/yourusername/go-image/internal/middleware"
	"github.com/yourusername/go-image/internal/service"
	"github.com/yourusername/go-image/internal/storage"
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化存储目录
	initStorageDirs()

	// 初始化配置服务
	configService, err := service.NewConfigService("config.json")
	if err != nil {
		log.Fatalf("初始化配置服务失败: %v", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 设置会话存储
	store := cookie.NewStore([]byte("change-this-to-a-random-string"))
	r.Use(sessions.Sessions("go-image-session", store))

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 设置静态文件路由
	r.Static("/static", "./static")

	// 初始化存储服务
	fileStorage, err := storage.NewLocalStorage("static/uploads")
	if err != nil {
		log.Fatalf("初始化本地存储失败: %v", err)
	}

	// 初始化图片服务
	imageService := service.NewImageService(fileStorage)

	// 初始化认证服务
	authService := service.NewAuthService()

	// 注册中间件
	r.Use(middleware.Logger())
	r.Use(middleware.ConfigMiddleware(configService))

	// 初始化路由
	r.GET("/init", api.InitPageHandler(configService))
	r.POST("/init", api.InitHandler(configService, authService))

	// 公共路由
	r.GET("/", api.HomeHandler)
	r.GET("/login", api.LoginPageHandler)
	r.POST("/login", api.LoginHandler)
	r.GET("/logout", api.LogoutHandler)

	// 需要认证的路由组
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// 图片上传
		auth.GET("/upload", api.UploadPageHandler)
		auth.POST("/upload", api.UploadHandler(imageService))

		// 图片管理
		auth.GET("/images", api.ListImagesHandler(imageService))
		auth.GET("/images/:id", api.GetImageHandler(imageService))
		auth.DELETE("/images/:id", api.DeleteImageHandler(imageService))
	}

	// 公共图片访问
	r.GET("/i/:id", api.ServeImageHandler(imageService))

	// 启动服务器
	port := 28080
	log.Printf("服务器启动在 http://localhost:%d", port)
	r.Run(fmt.Sprintf(":%d", port))
}

// 初始化存储目录
func initStorageDirs() {
	dirs := []string{
		"static/uploads",
		"static/uploads/thumbnails",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("无法创建目录 %s: %v", dir, err)
		}
	}
}
