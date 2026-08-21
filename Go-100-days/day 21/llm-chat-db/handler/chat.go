package handler

import (
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"llm-chat-db/service"
)

// ChatService 接口定义在"使用者"（handler）一侧：
// handler 只关心"需要哪些能力"，不关心具体实现（是 GORM 还是别的）。
type ChatService interface {
	BuildMessages(sid, q string) []map[string]string
	SaveMessage(sid, role, content string)
	StreamChat(ctx context.Context, messages []map[string]string, fn func(string)) error
}

// 编译期检查：确保 *service.ChatService 确实实现了本接口。
// 如果哪天 service 改方法签名导致不匹配，这里会直接编译报错，而不是运行期才炸。
var _ ChatService = (*service.ChatService)(nil)

// ChatHandler 只依赖 ChatService 接口，字段类型是接口
type ChatHandler struct {
	Svc ChatService
}

func NewChatHandler(svc ChatService) *ChatHandler {
	return &ChatHandler{Svc: svc}
}

func (h *ChatHandler) Chat(c *gin.Context) {
	q := c.Query("q")
	sid := c.DefaultQuery("sid", "default")
	if q == "" {
		c.String(400, "缺少参数 q")
		return
	}
	ctx := c.Request.Context()
	c.Header("Content-Type", "text/event-stream")

	h.Svc.SaveMessage(sid, "user", q)
	messages := h.Svc.BuildMessages(sid, q)

	var full strings.Builder
	err := h.Svc.StreamChat(ctx, messages, func(text string) {
		select {
		case <-ctx.Done():
			return // 用户断开，停止转发
		default:
			c.Writer.WriteString("data: " + text + "\n\n")
			c.Writer.Flush()
			full.WriteString(text)
		}
	})
	if err != nil {
		if ctx.Err() != nil {
			log.Println("[服务端] 用户断开，停止转发")
			c.Writer.WriteString("data: [断开]\n\n")
		} else {
			log.Printf("[服务端] 调用出错: %v", err)
			c.Writer.WriteString("data: [错误] " + err.Error() + "\n\n")
		}
		c.Writer.Flush()
		return
	}
	h.Svc.SaveMessage(sid, "assistant", full.String())
	log.Println("[服务端] 流结束，已存库")
}