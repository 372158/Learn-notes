package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // 启动一个 goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // 从通道 ch 中接收
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp ,err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbyters, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // 关闭响应体
	if err != nil {
		ch <- fmt.Sprint("while reading %s: %v", url ,err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbyters, url)
}
