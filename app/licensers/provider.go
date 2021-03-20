package licensers

import (
	"gold-rush/models"
)

type provider interface {
	GetLicenses() ([]models.LicenseFull, error)
	BuyLicense(payment models.PaymentForLicense) (models.License, error)
}
