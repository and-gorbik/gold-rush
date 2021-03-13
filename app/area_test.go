package app

import (
	"runtime"
	"testing"

	"gold-rush/app/mocks"
	"gold-rush/models"
)

func Test_PushPop_Concurrently(t *testing.T) {
	q := pushAllConcurrently(mocks.ExploredAreasGenerator(maxArea, areaSize))

	popAllConcurrently(q)
}

func pushAllConcurrently(areas <-chan models.ExploredArea) *AreaQueue {
	q := NewAreaQueue()

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			for {
				area, ok := <-areas
				if !ok {
					return
				}

				q.Push(area)
			}
		}(i)
	}

	return q
}

func popAllConcurrently(q *AreaQueue) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			for {
				q.Pop()
			}
		}(i)
	}
}
