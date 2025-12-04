package tool

import (
	"fmt"
	"sync"
)

type Task struct {
	Runnable func(workerId int)
}

func Distribute() {
	ch := make(chan Task, 10)
	var sum int
	var wg sync.WaitGroup
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
	close(ch)
	wg.Wait()
	fmt.Printf("最终结果: %d\n", sum)
}
