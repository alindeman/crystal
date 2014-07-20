// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package thrift

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
)

// (needed to ensure safety because of naive import list construction.)
var _ = math.MinInt32
var _ = thrift.ZERO
var _ = fmt.Printf

var GoUnusedProtection__ int

type Error int64

const (
	Error_NONE                    Error = 0
	Error_CLOCK_RUNNING_BACKWARDS Error = 1
)

func (p Error) String() string {
	switch p {
	case Error_NONE:
		return "Error_NONE"
	case Error_CLOCK_RUNNING_BACKWARDS:
		return "Error_CLOCK_RUNNING_BACKWARDS"
	}
	return "<UNSET>"
}

func ErrorFromString(s string) (Error, error) {
	switch s {
	case "Error_NONE":
		return Error_NONE, nil
	case "Error_CLOCK_RUNNING_BACKWARDS":
		return Error_CLOCK_RUNNING_BACKWARDS, nil
	}
	return Error(math.MinInt32 - 1), fmt.Errorf("not a valid Error string")
}

type Id []byte

type IdGenerationResult struct {
	Err Error `thrift:"err,1"`
	Id  Id    `thrift:"id,2"`
}

func NewIdGenerationResult() *IdGenerationResult {
	return &IdGenerationResult{
		Err: math.MinInt32 - 1, // unset sentinal value
	}
}

func (p *IdGenerationResult) IsSetErr() bool {
	return int64(p.Err) != math.MinInt32-1
}

func (p *IdGenerationResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *IdGenerationResult) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 1: %s")
	} else {
		p.Err = Error(v)
	}
	return nil
}

func (p *IdGenerationResult) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return fmt.Errorf("error reading field 2: %s")
	} else {
		p.Id = Id(v)
	}
	return nil
}

func (p *IdGenerationResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("IdGenerationResult"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *IdGenerationResult) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetErr() {
		if err := oprot.WriteFieldBegin("err", thrift.I32, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:err: %s", p, err)
		}
		if err := oprot.WriteI32(int32(p.Err)); err != nil {
			return fmt.Errorf("%T.err (1) field write error: %s", p)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:err: %s", p, err)
		}
	}
	return err
}

func (p *IdGenerationResult) writeField2(oprot thrift.TProtocol) (err error) {
	if p.Id != nil {
		if err := oprot.WriteFieldBegin("id", thrift.BINARY, 2); err != nil {
			return fmt.Errorf("%T write field begin error 2:id: %s", p, err)
		}
		if err := oprot.WriteBinary(p.Id); err != nil {
			return fmt.Errorf("%T.id (2) field write error: %s", p)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 2:id: %s", p, err)
		}
	}
	return err
}

func (p *IdGenerationResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("IdGenerationResult(%+v)", *p)
}