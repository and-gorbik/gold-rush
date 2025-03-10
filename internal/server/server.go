package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v3"
	"go.uber.org/zap"
)

const (
	defaultTimeout  = 1 * time.Second
	maxConnsPerHost = 1000
)

type serverError error

type Server struct {
	client *http.Client
	log    *zap.Logger
}

func Init(log *zap.Logger) *Server {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxConnsPerHost = 0
	t.MaxConnsPerHost = 0
	t.MaxIdleConnsPerHost = maxConnsPerHost

	client := &http.Client{
		Timeout:   defaultTimeout,
		Transport: t,
	}

	return &Server{client, log}
}

func (*Server) retryWithExponentialBackoff(f func() error) {
	b := newExponentialBackoff()
	var serr *serverError

	for {
		err := f()
		if err == nil || !errors.As(err, &serr) {
			return
		}

		<-time.After(b.NextBackOff())
	}
}

func (s *Server) processResponse(req *http.Request) ([]byte, error) {
	if req.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		if uerr, ok := err.(*url.Error); ok {
			if uerr.Temporary() || uerr.Timeout() {
				return nil, serverError(fmt.Errorf("do: %w", err))
			}
		}

		return nil, fmt.Errorf("do: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return nil, serverError(fmt.Errorf("status code: %d", resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return data, nil
}

func newExponentialBackoff() backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = time.Millisecond * 20
	b.MaxElapsedTime = 0 // retries never stops
	b.MaxInterval = time.Second * 5
	b.Reset()
	return b
}
