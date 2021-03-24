package app

import (
	"flag"
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

var (
	cfg *config.Config
)

func init() {
	path := flag.String("path", "config.yaml", "config path")
	flag.Parse()
	cfg = config.LoadFrom(*path)
}

func Run() {
	statusProvider := server.NewStatusProvider(cfg.StatusClient)
	retryDur := 10 * time.Millisecond
	for {
		if err := statusProvider.HealthCheck(); err == nil {
			break
		}
		<-time.After(retryDur)
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
