package table

import (
	"container/list"
	"../nodes"
	"fmt"
)


const BucketSize = 20;

type Contact struct {
	id nodes.NodeID
	ip string
	port int
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

func findElementInBucket(bucket *list.List, id nodes.NodeID) (ret nodes.NodeID){
	for i := range bucket {
		if id.Equals(bucket[i].(nodes.NodeID)) {
			ret = bucket[i].(nodes.NodeID)
			return
		}
	}
	ret = nil
	return
}

func (table *RoutingTable) Update(contact *Contact) {
	prefixLen := table.node.Distance(contact.id).PrefixLen()
	bucket := table.buckets[prefixLen]
	element := findElementInBucket(bucket, contact.id)

	if element != nil {
		if bucket.Len() <= BucketSize {
			bucket.PushFront(contact)
		} else {
			// Send ping to oldest contact
			oldestContact := bucket.Back().Value.(nodes.NodeID)
			fmt.Println(oldestContact)
			//TODO Replace print above with actual call to ping
		}
	}
}