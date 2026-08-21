package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// mockChatService 实现 ChatService 接口，但不连真实 DB / LLM。
// 测试通过控制 mock 的字段（chunks / streamErr）来模拟不同上游行为。
type mockChatService struct {
	chunks    []string // StreamChat 依次回调的内容
	streamErr error    // StreamChat 要返回的错误
	saved     []struct{ sid, role, content string } // 记录 SaveMessage 的调用
}

func (m *mockChatService) BuildMessages(sid, q string) []map[string]string {
	return nil
}

func (m *mockChatService) SaveMessage(sid, role, content string) {
	m.saved = append(m.saved, struct{ sid, role, content string }{sid, role, content})
}

func (m *mockChatService) StreamChat(ctx context.Context, messages []map[string]string, fn func(string)) error {
	if m.streamErr != nil {
		return m.streamErr
	}
	for _, ch := range m.chunks {
		fn(ch) // 模拟上游逐 chunk 产出
	}
	return nil
}

// 编译期检查：确保 mock 也满足接口
var _ ChatService = (*mockChatService)(nil)

func setupRouter(svc ChatService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/chat", NewChatHandler(svc).Chat)
	return r
}

func doChat(r *gin.Engine, query string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/chat?"+query, nil)
	r.ServeHTTP(w, req)
	return w
}

// 表驱动测试：正常流式的几个场景
func TestChatStreamSuccess(t *testing.T) {
	tests := []struct {
		name   string
		chunks []string
	}{
		{"单chunk", []string{"你"}},
		{"多chunk", []string{"你", "好", "世", "界"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockChatService{chunks: tt.chunks}
			r := setupRouter(m)

			w := doChat(r, "q=hello&sid=s1")
			if w.Code != http.StatusOK {
				t.Fatalf("期望 200, got %d", w.Code)
			}
			body := w.Body.String()
			for _, ch := range tt.chunks {
				if !strings.Contains(body, "data: "+ch) {
					t.Errorf("响应缺 chunk %q，实际: %q", ch, body)
				}
			}
			// 正常路径应保存 user + assistant 两条消息
			if len(m.saved) != 2 {
				t.Errorf("应保存 2 条消息, 实际 %d", len(m.saved))
			}
			if m.saved[0].role != "user" || m.saved[1].role != "assistant" {
				t.Errorf("消息顺序应为 user→assistant, 实际 %+v", m.saved)
			}
		})
	}
}

func TestChatMissingQ(t *testing.T) {
	m := &mockChatService{}
	r := setupRouter(m)

	w := doChat(r, "sid=s1") // 没传 q
	if w.Code != http.StatusBadRequest {
		t.Errorf("缺 q 应返回 400, got %d", w.Code)
	}
}

func TestChatStreamError(t *testing.T) {
	m := &mockChatService{streamErr: context.Canceled}
	r := setupRouter(m)

	w := doChat(r, "q=hello&sid=s1")
	body := w.Body.String()
	// handler 收到非 ctx 取消的错误时，应输出错误标记
	if !strings.Contains(body, "错误") {
		t.Errorf("流错误时响应应含错误标记, 实际: %q", body)
	}
}