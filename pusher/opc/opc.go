package opc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/semyon-dev/hackuniversity/pusher/model"
	"log"
)

var c *opcua.Client

// connect to opc and return data
func GetData() ([]byte, model.Data) {
	nodeIDs := []string{"ns=1;s=humidity", "ns=1;s=pressure", "ns=1;s=temphome", "ns=1;s=tempwork",
		"ns=1;s=levelph", "ns=1;s=levelco2", "ns=1;s=mass", "ns=1;s=water"}

	ctx := context.Background()

	c = opcua.NewClient("opc.tcp://semyonpc:4334/", opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	var data model.Data

	data.HUMIDITY = getValue(nodeIDs[0])
	data.PRESSURE = getValue(nodeIDs[1])
	data.TEMPHOME = getValue(nodeIDs[2])
	data.TEMPWORK = getValue(nodeIDs[3])
	data.LEVELPH = getValue(nodeIDs[4])
	data.LEVELCO2 = getValue(nodeIDs[5])
	data.MASS = getValue(nodeIDs[6])
	data.WATER = getValue(nodeIDs[7])

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Print(err)
	}
	return jsonData, data
}

func getValue(nodeID string) float64 {
	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		log.Printf("invalid node id: %v", err)
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
		log.Printf("Read failed: %s", err)
	}
	if resp.Results[0].Status != ua.StatusOK {
		log.Printf("Status not OK: %v", resp.Results[0].Status)
	}
	log.Printf("data from OPC %#v", resp.Results[0].Value.Value())
	return resp.Results[0].Value.Float()
}
