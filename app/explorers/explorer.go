package explorers

import (
	"log"
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

	// var posX, posY int
	points := make(chan Point, a.areasCount)

	go func() {
		defer close(points)

		for i := 0; i < MaxArea/a.areaSize; i++ {
			for j := 0; j < MaxArea/a.areaSize; j++ {
				points <- Point{X: i * a.areaSize, Y: j * a.areaSize}
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

		reconnectPeriod := 100 * time.Millisecond
		for {
			exploredArea, err := a.provider.Explore(area)
			if err == nil {
				queue.Push(exploredArea)
				break
			}

			<-time.After(reconnectPeriod)
			log.Printf("[explorer] dur: %v err: %v\n", reconnectPeriod, err)
			reconnectPeriod *= 2
		}
	}
}
