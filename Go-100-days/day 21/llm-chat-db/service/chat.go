package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"llm-chat-db/models"
)

// ChatService 掌管"聊天"的业务规则
type ChatService struct {
	DB     *gorm.DB
	Redis  *redis.Client
	APIKey string
	APIURL string
	Model  string
}

// NewChatService 构造 ChatService（依赖注入：外部把 db、redis、key 等传进来）
func NewChatService(g *gorm.DB, rdb *redis.Client, apiKey, apiURL, model string) *ChatService {
	return &ChatService{DB: g, Redis: rdb, APIKey: apiKey, APIURL: apiURL, Model: model}
}

// cacheKey 每个会话一个缓存 key，命名规范：对象:标识
func cacheKey(sid string) string {
	return "history:" + sid
}

// BuildMessages 查历史 + 拼当前问题（多轮记忆）
func (s *ChatService) BuildMessages(sid, q string) []map[string]string {
	messages := s.loadHistory(sid) // 从缓存或库拿历史
	messages = append(messages, map[string]string{"role": "user", "content": q})
	return messages
}

// loadHistory cache-aside：先查缓存 → miss 查库并回填；缓存出错则兜底查库
func (s *ChatService) loadHistory(sid string) []map[string]string {
	key := cacheKey(sid)
	cached, err := s.Redis.Get(context.Background(), key).Result()
	if err == redis.Nil { // 缓存未命中
		msgs := s.queryHistory(sid)
		if b, e := json.Marshal(msgs); e == nil {
			s.Redis.Set(context.Background(), key, b, 30*time.Second) // 回填 + 过期
		}
		fmt.Println("[缓存] miss，查库并回填", len(msgs), "条")
		return msgs
	}
	if err != nil { // 缓存本身出错，别拖垮请求，兜底查库
		log.Printf("读取缓存失败(%s): %v", key, err)
		return s.queryHistory(sid)
	}
	// 缓存命中
	var msgs []map[string]string
	if json.Unmarshal([]byte(cached), &msgs) != nil {
		return s.queryHistory(sid)
	}
	fmt.Println("[缓存] 命中", len(msgs), "条，没查库")
	return msgs
}

// queryHistory 从 MySQL 查该会话最近 20 条历史（原始实现）
func (s *ChatService) queryHistory(sid string) []map[string]string {
	var history []models.Message
	s.DB.Where("session_id = ?", sid).Order("created_at asc").Find(&history)
	start := 0
	if len(history) > 20 {
		start = len(history) - 20
	}
	msgs := make([]map[string]string, 0, len(history)-start)
	for _, m := range history[start:] {
		msgs = append(msgs, map[string]string{"role": m.Role, "content": m.Content})
	}
	fmt.Println("[DB] 查库，共", len(msgs), "条")
	return msgs
}

// SaveMessage 存一条消息；写入后失效缓存，保证下次不会读到旧数据
func (s *ChatService) SaveMessage(sid, role, content string) {
	s.DB.Create(&models.Message{SessionID: sid, Role: role, Content: content})
	s.Redis.Del(context.Background(), cacheKey(sid))
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
