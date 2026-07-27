package main

import (
	"fmt"
	"sync"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}

	//只需要检查到平方根即可
	for i := 2; i*i <= n; i++ {
		if n % i == 0 {
			return false
		}
	}
	return true
}

func main() {
	start := time.Now()

	//数字范围 1 -200000
	const maxNum = 200000
	//8个工作协程并发处理
	const numWokers = 8

	//1.创建任务管道和结果管道
	//存放带判断的数字
	taskChan := make(chan int, 100)
	
	//存放素数结果
	resultChan := make(chan int, 100)

	var wg sync.WaitGroup

	//2. 启动Worker协程池
	for i := 0; i < numWokers; i++ {
		wg.Add(1)
		go func (id int)  {
			defer wg.Done()
			for num := range taskChan {
				if isPrime(num){
					resultChan <- num
				}
			}
			//当前Worker 退出时打印日志（方便观察）
			fmt.Printf("Worker %d 退出\n", id)
		}(i)
	}

	//3.启动一个独立的协程，负责向任务管道发送数据
	go func() {
		for i := 0; i <= maxNum; i++ {
			taskChan <- i
		}
		close(taskChan)//所有任务发送完毕，关闭任务管道
		fmt.Println("所有任务已分配完毕")
	}()

	//4.启动一个独立的协程，等待所有的 worker 完成后关闭结果管道
	go func ()  {
			//等待所有Worker 退出
			wg.Wait() 
			close(resultChan)
			fmt.Println("所有 worker 已完成，结果管道关闭")
	}()

	//5.主协程收集并打印结果
	primes := make([]int , 0, 10000)//预分配容量提高性能
	for prime := range resultChan {
		primes = append(primes, prime)
	}

	//6.输出结果
	fmt.Printf("1-%d 中共有 %d 个素数\n", maxNum, len(primes))
	fmt.Printf("计算耗时：%v\n", time.Since(start))
}


