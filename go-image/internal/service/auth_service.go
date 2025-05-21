package service

import (
	"errors"
	"sync"
)

// User 表示系统用户
type User struct {
	Username string
	Password string
	IsAdmin  bool
}

// AuthService 处理用户认证相关的业务逻辑
type AuthService struct {
	users map[string]*User
	mutex sync.RWMutex
}

// NewAuthService 创建一个新的认证服务实例
func NewAuthService() *AuthService {
	service := &AuthService{
		users: make(map[string]*User),
	}

	// 从配置文件加载管理员账户
	configService, err := NewConfigService("")
	if err == nil && configService.IsInitialized() {
		config := configService.GetConfig()
		if config.AdminUsername != "" && config.AdminPassword != "" {
			service.AddUser(config.AdminUsername, config.AdminPassword, true)
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

// AddUser 添加新用户
func (s *AuthService) AddUser(username, password string, isAdmin bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[username]; exists {
		// 如果是初始化过程中添加管理员，则更新现有用户
		if isAdmin {
			s.users[username] = &User{
				Username: username,
				Password: password, // 在实际应用中应该使用加密密码
				IsAdmin:  isAdmin,
			}
			return nil
		}
		return errors.New("用户名已存在")
	}

	s.users[username] = &User{
		Username: username,
		Password: password, // 在实际应用中应该使用加密密码
		IsAdmin:  isAdmin,
	}

	return nil
}
