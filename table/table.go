package table

import (
	"container/list"
	"../nodes"
)


const BucketSize = 20;

type Contact struct {
	id nodes.NodeID;
}

type RoutingTable struct {
	node nodes.NodeID;
	buckets [nodes.IdLength*8]*list.List;
}

func NewRoutingTable(node nodes.NodeID) (ret RoutingTable) {
	for i := 0; i < nodes.IdLength * 8; i++ {
		ret.buckets[i] = list.New();
	}
	ret.node = node;
	return;
}