package explorers

import (
	"log"
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
	areasCount := 2500

	q := pushAllConcurrently(testmocks.ExploredAreasGenerator(maxArea, areaSize))

	popConcurrently := func() {
		go func() {
			for i := 0; i < areasCount; i++ {
				log.Println(q.PopOrWait())
			}
		}()
	}

	popConcurrently()
	popConcurrently()
	popConcurrently()

	<-time.After(1 * time.Second)
	assert.Equal(t, q.Length(), 0)
}

func Test_PopOrWait(t *testing.T) {
	q := NewAreaQueue()

	alwaysBlocking := func() <-chan struct{} {
		res := make(chan struct{})
		go func() {
			defer close(res)
			q.Push(models.ExploredArea{})
			q.PopOrWait()
			q.PopOrWait() // alway wait here
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
