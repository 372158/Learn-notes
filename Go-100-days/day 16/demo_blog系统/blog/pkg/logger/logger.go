package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

// Init 初始化日志
func Init() error {
	// 直接用 zap.NewProduction() 创建一个生产级别的 Logger
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	Log = logger
	return nil
}

// Sync 刷新日志缓冲
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
