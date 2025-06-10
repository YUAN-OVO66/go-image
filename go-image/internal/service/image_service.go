package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"go-image/internal/storage"
)

// ImageService 处理图片相关的业务逻辑
type ImageService struct {
	storage storage.Storage
}

// NewImageService 创建一个新的图片服务实例
func NewImageService(storage storage.Storage) *ImageService {
	return &ImageService{
		storage: storage,
	}
}

// UploadImage 处理图片上传
func (s *ImageService) UploadImage(userID string, file *multipart.FileHeader) (*storage.ImageInfo, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	// 读取文件内容到内存
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, src); err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	// 验证文件是否为图片
	fileContent := buffer.Bytes()
	contentType := detectContentType(fileContent)
	if !isImageType(contentType) {
		return nil, errors.New("不支持的文件类型，仅支持图片文件")
	}

	// 保存原始图片
	imageInfo, err := s.storage.Save(userID, file.Filename, contentType, bytes.NewReader(fileContent))
	if err != nil {
		return nil, fmt.Errorf("保存图片失败: %w", err)
	}

	// 生成缩略图
	if err := s.generateThumbnail(imageInfo, fileContent); err != nil {
		// 如果生成缩略图失败，记录错误但不影响上传
		fmt.Printf("生成缩略图失败: %v\n", err)
	}

	return imageInfo, nil
}

// GetImage 获取图片信息
func (s *ImageService) GetImage(userID string, id string) (*storage.ImageInfo, error) {
	return s.storage.Get(userID, id)
}

// DeleteImage 删除图片
func (s *ImageService) DeleteImage(userID string, id string) error {
	return s.storage.Delete(userID, id)
}

// ListImages 列出所有图片
func (s *ImageService) ListImages(userID string) ([]*storage.ImageInfo, error) {
	return s.storage.List(userID)
}

// 生成缩略图
func (s *ImageService) generateThumbnail(imageInfo *storage.ImageInfo, fileContent []byte) error {
	// 解码图片
	img, format, err := image.Decode(bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("解码图片失败: %w", err)
	}

	// 调整图片大小
	thumbnail := resize.Thumbnail(300, 300, img, resize.Lanczos3)

	// 创建缩略图文件
	thumbDir := filepath.Join("static/uploads/thumbnails")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return fmt.Errorf("创建缩略图目录失败: %w", err)
	}

	// 构建缩略图路径
	extension := filepath.Ext(imageInfo.Path)
	thumbPath := filepath.Join(thumbDir, fmt.Sprintf("%s_thumb%s", imageInfo.ID, extension))
	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return fmt.Errorf("创建缩略图文件失败: %w", err)
	}
	defer thumbFile.Close()

	// 根据原始图片格式保存缩略图
	switch format {
	case "jpeg":
		err = jpeg.Encode(thumbFile, thumbnail, nil)
	case "png":
		err = png.Encode(thumbFile, thumbnail)
	case "gif":
		err = gif.Encode(thumbFile, thumbnail, nil)
	default:
		// 默认使用JPEG格式
		err = jpeg.Encode(thumbFile, thumbnail, nil)
	}

	if err != nil {
		return fmt.Errorf("编码缩略图失败: %w", err)
	}

	// 更新图片信息，添加缩略图路径
	imageInfo.ThumbPath = filepath.Base(thumbPath)

	return nil
}

// 检测内容类型
func detectContentType(data []byte) string {
	return http.DetectContentType(data)
}

// 检查是否为支持的图片类型
func isImageType(contentType string) bool {
	supportedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	return supportedTypes[contentType]
}