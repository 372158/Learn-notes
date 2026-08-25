package tools

import (
	"encoding/json"
	"fmt"
)

// weatherArgs 这个工具要解析的参数形状
type weatherArgs struct {
	City string `json:"city"`
}

// 假数据库：城市 -> 天气（真实接口可换成外呼，这里演示工具组织方式）
var cityWeather = map[string]string{
	"北京": "26°C 晴",
	"上海": "28°C 多云",
	"广州": "30°C 雷阵雨",
}

// Weather 查天气工具
var Weather = Tool{
	Name:        "get_weather",
	Description: "查询指定城市的当前天气。可查城市：北京、上海、广州。传入 city 字段。",
	Call: func(argsJSON string) (string, error) {
		var args weatherArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("参数解析失败: %w", err)
		}
		v, ok := cityWeather[args.City]
		if !ok {
			return "", fmt.Errorf("暂无 %s 的天气数据", args.City)
		}
		return v, nil
	},
}