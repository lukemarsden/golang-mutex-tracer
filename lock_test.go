package muxtracer_test

import (
	sync "github.com/lukemarsden/golang-mutex-tracer"
	nativeSync "sync"
	"testing"
	"time"
)

const numRoutines = 16

func TestNew(t *testing.T) {
	l := sync.Mutex{}
	l.Lock()
	l.Unlock()
}

func TestNewEnabled(t *testing.T) {
	l := sync.Mutex{}
	l.EnableTracer()
	l.Lock()
	l.Unlock()
}

func TestNewDisabled(t *testing.T) {
	l := sync.Mutex{}
	l.DisableTracer()
	l.Lock()
	l.Unlock()
}

func TestNewEnabledDisabledHalfWay(t *testing.T) {
	l := sync.Mutex{}
	l.EnableTracer()
	l.Lock()
	l.DisableTracer()
	l.Unlock()
}

func TestNewEnabledDisabledEnd(t *testing.T) {
	l := sync.Mutex{}
	l.EnableTracer()
	l.Lock()
	l.Unlock()
	l.DisableTracer()
}

func TestNewEnableGlobal(t *testing.T) {
	l := sync.Mutex{}

	// enable globally
	sync.SetGlobalOpts(sync.Opts{
		Threshold: 100 * time.Millisecond,
		Enabled:   true,
	})

	l.Lock()
	time.Sleep(150 * time.Millisecond)
	l.Unlock()

	// reset again
	sync.ResetDefaults()
}

func TestNewEnabledHalfWay(t *testing.T) {
	l := sync.Mutex{}
	l.Lock()
	l.EnableTracer()
	l.Unlock()
	l.DisableTracer()
}

func TestNewEnabledShortDelay(t *testing.T) {
	l := sync.Mutex{}
	l.EnableTracer()
	l.Lock()
	time.Sleep(1 * time.Millisecond)
	l.Unlock()
	l.DisableTracer()
}

func TestNewEnabledLongDelay(t *testing.T) {
	l := sync.Mutex{}
	l.EnableTracer()
	l.Lock()
	time.Sleep(150 * time.Millisecond)
	l.Unlock()
	l.DisableTracer()
}

func TestNewEnabledAwaitLock(t *testing.T) {
	l := sync.Mutex{}
	l.EnableTracerWithOpts(sync.Opts{
		Threshold: 10 * time.Millisecond,
	})
	go func() {
		// concurrent await
		l.Lock()
		time.Sleep(5 * time.Millisecond)
		l.Unlock()
	}()
	l.Lock()
	time.Sleep(20 * time.Millisecond)
	l.Unlock()
	l.DisableTracer()
}

func TestNewEnabledId(t *testing.T) {
	l := sync.Mutex{}
	l.EnableTracerWithOpts(sync.Opts{
		Threshold: 10 * time.Millisecond,
		Id:        "testLock",
	})
	l.Lock()
	time.Sleep(20 * time.Millisecond)
	l.Unlock()
	l.DisableTracer()
}

func BenchmarkNativeLock(b *testing.B) {
	l := nativeSync.Mutex{}
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkTracerLockDisabled(b *testing.B) {
	l := sync.Mutex{}
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkTracerLockEnabled(b *testing.B) {
	l := sync.Mutex{}
	l.EnableTracer()
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkNativeLockWithConcurrency(b *testing.B) {
	l := nativeSync.Mutex{}
	wg := nativeSync.WaitGroup{}
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			for n := 0; n < b.N; n++ {
				l.Lock()
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkTracerLockDisabledWithConcurrency(b *testing.B) {
	l := sync.Mutex{}
	wg := nativeSync.WaitGroup{}
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			for n := 0; n < b.N; n++ {
				l.Lock()
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkTracerLockEnabledWithConcurrency(b *testing.B) {
	l := sync.Mutex{}
	l.EnableTracer()
	wg := nativeSync.WaitGroup{}
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			for n := 0; n < b.N; n++ {
				l.Lock()
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
