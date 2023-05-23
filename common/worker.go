package common

import (
	"net/http"
	"time"
)

type Result struct {
	Task    Task
	Resp    *http.Response
	Err     error
	Elapsed time.Duration
}

func Worker(id int, tasks <-chan Task, results chan<- Result, observer Observer) {
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
