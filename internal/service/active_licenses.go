package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	maxActiveLicenses = 10

	defaultTryIssueLicensePeriod = time.Millisecond * 100
	defaultMaxLicencesBuf        = 100
)

type LicenseService struct {
	server   handlers
	attempts map[string]int
	mx       sync.Mutex
	opts     LicenserOpts
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

type LicenserOpts struct {
	TryIssueLicensePeriod time.Duration
	MaxLicencesBuf        int
}

func (lo *LicenserOpts) SetDefaultsIfNeeded() {
	if lo.MaxLicencesBuf == 0 {
		lo.MaxLicencesBuf = defaultMaxLicencesBuf
	}

	if lo.TryIssueLicensePeriod == 0 {
		lo.TryIssueLicensePeriod = defaultTryIssueLicensePeriod
	}
}

func NewLicenceService(server handlers, opts LicenserOpts) *LicenseService {
	opts.SetDefaultsIfNeeded()

	ls := &LicenseService{
		attempts: make(map[string]int, maxActiveLicenses),
		server:   server,
		opts:     opts,
	}

	// ls.checkActualLicenses(licensesCheckPeriod)
	return ls
}

func (ls *LicenseService) StreamLicenses(ctx context.Context) <-chan License {
	licenses := make(chan license, maxActiveLicenses-1)
	licenseIDs := make(chan License, ls.opts.MaxLicencesBuf)

	go func() {
		defer close(licenses)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			li := ls.issueLicense()
			if li == nil {
				<-time.After(ls.opts.TryIssueLicensePeriod)
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

// func (ls *LicenseService) checkActualLicenses(dur time.Duration) {
// 	go func() {
// 		for {
// 			ls.mx.Lock()
// 			cur := len(ls.attempts)
// 			ls.mx.Unlock()

// 			if cur < 0 || cur > maxActiveLicenses {
// 				log.Fatal(fmt.Sprintf("inexpected number of active licenses: %d", cur))
// 			}

// 			<-time.After(dur)
// 		}
// 	}()
// }
