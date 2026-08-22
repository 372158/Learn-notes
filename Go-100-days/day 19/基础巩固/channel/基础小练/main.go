package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})

	go func(){
		for {
			select {
			case <- done:
				fmt.Println("worker exit")
				return
			case <-time.After(500 *time.Millisecond):
				fmt.Println("haerbeat") 
			}
		}
	}()

	time.Sleep(2 *time.Second)
	close(done)
	fmt.Println("main exit")
}