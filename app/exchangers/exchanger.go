package exchangers

import (
	"log"
	"time"

	"gold-rush/models"
)

type TreasuresExchanger struct {
	provider provider
	workers  int
	done     chan struct{}
}

func NewTreasuresExchanger(provider provider, workers int, treasures <-chan []string, coins chan<- int) *TreasuresExchanger {
	t := &TreasuresExchanger{
		provider: provider,
		workers:  workers,
		done:     make(chan struct{}),
	}

	go t.run(treasures, coins)
	return t
}

// на каждый сундук создается воркер, ответственный за обмен этого сундука на деньги
func (t *TreasuresExchanger) run(treasuresChan <-chan []string, coins chan<- int) {
	for treasures := range treasuresChan {
		for _, treasure := range treasures {
			go t.cash(treasure, coins)
		}
	}
}

func (t *TreasuresExchanger) cash(treasure string, coins chan<- int) {
	retryDur := 10 * time.Millisecond
	for {
		payment, err := t.provider.ExchangeTreasure(models.Treasure(treasure))
		if err == nil {
			for _, coin := range payment {
				coins <- coin
			}

			return
		}

		<-time.After(retryDur)
		log.Printf("[cash] dur: %v err: %v\n", retryDur, err)
		retryDur *= 2
	}
}

func (t *TreasuresExchanger) IncreaseWorkers(count int) {
}

func (t *TreasuresExchanger) ReduceWorkers(count int) {

}

func (t *TreasuresExchanger) Stop() {

}
