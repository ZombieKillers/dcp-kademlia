package kademlia

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
)

const BucketSize = 20

type RoutingTable struct {
	node      NodeID
	ipAddress string
	port      int
	buckets   [IdLength * 8]*list.List
}

func NewRoutingTable(contact Contact) (ret RoutingTable) {
	ret.node = contact.Id
	ret.ipAddress = contact.Ip
	ret.port = contact.Port
	for i := 0; i < IdLength*8; i++ {
		ret.buckets[i] = list.New()
	}
	return
}

func findElementInBucket(bucket *list.List, id NodeID) (ret *NodeID) {

	for e := bucket.Front(); e != nil; e = e.Next() {
		if id.Equals(e.Value.(*Contact).Id) {
			ret = &e.Value.(*Contact).Id
			return
		}
	}
	return
}

func (table *RoutingTable) Update(contact *Contact) {
	prefixLen := table.node.Distance(contact.Id).PrefixLen()
	bucket := table.buckets[prefixLen]
	element := findElementInBucket(bucket, contact.Id)

	if element == nil {
		if bucket.Len() <= BucketSize {
			bucket.PushFront(contact)
		} else {
			// Send ping to oldest contact
			oldestContact := bucket.Back().Value.(NodeID)
			fmt.Println(oldestContact)
			//TODO Replace print above with actual call to ping
		}
	}

	fmt.Println("Updated routing table")
	fmt.Println(table)
}

func (contact *Contact) String() (s string) {
	s = "("
	s += contact.Id.String() + ","
	s += contact.Ip + ","
	s += strconv.Itoa(contact.Port) + ")"
	return
}

func (table *RoutingTable) String() string {
	s := ""
	bucketNumbers := make([]int, 0)

	for bucketNumber, b := range table.buckets {
		if b.Len() > 0 {
			s += "[BUCKET " + strconv.Itoa(bucketNumber) + " ]\n"
			for el := b.Front(); el != nil; el = el.Next() {
				s += el.Value.(*Contact).String() + "\n"
			}
			bucketNumbers = append(bucketNumbers, bucketNumber)
		}
	}
	s += "Occupied buckets: "
	s += "[ "
	for _, bucket := range bucketNumbers {
		s += " " + strconv.Itoa(bucket)
	}
	s += " ]"
	return s
}

func copyToList(start, end *list.Element, vec *[]*ContactRecord, target NodeID) {
	for el := start; el != end; el = el.Next() {
		contact := el.Value.(*Contact)
		*vec = append(*vec, &ContactRecord{node: contact, sortKey: contact.Id.Distance(target)})
	}
}

func (table *RoutingTable) FindClosest(target NodeID, count int) (ret []*ContactRecord) {
	ret = make([]*ContactRecord, 0)
	prefixLen := table.node.Distance(target).PrefixLen()
	bucket := table.buckets[prefixLen]
	copyToList(bucket.Front(), nil, &ret, target)
	for i := 1; (prefixLen-i >= 0 || prefixLen+i < IdLength*8) && len(ret) < count; i++ {
		if prefixLen-i >= 0 {
			bucket = table.buckets[prefixLen-i]
			copyToList(bucket.Front(), nil, &ret, target)
		}
		if prefixLen+i < IdLength*8 {
			bucket = table.buckets[prefixLen+i]
			copyToList(bucket.Front(), nil, &ret, target)
		}
	}
	sort.Sort(ContactRecordList(ret))
	if len(ret) > count {
		ret = ret[0:count]
	}
	return
}
