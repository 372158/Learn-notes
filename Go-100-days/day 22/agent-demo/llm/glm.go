package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"agent-demo/tools"
)

// glmMessage 传给 glm 的单条消息结构
type glmMessage struct {
	Role       string        `json:"role"`
	Content    string        `json:"content"`
	ToolCalls  []glmToolCall `json:"tool_calls,omitempty"`
	ToolCallID string        `json:"tool_call_id,omitempty"`
}

type glmChatRequest struct {
	Model   string             `json:"model"`
	Message []glmMessage       `json:"messages"`
	Tools   []tools.ToolSchema `json:"tools,omitempty"`
}

type glmChatResponse struct {
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Choices []struct {
		Message glmMessage `json:"message"`
	} `json:"choices"`
}

// glmToolCall 模型返回/需要回传的一个工具调用
type glmToolCall struct {
	Type     string `json:"type"` // 固定 "function"
	ID       string `json:"id"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type GLM struct {
	APIKey string
	APIURL string
	Model  string
}

func NewGLM(apiKey, apiURL, model string) *GLM {
	return &GLM{APIKey: apiKey, APIURL: apiURL, Model: model}
}

// Chat 实现 llm.Model 接口（含工具调用）
func (g *GLM) Chat(ctx context.Context, message []Message, toolsSchemas []tools.ToolSchema) (*Completion, error) {
	msgs := make([]glmMessage, 0, len(message))
	for _, m := range message {
		gm := glmMessage{Role: m.Role, Content: m.Content, ToolCallID: m.ToolCallID}
		if len(m.AssistantToolCalls) > 0 {
			for _, tc := range m.AssistantToolCalls {
				gm.ToolCalls = append(gm.ToolCalls, glmToolCall{
					Type: "function",
					ID:   tc.ID,
					Function: struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					}{Name: tc.Name, Arguments: tc.Arguments},
				})
			}
		}
		msgs = append(msgs, gm)
	}

	body, _ := json.Marshal(glmChatRequest{
		Model:   g.Model,
		Message: msgs,
		Tools:   toolsSchemas,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", g.APIURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("构建请求：%w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求： %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var out glmChatResponse
	if jerr := json.Unmarshal(raw, &out); jerr != nil {
		return nil, fmt.Errorf("响应解析失败: %v，原文=%s", jerr, string(raw))
	}
	if out.Error != nil {
		return nil, fmt.Errorf("glm 返回错误: %s", out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return nil, fmt.Errorf("glm 未返回 choices，原文=%s", string(raw))
	}

	msg := out.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		tcs := make([]tools.ToolCall, 0, len(msg.ToolCalls))
		for _, tc := range msg.ToolCalls {
			tcs = append(tcs, tools.ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			})
		}
		return &Completion{ToolCalls: tcs}, nil
	}
	return &Completion{Content: msg.Content}, nil
}