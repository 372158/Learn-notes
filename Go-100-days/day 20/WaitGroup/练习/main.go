package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) // 每个循环登记 1 个任务
		go func(n int) {
			defer wg.Done()
			fmt.Printf(" 任务 %d 完成\n", n)
		}(i)
	}

	wg.Wait() //等待 5 个完成
	fmt.Println("全部完成")
}
