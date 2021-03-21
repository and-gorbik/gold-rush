package earners

import (
	"sync"
	"time"

	"gold-rush/config"
	"gold-rush/models"
)

const (
	TotalTreasuresCount = 490000
	MaxDepth            = 10
)

// point represents a unit of digging for worker
type point struct {
	X int
	Y int
	Z int
}

// type worker struct {
// 	done   chan struct{}
// 	result chan string
// }

type TreasuresEarner struct {
	// workers   map[string]worker
	// mx        sync.RWMutex
	// done      chan struct{}
	licenses  <-chan int
	treasures chan []string
	provider  provider
}

type queue interface {
	Push(area models.ExploredArea)
	PopOrWait() models.ExploredArea
}

func NewTreasuresEarner(cfg config.Entity, provider provider, areas queue, licenses <-chan int) *TreasuresEarner {
	e := &TreasuresEarner{
		provider:  provider,
		treasures: make(chan []string, TotalTreasuresCount),
		licenses:  licenses,
	}

	go e.run(areas)

	return e
}

func (t *TreasuresEarner) Treasures() <-chan []string {
	return t.treasures
}

func (t *TreasuresEarner) IncreaseWorkers(count int) {

}

func (t *TreasuresEarner) ReduceWorkers(count int) {

}

func (t *TreasuresEarner) Stop() {

}

// 1. получает новый участок из очереди
// 2. распределяет его между воркерами
// 3. воркер, ответственный за копание в точке (x,y), копает все 10 уровней
// 4. если все клады в области найдены, воркеры завершаются, мастер переходит к следующей области

func (t *TreasuresEarner) run(areas queue) {
	for {
		ea := areas.PopOrWait()
		var mx sync.RWMutex
		earnedCount := 0
		done := make(chan struct{})
		go func() {
			for {
				<-time.After(time.Second)
				mx.RLock()
				if earnedCount >= ea.Amount {
					mx.RUnlock()
					close(done)
					return
				}
				mx.RUnlock()
			}
		}()

		wg := new(sync.WaitGroup)
		wg.Add(ea.Area.SizeX * ea.Area.SizeY)
		for i := 0; i < ea.Area.SizeX; i++ {
			for j := 0; j < ea.Area.SizeY; j++ {
				go t.dig(done, wg, &mx, &earnedCount, point{
					X: ea.Area.PosX + i,
					Y: ea.Area.PosY + j,
					Z: 0,
				})
			}
		}

		wg.Wait()
	}
}

// Когда находит сокровище, добавляет его в канал и увеличивает общую для всех воркеров переменную.
// По этой переменной мастер узнает, когда все сокровища на территории будут собраны
// и завершит всех воркеров, чтобы они не делали лишней работы.
func (t *TreasuresEarner) dig(done <-chan struct{}, wg *sync.WaitGroup, mx sync.Locker, counter *int, p point) {
	defer wg.Done()

	work := func() {
		for z := 0; z < 10; z++ {
			license := <-t.licenses
			p.Z = z + 1

			retryDur := 10 * time.Millisecond
			for {
				treasures, err := t.provider.Dig(models.DigParams{
					LicenseID: license,
					PosX:      p.X,
					PosY:      p.Y,
					Depth:     p.Z,
				})
				if err == nil {
					t.treasures <- treasures
					mx.Lock()
					*counter++
					mx.Unlock()
					break
				}

				<-time.After(retryDur)
				retryDur *= 2
			}
		}
	}

	for {
		select {
		case <-done:
			return
		default:
		}

		work()
	}
}
