package app

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	"gold-rush/app/earners"
	"gold-rush/app/exchangers"
	"gold-rush/app/explorers"
	"gold-rush/app/licensers"
	"gold-rush/config"
	"gold-rush/server"
)

const (
	TotalCoinsCount = 23030000
)

func Run() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := flag.String("path", path.Join(dir, "config.yaml"), "config path")
	flag.Parse()
	cfg := config.LoadFrom(*path)

	statusProvider := server.NewStatusProvider(cfg.StatusClient)
	retryDur := 100 * time.Millisecond
	for {
		if err = statusProvider.HealthCheck(); err == nil {
			break
		}
		<-time.After(retryDur)
		log.Printf("[status] dur: %v err: %v\n", retryDur, err)
		retryDur *= 2
	}

	coins := make(chan int, TotalCoinsCount)

	explorer := explorers.NewAreaExplorer(
		server.NewExplorerProvider(cfg.Explorer.Client),
		cfg.Explorer.Workers, cfg.App.AreaSize,
	)

	licenser := licensers.NewLicenser(
		server.NewLicenserProvider(cfg.Licenser.Client),
		cfg.Licenser.Workers, coins,
	)

	// количество воркеров игнорируется
	earner := earners.NewTreasuresEarner(
		server.NewEarnerProvider(cfg.Earner.Client),
		cfg.Earner.Workers, explorer.Queue(), licenser.Lincenses(),
	)

	// количество воркеров игнорируется
	_ = exchangers.NewTreasuresExchanger(
		server.NewExchangerProvider(cfg.Exchanger.Client),
		cfg.Exchanger.Workers, earner.Treasures(), coins,
	)

	<-time.After(10 * time.Minute)
}
