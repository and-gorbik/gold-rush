package explorers

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gold-rush/app/testmocks"
	"gold-rush/models"
)

func Test_PushPop_Concurrently(t *testing.T) {
	maxArea := 3500
	areaSize := 70
	areasCount := maxArea * maxArea / areaSize / areaSize

	q := pushAllConcurrently(testmocks.ExploredAreasGenerator(maxArea, areaSize))

	counter := 0
	for i := 0; i < areasCount; i++ {
		q.PopOrWait()
		counter++
	}

	assert.Equal(t, areasCount, counter)
}

func Test_PopOrWait(t *testing.T) {
	q := NewAreaQueue()

	alwaysBlocking := func() <-chan struct{} {
		res := make(chan struct{})
		go func() {
			defer close(res)
			q.Push(models.ExploredArea{})
			q.PopOrWait()
		}()

		return res
	}()

	select {
	case <-alwaysBlocking:
		t.Fatal("no wait")
	case <-time.After(time.Second):
	}
}

func pushAllConcurrently(areas <-chan models.ExploredArea) *AreaQueue {
	q := NewAreaQueue()

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for area := range areas {
				q.Push(area)
			}
		}()
	}

	return q
}
