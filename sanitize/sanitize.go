//go:build !test
// +build !test

package main

import (
	"fmt"
	"log"
	c "observerGO/common"
	"sync"
	"time"
)

func main() {
	config := &c.Config{
		Timeout:         10 * time.Second,
		MaxRedirects:    0,
		RedirectSucceed: true,
		ConnLifetime:    10 * time.Minute,
	}

	observer := &c.HTTPObserver{
		Config:  config,
		Clients: make(map[string]*c.ClientInfo, 1000),
	}

	taskCh := make(chan c.Task, 1000)
	resultCh := make(chan c.Result, 1000)

	var wg sync.WaitGroup

	for i := 0; i < 750; i++ {
		wg.Add(1)
		go func(i int) {
			c.Worker(i, taskCh, resultCh, observer)
			wg.Done()
		}(i)
	}

	tasks, err := c.ReadTasksFromCsv("input.csv")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	successfulTasks := make([]c.Task, 0)

	for result := range resultCh {
		if result.Err != nil {
			fmt.Printf("Task: %+v\nError: %s\nElapsed: %s\n\n", result.Task, result.Err, result.Elapsed)
		} else {
			fmt.Printf("Task: %+v\nStatus: %s\nElapsed: %s\n\n", result.Task, result.Resp.Status, result.Elapsed)
			successfulTasks = append(successfulTasks, result.Task)
		}
	}

	err = c.WriteTasksToCsv("input.csv", successfulTasks)
	if err != nil {
		log.Fatal(err)
	}

}
