package crystal

import (
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

var (
	ClockRunningBackwards = errors.New("Clock is running backwards")
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

func (generator *IdGenerator) Generate() (Id, error) {
	id := [16]byte{}
	timeSource := generator.TimeSource
	if timeSource == nil {
		timeSource = time.Now
	}

	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	ts := timeSource()
	generatorTimeMs := uint64(generator.CurrentTime.UnixNano() / 1e6)
	currentTimeMs := uint64(ts.UnixNano() / 1e6)
	if generator.CurrentTime.IsZero() || currentTimeMs > generatorTimeMs {
		generator.CurrentTime = ts
		generator.Sequence = 0
	} else if currentTimeMs == generatorTimeMs {
		generator.Sequence++
	} else {
		return id, ClockRunningBackwards
	}

	// Timestamp (64 bits)
	binary.BigEndian.PutUint64(id[0:8], uint64(ts.UnixNano()/1e6))
	// Worker ID (48 bits)
	copy(id[8:14], generator.WorkerId[:])
	// Sequence (16 bits)
	binary.BigEndian.PutUint16(id[14:16], generator.Sequence)

	return id, nil
}
