package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 應用程式的配置結構體
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 伺服器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

// LogConfig 日誌配置
type LogConfig struct {
	Level string `mapstructure:"level"` // debug, info, error, fatal
}

// LoadConfig 載入配置，從檔案和環境變數中讀取
func LoadConfig() (*Config, error) {
	// 嘗試載入 .env 檔案
	// Load .env ignores system environment variables
	// Load dotenv does not overwrite existing variables from the system
	// Load will load all the variables from .env to the system's environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 檔案，將使用環境變數或預設值。")
	}

	viper.AddConfigPath(".")      // 設置配置檔案路徑
	viper.SetConfigName("config") // 配置檔案名稱 (不含副檔名)
	viper.SetConfigType("yaml")   // 配置檔案類型

	// 自動綁定環境變數，例如 SERVER_PORT 會自動映射到 server.port
	// 注意：viper.AutomaticEnv() 需要在 ReadInConfig() 之前調用
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true) // 允許環境變數為空字串

	// 綁定特定環境變數到 viper 的鍵
	viper.BindEnv("server.port", "APP_PORT")
	viper.BindEnv("database.dsn", "DB_DSN")
	viper.BindEnv("log.level", "LOG_LEVEL")

	// 讀取配置檔案
	if err := viper.ReadInConfig(); err != nil {
		// 如果檔案不存在，但環境變數存在，可以繼續
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("未找到 config.yaml 檔案，將嘗試使用環境變數。")
		} else {
			return nil, err // 其他讀取錯誤則直接返回
		}
	}

	var cfg Config
	// 將讀取的配置綁定到 Config 結構體
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 檢查必要的環境變數是否已設置，如果沒有，則使用默認值或報錯
	if cfg.Server.Port == "" {
		cfg.Server.Port = os.Getenv("APP_PORT") // 再次嘗試從環境變數獲取
		if cfg.Server.Port == "" {
			cfg.Server.Port = "8080" // 設置默認值
			log.Printf("APP_PORT 未設置，使用默認端口 %s\n", cfg.Server.Port)
		}
	}
	if cfg.Database.DSN == "" {
		cfg.Database.DSN = os.Getenv("DB_DSN") // 再次嘗試從環境變數獲取
		if cfg.Database.DSN == "" {
			log.Println("DB_DSN 未設置，資料庫連線將會失敗！")
		}
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = os.Getenv("LOG_LEVEL") // 再次嘗試從環境變數獲取
		if cfg.Log.Level == "" {
			cfg.Log.Level = "info" // 設置默認日誌級別
			log.Printf("LOG_LEVEL 未設置，使用默認級別 %s\n", cfg.Log.Level)
		}
	}

	return &cfg, nil
}
