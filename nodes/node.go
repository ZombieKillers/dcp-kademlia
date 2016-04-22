package nodes

type Node struct {
	Id int
}

func XORDistance(node1 Node, node2 Node) int{
	return node1.Id ^ node2.Id
}
