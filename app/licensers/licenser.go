package licensers

import (
	"sync"
	"time"

	"gold-rush/config"
	"gold-rush/models"
	"gold-rush/server"
)

const (
	MaxActiveLicenses = 10
)

var (
	getCoinTimeout = 1 * time.Millisecond
)

type Licenser struct {
	licenses             chan int
	licensesFromProvider chan models.License
	coins                chan int
	provider             provider
	workersCount         int

	prices       map[int]int // [cost]count
	bestPrice    int
	isCalculated bool
}

func NewLicenser(cfg config.Entity, coins <-chan int) *Licenser {
	licenser := &Licenser{
		licenses:             make(chan int, 100),
		licensesFromProvider: make(chan models.License, MaxActiveLicenses),
		provider:             server.NewLicenserProvider(cfg.Client),
		prices:               make(map[int]int),
		bestPrice:            1,
		workersCount:         cfg.Workers,
	}

	go licenser.run(coins)

	return licenser
}

func (l *Licenser) Lincenses() <-chan int {
	return l.licenses
}

func (l *Licenser) run(coins <-chan int) {
	go func() {
		for lic := range l.licensesFromProvider {
			for i := 0; i < lic.DigAllowed; i++ {
				l.licenses <- lic.ID
			}
		}
	}()

	var mx sync.Mutex
	for i := 0; i < l.workersCount; i++ {
		go l.buy(&mx, coins)
	}
}

func (l *Licenser) buy(mx *sync.Mutex, coins <-chan int) {
	// 1. Посчитать лучшую цену
	// 2. Набрать необходимое количество coin'ов в payment
	// 3. Если лишних coin'ов в канале нет (timeout), начать покупку бесплатной лицензии
	// 4. Иначе, если сумма набрана, купить платную лицензию

	price := l.calcBestPrice(mx)
	payment := make([]int, 0, price)

	for {
		select {
		case <-time.After(getCoinTimeout):
		case coin := <-coins:
			payment = append(payment, coin)
			if len(payment) != cap(payment) {
				continue
			}

			// buy license
		}

		// timeout or fulled payment
		break
	}

	retryDur := 10 * time.Millisecond
	for {
		license, err := l.provider.BuyLicense(payment)
		if err == nil {
			// prices statistics
			mx.Lock()
			if !l.isCalculated {
				l.prices[len(payment)] = license.DigAllowed
			}
			mx.Unlock()

			l.licensesFromProvider <- license
		}
		<-time.After(retryDur)
		retryDur *= 2
	}
}

func (l *Licenser) calcBestPrice(mx *sync.Mutex) int {
	mx.Lock()
	defer mx.Unlock()

	if l.isCalculated {
		return l.bestPrice
	}

	if len(l.prices) == 100 { // prices увеличивается воркерами
		bestCountsPerCost := float32(0)
		bestCost := 0
		for cost, count := range l.prices {
			if float32(count/cost) > bestCountsPerCost {
				bestCountsPerCost = float32(count / cost)
				bestCost = cost
			}
		}

		l.bestPrice = bestCost
		l.isCalculated = true
		return bestCost
	}

	return l.bestPrice + 1
}

func (l *Licenser) IncreaseWorkers(count int) {

}

func (l *Licenser) ReduceWorkers(count int) {

}

func (l *Licenser) Stop() {

}
