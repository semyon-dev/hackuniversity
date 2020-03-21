package opc

import (
	"context"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"log"
)

func GetData() {

	nodeIDs := []string{"ns=1;s=humidity", "ns=1;s=pressure", "ns=1;s=temphome", "ns=1;s=tempwork",
		"ns=1;s=levelph", "ns=1;s=levelco2", "ns=1;s=mass", "ns=1;s=water"}

	ctx := context.Background()

	c := opcua.NewClient("opc.tcp://185.251.90.101:4334", opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for _, nodeID := range nodeIDs {
		id, err := ua.ParseNodeID(nodeID)
		if err != nil {
			log.Fatalf("invalid node id: %v", err)
		}

		req := &ua.ReadRequest{
			MaxAge: 2000,
			NodesToRead: []*ua.ReadValueID{
				&ua.ReadValueID{NodeID: id},
			},
			TimestampsToReturn: ua.TimestampsToReturnBoth,
		}

		resp, err := c.Read(req)
		if err != nil {
			log.Fatalf("Read failed: %s", err)
		}
		if resp.Results[0].Status != ua.StatusOK {
			log.Fatalf("Status not OK: %v", resp.Results[0].Status)
		}
		log.Printf("%#v", resp.Results[0].Value.Value())
	}
}
