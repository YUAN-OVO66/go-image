package service

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Config 系统配置
type Config struct {
	StorageLimit   int64 `json:"storage_limit"`   // 存储空间限制（字节）
	CurrentStorage int64 `json:"current_storage"` // 当前已使用存储空间（字节）
}

// ConfigService 处理系统配置相关的业务逻辑
type ConfigService struct {
	config     *Config
	configPath string
	mutex      sync.RWMutex
}

// NewConfigService 创建一个新的配置服务实例
func NewConfigService(configPath string) (*ConfigService, error) {
	if configPath == "" {
		configPath = "config.json"
	}

	service := &ConfigService{
		config:     &Config{},
		configPath: configPath,
	}

	// 尝试加载配置文件
	if err := service.loadConfig(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return service, nil
}

// InitConfig 初始化配置
func (s *ConfigService) InitConfig(storageLimit int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.config.StorageLimit = storageLimit
	s.config.CurrentStorage = 0

	return s.saveConfig()
}

// GetConfig 获取当前配置
func (s *ConfigService) GetConfig() *Config {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 返回配置的副本以避免并发修改
	configCopy := *s.config
	return &configCopy
}

// UpdateStorageUsage 更新存储空间使用量
func (s *ConfigService) UpdateStorageUsage(size int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newSize := s.config.CurrentStorage + size
	if newSize > s.config.StorageLimit {
		return errors.New("超出存储空间限制")
	}

	s.config.CurrentStorage = newSize
	return s.saveConfig()
}

// loadConfig 从文件加载配置
func (s *ConfigService) loadConfig() error {
	file, err := os.Open(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，使用默认配置
			s.config = &Config{
				StorageLimit:   1024 * 1024 * 1024, // 默认1GB
				CurrentStorage: 0,
			}
			return s.saveConfig()
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(s.config)
}

// saveConfig 保存配置到文件
func (s *ConfigService) saveConfig() error {
	// 确保配置目录存在
	if err := os.MkdirAll(filepath.Dir(s.configPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(s.configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(s.config)
}
