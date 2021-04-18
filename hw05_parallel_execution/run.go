package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n == 0 {
		panic("Workers count should be more than zero")
	}

	if len(tasks) == 0 {
		return nil
	}

	countOfTasks := len(tasks)

	tasksCh := make(chan Task, len(tasks))
	errorCh := make(chan error, len(tasks))
	shutdown := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			newConsumer(shutdown, tasksCh, errorCh)
		}()
	}

	produce(tasks, tasksCh)

	defer func() {
		close(shutdown)
		wg.Wait()
	}()

	for err := range errorCh {
		if err != nil && m == 0 {
			panic("workers configured as error impossible")
		}

		if err != nil {
			m--
		}

		countOfTasks--

		if m == 0 {
			return ErrErrorsLimitExceeded
		}

		if countOfTasks == 0 {
			return nil
		}
	}

	return nil
}

func produce(tasks []Task, tasksChannel chan<- Task) {
	defer close(tasksChannel)

	for _, t := range tasks {
		tasksChannel <- t
	}
}

func newConsumer(shutdownChannel chan struct{}, tasks chan Task, errorCh chan error) {
	for {
		f, ok := <-tasks

		select {
		case <-shutdownChannel:
			return
		default:
		}

		if ok {
			err := f()
			errorCh <- err
		}
	}
}
