package infrastructure

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func job(a interface{}) interface{} {
	return a.(int) % 1000
}

func gen(done <-chan struct{}) <-chan interface{} {
	result := make(chan interface{})
	go func() {
		defer close(result)

		for {
			select {
			case <-done:
				return
			default:
			}

			rand.Seed(time.Now().UnixNano())
			result <- rand.Int()
			<-time.After(10 * time.Millisecond)
		}
	}()

	return result
}

func Test_IncreaseAndReduceWorkers_Success(t *testing.T) {
	done := make(chan struct{})
	input := gen(done)

	m := NewMaster()
	output, _ := m.Run(input, 5, job)

	go func() {
		for o := range output {
			log.Println(o)
		}
	}()

	go func() {
		for {
			<-time.After(1 * time.Second)
			m.ReduceWorkers(1)
		}
	}()

	go func() {
		for {
			<-time.After(2 * time.Second)
			m.IncreaseWorkers(1)
		}
	}()

	<-time.After(10 * time.Second)
	close(done)
}

func Test_CreateAndStopMaster_Success(t *testing.T) {
	done := make(chan struct{})
	input := gen(done)

	m := NewMaster()
	output, cancel := m.Run(input, 5, job)

	go func() {
		for o := range output {
			log.Println(o)
		}
	}()

	<-time.After(3 * time.Second)
	cancel()
	<-time.After(1 * time.Second)
	close(done)
}

func Test_CreateAndStopWorker_Success(t *testing.T) {
	done := make(chan struct{})
	input := gen(done)

	w1 := CreateWorker(input, job)
	w2 := CreateWorker(input, job)

	go func() {
		for {
			select {
			case res, ok := <-w1.result:
				if !ok {
					return
				}
				log.Println("Result from w1: ", res)
			case res, ok := <-w2.result:
				if !ok {
					return
				}
				log.Println("Result from w2: ", res)
			}
		}
	}()

	<-time.After(3 * time.Second)
	w1.Stop()
	w2.Stop()

	<-time.After(1 * time.Second)
	close(done)
}
