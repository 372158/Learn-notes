package main

import (
	"context"
	"fmt"
	"os"

	"agent-demo/agent"
	"agent-demo/llm"
	"agent-demo/tools"
)

func main() {
	m := llm.NewGLM(
		os.Getenv("LLM_API_KEY"),
		os.Getenv("LLM_API_URL"),
		os.Getenv("LLM_MODEL"),
	)
	mem := agent.NewMemMemory()
	ag := agent.New(m, []tools.Tool{tools.Calc, tools.Weather}, mem, 5)

	// 同一会话连续多轮，验证记忆
	sid := "memtest"
	questions := []string{
		"帮我算一下 2+3*4 等于多少",
		"我刚刚问的算术题的答案是多少？", // 依赖上轮记忆
		"北京今天天气怎么样？",
	}
	for _, q := range questions {
		fmt.Println("==== 问:", q)
		ans, err := ag.Run(context.Background(), sid, q)
		if err != nil {
			fmt.Println("Agent 出错:", err)
			continue
		}
		fmt.Println("▶ 回答:", ans)
	}
}