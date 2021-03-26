package server

import (
	"encoding/json"
	"log"
	"net/http"

	// jsoniter "github.com/json-iterator/go"

	"gold-rush/models"
)

func (p *StatusProvider) HealthCheck() (err error) {
	body, err := doRequest(p.client, http.MethodGet, healthCheckURL, nil)
	if err != nil {
		return
	}

	log.Println(string(body))
	return
}

func (p *BalanceProvider) GetBalance() (balance models.Balance, err error) {
	body, err := doRequest(p.client, http.MethodGet, balanceURL, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &balance); err != nil {
		return
	}

	return
}

func (p *LicenserProvider) GetLicenses() (licenses []models.LicenseFull, err error) {
	body, err := doRequest(p.client, http.MethodGet, licensesURL, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &licenses); err != nil {
		return
	}

	return
}

func (p *LicenserProvider) BuyLicense(payment models.PaymentForLicense) (license models.LicenseFull, err error) {
	body, err := doRequest(p.client, http.MethodPost, licensesURL, payment)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &license); err != nil {
		return
	}

	return
}

func (p *ExplorerProvider) Explore(area models.Area) (explored models.ExploredArea, err error) {
	body, err := doRequest(p.client, http.MethodPost, exploreURL, area)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &explored); err != nil {
		return
	}

	return
}

func (p *EarnerProvider) Dig(params models.DigParams) (tl models.TreasuresList, err error) {
	body, err := doRequest(p.client, http.MethodPost, digURL, params)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &tl); err != nil {
		return
	}

	return
}

func (p *ExchangerProvider) ExchangeTreasure(treasure models.Treasure) (pft models.PaymentForTreasure, err error) {
	body, err := doRequest(p.client, http.MethodPost, cashURL, treasure)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &pft); err != nil {
		return
	}

	return
}
