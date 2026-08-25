package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"agent-demo/agent"
	"agent-demo/auth"
)

// 演示用固定账号（教学项目不做用户库，真实项目要接数据库鉴权）
const (
	demoUser     = "admin"
	demoPassword = "123456"
)

// Handler 持有 Agent，把 Agent 能力封装成对外的 HTTP 接口
type Handler struct {
	Agent     *agent.Agent
	JWTSecret string
	TokenTTL  time.Duration
}

func New(a *agent.Agent, secret string, ttl time.Duration) *Handler {
	return &Handler{Agent: a, JWTSecret: secret, TokenTTL: ttl}
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login 校验账号密码，通过后签发 JWT 返回给客户端
func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误"})
		return
	}
	if req.Username != demoUser || req.Password != demoPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
		return
	}
	token, err := auth.Generate(h.JWTSecret, req.Username, h.TokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签发 token 失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "expires_in": int(h.TokenTTL.Seconds())})
}

type chatReq struct {
	Query     string `json:"query"`
	SessionID string `json:"session_id"`
}

// Chat 受 JWT 保护：用 SSE 把 Agent 的 thought/tool_call/tool_result/answer 事件边算边推给客户端（B1）
func (h *Handler) Chat(c *gin.Context) {
	var req chatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误"})
		return
	}
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query 不能为空"})
		return
	}
	sid := req.SessionID
	if sid == "" {
		sid = "default"
	}

	// SSE 需要这行头，否则浏览器/curl 不会一直挂着读
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Writer.Flush()

	ctx := c.Request.Context()
	err := h.Agent.RunStream(ctx, sid, req.Query, func(ev agent.Event) {
		// 客户端断开后别再往响应里写（写了也没人收，只会报错）
		if ctx.Err() != nil {
			return
		}
		line, _ := json.Marshal(ev)
		c.Writer.WriteString("data: " + string(line) + "\n\n")
		c.Writer.Flush() // 立刻刷给客户端，不回等攒够一份
	})

	// 收尾透出一个 done 事件，让前端知道"流结束了"；有错则带上 error
	if ctx.Err() == nil {
		end := map[string]any{"type": "done", "error": ""}
		if err != nil {
			end["error"] = err.Error()
		}
		line, _ := json.Marshal(end)
		c.Writer.WriteString("data: " + string(line) + "\n\n")
		c.Writer.Flush()
	}
}