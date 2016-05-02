package nodes

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

const IdLength = 20

type NodeID [IdLength]byte

func NewNodeId(data string) (ret NodeID, err error) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		fmt.Printf("%e\n", err)
		return;
	}

	for i := 0; i < IdLength; i++ {
		ret[i] = decoded[i]
	}
	return;
}

func NewRandomNodeId() (ret NodeID) {
	for i := 0; i < IdLength; i++ {
		ret[i] = uint8(rand.Intn(256))
	}
	return ret
}

func (node NodeID) String() string {
	return hex.EncodeToString(node[0:IdLength])
}

type Node struct {
	Id int
}

func XORDistance(node1 Node, node2 Node) int {
	return node1.Id ^ node2.Id
}
