package agent

import (
	"context"
	"fmt"

	"agent-demo/llm"
	"agent-demo/tools"
)

// Agent 协调者：持有 LLM、工具表和记忆，跑 ReAct 循环
type Agent struct {
	Model       llm.Model             // 用的 LLM
	ToolMap     map[string]tools.Tool // 工具名 -> 工具
	ToolSchemas []tools.ToolSchema    // 发给模型的工具"说明书"
	MaxSteps    int                   // 最多循环轮数，防死循环
	Memory      Memory                // 会话记忆（存历史，实现多轮）
}

// Event 一次 Agent 运行过程中的"事件"，用于流式透出给客户端（B1）。
type Event struct {
	Type      string `json:"type"`                 // thought | tool_call | tool_result | answer | error
	Step      int    `json:"step"`                 // 第几步（从 1 开始）
	ToolName  string `json:"tool_name,omitempty"`  // tool_call / tool_result 用：工具名
	Arguments string `json:"arguments,omitempty"`  // tool_call 用：模型给的工具参数(JSON)
	Result    string `json:"result,omitempty"`     // tool_result 用：工具执行结果
	Content   string `json:"content,omitempty"`    // thought / answer 用：说明文字或最终答案
	Error     string `json:"error,omitempty"`      // error 用：错误信息
}

// New 注入依赖，并为每个工具生成协议兼容的 function schema
func New(m llm.Model, ts []tools.Tool, mem Memory, maxSteps int) *Agent {
	schemas := make([]tools.ToolSchema, 0, len(ts))
	for _, t := range ts {
		schemas = append(schemas, toolSchemaFor(t))
	}
	return &Agent{
		Model:       m,
		ToolMap:     tools.Register(ts),
		ToolSchemas: schemas,
		MaxSteps:    maxSteps,
		Memory:      mem,
	}
}

// toolSchemaFor 按工具名给出参数说明书
func toolSchemaFor(t tools.Tool) tools.ToolSchema {
	td := tools.FunctionSchema{Name: t.Name, Description: t.Description}
	req := []string{}
	props := map[string]any{}
	switch t.Name {
	case "calculator":
		req = []string{"expr"}
		props = map[string]any{
			"expr": map[string]any{
				"type":        "string",
				"description": "要计算的数学表达式，如 1+2、3*(4-1)、10/4",
			},
		}
	case "get_weather":
		req = []string{"city"}
		props = map[string]any{
			"city": map[string]any{
				"type":        "string",
				"description": "要查询的城市名，如 北京、上海、广州",
			},
		}
	}
	return tools.ToolSchema{
		Type: "function",
		Function: tools.FunctionSchema{
			Name:        td.Name,
			Description: td.Description,
			Parameters: map[string]any{
				"type":       "object",
				"properties": props,
				"required":   req,
			},
		},
	}
}

// Run 兼容旧接口：跑完整轮只回最终答案文本（供 cmd 命令行程序用）
func (a *Agent) Run(ctx context.Context, sessionID, userInput string) (string, error) {
	var final string
	err := a.run(ctx, sessionID, userInput, func(ev Event) {
		if ev.Type == "answer" {
			final = ev.Content
		}
	})
	return final, err
}

// RunStream 流式版：途中每个关键节点都通过 emit 回调透出；
// 最终答案以 Type=="answer" 的事件给出，结束要么是 err!=nil 要么收到 answer。
func (a *Agent) RunStream(ctx context.Context, sessionID, userInput string, emit func(Event)) error {
	return a.run(ctx, sessionID, userInput, emit)
}

// run 核心逻辑：跑 ReAct 循环，并把每个关键节点以事件回调透出。
// 用回调而非返回切片，是想把"产出事件"和"怎么消费事件"解耦（B1 交给 SSE，命令行可无视）。
func (a *Agent) run(ctx context.Context, sessionID, userInput string, emit func(Event)) error {
	// 1. 加载该会话历史 + 塞入本次问题
	messages := a.Memory.Load(sessionID)
	if len(messages) == 0 {
		messages = append(messages, llm.Message{
			Role: "system", Content: "你是一个能调用工具来解决问题的助手。需要工具时就调用，不要编造结果。",
		})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userInput})
	// 记录本轮起始，便于结束只追加"这轮新增的消息"
	startIdx := len(messages)

	for step := 0; step < a.MaxSteps; step++ {
		// 客户端断开时 ctx 会被取消，及时停，别白调模型
		if ctx.Err() != nil {
			return ctx.Err()
		}

		comp, err := a.Model.Chat(ctx, messages, a.ToolSchemas)
		if err != nil {
			emit(Event{Type: "error", Step: step + 1, Error: err.Error()})
			return fmt.Errorf("第%d轮模型调用失败: %w", step+1, err)
		}

		if len(comp.ToolCalls) > 0 {
			var batch []llm.Message
			for _, tc := range comp.ToolCalls {
				tool, ok := a.ToolMap[tc.Name]
				if !ok {
					return fmt.Errorf("工具不存在: %s", tc.Name)
				}

				// 透出：模型这步"想"调哪个工具
				emit(Event{Type: "thought", Step: step + 1, Content: "模型决定调用工具 " + tc.Name})
				emit(Event{Type: "tool_call", Step: step + 1, ToolName: tc.Name, Arguments: tc.Arguments})

				result, err := tool.Call(tc.Arguments)
				if err != nil {
					return fmt.Errorf("工具 %s 执行失败: %w", tc.Name, err)
				}

				// 透出：工具执行的结果
				emit(Event{Type: "tool_result", Step: step + 1, ToolName: tc.Name, Result: result})

				batch = append(batch,
					llm.Message{
						Role:               "assistant",
						Content:            "",
						AssistantToolCalls: []tools.ToolCall{tc},
					},
					llm.Message{
						Role:       "tool",
						Content:    result,
						ToolCallID: tc.ID,
					},
				)
			}
			messages = append(messages, batch...)
			continue
		}

		if comp.Content != "" {
			// 最终回答也进历史，形成完整一问一答
			messages = append(messages, llm.Message{Role: "assistant", Content: comp.Content})
			a.Memory.Append(sessionID, messages[startIdx:])
			emit(Event{Type: "answer", Step: step + 1, Content: comp.Content})
			return nil
		}
		return fmt.Errorf("模型未返回内容")
	}
	return fmt.Errorf("达到最大步数 %d，Agent 未收敛", a.MaxSteps)
}