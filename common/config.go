package common

import (
	"os"
	"strconv"
)

// Config 配置结构
type Config struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Database string `json:"database"`
	Debug    bool   `json:"debug"`
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	config := &Config{
		Port:     8080,
		Host:     "localhost",
		Database: "default",
		Debug:    false,
	}

	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}

	if host := os.Getenv("HOST"); host != "" {
		config.Host = host
	}

	if database := os.Getenv("DATABASE"); database != "" {
		config.Database = database
	}

	if debug := os.Getenv("DEBUG"); debug != "" {
		if d, err := strconv.ParseBool(debug); err == nil {
			config.Debug = d
		}
	}

	return config
}
