//go:build !test

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	config := &Config{
		Timeout:         10 * time.Second,
		MaxRedirects:    0,
		RedirectSucceed: true,
		ConnLifetime:    10 * time.Minute,
	}

	observer := &HTTPObserver{
		config:  config,
		clients: make(map[string]*ClientInfo, 1000),
	}

	taskCh := make(chan Task, 20000)
	resultCh := make(chan Result, 20000)

	var wg sync.WaitGroup

	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func(i int) {
			worker(i, taskCh, resultCh, observer)
			wg.Done()
		}(i)
	}

	tasks, err := readTasksFromCsv("input_very_large.csv")
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		taskCh <- task
	}

	close(taskCh)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		if result.Err != nil {
			fmt.Printf("Task: %+v\nError: %s\nElapsed: %s\n\n", result.Task, result.Err, result.Elapsed)
		} else {
			fmt.Printf("Task: %+v\nStatus: %s\nElapsed: %s\n\n", result.Task, result.Resp.Status, result.Elapsed)
		}
	}
}
