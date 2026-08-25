package main

import (
	"context"
	"fmt"
	"os"

	"agent-demo/llm"
	"agent-demo/tools"
)

func main() {
	m := llm.NewGLM(
		os.Getenv("LLM_API_KEY"),
		os.Getenv("LLM_API_URL"),
		os.Getenv("LLM_MODEL"),
	)

	// 构造一个"查天气"工具的说明书（发给模型看）
	weather := tools.ToolSchema{
		Type: "function",
		Function: tools.FunctionSchema{
			Name:        "get_weather",
			Description: "查询指定城市的天气",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"city": map[string]any{"type": "string"},
				},
				"required": []string{"city"},
			},
		},
	}

	out, err := m.Chat(context.Background(), []llm.Message{
		{Role: "user", Content: "请问北京今天天气怎么样？"},
	}, []tools.ToolSchema{weather})
	if err != nil {
		fmt.Println("出错:", err)
		os.Exit(1)
	}

	fmt.Printf("返回 %d 个工具调用：\n", len(out.ToolCalls))
	for _, tc := range out.ToolCalls {
		fmt.Printf("  工具名=%s 参数=%s\n", tc.Name, tc.Arguments)
	}
	if len(out.ToolCalls) == 0 {
		fmt.Println("■ 文字回复:", out.Content)
	}
}