package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			count++
		}()
	}

	wg.Wait()
	fmt.Println("count=", count) //期望 1000 ？
}
