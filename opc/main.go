package main

import (
	"flag"
	"fmt"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"log"
)

var add = "opc.tcp://192.168.1.109:4334/UA/MyLittleServer"

func main() {
	endpoint := flag.String("endpoint", "opc.tcp://192.168.1.109:4334/UA/MyLittleServer", "OPC UA Endpoint URL")
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)

	//ctx := context.Background()

	endpoints, err := opcua.GetEndpoints(*endpoint)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(endpoints); i++ {
		fmt.Printf("%s \n", endpoints[i].EndpointURL)
	}
}
