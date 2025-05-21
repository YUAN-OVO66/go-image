# Go-Image 个人图床软件

## 项目简介

Go-Image 是一个使用 Go 语言开发的轻量级个人图床软件，提供图片上传、存储、管理和分享功能。

## 功能特点

- 图片上传：支持拖拽上传、粘贴上传和选择文件上传
- 图片管理：查看、删除和搜索已上传的图片
- 图片分享：生成图片链接，方便分享到其他平台
- 用户认证：基本的用户登录功能，保护您的图片安全
- 图片处理：自动生成缩略图，支持图片压缩

## 技术栈

- 后端：Go 语言 + Gin 框架
- 存储：本地文件系统
- 前端：HTML + CSS + JavaScript

## 目录结构

```
├── cmd/                # 应用程序入口
│   └── server/         # 服务器启动代码
├── configs/            # 配置文件
├── internal/           # 内部包
│   ├── api/            # API 处理器
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── service/        # 业务逻辑
│   └── storage/        # 存储接口
├── pkg/                # 可重用的库
├── static/             # 静态资源
│   ├── css/            # 样式文件
│   ├── js/             # JavaScript 文件
│   └── uploads/        # 上传的图片
├── templates/          # HTML 模板
├── go.mod              # Go 模块文件
├── go.sum              # Go 依赖校验
└── README.md           # 项目说明
```

## 安装和使用

### 前置条件

- Go 1.16 或更高版本

### 安装步骤

1. 克隆仓库
   ```
   git clone https://github.com/YUAN-OVO66/go-image.git
   cd go-image
   ```

2. 安装依赖
   ```
   go mod tidy
   ```

3. 运行应用
   ```
   go run cmd/server/main.go
   ```

4. 访问应用
   打开浏览器，访问 `http://localhost:28080`

## 配置说明

配置文件位于 `configs/config.yaml`，可以根据需要修改以下配置：

- 服务器端口
- 上传文件大小限制
- 存储路径
- 是否启用用户认证

## 许可证

Apache-2.0 license