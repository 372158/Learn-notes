package tools

import (
	"fmt"
	"strconv"
	"strings"
)

// calcArgs 这个工具要解析的参数形状——模型会按这个传
type calcArgs struct {
	Expr string `json:"expr"` // 比如 "1+2"、"3*(4-1)"
}

// Calc 计算工具：拿到表达式字符串，算出结果（极简递归下降，含边界保护）
var Calc = Tool{
	Name:        "calculator",
	Description: "计算一个数学表达式（只支持 + - * / 和括号与数字），如 12+34、3*(4-1)、10/4。传入 expr 字段。",
	Call: func(argsJSON string) (string, error) {
		var args calcArgs
		if err := ParseArgs(argsJSON, &args); err != nil {
			return "", err
		}
		if strings.TrimSpace(args.Expr) == "" {
			return "", fmt.Errorf("表达式为空")
		}
		val, err := evalSimple(args.Expr)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	},
}

// evalSimple 递归下降，支持 + - * / ( 数字，按优先级从左到右
func evalSimple(expr string) (float64, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}
	if len(tokens) == 0 {
		return 0, fmt.Errorf("无有效 token")
	}
	pos := 0


	var parseExpr func() (float64, error)
	var parseTerm func() (float64, error)
	var parseFactor func() (float64, error)

	parseExpr = func() (float64, error) {
		v, err := parseTerm()
		if err != nil {
			return 0, err
		}
		for pos < len(tokens) {
			if tokens[pos] != "+" && tokens[pos] != "-" {
				break
			}
			op := tokens[pos]
			pos++
			rhs, e := parseTerm()
			if e != nil {
				return 0, e
			}
			if op == "+" {
				v += rhs
			} else {
				v -= rhs
			}
		}
		return v, nil
	}
	parseTerm = func() (float64, error) {
		v, err := parseFactor()
		if err != nil {
			return 0, err
		}
		for pos < len(tokens) {
			if tokens[pos] != "*" && tokens[pos] != "/" {
				break
			}
			op := tokens[pos]
			pos++
			rhs, e := parseFactor()
			if e != nil {
				return 0, e
			}
			if op == "*" {
				v *= rhs
			} else {
				v /= rhs
			}
		}
		return v, nil
	}
	parseFactor = func() (float64, error) {
		if pos >= len(tokens) {
			return 0, fmt.Errorf("表达式不完整")
		}
		if tokens[pos] == "(" {
			pos++
			v, err := parseExpr()
			if err != nil {
				return 0, err
			}
			if pos >= len(tokens) || tokens[pos] != ")" {
				return 0, fmt.Errorf("缺少右括号")
			}
			pos++
			return v, nil
		}
		if !isNumber(tokens[pos]) {
			return 0, fmt.Errorf("无法识别的输入: %s", tokens[pos])
		}
		f, _ := strconv.ParseFloat(tokens[pos], 64)
		pos++
		return f, nil
	}

	return parseExpr()
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// tokenize 拆 token，遇到非数字/操作符的字符返回错误
func tokenize(s string) ([]string, error) {
	var out []string
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			out = append(out, string(cur))
			cur = nil
		}
	}
	for _, ch := range s {
		switch ch {
		case '+', '-', '*', '/', '(', ')':
			flush()
			out = append(out, string(ch))
		case ' ', '\n', '\t':
			continue
		default:
			// 只允许数字和小数点进入数字 token
			if !(ch >= '0' && ch <= '9' || ch == '.') {
				flush()
				return nil, fmt.Errorf("包含非法字符: %q", ch)
			}
			cur = append(cur, ch)
		}
	}
	flush()
	return out, nil
}