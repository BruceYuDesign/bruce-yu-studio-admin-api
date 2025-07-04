package database

import (
	"bruce-yu-studio-admin-api/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitPostgresDB 初始化 PostgreSQL 資料庫連線
func InitPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("無法連線到 PostgreSQL 資料庫: %v", err)
		return nil, err
	}
	logger.Info("成功連線到 PostgreSQL 資料庫。")

	// 可以設定連接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("無法獲取底層 SQL DB 對象: %v", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)  // 最大空閒連接數
	sqlDB.SetMaxOpenConns(100) // 最大開啟連接數

	return db, nil
}
