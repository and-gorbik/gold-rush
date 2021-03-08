package provider

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

type Provider struct {
	client *http.Client
}
