package storage

import (
	"io"
	"time"
)

// ImageInfo 图片信息
type ImageInfo struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Filename   string    `json:"filename"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	Path       string    `json:"path"`
	ThumbPath  string    `json:"thumb_path,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Storage 存储接口
type Storage interface {
	// Save 保存图片
	Save(userID string, filename string, mimeType string, content io.Reader) (*ImageInfo, error)

	// Get 获取图片信息
	Get(userID string, id string) (*ImageInfo, error)

	// Delete 删除图片
	Delete(userID string, id string) error

	// List 列出用户的所有图片
	List(userID string) ([]*ImageInfo, error)
}