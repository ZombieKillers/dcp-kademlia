package kademlia

import (
	"../table"
)

type Kademlia struct {
	routes *table.RoutingTable
	NetworkId string
}

func NewKademlia(self table.Contact, networkId string) (ret Kademlia) {
	ret = Kademlia{table.NewRoutingTable(self), networkId }
	return
}


