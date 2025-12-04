package tool

import (
	"fmt"
	"sync"
)

type Task struct {
	Runnable func(workerId int)
}

func Distribute(firstsum, times int) int {
	ch := make(chan Task, 10)
	var wg sync.WaitGroup
	sum := firstsum
	var lock sync.Mutex
	for id := range 10 {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for t := range ch {
				t.Runnable(workerId)
			}
		}(id)
	}
	for i := range times {
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
	return sum
}
