package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// 上游每个 chunk 的结构（只声明关心的字段，json 包自动匹配）
type llmChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func main() {
	apiKey := strings.TrimSpace(os.Getenv("LLM_API_KEY"))
	if apiKey == "" {
		panic("请先设置环境变量 LLM_API_KEY")
	}
	const url = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	const model = "glm-4-flash"

	r := gin.Default()

	r.GET("/chat", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.String(400, "缺少参数 q")
			return
		}
		ctx := c.Request.Context()
		c.Header("Content-Type", "text/event-stream")

		reqBody, _ := json.Marshal(map[string]any{
			"model": model,
			"messages": []map[string]string{
				{"role": "user", "content": q},
			},
			"stream": true,
		})

		// 空1：把请求的 ctx 绑到上游 HTTP 请求 → 用户断开时 API 调用自动取消
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
			defer close(ch) // 空2：goroutine 任何路径退出都关闭 → range 才能结束
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

		for text := range ch {
			select {
			case <-ctx.Done(): // 空3：用户断开 → 不再往死连接写
				fmt.Println("[服务端] 用户断开，停止转发")
				return
			default:
				c.Writer.WriteString("data: " + text + "\n\n")
				c.Writer.Flush()
			}
		}
		fmt.Println("[服务端] 流结束")
	})

	r.Run(":8080")
}
