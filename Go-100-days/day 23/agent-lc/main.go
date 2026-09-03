package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

func main() {
	ctx := context.Background()

	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	toolList := []tools.Tool{tools.Calculator{}, tools.Weather{}}

	agent, err := agents.NewOpenAIFunctionsAgent(llm, toolList, agents.WithMaxIterations(5))
	if err != nil {
		log.Fatal(err)
	}
	executor := agents.NewExecutor(agent)

	ans, err := executor.Run(ctx, "帮我算一下 2+3*4 等于多少")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ans)
}