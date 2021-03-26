package earners

import (
	"log"
	"sync"
	"time"

	"gold-rush/app/explorers"
	"gold-rush/infrastructure"
	"gold-rush/models"
)

const (
	MaxDepth = 10

	noSuchLicense = 403
	noTreasure    = 404

	wrongCoordinates = 1000
	wrongDepth       = 1001
)

var (
	startDurationTime = 100 * time.Millisecond
)

// point represents a unit of digging for worker
type point struct {
	X int
	Y int
}

type TreasuresEarner struct {
	treasures chan []string
	provider  provider
	workers   int
}

type queue interface {
	Push(area models.ExploredArea)
	PopOrWait() models.ExploredArea
}

func NewTreasuresEarner(provider provider, workers int, areas queue, licenses <-chan int) *TreasuresEarner {
	e := &TreasuresEarner{
		provider:  provider,
		treasures: make(chan []string, explorers.TotalTreasuresCount),
		workers:   workers,
	}

	go e.run(areas, licenses)

	return e
}

func (t *TreasuresEarner) Treasures() <-chan []string {
	return t.treasures
}

// 1. получает новый участок из очереди
// 2. распределяет его между воркерами
// 3. воркер, ответственный за копание в точке (x,y), копает все 10 уровней
// 4. если все клады в области найдены, воркеры завершаются, мастер переходит к следующей области
func (t *TreasuresEarner) run(areas queue, licenses <-chan int) {
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
				go t.dig(done, wg, &mx, &earnedCount, licenses, point{
					X: ea.Area.PosX + i,
					Y: ea.Area.PosY + j,
				})
			}
		}

		wg.Wait()
	}
}

// Когда находит сокровище, добавляет его в канал и увеличивает общую для всех воркеров переменную.
// По этой переменной мастер узнает, когда все сокровища на территории будут собраны
// и завершит всех воркеров, чтобы они не делали лишней работы.
func (t *TreasuresEarner) dig(done <-chan struct{}, wg *sync.WaitGroup, mx sync.Locker, counter *int, licenses <-chan int, p point) {
	defer wg.Done()

	for z := 1; z <= 10; {
		license := <-licenses
		retryDur := startDurationTime

		for {
			select {
			case <-done:
				return
			default:
			}

			treasures, err := t.provider.Dig(models.DigParams{
				LicenseID: license,
				PosX:      p.X,
				PosY:      p.Y,
				Depth:     z,
			})
			if err == nil {
				t.treasures <- treasures
				z++
				mx.Lock()
				*counter++
				mx.Unlock()
				break
			}

			e, ok := err.(infrastructure.ProviderError)
			if !ok {
				// repeat request with same params
				continue
			}

			if e.StatusCode >= 500 {
				// repeat request with same params after increased timeout
				log.Printf("[dig] dur: %v err: %v\n", retryDur, e)
				<-time.After(retryDur)
				retryDur *= 2
				continue
			}

			if e.Code == noSuchLicense {
				// repeat request with other license id
				log.Printf("[dig]: license: %d err: %v\n", license, e)
				break
			}

			if e.Code == noTreasure {
				z++
				break
			}

			log.Printf("[dig]: %v\n", e)
		}
	}
}

// 1. получает новый участок из очереди
// 2. получает максимальное количество активных лицензий
// 3. запускает n воркеров, каждому дается 10 возможностей копать, каждый копает свою точку
// 4. если все клады в области найдены, берется следующий участок из очереди
// func altRun(areas queue, workers int, licenses <-chan models.License) {
// 	activeLicenses := 0
// 	for {
// 		area := areas.PopOrWait()
// 		treasuresFound := 0
// 		var mx sync.Mutex

// 		for
// 	}
// }

// func alt_run(areas queue, workers int, licenses <-chan models.License) {
// 	for point := range pointGenerator(areas.PopOrWait()) {
// 		dig(point, actions)
// 	}
// }

// func actionGenerator(licenses <-chan int) (action [10]int, activeLicenses int) {
// 	licMap := make(map[int]struct{})
// 	for i := 0; i < len(action); i++ {
// 		id := <-licenses
// 		if _, ok := licMap[id]; !ok {
// 			licMap[id] = struct{}{}
// 		}
// 	}
// }

// func pointGenerator(area models.ExploredArea) <-chan point {
// 	points := make(chan point)
// 	go func() {
// 		defer close(points)
// 		for i := 0; i < area.Area.SizeX; i++ {
// 			for j := 0; j < area.Area.SizeY; j++ {
// 				points <- point{
// 					X: area.Area.PosX + i,
// 					Y: area.Area.PosY + j,
// 				}
// 			}
// 		}
// 	}()

// 	return points
// }
