package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"agent-demo/tools"
)

// glmMessage 传给 glm 的单条消息结构
type glmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type glmChatRequest struct {
	Model   string       `json:"model"`
	Message []glmMessage `json:"messages"`
}

type glmChatResponse struct {
	Choices []struct {
		Message glmMessage `json:"message"`
	} `json:"choices"`
}

// GLM 用 glm-4-flash 实现 Model 接口
type GLM struct {
	APIKey string
	APIURL string
	Model  string
}

func NewGLM(apiKey, apiURL, model string) *GLM {
	return &GLM{APIKey: apiKey, APIURL: apiURL, Model: model}
}

// Chat 实现 llm.Model 接口（只处理文字回复，工具调用留 A4）
func (g *GLM) Chat(ctx context.Context, message []Message, _ []tools.ToolSchema) (*Completion, error) {
	// 空A：把 []llm.Message 逐条转成 []glmMessage
	msgs := make([]glmMessage, 0, len(message))
	for _, m := range message {
		msgs = append(msgs, glmMessage{Role: m.Role, Content: m.Content})
	}

	body, _ := json.Marshal(glmChatRequest{
		Model:   g.Model,
		Message: msgs,
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

	// 空C：解码响应，从 choices[0].message.content 取人话
	var out glmChatResponse
	json.NewDecoder(resp.Body).Decode(&out)
	return &Completion{Content: out.Choices[0].Message.Content}, nil
}
