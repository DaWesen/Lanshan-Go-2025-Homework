package task

import (
	"sync"
)

type Task struct {
	tasks chan func()
	wg    sync.WaitGroup
}

func New(workerCount int) *Task {
	t := &Task{
		tasks: make(chan func(), 100),
	}
	for range workerCount {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			for task := range t.tasks {
				task()
			}
		}()
	}

	return t
}

func (t *Task) Submit(task func()) {
	t.tasks <- task
}
func (t *Task) Wait() {
	close(t.tasks)
	t.wg.Wait()
}
