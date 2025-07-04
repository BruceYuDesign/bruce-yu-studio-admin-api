package repository

import (
	"errors"

	"bruce-yu-studio-admin-api/internal/models"

	"gorm.io/gorm"
)

// userRepo 結構體實現 UserRepository 介面
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository 創建一個 UserRepository 實例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// CreateUser 創建新用戶
func (r *userRepo) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByID 根據 ID 獲取用戶
func (r *userRepo) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到記錄，返回 nil User, nil error
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根據 Email 獲取用戶
func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到記錄，返回 nil User, nil error
		}
		return nil, err
	}
	return &user, nil
}
