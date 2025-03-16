package service

import (
	"context"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func Test_Service(t *testing.T) {
	defer goleak.VerifyNone(t)

	MaxArea = 100
	ts := setupNewMockServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	Run(ctx, ts)
}
