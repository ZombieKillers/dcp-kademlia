package kademlia

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
)



const BucketSize = 20;

type Contact struct {
	Id NodeID
	Ip string
	Port int
}

func NewContact(id NodeID, ip string, port int) (Contact) {
	return Contact{id, ip, port}
}


type RoutingTable struct {
	node NodeID;
	ipAddress string
	port int
	buckets [IdLength*8]*list.List;
}

func NewRoutingTable(contact Contact) (ret RoutingTable) {
	ret.node = contact.Id
	ret.ipAddress = contact.Ip
	ret.port = contact.Port
	for i := 0; i < IdLength * 8; i++ {
		ret.buckets[i] = list.New();
	}
	return;
}

func findElementInBucket(bucket *list.List, id NodeID) (ret *NodeID){

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
}

func (contact *Contact) String() (s string){
	s = "( "
	s += "Node: " + contact.Id.String() + ", "
	s += "Ip: " + contact.Ip + ", "
	s += "Port: " + strconv.Itoa(contact.Port) + " ) "
	return
}

func (table *RoutingTable) String() string {
	s := ""
	bucketNumbers := make([]int, 0)

	for bucketNumber, b := range table.buckets {
		if b.Len() > 0 {
			s += b.Front().Value.(*Contact).String() + "\n"
			bucketNumbers = append(bucketNumbers, bucketNumber)
		}
	}
	s+=  "Occupied buckets: "
	s+= "[ "
	for _, bucket := range bucketNumbers {
		s += strconv.Itoa(bucket)
	}
	s += " ]"
	return s
}

type ContactRecordList []*ContactRecord

func (l ContactRecordList) Len() int { return len(l) }
func (l ContactRecordList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (a ContactRecordList) Less(i, j int) bool { return a[i].sortKey.Less(a[j].sortKey)}



type ContactRecord struct{
	node *Contact
	sortKey NodeID
}

func (contactRecord *ContactRecord) String() (s string) {
	s = "( Contact: " + contactRecord.node.String() + ", " + "Distance: " + contactRecord.sortKey.String() + " )"
	return
}

func (rec *ContactRecord) Less(other interface{}) bool {
	return rec.sortKey.Less(other.(*ContactRecord).sortKey);
}


func copyToList(start, end *list.Element, vec *[]*ContactRecord, target NodeID) {
	for el := start; el != end; el = el.Next() {
		contact := el.Value.(*Contact)
		*vec = append(*vec, &ContactRecord{node: contact, sortKey: contact.Id.Distance(target)})
	}
}

func (table *RoutingTable) FindClosest(target NodeID, count int) (ret []*ContactRecord){
	ret = make([]*ContactRecord, 0)
	prefixLen := table.node.Distance(target).PrefixLen()
	bucket := table.buckets[prefixLen]
	copyToList(bucket.Front(), nil, &ret, target)
	for i := 1; (prefixLen-i >= 0 || prefixLen+i < IdLength*8) && len(ret) < count; i++ {
		if prefixLen - i >= 0 {
			bucket = table.buckets[prefixLen - i];
			copyToList(bucket.Front(), nil, &ret, target);
		}
		if prefixLen + i < IdLength * 8 {
			bucket = table.buckets[prefixLen + i];
			copyToList(bucket.Front(), nil, &ret, target);
		}
	}
	sort.Sort(ContactRecordList(ret))
	if len(ret) > count {
		ret = ret[0:count]
	}
	return
}