package server

import (
	"net/http"
)

const (
	healthCheckURL = "/health-check"
	balanceURL     = "/balance"
	licensesURL    = "/licenses"
	exploreURL     = "/explore"
	digURL         = "/dig"
	cashURL        = "/cash"
)

type GoldRushServer struct {
	ExploreClient *http.Client
	DigClient     *http.Client
	CashClient    *http.Client
	LicenseClient *http.Client
	BalanceClient *http.Client
	StatusClient  *http.Client
}
