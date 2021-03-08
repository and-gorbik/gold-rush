package server

import (
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"gold-rush/models"
)

func (s *GoldRushServer) HealthCheck() (err error) {
	body, err := s.doRequest(http.MethodGet, healthCheckURL, nil)
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

func (s *GoldRushServer) GetBalance() (balance models.Balance, err error) {
	body, err := s.doRequest(http.MethodGet, balanceURL, nil)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &balance)
	return
}

func (s *GoldRushServer) GetLicenses() (licenses []models.License, err error) {
	body, err := s.doRequest(http.MethodGet, licensesURL, nil)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &licenses)
	return
}

func (s *GoldRushServer) BuyLicense(payment models.PaymentForLicense) (license models.License, err error) {
	body, err := s.doRequest(http.MethodPost, licensesURL, payment)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &license)
	return
}

func (s *GoldRushServer) Explore(area models.Area) (explored models.ExploredArea, err error) {
	body, err := s.doRequest(http.MethodPost, exploreURL, area)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &explored)
	return
}

func (s *GoldRushServer) Dig(params models.DigParams) (tl models.TreasuresList, err error) {
	body, err := s.doRequest(http.MethodPost, digURL, params)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &tl)
	return
}

func (s *GoldRushServer) ExchangeTreasure(treasure models.Treasure) (pft models.PaymentForTreasure, err error) {
	body, err := s.doRequest(http.MethodPost, cashURL, treasure)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(body, &pft)
	return
}
