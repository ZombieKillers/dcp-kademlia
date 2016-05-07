package main

import (
	"./nodes"
	"./server"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	n1 := nodes.NewRandomNodeId()
	fmt.Println(n1)

	fmt.Println("Trying some RPC stuff")

	if err := server.StartServer(); err != nil {
		panic(err)
	}
}
