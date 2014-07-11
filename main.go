package main

import (
	"./crystal"
	"git.apache.org/thrift.git/lib/go/thrift"
	cthrift "github.com/alindeman/crystal/thrift"
	"fmt"
	"os"
)

func main() {
	if os.Args[1] == "service" {
		generator := &crystal.IdGenerator{}
		crystal.ListenAndServeThrift(":8899", generator)
	} else if os.Args[1] == "client" {
		// TODO: Wrap this up somewhere else
		var transport thrift.TTransport
		transport, err := thrift.NewTSocket("127.0.0.1:8899")
		if err != nil {
			panic(err)
		}

		transport = thrift.NewTFramedTransport(transport)
		if err := transport.Open(); err != nil {
			panic(err)
		}
		defer transport.Close()

		client := cthrift.NewIdGenerationServiceClientFactory(transport, thrift.NewTBinaryProtocolFactoryDefault())
		res, err := client.Generate()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", res)
	}
}
