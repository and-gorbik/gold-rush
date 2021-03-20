package infrastructure

import (
	"net/http"
	"time"
)

func BuildHTTPClient(timeout time.Duration, maxIdleConns, maxConnsPerHost, maxIdleConnsPerHost int) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost

	return &http.Client{
		Timeout:   timeout,
		Transport: t,
	}
}
