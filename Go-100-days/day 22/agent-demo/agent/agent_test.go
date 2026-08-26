package agent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"agent-demo/llm"
	"agent-demo/tools"
)

// ---- 剧本假模型：不发网络请求，按预设顺序吐出结果 ----

// mockLLM scripted 假 LLM：实现 llm.Model，按 responses 顺序返回，并记录每次入参。
type mockLLM struct {
	responses []*llm.Completion
	idx       int
	allCalls  [][]llm.Message // 每次 Chat 收到的 messages（用于断言）
	err       error
}

func (m *mockLLM) Chat(_ context.Context, msgs []llm.Message, _ []tools.ToolSchema) (*llm.Completion, error) {
	cp := make([]llm.Message, len(msgs))
	copy(cp, msgs)
	m.allCalls = append(m.allCalls, cp) // 记录入参
	if m.err != nil {
		return nil, m.err
	}
	if m.idx >= len(m.responses) {
		return nil, errors.New("mockLLM：预设响应已用完")
	}
	r := m.responses[m.idx]
	m.idx++
	return r, nil
}

// call 返回第 n 次 Chat 收到的消息（0 起）
func (m *mockLLM) call(n int) []llm.Message { return m.allCalls[n] }

// newStubTool 返回一个只会返回固定结果/错误的假工具（mock Tool）
func newStubTool(name, result string, callErr error) tools.Tool {
	return tools.Tool{
		Name:        name,
		Description: "mock 工具",
		Call: func(_ string) (string, error) {
			return result, callErr
		},
	}
}

// 常用断言工具
func msgsHasRole(msgs []llm.Message, role string) bool {
	for _, m := range msgs {
		if m.Role == role {
			return true
		}
	}
	return false
}

func msgsContains(msgs []llm.Message, substr string) bool {
	for _, m := range msgs {
		if m.Role == "tool" && strings.Contains(m.Content, substr) {
			return true
		}
	}
	return false
}

// ---- 用例 ----

func TestRunDirectAnswer(t *testing.T) {
	// 模型第一轮就直接给最终答案
	model := &mockLLM{responses: []*llm.Completion{{Content: "你好，我是 Agent"}}}
	ag := New(model, []tools.Tool{tools.Calc, tools.Weather}, NewMemMemory(), 3)

	got, err := ag.Run(context.Background(), "s1", "打个招呼")
	if err != nil {
		t.Fatalf("Run 出错: %v", err)
	}
	if got != "你好，我是 Agent" {
		t.Fatalf("想要 %q，得到 %q", "你好，我是 Agent", got)
	}
	// 一次调用就答，未多跑
	if len(model.allCalls) != 1 {
		t.Fatalf("期望只调用模型 1 次，实际 %d", len(model.allCalls))
	}
}

func TestRunToolThenAnswer(t *testing.T) {
	// 第一轮：模型想调 calculator(1+2)；第二轮：给出最终答案
	model := &mockLLM{responses: []*llm.Completion{
		{ToolCalls: []tools.ToolCall{{ID: "t1", Name: "calculator", Arguments: `{"expr":"1+2"}`}}},
		{Content: "结果是 3"},
	}}
	ag := New(model, []tools.Tool{tools.Calc, tools.Weather}, NewMemMemory(), 3)

	got, err := ag.Run(context.Background(), "s1", "算 1+2")
	if err != nil {
		t.Fatalf("Run 出错: %v", err)
	}
	if got != "结果是 3" {
		t.Fatalf("想要 %q，得到 %q", "结果是 3", got)
	}
	// 第二轮入参里应带有 tool 角色的结果（calculator 算出 "3"）
	second := model.call(1)
	if !msgsHasRole(second, "tool") {
		t.Fatalf("第二轮应回传 tool 结果，实际 roles=%+v", second)
	}
	for _, m := range second {
		if m.Role == "tool" && m.Content != "3" {
			t.Fatalf("tool 结果应为 3，得到 %q", m.Content)
		}
	}
}

func TestRunToolErrorSelfCorrect(t *testing.T) {
	// 工具调用失败（B2 兜底）：第一轮模型想调一个总会失败的工具，
	// 不应中断，而应把错误喂回模型；第二轮模型给出答案。
	stub := newStubTool("boom", "", errors.New("炸了"))
	model := &mockLLM{responses: []*llm.Completion{
		{ToolCalls: []tools.ToolCall{{ID: "t1", Name: "boom", Arguments: "{}"}}},
		{Content: "处理了错误"},
	}}
	ag := New(model, []tools.Tool{stub}, NewMemMemory(), 3)

	got, err := ag.Run(context.Background(), "s1", "触发失败")
	if err != nil {
		t.Fatalf("工具失败不应中断 Agent，得到 err=%v", err)
	}
	if got != "处理了错误" {
		t.Fatalf("想要 %q，得到 %q", "处理了错误", got)
	}
	// 第二轮应把"工具执行失败"喂回给模型
	if !msgsContains(model.call(1), "工具执行失败") {
		t.Fatalf("第二轮应包含工具失败信息，实际=%+v", model.call(1))
	}
}

func TestRunExceedMaxSteps(t *testing.T) {
	// 模型总想调工具、从不给最终答案 → 应触发最大步数保护
	model := &mockLLM{responses: []*llm.Completion{
		{ToolCalls: []tools.ToolCall{{ID: "t1", Name: "calculator", Arguments: `{"expr":"1"}`}}},
		{ToolCalls: []tools.ToolCall{{ID: "t2", Name: "calculator", Arguments: `{"expr":"2"}`}}},
	}}
	ag := New(model, []tools.Tool{tools.Calc}, NewMemMemory(), 2)

	_, err := ag.Run(context.Background(), "s1", "一直调用")
	if err == nil || !strings.Contains(err.Error(), "最大步数") {
		t.Fatalf("应因达到最大步数报错，实际 err=%v", err)
	}
}

func TestRunStreamEventSequence(t *testing.T) {
	// 验证流式事件顺序：thought → tool_call → tool_result → answer（B1）
	model := &mockLLM{responses: []*llm.Completion{
		{ToolCalls: []tools.ToolCall{{ID: "t1", Name: "calculator", Arguments: `{"expr":"1+2"}`}}},
		{Content: "3"},
	}}
	ag := New(model, []tools.Tool{tools.Calc}, NewMemMemory(), 3)

	var types []string
	err := ag.RunStream(context.Background(), "s1", "算一下", func(ev Event) {
		types = append(types, ev.Type)
	})
	if err != nil {
		t.Fatalf("RunStream 出错: %v", err)
	}
	want := []string{"thought", "tool_call", "tool_result", "answer"}
	if len(types) != len(want) {
		t.Fatalf("事件数：想要 %v，得到 %v", want, types)
	}
	for i := range want {
		if types[i] != want[i] {
			t.Fatalf("事件顺序：想要 %v，得到 %v", want, types)
		}
	}
}