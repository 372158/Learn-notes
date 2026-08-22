package config

import "os"

// Config 集中管理所有配置项，从环境变量读取
type Config struct {
	APIKey    string // LLM API key
	APIURL    string // LLM 接口地址
	Model     string // LLM 模型名
	Port      string // HTTP 端口
	JWTSecret string // JWT 签名密钥
}

// Load 从环境变量构建配置
func Load() *Config {
	return &Config{
		APIKey:    os.Getenv("LLM_API_KEY"),
		APIURL:    os.Getenv("LLM_API_URL"),
		Model:     os.Getenv("LLM_MODEL"),
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
