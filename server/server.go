package server

import (
	"gold-rush/config"
	"gold-rush/infrastructure"
	"net/http"
	"time"
)

const (
	healthCheckURL = "/health-check"
	balanceURL     = "/balance"
	licensesURL    = "/licenses"
	exploreURL     = "/explore"
	digURL         = "/dig"
	cashURL        = "/cash"
)

type ExplorerProvider struct {
	client *http.Client
}

func NewExplorerProvider(cfg config.Client) *ExplorerProvider {
	return &ExplorerProvider{
		client: infrastructure.BuildHTTPClient(
			time.Duration(cfg.Timeout),
			cfg.MaxIdleConns,
			cfg.MaxConnsPerHost,
			cfg.MaxIdleConnsPerHost,
		),
	}
}

type EarnerProvider struct {
	client *http.Client
}

func NewEarnerProvider(cfg config.Client) *EarnerProvider {
	return &EarnerProvider{
		client: infrastructure.BuildHTTPClient(
			time.Duration(cfg.Timeout),
			cfg.MaxIdleConns,
			cfg.MaxConnsPerHost,
			cfg.MaxIdleConnsPerHost,
		),
	}
}

type ExchangerProvider struct {
	client *http.Client
}

func NewExchangerProvider(cfg config.Client) *ExchangerProvider {
	return &ExchangerProvider{
		client: infrastructure.BuildHTTPClient(
			time.Duration(cfg.Timeout),
			cfg.MaxIdleConns,
			cfg.MaxConnsPerHost,
			cfg.MaxIdleConnsPerHost,
		),
	}
}

type LicenserProvider struct {
	client *http.Client
}

func NewLicenserProvider(cfg config.Client) *LicenserProvider {
	return &LicenserProvider{
		client: infrastructure.BuildHTTPClient(
			time.Duration(cfg.Timeout),
			cfg.MaxIdleConns,
			cfg.MaxConnsPerHost,
			cfg.MaxIdleConnsPerHost,
		),
	}
}

type BalanceProvider struct {
	client *http.Client
}

func NewBalanceProvider(cfg config.Client) *BalanceProvider {
	return &BalanceProvider{
		client: infrastructure.BuildHTTPClient(
			time.Duration(cfg.Timeout),
			cfg.MaxIdleConns,
			cfg.MaxConnsPerHost,
			cfg.MaxIdleConnsPerHost,
		),
	}
}

type StatusProvider struct {
	client *http.Client
}

func NewStatusProvider(cfg config.Client) *StatusProvider {
	return &StatusProvider{
		client: infrastructure.BuildHTTPClient(
			time.Duration(cfg.Timeout),
			cfg.MaxIdleConns,
			cfg.MaxConnsPerHost,
			cfg.MaxIdleConnsPerHost,
		),
	}
}
