package main

import (
	"./nodes"
	"fmt"
	"math/rand"
	"time"
	"./kademlia"
	"./table"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	contact := table.NewContact(nodes.NewRandomNodeId(), "127.0.0.1", 12345)
	fmt.Println(contact.Port)
	k := kademlia.NewKademlia(contact, "1")


	fmt.Println("Trying some RPC stuff")

	k.StartServer()
}
