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
)

type handlers interface {
	IssueLicense(ctx context.Context, coins []int) server.License
	ExploreArea(ctx context.Context, area server.Area) server.ExploredArea
	Dig(ctx context.Context, params server.DigParams) (treasures []string)
	Cash(ctx context.Context, treasureID string) (coins []int)
}

func Run(h handlers) {
	points := NewExplorerService(h).ExplorePoints(explorers)

	ls := NewLicenceService(h)
	licenses := ls.StreamLicenses()

	ds := NewDiggerService(h, ls)
	treasures := ds.DigTreasures(points, licenses)

	cs := NewCashierService(h)
	cs.ExchangeForCash(cashiers, treasures)

	for {
		log.Println("Total coins: ", cs.TotalCoins())
		<-time.After(time.Second)
	}
}
