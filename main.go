package main

import (
	"./crystal"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	cthrift "github.com/alindeman/crystal/thrift"
	"os"
)

func main() {
	service := flag.Bool("service", false, "run a crystal service")
	ifaceName := flag.String("iface", "", "network interface name used for worker ID (e.g., eth0)")
	generate := flag.Int("generate", 0, "generate N IDs by connecting to the crystal service")
	port := flag.Int("port", 5829, "TCP port to listen or connect to the crystal service")

	flag.Parse()
	if *service {
		if *ifaceName == "" {
			fmt.Printf("-iface must be specified when running the service\n\n")
			flag.PrintDefaults()
			os.Exit(255)
		}

		workerId, err := crystal.WorkerIdFromNetworkInterfaceName(*ifaceName)
		if err != nil {
			fmt.Printf("Unable to determine worker ID from interface %v: %v\n", *ifaceName, err)
			os.Exit(1)
		}

		generator := &crystal.IdGenerator{
			WorkerId: workerId,
		}
		crystal.ListenAndServeThrift(fmt.Sprintf(":%d", *port), generator)
	} else if *generate > 0 {
		var transport thrift.TTransport
		transport, err := thrift.NewTSocket(fmt.Sprintf("127.0.0.1:%d", *port))
		if err != nil {
			fmt.Printf("Unable to connect to crystal service on port %v: %v\n", *port, err)
			os.Exit(1)
		}

		transport = thrift.NewTFramedTransport(transport)
		if err := transport.Open(); err != nil {
			fmt.Printf("Unable to open transport to crystal service: %v\n", err)
			os.Exit(1)
		}
		defer transport.Close()

		client := cthrift.NewIdGenerationServiceClientFactory(transport, thrift.NewTBinaryProtocolFactoryDefault())
		for i := 0; i < *generate; i++ {
			res, err := client.Generate()
			if err != nil {
				fmt.Printf("Unable to generate ID: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%v\n", res)
		}
	} else {
		flag.PrintDefaults()
	}
}
