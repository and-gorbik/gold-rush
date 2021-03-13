package app

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"gold-rush/app/mocks"
)

func Test_Explore_Success(t *testing.T) {
	a := App{provider: mocks.GoodProvider{}}
	q := a.explore(runtime.NumCPU())

	counter := 0
	for i := 0; i < areasCount; i++ {
		q.Pop()
		counter++
	}

	assert.Equal(t, areasCount, counter)
}

func Test_Explore_BadProvider(t *testing.T) {

}
