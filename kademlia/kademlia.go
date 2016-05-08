package kademlia

import (
	"../table"
	"../nodes"
	"net"
	"strings"
	"fmt"
	"errors"
	"strconv"
)

type Kademlia struct {
	self table.Contact
	ownServer *KademliaServer
	routes *table.RoutingTable
	NetworkId string
}

func (k *Kademlia) GetRoutes() (routes *table.RoutingTable){
	routes = k.routes
	return
}

func (k *Kademlia) StartServer() {
	k.ownServer = &KademliaServer{k}
	k.ownServer.StartServer(k.self)
}

func NewKademlia(self table.Contact, networkId string) (ret Kademlia) {
	routingTable := table.NewRoutingTable(self)
	ret = Kademlia{self: self, routes: &routingTable, NetworkId: networkId }
	return
}



// Server implementation

type KademliaServer struct {
	k *Kademlia
}

func (ks *KademliaServer) GetStates() []string{
	states := []string {"PING", "STORE", "FIND_NODE", "FIND_VALUE"}
	return states
}

func (ks *KademliaServer) getLocalAddress() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", ks.k.self.Ip + ":" + strconv.Itoa(ks.k.self.Port))

}

func (ks *KademliaServer) handlePing(message []string, address *net.UDPAddr) error {
	messageId, err := nodes.NewNodeId(message[1])
	if err != nil {
		return err
	}

	otherNodeId, err := nodes.NewNodeId(message[0])
	if err != nil {
		return err
	}

	LocalAddr, err := ks.getLocalAddress()
	if err != nil {
		return err
	}

	Conn, err := net.DialUDP("udp", LocalAddr, address)
	if err != nil {
		return err
	}

	reply := "PONG " + ks.k.self.Id.String() + " " + messageId.String()
	_, err = Conn.Write([]byte(reply))
	if err != nil {
		return  err
	}
	Conn.Close()

	newContact := table.NewContact(otherNodeId, address.IP.String(),address.Port)
	ks.k.routes.Update(&newContact)
	fmt.Println("Updated kademlia routing table")

	fmt.Println(ks.k.routes)
	return nil
}


func (ks *KademliaServer) HandleMessage(splitMessage []string, address *net.UDPAddr) error {
	states := ks.GetStates()
	if len(splitMessage) < 3 {
		return errors.New("Empty message or unrecognized RPC");
	}
	err := error(nil)
	procedure := splitMessage[0]
	switch procedure {
	case states[0]:
		err = ks.handlePing(splitMessage[1:], address)
		break
	case states[1]:
		fmt.Println("I got a store message!")
		break
	case states[2]:
		fmt.Println("I got a find node message!")
		break
	case states[3]:
		fmt.Println("I got a find value message!")
		break
	default:
		err = errors.New("RPC not found!")
	}

	if err != nil {
		fmt.Println("Error: ", err)
	}
	return nil
}

func (ks *KademliaServer) ListenForMessages(server *net.UDPConn) error {
	defer server.Close()

	fmt.Println("Listening for contacts here")
	buf := make([]byte, 1024)
	for {
		n,addr,err := server.ReadFromUDP(buf)
		msg := string(buf[0:n])
		fmt.Println("Received ", msg, " from ",addr)
		if err != nil {
			fmt.Println("Error: ",err)
			return err
		}

		ks.HandleMessage(strings.Split(msg, " "), addr)
	}
	return nil
}


func (ks *KademliaServer) StartServer(self table.Contact) error {
	fmt.Println("Port:", self.Port)
	ServerAddr, e := net.ResolveUDPAddr("udp",  ":" + strconv.Itoa(self.Port))
	if e != nil {
		return e
	}
	l, e := net.ListenUDP("udp", ServerAddr)
	if e != nil {
		return e
	}

	ks.ListenForMessages(l)
	return nil
}
