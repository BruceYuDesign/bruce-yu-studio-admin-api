package service

import "bruce-yu-studio-admin-api/internal/models"

// UserService 是一個介面，定義了用戶相關的業務邏輯操作
type UserService interface {
	RegisterUser(username, email, password string) (*models.User, error)
	GetUserProfile(id uint) (*models.User, error)
}
