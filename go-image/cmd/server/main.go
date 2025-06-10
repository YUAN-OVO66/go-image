package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go-image/internal/api"
	"go-image/internal/middleware"
	"go-image/internal/service"
	"go-image/internal/storage"
)

func main() {

	// 初始化存储目录
	initStorageDirs()

	// 初始化配置服务
	configService, err := service.NewConfigService("config.json")
	if err != nil {
		log.Fatalf("初始化配置服务失败: %v", err)
	}
	// 设置默认存储限制
	if err := configService.InitConfig(1024 * 1024 * 1024); err != nil { // 1GB
		log.Printf("设置默认存储限制失败: %v", err)
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

	// 公共路由
	r.GET("/", api.HomeHandler)
	r.GET("/register", api.RegisterPageHandler)
	r.POST("/register", api.RegisterHandler)
	r.GET("/login", api.LoginPageHandler)
	r.POST("/login", api.LoginHandler)
	r.GET("/logout", api.LogoutHandler)
	r.GET("/api-docs", api.APIDocsHandler)

	// 需要认证的路由组
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// 图片上传
		auth.GET("/upload", api.UploadPageHandler)
		auth.POST("/upload", api.UploadHandler(imageService))

		// 图片管理
		auth.GET("/images", api.ListImagesHandler(imageService, configService))
		auth.GET("/images/:id", api.GetImageHandler(imageService))
		auth.DELETE("/images/:id", api.DeleteImageHandler(imageService))
	}

	// API路由组
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.APIAuthRequired(authService))
	{
		// 图片上传
		apiGroup.POST("/upload", api.APIUploadHandler(imageService))
		// 图片列表
		apiGroup.GET("/images", api.APIListImagesHandler(imageService))
		// 删除图片
		apiGroup.DELETE("/images/:id", api.APIDeleteImageHandler(imageService))
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
