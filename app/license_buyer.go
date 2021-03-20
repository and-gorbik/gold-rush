package app

// import (
// 	"log"
// 	"sync"
// 	"time"

// 	"gold-rush/models"
// )

// const (
// 	getCoinTimeout = 1 * time.Millisecond
// )

// type LicenseBuyer struct {
// 	bestPrice  price
// 	curPrice	price
// 	mx	sync.Mutex
// 	curPayment []int
// 	provider   goldRushServer
// }

// type price struct {
// 	Cost     int
// 	Value    int
// 	Duration time.Duration
// }

// // Run maintains a whole process of licenses buying
// func (b *LicenseBuyer) Run(done <-chan struct{}, coins <-chan int, workers int) <-chan models.License {
// 	licenses := make(chan models.License, maxActiveLicenses)

// 	b.calculateBestPrice(workers)

// 	for i := 0; i < workers; i++ {
// 		go b.buyer(done, coins, licenses)
// 	}

// 	return licenses
// }

// func (b *LicenseBuyer) calculateBestPrice(coins <-chan int, licenses chan<- models.License, workers int) {
// 	done := make(chan struct{})
// 	wg := new(sync.WaitGroup)
// 	wg.Add(workers)

// 	for i := 0; i < workers; i++ {
// 		go func() {
// 			defer wg.Done()
// 			b.buyer(done, coins, licenses)
// 		}()
// 	}
// }

// func (b *LicenseBuyer) buyer(done <-chan struct{}, coins <-chan int, licenses chan<- models.License) {
// 	for {
// 		payment := b.curPayment
// 		if payment == nil {
// 			payment = make([]int, 0, b.bestPrice.Cost)
// 		}

// 		for {
// 			select {
// 			case coin := <-coins:
// 				payment = append(payment, coin)
// 				if len(payment) != cap(payment) {
// 					continue
// 				}
// 				// buy license

// 			case <-time.After(getCoinTimeout):
// 				// buy free license
// 				b.curPayment = payment
// 				payment = []int{}

// 			case <-done:
// 				return
// 			}

// 			break
// 		}

// 		for {
// 			license, err := b.provider.BuyLicense(models.PaymentForLicense(payment))
// 			if err != nil {
// 				if _, isBusiness := readError(err); isBusiness {
// 					log.Fatal(err)
// 				}

// 				continue
// 			}

// 			licenses <- license
// 			break
// 		}
// 	}
// }
