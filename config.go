package main

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 配置结构体
type Config struct {
	Port int `mapstructure:"port"` // 修正标签
}

// 加载配置
func LoadConfig(logger *zap.Logger) Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("读取配置失败", zap.Error(err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.Fatal("解析配置失败", zap.Error(err))
	}
	return config
}
