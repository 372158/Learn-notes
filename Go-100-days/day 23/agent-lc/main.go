package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	ctx := context.Background()

	// 代码里指定模型名（环境变量 OPENAI_MODEL 不设时的可靠做法）

	llm, err := openai.New(openai.WithModel("glm-4-flash"))
	if err != nil {
		log.Fatal(err)
	}

	// 构造消息：system + human 各一条
	msgs := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "你是一个简洁的助手，回答不要吵过两句话。"),
		llms.TextParts(llms.ChatMessageTypeHuman, "用一句话介绍 GO 语言"),
	}

	resp1, err := llm.GenerateContent(ctx, msgs)
	if err != nil {
		log.Fatal(err)
	}
	answer1 := resp1.Choices[0].Content
	fmt.Println("第一轮：", answer1)

	// 追加历史
	msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeAI, "你是一个简洁的助手，回答不要吵过两句话。"))
	msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeHuman, "“Go 和 Python 比有什么优势"))

	resp2, err := llm.GenerateContent(ctx, msgs)
	if err != nil {
		log.Fatal(err)
	}

	ansewer2 := resp2.Choices[0].Content
	fmt.Println("第二轮：", ansewer2)

}
