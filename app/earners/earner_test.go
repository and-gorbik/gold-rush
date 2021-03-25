package earners

import (
	"log"
	"testing"

	"gold-rush/app/explorers"
	"gold-rush/app/testmocks"
)

func Test_NewTreasuresEarner(t *testing.T) {
	p := &testmocks.GoodProvider{}
	queue := explorers.NewAreaExplorer(p, 100, 70).Queue()

	earner := NewTreasuresEarner(p, 10, queue, testmocks.LicensesGenerator())
	for t := range earner.Treasures() {
		log.Printf("treasures: %v\n", t)
	}
}
