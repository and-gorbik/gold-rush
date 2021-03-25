package testmocks

import (
	"math/rand"
	"time"
)

func LicensesGenerator() <-chan int {
	coins := make(chan int, 100)
	go func() {
		defer close(coins)
		for {
			rand.Seed(time.Now().UnixNano())
			coins <- rand.Int()
		}
	}()

	return coins
}
