package main

import (
	"./hash"
	"./nodes"
)

func main() {
	n1 := nodes.Node{3}
	n2 := nodes.Node{1023}

	dist := nodes.XORDistance(n1, n2)
	hash.SayHello(dist)

}
