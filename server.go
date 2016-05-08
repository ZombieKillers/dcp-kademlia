package kademlia

import (
	"fmt"
	"strings"
	"net"
	"strconv"
	"errors"
)


// Server implementation

type KademliaServer struct {
	contact *Contact
	PingContacts chan Contact
	Done chan bool
}

func (ks *KademliaServer) GetStates() []string{
	states := []string {"PING", "STORE", "FIND_NODE", "FIND_VALUE"}
	return states
}

func (ks *KademliaServer) getLocalAddress() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", ks.contact.Ip + ":" + strconv.Itoa(ks.contact.Port))

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

		if err = ks.HandleMessage(strings.Split(msg, " "), addr); err != nil{
			fmt.Println("[ERROR]", err)
			ks.Done <- true
		}
	}
	return nil
}


func (ks *KademliaServer) StartServer(self *Contact) error {
	ks.contact = self
	ks.Done = make(chan bool, 1)
	ks.PingContacts = make(chan Contact, 5)
	fmt.Println("Port:", self.Port)
	ServerAddr, e := net.ResolveUDPAddr("udp",  ":" + strconv.Itoa(self.Port))
	if e != nil {
		return e
	}
	l, e := net.ListenUDP("udp", ServerAddr)
	if e != nil {
		return e
	}

	go ks.ListenForMessages(l)
	return nil
}




