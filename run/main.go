package main

import (
	"../../dcp-kademlia"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Unix(1,2).UnixNano())
}

func main() {
	nodeId, _ := kademlia.NewNodeId("8ae869162642ab0a723f0bb6bf3e8c53398b90d2")
	contact := kademlia.NewContact(nodeId, "127.0.0.1", 12345)
	k := kademlia.NewKademlia(contact, "1")

	k.StartServer()
	<- k.ServerDone
}
