package licensers

import (
	"gold-rush/config"
	"gold-rush/server"
)

const (
	MaxActiveLicenses = 10
)

type Licenser struct {
	licenses chan int
	coins    chan int
	provider provider
}

func NewLicenser(cfg config.Entity, coins <-chan int) *Licenser {
	return &Licenser{
		provider: server.NewLicenserProvider(cfg.Client),
	}
}

func (l *Licenser) Lincenses() <-chan int {
	return l.licenses
}

func (l *Licenser) IncreaseWorkers(count int) {

}

func (l *Licenser) ReduceWorkers(count int) {

}

func (l *Licenser) Stop() {

}
