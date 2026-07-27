package main

//典型的 生产者---消费者模式
//1)开启一个writeData协程，向管道ch中写入50个整数
//2）开启一个readData协程，从管道ch中读取write DATa写入的数据
//3）注意writeData和readData协程都完成工作才能退出管道

import (
	"fmt"
	"sync"
	"time"
)

func writeData(ch chan int,wg *sync.WaitGroup) {
	defer wg.Done()// 协程结束是通知WaitGroup
	for i := 0; i < 50; i++ {
		ch <- i
		fmt.Printf("写入数据： %d",i)
		time.Sleep(10 * time.Millisecond) //模拟写入耗时，让读取看得更清楚
	}
	close(ch)
	fmt.Println("writeData 协程完成")
}

func readData(ch chan int,wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {//循环读取，知道管道关闭
		fmt.Printf("读取数据：%d\n",v)
		time.Sleep(20 * time.Millisecond) //模拟读取时间		
	}
	fmt.Printf("readData 协程完成")
}


func main() {
	//创建一个容量为10 的缓冲管道
	ch := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(2)//表示有两个协程需要等待


	//启动写协程
	go writeData(ch , &wg) 
	//启动读协程
	go readData(ch , &wg)

	//主线程等待两个协程全部完成
	wg.Wait()
	fmt.Println("主线程退出")
}
