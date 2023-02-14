package muxtracer

import (
	"sync"
	"sync/atomic"
)

type Mutex struct {
	lock sync.Mutex

	// internal trace fields
	threshold        atomic.Uint64 // 0 when disabled, else threshold in nanoseconds
	beginAwaitLock   atomic.Uint64 // start time in unix nanoseconds from start waiting for lock
	beginAwaitUnlock atomic.Uint64 // start time in unix nanoseconds from start waiting for unlock
	lockObtained     atomic.Uint64 // once we've entered the lock in unix nanoseconds
	id               []byte        // if set this will be printed as string
}

func (m *Mutex) Lock() {
	tracingThreshold := m.isTracing()
	if tracingThreshold != 0 {
		m.traceBeginAwaitLock()
	}

	// actual lock
	m.lock.Lock()

	if tracingThreshold != 0 {
		m.traceEndAwaitLock(tracingThreshold)
	}
}

func (m *Mutex) Unlock() {
	tracingThreshold := m.isTracing()
	if tracingThreshold != 0 {
		m.traceBeginAwaitUnlock()
	}

	// unlock
	m.lock.Unlock()

	if tracingThreshold != 0 {
		m.traceEndAwaitUnlock(tracingThreshold)
	}
}

func (m *Mutex) isTracing() Threshold {
	tracingThreshold := m.threshold.Load()
	if tracingThreshold == 0 {
		// always on?
		tracingThreshold = defaultThreshold.Load()
	}
	return Threshold(tracingThreshold)
}

func (m *Mutex) traceBeginAwaitLock() {
	m.beginAwaitLock.Store(now())
}

func (m *Mutex) traceEndAwaitLock(threshold Threshold) {
	ts := now() // first obtain the current time
	start := m.beginAwaitLock.Load()
	m.lockObtained.Store(ts)
	var took uint64
	if start < ts {
		// check for no overflow
		took = ts - start
	}
	if took >= uint64(threshold) {
		logViolation(Id(m.id), threshold, Actual(took), Now(ts), ViolationLock)
	}
}

func (m *Mutex) traceBeginAwaitUnlock() {
	m.beginAwaitUnlock.Store(now())
}

func (m *Mutex) traceEndAwaitUnlock(threshold Threshold) {
	ts := now() // first obtain the current time

	// lock obtained time (critical section)
	lockObtained := m.lockObtained.Load()
	var took uint64
	if lockObtained < ts {
		// check for no overflow
		took = ts - lockObtained
	}

	if took >= uint64(threshold) && lockObtained > 0 {
		// lockObtained = 0 when the tracer is enabled half way
		logViolation(Id(m.id), threshold, Actual(took), Now(ts), ViolationCritical)
	}
}
