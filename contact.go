package kademlia

type Contact struct {
	Id   NodeID
	Ip   string
	Port int
}

func NewContact(id NodeID, ip string, port int) Contact {
	return Contact{id, ip, port}
}

type ContactRecordList []*ContactRecord

func (l ContactRecordList) Len() int           { return len(l) }
func (l ContactRecordList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (a ContactRecordList) Less(i, j int) bool { return a[i].sortKey.Less(a[j].sortKey) }

type ContactRecord struct {
	node    *Contact
	sortKey NodeID
}

func NewContactRecord(cont *Contact, distance NodeID) ContactRecord {
	return ContactRecord{cont, distance}
}

func (contactRecord *ContactRecord) String() (s string) {
	s = "( Contact: " + contactRecord.node.String() + ", " + "Distance: " + contactRecord.sortKey.String() + " )"
	return
}

func (rec *ContactRecord) Less(other interface{}) bool {
	return rec.sortKey.Less(other.(*ContactRecord).sortKey)
}
