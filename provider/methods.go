package provider

import (
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"gold-rush/models"
)

func (p *Provider) HealthCheck() (err error) {
	body, err := p.doRequest(http.MethodGet, healthCheckURL, nil)
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

func (p *Provider) GetBalance() (balance models.Balance, err error) {
	body, err := p.doRequest(http.MethodGet, balanceURL, nil)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &balance)
	return
}

func (p *Provider) GetLicenses() (licenses []models.License, err error) {
	body, err := p.doRequest(http.MethodGet, licensesURL, nil)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &licenses)
	return
}

func (p *Provider) BuyLicense(payment models.PaymentForLicense) (license models.License, err error) {
	body, err := p.doRequest(http.MethodPost, licensesURL, payment)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &license)
	return
}

func (p *Provider) Explore(area models.Area) (explored models.ExploredArea, err error) {
	body, err := p.doRequest(http.MethodPost, exploreURL, area)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &explored)
	return
}

func (p *Provider) Dig(params models.DigParams) (tl models.TreasuresList, err error) {
	body, err := p.doRequest(http.MethodPost, digURL, params)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &tl)
	return
}

func (p *Provider) ExchangeTreasure(treasure models.Treasure) (pft models.PaymentForTreasure, err error) {
	body, err := p.doRequest(http.MethodPost, cashURL, treasure)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &pft)
	return
}
