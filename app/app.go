package app

import (
	"flag"
	"time"

	"gold-rush/app/earners"
	"gold-rush/app/exchangers"
	"gold-rush/app/explorers"
	"gold-rush/app/licensers"
	"gold-rush/config"
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
	coins := make(chan int, TotalCoinsCount)
	explorer := explorers.NewAreaExplorer(cfg.Explorer, cfg.App.AreaSize)
	licenser := licensers.NewLicenser(cfg.Licenser, coins)
	earner := earners.NewTreasuresEarner(cfg.Earner, explorer.Queue(), licenser.Lincenses())
	_ = exchangers.NewTreasuresExchanger(cfg.Exchanger, earner.Treasures(), coins)

	<-time.After(10 * time.Minute)
}
