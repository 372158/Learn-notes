package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseArgs 结构化解析工具参数（B2）：把模型给的参数 JSON 字符串强类型解析进结构体，
// 并给出可读报错。相比裸 json.Unmarshal，它对"参数为空/不是合法 JSON"给更友好的错误。
func ParseArgs(argsJSON string, into any) error {
	if strings.TrimSpace(argsJSON) == "" {
		return fmt.Errorf("参数为空：模型没传任何参数")
	}
	if err := json.Unmarshal([]byte(argsJSON), into); err != nil {
		return fmt.Errorf("参数不是合法 JSON: %v（原始=%s）", err, argsJSON)
	}
	return nil
}