package collectors

import (
	"log"
	"sync"

	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"
)

type Task[T any, R any] func(p *rod.Page, input T) (R, error)

type ConcurrentRunner[T any, R any] struct {
	scheduler sche.Scheduler
}

func NewConcurrentRunner[T any, R any](s sche.Scheduler) *ConcurrentRunner[T, R] {
	return &ConcurrentRunner[T, R]{
		scheduler: s,
	}
}

func (cr *ConcurrentRunner[T, R]) Run(inputs []T, task Task[T, R]) <-chan R {
	results := make(chan R, len(inputs))
	var wg sync.WaitGroup

	for _, input := range inputs {
		wg.Add(1)
		go func(input T) {
			defer wg.Done()

			p, err := cr.scheduler.Get()
			if err != nil {
				log.Println("Error getting page from scheduler:", err)
				return
			}
			defer cr.scheduler.Put(p)

			result, err := task(p, input)
			if err != nil {
				log.Println("Error executing task:", err)
				return
			}

			results <- result
		}(input)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
