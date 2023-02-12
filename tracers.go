package muxtracer

import (
	"log"
	"time"
)

func logViolation(id Id, threshold Threshold, actual Actual, now Now, violationType ViolationType) {
	var idStr string
	if id != nil {
		idStr = string(id) + " "
	}
	log.Printf("%sviolation %s section took %s %d (threshold %s)", idStr, violationType.String(), time.Duration(actual).String(), actual, time.Duration(threshold).String())
}

type Threshold uint64
type Actual uint64
type Now uint64

type TraceLocker interface {
	EnableTracer()
	DisableTracer()
	EnableTracerWithOpts(o Opts)
}

func (m *Mutex) EnableTracer() {
	m.EnableTracerWithOpts(obtainGlobalOpts())
}

func (m *Mutex) EnableTracerWithOpts(o Opts) {
	if o.Id != "" {
		m.id = []byte(o.Id)
	}
	m.threshold.Store(uint64(o.Threshold.Nanoseconds()))
}

func (m *Mutex) DisableTracer() {
	m.threshold.Store(0)
}

func (m *RWMutex) EnableTracer() {
	m.EnableTracerWithOpts(obtainGlobalOpts())
}

func (m *RWMutex) EnableTracerWithOpts(o Opts) {
	if o.Id != "" {
		m.id = []byte(o.Id)
	}
	m.threshold.Store(uint64(o.Threshold.Nanoseconds()))
}

func (m *RWMutex) DisableTracer() {
	m.threshold.Store(0)
}

type Id []byte
