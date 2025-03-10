package service

import (
	"context"

	"go.uber.org/zap"

	"gold-rush/internal/server"
)

const ()

type Service struct {
	handlers handlers
	log      *zap.Logger
}

func Init(h handlers, log *zap.Logger) *Service {
	return &Service{
		handlers: h,
		log:      log,
	}
}

type handlers interface {
	IssueLicense(ctx context.Context, coins []int) server.License
	ExploreArea(ctx context.Context, area server.Area) server.ExploredArea
	Dig(ctx context.Context, params server.DigParams) (treasures []string)
	Cash(ctx context.Context, treasureID string) (coins []int)
}
