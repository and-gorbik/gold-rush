package exchangers

import (
	"time"

	"gold-rush/config"
	"gold-rush/models"
	"gold-rush/server"
)

type TreasuresExchanger struct {
	provider provider
	done     chan struct{}
}

func NewTreasuresExchanger(cfg config.Entity, treasures <-chan []string, coins chan<- int) *TreasuresExchanger {
	t := &TreasuresExchanger{
		provider: server.NewExchangerProvider(cfg.Client),
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
	retryDur := time.Millisecond
	for {
		payment, err := t.provider.ExchangeTreasure(models.Treasure(treasure))
		if err == nil {
			for _, coin := range payment {
				coins <- coin
			}

			return
		}

		<-time.After(retryDur)
		retryDur *= 2
	}
}

func (t *TreasuresExchanger) IncreaseWorkers(count int) {
}

func (t *TreasuresExchanger) ReduceWorkers(count int) {

}

func (t *TreasuresExchanger) Stop() {

}
