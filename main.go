package main

import (
	"bruce-yu-studio-admin-api/config"
	"bruce-yu-studio-admin-api/internal/app"
	"bruce-yu-studio-admin-api/pkg/logger"
)

func main() {
	// 1. 載入配置
	cfg, err := config.LoadConfig()
	if err != nil {
		// 如果配置載入失敗，直接使用標準庫 log 報告致命錯誤並退出
		logger.Fatalf("無法載入配置: %v", err)
	}

	// 2. 初始化日誌系統 (使用我們自己的 logger)
	logger.InitLogger(cfg.Log.Level)
	logger.Info("配置載入成功並初始化日誌系統。")

	// 3. 創建應用程式實例
	// 所有的依賴注入和組件組裝邏輯都封裝在 NewApp 函數中
	application := app.NewApp(cfg)
	defer application.Close() // 確保在 main 函數結束時關閉資源

	// 4. 運行應用程式
	// Run 方法會阻塞直到伺服器停止或收到退出訊號
	if err := application.Run(); err != nil {
		logger.Error("應用程式運行失敗: %v", err)
	}

	logger.Info("應用程式已退出。")
}
