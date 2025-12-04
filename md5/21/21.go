package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sum := 0
	lock := sync.Mutex{}
	var times int
	fmt.Printf("输入多少次一千次加1")
	fmt.Scanf("%d", &times)
	for range times {
		go func() {
			for range 1000 {
				lock.Lock()
				sum += 1
				lock.Unlock()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Println(sum)
}
