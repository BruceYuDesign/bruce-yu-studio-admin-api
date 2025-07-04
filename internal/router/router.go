package router

import (
	"net/http"

	"bruce-yu-studio-admin-api/internal/handler"
	"bruce-yu-studio-admin-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 路由器並設定路由
func InitRouter(userHandler *handler.UserHandler) *gin.Engine {
	router := gin.New()

	// 全局中間件
	router.Use(gin.Logger())   // Gin 內建日誌中間件
	router.Use(gin.Recovery()) // Gin 內建恢復中間件，防止 Panic 導致伺服器崩潰

	// 自定義日誌中間件 (可選，已在 main 中初始化 logger)
	router.Use(func(c *gin.Context) {
		logger.Debug("請求進入: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next() // 執行下一個中間件或路由處理器
	})

	// 健康檢查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API 版本路由組
	apiV1 := router.Group("/api/v1")
	{
		// 用戶相關路由
		userGroup := apiV1.Group("/users")
		{
			userGroup.POST("/register", userHandler.RegisterUser)
			userGroup.GET("/:id", userHandler.GetUserProfile)
		}
	}

	return router
}
