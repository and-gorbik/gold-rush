package exchangers

import (
	"gold-rush/models"
)

type provider interface {
	ExchangeTreasure(treasure models.Treasure) (models.PaymentForTreasure, error)
}
