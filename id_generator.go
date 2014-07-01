package main

import (
	"encoding/binary"
	"sync"
	"time"
)

// A `TimeSource` returns the current time. `time.Now()` fits the bill, for
// example.
type TimeSource func() time.Time

// An Id is a 128 bit wide value. All values are encoded with big endian.
//
// * The first 64 bits encode milliseconds since unix epoch.
// * The next 48 bits encode the worker identifier, usually the MAC address of
//   the machine that generated the Id.
// * The final 16 bits encode a sequence to differentiate Ids generated in the
//   same millisecond.
type Id [16]byte

// A Worker ID is a 48 bit wide value, usually the MAC address of the machine
// that is generating IDs.
type WorkerId [6]byte

type IdGenerator struct {
	TimeSource  TimeSource
	CurrentTime time.Time
	WorkerId    WorkerId
	Sequence    uint16
	mutex       sync.Mutex
}

func (generator *IdGenerator) Generate() Id {
	timeSource := generator.TimeSource
	if timeSource == nil {
		timeSource = time.Now
	}

	// While we hold the lock, generate the current time and reset CurrentTime
	// and/or Sequence if needed.
	generator.mutex.Lock()
	ts := timeSource()
	generatorTimeMs := uint64(generator.CurrentTime.UnixNano() / 1e6)
	currentTimeMs := uint64(ts.UnixNano() / 1e6)
	if currentTimeMs > generatorTimeMs {
		generator.CurrentTime = ts
		generator.Sequence = 0
	} else {
		generator.Sequence++
	}
	sequence := generator.Sequence
	generator.mutex.Unlock()

	// Only the WorkerId property can be access safely here since we gave up the
	// lock. No other `generator` properties should be accessed.
	id := [16]byte{}
	// Timestamp (64 bits)
	binary.BigEndian.PutUint64(id[0:8], uint64(ts.UnixNano()/1e6))
	// Worker ID (48 bits)
	copy(id[8:14], generator.WorkerId[:])
	// Sequence (16 bits)
	binary.BigEndian.PutUint16(id[14:16], sequence)

	return id
}
