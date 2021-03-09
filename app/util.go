package app

import (
	"net/http"
	"time"

	"gold-rush/models"
)

func buildHTTPClient(timeout time.Duration) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost

	return &http.Client{
		Timeout:   timeout,
		Transport: t,
	}
}

func readError(err error) (string, bool) {
	if e, ok := err.(*models.BusinessError); ok {
		return e.Message, true
	}

	return err.Error(), false
}
