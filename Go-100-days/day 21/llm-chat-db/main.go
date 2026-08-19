package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Message 对应数据库表 messages
type Message struct {
	ID        uint   `gorm:"primaryKey"`
	SessionID string `gorm:"index"`
	Role      string
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
}

// 上游每个 chunk 的结构
type llmChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

var db *gorm.DB

func main() {
	apiKey := strings.TrimSpace(os.Getenv("LLM_API_KEY"))
	if apiKey == "" {
		panic("请先设置环境变量 LLM_API_KEY")
	}
	const url = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	const model = "glm-4-flash"

	var err error
	db, err = gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/chat_app?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}
	db.AutoMigrate(&Message{})

	r := gin.Default()

	r.GET("/chat", func(c *gin.Context) {
		q := c.Query("q")
		sid := c.DefaultQuery("sid", "default")
		if q == "" {
			c.String(400, "缺少参数 q")
			return
		}
		ctx := c.Request.Context()
		c.Header("Content-Type", "text/event-stream")

		// ① 先查该会话的历史消息（时间正序）
		var history []Message
		db.Where("session_id = ?", sid).Order("created_at asc").Find(&history)

		// ② 拼 messages：历史（最近 20 条）+ 当前问题
		messages := make([]map[string]string, 0, len(history)+1)
		start := 0
		if len(history) > 20 {
			start = len(history) - 20
		}
		for _, m := range history[start:] {
			messages = append(messages, map[string]string{"role": m.Role, "content": m.Content})
		}
		messages = append(messages, map[string]string{"role": "user", "content": q})

		// ③ 存用户消息
		db.Create(&Message{SessionID: sid, Role: "user", Content: q})

		reqBody, _ := json.Marshal(map[string]any{
			"model":    model,
			"messages": messages,
			"stream":   true,
		})

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
		if err != nil {
			c.String(500, "构建请求失败: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(500, "调用失败: %v", err)
			return
		}
		defer resp.Body.Close()

		ch := make(chan string, 16)

		go func() {
			defer close(ch)
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if !strings.HasPrefix(line, "data: ") {
					continue
				}
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					return
				}
				var ck llmChunk
				json.Unmarshal([]byte(data), &ck)
				if len(ck.Choices) > 0 && ck.Choices[0].Delta.Content != "" {
					ch <- ck.Choices[0].Delta.Content
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("[服务端] 读取上游流出错:", err)
			}
		}()

		var full strings.Builder
		for text := range ch {
			select {
			case <-ctx.Done():
				fmt.Println("[服务端] 用户断开，停止转发")
				return
			default:
				c.Writer.WriteString("data: " + text + "\n\n")
				c.Writer.Flush()
				full.WriteString(text)
			}
		}
		// ④ 存助手回复
		db.Create(&Message{SessionID: sid, Role: "assistant", Content: full.String()})
		fmt.Println("[服务端] 流结束，已存库")
	})

	r.Run(":8080")
}