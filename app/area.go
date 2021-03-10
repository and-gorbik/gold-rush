package app

import (
	"container/heap"
	"sync"

	"gold-rush/models"
)

type Point struct {
	X int
	Y int
}

type AreaQueue struct {
	pq   priorityQueue
	len  int
	mx   sync.Mutex
	cond sync.Cond
}

func NewAreaQueue() *AreaQueue {
	a := &AreaQueue{
		pq: make(priorityQueue, areasCount),
	}

	heap.Init(&a.pq)

	return a
}

func (a *AreaQueue) Push(elem models.ExploredArea) {
	if a.len == 0 {
		a.cond.L.Lock()
		a.cond.Broadcast()
		a.cond.L.Unlock()
	}

	a.mx.Lock()
	heap.Push(&a.pq, elem)
	a.len++
	a.mx.Unlock()
}

func (a *AreaQueue) Pop() models.ExploredArea {
	if a.len == 0 {
		a.cond.L.Lock()
		a.cond.Wait()
		a.cond.L.Unlock()
	}

	a.mx.Lock()
	defer a.mx.Unlock()
	a.len--

	return heap.Pop(&a.pq).(models.ExploredArea)
}

// очередь с приоритетом
type priorityQueue []*queueItem

type queueItem struct {
	models.ExploredArea
	index int
}

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Amount > pq[j].Amount
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*queueItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
