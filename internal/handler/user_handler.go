package handler

import (
	"net/http"
	"strconv"

	"bruce-yu-studio-admin-api/internal/service"
	"bruce-yu-studio-admin-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UserHandler 結構體
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 創建一個 UserHandler 實例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUserRequest 用戶註冊請求體
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterUser 處理用戶註冊 HTTP 請求
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("註冊請求參數綁定失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		logger.Error("用戶註冊業務邏輯失敗: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "用戶註冊成功",
		"user_id":   user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}

// GetUserProfile 處理獲取用戶資料 HTTP 請求
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Error("獲取用戶資料請求ID參數無效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的用戶 ID"})
		return
	}

	user, err := h.userService.GetUserProfile(uint(id))
	if err != nil {
		logger.Error("獲取用戶資料業務邏輯失敗: %v", err)
		if err.Error() == "用戶未找到" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服務內部錯誤"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}
