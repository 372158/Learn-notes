// 请编写一个程序，完成如下功能:
//  1) 在主线程(可以理解成进程)中，开启一个 goroutine, 该协程每隔 1 秒输出"hello,world"
//  2) 在主线程中也每隔一秒输出"hello,golang", 输出 10 次后，退出程序
//  3) 要求主线程和 goroutine 同时执行.
//  4) 画出主线程和协程执行流程图

package main

import (
	"fmt"
	"strconv"
	"time"
)

//编写一个函数，每个一秒输出"你好，吴帅！！！"
func test() {
	for i := 1; i <= 10; i++ {
		fmt.Println("test() 你好，吴帅！！！" + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}

func main() {
	go test()

	for i := 1; i <= 5; i++ {
		fmt.Println(" main() 你好，golang..." + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}