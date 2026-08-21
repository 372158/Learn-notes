package config

import "os"

// Config 集中管理所有配置项，从环境变量读取
type Config struct {
	DBDSN  string // 数据库连接串
	APIKey string // LLM API key
	APIURL string // LLM 接口地址
	Model  string // LLM 模型名
	Port   string // HTTP 端口
}

// Load 从环境变量构建配置（敏感信息走环境变量，不硬编码）
func Load() *Config {
	return &Config{
		DBDSN:  os.Getenv("DB_DSN"),
		APIKey: os.Getenv("LLM_API_KEY"),
		APIURL: os.Getenv("LLM_API_URL"),
		Model:  os.Getenv("LLM_MODEL"),
		Port:   os.Getenv("PORT"),
	}
}