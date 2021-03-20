package earners

import (
	"gold-rush/models"
)

type provider interface {
	Dig(params models.DigParams) (models.TreasuresList, error)
}
