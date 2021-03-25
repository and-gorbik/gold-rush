package explorers

import (
	"time"

	"gold-rush/models"
)

const (
	MaxArea             = 3500
	TotalTreasuresCount = 490000
)

type AreaExplorer struct {
	areasCount        int
	avgTreasuresCount int
	areaSize          int
	areas             *AreaQueue
	provider          provider
}

func NewAreaExplorer(provider provider, workers int, areaSize int) *AreaExplorer {
	areaCount := MaxArea * MaxArea / areaSize / areaSize

	a := &AreaExplorer{
		areaSize:          areaSize,
		areasCount:        areaCount,
		avgTreasuresCount: TotalTreasuresCount / areaCount,
		provider:          provider,
	}

	a.areas = a.explore(workers)

	return a
}

func (a *AreaExplorer) Queue() *AreaQueue {
	return a.areas
}

func (a *AreaExplorer) Stop() {

}

func (a *AreaExplorer) explore(workers int) *AreaQueue {

	var posX, posY int
	points := make(chan Point, a.areasCount)

	go func() {
		defer close(points)

		for i := 0; i < a.areasCount; i++ {
			points <- Point{X: posX, Y: posY}

			if posX >= MaxArea {
				posX, posY = 0, posY+a.areaSize
			} else {
				posX = posX + a.areaSize
			}
		}
	}()

	queue := NewAreaQueue()

	for i := 0; i < workers; i++ {
		go a.explorer(points, queue)
	}

	return queue
}

func (a *AreaExplorer) explorer(points <-chan Point, queue *AreaQueue) {
	for point := range points {
		area := models.Area{
			PosX:  point.X,
			PosY:  point.Y,
			SizeX: a.areaSize,
			SizeY: a.areaSize,
		}

		reconnectPeriod := 10 * time.Second
		for {
			exploredArea, err := a.provider.Explore(area)
			if err == nil {
				queue.Push(exploredArea)
				break
			}

			<-time.After(reconnectPeriod)
			reconnectPeriod *= 2
		}
	}
}
