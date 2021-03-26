package app

import (
	"log"
	"runtime"
	"testing"
	"time"

	"gold-rush/app/earners"
	"gold-rush/app/exchangers"
	"gold-rush/app/explorers"
	"gold-rush/app/licensers"
	"gold-rush/app/testmocks"
	"gold-rush/config"
)

func Test_Run_Success(t *testing.T) {
	p := testmocks.GoodProvider{}

	coins := make(chan int, TotalCoinsCount)
	explorer := explorers.NewAreaExplorer(p, 100, 70)
	licenser := licensers.NewLicenser(p, 10, coins)
	earner := earners.NewTreasuresEarner(p, 0, explorer.Queue(), licenser.Lincenses())
	_ = exchangers.NewTreasuresExchanger(p, 0, earner.Treasures(), coins)

	go func() {
		for {
			<-time.After(time.Second)
			log.Println("goroutine number: ", runtime.NumGoroutine())
		}
	}()

	<-time.After(10 * time.Minute)
}

func initTestConfig() *config.Config {
	return &config.Config{}
}
