package main

import (
	"encoding/binary"
	"testing"
	"time"
)

// Reconstructs a timestamp from the first 64 bits of an Id
func decodeTimestamp(id Id) time.Time {
	msSinceEpoch := int64(binary.BigEndian.Uint64(id[0:8]))
	return time.Unix(msSinceEpoch/1e3, (msSinceEpoch%1e3)*1e6)
}

func TestTimestampIsFirst64Bits(t *testing.T) {
	ts := time.Unix(12345689, 200*1e6)
	generator := &IdGenerator{
		TimeSource: func() time.Time { return ts },
	}

	id := generator.Generate()
	idTimestamp := decodeTimestamp(id)
	if idTimestamp != ts {
		t.Fatalf("Expected timestamp in ID to be %v, but was %v", ts, idTimestamp)
	}
}
