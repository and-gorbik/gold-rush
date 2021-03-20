package explorers

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
	mx   *sync.Mutex
	cond *sync.Cond
}

func NewAreaQueue() *AreaQueue {
	mx := new(sync.Mutex)

	a := &AreaQueue{
		pq:   make(priorityQueue, 0),
		mx:   mx,
		cond: sync.NewCond(mx),
	}

	heap.Init(&a.pq)

	return a
}

func (a *AreaQueue) Push(ea models.ExploredArea) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if a.len == 0 {
		a.cond.Broadcast()
	}

	heap.Push(&a.pq, &queueItem{ExploredArea: ea})
	a.len++
}

func (a *AreaQueue) PopOrWait() models.ExploredArea {
	a.mx.Lock()
	defer a.mx.Unlock()

	if a.len == 0 {
		a.cond.Wait()
	}

	a.len--

	return heap.Pop(&a.pq).(*queueItem).ExploredArea
}

func (a *AreaQueue) Length() int {
	a.mx.Lock()
	defer a.mx.Unlock()

	return a.len
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
