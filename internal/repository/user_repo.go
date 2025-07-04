package repository

import "bruce-yu-studio-admin-api/internal/models"

// UserRepository 是一個介面，定義了與用戶資料持久化相關的操作
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	// UpdateUser(user *models.User) error // 範例可以先不實現
	// DeleteUser(id uint) error          // 範例可以先不實現
}
