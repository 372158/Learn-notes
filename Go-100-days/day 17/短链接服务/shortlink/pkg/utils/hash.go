package utils

import (
	"crypto/md5"
	"math/big"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateShortCode 根据原始 URL 生成短码（base62 编码，字符空间 62）
func GenerateShortCode(url string, length int) string {
	// 用 MD5 生成 128 位哈希
	hash := md5.Sum([]byte(url))

	// 将哈希字节转为大整数，再转成 base62 字符串
	num := new(big.Int).SetBytes(hash[:])
	base := big.NewInt(62)
	remainder := new(big.Int)
	encoded := ""

	for num.Cmp(big.NewInt(0)) > 0 {
		num.DivMod(num, base, remainder)
		encoded = string(base62Chars[remainder.Int64()]) + encoded
	}

	// 长度不足则前补 '0'
	for len(encoded) < length {
		encoded = string(base62Chars[0]) + encoded
	}

	// 取前 length 位作为短码
	if length > len(encoded) {
		length = len(encoded)
	}
	return encoded[:length]
}
