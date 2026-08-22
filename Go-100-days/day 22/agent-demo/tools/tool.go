package tools

// ToolCall 模型通过 Function Calling 告诉我们"想调用哪个工具、传什么参数"。
// 注意：这是模型返回的 JSON 形状，不是 Go 结构体名。
type ToolCall struct {
	ID        string `json:"id"`
	Name      string `json:"name"`      // 要调用的工具名
	Arguments string `json:"arguments"` // 传给工具的参数字符串(JSON)，后面解析
}

// Tool 一个"能动手的工具"。A5 会写具体工具（计算器/查资料库等）。
type Tool struct {
	Name        string                              // 模型眼里这个工具叫啥
	Description string                              // 给模型看的说明书：什么时候该用我
	Call        func(argsJSON string) (string, error) // 真正执行工具；入参参数字符串(JSON)，出参结果字符串(JSON)
}

// Register 工具名 -> 工具对象映射，Agent 按名字分派执行
func Register(ts []Tool) map[string]Tool {
	m := make(map[string]Tool, len(ts))
	for _, t := range ts {
		m[t.Name] = t
	}
	return m
}

// ToolSchema 发送给模型的工具"说明书"（OpenAI 兼容的 function 格式）。
// A4 会教它怎么进请求体，这里先定义好形状。
type ToolSchema struct {
	Type     string         `json:"type"` // 固定 "function"
	Function FunctionSchema `json:"function"`
}

type FunctionSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"` // JSON Schema：声明参数名/类型/必填
}

// Schema 把 Tool 转成发给模型看的说明书（Name/Description 现成，Parameters 由具体工具提供）
func Schema(t Tool, params map[string]any) ToolSchema {
	return ToolSchema{
		Type: "function",
		Function: FunctionSchema{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  params,
		},
	}
}