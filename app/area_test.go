package app

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gold-rush/app/mocks"
	"gold-rush/models"
)

func Test_PushPop_Concurrently(t *testing.T) {
	q := pushAllConcurrently(mocks.ExploredAreasGenerator(maxArea, areaSize))

	counter := 0
	for i := 0; i < areasCount; i++ {
		q.Pop()
		counter++
	}

	assert.Equal(t, areasCount, counter)
}

func Test_PopWait(t *testing.T) {
	q := NewAreaQueue()

	alwaysBlocking := func() <-chan struct{} {
		res := make(chan struct{})
		go func() {
			defer close(res)
			q.Push(models.ExploredArea{})
			q.Pop()
			q.Pop()
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
