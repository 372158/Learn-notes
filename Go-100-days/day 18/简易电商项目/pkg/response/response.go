// pkg/response/response.go
// 统一响应格式工具包
// 学习要点：统一 API 响应格式，让前端处理更一致
//
// 为什么统一响应格式：
// 1. 前端不用每次判断不同接口的不同格式
// 2. 统一错误码便于排查问题
// 3. 代码更简洁，controller 中一行就能返回响应
//
// 标准格式：
// {
//   "code": 200,          // 业务状态码（200成功，非200失败）
//   "msg": "success",     // 提示信息
//   "data": {...}         // 数据，失败时为 null
// }

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"` // 业务状态码
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据（用 interface{} 可放任意类型）
}

// 定义常用业务码
// 学习要点：用常量定义魔法数字，代码可读性更好
const (
	CodeSuccess = 200 // 成功
	CodeError   = 500 // 失败
	CodeUnauthorized = 401 // 未授权（token 无效或过期）
)

// Success 成功响应
// 学习要点：封装快捷函数，controller 中调用 response.Success(c, data) 即可
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义提示信息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Error 失败响应
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: CodeError,
		Msg:  msg,
		Data: nil,
	})
}

// ErrorWithCode 失败响应（自定义业务码）
func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
