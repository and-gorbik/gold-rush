package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	licensesCheckDuration = time.Millisecond * 100
	licenseIssuePeriod    = time.Millisecond * 100
	maxActiveLicenses     = 10
	maxBuf                = 100
)

type LicenseService struct {
	server   handlers
	attempts map[string]int
	mx       sync.Mutex
}

type license struct {
	TrackID string
	ID      int
	Num     int
}

type License struct {
	TrackID string
	ID      int
}

func NewLicenceService(server handlers) *LicenseService {
	ls := &LicenseService{
		attempts: make(map[string]int, maxActiveLicenses),
		server:   server,
	}

	ls.checkActualLicenses(licensesCheckDuration)
	return ls
}

func (ls *LicenseService) StreamLicenses() <-chan License {
	licenses := make(chan license, maxActiveLicenses-1)
	licenseIDs := make(chan License, maxBuf)

	go func() {
		defer close(licenses)

		for {
			li := ls.issueLicense()
			if li == nil {
				<-time.After(licenseIssuePeriod)
				continue
			}

			licenses <- *li
		}
	}()

	go func() {
		defer close(licenseIDs)

		var wg sync.WaitGroup

		for lc := range licenses {
			wg.Add(1)
			go func(l license) {
				defer wg.Done()

				for i := 0; i < l.Num; i++ {
					licenseIDs <- License{
						TrackID: l.TrackID,
						ID:      l.ID,
					}
				}
			}(lc)
		}

		wg.Wait()
	}()

	return licenseIDs
}

func (ls *LicenseService) Done(trackID string) {
	ls.mx.Lock()
	ls.attempts[trackID]--
	if ls.attempts[trackID] <= 0 {
		delete(ls.attempts, trackID)
	}
	ls.mx.Unlock()
}

func (ls *LicenseService) issueLicense() *license {
	trackID := uuid.NewString()

	ls.mx.Lock()
	if len(ls.attempts) == maxActiveLicenses {
		ls.mx.Unlock()
		return nil
	}

	ls.attempts[trackID] = 0
	ls.mx.Unlock()

	l := ls.server.IssueLicense(context.Background(), []int{})
	log.Println("issued license: ", l)

	ls.mx.Lock()
	ls.attempts[trackID] = l.DigAllowed
	defer ls.mx.Unlock()

	return &license{trackID, l.ID, l.DigAllowed}
}

func (ls *LicenseService) checkActualLicenses(dur time.Duration) {
	go func() {
		for {
			ls.mx.Lock()
			cur := len(ls.attempts)
			ls.mx.Unlock()

			if cur < 0 || cur > maxActiveLicenses {
				log.Fatal(fmt.Sprintf("inexpected number of active licenses: %d", cur))
			}

			<-time.After(dur)
		}
	}()
}
