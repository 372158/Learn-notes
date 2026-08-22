package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs  <-chan int, results chan<- int, done <-chan struct{}) {
	for {
		select {
		case j := <-jobs:
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("worker%d 完成任务 %d\n", id, j)
			results <- j * 2
		case <- done :
			fmt.Printf("worker%d 下班\n", id)
			return
		}
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	done := make(chan struct{})

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results, done)
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
	}

	total := 0
	for i := 1;i <= 3; i++ {
		total += <- results
	}

	close(done)
	time.Sleep(200 *time.Millisecond) // 留时间看 worker 打印“下班”
	fmt.Println("前三个结果之和：", total)

}