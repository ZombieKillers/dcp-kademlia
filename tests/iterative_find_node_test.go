package rpc_test

import (
	"testing"
	"time"
	"../../dcp-kademlia"
	"fmt"
	"math/rand"
)

func TestIterativeFindNode(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	contact := kademlia.NewContact(kademlia.NewRandomNodeId(), "127.0.0.1", 15000)

	otherNodeId, _ := kademlia.NewNodeId("8ae869162642ab0a723f0bb6bf3e8c53398b90d2")
	otherContact := kademlia.NewContact(otherNodeId, "127.0.0.1", 12345)
	k := kademlia.NewKademlia(contact, "1")

	k.StartServer()
	k.Ping(&otherContact)

	resList := k.IterativeFindNode(otherNodeId, 3)
	fmt.Println("Got results from iterative find: ")
	for el := resList.Front(); el != nil; el = el.Next(){
		record := el.Value.(*kademlia.ContactRecord)
		fmt.Println(record)
	}

	k.ServerDone <- true
	<- k.ServerDone
}