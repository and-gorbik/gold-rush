package explorers

import (
	"gold-rush/models"
)

type provider interface {
	Explore(area models.Area) (models.ExploredArea, error)
}
