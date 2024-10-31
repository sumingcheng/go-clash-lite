package main

import (
	"go.uber.org/zap"
)

// 初始化日志
func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
