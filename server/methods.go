package server

import (
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"gold-rush/models"
)

func (p *StatusProvider) HealthCheck() (err error) {
	body, err := doRequest(p.client, http.MethodGet, healthCheckURL, nil)
	if err != nil {
		return
	}

	status := make(models.ServiceStatus)
	if err = jsoniter.Unmarshal(body, &status); err != nil {
		return
	}

	log.Println(status)
	return
}

func (p *BalanceProvider) GetBalance() (balance models.Balance, err error) {
	body, err := doRequest(p.client, http.MethodGet, balanceURL, nil)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &balance)
	return
}

func (p *LicenserProvider) GetLicenses() (licenses []models.LicenseFull, err error) {
	body, err := doRequest(p.client, http.MethodGet, licensesURL, nil)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &licenses)
	return
}

func (p *LicenserProvider) BuyLicense(payment models.PaymentForLicense) (license models.License, err error) {
	body, err := doRequest(p.client, http.MethodPost, licensesURL, payment)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &license)
	return
}

func (p *ExplorerProvider) Explore(area models.Area) (explored models.ExploredArea, err error) {
	body, err := doRequest(p.client, http.MethodPost, exploreURL, area)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &explored)
	return
}

func (p *EarnerProvider) Dig(params models.DigParams) (tl models.TreasuresList, err error) {
	body, err := doRequest(p.client, http.MethodPost, digURL, params)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &tl)
	return
}

func (p *ExchangerProvider) ExchangeTreasure(treasure models.Treasure) (pft models.PaymentForTreasure, err error) {
	body, err := doRequest(p.client, http.MethodPost, cashURL, treasure)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &pft)
	return
}
