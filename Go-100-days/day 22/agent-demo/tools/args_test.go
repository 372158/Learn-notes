package tools

import (
	"strings"
	"testing"
)

type args struct {
	City string `json:"city"`
}

func TestParseArgs(t *testing.T) {
	t.Run("合法JSON", func(t *testing.T) {
		var a args
		if err := ParseArgs(`{"city":"北京"}`, &a); err != nil {
			t.Fatalf("不应报错: %v", err)
		}
		if a.City != "北京" {
			t.Fatalf("city 应为 北京，得到 %q", a.City)
		}
	})
	t.Run("空参数兜底", func(t *testing.T) {
		var a args
		if err := ParseArgs("   ", &a); err == nil {
			t.Fatal("空参数应报错")
		}
	})
	t.Run("非法JSON兜底", func(t *testing.T) {
		var a args
		if err := ParseArgs("not-json", &a); err == nil || !strings.Contains(err.Error(), "合法 JSON") {
			t.Fatalf("应提示参数非法，实际 err=%v", err)
		}
	})
}