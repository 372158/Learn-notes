package config

import (
	"github.com/spf13/viper"
)

// Config 总配置结构体
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Shortlink ShortlinkConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"` // mapstructure 是 viper 用来解析配置的标签
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// ShortlinkConfig 短链接业务配置
type ShortlinkConfig struct {
	Domain string `mapstructure:"domain"`
	Length int    `mapstructure:"length"`
}

var AppConfig Config

// Load 加载配置文件
func Load() error {
	// 告诉 viper 去哪里找配置文件
	viper.SetConfigName("config")   // 文件名（不含扩展名）
	viper.SetConfigType("yaml")     // 文件类型
	viper.AddConfigPath(".")        // 在当前目录查找
	viper.AddConfigPath("./config") // 也可以在 ./config 目录查找

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 将配置映射到 AppConfig 结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	// 支持环境变量覆盖（方便 Docker 部署）
	viper.AutomaticEnv()

	return nil
}