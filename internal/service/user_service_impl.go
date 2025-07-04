package service

import (
	"errors"
	"fmt"

	"bruce-yu-studio-admin-api/internal/models"
	"bruce-yu-studio-admin-api/internal/repository"
	"bruce-yu-studio-admin-api/pkg/logger" // 使用日誌

	"golang.org/x/crypto/bcrypt"
)

// userService 結構體實現 UserService 介面
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 創建一個 UserService 實例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterUser 處理用戶註冊業務邏輯
func (s *userService) RegisterUser(username, email, password string) (*models.User, error) {
	// 1. 檢查用戶是否已存在 (業務規則)
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		logger.Error("查詢用戶失敗: %v", err)
		return nil, fmt.Errorf("服務內部錯誤: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("電子郵件已被註冊")
	}

	// 2. 密碼加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密碼加密失敗: %v", err)
		return nil, fmt.Errorf("密碼處理失敗: %w", err)
	}

	// 3. 創建用戶模型
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// 4. 透過 Repository 持久化用戶
	if err := s.userRepo.CreateUser(user); err != nil {
		logger.Error("創建用戶失敗: %v", err)
		return nil, fmt.Errorf("用戶註冊失敗: %w", err)
	}

	logger.Info("用戶註冊成功: ID=%d, Email=%s", user.ID, user.Email)
	return user, nil
}

// GetUserProfile 獲取用戶資料業務邏輯
func (s *userService) GetUserProfile(id uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		logger.Error("根據 ID 查詢用戶失敗: %v", err)
		return nil, fmt.Errorf("查詢用戶失敗: %w", err)
	}
	if user == nil {
		return nil, errors.New("用戶未找到")
	}
	return user, nil
}
