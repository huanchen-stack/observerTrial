//go:build test
// +build test

package main

import (
	"net/http"
	"testing"
	"time"
)

func TestObserver(t *testing.T) {
	config := &Config{
		MaxRedirects:    10,
		Timeout:         30 * time.Second,
		RedirectSucceed: true,
	}

	task := Task{IP: "93.184.216.34", Domain: "example.com", Endpoint: "/"}

	observer := &HTTPObserver{config: config}
	client, err := observer.EstablishConnection(task)
	if err != nil {
		t.Fatalf("Failed to establish connection: %v", err)
	}
	resp, err := observer.SendRequest(task, client)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}
