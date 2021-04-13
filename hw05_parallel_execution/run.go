package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	countOfTasks := len(tasks)

	tasksChannel := make(chan Task, len(tasks))
	errorChannel := make(chan error, len(tasks))
	shutdownChannel := make(chan interface{})

	var wg sync.WaitGroup

	for i := n; i != 0; i-- {
		wg.Add(1)
		go startConsumer(&wg, shutdownChannel, tasksChannel, errorChannel)
	}

	produce(tasks, tasksChannel)

	for err := range errorChannel {
		if err != nil {
			m--
		}

		countOfTasks--

		if m == 0 {
			forceClose(&wg, shutdownChannel, tasksChannel, errorChannel)
			return ErrErrorsLimitExceeded
		}

		if countOfTasks == 0 {
			forceClose(&wg, shutdownChannel, tasksChannel, errorChannel)
			return nil
		}
	}

	return nil
}

func forceClose(wg *sync.WaitGroup, shutdownChannel chan interface{}, tasksQueue chan Task, errorChannel chan error) {
	close(shutdownChannel)
	close(tasksQueue)
	wg.Wait()
	close(errorChannel)
}

func produce(tasks []Task, tasksChannel chan Task) {
	for _, t := range tasks {
		tasksChannel <- t
	}
}

func startConsumer(wg *sync.WaitGroup, shutdownChannel chan interface{}, tasksQueue chan Task, errorChannel chan error) {
	defer wg.Done()

	for {
		select {
		case <-shutdownChannel:
			return
		default:
		}

		f, ok := <-tasksQueue
		if ok {
			err := f()
			errorChannel <- err
		}
	}
}
