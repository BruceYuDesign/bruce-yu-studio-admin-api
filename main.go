package main

import (
	"bruce-yu-studio-admin-api/config"
	"bruce-yu-studio-admin-api/internal/app"
	"bruce-yu-studio-admin-api/pkg/logger"
)

func main() {
	// 1. 載入配置
	cfg, cfgErr := config.LoadConfig()
	if cfgErr != nil {
		logger.Fatalf("無法載入配置: %v", cfgErr)
	}

	// 2. 初始化日誌
	logger.InitLogger(cfg.Log.Level)
	logger.Info("配置載入成功並初始化日誌系統。")

	// 3. 創建應用程式實例
	appInstance := app.NewApp(cfg)
	defer appInstance.Close()

	// 4. 運行應用程式
	appErr := appInstance.Run()
	if appErr != nil {
		logger.Error("應用程式運行失敗: %v", appErr)
		return
	}

	logger.Info("應用程式已退出。")
}
