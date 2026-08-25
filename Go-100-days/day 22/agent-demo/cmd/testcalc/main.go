package main

import (
	"fmt"

	"agent-demo/tools"
)

func main() {
	for _, expr := range []string{"1+2", "3*(4-1)", "10/4", "2+3*4"} {
		// 模拟模型传来的参数字符串：{"expr":"..."}
		args := fmt.Sprintf(`{"expr":%q}`, expr)
		out, err := tools.Calc.Call(args)
		if err != nil {
			fmt.Printf("%-12s → 出错: %v\n", expr, err)
			continue
		}
		fmt.Printf("%-12s → %s\n", expr, out)
	}
}