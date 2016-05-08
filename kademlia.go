package kademlia

import (
	"fmt"
	"time"
)

type Kademlia struct {
	self      Contact
	server    *KademliaServer
	routes    *RoutingTable
	ServerDone chan bool
	NetworkId string
}

func NewKademlia(self Contact, networkId string) (ret Kademlia) {
	routingTable := NewRoutingTable(self)
	ret = Kademlia{self: self, routes: &routingTable,
		NetworkId: networkId, ServerDone: make(chan bool, 1)}
	return
}

func (k *Kademlia) GetRoutes() (routes *RoutingTable) {
	routes = k.routes
	return
}

func (k *Kademlia) StartServer() {
	k.server = &KademliaServer{}
	k.server.StartServer(&k.self)
	go k.handleRemoteResponses()
	fmt.Println("Listening for requests at:", k.self.Ip,":", k.self.Port )
}

func (k *Kademlia) handleRemoteResponses() {
	for {
		select {
		case replyContact := <-k.server.PingContacts:
			fmt.Println("Got contact from a ping")
			k.routes.Update(&replyContact)
		case request := <-k.server.FindNodeRequests:
			k.routes.Update(request.destination)
			k.server.FindNodeReplies <- k.routes.FindClosest(request.target, BucketSize)
		case done := <-k.server.Done:
			if done == true {
				k.ServerDone <- done
			}
		case err := <-k.server.Errors:
			fmt.Println("[SERVER_ERROR]", err)
		default:
			time.Sleep(time.Millisecond * 5)
		}
	}
}
