package exchangers

import (
	"gold-rush/config"
)

type TreasuresExchanger struct {
	provider provider
}

func NewTreasuresExchanger(cfg config.Entity, provider provider, coins chan<- int, treasures <-chan []string) *TreasuresExchanger {
	return &TreasuresExchanger{}
}

func (t *TreasuresExchanger) IncreaseWorkers(count int) {

}

func (t *TreasuresExchanger) ReduceWorkers(count int) {

}

func (t *TreasuresExchanger) Stop() {

}
