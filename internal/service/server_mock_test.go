package service

import (
	"context"
	"math/rand"
	"sync"

	"github.com/google/uuid"

	"gold-rush/internal/server"
)

type MockServer struct {
	mu             sync.Mutex
	licenses       []int8
	curLicenseIdx  int
	area           [100][100][10]bool
	costOfTreasure map[string]int
}

func setupNewMockServer() *MockServer {
	area := [100][100][10]bool{}
	total := 0
	for y := 0; y < len(area); y++ {
		for x := 0; x < len(area[0]); x++ {
			for z := 0; z < len(area[0][0]); z++ {
				if rand.Intn(100)%2 == 0 {
					area[x][y][z] = true
					total++
				}
			}
		}
	}

	costOfTreasure := make(map[string]int, total)
	for i := 0; i < total; i++ {
		costOfTreasure[uuid.NewString()] = rand.Intn(10)
	}

	licenses := make([]int8, 0, total/3+1)
	for i := 0; i < total/3+1; i++ {
		licenses = append(licenses, 3)
	}

	return &MockServer{
		costOfTreasure: costOfTreasure,
		licenses:       licenses,
		area:           area,
	}
}

func (ms *MockServer) IssueLicense(ctx context.Context, _ []int) server.License {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.licenses[ms.curLicenseIdx] == 0 {
		ms.curLicenseIdx++
	}

	ms.licenses[ms.curLicenseIdx]--

	return server.License{
		ID:         ms.curLicenseIdx,
		DigAllowed: int(ms.licenses[ms.curLicenseIdx]),
		DigUsed:    3 - int(ms.licenses[ms.curLicenseIdx]),
	}
}

func (ms *MockServer) ExploreArea(ctx context.Context, a server.Area) server.ExploredArea {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// for size = 1
	found := 0
	for _, isTreasure := range ms.area[a.PosX][a.PosY] {
		if isTreasure {
			found++
		}
	}

	return server.ExploredArea{Area: a, Amount: found}
}

func (ms *MockServer) Dig(ctx context.Context, p server.DigParams) (treasures []string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if !ms.area[p.PosX][p.PosY][p.Depth-1] {
		return []string{}
	}

	count := rand.Intn(3) + 1
	for i := 0; i < count; i++ {
		t := uuid.NewString()
		treasures = append(treasures, t)
		ms.costOfTreasure[t] = i
	}

	return treasures
}

func (ms *MockServer) Cash(ctx context.Context, treasureID string) (coins []int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	cost := ms.costOfTreasure[treasureID]

	for i := 0; i < cost; i++ {
		coins = append(coins, rand.Int())
	}

	return coins
}
