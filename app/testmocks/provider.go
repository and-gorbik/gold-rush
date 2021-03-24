package testmocks

import (
	"errors"
	"math/rand"
	"time"

	"gold-rush/models"
)

type GoodProvider struct{}

func (GoodProvider) HealthCheck() error {
	time.Sleep(time.Second)
	return nil
}

func (GoodProvider) GetBalance() (models.Balance, error) {
	time.Sleep(time.Second)
	return models.Balance{}, nil
}

func (GoodProvider) GetLicenses() ([]models.LicenseFull, error) {
	time.Sleep(time.Second)
	return []models.LicenseFull{}, nil
}

func (GoodProvider) BuyLicense(payment models.PaymentForLicense) (models.License, error) {
	time.Sleep(time.Second)
	capacity := BuyLicense(len(payment))
	rand.Seed(time.Now().UnixNano())
	return models.License{
		ID:         rand.Int(),
		DigAllowed: capacity,
	}, nil
}

func (GoodProvider) Explore(area models.Area) (models.ExploredArea, error) {
	time.Sleep(time.Millisecond)
	return models.ExploredArea{}, nil
}

func (GoodProvider) Dig(params models.DigParams) (models.TreasuresList, error) {
	time.Sleep(time.Millisecond)
	return models.TreasuresList{}, nil
}

func (GoodProvider) ExchangeTreasure(treasure models.Treasure) (models.PaymentForTreasure, error) {
	return models.PaymentForTreasure{}, nil
}

type BadExplorer struct {
	GoodProvider
}

func (BadExplorer) Explore(area models.Area) (models.ExploredArea, error) {
	time.Sleep(time.Millisecond)
	return models.ExploredArea{}, errors.New("error")
}

type BadDigger struct {
	GoodProvider
}

func (BadDigger) Dig(params models.DigParams) (models.TreasuresList, error) {
	time.Sleep(time.Second)
	return nil, errors.New("error")
}
