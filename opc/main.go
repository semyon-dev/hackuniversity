package main

import (
	"context"
	"flag"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/ua"
	"log"
)

var add = "opc.tcp://192.168.1.109:4334/UA/MyLittleServer"

func main() {
	var (
		endpoint = flag.String("endpoint", "opc.tcp://192.168.1.109:4334/UA/MyLittleServer", "OPC UA Endpoint URL")
		nodeID   = flag.String("node", "0", "NodeID to read")
	)
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)

	ctx := context.Background()

<<<<<<< HEAD
	var (
		nodeID   = flag.String("node", "", "NodeID to read")
	)
	flag.Parse()
	log.SetFlags(0)

	ctx := context.Background()

=======
>>>>>>> 85bf9a4e12376af0727a0847bee9f3374847c395
	c := opcua.NewClient(*endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Close()
<<<<<<< HEAD
=======
	c.RegisterNodes()
>>>>>>> 85bf9a4e12376af0727a0847bee9f3374847c395

	id, err := ua.ParseNodeID(*nodeID)
	if err != nil {
		log.Fatalf("invalid node id: %v", err)
	}

	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			&ua.ReadValueID{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
<<<<<<< HEAD
	}

	resp, err := c.Read(req)
	if err != nil {
		log.Fatalf("Read failed: %s", err)
	}
	if resp.Results[0].Status != ua.StatusOK {
		log.Fatalf("Status not OK: %v", resp.Results[0].Status)
	}
	log.Printf("%#v", resp.Results[0].Value.Value())




=======
	}

	resp, err := c.Read(req)
	if err != nil {
		log.Fatalf("Read failed: %s", err)
	}
	if resp.Results[0].Status != ua.StatusOK {
		log.Fatalf("Status not OK: %v", resp.Results[0].Status)
	}
	fmt.Println("resp.Results", resp.Results)
	fmt.Println("-------------------")
	log.Printf("%#v", resp.Results[0].Value.Value())
>>>>>>> 85bf9a4e12376af0727a0847bee9f3374847c395
}
