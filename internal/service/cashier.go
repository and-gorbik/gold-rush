package service

import (
	"context"
	"log"
	"sync"

	"go.uber.org/atomic"
)

type CashierService struct {
	server handlers
	total  *atomic.Int32
}

func NewCashierService(server handlers) *CashierService {
	return &CashierService{
		server: server,
		total:  atomic.NewInt32(0),
	}
}

func (cs *CashierService) TotalCoins() int {
	return int(cs.total.Load())
}

func (cs *CashierService) ExchangeForCash(workers int, treasures <-chan string) {
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for t := range treasures {
				// save them for buying licenses further
				coins := cs.cash(t)
				cs.total.Add(int32(len(coins)))
			}
		}()
	}

	wg.Wait()
}

func (cs *CashierService) cash(treasure string) (coins []int) {
	defer log.Println("got coins: ", coins)
	return cs.server.Cash(context.Background(), treasure)
}
