// pkg/config/config.go
// 配置加载模块，使用 Viper 读取 config.yaml 并解析到结构体
// 学习要点：Viper 是 Go 中最流行的配置管理库，支持 yaml/json/toml 等多种格式
//
// 为什么用结构体保存配置：
// 1. 类型安全：编译期就能发现拼写错误，而非运行时
// 2. 代码提示：IDE 能自动补全字段名
// 3. 集中管理：所有配置项一目了然

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
// 学习要点：结构体嵌套，把不同模块的配置分组管理
// yaml tag 用于和 config.yaml 中的字段名对应
type Config struct {
	App   AppConfig   `yaml:"app"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
	JWT   JWTConfig   `yaml:"jwt"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// DSN 返回 MySQL 连接字符串
// 学习要点：GORM 连接 MySQL 需要特定格式的 DSN（Data Source Name）
// 格式：用户名:密码@tcp(地址:端口)/数据库名?参数
func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.DBName)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"` // 过期时间，单位小时
}

// GlobalConfig 全局配置变量
// 学习要点：用全局变量保存配置，整个项目都能访问，避免到处传参
var GlobalConfig *Config

// Load 加载配置文件
// 学习要点：Viper 的标准使用流程
// 1. SetConfigName 设置配置文件名（不含扩展名）
// 2. SetConfigType 设置配置文件类型
// 3. AddConfigPath 添加查找路径（可多个）
// 4. ReadInConfig 读取配置文件
// 5. Unmarshal 把配置反序列化到结构体
func Load(path string) error {
	viper.SetConfigName("config")  // 配置文件名（不含扩展名）
	viper.SetConfigType("yaml")    // 配置文件类型
	viper.AddConfigPath(path)      // 查找路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 把配置解析到结构体
	// 学习要点：viper.Unmarshal 会根据 yaml tag 自动匹配字段
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	GlobalConfig = &cfg
	return nil
}
