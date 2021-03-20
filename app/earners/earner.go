package earners

import (
	"gold-rush/config"
	"gold-rush/models"
)

const (
	TotalTreasuresCount = 490000
	MaxDepth            = 10
)

type TreasuresEarner struct {
	treasures chan string
	areas     queue
	provider  provider
}

type queue interface {
	Push(area models.ExploredArea)
	PopOrWait() models.ExploredArea
}

func NewTreasuresEarner(cfg config.Entity, provider provider, areas queue, licenses <-chan int) *TreasuresEarner {
	return &TreasuresEarner{}
}

func (t *TreasuresEarner) Treasures() <-chan string {
	return t.treasures
}

func (t *TreasuresEarner) IncreaseWorkers(count int) {

}

func (t *TreasuresEarner) ReduceWorkers(count int) {

}

func (t *TreasuresEarner) Stop() {

}
