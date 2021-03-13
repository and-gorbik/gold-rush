package mocks

import (
	"math/rand"
	"time"

	"gold-rush/models"
)

var (
	maxRndAmount = 400
)

func ExploredAreasGenerator(maxArea, areaSize int) <-chan models.ExploredArea {
	areas := generateExploredAreas(maxArea, areaSize)

	result := make(chan models.ExploredArea, len(areas))
	go func() {
		defer close(result)
		for _, area := range areas {
			result <- area
		}
	}()

	return result
}

func generateExploredAreas(maxArea, areaSize int) []models.ExploredArea {
	areaCount := maxArea * maxArea / areaSize / areaSize
	areas := make([]models.ExploredArea, 0, areaCount)
	posX, posY := 0, 0
	for i := 0; i < areaCount; i++ {
		rand.Seed(time.Now().UnixNano())
		areas = append(areas, models.ExploredArea{
			Area: models.Area{
				PosX:  posX,
				PosY:  posY,
				SizeX: areaSize,
				SizeY: areaSize,
			},
			Amount: rand.Intn(maxRndAmount),
		})

		if posX >= maxArea {
			posX, posY = 0, posY+areaSize
		} else {
			posX = posX + areaSize
		}
	}

	return areas
}
