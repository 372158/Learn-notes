package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// 模拟一个"查天气"工具说明书
	tool := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        "get_weather",
			"description": "查询指定城市的天气",
			"parameters": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"city": map[string]any{"type": "string", "description": "城市名"},
				},
				"required": []string{"city"},
			},
		},
	}
	body, _ := json.Marshal(map[string]any{
		"model":    os.Getenv("LLM_MODEL"),
		"messages": []map[string]any{{"role": "user", "content": "请问北京今天天气怎么样？"}},
		"tools":    []any{tool},
		"stream":   false,
	})

	req, _ := http.NewRequest("POST", os.Getenv("LLM_API_URL"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LLM_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("出错:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	fmt.Println(string(raw))
}