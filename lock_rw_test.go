package muxtracer

import (
	"sync"
	"testing"
	"time"
)

func TestRWNew(t *testing.T) {
	l := RWMutex{}
	l.Lock()
	l.Unlock()
}

func TestRWNewEnabled(t *testing.T) {
	l := RWMutex{}
	l.EnableTracer()
	l.Lock()
	l.Unlock()
}

func TestRWNewDisabled(t *testing.T) {
	l := RWMutex{}
	l.DisableTracer()
	l.Lock()
	l.Unlock()
}

func TestRWNewEnabledDisabledHalfWay(t *testing.T) {
	l := RWMutex{}
	l.EnableTracer()
	l.Lock()
	l.DisableTracer()
	l.Unlock()
}

func TestRWNewEnabledDisabledEnd(t *testing.T) {
	l := RWMutex{}
	l.EnableTracer()
	l.Lock()
	l.Unlock()
	l.DisableTracer()
}

func TestRWNewEnableGlobal(t *testing.T) {
	l := RWMutex{}

	// enable globally
	SetGlobalOpts(Opts{
		Threshold: 100 * time.Millisecond,
		Enabled:   true,
	})

	l.Lock()
	time.Sleep(150 * time.Millisecond)
	l.Unlock()

	// reset again
	ResetDefaults()
}

func TestRWNewEnabledHalfWay(t *testing.T) {
	l := RWMutex{}
	l.Lock()
	l.EnableTracer()
	l.Unlock()
	l.DisableTracer()
}

func TestRWNewEnabledShortDelay(t *testing.T) {
	l := RWMutex{}
	l.EnableTracer()
	l.Lock()
	time.Sleep(1 * time.Millisecond)
	l.Unlock()
	l.DisableTracer()
}

func TestRWNewEnabledLongDelay(t *testing.T) {
	l := RWMutex{}
	l.EnableTracer()
	l.Lock()
	time.Sleep(150 * time.Millisecond)
	l.Unlock()
	l.DisableTracer()
}

func TestRWNewEnabledAwaitLock(t *testing.T) {
	l := RWMutex{}
	l.EnableTracerWithOpts(Opts{
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

func TestRWNewEnabledId(t *testing.T) {
	l := Mutex{}
	l.EnableTracerWithOpts(Opts{
		Threshold: 10 * time.Millisecond,
		Id:        "testRwLock",
	})
	l.Lock()
	time.Sleep(20 * time.Millisecond)
	l.Unlock()
	l.DisableTracer()
}

func BenchmarkRWNativeLock(b *testing.B) {
	l := RWMutex{}
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkRWTracerLockDisabled(b *testing.B) {
	l := RWMutex{}
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkRWTracerLockEnabled(b *testing.B) {
	l := RWMutex{}
	l.EnableTracer()
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkRWNativeLockWithConcurrency(b *testing.B) {
	l := RWMutex{}
	wg := sync.WaitGroup{}
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

func BenchmarkRWTracerLockDisabledWithConcurrency(b *testing.B) {
	l := RWMutex{}
	wg := sync.WaitGroup{}
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

func BenchmarkRWTracerLockEnabledWithConcurrency(b *testing.B) {
	l := RWMutex{}
	l.EnableTracer()
	wg := sync.WaitGroup{}
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
