package app

import (
	"log"
	"time"

	"gold-rush/models"
	"gold-rush/server"
)

const (
	maxArea             = 3500
	maxDepth            = 10
	totalTreasuresCount = 490000
	totalCoinsCount     = 23030000
)

var (
	areaSize          = 70
	areasCount        = maxArea * maxArea / areaSize / areaSize
	avgTreasuresCount = totalTreasuresCount / areasCount
)

var (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100

	exploreTimeout = 10 * time.Second
	digTimeout     = 1 * time.Second
	cashTimeout    = 2 * time.Second
	licenseTimeout = 3 * time.Second
	balanceTimeout = 4 * time.Second
	statusTimeout  = 1 * time.Second
)

type App struct {
	provider goldRushServer
}

func New() *App {
	p := &server.GoldRushServer{
		ExploreClient: buildHTTPClient(exploreTimeout),
		DigClient:     buildHTTPClient(digTimeout),
		CashClient:    buildHTTPClient(cashTimeout),
		LicenseClient: buildHTTPClient(licenseTimeout),
		BalanceClient: buildHTTPClient(balanceTimeout),
		StatusClient:  buildHTTPClient(statusTimeout),
	}

	return &App{provider: p}
}

func (a *App) explore(workers int) *AreaQueue {
	var posX, posY int
	points := make(chan Point, areasCount)

	go func() {
		for i := 0; i < areasCount; i++ {
			if posX >= maxArea {
				posX, posY = 0, posY+areaSize
			} else {
				posX = posX + areaSize
			}

			points <- Point{X: posX, Y: posY}
		}

		close(points)
	}()

	queue := new(AreaQueue)

	for i := 0; i < workers; i++ {
		go a.explorer(points, queue)
	}

	return queue
}

func (a *App) explorer(points <-chan Point, queue *AreaQueue) {
	for point := range points {
		area := models.Area{
			PosX:  point.X,
			PosY:  point.Y,
			SizeX: areaSize,
			SizeY: areaSize,
		}

		reconnectPeriod := time.Second
		for {
			exploredArea, err := a.provider.Explore(area)
			if err == nil {
				queue.Push(exploredArea)
				break
			}

			if msg, isBusiness := readError(err); isBusiness {
				log.Fatal(msg)
			}

			<-time.After(reconnectPeriod)
			reconnectPeriod *= 2
		}
	}
}

func (a *App) Run() {
	areas := a.explore(4)

	for i := 0; i < areasCount; {
		area := areas.Pop()
		if area.Amount < avgTreasuresCount {
			areas.Push(area)
			<-time.After(time.Millisecond)
			continue
		}

		i++
	}
}
