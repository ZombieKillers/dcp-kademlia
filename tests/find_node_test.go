package rpc_test

import (
	"testing"
	"../../dcp-kademlia"
	"fmt"
)

func TestFindNodeFromOtherKademlia(t *testing.T) {
	contact := kademlia.NewContact(kademlia.NewRandomNodeId(), "127.0.0.1", 15000)
	otherNodeId, _ := kademlia.NewNodeId("8ae869162642ab0a723f0bb6bf3e8c53398b90d2")
	otherContact := kademlia.NewContact(otherNodeId, "127.0.0.1", 12345)
	k := kademlia.NewKademlia(contact, "1")
	k.StartServer()
	done := make(chan []*kademlia.ContactRecord)
	req := kademlia.NewFindNodeRequest(&otherContact, otherNodeId)
	go k.FindNode(req, done)

	fmt.Println(<-done)

	k.ServerDone <- true
	<- k.ServerDone
}
