package explorers

import (
	"log"
	"time"

	"gold-rush/app/earners"
	"gold-rush/config"
	"gold-rush/infrastructure"
	"gold-rush/models"
	"gold-rush/server"
)

const (
	MaxArea = 3500
)

type AreaExplorer struct {
	areasCount        int
	avgTreasuresCount int
	areaSize          int
	areas             *AreaQueue
	provider          provider
}

func NewAreaExplorer(cfg config.Entity, areaSize int) *AreaExplorer {
	areaCount := MaxArea * MaxArea / areaSize / areaSize

	a := &AreaExplorer{
		areaSize:          areaSize,
		areasCount:        areaCount,
		avgTreasuresCount: earners.TotalTreasuresCount / areaCount,
		provider:          server.NewExplorerProvider(cfg.Client),
	}

	a.areas = a.explore(cfg.Workers)

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

		reconnectPeriod := time.Second
		for {
			exploredArea, err := a.provider.Explore(area)
			if err == nil {
				queue.Push(exploredArea)
				break
			}

			if msg, isBusiness := infrastructure.ReadError(err); isBusiness {
				log.Fatal(msg)
			}

			<-time.After(reconnectPeriod)
			reconnectPeriod *= 2
		}
	}
}
