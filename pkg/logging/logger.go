package logging

import (
	"log"
	"os"
)

// Logger 定义日志记录器结构体
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
	file        *os.File // 保存日志文件句柄
}

// NewLogger 创建并返回一个新的 Logger 实例
func NewLogger() *Logger {
	// 打开日志文件（如不存在则创建）
	file, err := os.OpenFile("certwatch.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return &Logger{
		infoLogger:  log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		file:        file,
	}
}

// Info 记录信息日志
func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
}

// Error 记录错误日志
func (l *Logger) Error(message string) {
	l.errorLogger.Println(message)
}

// Warn 记录警告日志
func (l *Logger) Warn(message string) {
	l.warnLogger.Println(message)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
	}
}
