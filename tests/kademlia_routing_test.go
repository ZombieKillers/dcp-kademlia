package rpc_test

import ("testing"
	"../nodes"
	"../table"
	"fmt"
)
func TestRoutingTableUpdate(t *testing.T) {

	contact1 := table.NewContact(nodes.NewRandomNodeId(), "192.168.0.1", 12345)
	contact2 := table.NewContact(nodes.NewRandomNodeId(), "192.168.0.1", 12345)
	fmt.Println("Printing the initial nodes: ")
	fmt.Println(contact1)
	fmt.Println(contact2)
	fmt.Println("====================\n")

	routingTable := table.NewRoutingTable(contact1.Id)
	routingTable.Update(contact2)
	fmt.Println(&routingTable)
}


func TestRoutingTableFind(t *testing.T) {
	contact1 := table.NewContact(nodes.NewRandomNodeId(), "192.168.0.1", 12345)
	contact2 := table.NewContact(nodes.NewRandomNodeId(), "192.168.0.1", 12345)
	fmt.Println("Printing the initial nodes: ")
	fmt.Println(contact1)
	fmt.Println(contact2)
	fmt.Println("====================\n")

	routingTable := table.NewRoutingTable(contact1.Id)
	routingTable.Update(contact2)

	fmt.Println(contact1.Id.Distance(contact2.Id).PrefixLen())
	closestNodes := routingTable.FindClosest(contact2.Id, 1)
	if len(closestNodes) != 1 {
		t.Error("Invalid length for the result!")
	}
	fmt.Println(closestNodes)

}