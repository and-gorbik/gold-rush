package service

import (
	"context"
	"log"
	"sync"

	"gold-rush/internal/server"
)

var (
	MaxArea = 3500
)

type ExplorerService struct {
	mu       sync.Mutex
	store    [][]int8
	explored chan Point
	server   handlers
}

type Point struct {
	Pos       Coord
	Treasures int
}

type Coord struct {
	X int
	Y int
}

func NewExplorerService(server handlers, pointsBuf int) *ExplorerService {
	store := make([][]int8, MaxArea)
	for i := 0; i < MaxArea; i++ {
		store[i] = make([]int8, MaxArea)
	}

	return &ExplorerService{
		server:   server,
		explored: make(chan Point, pointsBuf),
		store:    store,
	}
}

func (es *ExplorerService) ExplorePoints(ctx context.Context, workers int) <-chan Point {
	go func() {
		var wg sync.WaitGroup
		wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func(worker int) {
				defer wg.Done()

				// check that i is always less than maxArea
				for y := worker; y < MaxArea; y += workers {
					for x := 0; x < MaxArea; x++ {
						select {
						case <-ctx.Done():
							return
						default:
						}

						treasuresCount := es.explore(x, y, 1, 1)
						es.mu.Lock()
						es.store[y][x] = int8(treasuresCount)
						es.mu.Unlock()

						if treasuresCount > 0 {
							select {
							case <-ctx.Done():
								return
							case es.explored <- Point{Coord{x, y}, treasuresCount}:
							}
						}
					}
				}
			}(i)
		}

		wg.Wait()
		close(es.explored)
	}()

	return es.explored
}

func (es *ExplorerService) explore(x, y, sizeX, sizeY int) int {
	ea := es.server.ExploreArea(context.Background(), server.Area{
		PosX:  x,
		PosY:  y,
		SizeX: sizeX,
		SizeY: sizeY,
	})

	log.Println("explored area has treasures: ", ea)

	return ea.Amount
}
