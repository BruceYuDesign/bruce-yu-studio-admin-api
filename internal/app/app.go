package app

import (
	"context" // 新增這一行
	"fmt"
	"net/http"
	"os"
	"time"

	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bruce-yu-studio-admin-api/config"
	"bruce-yu-studio-admin-api/internal/handler"
	"bruce-yu-studio-admin-api/internal/models"
	"bruce-yu-studio-admin-api/internal/repository"
	"bruce-yu-studio-admin-api/internal/router"
	"bruce-yu-studio-admin-api/internal/service"
	"bruce-yu-studio-admin-api/pkg/database" // 我們自己的資料庫初始化包
	"bruce-yu-studio-admin-api/pkg/logger"
)

// Application 結構體包含應用程式的所有核心組件和配置
type Application struct {
	config *config.Config
	engine *gin.Engine
	db     *gorm.DB
	server *http.Server // 添加 http.Server 以便優雅關閉
}

// NewApp 是一個建構函數，用於組裝和初始化所有應用程式組件
func NewApp(cfg *config.Config) *Application {
	// 1. 資料庫初始化與連線
	db, err := database.InitPostgresDB(cfg.Database.DSN)
	if err != nil {
		logger.Fatalf("無法初始化資料庫連線: %v", err) // 使用 logger.Fatalf 確保應用程式啟動失敗時退出
	}

	// 2. 資料庫自動遷移 (僅限開發/測試環境，生產環境建議使用專業遷移工具)
	// 在實際專案中，這裡會被 migration 工具取代
	if err := db.AutoMigrate(&models.User{}); err != nil {
		logger.Fatalf("資料庫遷移失敗: %v", err)
	}
	logger.Info("資料庫自動遷移完成。")

	// 3. 依賴注入 (DI)：從最底層的依賴開始實例化
	// 創建 Repository 實例
	userRepo := repository.NewUserRepository(db)

	// 創建 Service 實例，注入其依賴的 Repository
	userService := service.NewUserService(userRepo)

	// 創建 Handler 實例，注入其依賴的 Service
	userHandler := handler.NewUserHandler(userService)

	// 4. 初始化 Gin 路由器並設置路由
	engine := router.InitRouter(userHandler)

	return &Application{
		config: cfg,
		engine: engine,
		db:     db,
	}
}

// Run 啟動應用程式的 HTTP 伺服器並處理優雅關閉
func (a *Application) Run() error {
	addr := fmt.Sprintf(":%s", a.config.Server.Port)
	a.server = &http.Server{
		Addr:    addr,
		Handler: a.engine,
	}

	// 啟動伺服器在一個 Go routine 中，以便主 Go routine 監聽訊號
	go func() {
		logger.Info("🚀 伺服器正在 %s 端口運行...", addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("伺服器啟動失敗: %v", err)
		}
	}()

	// 監聽作業系統訊號，用於優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 監聽 Ctrl+C 或 kill 命令
	<-quit                                               // 阻塞直到接收到訊號

	logger.Info("👋 接收到關閉訊號，正在關閉伺服器...")

	// 優雅關閉伺服器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 這裡使用 Gin 的 context
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("伺服器強制關閉: %v", err)
		return err
	}

	logger.Info("✅ 伺服器已優雅關閉。")
	return nil
}

// Close 方法用於應用程式退出前釋放資源
func (a *Application) Close() {
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil {
			logger.Error("獲取資料庫實例失敗: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("關閉資料庫連線失敗: %v", err)
		} else {
			logger.Info("資料庫連線已關閉。")
		}
	}
	// 可以添加其他資源的關閉邏輯，例如 Redis 連線等
}
