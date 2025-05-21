package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ImageInfo 存储图片的元数据
type ImageInfo struct {
	ID         string    `json:"id"`
	Filename   string    `json:"filename"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	Path       string    `json:"path"`
	ThumbPath  string    `json:"thumb_path,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Storage 定义存储接口
type Storage interface {
	Save(filename string, mimeType string, content io.Reader) (*ImageInfo, error)
	Get(id string) (*ImageInfo, error)
	Delete(id string) error
	List() ([]*ImageInfo, error)
}

// LocalStorage 实现本地文件系统存储
type LocalStorage struct {
	basePath string
	images   map[string]*ImageInfo
	mu       sync.RWMutex
	metaFile string
}

// NewLocalStorage 创建一个新的本地存储实例
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// 确保基础目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	// 确保缩略图目录存在
	thumbPath := filepath.Join(basePath, "thumbnails")
	if err := os.MkdirAll(thumbPath, 0755); err != nil {
		return nil, fmt.Errorf("创建缩略图目录失败: %w", err)
	}

	// 创建存储实例
	storage := &LocalStorage{
		basePath: basePath,
		images:   make(map[string]*ImageInfo),
		metaFile: filepath.Join(basePath, "metadata.json"),
	}

	// 加载元数据
	if err := storage.loadMetadata(); err != nil {
		// 如果文件不存在，不返回错误
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("加载元数据失败: %w", err)
		}
	}

	return storage, nil
}

// Save 保存图片到本地存储
func (s *LocalStorage) Save(filename string, mimeType string, content io.Reader) (*ImageInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 生成唯一ID
	id := uuid.New().String()

	// 确保文件名安全
	safeFilename := sanitizeFilename(filename)

	// 构建文件路径
	extension := filepath.Ext(safeFilename)
	if extension == "" {
		// 根据MIME类型添加默认扩展名
		switch mimeType {
		case "image/jpeg":
			extension = ".jpg"
		case "image/png":
			extension = ".png"
		case "image/gif":
			extension = ".gif"
		case "image/webp":
			extension = ".webp"
		default:
			extension = ".bin"
		}
	}

	// 构建存储路径
	relativePath := fmt.Sprintf("%s%s", id, extension)
	fullPath := filepath.Join(s.basePath, relativePath)

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入内容
	size, err := io.Copy(file, content)
	if err != nil {
		// 如果写入失败，删除文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	// 创建图片信息
	imageInfo := &ImageInfo{
		ID:         id,
		Filename:   safeFilename,
		Size:       size,
		MimeType:   mimeType,
		Path:       relativePath,
		UploadedAt: time.Now(),
	}

	// 保存到内存映射
	s.images[id] = imageInfo

	// 保存元数据
	if err := s.saveMetadata(); err != nil {
		return nil, fmt.Errorf("保存元数据失败: %w", err)
	}

	return imageInfo, nil
}

// Get 获取图片信息
func (s *LocalStorage) Get(id string) (*ImageInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	image, exists := s.images[id]
	if !exists {
		return nil, errors.New("图片不存在")
	}
	return image, nil
}

// Delete 删除图片
func (s *LocalStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	image, exists := s.images[id]
	if !exists {
		return errors.New("图片不存在")
	}

	// 删除文件
	fullPath := filepath.Join(s.basePath, image.Path)
	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	// 如果有缩略图，也删除
	if image.ThumbPath != "" {
		thumbPath := filepath.Join(s.basePath, "thumbnails", image.ThumbPath)
		os.Remove(thumbPath) // 忽略错误，因为缩略图可能不存在
	}

	// 从内存映射中删除
	delete(s.images, id)

	// 保存元数据
	if err := s.saveMetadata(); err != nil {
		return fmt.Errorf("保存元数据失败: %w", err)
	}

	return nil
}

// List 列出所有图片
func (s *LocalStorage) List() ([]*ImageInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var images []*ImageInfo
	for _, image := range s.images {
		images = append(images, image)
	}
	return images, nil
}

// sanitizeFilename 清理文件名，移除不安全字符
// saveMetadata 保存元数据到文件
func (s *LocalStorage) saveMetadata() error {
	data, err := json.MarshalIndent(s.images, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.metaFile, data, 0644)
}

// loadMetadata 从文件加载元数据
func (s *LocalStorage) loadMetadata() error {
	data, err := os.ReadFile(s.metaFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.images)
}

func sanitizeFilename(filename string) string {
	// 移除路径分隔符和其他不安全字符
	safe := strings.ReplaceAll(filename, "/", "")
	safe = strings.ReplaceAll(safe, "\\", "")
	safe = strings.ReplaceAll(safe, ":", "")
	safe = strings.ReplaceAll(safe, "*", "")
	safe = strings.ReplaceAll(safe, "?", "")
	safe = strings.ReplaceAll(safe, "\"", "")
	safe = strings.ReplaceAll(safe, "<", "")
	safe = strings.ReplaceAll(safe, ">", "")
	safe = strings.ReplaceAll(safe, "|", "")

	// 如果文件名为空，使用默认名称
	if safe == "" {
		safe = "image"
	}

	return safe
}
