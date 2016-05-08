package main

import (
	"fmt"
	"math/rand"
	"time"
	"../../dcp-kademlia"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	contact := kademlia.NewContact(kademlia.NewRandomNodeId(), "127.0.0.1", 12345)
	fmt.Println(contact.Port)
	k := kademlia.NewKademlia(contact, "1")


	fmt.Println("Trying some RPC stuff")

	k.StartServer()
}
