package infrastructure

import (
	"sync"

	"github.com/google/uuid"
)

const (
	BUF = 100
)

type Worker struct {
	done   chan struct{}
	result chan interface{}
}

type Master struct {
	workers map[string]Worker
	mx      sync.RWMutex
	in      <-chan interface{}
	job     func(interface{}) interface{}
	done    chan struct{}
}

func NewMaster() Master {
	return Master{
		workers: make(map[string]Worker),
	}
}

func (m *Master) Run(in <-chan interface{}, workers int, job func(interface{}) interface{}) (<-chan interface{}, func()) {
	m.in = in
	m.job = job

	m.IncreaseWorkers(workers)

	done := make(chan struct{})
	result := make(chan interface{}, BUF)

	go func() {
		defer close(result)

		for {
			m.mx.Lock()
			for id, w := range m.workers {
				select {
				case result <- <-w.result:
				case <-done:
					w.Stop()
					delete(m.workers, id)
				default:
				}
			}

			if len(m.workers) == 0 {
				m.mx.Unlock()
				return
			}

			m.mx.Unlock()
		}
	}()

	m.done = done

	cancel := func() {
		close(done)
	}

	return result, cancel
}

func (m *Master) IncreaseWorkers(num int) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for i := 0; i < num; i++ {
		m.workers[uuid.NewString()] = CreateWorker(m.in, m.job)
	}
}

func (m *Master) ReduceWorkers(num int) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for id := range m.workers {
		if num == 0 {
			break
		}

		close(m.workers[id].done)
		delete(m.workers, id)
		num--
	}

	if len(m.workers) == 0 {
		close(m.done)
	}
}

func CreateWorker(in <-chan interface{}, job func(interface{}) interface{}) Worker {
	w := Worker{
		done:   make(chan struct{}),
		result: make(chan interface{}),
	}

	go work(w.done, in, w.result, job)

	return w
}

func work(done <-chan struct{}, in <-chan interface{}, out chan<- interface{}, job func(interface{}) interface{}) {
	defer close(out)

	for {
		select {
		case <-done:
			return
		case out <- job(<-in):
		}
	}
}

func (w *Worker) Stop() {
	close(w.done)
}
