package main

import (
	"encoding/binary"
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

type IdGenerator struct {
	TimeSource TimeSource
}

func (generator *IdGenerator) Generate() Id {
	id := [16]byte{}

	ts := generator.TimeSource()
	binary.BigEndian.PutUint64(id[0:8], uint64(ts.UnixNano()/1e6))

	return id
}
