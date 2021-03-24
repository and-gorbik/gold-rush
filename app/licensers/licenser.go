package licensers

import (
	"sync"
	"time"

	"gold-rush/models"
)

const (
	MaxActiveLicenses = 10
)

var (
	getCoinTimeout = 1 * time.Millisecond
	maxCapacity    = 100
)

type Licenser struct {
	licenses             chan int
	licensesFromProvider chan models.License
	coins                chan int
	provider             provider
	workersCount         int
	payment              []int

	capacities   []int // prices[cost] = license capacity
	bestPrice    int
	isCalculated bool
}

func NewLicenser(provider provider, workers int, coins <-chan int) *Licenser {
	licenser := &Licenser{
		licenses:             make(chan int, 100),
		licensesFromProvider: make(chan models.License, MaxActiveLicenses),
		provider:             provider,
		capacities:           make([]int, 0, maxCapacity),
		bestPrice:            1,
		workersCount:         workers,
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
		go l.licenser(&mx, coins)
	}
}

// 1. Посчитать лучшую цену
// 2. Набрать необходимое количество coin'ов в payment
// 3. Если лишних coin'ов в канале нет (timeout), начать покупку бесплатной лицензии
// 4. Иначе, если сумма набрана, купить платную лицензию
func (l *Licenser) licenser(mx *sync.Mutex, coins <-chan int) {
	for {
		price := l.calcBestPrice(mx)
		l.buyLicense(mx, l.getPayment(mx, price))

		for {
			select {
			case <-time.After(getCoinTimeout):
			case coin := <-coins:
				mx.Lock()
				l.payment = append(l.payment, coin)
				if len(l.payment) != price {
					mx.Unlock()
					continue
				}
				mx.Unlock()
			}

			// timeout or fulled payment
			break
		}
	}
}

func (l *Licenser) buyLicense(mx *sync.Mutex, payment []int) {
	retryDur := 10 * time.Millisecond
	for {
		license, err := l.provider.BuyLicense(payment)
		if err == nil {
			// prices statistics
			mx.Lock()
			if !l.isCalculated {
				l.capacities = append(l.capacities, license.DigAllowed)
			}
			mx.Unlock()

			l.licensesFromProvider <- license
			return
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

	if len(l.capacities) >= maxCapacity { // capacities увеличивается воркерами
		bestCapsPerPrice := float32(0)
		bestPrice := 0
		for price, capacity := range l.capacities {
			// при одинаковом bestCapsPerPrice, выгодно взять больший price,
			// чтобы делать меньше запросов
			kf := float32(capacity / (price + 1))
			if kf >= bestCapsPerPrice {
				bestCapsPerPrice = kf
				bestPrice = price
			}
		}

		l.bestPrice = bestPrice
		l.isCalculated = true
		return bestPrice
	}

	l.bestPrice++
	return l.bestPrice
}

func (l *Licenser) getPayment(mx *sync.Mutex, count int) []int {
	mx.Lock()
	defer mx.Unlock()

	if len(l.payment) < count {
		return []int{}
	}

	result := make([]int, count)
	copy(result, l.payment[len(l.payment)-count:])
	l.payment = l.payment[:len(l.payment)-count]

	return result
}

func (l *Licenser) IncreaseWorkers(count int) {

}

func (l *Licenser) ReduceWorkers(count int) {

}

func (l *Licenser) Stop() {

}
