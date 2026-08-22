package llm

import (
	"context"

	"agent-demo/tools"
)

// Message 模型对话里的一则消息。role: system/user/assistant/tool
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
	// ToolCallID 当这则消息是对"某次工具调用结果"的回应时用到(assistant->tool 回执)
	ToolCallID string `json:"tool_call_id,omitempty"`
}

// ToolCallResult 一次工具调用的输出(执行完工具得到的)，要喂回给模型
type ToolCallResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"` // 工具执行返回的结果字符串(JSON)
}

// Completion 模型一次返回的结构：文字 + 是否想调工具
type Completion struct {
	Content   string           // 模型说的人话(若没调工具)
	ToolCalls []tools.ToolCall // 模型想要的工具调用(若想调工具)
}

// Model 所有大模型的统一"合同"。A3 用 glm-4-flash 实现它，A4 加上工具调用。
type Model interface {
	// Chat 发一轮对话，传历史+可选的工具，返回文字回复或工具调用请求。
	// 具体怎么接、怎么解析，A3/A4 实现。
	Chat(ctx context.Context, messages []Message, toolSchemas []tools.ToolSchema) (*Completion, error)
}