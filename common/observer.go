package common

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
	HTTPClient *http.Client
	LastUsed   time.Time
}

type Observer interface {
	EstablishConnection(task Task) (*http.Client, error)
	SendRequest(task Task, client *http.Client) (*http.Response, error)
	CleanupClients()
}

type HTTPObserver struct {
	Config       *Config
	Clients      map[string]*ClientInfo
	ClientsMutex sync.RWMutex
}

func (o *HTTPObserver) EstablishConnection(task Task) (*http.Client, error) {
	// TLS connection is only prepared in this function
	//      the actual TLS establishment is done in client.Do

	// reuse http.Client
	o.ClientsMutex.RLock()
	info, ok := o.Clients[task.IP]
	o.ClientsMutex.RUnlock()
	if ok && time.Since(info.LastUsed) < o.Config.ConnLifetime {
		o.ClientsMutex.Lock()
		info.LastUsed = time.Now()
		o.ClientsMutex.Unlock()
		return info.HTTPClient, nil
	}

	// new http.Client
	dailer := &net.Dialer{}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dailer.DialContext(ctx, network, task.IP+":443")
	}
	transport.MaxIdleConns = 0
	transport.IdleConnTimeout = 60 * time.Second

	client := &http.Client{
		Transport: transport,
		Timeout:   o.Config.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= o.Config.MaxRedirects {
				if o.Config.RedirectSucceed {
					return nil
				}
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// update LRU
	info = &ClientInfo{
		HTTPClient: client,
		LastUsed:   time.Now(),
	}

	o.ClientsMutex.Lock()
	o.Clients[task.IP] = info
	o.ClientsMutex.Unlock()

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
	for domain, info := range o.Clients {
		if time.Since(info.LastUsed) > o.Config.ConnLifetime {
			delete(o.Clients, domain)
		}
	}
}
