package agent

import (
	"context"

	"agent-demo/llm"
	"agent-demo/tools"
)

// Agent 协调者：持有 LLM 和工具表，把 A1 的 ReAct 循环跑起来。
// 它依赖 llm.Model 接口和 tools.Tool，具体实现(glm/计算器)由外部注入 —— 面向接口编程。
type Agent struct {
	Model  llm.Model                     // 用的 LLM
	ToolMap map[string]tools.Tool        // 工具名 -> 工具
	ToolSchemas []tools.ToolSchema       // 发给模型的工具"说明书"列表
	MaxSteps int                         // 最多循环轮数，防止死循环
}

// New 注入依赖，构造 Agent
func New(m llm.Model, ts []tools.Tool, maxSteps int) *Agent {
	return &Agent{
		Model:       m,
		ToolMap:     tools.Register(ts),
		ToolSchemas: make([]tools.ToolSchema, 0),
		MaxSteps:    maxSteps,
	}
}

// Run 执行一次对话（真正的 ReAct 循环 A6 实现，这里只占位）
func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {
	// A6 TODO: 1)把历史+工具说明书发给 llm
	//         2)若返回 ToolCall → 查到工具并执行 → 把结果喂回
	//         3)重复直到 model 给出纯文字回答
	return "", context.Canceled // 占位，A6 替换
}
