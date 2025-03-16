package service

import (
	"context"
	"log"
	"time"

	"gold-rush/internal/server"
)

const (
	explorers = 10
	cashiers  = 10

	bufPoints    = 100
	bufLicenses  = 100
	bufTreasures = 1000

	periodTryIssueLicense = 100 * time.Millisecond
)

type handlers interface {
	IssueLicense(ctx context.Context, coins []int) server.License
	ExploreArea(ctx context.Context, area server.Area) server.ExploredArea
	Dig(ctx context.Context, params server.DigParams) (treasures []string)
	Cash(ctx context.Context, treasureID string) (coins []int)
}

func Run(ctx context.Context, h handlers) {
	points := NewExplorerService(h, bufPoints).ExplorePoints(ctx, explorers)

	ls := NewLicenceService(h, LicenserOpts{
		TryIssueLicensePeriod: periodTryIssueLicense,
		MaxLicencesBuf:        bufLicenses,
	})
	licenses := ls.StreamLicenses(ctx)

	ds := NewDiggerService(h, ls)
	treasures := ds.DigTreasures(ctx, points, licenses, bufTreasures)

	cs := NewCashierService(h)
	cs.ExchangeForCash(ctx, cashiers, treasures)

	for {
		log.Println("Total coins: ", cs.TotalCoins())
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			// instead of graceful shutdown
			<-time.After(time.Second * 3)
			return
		}
	}
}
