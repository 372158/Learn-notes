package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
	"llm-chat-db/models"
)

// ChatService 掌管"聊天"的业务规则
type ChatService struct {
	DB     *gorm.DB
	APIKey string
	APIURL string
	Model  string
}

// NewChatService 构造 ChatService（依赖注入：外部把 db、key 等传进来）
func NewChatService(g *gorm.DB, apiKey, apiURL, model string) *ChatService {
	return &ChatService{DB: g, APIKey: apiKey, APIURL: apiURL, Model: model}
}

// BuildMessages 查历史 + 拼当前问题（多轮记忆）
func (s *ChatService) BuildMessages(sid, q string) []map[string]string {
	var history []models.Message
	s.DB.Where("session_id = ?", sid).Order("created_at asc").Find(&history)

	messages := make([]map[string]string, 0, len(history)+1)
	start := 0
	if len(history) > 20 {
		start = len(history) - 20
	}
	for _, m := range history[start:] {
		messages = append(messages, map[string]string{"role": m.Role, "content": m.Content})
	}
	messages = append(messages, map[string]string{"role": "user", "content": q})
	return messages
}

// SaveMessage 存一条消息
func (s *ChatService) SaveMessage(sid, role, content string) {
	s.DB.Create(&models.Message{SessionID: sid, Role: role, Content: content})
}

// StreamChat 调用 LLM，逐 chunk 通过回调 fn 返回；返回错误则中止
func (s *ChatService) StreamChat(ctx context.Context, messages []map[string]string, fn func(string)) error {
	reqBody, _ := json.Marshal(map[string]any{
		"model":    s.Model,
		"messages": messages,
		"stream":   true,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", s.APIURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("构建请求: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求: %w", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return nil
		}
		var ck struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		json.Unmarshal([]byte(data), &ck)
		if len(ck.Choices) > 0 && ck.Choices[0].Delta.Content != "" {
			fn(ck.Choices[0].Delta.Content)
		}
	}
	return scanner.Err()
}