package kademlia

import (
	"../table"
	"../nodes"
)

type Kademlia struct {
	RoutingTable *table.RoutingTable
}

func NewKademlia(myNode nodes.NodeID) Kademlia{
	ret := Kademlia{
		RoutingTable: table.NewRoutingTable(myNode),
	}
	return ret
}
