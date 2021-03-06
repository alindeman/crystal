package crystal

import (
	"bytes"
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

	id, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	idTimestamp := decodeTimestamp(id)
	if idTimestamp != ts {
		t.Fatalf("Expected timestamp in ID to be %v, but was %v", ts, idTimestamp)
	}
}

func TestWorkerIdIsNext48Bits(t *testing.T) {
	workerId := WorkerId{1, 2, 3, 4, 5, 6}
	generator := &IdGenerator{
		WorkerId: workerId,
	}

	id, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(id[8:14], workerId[:]) != 0 {
		t.Fatalf("Expected worker ID in ID to be %v, but was %v", workerId, id[8:14])
	}
}

func TestSequenceIsFinal16Bits(t *testing.T) {
	sequence := uint16(1234)
	ts := time.Unix(12345689, 200*1e6)
	generator := &IdGenerator{
		CurrentTime: ts,
		TimeSource:  func() time.Time { return ts },
		Sequence:    sequence,
	}

	sequenceBytes := make([]byte, 2)
	// the sequence will be incremented by one before an ID is generated
	binary.BigEndian.PutUint16(sequenceBytes, sequence+1)

	id, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(id[14:16], sequenceBytes) != 0 {
		t.Fatalf("Expected sequence in ID to be %v, but was %v", sequenceBytes, id[14:16])
	}
}

func TestSequenceIsIncrementedForSameTimestamp(t *testing.T) {
	ts := time.Unix(12345689, 200*1e6)
	generator := &IdGenerator{
		TimeSource: func() time.Time { return ts },
	}

	id, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(id[14:16], []byte{0, 0}) != 0 {
		t.Fatalf("Expected sequence in ID to be [0 0] on the first run, but was %v", id[14:16])
	}

	id2, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(id2[14:16], []byte{0, 1}) != 0 {
		t.Fatalf("Expected sequence in ID to be [0 1] on the second run, but was %v", id2[14:16])
	}
}

func TestSequenceIsResetWhenTimeMovesForward(t *testing.T) {
	ts := time.Unix(123456789, 200*1e6)
	generator := &IdGenerator{
		CurrentTime: ts,
		TimeSource: func() time.Time {
			ts = ts.Add(time.Millisecond) // each invocation, increment timestamp
			return ts
		},
	}

	id, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(id[14:16], []byte{0, 0}) != 0 {
		t.Fatalf("Expected sequence in ID to be [0 0] on the first run, but was %v", id[14:16])
	}

	id2, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(id2[14:16], []byte{0, 0}) != 0 {
		t.Fatalf("Expected sequence in ID to be [0 0] on the second run, but was %v", id2[14:16])
	}
}

func TestClockRunningBackwardsIsAnError(t *testing.T) {
	ts := time.Unix(123456789, 200*1e6)
	generator := &IdGenerator{
		CurrentTime: ts,
		TimeSource:  func() time.Time { return ts.Add(-time.Millisecond) },
	}

	_, err := generator.Generate()
	if err != ClockRunningBackwards {
		t.Fatalf("Expected a ClockRunningBackward error, but was %v", err)
	}
}
