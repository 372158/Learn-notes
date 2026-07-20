package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "请指定 URL\n用法: go run 获取URL.go <url>\n")
        os.Exit(1)
    }

    url := os.Args[1]

    resp, err := http.Get(url)
    if err != nil {
        fmt.Fprintf(os.Stderr, "请求失败: %v\n", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    // ========== 🎯 重点在这里：获取状态码 ==========
    // 1. 打印数字状态码（如 200, 404, 500）
    fmt.Printf("状态码: %d\n", resp.StatusCode)
    
    // 2. 打印描述性文本（如 "200 OK", "404 Not Found"）
    fmt.Printf("状态: %s\n", resp.Status)

    // 3. 【推荐】如果状态码不是 200，打印错误并退出，不读取网页内容
    if resp.StatusCode != http.StatusOK {
        fmt.Fprintf(os.Stderr, "错误：服务器返回非 200 状态码\n")
        os.Exit(1)
    }

    // 只有状态码是 200 时才读取内容
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "读取内容失败: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("网页内容长度: %d 字节\n", len(body))
    // fmt.Printf("%s", body) // 如果要打印内容，取消注释这行
}