package app

import (
	"gold-rush/models"
)

type provider interface {
	HealthCheck() error
	GetBalance() (models.Balance, error)
	GetLicenses() ([]models.License, error)
	BuyLicense(payment models.PaymentForLicense) (models.License, error)
	Explore(area models.Area) (models.ExploredArea, error)
	Dig(params models.DigParams) (models.TreasuresList, error)
	ExchangeTreasure(treasure models.Treasure) (models.PaymentForTreasure, error)
}
