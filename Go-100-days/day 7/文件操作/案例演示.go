package main

import (
	"bufio"
	"fmt"
	"os"
)

func main () {
	//创建一个新文件，写入内容5句 "你好，吴帅！！！"
	//1.打开文件 ./test.txt
	filePath := "./test.txt"
	file , err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err= %v", err)
		return
	}

	//及时关闭file句柄
	defer file.Close()
	//准备写入5句 "你好，吴帅！！！"
	str := "你好，吴帅！！！\n"
	//写入时，使用带缓存的* Writer
	writer := bufio.NewWriter(file)
	for i := 0;i < 5; i++ {
		writer.WriteString(str)
	}
	//因为writer是带缓存，因此在调用WriterString方法时，其实
	//内容实现写入到缓存的，所以需要调用Flush方法，将缓存的数据
	//真正写入到文件中， 否则文件中会没有数据
	writer.Flush()
}