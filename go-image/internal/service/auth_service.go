package service

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// User 表示系统用户
type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

// AuthService 认证服务
type AuthService struct {
	users map[string]*User
	mutex sync.RWMutex
	userFile string // 用户数据文件路径
}

// loadUsers 从文件加载用户数据
func (s *AuthService) loadUsers() error {
	// 确保data目录存在
	dataDir := filepath.Join("data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// 设置用户数据文件路径
	s.userFile = filepath.Join(dataDir, "users.json")

	// 读取文件
	file, err := os.Open(s.userFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码JSON数据
	return json.NewDecoder(file).Decode(&s.users)
}

// saveUsers 保存用户数据到文件
func (s *AuthService) saveUsers() error {
	// 创建文件
	file, err := os.Create(s.userFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 编码为JSON并写入文件
	return json.NewEncoder(file).Encode(s.users)
}

// NewAuthService 创建一个新的认证服务实例
func NewAuthService() *AuthService {
	service := &AuthService{
		users: make(map[string]*User),
	}
	
	// 从文件加载用户数据
	if err := service.loadUsers(); err != nil {
		// 如果文件不存在，创建空的用户数据
		if os.IsNotExist(err) {
			service.saveUsers()
		}
	}
	
	return service
}

// Authenticate 验证用户凭据
func (s *AuthService) Authenticate(username, password string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[username]
	if !exists || user.Password != password {
		return nil, errors.New("用户名或密码不正确")
	}

	return user, nil
}

// GetUser 获取用户信息
func (s *AuthService) GetUser(username string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, errors.New("用户不存在")
	}

	return user, nil
}

// Register 注册新用户
func (s *AuthService) Register(username, password string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[username]; exists {
		return errors.New("用户名已存在")
	}

	s.users[username] = &User{
		ID:        uuid.New().String(),
		Username:  username,
		Password:  password, // 在实际应用中应该使用加密密码
		CreatedAt: time.Now(),
	}

	// 保存用户数据到文件
	return s.saveUsers()
}
