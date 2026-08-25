package agent

import (
	"encoding/json"
	"sync"

	"agent-demo/llm"
)

// Memory 会话记忆接口：不同实现（内存/Redis/MySQL）替换用
type Memory interface {
	Load(sessionID string) []llm.Message           // 取出某个会话的历史
	Append(sessionID string, msgs []llm.Message)  // 追加一段到该会话
}

// MemMemory 进程内内存实现（练手够用；生产换 Redis/MySQL）
type MemMemory struct {
	mu  sync.Mutex
	dat map[string][]llm.Message
}

func NewMemMemory() *MemMemory {
	return &MemMemory{dat: make(map[string][]llm.Message)}
}

// Load 返回会话历史的深拷贝，避免调用方意外污染内部数据
func (m *MemMemory) Load(sid string) []llm.Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 深拷贝：把结构体通过 JSON 序列化还原（简单可靠的克隆手段）
	b, _ := json.Marshal(m.dat[sid])
	var cp []llm.Message
	json.Unmarshal(b, &cp)
	return cp
}

// Append 追加若干消息到会话
func (m *MemMemory) Append(sid string, msgs []llm.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dat[sid] = append(m.dat[sid], msgs...)
}