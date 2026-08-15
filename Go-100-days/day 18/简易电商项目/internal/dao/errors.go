// internal/dao/errors.go
// DAO 层公共错误定义
// 学习要点：用哨兵错误（sentinel error）统一错误标识，便于上层判断
//
// 为什么用 errors.New 定义全局错误变量：
// 1. 上层可用 errors.Is() 精确判断错误类型
// 2. 避免用字符串比较判断错误（易出错）
// 3. 错误信息集中管理

package dao

import "errors"

// ErrInsufficientStock 库存不足错误
var ErrInsufficientStock = errors.New("库存不足")
