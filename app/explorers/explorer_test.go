package explorers

import (
	"log"
	"testing"

	"gold-rush/app/testmocks"
)

func Test_NewAreaExplorer(t *testing.T) {
	p := &testmocks.GoodProvider{}

	areas := NewAreaExplorer(p, 100, 100).Queue()
	for {
		a := areas.PopOrWait()
		log.Printf("X: %d Y: %d Amount: %d\n", a.Area.PosX, a.Area.PosY, a.Amount)
	}
}
