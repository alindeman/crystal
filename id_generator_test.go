package main

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

	id := generator.Generate()
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

	id := generator.Generate()
	if bytes.Compare(id[8:14], workerId[:]) != 0 {
		t.Fatalf("Expected worker ID in ID to be %v, but was %v", workerId, id[8:14])
	}
}

func TestSequenceIsFinal16Bits(t *testing.T) {
	sequence := uint16(1234)
	generator := &IdGenerator{
		Sequence: sequence,
	}

	sequenceBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(sequenceBytes, sequence)

	id := generator.Generate()
	if bytes.Compare(id[14:16], sequenceBytes) != 0 {
		t.Fatalf("Expected sequence in ID to be %v, but was %v", sequenceBytes, id[14:16])
	}
}

func TestSequenceIsIncrementedForSameTimestamp(t *testing.T) {
	ts := time.Unix(12345689, 200*1e6)
	generator := &IdGenerator{
		TimeSource: func() time.Time { return ts },
	}

	id := generator.Generate()
	if bytes.Compare(id[14:16], []byte{0, 0}) != 0 {
		t.Fatalf("Expected sequence in ID to be [0 0] on the first run, but was %v", id[14:16])
	}

	id2 := generator.Generate()
	if bytes.Compare(id2[14:16], []byte{0, 1}) != 0 {
		t.Fatalf("Expected sequence in ID to be [0 1] on the second run, but was %v", id[14:16])
	}

}
