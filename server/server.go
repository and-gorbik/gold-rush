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
	client *http.Client
}
