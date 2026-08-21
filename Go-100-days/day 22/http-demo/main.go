package main

import (
	"fmt"
	"log"
	"net/http"
)

// 一个普通函数，签名满足 Handler 的要求（w 写响应，r 拿请求）
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "你好，我是原生 net/http 服务")
	fmt.Fprintf(w, "路径=%s 方法=%s\n", r.URL.Path, r.Method)
}

func main() {
	// 路由表：路径 -> 处理器
	mux := http.NewServeMux()

	// 方式1：普通函数用 http.HandlerFunc 包成 Handler
	mux.HandleFunc("/hello", helloHandler)

	// 方式2：直接内联匿名函数
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "命中默认路径 %s\n", r.URL.Path)
	})

	log.Println("原生服务器启动于 :8090")
	// 监听 + 分发循环（Gin 的 r.Run() 底层就是它）
	log.Fatal(http.ListenAndServe(":8090", mux))
}
