package kademlia

import (
	"fmt"
	"time"
)

type Kademlia struct {
	self Contact
	server *KademliaServer
	routes *RoutingTable
	NetworkId string
}


func NewKademlia(self Contact, networkId string) (ret Kademlia) {
	routingTable := NewRoutingTable(self)
	ret = Kademlia{self: self, routes: &routingTable, NetworkId: networkId }
	return
}

func (k *Kademlia) GetRoutes() (routes *RoutingTable){
	routes = k.routes
	return
}

func (k *Kademlia) StartServer() {
	k.server = &KademliaServer{}
	k.server.StartServer(&k.self)
	k.HandleRemoteResponses()
}

func (k *Kademlia) HandleRemoteResponses() {
	for {
		select {
		case replyContact := <- k.server.PingContacts:
			fmt.Println("Got contact from a ping")
			k.routes.Update(&replyContact)
		case done := <-k.server.Done:
			if done == true {
				fmt.Println("Exiting...")
				return
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}




