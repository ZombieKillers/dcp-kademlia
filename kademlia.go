package kademlia

import (
	"fmt"
	"time"
	"crypto/sha1"
	"io"
)

type Kademlia struct {
	self      Contact
	HashTable map[string]string
	server    *KademliaServer
	routes    *RoutingTable
	ServerDone chan bool
	NetworkId string
}

func NewKademlia(self Contact, networkId string) (ret Kademlia) {
	routingTable := NewRoutingTable(self)
	fmt.Println("Kademlia contact info:", self)
	ret = Kademlia{self: self, routes: &routingTable,
		NetworkId: networkId, ServerDone: make(chan bool, 1),
		HashTable: make(map[string]string)}
	return
}

func (k *Kademlia) GetRoutes() (routes *RoutingTable) {
	routes = k.routes
	return
}

func (k *Kademlia) StartServer() {
	k.server = &KademliaServer{}
	k.server.StartServer(&k.self)
	go k.handleRemoteResponses()
	fmt.Println("Listening for requests at:", k.self.Ip,":", k.self.Port )
}


func (k *Kademlia) updateRoutingTable(replyContact *Contact) {
	oldestNode, bucketNr := k.routes.Update(replyContact)
	if oldestNode != nil {
		reply := k.Ping(oldestNode)
		if reply != nil{
			bucket := k.routes.buckets[bucketNr]
			bucket.MoveToFront(bucket.Back())
		} else{
			bucket := k.routes.buckets[bucketNr]
			bucket.Remove(bucket.Back())
			bucket.PushFront(replyContact)
		}
	}
}

func (k *Kademlia) handleRemoteResponses() {
	for {
		select {
		case replyContact := <-k.server.PingContacts:
			k.updateRoutingTable(&replyContact)
		case request := <-k.server.FindNodeRequests:
			k.updateRoutingTable(request.destination)
			k.server.FindNodeReplies <- k.routes.FindClosest(request.target, BucketSize)
		case request := <- k.server.FindValueRequests:
			k.updateRoutingTable(request.destination)
			if val, ok := k.HashTable[request.target]; ok {
				fmt.Println(val)
				// TODO send the value
			} else {
				// TODO send closest nodes
			}
			
		case done := <-k.server.Done:
			if done == true {
				k.server.ServerHandle.Close()
				k.ServerDone <- done
			}
		case pair := <-k.server.KeyValuePairs:
			k.storePair(pair)
		case err := <-k.server.Errors:
			fmt.Println("[SERVER_ERROR]", err)
		default:
			time.Sleep(time.Millisecond * 5)
		}
	}
}

func (k *Kademlia) storePair(pair *KeyValuePair) {
	hash := k.getHashForValue(pair.Key)
	fmt.Printf("Hash: %x \n", hash)
	_, err := NewNodeId(fmt.Sprintf("%x", hash))
	if err != nil {
		fmt.Println(err)
	}
	k.HashTable[string(hash)] = pair.Value
}


func (k *Kademlia) getHashForValue(value string) (hash []byte) {
	sha := sha1.New()
	io.WriteString(sha, value)
	hash = sha.Sum(nil)
	return
}

func (k *Kademlia) StopServer(){
	k.server.ServerHandle.Close()
}
