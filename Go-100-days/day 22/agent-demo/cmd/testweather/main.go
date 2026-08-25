package main

import (
	"fmt"

	"agent-demo/tools"
)

func main() {
	for _, city := range []string{"北京", "上海", "杭州"} {
		args := fmt.Sprintf(`{"city":%q}`, city)
		out, err := tools.Weather.Call(args)
		if err != nil {
			fmt.Printf("%-6s → 出错: %v\n", city, err)
			continue
		}
		fmt.Printf("%-6s → %s\n", city, out)
	}
}