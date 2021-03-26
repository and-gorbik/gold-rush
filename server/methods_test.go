package server

import (
	"gold-rush/config"
	"gold-rush/models"
	"testing"
	"time"
)

func Test_BuyLicense(t *testing.T) {
	cfg := config.Client{
		MaxIdleConns:        100,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
		Timeout:             time.Duration(time.Second * 3),
	}

	host = "127.0.0.1"
	port = "4010"
	schema = "http"

	p := NewLicenserProvider(cfg)

	payment := models.PaymentForLicense{
		0, 1, 2, 5, 10, 100500, 8,
	}

	_, err := p.BuyLicense(payment)
	if err != nil {
		t.Fatal(err)
	}
}
