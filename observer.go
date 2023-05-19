package main

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	MaxRedirects    int
	Timeout         time.Duration
	RedirectSucceed bool
	ConnLifetime    time.Duration
}

type Task struct {
	IP       string
	Domain   string
	Endpoint string
}

type ClientInfo struct {
	httpClient *http.Client
	lastUsed   time.Time
}

type Observer interface {
	EstablishConnection(task Task) (*http.Client, error)
	SendRequest(task Task, client *http.Client) (*http.Response, error)
	CleanupClients()
}

type HTTPObserver struct {
	config       *Config
	clients      map[string]*ClientInfo
	clientsMutex sync.RWMutex
}

func (o *HTTPObserver) EstablishConnection(task Task) (*http.Client, error) {
	// TLS connection is only prepared in this function
	//      the actual TLS establishment is done in client.Do

	// reuse http.Client
	o.clientsMutex.RLock()
	info, ok := o.clients[task.IP]
	o.clientsMutex.RUnlock()
	if ok && time.Since(info.lastUsed) < o.config.ConnLifetime {
		o.clientsMutex.Lock()
		info.lastUsed = time.Now()
		o.clientsMutex.Unlock()
		return info.httpClient, nil
	}

	// new http.Client
	dailer := &net.Dialer{}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dailer.DialContext(ctx, network, task.IP+":443")
	}
	transport.MaxIdleConns = 10
	transport.IdleConnTimeout = 30 * time.Second

	client := &http.Client{
		Transport: transport,
		Timeout:   o.config.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= o.config.MaxRedirects {
				if o.config.RedirectSucceed {
					return nil
				}
				return http.ErrHandlerTimeout
			}
			return nil
		},
	}

	// update LRU
	info = &ClientInfo{
		httpClient: client,
		lastUsed:   time.Now(),
	}

	o.clientsMutex.Lock()
	o.clients[task.IP] = info
	o.clientsMutex.Unlock()

	return client, nil
}

func (o *HTTPObserver) SendRequest(task Task, client *http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://"+task.Domain+task.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	return resp, err
}

func (o *HTTPObserver) CleanupClients() {
	for domain, info := range o.clients {
		if time.Since(info.lastUsed) > o.config.ConnLifetime {
			delete(o.clients, domain)
		}
	}
}
