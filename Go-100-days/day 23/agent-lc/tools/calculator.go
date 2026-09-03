package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Calculator 计算工具：struct 实现 Tool 接口
type Calculator struct{}

func (c Calculator) Name() string {
	return "calculator"
}

func (c Calculator) Description() string {
	// 描述要写清楚：何时用 + 怎么传参（模型靠它决定是否调用、传什么）
	return `计算一个数学表达式（只支持 + - * / 和括号与数字），如 12+34、3*(4-1)、10/4。输入格式: {"expr": "表达式字符串"}`
}

func (c Calculator) Call(_ context.Context, input string) (string, error) {
	var args struct {
		Expr string `json:"expr"`
	}
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}
	if strings.TrimSpace(args.Expr) == "" {
		return "", fmt.Errorf("表达式为空")
	}
	val, err := evalSimple(args.Expr) // 你 day 22 的 evalSimple 可以搬过来
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(val, 'f', -1, 64), nil
}