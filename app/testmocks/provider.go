package testmocks

import (
	"errors"
	"log"
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

func (GoodProvider) BuyLicense(payment models.PaymentForLicense) (models.LicenseFull, error) {
	defer log.Println("BuyLicense()")

	time.Sleep(time.Second)
	capacity := BuyLicense(len(payment))
	if len(payment) == 0 {
		capacity = 3
	}
	rand.Seed(time.Now().UnixNano())
	return models.LicenseFull{
		License: models.License{
			ID:         rand.Int(),
			DigAllowed: capacity,
		},
		DigUsed: 0,
	}, nil
}

func (GoodProvider) Explore(area models.Area) (models.ExploredArea, error) {
	defer log.Println("Explore()")

	time.Sleep(time.Second)
	rand.Seed(time.Now().UnixNano())
	return models.ExploredArea{
		Area:   area,
		Amount: rand.Intn(400),
	}, nil
}

func (GoodProvider) Dig(params models.DigParams) (models.TreasuresList, error) {
	defer log.Println("Dig()")

	time.Sleep(time.Second)
	return models.TreasuresList(GenTreasuresList()), nil
}

func (GoodProvider) ExchangeTreasure(treasure models.Treasure) (models.PaymentForTreasure, error) {
	defer log.Println("ExchangeTreasure()")

	return models.PaymentForTreasure(GetPaymentForTreasure()), nil
}

// BAD ONES

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
