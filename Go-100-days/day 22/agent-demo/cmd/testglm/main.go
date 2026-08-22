package main

import (
	"context"
	"fmt"
	"os"

	"agent-demo/llm"
)

func main() {
	m := llm.NewGLM(
		os.Getenv("LLM_API_KEY"),
		os.Getenv("LLM_API_URL"),
		os.Getenv("LLM_MODEL"),
	)
	out, err := m.Chat(context.Background(), []llm.Message{
		{Role: "user", Content: "你好，一句话回复"},
	}, nil)
	if err != nil {
		fmt.Println("出错:", err)
		os.Exit(1)
	}
	fmt.Println("■ glm 回复:", out.Content)
}