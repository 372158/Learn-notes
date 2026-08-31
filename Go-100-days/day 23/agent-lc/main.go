package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	ctx := context.Background()

	llm, err := openai.New()

	if err != nil {
		log.Fatal(err)
	}

	// 1. 拼消息：system 明确要求 JSON 格式输出 objective
	msgs := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "你是翻译员。输出严格的 JSON:{\"translation\":\"译文\",\"confidence\":0到1}，不要输出其他文字。"),
		llms.TextParts(llms.ChatMessageTypeHuman, "The quick brown fox jumps over the lazy dog."),
	}

	//2. 带 JSON 模式调用
	resp, err := llm.GenerateContent(ctx, msgs, llms.WithJSONMode())
	if err != nil {
		log.Fatal(err)
	}
	raw := resp.Choices[0].Content
	fmt.Println("模型原始输出：", raw)

	// 3. 自己解析 + 兜底（这是 B2 的 parseAnswer 精神）
	var out struct {
		Translation string  `json:"translation"`
		Confidence  float64 `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		fmt.Println("⚠️ 解析失败，可能需要剥围栏:", err)
		return
	}
	fmt.Printf("译文=%s 置信度=%.2f\n", out.Translation, out.Confidence)
}
