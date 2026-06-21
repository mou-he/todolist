// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
)

// 定义程序所需的所有配置字段
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
type LogConfig struct {
	Level       string `mapstructure:"level"`
	OutputPath string `mapstructure:"output_path"`
	Format     string `mapstructure:"format"`
}

// Load 是唯一暴露给外部的初始化函数
func Load() (*Config, error) {
	// 1. 告诉 Viper 去哪里找配置文件
	viper.SetConfigName("config")    // 文件名（无后缀）
	viper.SetConfigType("yaml")      // 文件类型
	viper.AddConfigPath("./configs") // 文件路径

	// 2. 支持环境变量覆盖（优先级高于配置文件）
	viper.AutomaticEnv() // 自动匹配环境变量（如 SERVER_PORT）

	// 3. 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 4. 将配置映射到结构体（利用 mapstructure 标签）
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
