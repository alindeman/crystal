package crystal

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	cthrift "github.com/alindeman/crystal/thrift"
)

func ListenAndServeThrift(laddr string, generator *IdGenerator) error {
	var transport thrift.TServerTransport
	transport, err := thrift.NewTServerSocket(laddr)
	if err != nil {
		return err
	}

	handler := &IdGenerationServiceHandler{Generator: generator}
	processor := cthrift.NewIdGenerationServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(
		processor,
		transport,
		thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()),
		thrift.NewTBinaryProtocolFactoryDefault(),
	)
	return server.Serve()
}

type IdGenerationServiceHandler struct {
	Generator *IdGenerator
}

func (handler *IdGenerationServiceHandler) Generate() (*cthrift.IdGenerationResult, error) {
	resp := cthrift.NewIdGenerationResult()

	id, err := handler.Generator.Generate()
	if err == ClockRunningBackwards {
		resp.Err = cthrift.Error_CLOCK_RUNNING_BACKWARDS
	} else if err != nil {
		// An error we are not equipped to handle :(
		return resp, err
	} else {
		resp.Err = cthrift.Error_NONE
		resp.Id = id[:]
	}

	return resp, nil
}
