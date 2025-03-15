package service

import (
	"context"
	"log"

	"gold-rush/internal/server"
)

const (
	maxDepth        = 10
	maxTreasuresBuf = 1000
)

type DiggerService struct {
	server   handlers
	licenser licenser
}

type licenser interface {
	Done(trackID string)
}

func NewDiggerService(server handlers, licenser licenser) *DiggerService {
	return &DiggerService{
		server:   server,
		licenser: licenser,
	}
}

// один Digger копает одну шахту, потому что если их будет в шахте > 1, то будут тратиться
// лишние лицензии на уровни, где уже точно ничего нет
func (ds *DiggerService) DigTreasures(points <-chan Point, licenses <-chan License) <-chan string {
	treasures := make(chan string, maxTreasuresBuf)

	go func() {
		defer close(treasures)

		for point := range points {
			go func(p Point) {
				treasuresLeft := p.Treasures
				for depth := 1; depth <= maxDepth && treasuresLeft > 0; depth++ {
					license := <-licenses
					tt := ds.dig(license.ID, p.Pos.X, p.Pos.Y, depth)
					ds.licenser.Done(license.TrackID)
					treasuresLeft -= len(tt)

					for _, t := range tt {
						treasures <- t
					}
				}
			}(point)
		}
	}()

	return treasures
}

func (ds *DiggerService) dig(licenseID, x, y, z int) (treasures []string) {
	defer log.Println("digged treasures: ", treasures)
	return ds.server.Dig(context.Background(), server.DigParams{
		LicenseID: licenseID,
		PosX:      x,
		PosY:      y,
		Depth:     z,
	})
}
