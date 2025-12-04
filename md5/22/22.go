package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Runnable func(workerId int)
}

func main() {
	ch := make(chan Task, 10)
	var sum int
	var lock sync.Mutex
	for id := range 10 {
		go func(workerId int) {
			for t := range ch {
				t.Runnable(workerId)
			}
		}(id)
	}

	fmt.Print("输入初始数字: ")
	fmt.Scan(&sum)

	for i := range 20 {
		j := i
		task := Task{
			Runnable: func(workerId int) {
				lock.Lock()
				sum++
				current := sum
				lock.Unlock()
				fmt.Printf("workerId %v：task %v 自增，当前值: %v\n", workerId, j, current)
			},
		}
		ch <- task
	}

	time.Sleep(1 * time.Second)
	close(ch)
	fmt.Printf("最终结果: %d\n", sum)
}
