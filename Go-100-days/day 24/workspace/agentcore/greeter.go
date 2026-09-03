package agentcore

import "fmt"

func Greet(name string) string {
	return "哇哈哈哈, " + name + "（来自 agentcore 模块）"
}

var _ = fmt.Sprintf // 占位防止未用导入报错，可删