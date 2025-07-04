package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// Logger 簡單的日誌接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{}) // Fatal 級別，會導致程式退出
}

// logLevel 定義日誌級別
type logLevel int

const (
	logLevelDebug logLevel = iota
	logLevelInfo
	logLevelError
	logLevelFatal
)

var (
	defaultLogger *appLogger
	once          sync.Once
)

// appLogger 實現 Logger 接口
type appLogger struct {
	level logLevel
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
	fatal *log.Logger
}

// InitLogger 初始化日誌系統，應在應用程式啟動時調用
func InitLogger(level string) {
	once.Do(func() {
		minLevel := parseLogLevel(level)
		defaultLogger = &appLogger{
			level: minLevel,
			info:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
			error: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
			debug: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
			fatal: log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})
}

// Info 記錄資訊日誌
func Info(msg string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= logLevelInfo {
		defaultLogger.info.Output(2, fmt.Sprintf(msg, args...))
	}
}

// Error 記錄錯誤日誌
func Error(msg string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= logLevelError {
		defaultLogger.error.Output(2, fmt.Sprintf(msg, args...))
	}
}

// Debug 記錄調試日誌
func Debug(msg string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= logLevelDebug {
		defaultLogger.debug.Output(2, fmt.Sprintf(msg, args...))
	}
}

// Fatalf 記錄致命錯誤並退出程式
func Fatalf(msg string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= logLevelFatal {
		defaultLogger.fatal.Output(2, fmt.Sprintf(msg, args...))
	}
	os.Exit(1) // 致命錯誤後退出程式
}

// parseLogLevel 解析日誌級別字串
func parseLogLevel(level string) logLevel {
	switch strings.ToLower(level) {
	case "debug":
		return logLevelDebug
	case "info":
		return logLevelInfo
	case "error":
		return logLevelError
	case "fatal":
		return logLevelFatal
	default:
		return logLevelInfo // 默認 INFO 級別
	}
}
