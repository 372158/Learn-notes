// pkg/jwt/jwt.go
// JWT 工具包，用于生成和解析 token
// 学习要点：JWT（JSON Web Token）是无状态的认证方案
//
// 为什么用 JWT 而非 Session：
// 1. 无状态：服务端不存储，token 自带用户信息
// 2. 跨域友好：移动端、前后端分离场景适用
// 3. 易扩展：分布式系统不用共享 session 存储
//
// JWT 结构：header.payload.signature（用 . 分隔的三段）
// - header: 算法类型
// - payload: 存放用户信息（如 user_id, username）
// - signature: 签名，防止篡改

package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"simple-ecommerce/pkg/config"
)

// CustomClaims 自定义 Claims（载荷）
// 学习要点：嵌入 jwt.RegisteredClaims 获得标准字段（过期时间等）
// 再加自己的业务字段（用户 ID、用户名）
type CustomClaims struct {
	UserID   uint   `json:"user_id"`   // 用户 ID
	Username string `json:"username"`  // 用户名
	jwt.RegisteredClaims               // 嵌入标准字段（ExpiresAt 等）
}

// GenerateToken 生成 JWT token
// 学习要点：JWT 生成流程
// 1. 创建 Claims（载荷），填入用户信息和过期时间
// 2. 用密钥签名生成 token 字符串
// 3. 返回给客户端保存
func GenerateToken(userID uint, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 学习要点：过期时间 = 当前时间 + 配置的小时数
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GlobalConfig.JWT.Expire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			Subject:   username,                        // 主题（用户名）
		},
	}

	// 学习要点：jwt.NewWithClaims 创建 token
	// - jwt.SigningMethodHS256: 使用 HMAC-SHA256 算法
	// - claims: 载荷
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 用密钥签名，得到完整的 token 字符串
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

// ParseToken 解析 JWT token
// 学习要点：JWT 解析流程
// 1. 用密钥验证签名是否被篡改
// 2. 检查是否过期
// 3. 提取载荷中的用户信息
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 学习要点：jwt.ParseWithClaims 解析 token
	// 第三个参数是密钥验证函数，返回密钥用于验签
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否正确（防止算法攻击）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名算法错误")
		}
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// 类型断言，提取自定义 Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token 无效")
}
