package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := InitLogger()
	defer logger.Sync()

	config := LoadConfig(logger)

	server := NewServer(config, logger)
	if err := server.Start(); err != nil {
		logger.Fatal("服务器启动失败", zap.Error(err))
	}
}
