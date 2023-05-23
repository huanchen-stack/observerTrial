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

type stage struct {
	workerCount int
	config      *c.Config
	results     []c.Result
}

func main() {
	observer := &c.HTTPObserver{
		Clients: make(map[string]*c.ClientInfo, 1000),
	}

	stages := []stage{
		{
			workerCount: 1000,
			config: &c.Config{
				Timeout:         10 * time.Second,
				MaxRedirects:    0,
				RedirectSucceed: true,
				ConnLifetime:    10 * time.Minute,
			},
		},
		{
			workerCount: 1000,
			config: &c.Config{
				Timeout:         10 * time.Second,
				MaxRedirects:    0,
				RedirectSucceed: true,
				ConnLifetime:    10 * time.Minute,
			},
		},
		{
			workerCount: 3000,
			config: &c.Config{
				Timeout:         20 * time.Second,
				MaxRedirects:    0,
				RedirectSucceed: true,
				ConnLifetime:    10 * time.Minute,
			},
		},
		{
			workerCount: 5000,
			config: &c.Config{
				Timeout:         30 * time.Second,
				MaxRedirects:    0,
				RedirectSucceed: true,
				ConnLifetime:    10 * time.Minute,
			},
		},
	}

	tasks, err := c.ReadTasksFromCsv("input.csv")
	if err != nil {
		log.Fatal(err)
	}

	for i, stage := range stages {
		fmt.Printf("Running stage %d with %d workers...\n", i+1, stage.workerCount)

		start := time.Now()
		runStage(observer, tasks, &stage)
		duration := time.Since(start)

		successes, failures := countResults(stage.results)
		fmt.Printf("Stage %d completed in %s with %d successes and %d failures.\n", i+1, duration, successes, failures)
	}

}

func runStage(observer *c.HTTPObserver, tasks []c.Task, stage *stage) {
	observer.Config = stage.config

	taskCh := make(chan c.Task, 1000)
	resultCh := make(chan c.Result, 1000)

	var wg sync.WaitGroup

	for i := 0; i < stage.workerCount; i++ {
		wg.Add(1)
		go func(i int) {
			c.Worker(i, taskCh, resultCh, observer)
			wg.Done()
		}(i)
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

	for result := range resultCh {
		stage.results = append(stage.results, result)
	}

}

func countResults(results []c.Result) (successes, failures int) {
	for _, result := range results {
		if result.Err != nil {
			failures++
		} else {
			successes++
		}
	}
	return successes, failures
}
