package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

// Event 一次 Agent 运行过程中的"事件"，用于流式透出给客户端。
type Event struct {
	Type      string  `json:"type"`                // thought | tool_call | tool_result | answer | error
	Step      int     `json:"step"`                // 第几步（从 1 开始）
	ToolName  string  `json:"tool_name,omitempty"` // tool_call / tool_result 用：工具名
	Arguments string  `json:"arguments,omitempty"` // tool_call 用：模型给的工具参数(JSON)
	Result    string  `json:"result,omitempty"`    // tool_result 用：工具执行结果（或兜底时的错误）
	Content   string  `json:"content,omitempty"`   // thought / answer 用：说明文字或最终回答
	Answer    *Answer `json:"structured,omitempty"`// answer 用：结构化输出（B2）
	Note      string  `json:"note,omitempty"`      // 兜底说明等
	Error     string  `json:"error,omitempty"`     // error 用：错误信息
}

// Answer 结构化的最终输出（B2）：summary 结论 / detail 说明 / confidence 置信度(0~1)
type Answer struct {
	Summary    string  `json:"summary"`
	Detail     string  `json:"detail"`
	Confidence float64 `json:"confidence"`
}

// parseAnswer 把模型给的最终回答解析成结构化 Answer（B2）。
// 模型爱把 JSON 包在 ```json ... ``` 里或带前后缀文字，这里先截取首尾 {…} 再解析，能兜住围栏/杂音。
func parseAnswer(raw string) (*Answer, string) {
	s := strings.TrimSpace(raw)
	if i := strings.IndexByte(s, '{'); i >= 0 {
		if j := strings.LastIndexByte(s, '}'); j > i {
			s = s[i : j+1]
		}
	}
	var ans Answer
	if err := json.Unmarshal([]byte(s), &ans); err != nil {
		return nil, "非结构化回滚：模型未返回合法 JSON"
	}
	if ans.Summary == "" {
		return nil, "非结构化回滚：JSON 缺 summary 字段"
	}
	return &ans, ""
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

// RunStream 流式版：每个事件通过 emit 回调透出；最终答案以 Type=="answer" 给出。
func (a *Agent) RunStream(ctx context.Context, sessionID, userInput string, emit func(Event)) error {
	return a.run(ctx, sessionID, userInput, emit)
}

// run 核心逻辑：跑 ReAct 循环，并把关键节点以事件回调透出。
func (a *Agent) run(ctx context.Context, sessionID, userInput string, emit func(Event)) error {
	// 1. 加载该会话历史 + 塞入本次问题
	messages := a.Memory.Load(sessionID)
	if len(messages) == 0 {
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: "你是一个能调用工具来解决问题的助手。需要工具时就调用，不要编造结果。给出最终回答时，请严格按 JSON 输出：{\"summary\":\"一句话结论\",\"detail\":\"详细说明\",\"confidence\":0到1}。",
		})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userInput})
	startIdx := len(messages) // 记录本轮起始，便于结束时只追加"这轮新增的"

	for step := 0; step < a.MaxSteps; step++ {
		if ctx.Err() != nil {
			return ctx.Err() // 客户端断开，提前结束
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
					// "工具不存在"是程序 bug，属于不该发生的错，直接停
					return fmt.Errorf("工具不存在: %s", tc.Name)
				}

				emit(Event{Type: "thought", Step: step + 1, Content: "模型决定调用工具 " + tc.Name})
				emit(Event{Type: "tool_call", Step: step + 1, ToolName: tc.Name, Arguments: tc.Arguments})

				result, callErr := tool.Call(tc.Arguments)
				resultText := result
				if callErr != nil {
					// 错误兜底：工具执行失败【不中断 Agent】，
					// 把错误文本当 tool 结果喂回模型，让模型纠正后重试。
					resultText = "工具执行失败: " + callErr.Error()
				}
				emit(Event{Type: "tool_result", Step: step + 1, ToolName: tc.Name, Result: resultText})

				batch = append(batch,
					llm.Message{
						Role:               "assistant",
						Content:            "",
						AssistantToolCalls: []tools.ToolCall{tc},
					},
					llm.Message{Role: "tool", Content: resultText, ToolCallID: tc.ID},
				)
			}
			messages = append(messages, batch...)
			continue
		}

		if comp.Content != "" {
			messages = append(messages, llm.Message{Role: "assistant", Content: comp.Content})
			a.Memory.Append(sessionID, messages[startIdx:])
			ev := Event{Type: "answer", Step: step + 1, Content: comp.Content}
			// 结构化输出（B2）：尝试把最终回答按 JSON 解析进 Answer；
			// 解析失败或缺 summary 就兜底原样透出，绝不崩。
			if ans, note := parseAnswer(comp.Content); ans != nil {
				ev.Answer = ans
			} else {
				ev.Note = note
			}
			emit(ev)
			return nil
		}
		return fmt.Errorf("模型未返回内容")
	}
	return fmt.Errorf("达到最大步数 %d，Agent 未收敛", a.MaxSteps)
}