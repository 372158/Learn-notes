package main

import (
	"fmt"
	"math/rand"
	"time"
)
type Person struct {
	Name string
	Age int
	Address string
}

func setPerson() Person {
	//随机姓名池
	names := []string{"张三","李四","王五","赵六","孙七","周八","吴九","郑十","钱十一","冯十二"}
	//随机地址池
	address := []string{"北京","上海","广州","深圳","杭州","成都","武汉","西安","重庆","南京"}

	return Person{
		Name: names[rand.Intn(len(names))],
		Age:  rand.Intn(60) + 18, // 18~77岁
		Address: address[rand.Intn(len(address))],
	}
}

func main() {
	//1.初始话随机种子（保证每次运行结果不同）
	rand.Seed(time.Now().UnixNano())

	//2.创建channel （容量为10， 也可以为无缓冲，但这里用缓冲更高效）
	ch := make(chan Person,10)

	//3. 启动一个goroutine 生成10个Person并发送到channel
	go func() {
		for i := 0;i < 10; i++ {
			p := setPerson()
			ch <- p
		}
		close(ch)// 发送完关闭channel，通知接收方结束
	}()

	//4.遍历channel ，接受并打印channel，通知接收方结束
	fmt.Println("生成的10个 Person 信息如下：")
	count := 1
	for p := range ch {
		fmt.Printf("%d.姓名： %s, 年龄：%d, 地址：%s\n", count, p.Name, p.Age, p.Address)
		time.Sleep(time.Second)
		count++
	}
}