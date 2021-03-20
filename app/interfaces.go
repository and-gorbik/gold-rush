package app

import (
	"gold-rush/models"
)

type areaQueue interface {
	Push(area models.ExploredArea)
	PopOrWait() models.ExploredArea
}

type areaExplorer interface {
	Queue() areaQueue
	Stop()
}

type workersMaster interface {
	IncreaseWorkers(count int)
	ReduceWorkers(count int)
	Stop()
}

type licenser interface {
	workersMaster
	Lincenses() <-chan int
}

type treasuresEarner interface {
	workersMaster
	Treasures() <-chan string
}
