package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main(){
	conn, err := net.Dial("tcp","127.0.0.1:8888")
	if err != nil {
		fmt.Println("client dial err=",err)
		return
	}

	//功能一： 客户端可以发送单行数据，然后就退出
	//os.Stdin 代表标准输入【终端】
	reader := bufio.NewReader(os.Stdin)

	//从终端读取一行用户输入，并准备发送给服务器
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("readString err=",err)
	}


	//再将line发送给服务器
	n, err := conn.Write([]byte(line))
	if err != nil{
		fmt.Println("conn。Write err=", err)
	}

	fmt.Printf("客户端发送了%d 字节的数据，并退出", n)
}