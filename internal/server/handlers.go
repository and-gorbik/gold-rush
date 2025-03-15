package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func (s *Server) IssueLicense(ctx context.Context, coins []int) License {
	var license *License
	var err error

	s.retryWithExponentialBackoff("issueLicense", func() error {
		license, err = s.issueLicense(ctx, coins)
		return err
	})

	if err != nil {
		log.Fatalf("issue license failed unexpectedly: %v", err)
	}

	return *license
}

func (s *Server) ExploreArea(ctx context.Context, area Area) ExploredArea {
	var ea *ExploredArea
	var err error

	s.retryWithExponentialBackoff("exploreArea", func() error {
		ea, err = s.exploreArea(ctx, area)
		return err
	})

	if err != nil {
		log.Fatalf("explore area failed unexpectedly: %v", err)
	}

	return *ea
}

func (s *Server) Dig(ctx context.Context, params DigParams) []string {
	var treasures []string
	var err error

	s.retryWithExponentialBackoff("dig", func() error {
		treasures, err = s.dig(ctx, params)
		return err
	})

	if err != nil {
		log.Fatalf("dig failed unexpectedly: %v", err)
	}

	return treasures
}

func (s *Server) Cash(ctx context.Context, treasureID string) []int {
	var coins []int
	var err error

	s.retryWithExponentialBackoff("cash", func() error {
		coins, err = s.cash(ctx, treasureID)
		return err
	})

	if err != nil {
		log.Fatalf("cash failed unexpectedly: %v", err)
	}

	return coins
}

func (s *Server) issueLicense(ctx context.Context, coins []int) (*License, error) {
	data, err := jsoniter.Marshal(coins)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addr+"/licenses", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	data, err = s.processResponse(req)
	if err != nil {
		return nil, fmt.Errorf("process response: %w", err)
	}

	var li License
	if err := jsoniter.Unmarshal(data, &li); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &li, nil
}

func (s *Server) exploreArea(ctx context.Context, area Area) (*ExploredArea, error) {
	data, err := jsoniter.Marshal(area)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addr+"/explore", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	data, err = s.processResponse(req)
	if err != nil {
		return nil, fmt.Errorf("process response: %w", err)
	}

	var ea ExploredArea
	if err := jsoniter.Unmarshal(data, &ea); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &ea, nil
}

func (s *Server) dig(ctx context.Context, params DigParams) ([]string, error) {
	data, err := jsoniter.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addr+"/dig", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	data, err = s.processResponse(req)
	if err != nil {
		return nil, fmt.Errorf("process response: %w", err)
	}

	if data == nil {
		return []string{}, nil
	}

	var treasures []string
	if err := jsoniter.Unmarshal(data, &treasures); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return treasures, nil
}

func (s *Server) cash(ctx context.Context, treasureID string) ([]int, error) {
	data, err := jsoniter.Marshal(treasureID)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addr+"/cash", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	data, err = s.processResponse(req)
	if err != nil {
		return nil, fmt.Errorf("process response: %w", err)
	}

	var coins []int
	if err := jsoniter.Unmarshal(data, &coins); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return coins, nil
}
