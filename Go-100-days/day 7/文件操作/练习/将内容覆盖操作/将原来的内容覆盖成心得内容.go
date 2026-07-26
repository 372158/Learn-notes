package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//打开一个存在的文件 test.txt 将原来的内容覆盖为 
	// 沧海恨水碧悠悠，至今已过百万秋。
	// 怀忆凭栏意寞寞，离愁别绪梦沉沉。
	// 门殚户尽算浮生，功成垂败成若梦。
	// 今朝归来我非我，物是人非泪已空
	//创建一个新文件，写入"吴帅的一生："
	//1. 打开文件 ./test.txt
	filePath := "./test.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_TRUNC , 0666)
	if err != nil {
		fmt.Printf("open file err=%v", err)
		return
	}
	//及时关闭权柄
	defer file.Close()
	//准备写入的话
	str := "吴帅的一生：\r\n"
	//写入时，使用带缓存的 *Writer
	writer := bufio.NewWriter(file)
	writer.WriteString(str)
	//因为writer是带缓存的，因此在调用WriterString 方法时，其实
	//内容是先写入到缓存的，所哟需要调用Flush方法，将缓存的数据
	//真正写入到文件，否则文件中会没有数据
	writer.Flush()


}