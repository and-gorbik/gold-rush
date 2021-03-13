package mocks

import (
	"errors"
	"time"

	"gold-rush/models"
)

type GoodProvider struct{}

func (GoodProvider) HealthCheck() error {
	return nil
}

func (GoodProvider) GetBalance() (models.Balance, error) {
	return models.Balance{}, nil
}

func (GoodProvider) GetLicenses() ([]models.License, error) {
	return []models.License{}, nil
}

func (GoodProvider) BuyLicense(payment models.PaymentForLicense) (models.License, error) {
	return models.License{}, nil
}

func (GoodProvider) Explore(area models.Area) (models.ExploredArea, error) {
	time.Sleep(time.Millisecond)
	return models.ExploredArea{}, nil
}

func (GoodProvider) Dig(params models.DigParams) (models.TreasuresList, error) {
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