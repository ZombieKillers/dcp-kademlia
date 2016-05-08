package kademlia

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
		return
	}

	for i := 0; i < IdLength; i++ {
		ret[i] = decoded[i]
	}
	return
}

func NewRandomNodeId() (ret NodeID) {
	for i := 0; i < IdLength; i++ {
		ret[i] = uint8(rand.Intn(256))
	}
	return ret
}

func (node NodeID) Less(other interface{}) bool {
	for i := range node {
		if node[i] != other.(NodeID)[i] {
			return node[i] < other.(NodeID)[i]
		}
	}
	return false
}

func (node NodeID) Equals(other NodeID) bool {
	for i := 0; i < IdLength; i++ {
		if node[i] != other[i] {
			return false
		}
	}
	return true
}

func (node NodeID) String() string {
	return hex.EncodeToString(node[0:IdLength])
}

func (node NodeID) Distance(other NodeID) (ret NodeID) {
	for i := 0; i < IdLength; i++ {
		ret[i] = node[i] ^ other[i]
	}
	return
}

func (node NodeID) PrefixLen() (ret int) {
	for i := 0; i < IdLength; i++ {
		for j := 0; j < 8; j++ {
			if (node[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	return IdLength*8 - 1
}
