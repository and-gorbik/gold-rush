package licensers

import (
	"gold-rush/config"
)

const (
	MaxActiveLicenses = 10
)

type Licenser struct {
	licenses chan int
	coins    chan int
	provider provider
}

func NewLicenser(cfg config.Entity, provider provider, coins <-chan int) *Licenser {
	return &Licenser{
		provider: provider,
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
