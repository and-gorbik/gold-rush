package licensers

import (
	"gold-rush/app/testmocks"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_getPayment(t *testing.T) {
	l := Licenser{payment: []int{}}
	var mx sync.Mutex

	for i := 0; i < 2; i++ {
		go func() {
			for {
				mx.Lock()
				l.payment = append(l.payment, randSlice()...)
				mx.Unlock()
				time.Sleep(time.Second)
			}
		}()
	}

	for {
		rand.Seed(time.Now().UnixNano())
		price := rand.Intn(100)
		payment := l.getPayment(&mx, price)
		log.Println("price: ", price, "payment: ", payment)
		mx.Lock()
		log.Println("cap: ", cap(l.payment))
		mx.Unlock()
		time.Sleep(time.Second)
	}
}

func randSlice() []int {
	rand.Seed(time.Now().UnixNano())
	size := rand.Intn(10)
	res := make([]int, size)
	for i := range res {
		res[i] = rand.Intn(20)
	}

	return res
}

func Test_calcBestPrice(t *testing.T) {
	l := Licenser{capacities: make([]int, 0, maxCapacity)}
	var mx sync.Mutex

	go func() {
		for {
			mx.Lock()
			if !l.isCalculated {
				capacity := testmocks.BuyLicense(len(l.capacities))
				l.capacities = append(l.capacities, capacity)
				log.Println("Price: ", len(l.capacities), "Capacity: ", capacity)
			}
			mx.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		price := l.calcBestPrice(&mx)
		mx.Lock()
		log.Println("best price: ", price)
		mx.Unlock()
		time.Sleep(time.Second)
	}
}

func Test_NewLicenser(t *testing.T) {
	l := NewLicenser(&testmocks.GoodProvider{}, 10, testmocks.CoinsGenerator())
	for lic := range l.Lincenses() {
		log.Println(lic)
	}
}
