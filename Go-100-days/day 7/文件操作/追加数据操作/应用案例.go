package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//打开一个存在的文件，在原来的内容的位置追加内容
	//1.打开文件 ./test.txt
	filePath := "./test.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_APPEND,0666)
	if err != nil {
		fmt.Printf("open file err=%v", err)
		return 
	}
	//及时关闭句柄
	defer file.Close()
	

	//你准备写入的话
	str := `沧海恨水碧悠悠，至今已过百万秋。
			怀忆凭栏意寞寞，离愁别绪梦沉沉。
			门殚户尽算浮生，功成垂败成若梦。
			今朝归来我非我，物是人非泪已空`

	writer := bufio.NewWriter(file)
	writer.WriteString(str)

	writer.Flush()
}