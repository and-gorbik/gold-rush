package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v3"
	jsoniter "github.com/json-iterator/go"
)

const (
	defaultTimeout  = 1 * time.Second
	maxConnsPerHost = 1000
)

type serverError struct {
	err error
}

func (se serverError) Error() string {
	return se.err.Error()
}

type Server struct {
	addr   string
	client *http.Client
}

func Init(addr string, timeout time.Duration) *Server {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxConnsPerHost = 0
	t.MaxConnsPerHost = 0
	t.MaxIdleConnsPerHost = maxConnsPerHost

	client := &http.Client{
		Timeout:   timeout,
		Transport: t,
	}

	s := &Server{addr, client}

	if err := s.Healthcheck(); err != nil {
		log.Fatalf("server: %v\n", err)
	}

	return s
}

func (*Server) retryWithExponentialBackoff(name string, f func() error) {
	b := newExponentialBackoff()
	var serr serverError

	for {
		err := f()
		if err == nil || !errors.As(err, &serr) {
			// don't retry
			return
		}

		dur := b.NextBackOff()
		log.Printf("%s: %v; it will be retried after %v\n", err, name, dur)
		<-time.After(dur)
	}
}

type errorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Server) processResponse(req *http.Request) ([]byte, error) {
	if req.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		if uerr, ok := err.(*url.Error); ok {
			if uerr.Temporary() || uerr.Timeout() {
				return nil, serverError{fmt.Errorf("do: %w", err)}
			}
		}

		return nil, fmt.Errorf("do: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode >= http.StatusInternalServerError ||
		resp.StatusCode == http.StatusTooManyRequests {
		return nil, serverError{fmt.Errorf("status code: %d", resp.StatusCode)}
	}

	if resp.StatusCode != http.StatusOK {
		var errDetail errorDetail
		if err := jsoniter.Unmarshal(data, &errDetail); err != nil {
			return nil, fmt.Errorf("unmarshal err detail: %w", err)
		}

		return nil, fmt.Errorf("unexpected status code: %d, details: %s", resp.StatusCode, errDetail.Message)
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
