package logx

import (
	"context"
	"log"
)

// Info 信息日志
func Info(msg string, args ...interface{}) {
	log.Println("[INFO]", msg, args)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format+"\n", args...)
}

// Error 错误日志
func Error(msg string, args ...interface{}) {
	log.Println("[ERROR]", msg, args)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format+"\n", args...)
}

// ErrorfContext 带上下文的错误日志
func ErrorfContext(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[ERROR] "+format+"\n", args...)
}

// ErrorTrace 错误追踪日志
func ErrorTrace(msg string, args ...interface{}) {
	log.Println("[ERROR]", msg, args)
}

// Debug 调试日志
func Debug(msg string, args ...interface{}) {
	log.Println("[DEBUG]", msg, args)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format+"\n", args...)
}

// Warn 警告日志
func Warn(msg string, args ...interface{}) {
	log.Println("[WARN]", msg, args)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format+"\n", args...)
}
