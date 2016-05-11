package rpc_test

import (
	"testing"
	"time"
	"math/rand"
	"../../dcp-kademlia"
	"fmt"
)

func TestStore(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	contact := kademlia.NewContact(kademlia.NewRandomNodeId(), "127.0.0.1", 15000)

	otherNodeId, _ := kademlia.NewNodeId("8ae869162642ab0a723f0bb6bf3e8c53398b90d2")
	otherContact := kademlia.NewContact(otherNodeId, "127.0.0.1", 12345)
	k := kademlia.NewKademlia(contact, "1")
	k.StartServer()

	request := kademlia.KeyValuePair{Key: "hello", Value:"world"}
	request2 := kademlia.KeyValuePair{Key: "do", Value:"something"}
	done := make(chan bool)
	k.Store(&otherContact, &request, done)
	<- done
	k.Store(&otherContact, &request2, done)
	<- done


	val := k.IterativeFindValue("do", 3)
	fmt.Println(val)
	k.ServerDone <- true
}
