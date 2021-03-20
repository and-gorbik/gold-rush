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
	provider := server.New(cfg)

	coins := make(chan int, TotalCoinsCount)
	explorer := explorers.NewAreaExplorer(cfg.Explorer, provider, cfg.App.AreaSize)

	licenser := licensers.NewLicenser(cfg.Licenser, provider, coins)
	earner := earners.NewTreasuresEarner(cfg.Earner, provider, explorer.Queue(), licenser.Lincenses())
	_ = exchangers.NewTreasuresExchanger(cfg.Exchanger, provider, coins, earner.Treasures())

	<-time.After(10 * time.Minute)
}
