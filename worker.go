package main

import (
	"encoding/csv"
	"net/http"
	"os"
	"strings"
	"time"
)

func readTasksFromCsv(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for _, line := range lines {
		tasks = append(tasks, Task{
			IP:       strings.TrimSpace((line[0])),
			Domain:   strings.TrimSpace((line[1])),
			Endpoint: strings.TrimSpace(line[2]),
		})
	}

	return tasks, nil
}

type Result struct {
	Task    Task
	Resp    *http.Response
	Err     error
	Elapsed time.Duration
}

func worker(id int, tasks <-chan Task, results chan<- Result, observer Observer) {
	for task := range tasks {
		start := time.Now()
		client, err := observer.EstablishConnection(task)
		if err != nil {
			results <- Result{Task: task, Resp: nil, Err: err, Elapsed: time.Since(start)}
			continue
		}
		resp, err := observer.SendRequest(task, client)
		results <- Result{Task: task, Resp: resp, Err: err, Elapsed: time.Since(start)}
	}
}
