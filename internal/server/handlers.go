package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func (s *Server) IssueLicense(ctx context.Context, coins []int) License {
	var license *License
	var err error

	s.retryWithExponentialBackoff(func() error {
		license, err = s.issueLicense(ctx, coins)
		return err
	})

	if err != nil {
		s.log.Fatal("issue license failed unexpectedly", zap.Error(err))
	}

	return *license
}

func (s *Server) ExploreArea(ctx context.Context, area Area) ExploredArea {
	var ea *ExploredArea
	var err error

	s.retryWithExponentialBackoff(func() error {
		ea, err = s.exploreArea(ctx, area)
		return err
	})

	if err != nil {
		s.log.Fatal("explore area failed unexpectedly", zap.Error(err))
	}

	return *ea
}

func (s *Server) Dig(ctx context.Context, params DigParams) []string {
	var treasures []string
	var err error

	s.retryWithExponentialBackoff(func() error {
		treasures, err = s.dig(ctx, params)
		return err
	})

	if err != nil {
		s.log.Fatal("dig failed unexpectedly", zap.Error(err))
	}

	return treasures
}

func (s *Server) Cash(ctx context.Context, treasureID string) []int {
	var coins []int
	var err error

	s.retryWithExponentialBackoff(func() error {
		coins, err = s.cash(ctx, treasureID)
		return err
	})

	if err != nil {
		s.log.Fatal("cash failed unexpectedly", zap.Error(err))
	}

	return coins
}

func (s *Server) issueLicense(ctx context.Context, coins []int) (*License, error) {
	data, err := jsoniter.Marshal(coins)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/licenses", bytes.NewReader(data))
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/explore", bytes.NewReader(data))
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/dig", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	data, err = s.processResponse(req)
	if err != nil {
		return nil, fmt.Errorf("process response: %w", err)
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/cash", bytes.NewReader(data))
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
